package registration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMerger(t *testing.T) {
	c := newTestClient(t, "mgRs.json")
	res, err := c.Merger(context.Background(), Params{CorpCode: "00126380", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Len(t, res.General, 1)
	require.Len(t, res.IssuedSecurities, 1)
	require.Len(t, res.PartyCompanies, 1)
	assert.Equal(t, "20230410000111", res.General[0].RceptNo)
	assert.Equal(t, "흡수합병", res.General[0].Stn)
	assert.Equal(t, "1:0.5", res.General[0].RtVl)
	assert.Equal(t, "기명식 보통주", res.IssuedSecurities[0].Kndn)
	assert.Equal(t, "합병상대회사", res.PartyCompanies[0].Cmpnm)
}

func TestDivision(t *testing.T) {
	c := newTestClient(t, "dvRs.json")
	res, err := c.Division(context.Background(), Params{CorpCode: "00126380", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Len(t, res.General, 1)
	require.Len(t, res.IssuedSecurities, 1)
	require.Len(t, res.PartyCompanies, 1)
	assert.Equal(t, "20230510000222", res.General[0].RceptNo)
	assert.Equal(t, "인적분할", res.General[0].Stn)
	assert.Equal(t, "0.7:0.3", res.General[0].RtVl)
	assert.Equal(t, "분할신설회사", res.PartyCompanies[0].Cmpnm)
}

func TestStockExchangeTransfer(t *testing.T) {
	c := newTestClient(t, "extrRs.json")
	res, err := c.StockExchangeTransfer(context.Background(), Params{CorpCode: "00126380", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Len(t, res.General, 1)
	require.Len(t, res.IssuedSecurities, 1)
	require.Len(t, res.PartyCompanies, 1)
	assert.Equal(t, "20230610000333", res.General[0].RceptNo)
	assert.Equal(t, "주식의 포괄적 교환", res.General[0].Stn)
	assert.Equal(t, "1:0.8", res.General[0].RtVl)
	assert.Equal(t, "완전자회사", res.PartyCompanies[0].Cmpnm)
}
