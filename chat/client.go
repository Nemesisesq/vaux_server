package chat

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/websocket"
	"github.com/nemesisesq/vaux_server/models"
	log "github.com/sirupsen/logrus"
	"github.com/soveran/redisurl"
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
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Data struct {
	Type     string      `json:"type"`
	Paylaod  interface{} `json:"payload"`
	ThreadID interface{} `json:"thread_id"`
}

// Client is a middleman between the websocket connection and the hub (redis) connection.
type Client struct {
	user models.User

	subRedisConn redis.Conn

	pubRedisConn redis.Conn

	// The websocket connection.
	ws *websocket.Conn

	currentThread models.Thread

	// Buffered channel of outbound messages.
	out chan []byte
	in  chan []byte

	shutDown chan bool
}

func NewClient() Client {
	c := Client{}
	//c.user =  user
	c.out = make(chan []byte)
	c.in = make(chan []byte)
	c.shutDown = make(chan bool, 1)

	subRedisConn, err := redisurl.Connect()
	if err != nil {
		panic(err)
	}
	pubRedisConn, err := redisurl.Connect()

	if err != nil {
		panic(err)
	}

	c.subRedisConn = subRedisConn
	c.pubRedisConn = pubRedisConn

	return c

}

func (c *Client) Subscribe() {

	sub := &redis.PubSubConn{Conn: c.subRedisConn}

	threads, err := GetAllThreads(c)

	if err != nil {
		log.Panic(err)
	}

	for _, v := range threads {
		s := fmt.Sprintf("thread.%v", v.ID.String())
		sub.Subscribe(s)
	}

SUB:
	for {

		switch v := sub.Receive().(type) {
		case redis.Message:
			c.out <- v.Data
		case redis.Subscription:
			break
		case error:
			break SUB
			return
		default:
			log.Info("no data")
		}
	}
	return
}

func (c *Client) Publish() {

	pub := &redis.PubSubConn{Conn: c.pubRedisConn}
PUB:
	for {
		select {
		case data := <-c.in:

			d := Data{}

			err := json.Unmarshal(data, &d)

			if err != nil {
				log.Panic(err)
			}

			tmp, err := json.Marshal(d.Paylaod.([]interface{})[0])

			if err != nil {
				log.Panic(err)
			}
			message := &models.Message{}
			err = json.Unmarshal(tmp, &message)

			if err != nil {
				log.Panic(err)
			}

			channel := fmt.Sprintf("thread.%v", d.ThreadID)

			thread := models.GetThread(d.ThreadID.(string))

			message.User = c.user
			message.UserID = message.User.ID
			message.ThreadID = thread.ID
			message.Thread = thread

			err = message.Create()

			if err != nil {
				log.Panic(err)
			}

			j, err := json.Marshal(message)

			outString := string(j)
			//thread.AddMessage(*message)

			//Publish the message

			if err != nil {
				panic(err)
			}
			d.Paylaod = []interface{}{outString}

			data, err = json.Marshal(d)

			_, err = pub.Conn.Do("PUBLISH", channel, data)
			if err != nil {
				log.Panic(err)
			}

		case <-c.shutDown:
			break PUB
			return
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
	return
}

func (c *Client) Unsubscribe() {

	//c.pubSubConn.Unsubscribe()

}

func (c *Client) addUser() {
	//	Set The user name

	conn, err := redisurl.Connect()

	defer conn.Close()

	if err != nil {
		panic(err)
	}
	userKey := fmt.Sprint("online.", c.user.ID)
	val, err := conn.Do("SET", userKey, c.user.ID, "NX", "EX", "120")
	if val == nil {
		fmt.Println("User already online")
		os.Exit(1)
	}

	if err != nil {
		panic(err)
	}

	val, err = conn.Do("SADD", "users", c.user.ID)

	if err != nil {
		panic(err)
	}
	log.Info("Lets return some stuff")
	return
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.Unsubscribe()
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		messageType, message, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
			return
		}

		if messageType == 456 {

		}

		data := Data{}

		err = json.Unmarshal(message, &data)
		if err != nil {
			panic(err)
		}

		c.processData(data)

		//message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		//c.in <- message
	}
	return
}

// writePump pumps messages from the hub to the websocket wsection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.out:
			c.ws.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.ws.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.out)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.out)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.shutDown <- true
				return
			}
		}
	}
	return
}

func Connect(c buffalo.Context) error {
	serveWs(c)
	return nil
}

// serveWs handles websocket requests from the peer.
func serveWs(c buffalo.Context) {
	w := c.Response()
	r := c.Request()
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := NewClient()
	client.ws = conn

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.Subscribe()
	go client.Publish()
	go client.writePump()
	go client.readPump()

	//go client.broadcastThreads()

}
