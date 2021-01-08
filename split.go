package split

import (
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
)

func Split(file string, max uint64) error {
	return split(file, max)
}

func Unsplit(file string) error {
	return unsplit(file)
}

func split(path string, max uint64) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("split: error opening file %s: %w", path, err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("split: error getting file info %s: %w", path, err)
	}
	size := info.Size()

	parts := uint64(math.Ceil(float64(size) / float64(max)))
	digits := countDigits(parts)

	var written int64
	for i := uint64(0); i < parts; i++ {
		name := fmt.Sprintf("%s.part%0*d", path, digits, i)
		part, err := os.Create(name)
		if err != nil {
			return fmt.Errorf("split: error creating part %s: %w", name, err)
		}
		defer part.Close()

		// consider better logic so that io.CopyN doesn't return an EOF error
		n, err := io.CopyN(part, file, int64(max))
		if err != nil && err != io.EOF {
			return fmt.Errorf("split: error copying data: %w, bytes written: %d", err, written)
		}
		written += int64(n)
	}
	if size != written {
		err := fmt.Errorf("file size (%d) and written bytes (%d) do not match", size, written)
		return fmt.Errorf("split: error copying data: %w", err)
	}

	return nil
}

func unsplit(path string) error {
	glob := path + ".part*"
	parts, err := filepath.Glob(glob)
	if err != nil {
		return fmt.Errorf("unsplit: error finding parts %s: %w", glob, err)
	}

	s, err := os.Stat(parts[0])
	perm := s.Mode().Perm()

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, perm)
	if err != nil {
		return fmt.Errorf("unsplit: error creating file %s: %w", path, err)
	}
	defer file.Close()

	for _, v := range parts {
		p, err := os.Open(v)
		if err != nil {
			return fmt.Errorf("unsplit: error opening part %s: %w", v, err)
		}
		defer p.Close()

		_, err = io.Copy(file, p)
		if err != nil {
			return fmt.Errorf("unsplit: error copying data: %w", err)
		}
	}

	return nil
}
