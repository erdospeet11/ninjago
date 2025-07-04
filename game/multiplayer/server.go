// This file contains the server implementation of the multiplayer game.
package multiplayer

import (
	"image/color"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"szofttech.inf.elte.hu/szofttech-c-2024/group-08/ninjago/userinfo"
)

// ProtoGameInfo represents the game state information sent to the clients.
type User struct {
	userinfo.UserInfo              // The user information for the current player.
	IP                string       // The IP address of the client.
	Player            *ProtoPlayer // The player information for the current player.
}

// GameServer represents the server for the multiplayer game.
type GameServer struct {
	GameInfo ProtoGameInfo             // The game state information sent to the clients.
	Colors   []color.RGBA              // The colors available for the players.
	server   *gin.Engine               // The server instance.
	clients  map[*websocket.Conn]*User // The clients connected to the server.
}

// NewGameServer creates a new instance of GameServer and initializes the server settings and routes.
//
// Returns:
//   - *GameServer: A pointer to the newly created GameServer instance.
func NewGameServer() *GameServer {
	server := GameServer{
		server: gin.Default(),
		GameInfo: ProtoGameInfo{
			GameState:     GameStateLobby,
			Players:       make([]ProtoPlayer, 1, 10),
			Monsters:      make([]ProtoEntity, 0, 10),
			Bombs:         make([]ProtoEntity, 0, 10),
			Explosions:    make([]ProtoEntity, 0, 30),
			Boxes:         make([]ProtoEntity, 0, 10),
			StatusEffects: make([]ProtoEntity, 0, 10),
		},
		// Colors for the players Red, Green, Blue, Yellow, Cyan, Magenta, White
		Colors: []color.RGBA{
			{R: 255, G: 0, B: 0, A: 255},
			{R: 0, G: 255, B: 0, A: 255},
			{R: 0, G: 0, B: 255, A: 255},
			{R: 255, G: 255, B: 0, A: 255},
			{R: 0, G: 255, B: 255, A: 255},
			{R: 255, G: 0, B: 255, A: 255},
			{R: 255, G: 255, B: 255, A: 255}},
		clients: make(map[*websocket.Conn]*User),
	}

	server.server.GET("/", server.websocketHandler())

	return &server
}

// Run starts the game server and listens on port 8080 for incoming connections.
//
// Returns:
//   - func(): A function to close the server and all active client connections.
func (s *GameServer) Run() (Close func()) {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: s.server,
	}

	go func() { log.Println(srv.ListenAndServe()) }()
	log.Println("Server is listening")

	return func() {
		srv.Close()

		for conn := range s.clients {
			conn.Close()
		}
	}
}

// websocketHandler handles the WebSocket connections from clients.
//
// Returns:
//   - gin.HandlerFunc: A Gin handler function to manage WebSocket connections.
func (s *GameServer) websocketHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		defer func() {
			if client, exists := s.clients[conn]; exists && client != nil && client.Player != nil {
				client.Player.IsDead = true
				s.Colors = append(s.Colors, client.Player.Color)
			}

			conn.Close()
			delete(s.clients, conn)
		}()

		if err != nil {
			log.Println("Failed to upgrade connection", err)
			c.Writer.WriteHeader(http.StatusInternalServerError)

			return
		}

		s.clients[conn] = &User{UserInfo: userinfo.UserInfo{UserID: uuid.New().String()}}

		err = conn.WriteJSON(s.clients[conn].UserInfo)
		if err != nil {
			log.Println("Failed to send UUID", err)
			c.Writer.WriteHeader(http.StatusInternalServerError)

			return
		}

		s.clients[conn].IP = c.Request.RemoteAddr
		s.GameInfo.Players = append(s.GameInfo.Players, ProtoPlayer{
			X: float64(rand.Intn(14))*16 + 20,
			Y: float64(rand.Intn(14))*16 + 20,
		})
		s.clients[conn].Player = &s.GameInfo.Players[len(s.GameInfo.Players)-1]

		log.Println(s.clients[conn].Player, &s.GameInfo.Players[len(s.GameInfo.Players)-1])
		randColor := rand.Intn(len(s.Colors))
		s.clients[conn].Player.Color = s.Colors[randColor]
		s.Colors[randColor] = s.Colors[len(s.Colors)-1]
		s.Colors = s.Colors[:len(s.Colors)-1]

		var receivedMessage ProtoPlayer

		// Add panic recovery to handle websocket library panics
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered from panic in websocket handler: %v", r)
			}
		}()

		for {
			comStart := time.Now()

			// Use a separate function to handle the read operation with panic recovery
			err = func() (readErr error) {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("Recovered from websocket read panic: %v", r)
						readErr = &websocket.CloseError{Code: websocket.CloseAbnormalClosure, Text: "Connection failed"}
					}
				}()
				return conn.ReadJSON(&receivedMessage)
			}()

			if err != nil {
				log.Println("Failed to read from connection", err.Error())
				// If connection is closed or failed, exit the loop
				if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) ||
					websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Println("Connection closed, stopping server handler")
					return
				}
				// For other errors, wait a bit and try again
				time.Sleep(time.Millisecond * 100)
				continue
			}

			s.clients[conn].Username = receivedMessage.Username
			s.clients[conn].Player.Username = receivedMessage.Username
			s.clients[conn].Player.Control = receivedMessage.Control

			// Use a separate function to handle the write operation with panic recovery
			err = func() (writeErr error) {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("Recovered from websocket write panic: %v", r)
						writeErr = &websocket.CloseError{Code: websocket.CloseAbnormalClosure, Text: "Connection failed"}
					}
				}()
				return conn.WriteJSON(s.GameInfo)
			}()

			if err != nil {
				log.Println("Failed to write message to client", err)
				// If we can't write, the connection is likely closed
				if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) ||
					websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Println("Connection closed, stopping server handler")
					return
				}
			}

			s.clients[conn].Player.Ping = time.Since(comStart)

			// Add a small delay to prevent overwhelming the connection
			time.Sleep(time.Millisecond * 16) // ~60 FPS
		}
	}
}
