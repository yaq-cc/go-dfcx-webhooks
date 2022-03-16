package dfcxwebhooks

import (
	"encoding/json"
	"net/http"
)

type Message interface {
	isMessage()
}

type WebhookResponse struct {
	FulfillmentResponse *FulfillmentResponse `json:"fulfillmentResponse,omitempty"`
	PageInfo            *PageInfo            `json:"pageInfo,omitempty"`
	SessionInfo         *SessionInfo         `json:"sessionInfo,omitempty"`
	Payload             map[string]string    `json:"payload,omitempty"`
}

func (r *WebhookResponse) Respond(w http.ResponseWriter) error {
	err := json.NewEncoder(w).Encode(r)
	if err != nil {
		return err
	}
	return nil
}

func (r *WebhookResponse) AddMessage(m Message) *WebhookResponse {
	r.FulfillmentResponse.AddMessage(m)
	return r
}

type FulfillmentResponse struct {
	// https://pkg.go.dev/google.golang.org/genproto/googleapis/cloud/dialogflow/cx/v3beta1#ResponseMessage
	Messages      []Message `json:"messages,omitempty"`
	MergeBehavior int
}

func (fr *FulfillmentResponse) AddMessage(m Message) {
	fr.Messages = append(fr.Messages, m)
}

// Response Messages
type TextMessage struct {
	Text Text `json:"text,omitempty"`
}

func (t *TextMessage) isMessage() {}

type OutputAudioText struct {
	AllowPlaybackInterruption bool   `json:"allowPlaybackInterruption,omitempty"`
	Source                    Source `json:"source,omitempty"`
}

type OutputAudioTextMessage struct {
	OutputAudioText OutputAudioText `json:"outputAudioText,omitempty"`
}

func (o *OutputAudioTextMessage) isMessage() {}

type Source struct {
	Text string `json:"text,omitempty"`
	SSML string `json:"ssml,omitempty"`
}

type PayloadMessage struct {
	Payload map[string]interface{} `json:"payload,omitempty"`
}

func (r *PayloadMessage) isMessage() {}

type RichContentsMessage struct {
	Payload *RichContents `json:"payload"`
}

func (r *RichContentsMessage) isMessage() {}

type RichContents struct {
	RichContent [][]*RichContent `json:"richContent"`
}

func (cs *RichContents) AddContents(c *RichContent) {
	// How to properly check for initialization?
	cs.RichContent[0] = append(cs.RichContent[0], c)
}

func NewRichContentsMessage(c *RichContent) *RichContentsMessage {
	var rcm RichContentsMessage
	if c == nil {
		return &rcm
	}
	rcm.Payload.AddContents(c)
	return &rcm
}

type RichContent struct {
	Type  string `json:"type"`
	Title string `json:"title"`
	Text  string `json:"text"`
	Event *Event `json:"event,omitempty"`
	Icon  *Icon  `json:"icon,omitempty"`
	Link  string `json:"link,omitempty"`
}

type Icon struct {
	Color string `json:"color,omitempty"`
	Type  string `json:"type"`
}

type Event struct {
	Parameters   map[string]string `json:"parameters,omitempty"`
	Name         string            `json:"name"`
	LanguageCode string            `json:"languageCode"`
}

// Helpers
func NewTextResponse(msgs ...string) *WebhookResponse {
	return &WebhookResponse{
		FulfillmentResponse: &FulfillmentResponse{
			Messages: []Message{
				&TextMessage{
					Text: Text{
						Text: msgs,
					},
				},
			},
		},
	}
}

func (wr *WebhookResponse) TextResponse(w http.ResponseWriter, msgs ...string) {
	t := Text{
		Text:                      msgs,
		AllowPlaybackInterruption: true,
	}

	m := TextMessage{
		Text: t,
	}

	wr.FulfillmentResponse = &FulfillmentResponse{
		Messages:      []Message{&m},
		MergeBehavior: 0,
	}
	json.NewEncoder(w).Encode(wr)
}

func (wr *WebhookResponse) SSMLResponse(w http.ResponseWriter, msg string) {
	t := OutputAudioText{
		AllowPlaybackInterruption: true,
		Source: Source{
			SSML: msg,
		},
	}

	m := OutputAudioTextMessage{
		OutputAudioText: t,
	}

	wr.FulfillmentResponse = &FulfillmentResponse{
		Messages:      []Message{&m},
		MergeBehavior: 0,
	}
	json.NewEncoder(w).Encode(wr)
}
