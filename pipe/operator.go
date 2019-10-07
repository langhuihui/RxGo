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

//ElementAt
func ElementAt(index uint) Operator {
	return func(source Observable) Observable {
		return source.ElementAt(index)
	}
}

//Find
func Find(f func(interface{}) bool) Operator {
	return func(source Observable) Observable {
		return source.Find(f)
	}
}

//FindIndex
func FindIndex(f func(interface{}) bool) Operator {
	return func(source Observable) Observable {
		return source.FindIndex(f)
	}
}

//First
func First() Operator {
	return func(source Observable) Observable {
		return source.First()
	}
}

//Last
func Last() Operator {
	return func(source Observable) Observable {
		return source.Last()
	}
}

//Max
func Max() Operator {
	return func(source Observable) Observable {
		return source.Max()
	}
}

//Min
func Min() Operator {
	return func(source Observable) Observable {
		return source.Min()
	}
}

//Reduce
func Reduce(f func(interface{}, interface{}) interface{}) Operator {
	return func(source Observable) Observable {
		return source.Reduce(f)
	}
}

//Map
func Map(f func(interface{}) interface{}) Operator {
	return func(source Observable) Observable {
		return source.Map(f)
	}
}

//MapTo
func MapTo(data interface{}) Operator {
	return func(source Observable) Observable {
		return source.MapTo(data)
	}
}

//MergeMap
func MergeMap(f func(interface{}) Observable, resultSelector func(interface{}, interface{}) interface{}) Operator {
	return func(source Observable) Observable {
		return source.MergeMap(f, resultSelector)
	}
}

//MergeMapTo
func MergeMapTo(then Observable, resultSelector func(interface{}, interface{}) interface{}) Operator {
	return func(source Observable) Observable {
		return source.MergeMap(then, resultSelector)
	}
}

//SwitchMap
func SwitchMap(f func(interface{}) Observable, resultSelector func(interface{}, interface{}) interface{}) Operator {
	return func(source Observable) Observable {
		return source.SwitchMap(f, resultSelector)
	}
}

//SwitchMapTo
func SwitchMapTo(then Observable, resultSelector func(interface{}, interface{}) interface{}) Operator {
	return func(source Observable) Observable {
		return source.SwitchMapTo(then, resultSelector)
	}
}
