package session

import "goStudy/web"

type Manager struct {
	Store
	Propagator
	SessCtxKey string
}

func (m *Manager) InitSession(ctx *web.Context, id string) (Session, error) {
	session, err := m.Generate(ctx.Req.Context(), id)
	if err != nil {
		return nil, err
	}
	if err = m.Inject(id, ctx.Resp); err != nil {
		return nil, err
	}
	return session, err
}

// 为了加速获得session 把session 缓存到了ctx中
func (m *Manager) GetSession(ctx *web.Context) (Session, error) {
	if ctx.UserValues == nil {
		ctx.UserValues = make(map[string]any, 1)
	}
	val, ok := ctx.UserValues[m.SessCtxKey]
	if ok {
		return val.(Session), nil
	}
	id, err := m.Extract(ctx.Req)
	if err != nil {
		return nil, err
	}
	session, err := m.GET(ctx.Req.Context(), id)
	if err != nil {
		return nil, err
	}
	ctx.UserValues[m.SessCtxKey] = session
	return session, nil
}

func (m *Manager) RefreshSession(ctx *web.Context) (Session, error) {
	session, err := m.GetSession(ctx)
	if err != nil {
		return nil, err
	}
	//刷新存储时间
	if err := m.Refresh(ctx.Req.Context(), session.ID()); err != nil {
		return nil, err
	}

	//注入http resp
	if err := m.Inject(session.ID(), ctx.Resp); err != nil {
		return nil, err
	}
	return session, nil
}

func (m *Manager) RemoveSession(ctx *web.Context) error {
	sess, err := m.GetSession(ctx)
	if err != nil {
		return err
	}
	//从缓存中移除 内存或者redis
	err = m.Store.Remove(ctx.Req.Context(), sess.ID())
	if err != nil {
		return err
	}
	//http
	return m.Propagator.Remove(ctx.Resp)
}
