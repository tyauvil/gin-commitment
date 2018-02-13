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
		log.Fatal(err)
	}
	randint64 := int64(binary.BigEndian.Uint64(randbytes))
	rand.Seed(randint64)
	log.Println("Starting gin-commitment server...")
}

func loadMessages() []string {
	file, err := ioutil.ReadFile("./commit_messages.txt")
	if err != nil {
		log.Fatal(err)
	}
	messages := strings.Split(string(file), "\n")
	return messages
}

func randMessage(m []string) string {
	r := rand.Intn(len(m))
	return m[r]
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	messages := loadMessages()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": randMessage(messages),
		})
	})

	r.GET("/healthz", func(c *gin.Context) {
		c.String(200, "42")
	})

	r.GET("/commit.txt", func(c *gin.Context) {
		c.String(200, randMessage(messages))
	})

	return r
}

func main() {
	r := setupRouter()

	r.Run()
}
