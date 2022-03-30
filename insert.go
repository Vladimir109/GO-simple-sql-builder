package sql_builder

import (
	"context"
	"fmt"
	"strings"
)

const(
	INS_INTO   = "INTO"
	INS_IGNORE = "IGNORE"
)

func ValArgsInserts(dbRepo Rowsetable, fields []string, iterator func() Itemable) (query string, args []interface{}) {
	var addon string
	var isUpdating = false

	for {
		d := iterator()
		if d == nil {
			break
		}

		valPart := "("
		addon = "\nON DUPLICATE KEY UPDATE "
		for _, f := range fields {
			valPart += "?,"
			addon += f + "=VALUES(" + f + "),"
			args = append(args, d.GetVal(f))
			if dbRepo.GetInsertMode() == INS_INTO {
				isUpdating = isUpdating || d.IsNeedUpdate()
			}
		}
		query += valPart[:len(valPart)-1] + "),"
	}
	query = query[:len(query)-1]
	if isUpdating {
		query += addon[:len(addon)-1]
	}

	query = "INSERT " + dbRepo.GetInsertMode() + " " + dbRepo.GetTable() +
		"\n(" + strings.Join(fields, ",") + ") VALUES " + query

	return query, args
}

func SqlInsert(ctx context.Context, repo Rowsetable, query string, args ...interface{}) (newId uint64, err error) {
	var method = "SqlInsert"

	res, err := repo.GetConnection().ExecContext(ctx, query, args...)

	if err != nil {
		return 0, fmt.Errorf("%s() inserting has: %s", method, err.Error())
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s(): LastInsertId() has: %s", method, err.Error())
	}

	rowsCount, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("%s(): RowsAffected() has: %s", method, err.Error())
	}

	if rowsCount != int64(len(args)) {
		return uint64(id), fmt.Errorf("%s() has affected rows are not equal items count", method)
	}

	return uint64(id), nil
}
