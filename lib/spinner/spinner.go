package spinner

import (
	"fmt"
	"sync"
	"time"

	"github.com/hsk-kr/dev-setup-manager/lib/styles"
)

var frames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

type Spinner struct {
	message string
	stop    chan struct{}
	done    chan struct{}
	mu      sync.Mutex
}

func New(message string) *Spinner {
	return &Spinner{
		message: message,
		stop:    make(chan struct{}),
		done:    make(chan struct{}),
	}
}

func (s *Spinner) Start() {
	go func() {
		defer close(s.done)
		i := 0
		for {
			select {
			case <-s.stop:
				// Clear the spinner line
				fmt.Print("\r\033[K")
				return
			default:
				frame := frames[i%len(frames)]
				text := styles.LoadingText.Render(fmt.Sprintf("%s %s", frame, s.message))
				fmt.Printf("\r\033[K%s", text)
				i++
				time.Sleep(80 * time.Millisecond)
			}
		}
	}()
}

func (s *Spinner) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	select {
	case <-s.stop:
		// Already stopped
	default:
		close(s.stop)
	}
	<-s.done
}
