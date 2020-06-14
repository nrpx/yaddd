package config

// Конфигурация сервиса.
type Config struct {
	// Ключ доступа к API Яндекс.Коннект.
	PddToken string `json:"pddToken" yaml:"pddToken"`
	// Имя домена для управления его DNS-записями.
	DomainName string `json:"domain" yaml:"domain"`
	// Частота обновления данных (в формате crontab).
	Cron string `json:"cron" yaml:"cron"`
	// Адрес эл. почты для уведомлений.
	Email string `json:"email" yaml:"email"`
	// Текущий IP-адрес для определения A-записи.
	CurrentIP string `json:"-" yaml:"-"`
}
