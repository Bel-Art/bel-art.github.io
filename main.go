package main

import (
	"fmt"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)
	defer http.ListenAndServe(":8081", nil)
	fmt.Println("Server listenning on 8081")
}
