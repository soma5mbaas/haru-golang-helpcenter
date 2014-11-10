package controllers

import (
	"encoding/json"
	"github.com/martini-contrib/render"
	"github.com/streadway/amqp"
	"net/http"
)

type Email struct {
	Address string `json:"address"` // Email
	Title   string `json:"title"`   // 제목
	Body    string `json:"body"`    // 본문
}

func SendEmail(req *http.Request, mail Email, r render.Render, ch *amqp.Channel) {
	msg, _ := json.Marshal(mail)

	q, _ := ch.QueueDeclare(
		"email", // name
		true,    // durable
		false,   // delete when usused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)

	err := ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
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

type MongoExport struct {
	Address    string `json:"address"`    // Mail Address
	Collection string `json:"collection"` // DB.Collection
}

func ExportEmail(req *http.Request, mail MongoExport, r render.Render, ch *amqp.Channel) {
	msg, _ := json.Marshal(mail)

	q, _ := ch.QueueDeclare(
		"export", // name
		true,     // durable
		false,    // delete when usused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)

	err := ch.Publish( //RabbitMQ에 큐가 생성되어있어야 들어감.
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "json/plain",
			Body:        []byte(msg),
		})
	if err == nil {
		r.JSON(http.StatusOK, map[string]interface{}{"Export": mail})
	} else {
		r.JSON(http.StatusInternalServerError, err)
	}
}
