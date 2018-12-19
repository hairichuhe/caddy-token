package token

import (
	"caddy-token/utils/caddyutil"
	"net/http"
	"strings"

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
		if r.URL.Path == upsrc {
			config = h.Config.Scope[upsrc]
			goto inScope
		}
	}

	for _, proxysrc := range h.Config.ProxyScopes {
		newStr := proxysrc[1:len(proxysrc)]
		end := strings.Index(newStr, "/")
		newStr = "/" + newStr[0:end] + "/"
		if r.URL.Path == proxysrc {
			return h.Next.ServeHTTP(w, r)
		} else {
			if httpserver.Path(r.URL.Path).Matches(newStr) {
				if caddyutil.Nopass(w, r) {
					return 0, nil
				}
			}
			return h.Next.ServeHTTP(w, r)
		}
	}
	return h.Next.ServeHTTP(w, r)
inScope:

	switch r.Method {
	case "POST":
		if caddyutil.Nopass(w, r) {
			return 0, nil
		}
		caddyutil.UpLoad(w, r, config)
		return 0, nil
	case "OPTIONS":
		w.Header().Set("Access-Control-Allow-Origin", config.AllowOrigin)
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.WriteHeader(200)
		w.Write([]byte("true"))
		return 0, nil
	case "DELETE":
		if caddyutil.Nopass(w, r) {
			return 0, nil
		}
		caddyutil.DelFile(w, r, config)
		return 0, nil
	default:
		return h.Next.ServeHTTP(w, r)
	}
}
