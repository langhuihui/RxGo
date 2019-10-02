package rx

//Subject 可以手动发送数据
func Subject(input chan interface{}) Observable {
	return FromChan(input).Share()
}

//FromSlice 把Slice转成Observable
func FromSlice(slice []interface{}) Observable {
	return func(sink *Observer) error {
		for _, data := range slice {
			sink.Next(data)
			if sink.IsStopped() {
				return sink.err
			}
		}
		return sink.err
		//sink.Complete()
	}
}

//Of 发送一系列值
func Of(array ...interface{}) Observable {
	return FromSlice(array)
}

//FromChan 把一个chan转换成事件流
func FromChan(source chan interface{}) Observable {
	return func(sink *Observer) error {
		for {
			select {
			case <-sink.stop:
				return sink.err
			case data, ok := <-source:
				if ok {
					sink.Next(data)
				} else {
					return sink.err
				}
			}
		}
	}
}

//Never 永不回答
func Never() Observable {
	return func(sink *Observer) error {
		<-sink.stop
		return sink.err
	}
}

//Empty 不会发送任何数据，直接完成
func Empty() Observable {
	return func(sink *Observer) error {
		return sink.err
	}
}

//Throw 直接抛出一个错误
func Throw(err error) Observable {
	return func(sink *Observer) error {
		return err
		//sink.OnNext(&Event{
		//	err:     err,
		//	Target: sink,
		//})
	}
}
