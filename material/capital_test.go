package material

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPaidInCapitalIncrease(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/piicDecsn.json": "piicDecsn.json"})
	items, err := c.PaidInCapitalIncrease(context.Background(), MaterialParams{CorpCode: "00107598", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "남양유업", items[0].CorpName)
	assert.Equal(t, "주주우선공모증자", items[0].IcMthn)
	assert.Equal(t, "7,184,339,000", items[0].FdppOp)
}

func TestFreeCapitalIncrease(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/fricDecsn.json": "fricDecsn.json"})
	items, err := c.FreeCapitalIncrease(context.Background(), MaterialParams{CorpCode: "00117230", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "14,580,207", items[0].NstkOstkCnt)
	assert.Equal(t, "0.5", items[0].NstkAscntPsOstk)
	assert.Equal(t, "2,500", items[0].FvPs)
}
