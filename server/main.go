package main

import (
	"fmt"
	"log"
	"net/http"
	"golang-react-todo/router"
)

func main(){
	r := router.Router()
	fmt.Println("starting server on port 4000")
	log.Fatal(http.ListenAndServe(":4000",r))
}
