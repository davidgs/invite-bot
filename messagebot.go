package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)


type InvitePayload struct {
	Token    string `json:"token"`
	TeamID   string `json:"team_id"`
	APIAppID string `json:"api_app_id"`
	Event    struct {
		Type    string `json:"type"`
		Subtype string `json:"subtype"`
		Hidden  bool   `json:"hidden"`
		Message struct {
			BotID  string `json:"bot_id"`
			Type   string `json:"type"`
			Text   string `json:"text"`
			User   string `json:"user"`
			Team   string `json:"team"`
			Edited struct {
				User string `json:"user"`
				Ts   string `json:"ts"`
			} `json:"edited"`
			Attachments []struct {
				Text string `json:"text"`
				ID   int    `json:"id,omitempty"`
			} `json:"attachments"`
			Ts string `json:"ts"`
		} `json:"message"`
		Channel         string `json:"channel"`
		PreviousMessage struct {
			BotID       string `json:"bot_id"`
			Type        string `json:"type"`
			Text        string `json:"text"`
			User        string `json:"user"`
			Ts          string `json:"ts"`
			Team        string `json:"team"`
			Attachments []struct {
				Text       string `json:"text,omitempty"`
				ID         int    `json:"id"`
				CallbackID string `json:"callback_id,omitempty"`
				Fallback   string `json:"fallback,omitempty"`
				Actions    []struct {
					ID    string `json:"id"`
					Name  string `json:"name"`
					Text  string `json:"text"`
					Type  string `json:"type"`
					Value string `json:"value"`
					Style string `json:"style"`
				} `json:"actions,omitempty"`
			} `json:"attachments"`
		} `json:"previous_message"`
		EventTs     string `json:"event_ts"`
		Ts          string `json:"ts"`
		ChannelType string `json:"channel_type"`
	} `json:"event"`
	Type           string `json:"type"`
	EventID        string `json:"event_id"`
	EventTime      int    `json:"event_time"`
	Authorizations []struct {
		EnterpriseID        interface{} `json:"enterprise_id"`
		TeamID              string      `json:"team_id"`
		UserID              string      `json:"user_id"`
		IsBot               bool        `json:"is_bot"`
		IsEnterpriseInstall bool        `json:"is_enterprise_install"`
	} `json:"authorizations"`
	IsExtSharedChannel bool   `json:"is_ext_shared_channel"`
	EventContext       string `json:"event_context"`
}

func checkHeader(key string, data string) bool { // Test Written
	// Create a new HMAC by defining the hash type and the key (as byte array)
	h := hmac.New(sha256.New, []byte(SigningSecret))
	// Write Data to it
	h.Write([]byte(data))
	// Get result and encode as hexadecimal string
	sha := hex.EncodeToString(h.Sum(nil))
	comp := fmt.Sprintf("v0=%s", sha)
	return comp == key
}

func invites(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.Error(w, "GET Method not supported", 400)
	} else {
		key := r.Header.Get("X-Slack-Signature")

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		timestamp := r.Header.Get("X-Slack-Request-Timestamp")
		var pay = InvitePayload{}
		err = json.Unmarshal(body, &pay)
		if err != nil {
			panic(err)
		}
		signedData := fmt.Sprintf("v0:%s:%s", timestamp, string(body))
		if !checkHeader(key, signedData) {
			fmt.Println("Check Header Failed.")
		 	w.WriteHeader(400)
		 	return
		}
		w.WriteHeader(200)
		getRequestType(pay)
	}
}

// is a value in the array?
func isValueInList(value string, list []string) bool { // Test Written
	for _, v := range list {
		if strings.Contains(v, value) {
			return true
		}
	}
	return false
}

func getTS(data InvitePayload) string {
	return data.Event.Message.Ts
}

func getRequestType(dat InvitePayload) {
	reqBody, err := json.Marshal(dat)
	var prettyJSON bytes.Buffer
	_ = json.Indent(&prettyJSON, reqBody, "", "  ")
	log.Printf("Incoming message: %s\n", string(prettyJSON.Bytes()))
	if err != nil {
		log.Fatal(err)
	}
	if strings.Contains(dat.Event.Message.Text, "?") || strings.Contains(dat.Event.Message.Text, "help") || strings.Contains(dat.Event.Message.Text, "Help") {
		log.Printf("Help Requested: %s", dat.Event.Message.Text)
	}

	if strings.Contains(dat.Event.Message.Text, "requested to invite") {
		var prettyJSON bytes.Buffer
		_ = json.Indent(&prettyJSON, reqBody, "", "  ")
		// log.Println(string(prettyJSON.Bytes()))
		sender := dat.Event.Message.Text[strings.Index(dat.Event.Message.Text, "@")+1 : strings.Index(dat.Event.Message.Text, ">")]
		final_msg := ":avocado-heart: Sorry, direct invites are not allowed in this Slack. All members must go through the application process at: https://devrelcollective.fun"
		var ts = getTS(dat)
		reply_url := "https://slack.com/api/chat.postMessage"
		fmt.Println(reply_url)
		reqBody, err = json.Marshal(map[string]string{
			"channel":          sender,
			"replace_original": "false",
			"text":             final_msg,
			"username":         "InviteBot",
			"as_user":          "true",
			"message_ts":       ts,
		})
		if err != nil {
			log.Fatal(err)
		}
		var DefaultClient = &http.Client{}
		request, err := http.NewRequest("POST", reply_url, strings.NewReader(string(reqBody)))
		if err != nil {
			log.Fatal(err)
		}
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", "Bearer " + SlackSecret)
		request.Header.Set("Accept", "application/json")
		res, err := DefaultClient.Do(request)
		if err != nil {
			log.Fatal(err)
		}

		if res.StatusCode != 200 {
			log. Fatal(res.StatusCode)
		}
		message_ts := getTS(dat)
		channel := dat.Event.Channel
		reply_url = "https://slack.com/api/chat.postMessage"
		reqBody, err = json.Marshal(map[string]string{
			"channel":          channel,
			"replace_original": "true",
			"text":             ":avocado-heart: InviteBot Handled this via DM",
			"username":         "InviteBot",
			"thread_ts":        message_ts,
		})
		req, err := http.NewRequest("POST", reply_url, strings.NewReader(string(reqBody)))
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+SlackSecret)
		req.Header.Set("Accept", "application/json")
		res, err = DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		if res.StatusCode != 200 {
			log. Fatal(res.StatusCode)
		}
	}

}

func main() {
	fmt.Println("starting ... ")
	http.HandleFunc("/invites", invites)

	err := http.ListenAndServeTLS(":9932", "/home/davidgs/.node-red/combined", "/home/davidgs/.node-red/combined", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServeTLS: ", err)
	}
}
