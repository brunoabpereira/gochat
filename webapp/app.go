package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/golang-jwt/jwt/v5"
	"encoding/base64"
	// "encoding/json"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"fmt"
)

var base64HmacSecret = "26SrjQKKdr3Av2S04thIfsXcx4lSInVGjBYk5kUZrlSYFZfmGUZ9t9pcY8Rv8J2026SrjQKKdr3Av2S04thIfsXcx4lSInVGjBYk5kUZrlSYFZfmGUZ9t9pcY8Rv8J20"

func parseToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		hmacSecret, _ := base64.StdEncoding.DecodeString(base64HmacSecret)
		return hmacSecret, nil
	})

	return token, err
}

func verifyUser(c *gin.Context) (jwt.MapClaims, bool) {
	jwtid, err := c.Cookie("JWTID")

	if err != nil {
		log.Println("Cookie \"JWTID\" not set")
		return nil, false
	}

	token, err := parseToken(jwtid)
	if err != nil{
		log.Println(err)
		return nil, false
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && !token.Valid {
		log.Println("Cookie \"JWTID\" is not valid")
		return nil, false
	}

	return claims, true
}

func proxyHandler(upstream string) func (c *gin.Context){
	return func (c *gin.Context) {
		remote, err := url.Parse(upstream)
		if err != nil {
			log.Println(err)
			return
		}
	
		proxy := httputil.NewSingleHostReverseProxy(remote)
		proxy.Director = func(req *http.Request) {
			req.Header = c.Request.Header
			req.Host = remote.Host
			req.URL.Scheme = remote.Scheme
			req.URL.Host = remote.Host
			req.URL.Path = c.Request.URL.Path
		}
	
		log.Printf("Proxy request %s -> %s \n",c.Request.URL,remote.Host+c.Request.URL.Path)
	
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

/*
	handlers accessable to auth users
*/

func index(c *gin.Context) {
	_, auth := verifyUser(c)
	if !auth {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	c.Redirect(http.StatusFound, "/home")
}

func home(c *gin.Context) {
	_, auth := verifyUser(c)
	if !auth {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// session := sessions.Default(c)
	// session.Set("test","ok")
	// session.Save()

	c.HTML(http.StatusOK, "index.html", nil)
}

func logout(c *gin.Context) {
	_, auth := verifyUser(c)
	if !auth {
		return
	}

	cookie, err := c.Request.Cookie("JWTID")
	
	if err != nil {
		log.Println("Cookie \"JWTID\" not set")
		return
	}

	c.SetCookie(
		cookie.Name,
		cookie.Value,
		-1,
		cookie.Path,
		cookie.Domain,
		cookie.Secure,
		cookie.HttpOnly,
	)

	c.Redirect(http.StatusFound, "/")
}

/*
	API
*/

func users(c *gin.Context, db *gorm.DB) {
	claims, auth := verifyUser(c)
	if !auth {
		c.Status(http.StatusUnauthorized)
		return
	}
	
	if username, ok := claims["sub"]; !ok {
		return
	}else{
		response := make(map[string]interface{})

		userRes :=  make(map[string]interface{})
		err := db.Table("users").Select("userid","username","useremail").Where("username = ?",username).Take(&userRes).Error
		if err != nil {
			log.Println(err)
			return
		}
		response["username"] = userRes["username"]
		response["useremail"] = userRes["useremail"]
		
		channelsRes := make([]map[string]interface{},0)
		err = db.Raw("SELECT * FROM getUsersChannels(?)",userRes["username"]).Scan(&channelsRes).Error
		if err != nil {
			log.Println(err)
			return
		}
		response["channels"] = channelsRes

		c.JSON(200,response)
	}
}

/*	
	handlers not accessable to auth users
*/

func login(c *gin.Context) {
	_, auth := verifyUser(c)
	if auth {
		c.Redirect(http.StatusFound, "/")
		return
	}
	c.HTML(http.StatusOK, "login.html", nil)
}

func register(c *gin.Context) {
	_, auth := verifyUser(c)
	if auth {
		c.Redirect(http.StatusFound, "/")
		return
	}
	c.HTML(http.StatusOK, "register.html", nil)
}

func main() {
	var dsn string = "host=localhost user=gochat password=gochat dbname=gochat port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
		return
	}

	r := gin.Default()
	r.LoadHTMLGlob("chat-client/build/*.html")
	r.Static("/static", "./chat-client/build/static")

	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("session", store))

	r.GET("/", index)
	r.GET("/login", login)
	r.GET("/logout", logout)
	r.GET("/register", register)
	r.GET("/home", home)
	r.GET(
		"/api/users", 
		func(c *gin.Context){
			users(c,db)
		},
	)

	// proxy requests to auth service
	proxy := proxyHandler("http://localhost:9999")
	r.POST("/api/authorize", proxy)
	r.POST("/api/users", proxy)

	r.Run(":8000")
}