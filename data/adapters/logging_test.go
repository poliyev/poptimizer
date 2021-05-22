package adapters

import (
	"poptimizer/data/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestEventField(t *testing.T) {
	out := zap.String("event", "RowsAppended(a, b)")

	assert.Equal(t, out, EventField(&domain.RowsAppended{ID: domain.NewID("a", "b")}))
}

func TestTypeField(t *testing.T) {
	out := zap.String("test", "RowsAppended")

	assert.Equal(t, out, TypeField("test", &domain.RowsAppended{}))
}
