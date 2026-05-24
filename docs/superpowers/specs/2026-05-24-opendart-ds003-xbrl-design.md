# OpenDART DS003 정기보고서 재무정보 — XBRL 그룹 설계

- 작성일: 2026-05-24
- 모듈: `github.com/kenshin579/opendart`
- 범위: **DS003 XBRL 2개 API** (`report` 패키지 확장) — 재무제표 원본파일(ZIP) + 택사노미 양식

## 배경 & 목표

DS003 재무 핵심 JSON 5개는 main 에 머지됨(PR #7). 이 spec 은 DS003 의 나머지 — XBRL 2개다.
형태가 서로 다르다: 택사노미는 JSON list(기존 `getListParams` 재사용), 원본파일은 바이너리 ZIP
(기존 `httpclient.GetBytes` 재사용, disclosure.DownloadDocument 와 동일 패턴). 이 그룹 완료 시
DS003 은 7/7 완료된다.

## API 표면 (docs + 실 API 검증)

- **xbrlTaxonomy** (XBRL택사노미 재무제표양식): `GET /api/xbrlTaxonomy.json`, 파라미터 `sj_div`
  (재무제표구분 코드). JSON `list[]`. sj_div 는 BS1~4(재무상태표)/IS(손익)/CIS(포괄손익)/CF(현금흐름)/
  SCE(자본변동) 등 ~20+ 코드 — 전체 표가 크롤링 docs 에 없어 **평이한 string** 으로 받고 doc 으로 안내.
  (검증: BS1 519행, IS1 361행 등 status 000)
- **fnlttXbrl** (재무제표 원본파일): `GET /api/fnlttXbrl.xml`, 파라미터 `rcept_no`+`reprt_code`.
  **Zip FILE (binary)** — JSON 없음. (검증: 삼성 2023 사업보고서 rcept_no=20240312000736 → ~990KB ZIP, PK 매직)

## 아키텍처

`report` 패키지에 새 파일 `xbrl.go` 추가. `client.Report` 는 이미 root 에 와이어링됨 → root 변경 불필요.
기존 헬퍼만 재사용(`getListParams[T]`, `httpclient.GetBytes`). 새 추상화 없음.

```
report/
  xbrl.go        # TaxonomyItem + XbrlTaxonomy + DownloadXbrl (신규)
  xbrl_test.go   # (신규)
  testdata/      # xbrlTaxonomy fixture 추가
README.md        # (수정) DS003 커버리지에 XBRL
integration_test.go  # (수정) XbrlTaxonomy 통합 케이스
```

## 2개 메서드 (report/xbrl.go)

```go
// TaxonomyItem 은 XBRL 택사노미 재무제표양식 (xbrlTaxonomy) 한 건.
type TaxonomyItem struct {
	SjDiv     string `json:"sj_div"`     // 재무제표구분
	AccountId string `json:"account_id"` // 계정ID
	AccountNm string `json:"account_nm"` // 계정명
	BsnsDe    string `json:"bsns_de"`    // 기준일 (YYYYMMDD)
	LabelKor  string `json:"label_kor"`  // 한글 출력명
	LabelEng  string `json:"label_eng"`  // 영문 출력명
	DataTp    string `json:"data_tp"`    // 데이터 유형 (일부 행에는 없음)
	IfrsRef   string `json:"ifrs_ref"`   // IFRS Reference
}

// XbrlTaxonomy 는 표준 XBRL 재무제표 택사노미 양식을 조회한다.
// sjDiv 는 재무제표구분 코드(BS1~4 재무상태표 / IS 손익계산서 / CIS 포괄손익 / CF 현금흐름표 /
// SCE 자본변동표 등). 잘못된 코드는 *opendart.APIError 로 반환된다.
func (c *Client) XbrlTaxonomy(ctx context.Context, sjDiv string) ([]TaxonomyItem, error) {
	return getListParams[TaxonomyItem](ctx, c.http, "/api/xbrlTaxonomy.json", map[string]string{"sj_div": sjDiv})
}

// DownloadXbrl 은 접수번호(rceptNo)+보고서코드로 재무제표 원본 XBRL(ZIP) 을 그대로 반환한다.
// 압축 해제·파싱은 호출자 몫.
func (c *Client) DownloadXbrl(ctx context.Context, rceptNo string, reprtCode ReportCode) ([]byte, error) {
	return c.http.GetBytes(ctx, "/api/fnlttXbrl.xml", map[string]string{
		"rcept_no":   rceptNo,
		"reprt_code": string(reprtCode),
	})
}
```

## 에러 처리

기존 재사용: 데이터 없음 → `opendart.ErrNoData`(택사노미), 그 외 status → `*opendart.APIError`.
`DownloadXbrl` 은 `GetBytes` 가 본문이 JSON 에러 envelope(`{` 시작)이면 `*APIError` 로 변환(ZIP 은 PK 시작이라 정상 통과).

## 테스트 전략

- `report/xbrl_test.go`: 기존 `report/client_test.go` 의 `newTestClient` 재사용.
  - `XbrlTaxonomy`: 실 응답 fixture(BS1) 디코딩 → 대표 필드(sj_div/account_id/label_kor) 검증.
  - `DownloadXbrl`: httptest 로 바이너리("PK..." 바이트) 반환 검증 + JSON 에러(010) → `*httpclient.APIError`
    검증 (disclosure.DownloadDocument 테스트와 동일 방식). rcept_no/reprt_code 쿼리 주입 확인.
- `integration_test.go` 에 `XbrlTaxonomy` 통합 케이스 추가(`//go:build integration`, sj_div=BS1).
- fixture 는 실 API 로 캡처해 임베드(BS1 첫 항목).

## 컨벤션 (기존 유지)

- 모든 struct 필드에 한글 코멘트, 도메인 주석 한국어.
- 표준 net/http(httpclient 재사용), 응답 캐싱 없음, UTF-8.
- README "커버리지" DS003 줄에 "XBRL 택사노미·원본파일" 추가, "(예정)" 에서 DS003 XBRL 제거.

## 비범위 (후속 plan)

- DS004~DS006 카테고리(지분공시·주요사항보고서·증권신고서).
- DS002 개인별 보수 Ver 2.0 2종.
