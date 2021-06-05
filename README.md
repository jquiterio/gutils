gutils
======

Production-ready golang utils.

UUID
----

RFC4122 uuid generator 

Usage:
```go
import "github.com/jquiterio/gutils/uuid"

type User struct {
  ID uuid.UUID `json:"id"`
  Name string `json:"name"`
}

func newUser() User {
  var u User
  u.ID = uuid.New() // generates a new V4 uuid
  // for V4 uuid.NewV5(uuid.NamespaceURL, []byte("github.com"))
  u.Name = "Isaac"
  return u
}


// with gorm
func main() {
  var db *gorm.DB
  // db parameters... here
  var user User
  user := newUser()
  db.Create(&user)

  var userFromDB User
  db.First(&userFromDB)

  fmt.Println(user)
  fmt.Println(userFromDB)

}```
