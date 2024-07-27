package domain

type ID string

type URL struct {
	UserID  string
	Full    string
	Short   string
	Deleted bool
}
