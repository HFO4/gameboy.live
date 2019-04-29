package driver

type DisplayDriver interface {
	Init(*[160][144][3]uint8, string)
	Run(chan bool)
}
