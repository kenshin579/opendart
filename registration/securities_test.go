package registration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEquitySecurities(t *testing.T) {
	c := newTestClient(t, "estkRs.json")
	res, err := c.EquitySecurities(context.Background(), Params{CorpCode: "00107598", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Len(t, res.General, 1)
	require.Len(t, res.SecurityTypes, 1)
	require.Len(t, res.Underwriters, 1)
	require.Len(t, res.FundUsage, 2)
	require.Len(t, res.Sellers, 1)
	require.Len(t, res.RetailPutbackOption, 1)
	assert.Equal(t, "20230515002454", res.General[0].RceptNo)
	assert.Equal(t, "남양유업", res.General[0].CorpName)
	assert.Equal(t, "우선주", res.SecurityTypes[0].Stksen)
	assert.Equal(t, "NH투자증권", res.Underwriters[0].Actnmn)
	assert.Equal(t, "운영자금", res.FundUsage[0].Se)
	assert.Equal(t, "발행제비용", res.FundUsage[1].Se)
}
