package main

import (
	"context"
	"fmt"
	"log"

	"github.com/tanx-libs/tanx-connector-go/client"
)

type Entry struct {
	MsgHash string `json:"msg_hash"`
	R       string `json:"r"`
	S       string `json:"s"`
}

func main() {
	c, err := client.New(client.TESTNET)
	if err != nil {
		panic(err)
	}

	// first login
	ethAddr := "0xF58001619C165cDd20B5F7A0EDa072Fd13943002"
	ethPrivateKey := "e65c38b42af2e20540fde19d10bec7fb752ab58852e466151747abc08ae2494a"
	_, err = c.Login(context.TODO(), ethAddr, ethPrivateKey)
	if err != nil {
		panic(err)
	}

	// add correct valueshere to test the order
	// starkPrivateKey := "588ad6325783739b3806e27feebeb1c270d4c6875c29517ff1d689e777b13a6"
	// order, err := c.OrderCreate(context.TODO(), starkPrivateKey, client.OrderOptions{
	// 	Market:  "btcusdc",
	// 	OrdType: "market",
	// 	Side:    "buy",
	// 	Volume:  0.0001,
	// })
	// if err != nil {
	// 	log.Println(err)
	// }

	// order, err := c.OrdersList(context.TODO(), client.OrderListOptions{
	// 	Limit: 3,
	// })
	// if err != nil {
	// 	log.Println(err)
	// }
	order, err := c.OrderGet(context.TODO(), 1254)
	if err != nil {
		log.Println(err)
	}

	fmt.Printf("%+v", order)
}
