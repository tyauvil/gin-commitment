// Golang version of 

package main

import (
	//"fmt"
	"math/rand"
	"io/ioutil"
	"strings"

	"github.com/gin-gonic/gin"
)

//func name() string {
//}

func load_messages() []string {
	file, _ := ioutil.ReadFile("./commit_messages.txt")
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
	//names := []string{"Nick", "Steve", "Andy", "Qi", "Fanny", "Sarah", "Cord", "Todd", "Chris", "Pasha", "Gabe", "Tony", "Jason", "Randal", "Ali", "Kim", "Rainer", "Guillaume", "Kelan", "David", "John", "Stephen", "Tom", "Steven", "Jen", "Marcus", "Edy", "Rachel"}
	messages := load_messages()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": message(messages),
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
