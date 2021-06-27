package parser

import "github.com/ragpanda/model-ql/util"

type Identifier string

type CompileUnit struct {
	ModelList []*Model
}

func parseCompile(statements interface{}) (interface{}, error) {
	compileUnit := &CompileUnit{}
	statementsList := toIfaceSlice(statements)
	for _, stmts := range statementsList {
		switch stmt := stmts.([]interface{})[0].(type) {
		case *Model:
			compileUnit.ModelList = append(compileUnit.ModelList, stmt)
		default:

		}
	}

	return compileUnit, nil
}

type Model struct {
	Ident  string
	Select *SelectStmt
}

func parseModel(ident, selectData interface{}) (interface{}, error) {
	model := &Model{
		Ident:  parseString(toFirst(ident)),
		Select: nil,
	}

	model.Select = selectData.(*SelectStmt)

	return model, nil
}

type SelectStmt struct {
	FieldList    []*SelectField
	Model        string
	JoinClauses  []*Join
	WhereClauses *Where
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
			util.Info(nil, "%s", util.Display(item))
		}
	})

	Iter(where, func(v interface{}) {
		switch item := v.(type) {
		case *Where:
			selectStmt.WhereClauses = item
		}
	})

	selectStmt.Model = parseString(toFirst(model))

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

type Join struct {
	JoinTarget string
	Alias      *string
	Self       MiltiLevelField
	CompareOp  string
	Ref        MiltiLevelField
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

type Where struct {
}
