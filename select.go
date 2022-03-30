package sql_builder

import (
	"strings"
)

const(
	FieldSplitter = ","
	FieldDot = "."

	SELECT    = "SELECT "
	DISTINCT  = "DISTINCT "
	FROM      = "\nFROM "
	JOIN      = "\nINNER JOIN "
	INNER     = JOIN
	LEFT      = "\nLEFT JOIN "
	RIGHT     = "\nRIGHT JOIN "
	LEFT_OUT  = "\nLEFT OUTER JOIN "
	RIGHT_OUT = "\nRIGHT OUTER JOIN "
	WHERE     = "\nWHERE "
	GROUPBY   = "\nGROUP BY "
	HAVING    = "\nHAVING "
	ORDERBY   = "\nORDER BY "
	LIMIT     = "\nLIMIT "

	WITH_ROLLUP = " WITH ROLLUP"
)

// QueryPart -- simple sql query or part complex
// Has all sql parts for query. May be full simple query at this Rowset type
// Each Rowsettable releasing type objects disribes only one, own sql table
type QueryPart struct{
	fromType   string
	from       Rowsetable
	subquery   Query
	alias      string
	selects    []string
	wheres     []string
	groups     []string
	havings    []string
	ons        []string
	orders     string
	limits     string
	selectArgs []interface{}
	whereArgs  []interface{}
	havingArgs []interface{}
	onArgs     []interface{}
	isRollup   bool
	isDistinct bool
}

// SetAlias -- set alias for this Rowset type object
func (qp *QueryPart) SetAlias( alias string ) *QueryPart {
	qp.alias = alias
	return qp
}

// AliasedFields -- make string alias.field
func (qp *QueryPart) AliasedFields( fields ...string ) string {
	var addings = make([]string, len(fields))
	for i,f := range fields {
		addings[i] = qp.alias + FieldDot + f
	}
	return strings.Join(addings, FieldSplitter)
}

// SELECT -- add to select part list field names from this Rowset (table)
func (qp *QueryPart) SELECT( fields ...string ) *QueryPart {
	qp.selects = append(qp.selects, qp.AliasedFields(fields...))
	return qp
}

// DISTINCT -- add this sentense when SQl query will be stringed
func (qp *QueryPart) DISTINCT() *QueryPart {
	qp.isDistinct = true
	return qp
}

// COLUMN -- add "as is" string expression (sub select query too) as column to select part
func (qp *QueryPart) COLUMN( exp string, asName string, args ...interface{}) *QueryPart {
	qp.selects = append(qp.selects, "("+exp+") AS " + asName)
	qp.selectArgs = append(qp.selectArgs, args...)
	return qp
}

// WHERE -- add one condition to where sql part with this part alias or "as is"
// all where parts combined only with "AND" when query stringed.
// If need "OR", try Cond() before and insert into this "as is" where part.
func (qp *QueryPart) WHERE( field, cond string, args ...interface{}) *QueryPart {
	var resCond Condition
	if cond == AsIs {
		// as is: already prepared conditions with its args
		resCond.Cond = field
		resCond.Args = args
	}else{
		// make all condition
		resCond = Cond(qp.alias, field, cond, args...)
	}
	qp.wheres = append(qp.wheres, resCond.Cond)
	qp.whereArgs = append(qp.whereArgs, resCond.Args...)
	return qp
}

// F_WHERE -- alternate adding condition into where sql part from prepared Filter objects.
// This is useful to pass Filters from uplevel into this builder.
// all where parts combined only with "AND" when query stringed.
// If need "OR", try Cond() before and insert into this "as is" where part.
func (qp *QueryPart) F_WHERE(filters ...Filter) *QueryPart {
	for _,f := range filters {
		_ = qp.WHERE(f.Field, f.Cond, f.Args...)
	}
	return qp
}

// GROUPBY -- add field list to group by sql part with alias from this part
func (qp *QueryPart) GROUPBY( fields ...string ) *QueryPart {
	qp.groups = append(qp.groups, qp.AliasedFields( fields... ))
	return qp
}

// ROLLUP -- add roll up sentecse to group by part when query.String() called
func (qp *QueryPart) ROLLUP() *QueryPart {
	qp.isRollup = true
	return qp
}

