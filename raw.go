package goje

import (
	"strings"
)

// RawDelete Deletes entries with standard query
// This method dosen't support After,Before Triggers ...
func (handler *ContextHandler) RawDelete(Tablename string, Queries []QueryInterface) (int64, error) {
	query, args, err := ArgumentLessQueryBuilder(Delete, Tablename, nil, Queries)
	if err != nil {
		return -1, err
	}

	res, err := handler.DB.ExecContext(handler.Ctx, query, args...)
	if err != nil {
		return -1, err
	}

	return res.RowsAffected()
}

// RawUpdate update entries by map
// This method dosen't support After,Before Triggers ...
func (handler *ContextHandler) RawUpdate(Tablename string, Cols map[string]interface{}, Queries []QueryInterface) (int64, error) {
	if len(Cols) == 0 {
		return -1, ErrNoColsSetForUpdate
	}
	query := Update + " " + Tablename + " SET "
	var args []interface{}
	var items []string
	for key, val := range Cols {
		items = append(items, Tablename+"."+key+" = ?")
		args = append(args, val)
	}

	conditions, cargs, err := SQLConditionBuilder(Tablename, Queries)
	if err != nil {
		return -1, err
	}
	args = append(args, cargs...)

	res, err := handler.DB.ExecContext(handler.Ctx, query+strings.Join(items, ",")+conditions, args...)
	if err != nil {
		return -1, err
	}

	return res.RowsAffected()
}
