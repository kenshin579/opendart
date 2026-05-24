package report

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuditOpinion(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/accnutAdtorNmNdAdtOpinion.json": "accnutAdtorNmNdAdtOpinion.json"})
	items, err := c.AuditOpinion(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "삼정회계법인", items[0].Adtor)
	assert.Equal(t, "적정", items[0].AdtOpinion)
}

func TestAuditServiceContract(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/adtServcCnclsSttus.json": "adtServcCnclsSttus.json"})
	items, err := c.AuditServiceContract(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "삼정회계법인", items[0].Adtor)
	assert.Equal(t, "7,800", items[0].AdtCntrctDtlsMendng)
}

func TestNonAuditServiceContract(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/accnutAdtorNonAdtServcCnclsSttus.json": "accnutAdtorNonAdtServcCnclsSttus.json"})
	items, err := c.NonAuditServiceContract(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "ESG인증업무(국내종속기업)", items[0].ServcCn)
	assert.Equal(t, "삼정회계법인", items[0].Rm)
}

func TestOtherCorpInvestment(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/otrCprInvstmntSttus.json": "otrCprInvstmntSttus.json"})
	items, err := c.OtherCorpInvestment(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "삼성전기㈜", items[0].InvPrm)
	assert.Equal(t, "23.7", items[0].TrmendBlceQotaRt)
	assert.Equal(t, "445,244,000,000", items[0].TrmendBlceAcntbkAmount)
}

func TestPublicOfferingFundUsage(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/pssrpCptalUseDtls.json": "pssrpCptalUseDtls.json"})
	items, err := c.PublicOfferingFundUsage(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "2023-12-31", items[0].StlmDt)
	assert.Equal(t, "-", items[0].RsCptalUsePlanUseprps)
}

func TestPrivatePlacementFundUsage(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/prvsrpCptalUseDtls.json": "prvsrpCptalUseDtls.json"})
	items, err := c.PrivatePlacementFundUsage(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "2023-12-31", items[0].StlmDt)
	assert.Equal(t, "-", items[0].CptalUsePlan)
}
