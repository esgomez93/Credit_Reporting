package utils

import (
	"fmt"
	"math/rand"
	"time"
)

// generateRandomMeme is a placeholder for meme generation logic.
func GenerateRandomMeme(query string) string {
	memes := []string{
		"One does not simply walk into Mordor.",
		"Why can't programmers tell jokes? Because we don't get them.",
		"I would explain this to you, but it's in binary.",
	}

	if query != "" {
		// Incorporate the query into the meme for a bit of relevance
		memes = append(memes, fmt.Sprintf("When you search for '%s' and find the perfect meme.", query))
	}

	rand.Seed(time.Now().UnixNano())
	return memes[rand.Intn(len(memes))]
}
