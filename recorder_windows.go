//go:build windows

package audiorecorder

import "strconv"

func ffmpegArguments(deviceName string, duration int) []string {
	return []string{
		"-hide_banner",
		"-f",
		"dshow",
		"-i",
		"audio=" + deviceName,
		"-sample_rate",
		"44100",
		"-t",
		strconv.Itoa(duration),
		"-f",
		"wav",
		"-"}
}

func ffmpegArgumentsFile(deviceName string, fileName string, duration int) []string {

	return []string{
		"-hide_banner",
		"-f",
		"dshow",
		"-i",
		"audio=" + deviceName,
		"-sample_rate",
		"44100",
		"-t",
		strconv.Itoa(duration),
		"-y",
		fileName}
}
