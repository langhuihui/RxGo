# RxGo 非官方实现版本
[![Build Status](https://travis-ci.org/langhuihui/RxGo.svg?branch=master)](https://travis-ci.org/langhuihui/RxGo)
[![Go Report Card](https://goreportcard.com/badge/github.com/langhuihui/RxGo)](https://goreportcard.com/report/github.com/langhuihui/RxGo)
[![codecov](https://codecov.io/gh/langhuihui/RxGo/branch/master/graph/badge.svg)](https://codecov.io/gh/langhuihui/RxGo)

目标：代码精简，设计精妙，可读性强，实现优雅，占用系统资源低，性能强

**每一行代码都是深思熟虑……**

## 已实现的功能

### Observable

- [x] FromSlice
- [x] FromChan
- [x] Of 
- [x] Subject
- [x] Timeout
- [x] Interval
- [x] Merge
- [x] Concat
- [x] CombineLatest

### Operator

- [x] Take 
- [x] TakeWhile
- [x] TakeUntil
- [x] Skip 
- [x] SkipWhile 
- [x] SkipUntil
- [x] IgnoreElements
- [x] Share
- [x] StartWith

## 使用方法
### 链式调用方式
```go
import (
    . "github.com/langhuihui/RxGo/rx"
)
func main(){
    err := Of(1, 2, 3, 4).Take(2).Subscribe(ObserverFunc(func(event *Event) {
        
    }))
}
```
### 管道模式
```go
import (
    . "github.com/langhuihui/RxGo/rx"
    . "github.com/langhuihui/RxGo/pipe"
)
func main(){
    err := Of(1, 2, 3, 4).Pipe(Skip(1),Take(2)).Subscribe(ObserverFunc(func(event *Event) {
        
    }))
}
```
管道模式相比链式模式，具有操作符**可扩展性**，用户可以按照规则创建属于自己的操作符
```go
type Operator func(Observable) Observable
```
操作符只需要返回Operator这个类型即可

# 设计思想
## 总体方案
### 可观察对象（事件源）Observable
Observable 被定义成为一个函数，该函数含有一个类型为*Control的参数。
```go
type Observable func(*Control) error
```
任何事件源都是这样的一个函数，当调用该函数即意味着**订阅**了该事件源，入参为一个控制器，具体功能见下面

如果该函数返回nil，即意味着**事件流完成**

否则意味着**事件流异常**

### 控制器对象（由观察者提供传入可观察者）Control
```go
type Stop chan bool
type Control struct {
	observer Observer //缓存当前的Observer，后续可以被替换
	stop     Stop     //取消订阅的信号，只用来close
    err      error    //缓存当前错误
}
```
该控制器为一个结构体，其中observer记录了当前的observer，

在任何时候，如果关闭了stop这个channel，就意味着**取消订阅**。

由于Channel的close可以引发所有读取该Channel的阻塞行为唤醒，所以可以在不同层级复用该channel

并且，由于已经close的channel可以反复读取以取得是否close的状态信息，所以不需要再额外记录

Control对象为Observable和事件处理逻辑共同持有，是二者沟通的桥梁

### 观察者Observer
```go
type Event struct {
    data    interface{}
    control *Control
}
Observer interface {
    OnNext(*Event)
 }
```
观察者是一个接口，实现OnNext函数，当Observable数据推送到Observer中时，即调用了该函数

control属性用于存储当前发送事件的Control对象

这样做的好处是可以实现不同的观察者，比如函数或者channel
```go
type(
    ObserverFunc func(*Event)
    ObserverChan chan *Event
)
func (observer ObserverFunc) OnNext(event *Event) {
	observer(event)
}
func (observer ObserverChan) OnNext(event *Event) {
	observer <- event
}
```
## 设计细节 
未完待续