package spi

import "strings"

// Describe retrieves the result of 'desc table'.
//
// If includeHiddenColumns is true, the result includes hidden columns those name start with '_'
// such as "_RID" and "_ARRIVAL_TIME".
func DoDescribe(db Database, name string, includeHiddenColumns bool) (Description, error) {
	d := &TableDescription{}
	var tableType int
	var colCount int
	var colType int
	r := db.QueryRow("select name, type, flag, id, colcount from M$SYS_TABLES where name = ?", strings.ToUpper(name))
	if err := r.Scan(&d.Name, &tableType, &d.Flag, &d.Id, &colCount); err != nil {
		return nil, err
	}
	d.Type = TableType(tableType)

	rows, err := db.Query(`
		select
			name, type, length, id
		from
			M$SYS_COLUMNS
		where
			table_id = ? 
		order by id`, d.Id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		col := &ColumnDescription{}
		err = rows.Scan(&col.Name, &colType, &col.Length, &col.Id)
		if err != nil {
			return nil, err
		}
		if !includeHiddenColumns && strings.HasPrefix(col.Name, "_") {
			continue
		}
		col.Type = ColumnType(colType)
		d.Columns = append(d.Columns, col)
	}
	return d, nil
}

type Description interface {
	description()
}

func (td *TableDescription) description()  {}
func (cd *ColumnDescription) description() {}

// TableDescription is represents data that comes as a result of 'desc <table>'
type TableDescription struct {
	Name    string               `json:"name"`
	Type    TableType            `json:"type"`
	Flag    int                  `json:"flag"`
	Id      int                  `json:"id"`
	Columns []*ColumnDescription `json:"columns"`
}

// TypeString returns string representation of table type.
func (td *TableDescription) TypeString() string {
	return TableTypeDescription(td.Type, td.Flag)
}

// TableTypeDescription converts the given TableType and flag into string representation.
func TableTypeDescription(typ TableType, flag int) string {
	desc := "undef"
	switch typ {
	case LogTableType:
		desc = "Log Table"
	case FixedTableType:
		desc = "Fixed Table"
	case VolatileTableType:
		desc = "Volatile Table"
	case LookupTableType:
		desc = "Lookup Table"
	case KeyValueTableType:
		desc = "KeyValue Table"
	case TagTableType:
		desc = "Tag Table"
	}
	switch flag {
	case 1:
		desc += " (data)"
	case 2:
		desc += " (rollup)"
	case 4:
		desc += " (meta)"
	case 8:
		desc += " (stat)"
	}
	return desc
}

// columnDescription represents information of a column info.
type ColumnDescription struct {
	Id     uint64     `json:"id"`
	Name   string     `json:"name"`
	Type   ColumnType `json:"type"`
	Length int        `json:"length"`
}

// TypeString returns string representation of column type.
func (cd *ColumnDescription) TypeString() string {
	return ColumnTypeDescription(cd.Type)
}

// ColumnTypeDescription converts ColumnType into string.
func ColumnTypeDescription(typ ColumnType) string {
	switch typ {
	case Int16ColumnType:
		return "int16"
	case Uint16ColumnType:
		return "uint16"
	case Int32ColumnType:
		return "int32"
	case Uint32ColumnType:
		return "uint32"
	case Int64ColumnType:
		return "int64"
	case Uint64ColumnType:
		return "uint64"
	case Float32ColumnType:
		return "float"
	case Float64ColumnType:
		return "double"
	case VarcharColumnType:
		return "varchar"
	case TextColumnType:
		return "text"
	case ClobColumnType:
		return "clob"
	case BlobColumnType:
		return "blob"
	case BinaryColumnType:
		return "binary"
	case DatetimeColumnType:
		return "datetime"
	case IpV4ColumnType:
		return "ipv4"
	case IpV6ColumnType:
		return "ipv6"
	default:
		return "undef"
	}
}
