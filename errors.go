package tasks

import "github.com/pkg/errors"

func combineErrors(err1, err2 error) error {
	if err1 == nil {
		return err2
	}
	if err2 == nil {
		return err1
	}
	return errors.Wrap(err1, err2.Error())
}

var Cancel = errors.New("task canceled")
