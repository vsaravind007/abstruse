package server

import (
	"net/http"

	"github.com/bleenco/abstruse/pkg/assetfs"
	"github.com/bleenco/abstruse/server/api"
	"github.com/bleenco/abstruse/server/api/builds"
	"github.com/bleenco/abstruse/server/api/integration"
	"github.com/bleenco/abstruse/server/api/providers/github"
	"github.com/bleenco/abstruse/server/api/repos"
	"github.com/bleenco/abstruse/server/api/setup"
	"github.com/bleenco/abstruse/server/api/teams"
	"github.com/bleenco/abstruse/server/api/user"
	"github.com/bleenco/abstruse/server/api/workers"
	"github.com/bleenco/abstruse/server/websocket"
	"github.com/julienschmidt/httprouter"
)

// Router represents main HTTP router.
type Router struct {
	*httprouter.Router
}

// NewRouter returns new instance of main HTTP router.
func NewRouter() *Router {
	router := &Router{httprouter.New()}
	router.initWebsocket()
	router.initAPI()
	router.initUI()

	return router
}

func (r *Router) initAPI() {
	r.Router.GET("/api/setup/ready", setup.ReadyHandler)
	r.Router.POST("/api/user/login", user.LoginHandler)
	r.Router.GET("/api/version", api.AuthorizationMiddleware(FindVersionHandler))
	r.Router.GET("/api/integrations", api.AuthorizationMiddleware(integration.FindIntegrationsHandler))
	r.Router.GET("/api/integrations/:id", api.AuthorizationMiddleware(integration.FindIntegrationHandler))
	r.Router.GET("/api/integrations/:id/update", api.AuthorizationMiddleware(integration.UpdateGitHubIntegration))
	r.Router.GET("/api/integrations/:id/repos", api.AuthorizationMiddleware(integration.FetchIntegrationRepositoriesHandler))
	r.Router.POST("/api/integrations/github/add", api.AuthorizationMiddleware(integration.AddGitHubIntegration))
	r.Router.POST("/api/integrations/github/import/:id", api.AuthorizationMiddleware(github.ImportRepositoryHandler))
	r.Router.GET("/api/workers", api.AuthorizationMiddleware(workers.FindAllHandler))
	r.Router.GET("/api/repositories", api.AuthorizationMiddleware(repos.FetchRepositoriesHandler))
	r.Router.GET("/api/repositories/:id", api.AuthorizationMiddleware(repos.FetchRepositoryHandler))
	r.Router.GET("/api/repositories/:id/hooks", api.AuthorizationMiddleware(github.ListHooksHandler))
	r.Router.POST("/api/repositories/:id/hooks", api.AuthorizationMiddleware(github.CreateHookHandler))
	r.Router.POST("/api/builds/trigger", api.AuthorizationMiddleware(builds.TriggerBuildHandler))
	r.Router.GET("/api/teams", api.AuthorizationMiddleware(teams.FetchTeamsHandler))
	r.Router.GET("/api/teams/:id", api.AuthorizationMiddleware(teams.FetchTeamHandler))
	r.Router.POST("/api/teams/:id", api.AuthorizationMiddleware(teams.SaveTeamHandler))
	r.Router.GET("/api/users/personal", api.AuthorizationMiddleware(user.FetchPersonalInfo))
}

func (r *Router) initUI() {
	r.Router.NotFound = http.FileServer(&assetfs.AssetFS{
		Asset:    Asset,
		AssetDir: AssetDir,
		Prefix:   "web/abstruse/dist",
		Fallback: "index.html",
	})
}

func (r *Router) initWebsocket() {
	r.Router.GET("/ws", websocket.UpstreamHandler("127.0.0.1:7100"))
}
