package validate

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

var (
	emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,8}$`)
	urlRegex   = regexp.MustCompile(`^(https?:\/\/)?([\w\d\-]+(\.[\w\d\-]+)*\.[a-z]{2,6})(\/[\w\d\-.\/?=#]*)?$`)
	phoneRegex = regexp.MustCompile(`^\+?[0-9]{10,15}$`)
	ipRegex    = regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)
)

// Validator enhancements for nested structs and improved usability.
type Validator struct {
	data       interface{}
	rules      map[string][]Rule
	errors     map[string][]string
	validateFn map[string]func(interface{}) bool // Map of custom validation functions.
}

func NewValidator(data interface{}) *Validator {
	return &Validator{
		data:       data,
		rules:      make(map[string][]Rule),
		errors:     make(map[string][]string),
		validateFn: make(map[string]func(interface{}) bool),
	}
}

// Rule represents a validation rule.
type Rule struct {
	validate func(value interface{}, params ...interface{}) bool
	msg      func(field string, _ ...interface{}) string
	params   []interface{}
}

// AddRule adds a new rule for a field.
func (v *Validator) AddRule(field string, rule Rule) {
	v.rules[field] = append(v.rules[field], rule)
}

// Validate performs the validation on the provided data.
func (v *Validator) Validate() bool {
	val := reflect.ValueOf(v.data)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for field, rules := range v.rules {
		fieldVal := val.FieldByName(field)
		if !fieldVal.IsValid() {
			continue // Skip if the field doesn't exist.
		}

		for _, rule := range rules {
			if !rule.validate(fieldVal.Interface(), rule.params...) {
				v.errors[field] = append(v.errors[field], rule.msg(field, rule.params...))
			}
		}
	}

	return len(v.errors) == 0
}

// AddError adds a new error for a field.
func (v *Validator) AddError(field, msg string) {
	v.errors[field] = append(v.errors[field], msg)
}

// GetErrors returns the validation errors.
func (v *Validator) GetErrors() map[string][]string {
	return v.errors
}

// Required Helper functions to create common rules.
func Required() Rule {
	return Rule{
		validate: func(value interface{}, _ ...interface{}) bool {
			switch v := value.(type) {
			case string:
				return strings.TrimSpace(v) != ""
			default:
				return value != nil
			}
		},
		msg: func(field string, _ ...interface{}) string {
			return field + " is required"
		},
	}
}

func Email() Rule {
	return Rule{
		validate: func(value interface{}, _ ...interface{}) bool {
			str, ok := value.(string)
			return ok && emailRegex.MatchString(str)
		},
		msg: func(field string, _ ...interface{}) string {
			return field + " is not a valid email address"
		},
	}
}

func StrongPassword() Rule {
	return Rule{
		validate: func(value interface{}, _ ...interface{}) bool {
			password, ok := value.(string)
			if !ok || len(password) < 8 {
				return false
			}
			hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
			hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
			hasDigit := regexp.MustCompile(`\d`).MatchString(password)
			hasSpecial := regexp.MustCompile(`[\W_]`).MatchString(password)
			return hasUpper && hasLower && hasDigit && hasSpecial
		},
		msg: func(field string, _ ...interface{}) string {
			return field + " is not a strong enough password"
		},
	}
}

func URL() Rule {
	return Rule{
		validate: func(value interface{}, _ ...interface{}) bool {
			str, ok := value.(string)
			return ok && urlRegex.MatchString(str)
		},
		msg: func(field string, _ ...interface{}) string {
			return field + " is not a valid URL"
		},
	}
}

func PhoneNumber() Rule {
	// Adjust the regex according to the phone number formats you want to support
	return Rule{
		validate: func(value interface{}, _ ...interface{}) bool {
			str, ok := value.(string)
			return ok && phoneRegex.MatchString(str)
		},
		msg: func(field string, _ ...interface{}) string {
			return field + " is not a valid phone number"
		},
	}
}

func Min(min int) Rule {
	return Rule{
		validate: func(value interface{}, _ ...interface{}) bool {
			switch v := value.(type) {
			case string:
				return len(v) >= min
			case int, int8, int16, int32, int64, float32, float64:
				return reflect.ValueOf(value).Convert(reflect.TypeOf(float64(0))).Float() >= float64(min)
			default:
				return false
			}
		},
		msg: func(field string, _ ...interface{}) string {
			return fmt.Sprintf("%s must be at least %d", field, min)
		},
	}
}

func Max(max int) Rule {
	return Rule{
		validate: func(value interface{}, _ ...interface{}) bool {
			switch v := value.(type) {
			case string:
				return len(v) <= max
			case int, int8, int16, int32, int64, float32, float64:
				return reflect.ValueOf(value).Convert(reflect.TypeOf(float64(0))).Float() <= float64(max)
			default:
				return false
			}
		},
		msg: func(field string, _ ...interface{}) string {
			return fmt.Sprintf("%s must not exceed %d", field, max)
		},
	}
}

func IP() Rule {
	return Rule{
		validate: func(value interface{}, _ ...interface{}) bool {
			str, ok := value.(string)
			return ok && ipRegex.MatchString(str)
		},
		msg: func(field string, _ ...interface{}) string {
			return field + " is not a valid IP address"
		},
	}
}

func In(values ...interface{}) Rule {
	return Rule{
		validate: func(value interface{}, _ ...interface{}) bool {
			for _, v := range values {
				if reflect.DeepEqual(v, value) {
					return true
				}
			}
			return false
		},
		msg: func(field string, _ ...interface{}) string {
			return field + " is not a valid value"
		},
	}
}

func NotIn(values ...interface{}) Rule {
	return Rule{
		validate: func(value interface{}, _ ...interface{}) bool {
			for _, v := range values {
				if reflect.DeepEqual(v, value) {
					return false
				}
			}
			return true
		},
		msg: func(field string, _ ...interface{}) string {
			return field + " is not a valid value"
		},
	}
}

// Custom - Create a custom validation rule.
// Usage:
// 	v.AddRule("Field", validate.Custom(func(value interface{}, _ ...interface{}) bool {
//		return value == "some value"
//	}, "Field must be some value"))

func Custom(fn func(interface{}, ...interface{}) bool, msg string) Rule {
	return Rule{
		validate: fn,
		msg: func(_ string, _ ...interface{}) string {
			return msg
		},
	}
}
