package mutagen

import (
	"context"
	"errors"

	"github.com/docker/compose/v2/pkg/progress"
)

// statusUpdater provides an adapter for feeding Mutagen-related events to the
// Compose progress writer. It also implements the
// github.com/mutagen-io/mutagen/pkg/prompting.Prompter interface to provide
// message-only prompting.
type statusUpdater struct {
	// writer is the underlying Compose progress writer.
	writer progress.Writer
	// eventID is the identifier to use for events.
	eventID string
}

// newStatusUpdater extracts the Compose progress writer from the specified
// context and constructs a new statusUpdater.
func newStatusUpdater(ctx context.Context, eventID string) *statusUpdater {
	return &statusUpdater{writer: progress.ContextWriter(ctx), eventID: eventID}
}

// working registers a normal working event.
func (u *statusUpdater) working(description string) {
	u.writer.Event(progress.NewEvent(u.eventID, progress.Working, description))
}

// error registers an error event.
func (u *statusUpdater) error(err error) {
	u.writer.Event(progress.NewEvent(u.eventID, progress.Error, "Error: "+err.Error()))
}

// done registers a done event.
func (u *statusUpdater) done(description string) {
	u.writer.Event(progress.NewEvent(u.eventID, progress.Done, description))
}

// Message implements
// github.com/mutagen-io/mutagen/pkg/prompting.Prompter.Message.
func (u *statusUpdater) Message(message string) error {
	u.working(message)
	return nil
}

// Prompt implements
// github.com/mutagen-io/mutagen/pkg/prompting.Prompter.Prompt.
func (u *statusUpdater) Prompt(_ string) (string, error) {
	return "", errors.New("prompting not supported")
}
