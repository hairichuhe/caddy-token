package token

import (
	"caddy-token/utils/caddyutil"
	"fmt"
	"net/http"

	"github.com/mholt/caddy/caddyhttp/httpserver"
)

type Handler struct {
	Next   httpserver.Handler
	Config HandlerConfiguration
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	var (
		config *caddyutil.Config
	)

	for _, upsrc := range h.Config.UpFileScopes {
		if httpserver.Path(r.URL.Path).Matches(upsrc) {
			config = h.Config.Scope[upsrc]
			goto inScope
		}
	}

	for _, proxysrc := range h.Config.ProxyScopes {
		fmt.Println(proxysrc)
		fmt.Println(r.URL.Path)
		if httpserver.Path(r.URL.Path).Matches(proxysrc) {

			if caddyutil.Nopass(w, r) {
				return http.StatusForbidden, nil
			}
			return h.Next.ServeHTTP(w, r)
		}
	}
	fmt.Println("到这里了！")
	return h.Next.ServeHTTP(w, r)
inScope:

	switch r.Method {
	case "POST":
		if caddyutil.Nopass(w, r) {
			return http.StatusForbidden, nil
		}
		caddyutil.UpLoad(w, r, config)
		return 0, nil
	default:
		return h.Next.ServeHTTP(w, r)
	}
}
