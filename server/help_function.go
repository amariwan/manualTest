package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"net/url"
	"os"
	"strings"

	//"net/smtp"
	"runtime"
	"time"

	"github.com/apsdehal/go-logger"
	c "github.com/ostafen/clover/v2"
)

var dateOfToScanFile time.Time

// Parse the command line arguments and return the port number, whether to use TLS, whether to use a
// proxy, and the proxy port number.
func parseCommandLineArguments(logLevel *logger.Logger) (int, bool, bool, int) {
	var level int
	var isToBeRestored bool
	var isLinux bool
	var port int
	if runtime.GOOS == "Linux" {
		isLinux = true
	} else {
		isLinux = false
	}

	flag.IntVar(&level, "l", 2, "LogLevel")
	flag.IntVar(&port, "p", 3000, "Port")
	flag.BoolVar(&isToBeRestored, "r", false, "if the database should be restored (default false)")
	flag.Parse()

	logLevel.DebugF("level is---- %v\n", level)

	return level, isToBeRestored, isLinux, port
}

// It converts an integer to a log level
func int2LogLevel(level int) logger.LogLevel {
	switch level {
	case 1:
		return logger.CriticalLevel
	case 2:
		return logger.ErrorLevel
	case 3:
		return logger.WarningLevel
	case 4:
		return logger.NoticeLevel
	case 5:
		return logger.InfoLevel
	case 6:
		return logger.DebugLevel
	default:
		return logger.ErrorLevel
	}
}

// > This function initializes the database
func databaseInitialization(db *c.DB, isToBeRestored bool, isLinux bool, logLevel *logger.Logger) error {
	var err error
	err = createTableCollection(db, isToBeRestored, logLevel) // Create Database Tables
	if err != nil {
		logLevel.Error(err.Error())
	}
	return err
}

//func readFromFileAndInsertInDataBase(debugLevel *logger.Logger, db *c.DB) error {
//	path := "test.json"
//	content, err := ioutil.ReadFile(path)
//	if err != nil {
//		debugLevel.Fatal("Error when opening file: " + err.Error())
//	}
//
//	// Now let's unmarshall the data into `payload`
//	payload := &OrderAnswerJenkins{}
//	err = json.Unmarshal(content, payload)
//	if err != nil {
//		debugLevel.Fatal("Error during Unmarshal(): " + err.Error())
//	}
//
//	insertIn2DatabaseTable(db, debugLevel, payload)
//	return nil
//}

// This function inserts the data from the configData struct into the database table
func insertIn2DatabaseTable(db *c.DB, logLevel *logger.Logger, configData OrderAnswerJenkins) OrderAnswerJenkins {
	logLevel.DebugF("Inserting in the dataBase ...")
	doc := c.NewDocumentOf(configData)
	mapDocId, _ := db.InsertOne("assignment", doc)
	updates := make(map[string]interface{})
	updates["Id"] = mapDocId
	db.Update(c.NewQuery("assignment").Where(c.Field("Id").Eq("")), updates)

	configData.Id = mapDocId
	configData.CreatedAt = time.Now().String()
	configData.UpdatedAt = time.Now().String()
	fmt.Println(mapDocId)
	for i := 0; i < len(configData.Testers); i++ {
		from := "dev@aland-mariwan.de"
		password := "njxqreyjxevcsvpd"
		auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")
		url, err := url.Parse("http://localhost:3000/cfg/?id=1&tester=''")
		if err != nil {
			logLevel.Error(err.Error())
		}
		url.Scheme = "http"
		url.Host = "localhost:3000"
		q := url.Query()
		q.Set("id", configData.Id)
		q.Set("tester", configData.Testers[i].Email)
		url.RawQuery = q.Encode()
		fmt.Println(url)

		//// Here we do it all: connect to our server, set up a message and send it

		to := []string{configData.Testers[i].Email}
		msg := []byte("To: " + configData.Testers[i].Email + "\r\n" +
			configData.Mail.Subject + "\r\n" +
			"\r\n" +
			configData.Mail.Body + "\r\n" +
			"go to the Link\r\n" + fmt.Sprintf("%v", url))
		err = smtp.SendMail("smtp.gmail.com:587", auth, from, to, msg)
		if err != nil {
			logLevel.Error(err.Error())
		}
	}
	return configData
}

