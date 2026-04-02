// all type of validation (input, email , name , etc..)

package validator

import (
	"regexp"
	"strings"
)

// regex for email

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// valid function

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(field, message string) {
	if _, exists := v.Errors[field]; !exists {
		v.Errors[field] = message
	}
}

// checking the field and message

func (v *Validator) Check(ok bool, field, message string) {
	if !ok {
		v.AddError(field, message)
	}
}

// required string function to validate the input string

func (v *Validator) RequiredString(value, field string) {
	v.Check(strings.TrimSpace(value) != "", field, field+" is required")
}

func (v *Validator) MinLength(value, field string, min int) {
	v.Check(len(strings.TrimSpace(value)) >= min, field, field+" must be at least "+string(rune('0'+min))+" characters")
}

// email validation

func (v *Validator) ValidEmail(email, field string) {
	v.Check(emailRegex.MatchString(email), field, "invalid email address")
}

// date validation

func (v *Validator) ValidDate(date, field string) {
	matched, _ := regexp.MatchString(`^\d{4}-\d{2}-\d{2}$`, date)
	v.Check(matched, field, field+" must be in YYYY-MM-DD format")
}

// values must be greater than zero

func (v *Validator) PositiveFloat(value float64, field string) {
	v.Check(value > 0, field, field+" must be greater than 0")
}

// a selector function like, chosen thing must be one of them

func (v *Validator) OneOf(value, field string, allowed ...string) {
	for _, a := range allowed {
		if value == a {
			return
		}
	}

	v.AddError(field, field+" must be one of: "+strings.Join(allowed, ","))
}
