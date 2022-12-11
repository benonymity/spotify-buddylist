// by benonymity on 12.3.22

package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
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

// Data models
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

type UserActivity struct {
	Activity []Activity `json:"activity"`
}

type Activity struct {
	Timestamp int   `json:"timestamp"`
	Duration  int   `json:"duration"`
	Track     Track `json:"track"`
}

type Friend struct {
	Timestamp int   `json:"timestamp"`
	Duration  int   `json:"duration"`
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

type TrackInfo struct {
	Album struct {
		AlbumType string `json:"album_type"`
		Artists   []struct {
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href string `json:"href"`
			ID   string `json:"id"`
			Name string `json:"name"`
			Type string `json:"type"`
			URI  string `json:"uri"`
		} `json:"artists"`
		AvailableMarkets []string `json:"available_markets"`
		ExternalUrls     struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Href   string `json:"href"`
		ID     string `json:"id"`
		Images []struct {
			Height int    `json:"height"`
			URL    string `json:"url"`
			Width  int    `json:"width"`
		} `json:"images"`
		Name                 string `json:"name"`
		ReleaseDate          string `json:"release_date"`
		ReleaseDatePrecision string `json:"release_date_precision"`
		TotalTracks          int    `json:"total_tracks"`
		Type                 string `json:"type"`
		URI                  string `json:"uri"`
	} `json:"album"`
	Artists []struct {
		ExternalUrls struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Href string `json:"href"`
		ID   string `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
		URI  string `json:"uri"`
	} `json:"artists"`
	AvailableMarkets []string `json:"available_markets"`
	DiscNumber       int      `json:"disc_number"`
	DurationMs       int      `json:"duration_ms"`
	Explicit         bool     `json:"explicit"`
	ExternalIds      struct {
		Isrc string `json:"isrc"`
	} `json:"external_ids"`
	ExternalUrls struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Href        string `json:"href"`
	ID          string `json:"id"`
	IsLocal     bool   `json:"is_local"`
	Name        string `json:"name"`
	Popularity  int    `json:"popularity"`
	PreviewURL  string `json:"preview_url"`
	TrackNumber int    `json:"track_number"`
	Type        string `json:"type"`
	URI         string `json:"uri"`
}

// Unified API call stack
func call(url, header string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Got error " + err.Error())
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
		log.Println("Got error " + err.Error())
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

func getTrackInfo(id string) (TrackInfo, error) {
	url := fmt.Sprintf("https://api.spotify.com/v1/tracks/%s", id)
	response, err := call(url, "Auth")
	var resp_struct TrackInfo
	if err != nil {
		refreshToken()
		return resp_struct, err
	} else {
		json.Unmarshal([]byte(response), &resp_struct)
		return resp_struct, nil
	}
}

// HTTP response handlers
func userActivityResponse(w http.ResponseWriter, r *http.Request) {
	user := mux.Vars(r)["user"]
	page_str := mux.Vars(r)["page"]
	page, err := strconv.Atoi(page_str)
	if err != nil {
		page = -1
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	uris := getUserUris()
	for i := 0; i < len(uris); i++ {
		if strings.Split(uris[i], ":")[2] == user {
			resp, err := getUserCachedActivity(user, page)
			handleErr(err)
			resp_json, err := json.Marshal(resp)
			handleErr(err)
			fmt.Fprintf(w, string(resp_json))
			return
		}
	}
	fmt.Fprintf(w, "User does not exist")
	return
}

func latestActivityResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	resp, err := getLatestCachedActivity()
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
	router := mux.NewRouter()
	router.HandleFunc("/api/latest", latestActivityResponse)
	router.HandleFunc("/api/{user}/{page}", userActivityResponse)
	router.HandleFunc("/config", configResponse)
	if _, err := os.Stat("/.dockerenv"); errors.Is(err, os.ErrNotExist) {
		router.PathPrefix("/").Handler(http.FileServer(http.Dir("../frontend/dist")))
	} else {
		router.PathPrefix("/").Handler(http.FileServer(http.Dir("/dist")))
	}
	log.Println("Serving webpage!")
	http.ListenAndServe(":10000", router)
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
		cacheActivity()
		go doEvery(5*time.Second, cacheActivity)
	}
}

func refresh() {
	update := updateUserDb()
	if update {
		updateActivityDbs()
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
			cmd := fmt.Sprintf(`
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
	user_rows.Close()
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
	user_rows.Close()
	return users
}

func getUserUrisTables() ([]string, []string) {
	user_rows, err := db.Query("SELECT user_uri, user_table from users")
	handleErr(err)
	var user_uri string
	var user_table string
	uris := make([]string, 0)
	users := make([]string, 0)
	for user_rows.Next() {
		err = user_rows.Scan(&user_uri, &user_table)
		handleErr(err)
		uris = append(uris, user_uri)
		users = append(users, user_table)
	}
	user_rows.Close()
	return uris, users
}

func updateActivityDbs() {
	users := getUserTables()
	table_rows, err := db.Query("SELECT name from sqlite_schema where type ='table' AND name NOT LIKE 'sqlite_%'")
	handleErr(err)
	var table_name string
	tables := make([]string, 0)
	for table_rows.Next() {
		err = table_rows.Scan(&table_name)
		handleErr(err)
		tables = append(tables, table_name)
	}
	table_rows.Close()
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
		content_index INT,
		duration_ms INT
	)
	`
	cmd := fmt.Sprintf(str, user)
	_, err := db.Exec(cmd)
	handleErr(err)
}

func cacheActivity() {
	uris, users := getUserUrisTables()
	resp, err := callActivity()
	handleErr(err)
	for i := 0; i < len(users); i++ {
		for j := 0; j < len(resp.Friends); j++ {
			if uris[i] == resp.Friends[j].User.URI {
				cmd := fmt.Sprintf(`
				INSERT OR REPLACE INTO %s (
					timestamp,
					track_uri,
					track_name,
					track_image,
					album_uri,
					album_name,
					artist_uri,
					artist_name,
					context_uri,
					context_name,
					content_index,
					duration_ms
				) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
				`, users[i])
				stmt, err := db.Prepare(cmd)
				handleErr(err)
				id := strings.Split(resp.Friends[j].Track.URI, ":")[2]
				track_info, err := getTrackInfo(id)
				handleErr(err)
				length := track_info.DurationMs
				_, err = stmt.Exec(
					resp.Friends[j].Timestamp,
					resp.Friends[j].Track.URI,
					resp.Friends[j].Track.Name,
					resp.Friends[j].Track.ImageURL,
					resp.Friends[j].Track.Album.URI,
					resp.Friends[j].Track.Album.Name,
					resp.Friends[j].Track.Artist.URI,
					resp.Friends[j].Track.Artist.Name,
					resp.Friends[j].Track.Context.URI,
					resp.Friends[j].Track.Context.Name,
					resp.Friends[j].Track.Context.Index,
					length,
				)
				handleErr(err)
			}
		}
	}
}

func getLatestCachedActivity() (FriendActivity, error) {
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
		var friend Friend
		for rows.Next() {
			err = rows.Scan(&friend.User.URI, &friend.User.Name, &friend.User.ImageURL)
			handleErr(err)
		}
		rows.Close()
		cmd = fmt.Sprintf("SELECT * from %s ORDER BY timestamp DESC LIMIT 1", tables[i])
		rows, err = db.Query(cmd)
		handleErr(err)
		for rows.Next() {
			var duration_ms sql.NullInt64
			err = rows.Scan(
				&friend.Timestamp,
				&friend.Track.URI,
				&friend.Track.Name,
				&friend.Track.ImageURL,
				&friend.Track.Album.URI,
				&friend.Track.Album.Name,
				&friend.Track.Artist.URI,
				&friend.Track.Artist.Name,
				&friend.Track.Context.URI,
				&friend.Track.Context.Name,
				&friend.Track.Context.Index,
				&duration_ms,
			)
			handleErr(err)
			if duration_ms.Valid {
				duration_ms.Int64 = 300000
			}
			friend.Duration = int(duration_ms.Int64)
			activity = append(activity, friend)
		}
		rows.Close()
	}
	return FriendActivity{Friends: activity}, nil
}

func getUserCachedActivity(user string, page int) (UserActivity, error) {
	var cmd string
	if page == -1 {
		cmd = fmt.Sprintf("SELECT * from user%s", user)
	} else {
		cmd = fmt.Sprintf("SELECT * from user%s ORDER BY timestamp DESC LIMIT 20 OFFSET %v", user, page*20)
	}
	println(cmd)
	rows, err := db.Query(cmd)
	handleErr(err)
	activity := make([]Activity, 0)
	for rows.Next() {
		var friend Activity
		var duration_ms sql.NullInt64
		err = rows.Scan(
			&friend.Timestamp,
			&friend.Track.URI,
			&friend.Track.Name,
			&friend.Track.ImageURL,
			&friend.Track.Album.URI,
			&friend.Track.Album.Name,
			&friend.Track.Artist.URI,
			&friend.Track.Artist.Name,
			&friend.Track.Context.URI,
			&friend.Track.Context.Name,
			&friend.Track.Context.Index,
			&duration_ms,
		)
		handleErr(err)
		if duration_ms.Int64 == 0 {
			duration_ms.Int64 = 300000
		}
		friend.Duration = int(duration_ms.Int64)
		activity = append(activity, friend)
	}
	rows.Close()
	return UserActivity{Activity: activity}, nil

}

// Utility functions
func handleErr(error error) {
	if error != nil {
		log.Println("Error: ", error)
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

func doEvery(d time.Duration, f func()) {
	for range time.Tick(d) {
		f()
	}
}
