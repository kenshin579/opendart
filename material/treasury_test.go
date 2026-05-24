package material

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTreasuryStockAcquisition(t *testing.T) {
	c := newTestClient(t, map[string]string{
		"/api/tsstkAqDecsn.json": "tsstkAqDecsn.json",
	})

	items, err := c.TreasuryStockAcquisition(context.Background(), MaterialParams{CorpCode: "00126380", BgnDe: "20240101", EndDe: "20241231"})
	require.NoError(t, err)
	require.Len(t, items, 1)

	got := items[0]
	assert.Equal(t, "20241118000328", got.RceptNo)
	assert.Equal(t, "50,144,628", got.AqplnStkOstk)
	assert.Equal(t, "유가증권시장을 통한 장내 매수", got.AqMth)
	assert.Equal(t, "-", got.AqWtnDivOstkRt)
	assert.Equal(t, "6,525,123", got.D1ProdlmOstk)
}
