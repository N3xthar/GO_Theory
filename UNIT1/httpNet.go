package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {

	fmt.Printf("Accessing data across internet\n")

	for _, url := range os.Args[1:] {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("Error while getting the url %s: %v\n", url, err)
			continue
		}

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error while reading the string")
			resp.Body.Close()
			continue
		}

		resp.Body.Close()

		// Fixed: Changed Println to Printf so %s prints the text, not raw bytes
		fmt.Printf("%s\n", data)
	}
}
