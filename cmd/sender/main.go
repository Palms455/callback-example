package main

import (
	"bytes"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
}

func main() {
	callbackAddr, _ := os.LookupEnv("CALLBACK_ADDR")
	senderAddr, _ := os.LookupEnv("SENDER_ADDR")


	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	go func() {
		client := &http.Client{Timeout: 10 * time.Second}

		for {
			time.Sleep(5 * time.Second)

			ids := make([]string, rnd.Int31n(300))
			for i := range ids {
				ids[i] = strconv.Itoa((rnd.Int()) % 100)
				if ids[i] == "0" {
					ids[i] = "1"
				}
			}

			log.Printf("=> send %d object IDs\n", len(ids))

			body := bytes.NewBuffer([]byte(fmt.Sprintf(`{"object_ids":[%s]}`, strings.Join(ids, ","))))
			resp, err := client.Post(callbackAddr, "application/json", body)
			if err != nil {
				fmt.Println(err)
				continue
			}
			_ = resp.Body.Close()
		}
	}()

	http.HandleFunc("/objects/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Duration(rnd.Int63n(4000)+300) * time.Millisecond)

		idRaw := strings.TrimPrefix(r.URL.Path, "/objects/")
		id, err := strconv.Atoi(idRaw)
		if err != nil {
			http.Error(w, "invalid url", http.StatusBadRequest)
			return
		}

		_, err = w.Write([]byte(fmt.Sprintf(`{"id":%d,"online":%v}`, id, rnd.Int()%2 == 0)))
		if err != nil {
			log.Println("Err durig write response: %s", err)
		}
	})
	go func() { _ = http.ListenAndServe(senderAddr, nil) }()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	fmt.Println("exit")
}
