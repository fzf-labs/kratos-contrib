package bootstrap

import "os"

type Service struct {
	Name     string
	Version  string
	ID       string
	Metadata map[string]string
}

func NewService(name, version, id string, metadata map[string]string) *Service {
	if id == "" {
		id, _ = os.Hostname()
	}
	if metadata == nil {
		metadata = map[string]string{}
	}
	return &Service{
		Name:     name,
		Version:  version,
		ID:       id,
		Metadata: metadata,
	}
}

// GetName 获取服务名
func (s *Service) GetName() string {
	return s.Name
}

func (s *Service) GetVersion() string {
	return s.Version
}

// GetInstanceID 获取实例ID
func (s *Service) GetInstanceID() string {
	return s.ID
}

// GetMetadata 获取元数据
func (s *Service) GetMetadata() map[string]string {
	return s.Metadata
}


