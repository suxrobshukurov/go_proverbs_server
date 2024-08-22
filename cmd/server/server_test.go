package main

import (
	"bufio"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func Test_handleConn(t *testing.T) {
	proverbs := []string{
		"Concurrency is not parallelism.",
		"Don't communicate by sharing memory, share memory by communicating.",
		"Channels orchestrate; mutexes serialize.",
	}
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	go handleConn(server, proverbs)

	reader := bufio.NewReader(client)
	timeout := 5 * time.Second

	for i := 0; i < len(proverbs); i++ {
		client.SetReadDeadline(time.Now().Add(timeout))
		line, err := reader.ReadString('\n')
		if err != nil {
			t.Fatalf("Failed to read from connection: %v", err)
		}
		line = strings.TrimSpace(line)

		if !contains(proverbs, line) {
			t.Errorf("Unexpected proverb: %s", line)
		}
	}
}

func contains(proverbs []string, proverb string) bool {
	for _, p := range proverbs {
		if p == proverb {
			return true
		}
	}

	return false
}

func Test_fetchProverbs(t *testing.T) {
	mockData := `
		<h3><a href="#">Proverb 1</a></h3>
		<h3><a href="#">Proverb 2</a></h3>
		<h3><a href="#">Proverb 3</a></h3>
	`
	// Используем httptest для создания тестового HTTP-сервера
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, mockData)
	}))
	defer server.Close()

	// Вызываем fetchProverbs с URL тестового сервера
	proverbs, err := fetchProverbs(server.URL)
	if err != nil {
		t.Fatalf("Failed to fetch proverbs: %v", err)
	}
	expected := []string{"Proverb 1", "Proverb 2", "Proverb 3"}

	if len(proverbs) != len(expected) {
		t.Fatalf("Expected %d proverbs, got %d", len(expected), len(proverbs))
	}

	for i, proverb := range proverbs {
		if proverb != expected[i] {
			t.Errorf("Expected '%s', got '%s'", expected[i], proverb)
		}
	}

}
