

type GetUsersUserData struct{
  Name string `json:"name"`
}

type GetUsersData struct{
  Users GetUsersUserData[] `json:"users"`
}

