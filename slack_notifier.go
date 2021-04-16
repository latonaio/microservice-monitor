package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type MessageContents struct {
	MicroserviceName string
	DeviceName       string
	IPv4             string
	AlertTime        jsonTime
	AlertLog         string
}

func (mc MessageContents) MapToAttachiments() Attachments {
	return Attachments{
		FallBack:   fmt.Sprintf("[%s]: %s resource alert", mc.DeviceName, mc.MicroserviceName),
		Color:      "#36a64f",
		PreText:    fmt.Sprintf("[%s]: %s resources over threshold", mc.DeviceName, mc.MicroserviceName),
		AuthorName: fmt.Sprintf("%s(%s)", mc.DeviceName, mc.IPv4),
		Title:      "Microservice Monitor Alert",
		Text:       mc.AlertLog,
		Fields:     []Field{},
		Footer:     "Occurred",
		TimeStamp:  mc.AlertTime.Unix(),
	}
}

type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

type Attachments struct {
	FallBack   string  `json:"fallback"`
	Color      string  `json:"color"`
	PreText    string  `json:"pretext"`
	AuthorName string  `json:"author_name"`
	Title      string  `json:"title"`
	Text       string  `json:"text"`
	Fields     []Field `json:"fields"`
	Footer     string  `json:"footer"`
	TimeStamp  int64   `json:"ts"`
}

type Slack struct {
	Attachments []Attachments `json:"attachments"`
	Username    string        `json:"uesrname"`
	IconEmoji   string        `json:"icon_emoji"`
	IconURL     string        `json:"icon_url"`
	Channel     string        `json:"channel"`
}

func Notify(msg MessageContents, env Env) error {
	params := Slack{
		Attachments: []Attachments{msg.MapToAttachiments()},
		Username:    "microservice_monitor",
		IconEmoji:   "",
		IconURL:     "",
		Channel:     "",
	}

	jsonparams, _ := json.Marshal(params)
	resp, err := http.PostForm(
		env.SlackIncomingURL,
		url.Values{"payload": {string(jsonparams)}},
	)
	if err != nil {
		return err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	fmt.Println(body)

	return nil
}
