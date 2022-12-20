package curvebs

import (
	"github.com/opencurve/curve-manager/api/curvebs/core"
	"github.com/opencurve/curve-manager/api/curvebs/manager"
	"github.com/opencurve/curve-manager/api/curvebs/user"
	"github.com/opencurve/pigeon"
)

func NewServer() *pigeon.HTTPServer {
	server := pigeon.NewHTTPServer("curvebs")
	server.Initer(core.Init)
	server.Route("/curvebs",
		core.Rewrite,
		manager.Entrypoint,
		user.Entrypoint)
	server.DefaultRoute(core.Default)
	return server
}

func main() {
	pigeon.Serve(NewServer())
}
