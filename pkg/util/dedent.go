package util

import (
	"github.com/ekalinin/go-textwrap"
)

func Dedent(s string) string {
	return textwrap.Dedent(s)
}