# OpenDART DS002 감사·자금·출자 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** DS002 정기보고서 주요정보의 감사·자금·출자 6개 API를 `report` 패키지에 추가한다.

**Architecture:** PR #3/#4에서 확립한 `report` 패키지의 제네릭 `getList[T]` + `ReportParams` 패턴을 그대로 재사용한다. 6개 모두 표준 요청(`corp_code`+`bsns_year`+`reprt_code`)과 list 응답이라 새 추상화가 없다. 새 파일 `report/audit.go` 에 item struct + 한 줄 메서드를 추가한다. `client.Report` 는 이미 root 에 와이어링되어 있어 root 변경은 없다.

**Tech Stack:** Go 1.25+ (제네릭), 표준 net/http (internal/httpclient 재사용), testify.

**Spec:** `docs/superpowers/specs/2026-05-24-opendart-ds002-audit-design.md`

**검증된 사실 (실 API, 삼성전자 00126380 / 2023 / 11011):** 6개 모두 status 000 + list 반환(감사 3개 각 3행, 타법인 출자 143행, 공모/사모 자금 각 1행). 숫자 콤마 문자열, 빈 값 "-". 아래 fixture 는 실 응답 첫 항목 (감사 항목의 긴 다중행 텍스트는 가독성 위해 대표 단일행으로 축약 — 실제 필드 키/구조는 동일).

