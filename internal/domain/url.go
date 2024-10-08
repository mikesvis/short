// Модуль доменных сущностей.
package domain

// ID в виде строки.
type ID string

// Сущность URL в домене.
type URL struct {
	// ID пользователя.
	UserID string

	// Полный URL.
	Full string

	// Короткий ключ.
	Short string

	// Флаг удаленного элемента.
	Deleted bool
}
