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
func (handler *ContextHandler) RawUpdate(Tablename string, Cols map[string]interface{}, Queries ...QueryInterface) (int64, error) {
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

// RawBulkInsert insert multiple entries by []map[column name]value
// This method dosen't support After,Before Triggers ...
func (handler *ContextHandler) RawBulkInsert(Tablename string, Rows []map[string]interface{}) (int64, error) {
	if len(Rows) == 0 {
		return -1, ErrNoRowsForInsert
	}

	query := Insert + " INTO " + Tablename
	var args []interface{}
	var columnNames []string

	for index, row := range Rows {
		//use first index as column name index
		if index == 0 {
			for colName, _ := range row {
				columnNames = append(columnNames, colName)
			}
			if len(columnNames) == 0 {
				return -1, ErrNoRowsColsForInsert
			}
		}

		//put arguments attiontion to column names that fetched from index 0
		for _, colName := range columnNames {
			if arg, ok := row[colName]; ok {
				args = append(args, arg)
			} else {
				args = append(args, nil)
			}

		}
	}

	eachRowArgs := strings.Repeat(",?", len(columnNames))
	eachRowArgs = ",(" + eachRowArgs[1:] + ")"
	values := strings.Repeat(eachRowArgs, len(Rows))
	values = values[1:]

	res, err := handler.DB.ExecContext(handler.Ctx, query+"("+strings.Join(columnNames, ",")+") VALUES "+values, args...)
	if err != nil {
		return -1, err
	}

	return res.RowsAffected()
}
