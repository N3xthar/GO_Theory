package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc("/catch-me", handler)
	http.HandleFunc("/ping", ping)

	fmt.Println("Server is running on :8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {

	// Welcome message
	fmt.Fprintln(w, "Welcome Dude")
	fmt.Fprintln(w, "================")

	// Read query parameter
	val := r.URL.Query().Get("bigData")

	// Check if parameter exists
	if val == "" {
		fmt.Fprintln(w, "Usage:")
		fmt.Fprintln(w, "http://localhost:8080/catch-me?bigData=5")
		return
	}

	// Convert string to integer
	count, err := strconv.Atoi(val)
	if err != nil {
		http.Error(w, "bigData must be an integer", http.StatusBadRequest)
		return
	}

	// Print Welcome count times
	for i := 1; i <= count; i++ {
		fmt.Fprintf(w, "Welcome %d\n", i)
	}
}

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello Aman Deep, how are you my bro?")
}
