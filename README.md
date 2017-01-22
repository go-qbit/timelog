# timelog

[![Build Status](https://travis-ci.org/go-qbit/timelog.svg?branch=master)](https://travis-ci.org/go-qbit/timelog)

## Installation
    go get github.com/go-qbit/timelog

## Usage
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
    
The code will produce:

	[100.0000%] Level 1		(1.000399975s: 0s ⟼ 1.000399975s)
		[0.0001%] Working		(824ns: 0s ⟼ 824ns)
		[79.9963%] Level 2.1		(800.282978ms: 824ns ⟼ 800.283802ms)
			[62.4931%] Working		(500.121338ms: 0s ⟼ 500.121338ms)
			[37.5065%] Level 3		(300.157855ms: 500.121338ms ⟼ 800.279193ms)
			[0.0005%] Working		(3.785µs: 800.279193ms ⟼ 800.282978ms)
		[0.0001%] Working		(765ns: 800.283802ms ⟼ 800.284567ms)
		[20.0032%] Level 2.2		(200.11168ms: 800.284567ms ⟼ 1.000396247s)
		[0.0004%] Working		(3.728µs: 1.000396247s ⟼ 1.000399975s)

