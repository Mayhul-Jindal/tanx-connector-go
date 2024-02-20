package main

import (
	"context"
	"fmt"

	"github.com/tanx-libs/tanx-connector-go/client"
)

func main() {
	c, err := client.New(client.MAINET)
	if err != nil {
		panic(err)
	}

	// first login
	ethAddr := "0xF58001619C165cDd20B5F7A0EDa072Fd13943002"
	ethPrivateKey := ""
	_, _, err = c.Login(context.TODO(), ethAddr, ethPrivateKey)
	if err != nil {
		panic(err)
	}

	// fmt.Println(nonce, jwt)


	// add correct valueshere to test the order
	starkPrivateKey := ""
	order, err := c.OrderCreate(context.TODO(), starkPrivateKey, "ethusdt", "limit", 100, "buy", 100)
	if err != nil {
		panic(err)
	}

	fmt.Println(order)
}
