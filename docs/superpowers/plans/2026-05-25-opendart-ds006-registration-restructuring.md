# OpenDART DS006 증권신고서 Sub-2 «신고» Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** DS006 증권신고서 Sub-2 «신고» 3개 API(합병/분할/주식의포괄적교환·이전)를 기존 `registration` 패키지에 추가한다. (완료 시 DS006 6/6.)

**Architecture:** 기존 `registration.Client` + `Params` + `httpclient.GetGroups`(그룹형 디코더) 재사용. 신규 파일 `registration/restructuring.go`. 세 엔드포인트 그룹 스키마 동일 → 공유 item 타입 3종(`RestructuringGeneralItem` 17, `RestructuringIssuedSecurityItem` 9, `RestructuringPartyCompanyItem` 10), wrapper 3종(엔드포인트별), 메서드 3개(GetGroups→title switch→json.Unmarshal). root `opendart` 패키지 변경 없음(`client.Registration` 기존).

**Tech Stack:** Go, 표준 net/http(`internal/httpclient`), `encoding/json`, testify, httptest. fixture 는 실 API 캡처 권장(불가 시 docs 스키마 일치 샘플, 그룹 배열 형태 유지).

**SINGLE SOURCE OF TRUTH for struct definitions:** `docs/superpowers/specs/2026-05-25-opendart-ds006-registration-restructuring-design.md` 에 3개 item struct + 3개 wrapper + Merger 메서드가 EXACT Go 코드로 정의돼 있다. Division/StockExchangeTransfer 메서드는 spec 에서 Merger 와 동일 패턴으로 축약(`...`)돼 있으나 이 plan 에 완전한 코드를 명시한다. 각 Task 의 item/wrapper struct 는 spec 동일 이름을 그대로 복사한다(필드 가감 금지).

---

## File Structure

- Create: `registration/restructuring.go` — 3 공유 item + 3 wrapper + 3 메서드.
- Create: `registration/restructuring_test.go` — fixture 디코딩 테스트(기존 `newTestClient` 재사용).
- Create: `registration/testdata/{mgRs,dvRs,extrRs}.json` — 3 fixture.
- Modify: `integration_test.go` — 통합 케이스 2개(`//go:build integration`).
- Modify: `README.md` — DS006 커버리지에 합병·분할·교환이전 추가, 예정 줄에서 DS006 제거.

기존 컨벤션(변경 금지):
- `registration/securities.go` 메서드 패턴(Sub-1): `func (c *Client) EquitySecurities(ctx, p Params) (*EquitySecuritiesRegistration, error)` 가 `httpclient.GetGroups(ctx, c.http, "/api/estkRs.json", p.toMap())` 호출 후 `switch g.Title` 로 `json.Unmarshal(g.List, &out.<field>)`.
- `registration/client.go`: `Client{http}`, `Params{CorpCode,BgnDe,EndDe}`+`toMap()`.
- `registration/client_test.go`: `newTestClient(t, fixture string) *Client` (지정 fixture 서빙).
- `registration/securities.go` import: `context`, `encoding/json`, `github.com/kenshin579/opendart/internal/httpclient`.
- 모든 응답 필드 string + 한글 코멘트. UTF-8. testify.

`go test` 의 작업 디렉터리는 항상 `cd /Users/user/src/workspace_moneyflow/opendart`.

---

## Task 1: 공유 item 3종 + 합병 (Merger) — mgRs

**Files:** Create `registration/restructuring.go`, `registration/restructuring_test.go`, `registration/testdata/mgRs.json`.

이 Task 가 공유 item 3종을 정의하고 첫 메서드(Merger)까지 동작시킨다.

