package validators

import (
	"fmt"
)

func Min(limit int) func(v string) error {
	return func(v string) error {
		if len(v) < limit {
			return fmt.Errorf("the value must be at least %d characters long", limit)
		}
		return nil
	}
}

func Max(limit int) func(v string) error {
	return func(v string) error {
		if len(v) > limit {
			return fmt.Errorf("the value is limited to %d characters", limit)
		}
		return nil
	}
}

func ItemName(name string) error {
	if err := Min(3)(name); err != nil {
		return err
	}
	return Max(50)(name)
}
