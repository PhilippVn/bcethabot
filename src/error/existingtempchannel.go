package customerror

import (
	"fmt"
)

type ExistingTempChannelError struct {
	channelName         string
	secondsSinceCreated int
}

func NewExistingTempChannelError(channelname string, secondsincecreated int) *ExistingTempChannelError {
	return &ExistingTempChannelError{
		channelName:         channelname,
		secondsSinceCreated: secondsincecreated,
	}
}

func (ce *ExistingTempChannelError) Error() string {
	return fmt.Sprintf(":x: You already created a temp channel: `%s` %v seconds ago", ce.channelName, ce.secondsSinceCreated)
}
