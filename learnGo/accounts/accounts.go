package accounts

import (
	"errors"
	"fmt"
)

// Account struct
type Account struct {
	owner   string
	balance int
}

// NewAcoount creates Account
func NewAccount(owner string) *Account {
	account := Account{owner: owner, balance: 0}
	return &account
}

// Deposit from your account
func (a *Account) Deposit(amount int) {
	a.balance += amount
}

// Withdraws from your account

func (a *Account) Withdraw(amount int) error {
	if a.balance < amount {
		return errors.New("can't withdraw")
	}
	a.balance -= amount
	return nil
}

func (a *Account) Balance() int {
	return a.balance
}

// ChangeOwner of the account
func (a *Account) ChangeOwner(newOwner string) {
	a.owner = newOwner
}

func (a *Account) Owner() string {
	return a.owner
}

func (a *Account) String() string {
	return fmt.Sprint(a.Owner(), "'s account.\nbalance: ", a.Balance())
}
