package remokon

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/stk132/pmstr"
	"reflect"
)

var (
	ErrNotStructType = errors.New("must be struct type")
)

type Client struct {
	p *pmstr.Pmstr
}

type remokonValue struct {
	name string
	parameterPath string
	typ reflect.Type
}

func New(sess *session.Session) *Client {
	p := pmstr.New(ssm.New(sess))
	return &Client{p}
}

func (c *Client) Load(i interface{}) error {
	e := reflect.ValueOf(i).Elem()
	typ := e.Type()
	if typ.Kind() != reflect.Struct {
		return ErrNotStructType
	}

	fields := []*remokonValue{}
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		p, ok := f.Tag.Lookup("parameter_store")
		if ok {
			v := &remokonValue{
				typ: f.Type,
				name: f.Name,
				parameterPath: p,
			}
			fields = append(fields, v)
		}
	}

	for _, f := range fields {
		v, err := c.p.Get(f.parameterPath).AsString()
		if err != nil {
			return err
		}

		e.FieldByName(f.name).SetString(v)
	}

	return nil
}