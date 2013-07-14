package main

import (
    "fmt"
    "net/http"
    "github.com/gorilla/mux"
)


func GetIntFunc() func() int{
  return func () int {
    return 1
  }
}
func GetHandlerFunc() func(http.ResponseWriter,*http.Request){
  return func (w http.ResponseWriter, r *http.Request)  {
    fmt.Fprintf(w, "Hi there, I lovezzz %s!", r.URL.Path[1:])
  }
}

type WebServer struct{

}

// type Handler interface {
//     ServeHTTP(http.ResponseWriter, *http.Request)
// }

type GetUsersEndpoint struct{

}
func (this *GetUsersEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I asasdsadas %s!", r.URL.Path[1:])

}

func Wrap(s *Server, h func(*Server,http.ResponseWriter,*http.Request)) func(http.ResponseWriter,*http.Request) {
  return func (w http.ResponseWriter, r *http.Request)  {
    h(s,w,r)
  }
}

func (this *WebServer) Start(s *Server, h func(*Server,http.ResponseWriter,*http.Request)) {
  // fmt.Print("Start")
  // fmt.Printf("%v\n", GetIntFunc()())

  r := mux.NewRouter()
  // r.HandleFunc("/", GetHandlerFunc())
  r.HandleFunc("/", Wrap(s,h))
  // gue := &GetUsersEndpoint{}
  // r.Handle("/", gue)
  http.Handle("/", r)
  // http.HandleFunc("/", GetHandlerFunc())
  http.ListenAndServe(":8080", nil)


}

type Server struct{

}

func (this *Server)Start() {
  webServer := WebServer{}
  webServer.Start(this, (*Server).handler2)
}



func (this *Server) handler2(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi therezzz, I love %s!", r.URL.Path[1:])
}



var a = 1

func handler(w http.ResponseWriter, r *http.Request) {
  server := Server{}

  server.handler2(w,r)
    //fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {


  server := Server{}
  server.Start()
}