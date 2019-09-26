package rx

//FromSlice 把Slice转成Observable
func FromSlice(slice []interface{}) Observable {
	return func(sink Sink) {
		var observer Observer
		var ok bool
		for _, data := range slice {
			if observer, ok = <-sink; ok {
				observer.Next(data)
			} else {
				return
			}
		}
		if observer != nil {
			observer.Complete()
		}
	}
}

//Of 发送一系列值
func Of(array ...interface{}) Observable {
	return FromSlice(array)
}

//FromChan 把一个chan转换成事件流
func FromChan(source chan interface{}) Observable {
	return func(sink Sink) {
		for observer := range sink {
			if data, ok := <-source; ok {
				observer.Next(data)
			} else {
				observer.Complete()
				break
			}
		}
	}
}

//Never 永不回答
func Never() Observable {
	return func(sink Sink) {
		for range sink {
		}
	}
}

//Empty 直接完成
func Empty() Observable {
	return func(sink Sink) {
		for observer := range sink {
			observer.Complete()
			break
		}
	}
}
