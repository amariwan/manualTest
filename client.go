package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

func runEvent(ws *websocket.Conn, st string) {
	if _, err := ws.Write([]byte(st)); err != nil {
		log.Fatal(err)
	}
	var msg = make([]byte, 50000)
	var n int
	var err error
	if n, err = ws.Read(msg); err != nil {
		log.Fatal(err)
	}
	str := string(msg)
	fmt.Println(str)
	fmt.Printf("Received: %s.\n", msg[:n])
}

type Tester struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type OrderRequestJenkins struct {
	CallbackUrl string      `json:"callbackUrl"`
	Creator     interface{} `json:"creator"`
	Description string      `json:"description"`
	ExpireDate  string      `json:"expireDate"`
	Mail        interface{} `json:"mail"`
	Tasks       interface{} `json:"tasks"`
	Testers     []Tester    `json:"testers"`
}

func main() {
	post("http://localhost:3000/create", []byte(`{
							    "callbackUrl": "http://localhost:3000/testCallback",
							    "creator": {
							        "email": "dev@aland-mariwan.de",
							        "name": "Aland Mariwan"
							    },
							    "description": "Do as I say or I will smash your face in a delicious pumpkin pie",
							    "expireDate": "2024-12-13T13:58:12.302Z",
							    "mail": {
							        "body": "Bitte teste die neue Firmware, bli bla blub",
							        "subject": "Es gibt einen neuen Test"
							    },
							    "tasks": [
							        {
							            "description": "Go jogging for 1 minute and 21 seconds."
							        },
							        {
							            "description": "Relax for 2 seconds"
							        },
							        {
							            "description": "Work until death. (But click \"completed\" just before you die)"
							        }
							    ],
							    "testers": [
							        {
							            "email": "dev@aland-mariwan.de",
							            "name": "Aland"
							        }
							    ]
							}`))

	//jenkins()
	//flutter()
}

func flutter() {
	origin := "http://localhost/"
	url := "ws://:3000/ws"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	str := `{
		"event":  "getValue",
		"callbackUrl": "callbackUrl",
		"creator": {
	   	 "email": "dev@aland-mariwan.de",
	   	 "name": "Aland Marwan"
	   	},
		"description": "description",
		"expireDate": "expireDate",
		"mail": {
	        "body": "Bitte teste die neue Firmware, bli bla blub",
	        "subject": "Es gibt einen neuen Test"
	      },
		"tasks": [
	         {
	           "description": "Go jogging for 1 minute and 21 seconds.",
	           "result": "",
	           "tester": ""
	         },
	         {
	           "description": "Relax for 2 seconds",
	           "result": "",
	           "tester": ""
	         },
	         {
	           "description": "Work until death. (But click \"completed\" just before you die)",
	           "result": "",
	           "tester": ""
	         }
	       ],
		"testers": [
	        {
	          "email": "noe@aland-marwan.de",
	          "name": "Johannes"
	        },
	        {
	          "email": "dev@aland-mariwan.de",
	          "name": "Gimpel 2"
	        }
	      ],
		"id": "c448447d-f01d-406f-91ee-cf8a6ae6fd3c",
		"createdAt": "createdAt"
		}`
	runEvent(ws, str)

	str3 := `{
		"event":  "updateValue",
		"callbackUrl": "callbackUrl",
		"creator": {
	   	 "email": "dev@aland-mariwan.de",
	   	 "name": "Aland Marwan"
	   	},
		"description": "description",
		"expireDate": "expireDate",
		"mail": {
	        "body": "Bitte teste die neue Firmware, bli bla blub",
	        "subject": "Es gibt einen neuen Test"
	      },
		"tasks": [
	         {
	           "description": "Go jogging for 1 minute and 21 seconds.",
	           "result": "",
	           "tester": ""
	         },
	         {
	           "description": "Relax for 2 seconds",
	           "result": "",
	           "tester": ""
	         },
	         {
	           "description": "Work until death. (But click \"completed\" just before you die)",
	           "result": "",
	           "tester": ""
	         }
	       ],
		"testers": [
	        {
	          "email": "niemad@aland-marwan.de",
	          "name": "Johannes"
	        },
	        {
	          "email": "dev@aland-mariwan.de",
	          "name": "Gimpel 2"
	        }
	      ],
		"id": "76f87a21-b5cd-4523-8c39-c37535a7db08",
		"createdAt": "createdAt"
		}`
	runEvent(ws, str3)

	str2 := `{
		"event":  "getValues",
		"callbackUrl": "callbackUrl",
		"creator": {
	   	 "email": "dev@aland-mariwan.de",
	   	 "name": "Aland Marwan"
	   	},
		"description": "description",
		"expireDate": "expireDate",
		"mail": {
	        "body": "Bitte teste die neue Firmware, bli bla blub",
	        "subject": "Es gibt einen neuen Test"
	      },
		"tasks": [
	         {
	           "description": "Go jogging for 1 minute and 21 seconds.",
	           "result": "",
	           "tester": ""
	         },
	         {
	           "description": "Relax for 2 seconds",
	           "result": "",
	           "tester": ""
	         },
	         {
	           "description": "Work until death. (But click \"completed\" just before you die)",
	           "result": "",
	           "tester": ""
	         }
	       ],
		"testers": [
	        {
	          "email": "noe@aland-marwan.de",
	          "name": "Johannes"
	        },
	        {
	          "email": "dev@aland-mariwan.de",
	          "name": "Gimpel 2"
	        }
	      ],
		"id": "c448447d-f01d-406f-91ee-cf8a6ae6fd3c",
		"createdAt": "createdAt"
		}`

	runEvent(ws, str2)

	//str1 := `{"event": "allOrder","userName":"Rabiyya","password":"myroot", "data": "Good bye Aland\n"}`
	//runEvent(ws, str1)
	//str := `{"event": "message","userName":"","password":"", "data": "Good bye Aland\n"}`
	//runEvent(ws, str)
}

func jenkins() {
	httpposturl := "http://localhost:3000/jenkins"
	fmt.Println("HTTP JSON POST URL:", httpposturl)

	var jsonData = []byte(`{
		"callbackUrl": "callbackUrl",
		"creator": {
	   	 "email": "dev@aland-mariwan.de",
	   	 "name": "Aland Marwan"
	   	},
		"mail": {
	        "body": "Bitte teste die neue Firmware, bli bla blub",
	        "subject": "Es gibt einen neuen Test"
	      },
		"tasks": [
	         {
	           "description": "Go jogging for 1 minute and 21 seconds.",
	           "result": "",
	           "tester": ""
	         },
	         {
	           "description": "Relax for 2 seconds",
	           "result": "",
	           "tester": ""
	         },
	         {
	           "description": "Work until death. (But click \"completed\" just before you die)",
	           "result": "",
	           "tester": ""
	         }
	       ],
		"testers": [
	        {
	          "email": "noe@aland-marwan.de",
	          "name": "Johannes"
	        },
	        {
	          "email": "dev@aland-mariwan.de",
	          "name": "Gimpel 2"
	        }
	      ]

		}`)
	request, error := http.NewRequest("POST", httpposturl, bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()

	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println("response Body:", string(body))
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
