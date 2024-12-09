package closer

import (
	"log"
	"os"
	"os/signal"
	"sync"
)

var globalCloser = NewCloser()

// Closer — структура для работы с пулом функций, закрывающих ресурсы
type Closer struct {
	mu    sync.Mutex
	once  sync.Once
	done  chan struct{}
	funcs []func() error
}

// NewCloser принимает набор сигналов операционной системы, при срабатывании одного из них будет вызвана функция CloseAll
func NewCloser(sig ...os.Signal) *Closer {
	c := &Closer{done: make(chan struct{})}
	if len(sig) > 0 {
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

// CloseAll проходит по пулу функций и закрывает ресурсы
func CloseAll() {
	globalCloser.CloseAll()
}

// Add добавляет функции закрытия ресурса в пул функций
func Add(fn ...func() error) {
	globalCloser.Add(fn...)
}

// Wait блокирует выполнение пока канал не закроется в функции CloseAll
func Wait() {
	globalCloser.Wait()
}

// CloseAll проходит по пулу функций и закрывает ресурсы
func (c *Closer) CloseAll() {
	c.once.Do(func() {
		defer close(c.done)

		c.mu.Lock()
		funcs := c.funcs
		c.funcs = nil
		c.mu.Unlock()

		errs := make(chan error, len(funcs))
		for _, fn := range funcs {
			errs <- fn()
		}

		for i := 0; i < cap(errs); i++ {
			if err := <-errs; err != nil {
				log.Println("error on closing resource")
			}
		}
	})
}

// Add добавляет функцию закрытия ресурса в пул
func (c *Closer) Add(fn ...func() error) {
	c.mu.Lock()
	c.funcs = append(c.funcs, fn...)
	c.mu.Unlock()
}

// Wait блокирует выполнение пока канал не закроется в функции CloseAll
func (c *Closer) Wait() {
	<-c.done
}
