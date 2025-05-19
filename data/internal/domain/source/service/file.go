package service

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/larek-tech/diploma/data/internal/domain/file"
	"github.com/larek-tech/diploma/data/internal/domain/source"
)

func getFileExtension(filename string) string {
	ext := filepath.Ext(filename)
	if ext == "" {
		return "unknown"
	}
	return ext[1:] // Remove the dot
}

func (s Service) createFile(src *source.Source, msg source.DataMessage) (*file.File, error) {
	f := file.NewFile(src.ID, getFileExtension(msg.Title))
	f.Filename = msg.Title
	f.Raw = msg.Content
	return f, nil
}

func (s Service) createArchive(src *source.Source, msg source.DataMessage) ([]*file.File, error) {
	// create temp file for archive
	tmpZipFile, err := os.CreateTemp("", "archive-*.zip")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp zip file: %w", err)
	}
	defer tmpZipFile.Close()

	defer os.Remove(tmpZipFile.Name())

	// write archive to temp file
	if _, err := tmpZipFile.Write(msg.Content); err != nil {
		return nil, fmt.Errorf("failed to write archive content: %w", err)
	}

	tmpDir, err := os.MkdirTemp("", "unzipped-archive-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	archive, err := zip.OpenReader(tmpZipFile.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to open archive: %w", err)
	}
	defer archive.Close()

	var files []*file.File

	for _, f := range archive.File {
		filePath := filepath.Join(tmpDir, f.Name)
		if !strings.HasPrefix(filePath, filepath.Clean(tmpDir)+string(os.PathSeparator)) {
			return nil, fmt.Errorf("invalid file path: %s", filePath)
		}
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
				return nil, fmt.Errorf("failed to create dir: %w", err)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return nil, fmt.Errorf("failed to create parent dir: %w", err)
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return nil, fmt.Errorf("failed to create file: %w", err)
		}

		fileInArchive, err := f.Open()
		if err != nil {
			dstFile.Close()
			return nil, fmt.Errorf("failed to open file in archive: %w", err)
		}

		content, err := io.ReadAll(fileInArchive)
		if err != nil {
			fileInArchive.Close()
			dstFile.Close()
			return nil, fmt.Errorf("failed to read file content: %w", err)
		}

		if _, err := dstFile.Write(content); err != nil {
			fileInArchive.Close()
			dstFile.Close()
			return nil, fmt.Errorf("failed to write file content: %w", err)
		}

		dstFile.Close()
		fileInArchive.Close()

		newFile := file.NewFile(src.ID, getFileExtension(f.Name))
		newFile.Filename = f.Name
		newFile.Raw = content
		files = append(files, newFile)
	}

	return files, nil
}