- [ ] **Step 1: fixture `registration/testdata/mgRs.json`** (그룹형, 3그룹, 각 list 1건; 일반사항 17 키 전부 포함)
```json
{
  "status": "000",
  "message": "정상",
  "group": [
    {"title": "일반사항", "list": [
      {"rcept_no": "20230410000111", "corp_cls": "Y", "corp_code": "00126380", "corp_name": "테스트합병", "stn": "흡수합병", "bddd": "2023년 04월 10일", "ctrd": "2023년 04월 10일", "gmtsck_shddstd": "2023년 05월 10일", "ap_gmtsck": "2023년 06월 10일", "aprskh_pd_bgd": "2023년 06월 10일", "aprskh_pd_edd": "2023년 06월 30일", "aprskh_prc": "65,000", "mgdt_etc": "2023년 07월 15일", "rt_vl": "1:0.5", "exevl_int": "삼일회계법인", "grtmn_etc": "-", "rpt_rcpn": "20230410003707"}
    ]},
    {"title": "발행증권", "list": [
      {"rcept_no": "20230410000111", "corp_cls": "Y", "corp_code": "00126380", "corp_name": "테스트합병", "kndn": "기명식 보통주", "cnt": "1,000,000", "fv": "5,000", "slprc": "65,000", "slta": "65,000,000,000"}
    ]},
    {"title": "당사회사에관한사항", "list": [
      {"rcept_no": "20230410000111", "corp_cls": "Y", "corp_code": "00126380", "corp_name": "테스트합병", "cmpnm": "합병상대회사", "sen": "소멸회사", "tast": "500,000,000,000", "cpt": "50,000,000,000", "isstk_knd": "기명식 보통주", "isstk_cnt": "2,000,000"}
    ]}
  ]
}
```

- [ ] **Step 2: create `registration/restructuring_test.go`** (`package registration` + import `context`,`testing`,`github.com/stretchr/testify/assert`,`github.com/stretchr/testify/require`)
```go
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
```

- [ ] **Step 3: run `go test ./registration/ -run TestMerger`** — confirm FAIL (Merger/types undefined).

- [ ] **Step 4: create `registration/restructuring.go`** — header:
```go
package registration

import (
	"context"
	"encoding/json"

	"github.com/kenshin579/opendart/internal/httpclient"
)
```
Then copy VERBATIM from spec: `RestructuringGeneralItem` (17), `RestructuringIssuedSecurityItem` (9), `RestructuringPartyCompanyItem` (10), `MergerRegistration` wrapper, `Merger` method (full code in spec). Do NOT add Division/StockExchangeTransfer yet.

- [ ] **Step 5: run test** — PASS.
- [ ] **Step 6: `gofmt -l registration/restructuring.go registration/restructuring_test.go`** — no output (else `gofmt -w`).
- [ ] **Step 7: verify field counts** — `RestructuringGeneralItem` 17, `RestructuringIssuedSecurityItem` 9, `RestructuringPartyCompanyItem` 10 json tags.
- [ ] **Step 8: real API BEST EFFORT** — `$OPENDART_API_KEY` set → `curl -s "https://opendart.fss.or.kr/api/mgRs.json?crtfc_key=$OPENDART_API_KEY&corp_code=00126380&bgn_de=20180101&end_de=20241231"`. 서버 TLS1.2 RSA. clean `"status":"000"` + group 있으면 fixture 교체 + assert 갱신 후 step 5 재실행. 013 이면 합병 사례 종목 1~2개 더 시도(예: 00164779). 모두 013 이면 샘플 유지. 최대 2회.
- [ ] **Step 9: commit (no push)**
```bash
cd /Users/user/src/workspace_moneyflow/opendart
git add registration/restructuring.go registration/restructuring_test.go registration/testdata/mgRs.json
git commit -m "feat(registration): add DS006 Merger (합병) + Restructuring 공유 타입"
```

---

## Task 2: 분할 (Division) — dvRs

**Files:** Modify `registration/restructuring.go`, `registration/restructuring_test.go`; Create `registration/testdata/dvRs.json`.

공유 item 타입(Task 1) 재사용 — 신규 item 없음. wrapper + 메서드만 추가.

