package client

func GetDSKeyStorePath(username string) string {
	return "DSKEY/" + username
}

func GetPKKeyStorePath(username string) string {
	return "PK/" + username
}
