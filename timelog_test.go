package timelog

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTL(t *testing.T) {
	ctx := Start(context.Background(), "Level 1")
	tl := Get(ctx)

	ctx = Start(ctx, "Level 2.1")
	ctx = Finish(ctx)

	ctx = Start(ctx, "Level 2.2")
	Start(ctx, "Level 3") // Unfinished
	ctx = Finish(ctx)

	Finish(ctx)

	assert.NotNil(t, tl)
	assert.Len(t, tl.children, 2)
	assert.True(t, tl.start.Before(tl.finish))
	assert.Equal(t, tl.children[1].finish, tl.children[1].children[0].finish)

	action := tl.Analyze()

	rowHeaderRe := `\[\d+\.\d+%\] `
	rowFooterRe := `\t\t\(\S+: \S+ ⟼ \S+\)\n`
	assert.Regexp(t, "^"+rowHeaderRe+"Level 1"+rowFooterRe+
		"	"+rowHeaderRe+"Working"+rowFooterRe+
		"	"+rowHeaderRe+"Level 2.1"+rowFooterRe+
		"	"+rowHeaderRe+"Working"+rowFooterRe+
		"	"+rowHeaderRe+"Level 2.2"+rowFooterRe+
		"		"+rowHeaderRe+"Working"+rowFooterRe+
		"		"+rowHeaderRe+"Level 3"+rowFooterRe+
		"	"+rowHeaderRe+"Working"+rowFooterRe+
		"", action.String())
}
