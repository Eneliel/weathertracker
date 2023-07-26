package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	weather "weather/internal"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Start(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", nil)
}

func Regist(c echo.Context) error {
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
		log.Printf("Проблема с записью в БД:%s", err)
		return c.Render(http.StatusOK, "errpage.html", err)
	}
	defer ins_db.Close()
	return c.Render(http.StatusOK, "regist_done.html", username)
}

func Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	u := User{}

	db, err := sql.Open("mysql", "root:@/test_auth")
	if err != nil {
		log.Printf("Проблема со входом в БД:%s", err)
		return c.Render(http.StatusOK, "errorpage.html", err)
	}
	defer db.Close()

	sel, err := db.Query("SELECT `username`, `password` FROM `user`")
	if err != nil {
		log.Printf("Проблема с запросом SELECT:%s", err)
		return c.Render(http.StatusOK, "errpage.html", err)
	}
	for sel.Next() {
		err = sel.Scan(&u.Username, &u.Password)
		if err != nil {
			log.Printf("Ошибка при чтении с БД:%s", err)
			return c.Render(http.StatusOK, "errpage.html", err)
		}
		if username == u.Username && password == u.Password {
			return c.Render(http.StatusOK, "login_done.html", username)
		}
		if username == u.Username && password != u.Password {
			return c.Render(http.StatusOK, "errpage.html", "Не правильный пароль!")
		}
	}
	return c.Render(http.StatusOK, "errpage.html", "Пользователь не зарегистрирован!")
}

func Weather(c echo.Context) error {
	city := c.FormValue("city")
	w, err := weather.Query(city)
	if err != nil {
		return c.Render(http.StatusOK, "errpage.html", err)
	}
	return c.Render(http.StatusOK, "main.html", w)
}

func MainPage(c echo.Context) error {
	return c.Render(http.StatusOK, "main.html", nil)
}
