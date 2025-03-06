package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		email    string
		expected bool
	}{
		// ✅ Valid email cases
		{"user@example.com", true},
		{"1234567890@example.com", true},
		{"user_name@example.com", true},
		{"user@sub.example.com", true},
		{"\"quoted@name\"@example.com", true}, // Quoted local part

		// ❌ Invalid email cases
		{"plainaddress", false},                     // No @ symbol
		{"@missingusername.com", false},             // Missing username
		{"username@.com", false},                    // Invalid domain
		{"username@com", false},                     // No TLD
		{"username@domain..com", false},             // Double dot in domain
		{"user name@example.com", false},            // Space in email
		{"user@exam_ple.com", false},                // Underscore in domain
		{"user@-example.com", false},                // Leading hyphen in domain
		{"user@example-.com", false},                // Trailing hyphen in domain
		{"user@.example.com", false},                // Leading dot in domain
		{"user@invalid_domain.com", false},          // Underscore in domain name
        {"user@example", false},                     // Missing domain extension
	}

	for _, test := range tests {
		t.Run(test.email, func(t *testing.T) {
			assert.Equal(t, test.expected, IsEmailValid(test.email), "Failed for email: %s", test.email)
		})
	}
}
