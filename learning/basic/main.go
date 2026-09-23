package main

import (
	"fmt"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, "Hello from BE_GKMI_NTC_GO")
}

func main() {

	http.HandleFunc("/hello", helloHandler)

	fmt.Println("Server running on http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Server error:", err)
	}
}
