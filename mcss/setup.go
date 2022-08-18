package mcss

import (
	"github.com/qmuntal/stateless"
)

// "enum" of valid triggers (actions)
type Trigger struct {
	start           string
	startSuccess    string
	startFailure    string
	startCancel     string
	stop            string
	stopSuccess     string
	stopFailure     string
	lock            string
	unlock          string
	becomeHealthy   string
	becomeUnhealthy string
}

var t = &Trigger{
	start:           "start",
	startSuccess:    "startSuccess",
	startFailure:    "startFailure",
	startCancel:     "startCancel",
	stop:            "stop",
	stopSuccess:     "stopSuccess",
	stopFailure:     "stopFailure",
	lock:            "lock",
	unlock:          "unlock",
	becomeHealthy:   "becomeHealthy",
	becomeUnhealthy: "becomeUnhealthy",
}

// "enum" of valid states
type State struct {
	stopped     string
	starting    string
	healthy     string
	unhealthy   string
	stopping    string
	lockStopped string
}

var s = &State{
	stopped:     "stopped",
	starting:    "starting",
	healthy:     "healthy",
	unhealthy:   "unhealthy",
	stopping:    "stopping",
	lockStopped: "lockStopped",
}

func MakeStateMachine() *stateless.StateMachine {
	// initial state is stopped and ready, from here we can start or lock it stopped
	mcss := stateless.NewStateMachine(s.stopped)
	mcss.Configure(s.stopped).
		Permit(t.start, s.starting).
		Permit(t.lock, s.lockStopped)

	// from lockStopped we can only release the lock
	mcss.Configure(s.lockStopped).
		Permit(t.unlock, s.stopped)

	// from a starting state we can succeed or fail, failure could happen from a process or user cancellation
	mcss.Configure(s.starting).
		Permit(t.startSuccess, s.healthy).
		Permit(t.startFailure, s.stopped).
		Permit(t.startCancel, s.stopped)

	// from healthy it might become unhealthy (but still running), or begin stopping
	mcss.Configure(s.healthy).
		Permit(t.stop, s.stopping).
		Permit(t.becomeUnhealthy, s.unhealthy)

	// from unhealthy we can become healthy (idk how) or begin to stopping
	mcss.Configure(s.unhealthy).
		Permit(t.becomeHealthy, s.healthy).
		Permit(t.stop, s.stopping)

	// from stopping I'm not sure what failure would look like or how to handle it, but "unhealthy" is
	// probably the correct state to go into
	mcss.Configure(s.stopping).
		Permit(t.stopSuccess, s.stopped).
		Permit(t.stopFailure, s.unhealthy)

	return mcss
}
