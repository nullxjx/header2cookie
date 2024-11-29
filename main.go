package header2cookie

import (
	"context"
	"fmt"
	"net/http"
)

// Config the plugin configuration.
type Config struct {
	Cookie []string `json:"cookie,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Cookie: []string{},
	}
}

type CookieManager struct {
	next   http.Handler
	Config *Config
	name   string
}

func (c *CookieManager) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	fmt.Printf("header2cookie receive request, cookie: %v, len: %v\n", c.Config.Cookie, len(c.Config.Cookie))
	// 把header中的一些值设置到cookie里面
	for _, key := range c.Config.Cookie {
		if value := req.Header.Get(key); value != "" {
			http.SetCookie(rw, &http.Cookie{
				Name:     key,
				Value:    value,
				Path:     "/",
				HttpOnly: true,
			})
			// 手动将新设置的 Cookie 添加到请求中
			req.AddCookie(&http.Cookie{
				Name:     key,
				Value:    value,
				Path:     "/",
				HttpOnly: true,
			})
		}
	}
	fmt.Printf("request cookie: %v\n", req.Cookies())
	// 继续处理下一个 Handler
	c.next.ServeHTTP(rw, req)
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &CookieManager{
		Config: config,
		next:   next,
		name:   name,
	}, nil
}
