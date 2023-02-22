package dash

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/strangefru1t/shadowchat/internal/config"
	"github.com/strangefru1t/shadowchat/internal/psql"
	"golang.org/x/crypto/bcrypt"
)

func Dashboard(c *gin.Context) {
	session, err := config.Store.Get(c.Request, "session")
	if err != nil {
		c.HTML(200, "login.html", nil)
		return
	}
	val, ok := session.Values["user"]
	if !ok {
		c.HTML(http.StatusForbidden, "login.html", nil)
		log.Println(val)
		return
	}
	c.HTML(200, "dashboard.html", nil)
}
func LoginPOST(c *gin.Context) {
	session, err := config.Store.Get(c.Request, "session")
	if err != nil {
		c.HTML(200, "login.html", nil)
		return
	}
	val, ok := session.Values["user"]
	if !ok {
		c.HTML(http.StatusForbidden, "login.html", nil)
		log.Println(val)
		return
	}
	c.HTML(200, "dashboard.html", nil)
}
func Login(c *gin.Context) {
	config.Store.Options.HttpOnly = true
	config.Store.Options.Secure = true
	pwhash, err := psql.PassHash(c.PostForm("username"), psql.DBPOOL)
	if err == nil {
		err = bcrypt.CompareHashAndPassword([]byte(pwhash), []byte(c.PostForm("password")))
		if err == nil {
			session, _ := config.Store.Get(c.Request, "session")
			session.Values["user"] = c.PostForm("username")
			session.Save(c.Request, c.Writer)
			c.HTML(200, "dashboard.html", nil)
			return
		} else {
			c.HTML(http.StatusUnauthorized, "login.html", nil)
			return
		}
	} else {
		//log.Println("Login attempted for nonexistent user: " + c.PostForm("username"))
		c.HTML(http.StatusUnauthorized, "login.html", nil)
		return
	}
}
