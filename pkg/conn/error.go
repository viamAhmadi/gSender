package conn

type Error struct {
	Msg         string
	Destination []byte
}