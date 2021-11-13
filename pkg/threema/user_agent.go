package threema

import "strings"

// ValidUserAgentSubstring contains a string which must be part of the "User-Agent" header
const ValidUserAgentSubstring = "Threema"

// ValidateUserAgent checks if a user agent meets the requirements
func ValidateUserAgent(userAgent string) error {
	if !strings.Contains(userAgent, ValidUserAgentSubstring) {
		return &ErrInvalidUserAgent{userAgent}
	}

	return nil
}
