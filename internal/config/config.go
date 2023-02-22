package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/gorilla/sessions"
	"github.com/strangefru1t/shadowchat/internal/templates"
)

type ConfigJSON struct {
	XMRURL           string `json:"RPCURL"`
	DBURL            string `json:"DBURL"`
	DBName           string `json:"DBName"`
	DBHost           string `json:"DBHost"`
	DBUser           string `json:"DBUser"`
	DBPass           string `json:"DBPass"`
	MempoolInterval  int    `json:"MempoolCacheIntervalMiliseconds"`
	GetPriceInterval int    `json:"MoneroGetPriceIntervalSeconds"`
	UnpaidExpiration int    `json:"DeleteUnpaidChatsAfterMinutes"`
	StreamlabsToken  string `json:"StreamlabsAccessToken"`
	StreamlabsImage  string `json:"StreamlabsCustomSound"`
	StreamlabsSound  string `json:"StreamlabsCustomImage"`
	USDConversion    bool   `json:"StreamlabsAlertInUSD"`
	CookieSigningKey string `json:"CookieSigningKey"`
}
type WebJSON struct {
	MinDono             float64 `json:"MinimumDonationXMR"`
	MaximumMessageChars int     `json:"MaximumMessageChars"`
	XMRUSD              int     `json:"XMRUSD"`
}

var Settings *ConfigJSON
var Web *WebJSON
var Store *sessions.CookieStore

func Load() {
	file, err := ioutil.ReadFile("./config.json")
	if _, err := os.Stat("html"); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir("html", os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}
	if err != nil {
		ioutil.WriteFile("./config.json", []byte(templates.ConfigJSON), 0600)
	}
	file, err = ioutil.ReadFile("./config.json")
	_, err = ioutil.ReadFile("./html/index.html")
	if err != nil {
		ioutil.WriteFile("./html/index.html", []byte(templates.IndexHTML), 0640)
	}
	_, err = ioutil.ReadFile("./rpc.conf")
	if err != nil {
		ioutil.WriteFile("./rpc.conf", []byte(templates.RPCCONF), 0640)
	}
	_, err = ioutil.ReadFile("./html/login.html")
	if err != nil {
		ioutil.WriteFile("./html/login.html", []byte(templates.LoginHTML), 0640)
	}
	_, err = ioutil.ReadFile("./html/dashboard.html")
	if err != nil {
		ioutil.WriteFile("./html/dashboard.html", []byte(templates.DashboardHTML), 0640)
	}
	_, err = ioutil.ReadFile("./html/donogoalwidget.html")
	if err != nil {
		ioutil.WriteFile("./html/donogoalwidget.html", []byte(templates.DonoGoalWidgetHTML), 0640)
	}
	_, err = ioutil.ReadFile("./html/scapi.js")
	if err != nil {
		ioutil.WriteFile("./html/scapi.js", []byte(templates.SCAPIJS), 0640)
	}
	_, err = ioutil.ReadFile("./html/style.css")
	if err != nil {
		ioutil.WriteFile("./html/style.css", []byte(templates.StyleCSS), 0640)
	}
	err = json.Unmarshal(file, &Settings)
	err = json.Unmarshal(file, &Web)
	if err != nil {
		log.Fatal("Error processing config file", err.Error())
	}
	if Settings.DBPass != "" {
		Settings.DBURL = fmt.Sprintf("postgres://%s:%s@%s/%s", Settings.DBUser, Settings.DBPass, Settings.DBHost, Settings.DBName)
	} else {
		log.Fatal("Please set database password in config.json")
	}
	Store = sessions.NewCookieStore([]byte(Settings.CookieSigningKey))
}
