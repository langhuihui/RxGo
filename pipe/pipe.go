package pipe

import . "github.com/langhuihui/RxGo/rx"

//Pipe 管道操作，可以将一组操作传入管道中，最后一个参数如果是Observer类型的参数的话就会激活事件流
func Pipe(source Observable, cbs ...Operator) Observable {
	return source.Pipe(cbs...)
}
