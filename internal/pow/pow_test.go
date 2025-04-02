package pow

import (
	"strconv"
	"testing"
	"time"
)

func TestPowVerification(t *testing.T) {
	challenge := GenerateChallenge("test", 3, 60)

	var nonce string
	found := false
	for i := 0; i < 10000000; i++ {
		nonce = strconv.Itoa(i)
		if challenge.Verify(nonce) {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("Failed to find a valid nonce for the test")
	}

	if !challenge.Verify(nonce) {
		t.Errorf("Nonce %s failed verification.", nonce)
	}
}

func TestChallengeTTL(t *testing.T) {
	challenge := GenerateChallenge("test", 2, 1)
	time.Sleep(2 * time.Second)
	if time.Since(challenge.Timestamp) < time.Duration(challenge.TTL)*time.Second {
		t.Error("Failed to find a valid nonce for the test")
	}
}
