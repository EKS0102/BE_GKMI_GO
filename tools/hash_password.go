package main

import (
	"fmt"
	"log"

	"BE_GKMI_NTC_GO/internal/security"
)

func main() {
	hash, err := security.HashPassword("rahasia123")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(hash)
}
