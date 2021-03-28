package helpers

import (
	"fmt"
)

func EmojiToID(name, id string) string {
	return fmt.Sprintf(":%s:%s", name, id)
}
