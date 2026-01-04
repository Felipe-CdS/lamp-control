package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand/v2"
	"net/http"
	"os"
)

func main() {
	b := make([]byte, 16)

	binary.BigEndian.PutUint64(b, uint64(rand.Uint()))
	binary.BigEndian.PutUint64(b[8:], uint64(rand.Uint()))

	resp, err := http.Post(
		"http://192.168.1.2/app/handshake1",
		"application/octet-stream",
		bytes.NewReader(b),
	)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(resp)
}
