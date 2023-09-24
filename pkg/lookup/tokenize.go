package lookup

import "github.com/tiktoken-go/tokenizer/codec"

var tokenizer = codec.NewCl100kBase()

func Tokenize(s string) ([]string, error) {
	_, strings, err := tokenizer.Encode(s)
	return strings, err
}
