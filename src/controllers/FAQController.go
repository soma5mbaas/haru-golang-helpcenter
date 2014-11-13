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

type Faq struct {
	Id       string `bson:"_id,omitempty" json:"_id"`   // UUID
	Title    string `bson:"title" json:"title"`         // 제목
	Body     string `bson:"body" json:"body"`           // 본문
	Category string `bson:"category" json:"category"`   // Category
	Platform string `bson:"platform" json:"platform"`   // Platform
	Time     int64  `bson:"time,omitempty" json:"time"` // 시간
}

func CreateFaq(req *http.Request, params martini.Params, fa Faq, r render.Render, db *mgo.Database, f *log.Logger) {
	appid := req.Header.Get("Application-Id")
	if appid == "" {
		r.JSON(handlers.HttpErr(http.StatusNotFound, "insert to Application-Id"))
		return
	}

	fa.Time = time.Now().Unix()
	fa.Id = uuid.New()

	CollectionName := handlers.CollectionNameFAQ(appid)
	// if count, _ := db.C(CollectionName).Count(); count > 0 {
	// 	if err := db.Session.DB("admin").Run(bson.M{"shardCollection": "haru" + "." + CollectionName, "key": bson.M{"_id": 1}}, nil); err != nil {
	// 		f.Println(CollectionName+" Sharde Fail :", err)
	// 	} else {
	// 		f.Println(CollectionName+" Sharde ok :", err)
	// 	}
	// }

	if err := db.C(CollectionName).Insert(fa); err != nil {
		r.JSON(handlers.HttpErr(http.StatusNotFound, err.Error()))
		return
	}
	r.JSON(http.StatusOK, map[string]interface{}{"Faq": fa})
}

func ReadListCategoryFaq(req *http.Request, params martini.Params, r render.Render, db *mgo.Database) {
	appid := req.Header.Get("Application-Id")
	if appid == "" {
		r.JSON(http.StatusNotFound, "insert to Application-Id")
		return
	}

	var faqs []Faq
	rawId := params["category"]
	CollectionName := handlers.CollectionNameFAQ(appid)
	err := db.C(CollectionName).Find(bson.M{"category": rawId}).All(&faqs)
	if err != nil {
		r.JSON(handlers.HttpErr(http.StatusNotFound, err.Error()))
		return
	}
	if len(faqs) == 0 {
		r.JSON(http.StatusOK, map[string]interface{}{"return": bson.D{}})
		return
	}
	r.JSON(http.StatusOK, map[string]interface{}{"return": faqs})
}

func ReadListFaq(req *http.Request, r render.Render, db *mgo.Database) {
	appid := req.Header.Get("Application-Id")
	if appid == "" {
		r.JSON(handlers.HttpErr(http.StatusNotFound, "insert to Application-Id"))
		return
	}

	var faqs []Faq
	CollectionName := handlers.CollectionNameFAQ(appid)
	err := db.C(CollectionName).Find(bson.M{}).All(&faqs)
	if err != nil {
		r.JSON(handlers.HttpErr(http.StatusNotFound, err.Error()))
		return
	}
	if len(faqs) == 0 {
		r.JSON(http.StatusOK, map[string]interface{}{"return": bson.D{}})
		return
	}
	r.JSON(http.StatusOK, map[string]interface{}{"return": faqs})
}

func ReadIdFaq(req *http.Request, params martini.Params, r render.Render, db *mgo.Database) {

	appid := req.Header.Get("Application-Id")
	if appid == "" {
		r.JSON(handlers.HttpErr(http.StatusNotFound, "insert to Application-Id"))
		return
	}

	var fa Faq
	rawId := params["id"]
	CollectionName := handlers.CollectionNameFAQ(appid)
	err := db.C(CollectionName).Find(bson.M{"_id": rawId}).One(&fa)
	if err != nil {
		r.JSON(handlers.HttpErr(http.StatusNotFound, err.Error()))
		return
	}

	r.JSON(http.StatusOK, fa)
}

func UpdateFaq(req *http.Request, params martini.Params, fa Faq, r render.Render, db *mgo.Database) {
	appid := req.Header.Get("Application-Id")
	if appid == "" {
		r.JSON(handlers.HttpErr(http.StatusNotFound, "insert to Application-Id"))
		return
	}

	fa.Id = ""
	rawId := params["id"]

	colQuerier := bson.M{"_id": rawId}
	change := bson.M{"$set": fa}
	CollectionName := handlers.CollectionNameFAQ(appid)
	err := db.C(CollectionName).Update(colQuerier, change)
	if err != nil {
		r.JSON(handlers.HttpErr(http.StatusNotFound, "NotFound "+rawId))
		return
	}

	r.JSON(http.StatusOK, "UPDATE_OK")
}

func DeleteFaq(req *http.Request, params martini.Params, r render.Render, db *mgo.Database) {
	appid := req.Header.Get("Application-Id")
	if appid == "" {
		r.JSON(handlers.HttpErr(http.StatusNotFound, "insert to Application-Id"))
		return
	}

	rawId := params["id"]
	CollectionName := handlers.CollectionNameFAQ(appid)
	err := db.C(CollectionName).Remove(bson.M{"_id": rawId})
	if err != nil {
		r.JSON(handlers.HttpErr(http.StatusNotFound, "NotFound "+rawId))
		return
	}

	r.JSON(http.StatusOK, "DELETE_OK")
}
