package imports

import (
	_ "gin/app/command"
	_ "gin/pkg/cli/db"
	_ "gin/pkg/cli/event"
	_ "gin/pkg/cli/grpc"
	_ "gin/pkg/cli/job"
	_ "gin/pkg/cli/jwt"
	_ "gin/pkg/cli/make"
	_ "gin/pkg/cli/mcp"
	_ "gin/pkg/cli/permission"
	_ "gin/pkg/cli/queue"
	_ "gin/pkg/cli/route"
)
