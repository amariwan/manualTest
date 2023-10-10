package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/apsdehal/go-logger"
	"github.com/gorilla/websocket"
	c "github.com/ostafen/clover/v2"
)

func example(w http.ResponseWriter, r *http.Request) {
	setLogLevel, errorWithLogin := logger.New("test", 1, os.Stdout)
	if errorWithLogin != nil {
		panic(errorWithLogin) // Check for error
	}
	db, _ := c.Open("test-db")
	collection, err := db.HasCollection("assignment")
	if err != nil {
		panic(err)
	}

	if collection {
		db.DropCollection("assignment")
	}
	err = db.CreateCollection("assignment")
	if err != nil {
		panic(err)
	}
	var event Order
	event.Name = "getValues"
	event.CallbackUrl = "callbackUrl"
	var creator Tester
	creator.Email = "test@gmail.com"
	creator.Name = "Ersteller"
	event.Creator = creator
	event.Description = "Test Beschreibung"
	event.ExpireDate = "2025-12-13T13:58:12.302Z"
	var mail MailStruck
	mail.Body = "Bitte teste die neue Firmware"
	mail.Subject = "Es gibt einen neuen Test"
	event.Mail = mail
	var tasks Task
	tasks.Description = "first Task"
	event.Tasks = []Task{tasks, tasks}
	event.Testers = []Tester{creator}
	var sig Signature
	sig.Index = 0
	sig.CreatedAt = "2023-01-18T14:34:21.257Z"
	for i := 0; i < 10; i++ {
		var orderRequestJenkins OrderAnswerJenkins
		request, err := event.Raw()
		if err != nil {
			panic(err)
		}
		orderRequest, err := orderRequestJenkins.OrderAnswerJenkinsFromRaw(request)
		if err != nil {
			panic(err)
		}
		insertIn2DatabaseTable(db, setLogLevel, *orderRequest)
	}
	debugLevel := int2LogLevel(6)
	setLogLevel.SetLogLevel(debugLevel)
	ws, err := NewWebSocket(w, r, setLogLevel, db)
	if err != nil {
		panic(err)
	}

	ws.On("getValues", func(event *Order) {
		answer := getAllOrder(db, event, 0)
		ws.Out <- answer.Raw()
	})

}
func TestExample(t *testing.T) {
	// Create test server with the example handler.
	s := httptest.NewServer(http.HandlerFunc(example))
	defer s.Close()
	// Convert http://127.0.0.1 to ws://127.0.0.
	u := "ws" + strings.TrimPrefix(s.URL, "http")

	// Connect to the server
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ws.Close()

	// Send message to server, read response and check to see if it's what we expect.
	for i := 0; i < 10; i++ {

		if err := ws.WriteMessage(websocket.TextMessage, []byte(`{
	    "event": "getValues"
	}`)); err != nil {
			t.Fatalf("%v", err)
		}
		_, p, err := ws.ReadMessage()
		if err != nil {
			t.Fatalf("%v", err)
		}
		var event OrderAnswerFlutter
		err = json.Unmarshal([]byte(p), &event)
		if err != nil {
			t.Fatalf("%v", err)
			return
		}
		if len(event.Orders) != 10 {
			t.Fatalf("%v", err)
			return
		}
	}
}
