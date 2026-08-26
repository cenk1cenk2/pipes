package cli

import (
	"sync"

	"github.com/cenk1cenk2/plumber/v6"
)

// Flags are registered from package level variables, which is why the registry is
// package level as well. The lock is here because a pipe may build its flags from
// more than one goroutine, not because the set is expected to change at runtime.
var (
	secretsLock sync.Mutex
	secrets     []*string
)

// MarkSecret records a flag destination as sensitive and hands it back, so it
// reads inline as the flag Destination. Validated masks whatever landed in it.
func MarkSecret(dst *string) *string {
	secretsLock.Lock()
	defer secretsLock.Unlock()

	secrets = append(secrets, dst)

	return dst
}

// Validated applies the defaults to the pipe, validates it, then keeps every
// value marked with MarkSecret out of the logs. Pipes call this instead of
// Validate so a new secret flag is masked by the act of declaring it.
func Validated(p *plumber.Plumber, pipe any) error {
	if err := p.Validate(pipe); err != nil {
		return err
	}

	secretsLock.Lock()
	defer secretsLock.Unlock()

	for _, secret := range secrets {
		if *secret != "" {
			p.AppendSecrets(*secret)
		}
	}

	return nil
}
