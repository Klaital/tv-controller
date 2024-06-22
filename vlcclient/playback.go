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

func (c Client) TrackAhead() error {
	err := c.Do("/requests/status.json", map[string]string{
		"command": "pl_next",
	}, nil)
	if err != nil {
		return fmt.Errorf("TrackAhead: %w", err)
	}
	// Success!
	return nil
}

func (c Client) TrackBack() error {
	err := c.Do("/requests/status.json", map[string]string{
		"command": "pl_previous",
	}, nil)
	if err != nil {
		return fmt.Errorf("TrackBack: %w", err)
	}
	// Success!
	return nil
}

func (c Client) Loop() error {
	err := c.Do("/requests/status.json", map[string]string{
		"command": "pl_loop",
	}, nil)
	if err != nil {
		return fmt.Errorf("Loop: %w", err)
	}
	// Success!
	return nil
}

func (c Client) Random() error {
	err := c.Do("/requests/status.json", map[string]string{
		"command": "pl_random",
	}, nil)
	if err != nil {
		return fmt.Errorf("Random: %w", err)
	}
	// Success!
	return nil
}
