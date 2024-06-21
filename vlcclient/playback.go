package vlcclient

import "fmt"

func (c Client) PlayPause() error {
	err := c.Do("/requests/status.json", map[string]string{
		"command": "pl_pause",
	}, nil)
	if err != nil {
		return fmt.Errorf("PlayPause: %w", err)
	}
	// Success!
	return nil
}
