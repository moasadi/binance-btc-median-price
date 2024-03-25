// main.go
package main

import (
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/moasadi/binance-trade/api/application"
	"github.com/moasadi/binance-trade/api/infrastructure"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan float64)
var mutex = &sync.Mutex{}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	wsURL := os.Getenv("BINANCE_WS_URL")

	if wsURL == "" {
		log.Fatal("BINANCE_WS_URL must be set")
	}

	tradeService := createTradeService(wsURL)
	tradeApp := createTradeApp(tradeService)

	medianPriceChan := make(chan float64)
	go func() {
		for {
			medianPrice := <-medianPriceChan
			broadcast <- medianPrice
		}
	}()

	go func() {
		err := tradeApp.Run(medianPriceChan)
		if err != nil {
			log.Fatal(err)
		}
	}()

	r := gin.Default()
	r.LoadHTMLFiles("templates/index.html")
	r.GET("/ws", func(c *gin.Context) {

		handleConnections(c.Writer, c.Request)
	})
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	go handleMessages()
	r.Run()
}
func createTradeService(wsURL string) *infrastructure.TradeService {
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	return infrastructure.NewTradeService(conn)
}

func createTradeApp(service *infrastructure.TradeService) *application.TradeApp {
	return application.NewTradeApp(service)
}
func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()
	mutex.Lock()
	clients[ws] = true
	mutex.Unlock()
	for {
		var msg float64
		err := ws.ReadJSON(&msg)
		if err != nil {
			mutex.Lock()
			delete(clients, ws)
			mutex.Unlock()
			break
		}
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		mutex.Lock()
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
		mutex.Unlock()
	}
}
