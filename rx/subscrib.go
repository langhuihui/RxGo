package rx

//Subscribe 对外同步订阅Observable，可以在收到的事件中访问Observer的Dispose函数来终止事件流
func (ob Observable) Subscribe(onNext NextHandler) error {
	return ob(&Observer{next: onNext})
}

//SubscribeAsync 对外异步订阅模式，返回一个可用于终止事件流的控制器,可以在未收到数据时也能终止事件流
func (ob Observable) SubscribeAsync(onNext NextHandler, onError func(error), onComplete func()) *Observer {
	source := &Observer{next: onNext}
	go func() {
		err := ob(source)
		if !source.disposed {
			if err != nil {
				onError(err)
			} else {
				onComplete()
			}
		}
	}()
	return source
}
