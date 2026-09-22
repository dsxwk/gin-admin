package cli

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/progress"
)

// StartProgress 启动命令进度
func StartProgress(message string, total int) (*progress.Tracker, <-chan struct{}) {
	writer := progress.NewWriter()
	writer.SetAutoStop(true)
	writer.SetOutputWriter(os.Stdout)
	writer.SetStyle(progress.StyleBlocks)

	tracker := &progress.Tracker{
		Message: message,
		Total:   int64(total),
	}
	writer.AppendTracker(tracker)

	done := make(chan struct{})
	go func() {
		writer.Render()
		close(done)
	}()

	return tracker, done
}

// StopProgress 停止命令进度
func StopProgress(tracker *progress.Tracker, done <-chan struct{}, success bool) {
	if success {
		tracker.MarkAsDone()
	} else {
		tracker.MarkAsErrored()
	}
	<-done
}
