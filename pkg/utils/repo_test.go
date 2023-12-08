package utils

import (
	"github.com/romandnk/todo/internal/constant"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSetPlaceholders(t *testing.T) {
	type input struct {
		placeholder string
		quantity    int
	}

	testCases := []struct {
		name           string
		input          input
		expectedOutput string
		expectedError  error
	}{
		{
			name: "$ placeholder",
			input: input{
				placeholder: "$",
				quantity:    4,
			},
			expectedOutput: "($1, $2, $3, $4)",
			expectedError:  nil,
		},
		{
			name: "? placeholder",
			input: input{
				placeholder: "?",
				quantity:    4,
			},
			expectedOutput: "(?1, ?2, ?3, ?4)",
			expectedError:  nil,
		},
		{
			name: "zero quantity",
			input: input{
				placeholder: "$",
				quantity:    0,
			},
			expectedOutput: "",
			expectedError:  constant.ErrNonPositiveQuantity,
		},
		{
			name: "negative quantity",
			input: input{
				placeholder: "$",
				quantity:    -1,
			},
			expectedOutput: "",
			expectedError:  constant.ErrNonPositiveQuantity,
		},
		{
			name: "empty placeholder",
			input: input{
				placeholder: "",
				quantity:    4,
			},
			expectedOutput: "",
			expectedError:  constant.ErrEmptyPlaceholder,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			output, err := SetPlaceholders(tc.input.placeholder, tc.input.quantity)
			require.ErrorIs(t, err, tc.expectedError)
			require.Equal(t, tc.expectedOutput, output)
		})
	}
}
