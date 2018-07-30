package medb

import (
	"database/sql"
)

// Row 单行
type Row struct {
	*sql.Row
}

// Rows 数据行
type Rows struct {
	*sql.Rows
	err     error
	columns map[string]int // 对应数据库的列名和列序号
}

// Err 返回错误信息
func (r *Rows) Err() error {
	return r.err
}

func (r *Rows) ColumnsMap() (map[string]int, error) {
	if r.columns == nil {
		cols, err := r.Columns()
		if err != nil {
			return nil, err
		}
		r.columns = make(map[string]int, len(cols))
		for i, col := range cols {
			r.columns[col] = i
		}
	}
	return r.columns, nil
}

// ScanNext 组合scan和next
func (r *Rows) ScanNext(dest ...interface{}) error {
	if r.err != nil {
		return r.err
	}
	defer r.Close()
	if r.Next() {
		err := r.Scan(dest...)
		if err != nil {
			return err
		}
	}
	return nil
}

// ScanStrings 解析到字符串数组中
func (r *Rows) ScanStrings() ([][]string, error) {
	return r.scanStrings(true)
}

// ScanStrings 解析到字符串数组中,并且处理Null的情况
func (r *Rows) ScanStrings2() ([][]string, error) {
	return r.scanStrings(false)
}

// ScanStrings 解析到字符串数组中
func (r *Rows) scanStrings(handleNull bool) ([][]string, error) {
	if r.err != nil {
		return nil, r.err
	}
	colM, err := r.ColumnsMap()
	if err != nil {
		return nil, err
	}
	ret := make([][]string, 0)
	defer r.Close()

	if handleNull {
		for r.Next() {
			rd := make([]sql.NullString, len(colM))
			err = r.Scan(NullStringsPtrToInterfaces(rd)...)
			if err != nil {
				return nil, err
			}
			ret = append(ret, NullStringsPtrToStrings(rd))
		}
	} else {
		for r.Next() {
			rd := make([]string, len(colM))
			err = r.Scan(StringsPtrToInterfaces(rd)...)
			if err != nil {
				return nil, err
			}
			ret = append(ret, rd)
		}
	}
	return ret, nil
}

// ScanStringsOne 解析到字符串数组中,但是最多解析一条
func (r *Rows) ScanStringsOne() ([]string, error) {
	return r.scanStringsOne(true)
}

// ScanStringsOne2 解析到字符串数组中,但是最多解析一条
func (r *Rows) ScanStringsOne2() ([]string, error) {
	return r.scanStringsOne(false)
}

// scanStringsOne 解析到字符串数组中,但是最多解析一条
func (r *Rows) scanStringsOne(handleNull bool) ([]string, error) {
	if r.err != nil {
		return nil, r.err
	}
	colM, err := r.ColumnsMap()
	if err != nil {
		return nil, err
	}
	rd := make([]string, len(colM))
	defer r.Close()
	if r.Next() {
		if handleNull {
			nrd := make([]sql.NullString, len(colM))
			err = r.Scan(NullStringsPtrToInterfaces(nrd)...)
			rd = NullStringsPtrToStrings(nrd)
		} else {
			err = r.Scan(StringsPtrToInterfaces(rd)...)
		}
	}
	return rd, err
}

// ScanMap 解析到map中
func (r *Rows) ScanMap() (map[string][]string, error) {
	return r.scanMap(true)
}

// ScanMap 解析到map中
func (r *Rows) ScanMap2() (map[string][]string, error) {
	return r.scanMap(false)
}

// scanMap 解析到map中
func (r *Rows) scanMap(handleNull bool) (map[string][]string, error) {
	if r.err != nil {
		return nil, r.err
	}
	colM, err := r.ColumnsMap()
	if err != nil {
		return nil, err
	}
	ret := make(map[string][]string)
	defer r.Close()
	if handleNull {
		for r.Next() {

			rd := make([]sql.NullString, len(colM))
			err = r.Scan(NullStringsPtrToInterfaces(rd)...)
			if err != nil {
				return nil, err
			}
			for col, ind := range colM {
				ret[col] = append(ret[col], rd[ind].String)
			}
		}
	} else {
		for r.Next() {
			rd := make([]string, len(colM))
			err = r.Scan(StringsPtrToInterfaces(rd)...)
			if err != nil {
				return nil, err
			}
			for col, ind := range colM {
				ret[col] = append(ret[col], rd[ind])
			}
		}
	}
	return ret, nil
}

// ScanMapOne 解析到map中,但是最多解析一条
func (r *Rows) ScanMapOne() (map[string]string, error) {
	return r.scanMapOne(true)
}

// ScanMapOne 解析到map中,但是最多解析一条
func (r *Rows) ScanMapOne2() (map[string]string, error) {
	return r.scanMapOne(false)
}

// scanMapOne 解析到map中,但是最多解析一条
func (r *Rows) scanMapOne(handleNull bool) (map[string]string, error) {
	if r.err != nil {
		return nil, r.err
	}

	colM, err := r.ColumnsMap()
	if err != nil {
		return nil, err
	}
	ret := make(map[string]string)
	defer r.Close()
	if r.Next() {
		if handleNull {
			rd := make([]sql.NullString, len(colM))
			err = r.Scan(NullStringsPtrToInterfaces(rd)...)
			if err != nil {
				return nil, err
			}
			for col, ind := range colM {
				ret[col] = rd[ind].String
			}
		} else {
			rd := make([]string, len(colM))
			err = r.Scan(StringsPtrToInterfaces(rd)...)
			if err != nil {
				return nil, err
			}
			for col, ind := range colM {
				ret[col] = rd[ind]
			}
		}
	}
	return ret, nil
}
