

type RegisterUserData interface{
  name string
}


func RegisterUser(data RegisterUserData) {
  //validation
  user := new User()
  user.name = data.name

}