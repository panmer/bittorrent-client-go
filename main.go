package main

import (
	// "encoding/json"
	"fmt"
	"os"
	"torrent-client-go/decode"
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
			decode.Command(bencodedValue)
		case "info":
			fmt.Println("info command");
		// case "peers":
		// 	peersCommand.Command(bencodedValue)
	}
}