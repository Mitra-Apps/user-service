// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: proto/user/user.proto

package user

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on User with the rules defined in the proto
// definition for this message. If any rules are violated, the first error
// encountered is returned, or nil if there are no violations.
func (m *User) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on User with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in UserMultiError, or nil if none found.
func (m *User) ValidateAll() error {
	return m.validate(true)
}

func (m *User) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	// no validation rules for Username

	// no validation rules for Password

	// no validation rules for Email

	// no validation rules for PhoneNumber

	// no validation rules for AvatarImageId

	// no validation rules for AccessToken

	// no validation rules for IsActive

	// no validation rules for Name

	// no validation rules for Address

	if len(errors) > 0 {
		return UserMultiError(errors)
	}

	return nil
}

// UserMultiError is an error wrapping multiple validation errors returned by
// User.ValidateAll() if the designated constraints aren't met.
type UserMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UserMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UserMultiError) AllErrors() []error { return m }

// UserValidationError is the validation error returned by User.Validate if the
// designated constraints aren't met.
type UserValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UserValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UserValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UserValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UserValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UserValidationError) ErrorName() string { return "UserValidationError" }

// Error satisfies the builtin error interface
func (e UserValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUser.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UserValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UserValidationError{}

// Validate checks the field values on Role with the rules defined in the proto
// definition for this message. If any rules are violated, the first error
// encountered is returned, or nil if there are no violations.
func (m *Role) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Role with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in RoleMultiError, or nil if none found.
func (m *Role) ValidateAll() error {
	return m.validate(true)
}

func (m *Role) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	// no validation rules for RoleName

	// no validation rules for Description

	// no validation rules for IsActive

	if all {
		switch v := interface{}(m.GetPermission()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, RoleValidationError{
					field:  "Permission",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, RoleValidationError{
					field:  "Permission",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetPermission()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return RoleValidationError{
				field:  "Permission",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return RoleMultiError(errors)
	}

	return nil
}

// RoleMultiError is an error wrapping multiple validation errors returned by
// Role.ValidateAll() if the designated constraints aren't met.
type RoleMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m RoleMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m RoleMultiError) AllErrors() []error { return m }

// RoleValidationError is the validation error returned by Role.Validate if the
// designated constraints aren't met.
type RoleValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e RoleValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e RoleValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e RoleValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e RoleValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e RoleValidationError) ErrorName() string { return "RoleValidationError" }

// Error satisfies the builtin error interface
func (e RoleValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRole.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = RoleValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = RoleValidationError{}

// Validate checks the field values on ListRole with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *ListRole) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ListRole with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in ListRoleMultiError, or nil
// if none found.
func (m *ListRole) ValidateAll() error {
	return m.validate(true)
}

func (m *ListRole) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetRoles() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, ListRoleValidationError{
						field:  fmt.Sprintf("Roles[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ListRoleValidationError{
						field:  fmt.Sprintf("Roles[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ListRoleValidationError{
					field:  fmt.Sprintf("Roles[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return ListRoleMultiError(errors)
	}

	return nil
}

// ListRoleMultiError is an error wrapping multiple validation errors returned
// by ListRole.ValidateAll() if the designated constraints aren't met.
type ListRoleMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ListRoleMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ListRoleMultiError) AllErrors() []error { return m }

// ListRoleValidationError is the validation error returned by
// ListRole.Validate if the designated constraints aren't met.
type ListRoleValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListRoleValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListRoleValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListRoleValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListRoleValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListRoleValidationError) ErrorName() string { return "ListRoleValidationError" }

// Error satisfies the builtin error interface
func (e ListRoleValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListRole.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListRoleValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListRoleValidationError{}

// Validate checks the field values on UserLoginRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *UserLoginRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UserLoginRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// UserLoginRequestMultiError, or nil if none found.
func (m *UserLoginRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *UserLoginRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if err := m._validateEmail(m.GetEmail()); err != nil {
		err = UserLoginRequestValidationError{
			field:  "Email",
			reason: "value must be a valid email address",
			cause:  err,
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if l := utf8.RuneCountInString(m.GetPassword()); l < 6 || l > 8 {
		err := UserLoginRequestValidationError{
			field:  "Password",
			reason: "value length must be between 6 and 8 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return UserLoginRequestMultiError(errors)
	}

	return nil
}

func (m *UserLoginRequest) _validateHostname(host string) error {
	s := strings.ToLower(strings.TrimSuffix(host, "."))

	if len(host) > 253 {
		return errors.New("hostname cannot exceed 253 characters")
	}

	for _, part := range strings.Split(s, ".") {
		if l := len(part); l == 0 || l > 63 {
			return errors.New("hostname part must be non-empty and cannot exceed 63 characters")
		}

		if part[0] == '-' {
			return errors.New("hostname parts cannot begin with hyphens")
		}

		if part[len(part)-1] == '-' {
			return errors.New("hostname parts cannot end with hyphens")
		}

		for _, r := range part {
			if (r < 'a' || r > 'z') && (r < '0' || r > '9') && r != '-' {
				return fmt.Errorf("hostname parts can only contain alphanumeric characters or hyphens, got %q", string(r))
			}
		}
	}

	return nil
}

func (m *UserLoginRequest) _validateEmail(addr string) error {
	a, err := mail.ParseAddress(addr)
	if err != nil {
		return err
	}
	addr = a.Address

	if len(addr) > 254 {
		return errors.New("email addresses cannot exceed 254 characters")
	}

	parts := strings.SplitN(addr, "@", 2)

	if len(parts[0]) > 64 {
		return errors.New("email address local phrase cannot exceed 64 characters")
	}

	return m._validateHostname(parts[1])
}

// UserLoginRequestMultiError is an error wrapping multiple validation errors
// returned by UserLoginRequest.ValidateAll() if the designated constraints
// aren't met.
type UserLoginRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UserLoginRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UserLoginRequestMultiError) AllErrors() []error { return m }

// UserLoginRequestValidationError is the validation error returned by
// UserLoginRequest.Validate if the designated constraints aren't met.
type UserLoginRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UserLoginRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UserLoginRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UserLoginRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UserLoginRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UserLoginRequestValidationError) ErrorName() string { return "UserLoginRequestValidationError" }

// Error satisfies the builtin error interface
func (e UserLoginRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUserLoginRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UserLoginRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UserLoginRequestValidationError{}

// Validate checks the field values on UserRegisterRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *UserRegisterRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UserRegisterRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// UserRegisterRequestMultiError, or nil if none found.
func (m *UserRegisterRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *UserRegisterRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if err := m._validateEmail(m.GetEmail()); err != nil {
		err = UserRegisterRequestValidationError{
			field:  "Email",
			reason: "value must be a valid email address",
			cause:  err,
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if l := utf8.RuneCountInString(m.GetPassword()); l < 6 || l > 8 {
		err := UserRegisterRequestValidationError{
			field:  "Password",
			reason: "value length must be between 6 and 8 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if l := utf8.RuneCountInString(m.GetName()); l < 1 || l > 255 {
		err := UserRegisterRequestValidationError{
			field:  "Name",
			reason: "value length must be between 1 and 255 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if l := utf8.RuneCountInString(m.GetPhoneNumber()); l < 9 || l > 14 {
		err := UserRegisterRequestValidationError{
			field:  "PhoneNumber",
			reason: "value length must be between 9 and 14 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	// no validation rules for Address

	if len(errors) > 0 {
		return UserRegisterRequestMultiError(errors)
	}

	return nil
}

func (m *UserRegisterRequest) _validateHostname(host string) error {
	s := strings.ToLower(strings.TrimSuffix(host, "."))

	if len(host) > 253 {
		return errors.New("hostname cannot exceed 253 characters")
	}

	for _, part := range strings.Split(s, ".") {
		if l := len(part); l == 0 || l > 63 {
			return errors.New("hostname part must be non-empty and cannot exceed 63 characters")
		}

		if part[0] == '-' {
			return errors.New("hostname parts cannot begin with hyphens")
		}

		if part[len(part)-1] == '-' {
			return errors.New("hostname parts cannot end with hyphens")
		}

		for _, r := range part {
			if (r < 'a' || r > 'z') && (r < '0' || r > '9') && r != '-' {
				return fmt.Errorf("hostname parts can only contain alphanumeric characters or hyphens, got %q", string(r))
			}
		}
	}

	return nil
}

func (m *UserRegisterRequest) _validateEmail(addr string) error {
	a, err := mail.ParseAddress(addr)
	if err != nil {
		return err
	}
	addr = a.Address

	if len(addr) > 254 {
		return errors.New("email addresses cannot exceed 254 characters")
	}

	parts := strings.SplitN(addr, "@", 2)

	if len(parts[0]) > 64 {
		return errors.New("email address local phrase cannot exceed 64 characters")
	}

	return m._validateHostname(parts[1])
}

// UserRegisterRequestMultiError is an error wrapping multiple validation
// errors returned by UserRegisterRequest.ValidateAll() if the designated
// constraints aren't met.
type UserRegisterRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UserRegisterRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UserRegisterRequestMultiError) AllErrors() []error { return m }

// UserRegisterRequestValidationError is the validation error returned by
// UserRegisterRequest.Validate if the designated constraints aren't met.
type UserRegisterRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UserRegisterRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UserRegisterRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UserRegisterRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UserRegisterRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UserRegisterRequestValidationError) ErrorName() string {
	return "UserRegisterRequestValidationError"
}

// Error satisfies the builtin error interface
func (e UserRegisterRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUserRegisterRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UserRegisterRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UserRegisterRequestValidationError{}

// Validate checks the field values on SuccessResponse with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *SuccessResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on SuccessResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// SuccessResponseMultiError, or nil if none found.
func (m *SuccessResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *SuccessResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetData()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, SuccessResponseValidationError{
					field:  "Data",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, SuccessResponseValidationError{
					field:  "Data",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetData()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return SuccessResponseValidationError{
				field:  "Data",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return SuccessResponseMultiError(errors)
	}

	return nil
}

// SuccessResponseMultiError is an error wrapping multiple validation errors
// returned by SuccessResponse.ValidateAll() if the designated constraints
// aren't met.
type SuccessResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m SuccessResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m SuccessResponseMultiError) AllErrors() []error { return m }

// SuccessResponseValidationError is the validation error returned by
// SuccessResponse.Validate if the designated constraints aren't met.
type SuccessResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e SuccessResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e SuccessResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e SuccessResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e SuccessResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e SuccessResponseValidationError) ErrorName() string { return "SuccessResponseValidationError" }

// Error satisfies the builtin error interface
func (e SuccessResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sSuccessResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = SuccessResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = SuccessResponseValidationError{}

// Validate checks the field values on GetUsersRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *GetUsersRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetUsersRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetUsersRequestMultiError, or nil if none found.
func (m *GetUsersRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetUsersRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return GetUsersRequestMultiError(errors)
	}

	return nil
}

// GetUsersRequestMultiError is an error wrapping multiple validation errors
// returned by GetUsersRequest.ValidateAll() if the designated constraints
// aren't met.
type GetUsersRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetUsersRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetUsersRequestMultiError) AllErrors() []error { return m }

// GetUsersRequestValidationError is the validation error returned by
// GetUsersRequest.Validate if the designated constraints aren't met.
type GetUsersRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetUsersRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetUsersRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetUsersRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetUsersRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetUsersRequestValidationError) ErrorName() string { return "GetUsersRequestValidationError" }

// Error satisfies the builtin error interface
func (e GetUsersRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetUsersRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetUsersRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetUsersRequestValidationError{}

// Validate checks the field values on GetUsersResponse with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *GetUsersResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetUsersResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetUsersResponseMultiError, or nil if none found.
func (m *GetUsersResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *GetUsersResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetUsers() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, GetUsersResponseValidationError{
						field:  fmt.Sprintf("Users[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, GetUsersResponseValidationError{
						field:  fmt.Sprintf("Users[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return GetUsersResponseValidationError{
					field:  fmt.Sprintf("Users[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return GetUsersResponseMultiError(errors)
	}

	return nil
}

// GetUsersResponseMultiError is an error wrapping multiple validation errors
// returned by GetUsersResponse.ValidateAll() if the designated constraints
// aren't met.
type GetUsersResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetUsersResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetUsersResponseMultiError) AllErrors() []error { return m }

// GetUsersResponseValidationError is the validation error returned by
// GetUsersResponse.Validate if the designated constraints aren't met.
type GetUsersResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetUsersResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetUsersResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetUsersResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetUsersResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetUsersResponseValidationError) ErrorName() string { return "GetUsersResponseValidationError" }

// Error satisfies the builtin error interface
func (e GetUsersResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetUsersResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetUsersResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetUsersResponseValidationError{}
