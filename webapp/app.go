package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"log"
)

func index(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func register(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

func home(c *gin.Context) {
	c.HTML(http.StatusOK, "home.html", nil)
}

type AuthRequest struct {
	Username	string    `form:"username"`
	Password	string    `form:"password"`
}

func login(c *gin.Context) {
	var authReq AuthRequest
	if c.ShouldBind(&authReq) == nil {
		log.Printf("Username:%s, Password:%s\n",authReq.Username,authReq.Password)
	}
	c.String(http.StatusOK, "Success")
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	r.GET("/", index)
	r.GET("/register", register)
	r.GET("/home", home)
	r.POST("/login", login)

	r.Run(":8000")
}