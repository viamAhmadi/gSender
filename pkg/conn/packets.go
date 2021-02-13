package conn

import (
	"fmt"
	gUtil "github.com/viamAhmadi/gReceiver2/pkg/util"
	"strconv"
)

func ConvertToSendConn(b []byte) (*SendConn, error) {
	if cap(b) < 52 {
		return nil, ErrConvertToModel
	}
	count, err := strconv.Atoi(gUtil.RemoveAdditionalCharacters(b[48:53]))
	if err != nil {
		return nil, err
	}
	return NewSendConn(gUtil.RemoveAdditionalCharacters(b[1:28]), string(b[28:48]), count), nil
}

func SerializeSendConn(destination string, count int, id string) []byte {
	return []byte(fmt.Sprintf("c%s%s%s", gUtil.ConvertDesToBytes(destination), id, gUtil.ConvertIntToBytes(count)))
}
