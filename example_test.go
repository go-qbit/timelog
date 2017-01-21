package timelog_test

import (
	"context"
	"os"
	"time"

	"github.com/go-qbit/timelog"
)

func ExampleAction_Print() {
	// Initialize TimeLog
	ctx := timelog.Start(context.Background(), "Level 1")

	// Do some actions
	Action1(ctx)
	Action2(ctx)

	// Finalize TimeLog
	timelog.Finish(ctx)

	// Print information
	timelog.Get(ctx).Analyze().Print(os.Stdout, "\t")
}

func Action1(ctx context.Context) {
	ctx = timelog.Start(ctx, "Level 2.1")
	defer timelog.Finish(ctx)

	time.Sleep(500 * time.Millisecond)

	Action3(ctx)
}

func Action2(ctx context.Context) {
	ctx = timelog.Start(ctx, "Level 2.2")
	defer timelog.Finish(ctx)

	time.Sleep(200 * time.Millisecond)
}

func Action3(ctx context.Context) {
	ctx = timelog.Start(ctx, "Level 3")
	defer timelog.Finish(ctx)

	time.Sleep(300 * time.Millisecond)
}
