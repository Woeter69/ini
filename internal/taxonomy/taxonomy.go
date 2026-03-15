package taxonomy

type Category struct {
	ID          string // "game", "web", "ai"
	DisplayName string // "Games", "Web & Internet", "AI & Machine Learning"
}

// Global project types defined by the system.
// Any language can choose to support any slice of these.
var Categories = []Category{
	{"basic", "Basic / Standard Application"},
	{"app", "Full Application Structure"},
	{"web", "Web & Internet"},
	{"api", "API & Web Services"},
	{"mobile", "Mobile"},
	{"desktop", "Desktop"},
	{"game", "Games"},
	{"ai", "AI & Machine Learning"},
	{"data", "Data & Databases"},
	{"devops", "DevOps & Infrastructure"},
	{"network", "Networking"},
	{"security", "Security & Cryptography"},
	{"os", "Operating Systems & Low-level"},
	{"embedded", "Embedded & IoT"},
	{"lang", "Programming Languages & Compilers"},
	{"finance", "Finance & Trading"},
	{"comm", "Communication"},
	{"script", "Automation & Scripting"},
	{"monitor", "Observability & Monitoring"},
	{"stream", "Messaging & Streaming"},
	{"science", "Science & Research"},
	{"media", "Media & Content"},
	{"web3", "Blockchain & Web3"},
	{"graphics", "AR / VR / Graphics"},
	{"edu", "Education & Productivity"},
	{"business", "Business & Enterprise"},
	{"cli", "Command Line Interface (CLI)"},
	{"math", "Mathematics & Simulation"},
	{"stats", "Statistics & Analysis"},
}

// Aliases maps synonym IDs to canonical IDs.
var Aliases = map[string]string{
	"db":          "data",
	"storage":     "data",
	"iot":         "embedded",
	"interactive": "cli",
	"stat":        "stats",
	"mac":         "desktop",
	"ios":         "mobile",
}

// Canonical returns the canonical ID for a given ID (resolving aliases).
func Canonical(id string) string {
	if canonical, ok := Aliases[id]; ok {
		return canonical
	}
	return id
}

// GetName returns the display name for a given taxonomy ID.
func GetName(id string) string {
	id = Canonical(id)
	for _, c := range Categories {
		if c.ID == id {
			return c.DisplayName
		}
	}
	return id
}

// IsValid checks if the provided ID (or its alias) exists in the global taxonomy.
func IsValid(id string) bool {
	id = Canonical(id)
	for _, c := range Categories {
		if c.ID == id {
			return true
		}
	}
	return false
}
