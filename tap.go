package wiretap

import (
	"bytes"
	"io"
	"io/ioutil"
)

func tap(input io.ReadCloser) ([]byte, io.ReadCloser) {
	output, _ := ioutil.ReadAll(input)

	reader := bytes.NewReader(output)
	readerCloser := ioutil.NopCloser(reader)

	return output, readerCloser
}
