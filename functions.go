package tasks

import "context"

// functions is a collection of functions to be run on different stages of the task.
type functions[T any] struct {
	fns []func(state *T, ctx context.Context) error
}

// Add adds a function to the collection.
func (fns *functions[T]) Add(fn func(state *T, ctx context.Context) error) {
	fns.fns = append(fns.fns, fn)
}

// Run runs the functions in the collection, and returns the first error encountered.
func (fns *functions[T]) Run(state *T, ctx context.Context) error {
	done := ctx.Done()

	for _, fn := range fns.fns {
		select {
		case <-done:
			return nil
		default:
		}

		if err := fn(state, ctx); err != nil {
			return err
		}
	}
	return nil
}

// RunParallel runs the functions in the collection in parallel, and returns the first error encountered.
func (fns *functions[T]) RunParallel(state *T, ctx context.Context) error {
	done := ctx.Done()

	errs := make(chan error, len(fns.fns))
	for _, fn := range fns.fns {
		go func(fn func(state *T, ctx context.Context) error) {
			select {
			case <-done:
				errs <- nil
			default:
				errs <- fn(state, ctx)
			}
		}(fn)
	}

	for i := 0; i < len(fns.fns); i++ {
		select {
		case <-done:
			return nil
		case err := <-errs:
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// errorFunctions is a collection of functions to be run on different stages of the task.
type errorFunctions[T any] struct {
	fns []func(state *T, ctx context.Context, err error)
}

// Add adds a function to the collection.
func (fns *errorFunctions[T]) Add(fn func(state *T, ctx context.Context, err error)) {
	fns.fns = append(fns.fns, fn)
}

// Run runs the functions in the collection.
func (fns *errorFunctions[T]) Run(state *T, ctx context.Context, err error) {
	done := ctx.Done()

	for _, fn := range fns.fns {
		select {
		case <-done:
			break
		default:
			fn(state, ctx, err)
		}
	}
	return
}
