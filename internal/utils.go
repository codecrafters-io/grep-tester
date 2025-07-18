package internal

import (
	"strconv"

	"github.com/codecrafters-io/tester-utils/random"
)

// FRUITS and VEGETABLES are used to generate test file contents
var FRUITS = []string{"apple", "banana", "blackberry", "blueberry", "cherry", "grape", "lemon", "mango", "orange", "pear", "pineapple", "plum", "raspberry", "strawberry", "watermelon"}
var VEGETABLES = []string{"carrot", "onion", "potato", "tomato", "broccoli", "cauliflower", "cabbage", "lettuce", "spinach", "asparagus", "peas", "corn", "zucchini", "pumpkin"}

// randomFilePrefix returns 4 digit random prefix for test files
func randomFilePrefix() string {
	return strconv.Itoa(random.RandomInt(1000, 10000))
}
