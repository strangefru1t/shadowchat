package api

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/strangefru1t/shadowchat/internal/config"
	"github.com/strangefru1t/shadowchat/internal/moneroRPC"
	"github.com/strangefru1t/shadowchat/internal/psql"
	"github.com/strangefru1t/shadowchat/internal/relayto"
	"golang.org/x/net/websocket"
)

type MessagePostJSON struct {
	Name    string  `json:"name"`
	Amount  float64 `json:"amount"`
	Message string  `json:"message"`
}
type InvoiceJSON struct {
	ID      string  `json:"id"`
	Address string  `json:"address"`
	QR      string  `json:"qr"`
	Name    string  `json:"name"`
	Message string  `json:"message"`
	Amount  float64 `json:"amount"`
}

func Auth(c *gin.Context) {
	session, _ := config.Store.Get(c.Request, "session")
	_, ok := session.Values["user"]
	if !ok {
		c.HTML(http.StatusForbidden, "login.html", nil)
		//	c.AbortWithStatusJSON(http.StatusForbidden, nil)
		c.Abort()
		return
	}
	c.Next()
}
func ChatLog(c *gin.Context) {
	session, _ := config.Store.Get(c.Request, "session")
	val := session.Values["user"]
	log.Println("Returning chat log for user: " + fmt.Sprint(val))
	c.IndentedJSON(http.StatusOK, psql.ChatLog(psql.DBPOOL))
}
func Settings(c *gin.Context) {
	session, _ := config.Store.Get(c.Request, "session")
	u := session.Values["user"]
	s, err := psql.UserSettings(fmt.Sprint(u), psql.DBPOOL)
	if err == nil {
		c.IndentedJSON(http.StatusOK, s)
	} else {

		c.IndentedJSON(http.StatusNotFound, nil)

	}
}
func SettingsUpdate(c *gin.Context) {
	session, _ := config.Store.Get(c.Request, "session")
	u := session.Values["user"]
	var s psql.UserLiveSettings
	if err := c.BindJSON(&s); err != nil {
		println(err)
		return
	}
	if err := psql.UpdateUserSettings(fmt.Sprint(u), s, psql.DBPOOL); err != nil {
		println(err)
		c.IndentedJSON(http.StatusOK, `{"message":"Failed"}`)
		return
	}
	log.Println(s)
	c.IndentedJSON(http.StatusOK, `{"message":"Successfully Updated"}`)
}

func VerifyID(ws websocket.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.IsWebsocket() && len(c.Query("id")) == 16 {
			balance, err := psql.PayIDBalance(c.Query("id"), psql.DBPOOL)
			if err == nil && balance == 0 {
				ws.ServeHTTP(c.Writer, c.Request)
			} else {
				c.AbortWithStatus(http.StatusBadRequest)
			}
		} else {
			c.AbortWithStatus(http.StatusBadRequest)
		}
	}
}
func Check(ws *websocket.Conn) {
	interval := time.Duration(1200) * time.Millisecond
	tk := time.NewTicker(interval)
	for range tk.C {
		_, err := ws.Write([]byte("KEEPALIVE"))
		if err != nil {
			ws.Close()
			break
		}
		balance := moneroRPC.CheckIDMempool(ws.Request().FormValue("id"))
		if balance > 0 {
			ws.Write([]byte(`{"type":"confirmation", "message":` + fmt.Sprint(balance) + `}`))
			psql.MarkPaid(ws.Request().FormValue("id"), balance, psql.DBPOOL)
			var message, name string = psql.ChatData(ws.Request().FormValue("id"), psql.DBPOOL)
			go relayto.Streamlabs(psql.AccessToken("admin", psql.DBPOOL), ws.Request().Header.Get("X-Real-IP"), balance, message, name)
			break
		}
	}
}
func Limits(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.IndentedJSON(http.StatusOK, config.Web)
}

type DonoWidget struct {
	Goal     float64
	Received float64
}

func DonationGoalWidget(c *gin.Context) {
	session, _ := config.Store.Get(c.Request, "session")
	_, ok := session.Values["user"]
	if !ok {
		c.HTML(http.StatusUnauthorized, "login.html", nil)
		return
	}
	var d DonoWidget
	d.Received, d.Goal = psql.DonoGoal(240, psql.DBPOOL)
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.HTML(200, "donogoalwidget.html", d)
}
func Pay(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	var m MessagePostJSON
	if err := c.BindJSON(&m); err != nil {
		println(err)
		return
	}
	var j InvoiceJSON
	if m.Name == "" {
		m.Name = "Anonymous"
	}
	if m.Amount < config.Web.MinDono {
		m.Amount = config.Web.MinDono
	}
	j.Address, j.ID, j.QR = moneroRPC.XMRIntAddrPayID(m.Amount)
	if len(m.Name) > 25 {
		m.Name = m.Name[0:25]
	}
	if len(m.Message) > config.Web.MaximumMessageChars {
		m.Message = m.Message[0:config.Web.MaximumMessageChars]
	}
	j.Message, j.Name, j.Amount = strings.Replace(m.Message, "\n", "", -1), strings.Replace(m.Name, "\n", "", -1), m.Amount
	c.IndentedJSON(http.StatusOK, j)
	psql.StoreSuperchat(j.ID, c.ClientIP(), j.Message, j.Name, psql.DBPOOL)
}
