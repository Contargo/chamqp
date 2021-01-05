package exchange_declare

import (
	"github.com/Contargo/chamqp"
	"github.com/streadway/amqp"
)

var Defaults = chamqp.ExchangeDeclareSpec{
	Name:       "",
	Kind:       "topic",
	Durable:    true,
	AutoDelete: false,
	Internal:   false,
	NoWait:     false,
	Args:       nil,
	ErrorChan:  nil,
}

type NameDecl struct {
	exchangeDeclarationSpec chamqp.ExchangeDeclareSpec
}

func (n NameDecl) WithKind(kind string) KindDecl {
	n.exchangeDeclarationSpec.Kind = kind
	return KindDecl{n}
}

func (n NameDecl) WithDefaultKind() KindDecl {
	return KindDecl{n}
}

func (n NameDecl) Defaults() End {
	return End{n}
}

type KindDecl struct {
	nameDecl NameDecl
}

func (k KindDecl) WithDurable(durable bool) AutoDeleteDecl {
	k.nameDecl.exchangeDeclarationSpec.Durable = durable
	return AutoDeleteDecl{k.nameDecl}
}

func (k KindDecl) WithDefaultDurable() AutoDeleteDecl {
	return AutoDeleteDecl{k.nameDecl}
}

func (k KindDecl) Defaults() End {
	return End{k.nameDecl}
}

type AutoDeleteDecl struct {
	nameDecl NameDecl
}

func (a AutoDeleteDecl) WithAutoDelete(autodelete bool) InternalDecl {
	a.nameDecl.exchangeDeclarationSpec.AutoDelete = autodelete
	return InternalDecl{a.nameDecl}
}

func (a AutoDeleteDecl) WithDefaultAutoDelete() InternalDecl {
	return InternalDecl{a.nameDecl}
}

func (a AutoDeleteDecl) Defaults() End {
	return End{a.nameDecl}	
}

type InternalDecl struct {
	nameDecl NameDecl
}

func (i InternalDecl) WithInternal(internal bool) NoWaitDecl {
	i.nameDecl.exchangeDeclarationSpec.Internal = internal
	return NoWaitDecl{i.nameDecl}
}

func (i InternalDecl) WithDefaultInternal() NoWaitDecl {
	return NoWaitDecl{i.nameDecl}
}

func (i InternalDecl) Defaults() End {
	return End{i.nameDecl}
}

type NoWaitDecl struct {
	nameDecl NameDecl
}

func (n NoWaitDecl) WithNoWait(nowait bool) ArgsDecl {
	n.nameDecl.exchangeDeclarationSpec.NoWait = nowait
	return ArgsDecl{n.nameDecl}
}

func (n NoWaitDecl) WithDefaultNoWait() ArgsDecl {
	return ArgsDecl{n.nameDecl}
}

func (n NoWaitDecl) Defaults() End {
	return End{n.nameDecl}
}

type ArgsDecl struct {
	nameDecl NameDecl
}

func (a ArgsDecl) WithArgs(args amqp.Table) ErrorChanDecl {
	a.nameDecl.exchangeDeclarationSpec.Args = args
	return ErrorChanDecl{a.nameDecl}
}

func (a ArgsDecl) Defaults() End {
	return End{a.nameDecl}
}

type ErrorChanDecl struct {
	nameDecl NameDecl
}

func (e ErrorChanDecl) WithErrorChan(errorChan chan error) End {
	e.nameDecl.exchangeDeclarationSpec.ErrorChan = errorChan
	return End{e.nameDecl}
}

func (e ErrorChanDecl) Defaults() End {
	return End{e.nameDecl}
}

type End struct {
	nameDecl NameDecl
}

func (e End) BuildSpec() chamqp.ExchangeDeclareSpec {
	return e.nameDecl.exchangeDeclarationSpec
}

func (e End) Build(ch *chamqp.Channel) {
	ch.ExchangeDeclareWithSpec(e.BuildSpec())
}

func DeclareExchange(exchangeName string) NameDecl {
	filledDefaults := Defaults
	filledDefaults.Name = exchangeName
	return NameDecl{filledDefaults}
}
