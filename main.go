package main

import (
	"encoding/json"
	"fmt"
	"os"
	"decode/decode"
)

func main() {
	fmt.Fprintln(os.Stderr, "Logging from program!");

	if len(os.Args) != 3 {
		fmt.Println("<command> <bencoded_value>");
		return;
	}

	command := os.Args[1];
	bencodedValue := os.Args[2];

	switch command {
		case "decode":
			decodeCommand.Command(bencodedValue)
		// case "info":
		// 	infoCommand.Command(bencodedValue)
		// case "peers":
		// 	peersCommand.Command(bencodedValue)
	}
}