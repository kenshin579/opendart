package report

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSingleAccount(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/fnlttSinglAcnt.json": "fnlttSinglAcnt.json"})
	items, err := c.SingleAccount(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "유동자산", items[0].AccountNm)
	assert.Equal(t, "연결재무제표", items[0].FsNm)
	assert.Equal(t, "195,936,557,000,000", items[0].ThstrmAmount)
	assert.Equal(t, "00126380", items[0].CorpCode)
}

func TestMultiAccount(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/fnlttMultiAcnt.json": "fnlttMultiAcnt.json"})
	items, err := c.MultiAccount(context.Background(), ReportParams{CorpCode: "00126380,00164779", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "유동자산", items[0].AccountNm)
	assert.Equal(t, "00126380", items[0].CorpCode)
}
