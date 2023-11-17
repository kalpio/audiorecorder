package audiorecorder

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"strings"
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
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), "Press [q] to stop") {
				log.Println("Start recording...")
			}
			if strings.Contains(scanner.Text(), "size=") {
				log.Println("Stop recording")
			}
		}
	}()

	result := &domain.Record{}
	_, err = io.Copy(result, stderr)
	if err != nil {
		return nil, err
	}

	if err := cmd.Wait(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r recorder) RecordFile(ctx context.Context, device string, duration int) (*domain.RecordFile, error) {
	fileName := path.Join(os.TempDir(), fmt.Sprintf("%s.wav", r.randomString(10)))
	cmd := exec.CommandContext(ctx, r.ffmpeg, ffmpegArgumentsFile(device, fileName, duration)...)
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), "Press [q] to stop") {
				log.Println("Start recording...")
			}
			if strings.Contains(scanner.Text(), "size=") {
				log.Println("Stop recording")
			}
		}
	}()

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
