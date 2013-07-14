



func UserOnRegistered(data UserData) {
  eventStore.save("user", "registered", data)

  user.OnUserRegistered(data)
}


func HandleReplayEvent(eventName string, data interface{}) {

}