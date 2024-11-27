package main

import (
	"fmt"

	"github.com/newtoallofthis123/noob_social/routes"
)

func main() {
	fmt.Println("Listening on port 8080")

	api := routes.New()

	err := api.Start()
	if err != nil {
		panic(err)
	}
}
