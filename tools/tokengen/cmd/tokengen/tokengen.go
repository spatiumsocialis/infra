package main

import (
	"flag"
	"fmt"

	"github.com/spatiumsocialis/infra/pkg/common"
	"github.com/spatiumsocialis/infra/pkg/common/auth"
)

func main() {
	common.LoadEnv()
	uid := flag.String("u", auth.TestUID, "uid to generate token for")
	flag.Parse()
	token, err := auth.GenerateToken(*uid)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(token)

}
