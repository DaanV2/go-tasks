package tasks

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IfElse(t *testing.T) {
	t.Run("If first is successful, second will not be called", func(t *testing.T) {
		var called bool
		err := IfElse(
			func(state *int, ctx context.Context) error {
				return nil
			},
			func(state *int, ctx context.Context) error {
				called = true
				return nil
			},
		)(nil, context.Background())
		assert.NoError(t, err)

		if called {
			t.Error("Second function was called")
		}
	})

	t.Run("If first is not successful, second will be called", func(t *testing.T) {
		var called bool

		err := IfElse(
			func(state *int, ctx context.Context) error {
				return Cancel
			},
			func(state *int, ctx context.Context) error {
				called = true
				return nil
			},
		)(nil, context.Background())
		assert.NoError(t, err)

		if !called {
			t.Error("Second function was not called")
		}
	})
}
