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

type FaqCategory struct {
	Id       string `bson:"_id,omitempty" json:"_id"`   // UUID
	Category string `bson:"category" json:"category"`   // Category
	Time     int64  `bson:"time,omitempty" json:"time"` // 시간
}

func CreateFaqCategory(req *http.Request, params martini.Params, fa FaqCategory, r render.Render, db *mgo.Database, f *log.Logger) {

	appid := req.Header.Get("Application-Id")
	if appid == "" {
		r.JSON(handlers.HttpErr(http.StatusNotFound, "insert to Application-Id"))
		return
	}

	fa.Time = time.Now().Unix()
	fa.Id = uuid.New()

	CollectionName := handlers.CollectionNameFAQCategory(appid)
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

func ReadListFaqCategory(req *http.Request, r render.Render, db *mgo.Database) {
	appid := req.Header.Get("Application-Id")
	if appid == "" {
		r.JSON(handlers.HttpErr(http.StatusNotFound, "insert to Application-Id"))
		return
	}

	var faqs []FaqCategory
	CollectionName := handlers.CollectionNameFAQCategory(appid)
	err := db.C(CollectionName).Find(bson.M{}).All(&faqs)

	if err != nil {
		r.JSON(handlers.HttpErr(http.StatusNotFound, err.Error()))
		return
	}
	if faqs == nil {
		r.JSON(http.StatusOK, map[string]interface{}{"return": bson.D{}})
		return
	}
	r.JSON(http.StatusOK, map[string]interface{}{"return": faqs})
}

func ReadIdFaqCategory(req *http.Request, params martini.Params, r render.Render, db *mgo.Database) {

	appid := req.Header.Get("Application-Id")
	if appid == "" {
		r.JSON(handlers.HttpErr(http.StatusNotFound, "insert to Application-Id"))
		return
	}

	var fa FaqCategory
	rawId := params["id"]
	CollectionName := handlers.CollectionNameFAQCategory(appid)
	err := db.C(CollectionName).Find(bson.M{"_id": rawId}).One(&fa)

	if err != nil {
		r.JSON(handlers.HttpErr(http.StatusNotFound, err.Error()))
		return
	}

	r.JSON(http.StatusOK, fa)
}

func CountFaqCategory(req *http.Request, params martini.Params, r render.Render, db *mgo.Database) {

	appid := req.Header.Get("Application-Id")
	if appid == "" {
		r.JSON(handlers.HttpErr(http.StatusNotFound, "insert to Application-Id"))
		return
	}

	var faqs []FaqCategory
	rawId := params["id"]
	CollectionName := handlers.CollectionNameFAQCategory(appid)
	err := db.C(CollectionName).Find(bson.M{"category": rawId}).All(&faqs)
	if err != nil {
		r.JSON(handlers.HttpErr(http.StatusNotFound, err.Error()))
		return
	}
	r.JSON(http.StatusOK, map[string]interface{}{"return": len(faqs)})
}

func UpdateFaqCategory(req *http.Request, params martini.Params, fa FaqCategory, r render.Render, db *mgo.Database) {
	appid := req.Header.Get("Application-Id")
	if appid == "" {
		r.JSON(handlers.HttpErr(http.StatusNotFound, "insert to Application-Id"))
		return
	}

	fa.Id = ""
	rawId := params["id"]

	colQuerier := bson.M{"_id": rawId}
	change := bson.M{"$set": fa}
	CollectionName := handlers.CollectionNameFAQCategory(appid)
	err := db.C(CollectionName).Update(colQuerier, change)
	if err != nil {
		r.JSON(handlers.HttpErr(http.StatusNotFound, "NotFound "+rawId))
		return
	}

	r.JSON(http.StatusOK, "UPDATE_OK")
}

func DeleteFaqCategory(req *http.Request, params martini.Params, r render.Render, db *mgo.Database) {
	appid := req.Header.Get("Application-Id")
	if appid == "" {
		r.JSON(handlers.HttpErr(http.StatusNotFound, "insert to Application-Id"))
		return
	}

	rawId := params["id"]
	CollectionName := handlers.CollectionNameFAQCategory(appid)
	err := db.C(CollectionName).Remove(bson.M{"_id": rawId})
	if err != nil {
		r.JSON(handlers.HttpErr(http.StatusNotFound, "NotFound "+rawId))
		return
	}

	r.JSON(http.StatusOK, "DELETE_OK")
}
