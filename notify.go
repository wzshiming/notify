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
func Once(signal os.Signal, fun func()) func() {
	return std.Once(signal, fun)
}

// OnSlice system signal callback.
func OnSlice(signals []os.Signal, fun func()) func() {
	return std.OnSlice(signals, fun)
}

// OnceSlice system signal callback.
func OnceSlice(signals []os.Signal, fun func()) func() {
	return std.OnceSlice(signals, fun)
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

func (n *notify) Once(sig os.Signal, fun func()) func() {
	return warpOnceFunc(n.once(sig, warpOnceFunc(fun)))
}

func (n *notify) OnceSlice(sigs []os.Signal, fun func()) func() {
	switch len(sigs) {
	case 0:
		return func() {}
	case 1:
		return n.Once(sigs[0], fun)
	default:
		offs := make([]func(), 0, len(sigs))
		off := warpOnceFunc(func() {
			for _, off := range offs {
				off()
			}
		})
		funAndOff := warpOnceFunc(func() {
			fun()
			off()
		})
		for _, sig := range sigs {
			offs = append(offs, n.on(sig, funAndOff))
		}
		return off
	}
}

func (n *notify) On(sig os.Signal, fun func()) func() {
	return warpOnceFunc(n.on(sig, fun))
}

func (n *notify) OnSlice(sigs []os.Signal, fun func()) func() {
	switch len(sigs) {
	case 0:
		return func() {}
	case 1:
		return n.On(sigs[0], fun)
	default:
		offs := make([]func(), 0, len(sigs))
		for _, sig := range sigs {
			offs = append(offs, n.on(sig, fun))
		}
		return warpOnceFunc(func() {
			for _, off := range offs {
				off()
			}
		})
	}
}

func (n *notify) on(sig os.Signal, fun func()) func() {
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

func (n *notify) once(sig os.Signal, fun func()) func() {
	off := func() {}
	off = n.on(sig, func() {
		fun()
		off()
	})
	return off
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
		n.step(sig)
	}
}

func (n *notify) step(sig os.Signal) {
	n.mut.Lock()
	funcs := n.event[sig]
	n.mut.Unlock()
	for _, fun := range funcs {
		fun()
	}
}

func warpOnceFunc(fun func()) func() {
	var once sync.Once
	return func() {
		once.Do(fun)
	}
}
