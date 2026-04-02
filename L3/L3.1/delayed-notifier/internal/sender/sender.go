package sender

import (
	"errors"
	"fmt"
)

func Send(channel, to, message string) error {
	switch channel {
	case "email":
		fmt.Println("EMAIL:", to, message)
	case "telegram":
		fmt.Println("TG:", to, message)
	default:
		return errors.New("unknown channel")
	}

	return nil
}
