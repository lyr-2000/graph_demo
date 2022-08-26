package main

import (
	"encoding/json"
	"fmt"
)

type A struct {
	B string `json:"b"`
}

func (u *A) MarkshalJSON() ([]byte, error) {
	return []byte("{b:1}"), nil
}
func do2() {
	do3()
	fmt.Println("---")
}
func do3() {

}
func do1() {

	do2()
}

//go:generate pwd
func main() {
	pwd := "aaa"
	json.Marshal(pwd)
}
