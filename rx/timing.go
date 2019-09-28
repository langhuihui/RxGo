package rx

import "time"

//Timeout 在一段时间后发送数据
func Timeout(duration time.Duration) Observable {
	return func(sink *Control) {
		timeout := time.After(duration)
		for {
			select {
			case observer, ok := <-sink.control:
				if !ok {
					return
				}
				sink.observer = observer
			case data := <-timeout:
				sink.Next(data)
				sink.Complete()
				return
			}
		}
	}
}

//Interval 按照一定频率持续发送数据
func Interval(duration time.Duration) Observable {
	return func(sink *Control) {
		interval := time.NewTicker(duration)
		i := 0
		defer interval.Stop()
		for {
			select {
			case observer, ok := <-sink.control:
				if !ok {
					return
				}
				sink.observer = observer
			case <-interval.C:
				i++
				sink.Next(i)
			}
		}
	}
}
