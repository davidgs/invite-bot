package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
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

type MessageStruct struct {
	ClientMsgID     string `json:"ClientMsgID"`
	Type            string `json:"Type"`
	User            string `json:"User"`
	Text            string `json:"Text"`
	ThreadTimeStamp string `json:"ThreadTimeStamp"`
	TimeStamp       string `json:"TimeStamp"`
	Channel         string `json:"Channel"`
	ChannelType     string `json:"ChannelType"`
	EventTimeStamp  string `json:"EventTimeStamp"`
	UserTeam        string `json:"UserTeam"`
	SourceTeam      string `json:"SourceTeam"`
	Message         string `json:"Message"`
	PreviousMessage string `json:"PreviousMessage"`
	Edited          string `json:"Edited"`
	SubType         string `json:"SubType"`
	BotID           string `json:"BotID"`
	Username        string `json:"Username"`
	Icons           string `json:"Icons"`
	Upload          bool   `json:"Upload"`
	Files           []any  `json:"Files"`
	Attachments     []any  `json:"Attachments"`
	Root            string `json:"Root"`
}

type Payload struct {
	Token               string `json:"token"`
	TeamID              string `json:"team_id"`
	ContextTeamID       string `json:"context_team_id"`
	ContextEnterpriseID any    `json:"context_enterprise_id"`
	APIAppID            string `json:"api_app_id"`
	Event               struct {
		ClientMsgID string `json:"client_msg_id"`
		Type        string `json:"type"`
		Text        string `json:"text"`
		User        string `json:"user"`
		Ts          string `json:"ts"`
		Blocks      []struct {
			Type     string `json:"type"`
			BlockID  string `json:"block_id"`
			Elements []struct {
				Type     string `json:"type"`
				Elements []struct {
					Type   string `json:"type"`
					UserID string `json:"user_id,omitempty"`
					Text   string `json:"text,omitempty"`
				} `json:"elements"`
			} `json:"elements"`
		} `json:"blocks"`
		Team        string `json:"team"`
		Channel     string `json:"channel"`
		EventTs     string `json:"event_ts"`
		ChannelType string `json:"channel_type"`
	} `json:"event"`
	Type           string `json:"type"`
	EventID        string `json:"event_id"`
	EventTime      int    `json:"event_time"`
	Authorizations []struct {
		EnterpriseID        any    `json:"enterprise_id"`
		TeamID              string `json:"team_id"`
		UserID              string `json:"user_id"`
		IsBot               bool   `json:"is_bot"`
		IsEnterpriseInstall bool   `json:"is_enterprise_install"`
	} `json:"authorizations"`
	IsExtSharedChannel bool   `json:"is_ext_shared_channel"`
	EventContext       string `json:"event_context"`
}


type ChallengeResponse struct {
	Challenge string
}

func invites(w http.ResponseWriter, r *http.Request) {
	signingSecret := os.Getenv("SLACK_SIGNING_SECRET")
	body, err := io.ReadAll(r.Body)
	if err != nil {
	  log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sv, err := slack.NewSecretsVerifier(r.Header, signingSecret)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := sv.Write(body); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := sv.Ensure(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	msgT := InvitePayload{}
	err = json.Unmarshal([]byte(body), &msgT)
	if err != nil {
		log.Println(err)
	}
	output, err := json.MarshalIndent(msgT, "", "   ")
	if err != nil {
		log.Println(err)
	}
	log.Println("InvitePayload", string(output))
	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	eventsAPICallbackEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
	if err != nil {
		log.Println(err)
	}
	log.Println("Type: ", eventsAPIEvent.Type)
	log.Println("SubType: ", eventsAPIEvent.SubType)
	log.Println("APICallbackEventType: ", eventsAPICallbackEvent.Type)
	// respond to a challenge.
	if eventsAPIEvent.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(body), &r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text")
		w.Write([]byte(r.Challenge))
	}
	log.Println("InnerEvent: ", eventsAPIEvent.InnerEvent.Type)
	if eventsAPIEvent.InnerEvent.Type == "message" {
		  var p2 = slackevents.MessageEvent{}
			var pay = Payload{}
			err := json.Unmarshal(body, &pay)
			if err != nil {
				log.Println(err)
			}
			err = json.Unmarshal(body, &p2)
			if err != nil {
				log.Println(err)
			}
			msgText := pay.Event.Text
			msgID := pay.EventID
			log.Println(msgID)
			log.Printf("Incoming Message: %s\n", msgText)
			if strings.Contains(msgText, "?") || strings.Contains(msgText, "help") || strings.Contains(msgText, "Help") {
				log.Printf("Help Requested: %s", msgText)
			}
			if p2.Attachments != nil {
				a0 := p2.Attachments[0].Text
				if strings.Contains(a0, "denied this request") {
					log.Println(a0)
					handleInvite(pay)
				}
			}
			if strings.Contains(msgText, "denied this request") {
				log.Printf("Denied %s\n", pay)
				handleInvite(pay)

			}
	}
	if r.Method == "GET" {
		http.Error(w, "GET Method not supported", 400)
	}  else {
		w.WriteHeader(200)
	}
}


func handleInvite(data Payload) {
	apiToken := os.Getenv("SLACK_SECRET")
	var final_msg = ":avocado-heart: Sorry, <@" + data.Event.User + ">, but direct invites are not allowed in this Slack. All members must go through the application process at: https://devrelcollective.fun We appreciate your understanding."
	// 	var ts = getTS(dat)
	reply_url := "https://slack.com/api/chat.postMessage"
	reqBody, err := json.Marshal(map[string]string{
		"channel":          data.Event.User,
		"replace_original": "false",
		"text":             final_msg,
		"username":         "InviteBot",
		"as_user":          "true",
		"message_ts":       data.Event.Ts,
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
	request.Header.Set("Authorization", "Bearer " + apiToken)
	request.Header.Set("Accept", "application/json")
	res, err := DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != 200 {
		log.Fatal(res.StatusCode)
	}
		reqBody, err = json.Marshal(map[string]string{
			"channel":          "G0A7K9GPN",
			"replace_original": "true",
			"text":             ":avocado-heart: InviteBot Handled this via DM",
			"username":         "InviteBot",
			"thread_ts":        data.Event.EventTs,
		})
		if err != nil {
			log.Fatal(err)
		}
		req, err := http.NewRequest("POST", reply_url, strings.NewReader(string(reqBody)))
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+apiToken)
		req.Header.Set("Accept", "application/json")
		res, err = DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		if res.StatusCode != 200 {
			log.Fatal(res.StatusCode)
		}
}

func main() {
	f, err := os.OpenFile("/var/log/invitebot.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
    log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.SetFlags(log.Lshortfile)
	log.Println("starting ... ")
	http.HandleFunc("/", invites)

	err = http.ListenAndServe(":3333", nil)
	if err != nil {
		log.Fatal("ListenAndServeTLS: ", err)
	}
	log.Println("Server up and running on port 3333")
}
