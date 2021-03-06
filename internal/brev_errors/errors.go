package brev_errors

type BrevError interface {

	// Error returns a user-facing string explaining the error
	Error() string

	// Directive returns a user-facing string explaining how to overcome the error
	Directive() string
}

type SuppressedError struct{}

func (e *SuppressedError) Directive() string {
	return ""
}

func (e *SuppressedError) Error() string {
	return ""
}

type CredentialsFileNotFound struct{}

func (e *CredentialsFileNotFound) Directive() string {
	return "run `brev login`"
}

func (e *CredentialsFileNotFound) Error() string {
	return "credentials file not found"
}

type LocalProjectFileNotFound struct{}

func (e *LocalProjectFileNotFound) Directive() string {
	return "run `brev init`"
}

func (e *LocalProjectFileNotFound) Error() string {
	return "local project file not found"
}

type LocalEndpointFileNotFound struct{}

func (e *LocalEndpointFileNotFound) Directive() string {
	return "run `brev init`"
}

func (e *LocalEndpointFileNotFound) Error() string {
	return "local endpoint file not found"
}

type InitExistingProjectFile struct{}

func (e *InitExistingProjectFile) Directive() string {
	return "move to a new directory or delete the local .brev directory"
}

func (e *InitExistingProjectFile) Error() string {
	return "`brev init` called in a directory with an existing project file"
}

type InitExistingEndpointsFile struct{}

func (e *InitExistingEndpointsFile) Directive() string {
	return "move to a new directory or delete the local .brev directory"
}

func (e *InitExistingEndpointsFile) Error() string {
	return "init called in a directory with an existing endpoints file"
}

type CotterClientError struct{}

func (e *CotterClientError) Directive() string {
	return "run `brev login`"
}

func (e *CotterClientError) Error() string {
	return "invalid refresh token reported by auth server"
}

type CotterServerError struct{}

func (e *CotterServerError) Directive() string {
	return "wait for 60 seconds and run `brev login`"
}

func (e *CotterServerError) Error() string {
	return "internal error reported by auth server"
}
