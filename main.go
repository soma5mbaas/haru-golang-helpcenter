package main

import (
	"./src/controllers"
	"./src/encoding"
	"./src/models"
	"./src/utility"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	martini.Env = martini.Prod
	m := martini.Classic()
	m.Use(render.Renderer())
	m.Use(encoding.MapEncoder)
	m.Use(models.InitDB())

	f := utility.InitLogger(m)
	defer f.Close()

	m.Group("/notice", func(r martini.Router) {
		r.Get("/list", controllers.ReadListNotice)
		r.Get("/:id", controllers.ReadIdNotice)
		r.Post("/add", binding.Json(controllers.Notice{}), controllers.CreateNotice)
		r.Put("/:id", binding.Json(controllers.Notice{}), controllers.UpdateNotice)
		r.Delete("/:id", binding.Json(controllers.Notice{}), controllers.DeleteNotice)
	})

	m.Group("/faq/category", func(r martini.Router) {
		r.Get("/list", controllers.ReadListFaqCategory)
		r.Get("/:id", controllers.ReadIdFaqCategory)
		r.Post("/add", binding.Json(controllers.FaqCategory{}), controllers.CreateFaqCategory)
		r.Put("/:id", binding.Json(controllers.FaqCategory{}), controllers.UpdateFaqCategory)
		r.Delete("/:id", binding.Json(controllers.FaqCategory{}), controllers.DeleteFaqCategory)
	})

	m.Group("/faq", func(r martini.Router) {
		r.Get("/list", controllers.ReadListFaq)
		r.Get("/list/:id", binding.Json(controllers.Faq{}), controllers.ReadListCategoryFaq)
		r.Get("/:id", controllers.ReadIdFaq)
		r.Post("/add", binding.Json(controllers.Faq{}), controllers.CreateFaq)
		r.Put("/:id", binding.Json(controllers.Faq{}), controllers.UpdateFaq)
		r.Delete("/:id", binding.Json(controllers.Faq{}), controllers.DeleteFaq)
	})

	m.Group("/email", func(r martini.Router) {
		r.Post("/send", binding.Json(controllers.Faq{}), controllers.CreateFaq)
	})

	m.Run()
}
