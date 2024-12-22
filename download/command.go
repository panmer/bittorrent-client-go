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


func DownloadPiece(tcpConnection *net.TCPConn, pidx int, totalBlocks int,
					pieceLength int, pieceReceivedIndex int, pieceData []byte,
					downloadPath string, Info *torrent.InfoData) []byte {
		metadata, err := infoCommand.LoadTorrentFile(bencodedValue)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		peers := peers.Command(bencodedValue)

		pieceData := make([]byte, 0)
		pidx, _ := strconv.Atoi(pieceIndex)

		tcpConnection := tcp.ConnectTCP(bencodedValue, peers[pidx%len(peers)])
	
		pieceLength := metadata.Info.PieceLength

		for {
			messageLength := make([]byte, 4)

			_, err := io.ReadFull(tcpConnection, messageLength)
			if err != nil {
				fmt.Println("Error reading messageLength", err)
				retry(pidx)
				return nil
			}

			length := binary.BigEndian.Uint32(messageLength)
			if length == 0 {
				fmt.Println("Keep alive message eceived")
				continue
			}

			messageID := make([]byte, 1)
			_, err = io.ReadFull(tcpConnection, messageID)
			if err != nil {
				fmt.Println("Error read messageID", err)
				return nil
			}
			id := uint8(messageID[0])

			switch id {
			case 5:
				fmt.Println("Received bitfield message")
				payload := make([]byte, length-1)
				_, err := io.ReadFull(tcpConnection, payload)
				if err != nil {
					fmt.Println("error reading bitfieldPayload", err)
					retry(pidx)
					return nil
				}
				interested := []byte{0, 0, 0, 1, 2}
				_, err = tcpConnection.Write(interested)
				if err != nil {
					fmt.Println("Error sending interested message:", err)
					retry(pidx)
					return nil
				}
	
			case 1:
				fmt.Println("Unchoke message received")
				for i := 0; i < totalBlocks; i++ {
					blockSize := 16 * 1024
					if i == totalBlocks-1 {
						blockSize = pieceLength % (16 * 1024)
					}
	
					request := make([]byte, 17)
					binary.BigEndian.PutUint32(request[0:4], 13)                  // Message length (13 bytes)
					request[4] = 6                                                // Message ID (request)
					binary.BigEndian.PutUint32(request[5:9], uint32(pidx))    // Piece index
					binary.BigEndian.PutUint32(request[9:13], uint32(i*16*1024))  // Begin offset
					binary.BigEndian.PutUint32(request[13:17], uint32(blockSize)) // Block length
	
					_, err = tcpConnection.Write(request)
					if err != nil {
						fmt.Printf("Error sending request for block %d: %v\n", i+1, err)
						retry(pidx)
						return nil
					}
				}
			case 7:
				header := make([]byte, 8)
				_, err := io.ReadFull(tcpConnection, header)
				if err != nil {
					fmt.Println("Error reading piece header:", err)
					retry(pidx)
					return nil
				}
				index := binary.BigEndian.Uint32(header[0:4])
				if int(index) != pidx {
					fmt.Printf("Wrong piece index received. Expected %d, got %d\n", pidx, index)
					retry(pidx)
					return nil
				}
	
				blockSize := 16 * 1024
				if pieceReceivedIndex == totalBlocks-1 {
					blockSize = pieceLength % (16 * 1024)
				}
	
				dataBuff := make([]byte, blockSize)
				_, err = io.ReadFull(tcpConnection, dataBuff)
				if err != nil {
					fmt.Printf("Error reading piece data (block %d): %v\n", pieceReceivedIndex, err)
					retry(pidx)
					return nil
				}
	
				pieceData = append(pieceData, dataBuff...)
				pieceReceivedIndex++
				fmt.Printf("Received block %d of %d (size: %d bytes)\n", pieceReceivedIndex, totalBlocks, blockSize)
	
				if pieceReceivedIndex == totalBlocks {
					receivedPieceHash := sha1.Sum(pieceData)
					expectedHash := Info.Pieces[pidx*20 : (pidx+1)*20]
	
					if bytes.Equal(receivedPieceHash[:], []byte(expectedHash)) {
						fmt.Println("Piece hash verified successfully")
						if downloadPath != "" {
							err := queueSavePieces(pieceData, downloadPath)

							if err != nil {
								fmt.Println("Error saving piece to file:", err)
								return nil
							}
							fmt.Println("Piece saved successfully")
						}
						return pieceData
					} else {
						fmt.Println("Piece hash verification failed")
						return nil
					}
				}
			}
		}
	}
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