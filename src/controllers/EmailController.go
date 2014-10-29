package controllers

import (
	"encoding/json"
	"github.com/martini-contrib/render"
	"github.com/streadway/amqp"
	"net/http"
)

type Email struct {
	Address string `json:"address"` // UUID
	Title   string `json:"title"`   // 제목
	Body    string `json:"body"`    // 본문
}

func SendEmail(req *http.Request, mail Email, r render.Render, ch *amqp.Channel) {
	msg, _ := json.Marshal(mail)

	err := ch.Publish(
		"",      // exchange
		"email", // routing key
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType: "json/plain",
			Body:        []byte(msg),
		})

	if err == nil {
		r.JSON(http.StatusOK, map[string]interface{}{"Email": mail})
	} else {
		r.JSON(http.StatusInternalServerError, err)
	}
}
