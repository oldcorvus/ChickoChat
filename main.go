package main

import (
	"chicko_chat/controllers"
	"chicko_chat/database"
	"chicko_chat/models"
	"chicko_chat/websocket"

	"flag"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	mongoURI := flag.String("mongoURI", "mongodb://mongo:27017", "Database hostname url")
	enableCredentials := flag.Bool("enableCredentials", false, "Enable the use of credentials for mongo connection")
	flag.Parse()

	db := database.ConnectDatabse(*mongoURI, *enableCredentials)
	router := gin.Default()
	controller := controllers.Controller{
		DB: db,
	}
	manager := &websocket.BrokerManager{
		Brokers: make(map[*data.Broker]bool),
		DB:      db,
	}
	websocketServer := &websocket.WsServer{
		Manager: manager,
	}
	router.LoadHTMLGlob("templates/*")
	router.GET("/", func(c *gin.Context) {

		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Sample Front",
			"name":  "Moel",
		})
	})
	router.GET("/chat/", controller.JoinRoom)
	router.POST("/start/", controller.StartConversationApi)
	router.POST("/user-rooms/", controller.GetUserRoomsApi)
	router.POST("/create-room/", controller.CreateRoomApi)
	router.POST("/room-history/", controller.RoomHistoryApi)
	router.POST("/add-user-room/", controller.AddUserToRoomApi)
	router.POST("/room-user-details/", controller.GetUserDetailsRoomApi)
	router.GET("/ws/:roomId/:userId/", func(c *gin.Context) {
		roomId := c.Param("roomId")
		userId := c.Param("userId")
		websocketServer.ServeWs(c.Writer, c.Request, roomId, userId)
	})
	router.Run()

}
