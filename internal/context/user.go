package context

// вообще это скорее не доменная сущность а вариант реализации
// пока не придумал куда это приткнуть, пусть тут полежит
type ContextKey string

const ContextUserKey ContextKey = "UserID"