// HAVING -- add conditions into having sql part.
// all where parts combined only with "AND" when query stringed.
// If need "OR", try Cond() before and insert into this "as is" where part.
func (qp *QueryPart) HAVING( field, cond string, args ...interface{}) *QueryPart {
	if cond == AsIs {
		// as is: prepared conditions with its args
		qp.havings = append(qp.havings, field)
		qp.havingArgs = append(qp.havingArgs, args...)
		return qp
	}

	where := Cond(qp.alias, field, cond, args...)

	qp.havings = append( qp.havings, where.Cond)
	qp.havingArgs = append(qp.havingArgs, args...)
	return qp
}

// F_HAVING -- alternate adding condition into having sql part from prepared Filter objects.
// This is useful to pass Filters from uplevel into this builder.
// all where parts combined only with "AND" when query stringed.
// If need "OR", try Cond() before and insert into this "as is" where part.
func (qp *QueryPart) F_HAVING(filters ...Filter) *QueryPart {
	for _,f := range filters {
		_ = qp.HAVING(f.Field, f.Cond, f.Args...)
	}
	return qp
}

// ORDERBY -- add ordering strings "as is" into this sql part
func (qp *QueryPart) ORDERBY(orders ...string) *QueryPart {
	if qp.orders != "" {
		qp.orders += FieldSplitter
	}
	qp.orders += strings.Join(orders, FieldSplitter)
	return qp
}

// LIMIT -- calculating offset and limit from Paginator object in one place
func (qp *QueryPart) LIMIT(pages Paginator) *QueryPart {
	qp.limits = pages.Limit()
	return qp
}

// ON -- adding conditions into JOIN..ON sql part.
// base t1.f1 = t2.f2 from other QueryPart need adding "as is" variant.
// all where parts combined only with "AND" when query stringed.
// If need "OR", try Cond() before and insert into this "as is" where part.
func (qp *QueryPart) ON( field, cond string, args ...interface{}) *QueryPart {
	var on Condition

	on.Args = args
	if cond == AsIs {
		// as is: prepared conditions with its args
		on.Cond = field
	}else{
		on = Cond(qp.alias, field, cond, args...)
	}
	qp.ons = append( qp.ons, on.Cond)
	qp.onArgs = append(qp.onArgs, args...)
	return qp
}

// F_ON -- alternate adding condition into JOIN..ON sql part from prepared Filter objects.
// This is useful to pass Filters from uplevel into this builder.
// all where parts combined only with "AND" when query stringed.
// If need "OR", try Cond() before and insert into this "as is" where part.
func (qp *QueryPart) F_ON(filters ...Filter) *QueryPart {
	for _,f := range filters {
		_ = qp.ON(f.Field, f.Cond, f.Args...)
	}
	return qp
}

// SubQuery -- create internal Query as sub sql query in from or join sections
// returns own QueryPart for further add sql parts into SubQuery
func (qp *QueryPart) SubQuery(rowset Rowsetable, alias string) *QueryPart {
	if alias == "" {
		alias = rowset.GetAlias()
	}
	qp.subquery = append(qp.subquery, QueryPart{fromType: FROM, from: rowset, alias: alias} )
	return &(qp.subquery[len(qp.subquery)-1])
}

// Query -- list of QueryPart. This is a ful complex sql query
type Query []QueryPart

// NewQuery -- create Query from Rowsetable (table type objects) with one QueryPart and set alias for this QueryPart 
func NewQuery(rowset Rowsetable, alias string) Query {
	if alias == "" {
		alias = rowset.GetAlias()
	}
	return Query{ {fromType: FROM, from: rowset, alias: alias} }
}

// AddFROM -- add second Rowsetable part with own alias as from sql Decart multiply.
// Returns pointer for new part for the next use all QueryPart methods.
func (q *Query) AddFROM(rowset Rowsetable, alias string) *QueryPart {
	if alias == "" {
		alias = rowset.GetAlias()
	}
	*q = append(*q, QueryPart{fromType: FROM, from: rowset, alias: alias} )
	return &((*q)[len(*q)-1])
}

