package botlogic

func heartbeatHandler() {
	select {
	case HeartbeatReceived <- struct{}{}:
	default:
	}
}
