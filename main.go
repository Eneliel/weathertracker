package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
)

type User struct {
	Username string
	Password string
}

func main() {

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", nil)
	})

	e.POST("/register", func(c echo.Context) error {
		username := c.FormValue("username")
		password := c.FormValue("password")

		db, err := sql.Open("mysql", "root:@/test_auth")
		if err != nil {
			log.Printf("Проблема со входом в БД:%s", err)
			return c.Render(http.StatusOK, "errpage.html", err)
		}
		defer db.Close()
		insert := fmt.Sprintf("INSERT INTO `user`(`username`, `password`) VALUES ('%s','%s')", username, password)
		ins_db, err := db.Query(insert)
		if err != nil {
			log.Printf("Проблема с Записью:%s", err)
			return c.Render(http.StatusOK, "errpage.html", err)
		}
		defer ins_db.Close()

		return c.Render(http.StatusOK, "register_done.html", username)
	})

	e.POST("/login", func(c echo.Context) error {
		username := c.FormValue("username")
		password := c.FormValue("password")

		u := User{}
		db, err := sql.Open("mysql", "root:@/test_auth")
		if err != nil {
			log.Printf("Проблема со входом в БД:%s", err)
			return c.Render(http.StatusOK, "errpage.html", err)
		}
		defer db.Close()
		sel, err := db.Query("SELECT `username`, `password` FROM `user`")
		if err != nil {
			log.Printf("Ошибка с запросом SELECT:%s", err)
			return c.Render(http.StatusOK, "errpage.html", err)
		}
		for sel.Next() {
			err = sel.Scan(&u.Username, &u.Password)
			if err != nil {
				log.Printf("Ошибка при чтении с БД:%s", err)
				return c.Render(http.StatusOK, "errpage.html", err)
			}
			if u.Username == username && u.Password == password {
				return c.Render(http.StatusOK, "login_done.html", u)
			}
		}

		return c.Render(http.StatusOK, "errpage.html", "Пользователь не найден")
	})

	e.POST("/weather", func(c echo.Context) error {
		city := c.FormValue("city")
		w, err := Query(city)
		if err != nil {
			return c.Render(http.StatusOK, "errpage.html", err)
		}

		return c.Render(http.StatusOK, "main.html", w)
	})

	e.GET("/main", func(c echo.Context) error {
		return c.Render(http.StatusOK, "main.html", nil)
	})

	e.Renderer = &TemplateRenderer{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}

	e.Start(":8080")
}

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
