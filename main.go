package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/chat/v1"
	"google.golang.org/appengine" // Required external App Engine library
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

type MessagePayload struct {
	Text string `json:"text"`
}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	ctx := appengine.NewContext(r)

	// Set Headers
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	log.Infof(ctx, "Endpoint reached "+r.URL.Path)
	// Check Endpoint for Secure Endpoint
	if r.URL.Path != "/"+os.Getenv("SECURE_ENDPOINT") {
		http.Error(w, "Bad Request", http.StatusForbidden)
		return
	}

	if r.Header.Get("Authorization") != "Bearer "+os.Getenv("AUTHORIZATION_HEADER") {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	destination := r.Header.Get("Destination")
	space := r.Header.Get("Space")
	threadKey := r.Header.Get("ThreadKey")
	thread := r.Header.Get("Thread")
	if destination == "" {
		http.Error(w, "No Destination Specified", http.StatusBadRequest)
		return
	}
	if space == "" {
		http.Error(w, "No Space Specified", http.StatusBadRequest)
		return
	}

	// Set Context to appengine context

	// Read Body into Bytes Array
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Errorf(ctx, "Error Reading Body "+err.Error())
		http.Error(w, "Error Reading Body", http.StatusInternalServerError)
		return
	}
	log.Infof(ctx, "Body: %+v", string(b))
	var mp MessagePayload
	err = json.Unmarshal(b, &mp)
	if err != nil {
		log.Errorf(ctx, "Error Unmarshalling Message Payload", err)
		http.Error(w, "Error Reading Payload", http.StatusInternalServerError)
		return
	}

	if strings.ToLower(destination) == "google" {
		msg := chat.Message{Text: mp.Text}
		if thread != "" {
			msg.Thread = &chat.Thread{Name: thread}
		}
		err = postToGoogleRoom(ctx, msg, space, threadKey)
	} else if strings.ToLower(destination) == "slack" {
		err = postToSlackRoom(ctx, chat.Message{Text: mp.Text}, space)
	} else if strings.ToLower(destination) == "messenger" {
		err = postToMessengerRoom(ctx, chat.Message{Text: mp.Text}, "", "")
	}
	if err != nil {
		log.Errorf(ctx, "Error Posting to Room", err)
		http.Error(w, "Error Sending Alert", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Success"))
}

func main() {
	http.HandleFunc("/", indexHandler)
	appengine.Main() // Starts the server to receive requests
}

// Helper Function to cut down on code redundancy
func postToGoogleRoom(ctx context.Context, payload chat.Message, space string, threadKey string) error {
	url := "https://chat.googleapis.com/v1/spaces/" + space + "/messages?threadKey=" + threadKey
	client, err := google.DefaultClient(ctx, "https://www.googleapis.com/auth/chat.bot")
	if err != nil {
		log.Errorf(ctx, "Error Getting Default Token Source", err)
		return err
	}
	body, err := json.Marshal(payload)
	resp, err := client.Post(url, "application/json; charset=utf-8", bytes.NewBuffer(body))
	if err != nil {
		log.Infof(ctx, "Error In Post to Room %+v", err)
		return err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	log.Infof(ctx, "Byte to String %v", string(b))
	if err != nil {
		return err
	}

	return nil

}

type SlackPayload struct {
	Text    string `json:"text"`
	Channel string `json:"channel"`
}
type SlackResponse struct {
	Ok      bool         `json:"ok"`
	Channel string       `json:"channel"`
	TS      string       `json:"ts"`
	Message SlackMessage `json:"message"`
}
type SlackMessage struct {
	Text     string `json:"text"`
	Username string `json:"username"`
	BotID    string `json:"bot_id"`
	Type     string `json:"type"`
	SubType  string `json:subtype"`
	TS       string `json:"ts"`
}

func postToSlackRoom(ctx context.Context, payload chat.Message, space string) error {
	url := "https://slack.com/api/chat.postMessage"
	client := urlfetch.Client(ctx)

	sp := SlackPayload{Text: payload.Text, Channel: space}
	body, err := json.Marshal(sp)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		log.Errorf(ctx, "Error Creating new request %v", err)
		return err
	}
	req.Header.Add("Authorization", "Bearer "+os.Getenv("SLACK_TOKEN"))
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf(ctx, "Error In Post to Room %+v", err)
		return err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	log.Infof(ctx, "Byte to String %v", string(b))
	if err != nil {
		return err
	}
	var sr SlackResponse
	err = json.Unmarshal(b, &sr)
	if err != nil {
		log.Errorf(ctx, "Unable unmarshall slack response", err)
		return err
	}
	// log.Infof(ctx,"")
	return nil
}

func postToMessengerRoom(ctx context.Context, payload chat.Message, space string, threadKey string) error {
	return nil
}
