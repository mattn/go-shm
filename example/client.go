package main

import (
	"github.com/mattn/go-shm"
	"log"
	"os"
	"strings"
)

func main() {
	if len(os.Args) == 0 {
		return
	}
	sm, err := shm.Attach("mattn", 200)
	if err != nil {
		log.Fatal(err)
	}

	for n := 0; n < 200; n++ {
		sm.Data()[n] = 0
	}
	for n, c := range []byte(strings.Join(os.Args[1:], " ")) {
		sm.Data()[n] = c
	}

	err = sm.Detatch()
	if err != nil {
		log.Fatal(err)
	}
}
