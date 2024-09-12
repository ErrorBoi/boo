package utils

import (
	"testing"

	"github.com/errorboi/boo/internal/locale"
)

func TestGetBalance(t *testing.T) {
	// create several test cases
	tests := []struct {
		name     string
		balance  int64
		expected string
	}{
		{
			name:     "below 1k",
			balance:  999,
			expected: "💰 *Balance*: 999 $eBOO",
		},
		{
			name:     "k",
			balance:  1000,
			expected: "💰 *Balance*: 1000 $eBOO",
		},
		{
			name:     "several k",
			balance:  123456,
			expected: "💰 *Balance*: 123456 $eBOO",
		},
		{
			name:     "M",
			balance:  1000000,
			expected: "💰 *Balance*: 1000k $eBOO",
		},
		{
			name:     "several M",
			balance:  123456789,
			expected: "💰 *Balance*: 123456k $eBOO",
		},
		{
			name:     "B",
			balance:  1000000000,
			expected: "💰 *Balance*: 1000M $eBOO",
		},
		{
			name:     "several B",
			balance:  123456789012,
			expected: "💰 *Balance*: 123456M $eBOO",
		},
	}

	// iterate over test cases
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// call the function
			result := GetBalanceText(test.balance, locale.English)

			// compare the result with the expected value
			if result != test.expected {
				t.Errorf("expected %s, got %s", test.expected, result)
			}
		})
	}
}
