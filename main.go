package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"log"
	"firebase.google.com/go"
	"google.golang.org/api/option"
	"golang.org/x/net/context"
	"cloud.google.com/go/firestore"
	"github.com/gin-contrib/cors"
)

func main() {
	// Use a service account
	ctx := context.Background()
	sa := option.WithCredentialsFile("kunci.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	//log.Println(port)

	router := gin.New()
	router.Use(cors.Default())
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*")

	router.GET("/data", func(c *gin.Context) {
		data, err := getKata(ctx, client)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
		}

		c.JSON(http.StatusOK, data)
	})
	router.GET("/jtweets", func(c *gin.Context) {
		data := getJtweets(ctx, client)
		c.JSON(http.StatusOK, data)
	})

	router.GET("/jkata", func(c *gin.Context) {
		data := getJKata(ctx, client)
		c.JSON(http.StatusOK, data)
	})

	router.GET("/index",func(c *gin.Context) {
		c.HTML(200,"index.html","")
	})

	//Run
	router.Run(":" + port)
}

func getKata(ctx context.Context, client *firestore.Client) (map[string]interface{}, error) {
	ref, err := client.Collection("kata_kunci").Doc("kata").Get(ctx)
	if err != nil {
		log.Printf("Failed: %v", err)
	}

	return ref.Data(), err
}

func getJtweets(ctx context.Context, client *firestore.Client) interface{} {
	ref, err := client.Collection("lasttweet").Doc("tweet").Get(ctx)
	if err != nil {
		log.Printf("Failed: %v", err)
	}

	j, err:=ref.DataAt("total")
	if err != nil {
		log.Printf("Failed: %v", err)
	}

	return j
}

func getJKata(ctx context.Context, client *firestore.Client) interface{} {
	ref, err := client.Collection("kata_kunci").Doc("banyak").Get(ctx)
	if err != nil {
		log.Printf("Failed: %v", err)
	}

	j, err:=ref.DataAt("total")
	if err != nil {
		log.Printf("Failed: %v", err)
	}

	return j
}
