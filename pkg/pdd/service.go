package pdd

type Service int

const (
	serviceUnknown Service = iota
	serviceDNS
	serviceDomain
)

func (s Service) String() string {
	return [...]string{
		"unknown",
		"dns",
		"domain",
	}[s]
}

func (s Service) IsUnknown() bool {
	return s == serviceUnknown
}
