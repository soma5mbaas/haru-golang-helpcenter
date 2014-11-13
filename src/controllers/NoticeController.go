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
	Id        string `bson:"_id,omitempty" json:"_id"`              // UUID
	Title     string `bson:"title" binding:"required" json:"title"` // 제목
	Body      string `bson:"body" binding:"required" json:"body"`   // 본문
	Time      int64  `bson:"time,omitempty" json:"time"`            // 시간
	Reception bool   `bson:"reception,omitempty" json:"reception"`  // 읽기여부
	URL       string `bson:"url,omitempty" json:"url"`              // Image URL
}

func CreateNotice(req *http.Request, params martini.Params, notice Notice, r render.Render, db *mgo.Database, f *log.Logger) {
	appid := req.Header.Get("Application-Id")
	if appid == "" {
		r.JSON(handlers.HttpErr(http.StatusNotFound, "insert to Application-Id"))
		return
	}

	notice.Time = time.Now().Unix()
	notice.Id = uuid.New()
	notice.Reception = false

	CollectionName := handlers.CollectionNameNotice(appid)
	// if count, _ := db.C(CollectionName).Count(); count > 0 {
	// 	if err := db.Session.DB("admin").Run(bson.M{"shardCollection": "haru" + "." + CollectionName, "key": bson.M{"_id": 1}}, nil); err != nil {
	// 		f.Println(CollectionName+" Sharde Fail :", err)
	// 	} else {
	// 		f.Println(CollectionName+" Sharde ok :", err)
	// 	}
	// }

	//list가 nil이면 empty list 반환
	if err := db.C(CollectionName).Insert(notice); err != nil {
		r.JSON(handlers.HttpErr(http.StatusNotFound, err.Error()))
		return
	}
	r.JSON(http.StatusOK, map[string]interface{}{"Notice": notice})
}

func ReadListNotice(req *http.Request, r render.Render, db *mgo.Database) {
	appid := req.Header.Get("Application-Id")
	if appid == "" {
		r.JSON(handlers.HttpErr(http.StatusNotFound, "insert to Application-Id"))
		return
	}
	colQuerier := bson.M{}
	//읽기여부 true로 바꿔준다.
	change := bson.M{"$set": bson.M{"reception": true}}
	CollectionName := handlers.CollectionNameNotice(appid)
	if _, err := db.C(CollectionName).UpdateAll(colQuerier, change); err != nil {
		r.JSON(handlers.HttpErr(http.StatusNotFound, err.Error()))
		return
	}
	var notices []Notice
	if err := db.C(CollectionName).Find(bson.M{}).Sort("-time").All(&notices); err != nil {
		r.JSON(handlers.HttpErr(http.StatusNotFound, err.Error()))
		return
	}

	// /////////
	// change := mgo.Change{
	// 	Update:    bson.M{"$set": bson.M{"reception": true}},
	// 	ReturnNew: false,
	// }
	// a, b := db.C(CollectionName).Find(nil).Select(bson.M{"_id": 1}).Apply(change, &notices)
	// fmt.Println(a, b)
	//////
	// expressions := make([]bson.M, len(notices))
	// for i := 0; i < len(notices); i++ {
	// 	notices[i].Reception = true
	// 	b, _ := json.Marshal(notices[i])
	// 	expressions[i] = b
	// }

	// tag := bson.M{}
	// tags := []bson.M{}
	// iter := db.C(CollectionName).Pipe(expressions).Iter()
	// // BUG: mgo is not returning consistent results. It seems to be only broken with Pipe().
	// for iter.Next(&tag) {
	// 	tags = append(tags, tag["_id"].(bson.M))
	// }
	// if err := iter.Close(); err != nil {
	// 	log.Fatalln(err)
	// 	return
	// }

	// pipe := db.C(CollectionName).Pipe(expressions)
	// iter := pipe.All(expressions)

	if notices == nil {
		r.JSON(http.StatusOK, map[string]interface{}{"return": bson.D{}})
		return
	}

	r.JSON(http.StatusOK, map[string]interface{}{"return": notices})
}

func ReadIdNotice(req *http.Request, params martini.Params, r render.Render, db *mgo.Database) {
	appid := req.Header.Get("Application-Id")
	if appid == "" {
		r.JSON(handlers.HttpErr(http.StatusNotFound, "insert to Application-Id"))
		return
	}
	rawId := params["id"]
	var notices Notice
	CollectionName := handlers.CollectionNameNotice(appid)
	if err := db.C(CollectionName).Find(bson.M{"_id": string(rawId)}).One(&notices); err != nil {
		r.JSON(handlers.HttpErr(http.StatusNotFound, "NotFound "+rawId))
		return
	}
	r.JSON(http.StatusOK, map[string]interface{}{"return": notices})
}

func UpdateNotice(req *http.Request, params martini.Params, notice Notice, r render.Render, db *mgo.Database) {
	appid := req.Header.Get("Application-Id")
	if appid == "" {
		r.JSON(handlers.HttpErr(http.StatusNotFound, "insert to Application-Id"))
		return
	}
	rawId := params["id"]
	notice.Id = ""
	colQuerier := bson.M{"_id": rawId}
	change := bson.M{"$set": notice}
	CollectionName := handlers.CollectionNameNotice(appid)
	if err := db.C(CollectionName).Update(colQuerier, change); err != nil {
		r.JSON(handlers.HttpErr(http.StatusNotFound, "NotFound "+rawId))
		return
	}

	r.JSON(http.StatusOK, "UPDATE_OK")
}

func DeleteNotice(c http.ResponseWriter, req *http.Request, params martini.Params, r render.Render, db *mgo.Database) {
	appid := req.Header.Get("Application-Id")
	if appid == "" {
		r.JSON(handlers.HttpErr(http.StatusNotFound, "insert to Application-Id"))
		return
	}
	rawId := params["id"]
	CollectionName := handlers.CollectionNameNotice(appid)
	if err := db.C(CollectionName).Remove(bson.M{"_id": rawId}); err != nil {
		r.JSON(handlers.HttpErr(http.StatusNotFound, "NotFound "+rawId))
		return
	}

	r.JSON(http.StatusOK, "DELETE_OK")
}
