package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)



type Endpointer interface{
  HandleHttp(w http.ResponseWriter, r *http.Request)
  GetRoute() (string, string)
}

type EventHandler interface{
  HandleEvent(event interface{})
}

type Validater interface{
  ValidateEvent(event interface{}) (bool, string)
}


var eventHandlers = make([]EventHandler,0)
var validaters = make([]Validater,0)

func main() {



  validateUniqueUsername := &ValidateUniqueUsername{}
  eventHandlers = append(eventHandlers, validateUniqueUsername)

  validaters = append(validaters, validateUniqueUsername)

	r := mux.NewRouter()

  getIndexEndpoint := GetIndexEndpoint{}
  HandleEndpoint(&getIndexEndpoint, r)

  registerUserEndpoint := RegisterUserEndpoint{}
  HandleEndpoint(&registerUserEndpoint, r)

  getUsersEndpoint := GetUsersEndpoint{}
  eventHandlers = append(eventHandlers, &getUsersEndpoint)
  HandleEndpoint(&getUsersEndpoint, r)

  changeUsernameEndpoint := ChangeUsernameEndpoint{}
  HandleEndpoint(&changeUsernameEndpoint, r)

  getUserByUsernameEndpoint := GetUserByUsernameEndpoint{}
  eventHandlers = append(eventHandlers, &getUserByUsernameEndpoint)
  HandleEndpoint(&getUserByUsernameEndpoint, r)



  //load all from data
  files, err := ioutil.ReadDir("data")

  for _, file := range files {
    if strings.HasSuffix(file.Name(), ".json") {

      filebyte, error := ioutil.ReadFile("data/" + file.Name())
      if error != nil {
        log.Fatal("Could not read file " + file.Name() + " to parse")
        os.Exit(1)
      }

      var se StoredEvent
      json.Unmarshal(filebyte, &se)

      event := Event{Name: se.Name, HappenedOn : se.HappenedOn}

      if se.MetadataTypeName == "UserRegisteredEvent"{
        var userRegisteredEvent UserRegisteredEvent
        json.Unmarshal([]byte(se.Metadata), &userRegisteredEvent)
        event.Metadata = userRegisteredEvent
      } else if se.MetadataTypeName == "ChangeUsernameData" {
        var changeUsernameData ChangeUsernameData
        json.Unmarshal([]byte(se.Metadata), &changeUsernameData)
        event.Metadata = changeUsernameData
      } else {
        fmt.Printf("\nResult err: %s\n %#v\n\n", file.Name(), se)
        panic("What?? " + se.MetadataTypeName)
      }

      // sourceType := reflect.TypeOf(e)

   //    fmt.Printf("Result err: %#v", sourceType)
      // fmt.Printf("Result err: %s", sourceType.Name())
      fmt.Printf("Result event: %s", event)
      PublishEventWithoutValidation(event.Name, event.Metadata)
      eventCounter ++

    }
  }


	http.Handle("/", r)
	err2 := http.ListenAndServe(":8080", nil)
	if err2 != nil {
		log.Fatal(err)
	}

}


func HandleEndpoint(endpointer Endpointer, r *mux.Router) {
  path, method := endpointer.GetRoute()
  r.HandleFunc(path, Wrap(endpointer)).Methods(method)


}

func Wrap(endpointer Endpointer) func(http.ResponseWriter,*http.Request) {
  return func (w http.ResponseWriter, r *http.Request)  {
   endpointer.HandleHttp(w, r)
  }
}



type Event struct {
  Name       string
  Metadata   interface{}
  HappenedOn time.Time
}

type StoredEvent struct {
	Name       string
	Metadata   string
	HappenedOn time.Time
  MetadataTypeName string
}

var eventCounter = 0

