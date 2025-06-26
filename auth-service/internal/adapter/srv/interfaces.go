package srv

type Server interface {
	Start() error
	Stop(ctx context.Context) error
}
