package test

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestRand(t *testing.T) {
	intn := rand.Intn(10)
	fmt.Println(intn)
}
