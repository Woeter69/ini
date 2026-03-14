package handler

import "fmt"

// ProjectConfig holds the configuration for a new project.
type ProjectConfig struct {
	Name     string
	Path     string // Resolved absolute path where project will be created
	Language string // Language key (e.g. "python", "node")
	Type     string // Type of project from global taxonomy (e.g. "web", "game"). Default "basic"
	Git      bool   // Whether to initialize a git repo
}

// Handler is the interface that all language handlers must implement.
type Handler interface {
	// Name returns the display name of the language.
	Name() string

	// Validate checks if the required toolchain is installed.
	Validate() error

	// Init creates the project with the given config.
	Init(config ProjectConfig) error
}

// TypedHandler is an optional interface for handlers that support specific project types.
// If a handler doesn't implement this, it implicitly only supports the "basic" type.
type TypedHandler interface {
	SupportedTypes() []string
}

// registry stores all registered handlers
var registry = map[string]Handler{}

// Register adds a handler to the registry.
func Register(name string, h Handler) {
	registry[name] = h
}

// Get retrieves a handler by name.
func Get(name string) (Handler, error) {
	h, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("no handler registered for %q", name)
	}
	return h, nil
}
