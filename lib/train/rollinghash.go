package main

import (
	"fmt"
	/*"github.com/chmduquesne/rollinghash/buzhash64"*/
	"github.com/syncthing/syncthing/lib/rand"
)

func main() {
	/*data := []byte("here is some data to roll on")
	h := buzhash64.New()
	n := 16

	// Initialize the rolling window
	h.Write(data[:n])

	for _, c := range(data[n:]) {

		// Slide the window and update the hash
		h.Roll(c)

		// Get the updated hash value
		fmt.Println(h.Sum64())
	}*/

	var name[] int;
	name = append(name, 1,2,3,4,5)
	for i := range name {
		j := rand.Intn(i + 1)
		name[i], name[j] = name[j], name[i]
	}

	for i := range name {
		fmt.Println(name[i])
	}
}