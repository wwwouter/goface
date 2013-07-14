import "github.com/gorilla/mux"

func main() {
  app:=app.create()
  server := new Server()
  server.app = app
}



type Server struct  {
  app App

  func registerUser() `route-path:users` `route-method:POST`
}


routeToMethodWrapper dingie (mag een ref hebben naar mux)
serialize param to input
serialeze json body to struct

func (this *Server) init() {

  r := mux.NewRouter()
  r.HandleFunc("/", this.HomeHandler).Methods("GET")
  r.HandleFunc("/users", UsersHandler).Methods("GET")
  r.HandleFunc("/users", UsersHandler).Methods("POST")
  http.Handle("/", r)

  util.SetupHttpHandlers(this)
  app.Events.OnUserRegistered(this.OnUserRegistered)

  util.AddHttpHandler("users", "GET", this.getUsers)
  util.Start(this)

  //userDetails gorest.EndPoint `method:"GET" path:"/users/{Id:int}" output:"User"`

// http://godoc.org/bitbucket.org/rj/httprouter-go#_example_NewBuilder
//   b.HandleFunc("/login", PostMethod, func(w http.ResponseWriter, r *http.Request) {
//     fmt.Fprintf(w, "<html><body>You're logged in!</body></html>")
// })

  // https://github.com/sdming/wk
  // server.RouteTable.Get("/data/int/{p0}?").ToFunc(model.DataByInt)

}

type OnUserRegisteredData struct{
  name string
}


func (this *Server) OnUserRegistered(data UserEvent.RegisteredData) {
  dataOut := new OnUserRegisteredData()
  dataOut.name = data.name
  clientConnections.notify("user", "registered", dataOut)
}

type RegisterUserData struct{
  Name string `json:"name"` `validation:"required"` `validation:"max-length=10`
}


`routes-path:"user/"` `routes-method:"POST"`
func (this *Server) registerUser(request *HttpRequest, response *HttpResponse, user RegisterUserData) {
  error = app.Commands.User.RegisterUser(user)
  if(error){
    logger.error("Error on registering user", error)
    return response.status(500).send("Something went wrong")
  }
  response.status(200).send("Ok")
}


type Endpoint struct{

}

func (this *Endpoint) Register() {
  util.route(this.getUsers, "users", "GET")
  util.route(this.getUsers, "users", "POST")
}

func (this *Server) getUsers(request *HttpRequest, response *HttpResponse) {
  users := this.app.WebQueryStore.GetUsers()
  response.json(users)
}


func (this * Server) RegisterGetUsers() path string, method string,handler func{
  return "users", "GET", func getUsers(request *HttpRequest, response *HttpResponse) {
    users := this.app.WebQueryStore.GetUsers()
    response.json(users)
  }
}


type POST struct{
  value string `route-info:method`
}
type GET struct{
  value string `route-info:method`
}
type PathUsers struct{
  value string `route-path:users`
}


func (this *Server) getUsers(method GET, path PathUsers, request *HttpRequest, response *HttpResponse) {
  users := this.app.WebQueryStore.GetUsers()
  response.json(users)
}



func (this *Server) getUsers_GET_users(request *HttpRequest, response *HttpResponse) {
  users := this.app.WebQueryStore.GetUsers()
  response.json(users)
}



func (this *Server) getUsers(request *HttpRequest, response *HttpResponse) {
  users := this.app.WebQueryStore.GetUsers()
  response.json(users)
}

func (this *Server) Get_Users(request *HttpRequest, response *HttpResponse) {
  users := this.app.WebQueryStore.GetUsers()
  response.json(users)
}

//  users/:id/posts
func (this *Server) Get_Users_id_Posts(request *HttpRequest, response *HttpResponse) {
  users := this.app.WebQueryStore.GetUsers()
  response.json(users)
}


func (this *Server) SetupHandlers() {
  handlers := new Handlers()

  func getUsers(request *HttpRequest, response *HttpResponse) {
    users := this.app.WebQueryStore.GetUsers()
    response.json(users)
  }
  handlers.add("users", "GET", getUsers)
}




handlers := new Handlers()

func (this *Server) getUsers(request *HttpRequest, response *HttpResponse) {
  users := this.app.WebQueryStore.GetUsers()
  response.json(users)
}
handlers.add("users", "GET", getUsers)

