package store

import (
	"bytes"
	"encoding/binary"
	"io"
)

var lenUUID = 36

func WriteUUID(w io.Writer, s string) error {
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

func ReadUUID(r io.Reader) (string, error) {
	var b = make([]byte, lenUUID)
	_, err := r.Read(b)
	if err != nil {
		return "", err
	}
	return string(bytes.Trim(b, "\x00")), nil
}

func WriteInt64(w io.Writer, i int64) error {
	return binary.Write(w, binary.BigEndian, uint64(i))
}

func ReadInt64(r io.Reader) (int64, error) {
	var b = make([]byte, 8)
	_, err := r.Read(b)
	if err != nil {
		return 0, err
	}
	return int64(binary.BigEndian.Uint64(b)), nil
}
