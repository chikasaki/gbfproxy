package mypool

import "time"

// OneTimePool : only put once and take once
type OneTimePool struct {
	eChan   chan interface{}
	initial func() interface{}
}

func NewOneTimePool(fullSize int, New func() interface{}) *OneTimePool {
	ans := &OneTimePool{
		eChan:   make(chan interface{}, fullSize),
		initial: New,
	}
	for fullSize > 0 {
		go func() {
			for {
				if e := New(); e != nil {
					ans.eChan <- e
				} else {
					//error, sleep to wait resource release
					time.Sleep(100 * time.Millisecond)
				}
			}
		}()
		fullSize--
	}
	return ans
}

func (p *OneTimePool) Take() interface{} {
	select {
	case ans := <-p.eChan:
		return ans
	default:
		return p.initial()
	}
}
