package pipe

import . "../rx"

func Take(count int) Operator {
	return func(source Observable) Observable {
		return source.Take(count)
	}
}

func TakeUntil(until Observable) Operator {
	return func(source Observable) Observable {
		return source.TakeUntil(until)
	}
}
func TakeWhile(f func(interface{}) bool) Operator {
	return func(source Observable) Observable {
		return source.TakeWhile(f)
	}
}
func Skip(count int) Operator {
	return func(source Observable) Observable {
		return source.Skip(count)
	}
}
func SkipUntil(until Observable) Operator {
	return func(source Observable) Observable {
		return source.SkipUntil(until)
	}
}
func SkipWhile(f func(interface{}) bool) Operator {
	return func(source Observable) Observable {
		return source.SkipWhile(f)
	}
}
func Share() Operator {
	return func(source Observable) Observable {
		return source.Share()
	}
}

func StartWith(xs ...interface{}) Operator {
	return func(source Observable) Observable {
		return source.StartWith(xs...)
	}
}
func IgnoreElements() Operator {
	return func(source Observable) Observable {
		return source.IgnoreElements()
	}
}
