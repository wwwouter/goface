


func AUserCanBeCreated(t *testing.T) {
  user := new User()
  if &user == nil {
    t.Error("user not set")
  }
}




// func WhenAUserIsCreatedItShouldSentCreatedEvent(t *testing.T) {
//   user = new User()
//   if &a == nil {
//     t.Error("a not set")
//   }
// }