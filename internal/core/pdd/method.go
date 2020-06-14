package pdd

import "net/http"

// Метод API.
type Method int

const (
	methodUnknown Method = iota
	methodDomains
	methodList
	methodAdd
	methodEdit
)

// Возвращает имя метода для использования в URL.
func (m Method) String() string {
	return [...]string{
		"unknown",
		"domains",
		"list",
		"add",
		"edit",
	}[m]
}

func (m Method) HTTPMethod() string {
	if m == methodUnknown {
		return methodUnknown.String()
	}

	switch m {
	case methodDomains, methodList:
		return http.MethodGet
	default:
		return http.MethodPost
	}
}

func (m Method) IsUnknown() bool {
	return m == methodUnknown
}
