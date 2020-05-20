package main

type contentTypeInvalid struct{}

func (e *contentTypeInvalid) Error() string {
	return "Content type invalid"
}
