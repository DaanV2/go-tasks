package tasks

import (
	"context"
)

// hooks is a collection of functions to be run on different stages of the task.
type hooks[T any] struct {
	todo    functions[T]
	then    functions[T]
	finally functions[T]
	onerror errorFunctions[T]
}

// Task is a collection of functions to be run in a specific order.
type Task[T any] struct {
	//Collection of functions to run on different stages of the task
	hooks[T]

	//State to be passed to the functions
	state *T
}

// New returns a new task with a new state.
func New[T any]() *Task[T] {
	var state T
	return NewWith(&state)
}

// NewWith returns a new task with the given state.
func NewWith[T any](state *T) *Task[T] {
	return &Task[T]{state: state}
}

// State returns the state of the task.
func (t *Task[T]) State() *T {
	return t.state
}

// CopyFor returns a new task with the same functions as the current task, but with a new state.
func (t *Task[T]) CopyFor(state *T) *Task[T] {
	return &Task[T]{
		state: state,
		hooks: t.hooks,
	}
}

// Run runs the task, and returns the first error encountered.
func (t *Task[T]) Run(ctx context.Context) (err error) {
	//Recover
	defer func() {
		if r := recover(); r != nil {
			if recovered, ok := r.(error); ok {
				t.onerror.Run(t.state, ctx, recovered)
			}
		}
	}()
	//Finally
	defer func() {
		finalErr := t.finally.Run(t.state, ctx)
		err = combineErrors(err, finalErr)
	}()

	// Do
	parallelCtx, cancel := context.WithCancel(ctx)

	//If any returns an error, or just done, cancel the parallel context and run the error functions
	err = t.todo.RunParallel(t.state, parallelCtx)
	cancel()
	ctx = context.WithValue(ctx, "error", err)

	if err != nil {
		t.onerror.Run(t.state, ctx, err)
		return
	}

	// Then
	err = t.then.Run(t.state, ctx)
	return
}

// Do adds a function to be called when the task is run. The function will be run in parallel with other functions added with Do.
func (t *Task[T]) Do(fns func(state *T, ctx context.Context) error) *Task[T] {
	t.todo.Add(fns)
	return t
}

// Then adds a function to be called after all functions added with Do have been called. and no errors have been returned.
func (t *Task[T]) Then(fns func(state *T, ctx context.Context) error) *Task[T] {
	t.then.Add(fns)
	return t
}

// Finally adds a function to be called after all other functions have been called.
func (t *Task[T]) Finally(fns func(state *T, ctx context.Context) error) *Task[T] {
	t.finally.Add(fns)
	return t
}

// OnError adds a function to be called if any of the previous functions return an error.
func (t *Task[T]) OnError(fns func(state *T, ctx context.Context, err error)) *Task[T] {
	t.onerror.Add(fns)
	return t
}

// Chain adds a task to be run after the current task and returns the chained task.
func (t *Task[T]) Chain(next *Task[T]) *Task[T] {
	t.Finally(func(state *T, ctx context.Context) error {
		return next.Run(ctx)
	})

	return next
}
