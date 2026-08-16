package random

import (
	"math/rand"
	"strings"
	"time"
)

func NewRandomAlias(randAliasLength int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	ret := strings.Builder{}

	symbols := "AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz0123456789"

	for i := 0; i < randAliasLength; i++ {
		ret.WriteByte(symbols[rnd.Intn(len(symbols))])
	}
	return ret.String()
}