// It creates a table collection.
func createTableCollection(db *c.DB, isToBeRestored bool, logLevel *logger.Logger) error {
	collection, err := db.HasCollection("assignment")
	if err != nil {
		logLevel.Error(err.Error())
	}
	if isToBeRestored {
		if collection {
			db.DropCollection("assignment")
		}
	}

	logLevel.DebugF("Test table created")
	// Create a new collection named "assignment"
	err = db.CreateCollection("assignment")
	return err
}

func insertOrderIn2Database(req *http.Request, setLogLevel *logger.Logger, database *c.DB) []byte {
	orderRequestJenkins := OrderAnswerJenkins{}
	err := json.NewDecoder(req.Body).Decode(&orderRequestJenkins)
	if err != nil {
		setLogLevel.Error(err.Error())
	}
	orderRequestJenkins.CreatedAt = time.Now().String()
	answer := insertIn2DatabaseTable(database, setLogLevel, orderRequestJenkins)
	answerJson, errJson := json.Marshal(answer)
	if errJson != nil {
		setLogLevel.ErrorF(errJson.Error())
	}
	return answerJson
}

func createOrder(req *http.Request, setLogLevel *logger.Logger, w http.ResponseWriter, database *c.DB) OrderAnswerFlutter {
	orderRequestJenkins := OrderAnswerJenkins{}
	err := json.NewDecoder(req.Body).Decode(&orderRequestJenkins)
	if err != nil {
		setLogLevel.Error(err.Error())
	}
	post("http://localhost:3000/jenkins", orderRequestJenkins.Raw())
	w.Write(orderRequestJenkins.Raw())
	docs, _ := database.FindAll(c.NewQuery("assignment"))
	answer := OrderAnswerFlutter{Name: "getValue", Orders: []Order{}}
	for _, doc := range docs {
		order := &OrderDB{}
		doc.Unmarshal(order)
		e := Order{}
		answer.Orders = append(answer.Orders, *getEFromDB(&e, order))

	}
	return answer
}

func updateValue(e *Order, err error, database *c.DB, setLogLevel *logger.Logger, index int) (error, *OrderAnswerFlutter, int) {
	update := make(map[string]interface{})

	update["Id"] = e.Id
	update["UpdatedAt"] = time.Now().String()
	update["CallbackUrl"] = e.CallbackUrl
	update["Creator"] = e.Creator
	update["Description"] = e.Description
	update["ExpireDate"] = e.ExpireDate
	update["Mail"] = e.Mail
	update["Tasks"] = e.Tasks
	update["Testers"] = e.Testers
	update["CreatedAt"] = e.CreatedAt
	update["TestComplete"] = e.TestComplete
	update["UpdatedAt"] = e.UpdatedAt
	err = database.Update(c.NewQuery("assignment").Where(c.Field("Id").Eq(e.Id)), update)
	if err != nil {
		setLogLevel.Error(err.Error())
	}
	answer := &OrderAnswerFlutter{Name: e.Name, Orders: []Order{*e}, Stump: Signature{Index: index, CreatedAt: time.Now().String()}}
	index++
	return err, answer, index
}

func sendEmail(e *Order, err error, setLogLevel *logger.Logger) {
	from := "dev@aland-mariwan.de"
	password := "njxqreyjxevcsvpd"
	auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")

	to := []string{e.Testers[0].Email}
	msg := []byte("To: " + e.Testers[0].Email + "\r\n" +
		"\r\n Thes Test id:\r\n" + fmt.Sprintf("%v", e.Id) +
		" mit Beschreibung \r\n" + fmt.Sprintf("%v", e.Description) + "\r\n" +
		"ist nicht mehr zu Ihnen gemerkt\r\n")
	err = smtp.SendMail("smtp.gmail.com:587", auth, from, to, msg)
	if err != nil {
		setLogLevel.Error(err.Error())
	}
}

