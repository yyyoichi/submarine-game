package battle

import (
	"bytes"
	"encoding/binary"
	"io"
	"time"
)

func getGameModelKey(gameId string) ([]byte, error) {
	var buf bytes.Buffer
	_, err := buf.WriteString("G")
	if err != nil {
		return nil, err
	}

	err = writeUUID(&buf, gameId)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func getActionModelKey(gameId string, playerId string, timestamp time.Time) ([]byte, error) {
	var buf bytes.Buffer
	_, err := buf.WriteString("A")
	if err != nil {
		return nil, err
	}

	err = writeUUID(&buf, gameId)
	if err != nil {
		return nil, err
	}

	err = writeInt64(&buf, timestamp.UnixNano())
	if err != nil {
		return nil, err
	}

	err = writeUUID(&buf, playerId)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func getActionModelQueryKey(gameId string) ([]byte, error) {
	var buf bytes.Buffer
	_, err := buf.WriteString("A")
	if err != nil {
		return nil, err
	}
	err = writeUUID(&buf, gameId)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

var lenUUID = 36

func writeUUID(w io.Writer, s string) error {
	var b = make([]byte, lenUUID)
	if lenUUID < len(s) {
		// 末尾から書き込み
		_ = copy(b, []byte(s)[len(s)-lenUUID:])
	} else {
		_ = copy(b, []byte(s))
	}
	_, err := w.Write(b)
	return err
}

// func readUUID(r io.Reader) (string, error) {
// 	var b = make([]byte, lenUUID)
// 	_, err := r.Read(b)
// 	if err != nil {
// 		return "", err
// 	}
// 	return string(bytes.Trim(b, "\x00")), nil
// }

func writeInt64(w io.Writer, i int64) error {
	return binary.Write(w, binary.BigEndian, uint64(i))
}

// func readInt64(r io.Reader) (int64, error) {
// 	var b = make([]byte, 8)
// 	_, err := r.Read(b)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return int64(binary.BigEndian.Uint64(b)), nil
// }
