package main

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var pubKey *rsa.PublicKey

func readJWTKey(jwtKeyFilename string) *rsa.PublicKey {
	raw, err := os.ReadFile(jwtKeyFilename)
	if err != nil {
		panic("failed to read public key file" + err.Error())
	}
	pub, err := x509.ParsePKIXPublicKey(raw)
	if err != nil {
		panic("failed to parse DER encoded public key: " + err.Error())
	}
	return pub.(*rsa.PublicKey)
}

func parseToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return pubKey, nil
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

/*
	handlers accessible to auth users
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

	c.HTML(http.StatusOK, "chatclient.html", nil)
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
	c.HTML(http.StatusOK, "index.html", nil)
}

func getEnvVar(name string, dflt string) string {
	if val, ok := os.LookupEnv(name); ok {
		return val
	}
	return dflt
}

func main() {
	htmlDir := getEnvVar("HTML_DIR", "chat-client/build")
	jwtKeyFilename := getEnvVar("JWTKEY_FILENAME", "public_key.der")
	dbHost := getEnvVar("POSTGRES_HOST", "localhost")
	dbPort := getEnvVar("POSTGRES_PORT", "5432")
	dbName := getEnvVar("POSTGRES_DB", "gochat")
	dbUser := getEnvVar("POSTGRES_USERNAME", "gochat")
	dbPassword := getEnvVar("POSTGRES_PASSWORD", "gochat")

	pubKey = readJWTKey(jwtKeyFilename);

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName, dbPort,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
		return
	}

	r := gin.Default()
	r.LoadHTMLGlob(fmt.Sprintf("%s/*.html",htmlDir))
	r.Static("/static", fmt.Sprintf("%s/static",htmlDir))

	r.GET("/", index)
	r.GET("/login", login)
	r.GET("/logout", logout)
	r.GET("/home", home)
	r.GET(
		"/api/users", 
		func(c *gin.Context){
			users(c,db)
		},
	)

	r.Run(":8000")
}