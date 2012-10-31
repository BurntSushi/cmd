package cmd

import (
	"os/exec"
	"runtime"
	"sync"
)

// Commands allows any kind of command with a "Run() error" method to be used
// with the pool. (i.e., you aren't forced to use this packages Command type.)
type Commander interface {
	Run() error
}

// Commands is a list of values that implement the Commander interface.
// This is used as the list of commands to be executed in a pool.
type Commands []Commander

// NewCommands is a convenience function for creating a list of Commanders from
// a list of *Command.
func NewCommands(cmds []*Command) Commands {
	lst := make([]Commander, len(cmds))
	for i, cmd := range cmds {
		lst[i] = Commander(cmd)
	}
	return lst
}

// NewCmds is a convenience function for creating a list of Commanders from
// a list of *exec.Cmd.
func NewCmds(cmds []*exec.Cmd) Commands {
	lst := make([]Commander, len(cmds))
	for i, cmd := range cmds {
		lst[i] = Commander(cmd)
	}
	return lst
}

// RunMany creates a pool with a number of workers specified by "workers".
// If "workers" is less than 1, then the value of GOMAXPROCS is used.
// Every command in "cmds" is executed once by a single worker.
// A list of errors corresponding to the list of 'cmds' is returned, where the 
// length of the list of errors is always equivalent to the length of 'cmds'.
//
// A convenient way to use this method, given a list of *Command:
//
//	errs := NewCommands(commands).RunMany(0)
func (cmds Commands) RunMany(workers int) []error {
	if workers < 1 {
		workers = runtime.GOMAXPROCS(0)
	}
	errs := make([]error, len(cmds))
	jobs := make(chan int, workers)
	wg := new(sync.WaitGroup)

	for i := 0; i < workers; i++ {
		go func() {
			wg.Add(1)
			defer wg.Done()

			for job := range jobs {
				if err := cmds[job].Run(); err != nil {
					errs[job] = err
				}
			}
		}()
	}
	for i := range cmds {
		jobs <- i
	}
	close(jobs)

	wg.Wait()
	return errs
}
