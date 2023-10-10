package main

import (
	"encoding/json"
)

type Tester struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type TaskJenkinsRequest struct {
	Description string `json:"description"`
}

type OrderRequestJenkins struct {
	CallbackUrl string               `json:"callbackUrl"`
	Creator     Tester               `json:"creator"`
	Description string               `json:"description"`
	ExpireDate  string               `json:"expireDate"`
	Mail        MailStruck           `json:"mail"`
	Tasks       []TaskJenkinsRequest `json:"tasks"`
	Testers     []Tester             `json:"testers"`
}

type MailStruck struct {
	Body    string `json:"body"`
	Subject string `json:"subject"`
}

type Task struct {
	Description string `json:"description"`
	Result      string `json:"result"`
	Tester      string `json:"tester"`
	Comment     string `json:"comment"`
}

func (p Task) Equals(other Task) bool {
	return p.Description == other.Description && p.Result == other.Result && p.Tester == other.Tester && p.Comment == other.Comment
}

type OrderAnswerJenkins struct {
	CallbackUrl  string     `json:"callbackUrl"`
	Creator      Tester     `json:"creator"`
	Description  string     `json:"description"`
	ExpireDate   string     `json:"expireDate"`
	Mail         MailStruck `json:"mail"`
	Tasks        []Task     `json:"tasks"`
	Testers      []Tester   `json:"testers"`
	Id           string     `json:"id"`
	CreatedAt    string     `json:"createdAt"`
	UpdatedAt    string     `json:"updatedAt"`
	TestComplete bool       `json:"testComplete"`
}

func (x *OrderAnswerJenkins) OrderAnswerJenkinsFromRaw(rawData []byte) (*OrderAnswerJenkins, error) {
	output := &OrderAnswerJenkins{}
	err := json.Unmarshal(rawData, output)
	return output, err
}

type OrderAnswerFlutter struct {
	Name   string    `json:"event"`
	Orders []Order   `json:"orders"`
	Stump  Signature `json:"stump"`
}

type Order struct {
	Name         string     `json:"event"`
	CallbackUrl  string     `json:"callbackUrl"`
	Creator      Tester     `json:"creator"`
	Description  string     `json:"description"`
	ExpireDate   string     `json:"expireDate"`
	Mail         MailStruck `json:"mail"`
	Tasks        []Task     `json:"tasks"`
	Testers      []Tester   `json:"testers"`
	TestComplete bool       `json:"testComplete"`
	Id           string     `json:"id"`
	CreatedAt    string     `json:"createdAt"`
	UpdatedAt    string     `json:"updatedAt"`
	Stump        Signature  `json:"stump"`
}

func (e *Order) RawOrder() []byte {
	raw, _ := json.Marshal(e)
	return raw
}

type Signature struct {
	Index     int    `json:"index"`
	CreatedAt string `json:"createdAt"`
}

type OrderDB struct {
	CallbackUrl  string     `clover:"callbackUrl"`
	Creator      Tester     `clover:"creator"`
	Description  string     `clover:"description"`
	ExpireDate   string     `clover:"expireDate"`
	Mail         MailStruck `clover:"mail"`
	Tasks        []Task     `clover:"tasks"`
	Testers      []Tester   `clover:"testers"`
	TestComplete bool       `clover:"testComplete"`
	Id           string     `clover:"id"`
	CreatedAt    string     `clover:"createdAt"`
	UpdatedAt    string     `clover:"updatedAt"`
}

func tasksFromTaskJenkinsRequest(task []Task) []Task {
	output := []Task{}
	for i := 0; i < len(task); i++ {
		output = append(output, Task{Description: task[i].Description, Result: "", Tester: "", Comment: ""})
	}

	return output
}

type EventHandler func(*Order)

// A method that is attached to the OrderAnswerFlutter struct. It returns a byte array.
func (e *OrderAnswerFlutter) Raw() []byte {
	raw, _ := json.Marshal(e)
	return raw
}

// A function that takes a byte array and returns a pointer to an OrderAnswerFlutter and an error.
func (x *OrderAnswerFlutter) NewEventFromRaw(rawData []byte) (*OrderAnswerFlutter, error) {
	output := &OrderAnswerFlutter{}
	err := json.Unmarshal(rawData, output)
	return output, err
}

// A method attached to the OrderAnswerJenkins struct. It returns a byte array.
func (e *OrderAnswerJenkins) Raw() []byte {
	raw, _ := json.Marshal(e)
	return raw
}

// A method attached to the OrderRequestJenkins struct. It takes a byte array and returns a pointer to
// an OrderRequestJenkins and an error.
func (x *OrderRequestJenkins) NewEventFromRaw(rawData []byte) (*OrderRequestJenkins, error) {
	output := &OrderRequestJenkins{}
	err := json.Unmarshal(rawData, output)
	return output, err
}

// A method attached to the Order struct. It returns a byte array and an error.
func (e *Order) Raw() ([]byte, error) {
	raw, err := json.Marshal(e)
	return raw, err
}

// A method attached to the Order struct. It takes a byte array and returns a pointer to an Order and
// an error.
func (x *Order) NewEventFromRaw(rawData []byte) (*Order, error) {
	output := &Order{}
	err := json.Unmarshal(rawData, output)
	return output, err
}

// It takes a Tester and returns a Tester
func cloneTester(tester Tester) Tester {
	return Tester{Email: tester.Email, Name: tester.Name}
}

// It clones a mail.
func cloneMail(mail MailStruck) MailStruck {
	return MailStruck{Body: mail.Body, Subject: mail.Subject}
}

// It takes a value of type `Signature` and returns a copy of that value
func cloneSignature(value Signature) Signature {
	return Signature{Index: value.Index, CreatedAt: value.CreatedAt}
}

// It takes a slice of Task structs and returns a slice of Task structs
func cloneTasks(tasks []Task) []Task {
	output := []Task{}
	for i := 0; i < len(tasks); i++ {
		output = append(output, Task{Description: tasks[i].Description, Result: tasks[i].Result, Tester: tasks[i].Tester, Comment: tasks[i].Comment})
	}
	return output
}

// It clones a slice of `Tester`s
func cloneTesters(testers []Tester) []Tester {
	output := []Tester{}
	for i := 0; i < len(testers); i++ {
		output = append(output, Tester{Email: testers[i].Email, Name: testers[i].Name})
	}
	return output
}
