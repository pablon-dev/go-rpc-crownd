package main

import (
	"fmt"

	crownd "github.com/pablon-dev/go-rpc-crownd"
)

func main() {
	client, err := crownd.NewClientWithSSL("localhost", 9341, "pablon", "3hRbcibf1oFG5KhgMQNM76Em3sjVd2EiFzVfWc4VocWv", 3)
	if err != nil {
		panic(err)
	}
	resp, err := client.GetInfo()
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
