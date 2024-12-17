package service

const (
	possibleChars = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890!@#$%^&*()-=_+[]{};<>?"
)

var (
	keyPrefixes   = []string{}
	valuePrefixes = []string{"encrypt-"}
	decryptPrefix = "decrypt-"

	encKeyLength = 32
	maxRetry     = 20

	pgErrConstrationID = "#23505"
	hashConstraintKey  = "_key"
	tokenConstraintKey = "_pkey"
)
