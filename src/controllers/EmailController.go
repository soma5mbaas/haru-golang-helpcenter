package controllers

import (
	"../../src"
	"../handlers"
	"encoding/json"
	"github.com/martini-contrib/render"
	"github.com/streadway/amqp"
	"gopkg.in/mgo.v2"
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

func ExportEmail(req *http.Request, mail MongoExport, r render.Render, ch *amqp.Channel, db *mgo.Database) {
	appid := req.Header.Get("Application-Id")
	if appid == "" {
		r.JSON(handlers.HttpErr(http.StatusNotFound, "insert to Application-Id"))
		return
	}

	mail.Collection = handlers.CollectionTable(mail.Collection, appid)

	if count, err := db.C(mail.Collection).Find(nil).Count(); err != nil || count == 0 {
		r.JSON(http.StatusNotFound, map[string]interface{}{"Export": "Not Found Collection"})
		return
	}

	msg, _ := json.Marshal(mail)

	q, _ := ch.QueueDeclare(
		config.NAMESPACE+":"+"export", // name
		true,  // durable
		false, // delete when usused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
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
