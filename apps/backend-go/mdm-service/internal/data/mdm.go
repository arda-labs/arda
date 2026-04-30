package data

import (
	"context"
	"fmt"

	"github.com/arda-labs/arda/arda-be-go/services/mdm-service/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
)

type MdmRepo struct {
	data *Data
	log  *log.Helper
}

func NewMdmRepo(data *Data, logger log.Logger) biz.MdmRepo {
	return &MdmRepo{data: data, log: log.NewHelper(logger)}
}

func softDelete(ctx context.Context, data *Data, table, id string) error {
	tag, err := data.db.Pool.Exec(ctx,
		fmt.Sprintf(`UPDATE %s SET deleted_at = now(), status = 'DELETED', updated_at = now() WHERE id = $1 AND deleted_at IS NULL`, table),
		id,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return biz.ErrNotFound
	}
	return nil
}
