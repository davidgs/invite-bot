package main

import (
	// "fmt"
	"testing"
	// "encoding/json"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// const SlackSecret = "xoxb-2633077179-1942456225123-fZ7f0JIRbEAxyWGeIb4OzQb3"
const ConfigString = `v0:1618860506:{"token":"VXTC152lQdTSMcRL7caqvW8E","team_id":"T02JM2959","api_app_id":"A01TQD46ZMZ","event":{"client_msg_id":"a6fc2100-fa51-411e-80d0-96abc55c81c2","type":"message","text":"<@U93FPEB0B> I\u2019m not sure if we\u2019re currently using Zapier for anything. It used to be our way of greeting people but it doesn\u2019t work anymore. I don\u2019t *think* we\u2019re using it for anything else Slack-related. Maybe bring it up at the meeting on Monday to double-check just in case?","user":"U02UAUXM0","ts":"1618860506.005800","team":"T02JM2959","blocks":[{"type":"rich_text","block_id":"tr9O","elements":[{"type":"rich_text_section","elements":[{"type":"user","user_id":"U93FPEB0B"},{"type":"text","text":" I\u2019m not sure if we\u2019re currently using Zapier for anything. It used to be our way of greeting people but it doesn\u2019t work anymore. I don\u2019t "},{"type":"text","text":"think","style":{"bold":true}},{"type":"text","text":" we\u2019re using it for anything else Slack-related. Maybe bring it up at the meeting on Monday to double-check just in case?"}]}]}],"thread_ts":"1618837624.002600","parent_user_id":"U9WKRHW8Z","channel":"G0A7K9GPN","event_ts":"1618860506.005800","channel_type":"group"},"type":"event_callback","event_id":"Ev01UBSQ6M0F","event_time":1618860506,"authorizations":[{"enterprise_id":null,"team_id":"T02JM2959","user_id":"U93FPEB0B","is_bot":false,"is_enterprise_install":false}],"is_ext_shared_channel":false,"event_context":"1-message-T02JM2959-G0A7K9GPN"}
`

func TestCheckHeader(t *testing.T) {
	h := hmac.New(sha256.New, []byte(SigningSecret))
	h.Write([]byte(ConfigString))
	sha := hex.EncodeToString(h.Sum(nil))
	input := fmt.Sprintf("v0=%s", sha)
	result := checkHeader(input, ConfigString)
	// v0=892bd3d0dd5f74e60b16f0c4795fface839448d03092e23aa5ece0ef2d6996f7
	if !result {
		t.Errorf("checkHeader Failed got %v", result)
	}
}