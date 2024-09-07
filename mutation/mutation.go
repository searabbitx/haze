package mutation

import (
	"errors"
	"github.com/kamil-s-solecki/haze/http"
)

var mutationNotApplicable = errors.New("Mutation not applicable")

type Mutable string

type Mutation func(http.Request, Mutable) (http.Request, error)

const (
	Path      Mutable = "Path"
	Parameter Mutable = "Parameter"
)

func SingleQuotes(rq http.Request, mutable Mutable) (mrq http.Request, err error) {
	mrq = rq
	err = nil
	switch mutable {
	case Parameter:
		mrq = rq.WithQuery(rq.Query + "'")
	case Path:
		mrq = rq.WithPath(rq.Path + "'")
	default:
		err = mutationNotApplicable
	}
	return
}

// REFACTOR NOW!!

func DoubleQuotes(rq http.Request, mutable Mutable) (mrq http.Request, err error) {
	mrq = rq
	err = nil
	switch mutable {
	case Parameter:
		mrq = rq.WithQuery(rq.Query + "\"")
	case Path:
		mrq = rq.WithPath(rq.Path + "\"")
	default:
		err = mutationNotApplicable
	}
	return
}

func AllMutations() []Mutation {
	return []Mutation{SingleQuotes, DoubleQuotes}
}

func AllMutatables() []Mutable {
	return []Mutable{Path, Parameter}
}

func Mutate(rq http.Request, mutations []Mutation, mutables []Mutable) []http.Request {
	result := []http.Request{}
	for _, mutation := range mutations {
		for _, mutable := range mutables {
			mrq, err := mutation(rq, mutable)
			if err == nil {
				result = append(result, mrq)
			}
		}
	}
	return result
}
