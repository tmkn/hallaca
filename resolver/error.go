package resolver

import "fmt"

type ResolverError struct {
	Msg string
	Pkg *ResolverQueueItem
}

func (e *ResolverError) Error() string {
	if e.Pkg != nil {
		return fmt.Sprintf(e.Msg, e.Pkg.Name)
	}
	return e.Msg
}
