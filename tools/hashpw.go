package main

import (
 "fmt"
 "golang.org/x/crypto/bcrypt"
)

func main() {
 pw := "test1234" // поменяй при желании
 h, _ := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
 fmt.Println(string(h))
}
