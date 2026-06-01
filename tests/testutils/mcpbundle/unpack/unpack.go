// Copyright 2026 The MathWorks, Inc.

package unpack

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type WriteCloser interface {
	Write(p []byte) (int, error)
	Close() error
}

type FileSystem interface {
	MkdirAll(path string, perm fs.FileMode) error
	OpenFile(name string, flag int, perm fs.FileMode) (WriteCloser, error)
}

type Unpacker struct {
	fs FileSystem
}

func New() *Unpacker {
	return newUnpacker(osFileSystem{})
}

func newUnpacker(fs FileSystem) *Unpacker {
	return &Unpacker{fs: fs}
}

func (u *Unpacker) Unpack(archive io.ReaderAt, size int64, destDir string) error {
	r, err := zip.NewReader(archive, size)
	if err != nil {
		return fmt.Errorf("reading archive: %w", err)
	}

	for _, f := range r.File {
		targetPath := filepath.Join(destDir, f.Name) //nolint:gosec // Validated on next line
		cleanTarget := filepath.Clean(targetPath)
		cleanDest := filepath.Clean(destDir)
		if cleanTarget != cleanDest && !strings.HasPrefix(cleanTarget, cleanDest+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path in archive: %s", f.Name)
		}
		if f.FileInfo().IsDir() {
			if err := u.fs.MkdirAll(targetPath, 0750); err != nil {
				return fmt.Errorf("creating directory %s: %w", f.Name, err)
			}
			continue
		}

		if err := u.fs.MkdirAll(filepath.Dir(targetPath), 0750); err != nil {
			return fmt.Errorf("creating parent for %s: %w", f.Name, err)
		}

		if err := u.extractFile(f, targetPath); err != nil {
			return err
		}
	}
	return nil
}

func (u *Unpacker) extractFile(f *zip.File, targetPath string) error {
	rc, err := f.Open()
	if err != nil {
		return fmt.Errorf("opening %s in archive: %w", f.Name, err)
	}
	defer func() { _ = rc.Close() }()

	out, err := u.fs.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Mode())
	if err != nil {
		return fmt.Errorf("creating %s: %w", f.Name, err)
	}
	defer func() { _ = out.Close() }()

	if _, err = io.Copy(out, rc); err != nil { //nolint:gosec // Size is bounded by trusted CI-built archive
		return fmt.Errorf("copying %s: %w", f.Name, err)
	}
	return nil
}

type osFileSystem struct{}

func (osFileSystem) MkdirAll(path string, perm fs.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (osFileSystem) OpenFile(name string, flag int, perm fs.FileMode) (WriteCloser, error) {
	return os.OpenFile(name, flag, perm) //nolint:gosec // Preserves archive permissions
}
