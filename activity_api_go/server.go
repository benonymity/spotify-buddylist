package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"errors"
)

var sp_dc = os.Getenv("SP_DC")
var token string

type TokenResponse struct {
	ClientId                         string `json:"clientId"`
	AccessToken                      string `json:"accessToken"`
	AccessTokenExpirationTimestampMs int    `json:"accessTokenExpirationTimestampMs"`
	IsAnonymous                      bool   `json:"isAnonymous"`
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

func call(url, header string) (string, bool) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Got error " + err.Error())
		return err.Error(), true
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
		return err.Error(), true
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err.Error(), true
	} else if response.StatusCode != 200 {
		// fmt.Println("Bad status code")
		return string(responseData), true
	} else {
		return string(responseData), false
	}
}

func refreshToken() {
	response, err := call("https://open.spotify.com/get_access_token?reason=transport&productType=web_player", "Cookie")
	if err {
		fmt.Println("Error refreshing token!")
	} else {
		var resp_json TokenResponse
		json.Unmarshal([]byte(response), &resp_json)
		token = resp_json.AccessToken
	}
}

func latestActivity(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	response, err := call("https://guc-spclient.spotify.com/presence-view/v1/buddylist", "Auth")
	if err {
		refreshToken()
	} else {
		var resp_struct FriendActivity
		json.Unmarshal([]byte(response), &resp_struct)
		resp_json, err := json.Marshal(resp_struct)
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Fprintf(w, string(resp_json))
	}
}

func handleRequests() {
	if _, err := os.Stat("/.dockerenv"); errors.Is(err, os.ErrNotExist) {
		fs := http.FileServer(http.Dir("../activity_frontend/dist"))
		http.Handle("/", fs)
	} else {
		fs := http.FileServer(http.Dir("/dist"))
		http.Handle("/", fs)
	}
	http.HandleFunc("/api/latest", latestActivity)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {
	if sp_dc != "" {
		handleRequests()
	} else {
		println("Please provide your sp_dc token")
		os.Exit(1)
	}
}
