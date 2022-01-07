package common

type APIChecker interface {
	Check(operation string, apikey string) (bool, string, error)
}

type AllowAll struct{}

func (a AllowAll) Check(operation string, apikey string) (bool, string, error) { return true, "", nil }
