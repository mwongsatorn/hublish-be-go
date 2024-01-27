package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
	"time"
)

func GenerateSlug(title string) string {

	firstPartSlug := strings.Join(strings.Fields(title), "-")
	firstPartSlug = url.PathEscape(firstPartSlug)
	firstPartSlug = strings.ToLower(firstPartSlug)

	secondPartSlug := fmt.Sprintf("%s%d", firstPartSlug, time.Now().UnixNano()/int64(time.Millisecond))
	hasher := sha256.New()
	hasher.Write([]byte(secondPartSlug))
	hashedSecondPart := hex.EncodeToString(hasher.Sum(nil))[:8]

	finalSlug := fmt.Sprintf("%s-%s", firstPartSlug, hashedSecondPart)

	return finalSlug
}
