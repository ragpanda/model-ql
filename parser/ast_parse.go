package parser

import (
	"github.com/ragpanda/model-ql/util"
)


func parseCompile(statements interface{}) (interface{}, error) {
	compileUnit := &CompileUnit{}
	statementsList := toIfaceSlice(statements)
	for _, stmts := range statementsList {
		switch stmt := stmts.([]interface{})[0].(type) {
		case *View:
			compileUnit.ViewList = append(compileUnit.ViewList, stmt)
		default:

		}
	}

	return compileUnit, nil
}

func parseModel(ident, selectData interface{}) (interface{}, error) {
	model := &View{
		Ident:  parseString(toFirst(ident)),
		Select: nil,
	}

	model.Select = selectData.(*SelectStmt)

	return model, nil
}

func parseSelectStmt(fieldList, model, join, where interface{}) (interface{}, error) {
	selectStmt := &SelectStmt{}

	Iter(fieldList, func(v interface{}) {
		switch item := v.(type) {
		case *SelectField:
			selectStmt.FieldList = append(
				selectStmt.FieldList, item)
		}
	})

	Iter(join, func(v interface{}) {
		switch item := v.(type) {
		case *Join:
			selectStmt.JoinClauses = append(
				selectStmt.JoinClauses, item)
		default:
			util.Debug(nil, "%s", util.Display(item))
		}
	})

	Iter(where, func(v interface{}) {
		switch item := v.(type) {
		case *WhereList:
			selectStmt.WhereClauses = item
		}
	})

	selectStmt.View = parseString(toFirst(model))

	return selectStmt, nil
}

type SelectField struct {
	Field string
	Alias *string
}

func parseSelectField(name, alias interface{}) (interface{}, error) {
	sf := &SelectField{
		Field: parseString(name),
		Alias: nil,
	}

	if alias != nil {
		Iter(alias, func(x interface{}) {
			switch item := x.(type) {
			case Identifier:
				aliasStr := string(item)
				sf.Alias = &aliasStr
			}
		})

	}

	return sf, nil
}


func parseJoin(target, self, compare, ref interface{}) (interface{}, error) {
	t := target.(*JoinTarget)

	join := &Join{
		JoinTarget: t.Model,
		Alias:      t.Alias,
		Self:       *self.(*MiltiLevelField),
		CompareOp:  parseString(compare),
		Ref:        *ref.(*MiltiLevelField),
	}

	return join, nil
}

type JoinTarget struct {
	Model string
	Alias *string
}

func parseJoinTarget(name, alias interface{}) (interface{}, error) {

	jt := &JoinTarget{
		Model: parseString(name),
		Alias: nil,
	}

	if alias != nil {
		Iter(alias, func(x interface{}) {
			switch item := x.(type) {
			case Identifier:
				aliasStr := string(item)
				jt.Alias = &aliasStr
			}
		})

	}

	return jt, nil
}

type MiltiLevelField struct {
	Ident []string
}

func parseMiltiLevelField(idents interface{}) (interface{}, error) {
	r := &MiltiLevelField{}

	Iter(idents, func(x interface{}) {
		switch item := x.(type) {
		case Identifier:
			r.Ident = append(r.Ident, string(item))
		}
	})

	return r, nil
}

type WhereList struct {
	Must   []*WhereList
	Should []*WhereList
	Not    []*WhereList

	CompareOp *CompareOp
	Field     *MiltiLevelField
	Condition interface{}
}

func parseWhereList(first, tails interface{}) (interface{}, error) {
	where := &WhereList{
		Must:   []*WhereList{},
		Should: []*WhereList{},
		Not:    []*WhereList{},
	}

	if tails == nil {
		where.Must = append(where.Must, first.(*WhereList))
		return where, nil
	}

	var dispatchOp = func(op BoolOp, sub *WhereList) {
		switch op {
		case "and":
			where.Must = append(where.Must, sub)
		case "or":
			where.Should = append(where.Should, sub)
		case "not":
			where.Not = append(where.Not, sub)
		}
	}

	util.Debug(nil, "tails us: %s", util.Display(tails))

	op := BoolOp("")
	for _, item := range toIfaceSlice(tails) {

		Iter(item, func(item interface{}) {
			switch item := item.(type) {
			case BoolOp:
				util.Debug(nil, "BoolOp: %s", util.Display(item))
				if string(op) == "" {
					op = item
					dispatchOp(op, first.(*WhereList))
				}
			case *WhereList:
				dispatchOp(op, item)
			}

			util.Debug(nil, "parseWhereList: %s", util.Display(item))
		})

	}

	return where, nil
}

func parseWhereCondition(field, op, condition interface{}) (interface{}, error) {
	where := &WhereList{}
	opType := CompareOp(parseString(op))

	where.Field = field.(*MiltiLevelField)
	where.CompareOp = &opType
	where.Condition = condition

	util.Debug(nil, "parseWhereCondition: %s ", util.Display(where))
	return where, nil
}
