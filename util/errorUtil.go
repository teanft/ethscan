package util

import "fmt"

func NewErr(Msg string, err error) error {
	return fmt.Errorf("%s: %w", Msg, err)
}
