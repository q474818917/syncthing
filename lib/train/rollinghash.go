package main

import (
	"fmt"
	"github.com/chmduquesne/rollinghash/buzhash64"
)

func main() {
	data := []byte("here is some data to roll on")
	h := buzhash64.New()
	n := 16

	// Initialize the rolling window
	h.Write(data[:n])

	for _, c := range(data[n:]) {

		// Slide the window and update the hash
		h.Roll(c)

		// Get the updated hash value
		fmt.Println(h.Sum64())
	}
}