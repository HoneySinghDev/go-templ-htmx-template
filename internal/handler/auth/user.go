package auth

import "github.com/HoneySinghDev/go-templ-htmx-template/pkg/kit/validate"

type UserCreds struct {
	EmailID         string
	Password        string
	ConfirmPassword string
}

func (u *UserCreds) Validate() (map[string][]string, bool) {
	v := validate.NewValidator(u)
	v.AddRule("EmailID", validate.Email())
	v.AddRule("Password", validate.StrongPassword())

	// This is a custom rule, not provided by the package
	if u.ConfirmPassword != "" && u.Password != u.ConfirmPassword {
		v.AddError("Password", "Passwords do not match")
	}

	if ok := v.Validate(); !ok {
		return v.GetErrors(), ok
	}

	return nil, true
}
