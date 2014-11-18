package controllers

import (
	"../../src"
	"../handlers"
	"code.google.com/p/go-uuid/uuid"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"

	"net/url"
	"time"
)

type QnA struct {
	Id           string `bson:"_id,omitempty" json:"_id"`             // UUID
	EmailAddress string `bson:"emailaddress" json:"emailaddress"`     // Email
	Body         string `bson:"body" json:"body"`                     // 본문
	Time         int64  `bson:"time" json:"time"`                     // 시간
	Reception    bool   `bson:"reception,omitempty" json:"reception"` // 읽기여부
	Comment      string `bson:"comment" json:"comment"`               // Comment
	CommentTime  int64  `bson:"commenttime" json:"commenttime"`       // Comment 단 시간
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

	//web notify
	var request string = config.DASHBOARD_WEB + "/qna/webhook?appid="
	request += appid + "&"
	v := url.Values{}
	v.Add("body", qna.Body)
	request += v.Encode()
	_, Geterr := http.Get(request)
	if Geterr != nil {
		log.Println(Geterr)
	}

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

	if qnas == nil {
		r.JSON(http.StatusOK, map[string]interface{}{"return": bson.D{}})
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

	if qnas == nil {
		r.JSON(http.StatusOK, map[string]interface{}{"return": bson.D{}})
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

//////////////
func ReadCountQnA(req *http.Request, params martini.Params, r render.Render, db *mgo.Database) {

	appid := req.Header.Get("Application-Id")
	if appid == "" {
		r.JSON(handlers.HttpErr(http.StatusNotFound, "insert to Application-Id"))
		return
	}

	var qnas []QnA
	CollectionName := handlers.CollectionNameQnA(appid)
	if err := db.C(CollectionName).Find(nil).All(&qnas); err != nil { //현재 Month로 카운트 하는거 나중에 구현하기!!
		r.JSON(handlers.HttpErr(http.StatusNotFound, err.Error()))
		return
	}

	r.JSON(http.StatusOK, map[string]interface{}{"return": len(qnas)})
}

//////////

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
