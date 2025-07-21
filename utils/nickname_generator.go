package utils

import (
	"fmt"
	"math/rand"
)

var adjactives = []string{"Blue", "Crazy", "Fast", "Silent", "Wild", "Fuzzy", "Funny", "Micro", "Green"}
var noons = []string{"Tiger", "Banana", "Fox", "Panda", "Llama", "Eagle", "Bear", "Pigeon", "ButterFly"}

func GenerateNickname() string {
	adj := adjactives[rand.Intn(len(adjactives))]
	noon := noons[rand.Intn(len(noons))]
	num := rand.Intn(100)

	return fmt.Sprintf("%s%s%d", adj, noon, num)
}
