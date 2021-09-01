package handler

import "os"

const DEBUG = "DEBUG"

type handler interface {
	Init() error
	Start()
}

func New() handler {
	var h handler
	if onLambda() {
		h = &lambdaHandler{}
	} else {
		h = &localHandler{}
	}
	return h
}

func onLambda() bool {
	return os.Getenv(DEBUG) == ""
}
