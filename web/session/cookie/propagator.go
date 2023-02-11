package cookie

import "net/http"

type PropagatorOption func(propagator *Propagator)

func WithCookieOption(opt func(cookie *http.Cookie)) PropagatorOption {
	return func(propagator *Propagator) {
		propagator.cookieOption = opt
	}
}

func NewPropagator(cookieName string, opts PropagatorOption) *Propagator {
	res := &Propagator{
		cookieName: cookieName,
		cookieOption: func(c *http.Cookie) {
		},
	}
	opts(res)
	return res
}

type Propagator struct {
	cookieName   string
	cookieOption func(c *http.Cookie)
}

func (p *Propagator) Inject(id string, writer http.ResponseWriter) error {
	cookie := &http.Cookie{
		Name:  p.cookieName,
		Value: id,
	}
	p.cookieOption(cookie)
	http.SetCookie(writer, cookie)
	return nil
}

func (p *Propagator) Extract(req *http.Request) (string, error) {
	cookie, err := req.Cookie(p.cookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func (p *Propagator) Remove(writer http.ResponseWriter) error {
	cookie := &http.Cookie{
		Name:   p.cookieName,
		MaxAge: -1,
	}
	p.cookieOption(cookie)
	http.SetCookie(writer, cookie)
	return nil
}
