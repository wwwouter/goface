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

func main() {

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

      if se.MetadataTypeName == "RegisterUserData"{
        var registerUserData RegisterUserData
        json.Unmarshal([]byte(se.Metadata), &registerUserData)
        event.Metadata = registerUserData
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
			PublishEvent(event.Name, event.Metadata)

		}
	}

	//sort (date)
	//replay events

	r := mux.NewRouter()
	r.HandleFunc("/", GetIndex)
	r.HandleFunc("/reg", RegisterUser)
	r.HandleFunc("/users", GetUsers)
  r.HandleFunc("/change", ChangeUsername)
	r.HandleFunc("/user", GetUserByUsername)
	http.Handle("/", r)
	err2 := http.ListenAndServe(":8080", nil)
	if err2 != nil {
		log.Fatal(err)
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

func PublishEvent(name string, metadata interface{}) {


	StoreEvent(Event{Name: name, Metadata: metadata, HappenedOn: time.Now()})

	if name == "UserRegistered" {
    GetUsersHandleRegisterUser(metadata.(RegisterUserData))
    GetUserByUsernameHandleRegisterUser(metadata.(RegisterUserData))
		ValidateUniqueUsernameHandleRegisterUser(metadata.(RegisterUserData))
	}
	if name == "UsernameChanged" {
    GetUsersHandleChangeUsername(metadata.(ChangeUsernameData))
    GetUserByUsernameHandleChangeUsername(metadata.(ChangeUsernameData))
		ValidateUniqueUsernameHandleChangeUsername(metadata.(ChangeUsernameData))
	}
}


// ** Register User

type RegisterUserData struct {
	Username string
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	//get username from query string / form / whatever
	username := "Kílian" + strconv.FormatInt(time.Now().Unix(), 10)

  if ValidateUniqueUsername(username){
      PublishEvent("UserRegistered", RegisterUserData{Username: username})
      fmt.Fprintf(w, "Done!")
  } else{
    fmt.Fprintf(w, "Could not create user: user already exists")

  }

	//validation


}


// ** GetUsers

var getUsers_store = make([]GetUsersData, 0)

type GetUsersData struct {
	Username string
}

func GetUsersHandleRegisterUser(data RegisterUserData) {
	getUsers_store = append(getUsers_store, GetUsersData{Username: data.Username})
}

func GetUsersHandleChangeUsername(data ChangeUsernameData) {

	for i, user := range getUsers_store {
		if user.Username == data.OriginalName {
			user.Username = data.NewName
			getUsers_store[i] = user
		}
	}
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	data2, err := json.Marshal(getUsers_store)
	if err != nil {
		fmt.Printf("Result err: %s", err)
	}
	fmt.Fprintf(w, string(data2))
}


// ** ChangeUsername

type ChangeUsernameData struct {
	OriginalName string
	NewName      string
}

func ChangeUsername(w http.ResponseWriter, r *http.Request) {

	originalName := getUsers_store[0].Username
	newName := "NEW!" + strconv.FormatInt(time.Now().Unix(), 10)

  if ValidateUniqueUsername(newName){
      PublishEvent("UsernameChanged", ChangeUsernameData{OriginalName: originalName, NewName: newName})
      fmt.Fprintf(w, "Done!")
  } else{
    fmt.Fprintf(w, "Could not change user: user already exists")

  }


}


// ** GetUserByUsername

var getUserByUsername_store = make([]GetUserByUsernameData, 0)

type GetUserByUsernameData struct {
  Username string
  Html string
}
func GetUserByUsername(w http.ResponseWriter, r *http.Request) {
  username := getUsers_store[0].Username

  var foundUser *GetUserByUsernameData
  for _, user := range getUserByUsername_store {
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

func CreateGetUserByUsernameHtml(username string) string {
  return "<h1>" + username + "</h1>"
}

func GetUserByUsernameHandleRegisterUser(data RegisterUserData) {
  getUserByUsername_store = append(getUserByUsername_store, GetUserByUsernameData{Username: data.Username,
    Html: CreateGetUserByUsernameHtml(data.Username)})
}

func GetUserByUsernameHandleChangeUsername(data ChangeUsernameData) {

  for i, user := range getUserByUsername_store {
    if user.Username == data.OriginalName {
      user.Username = data.NewName
      user.Html = CreateGetUserByUsernameHtml(data.NewName)
      getUserByUsername_store[i] = user
    }
  }
}


// ** GetIndex
func GetIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "Register by clicking here <a href='reg'>HEERRREEE</a>.")

}



// ** validator

var validateUniqueUsername_store = make([]string, 0)


func ValidateUniqueUsername(usernameToValidate string) bool {


  for _, username := range validateUniqueUsername_store {
    if username == usernameToValidate {
      return false
    }
  }

  return true
}

func ValidateUniqueUsernameHandleRegisterUser(data RegisterUserData) {
  validateUniqueUsername_store = append(validateUniqueUsername_store, data.Username)
}

func ValidateUniqueUsernameHandleChangeUsername(data ChangeUsernameData) {
  for i, username := range validateUniqueUsername_store {
    if username == data.OriginalName {
      validateUniqueUsername_store[i] = data.NewName
    }
  }
}

