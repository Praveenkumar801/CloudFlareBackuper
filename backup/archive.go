package backup

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

func CreateArchive(folders []string, outputPath string) error {

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create archive file: %w", err)
	}
	defer outFile.Close()

	gzipWriter := gzip.NewWriter(outFile)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	for _, folder := range folders {
		if err := addToArchive(tarWriter, folder); err != nil {
			return fmt.Errorf("failed to add %s to archive: %w", folder, err)
		}
	}

	return nil
}

func addToArchive(tarWriter *tar.Writer, sourcePath string) error {

	_, err := os.Stat(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to stat %s: %w", sourcePath, err)
	}

	return filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return fmt.Errorf("failed to create tar header: %w", err)
		}

		header.Name = path

		if err := tarWriter.WriteHeader(header); err != nil {
			return fmt.Errorf("failed to write tar header: %w", err)
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("failed to open file %s: %w", path, err)
			}

			_, copyErr := io.Copy(tarWriter, file)
			file.Close() // Close immediately after copying, not deferred

			if copyErr != nil {
				return fmt.Errorf("failed to write file content: %w", copyErr)
			}
		}

		return nil
	})
}

func GenerateBackupFilename(prefix string) string {
	timestamp := time.Now().Format("20060102-150405")
	return fmt.Sprintf("%s-%s.tar.gz", prefix, timestamp)
}
