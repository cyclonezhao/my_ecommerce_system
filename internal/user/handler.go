package user

import (
	"fmt"
	"net/http"
)

func SayHello(writer http.ResponseWriter, request *http.Request) {
	name := request.FormValue("name")
	msg := fmt.Sprintf("Hi, I am %s", name)
	fmt.Println(msg)
	fmt.Fprintf(writer, msg)

}
