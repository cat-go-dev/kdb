package compute

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandType(t *testing.T) {
	tests := []struct {
		name      string
		cType     CommandType
		extpected string
	}{
		{name: "is get type", cType: Get, extpected: "GET"},
		{name: "is set type", cType: Set, extpected: "SET"},
		{name: "is del type", cType: Del, extpected: "DEL"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.extpected {
			case "GET":
				assert.True(t, tt.cType.IsGet())
			case "SET":
				assert.True(t, tt.cType.IsSet())
			case "DEL":
				assert.True(t, tt.cType.IsDel())
			}
		})
	}
}