// JOIN -- add new Query item as joined QueryPart.
// Returns pointer for new part for the next use all QueryPart methods.
func (q *Query) JOIN(join string, rowset Rowsetable, alias string) *QueryPart {
	if alias == "" {
		alias = rowset.GetAlias()
	}
	*q = append(*q, QueryPart{fromType: join, from: rowset, alias: alias} )
	return &((*q)[len(*q)-1])
}

// String -- make sql query text from this Query.
// The all parts data are combined with each other SQL-parts.
// if any sql part is empty it does not include into result! This is feature for make only need part from
// ful sql sentense if need.
func (q Query) String() string {
	var selects,froms,joins,wheres,groups,havings,orders,limit string
	var isRollup bool
	var isDistinct bool
	for _,qp := range q {
		if len(qp.selects) > 0 {
			selects += strings.Join(qp.selects, FieldSplitter) + FieldSplitter
		}
		if len(qp.groups) > 0 {
			groups += strings.Join(qp.groups, FieldSplitter) + FieldSplitter
		}

		table := ""
		if qp.subquery != nil {
			table = "(\n" + qp.subquery.String() + "\n)"
		}else if qp.from != nil {
			table = qp.from.GetTable()
		}else{
			table = "!rowset or subquery is need!"
		}
		if table != "" {
			table += " " + qp.alias
		}
		if qp.fromType == FROM {
			froms += table + FieldSplitter
		} else if qp.fromType != "" {
			// anyfrom types is join kinds
			joins += qp.fromType + table + " ON "
			// ON WHERE собирается здесь же в joins!
			if len(qp.ons) > 0 {
				joins += strings.Join(qp.ons, AND)
			}
		}

		if len(qp.wheres) > 0 {
			if wheres != "" { // concat with previous where
				wheres += AND
			}
			wheres += " " + strings.Join(qp.wheres, AND)
		}
		if len(qp.havings) > 0 {
			if havings != "" { // concat with previous where
				havings += AND
			}
			havings += " " + strings.Join(qp.havings, AND)
		}
		if qp.orders != "" {
			orders += qp.orders + FieldSplitter
		}
		if qp.limits != "" {
			limit += qp.limits + FieldSplitter
		}
		if qp.isRollup && !isRollup {
			isRollup = true
		}
		if qp.isDistinct && !isDistinct {
			isDistinct = true
		}
	}
	selType := SELECT
	if isDistinct {
		selType += DISTINCT
	}

	if selects != "" { selects = selType + selects[:len(selects)-1] }
	if froms != ""   { froms   = FROM + froms[:len(froms)-1] }
	if wheres != ""  { wheres  = WHERE + wheres }
	if groups != ""  { groups  = GROUPBY + groups[:len(groups)-1] }
	if isRollup      { groups += WITH_ROLLUP }
	if havings != "" { havings = HAVING + havings }
	if orders != ""  { orders  = ORDERBY + orders[:len(orders)-1] }
	if limit != ""   { limit   = LIMIT + limit[:len(limit)-1] }

	return selects + froms + joins + wheres + groups + havings + orders + limit
}

// MakeArgs -- as String() are combines all arguments into []interface{} into sequence for sql parts
// popendicularity all QueryPart items. This is garants that all arguments will be on each places.
func (q Query) MakeArgs() []interface{} {
	var result,selects,subargs,ons,wheres,havings []interface{}
	for _,qp := range q {
		selects = append(selects, qp.selectArgs...)

		if len(qp.subquery)>0 {
			subs := qp.subquery.MakeArgs()
			if qp.fromType == FROM || qp.fromType == "" {
				// this is a subquery in from part (before join parts)
				subargs = append(subargs, subs...)
			}else{
				// join sub query args must be before ON argument part!
				ons = append(ons, subs...)
			}
		}
		ons     = append(ons, qp.onArgs...)
		wheres  = append(wheres, qp.whereArgs...)
		havings = append(havings, qp.havingArgs...)
	}
	result = append(result, selects...)
	result = append(result, subargs...)
	result = append(result, ons...)
	result = append(result, wheres...)
	result = append(result, havings...)
	return result
}
