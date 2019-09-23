package curator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFixedEnsembleProvider(t *testing.T) {
	p := NewFixedEnsembleProvider("connStr")

	assert.NotNil(t, p)

	assert.NoError(t, p.Start())

	assert.Equal(t, "connStr", p.ConnectionString())

	assert.NoError(t, p.Close())
}
