package utils

import (
	"fmt"
	"reflect"
)

type Result[T any] struct {
	has bool
	val *T
	err error
}

func (r Result[T]) String() string {

	obj := new(T)

	return fmt.Sprintf("#result[%s]{has: %v}", reflect.ValueOf(obj).Type().Name(), r.has)
}

func (opt Result[T]) IsOk() bool {
	return opt.has
}

func (opt Result[T]) Unwrap() T {

	if !opt.has {
		panic("unwrapping failed result")
	}

	return *opt.val
}

func (opt Result[T]) UnwrapClassic() (*T, error) {
	return opt.val, opt.err
}

func (opt Result[T]) UnwrapOrDefault(def T) T {

	if opt.IsOk() {
		return opt.Unwrap()
	} else {
		return def
	}

}

func (opt Result[T]) UnwrapOrPanic(msg string) T {
	if opt.IsOk() {
		return opt.Unwrap()
	} else {
		panic(msg)
	}
}

func (opt Result[T]) UnwrapError() error {
	return opt.err
}

func (opt Result[T]) Accept(f func(*T)) (nextResult Result[T]) {

	defer RecoverPanic(func(pe *PanicError) {
		nextResult = ResultFailed[T](pe)
	})

	if opt.IsOk() {
		f(opt.val)
	}
	return opt
}

func (opt Result[T]) Success(f func()) (nextResult Result[T]) {
	if opt.IsOk() {
		f()
	}
	return
}

func (opt Result[T]) Then(f func(*T) *Result[T]) (nextResult Result[T]) {

	defer RecoverPanic(func(pe *PanicError) {
		nextResult.SetFail(pe)
	})

	if opt.IsOk() {
		result := f(opt.val)
		if result != nil {
			return *result
		}
	}
	return opt
}

func (opt Result[T]) Fail(f func(e error)) Result[T] {
	if !opt.IsOk() {
		f(opt.err)
	}
	return opt
}

func (opt *Result[T]) SetOk(v T) *Result[T] {
	opt.val = &v
	opt.has = true
	return opt
}

func (opt *Result[T]) SetFail(v error) *Result[T] {
	opt.err = v
	opt.has = false
	return opt
}

func ResultOk[T any](v T) Result[T] {
	res := Result[T]{}
	res.SetOk(v)

	return res
}

func ToResult(err error) Result[any] {
	if err != nil {
		return ResultFailed[any](err)
	} else {
		return ResultOk[any](nil)
	}
}

func ResultFailed[T any](err error) Result[T] {
	res := Result[T]{}
	res.SetFail(err)
	return res
}
