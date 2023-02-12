package spi

type QueryContext struct {
	DB           Database
	OnFetchStart func(Columns)
	OnFetch      func(rownum int64, values []any) bool
	OnFetchEnd   func()
}

func DoQuery(ctx *QueryContext, sqlText string, args ...any) error {
	rows, err := ctx.DB.Query(sqlText, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	if !rows.IsFetchable() {
		return nil
	}

	cols, err := rows.Columns()
	if err != nil {
		return err
	}
	if ctx.OnFetchStart != nil {
		ctx.OnFetchStart(cols)
	}
	if ctx.OnFetchEnd != nil {
		defer ctx.OnFetchEnd()
	}

	var nrow int64
	rec := cols.MakeBuffer()
	for rows.Next() {
		err = rows.Scan(rec...)
		if err != nil {
			return err
		}
		nrow++
		if ctx.OnFetch != nil && !ctx.OnFetch(nrow, rec) {
			break
		}
	}
	return nil
}
