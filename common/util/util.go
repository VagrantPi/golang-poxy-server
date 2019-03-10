package util

import "math/rand"

// GenRandString - get rend string
func GenRandString(times int) (code string) {
	for i := 0; i < times; i++ {
		c := rand.Int() % 62
		if c >= 10 && c < 36 {
			c += 55
		} else if c >= 36 && c < 62 {
			c += 61
		} else {
			c += 48
		}
		code += string(rune(c))
	}
	return
}
