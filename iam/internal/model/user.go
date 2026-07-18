package model

// NotificationMethod описывает канал связи, выбранный пользователем
type NotificationMethod struct {
	ProviderName string
	Target       string
}

// User - чистая доменная модель пользователя для бизнес-логики
type User struct {
	UUID                string
	Login               string
	PasswordHash        string
	Email               string
	NotificationMethods []NotificationMethod
}
