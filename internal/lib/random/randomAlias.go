package random

import (
	"math/rand"
	"strings"
	"time"
)

func NewRandomAlias(randAliasLength int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	ret := strings.Builder{}
	for i := 0; i < randAliasLength; i++ {
		ret.WriteByte(byte(rnd.Intn(90-48) + 48)) // 90 символ ascii это большая Z
	}
	return ret.String()
}
