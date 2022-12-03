// by benonymity on 12.3.22

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
	db, err = sql.Open("sqlite3", "./activity.db")
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
	Friends []Friend `json:"friends"`
}

type Friend struct {
	Timestamp int64 `json:"timestamp"`
	User      User  `json:"user"`
	Track     Track `json:"track"`
}

type User struct {
	URI      string `json:"uri"`
	Name     string `json:"name"`
	ImageURL string `json:"imageUrl"`
}

type Track struct {
	URI      string  `json:"uri"`
	Name     string  `json:"name"`
	ImageURL string  `json:"imageUrl"`
	Album    Album   `json:"album"`
	Artist   Artist  `json:"artist"`
	Context  Context `json:"context"`
}

type Album struct {
	URI  string `json:"uri"`
	Name string `json:"name"`
}

type Artist struct {
	URI  string `json:"uri"`
	Name string `json:"name"`
}

type Context struct {
	URI   string `json:"uri"`
	Name  string `json:"name"`
	Index int    `json:"index"`
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

func callActivity() (FriendActivity, error) {
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
	resp, err := cachedActivity()
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
		str := `
		CREATE TABLE IF NOT EXISTS users(
			user_uri TEXT NOT NULL PRIMARY KEY, 
			user_name TEXT,
			user_image TEXT,
			user_table TEXT
		)
		`
		_, err := db.Exec(str)
		handleErr(err)
		update := updateUserDb()
		if update {
			updateActivityDbs()
		}

	}
}

func updateUserDb() bool {
	callActivity()
	resp, err := callActivity()
	handleErr(err)
	cmd := `
	INSERT OR REPLACE INTO users (
		user_uri,
		user_name,
		user_image,
		user_table
	) values(?, ?, ?, ?)
	`
	stmt, err := db.Prepare(cmd)
	handleErr(err)
	users := getUserUris()
	updated := false
	resp_users := make([]string, 0)
	// Update users table with any new followers
	for i := 0; i < len(resp.Friends); i++ {
		user := resp.Friends[i].User
		if !contains(users, user.URI) {
			user_table := "user" + strings.Split(user.URI, ":")[2]
			_, err := stmt.Exec(user.URI, user.Name, user.ImageURL, user_table)
			handleErr(err)
			updated = true
		}
		resp_users = append(resp_users, user.URI)
	}
	// Trim users table for any unfollowed users
	for i := 0; i < len(users); i++ {
		if !contains(resp_users, users[i]) {
			cmd := fmt.Sprint(`
				DELETE FROM users
				WHERE user_uri = %s
			`, users[i])
			_, err := db.Exec(cmd)
			handleErr(err)
		}
	}
	return updated
}

func getUserUris() []string {
	user_rows, err := db.Query("SELECT user_uri from users")
	handleErr(err)
	var user_table string
	users := make([]string, 0)
	for user_rows.Next() {
		err = user_rows.Scan(&user_table)
		handleErr(err)
		users = append(users, user_table)
	}
	return users
}

func getUserTables() []string {
	user_rows, err := db.Query("SELECT user_table from users")
	handleErr(err)
	var user_table string
	users := make([]string, 0)
	for user_rows.Next() {
		err = user_rows.Scan(&user_table)
		handleErr(err)
		users = append(users, user_table)
	}
	return users
}

func updateActivityDbs() {
	user_rows, err := db.Query("SELECT user_table from users")
	handleErr(err)
	var user_table string
	users := make([]string, 0)
	for user_rows.Next() {
		err = user_rows.Scan(&user_table)
		handleErr(err)
		users = append(users, user_table)
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
		context_name TEXT,
		content_index INT
	)
	`
	cmd := fmt.Sprintf(str, user)
	_, err := db.Exec(cmd)
	handleErr(err)
}

func cacheActivity() {
	users := getUserTables()
	cmd := fmt.Sprintf(`
	INSERT OR REPLACE INTO %s (
		user_uri,
		user_name,
		user_image,
		user_table
	) values(?, ?, ?, ?)
	`)
	stmt, err := db.Prepare(cmd)
	resp, err := callActivity()
	handleErr(err)
	for i := 0; i < len(users); i++ {
		fmt.Println(resp.Friends[i].User.Name)
		_, err := stmt.Exec(resp)
		handleErr(err)
	}
}

func getCachedActivity() (FriendActivity, error) {
	user_rows, err := db.Query("SELECT user_table from users")
	handleErr(err)
	var user_table string
	tables := make([]string, 0)
	for user_rows.Next() {
		err = user_rows.Scan(&user_table)
		handleErr(err)
		tables = append(tables, user_table)
	}
	user_rows.Close()
	activity := make([]Friend, 0)
	if len(tables) == 0 {
		return FriendActivity{}, errors.New("No users in table")
	}
	for i := 0; i < len(tables); i++ {
		cmd := fmt.Sprintf(`SELECT user_uri, user_name, user_image from users WHERE user_table = "%s"`, tables[i])
		rows, err := db.Query(cmd)
		handleErr(err)
		var user_uri string
		var user_name string
		var user_image string
		for rows.Next() {
			err = rows.Scan(&user_uri, &user_name, &user_image)
		}
		handleErr(err)
		rows.Close()
		cmd = fmt.Sprintf("SELECT * from %s ORDER BY timestamp DESC LIMIT 1", tables[i])
		for rows.Next() {
			rows, err = db.Query(cmd)
		}
		handleErr(err)
		var timestamp int64
		var track_uri string
		var track_name string
		var track_image string
		var album_uri string
		var album_name string
		var artist_uri string
		var artist_name string
		var context_uri string
		var context_name string
		var context_index int
		err = rows.Scan(
			&timestamp,
			&track_uri,
			&track_name,
			&track_image,
			&album_uri,
			&album_name,
			&artist_uri,
			&artist_name,
			&context_uri,
			&context_name,
			&context_index,
		)
		println(user_name)
		println(track_name)
		handleErr(err)
		rows.Close()
		friend := Friend{
			Timestamp: timestamp,
			User: User{
				URI:      user_uri,
				Name:     user_name,
				ImageURL: user_image,
			},
			Track: Track{
				URI:      track_uri,
				Name:     track_name,
				ImageURL: track_image,
				Album: Album{
					URI:  album_uri,
					Name: album_name,
				},
				Artist: Artist{
					URI:  artist_uri,
					Name: artist_name,
				},
				Context: Context{
					URI:   context_uri,
					Name:  context_name,
					Index: context_index,
				},
			},
		}
		activity = append(activity, friend)
	}
	return FriendActivity{Friends: activity}, nil
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
