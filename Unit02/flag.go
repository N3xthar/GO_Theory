package main

import (
	"flag"
	"fmt"
)

func main() {
	// command line arguments
	// // third variable is helper text for the flag
	name := flag.String("name", "Guest", "insert your name ")
	age := flag.Int("age", 18, "enter your age ")
	admin := flag.Bool("admin", false, "Is this admin")

	// parse the command line arguments
	flag.Parse()

	// derefrence
	fmt.Println("Name :", *name)
	fmt.Println("age :", *age)
	fmt.Println("Admin : ", *admin)
}
