package rx

//Subject 可以手动发送数据
func Subject(input chan interface{}) Observable {
	return FromChan(input).Share()
}

//FromSlice 把Slice转成Observable
func FromSlice(slice []interface{}) Observable {
	return func(sink *Control) {
		for _, data := range slice {
			sink.Next(data)
			if !sink.Check() {
				return
			}
		}
		sink.Complete()
	}
}

//Of 发送一系列值
func Of(array ...interface{}) Observable {
	return FromSlice(array)
}

//FromChan 把一个chan转换成事件流
func FromChan(source chan interface{}) Observable {
	return func(sink *Control) {
		for {
			select {
			case observer, ok := <-sink.control:
				if !ok {
					return
				}
				sink.observer = observer
			case data, ok := <-source:
				if ok {
					sink.Next(data)
				} else {
					sink.Complete()
				}
			}
		}
	}
}

//Never 永不回答
func Never() Observable {
	return func(sink *Control) {
		for sink.Check() {

		}
	}
}

//Empty 直接完成
func Empty() Observable {
	return func(sink *Control) {
		sink.Complete()
	}
}

func Throw(err error) Observable {
	return func(sink *Control) {
		sink.Push(&Event{
			err:     err,
			control: sink,
		})
	}
}
