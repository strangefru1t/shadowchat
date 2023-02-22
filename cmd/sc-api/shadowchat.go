package main

import (
	"github.com/gin-gonic/gin"
	"github.com/strangefru1t/shadowchat/internal/api"
	"github.com/strangefru1t/shadowchat/internal/config"
	dash "github.com/strangefru1t/shadowchat/internal/dashboard"
	"github.com/strangefru1t/shadowchat/internal/moneroRPC"
	"github.com/strangefru1t/shadowchat/internal/psql"
)

func main() {
	config.Load()
	moneroRPC.InitXMRPrice()
	psql.InitPool()
	defer psql.DBPOOL.Close()
	// Background processes
	go psql.CacheUserLimits("admin", psql.DBPOOL)
	go moneroRPC.UpdateXMRPrice(config.Settings.GetPriceInterval)
	go moneroRPC.CacheMempool(config.Settings.MempoolInterval)
	go psql.ProcessUnpaidChats(psql.DBPOOL)
	go psql.ClearUnpaidChats(config.Settings.UnpaidExpiration, psql.DBPOOL)
	//l, _ := os.Create(".shadowchat.log")
	//gin.DefaultWriter = io.MultiWriter(os.Stdout, l)
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.LoadHTMLGlob("html/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	r.GET("/style.css", func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/css")
		c.HTML(200, "style.css", nil)
	})
	r.GET("/scapi.js", func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/javascript")
		c.HTML(200, "scapi.js", nil)
	})

	ar := r.Group("/api", api.Auth)
	wr := r.Group("/overlay", api.Auth)

	r.POST("/pay", api.Pay)
	r.GET("/dashboard", dash.Dashboard)
	r.GET("/login", dash.Login)
	r.POST("/login", dash.LoginPOST)
	r.POST("/dashboard", dash.Login)
	r.GET("/limits", api.Limits)
	ar.GET("/chatlog", api.ChatLog)
	ar.GET("/settings", api.Settings)
	ar.POST("/settings", api.SettingsUpdate)
	wr.GET("/goal", api.DonationGoalWidget)
	wr.POST("/goal", api.DonationGoalWidget)
	r.GET("/check", api.VerifyID(api.Check))
	r.Run(":8000")

	//moneroRPC.CacheMempool()
	//fmt.Println(moneroRPC.MEMPOOL.Result.Pool)

}
