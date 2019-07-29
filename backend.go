package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Backend struct {
	Name      string
	Path      string // directory of the backend
	Port      int
	Conn      net.Conn
	Sentences []string
	Config    *BackendConfiguration
}

type BackendRequest struct {
	Sentence string   `json:"sentence"`
	Args     []string `json:"args"`
}

type BackendResponse struct {
	Text  string            `json:"text"`
	Set   map[string]string `json:"set"`
	Error string            `json:"error"`
}

type BackendConfiguration struct {
	Run string
}

func NewBackend(name string) *Backend {
	return &Backend{
		Name: name,
	}
}

func (backend *Backend) connect() (err error) {
	to := fmt.Sprintf("127.0.0.1:%d", backend.Port)
	backend.Conn, err = net.Dial("tcp", to)

	return
}

func (backend *Backend) sendRaw(body string) (string, error) {
	_, err := fmt.Fprintln(backend.Conn, body)
	if err != nil {
		fmt.Printf("%s ! %v\n", backend.Name, err)
		return "", err
	}

	response, err := bufio.NewReader(backend.Conn).ReadString('\n')
	if err != nil {
		return "", err
	}

	return response, err
}

func (backend *Backend) send(request *BackendRequest) (*BackendResponse, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	responseData, err := backend.sendRaw(string(jsonData))
	if err != nil {
		return nil, err
	}

	var response *BackendResponse
	err = json.Unmarshal([]byte(responseData), &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (backend *Backend) loadSentences() error {
	response, err := backend.sendRaw(`{"special":"sentences"}`)
	if err != nil {
		return err
	}

	var parsedResponse struct {
		Sentences []string
	}
	err = json.Unmarshal([]byte(response), &parsedResponse)
	if err != nil {
		return err
	}

	backend.Sentences = parsedResponse.Sentences

	return nil
}

func (backend *Backend) backendDirectories() []string {
	directories := os.Getenv("BENTO_BACKEND")
	if directories == "" {
		// Default to the current directory if none are provided.
		directories = "."
	}

	return strings.Split(directories, ":")
}

func (backend *Backend) findBackend() error {
	// We always use the first backend that matches the name. Even if there are
	// other backends by the same name that would have otherwise been discovered
	// in the future.
	dirs := backend.backendDirectories()
	for _, dir := range dirs {
		path := fmt.Sprintf("%s/%s", dir, backend.Name)
		fs, err := os.Stat(path)
		if err == nil && fs.IsDir() {
			backend.Path = path
			return nil
		}
	}

	return fmt.Errorf("no such backend %s in any path: %s",
		backend.Name, strings.Join(dirs, ":"))
}

func (backend *Backend) readBentoConfig() error {
	data, err := ioutil.ReadFile(backend.Path + "/bento.json")
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &backend.Config)
	if err != nil {
		return err
	}

	return nil
}

func (backend *Backend) Start() error {
	err := backend.findBackend()
	if err != nil {
		return err
	}

	err = backend.readBentoConfig()
	if err != nil {
		return err
	}

	// TODO: Need a more reliable way to pick the port.
	backend.Port = 50000 + rand.Intn(1000)

	// TODO: This is problematic for arguments that have spaces.
	cmdParts := strings.Split(backend.Config.Run, " ")

	// It starts in the background.
	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
	cmd.Env = append(cmd.Env, fmt.Sprintf("BENTO_PORT=%d", backend.Port))
	cmd.Dir = backend.Path

	if err := cmd.Start(); err != nil {
		return err
	}

	// TODO: Remove this?
	time.Sleep(time.Second)

	err = backend.connect()
	if err != nil {
		return err
	}

	err = backend.loadSentences()
	if err != nil {
		return err
	}

	return nil
}

func (backend *Backend) String() string {
	return backend.Name
}
