package main

import "fmt"

type contentTypeInvalid struct{}

func (e *contentTypeInvalid) Error() string {
	return fmt.Sprint("Content type invalid")
}