- [ ] **Step 1: fixture `registration/testdata/dvRs.json`** (3그룹, 스키마 = Restructuring* 동일)
```json
{
  "status": "000",
  "message": "정상",
  "group": [
    {"title": "일반사항", "list": [
      {"rcept_no": "20230510000222", "corp_cls": "Y", "corp_code": "00126380", "corp_name": "테스트분할", "stn": "인적분할", "bddd": "2023년 05월 10일", "ctrd": "2023년 05월 10일", "gmtsck_shddstd": "2023년 06월 10일", "ap_gmtsck": "2023년 07월 10일", "aprskh_pd_bgd": "2023년 07월 10일", "aprskh_pd_edd": "2023년 07월 30일", "aprskh_prc": "70,000", "mgdt_etc": "2023년 08월 01일", "rt_vl": "0.7:0.3", "exevl_int": "안진회계법인", "grtmn_etc": "-", "rpt_rcpn": "20230510003707"}
    ]},
    {"title": "발행증권", "list": [
      {"rcept_no": "20230510000222", "corp_cls": "Y", "corp_code": "00126380", "corp_name": "테스트분할", "kndn": "기명식 보통주", "cnt": "500,000", "fv": "5,000", "slprc": "70,000", "slta": "35,000,000,000"}
    ]},
    {"title": "당사회사에관한사항", "list": [
      {"rcept_no": "20230510000222", "corp_cls": "Y", "corp_code": "00126380", "corp_name": "테스트분할", "cmpnm": "분할신설회사", "sen": "신설회사", "tast": "300,000,000,000", "cpt": "30,000,000,000", "isstk_knd": "기명식 보통주", "isstk_cnt": "500,000"}
    ]}
  ]
}
```

- [ ] **Step 2: append test to `registration/restructuring_test.go`**
```go
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
```

- [ ] **Step 3: run `go test ./registration/ -run TestDivision`** — FAIL (undefined).
- [ ] **Step 4: append to `registration/restructuring.go`**:
```go
// DivisionRegistration 은 분할 증권신고서(dvRs)의 그룹별 항목.
type DivisionRegistration struct {
	General          []RestructuringGeneralItem        // 일반사항
	IssuedSecurities []RestructuringIssuedSecurityItem // 발행증권
	PartyCompanies   []RestructuringPartyCompanyItem   // 당사회사에관한사항
}

// Division 은 분할 증권신고서(DS006)를 조회한다.
func (c *Client) Division(ctx context.Context, p Params) (*DivisionRegistration, error) {
	groups, err := httpclient.GetGroups(ctx, c.http, "/api/dvRs.json", p.toMap())
	if err != nil {
		return nil, err
	}
	out := &DivisionRegistration{}
	for _, g := range groups {
		var derr error
		switch g.Title {
		case "일반사항":
			derr = json.Unmarshal(g.List, &out.General)
		case "발행증권":
			derr = json.Unmarshal(g.List, &out.IssuedSecurities)
		case "당사회사에관한사항":
			derr = json.Unmarshal(g.List, &out.PartyCompanies)
		}
		if derr != nil {
			return nil, derr
		}
	}
	return out, nil
}
```
- [ ] **Step 5: run test** — PASS.
- [ ] **Step 6: gofmt** — `gofmt -l registration/` 출력 없으면 OK.
- [ ] **Step 7: real API BEST EFFORT** — `curl -s "https://opendart.fss.or.kr/api/dvRs.json?crtfc_key=$OPENDART_API_KEY&corp_code=00126380&bgn_de=20180101&end_de=20241231"`. clean 000+group → 교체+갱신, 안되면 샘플 유지. 1회.
- [ ] **Step 8: commit (no push)**
```bash
cd /Users/user/src/workspace_moneyflow/opendart
git add registration/restructuring.go registration/restructuring_test.go registration/testdata/dvRs.json
git commit -m "feat(registration): add DS006 Division (분할)"
```

---

## Task 3: 주식의포괄적교환·이전 (StockExchangeTransfer) — extrRs

**Files:** Modify `registration/restructuring.go`, `registration/restructuring_test.go`; Create `registration/testdata/extrRs.json`.

공유 item 타입 재사용 — 신규 item 없음.

