package torrent

type InfoData struct {
	Length int `bencode:"length"`
	Name string `bencode:"name"`
	PieceLength int `bencode:"piece length"`
	Pieces string `bencode:"pieces"`
}
type Torrent struct {
	Announce string `bencode:"announce"`
	Info InfoData `bencode:"info"`
}

type TrackerResponse struct {
	Peers string `bencode:"peers"`
	Interval int`bencode:"interval"`
}

type TCPRequest struct {
	Length   uint8
	Protocol [19]byte
	Reserve  [8]byte
	InfoHash [20]byte
	PeerID   [20]byte
}