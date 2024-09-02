package closer

import (
	"os"
	"os/signal"
	"sync"

	"go.uber.org/zap"

	"github.com/Genvekt/cli-chat/libraries/logger/pkg/logger"
)

var globalCloser = New()

// Add is global function to add closing func to closer
func Add(f ...func() error) {
	globalCloser.Add(f...)
}

// Wait is global function to wait of closer closing
func Wait() {
	globalCloser.Wait()
}

// CloseAll is global function of closer closing
func CloseAll() {
	globalCloser.CloseAll()
}

// Closer collects functions that are executed on graceful shutdown
type Closer struct {
	mu    sync.Mutex
	once  sync.Once
	done  chan struct{}
	funcs []func() error
}

// New initialises closer and targets it on provided os signals
func New(sig ...os.Signal) *Closer {
	c := &Closer{done: make(chan struct{})}

	if len(sig) > 1 {
		go func() {
			ch := make(chan os.Signal, 1)

			signal.Notify(ch, sig...)
			<-ch
			signal.Stop(ch)

			c.CloseAll()
		}()
	}

	return c
}

// Add is used to add more functions to closer
func (c *Closer) Add(f ...func() error) {
	c.mu.Lock()
	c.funcs = append(c.funcs, f...)
	c.mu.Unlock()
}

// CloseAll executes all collected functions once
func (c *Closer) CloseAll() {
	c.once.Do(func() {
		defer close(c.done)

		c.mu.Lock()
		funcs := c.funcs
		c.funcs = nil
		c.mu.Unlock()

		errs := make(chan error, len(funcs))
		for _, f := range funcs {
			go func(f func() error) {
				errs <- f()
			}(f)
		}

		for i := 0; i < cap(errs); i++ {
			if err := <-errs; err != nil {
				logger.Error("closer func returned error", zap.Error(err))
			}
		}
	})
}

// Wait hangs until CloseAll is finished
func (c *Closer) Wait() {
	<-c.done
}