**기존 재사용 심볼 (PR #3/#4, report 패키지):** `Client`, `ReportParams{CorpCode,BsnsYear,ReprtCode}`, `ReportCode`/`AnnualReport`, `getList[T](ctx, c.http, path, p)`, `report/client_test.go` 의 `newTestClient(t, routes map[string]string) *Client`.

---

## File Structure

```
report/
  audit.go         # 6개 메서드 + item struct (신규)
  audit_test.go    # 6개 fixture 테스트 (신규, newTestClient 재사용)
  testdata/        # 6개 실 응답 JSON fixture 추가
README.md          # (수정) DS002 커버리지에 감사·자금·출자
integration_test.go  # (수정) AuditOpinion 통합 케이스
```

---

### Task 1: 감사의견·감사용역·비감사용역 (3 엔드포인트)

**Files:**
- Create: `report/audit.go`, `report/audit_test.go`
- Create: `report/testdata/accnutAdtorNmNdAdtOpinion.json`, `report/testdata/adtServcCnclsSttus.json`, `report/testdata/accnutAdtorNonAdtServcCnclsSttus.json`

- [ ] **Step 1: fixture 작성**

`report/testdata/accnutAdtorNmNdAdtOpinion.json`:
```json
{
    "status": "000",
    "message": "정상",
    "list": [
        {
            "rcept_no": "20240312000736",
            "corp_cls": "Y",
            "corp_code": "00126380",
            "corp_name": "삼성전자",
            "bsns_year": "제55기(당기)",
            "adtor": "삼정회계법인",
            "adt_opinion": "적정",
            "adt_reprt_spcmnt_matter": "-",
            "emphs_matter": "해당사항 없음",
            "core_adt_matter": "메모리 반도체 재고자산 순실현가치 평가",
            "stlm_dt": "2023-12-31"
        }
    ]
}
```

`report/testdata/adtServcCnclsSttus.json`:
```json
{
    "status": "000",
    "message": "정상",
    "list": [
        {
            "rcept_no": "20240312000736",
            "corp_cls": "Y",
            "corp_code": "00126380",
            "corp_name": "삼성전자",
            "bsns_year": "제55기(당기)",
            "adtor": "삼정회계법인",
            "cn": "별도 및 연결 재무제표에 대한 감사",
            "mendng": "-",
            "tot_reqre_time": "-",
            "adt_cntrct_dtls_mendng": "7,800",
            "adt_cntrct_dtls_time": "85,700",
            "real_exc_dtls_mendng": "7,800",
            "real_exc_dtls_time": "85,036",
            "stlm_dt": "2023-12-31"
        }
    ]
}
```

`report/testdata/accnutAdtorNonAdtServcCnclsSttus.json`:
```json
{
    "status": "000",
    "message": "정상",
    "list": [
        {
            "rcept_no": "20240312000736",
            "corp_cls": "Y",
            "corp_code": "00126380",
            "corp_name": "삼성전자",
            "bsns_year": "제55기(당기)",
            "cntrct_cncls_de": "2023.05",
            "servc_cn": "ESG인증업무(국내종속기업)",
            "servc_exc_pd": "2023.05~2023.07",
            "servc_mendng": "25",
            "rm": "삼정회계법인",
            "stlm_dt": "2023-12-31"
        }
    ]
}
```

- [ ] **Step 2: 실패하는 테스트 작성** — `report/audit_test.go` (기존 `newTestClient` 재사용):
```go
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
```

- [ ] **Step 3: 테스트 실패 확인**

Run: `go test ./report/ -run 'TestAuditOpinion|TestAuditServiceContract|TestNonAuditServiceContract' -v`
Expected: FAIL — `undefined: ... AuditOpinion` 등.

- [ ] **Step 4: 구현** — `report/audit.go`:
```go
package report

import "context"

// AuditOpinionItem 은 회계감사인의 명칭 및 감사의견 (accnutAdtorNmNdAdtOpinion) 한 건.
type AuditOpinionItem struct {
	RceptNo              string `json:"rcept_no"`                // 접수번호
	CorpCls              string `json:"corp_cls"`                // 법인구분 (Y/K/N/E)
	CorpCode             string `json:"corp_code"`               // 고유번호
	CorpName             string `json:"corp_name"`               // 회사명
	BsnsYear             string `json:"bsns_year"`               // 사업연도
	Adtor                string `json:"adtor"`                   // 감사인
	AdtOpinion           string `json:"adt_opinion"`             // 감사의견
	AdtReprtSpcmntMatter string `json:"adt_reprt_spcmnt_matter"` // 감사보고서 특기사항
	EmphsMatter          string `json:"emphs_matter"`            // 강조사항 등
	CoreAdtMatter        string `json:"core_adt_matter"`         // 핵심감사사항
	StlmDt               string `json:"stlm_dt"`                 // 결산기준일
}

// AuditOpinion 은 회계감사인의 명칭 및 감사의견을 조회한다.
func (c *Client) AuditOpinion(ctx context.Context, p ReportParams) ([]AuditOpinionItem, error) {
	return getList[AuditOpinionItem](ctx, c.http, "/api/accnutAdtorNmNdAdtOpinion.json", p)
}

// AuditServiceContractItem 은 감사용역체결현황 (adtServcCnclsSttus) 한 건.
type AuditServiceContractItem struct {
	RceptNo             string `json:"rcept_no"`               // 접수번호
	CorpCls             string `json:"corp_cls"`               // 법인구분 (Y/K/N/E)
	CorpCode            string `json:"corp_code"`              // 고유번호
	CorpName            string `json:"corp_name"`              // 회사명
	BsnsYear            string `json:"bsns_year"`              // 사업연도
	Adtor               string `json:"adtor"`                  // 감사인
	Cn                  string `json:"cn"`                     // 내용
	Mendng              string `json:"mendng"`                 // 보수
	TotReqreTime        string `json:"tot_reqre_time"`         // 총소요시간
	AdtCntrctDtlsMendng string `json:"adt_cntrct_dtls_mendng"` // 감사계약내역(보수)
	AdtCntrctDtlsTime   string `json:"adt_cntrct_dtls_time"`   // 감사계약내역(시간)
	RealExcDtlsMendng   string `json:"real_exc_dtls_mendng"`   // 실제수행내역(보수)
	RealExcDtlsTime     string `json:"real_exc_dtls_time"`     // 실제수행내역(시간)
	StlmDt              string `json:"stlm_dt"`                // 결산기준일
}

// AuditServiceContract 는 감사용역체결현황을 조회한다.
func (c *Client) AuditServiceContract(ctx context.Context, p ReportParams) ([]AuditServiceContractItem, error) {
	return getList[AuditServiceContractItem](ctx, c.http, "/api/adtServcCnclsSttus.json", p)
}

// NonAuditServiceContractItem 은 회계감사인과의 비감사용역 계약체결 현황 (accnutAdtorNonAdtServcCnclsSttus) 한 건.
type NonAuditServiceContractItem struct {
	RceptNo       string `json:"rcept_no"`        // 접수번호
	CorpCls       string `json:"corp_cls"`        // 법인구분 (Y/K/N/E)
	CorpCode      string `json:"corp_code"`       // 고유번호
	CorpName      string `json:"corp_name"`       // 회사명
	BsnsYear      string `json:"bsns_year"`       // 사업연도
	CntrctCnclsDe string `json:"cntrct_cncls_de"` // 계약체결일
	ServcCn       string `json:"servc_cn"`        // 용역내용
	ServcExcPd    string `json:"servc_exc_pd"`    // 용역수행기간
	ServcMendng   string `json:"servc_mendng"`    // 용역보수
	Rm            string `json:"rm"`              // 비고
	StlmDt        string `json:"stlm_dt"`         // 결산기준일
}

// NonAuditServiceContract 는 회계감사인과의 비감사용역 계약체결 현황을 조회한다.
func (c *Client) NonAuditServiceContract(ctx context.Context, p ReportParams) ([]NonAuditServiceContractItem, error) {
	return getList[NonAuditServiceContractItem](ctx, c.http, "/api/accnutAdtorNonAdtServcCnclsSttus.json", p)
}
```

- [ ] **Step 5: 테스트 통과 확인**

Run: `go test ./report/ -v`
Expected: 전체 PASS (기존 15개 + 신규 3개, 회귀 없음). `go vet ./report/` clean.

- [ ] **Step 6: Commit**

```bash
git add report/audit.go report/audit_test.go report/testdata/
git commit -m "feat(report): 회계감사인 명칭·감사의견·감사용역·비감사용역"
```

---

### Task 2: 타법인 출자·공모자금·사모자금 (3 엔드포인트)

**Files:**
- Modify: `report/audit.go` (3개 메서드+struct 추가)
- Modify: `report/audit_test.go` (3개 테스트 추가)
- Create: `report/testdata/otrCprInvstmntSttus.json`, `report/testdata/pssrpCptalUseDtls.json`, `report/testdata/prvsrpCptalUseDtls.json`

- [ ] **Step 1: fixture 작성**

`report/testdata/otrCprInvstmntSttus.json`:
```json
{
    "status": "000",
    "message": "정상",
    "list": [
        {
            "rcept_no": "20240312000736",
            "corp_cls": "Y",
            "corp_code": "00126380",
            "corp_name": "삼성전자",
            "inv_prm": "삼성전기㈜",
            "frst_acqs_de": "1977.01.01",
            "invstmnt_purps": "경영참여",
            "frst_acqs_amount": "250,000,000",
            "bsis_blce_qy": "17,693,000",
            "bsis_blce_qota_rt": "23.7",
            "bsis_blce_acntbk_amount": "445,244,000,000",
            "incrs_dcrs_acqs_dsps_qy": "-",
            "incrs_dcrs_acqs_dsps_amount": "-",
            "incrs_dcrs_evl_lstmn": "-",
            "trmend_blce_qy": "17,693,000",
            "trmend_blce_qota_rt": "23.7",
            "trmend_blce_acntbk_amount": "445,244,000,000",
            "recent_bsns_year_fnnr_sttus_tot_assets": "11,657,872,000,000",
            "recent_bsns_year_fnnr_sttus_thstrm_ntpf": "450,482,000,000",
            "stlm_dt": "2023-12-31"
        }
    ]
}
```

`report/testdata/pssrpCptalUseDtls.json`:
```json
{
    "status": "000",
    "message": "정상",
    "list": [
        {
            "rcept_no": "20240312000736",
            "corp_cls": "Y",
            "corp_code": "00126380",
            "corp_name": "삼성전자",
            "se_nm": "-",
            "tm": "-",
            "pay_de": "-",
            "pay_amount": "-",
            "on_dclrt_cptal_use_plan": "-",
            "real_cptal_use_sttus": "-",
            "rs_cptal_use_plan_useprps": "-",
            "rs_cptal_use_plan_prcure_amount": "-",
            "real_cptal_use_dtls_cn": "-",
            "real_cptal_use_dtls_amount": "-",
            "dffrnc_occrrnc_resn": "-",
            "stlm_dt": "2023-12-31"
        }
    ]
}
```

`report/testdata/prvsrpCptalUseDtls.json`:
```json
{
    "status": "000",
    "message": "정상",
    "list": [
        {
            "rcept_no": "20240312000736",
            "corp_cls": "Y",
            "corp_code": "00126380",
            "corp_name": "삼성전자",
            "se_nm": "-",
            "tm": "-",
            "pay_de": "-",
            "pay_amount": "-",
            "cptal_use_plan": "-",
            "real_cptal_use_sttus": "-",
            "mtrpt_cptal_use_plan_useprps": "-",
            "mtrpt_cptal_use_plan_prcure_amount": "-",
            "real_cptal_use_dtls_cn": "-",
            "real_cptal_use_dtls_amount": "-",
            "dffrnc_occrrnc_resn": "-",
            "stlm_dt": "2023-12-31"
        }
    ]
}
```

- [ ] **Step 2: 실패하는 테스트 추가** — `report/audit_test.go` 에 추가:
```go
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
```

- [ ] **Step 3: 테스트 실패 확인**

Run: `go test ./report/ -run 'TestOtherCorpInvestment|TestPublicOfferingFundUsage|TestPrivatePlacementFundUsage' -v`
Expected: FAIL — `undefined: ... OtherCorpInvestment` 등.

- [ ] **Step 4: 구현** — `report/audit.go` 에 추가:
```go
// OtherCorpInvestmentItem 은 타법인 출자현황 (otrCprInvstmntSttus) 한 건.
type OtherCorpInvestmentItem struct {
	RceptNo                           string `json:"rcept_no"`                                // 접수번호
	CorpCls                           string `json:"corp_cls"`                                // 법인구분 (Y/K/N/E)
	CorpCode                          string `json:"corp_code"`                               // 고유번호
	CorpName                          string `json:"corp_name"`                               // 회사명
	InvPrm                            string `json:"inv_prm"`                                 // 법인명 (피출자 법인)
	FrstAcqsDe                        string `json:"frst_acqs_de"`                            // 최초 취득 일자
	InvstmntPurps                     string `json:"invstmnt_purps"`                          // 출자 목적
	FrstAcqsAmount                    string `json:"frst_acqs_amount"`                        // 최초 취득 금액
	BsisBlceQy                        string `json:"bsis_blce_qy"`                            // 기초 잔액 수량
	BsisBlceQotaRt                    string `json:"bsis_blce_qota_rt"`                       // 기초 잔액 지분율
	BsisBlceAcntbkAmount              string `json:"bsis_blce_acntbk_amount"`                 // 기초 잔액 장부 가액
	IncrsDcrsAcqsDspsQy               string `json:"incrs_dcrs_acqs_dsps_qy"`                 // 증가 감소 취득 처분 수량
	IncrsDcrsAcqsDspsAmount           string `json:"incrs_dcrs_acqs_dsps_amount"`             // 증가 감소 취득 처분 금액
	IncrsDcrsEvlLstmn                 string `json:"incrs_dcrs_evl_lstmn"`                    // 증가 감소 평가 손액
	TrmendBlceQy                      string `json:"trmend_blce_qy"`                          // 기말 잔액 수량
	TrmendBlceQotaRt                  string `json:"trmend_blce_qota_rt"`                     // 기말 잔액 지분율
	TrmendBlceAcntbkAmount            string `json:"trmend_blce_acntbk_amount"`               // 기말 잔액 장부 가액
	RecentBsnsYearFnnrSttusTotAssets  string `json:"recent_bsns_year_fnnr_sttus_tot_assets"`  // 최근 사업연도 재무현황 총자산
	RecentBsnsYearFnnrSttusThstrmNtpf string `json:"recent_bsns_year_fnnr_sttus_thstrm_ntpf"` // 최근 사업연도 재무현황 당기순이익
	StlmDt                            string `json:"stlm_dt"`                                 // 결산기준일
}

// OtherCorpInvestment 는 타법인 출자현황을 조회한다.
func (c *Client) OtherCorpInvestment(ctx context.Context, p ReportParams) ([]OtherCorpInvestmentItem, error) {
	return getList[OtherCorpInvestmentItem](ctx, c.http, "/api/otrCprInvstmntSttus.json", p)
}

// PublicOfferingFundUsageItem 은 공모자금의 사용내역 (pssrpCptalUseDtls) 한 건.
type PublicOfferingFundUsageItem struct {
	RceptNo                    string `json:"rcept_no"`                        // 접수번호
	CorpCls                    string `json:"corp_cls"`                        // 법인구분 (Y/K/N/E)
	CorpCode                   string `json:"corp_code"`                       // 고유번호
	CorpName                   string `json:"corp_name"`                       // 회사명
	SeNm                       string `json:"se_nm"`                           // 구분
	Tm                         string `json:"tm"`                              // 회차
	PayDe                      string `json:"pay_de"`                          // 납입일
	PayAmount                  string `json:"pay_amount"`                      // 납입금액
	OnDclrtCptalUsePlan        string `json:"on_dclrt_cptal_use_plan"`         // 신고서상 자금사용 계획
	RealCptalUseSttus          string `json:"real_cptal_use_sttus"`            // 실제 자금사용 현황
	RsCptalUsePlanUseprps      string `json:"rs_cptal_use_plan_useprps"`       // 증권신고서 등의 자금사용 계획(사용용도)
	RsCptalUsePlanPrcureAmount string `json:"rs_cptal_use_plan_prcure_amount"` // 증권신고서 등의 자금사용 계획(조달금액)
	RealCptalUseDtlsCn         string `json:"real_cptal_use_dtls_cn"`          // 실제 자금사용 내역(내용)
	RealCptalUseDtlsAmount     string `json:"real_cptal_use_dtls_amount"`      // 실제 자금사용 내역(금액)
	DffrncOccrrncResn          string `json:"dffrnc_occrrnc_resn"`             // 차이발생 사유 등
	StlmDt                     string `json:"stlm_dt"`                         // 결산기준일
}

// PublicOfferingFundUsage 는 공모자금의 사용내역을 조회한다.
func (c *Client) PublicOfferingFundUsage(ctx context.Context, p ReportParams) ([]PublicOfferingFundUsageItem, error) {
	return getList[PublicOfferingFundUsageItem](ctx, c.http, "/api/pssrpCptalUseDtls.json", p)
}

// PrivatePlacementFundUsageItem 은 사모자금의 사용내역 (prvsrpCptalUseDtls) 한 건.
type PrivatePlacementFundUsageItem struct {
	RceptNo                       string `json:"rcept_no"`                          // 접수번호
	CorpCls                       string `json:"corp_cls"`                          // 법인구분 (Y/K/N/E)
	CorpCode                      string `json:"corp_code"`                         // 고유번호
	CorpName                      string `json:"corp_name"`                         // 회사명
	SeNm                          string `json:"se_nm"`                             // 구분
	Tm                            string `json:"tm"`                                // 회차
	PayDe                         string `json:"pay_de"`                            // 납입일
	PayAmount                     string `json:"pay_amount"`                        // 납입금액
	CptalUsePlan                  string `json:"cptal_use_plan"`                    // 자금사용 계획
	RealCptalUseSttus             string `json:"real_cptal_use_sttus"`              // 실제 자금사용 현황
	MtrptCptalUsePlanUseprps      string `json:"mtrpt_cptal_use_plan_useprps"`      // 주요사항보고서의 자금사용 계획(사용용도)
	MtrptCptalUsePlanPrcureAmount string `json:"mtrpt_cptal_use_plan_prcure_amount"` // 주요사항보고서의 자금사용 계획(조달금액)
	RealCptalUseDtlsCn            string `json:"real_cptal_use_dtls_cn"`            // 실제 자금사용 내역(내용)
	RealCptalUseDtlsAmount        string `json:"real_cptal_use_dtls_amount"`        // 실제 자금사용 내역(금액)
	DffrncOccrrncResn             string `json:"dffrnc_occrrnc_resn"`               // 차이발생 사유 등
	StlmDt                        string `json:"stlm_dt"`                           // 결산기준일
}

// PrivatePlacementFundUsage 는 사모자금의 사용내역을 조회한다.
func (c *Client) PrivatePlacementFundUsage(ctx context.Context, p ReportParams) ([]PrivatePlacementFundUsageItem, error) {
	return getList[PrivatePlacementFundUsageItem](ctx, c.http, "/api/prvsrpCptalUseDtls.json", p)
}
```

- [ ] **Step 5: 테스트 통과 확인**

Run: `go test ./report/ -v`
Expected: 전체 PASS (Task 1 포함 신규 6개 + 기존). `go vet ./report/` clean.

- [ ] **Step 6: Commit**

```bash
git add report/audit.go report/audit_test.go report/testdata/
git commit -m "feat(report): 타법인 출자·공모자금·사모자금 사용내역"
```

---

### Task 3: README 커버리지 · 통합 테스트 · 최종 검증

**Files:**
- Modify: `README.md`
- Modify: `integration_test.go`

- [ ] **Step 1: README 커버리지 갱신** — `README.md` 의 DS002 줄을 다음으로 교체:
```markdown
- DS002 정기보고서 주요정보: 증자(감자) · 배당 · 자기주식 · 주식총수 · 최대주주 · 최대주주변동 · 소액주주 현황 · 증권 발행실적 · 미상환 잔액(회사채/기업어음/단기사채/신종자본증권/조건부자본증권) · 감사의견 · 감사/비감사용역 · 타법인 출자 · 공모/사모자금 사용내역
```
(바로 다음 줄 `- (예정) DS002 나머지 · DS003~DS006` 은 그대로 둔다.)

- [ ] **Step 2: 통합 테스트 추가** — `integration_test.go` 에 함수 추가 (기존 `//go:build integration` · `report` import 유지):
```go
func TestIntegration_AuditOpinion(t *testing.T) {
	c, err := NewClientFromEnv(WithCorpCodeCacheDir(t.TempDir()))
	require.NoError(t, err)

	corp, err := c.ResolveCorpCode(context.Background(), "005930")
	require.NoError(t, err)

	items, err := c.Report.AuditOpinion(context.Background(), report.ReportParams{
		CorpCode:  corp,
		BsnsYear:  "2023",
		ReprtCode: report.AnnualReport,
	})
	require.NoError(t, err)
	require.NotEmpty(t, items)
}
```

- [ ] **Step 3: 통합 빌드 확인 (기본 빌드 제외)**

Run: `go vet -tags integration ./...`
Expected: clean.
Run: `go test ./...`
Expected: 전체 PASS, integration 미실행.

- [ ] **Step 4: 최종 전체 검증**

Run:
```bash
go build ./...
go vet ./...
go test ./...
gofmt -l . | grep -v '^scripts/crawl' || echo "clean"
```
Expected: build/vet 성공, 전체 PASS, gofmt 신규 파일 차이 없음("clean").

- [ ] **Step 5: Commit**

```bash
git add README.md integration_test.go
git commit -m "docs(report): DS002 감사·자금·출자 커버리지 + 통합 테스트"
```

---

## Self-Review Notes

- **Spec coverage:** 6개 메서드+struct = Task1(3)+Task2(3) · 테스트(fixture) = Task1·2 · README 커버리지 = Task3 · 통합 테스트 = Task3. root 와이어링은 PR #3 에서 완료(변경 없음). 모두 매핑됨.
- **Type consistency:** 6개 `XItem`/메서드(`AuditOpinion`/`AuditServiceContract`/`NonAuditServiceContract`/`OtherCorpInvestment`/`PublicOfferingFundUsage`/`PrivatePlacementFundUsage`), 시그니처 `(ctx, ReportParams) ([]XItem, error)` 일관. `getList[T]`/`ReportParams`/`AnnualReport`/`newTestClient` 는 PR #3/#4 기존 심볼 재사용. 필드명·json 태그는 캡처한 실 응답과 1:1.
- **검증된 fixture:** 6개 모두 실 API(삼성전자/2023/사업보고서) 응답 첫 항목. 감사 항목의 긴 다중행 텍스트는 가독성 위해 대표 단일행으로 축약(키/구조 동일). 숫자 콤마 string, 빈 값 "-".
- **새 추상화 없음:** 기존 제네릭 getList 재사용만. root 변경 없음(client.Report 기존).
