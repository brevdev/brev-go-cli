package brev_errors

import "fmt"

type CredentialsFileNotFound struct{}

func (e *CredentialsFileNotFound) Error() string {
	return fmt.Sprintf("No credentials file found")
}
