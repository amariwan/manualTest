package main

import (
	"log"
	"net/http"

	"github.com/apsdehal/go-logger"
	"github.com/gorilla/websocket"
	c "github.com/ostafen/clover/v2"
)

// A slice of pointers to websocket connections.
var clients []*websocket.Conn

// Defining a constant.
const (
	MAX_MSG_SIZE = 100000
)

// Creating a new upgrader object.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Creating a new struct called WebSocket.
type WebSocket struct {
	Conn    *websocket.Conn
	Out     chan []byte
	In      chan []byte
	Events  map[string]EventHandler
	IsLogin bool
}

// It creates a new websocket connection.
func NewWebSocket(w http.ResponseWriter, r *http.Request, debugLevel *logger.Logger, database *c.DB) (*WebSocket, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	clients = append(clients, conn)
	if err != nil {
		log.Printf("[ERROR | SOCKET CONNECT] %v", err)
		return nil, err
	}
	// conn.SetWriteDeadline(time.Now().Add(MSG_TIMEOUT))
	ws := &WebSocket{
		Conn:    conn,
		Out:     make(chan []byte),
		In:      make(chan []byte),
		Events:  make(map[string]EventHandler),
		IsLogin: false,
	}
	go ws.Reader(debugLevel, database)
	go ws.Writer(debugLevel)
	return ws, nil
}

// A method of the WebSocket struct. It takes a logger and a database as parameters.
func (ws *WebSocket) Reader(debugLevel *logger.Logger, database *c.DB) {
	defer func() {
		disconnectClient(ws)
	}()

	for {
		_, message, err := ws.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[ERROR] %v", err)
			}
			break
		}
		orderFromFlutter := &Order{}
		event, err := orderFromFlutter.NewEventFromRaw(message)
		if err != nil {
			log.Printf("[ERROR | MSG] %v", err)
		} else {
			debugLevel.DebugF("[MSG] %v", event)
		}
		if event.Name == "updateValue" || event.Name == "protokol" {
			answerToFlutter := &Order{
				Name:         event.Name,
				CallbackUrl:  event.CallbackUrl,
				Creator:      cloneTester(event.Creator),
				Description:  event.Description,
				ExpireDate:   event.ExpireDate,
				Mail:         cloneMail(event.Mail),
				Tasks:        cloneTasks(event.Tasks),
				Testers:      cloneTesters(event.Testers),
				TestComplete: event.TestComplete,
				Id:           event.Id,
				CreatedAt:    event.CreatedAt,
				UpdatedAt:    event.UpdatedAt,
				Stump:        cloneSignature(event.Stump),
			}
			if action, ok := ws.Events[event.Name]; ok {
				action(answerToFlutter)
			}
		} else {
			docs, _ := database.FindAll(c.NewQuery("assignment").Where(c.Field("_id").Like(event.Id)))
			answerToFlutter := &Order{
				Name: event.Name,
			}
			if len(docs) > 0 {
				docs[len(docs)-1].Unmarshal(answerToFlutter)
			}
			answerToFlutter.Name = event.Name
			if action, ok := ws.Events[event.Name]; ok {
				action(answerToFlutter)
			}
		}

	}
}

// > This function is called when a client disconnects from the server
func disconnectClient(ws *WebSocket) {
	for i, client := range clients {
		if client == ws.Conn {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}

	ws.Conn.Close()

}

// A method of the WebSocket struct. It takes a logger as a parameter.
func (ws *WebSocket) Writer(debugLevel *logger.Logger) {
	for {
		select {
		case message, ok := <-ws.Out:
			if !ok {
				//log.Printf("hier ist nessage %v", ok)
				ws.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			writer, err := ws.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				//log.Printf("hier ist nessage in err %v", err)
				return
			}
			//log.Printf("hier ist nessage %v", message)
			writer.Write(message)
			writer.Close()
		}
	}
}

// It takes a byte array, and sends it to all the clients connected to the server
func broadcastMessage(message []byte) {
	for _, client := range clients {
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("[ERROR] %v", err)
		}
	}
}

// A method of the WebSocket struct. It takes an event name and an EventHandler function as parameters.
// It then adds the event name and the EventHandler function to the Events map of the WebSocket struct.
func (ws *WebSocket) On(event string, action EventHandler) *WebSocket {
	ws.Events[event] = action
	return ws
}
