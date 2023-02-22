package psql

import (
	"context"
	"crypto/md5"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/strangefru1t/shadowchat/internal/config"
	"github.com/strangefru1t/shadowchat/internal/moneroRPC"
	"github.com/strangefru1t/shadowchat/internal/relayto"
	"golang.org/x/crypto/bcrypt"
)

type PaidChat struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Message   string    `json:"message"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}

var DBPOOL *pgxpool.Pool

func InitPool() {
	var err error
	DBPOOL, err = pgxpool.New(context.Background(), config.Settings.DBURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	CreateTableIfNotExists(DBPOOL)
}

func ProcessUnpaidChats(dbpool *pgxpool.Pool) {
	interval := time.Duration(15) * time.Second
	tk := time.NewTicker(interval)
	for range tk.C {
		rows, err := dbpool.Query(context.Background(), "SELECT PayID, Name, Message FROM scapi WHERE Received = 0;")
		if err != nil {
			fmt.Println(err.Error())
		}
		defer rows.Close()
		for rows.Next() {
			var id, name, message string
			err := rows.Scan(&id, &name, &message)
			if err != nil {
				log.Println(err.Error())

			}
			var bal float64 = moneroRPC.CheckID(id)
			if bal > 0 {
				MarkPaid(id, bal, DBPOOL)
				log.Println("Confirmed transacton " + id)
				relayto.Streamlabs(AccessToken("admin", DBPOOL), "noip", bal, message, name)
			}
		}
	}
}
func ClearUnpaidChats(expiration int, dbpool *pgxpool.Pool) {
	interval := time.Duration(10) * time.Second
	tk := time.NewTicker(interval)
	for range tk.C {
		_, err := dbpool.Exec(context.Background(), "DELETE FROM scapi WHERE Received = 0 AND time < NOW() - INTERVAL '"+fmt.Sprint(expiration)+" Minutes';")
		if err != nil {
			fmt.Println(err.Error())
		}
	}

}
func UpdatePassword(newpw string) {

	var adminpass []byte
	adminpass, _ = bcrypt.GenerateFromPassword([]byte(newpw), bcrypt.DefaultCost)

	_, err := DBPOOL.Exec(context.Background(), "UPDATE users SET pwhash = $1;", adminpass)
	if err != nil {
		fmt.Println(err)
	}

}
func CreateTableIfNotExists(dbpool *pgxpool.Pool) {
	var adminpass []byte
	adminpass, _ = bcrypt.GenerateFromPassword([]byte(config.Settings.DBPass), bcrypt.DefaultCost)
	p, err := dbpool.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS scapi (id serial PRIMARY KEY, PayID VARCHAR, UserID VARCHAR, Name VARCHAR, Message VARCHAR, Received DOUBLE PRECISION DEFAULT 0, Time TIMESTAMP WITH TIME ZONE);")
	_, err = dbpool.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS users ( id serial PRIMARY KEY, Username VARCHAR, pwhash VARCHAR, MinDono DOUBLE PRECISION DEFAULT 0.015, DonoGoal DOUBLE PRECISION DEFAULT 0.15,DonoGoalHistHours NUMERIC DEFAULT 24,MaxChars NUMERIC DEFAULT 120, StreamlabsToken VARCHAR DEFAULT '', AlertInUSD BOOL DEFAULT true);")
	_, err = dbpool.Exec(context.Background(), "INSERT INTO users (id, username, pwhash) VALUES (1, 'admin', $1) ON CONFLICT (id) DO NOTHING;", string(adminpass))
	if err != nil {
		fmt.Println(err, p)
	} else {
		log.Println("DATABASE OK")
	}
}
func StoreSuperchat(id string, ip string, message string, name string, dbpool *pgxpool.Pool) {
	if len(message) > config.Web.MaximumMessageChars {
		message = message[0:config.Web.MaximumMessageChars]
	}
	if len(name) > 25 {
		name = name[0:25]
	}
	_, err := dbpool.Exec(context.Background(), "INSERT INTO scapi (Time, Message, Name, PayID, UserID) VALUES (CURRENT_TIMESTAMP(0), $1, $2, $3, $4);", strings.Replace(message, "\n", "", -1), strings.Replace(name, "\n", "", -1), id, fmt.Sprintf("%x", md5.Sum([]byte(ip)))[0:12])
	if err != nil {
		fmt.Println(err)
	}
}
func MarkPaid(payid string, received float64, dbpool *pgxpool.Pool) {
	_, err := dbpool.Exec(context.Background(), "UPDATE scapi SET received = $1 WHERE PayID = $2;", received, payid)
	if err != nil {
		fmt.Println(err)
	}
}
func PayIDBalance(payid string, dbpool *pgxpool.Pool) (float64, error) {
	var balance float64
	err := dbpool.QueryRow(context.Background(), "SELECT Received FROM scapi WHERE PayID = $1;", payid).Scan(&balance)
	if err != nil {
		return 0.0, err
	}
	return balance, err
}
func DonoGoal(hours int, dbpool *pgxpool.Pool) (float64, float64) {
	var goal, received float64

	//err := dbpool.QueryRow(context.Background(), "select sum(received) from scapi;").Scan(&goal)
	//err := dbpool.QueryRow(context.Background(), "select sum(received) from scapi WHERE time > NOW() - INTERVAL '"+fmt.Sprint(hours)+" Hours';").Scan(&goal)
	err := dbpool.QueryRow(context.Background(), "select sum(received),(select donogoal from users where username = 'admin') from scapi WHERE scapi.time > NOW() - INTERVAL '"+fmt.Sprint(hours)+" Hours';").Scan(&received, &goal)
	if err != nil {
		log.Println(err.Error())
		return 0.0, 0.0
	}
	return math.Round(received*1000) / 1000, goal
}

type UserLiveSettings struct {
	MinDono           float64 `json:"mindono"`
	MaxChars          int     `json:"maxchars"`
	StreamLabsToken   string  `json:"streamlabstoken"`
	NewPW             string  `json:"password"`
	AlertInUSD        bool    `json:"alertinusd"`
	DonoGoal          float64 `json:"donogoal"`
	DonoGoalHistHours int     `json:"donogoalhisthours"`
}

func CacheUserLimits(user string, dbpool *pgxpool.Pool) {
	interval := time.Duration(4) * time.Second
	tk := time.NewTicker(interval)
	for range tk.C {
		dbpool.QueryRow(context.Background(), "select mindono, maxchars, alertinusd from users where username = $1;", user).Scan(&config.Web.MinDono, &config.Web.MaximumMessageChars, &config.Settings.USDConversion)
	}
}
func AccessToken(user string, dbpool *pgxpool.Pool) string {
	var t string
	dbpool.QueryRow(context.Background(), "select streamlabstoken from users where username = $1;", user).Scan(&t)
	return t
}

func UpdateUserSettings(user string, edit UserLiveSettings, dbpool *pgxpool.Pool) error {
	_, err := dbpool.Exec(context.Background(), "UPDATE users SET mindono = $1, maxchars = $2, streamlabstoken = $3, alertinusd = $4, donogoal = $5, donogoalhisthours = $6 WHERE username = $7;", edit.MinDono, edit.MaxChars, edit.StreamLabsToken, edit.AlertInUSD, edit.DonoGoal, edit.DonoGoalHistHours, user)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	if edit.NewPW != "" {
		fmt.Println(edit.NewPW)
		UpdatePassword(edit.NewPW)
	}
	return nil
}

func UserSettings(user string, dbpool *pgxpool.Pool) (UserLiveSettings, error) {
	var s UserLiveSettings
	err := dbpool.QueryRow(context.Background(), "select mindono,maxchars,streamlabstoken,alertinusd,donogoal,donogoalhisthours from users where username = $1;", user).Scan(&s.MinDono, &s.MaxChars, &s.StreamLabsToken, &s.AlertInUSD, &s.DonoGoal, &s.DonoGoalHistHours)
	if err != nil {
		log.Println(err.Error())
		return s, err
	}
	return s, nil
}
func ChatData(payid string, dbpool *pgxpool.Pool) (string, string) {
	var message, name string
	err := dbpool.QueryRow(context.Background(), "SELECT Message,Name FROM scapi WHERE PayID = $1;", payid).Scan(&message, &name)
	if err != nil {
		return "error", "error"
	}
	return message, name
}
func ChatLog(dbpool *pgxpool.Pool) []PaidChat {
	var cs []PaidChat
	rows, err := dbpool.Query(context.Background(), "SELECT UserID, Name, Message, Received, time FROM scapi WHERE Received > 0;")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var c PaidChat
		err := rows.Scan(&c.ID, &c.Name, &c.Message, &c.Amount, &c.Timestamp)
		if err != nil {
			log.Println(err.Error())

		}
		cs = append(cs, c)
	}
	return cs
}
func PassHash(username string, dbpool *pgxpool.Pool) (string, error) {
	var passhash string
	err := dbpool.QueryRow(context.Background(), "SELECT pwhash FROM users WHERE username = $1;", username).Scan(&passhash)
	if err != nil {
		return "", err
	} else {
		return passhash, nil

	}
}
