package types

import (
	"io"

	"github.com/lyraproj/issue/issue"
	"github.com/lyraproj/pcore/errors"
	"github.com/lyraproj/pcore/eval"
)

type TypeReferenceType struct {
	typeString string
}

var TypeReferenceMetaType eval.ObjectType

func init() {
	TypeReferenceMetaType = newObjectType(`Pcore::TypeReference`,
		`Pcore::AnyType {
	attributes => {
		type_string => String[1]
	}
}`, func(ctx eval.Context, args []eval.Value) eval.Value {
			return newTypeReferenceType2(args...)
		})
}

func DefaultTypeReferenceType() *TypeReferenceType {
	return typeReferenceTypeDefault
}

func NewTypeReferenceType(typeString string) *TypeReferenceType {
	return &TypeReferenceType{typeString}
}

func newTypeReferenceType2(args ...eval.Value) *TypeReferenceType {
	switch len(args) {
	case 0:
		return DefaultTypeReferenceType()
	case 1:
		if str, ok := args[0].(stringValue); ok {
			return &TypeReferenceType{string(str)}
		}
		panic(NewIllegalArgumentType(`TypeReference[]`, 0, `String`, args[0]))
	default:
		panic(errors.NewIllegalArgumentCount(`TypeReference[]`, `0 - 1`, len(args)))
	}
}

func (t *TypeReferenceType) Accept(v eval.Visitor, g eval.Guard) {
	v(t)
}

func (t *TypeReferenceType) Default() eval.Type {
	return typeReferenceTypeDefault
}

func (t *TypeReferenceType) Equals(o interface{}, g eval.Guard) bool {
	if ot, ok := o.(*TypeReferenceType); ok {
		return t.typeString == ot.typeString
	}
	return false
}

func (t *TypeReferenceType) Get(key string) (eval.Value, bool) {
	switch key {
	case `type_string`:
		return stringValue(t.typeString), true
	default:
		return nil, false
	}
}

func (t *TypeReferenceType) IsAssignable(o eval.Type, g eval.Guard) bool {
	tr, ok := o.(*TypeReferenceType)
	return ok && t.typeString == tr.typeString
}

func (t *TypeReferenceType) IsInstance(o eval.Value, g eval.Guard) bool {
	return false
}

func (t *TypeReferenceType) MetaType() eval.ObjectType {
	return TypeReferenceMetaType
}

func (t *TypeReferenceType) Name() string {
	return `TypeReference`
}

func (t *TypeReferenceType) CanSerializeAsString() bool {
	return true
}

func (t *TypeReferenceType) SerializationString() string {
	return t.String()
}

func (t *TypeReferenceType) String() string {
	return eval.ToString2(t, None)
}

func (t *TypeReferenceType) Parameters() []eval.Value {
	if *t == *typeReferenceTypeDefault {
		return eval.EmptyValues
	}
	return []eval.Value{stringValue(t.typeString)}
}

func (t *TypeReferenceType) Resolve(c eval.Context) eval.Type {
	r := c.ParseType2(t.typeString)
	if rt, ok := r.(eval.ResolvableType); ok {
		if tr, ok := rt.(*TypeReferenceType); ok && t.typeString == tr.typeString {
			panic(eval.Error(eval.UnresolvedType, issue.H{`typeString`: t.typeString}))
		}
		r = rt.Resolve(c)
	}
	return r
}

func (t *TypeReferenceType) ToString(b io.Writer, s eval.FormatContext, g eval.RDetect) {
	TypeToString(t, b, s, g)
}

func (t *TypeReferenceType) PType() eval.Type {
	return &TypeType{t}
}

func (t *TypeReferenceType) TypeString() string {
	return t.typeString
}

var typeReferenceTypeDefault = &TypeReferenceType{`UnresolvedReference`}