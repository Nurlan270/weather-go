package view

type BaseViewData struct {
	Title string
	Login string
	Old   OldValues
	Error ErrorData
}

func NewBaseViewData(title, login string) *BaseViewData {
	return &BaseViewData{
		Title: title,
		Login: login,
		Old:   make(OldValues),
	}
}

type ErrorViewData struct {
	Title   string
	Login   string
	Message string
}

func NewErrorViewData(title, message, login string) *ErrorViewData {
	return &ErrorViewData{
		Title:   title,
		Login:   login,
		Message: message,
	}
}

type ErrorData struct {
	Message string
	Items   ErrorMessages
}
