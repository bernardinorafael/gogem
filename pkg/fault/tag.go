package fault

import "errors"

type Tag string

// Tag is a string that represents the type of error
//
// You can create your own tags by defining a new Tag constant
// and adding it to the list below
// Example:
//
//	const (
//		MyTag fault.Tag = "MY_TAG"
//		MyTag2 fault.Tag = "MY_TAG2"
//		MyTag3 fault.Tag = "MY_TAG3"
//		MyTag4 fault.Tag = "MY_TAG4"
//	)
const (
	UNTAGGED              Tag = "UNTAGGED"
	BAD_REQUEST           Tag = "BAD_REQUEST"
	NOT_FOUND             Tag = "NOT_FOUND"
	INTERNAL_SERVER_ERROR Tag = "INTERNAL_SERVER_ERROR"
	UNAUTHORIZED          Tag = "UNAUTHORIZED"
	FORBIDDEN             Tag = "FORBIDDEN"
	CONFLICT              Tag = "CONFLICT"
	TOO_MANY_REQUESTS     Tag = "TOO_MANY_REQUESTS"
	UNPROCESSABLE_ENTITY  Tag = "UNPROCESSABLE_ENTITY"
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
	if err == nil {
		return UNTAGGED
	}

	for err != nil {
		e, ok := err.(*Fault)
		if ok && e.Tag != "" {
			return e.Tag
		}
		err = errors.Unwrap(err)
	}

	return UNTAGGED
}
