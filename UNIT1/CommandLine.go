package main

import (
	"fmt"
	"os"
)

func main() {
	var s string
	for i := 1; i < len(os.Args); i++ {
		s = s + os.Args[i]
	}
	fmt.Println(s)
	fmt.Println("\n\n")

	//  Exercise 1.1:
	//Modify the echo program to also print os.Args[0], the name of the command that invoked it.
	fmt.Println(os.Args[0])
	// Exercise 1.2:
	// Exercise 1.2: Modify the echo program to print the index and value of each of its arguments, one per line
	for index, value := range os.Args {
		fmt.Println("Index is this ", index, "Value is this ", value, "\n")
	}

	// Exercise 1.3: Experiment to measure the difference in running time between our potentially
// inefficient versions and the one that uses strings.Join. (Section 1.6 illustrates part of the
// time package, and Section 11.4 shows how to write benchmark tests for systematic per-
// formance evaluation.)
	
		// Using 
		fmt.Println("\n\n")
		fmt.Println(os.Args[1:]," ")
	}
