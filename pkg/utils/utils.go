package pkg

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func ExecCommandInDir(inDir string, name string, args ...string) (string, error) {
	c := exec.Command(name, args...)
	if inDir != "" {
		c.Dir = inDir
	}

	return "", DisplayCommandOutput(c)
}

func DisplayCommandOutput(c *exec.Cmd) error {
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	return c.Start()
}

func ExecCommand(name string, args ...string) (string, error) {
	return ExecCommandInDir("", name, args...)
}

func ToClipboard(content []byte) {
	arch := runtime.GOOS
	var copyCmd *exec.Cmd

	if arch == "darwin" {
		copyCmd = exec.Command("pbcopy")
	} else if arch == "linux" {
		copyCmd = exec.Command("xclip", "-selection", "c")
	}

	in, err := copyCmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := copyCmd.Start(); err != nil {
		log.Fatal(err)
	}

	if _, err := in.Write([]byte(content)); err != nil {
		log.Fatal(err)
	}

	if err := in.Close(); err != nil {
		log.Fatal(err)
	}

	copyCmd.Wait()
}

func OpenBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}

func DownloadFile(url string, destinationFile string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(destinationFile)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	out.Chmod(0755)
	return err
}

func Expand(archive string, destinationDir string) error {
	if strings.HasSuffix(archive, ".zip") {
		return unzip(archive, destinationDir)
	} else if strings.HasSuffix(archive, ".tar.gz") {
		return untarGz(archive, destinationDir)
	} else {
		return fmt.Errorf("Unknown archive format")
	}
}

func unzip(archive string, destinationDir string) error {
	r, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		filePath := filepath.Join(destinationDir, f.Name)

		if !strings.HasPrefix(filePath, filepath.Clean(destinationDir)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: Illegal file path", filePath)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}

func untarGz(archive string, destinationDir string) error {
	reader, err := os.Open(archive)
	if err != nil {
		return err
	}
	defer reader.Close()

	gz, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	defer gz.Close()

	t := tar.NewReader(gz)
	for {
		header, err := t.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		if header.FileInfo().IsDir() {
			err = os.MkdirAll(filepath.Join(destinationDir, header.Name), 0755)
			if err != nil {
				return err
			}
		} else {
			if strings.ContainsRune(header.Name, os.PathSeparator) {
				idx := strings.LastIndex(header.Name, string(os.PathSeparator))
				parent := header.Name[:idx]
				err = os.MkdirAll(filepath.Join(destinationDir, parent), 0755)

				if err != nil {
					return err
				}
			}

			file, err := os.Create(filepath.Join(destinationDir, header.Name))
			if err != nil {
				return err
			}
			file.Chmod(0755)

			_, err = io.Copy(file, t)
			file.Close()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func IsLinux() bool {
	return runtime.GOOS == "linux"
}

func IsOSX() bool {
	return runtime.GOOS == "darwin"
}
