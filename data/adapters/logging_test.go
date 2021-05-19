package adapters

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"poptimizer/data/domain"
	"testing"
)

func TestTypeField(t *testing.T) {
	out := zap.String("type", "RowsAppended")

	assert.Equal(t, out, TypeField(&domain.RowsAppended{}))
}
