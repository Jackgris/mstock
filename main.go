package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

func main() {
	m := martini.Classic()
	// We change the limiters to avoid problems with AngularJs
	m.Use(render.Renderer(render.Options{
		Delims: render.Delims{"<<<", ">>>"},
	}))

	m.Get("/", func(r render.Render) {
		r.HTML(200, "home", "prueba")
	})

	m.RunOnAddr(":8080")
	m.Run()
}
