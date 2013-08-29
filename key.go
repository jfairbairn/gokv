package gokv

import (
	"fmt"
	"regexp"
)

func keyError(k string) error {
	isValid := validKeyRegex.MatchString(k)
	if !isValid {
		return &BadKey{key: k}
	}
	return nil
}

type BadKey struct {
	key string
}

func (b *BadKey) Error() string {
	return fmt.Sprintf("Bad key %s", b.key)
}

var validKeyRegex *regexp.Regexp

func init() {
	validKeyRegex = regexp.MustCompile("^[A-z0-9\\.\\-_]+$")
}
