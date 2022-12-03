package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var sp_dc string
var token string
var db *sql.DB

// Main functions
func main() {
	var err error
	db, err = sql.Open("sqlite3", "file:activity.db")
	handleErr(err)
	get_spdc()
	handleDb()
	refreshToken()
	if sp_dc != "" {
		handleRequests()
	} else {
		println("Please provide your sp_dc token")
		os.Exit(1)
	}
}

type TokenResponse struct {
	ClientId                         string `json:"clientId"`
	AccessToken                      string `json:"accessToken"`
	AccessTokenExpirationTimestampMs int    `json:"accessTokenExpirationTimestampMs"`
	IsAnonymous                      bool   `json:"isAnonymous"`
}

type Config struct {
	SearchEnabled string
}

type FriendActivity struct {
	Friends []struct {
		Timestamp int64 `json:"timestamp"`
		User      struct {
			URI      string `json:"uri"`
			Name     string `json:"name"`
			ImageURL string `json:"imageUrl"`
		} `json:"user"`
		Track struct {
			URI      string `json:"uri"`
			Name     string `json:"name"`
			ImageURL string `json:"imageUrl"`
			Album    struct {
				URI  string `json:"uri"`
				Name string `json:"name"`
			} `json:"album"`
			Artist struct {
				URI  string `json:"uri"`
				Name string `json:"name"`
			} `json:"artist"`
			Context struct {
				URI   string `json:"uri"`
				Name  string `json:"name"`
				Index int    `json:"index"`
			} `json:"context"`
		} `json:"track"`
	} `json:"friends"`
}

// Unified API call stack
func call(url, header string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Got error " + err.Error())
		return err.Error(), errors.New("Error creating request")
	}
	if header == "Cookie" {
		cookie := "sp_dc=" + sp_dc
		req.Header.Add("Cookie", cookie)
	}
	if header == "Auth" {
		cookie := "Bearer " + token
		req.Header.Add("Authorization", cookie)
	}
	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Got error " + err.Error())
		return err.Error(), errors.New("Request error")
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	} else if response.StatusCode != 200 {
		return "", errors.New("Bad status code: " + string(responseData))
	} else {
		return string(responseData), nil
	}
}

// Misc functions
func get_spdc() {
	if _, err := os.Stat("/.dockerenv"); errors.Is(err, os.ErrNotExist) {
		file, err := os.Open("../sp_dc.txt")
		handleErr(err)
		txt, err := ioutil.ReadAll(file)
		handleErr(err)
		sp_dc = strings.TrimSpace(string(txt))
	} else {
		sp_dc = strings.TrimSpace(os.Getenv("SP_DC"))
	}
}

func refreshToken() {
	response, err := call("https://open.spotify.com/get_access_token?reason=transport&productType=web_player", "Cookie")
	handleErr(err)
	var resp_json TokenResponse
	json.Unmarshal([]byte(response), &resp_json)
	token = resp_json.AccessToken
}

func getActivity() (FriendActivity, error) {
	response, err := call("https://guc-spclient.spotify.com/presence-view/v1/buddylist", "Auth")
	var resp_struct FriendActivity
	if err != nil {
		refreshToken()
		return resp_struct, err
	} else {
		json.Unmarshal([]byte(response), &resp_struct)
		return resp_struct, nil
	}
}

// HTTP response handlers
func activityResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	resp, err := getActivity()
	handleErr(err)
	resp_json, err := json.Marshal(resp)
	handleErr(err)
	fmt.Fprintf(w, string(resp_json))
}

func configResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	if _, err := os.Stat("/.dockerenv"); errors.Is(err, os.ErrNotExist) {
		file, err := os.Open("../config.txt")
		handleErr(err)
		txt, err := ioutil.ReadAll(file)
		handleErr(err)
		fmt.Fprintf(w, strings.TrimSpace(string(txt)))
	} else {
		configs := os.Getenv("CONFIG")
		fmt.Fprintf(w, string(configs))
	}
}

func handleRequests() {
	if _, err := os.Stat("/.dockerenv"); errors.Is(err, os.ErrNotExist) {
		fs := http.FileServer(http.Dir("../frontend/dist"))
		http.Handle("/", fs)
	} else {
		fs := http.FileServer(http.Dir("/dist"))
		http.Handle("/", fs)
	}
	http.HandleFunc("/api/latest", activityResponse)
	http.HandleFunc("/config", configResponse)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

// Database functions
func handleDb() {
	if db != nil {
		// if _, err := os.Stat("./activity.db"); errors.Is(err, os.ErrNotExist) {
		str := `
		CREATE TABLE IF NOT EXISTS users(
			user_uri TEXT NOT NULL PRIMARY KEY, 
			user_name TEXT,
			user_image TEXT
		)
		`
		_, err := db.Exec(str)
		handleErr(err)
		// }
		updateUserDb()
		updateActivityDbs()
	}
}

func updateUserDb() {
	getActivity()
	resp, err := getActivity()
	handleErr(err)
	cmd := `
	INSERT OR REPLACE INTO users (
		user_uri,
		user_name,
		user_image
	) values(?, ?, ?)
	`
	stmt, err := db.Prepare(cmd)
	handleErr(err)
	for i := 0; i < len(resp.Friends); i++ {
		user := resp.Friends[i].User
		_, err := stmt.Exec(user.URI, user.Name, user.ImageURL)
		handleErr(err)
	}
}

func updateActivityDbs() {
	user_rows, err := db.Query("SELECT user_name from users")
	handleErr(err)
	var user_id string
	users := make([]string, 0)
	for user_rows.Next() {
		err = user_rows.Scan(&user_id)
		handleErr(err)
		users = append(users, user_id)
	}
	table_rows, err := db.Query("SELECT name from sqlite_schema where type ='table' AND name NOT LIKE 'sqlite_%'")
	handleErr(err)
	var table_name string
	tables := make([]string, 0)
	for table_rows.Next() {
		err = table_rows.Scan(&table_name)
		handleErr(err)
		tables = append(tables, table_name)
	}
	for i := 0; i < len(users); i++ {
		user := strings.ReplaceAll(users[i], " ", "")
		user = strings.ReplaceAll(user, ".", "")
		if !contains(tables, user) {
			createActivityDb(user)
		}
	}
}

func createActivityDb(user string) {
	str := `
	CREATE TABLE IF NOT EXISTS %s (
		timestamp TEXT NOT NULL PRIMARY KEY UNIQUE,
		track_uri TEXT NOT NULL,
		track_name TEXT,
		track_image TEXT,
		album_uri TEXT,
		album_name TEXT,
		artist_uri TEXT,
		artist_name TEXT,
		context_uri TEXT,
		context_name TEXT
	)
	`
	cmd := fmt.Sprintf(str, user)
	_, err := db.Exec(cmd)
	handleErr(err)
}

// Utility functions
func handleErr(error error) {
	if error != nil {
		fmt.Println(error)
	}
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
