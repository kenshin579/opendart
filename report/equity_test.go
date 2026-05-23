package report

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDividend(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/alotMatter.json": "alotMatter.json"})
	items, err := c.Dividend(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "주당액면가액(원)", items[0].Se)
	assert.Equal(t, "100", items[0].Thstrm)
	assert.Equal(t, "2023-12-31", items[0].StlmDt)
}
