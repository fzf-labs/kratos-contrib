package bootstrap

import "os"

type Service struct {
	Name     string
	Version  string
	ID       string
	Metadata map[string]string
}

func NewService(name, version, id string) *Service {
	if id == "" {
		id, _ = os.Hostname()
	}
	return &Service{
		Name:     name,
		Version:  version,
		ID:       id,
		Metadata: map[string]string{},
	}
}

func (s *Service) GetInstanceID() string {
	return s.ID + "." + s.Name
}

func (s *Service) SetMataData(k, v string) {
	s.Metadata[k] = v
}
