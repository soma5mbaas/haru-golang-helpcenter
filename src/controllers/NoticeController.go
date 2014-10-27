package controllers

import (
	"code.google.com/p/go-uuid/uuid"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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

func CreateNotice(params martini.Params, notice Notice, r render.Render, db *mgo.Database) {
	notice.Time = time.Now().Unix()
	notice.Id = uuid.New()
	notice.Reception = false
	err := db.C("notice").Insert(notice)
	r.JSON(http.StatusOK, map[string]interface{}{"err": err, "Notice": notice})
}

func ReadListNotice(r render.Render, db *mgo.Database) {
	var notices []Notice
	err := db.C("notice").Find(bson.M{}).All(&notices)
	if err != nil {
		r.JSON(http.StatusNotFound, err)
		return
	}
	r.JSON(http.StatusOK, notices)
}

func ReadIdNotice(params martini.Params, r render.Render, db *mgo.Database) {
	rawId := params["id"]
	var notices Notice
	err := db.C("notice").Find(bson.M{"_id": string(rawId)}).One(&notices)
	if err != nil {
		r.JSON(http.StatusNotFound, err)
		return
	}
	r.JSON(http.StatusOK, notices)
}

func UpdateNotice(params martini.Params, notice Notice, r render.Render, db *mgo.Database) {
	rawId := params["id"]
	notice.Id = ""
	colQuerier := bson.M{"_id": rawId}
	change := bson.M{"$set": notice}
	err := db.C("notice").Update(colQuerier, change)
	if err != nil {
		r.JSON(http.StatusNotFound, err)
		return
	}

	r.JSON(http.StatusOK, "UPDATE_OK")
}

func DeleteNotice(params martini.Params, r render.Render, db *mgo.Database) {
	rawId := params["id"]
	err := db.C("notice").Remove(bson.M{"_id": rawId})
	if err != nil {
		r.JSON(http.StatusNotFound, "NotFound "+rawId)
		return
	}
	r.JSON(http.StatusOK, "DELETE_OK")
}
