package horloge

import (
	"github.com/gin-gonic/gin"
	melody "gopkg.in/olahol/melody.v1"
)

type Server struct {
	websocketHandler *melody.Melody
	httpServer       *gin.Engine
}

func NewWebsocketServer() *Server {
	r := gin.Default()
	m := melody.New()

	r.GET("/", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	r.GET("/ping", func(c *gin.Context) {
		c.Data(200, "text/plain", []byte("pong"))
	})

	return &Server{
		websocketHandler: m,
		httpServer:       r,
	}
}

func (s *Server) Publish(message string) {
	s.websocketHandler.Broadcast([]byte(message))
}

func (s *Server) Run(addr string) {
	s.httpServer.Run(addr)
}
