package audit

// Service ..
type Service struct {
	RoleDao *SysAudit
}

// NewService ..
func NewService() *Service {
	return &Service{
		RoleDao: &SysAudit{},
	}
}
