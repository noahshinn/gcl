package utils

func Map[I any, O any](ins []I, fn func(in I, idx int) O) []O {
	outs := make([]O, len(ins))
	for idx, in := range ins {
		outs[idx] = fn(in, idx)
	}
	return outs
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
