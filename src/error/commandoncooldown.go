package customerror

import "fmt"

type CommandOnCooldownError struct {
	seconds_left int
}

func NewCommandOnCooldownErrorError(secondsLeft int) *CommandOnCooldownError {
	return &CommandOnCooldownError{seconds_left: secondsLeft}
}

func (ce *CommandOnCooldownError) Error() string {
	return fmt.Sprintf(":hourglass: This Command is on cooldown! (%vs left)", ce.seconds_left)
}
