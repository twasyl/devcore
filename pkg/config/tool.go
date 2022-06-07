package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	du "io.twasyl/devcore/pkg/utils"
)

type Tool struct {
	Name            string
	Description     string
	CommandLineName string
	DefaultVersion  string
	OS              func() string
	Install         func(version string) error
}

// DefaultSupportedTools represents the tools that can be installed using `devcore tools install` with their default
// version.
var DefaultSupportedTools []Tool

func bat() Tool {
	tool := Tool{
		Name:            "bat",
		Description:     "A cat like tool, but more powerful",
		CommandLineName: "bat",
		DefaultVersion:  "0.20.0",
		OS: func() string {
			if runtime.GOOS == "darwin" {
				return "apple-darwin"
			} else if runtime.GOOS == "windows" {
				return "pc-windows-gnu"
			} else {
				return "unknown-linux-gnu"
			}
		},
	}
	tool.Install = func(version string) error {
		var archive = filepath.Join(os.TempDir(), "bat.tar.gz")
		osName := tool.OS()

		err := du.DownloadFile(fmt.Sprintf("https://github.com/sharkdp/bat/releases/download/v%s/bat-v%s-x86_64-%s.tar.gz", version, version, osName), archive)
		if err != nil {
			return err
		}

		destinationDir := filepath.Join(os.TempDir(), "bat")
		err = os.Mkdir(destinationDir, 0755)
		if err != nil {
			return err
		}

		err = du.Expand(archive, destinationDir)
		if err != nil {
			return err
		}

		err = os.Rename(filepath.Join(destinationDir, fmt.Sprintf("bat-v%s-x86_64-%s", version, osName), "bat"), fmt.Sprintf("/usr/local/bin/%s", tool.CommandLineName))
		if err != nil {
			return err
		}

		err = os.Remove(archive)
		if err != nil {
			return err
		}

		return os.RemoveAll(destinationDir)
	}
	return tool
}

func geckodriver() Tool {
	tool := Tool{
		Name:            "geckodriver",
		Description:     "A driver to be used by Selenium",
		CommandLineName: "geckodriver",
		DefaultVersion:  "0.31.0",
		OS: func() string {
			if runtime.GOOS == "darwin" {
				return "macos"
			} else if runtime.GOOS == "linux" {
				return "linux64"
			} else {
				return "win64"
			}
		},
	}
	tool.Install = func(version string) error {
		var archive = filepath.Join(os.TempDir(), "geckodriver.tar.gz")

		err := du.DownloadFile(fmt.Sprintf("https://github.com/mozilla/geckodriver/releases/download/v%s/geckodriver-v%s-%s.tar.gz", version, version, tool.OS()), archive)
		if err != nil {
			return err
		}

		destinationDir := filepath.Join(os.TempDir(), "geckodriver")
		err = os.Mkdir(destinationDir, 0755)
		if err != nil {
			return err
		}

		err = du.Expand(archive, destinationDir)
		if err != nil {
			return err
		}

		err = os.Rename(filepath.Join(destinationDir, "geckodriver"), fmt.Sprintf("/usr/local/bin/%s", tool.CommandLineName))
		if err != nil {
			return err
		}

		err = os.Remove(archive)
		if err != nil {
			return err
		}

		return os.RemoveAll(destinationDir)
	}

	return tool
}

func gh() Tool {
	tool := Tool{
		Name:            "gh",
		Description:     "GitHub CLI tool for interacting with GitHub",
		CommandLineName: "gh",
		DefaultVersion:  "2.9.0",
		OS: func() string {
			if runtime.GOOS == "darwin" {
				return "macOS"
			} else {
				return runtime.GOOS
			}
		},
	}
	tool.Install = func(version string) error {
		var archive = filepath.Join(os.TempDir(), "gh.tar.gz")
		osName := tool.OS()

		err := du.DownloadFile(fmt.Sprintf("https://github.com/cli/cli/releases/download/v%s/gh_%s_%s_amd64.tar.gz", version, version, osName), archive)
		if err != nil {
			return err
		}

		destinationDir := filepath.Join(os.TempDir(), "github")
		err = os.Mkdir(destinationDir, 0755)
		if err != nil {
			return err
		}

		err = du.Expand(archive, destinationDir)
		if err != nil {
			return err
		}

		err = os.Rename(filepath.Join(destinationDir, fmt.Sprintf("gh_%s_%s_amd64", version, osName), "bin", "gh"), fmt.Sprintf("/usr/local/bin/%s", tool.CommandLineName))
		if err != nil {
			return err
		}

		err = os.Remove(archive)
		if err != nil {
			return err
		}

		return os.RemoveAll(destinationDir)
	}
	return tool
}

