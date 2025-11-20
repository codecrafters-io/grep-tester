package utils

import (
	"strconv"
	"strings"

	"github.com/codecrafters-io/tester-utils/random"
)

// FRUITS, VEGETABLES, and ANIMALS are used to generate test file contents
var FRUITS = []string{"apple", "banana", "blackberry", "blueberry", "cherry", "grape", "lemon", "mango", "orange", "pear", "pineapple", "plum", "raspberry", "strawberry", "watermelon"}
var VEGETABLES = []string{"carrot", "onion", "potato", "tomato", "broccoli", "cauliflower", "cabbage", "lettuce", "spinach", "asparagus", "pea", "corn", "zucchini", "pumpkin"}
var ANIMALS = []string{"cat", "dog", "elephant", "fox", "giraffe", "horse", "lion", "monkey", "panda", "rabbit", "tiger", "wolf", "zebra"}

// RandomFilePrefix returns 4 digit random prefix for test files
func RandomFilePrefix() string {
	return strconv.Itoa(random.RandomInt(1000, 10000))
}

func RandomWordsWithoutSubstrings(n int) []string {
loop:
	for {
		words := random.RandomWords(n)

		for i := 0; i < len(words); i++ {
			for j := i + 1; j < len(words); j++ {
				if strings.Contains(words[j], words[i]) || strings.Contains(words[i], words[j]) {
					continue loop
				}
			}
		}

		return words
	}
}
