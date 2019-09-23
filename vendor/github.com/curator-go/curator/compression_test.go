package curator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGzipCompressionProvider(t *testing.T) {
	p := NewGzipCompressionProvider()

	assert.NotNil(t, p)

	data, err := p.Compress("/node", []byte("data"))

	assert.Equal(t, 28, len(data))
	assert.NoError(t, err)

	data, err = p.Decompress("/node", data)

	assert.Equal(t, "data", string(data))
	assert.NoError(t, err)
}

func TestLZ4CompressionProvider(t *testing.T) {
	p := NewLZ4CompressionProvider()

	assert.NotNil(t, p)

	data, err := p.Compress("/node", []byte("data"))

	assert.Equal(t, 9, len(data))
	assert.NoError(t, err)

	data, err = p.Decompress("/node", data)

	assert.Equal(t, "data", string(data))
	assert.NoError(t, err)
}
