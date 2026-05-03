package pkg

import (
	"fmt"
	"math/rand/v2"
)

func GenerateRandomString(name string) string {
	randomNum := rand.IntN(9000) + 1000
	digit4 := fmt.Sprintf("%04d", randomNum)
	return fmt.Sprintf("%s%s", name, digit4)
}
