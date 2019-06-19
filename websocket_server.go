package horloge

import (
	"github.com/gin-gonic/gin"
	melody "gopkg.in/olahol/melody.v1"
)

type WebsocketServer struct {
	websocketHandler *melody.Melody
	httpServer       *gin.Engine
}

// NewWebsocketServer Creates a new WebsocketServer instance
func NewWebsocketServer() *WebsocketServer {
	r := gin.Default()
	m := melody.New()

	r.GET("/", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	r.GET("/ping", func(c *gin.Context) {
		c.Data(200, "text/plain", []byte("pong"))
	})

	return &WebsocketServer{
		websocketHandler: m,
		httpServer:       r,
	}
}

// Publish Emits a message on the websocket connection
func (s *WebsocketServer) Publish(message string) {
	s.websocketHandler.Broadcast([]byte(message))
}

// Run Runs the server
func (s *WebsocketServer) Run(addr string) {
	s.httpServer.Run(addr)
}
