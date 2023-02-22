package moneroRPC

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/skip2/go-qrcode"
	"github.com/strangefru1t/shadowchat/internal/config"
)

type TransfersJSON struct {
	ID      string `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		In []struct {
			Address         string  `json:"address"`
			Amount          int64   `json:"amount"`
			Amounts         []int64 `json:"amounts"`
			Confirmations   int     `json:"confirmations"`
			DoubleSpendSeen bool    `json:"double_spend_seen"`
			Fee             int     `json:"fee"`
			Height          int     `json:"height"`
			Locked          bool    `json:"locked"`
			Note            string  `json:"note"`
			PaymentID       string  `json:"payment_id"`
			SubaddrIndex    struct {
				Major int `json:"major"`
				Minor int `json:"minor"`
			} `json:"subaddr_index"`
			SubaddrIndices []struct {
				Major int `json:"major"`
				Minor int `json:"minor"`
			} `json:"subaddr_indices"`
			SuggestedConfirmationsThreshold int    `json:"suggested_confirmations_threshold"`
			Timestamp                       int    `json:"timestamp"`
			Txid                            string `json:"txid"`
			Type                            string `json:"type"`
			UnlockTime                      int    `json:"unlock_time"`
		} `json:"in"`
		Pool []struct {
			Address         string  `json:"address"`
			Amount          int64   `json:"amount"`
			Amounts         []int64 `json:"amounts"`
			DoubleSpendSeen bool    `json:"double_spend_seen"`
			Fee             int     `json:"fee"`
			Height          int     `json:"height"`
			Locked          bool    `json:"locked"`
			Note            string  `json:"note"`
			PaymentID       string  `json:"payment_id"`
			SubaddrIndex    struct {
				Major int `json:"major"`
				Minor int `json:"minor"`
			} `json:"subaddr_index"`
			SubaddrIndices []struct {
				Major int `json:"major"`
				Minor int `json:"minor"`
			} `json:"subaddr_indices"`
			SuggestedConfirmationsThreshold int    `json:"suggested_confirmations_threshold"`
			Timestamp                       int    `json:"timestamp"`
			Txid                            string `json:"txid"`
			Type                            string `json:"type"`
			UnlockTime                      int    `json:"unlock_time"`
		} `json:"pool"`
	} `json:"result"`
}

type IntegratedAddrJSON struct {
	Result struct {
		IntegratedAddress string `json:"integrated_address"`
		PaymentID         string `json:"payment_id"`
	}
}

type PaymentJSON struct {
	ID      string `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		Payments []struct {
			Address      string `json:"address"`
			Amount       int64  `json:"amount"`
			BlockHeight  int    `json:"block_height"`
			Locked       bool   `json:"locked"`
			PaymentID    string `json:"payment_id"`
			SubaddrIndex struct {
				Major int `json:"major"`
				Minor int `json:"minor"`
			} `json:"subaddr_index"`
			TxHash     string `json:"tx_hash"`
			UnlockTime int    `json:"unlock_time"`
		} `json:"payments"`
	} `json:"result"`
}
type MoneroPrice struct {
	Monero struct {
		Usd float64 `json:"usd"`
	} `json:"monero"`
}

var MEMPOOL *TransfersJSON = &TransfersJSON{}

func InitXMRPrice() {
	r, _ := http.NewRequest("GET", "https://api.coingecko.com/api/v3/simple/price?ids=monero&vs_currencies=usd", nil)
	r.Header.Set("Content-Type", "application/json")
	xmprice, err := http.DefaultClient.Do(r)
	if err != nil {
		log.Println(err.Error())
		return
	}
	resp := &MoneroPrice{}
	if err := json.NewDecoder(xmprice.Body).Decode(resp); err != nil {
		log.Println(err.Error())
	}
	config.Web.XMRUSD = int(resp.Monero.Usd)
}
func UpdateXMRPrice(seconds int) {
	interval := time.Duration(seconds) * time.Second
	tk := time.NewTicker(interval)
	for range tk.C {
		req, _ := http.NewRequest("GET", "https://api.coingecko.com/api/v3/simple/price?ids=monero&vs_currencies=usd", nil)
		req.Header.Set("Content-Type", "application/json")
		xmprice, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println(err.Error())
		}
		if err == nil {
			resp := &MoneroPrice{}
			if err := json.NewDecoder(xmprice.Body).Decode(resp); err != nil {
				log.Println(err.Error())
			}
			config.Web.XMRUSD = int(resp.Monero.Usd)
			log.Println("Updated XMR/USD Conversion Price: $" + fmt.Sprint(config.Web.XMRUSD))
		}
	}

}
func CacheMempool(miliseconds int) {
	interval := time.Duration(miliseconds) * time.Millisecond
	tk := time.NewTicker(interval)
	for range tk.C {
		payload := strings.NewReader(`{"jsonrpc":"2.0","id":"0","method":"get_transfers","params":{"pool":true}}`)
		req, err := http.NewRequest("POST", config.Settings.XMRURL, payload)
		if err == nil {
			req.Header.Set("Content-Type", "application/json")
			res, err := http.DefaultClient.Do(req)
			if err == nil {
				if err := json.NewDecoder(res.Body).Decode(MEMPOOL); err != nil {
					fmt.Println(err.Error())
				}
			} else {
				log.Println("Connection to rpc wallet failed")
			}
		}
	}
}

func CheckIDMempool(payid string) float64 {
	for _, tx := range MEMPOOL.Result.Pool {
		if payid == tx.PaymentID {
			return (float64(tx.Amount) / 1000000000000)
		}
	}
	return 0.0
}
func CheckID(payid string) float64 {
	payload := strings.NewReader(fmt.Sprintf(`{"jsonrpc":"2.0","id":"0","method":"get_payments","params":{"payment_id":"%s"}}`, payid))
	req, err := http.NewRequest("POST", config.Settings.XMRURL, payload)
	if err == nil {
		req.Header.Set("Content-Type", "application/json")
		res, err := http.DefaultClient.Do(req)
		if err == nil {
			resp := &PaymentJSON{}
			if err := json.NewDecoder(res.Body).Decode(resp); err != nil {
				fmt.Println(err.Error())
			} else {
				if resp.Result.Payments != nil {
					if payid == resp.Result.Payments[0].PaymentID {
						return (float64(resp.Result.Payments[0].Amount) / 1000000000000)
					}

				}
				return 0.0
			}
		} else {
			log.Println("Connection to rpc wallet failed")
		}
	}
	return 0.00
}

func XMRIntAddrPayID(amount float64) (string, string, string) {
	payload := strings.NewReader(`{"jsonrpc":"2.0","id":"0","method":"make_integrated_address"}`)
	req, err := http.NewRequest("POST", config.Settings.XMRURL, payload)
	if err == nil {
		req.Header.Set("Content-Type", "application/json")
		res, err := http.DefaultClient.Do(req)
		if err == nil {
			resp := &IntegratedAddrJSON{}
			if err := json.NewDecoder(res.Body).Decode(resp); err != nil {
				fmt.Println(err.Error())
			} else {
				return resp.Result.IntegratedAddress, resp.Result.PaymentID, QR(resp.Result.IntegratedAddress, amount)
			}
		} else {
			log.Println("error connecting to rpc wallet")
		}
	}
	return "", "", ""
}
func QR(integratedaddress string, amount float64) string {
	q, _ := qrcode.Encode(fmt.Sprintf("monero:%s?tx_amount=%s", integratedaddress, fmt.Sprint(amount)), qrcode.Low, 420)
	return base64.StdEncoding.EncodeToString(q)
}
