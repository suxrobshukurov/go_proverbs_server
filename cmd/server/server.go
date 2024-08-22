package main

import (
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"regexp"
	"time"
)

const (
	addr           = "0.0.0.0:12345"
	proto          = "tcp4"
	proverbsURL    = "https://go-proverbs.github.io/"
	tickerInterval = 3 * time.Second
)

func handleConn(conn net.Conn, proverbs []string) {
	defer conn.Close()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	t := time.NewTicker(tickerInterval)
	defer t.Stop()

	for range t.C {
		conn.Write([]byte(proverbs[r.Intn(len(proverbs))] + "\n"))
	}
}

func fetchProverbs(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`<h3><a [^>]*>([^<]+)</a></h3>`)
	matches := re.FindAllStringSubmatch(string(body), -1)

	if len(matches) == 0 {
		return nil, err
	}

	proverbs := make([]string, 0, len(matches))
	for _, match := range matches {
		if len(match) > 1 {
			proverbs = append(proverbs, match[1])
		}
	}

	return proverbs, nil
}

func main() {

	proverbs, err := fetchProverbs(proverbsURL)
	if err != nil || len(proverbs) == 0 {
		log.Fatalf("Failed to fetch proverbs: %v", err)
	} else {
		log.Printf("Fetched %d proverbs", len(proverbs))
	}

	listener, err := net.Listen(proto, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	log.Printf("Server listening on %s", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Failed to accept connection:", err)
			continue
		}
		go handleConn(conn, proverbs)
	}
}
