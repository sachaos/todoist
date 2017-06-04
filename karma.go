package main

import (
	"fmt"

	"github.com/urfave/cli"
)

func Karma(c *cli.Context) error {
	client := GetClient(c)

	fmt.Println(client.Store.User.Karma)
	return nil
}
