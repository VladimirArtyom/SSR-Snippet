package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

const (

	BLANK_MESSAGE = "Cannot be Blank"
	MAX_CHAR_MESSAGE = "Cannot be more than %d characters"
	NOT_IN_OPTIONS = "Not in Options %d"
	INVALID_EMAIL = "Invalid Email"
	MIN_CHAR_MESSAGE = "Must be at least %d characters; yours %d"
	DUPLICATE_EMAIL = "Email address is already taken"

)

var EmailRX *regexp.Regexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
type Validator struct {
	
	FieldError map[string]string
}


func IsNotBlank(str string) bool {

	return strings.TrimSpace(str) != ""

}

func IsNotMaxChars(str string, n int) bool {
	return utf8.RuneCountInString(str) <= n
}

func PermittedInt(value int, permittedValues ...int) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}

func Matches(value string, regex regexp.Regexp) bool {
	return regex.MatchString(value)
}

func MinChars(value string, min int) bool {
	return utf8.RuneCountInString(value) >= min
}


func (v *Validator) CheckField(ok bool, key string, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

func (v *Validator) AddFieldError(key string, value string) {
	if len(v.FieldError) == 0 {
		v.FieldError = make(map[string]string)
	}
	
	_, exists := v.FieldError[key]
	if !exists {
		v.FieldError[key] = value
	} 
}

func (v *Validator) IsValid() bool {
	return len(v.FieldError) == 0
}

