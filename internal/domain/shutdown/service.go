package shutdown

type Service struct {
	shutdownRequest *chan struct{}
}

func NewService(ch *chan struct{}) *Service {
	return &Service{
		shutdownRequest: ch,
	}
}

func (s Service) Shutdown() {
	select {
	case *s.shutdownRequest <- struct{}{}:
	default:
	}
}
