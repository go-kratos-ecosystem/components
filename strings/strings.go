package strings

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"math/rand"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

const (
	// randomLetters is the letters used in Random.
	randomLetters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// uuidPattern is the pattern used in IsUuid.
	uuidPattern = "^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"
)

// Is returns true if the value matches the pattern.
// The pattern can contain the wildcard character *.
//
// Example:
//
//	Is("*.example.com", "www.example.com") // true
//	Is("*.example.com", "example.com") // false
func Is(pattern, value string) bool {
	if pattern == value {
		return true
	}

	pattern = strings.ReplaceAll(pattern, "*", ".*")

	match, err := regexp.Match(pattern, []byte(value))
	if err != nil {
		return false
	}

	return match
}

// InSlice checks if a string is in a string slice.
//
// Example:
//
//	InSlice([]string{"1", "2"}, "1") // true
func InSlice(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}

	return false
}

// Md5 returns the md5 hash of a string.
//
// Example:
//
//	Md5("abc") // 900150983cd24fb0d6963f7d28e17f72
func Md5(s string) string {
	sm := md5.Sum([]byte(s))

	return fmt.Sprintf("%x", sm)
}

// Sha1 returns the sha1 hash of a string.
//
// Example:
//
//	Sha1("abc") // a9993e364706816aba3e25717850c26c9cd0d89d
func Sha1(s string) string {
	sm := sha1.Sum([]byte(s))

	return fmt.Sprintf("%x", sm)
}

// Reverse returns a string with its characters in reverse order.
//
// Example:
//
//	Reverse("abc") // "cba"
func Reverse(s string) string {
	var reversed string

	for _, v := range s {
		reversed = string(v) + reversed
	}

	return reversed
}

// Replace replaces all occurrences of one substring with another.
//
// Example:
//
//	Replace("aabbcc", "a", "b") // "bbbbcc"
func Replace(s, from, to string) string {
	return strings.NewReplacer(from, to).Replace(s)
}

// Shuffle returns a string with its characters in random order.
//
// Example:
//
//	Shuffle("abc") // "bca"
func Shuffle(s string) string {
	ss := strings.Split(s, "")

	rand.Shuffle(len(ss), func(i, j int) {
		ss[i], ss[j] = ss[j], ss[i]
	})

	return strings.Join(ss, "")
}

// Random returns a random string with the specified length.
//
// Example:
//
//	Random(10) // "qujrlkhyqr"
func Random(length int) string {
	letters := []rune(randomLetters)
	lettersLength := len(letters)

	b := make([]rune, length)

	for i := range b {
		b[i] = letters[rand.Intn(lettersLength)]
	}

	return string(b)
}

// Len returns the length of a string.
// Support chinese characters.
//
// Example:
//
//	Len("abc") // 3
//	Len("张三李四") // 4
func Len(s string) int {
	return len([]rune(s))
}

// IsUUID returns true if the string is a valid UUID.
func IsUUID(str string) bool {
	match, err := regexp.MatchString(uuidPattern, str)

	return err == nil && match
}

// UUID generate uuid string
func UUID() string {
	return uuid.New().String()
}
