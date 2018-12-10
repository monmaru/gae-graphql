package gql

import "github.com/graphql-go/graphql"

func concurrentResolve(fn graphql.FieldResolveFn) graphql.FieldResolveFn {
	return func(params graphql.ResolveParams) (interface{}, error) {
		type result struct {
			data interface{}
			err  error
		}
		ch := make(chan *result, 1)
		go func() {
			defer close(ch)
			data, err := fn(params)
			ch <- &result{data: data, err: err}
		}()
		return func() (interface{}, error) {
			r := <-ch
			return r.data, r.err
		}, nil
	}
}
