package fault

import "errors"

type Tag string

const (
	Untagged            Tag = "UNTAGGED"
	BadRequest          Tag = "BAD_REQUEST"
	NotFound            Tag = "NOT_FOUND"
	InternalServerError Tag = "INTERNAL_SERVER_ERROR"
	Unauthorized        Tag = "UNAUTHORIZED"
	Forbidden           Tag = "FORBIDDEN"
	Conflict            Tag = "CONFLICT"
	TooManyRequests     Tag = "TOO_MANY_REQUESTS"
	ValidationError     Tag = "VALIDATION"
	UnprocessableEntity Tag = "UNPROCESSABLE_ENTITY"
	DB                  Tag = "DATABASE"
	TX                  Tag = "DB_TRANSACTION"
)

// GetTag returns the first tag of the error
//
// Example:
//
//	err := fault.NewBadRequest("invalid request")
//	tag := fault.GetTag(err)
//	fmt.Println(tag) // Output: BAD_REQUEST
//
// Example with switch:
//
//	switch fault.GetTag(err) {
//	case fault.BAD_REQUEST:
//		fmt.Println("bad request")
//	case fault.NOT_FOUND:
//		fmt.Println("not found")
//	default:
//		fmt.Println("unknown error")
//	}
func GetTag(err error) Tag {
	var f *Fault
	if errors.As(err, &f) && f.Tag != "" {
		return f.Tag
	}
	return Untagged
}
