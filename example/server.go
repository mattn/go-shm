package main

import (
	"fmt"
	"github.com/mattn/go-shm"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	sm, err := shm.New("mattn", 200)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		sc := make(chan os.Signal)
		signal.Notify(sc, os.Interrupt)
		<-sc
		err = sm.Rm()
		if err != nil {
			log.Fatal(err)
		}
	}()

	defer func() {
		recover()
	}()

	for n := 0; n < 100; n++ {
		fmt.Println(string(sm.Data()[:200]))
		time.Sleep(1 * time.Second)
	}
}
