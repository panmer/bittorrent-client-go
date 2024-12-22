package info

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"torrent-client-go/torrent"
	"github.com/jackpal/bencode-go"
)
func LoadTorrentFile(torrentFilePath string) (*torrent.Torrent, error) {
	torrentFile, err := os.Open(torrentFilePath)

	if err != nil {
		return nil, fmt.Errorf("Error open file %s: %v", torrentFilePath, err)
	}
	defer torrentFile.Close()

	torrentData, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("Error read file %s: %v", torrentFilePath, err)
	}

	buffer := bytes.NewReader(torrentData)
	metadata := torrent.Torrent{}

	if err := bencode.Unmarshal(buffer, &metadata); err != nil {
		return nil, fmt.Errorf("Error unmarshal torrent data: %v", err)
	}

	return &metadata, nil
}

func GenHash(torrentInfo torrent.InfoData) ([20]byte, error) {
	var infoBuffer bytes.Buffer

	err := bencode.Marshal(&infoBuffer, torrentInfo)

	if err != nil {
		return [20]byte{}, fmt.Errorf("Error encoding: %v", err)
	}
	hash := sha1.Sum(infoBuffer.Bytes())

	return hash, nil
}