package tasks

import "context"

// IfElse runs the given function, and if it returns an error runs the next function.
func IfElse[T any](do func(state *T, ctx context.Context) error, ifError func(state *T, ctx context.Context) error) func(state *T, ctx context.Context) error {
	return func(state *T, ctx context.Context) error {
		if err := do(state, ctx); err != nil {
			ctx = context.WithValue(ctx, "error", err)
			return ifError(state, ctx)
		}
		return nil
	}
}

// Serial runs the given functions in order, passing the state to each function.
func Serial[T any](tasks ...func(state *T, ctx context.Context) error) func(state *T, ctx context.Context) error {
	return func(state *T, ctx context.Context) error {
		for _, task := range tasks {
			if err := task(state, ctx); err != nil {
				return err
			}
		}
		return nil
	}
}
