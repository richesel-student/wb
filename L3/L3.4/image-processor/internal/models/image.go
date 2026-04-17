package models

import "github.com/go-playground/validator/v10"

var validate = validator.New()

type Status string

const (
	StatusQueued     Status = "QUEUED"
	StatusProcessing Status = "PROCESSING"
	StatusDone       Status = "DONE"
	StatusFailed     Status = "FAILED"
)

type Image struct {
	ID            string `validate:"required,uuid"`
	Status        Status `validate:"required,oneof=QUEUED PROCESSING DONE FAILED"`
	OriginalPath  string `validate:"required"`
	ProcessedPath string
}

func (i *Image) Validate() error {
	return validate.Struct(i)
}
