package errors

import (
	"fmt"
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestAPIError(t *testing.T) {
	newErr := New("test error")
	apiErr := ErrIllegalObjectName
	err := apiErr.WrapErrNoStack(newErr)

	assert.Equal(t, Is(err, ErrIllegalObjectName), true) // because not with stack
	assert.Equal(t, Is(err, newErr), true)

	fmt.Printf("%v\n\n", err)
	fmt.Printf("%+v\n\n", err)
	fmt.Println(err.Error())
	fmt.Println()

	stackErr := apiErr.WrapErr(newErr)
	fmt.Printf("%v\n\n", stackErr)
	fmt.Printf("%+v\n\n", stackErr)
	fmt.Println(stackErr.Error())

	apiErr2 := &APIError{}
	assert.Equal(t, As(err, &apiErr2), true)
	assert.Equal(t, Is(apiErr2, ErrIllegalObjectName), true)
	assert.Equal(t, apiErr2, ErrIllegalObjectName)
}
