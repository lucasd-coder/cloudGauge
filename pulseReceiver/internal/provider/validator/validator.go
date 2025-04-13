package validator

import (
	"sync"

	"github.com/go-playground/validator/v10"
)

type Validation struct {
	once     sync.Once
	validate *validator.Validate
}

func NewValidation() *Validation {
	return &Validation{}
}

func (v *Validation) ValidateStruct(s interface{}) error {
	v.lazyInit()

	return v.validate.Struct(s)
}

func (v *Validation) lazyInit() {
	v.once.Do(func() {
		v.validate = validator.New()
	})
}
