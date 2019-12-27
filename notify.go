package notify

import (
	"os"
	"os/signal"
	"sync"
)

var std = newNotify()

// On system signal callback.
func On(signal os.Signal, fun func()) {
	std.On(signal, fun)
}

type notify struct {
	ch    chan os.Signal
	event map[os.Signal][]func()
	once  sync.Once
}

func newNotify() *notify {
	return &notify{
		event: map[os.Signal][]func(){},
	}
}

func (n *notify) On(signal os.Signal, fun func()) {
	_, ok := n.event[signal]
	if !ok {

		n.reset()
	}
	n.event[signal] = append(n.event[signal], fun)
}

func (n *notify) reset() {
	if n.ch == nil {
		n.once.Do(func() {
			n.ch = make(chan os.Signal)
			go n.run()
		})
	}

	sigs := make([]os.Signal, 0, len(n.event))
	for sig := range n.event {
		sigs = append(sigs, sig)
	}
	signal.Notify(n.ch, sigs...)
}

func (n *notify) run() {
	for signal := range n.ch {
		n.on(signal)
	}
}

func (n *notify) on(signal os.Signal) {
	for _, fun := range n.event[signal] {
		fun()
	}
}
