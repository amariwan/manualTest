package main

import (
	//"bufio"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"time"

	"github.com/apsdehal/go-logger"
	"github.com/gorilla/handlers"
)

var (
	//go:embed flutterCfg
	flutterCfgFs embed.FS
	//pages        = map[string]string{
	//	"/cfg": "flutterCfg/index.html",
	//}
)

func main() {
	var index = 0
	var errorWithLogin error
	setLogLevel, errorWithLogin := logger.New("test", 1, os.Stdout)
	if errorWithLogin != nil {
		panic(errorWithLogin) // Check for error
	}
	setLogLevel.SetLogLevel(logger.ErrorLevel)
	level, isToBeRestored, isLinux, port := parseCommandLineArguments(setLogLevel)
	debugLevel := int2LogLevel(level)
	setLogLevel.SetLogLevel(debugLevel)
	mux := http.NewServeMux()
	client, err := fs.Sub(flutterCfgFs, "flutterCfg")
	if err != nil {
		panic(errorWithLogin) // Check for error
	}
	mux.HandleFunc("/cfg", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/cfg/index.html", http.StatusSeeOther)
	})
	mux.Handle("/cfg/", http.StripPrefix("/cfg/", handlers.CombinedLoggingHandler(os.Stdout,
		http.FileServer(http.FS(client)))))

	setLogLevel.DebugF("Creating db...")
	database := ConnectClover()
	defer database.Close()
	databaseInitialization(database, isToBeRestored, isLinux, setLogLevel)

	mux.Handle("/", http.FileServer(http.Dir("./assets")))
	mux.HandleFunc("/jenkins", func(w http.ResponseWriter, req *http.Request) {
		// `answerJson := insertOrderIn2Database(req, setLogLevel, database)` is a function call that inserts
		// an order received from a Jenkins server into a database table and returns a JSON response
		// containing information about the success or failure of the operation. The function takes three
		// arguments: `req`, which is an HTTP request object containing the order data; `setLogLevel`, which
		// is a logger object used for debugging and error reporting; and `database`, which is a database
		// object used to interact with the database. The function returns a JSON response containing
		// information about the success or failure of the operation.
		answerJson := insertOrderIn2Database(req, setLogLevel, database)
		w.Write(answerJson)
	})
	mux.HandleFunc("/create", func(w http.ResponseWriter, req *http.Request) {
		// The above code is calling the function `createOrder` with four arguments: `req` (presumably a
		// request object), `setLogLevel` (a function for setting the log level), `w` (presumably a writer
		// object), and `database` (presumably a database connection object). The return value of
		// `createOrder` is being assigned to the variable `answer`.
		answer := createOrder(req, setLogLevel, w, database)
		broadcastMessage((answer).Raw())
	})
	mux.HandleFunc("/testCallback", func(w http.ResponseWriter, req *http.Request) {
		orderRequestJenkins := OrderAnswerJenkins{}
		err := json.NewDecoder(req.Body).Decode(&orderRequestJenkins)
		if err != nil {
			setLogLevel.Error(err.Error())
		}
		w.Write(orderRequestJenkins.Raw())

	})
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws, err := NewWebSocket(w, r, setLogLevel, database)
		if err != nil {
			setLogLevel.Error(err.Error())
		}
		ws.On("updateValue", func(e *Order) {

			var answer *OrderAnswerFlutter
			err, answer, index = updateValue(e, err, database, setLogLevel, index)
			broadcastMessage((answer).Raw())
		})
		ws.On("getValues", func(e *Order) {
			answer := getAllOrder(database, e, index)
			index++
			ws.Out <- (answer).Raw()
		})
		ws.On("sendEmail", func(e *Order) {
			// `sendEmail` is a function that sends an email to a tester specified in the `Order` object `e`.
			// The email contains information about the test ID and description. The function uses the Gmail
			// SMTP server to send the email and requires authentication with a Gmail account.
			sendEmail(e, err, setLogLevel)
		})
		ws.On("getValue", func(e *Order) {

			answer := &OrderAnswerFlutter{Name: e.Name, Orders: []Order{*e}, Stump: Signature{Index: index, CreatedAt: time.Now().String()}}
			index++
			ws.Out <- (answer).Raw()
		})
		ws.On("delete", func(e *Order) {
			err = database.DeleteById("assignment", e.Id)
			if err != nil {
				setLogLevel.Error(err.Error())
			}
			answer := getAllOrder(database, e, index)
			// `broadcastMessage((answer).Raw())` is sending a message to all connected WebSocket clients,
			// containing the raw JSON data of the `answer` variable.
			broadcastMessage((answer).Raw())
		})
		ws.On("protokol", func(e *Order) {
			answer := sendProtocol(e, database, setLogLevel, index)
			broadcastMessage((answer).Raw())
		})
	})
	http.ListenAndServe(":"+fmt.Sprint(port), mux)
}
