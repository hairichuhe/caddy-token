package token

import (
	"caddy-token/utils/caddyutil"
	"os"

	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/httpserver"
)

// Init initializes the plugin
func init() {
	caddy.RegisterPlugin("token", caddy.Plugin{
		ServerType: "http",
		Action:     Setup,
	})
}

// HandlerConfiguration is the result of directives found in a 'Caddyfile'.
//
// Can be modified at runtime, except for values that are marked as 'read-only'.
type HandlerConfiguration struct {
	// Prefixes on which Caddy activates this plugin (read-only).
	//
	// UpFile paths
	UpFileScopes []string

	// proxy paths
	ProxyScopes []string

	// Maps scopes (paths) to their own and potentially differently configurations.
	Scope map[string]*caddyutil.Config
}

// Setup parses the token and returns the middleware handler.
func Setup(c *caddy.Controller) error {
	config, err := parseCaddyConfig(c)
	if err != nil {
		return err
	}

	httpserver.GetConfig(c).AddMiddleware(func(next httpserver.Handler) httpserver.Handler {
		return &Handler{
			Next:   next,
			Config: *config,
		}
	})

	return nil
}

func parseCaddyConfig(c *caddy.Controller) (*HandlerConfiguration, error) {
	siteConfig := &HandlerConfiguration{
		UpFileScopes: make([]string, 0, 1),
		ProxyScopes:  make([]string, 0, 1),
		Scope:        make(map[string]*caddyutil.Config),
	}

	for c.Next() {
		config := caddyutil.DefaultConfig()

		scopes := c.RemainingArgs() // most likely only one path; but could be more
		if len(scopes) != 2 {
			return siteConfig, c.ArgErr()
		}
		siteConfig.UpFileScopes = append(siteConfig.UpFileScopes, scopes[0])
		siteConfig.ProxyScopes = append(siteConfig.ProxyScopes, scopes[1])

		for c.NextBlock() {
			key := c.Val()
			switch key {
			case "avatar_src":
				if !c.NextArg() {
					return siteConfig, c.ArgErr()
				}
				// must be a directory
				avatarToPath := c.Val()
				avatarfinfo, err := os.Stat(avatarToPath)
				if err != nil {
					return siteConfig, c.Err(err.Error())
				}
				if !avatarfinfo.IsDir() {
					return siteConfig, c.ArgErr()
				}
				config.AvatarSrc = avatarToPath
			case "file_src":
				if !c.NextArg() {
					return siteConfig, c.ArgErr()
				}
				// must be a directory
				fileToPath := c.Val()
				filefinfo, err := os.Stat(fileToPath)
				if err != nil {
					return siteConfig, c.Err(err.Error())
				}
				if !filefinfo.IsDir() {
					return siteConfig, c.ArgErr()
				}
				config.FileSrc = fileToPath
			case "allow_origin":
				if c.NextArg() {
					config.AllowOrigin = c.Val()
				}
			}
		}

		if config.AvatarSrc == "" {
			return siteConfig, c.Errf("请配置头像储存路径（“avatar_src”）！")
		}

		if config.FileSrc == "" {
			return siteConfig, c.Errf("请配置文件储存路径（“file_src”）！")
		}
		siteConfig.Scope[scopes[0]] = config
	}

	return siteConfig, nil
}
