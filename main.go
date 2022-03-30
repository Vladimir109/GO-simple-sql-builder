package main

import (
	"fmt"
	sb "gitlab.private.kvado.ru/golib/sql-builder/sql_builder"
	"time"
)

// ------------------------ Itemable class for sql table `example1`: ------------------------------- //

const(
	Ex1fieldId = "id"
	Ex1fieldName = "name"
	Ex1fieldDate = "for_date"
)

type Ex1rowDB struct{
	Id      uint64    `db:"id"`
	Name    string    `db:"name"`
	ForDate time.Time `db:"for_date"`
}

func (item *Ex1rowDB) IsNeedUpdate() bool {
	if item.Id > 0 {
		return true
	}
	return false
}

func (item *Ex1rowDB) GetVal( field string ) interface{} {
	switch field {
	case Ex1fieldId:   return item.Id
	case Ex1fieldName: return item.Name
	case Ex1fieldDate: return item.ForDate
	}
	return nil
}

type Ex1Rowset struct{
	sb.DbMock
}
func NewEx1Rowset() *Ex1Rowset { return &Ex1Rowset{} }

func (e Ex1Rowset) GetTable()      string     { return "example1"  }
func (e Ex1Rowset) GetAlias()      string     { return "ex1"       }
func (e Ex1Rowset) GetInsertMode() string     { return sb.INS_INTO }
func (e Ex1Rowset) GetConnection() sb.Sqlable { return nil         }
func (e Ex1Rowset) SetInsertMode(mode string) {}

// ------------------------ Itemable class for sql table `example2`: ------------------------------- //

const(
	Ex2fieldId    = "id"
	Ex2fieldFkEx1 = "fkey_ex1"
	Ex2fieldWare  = "ware"
	Ex2fieldCost  = "cost"
)

type Ex2rowDB struct{
	Id       uint64  `db:"id"`
	FkeyEx1  uint64  `db:"fkey_ex1"`
	WareName string  `db:"ware"`
	CostVal  float64 `db:"cost"`
}

func (item *Ex2rowDB) IsNeedUpdate() bool {
	if item.Id > 0 {
		return true
	}
	return false
}

func (item *Ex2rowDB) GetVal( field string ) interface{} {
	switch field {
	case Ex2fieldId:    return item.Id
	case Ex2fieldFkEx1: return item.FkeyEx1
	case Ex2fieldWare:  return item.WareName
	case Ex2fieldCost:  return item.CostVal
	}
	return nil
}

type Ex2Rowset struct{
	sb.DbMock
}
func NewEx2Rowset() *Ex2Rowset { return &Ex2Rowset{} }

func (e Ex2Rowset) GetTable()      string     { return "sells"     }
func (e Ex2Rowset) GetAlias()      string     { return "s"         }
func (e Ex2Rowset) GetInsertMode() string     { return sb.INS_INTO }
func (e Ex2Rowset) GetConnection() sb.Sqlable { return nil         }
func (e Ex2Rowset) SetInsertMode(mode string) {}

func main(){

	dbEx1 := NewEx1Rowset()
	dbEx2 := NewEx2Rowset()

	Example(dbEx1, dbEx2)
}

func Example(tbl1,tbl2 sb.Rowsetable){
	item := Ex1rowDB{ Id: 1, Name: "Pupkin Vasiliy", ForDate: time.Now() } // inserting twice, see iterator!
	num := 0
	sql,args := sb.ValArgsInserts(tbl1, []string{Ex1fieldId, Ex1fieldName, Ex1fieldDate},
		func() sb.Itemable {
			if num == 2 { return nil }
			num++
			return &item
	})
	fmt.Println(fmt.Sprintf("\nsql: %s\nargs:%#v", sql, args))

	query := sb.NewQuery(tbl1, "t1")
	query[0].SELECT( Ex1fieldId, Ex1fieldName ).DISTINCT().
		F_WHERE(
			sb.NewFilter(Ex1fieldDate, sb.Eq, "2022-02-24"),
			sb.NewFilter(Ex1fieldName, sb.Like, "Pupk%"),
		).
		GROUPBY( Ex1fieldDate ).
		HAVING("(COUNT(t1.id) > ?)", sb.AsIs, 3)

	query.JOIN(sb.LEFT, tbl2, "tj").ON("t1.id = tj.fkey_ex1", sb.AsIs).
		SubQuery(tbl1, "t2").
		COLUMN(sb.COUNT("t2", Ex2fieldWare), "wares_count").
		COLUMN("MAX(t2.cost)", "maxCost").
		GROUPBY( Ex2fieldFkEx1 ).ROLLUP().
		HAVING( Ex2fieldWare, sb.In, "ware1","ware2","ware3")

	query.AddFROM(tbl2, "sub").SELECT("name").WHERE("sub.id IS NOT NULL", sb.AsIs).
		SubQuery(tbl2, "s2").SELECT("id","name").WHERE("id", sb.Gt, 3)

	fmt.Println(fmt.Sprintf("\n\nsql: %s\nargs:%#v", query.String(), query.MakeArgs()))

	qpart := &sb.QueryPart{}
	qpart.SetAlias("qq").
		WHERE("name", sb.Like, "%nam%").
		ORDERBY("qq.id ASC", "qq.name DESC").
		LIMIT(sb.Paginator{ Page: 2, Onpage:  24 })
	query = sb.Query{ *qpart }
	fmt.Println(fmt.Sprintf("\n\nsql2: %s\nargs:%#v", query.String(), query.MakeArgs()))
}
