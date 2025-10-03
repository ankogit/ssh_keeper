package ssh

import "ssh-keeper/internal/models"

// SSHClientInterface определяет интерфейс для SSH клиентов
type SSHClientInterface interface {
	Connect() error
	GetConnectionString() string
}

// ClientFactory создает соответствующий SSH клиент на основе типа аутентификации
type ClientFactory struct{}

// NewClientFactory создает новую фабрику клиентов
func NewClientFactory() *ClientFactory {
	return &ClientFactory{}
}

// CreateClient создает SSH клиент на основе типа аутентификации
func (cf *ClientFactory) CreateClient(conn *models.Connection) SSHClientInterface {
	if conn.HasPassword {
		return NewPasswordClient(conn)
	}
	return NewKeyClient(conn)
}

// CreateKeyClient создает клиент для аутентификации по ключу
func (cf *ClientFactory) CreateKeyClient(conn *models.Connection) *KeyClient {
	return NewKeyClient(conn)
}

// CreatePasswordClient создает клиент для аутентификации по паролю
func (cf *ClientFactory) CreatePasswordClient(conn *models.Connection) *PasswordClient {
	return NewPasswordClient(conn)
}
