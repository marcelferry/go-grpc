package main

import (
    "fmt"
    "net/http"
)

func hello(res http.ResponseWriter, req *http.Request) {
    fmt.Println("/hello endpoint called")
    fmt.Fprintf(res, "hello\n")
}

func main() {
    http.HandleFunc("/hello", hello)
    fmt.Println("Server up and listening...")
    http.ListenAndServe(":80", nil)
}