package controllers

import (
	"../handlers"
	"code.google.com/p/go-uuid/uuid"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"time"
)

type QnA struct {
	Id           string `bson:"_id,omitempty"`        // UUID
	EmailAddress string `bson:"emailaddress"`         // Email
	Body         string `bson:"body" `                // 본문
	Category     string `bson:"category"`             // Category
	Time         int64  `bson:"time" `                // 시간
	Reception    bool   `bson:"reception,omitempty" ` // 읽기여부
	Comment      string `bson:"comment" `             // Comment
	CommentTime  int64  `bson:"commenttime" `         // Comment 단 시간
}
type Comment struct {
	Content string `bson:"content" ` // Comment
}

func CreateQnA(req *http.Request, params martini.Params, qna QnA, r render.Render, db *mgo.Database, f *log.Logger) {
	appid := req.Header.Get("Application-Id")
	if appid == "" {
		r.JSON(http.StatusNotFound, "insert to Application-Id")
		return
	}

	qna.Time = time.Now().Unix()
	qna.Id = uuid.New()
	qna.Reception = false
	if qna.Comment != "" {
		qna.CommentTime = qna.Time
	}

	CollectionName := handlers.CollectionNameQnA(appid)
	// if count, _ := db.C(CollectionName).Count(); count > 0 {
	// 	if err := db.Session.DB("admin").Run(bson.M{"shardCollection": "haru" + "." + CollectionName, "key": bson.M{"_id": 1}}, nil); err != nil {
	// 		f.Println(CollectionName+" Sharde Fail :", err)
	// 	} else {
	// 		f.Println(CollectionName+" Sharde ok :", err)
	// 	}
	// }

	if err := db.C(CollectionName).Insert(qna); err != nil {
		r.JSON(handlers.HttpErr(http.StatusNotFound, err.Error()))
		return
	}
	r.JSON(http.StatusOK, map[string]interface{}{"QnA": qna})
}

func AddcommentFaq(req *http.Request, params martini.Params, com Comment, r render.Render, db *mgo.Database, f *log.Logger) {
	appid := req.Header.Get("Application-Id")
	if appid == "" {
		r.JSON(handlers.HttpErr(http.StatusNotFound, "insert to Application-Id"))
		return
	}

	rawId := params["id"]
	colQuerier := bson.M{"_id": rawId}
	change := bson.M{"$set": bson.M{"comment": com.Content, "commenttime": time.Now().Unix()}}
	CollectionName := handlers.CollectionNameQnA(appid)
	if err := db.C(CollectionName).Update(colQuerier, change); err != nil {
		r.JSON(http.StatusNotFound, "NotFound "+rawId)
		return
	}

	r.JSON(http.StatusOK, map[string]interface{}{"UPDATE_OK Comment": com})
}

func ReadListUserQnA(req *http.Request, params martini.Params, r render.Render, db *mgo.Database) {
	appid := req.Header.Get("Application-Id")
	if appid == "" {
		r.JSON(handlers.HttpErr(http.StatusNotFound, "insert to Application-Id"))
		return
	}

	var qnas []QnA
	rawId := params["id"]
	CollectionName := handlers.CollectionNameQnA(appid)
	if err := db.C(CollectionName).Find(bson.M{"emailaddress": rawId}).Sort("-time").All(&qnas); err != nil {
		r.JSON(http.StatusNotFound, err)
		return
	}

	r.JSON(http.StatusOK, map[string]interface{}{"return": qnas})
}

func ReadListQnA(req *http.Request, r render.Render, db *mgo.Database) {
	appid := req.Header.Get("Application-Id")
	if appid == "" {
		r.JSON(handlers.HttpErr(http.StatusNotFound, "insert to Application-Id"))
		return
	}
	////
	colQuerier := bson.M{}
	change := bson.M{"$set": bson.M{"reception": true}}
	CollectionName := handlers.CollectionNameQnA(appid)
	if _, err := db.C(CollectionName).UpdateAll(colQuerier, change); err != nil {
		r.JSON(handlers.HttpErr(http.StatusNotFound, err.Error()))
		return
	}
	var qnas []QnA
	if err := db.C(CollectionName).Find(bson.M{}).Sort("-time").All(&qnas); err != nil {
		r.JSON(handlers.HttpErr(http.StatusNotFound, err.Error()))
		return
	}

	r.JSON(http.StatusOK, map[string]interface{}{"return": qnas})
}

func ReadIdQnA(req *http.Request, params martini.Params, r render.Render, db *mgo.Database) {

	appid := req.Header.Get("Application-Id")
	if appid == "" {
		r.JSON(handlers.HttpErr(http.StatusNotFound, "insert to Application-Id"))
		return
	}

	var qna QnA
	rawId := params["id"]
	CollectionName := handlers.CollectionNameQnA(appid)
	if err := db.C(CollectionName).Find(bson.M{"_id": rawId}).One(&qna); err != nil {
		r.JSON(handlers.HttpErr(http.StatusNotFound, err.Error()))
		return
	}

	r.JSON(http.StatusOK, qna)
}

func UpdateQnA(req *http.Request, params martini.Params, qna QnA, r render.Render, db *mgo.Database) {
	appid := req.Header.Get("Application-Id")
	if appid == "" {
		r.JSON(handlers.HttpErr(http.StatusNotFound, "insert to Application-Id"))
		return
	}

	qna.Id = ""
	rawId := params["id"]
	colQuerier := bson.M{"_id": rawId}
	change := bson.M{"$set": qna}
	CollectionName := handlers.CollectionNameQnA(appid)
	if err := db.C(CollectionName).Update(colQuerier, change); err != nil {
		r.JSON(handlers.HttpErr(http.StatusNotFound, "NotFound "+rawId))
		return
	}

	r.JSON(http.StatusOK, "UPDATE_OK")
}

func DeleteQnA(req *http.Request, params martini.Params, r render.Render, db *mgo.Database) {
	appid := req.Header.Get("Application-Id")
	if appid == "" {
		r.JSON(handlers.HttpErr(http.StatusNotFound, "insert to Application-Id"))
		return
	}

	rawId := params["id"]
	CollectionName := handlers.CollectionNameQnA(appid)
	if err := db.C(CollectionName).Remove(bson.M{"_id": rawId}); err != nil {
		r.JSON(handlers.HttpErr(http.StatusNotFound, "NotFound "+rawId))
		return
	}

	r.JSON(http.StatusOK, "DELETE_OK")
}
