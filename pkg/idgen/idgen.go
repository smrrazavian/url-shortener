// Package idgen provides functionality for generating short unique IDs.
package idgen

import (
	"strings"
	"sync"
	"time"
)

const (
	IdLength = 6
	charset  = "23456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

// Generator is used internally to create unique IDs
type Generator struct {
	lastTimeStamp int64
	counter       uint32
	mutex         sync.Mutex
}

// singleton instance used by the package
var defaultGenerator = &Generator{
	lastTimeStamp: 0,
	counter:       0,
}

// GenerateID creates a new unique ID for URL shortening
// This is the public function users of the package will call
func GenerateID() (string, error) {
	return defaultGenerator.generate()
}

// encodeToBase58 converts a number to a base58 string using the charset
func encodeToBase58(num int64) string {
	if num == 0 {
		return string(charset[0])
	}

	base := int64(len(charset))
	encoded := make([]byte, 0, IdLength)

	for num > 0 {
		remainder := num % base
		encoded = append(encoded, charset[remainder])
		num /= base
	}

	for i, j := 0, len(encoded)-1; i < j; i, j = i+1, j-1 {
		encoded[i], encoded[j] = encoded[j], encoded[i]
	}

	return string(encoded)
}

// generate is the internal implementation
func (g *Generator) generate() (string, error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	now := time.Now().UnixMicro()

	if now == g.lastTimeStamp {
		g.counter++
	} else {
		g.lastTimeStamp = now
		g.counter = 0
	}

	combined := (now << 20) | int64(g.counter)
	id := encodeToBase58(combined)

	if len(id) < IdLength {
		padding := string(charset[0])
		id = strings.Repeat(padding, IdLength-len((id))) + id
	} else if len(id) > IdLength {
		id = id[len(id)-IdLength:]
	}

	return id, nil
}
