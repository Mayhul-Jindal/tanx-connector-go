package main

import (
	"fmt"

	"github.com/tanx-libs/tanx-connector-go/client"
)

func main() {
	// # pvt key, hash


	// dont add 0x
	fmt.Println(client.Sign("504e27c7c7c3aba8104ae3c50831406a458c08f6c9299d0fafb05a60ab8d51", "27dd634b2534f618d1b0d2bdd67f4ce05202447f6a79a8e778b3863bd01e68c"))
}
