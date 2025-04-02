package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/AnatolyPoluyaktov/go-wisdom-pow/config"
	"github.com/AnatolyPoluyaktov/go-wisdom-pow/internal/pow"
	"github.com/AnatolyPoluyaktov/go-wisdom-pow/internal/quotes"
	"github.com/google/uuid"
	"log"
	"net"
	"time"
)

func handleConnection(conn net.Conn, cfg config.ServerConfig) {
	defer conn.Close()

	clientID := uuid.New().String()
	challenge := pow.GenerateChallenge(clientID, cfg.POW.Difficulty, cfg.POW.TTL)

	challengeJSON, err := json.Marshal(challenge)
	if err != nil {
		log.Printf("Error marshaling challenge: %v", err)
		return
	}
	_, err = conn.Write(append(challengeJSON, '\n'))
	if err != nil {
		log.Printf("Error sending challenge: %v", err)
		return
	}

	scanner := bufio.NewScanner(conn)
	attempt := 0
	valid := false

	for scanner.Scan() && attempt < cfg.POW.MaxAttempts {
		attempt++

		if time.Since(challenge.Timestamp) > time.Duration(challenge.TTL)*time.Second {
			log.Printf("Challenge for client %s has expired", clientID)
			conn.Write([]byte("Challenge expired\n"))
			return
		}

		var clientMsg struct {
			Nonce string `json:"nonce"`
		}
		err = json.Unmarshal(scanner.Bytes(), &clientMsg)
		if err != nil {
			log.Printf("Error decoding client message: %v", err)
			conn.Write([]byte("Invalid message format\n"))
			continue
		}
		if !challenge.Verify(clientMsg.Nonce) {
			log.Printf("Attempt %d: invalid solution from client %s", attempt, clientID)
			conn.Write([]byte(fmt.Sprintf("Invalid solution, attempt %d/%d\n", attempt, cfg.POW.MaxAttempts)))
			continue
		}
		valid = true
		break
	}

	if !valid {
		log.Printf("Client %s has exhausted all attempts", clientID)
		conn.Write([]byte("Maximum number of attempts exceeded\n"))
		return
	}

	quote := quotes.GetRandomQuote()
	_, err = conn.Write([]byte(fmt.Sprintf("Quote of the day: %s\n", quote)))
	if err != nil {
		log.Printf("Error sending quote to client %s: %v", clientID, err)
		return
	}
	log.Printf("Client %s successfully received a quote", clientID)
}

func main() {
	cfg, err := config.LoadConfig("./config")
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	listenAddr := fmt.Sprintf("0.0.0.0:%d", cfg.Server.Port)
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	log.Printf("Server is running on %s", listenAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Connection error: %v", err)
			continue
		}
		go handleConnection(conn, cfg.Server)
	}
}
