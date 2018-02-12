// Golang version of https://github.com/ngerakines/commitment

package main

import (
	cryptorand "crypto/rand"
	"encoding/binary"
	"io/ioutil"
	"log"
	"math/rand"
	"strings"

	"github.com/gin-gonic/gin"
)

func init() {
	randbytes := make([]byte, 8)
	_, err := cryptorand.Read(randbytes)
	if err != nil {
		log.Println("error:", err)
		return
	}
	randint64 := int64(binary.BigEndian.Uint64(randbytes))
	rand.Seed(randint64)
	log.Println("Starting gin-commitment server...")
}

func load_messages() []string {
	file, err := ioutil.ReadFile("./commit_messages.txt")
	if err != nil {
		log.Fatal(err)
	}
	messages := strings.Split(string(file), "\n")
	return messages
}

func message(m []string) string {
	l := len(m)
	r := rand.Intn(l)
	return m[r]
}

func main() {
	r := gin.Default()
	messages := load_messages()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": message(messages),
		})
	})

	r.GET("/healthz", func(c *gin.Context) {
		c.String(200, "42")
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
