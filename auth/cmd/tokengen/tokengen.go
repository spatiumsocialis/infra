package main

import (
	"fmt"

	"github.com/safe-distance/socium-infra/auth"
)

func main() {
	token, err := auth.GenerateToken(auth.TestUID)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(token)

}
