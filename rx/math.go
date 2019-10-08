package rx

//Count 计算元素总数然后返回
func (ob Observable) Count() Observable {
	return func(sink *Observer) error {
		count := 0
		defer func() {
			sink.Next(count)
		}()
		return ob(sink.New3(NextFunc(func(event *Event) {
			count++
		})))
	}
}

//Max 完成时返回最大值
func (ob Observable) Max() Observable {
	return func(sink *Observer) error {
		var max int
		defer func() {
			sink.Next(max)
		}()
		return ob(sink.New3(NextFunc(func(event *Event) {
			if data := event.Data.(int); data > max {
				max = data
			}
		})))
	}
}

//Min 完成时返回最小值
func (ob Observable) Min() Observable {
	return func(sink *Observer) error {
		var min int
		defer func() {
			sink.Next(min)
		}()
		return ob(sink.New3(NextFunc(func(event *Event) {
			if data := event.Data.(int); data < min {
				min = data
			}
		})))
	}
}

//Reduce 完成时返回累加值
func (ob Observable) Reduce(f func(interface{}, interface{}) interface{}) Observable {
	return func(sink *Observer) error {
		var aac interface{}
		defer func() {
			sink.Next(aac)
		}()
		aacNext := func(event *Event) {
			aac = f(aac, event.Data)
		}
		return ob(sink.New3(NextFunc(func(event *Event) {
			aac = event.Data
			event.Target.next = NextFunc(aacNext)
		})))
	}
}
