package taxonomy

type Category struct {
	ID          string // "game", "web", "ai"
	DisplayName string // "Games", "Web & Internet", "AI & Machine Learning"
}

// Global project types defined by the system.
// Any language can choose to support any slice of these.
var Categories = []Category{
	{"basic", "Basic / Standard Application"},
	{"web", "Web & Internet"},
	{"mobile", "Mobile"},
	{"desktop", "Desktop"},
	{"game", "Games"},
	{"ai", "AI & Machine Learning"},
	{"data", "Data"},
	{"db", "Databases & Storage"},
	{"devops", "DevOps & Infrastructure"},
	{"network", "Networking"},
	{"security", "Security & Cryptography"},
	{"os", "Operating Systems & Low-level"},
	{"lang", "Programming Languages & Compilers"},
	{"finance", "Finance & Trading"},
	{"comm", "Communication"},
	{"script", "Automation & Scripting"},
	{"monitor", "Observability & Monitoring"},
	{"stream", "Messaging & Streaming"},
	{"science", "Science & Research"},
	{"media", "Media & Content"},
	{"iot", "IoT & Hardware"},
	{"web3", "Blockchain & Web3"},
	{"graphics", "AR / VR / Graphics"},
	{"edu", "Education & Productivity"},
}

// GetName returns the display name for a given taxonomy ID.
func GetName(id string) string {
	for _, c := range Categories {
		if c.ID == id {
			return c.DisplayName
		}
	}
	return id
}

// IsValid checks if the provided ID exists in the global taxonomy.
func IsValid(id string) bool {
	for _, c := range Categories {
		if c.ID == id {
			return true
		}
	}
	return false
}