func openshiftClient() Tool {
	tool := Tool{
		Name:            "openshift-client",
		Description:     "Interact with Openshift in the CLI",
		CommandLineName: "oc",
		DefaultVersion:  "4.10.10",
		OS: func() string {
			if runtime.GOOS == "darwin" {
				return "mac"
			} else {
				return runtime.GOOS
			}
		},
	}
	tool.Install = func(version string) error {
		var archive = filepath.Join(os.TempDir(), "openshift-client.tar.gz")
		err := du.DownloadFile(fmt.Sprintf("https://mirror.openshift.com/pub/openshift-v4/x86_64/clients/ocp/%s/openshift-client-%s-%s.tar.gz", version, tool.OS(), version), archive)
		if err != nil {
			return err
		}

		destinationDir := filepath.Join(os.TempDir(), "openshift-client")
		err = os.Mkdir(destinationDir, 0755)
		if err != nil {
			return err
		}

		err = du.Expand(archive, destinationDir)
		if err != nil {
			return err
		}

		err = os.Rename(filepath.Join(destinationDir, "oc"), fmt.Sprintf("/usr/local/bin/%s", tool.CommandLineName))
		if err != nil {
			return err
		}

		err = os.Remove(archive)
		if err != nil {
			return err
		}

		return os.RemoveAll(destinationDir)
	}
	return tool
}

func openshiftInstall() Tool {
	tool := Tool{
		Name:            "openshift-install",
		Description:     "Openshift installer",
		CommandLineName: "openshift-install",
		DefaultVersion:  "4.9.12",
		OS: func() string {
			if runtime.GOOS == "darwin" {
				return "mac"
			} else {
				return runtime.GOOS
			}
		},
	}
	tool.Install = func(version string) error {
		var archive = filepath.Join(os.TempDir(), "openshift-install.tar.gz")

		err := du.DownloadFile(fmt.Sprintf("https://mirror.openshift.com/pub/openshift-v4/x86_64/clients/ocp/%s/openshift-install-%s-%s.tar.gz", version, tool.OS(), version), archive)
		if err != nil {
			return err
		}

		destinationDir := filepath.Join(os.TempDir(), "openshift-install")
		err = os.Mkdir(destinationDir, 0755)
		if err != nil {
			return err
		}

		err = du.Expand(archive, destinationDir)
		if err != nil {
			return err
		}

		err = os.Rename(filepath.Join(destinationDir, "openshift-install"), fmt.Sprintf("/usr/local/bin/%s", tool.CommandLineName))
		if err != nil {
			return err
		}

		err = os.Remove(archive)
		if err != nil {
			return err
		}

		return os.RemoveAll(destinationDir)
	}
	return tool
}

