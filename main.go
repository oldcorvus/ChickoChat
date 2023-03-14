package main

import (
	"chicko_chat/controllers"
	"chicko_chat/database"
	"chicko_chat/websocket"
	"chicko_chat/models"
	"net/http"
	"context"
	"flag"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)



func main() {
	var addr = flag.String("addr", ":8080", "The addr of the application.")
	mongoURI := flag.String("mongoURI", "mongodb://localhost:27017", "Database hostname url")
	mongoDatabase := flag.String("mongoDatabase", "chicko_chat", "Database name")
	enableCredentials := flag.Bool("enableCredentials", false, "Enable the use of credentials for mongo connection")
	flag.Parse()
	_ = addr
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
	manager := &websocket.BrokerManager{
		Brokers :      make(map[*data.Broker]bool),
		DB : db,
	}
	websocketServer := &websocket.WsServer{
		Manager : manager,
	}
	router.LoadHTMLGlob("templates/*")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Sample Front",
			"name" : "Moel",
		})
	})
	router.POST("/start/", controller.StartConversationApi) 
	router.POST("/user-rooms/", controller.GetUserRoomsApi) 
	router.POST("/create-room/", controller.CreateRoomApi) 
	router.POST("/room-history/", controller.RoomHistoryApi)
	router.POST("/add-user-room/", controller.AddUserToRoomApi)
	router.GET("/ws/:roomId/:userId/", func(c *gin.Context) {
		roomId := c.Param("roomId")
		userId := c.Param("userId")
		websocketServer.ServeWs(c.Writer, c.Request, roomId, userId)
	 })


	 router.Run()


}
