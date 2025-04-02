package quotes

import "math/rand"

var Quotes = []string{
	"Be yourself; everyone else is already taken.",
	"Life is what happens when you're busy making other plans.",
	"Success is not final, failure is not fatal: it is the courage to continue that counts.",
	"If you want to live a happy life, tie it to a goal, not to people or things.",
}

func GetRandomQuote() string {
	return Quotes[rand.Intn(len(Quotes))]
}
