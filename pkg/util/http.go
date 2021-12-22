package util

import "net/http"

const terminal = "terminal"

// GetTerminal Fetch the security token set by the client.
func GetTerminal(r *http.Request) (terminal string) {
	terminal = r.Header.Get(terminal)
	if terminal != "" {
		return terminal
	}
	terminal = r.Form.Get(terminal)
	if terminal != "" {
		return terminal
	}
	return ""
}

func GetIpAddr(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}
