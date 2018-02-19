package main

import (
	"bytes"
	cryptorand "crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/acme/autocert"
)

var SourceBranch = "unset"
var SourceCommit = "unset"
var GolangVersion = "unset"

type Format struct {
	Fname  string
	FnameU string
	FnameL string
	NumL   string
	NumM   string
	NumH   string
}

func loadLines(file string) []string {
	f, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(f), "\n")
	return lines
}

func strIntRangeRand(min, max int) string {
	x := rand.Intn(max-min) + min
	return fmt.Sprint(x)
}

func sha256sum(text string) string {
	h := sha256.New()
	h.Write([]byte(text))
	sha256hash := hex.EncodeToString(h.Sum(nil))
	return sha256hash
}

func randMessage() map[string]string {
	rand.Seed(42) //make the messages deterministic for permalink
	names := loadLines("./static/names.txt")
	messages := loadLines("./static/commit_messages.txt")
	mapMsg := make(map[string]string)

	for _, v := range messages {
		rn := rand.Intn(len(names))
		fn := names[rn]
		fnu := strings.ToUpper(fn)
		fnl := strings.ToLower(fn)
		nl := strIntRangeRand(1, 10)
		nm := strIntRangeRand(20, 75)
		nh := strIntRangeRand(50, 99)
		t := Format{fn, fnu, fnl, nl, nm, nh}
		var b bytes.Buffer
		tmpl, err := template.New("wtc").Parse(v)
		if err != nil {
			log.Fatal(err)
		}
		err = tmpl.Execute(&b, t)
		if err != nil {
			log.Fatal(err)
		}
		msg := b.String()
		sha256msg := sha256sum(msg)
		mapMsg[sha256msg[:8]] = msg
	}
	return mapMsg
}

func isHex(s string) bool {
	r := regexp.MustCompile(`^[a-f0-9]+$`).MatchString
	if !r(s) {
		log.Println("Error parsing sha256 sum.")
		return false
	}
	return true
}

func randomInit() {
	randbytes := make([]byte, 8)
	_, err := cryptorand.Read(randbytes)
	if err != nil {
		log.Fatal(err)
	}
	randint64 := int64(binary.BigEndian.Uint64(randbytes))
	rand.Seed(randint64)
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.LoadHTMLGlob("./static/index.tmpl")
	msgs := randMessage()
	keys := []string{}
	for k := range msgs {
		keys = append(keys, k)
	}

	r.GET("/", func(c *gin.Context) {
		key := keys[rand.Intn(len(keys))]
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"message":   msgs[key],
			"permalink": key,
		})
	})

	r.GET("/p/:sha", func(c *gin.Context) {
		input := c.Param("sha")
		if isHex(input) && len(input) == 8 {
			if k, ok := msgs[input]; ok {
				c.HTML(http.StatusOK, "index.tmpl", gin.H{
					"message":   k,
					"permalink": msgs[k],
				})
			}
		} else {
			c.String(http.StatusBadRequest, "400 Bad Request")
		}
	})

	r.GET("/commit.txt", func(c *gin.Context) {
		key := keys[rand.Intn(len(keys))]
		c.String(http.StatusOK, msgs[key])
	})

	r.GET("/json", func(c *gin.Context) {
		key := keys[rand.Intn(len(keys))]
		c.JSON(http.StatusOK, gin.H{
			"message": msgs[key],
		})
	})

	r.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	r.GET("/robots.txt", func(c *gin.Context) {
		c.String(http.StatusOK, "User-agent: *\nDisallow: /")
	})

	return r
}

func main() {
	randomInit()
	r := setupRouter()
	domain := os.Getenv("DOMAIN")
	if os.Getenv("TLS") == "true" {
		m := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(domain),
			Cache:      autocert.DirCache("./cache"),
		}
		s := &http.Server{
			Handler: m.HTTPHandler(nil),
			Addr:    ":80",
		}
		go s.ListenAndServe()
		log.Fatal(autotls.RunWithManager(r, &m))
	} else {
		r.Run()
	}
}

func init() {
	log.Println("Starting gin-commitment release: ", SourceBranch)
	log.Println("Git SHA: ", SourceCommit)
	log.Println("Go version: ", GolangVersion)
}
