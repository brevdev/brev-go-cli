package brev_errors

import "fmt"

type CredentialsFileNotFound struct{}

func (e *CredentialsFileNotFound) Error() string {
	return fmt.Sprintf("Credentials file not found")
}

type GlobalProjectPathsFileNotFound struct{}

func (e *GlobalProjectPathsFileNotFound) Error() string {
	return fmt.Sprintf("Global project paths file not found")
}

type LocalProjectFileNotFound struct{}

func (e *LocalProjectFileNotFound) Error() string {
	return fmt.Sprintf("Local project file not found")
}

type LocalEndpointFileNotFound struct{}

func (e *LocalEndpointFileNotFound) Error() string {
	return fmt.Sprintf("Local endpoint file not found")
}

type InitExistingProjectFile struct{}

func (e *InitExistingProjectFile) Error() string {
	return fmt.Sprintf("Init called in a directory with an existing project file")
}

type InitExistingEndpointsFile struct{}

func (e *InitExistingEndpointsFile) Error() string {
	return fmt.Sprintf("Init called in a directory with an existing endpoints file")
}

type CotterClientError struct{}

func (e *CotterClientError) Error() string {
	return fmt.Sprintf("Invalid refresh token reported by auth server")
}

type CotterServerError struct{}

func (e *CotterServerError) Error() string {
	return fmt.Sprintf("Internal error reported by auth server")
}