package main

import (
	"bytes"
	"compress/gzip"
	"io"
	"log"
)

func UGZipBytes(data []byte) []byte {
	var out bytes.Buffer
	var in bytes.Buffer

	in.Write(data)

	r, err := gzip.NewReader(&in)
	if err != nil {
		log.Println(err)
	}
	defer r.Close()

	_, err = io.Copy(&out, r)
	if err != nil {
		log.Println(err)
	}

	return out.Bytes()
}
