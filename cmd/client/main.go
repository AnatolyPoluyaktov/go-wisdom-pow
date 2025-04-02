package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/AnatolyPoluyaktov/go-wisdom-pow/config"
	"github.com/AnatolyPoluyaktov/go-wisdom-pow/internal/pow"
	"log"
	"net"
	"strconv"
	"time"
)

func solveChallenge(challenge pow.Challenge) (string, error) {
	nonce := 0
	for {
		nonceStr := strconv.Itoa(nonce)
		hash := sha256.Sum256([]byte(challenge.Data + nonceStr))
		hashStr := hex.EncodeToString(hash[:])
		prefix := ""
		for i := 0; i < challenge.Difficulty; i++ {
			prefix += "0"
		}
		if hashStr[:challenge.Difficulty] == prefix {
			return nonceStr, nil
		}
		nonce++
		if nonce%1000000 == 0 {
			log.Printf("Searching for nonce, current value: %d", nonce)
		}
	}
}

func main() {
	cfg, err := config.LoadConfig("./config")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	serverAddr := cfg.Client.ServerAddr
	var conn net.Conn
	maxRetryCount := 5

	for i := 0; i < maxRetryCount; i++ {
		conn, err = net.Dial("tcp", serverAddr)
		if err == nil {
			break
		}
		log.Printf("Connection attempt %d/%d failed: %v", i+1, maxRetryCount, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("Failed to connect to the server: %v", err)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	challengeLine, err := reader.ReadBytes('\n')
	if err != nil {
		log.Fatalf("Error reading challenge: %v", err)
	}
	var challenge pow.Challenge
	err = json.Unmarshal(challengeLine, &challenge)
	if err != nil {
		log.Fatalf("Error decoding challenge: %v", err)
	}

	log.Printf("Challenge received: %+vv", challenge)

	nonce, err := solveChallenge(challenge)
	if err != nil {
		log.Fatalf("Error solving challenge: %v", err)
	}

	solution := map[string]string{"nonce": nonce}
	solutionJSON, _ := json.Marshal(solution)
	_, err = conn.Write(append(solutionJSON, '\n'))
	if err != nil {
		log.Fatalf("Error sending solution: %v", err)
	}

	response, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error reading server response: %v", err)
	}
	fmt.Printf("Server response: %s", response)
}
