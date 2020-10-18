package mymigrate

func resetMigrations() {
	migrations = make(map[string]mig)
}

func resetAppliedFunc() {
	getApplied = defaultAppliedFunc
}

func resetMarkAppliedFunc() {
	markApplied = defaultMarkAppliedFunc
}

func resetDownFunc() {
	down = defaultDownFunc
}
