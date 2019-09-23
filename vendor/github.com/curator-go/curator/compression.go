package curator

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"

	"github.com/bkaradzic/go-lz4"
)

var (
	CompressionProviders = map[string]CompressionProvider{
		"gzip": NewGzipCompressionProvider(),
		"lz4":  NewLZ4CompressionProvider(),
	}
)

type CompressionProvider interface {
	Compress(path string, data []byte) ([]byte, error)

	Decompress(path string, compressedData []byte) ([]byte, error)
}

type GzipCompressionProvider struct {
	level int
}

func NewGzipCompressionProvider() *GzipCompressionProvider {
	return NewGzipCompressionProviderWithLevel(gzip.DefaultCompression)
}

func NewGzipCompressionProviderWithLevel(level int) *GzipCompressionProvider {
	return &GzipCompressionProvider{level: level}
}

func (c *GzipCompressionProvider) Compress(path string, data []byte) ([]byte, error) {
	var buf bytes.Buffer

	if writer, err := gzip.NewWriterLevel(&buf, c.level); err != nil {
		return nil, err
	} else if _, err := writer.Write(data); err != nil {
		return nil, err
	} else if err := writer.Close(); err != nil {
		return nil, err
	} else {
		return buf.Bytes(), nil
	}
}

func (c *GzipCompressionProvider) Decompress(path string, compressedData []byte) ([]byte, error) {
	buf := bytes.NewBuffer(compressedData)

	if reader, err := gzip.NewReader(buf); err != nil {
		return nil, err
	} else if data, err := ioutil.ReadAll(reader); err != nil {
		return nil, err
	} else {
		return data, err
	}
}

type LZ4CompressionProvider struct{}

func NewLZ4CompressionProvider() *LZ4CompressionProvider {
	return &LZ4CompressionProvider{}
}

func (c *LZ4CompressionProvider) Compress(path string, data []byte) ([]byte, error) {
	return lz4.Encode(nil, data)
}

func (c *LZ4CompressionProvider) Decompress(path string, compressedData []byte) ([]byte, error) {
	return lz4.Decode(nil, compressedData)
}
