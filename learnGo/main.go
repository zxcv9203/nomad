package main

import (
	"fmt"
	"log"

	"github.com/zxcv9203/nomad/learnGo/accounts"
)

func main() {
	account := accounts.NewAccount("yongckim")
	account.Deposit(10)
	fmt.Println(account.Balance())
	err := account.Withdraw(110)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(account.Balance())
}
