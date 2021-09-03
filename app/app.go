package app

import "context"

type fnAction func(ctx context.Context) error
type faction func() error

//App base app methods
type App interface {
	//SetactionFunc set action function to be run
	SetactionFunc(fn fnAction)
	//Run execute app.action method
	Run() error
}

//NewApp returns new instance App
func NewApp(ctx context.Context) App {
	return &app{ctx: ctx}
}

type app struct {
	action faction
	ctx    context.Context
}

//SetactionFunc implementing App.SetactionFunc
func (app *app) SetactionFunc(fn fnAction) {
	app.action = func() error {
		return fn(app.ctx)
	}
}

//Run implementing App.run
func (app *app) Run() error {
	return app.action()
}
