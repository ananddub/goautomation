package process

import (
	"context"
	"database/sql"
	"fmt"
	"goautomation/client"
	"goautomation/db/sqlc"
)

const (
	width  = 0.62
	height = 0.92
)

func InsertWindow(rmq *client.AutomationRPCClient, s *sqlc.Queries, ctx context.Context) (err error) {
	resp, err := rmq.GetPixelColor(
		width,
		height,
	)
	if err != nil {
		return err
	}

	for windowID, data := range resp.Windows {
		_, err := s.CreateAutomationLog(
			ctx,
			sqlc.CreateAutomationLogParams{
				Name:       rmq.User,
				Title:      windowID,
				ActionType: "get_pixel_color",
				XPosition:  sql.NullFloat64{Float64: data.X, Valid: true},
				YPosition:  sql.NullFloat64{Float64: data.Y, Valid: true},
				PixelColor: sql.NullString{
					String: fmt.Sprintf("(%d, %d, %d)", data.RGB.R, data.RGB.G, data.RGB.B),
					Valid:  true,
				},
				Success: sql.NullBool{Bool: resp.Success, Valid: true},
			},
		)
		if err != nil {
			return err
		}
	}
	return nil
}
