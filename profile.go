package zsx

import (
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/exp/slices"
	"gopkg.in/ini.v1"
)

// https://docs.aws.amazon.com/sdkref/latest/guide/file-location.html
// https://docs.aws.amazon.com/sdkref/latest/guide/file-format.html
// https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html

func ListProfiles() ([]string, error) {
	configs, err := loadConfigs()
	if err != nil {
		return nil, err
	}
	credentials, err := loadCredentials()
	if err != nil {
		return nil, err
	}
	return mergedProfiles(configs, credentials), nil
}

type Config struct {
	name string
}

func loadConfigs() ([]*Config, error) {
	f := func(s *ini.Section) (*Config, bool) {
		if name, found := cutSectionProfile(s.Name()); found {
			c := &Config{
				name,
			}
			return c, true
		}
		return nil, false
	}
	path, err := resolveFilePath(".aws", "config")
	if err != nil {
		return nil, err
	}
	return load(path, f)
}

func cutSectionProfile(name string) (string, bool) {
	if s, found := strings.CutPrefix(name, "profile "); found {
		return s, true
	}
	if s, found := strings.CutPrefix(name, "sso-session "); found {
		return s, true
	}
	return "", false
}

type Credential struct {
	name string
}

func loadCredentials() ([]*Credential, error) {
	f := func(s *ini.Section) (*Credential, bool) {
		c := &Credential{
			s.Name(),
		}
		return c, true
	}
	path, err := resolveFilePath(".aws", "credentials")
	if err != nil {
		return nil, err
	}
	return load(path, f)
}

func resolveFilePath(elem ...string) (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	elem = append([]string{dir}, elem...)
	return filepath.Join(elem...), nil
}

func load[T any](path string, f func(*ini.Section) (T, bool)) ([]T, error) {
	cfg, err := ini.Load(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	ss := cfg.Sections()

	xs := make([]T, 0, len(ss))
	for _, s := range ss {
		name := s.Name()
		if name == ini.DefaultSection {
			continue
		}
		if x, ok := f(s); ok {
			xs = append(xs, x)
		}
	}
	return xs, nil
}

func mergedProfiles(configs []*Config, credentials []*Credential) []string {
	profiles := make([]string, 0)
	for _, c := range configs {
		if slices.Contains(profiles, c.name) {
			continue
		}
		profiles = append(profiles, c.name)
	}
	for _, c := range credentials {
		if slices.Contains(profiles, c.name) {
			continue
		}
		profiles = append(profiles, c.name)
	}
	return profiles
}
