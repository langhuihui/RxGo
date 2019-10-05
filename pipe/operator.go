package pipe

import (
	. "github.com/langhuihui/RxGo/rx"
	"time"
)

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

//Do
func Do(f func(interface{})) Operator {
	return func(source Observable) Observable {
		return source.Do(f)
	}
}

//Filter
func Filter(f func(interface{}) bool) Operator {
	return func(source Observable) Observable {
		return source.Filter(f)
	}
}

//Distinct
func Distinct() Operator {
	return func(source Observable) Observable {
		return source.Distinct()
	}
}

//DistinctUntilChanged
func DistinctUntilChanged() Operator {
	return func(source Observable) Observable {
		return source.DistinctUntilChanged()
	}
}

//Debounce
func Debounce(f func(interface{}) Observable) Operator {
	return func(source Observable) Observable {
		return source.Debounce(f)
	}
}

//DebounceTime
func DebounceTime(duration time.Duration) Operator {
	return func(source Observable) Observable {
		return source.DebounceTime(duration)
	}
}

//Throttle
func Throttle(f func(interface{}) Observable) Operator {
	return func(source Observable) Observable {
		return source.Throttle(f)
	}
}

//ThrottleTime
func ThrottleTime(duration time.Duration) Operator {
	return func(source Observable) Observable {
		return source.ThrottleTime(duration)
	}
}
