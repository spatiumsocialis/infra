package main

import (
	"flag"
	"fmt"

	"github.com/safe-distance/socium-infra/pkg/common/auth"
)

func main() {
	uid := flag.String("u", auth.TestUID, "uid to generate token for")
	flag.Parse()
	token, err := auth.GenerateToken(*uid)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(token)

}
