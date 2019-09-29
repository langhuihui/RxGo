package pipe

import . "../rx"

//Take 获取最多count数量的事件，然后完成
func Take(count int) Operator {
	return func(source Observable) Observable {
		return source.Take(count)
	}
}

//Skip 跳过若干个数据
func Skip(count int) Operator {
	return func(source Observable) Observable {
		return source.Skip(count)
	}
}

func Share() Operator {
	return func(source Observable) Observable {
		return source.Share()
	}
}

func TakeUntil(until Observable) Operator {
	return func(source Observable) Observable {
		return source.TakeUntil(until)
	}
}
func StartWith(xs ...interface{}) Operator {
	return func(source Observable) Observable {
		return source.StartWith(xs...)
	}
}
