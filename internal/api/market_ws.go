package api

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

var jwtSecret = []byte(getJWTSecretOrDefault("justdoit"))

func getJWTSecretOrDefault(def string) string {
	if s := os.Getenv("JWT_SECRET"); s != "" {
		return s
	}
	return def
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//client manager

type clienthub struct {
	clients map[*websocket.Conn]bool
	lock    sync.Mutex
}

func newHub() *clienthub {
	return &clienthub{
		clients: make(map[*websocket.Conn]bool),
	}
}

func (h *clienthub) add(c *websocket.Conn) {
	h.lock.Lock()
	h.clients[c] = true
	h.lock.Unlock()
}

func (h *clienthub) remove(c *websocket.Conn) {
	h.lock.Lock()
	delete(h.clients, c)
	h.lock.Unlock()
}

func (h *clienthub) broadcast(v interface{}) {
	h.lock.Lock()
	defer h.lock.Unlock()

	for c := range h.clients {
		c.SetWriteDeadline(time.Now().Add(5 * time.Second))
		if err := c.WriteJSON(v); err != nil {
			c.Close()
			delete(h.clients, c)
		}
	}
}

//Price message we send to clients

type PriceMessage struct {
	Price float64 `json:"price"`
	Time  int64   `json:"time"`
}

func SetupMarketWebS(r *gin.Engine, priceChan <-chan float64) {
	hub := newHub()

	r.GET("/ws/price", func(c *gin.Context) {

		tokenStr := c.Query("token")

		if tokenStr == "" {
			auth := c.GetHeader("Authorization")
			if strings.HasPrefix(auth, "Bearer ") {
				tokenStr = strings.TrimPrefix(auth, "Bearer ")
			}
		}

		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		//verify token

		_, err := verifyJWT(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token!"})
			return
		}

		//upgrade to websocket
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}
		hub.add(conn)

		//read pump
		go func() {
			defer func() {
				hub.remove(conn)
				conn.Close()
			}()
			conn.SetReadLimit(512)
			_ = conn.SetReadDeadline(time.Now().Add(60 * time.Second))
			conn.SetPongHandler(func(string) error {
				conn.SetReadDeadline(time.Now().Add(60 * time.Second))
				return nil
			})
			for {
				if _, _, err := conn.ReadMessage(); err != nil {
					return
				}
			}
		}()

	})
	//broadcaster go routine
	go func() {
		for p := range priceChan {
			msg := PriceMessage{Price: p, Time: time.Now().Unix()}
			hub.broadcast(msg)
		}
	}()

	//verify JWT token
}
func verifyJWT(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	return claims, nil
}
