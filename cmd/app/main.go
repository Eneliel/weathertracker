package main

import (
	"html/template"
	"io"
	handler "weather/cmd/handlers"

	"github.com/labstack/echo"
)

func main() {

	e := echo.New()

	e.GET("/", handler.Start)
	e.POST("/regist", handler.Regist)
	e.POST("/login", handler.Login)
	e.GET("/main", handler.MainPage)
	e.POST("/weather", handler.Weather)

	e.Renderer = &TemplateRenderer{
		temp: template.Must(template.ParseGlob("templates/*.html")),
	}
	e.Start(":8080")
}

type TemplateRenderer struct {
	temp *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.temp.ExecuteTemplate(w, name, data)
}