func sendProtocol(e *Order, database *c.DB, setLogLevel *logger.Logger, index int) *OrderAnswerFlutter {
	update := OrderAnswerJenkins{}
	info, err := e.Raw()
	callbackUrl := "http://localhost:3000/testCallback"

	request, error := http.NewRequest("POST", callbackUrl, bytes.NewBuffer(info))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()
	ioutil.ReadAll(response.Body)

	json.Unmarshal(info, &update)
	f, err := os.Create("test.json")
	if err != nil {
		panic(error)
	}

	defer f.Close()
	f.Write([]byte(update.Raw()))

	update2 := make(map[string]interface{})
	update2["Id"] = e.Id
	update2["Update2dAt"] = time.Now().String()
	update2["CallbackUrl"] = e.CallbackUrl
	update2["Creator"] = e.Creator
	update2["Description"] = e.Description
	update2["ExpireDate"] = e.ExpireDate
	update2["Mail"] = e.Mail
	update2["Tasks"] = e.Tasks
	update2["Testers"] = e.Testers
	update2["CreatedAt"] = e.CreatedAt
	update2["TestComplete"] = e.TestComplete
	update2["UpdatedAt"] = e.UpdatedAt
	update2["Id"] = e.Id
	update2["UpdatedAt"] = time.Now().String()
	update2["CallbackUrl"] = e.CallbackUrl
	update2["Creator"] = e.Creator
	update2["Description"] = e.Description
	update2["ExpireDate"] = e.ExpireDate
	update2["Mail"] = e.Mail
	update2["Tasks"] = e.Tasks
	update2["Testers"] = e.Testers
	update2["CreatedAt"] = e.CreatedAt
	update2["TestComplete"] = e.TestComplete
	update2["UpdatedAt"] = e.UpdatedAt
	update2["Id"] = e.Id
	update2["UpdatedAt"] = time.Now().String()
	update2["CallbackUrl"] = e.CallbackUrl
	update2["Creator"] = e.Creator
	update2["Description"] = e.Description
	update2["ExpireDate"] = e.ExpireDate
	update2["Mail"] = e.Mail
	update2["Tasks"] = e.Tasks
	update2["Testers"] = e.Testers
	update2["CreatedAt"] = e.CreatedAt
	update2["TestComplete"] = e.TestComplete
	update2["UpdatedAt"] = e.UpdatedAt
	err = database.Update(c.NewQuery("assignment").Where(c.Field("Id").Eq(e.Id)), update2)
	if err != nil {
		setLogLevel.Error(err.Error())
	}
	answer := &OrderAnswerFlutter{Name: e.Name, Orders: []Order{*e}, Stump: Signature{Index: index, CreatedAt: time.Now().String()}}
	index++
	return answer
}

func getAllOrder(database *c.DB, e *Order, index int) OrderAnswerFlutter {
	docs, _ := database.FindAll(c.NewQuery("assignment"))
	answer := OrderAnswerFlutter{Name: e.Name, Orders: []Order{}, Stump: Signature{Index: index, CreatedAt: time.Now().String()}}
	for _, doc := range docs {
		order := &OrderDB{}
		err := doc.Unmarshal(order)
		if err != nil {
			panic(err)
		}
		answer.Orders = append(answer.Orders, *getEFromDB(e, order))
		fmt.Println(order)
	}
	return answer
}

func createResponse(database *c.DB, e *Order) {
}

func join2Array(array1 []int, array2 []int) []int {
	for i := 0; i < len(array1); i++ {
		array2 = append(array2, array1[i])
	}
	return array2
}

func getEFromDB(e *Order, order2 *OrderDB) *Order {
	e.TestComplete = order2.TestComplete
	e.CallbackUrl = order2.CallbackUrl
	e.CreatedAt = order2.CreatedAt
	e.Creator = cloneTester(order2.Creator)
	e.Description = order2.Description
	e.ExpireDate = order2.ExpireDate
	e.Id = order2.Id
	e.Mail = cloneMail(order2.Mail)
	e.Tasks = cloneTasks(order2.Tasks)
	e.Testers = cloneTesters(order2.Testers)
	e.UpdatedAt = order2.UpdatedAt
	return e
}

func LoadAndAddToRoot(funcMap template.FuncMap, rootTemplate *template.Template) error {
	// This solution is CC by share alike 4.0
	// Copyright Rik-777
	// https://stackoverflow.com/a/50581032
	err := fs.WalkDir(flutterCfgFs, ".", func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if !d.IsDir() && strings.HasSuffix(path, ".html") {
			data, readErr := flutterCfgFs.ReadFile(path)
			if readErr != nil {
				return readErr
			}
			t := rootTemplate.New(path).Funcs(funcMap)
			if _, parseErr := t.Parse(string(data)); parseErr != nil {
				return parseErr
			}
		}
		return nil
	})
	return err
}

func post(httpPostUrl string, toSendMessage []byte) []byte {
	request, error := http.NewRequest("POST", httpPostUrl, bytes.NewBuffer(toSendMessage))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println("response Body:", string(body))
	return body
}
