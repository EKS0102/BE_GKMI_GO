package database

import (
	"context"
	"fmt"
)

func TestConnection() error {

	conn, err := Connect()
	if err != nil {
		return err
	}

	defer conn.Close(context.Background())

	err = conn.Ping(context.Background())
	if err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	fmt.Println("PostgreSQL connection successful")

	return nil
}
