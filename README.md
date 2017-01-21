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

	1.000453077s	Level 1
		1.277µs	Working
		800.303067ms	Level 2.1
			500.166711ms	Working
			300.122218ms	Level 3
			14.138µs	Working
		843ns	Working
		200.143999ms	Level 2.2
		3.891µs	Working
