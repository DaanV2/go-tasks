package tasks

import "context"

type functions[T any] struct {
	fns []func(state *T, ctx context.Context) error
}

func (fns *functions[T]) Add(fn func(state *T, ctx context.Context) error) {
	fns.fns = append(fns.fns, fn)
}

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

type errorFunctions[T any] struct {
	fns []func(state *T, ctx context.Context, err error)
}

func (fns *errorFunctions[T]) Add(fn func(state *T, ctx context.Context, err error)) {
	fns.fns = append(fns.fns, fn)
}

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
