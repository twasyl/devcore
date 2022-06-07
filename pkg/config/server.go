package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	du "io.twasyl/devcore/pkg/utils"
)

// Server represents an application server that can be installed and used.
type Server struct {
	Name           string
	DefaultVersion string
	OS             func() string
	Install        func(version string) error
}

// DefaultSupportedServers represents the servers that can be installed using `devcore servers install` with their default
// version.
var DefaultSupportedServers []Server

func init() {
	DefaultSupportedServers = []Server{
		{
			Name:           "tomcat",
			DefaultVersion: "10.0.10",
			Install: func(version string) error {
				tomcatsDir := filepath.Join(Config.ServersDir, "tomcat")
				versionDir := filepath.Join(tomcatsDir, version)

				if _, err := os.Stat(versionDir); os.IsNotExist(err) == false {
					return errors.New(fmt.Sprintf("Tomcat %s already present at %s", version, versionDir))
				}

				if _, err := os.Stat(tomcatsDir); os.IsNotExist(err) {
					err = os.MkdirAll(tomcatsDir, 0755)
					if err != nil {
						return nil
					}
				}

				var archive = filepath.Join(os.TempDir(), "tomcat.zip")
				majorVersion := version[0:strings.Index(version, ".")]
				err := du.DownloadFile(fmt.Sprintf("https://archive.apache.org/dist/tomcat/tomcat-%s/v%s/bin/apache-tomcat-%s.zip", majorVersion, version, version), archive)
				if err != nil {
					return err
				}

				err = du.Expand(archive, tomcatsDir)
				if err != nil {
					return err
				}

				err = os.Rename(filepath.Join(tomcatsDir, fmt.Sprintf("apache-tomcat-%s", version)), versionDir)
				if err != nil {
					return err
				}

				err = os.Remove(archive)
				return err
			},
		},
	}
}

func FindServer(name string) (Server, error) {
	for _, server := range DefaultSupportedServers {
		if server.Name == name {
			return server, nil
		}
	}
	return Server{}, &ServerNotFound{Name: name}
}
