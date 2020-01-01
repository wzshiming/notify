package notify

import (
	"os"
	"os/signal"
	"sync"
)

var std = newNotify()

// On system signal callback.
func On(signal os.Signal, fun func()) func() {
	return std.On(signal, fun)
}

// Once system signal callback.
func Once(signal os.Signal, fun func()) {
	std.Once(signal, fun)
}

type notify struct {
	ch    chan os.Signal
	size  int
	event map[os.Signal]map[int]func()
	mut   sync.Mutex
}

func newNotify() *notify {
	return &notify{
		event: map[os.Signal]map[int]func(){},
	}
}

func (n *notify) Once(sig os.Signal, fun func()) {
	off := func() {}
	off = n.On(sig, func() {
		fun()
		off()
	})
}

func (n *notify) On(sig os.Signal, fun func()) func() {
	n.mut.Lock()
	defer n.mut.Unlock()
	_, ok := n.event[sig]
	if !ok {
		n.init(sig)
	}
	n.size++
	i := n.size
	n.event[sig][i] = fun
	return func() {
		n.off(sig, i)
	}
}

func (n *notify) off(sig os.Signal, i int) {
	n.mut.Lock()
	defer n.mut.Unlock()
	_, ok := n.event[sig]
	if !ok {
		return
	}
	delete(n.event[sig], i)
	if len(n.event[sig]) == 0 {
		delete(n.event, sig)
		n.reset()
	}
}

func (n *notify) init(sig os.Signal) {
	if n.ch == nil {
		n.ch = make(chan os.Signal)
		go n.run()
	}
	n.event[sig] = map[int]func(){}
	n.reset()
}

func (n *notify) reset() {
	if len(n.event) == 0 {
		signal.Stop(n.ch)
		close(n.ch)
		n.ch = nil
		return
	}
	sigs := make([]os.Signal, 0, len(n.event))
	for sig := range n.event {
		sigs = append(sigs, sig)
	}
	signal.Notify(n.ch, sigs...)
}

func (n *notify) run() {
	for sig := range n.ch {
		n.on(sig)
	}
}

func (n *notify) on(sig os.Signal) {
	n.mut.Lock()
	defer n.mut.Unlock()
	for _, fun := range n.event[sig] {
		fun()
	}
}
