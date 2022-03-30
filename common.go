package sql_builder

import "strconv"

type Itemable interface {
	// GetVal -- return value from DbItemable::field
	GetVal(field string) interface{}
	// IsNeedUpdate -- Is this data item need ON DUPLICATE KEY UPDATE?
	IsNeedUpdate() bool
}

// Rowsetable -- интерфейс одной таблички БД, []DbItemable?
type Rowsetable interface {
	// GetTable -- возвращает имя таблички из константы
	GetTable() string
	// GetAlias -- возвращает алиас на табличку из константы
	GetAlias() string
	// GetInsertMode -- return string {INTO | IGNORE}
	GetInsertMode() string
	// SetInsertMode -- set {IGNORE|INTO} mode
	SetInsertMode( mode string )
	// GetItem -- return DbItemable data struct from args
	// GetItem(num int, args ...interface{}) Itemable
	// GetConnection -- return DB connection
	GetConnection() Sqlable
}

type Paginator struct{
	AllRows uint64
	Page    uint
	Onpage  uint
}

func (pages Paginator) Limit() string {
	return strconv.Itoa( int((pages.Page-1)*pages.Onpage) ) + "," + strconv.Itoa( int(pages.Onpage) )
}
