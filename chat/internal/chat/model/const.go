package model

const (
	// ChatDefaultTitle is a default title for chat.
	ChatDefaultTitle string = "Новый чат"
)

// ResponseStatus enum of statuses processing response.
type ResponseStatus uint8

const (
	_ ResponseStatus = iota
	// StatusCreated empty response created.
	StatusCreated
	// StatusProcessing response is being generated.
	StatusProcessing
	// StatusSuccess response is successfully generated.
	StatusSuccess
	// StatusError response generation failed.
	StatusError
	// StatusCanceled response generation was canceled.
	StatusCanceled
)
