package decode

import (
	"fmt"
	"strings"
	"encoding/json"
	"github.com/jackpal/bencode-go"
)

func decodeBencode(bencodedString string) (interface {}, error) {
	fmt.Println(bencodedString)
	data, err := bencode.Decode(strings.NewReader(bencodedString))

	return data, err
}

func Command(bencodedValue string) {
	// bencodedString := "Encode Command"
	decoded, err := decodeBencode(bencodedValue)
	if err != nil {
		fmt.Println("Error decoding bencoded:", err);
		return
	}

	serialized, err := json.MarshalIndent(decoded, "", "  ")
	if err != nil {
		fmt.Println("Error serialize to JSON:", err);
		return
	}

	fmt.Println(string(serialized))
}