func StoreEvent(event Event) {
	to_file := "data/" + strconv.Itoa(eventCounter) + ".json"

	file_handle, err := os.OpenFile(to_file, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	/* Automatically close when we finish in this function, consider
	 * with open(to_string) as file_handle. */
	defer file_handle.Close()

  metadataType := reflect.TypeOf(event.Metadata)

  metadataJson, err := json.Marshal(event.Metadata)
  if err != nil {
    fmt.Printf("Result err: %s", err)
  }


  storedEvent := StoredEvent{Name: event.Name, HappenedOn : event.HappenedOn,
    MetadataTypeName : metadataType.Name(), Metadata : string(metadataJson) }

	data2, err := json.Marshal(storedEvent)
	if err != nil {
		fmt.Printf("Result err: %s", err)
	}

	if _, err = file_handle.Write(data2); err != nil {
		panic(err)
	}

	eventCounter++
}

func PublishEventWithoutValidation(name string, metadata interface {}) {

  for _, eventHandler := range eventHandlers{
    eventHandler.HandleEvent(metadata)
  }
}

func PublishEvent(name string, metadata interface{}) (success bool, message string) {
  for _, validater := range validaters{
      validationSuccess, validationMessage := validater.ValidateEvent(metadata)
      if ! validationSuccess{
        return false, validationMessage
      }
  }

  StoreEvent(Event{Name: name, Metadata: metadata, HappenedOn: time.Now()})

  PublishEventWithoutValidation(name, metadata)

  return true, ""
}


// **** Events

type UserRegisteredEvent struct{
  Username string
}

type ChangeUsernameData struct {
  OriginalName string
  NewName      string
}


// **** Endpoints


// ***** GetIndex

type GetIndexEndpoint struct{

}

func (this *GetIndexEndpoint) HandleHttp(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Content-Type", "text/html")
  fmt.Fprintf(w, "Register by clicking here <a href='reg'>HEERRREEE</a>.")
}

func(this *GetIndexEndpoint) GetRoute() (string, string){
  return "/", "GET"
}


// **** RegisterUser

type RegisterUserEndpoint struct{

}

func (this *RegisterUserEndpoint) HandleHttp(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Content-Type", "text/html")
  //get username from query string / form / whatever
  username := "username"


  success, message := PublishEvent("UserRegistered", UserRegisteredEvent{Username: username})
  if success{
      fmt.Fprintf(w, "Done!")
  } else{
    fmt.Fprintf(w, "Could not create user: " + message)

  }
}

func(this *RegisterUserEndpoint) GetRoute() (string, string){
  return "/reg", "GET"
}


// ** GetUsers



type GetUsersData struct {
  Username string
}

type GetUsersEndpoint struct{
  store []GetUsersData
}


func (this *GetUsersEndpoint) HandleHttp(w http.ResponseWriter, r *http.Request){
  data2, err := json.Marshal(this.store)
  if err != nil {
    fmt.Printf("Result err: %s", err)
  }
  fmt.Fprintf(w, string(data2))
}

func(this *GetUsersEndpoint) GetRoute() (string, string){
  return "/users", "GET"
}

func (this *GetUsersEndpoint)HandleEvent(event interface{}) {
  switch data := event.(type) {
    case UserRegisteredEvent:
      this.store = append(this.store, GetUsersData{Username: data.Username})
    case ChangeUsernameData:
      for i, user := range this.store {
        if user.Username == data.OriginalName {
          user.Username = data.NewName
          this.store[i] = user
          break
        }
      }
  }
}


// ** ChangeUsername


type ChangeUsernameEndpoint struct{
}

func (this *ChangeUsernameEndpoint) HandleHttp(w http.ResponseWriter, r *http.Request){
  originalName := "username"
  newName := "NEW!" + strconv.FormatInt(time.Now().Unix(), 10)

  success, message := PublishEvent("UsernameChanged", ChangeUsernameData{OriginalName: originalName, NewName: newName})
  if success{
      fmt.Fprintf(w, "Done!")
  } else {
    fmt.Fprintf(w, "Could not change user: " + message)
  }
}

func(this *ChangeUsernameEndpoint) GetRoute() (string, string){
  return "/change", "GET"
}




// ** GetUserByUsername

type GetUserByUsernameData struct {
  Username string
  Html string
}

type GetUserByUsernameEndpoint struct{
  store []GetUserByUsernameData
}


func (this *GetUserByUsernameEndpoint) HandleHttp(w http.ResponseWriter, r *http.Request){
 username := "username"

  var foundUser *GetUserByUsernameData
  for _, user := range this.store {
    if user.Username == username{
      foundUser = &user
      break
    }
  }
  w.Header().Set("Content-Type", "text/html")
  if foundUser == nil {
    fmt.Fprintf(w, "User not found.")

  } else{

    fmt.Fprintf(w, foundUser.Html)
  }
}

func(this *GetUserByUsernameEndpoint) GetRoute() (string, string){
  return "/user", "GET"
}


func CreateGetUserByUsernameHtml(username string) string {
  return "<h1>" + username + "</h1>"
}

func (this *GetUserByUsernameEndpoint)HandleEvent(event interface{}) {
  switch data := event.(type) {
    case UserRegisteredEvent:
      this.store = append(this.store, GetUserByUsernameData{Username: data.Username,
        Html: CreateGetUserByUsernameHtml(data.Username)})
    case ChangeUsernameData:
      for i, user := range this.store {
        if user.Username == data.OriginalName {
          user.Username = data.NewName
          user.Html = CreateGetUserByUsernameHtml(data.NewName)
          this.store[i] = user
          break
        }
      }
  }
}









// ******* Validators



type ValidateUniqueUsername struct{
  store []string
}



func (this *ValidateUniqueUsername) ValidateEvent(event interface{}) (bool, string) {
  switch data := event.(type) {
    case UserRegisteredEvent:
      return this.ValidateUniqueUsername(data.Username)
    case ChangeUsernameData:
      return this.ValidateUniqueUsername(data.NewName)
  }
  return true, ""
}

func (this *ValidateUniqueUsername) ValidateUniqueUsername(usernameToValidate string) (bool, string) {
  usernameToValidate = strings.TrimSpace(strings.ToLower(usernameToValidate))
  for _, username := range this.store {
    if username == usernameToValidate {
      return false, "Username already exists."
    }
  }
  return true, ""
}

func (this *ValidateUniqueUsername) HandleEvent(event interface{}) {
  switch data := event.(type) {
    case UserRegisteredEvent:
      username := strings.TrimSpace(strings.ToLower(data.Username))
      this.store = append(this.store, username)
    case ChangeUsernameData:
      for i, username := range this.store {
        if username == data.OriginalName {
          username := strings.TrimSpace(strings.ToLower(data.NewName))
          this.store[i] = username
          break
        }
      }
  }

}

