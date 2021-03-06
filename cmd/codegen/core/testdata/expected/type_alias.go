// Code generated by rest/model/codegen.go. DO NOT EDIT.

package model

import "github.com/evergreen-ci/evergreen/rest/model"

type APIStructWithAliased struct {
	Foo *string `json:"foo"`
	Bar *string `json:"bar"`
}

// APIStructWithAliasedBuildFromService takes the model.StructWithAliased DB struct and
// returns the REST struct *APIStructWithAliased with the corresponding fields populated
func APIStructWithAliasedBuildFromService(t model.StructWithAliased) *APIStructWithAliased {
	m := APIStructWithAliased{}
	m.Foo = StringStringPtr(t.Foo)
	m.Bar = StringStringPtr(t.Bar)
	return &m
}

// APIStructWithAliasedToService takes the APIStructWithAliased REST struct and returns the DB struct
// *model.StructWithAliased with the corresponding fields populated
func APIStructWithAliasedToService(m APIStructWithAliased) *model.StructWithAliased {
	out := &model.StructWithAliased{}
	out.Foo = StringPtrString(m.Foo)
	out.Bar = StringPtrString(m.Bar)
	return out
}
