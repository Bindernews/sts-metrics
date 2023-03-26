package orm

import (
	"context"
	"fmt"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

type CardSpec struct {
	// Base card name (without upgrades)
	Card string
	// Upgrade count
	Upgrades int
}

func (dst *CardSpec) DecodeBinary(ci *pgtype.ConnInfo, buf []byte) error {
	return (pgtype.CompositeFields{&dst.Card, &dst.Upgrades}).DecodeBinary(ci, buf)
}

func (src CardSpec) EncodeBinary(ci *pgtype.ConnInfo, buf []byte) (newBuf []byte, err error) {
	return (pgtype.CompositeFields{src.Card, int32(src.Upgrades)}).EncodeBinary(ci, buf)
}

func (dst *CardSpec) DecodeText(ci *pgtype.ConnInfo, buf []byte) error {
	return (pgtype.CompositeFields{&dst.Card, &dst.Upgrades}).DecodeText(ci, buf)
}

func (src CardSpec) EncodeText(ci *pgtype.ConnInfo, buf []byte) (newBuf []byte, err error) {
	return (pgtype.CompositeFields{src.Card, int32(src.Upgrades)}).EncodeText(ci, buf)
}

func (dst *CardSpec) Set(src any) error {
	if v, ok := src.(CardSpec); ok {
		*dst = v
		return nil
	} else {
		return fmt.Errorf("could not cast %T to %T", src, dst)
	}
}

func (src *CardSpec) Get() any {
	return *src
}

func (src *CardSpec) AssignTo(dst any) error {
	if dst_, ok := dst.(*CardSpec); ok {
		*dst_ = *src
		return nil
	} else {
		return fmt.Errorf("could not assign %T to %T", dst, src)
	}
}

func (CardSpec) RegisterType(ctx context.Context, c *pgx.Conn) error {
	name := "card_spec_io"
	var oid, arrayoid uint32
	if err := New(c).GetOID(ctx, "card_spec_io", &oid, &arrayoid); err != nil {
		return err
	}
	ci := c.ConnInfo()
	ct, err := pgtype.NewCompositeType(name, []pgtype.CompositeTypeField{
		{Name: "card", OID: pgtype.TextOID},
		{Name: "upg", OID: pgtype.Int4OID},
	}, ci)
	if err != nil {
		return err
	}
	arrayType := pgtype.NewArrayType(ct.TypeName(), oid, func() pgtype.ValueTranscoder {
		return &CardSpec{}
	})
	ci.RegisterDataType(pgtype.DataType{Value: ct, Name: ct.TypeName(), OID: oid})
	ci.RegisterDataType(pgtype.DataType{Value: arrayType, Name: arrayType.TypeName(), OID: arrayoid})
	return nil
}

func (q *Queries) CardSpecToId(ctx context.Context, specs []CardSpec) ([]int32, error) {
	rows, _ := q.db.Query(ctx, `SELECT card_spec_to_id($1)`, specs)
	defer rows.Close()
	ids := make([]int32, 0)
	for rows.Next() {
		var id int32
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (q *Queries) CardSpecAdd(ctx context.Context, specs []CardSpec) error {
	_, err := q.db.Exec(ctx, `SELECT card_spec_add($1::card_spec_io[])`, specs)
	return err
}

func (q *Queries) GetOID(ctx context.Context, typename string, oid *uint32, arrayOid *uint32) error {
	return q.db.QueryRow(ctx, `SELECT oid, typarray from pg_catalog.pg_type WHERE typname = $1`, typename).Scan(oid, arrayOid)
}
