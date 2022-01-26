package settings

var Decoding = initDecoding()

func initDecoding() *decoding {
	return &decoding{
		HardwareAccel:  "auto",
		TonemapHWAccel: true,
		ScalingHWAccel: "none",
	}
}

type decoding struct {
	HardwareAccel  string
	TonemapHWAccel bool
	ScalingHWAccel string
}
