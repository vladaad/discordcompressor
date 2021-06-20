package settings

var Decoding = initDecoding()

func initDecoding() *decoding {
	return &decoding{
		HardwareAccel:       "auto",
	}
}

type decoding struct {
	HardwareAccel    string
}