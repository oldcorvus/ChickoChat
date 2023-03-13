package main

import (
	"chicko_chat/controllers"
	"chicko_chat/database"
	"chicko_chat/log"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"path/filepath"
	"sync"
	"text/template"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// represents a single template
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// ServeHTTP handles the HTTP request.
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", ":8080", "The addr of the application.")
	mongoURI := flag.String("mongoURI", "mongodb://localhost:27017", "Database hostname url")
	mongoDatabase := flag.String("mongoDatabase", "chicko_chat", "Database name")
	enableCredentials := flag.Bool("enableCredentials", false, "Enable the use of credentials for mongo connection")
	flag.Parse()

	_ = mongoDatabase
	// Create mongo client configuration
	co := options.Client().ApplyURI(*mongoURI)
	if *enableCredentials {
		co.Auth = &options.Credential{
			Username: os.Getenv("MONGODB_USERNAME"),
			Password: os.Getenv("MONGODB_PASSWORD"),
		}
	}

	// Establish database connection
	client, err := mongo.NewClient(co)
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	db := &database.ChatDatabase{
		Users:    client.Database("chicko_chat").Collection("users"),
		Messages: client.Database("chicko_chat").Collection("messages"),
		Rooms:    client.Database("chicko_chat").Collection("rooms"),
	}


	log.Printf("Database connection established")


	router := gin.Default()
	controller := controllers.Controller{
		DB: db,
	}
	websocketServer := &websocket.WsServer{
		controller : controller,

	}
	router.POST("/start", controller.StartConversationApi) 
	router.POST("/user-rooms/", controller.GetUserRoomsApi) 
	router.POST("/create-room/", controller.CreateRoomApi) 
	router.POST("/room-history/", controller.RoomHistoryApi)
	router.POST("/add-user-room/", controller.AddUserToRoomApi)
	router.GET("/ws/:roomId/:userId/", func(c *gin.Context) {
		roomId := c.Param("roomId")
		userId := c.Param("userId")
		websocketServer.serveWs(c.Writer, c.Request, roomId, userId)
	 })


	router.Run()


}