- [ ] **Step 1: fixture `registration/testdata/extrRs.json`** (3그룹)
```json
{
  "status": "000",
  "message": "정상",
  "group": [
    {"title": "일반사항", "list": [
      {"rcept_no": "20230610000333", "corp_cls": "Y", "corp_code": "00126380", "corp_name": "테스트교환이전", "stn": "주식의 포괄적 교환", "bddd": "2023년 06월 10일", "ctrd": "2023년 06월 10일", "gmtsck_shddstd": "2023년 07월 10일", "ap_gmtsck": "2023년 08월 10일", "aprskh_pd_bgd": "2023년 08월 10일", "aprskh_pd_edd": "2023년 08월 30일", "aprskh_prc": "68,000", "mgdt_etc": "2023년 09월 01일", "rt_vl": "1:0.8", "exevl_int": "삼정회계법인", "grtmn_etc": "-", "rpt_rcpn": "20230610003707"}
    ]},
    {"title": "발행증권", "list": [
      {"rcept_no": "20230610000333", "corp_cls": "Y", "corp_code": "00126380", "corp_name": "테스트교환이전", "kndn": "기명식 보통주", "cnt": "800,000", "fv": "5,000", "slprc": "68,000", "slta": "54,400,000,000"}
    ]},
    {"title": "당사회사에관한사항", "list": [
      {"rcept_no": "20230610000333", "corp_cls": "Y", "corp_code": "00126380", "corp_name": "테스트교환이전", "cmpnm": "완전자회사", "sen": "완전자회사", "tast": "400,000,000,000", "cpt": "40,000,000,000", "isstk_knd": "기명식 보통주", "isstk_cnt": "1,000,000"}
    ]}
  ]
}
```

- [ ] **Step 2: append test to `registration/restructuring_test.go`**
```go
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
```

- [ ] **Step 3: run `go test ./registration/ -run TestStockExchangeTransfer`** — FAIL (undefined).
- [ ] **Step 4: append to `registration/restructuring.go`**:
```go
// StockExchangeTransferRegistration 은 주식의포괄적교환·이전 증권신고서(extrRs)의 그룹별 항목.
type StockExchangeTransferRegistration struct {
	General          []RestructuringGeneralItem        // 일반사항
	IssuedSecurities []RestructuringIssuedSecurityItem // 발행증권
	PartyCompanies   []RestructuringPartyCompanyItem   // 당사회사에관한사항
}

// StockExchangeTransfer 는 주식의포괄적교환·이전 증권신고서(DS006)를 조회한다.
func (c *Client) StockExchangeTransfer(ctx context.Context, p Params) (*StockExchangeTransferRegistration, error) {
	groups, err := httpclient.GetGroups(ctx, c.http, "/api/extrRs.json", p.toMap())
	if err != nil {
		return nil, err
	}
	out := &StockExchangeTransferRegistration{}
	for _, g := range groups {
		var derr error
		switch g.Title {
		case "일반사항":
			derr = json.Unmarshal(g.List, &out.General)
		case "발행증권":
			derr = json.Unmarshal(g.List, &out.IssuedSecurities)
		case "당사회사에관한사항":
			derr = json.Unmarshal(g.List, &out.PartyCompanies)
		}
		if derr != nil {
			return nil, derr
		}
	}
	return out, nil
}
```
- [ ] **Step 5: run test** — PASS.
- [ ] **Step 6: gofmt** — `gofmt -l registration/` 출력 없으면 OK.
- [ ] **Step 7: real API BEST EFFORT** — `curl -s "https://opendart.fss.or.kr/api/extrRs.json?crtfc_key=$OPENDART_API_KEY&corp_code=00126380&bgn_de=20180101&end_de=20241231"`. clean 000+group → 교체+갱신, 안되면 샘플 유지. 1회.
- [ ] **Step 8: commit (no push)**
```bash
cd /Users/user/src/workspace_moneyflow/opendart
git add registration/restructuring.go registration/restructuring_test.go registration/testdata/extrRs.json
git commit -m "feat(registration): add DS006 StockExchangeTransfer (주식의포괄적교환·이전)"
```

---

## Task 4: 통합 테스트 + README

**Files:** Modify `integration_test.go`, `README.md`.

- [ ] **Step 1: 통합 테스트 추가**

