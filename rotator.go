package goolog2

import (
	"container/heap"
	"sync"
	"time"
)

// The rotator interface. The rotator is usually connected to one particular logger.
// If a rotator is added to logsRotatorStarter,  methods NeedRotate + Rotate are runed
// in separate goroutine. The start time of boths method is determined by the method GetNextCheckTime.
type LogRotator interface {
	// return if this file need rotation
	NeedRotate(timesrc TimeSource) bool
	// Rotate this file
	Rotate(timesrc TimeSource)
	// Return the time of next check. It is called after logsRotator.Add and after NeedRotate/Rotate
	GetNextCheckTime(timesrc TimeSource) time.Time
}

// FileHolder with the LogRotator
type RotatableFileHolder interface {
	LogRotator
	FileHolder
}

type logsRotatorStarter struct {
	timesrc     TimeSource
	nextCheck   time.Time
	rotators    *rotatorHeap // heap - first item is always the first scheduled ation
	destroyChan chan struct{}
	addNewChan  chan LogRotator
	once        sync.Once
}

func newRotators(timesrc TimeSource) logsRotatorStarter {
	return logsRotatorStarter{timesrc: timesrc, rotators: &rotatorHeap{}}
}

// Add the rotator to starter. See LogRotator inteface.
func (this *logsRotatorStarter) Add(rotator LogRotator) {
	this.once.Do(func() {
		this.nextCheck = this.timesrc.Now().Add(time.Hour)
		this.addNewChan = make(chan LogRotator)
		this.destroyChan = make(chan struct{})
		go this.mainThread()
	})
	this.addNewChan <- rotator
}

// See function AfterChangeMockedTime
func (this *logsRotatorStarter) OnMockedTimeChanged(wait bool) {
	if this.rotators == nil || len(*this.rotators) == 0 {
		return
	}
	this.Add(nil)
	if wait {
		// wait to next mainThread step
		for this.nextCheck.Before(this.timesrc.Now()) {
			time.Sleep(10 * time.Microsecond)
		}
	}
}

func (this *logsRotatorStarter) Destroy() {
	this.destroyChan <- struct{}{}
}

func (this *logsRotatorStarter) mainThread() {
	heap.Init(this.rotators)
	var timer *time.Timer
	for {
		timer = this.sleepTo(this.nextCheck, timer)
		select {
		case <-this.destroyChan:
			this.abortSleep(timer)
			return
		case newRotator := <-this.addNewChan:
			if newRotator != nil {
				checkTime := newRotator.GetNextCheckTime(this.timesrc)
				if checkTime.Before(this.nextCheck) {
					this.nextCheck = checkTime
				}
				heap.Push(this.rotators, &rotatorWithTime{rotator: newRotator, nextCheckTime: checkTime})
			}
		case <-timer.C:
			this.step()
		}
	}
}

func (this *logsRotatorStarter) sleepTo(awakeTime time.Time, timer *time.Timer) *time.Timer {
	// calculate duration
	now := this.timesrc.Now()
	duration := 0 * time.Second
	if now.Before(awakeTime) {
		duration = awakeTime.Sub(now)
	}
	if timer == nil {
		// first call - create timer
		return time.NewTimer(duration)
	}
	// next calls - reuse this timer
	this.abortSleep(timer)
	timer.Reset(duration)
	return timer
}

func (this *logsRotatorStarter) abortSleep(timer *time.Timer) {
	if timer != nil {
		timer.Stop()
		// read chan if it is possible
		select {
		case <-timer.C:
		default:
		}
	}
}

func (this *logsRotatorStarter) step() {
	if len(*this.rotators) == 0 {
		this.nextCheck = this.timesrc.Now().Add(time.Hour)
		return
	}
	first := (*this.rotators)[0]
	now := this.timesrc.Now()
	for !now.Before(first.nextCheckTime) {
		// rotate (if it is required)
		if first.rotator.NeedRotate(this.timesrc) {
			first.rotator.Rotate(this.timesrc)
		}
		first.nextCheckTime = first.rotator.GetNextCheckTime(this.timesrc)
		// reorder rotators
		heap.Fix(this.rotators, 0)
		first = (*this.rotators)[0]
	}
	this.nextCheck = first.nextCheckTime
}

type rotatorWithTime struct {
	rotator       LogRotator
	nextCheckTime time.Time
}

// It implements heap.Interface. The first item is always the first scheduled action.
type rotatorHeap []*rotatorWithTime

func (this *rotatorHeap) Len() int {
	return len(*this)
}

func (this *rotatorHeap) Less(i, j int) bool {
	return (*this)[i].nextCheckTime.Before((*this)[j].nextCheckTime)
}

func (this *rotatorHeap) Swap(i, j int) {
	(*this)[i], (*this)[j] = (*this)[j], (*this)[i]
}

func (this *rotatorHeap) Push(x interface{}) {
	item := x.(*rotatorWithTime)
	*this = append(*this, item)
}

func (this *rotatorHeap) Pop() interface{} {
	panic("unused")
}
