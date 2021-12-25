package main

import "fmt"

func main() {
	c := make(chan string)

	people := [2]string{"yongckim", "nico"}
	for _, person := range people {
		go isSexy(person, c)
	}
	for i := 0; i < len(people); i++ {
		fmt.Println(i)
		fmt.Println(<-c)
	}

}

func isSexy(person string, c chan string) {
	c <- person + " is sexy"
}