`integration_test.go` 를 먼저 읽어 패턴 확인(`//go:build integration`, `package opendart`, `NewClientFromEnv(WithCorpCodeCacheDir(t.TempDir()))`, `ErrNoData` 직접 참조, `registration` import 이미 존재 — Sub-1 에서 추가됨). 파일 끝에 2개 추가:
```go
func TestIntegration_Merger(t *testing.T) {
	c, err := NewClientFromEnv(WithCorpCodeCacheDir(t.TempDir()))
	require.NoError(t, err)
	corp, err := c.ResolveCorpCode(context.Background(), "005930")
	require.NoError(t, err)
	res, err := c.Registration.Merger(context.Background(), registration.Params{CorpCode: corp, BgnDe: "20180101", EndDe: "20241231"})
	if errors.Is(err, ErrNoData) {
		t.Skip("해당 기간 합병 증권신고서 데이터 없음")
	}
	require.NoError(t, err)
	require.NotNil(t, res)
	for _, it := range res.General {
		require.NotEmpty(t, it.RceptNo)
	}
}

func TestIntegration_Division(t *testing.T) {
	c, err := NewClientFromEnv(WithCorpCodeCacheDir(t.TempDir()))
	require.NoError(t, err)
	corp, err := c.ResolveCorpCode(context.Background(), "005930")
	require.NoError(t, err)
	res, err := c.Registration.Division(context.Background(), registration.Params{CorpCode: corp, BgnDe: "20180101", EndDe: "20241231"})
	if errors.Is(err, ErrNoData) {
		t.Skip("해당 기간 분할 증권신고서 데이터 없음")
	}
	require.NoError(t, err)
	require.NotNil(t, res)
	for _, it := range res.General {
		require.NotEmpty(t, it.RceptNo)
	}
}
```

- [ ] **Step 2: 통합 빌드 확인** — `cd /Users/user/src/workspace_moneyflow/opendart && go vet -tags integration ./...` → 출력 없음.
- [ ] **Step 3: 통합 테스트 실행(키 있으면)** — `go test -tags integration -run "TestIntegration_Merger|TestIntegration_Division" ./...` → PASS 또는 SKIP.

- [ ] **Step 4: README 커버리지 갱신**

`README.md` 의 DS006 줄(현재 `- DS006 증권신고서 주요정보: 지분증권 · 채무증권 · 증권예탁증권`)에 ` · 합병 · 분할 · 주식의포괄적교환·이전` 을 추가. 그리고 `(예정)` 줄(현재 `(예정) DS006 나머지(합병/분할/주식의포괄적교환·이전) · DS002 개인별 보수 Ver2.0`)을 `(예정) DS002 개인별 보수 Ver2.0` 로 변경(DS006 6/6 완료). 실제 파일 문구 확인해 동등 편집.

- [ ] **Step 5: 전체 게이트** — `cd /Users/user/src/workspace_moneyflow/opendart && go build ./... && go test ./... && gofmt -l registration/ integration_test.go` → 빌드 OK, 전체 PASS, gofmt 출력 없음.
- [ ] **Step 6: README UTF-8** — `file -I README.md` → `charset=utf-8`.
- [ ] **Step 7: 커밋** — `git add integration_test.go README.md && git commit -m "test(registration): add DS006 신고 통합 테스트 + README 커버리지 (DS006 6/6 완료)"`.

---

## Self-Review (작성자 점검 결과)

**1. Spec coverage:** spec 의 3개 메서드(Merger/Division/StockExchangeTransfer) → Task 1/2/3. 공유 item 3종은 Task 1. spec 에서 축약(`...`)된 Division/StockExchangeTransfer 메서드를 Task 2/3 에 완전한 코드로 명시. 통합·README = Task 4. 누락 없음.

**2. Placeholder scan:** TBD/TODO 없음. 인프라 재사용, 메서드·fixture·test 완전. item/wrapper struct body 는 committed spec(EXACT Go 코드) 참조. 필드 수 명시(17/9/10).

**3. Type consistency:** 메서드/wrapper/item 이름 spec 표와 1:1. 세 wrapper 가 동일 3 공유 item 타입 참조. fixture group title(일반사항/발행증권/당사회사에관한사항)이 메서드 switch case 와 정확히 일치. 통합 테스트는 같은 패키지라 `ErrNoData` 직접 참조, registration import 기존.
