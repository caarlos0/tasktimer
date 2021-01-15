package ui

type errMsg struct{ error }

func (e errMsg) Error() string { return e.error.Error() }
