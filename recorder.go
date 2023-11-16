package audiorecorder

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/kalpio/audiorecorder/domain"
)

type recorder struct {
	ffmpeg string
	rand   *rand.Rand
}

func NewRecorder(ffmpeg string) *recorder {
	return &recorder{
		ffmpeg: ffmpeg,
		rand:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (r recorder) Record(ctx context.Context, device string, duration int) (*domain.Record, error) {
	cmd := exec.CommandContext(ctx, r.ffmpeg, ffmpegArguments(device, duration)...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	log.Println("Start recording...")

	result := &domain.Record{}
	_, err = io.Copy(result, stdout)
	if err != nil {
		return nil, err
	}

	if err := cmd.Wait(); err != nil {
		return nil, err
	}

	log.Println("Stop recording")

	return result, nil
}

func (r recorder) RecordFile(ctx context.Context, device string, duration int) (*domain.RecordFile, error) {
	fileName := path.Join(os.TempDir(), fmt.Sprintf("%s.wav", r.randomString(10)))
	cmd := exec.CommandContext(ctx, r.ffmpeg, ffmpegArgumentsFile(device, fileName, duration)...)
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	log.Println("Start recording...")

	if err := cmd.Wait(); err != nil {
		return nil, err
	}

	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			return // TODO: log error
		}
	}()

	record := &domain.Record{}
	_, err = io.Copy(record, f)
	if err != nil {
		return nil, err
	}

	return &domain.RecordFile{
		Path:   f.Name(),
		Record: record}, nil
}

var letterRune = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func (r recorder) randomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRune[r.rand.Intn(len(letterRune))]
	}

	return string(b)
}
