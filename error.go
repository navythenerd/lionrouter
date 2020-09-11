package lionrouter

import "errors"

var (
	ErrUnknownHTTPMethod = errors.New("unknown or unsupported http method")
	ErrNilHandler        = errors.New("nil handler cannot be assigned")
	ErrAlreadyAssigned   = errors.New("cannot reassign handler")
	ErrAssignment        = errors.New("cannot assign child/leaf to node")
	ErrNoHandler         = errors.New("no handler for http method assigned")
)
