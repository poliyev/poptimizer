package adapters

import (
	"poptimizer/data/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestTypeField(t *testing.T) {
	out := zap.String("type", "RowsAppended")

	assert.Equal(t, out, TypeField(&domain.RowsAppended{}))
}
