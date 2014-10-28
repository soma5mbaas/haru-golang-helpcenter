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

type Notice struct {
	Id        string `bson:"_id,omitempty"`            // UUID
	Title     string `bson:"title" binding:"required"` // 제목
	Body      string `bson:"body" binding:"required"`  // 본문
	Time      int64  `bson:"time,omitempty" `          // 시간
	Reception bool   `bson:"reception,omitempty" `     // 읽기여부
	URL       string `bson:"url,omitempty" `           // Image URL
}

func CreateNotice(req *http.Request, params martini.Params, notice Notice, r render.Render, db *mgo.Database, f *log.Logger) {
	appid := req.Header.Get("Application-Id")
	if appid == "" {
		r.JSON(http.StatusNotFound, "insert to Application-Id")
		return
	}

	notice.Time = time.Now().Unix()
	notice.Id = uuid.New()
	notice.Reception = false

	CollectionName := handlers.CollectionNameNotice(appid)
	if count, _ := db.C(CollectionName).Count(); count > 0 {
		if err := db.Session.DB("admin").Run(bson.M{"shardCollection": "haru" + "." + CollectionName, "key": bson.M{"_id": 1}}, nil); err != nil {
			f.Println(CollectionName+" Sharde Fail :", err)
		} else {
			f.Println(CollectionName+" Sharde ok :", err)
		}
	}

	if err := db.C(CollectionName).Insert(notice); err != nil {
		r.JSON(http.StatusNotFound, err)
		return
	}
	r.JSON(http.StatusOK, map[string]interface{}{"Notice": notice})
}

func ReadListNotice(req *http.Request, r render.Render, db *mgo.Database) {
	appid := req.Header.Get("Application-Id")
	if appid == "" {
		r.JSON(http.StatusNotFound, "insert to Application-Id")
		return
	}
	var notices []Notice
	CollectionName := handlers.CollectionNameNotice(appid)

	if err := db.C(CollectionName).Find(bson.M{}).Sort("-time").All(&notices); err != nil {
		r.JSON(http.StatusNotFound, err)
		return
	}
	r.JSON(http.StatusOK, notices)
}

func ReadIdNotice(req *http.Request, params martini.Params, r render.Render, db *mgo.Database) {
	appid := req.Header.Get("Application-Id")
	if appid == "" {
		r.JSON(http.StatusNotFound, "insert to Application-Id")
		return
	}
	rawId := params["id"]
	var notices Notice
	CollectionName := handlers.CollectionNameNotice(appid)

	if err := db.C(CollectionName).Find(bson.M{"_id": string(rawId)}).One(&notices); err != nil {
		r.JSON(http.StatusNotFound, "NotFound "+rawId)
		return
	}
	r.JSON(http.StatusOK, notices)
}

func UpdateNotice(req *http.Request, params martini.Params, notice Notice, r render.Render, db *mgo.Database) {
	appid := req.Header.Get("Application-Id")
	if appid == "" {
		r.JSON(http.StatusNotFound, "insert to Application-Id")
		return
	}
	rawId := params["id"]
	notice.Id = ""
	colQuerier := bson.M{"_id": rawId}
	change := bson.M{"$set": notice}
	CollectionName := handlers.CollectionNameNotice(appid)

	if err := db.C(CollectionName).Update(colQuerier, change); err != nil {
		r.JSON(http.StatusNotFound, "NotFound "+rawId)
		return
	}

	r.JSON(http.StatusOK, "UPDATE_OK")
}

func DeleteNotice(req *http.Request, params martini.Params, r render.Render, db *mgo.Database) {
	appid := req.Header.Get("Application-Id")
	if appid == "" {
		r.JSON(http.StatusNotFound, "insert to Application-Id")
		return
	}
	rawId := params["id"]
	CollectionName := handlers.CollectionNameNotice(appid)

	if err := db.C(CollectionName).Remove(bson.M{"_id": rawId}); err != nil {
		r.JSON(http.StatusNotFound, "NotFound "+rawId)
		return
	}
	r.JSON(http.StatusOK, "DELETE_OK")
}
