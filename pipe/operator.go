package pipe

import . "github.com/langhuihui/RxGo/rx"

//Take
func Take(count uint) Operator {
	return func(source Observable) Observable {
		return source.Take(count)
	}
}

//TakeUntil
func TakeUntil(until Observable) Operator {
	return func(source Observable) Observable {
		return source.TakeUntil(until)
	}
}

//TakeWhile
func TakeWhile(f func(interface{}) bool) Operator {
	return func(source Observable) Observable {
		return source.TakeWhile(f)
	}
}

//Skip
func Skip(count uint) Operator {
	return func(source Observable) Observable {
		return source.Skip(count)
	}
}

//SkipUntil
func SkipUntil(until Observable) Operator {
	return func(source Observable) Observable {
		return source.SkipUntil(until)
	}
}

//SkipWhile
func SkipWhile(f func(interface{}) bool) Operator {
	return func(source Observable) Observable {
		return source.SkipWhile(f)
	}
}

//Share
func Share() Operator {
	return func(source Observable) Observable {
		return source.Share()
	}
}

//StartWith
func StartWith(xs ...interface{}) Operator {
	return func(source Observable) Observable {
		return source.StartWith(xs...)
	}
}

//IgnoreElements
func IgnoreElements() Operator {
	return func(source Observable) Observable {
		return source.IgnoreElements()
	}
}
