package lionrouter

import (
	"net/http"
)

const (
	leafHandlerGet int = iota
	leafHandlerPost
	leafHandlerPut
	leafHandlerPatch
	leafHandlerDelete
	leafHandlerHead
	leafHandlerOptions
)

type leaf struct {
	handler []http.Handler
}

func newLeaf() *leaf {
	return &leaf{
		handler: make([]http.Handler, 0),
	}
}

func handlerMethodFromString(method string) (int, error) {
	switch method {
	case http.MethodGet:
		return leafHandlerGet, nil
	case http.MethodPost:
		return leafHandlerPost, nil
	case http.MethodPut:
		return leafHandlerPut, nil
	case http.MethodPatch:
		return leafHandlerPatch, nil
	case http.MethodDelete:
		return leafHandlerDelete, nil
	case http.MethodHead:
		return leafHandlerHead, nil
	case http.MethodOptions:
		return leafHandlerOptions, nil
	default:
		return 0, ErrUnknownHTTPMethod
	}
}

func (l *leaf) addHandler(method string, handler http.Handler) error {
	handlerMethod, err := handlerMethodFromString(method)

	if err != nil {
		return err
	}

	if len(l.handler) <= handlerMethod {
		sliceLen := handlerMethod - len(l.handler) + 1
		l.handler = append(l.handler, make([]http.Handler, sliceLen)...)
	}

	if l.handler[handlerMethod] != nil {
		return ErrAlreadyAssigned
	}

	l.handler[handlerMethod] = handler
	return nil
}

func (l *leaf) unsetHandler(method string) error {
	handlerMethod, err := handlerMethodFromString(method)

	if err != nil {
		return err
	}

	if len(l.handler) <= handlerMethod {
		return ErrNoHandler
	}

	if l.handler[handlerMethod] == nil {
		return ErrNoHandler
	}

	l.handler[handlerMethod] = nil

	if len(l.handler)-1 == handlerMethod {
		l.handler = l.handler[:len(l.handler)-1]
	}

	return nil
}

func (l *leaf) getHandler(method string) (http.Handler, error) {
	handlerMethod, err := handlerMethodFromString(method)

	if err != nil {
		return nil, err
	}

	if len(l.handler) > handlerMethod {
		return l.handler[handlerMethod], nil
	}

	return nil, nil
}
