package config

// GitHubStoreRepo configures where data goes
type GitHubStoreRepo struct {
	Owner string `env:"GITHUB_STORE_OWNER"`
	Repo  string `env:"GITHUB_STORE_REPO"`
}
