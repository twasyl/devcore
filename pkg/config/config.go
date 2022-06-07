//go:generate go run -v gen_config.go

package config

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

// DevCoreConfig represents the configuration of the CLI
type DevCoreConfig struct {
	DockerCompose       DockerCompose     `json:"docker-compose"`
	DefaultToolsVersion map[string]string `json:"default-tools-version"`
	ProjectsDir         string            `json:"projects-dir"`
	ServersDir          string            `json:"servers-dir"`
	Jenkins             Jenkins           `json:"jenkins"`
}

type DockerCompose struct {
	CurrentContext string                 `json:"current-context"`
	Contexts       []DockerComposeContext `json:"contexts"`
}

// DockerComposeContext describes a docker compose environment
type DockerComposeContext struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	File        string `json:"file"`
}

// Jenkins describes the configuration of the Jenkins command
type Jenkins struct {
	Cli            string           `json:"cli"`
	CurrentContext string           `json:"current-context"`
	Contexts       []JenkinsContext `json:"contexts"`
}

type JenkinsContext struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	War         string   `json:"war"`
	JenkinsHome string   `json:"jenkins-home"`
	JavaHome    string   `json:"java-home"`
	Options     []string `json:"options"`
	JVMOptions  []string `json:"jvm-options"`
	Pid         int      `json:"pid"`
}

var Config = DevCoreConfig{}

// Load loads the CLI configuration into the Config struct.
func Load() error {
	ensureConfigFileSystemElements()
	file, err := os.Open(configFile())
	if err != nil {
		return err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&Config)
	if err != nil {
		return err
	}

	Config.fillDefaultToolsVersion()
	Config.fillDefaultProjectsDir()
	Config.fillDefaultServersDir()

	return nil
}

func (c *DevCoreConfig) fillDefaultToolsVersion() {
	if c.DefaultToolsVersion == nil {
		c.DefaultToolsVersion = make(map[string]string)
	}

	for _, tool := range DefaultSupportedTools {
		if _, exists := c.DefaultToolsVersion[tool.Name]; !exists {
			c.DefaultToolsVersion[tool.Name] = tool.DefaultVersion
		}
	}
}

func (c *DevCoreConfig) fillDefaultProjectsDir() {
	if c.ProjectsDir == "" {
		homeDir, _ := os.UserHomeDir()
		c.ProjectsDir = filepath.Join(homeDir, "Projects")
	}
}

func (c *DevCoreConfig) fillDefaultServersDir() {
	if c.ServersDir == "" {
		homeDir, _ := os.UserHomeDir()
		c.ServersDir = filepath.Join(homeDir, "Servers")
	}
}

// FindContextByName looks in the config for a DockerComposeContext named with the desired one.
func (c *DockerCompose) FindContextByName(name string) (DockerComposeContext, error) {
	for _, context := range c.Contexts {
		if context.Name == name {
			return context, nil
		}
	}
	return DockerComposeContext{}, &DockerComposeContextNotFound{name}
}

// AddContext will add the given context to the configuration.
func (c *DockerCompose) AddContext(context DockerComposeContext) {
	c.Contexts = append(c.Contexts, context)
}

// DeleteContext will remove the given context from the configuration.
func (c *DockerCompose) DeleteContext(toDelete DockerComposeContext) error {
	for index, context := range c.Contexts {
		if context.Name == toDelete.Name {
			contexts := c.Contexts
			c.Contexts = append(contexts[:index], contexts[index+1:]...)
			return nil
		}
	}
	return &DockerComposeContextNotFound{toDelete.Name}
}

// Dir returns the name of the folder containing the Docker compose file.
func (c *DockerComposeContext) Dir() string {
	return filepath.Dir(c.File)
}

// FindContextByName looks in the config for a DockerComposeContext named with the desired one.
func (j *Jenkins) FindContextByName(name string) (JenkinsContext, error) {
	for _, context := range j.Contexts {
		if context.Name == name {
			return context, nil
		}
	}
	return JenkinsContext{}, &JenkinsContextNotFound{name}
}

// AddContext will add the given context to the configuration.
func (j *Jenkins) AddContext(context JenkinsContext) {
	j.Contexts = append(j.Contexts, context)
}

// DeleteContext will remove the given context from the configuration.
func (j *Jenkins) DeleteContext(toDelete JenkinsContext) error {
	for index, context := range j.Contexts {
		if context.Name == toDelete.Name {
			contexts := j.Contexts
			j.Contexts = append(contexts[:index], contexts[index+1:]...)
			return nil
		}
	}
	return &JenkinsContextNotFound{toDelete.Name}
}

func (j *Jenkins) UpdateContext(c JenkinsContext) error {
	for index := range j.Contexts {
		if j.Contexts[index].Name == c.Name {
			j.Contexts[index].Description = c.Description
			j.Contexts[index].JavaHome = c.JavaHome
			j.Contexts[index].JenkinsHome = c.JenkinsHome
			j.Contexts[index].War = c.War
			j.Contexts[index].Pid = c.Pid
			return nil
		}
	}
	return &JenkinsContextNotFound{c.Name}
}

// Save will save the CLI configuration to the file system.
func Save() error {
	ensureConfigFileSystemElements()
	_, err := os.Stat(configFile())
	if err != nil {
		return err
	}

	buffer := bytes.Buffer{}
	err = json.NewEncoder(&buffer).Encode(Config)
	if err != nil {
		return err
	}
	return os.WriteFile(configFile(), buffer.Bytes(), 0755)
}

func ensureConfigFileSystemElements() {
	if _, err := os.Stat(configDir()); os.IsNotExist(err) {
		err := os.Mkdir(configDir(), 0755)
		if err != nil {
			log.Fatal("Can not create configuration directory", err)
		}
	}

	if _, err := os.Stat(configFile()); os.IsNotExist(err) {
		file, err := os.Create(configFile())
		if err != nil {
			log.Fatal("Can not create configuration file", err)
		} else {
			defer file.Close()
			file.Write([]byte("{}"))
		}
	}
}

func configDir() string {
	if homeDir, err := os.UserHomeDir(); err != nil {
		log.Fatal("Can not determine user home dir", err)
		return ""
	} else {
		return filepath.Join(homeDir, ".devcore")
	}
}

func configFile() string {
	return filepath.Join(configDir(), "config.json")
}
