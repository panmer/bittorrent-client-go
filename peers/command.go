package peers

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"encoding/binary"

	"torrent-client-go/torrent"
	"torrent-client-go/info"

	"github.com/jackpal/bencode-go"
)

func GetPeersFromTracker(trackerURL string, infoHash [20]byte,
							metadata *torrent.Torrent) ([]string, error) {
	hashStr := string(infoHash[:])
	baseURL, _ := url.Parse(trackerURL)

	params := url.Values{}
	params.Add("info_hash", hashStr)

	params.Add("port", "6881")
	params.Add("uploaded", "48")
	params.Add("downloaded", "48")
	params.Add("peer_id", "abcedfgh123456")

	if metadata != nil {
		params.Add("left", strconv.Itoa(metadata.Info.Length))
	} else {
		params.Add("left", "9999")
	}
	params.Add("compact", "1")

	baseURL.RawQuery = params.Encode()

	resp, err := http.Get(baseURL.String())
	if err != nil {
		return nil, fmt.Errorf("Error fetching trackerURL : %v", err)
	}
	defer resp.Body.Close()

	trackerResp := torrent.TrackerResponse{}
	if err := bencode.Unmarshal(resp.Body, &trackerResp); err != nil {
		return nil, fmt.Errorf("Error Unmarshal Tracker Response: %v", err)
	}

	var peers []string = {}
	peerData := []byte(trackerResp.Peers)

	for i := 0; i < len(peerData); i += 6 {
		ip := fmt.Sprintf("%d.%d.%d.%d", peerData[i], peerData[i+1], peerData[i+2], peerData[i+3])

		port := binary.BigEndian.Uint16(peerData[i+4 : i+4+2])
		peers = append(peers, fmt.Sprintf("%s:%d", ip, port))
	}

	return peers, nil
}


func Command(bencodedValue string) []string {
	metadata, err := info.LoadTorrentFile(bencodedValue)
	if err != nil {
		fmt.Println(err)
		return []string {}
	}

	infoHash, err := info.GenHash(metadata.Info)
	if err != nil {
		fmt.Println(err)
		return []string {}
	}

	peers, err := GetPeersFromTracker(metadata.Announce, infoHash, metadata)
	if err != nil {
		fmt.Println(err)
		return []string {}
	}

	fmt.Println("All Peers:", peers)
	return peers
}