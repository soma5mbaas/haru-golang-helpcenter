package controllers

// import (
// 	"../handlers"
// 	"code.google.com/p/go-uuid/uuid"
// 	"github.com/go-martini/martini"
// 	"github.com/martini-contrib/render"
// 	"gopkg.in/mgo.v2"
// 	"gopkg.in/mgo.v2/bson"
// 	"net/http"
// 	"time"
// )

// type Email struct {
// 	Id        string `bson:"_id,omitempty"`            // UUID
// 	Title     string `bson:"title" binding:"required"` // 제목
// 	Body      string `bson:"body" binding:"required"`  // 본문
// 	Time      int64  `bson:"time,omitempty" `          // 시간
// 	Reception bool   `bson:"reception,omitempty" `     // 읽기여부
// 	URL       string `bson:"url,omitempty" `           // Image URL
// }

// func CreateNotice(req *http.Request, params martini.Params, notice Notice, r render.Render, db *mgo.Database) {
// 	appid := req.Header.Get("Application-Id")
// 	if appid == "" {
// 		r.JSON(http.StatusNotFound, "insert to Application-Id")
// 		return
// 	}

// 	notice.Time = time.Now().Unix()
// 	notice.Id = uuid.New()
// 	notice.Reception = false
// 	CollectionName := handlers.CollectionNameNotice(appid)
// 	err := db.C(CollectionName).Insert(notice)
// 	if err != nil {
// 		r.JSON(http.StatusNotFound, err)
// 		return
// 	}
// 	r.JSON(http.StatusOK, map[string]interface{}{"err": err, "Notice": notice})
// }
