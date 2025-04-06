package ziggurat

import "math/rand/v2"

type globalRand struct{}

func (g globalRand) Uint64() uint64 {
	return rand.Uint64()
}
