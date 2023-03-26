package orm

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

const testEnvConn = "POSTGRES_CONN"

func TestCardSpecRegister(t *testing.T) {
	ctx := context.Background()
	if err := godotenv.Load("../.env"); err != nil {
		t.Fatal(err)
	}
	c, err := pgx.Connect(ctx, os.Getenv(testEnvConn))
	if err != nil {
		t.Fatal(err)
	}
	if err := (CardSpec{}).RegisterType(ctx, c); err != nil {
		t.Fatal(err)
	}

	{
		var card string
		var upg int
		err = c.QueryRow(ctx, `select ('ab',1)::card_spec_io`).Scan([]any{&card, &upg})
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, card, "ab")
		assert.Equal(t, upg, 1)
	}

	{
		act := []CardSpec{}
		exp := []CardSpec{
			{Card: "ab", Upgrades: 1},
			{Card: "cd", Upgrades: 2},
		}
		err = c.QueryRow(ctx, `select array[('ab',1),('cd',2)]::card_spec_io[]`).Scan(&act)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, exp, act)
	}

	{
		inp := CardSpec{Card: "ab", Upgrades: 1}
		out := CardSpec{}
		err = c.QueryRow(ctx, `select $1::card_spec_io`, inp).Scan(&out)
		assert.NoError(t, err)
		assert.Equal(t, inp, out)
	}

	{
		inp := []CardSpec{
			{Card: "ab", Upgrades: 1},
			{Card: "cd", Upgrades: 2},
		}
		out := CardSpec{}
		err = c.QueryRow(ctx, `select * from unnest($1::card_spec_io[]) limit 1`, inp).Scan(&out.Card, &out.Upgrades)
		assert.NoError(t, err)
		assert.Equal(t, out, inp[0])
	}
}
