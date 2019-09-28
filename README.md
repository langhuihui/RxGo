# RxGo 非官方实现版本

目标：代码精简，可读性强，实现优雅，占用系统资源低，性能强

## 使用方法
### 链式调用方式
```go
import (
    "github.com/langhuihui/RxGo/rx"
)
func main(){
    err := rx.Of(1, 2, 3, 4).Take(2).SubscribeSync(func(data interface{}) {
        
    })
}
```
### 管道模式
```go
import (
    "github.com/langhuihui/RxGo/rx"
)
func main(){
    err := rx.Of(1, 2, 3, 4).Pipe(rx.Skip(1),rx.Take(2)).SubscribeSync(func(data interface{}) {
        
    }))
}
```

# 设计思想
## 总体方案
### 可观察对象（事件源）Observable
Observable 被定义成为一个函数，该函数含有一个类型为*Control的参数。
```go
type Observable func(*Control)
```
任何事件源都是这样的一个函数，当调用该函数即意味着**订阅**了该事件源，入参为一个控制器，具体功能见下面

### 控制器对象（由观察者提供传入可观察者）Control
```go
type Control struct {
	observer Observer
	control  chan Observer
	closed   bool
}
```
该控制器为一个结构体，其中observer记录了当前的observer，control为一个channel用于取消订阅以及修改observer的功能，closed代表已经被取消订阅了
在任何时候，如果关闭了control这个channel，就意味着**取消订阅**。如果向control写入新的Observer，意味着用于接受数据的函数改变了，这个功能是为了让实现更灵活，更简洁而设计
Control对象为Observable和事件处理逻辑共同持有，是二者沟通的桥梁

### 观察者Observer
```go
type Observer func(interface{}, error)
```
观察者仅仅只是一个函数，当Observable数据推送到Observer中时，即调用了该函数
如果传入的第二个参数等于特殊的Complete变量，即意味着**事件流完成**
如果传入的第二个参数不等于nil，即意味着**事件流异常**
其他情况，意味着**事件流事件**

### 数据传输控制
为了方便控制内部数据的流动，采用读取control一次，推送一次的方式
大致示意如下
```go
observer:=<-control
observer.Next(data)
```
为了高效灵活的控制，我们在Control的实现中加入了自动写入上一个Observer的方式，除非channel中有来自下游提供的新的Observer
，具体实现方法见下面

## 设计细节 
### 创建控制器
```go
func NewControl(observer Observer) *Control {
	return &Control{
		observer: observer,
		control:  make(chan Observer, 1),
		closed:   false,
	}
}
```
说明：control为带缓冲的channel，这里的设计是为了防止死锁
因为大部分更换Observer的操作都是发生在Observer函数内部，即事件处理逻辑内部，此时如果执行control<-阻塞，会导致Observable内部阻塞

### 推送数据
最精髓的设计原理就在下面这个函数中
```go
func (c *Control) Push(data interface{}, err ...error) {
	if len(err) == 1 && err[0] != nil {
		c.Stop()
		c.observer.Error(err[0])
	} else {
		c.observer.Next(data)
		if len(c.control) == 0 {
			c.control <- c.observer
		}
	}
}
```
Observable可以调用push函数推送数据，如果control中没有新的Observer，则写入上一个Observer。这里是关键逻辑
原本可以在下游不停的向control写入Observer，但由于大多数情况是相同的Observer，所以会导致代码冗长，所以在Control中缓存了一个Observer，用于重复使用