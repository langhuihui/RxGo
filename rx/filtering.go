package rx

//Take 获取最多count数量的事件，然后完成
func (ob Observable) Take(count int) Observable {
	return func(sink Sink) {
		remain := count
		source := ob.NewSink()
		defer close(source)
		for observer := range sink {
			if remain--; remain < 0 {
				observer.Complete()
				return
			}
			source <- observer
		}
	}
}

//Skip 跳过若干个数据
func (ob Observable) Skip(count int) Observable {
	return func(sink Sink) {
		source := ob.NewSink()
		defer close(source)
		for remain := count; remain > 0; remain-- {
			source <- func(data interface{}, err error) {
				if err != nil {
					if observer, ok := <-sink; ok {
						observer.Error(err)
					}
					return
				}
			}
		}
		for observer := range sink {
			source <- observer
		}
	}
}
