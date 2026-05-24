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
