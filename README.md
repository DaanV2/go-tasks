# Go Tasks

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/DaanV2/go-tasks)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/DaanV2/go-tasks)
[![🐹 Golang](https://github.com/DaanV2/go-tasks/actions/workflows/go-checks.yml/badge.svg)](https://github.com/DaanV2/go-tasks/actions/workflows/go-checks.yml)

## Examples
```golang
//State object
type User struct {
	Email  string
	Name   string
	Age    int
	Scopes []string
}

//Example 1

//Create a new task around the state object
task := tasks.New[User]()
task.
	//Add a tasks to the chain that will execute in order
	Do(RetrieveFromCache).
	Do(RetrieveFromDatabase).
	//What to do if the tasks is successful
	Then(UpdateCache).
	// What to do when the task is done
	Finally(func(user *User, ctx context.Context) error {
		zap.S().Infof("User: %+v", user)
		return nil
	}).
	//What to do if the tasks fails
	OnError(func(user *User, ctx context.Context, err error) {
		zap.S().Errorf("Error: %+v", err)
	})

err := task.Run(context.Background())

//Example 2
task := NewWith[User](&User{Email: "foo@bar.com"})
task.
	// If the first function returns an error, run the second function.
	Do(IfElse(RetrieveFromCache, RetrieveFromDatabase)).
	Then(UpdateCache)

//Example 3
task := tasks.New[User]()
task.
	Do(RetrieveFromCache).
	Do(RetrieveFromDatabase).
	Then(UpdateCache).
	Finally(func(user *User, ctx context.Context) error {
		zap.S().Infof("User: %+v", user)
		return nil
	}).
	OnError(func(user *User, ctx context.Context, err error) {
		zap.S().Errorf("Error: %+v", err)
	})

newTask := task.CopyFor(&User{Email: "foo@bar.com"})

//Example 3
next := task.Chain(New[User]())
next.
	Do(ValidateUser).
	Finally(func(user *User, ctx context.Context) error {
		err := ctx.Value("error").(error)
		if err != nil {
			zap.S().Errorf("Error: %+v", err)
		}

		return nil
	})

```
