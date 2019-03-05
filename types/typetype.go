package types

import (
	"io"

	"github.com/lyraproj/pcore/errors"
	"github.com/lyraproj/pcore/eval"
)

type TypeType struct {
	typ eval.Type
}

var typeTypeDefault = &TypeType{typ: anyTypeDefault}

var TypeMetaType eval.ObjectType

func init() {
	TypeMetaType = newObjectType(`Pcore::TypeType`,
		`Pcore::AnyType {
	attributes => {
		type => {
			type => Optional[Type],
			value => Any
		},
	}
}`, func(ctx eval.Context, args []eval.Value) eval.Value {
			return newTypeType2(args...)
		})

	newGoConstructor(`Type`,
		func(d eval.Dispatch) {
			d.Param(`String`)
			d.Function(func(c eval.Context, args []eval.Value) eval.Value {
				return c.ParseType(args[0])
			})
		},
		func(d eval.Dispatch) {
			d.Param2(TypeObjectInitHash)
			d.Function(func(c eval.Context, args []eval.Value) eval.Value {
				return newObjectType3(``, nil, args[0].(eval.OrderedMap)).Resolve(c)
			})
		})
}

func DefaultTypeType() *TypeType {
	return typeTypeDefault
}

func NewTypeType(containedType eval.Type) *TypeType {
	if containedType == nil || containedType == anyTypeDefault {
		return DefaultTypeType()
	}
	return &TypeType{containedType}
}

func newTypeType2(args ...eval.Value) *TypeType {
	switch len(args) {
	case 0:
		return DefaultTypeType()
	case 1:
		if containedType, ok := args[0].(eval.Type); ok {
			return NewTypeType(containedType)
		}
		panic(NewIllegalArgumentType(`Type[]`, 0, `Type`, args[0]))
	default:
		panic(errors.NewIllegalArgumentCount(`Type[]`, `0 or 1`, len(args)))
	}
}

func (t *TypeType) ContainedType() eval.Type {
	return t.typ
}

func (t *TypeType) Accept(v eval.Visitor, g eval.Guard) {
	v(t)
	t.typ.Accept(v, g)
}

func (t *TypeType) Default() eval.Type {
	return typeTypeDefault
}

func (t *TypeType) Equals(o interface{}, g eval.Guard) bool {
	if ot, ok := o.(*TypeType); ok {
		return t.typ.Equals(ot.typ, g)
	}
	return false
}

func (t *TypeType) Generic() eval.Type {
	return NewTypeType(eval.GenericType(t.typ))
}

func (t *TypeType) Get(key string) (value eval.Value, ok bool) {
	switch key {
	case `type`:
		return t.typ, true
	}
	return nil, false
}

func (t *TypeType) IsAssignable(o eval.Type, g eval.Guard) bool {
	if ot, ok := o.(*TypeType); ok {
		return GuardedIsAssignable(t.typ, ot.typ, g)
	}
	return false
}

func (t *TypeType) IsInstance(o eval.Value, g eval.Guard) bool {
	if ot, ok := o.(eval.Type); ok {
		return GuardedIsAssignable(t.typ, ot, g)
	}
	return false
}

func (t *TypeType) MetaType() eval.ObjectType {
	return TypeMetaType
}

func (t *TypeType) Name() string {
	return `Type`
}

func (t *TypeType) Parameters() []eval.Value {
	if t.typ == DefaultAnyType() {
		return eval.EmptyValues
	}
	return []eval.Value{t.typ}
}

func (t *TypeType) Resolve(c eval.Context) eval.Type {
	t.typ = resolve(c, t.typ)
	return t
}

func (t *TypeType) CanSerializeAsString() bool {
	return canSerializeAsString(t.typ)
}

func (t *TypeType) SerializationString() string {
	return t.String()
}

func (t *TypeType) String() string {
	return eval.ToString2(t, None)
}

func (t *TypeType) PType() eval.Type {
	return &TypeType{t}
}

func (t *TypeType) ToString(b io.Writer, s eval.FormatContext, g eval.RDetect) {
	TypeToString(t, b, s, g)
}