package hardware

import (
	"os/exec"
)

func CudaCheck(ffmpegExec string) error {
	var options []string
	options = append(options, "-y", "-hwaccel", "cuda", "-f", "lavfi", "-i", "nullsrc=2560x1440")
	options = append(options, "-vframes", "1", "-an")
	options = append(options, "-vf", "hwupload_cuda,scale_cuda=1920:1080:lanczos:p010,hwdownload,format=p010le,format=yuv420p10le")
	options = append(options, "-f", "null", "-")

	cmd := exec.Command(ffmpegExec, options...)

	err := cmd.Start()
	err = cmd.Wait()
	return err
}
