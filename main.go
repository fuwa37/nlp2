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
	"google.golang.org/api/iterator"
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
	router.GET("/all", func(c *gin.Context) {
		data := getAll(ctx, client)
		c.JSON(200, data)
	})

	router.GET("/last", func(c *gin.Context) {
		data := getLast(ctx, client)
		c.JSON(200, data)
	})

	router.GET("/jtweets", func(c *gin.Context) {
		data := getJtweets(ctx, client)
		c.JSON(http.StatusOK, data)
	})

	router.GET("/jkata", func(c *gin.Context) {
		data := getJKata(ctx, client)
		c.JSON(http.StatusOK, data)
	})

	router.GET("/index", func(c *gin.Context) {
		c.HTML(200, "index.html", "")
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

func getLast(ctx context.Context, client *firestore.Client) map[string]interface{} {
	ref, err := client.Collection("lasttweet").Doc("last").Get(ctx)
	if err != nil {
		log.Printf("Failed: %v", err)
	}

	last, err := ref.DataAt("ref")
	if err != nil {
		log.Printf("Failed: %v", err)
	}

	lastref := last.(*firestore.DocumentRef)

	dataref, err := lastref.Get(ctx)
	if err != nil {
		log.Printf("Failed: %v", err)
	}

	data := dataref.Data()

	return data
}

func getJtweets(ctx context.Context, client *firestore.Client) interface{} {
	ref, err := client.Collection("lasttweet").Doc("tweet").Get(ctx)
	if err != nil {
		log.Printf("Failed: %v", err)
	}

	j, err := ref.DataAt("total")
	if err != nil {
		log.Printf("Failed: %v", err)
	}

	return j
}

func getAll(ctx context.Context, client *firestore.Client) map[string]interface{} {
	iter := client.Collection("Tweet").Documents(ctx)
	type data struct {
		username string `json:"username"`
		text     string `json:"text"`
		time     string `json:"time"`
	}
	j := make(map[string]interface{})
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Failed to iterate: %v", err)
		}
		j[doc.Ref.ID] = doc.Data()
	}

	return j
}

func getJKata(ctx context.Context, client *firestore.Client) interface{} {
	ref, err := client.Collection("kata_kunci").Doc("banyak").Get(ctx)
	if err != nil {
		log.Printf("Failed: %v", err)
	}

	j, err := ref.DataAt("total")
	if err != nil {
		log.Printf("Failed: %v", err)
	}

	return j
}
