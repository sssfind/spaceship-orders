package order

import (
	"context"
	"fmt"

	"order/internal/model"
)

func (r *repo) Delete(ctx context.Context, orderUUID string) error {
	query := `DELETE FROM orders WHERE order_uuid = $1`

	cmdTag, err := r.db.Exec(ctx, query, orderUUID)
	if err != nil {
		return fmt.Errorf("repository: delete order: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return model.ErrOrderNotFound
	}

	return nil
}
