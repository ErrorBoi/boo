package random

import (
	"math/rand"
	"time"
)

// Randomizer is a service for generating random strings.
type Randomizer struct {
	rnd *rand.Rand
}

// NewRandomizer initializes a new Randomizer instance.
func NewRandomizer() *Randomizer {
	// Seed the random number generator
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	return &Randomizer{rnd: rnd}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// String generates a random string of a specified length.
func (r *Randomizer) String(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[r.rnd.Intn(len(letterRunes))]
	}
	return string(b)
}
