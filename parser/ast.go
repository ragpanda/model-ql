package parser

import (
	"context"

	"github.com/spf13/cast"
)

type CompileUnit struct {
	ViewList []*View

	viewMap map[string]*View
}

func (self *CompileUnit) process(ctx context.Context) error {
	self.viewMap = make(map[string]*View, 0)
	for _, view := range self.ViewList {
		self.viewMap[view.Ident] = view
		err := view.process(ctx)
		if err != nil {
			return err
		}
	}

	return nil

}

type View struct {
	Ident  string
	Select *SelectStmt
}

func (self *View) process(ctx context.Context) error {
	return self.Select.process(ctx)
}

type SelectStmt struct {
	FieldList    []*SelectField
	View         string
	JoinClauses  []*Join
	WhereClauses *WhereList

	// IR

	// Field
	QueryFromView func() *View
	GetViewType   func() *Type
}

func (self *SelectStmt) process(ctx context.Context) error {

	hasMatchAll := false
	for _, f := range self.FieldList {
		if f.Field == "*" {
			hasMatchAll = true
		}
	}

	if hasMatchAll {
		return NewSemanticError("Match can't exist with field on same time")
	}

	self.QueryFromView = ResolveView(ctx, self.View)
	self.GetViewType = func() *Type {
		originView := self.QueryFromView()
		t := &Type{
			Name:        "SelectStmtVirtal_" + self.View + "_" + cast.ToString(GetCurrentCount()),
			TEnum:       CustomStruct,
			Field:       []*FieldItem{},
			Enum:        []*EnumItem{},
			KeyType:     &Type{},
			ValueType:   &Type{},
			Annotations: []*Annotation{},
		}

		parentView := originView.Select.GetViewType()

		if hasMatchAll {
			t.Field = parentView.Field
		} else {
			// Compare with select field
			for _, selectField := range self.FieldList {
				for _, originModelField := range parentView.Field {
					if selectField.Field == originModelField.Name {
						t.Field = append(t.Field, originModelField) // TODO do Copy
					}
				}
			}
		}

		return t
	}

	for _, joinClause := range self.JoinClauses {
		joinClause.process(ctx)
	}

	return nil

}

type Join struct {
	JoinTarget string
	Alias      *string
	Self       MiltiLevelField
	CompareOp  string
	Ref        MiltiLevelField

	// IR
	JoinTargetRef *View
}

func (self *Join) process(ctx context.Context) error {
	return nil
}
