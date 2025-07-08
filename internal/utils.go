package internal

import (
	"strconv"

	"github.com/codecrafters-io/tester-utils/random"
)

var FRUITS = []string{"apple", "banana", "blackberry", "blueberry", "cherry", "grape", "lemon", "mango", "orange", "pear", "pineapple", "plum", "raspberry", "strawberry", "watermelon"}
var VEGETABLES = []string{"carrot", "onion", "potato", "tomato", "broccoli", "cauliflower", "cabbage", "lettuce", "spinach", "asparagus", "peas", "corn", "zucchini", "pumpkin"}

func randomFilePrefix() string {
	return strconv.Itoa(random.RandomInt(1000, 10000))
}
