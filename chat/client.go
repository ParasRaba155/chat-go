package chat

import (
	"bytes"
	"context"
	"net/http"
	"time"

	http_util "app/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// client is holds the websocket connection and contains details about the user
type client struct {
	room *room

	conn *websocket.Conn

	send chan []byte

	// email address of the sender
	senderMail string

	// email address of the reciver
	receiverMail string

	redis *redis.Client
}

func (c *client) readPump(l *zap.Logger) {
	defer func() {
		c.room.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				l.Error("reading socket error", zap.Error(err))
			}
			break
		}
		l.Debug("subsribe step 1")
		subscriber := c.redis.Subscribe(context.TODO(), c.room.id.String())
		l.Debug("subsribe step 2")
		redisMsg, err := subscriber.Receive(context.TODO())
		l.Debug("subscirbe step 3")
		if err != nil {
			l.Error("redis msg recive error", zap.Error(err))
		}
		l.Debug("redis msg", zap.Any("redis_msg", redisMsg))
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.room.broadcast <- message
	}
}

// writePump pumps message from the hub to the web socket connection
//
// The application runs writePump in per-connection goroutine. the application
// ensures that there is at most one writer on a connection by executing all
// writes from this goroutine
func (c *client) writePump(l *zap.Logger) {
	ticker := time.NewTicker(pingPeriod)
	l.Debug("started Writing")

	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			l.Debug("hi from select")
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			defer w.Close()

			w.Write(message)

			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte("\n"))
				msg := <-c.send
				w.Write(msg)
				c.redis.Publish(context.TODO(), c.room.id.String(), msg)
				l.Debug("published the message")
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			l.Debug("Ping")
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func Handler(l *zap.Logger, rc *redis.Client, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		l.Error("upgrading to websocket failed", zap.Error(err))
		w.Write([]byte("could not process the websocket request"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	roomID, err := uuid.NewRandom()
	if err != nil {
		l.Error("error generating room id", zap.Error(err))
		w.Write([]byte(""))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	room := NewRoom(roomID)
	go room.Run(l)

	receiver, ok := mux.Vars(r)["receiver"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no receiver"))
		return
	}

	user := http_util.GetUserFromRequestContext(r)

	client := &client{
		room:        room,
		conn:        conn,
		send:        make(chan []byte, 256),
		senderMail:  user.Email,
		receiverMail: receiver,
		redis:       rc,
	}
	l.Debug("registering client in room")

	client.room.register <- client
	l.Debug("registered client in room")

	go client.writePump(l)
	go client.readPump(l)
}
