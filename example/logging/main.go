package main

import (
	"context"
	"time"

	"github.com/adipurnama/go-toolkit/log"
)

func main() {
	log.SetupWithLogfmtOutput(time.UTC)
	ctx := context.Background()

	// log.DebugCtx(ctx).Stack().Caller().AnErr("error", definitelyError()).Msg("debug message with error")
	// log.DebugCtx(ctx).Err(definitelyError()).Str("field_here", "whatever").Msg("debug message")
	log.DebugCtx(ctx).Str("field_here", "whatever").Msg("debug message - no error")

	// log.Debug().Stack().Caller().Err(definitelyError()).Msg("debug message with error")
	log.Debug().Err(definitelyError()).Str("field_here", "whatever").Msg("debug message")
	log.Debug().Str("field_here", "whatever").Msg("debug message - no error")
}
