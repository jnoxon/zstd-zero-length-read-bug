package poc

import (
	"bytes"
	"io"
	"testing"

	"github.com/klauspost/compress/zstd"
)

// CursedReader wraps a reader and returns zero bytes every other read.
// This is used to test the ability of the consumer to handle empty reads without EOF,
// which can happen when reading from a network connection.
type CursedReader struct {
	io.Reader
	numReads int
}

func (r *CursedReader) Read(p []byte) (n int, err error) {
	r.numReads++
	if r.numReads%2 == 0 {
		return 0, nil
	}

	return r.Reader.Read(p)
}

func NewCursedReader(r io.Reader) *CursedReader { return &CursedReader{Reader: r} }

func zeroLengthReadTest() error {
	// create a buffer with some data; the data doesn't matter
	srcBuf := bytes.NewBuffer(make([]byte, 1024))
	dstBuf := &bytes.Buffer{}

	// compress the data
	enc, err := zstd.NewWriter(dstBuf, zstd.WithEncoderLevel(zstd.SpeedFastest))
	if err != nil {
		return err
	}
	_, err = enc.Write(srcBuf.Bytes())
	if err != nil {
		return err
	}
	err = enc.Close()
	if err != nil {
		return err
	}

	// now let's read it back using a cursed reader
	dec, err := zstd.NewReader(NewCursedReader(dstBuf))
	if err != nil {
		return err
	}
	_, err = io.Copy(io.Discard, dec)
	if err != nil {
		return err
	}

	return nil
}

func TestZeroLengthReads(t *testing.T) {
	if err := zeroLengthReadTest(); err != nil {
		t.Error(err)
	}
}
