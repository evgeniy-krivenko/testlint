package activities

import (
	"context"
	"crypto/rand"
	"math/big"
	"strings"
)

func GenerateCode(_ context.Context, codeSize int) (string, error) {
	var buf strings.Builder
	buf.Grow(codeSize)
	for i := 0; i < codeSize; i++ {
		res, err := rand.Int(rand.Reader, big.NewInt(9))
		if err != nil {
			return "", err
		}
		buf.WriteString(res.String())
	}
	return buf.String(), nil
}
