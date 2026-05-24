# OpenDART DS002 정기보고서 주요정보 — 감사·자금·출자 그룹 설계

- 작성일: 2026-05-24
- 모듈: `github.com/kenshin579/opendart`
- 범위: **DS002 감사·자금·출자 6개 API** (`report` 패키지 확장)

## 배경 & 목표

DS002 공통 추상화(`report` 패키지: `getList`/`ReportParams`)와 지분·주식·배당 7개(PR #3),
증권 발행·미상환 6개(PR #4)는 main 에 머지됨. 이 spec 은 DS002 세 번째 그룹인 **감사·자금·출자
6개**다. 6개 모두 표준 요청(`corp_code`+`bsns_year`+`reprt_code`)과 list 응답을 가지므로 **새
추상화 없이** 기존 `getList[T]` 패턴을 재사용한다 — 각 엔드포인트 = item struct + 한 줄 메서드.
이 그룹 완료 시 DS002 는 19/30 (임원·보수 11개만 남음).

## API 표면 (docs 기반 사실)

- 요청: 전 API 공통 `crtfc_key`(자동 주입) + `corp_code` + `bsns_year` + `reprt_code` (추가 파라미터 없음).
- 응답: 공통 envelope(`status`/`message`) + `list[]`. 숫자는 콤마 문자열, 빈 값 "-".
- 감사 관련 3개(감사의견/감사용역/비감사용역)는 응답 list 항목에 `bsns_year` 필드 포함.

## 아키텍처

`report` 패키지에 새 파일 `audit.go` 추가 (기존 `equity.go`/`securities.go` 와 동일하게 그룹당 1파일).
`client.Report` 는 이미 root 에 와이어링됨 → **root 변경 불필요**. 메서드는 기존 `getList[T]`/
`ReportParams` 재사용.

```
report/
  audit.go        # 6개 메서드 + item struct (신규)
  audit_test.go   # 6개 fixture 테스트 (신규, newTestClient 재사용)
  testdata/       # 6개 실 응답 JSON fixture 추가
README.md         # (수정) DS002 커버리지에 감사·자금·출자
integration_test.go  # (수정) AuditOpinion 통합 케이스
```

## 6개 엔드포인트 (report/audit.go)

각 메서드: `func (c *Client) X(ctx, p ReportParams) ([]XItem, error) { return getList[XItem](ctx, c.http, "<path>", p) }`.

| 메서드 | 한글 | 엔드포인트 |
|--------|------|-----------|
| `AuditOpinion` | 회계감사인의 명칭 및 감사의견 | `/api/accnutAdtorNmNdAdtOpinion.json` |
| `AuditServiceContract` | 감사용역체결현황 | `/api/adtServcCnclsSttus.json` |
| `NonAuditServiceContract` | 회계감사인과의 비감사용역 계약체결 현황 | `/api/accnutAdtorNonAdtServcCnclsSttus.json` |
| `OtherCorpInvestment` | 타법인 출자현황 | `/api/otrCprInvstmntSttus.json` |
| `PublicOfferingFundUsage` | 공모자금의 사용내역 | `/api/pssrpCptalUseDtls.json` |
| `PrivatePlacementFundUsage` | 사모자금의 사용내역 | `/api/prvsrpCptalUseDtls.json` |

```go
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

// OtherCorpInvestmentItem 은 타법인 출자현황 (otrCprInvstmntSttus) 한 건.
type OtherCorpInvestmentItem struct {
	RceptNo                            string `json:"rcept_no"`                                  // 접수번호
	CorpCls                            string `json:"corp_cls"`                                  // 법인구분 (Y/K/N/E)
	CorpCode                           string `json:"corp_code"`                                 // 고유번호
	CorpName                           string `json:"corp_name"`                                 // 회사명
	InvPrm                             string `json:"inv_prm"`                                   // 법인명 (피출자 법인)
	FrstAcqsDe                         string `json:"frst_acqs_de"`                              // 최초 취득 일자
	InvstmntPurps                      string `json:"invstmnt_purps"`                            // 출자 목적
	FrstAcqsAmount                     string `json:"frst_acqs_amount"`                          // 최초 취득 금액
	BsisBlceQy                         string `json:"bsis_blce_qy"`                              // 기초 잔액 수량
	BsisBlceQotaRt                     string `json:"bsis_blce_qota_rt"`                         // 기초 잔액 지분율
	BsisBlceAcntbkAmount               string `json:"bsis_blce_acntbk_amount"`                   // 기초 잔액 장부 가액
	IncrsDcrsAcqsDspsQy                string `json:"incrs_dcrs_acqs_dsps_qy"`                   // 증가 감소 취득 처분 수량
	IncrsDcrsAcqsDspsAmount            string `json:"incrs_dcrs_acqs_dsps_amount"`               // 증가 감소 취득 처분 금액
	IncrsDcrsEvlLstmn                  string `json:"incrs_dcrs_evl_lstmn"`                      // 증가 감소 평가 손액
	TrmendBlceQy                       string `json:"trmend_blce_qy"`                            // 기말 잔액 수량
	TrmendBlceQotaRt                   string `json:"trmend_blce_qota_rt"`                       // 기말 잔액 지분율
	TrmendBlceAcntbkAmount             string `json:"trmend_blce_acntbk_amount"`                 // 기말 잔액 장부 가액
	RecentBsnsYearFnnrSttusTotAssets   string `json:"recent_bsns_year_fnnr_sttus_tot_assets"`    // 최근 사업연도 재무현황 총자산
	RecentBsnsYearFnnrSttusThstrmNtpf  string `json:"recent_bsns_year_fnnr_sttus_thstrm_ntpf"`   // 최근 사업연도 재무현황 당기순이익
	StlmDt                             string `json:"stlm_dt"`                                   // 결산기준일
}

// PublicOfferingFundUsageItem 은 공모자금의 사용내역 (pssrpCptalUseDtls) 한 건.
type PublicOfferingFundUsageItem struct {
	RceptNo                   string `json:"rcept_no"`                     // 접수번호
	CorpCls                   string `json:"corp_cls"`                     // 법인구분 (Y/K/N/E)
	CorpCode                  string `json:"corp_code"`                    // 고유번호
	CorpName                  string `json:"corp_name"`                    // 회사명
	SeNm                      string `json:"se_nm"`                        // 구분
	Tm                        string `json:"tm"`                           // 회차
	PayDe                     string `json:"pay_de"`                       // 납입일
	PayAmount                 string `json:"pay_amount"`                   // 납입금액
	OnDclrtCptalUsePlan       string `json:"on_dclrt_cptal_use_plan"`      // 신고서상 자금사용 계획
	RealCptalUseSttus         string `json:"real_cptal_use_sttus"`         // 실제 자금사용 현황
	RsCptalUsePlanUseprps     string `json:"rs_cptal_use_plan_useprps"`    // 증권신고서 등의 자금사용 계획(사용용도)
	RsCptalUsePlanPrcureAmount string `json:"rs_cptal_use_plan_prcure_amount"` // 증권신고서 등의 자금사용 계획(조달금액)
	RealCptalUseDtlsCn        string `json:"real_cptal_use_dtls_cn"`       // 실제 자금사용 내역(내용)
	RealCptalUseDtlsAmount    string `json:"real_cptal_use_dtls_amount"`   // 실제 자금사용 내역(금액)
	DffrncOccrrncResn         string `json:"dffrnc_occrrnc_resn"`          // 차이발생 사유 등
	StlmDt                    string `json:"stlm_dt"`                      // 결산기준일
}

// PrivatePlacementFundUsageItem 은 사모자금의 사용내역 (prvsrpCptalUseDtls) 한 건.
type PrivatePlacementFundUsageItem struct {
	RceptNo                      string `json:"rcept_no"`                       // 접수번호
	CorpCls                      string `json:"corp_cls"`                       // 법인구분 (Y/K/N/E)
	CorpCode                     string `json:"corp_code"`                      // 고유번호
	CorpName                     string `json:"corp_name"`                      // 회사명
	SeNm                         string `json:"se_nm"`                          // 구분
	Tm                           string `json:"tm"`                             // 회차
	PayDe                        string `json:"pay_de"`                         // 납입일
	PayAmount                    string `json:"pay_amount"`                     // 납입금액
	CptalUsePlan                 string `json:"cptal_use_plan"`                 // 자금사용 계획
	RealCptalUseSttus            string `json:"real_cptal_use_sttus"`           // 실제 자금사용 현황
	MtrptCptalUsePlanUseprps     string `json:"mtrpt_cptal_use_plan_useprps"`   // 주요사항보고서의 자금사용 계획(사용용도)
	MtrptCptalUsePlanPrcureAmount string `json:"mtrpt_cptal_use_plan_prcure_amount"` // 주요사항보고서의 자금사용 계획(조달금액)
	RealCptalUseDtlsCn           string `json:"real_cptal_use_dtls_cn"`         // 실제 자금사용 내역(내용)
	RealCptalUseDtlsAmount       string `json:"real_cptal_use_dtls_amount"`     // 실제 자금사용 내역(금액)
	DffrncOccrrncResn            string `json:"dffrnc_occrrnc_resn"`            // 차이발생 사유 등
	StlmDt                       string `json:"stlm_dt"`                        // 결산기준일
}
```

각 메서드 (audit.go):
```go
func (c *Client) AuditOpinion(ctx context.Context, p ReportParams) ([]AuditOpinionItem, error) {
	return getList[AuditOpinionItem](ctx, c.http, "/api/accnutAdtorNmNdAdtOpinion.json", p)
}
func (c *Client) AuditServiceContract(ctx context.Context, p ReportParams) ([]AuditServiceContractItem, error) {
	return getList[AuditServiceContractItem](ctx, c.http, "/api/adtServcCnclsSttus.json", p)
}
func (c *Client) NonAuditServiceContract(ctx context.Context, p ReportParams) ([]NonAuditServiceContractItem, error) {
	return getList[NonAuditServiceContractItem](ctx, c.http, "/api/accnutAdtorNonAdtServcCnclsSttus.json", p)
}
func (c *Client) OtherCorpInvestment(ctx context.Context, p ReportParams) ([]OtherCorpInvestmentItem, error) {
	return getList[OtherCorpInvestmentItem](ctx, c.http, "/api/otrCprInvstmntSttus.json", p)
}
func (c *Client) PublicOfferingFundUsage(ctx context.Context, p ReportParams) ([]PublicOfferingFundUsageItem, error) {
	return getList[PublicOfferingFundUsageItem](ctx, c.http, "/api/pssrpCptalUseDtls.json", p)
}
func (c *Client) PrivatePlacementFundUsage(ctx context.Context, p ReportParams) ([]PrivatePlacementFundUsageItem, error) {
	return getList[PrivatePlacementFundUsageItem](ctx, c.http, "/api/prvsrpCptalUseDtls.json", p)
}
```

## 에러 처리

기존 재사용: 데이터 없음 → `opendart.ErrNoData`, 그 외 status → `*opendart.APIError`.

## 테스트 전략

- `report/audit_test.go`: 기존 `report/client_test.go` 의 `newTestClient` 재사용.
- 6개 메서드 각각: 실 응답 JSON fixture 디코딩 → 대표 필드 매핑 검증.
- fixture 는 실 API 로 캡처해 임베드(계획 작성 단계). 데이터 없는 종목은 데이터 있는 회사/연도로 캡처.
- `integration_test.go` 에 `AuditOpinion` 통합 케이스 추가(`//go:build integration`).

## 컨벤션 (기존 유지)

- 모든 item struct 필드에 한글 코멘트, 도메인 주석 한국어.
- 표준 net/http(httpclient 재사용), 응답 캐싱 없음, 숫자 coercion 없음(콤마 string 유지), UTF-8.
- README "커버리지" DS002 줄에 "감사의견·감사용역·비감사용역·타법인 출자·공모/사모자금 사용내역" 추가.

## 비범위 (후속 plan)

- DS002 마지막 그룹: 임원·보수(11개) — 임원/직원/미등기임원 보수/사외이사 변동 + 이사·감사 전체 보수(3) + 개인별 보수(±Ver2.0, 4).
- DS003~DS006 카테고리.
- 신규 예제(기존 `examples/report` 로 충분).
