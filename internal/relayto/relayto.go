package relayto

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"net/url"

	"github.com/strangefru1t/shadowchat/internal/config"
)

func Streamlabs(access_token string, ip string, amount float64, message string, name string) {
	if access_token != "" {
		donation := url.Values{}
		donation.Add("name", name)
		donation.Add("message", message)
		donation.Add("identifier", fmt.Sprintf("%x", md5.Sum([]byte(ip)))[0:12])
		donation.Add("amount", fmt.Sprint(float64(config.Web.XMRUSD)*amount))
		donation.Add("currency", "USD")
		if config.Settings.USDConversion == false {
			donation.Add("skip_alert", "yes")
		}
		donation.Add("access_token", access_token)
		_, err := http.PostForm("https://streamlabs.com/api/v1.0/donations", donation)
		if err != nil {
			fmt.Println(err)
		}
		if config.Settings.USDConversion == false {
			alert := url.Values{}
			alert.Add("type", "donation")
			alert.Add("message", name+" sent "+fmt.Sprint(amount)+" XMR")
			alert.Add("user_message", message)
			alert.Add("special_text_color", "red")
			alert.Add("image_href", config.Settings.StreamlabsImage)
			alert.Add("sound_href", config.Settings.StreamlabsSound)
			alert.Add("access_token", access_token)
			_, err := http.PostForm("https://streamlabs.com/api/v1.0/alerts", alert)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
