package main

import (
	"fmt"
	"net/http"
)

func main() {
	port := "8081"
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)
	defer http.ListenAndServe(":"+port, nil)
	fmt.Println("Server listenning on http://localhost:" + port)
}