func init() {

	DefaultSupportedTools = []Tool{
		bat(),
		{
			Name:            "dive",
			Description:     "A tool for exploring container in depth",
			CommandLineName: "dive",
			DefaultVersion:  "0.10.0",
			Install: func(version string) error {
				var archive = filepath.Join(os.TempDir(), "dive.tar.gz")
				err := du.DownloadFile(fmt.Sprintf("https://github.com/wagoodman/dive/releases/download/v%s/dive_%s_%s_amd64.tar.gz", version, version, runtime.GOOS), archive)
				if err != nil {
					return err
				}

				destinationDir := filepath.Join(os.TempDir(), "dive")
				err = os.Mkdir(destinationDir, 0755)
				if err != nil {
					return err
				}

				err = du.Expand(archive, destinationDir)
				if err != nil {
					return err
				}

				err = os.Rename(filepath.Join(destinationDir, "dive"), "/usr/local/bin/dive")
				if err != nil {
					return err
				}
				err = os.Chmod("/usr/local/bin/dive", 0755)
				if err != nil {
					return err
				}

				err = os.Remove(archive)
				if err != nil {
					return err
				}

				return os.RemoveAll(destinationDir)
			},
		},
		geckodriver(),
		gh(),
		{
			Name:            "helm",
			Description:     "The package manager for Kubernetes",
			CommandLineName: "helm",
			DefaultVersion:  "3.8.2",
			Install: func(version string) error {
				var archive = filepath.Join(os.TempDir(), "helm.tar.gz")
				err := du.DownloadFile(fmt.Sprintf("https://get.helm.sh/helm-v%s-%s-amd64.tar.gz", version, runtime.GOOS), archive)
				if err != nil {
					return err
				}

				destinationDir := filepath.Join(os.TempDir(), "helm")
				err = os.Mkdir(destinationDir, 0755)
				if err != nil {
					return err
				}

				err = du.Expand(archive, destinationDir)
				if err != nil {
					return err
				}

				err = os.Rename(filepath.Join(destinationDir, fmt.Sprintf("%s-amd64", runtime.GOOS), "helm"), "/usr/local/bin/helm")
				if err != nil {
					return err
				}

				err = os.Remove(archive)
				if err != nil {
					return err
				}

				return os.RemoveAll(destinationDir)
			},
		},
		{
			Name:            "jq",
			Description:     "jq is a lightweight and flexible command-line JSON processor",
			CommandLineName: "jq",
			DefaultVersion:  "1.6",
			OS: func() string {
				if runtime.GOOS == "darwin" {
					return "macOS"
				} else {
					return runtime.GOOS
				}
			},
			Install: func(version string) error {
				var jqDownloadFilename string
				var destinationFile string
				if runtime.GOOS == "darwin" {
					jqDownloadFilename = "jq-osx-amd64"
					destinationFile = "/usr/local/bin/jq"
				} else if runtime.GOOS == "linux" {
					jqDownloadFilename = "jq-linux64"
					destinationFile = "/usr/local/bin/jq"
				} else {
					jqDownloadFilename = "jq-win64.exe"
					destinationFile = "C:\\jq.exe"
				}

				return du.DownloadFile(fmt.Sprintf("https://github.com/stedolan/jq/releases/download/jq-%s/%s", version, jqDownloadFilename), destinationFile)
			},
		},
		{
			Name:            "kind",
			Description:     "kind is a tool for running local Kubernetes clusters using Docker container \"nodes\"",
			CommandLineName: "kind",
			DefaultVersion:  "0.12.0",
			Install: func(version string) error {
				return du.DownloadFile(fmt.Sprintf("https://github.com/kubernetes-sigs/kind/releases/download/v%s/kind-%s-amd64", version, runtime.GOOS), "/usr/local/bin/kind")
			},
		},
		{
			Name:            "kubectl",
			Description:     "The Kubernetes command-line tool",
			CommandLineName: "kubectl",
			DefaultVersion:  "1.23.6",
			Install: func(version string) error {
				return du.DownloadFile(fmt.Sprintf("https://dl.k8s.io/release/v%s/bin/%s/amd64/kubectl", version, runtime.GOOS), "/usr/local/bin/kubectl")
			},
		},
		{
			Name:            "kustomize",
			Description:     "Kubernetes native configuration management",
			CommandLineName: "kustomize",
			DefaultVersion:  "4.5.4",
			Install: func(version string) error {
				var archive = filepath.Join(os.TempDir(), "kustomize.tar.gz")
				err := du.DownloadFile(fmt.Sprintf("https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize/v%s/kustomize_v%s_%s_amd64.tar.gz", version, version, runtime.GOOS), archive)
				if err != nil {
					return err
				}

				destinationDir := "/usr/local/bin"

				err = du.Expand(archive, destinationDir)
				if err != nil {
					return err
				}

				err = os.Remove(archive)
				if err != nil {
					return err
				}

				return os.Chmod(filepath.Join(destinationDir, "kustomize"), 0755)
			},
		},
		{
			Name:            "maven",
			Description:     "Apache Maven is a software project management and comprehension tool",
			CommandLineName: "mvn",
			DefaultVersion:  "3.8.5",
			Install: func(version string) error {
				var archive = filepath.Join(os.TempDir(), "maven.zip")
				err := du.DownloadFile(fmt.Sprintf("https://dlcdn.apache.org/maven/maven-3/%s/binaries/apache-maven-%s-bin.zip", version, version), archive)
				if err != nil {
					return err
				}

				var destinationDir string
				if runtime.GOOS == "darwin" {
					destinationDir = "/Library/maven"
				} else if runtime.GOOS == "linux" {
					destinationDir = "/usr/local/maven"
				} else {
					destinationDir = "C:\\maven"
				}

				if _, err := os.Stat(destinationDir); os.IsNotExist(err) {
					err = os.Mkdir(destinationDir, 0755)
					if err != nil {
						return err
					}
				}

				err = du.Expand(archive, destinationDir)
				if err != nil {
					return err
				}

				err = os.Rename(filepath.Join(destinationDir, fmt.Sprintf("apache-maven-%s", version)), filepath.Join(destinationDir, version))
				if err != nil {
					return err
				}

				if runtime.GOOS == "darwin" {
					mavenPathD := "/etc/paths.d/maven"
					if _, err := os.Stat(mavenPathD); os.IsNotExist(err) {
						if _, err := os.Create(mavenPathD); err != nil {
							return err
						}
					}

					if err := os.WriteFile(mavenPathD, []byte(fmt.Sprintf("/Library/maven/%s/bin/mvn", version)), 0755); err != nil {
						return err
					}
				} else if runtime.GOOS == "linux" {
					if err := os.Symlink(fmt.Sprintf("/usr/local/maven/%s/bin/mvn", version), "/usr/local/bin/mvn"); err != nil {
						return err
					}
				}

				return os.Remove(archive)
			},
		},
		{
			Name:            "minishift",
			Description:     "Minishift is a tool that helps you run OpenShift locally by running a single-node OpenShift cluster inside a VM",
			CommandLineName: "minishift",
			DefaultVersion:  "1.34.3",
			Install: func(version string) error {
				var archive = filepath.Join(os.TempDir(), "minishift.tar.gz")
				err := du.DownloadFile(fmt.Sprintf("https://github.com/minishift/minishift/releases/download/v%s/minishift-%s-%s-amd64.tgz", version, version, runtime.GOOS), archive)
				if err != nil {
					return err
				}

				destinationDir := filepath.Join(os.TempDir(), "minishift")
				err = os.Mkdir(destinationDir, 0755)
				if err != nil {
					return err
				}

				err = du.Expand(archive, destinationDir)
				if err != nil {
					return err
				}

				err = os.Rename(filepath.Join(destinationDir, fmt.Sprintf("minishift-%s-%s-amd64", version, runtime.GOOS), "minishift"), "/usr/local/bin/minishift")
				if err != nil {
					return err
				}
				err = os.Chmod("/usr/local/bin/minishift", 0755)
				if err != nil {
					return err
				}

				err = os.Remove(archive)
				if err != nil {
					return err
				}

				return os.RemoveAll(destinationDir)
			},
		},
		openshiftClient(),
		openshiftInstall(),
		{
			Name:            "operator-sdk",
			Description:     "The Operator SDK provides the tools to build, test, and package Operators",
			CommandLineName: "operator-sdk",
			DefaultVersion:  "1.9.0",
			Install: func(version string) error {
				return du.DownloadFile(fmt.Sprintf("https://github.com/operator-framework/operator-sdk/releases/download/v%s/operator-sdk_%s_amd64", version, runtime.GOOS), "/usr/local/bin/operator-sdk")
			},
		},
		{
			Name:            "terraform",
			Description:     "Terraform is an open-source infrastructure as code software tool that enables you to safely and predictably create, change, and improve infrastructure",
			CommandLineName: "terraform",
			DefaultVersion:  "1.1.9",
			Install: func(version string) error {
				var archive = filepath.Join(os.TempDir(), "terraform.zip")
				err := du.DownloadFile(fmt.Sprintf("https://releases.hashicorp.com/terraform/%s/terraform_%s_%s_amd64.zip", version, version, runtime.GOOS), archive)
				if err != nil {
					return err
				}

				destinationDir := filepath.Join(os.TempDir(), "terraform")
				err = os.Mkdir(destinationDir, 0755)
				if err != nil {
					return err
				}

				err = du.Expand(archive, destinationDir)
				if err != nil {
					return err
				}

				err = os.Rename(filepath.Join(destinationDir, "terraform"), "/usr/local/bin/terraform")
				if err != nil {
					return err
				}

				err = os.Remove(archive)
				if err != nil {
					return err
				}

				return os.RemoveAll(destinationDir)
			},
		},
		{
			Name:            "yq",
			Description:     "yq is a lightweight and portable command-line YAML, JSON and XML processor",
			CommandLineName: "yq",
			DefaultVersion:  "4.24.5",
			Install: func(version string) error {
				return du.DownloadFile(fmt.Sprintf("https://github.com/mikefarah/yq/releases/download/v%s/yq_%s_amd64", version, runtime.GOOS), "/usr/local/bin/yq")
			},
		},
	}
}

func FindTool(name string) (Tool, error) {
	for _, tool := range DefaultSupportedTools {
		if tool.Name == name {
			return tool, nil
		}
	}
	return Tool{}, &ToolNotFound{Name: name}
}
