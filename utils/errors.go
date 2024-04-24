package utils

func IsSameError(err1, err2 error) bool {
	if err1 == nil || err2 == nil {
		return false
	}
	return err1.Error() == err2.Error()
}
