package utils

import (
	"fmt"
	"time"
)

func Map[I any, O any](ins []I, fn func(in I, idx int) O) []O {
	outs := make([]O, len(ins))
	for idx, in := range ins {
		outs[idx] = fn(in, idx)
	}
	return outs
}

func BuildTimeDiffDisplay(t time.Time, prefix string) string {
	prefix += " "
	diff := time.Since(t)
	if diff < time.Minute {
		return BOLD_MAGENTA + prefix + "just now" + RESET
	} else if diff < time.Hour {
		return BOLD_MAGENTA + prefix + fmt.Sprintf("%dm ago", int(diff.Minutes())) + RESET
	} else if diff < time.Hour*24 {
		return BOLD_MAGENTA + prefix + fmt.Sprintf("%dh ago", int(diff.Hours())) + RESET
	} else if diff < time.Hour*72 {
		return MAGENTA + prefix + fmt.Sprintf("%dh ago", int(diff.Hours()/24)) + RESET
	} else {
		return MAGENTA + prefix + fmt.Sprintf("%dd ago", int(diff.Hours()/24)) + RESET
	}
}

const (
	YELLOW       = "\033[33;1m"
	CYAN         = "\033[36;1m"
	GREEN        = "\033[32;1m"
	RED          = "\033[31;1m"
	MAGENTA      = "\033[35m"
	BOLD_MAGENTA = "\033[35;1m"
	RESET        = "\033[0m"
)
