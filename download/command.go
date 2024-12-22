package download

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"torrent-client-go/torrent"
	"torrent-client-go/info"
)


func queueSavePieces(pieceData []byte, downloadPath string) error {
	file, err := os.OpenFile(downloadPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	_, err = file.Write(pieceData)
	if err != nil {
		return fmt.Errorf("error writing piece data to file: %v", err)
	}

	return nil
}
func queueAddPieces(numPieces int) {
	for i := 0; i < numPieces; i++ {
		queue.Push(i)
	}
}


func DownloadPiece(bencodedValue string, downloadPath string, pieceIndex string) []byte {
}

func DownloadFile(bencodedValue string, downloadPath string) {
	metadata, err := info.LoadTorrentFile(bencodedValue)

	if err != nil {
		fmt.Println("Error load .torrent file", bencodedValue)
		return
	}

	file := make([]byte, 0)
	numPieces := len(metadata.Info.Pieces) / 20

	queueAddPieces(numPieces)

	for !queue.Empty() {
		pidx := queue.Front(); queue.Pop()

		pieceData := DownloadPiece(bencodedValue, "", strconv.Itoa(pidx))
		file = append(file, pieceData...)
	}

	err = queueSavePieces(file, downloadPath)
	if err != nil {
		fmt.Println("Error save peices of ", downloadPath)
		return
	}

	fmt.Println("File Saved successfully")

}