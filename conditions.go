package sql_builder

import "strings"

const(
	AsIs = "asis"
	Gt = " > "
	Ge = " >= "
	Eq = " = "
	Ne = " <> "
	Le = " <= "
	Lt = " < "
	In = " IN"
	NotIn = " NOT IN"
	Like = " LIKE "
	NotLike = " NOT LIKE "
	Exists = " EXISTS "
	NotExists = " NOT EXISTS "
	AND = " AND "
	OR = " OR "
)

type Condition struct{
	Cond string
	Args []interface{}
}

func Cond( alias, field, cond string, args ...interface{}) Condition {
	lenArgs := len(args)
	if alias != "" {
		alias += "." + field
	}
	where := alias + cond
	if lenArgs == 1 {
		where = alias + cond +"(?)"
	}else{
		q := strings.Repeat("?,", lenArgs)
		where = alias + cond +"("+q[:len(q)-1]+")"
	}
	where = "(" + where + ")"
	return Condition{ Cond: where, Args: args }
}

func AndOr( andOr string, conds ...Condition ) Condition {
	var resCond = Condition{}
	var where = make([]string, len(conds))
	for i,c := range conds {
		where[i] = c.Cond
		resCond.Args = append(resCond.Args, c.Args...)
	}
	resCond.Cond = "(" + strings.Join(where, ")" + andOr + "(") + ")"
	return resCond
}

func Or( conds ...Condition ) Condition {
	return AndOr(OR, conds...)
}

func And( conds ...Condition ) Condition {
	return AndOr(AND, conds...)
}

// Filter -- элемент фильтров данных
type Filter struct {
	Field string
	Cond  string
	Args  []interface{}
}

// NewFilter -- генератор фильтра из базовых деталек
func NewFilter(field, cond string, args ...interface{}) Filter {
	return Filter{
		Field: field,
		Cond:  cond,
		Args:  args,
	}
}

func (f Filter) AsCond(alias string) Condition {
	return Cond(alias, f.Field, f.Cond, f.Args...)
}
