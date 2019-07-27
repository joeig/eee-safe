package threema

import "testing"

func TestValidateUserAgent(t *testing.T) {
	t.Run("UserAgentOK", func(t *testing.T) {
		userAgent := "Threema 1.2.3.4"
		if err := ValidateUserAgent(userAgent); err != nil {
			t.Errorf("User agent validation failed wrongly: %s", userAgent)
		}
	})
	t.Run("UserAgentNOK", func(t *testing.T) {
		userAgent := "foo"
		if err := ValidateUserAgent(userAgent); err == nil {
			t.Errorf("User agent validation succeeded wrongly: %s", userAgent)
		}
	})
}
