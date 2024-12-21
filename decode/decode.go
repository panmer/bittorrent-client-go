package decode

import (
	"fmt"
	"strings"
	"encoding/json"
	"github.com/jackpal/bencode-go"
)

func decodeBencode(bencodedString string) (interface {}, error) {
	data, err := bencode.Decode(strings.NewReader(bencodedString))
	return data, err
}

func Command(bencodedValue string) {
	decoded, err := decodeBencode(bencodedValue)

	if err != nil {
		fmt.Println("Error decoding bencoded value:", err)
		return
	}

	fmt.Println(string(jsonOutput))
}