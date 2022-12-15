package repo

import (
	"context"
	"ehdw/smartiko-test/src/logic/service"
	"ehdw/smartiko-test/src/logic/service/model"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type repo struct {
	pg *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) service.Repository {
	return repo{pg: db}
}

var createEmptyRecords = `
;with devid as (select id from devices where eui = $1)

insert into params_records (device_id, flag, change_time) VALUES 
((select * from devid), 1,  '0001-01-01T00:00:00Z'),
((select * from devid), 2,  '0001-01-01T00:00:00Z'),
((select * from devid), 3,  '0001-01-01T00:00:00Z')
ON CONFLICT DO NOTHING`

func (r repo) CreateDevice(ctx context.Context, devID string) (int, error) {
	batch := &pgx.Batch{}
	batch.Queue(`insert into devices (eui) VALUES ($1) ON CONFLICT DO NOTHING RETURNING id`, devID)
	batch.Queue(createEmptyRecords, devID)

	tx, _ := r.pg.Begin(ctx)
	defer tx.Rollback(ctx)
	res := tx.SendBatch(ctx, batch)
	defer res.Close()
	var id int
	err := res.QueryRow().Scan(&id)
	if err != nil && err.Error() != pgx.ErrNoRows.Error() {
		return 0, err
	}
	res.Exec()
	res.Close()
	tx.Commit(ctx)
	return id, nil
}

func (r repo) GetDevice(ctx context.Context, devID string) (model.Device, error) {
	query := `select d.eui, p.flag, p.value, p.change_time from devices as d inner join params_records as p on p.device_id = d.id where d.eui = $1`
	rows, _ := r.pg.Query(ctx, query, devID)
	defer rows.Close()
	d := model.Device{ID: devID}
	for rows.Next() {
		f := &model.Flag{}
		rows.Scan(&d.ID, &f.Number, &f.Value, &f.ChangeTime)
		d.Flags = append(d.Flags, f)
	}
	return d, nil
}

func (r repo) GetDevices(ctx context.Context) ([]model.Device, error) {
	return []model.Device{}, service.ErrNotImplemented
}

func (r repo) DeleteDevice(ctx context.Context, devID string) error {
	tx, _ := r.pg.Begin(ctx)
	defer tx.Rollback(ctx)

	_, err := tx.Exec(ctx, `delete from params_records where device_id = (select id from devices where eui = $1);`, devID)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `delete from devices where eui = $1`, devID)
	if err != nil {
		return err
	}
	tx.Commit(ctx)
	return nil
}

func (r repo) ModifyFlags(ctx context.Context, devID string, flags []*model.Flag) error {
	query := `update params_records set change_time = $1, "value" = $2 where flag = $3 and device_id = (select id from devices where eui = $4)`
	batch := &pgx.Batch{}
	for _, f := range flags {
		batch.Queue(query, f.ChangeTime, f.Value, f.Number, devID)
	}
	tx, _ := r.pg.Begin(ctx)
	defer tx.Rollback(ctx)
	res := tx.SendBatch(ctx, batch)
	defer res.Close()
	for i := 0; i < batch.Len(); i++ {
		_, err := res.Exec()
		if err != nil {
			return err
		}
	}
	res.Close()
	tx.Commit(ctx)
	return nil
}

func (r repo) GetAllEnabledDevices(ctx context.Context) ([]string, error) {
	query := `select eui from devices`
	rows, err := r.pg.Query(ctx, query)
	if err != nil && err.Error() != pgx.ErrNoRows.Error() {
		return nil, err
	}
	defer rows.Close()
	var res []string
	for rows.Next() {
		var tmp string
		err := rows.Scan(&tmp)
		if err != nil {
			return nil, err
		}
		res = append(res, tmp)
	}
	return res, nil
}
