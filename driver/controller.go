package driver

type ControllerDriver interface {
	InitStatus(*byte)
	UpdateInput() bool
}
