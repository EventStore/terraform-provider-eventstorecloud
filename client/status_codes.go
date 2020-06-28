package client

import "fmt"

func translateStatusCode(status int, activity string) error {
	return fmt.Errorf("error %s", activity)
}

