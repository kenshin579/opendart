package material

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOtherAssetTransferPutbackOption(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/astInhtrfEtcPtbkOpt.json": "astInhtrfEtcPtbkOpt.json"})
	items, err := c.OtherAssetTransferPutbackOption(context.Background(), MaterialParams{CorpCode: "00126380", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	got := items[0]
	assert.Equal(t, "20230315000111", got.RceptNo)
	assert.Equal(t, "자산양수도(기타)에 해당", got.RpRsn)
	assert.Equal(t, "120,000,000,000", got.AstInhtrfPrc)
}

func TestStockExchangeTransfer(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/stkExtrDecsn.json": "stkExtrDecsn.json"})
	items, err := c.StockExchangeTransfer(context.Background(), MaterialParams{CorpCode: "00126380", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	got := items[0]
	assert.Equal(t, "20230420000222", got.RceptNo)
	assert.Equal(t, "주식의 포괄적 교환", got.ExtrSen)
	assert.Equal(t, "대상법인", got.ExtrTgcmpCmpnm)
	assert.Equal(t, "1:0.5", got.ExtrRt)
	assert.Equal(t, "2023년 06월 10일", got.ExtrscGmtsckPrd)
	assert.Equal(t, "제출", got.RsSmAtn)
}
