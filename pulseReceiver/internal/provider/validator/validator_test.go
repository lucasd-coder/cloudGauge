package validator_test

import (
	"testing"

	"github.com/lucasd-coder/pulseReceiver/internal/provider/validator"
	"github.com/stretchr/testify/require"
)

type TestStruct struct {
	Name  string `validate:"required"`
	Email string `validate:"required,email"`
}

func TestValidateStruct_ValidInput(t *testing.T) {
	v := validator.NewValidation()

	valid := TestStruct{
		Name:  "Lucas",
		Email: "lucas@example.com",
	}

	err := v.ValidateStruct(valid)
	require.NoError(t, err)
}

func TestValidateStruct_InvalidInput(t *testing.T) {
	v := validator.NewValidation()

	invalid := TestStruct{
		Name:  "",
		Email: "not-an-email",
	}

	err := v.ValidateStruct(invalid)
	require.Error(t, err)
}
