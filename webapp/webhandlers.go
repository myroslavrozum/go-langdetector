package webapp

import (
	"go-langdetector/algos"
	"go-langdetector/constants"
	"go-langdetector/crawler"
	"go-langdetector/trainer"
	"log"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func serveIndexPage(c *gin.Context) {
	type Webdata struct {
		Title              string
		SupportedLanguages map[string]string
	}

	webdata := Webdata{
		Title:              "language Detectur",
		SupportedLanguages: make(map[string]string),
	}

	for shortName, v := range constants.UrlDictionary {
		fullName := v[0]
		webdata.SupportedLanguages[shortName] = fullName
	}

	c.HTML(http.StatusOK, "index.html.tpl", webdata)
}

func Detect(trigrammes map[string]map[string]float64) gin.HandlerFunc {
	return func(c *gin.Context) {
		contentToCheck := c.PostForm("content")

		data := ""

		if strings.HasPrefix(contentToCheck, `http://`) || strings.HasPrefix(contentToCheck, `https://`) {
			var err error
			data, err = crawler.GetTextFromURL(contentToCheck)
			if err != nil {
				log.Printf("Error getting data from URL %s: %v", contentToCheck, err)
				log.Println("Setting `data` to: ", contentToCheck)
				data = contentToCheck
			}
		} else {
			data = contentToCheck
		}

		trigrammes2investigate := trainer.ExtractTrigrammesFromText(data)
		log.Printf("Got %d trigrammes for: \"%s\"", len(trigrammes2investigate), contentToCheck)

		distances := make(map[string]float64)
		var minD float64 = math.MaxFloat64
		var minLangFull string
		var minLang string
		for lang, v := range constants.UrlDictionary {
			d := algos.CalculateColsineDistances(trigrammes[lang], trigrammes2investigate)
			distances[lang] = d
			log.Printf("Calculated distance to %s is %f", v[0], d)
			if d < minD {
				minD = d
				minLangFull = v[0]
				minLang = lang
			}
		}
		log.Printf("Minimum distance (%f) is to %s", minD, minLangFull)

		c.JSON(http.StatusOK, gin.H{
			"status":      "posted",
			"minD":        minD,
			"minLangFull": minLangFull,
			"minLang":     minLang,
		})
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// CheckOrigin allows cross-origin requests. Adjust for production security.
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	conn     *websocket.Conn
	send     chan []byte
	mu       sync.Mutex // Ensures closing happens only once
	isClosed bool
}

func (c *Client) writePump() {
	defer func() {
		c.mu.Lock()
		if !c.isClosed {
			c.conn.Close()
			c.isClosed = true
		}
		c.mu.Unlock()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				// The channel was closed by the hub.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// Write the message to the websocket
			c.conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func wsServe(logger chan string) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("Upgrade error:", err)
			return
		}

		client := &Client{
			conn: conn,
			send: make(chan []byte, 256), // Buffered channel for async queuing
		}

		// Start the async writer in a separate goroutine
		go client.writePump()

		// --- Example of generating and sending async messages ---
		// In a real application, messages would come from a central broadcast hub.
		go func() {
			for {
				time.Sleep(2 * time.Second)
				l := <-logger
				select {
				case client.send <- []byte([]byte(l)):
				default:
					// If the send channel is full or closed, break out of the loop
					log.Println("Failed to send message, client disconnected?")
					return
				}
			}
		}()
	}
}
