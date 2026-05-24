# OpenDART DS005 주요사항보고서 주요정보 — 양수도 Sub-2 «기타» 설계

- 작성일: 2026-05-24
- 모듈: `github.com/kenshin579/opendart`
- 범위: **DS005 양수도 그룹 중 Sub-2 2개 API** (`material` 패키지 확장)

## 배경 & 목표

DS005 양수도 그룹(10개)을 2 sub-group 으로 분할(A안). Sub-1 «실물 양수도» 8개는 PR #14 main 머지됨.
이 spec 은 **Sub-2 «기타» 2개**다(자산양수도(기타)·풋백옵션 + 주식교환·이전 결정). 기존 `MaterialParams`
+ `httpclient.GetList[T]` 재사용, root 변경 없음. 신규 파일 `material/transfer_etc.go`.

## API 표면 (docs 기반 사실)

- 2개 모두 동일 요청 `corp_code`+`bgn_de`+`end_de` (= `MaterialParams`), JSON `list[]`.
- 자산양수도(기타)·풋백옵션(astInhtrfEtcPtbkOpt 6필드 — 머리 4 + rp_rsn + ast_inhtrf_prc), 주식교환·이전(stkExtrDecsn 56필드).
- 값은 문자열(금액 콤마, 비율, 빈 값 "-").

## 아키텍처

```
material/
  transfer_etc.go       # 2개 메서드 + item struct (신규)
  transfer_etc_test.go  # 2개 fixture 테스트 (신규)
  testdata/             # 2개 fixture
README.md               # (수정) DS005 커버리지에 양수도 Sub-2
integration_test.go     # (수정) 통합 케이스 1~2개 (ErrNoData skip)
```

각 메서드: `func (c *Client) X(ctx, p MaterialParams) ([]XItem, error) { return httpclient.GetList[XItem](ctx, c.http, "<path>", p.toMap()) }`.

## 2개 메서드 (material/transfer_etc.go)

| 메서드 | 한글 | 엔드포인트 | 필드 |
|--------|------|-----------|------|
| `OtherAssetTransferPutbackOption` | 자산양수도(기타), 풋백옵션 | `/api/astInhtrfEtcPtbkOpt.json` | 6 |
| `StockExchangeTransfer` | 주식교환·이전 결정 | `/api/stkExtrDecsn.json` | 56 |

```go
// OtherAssetTransferPutbackOptionItem 은 자산양수도(기타), 풋백옵션 (astInhtrfEtcPtbkOpt) 한 건.
type OtherAssetTransferPutbackOptionItem struct {
	RceptNo      string `json:"rcept_no"`        // 접수번호
	CorpCls      string `json:"corp_cls"`        // 법인구분 (Y/K/N/E)
	CorpCode     string `json:"corp_code"`       // 고유번호
	CorpName     string `json:"corp_name"`       // 회사명
	RpRsn        string `json:"rp_rsn"`          // 보고 사유
	AstInhtrfPrc string `json:"ast_inhtrf_prc"`  // 자산양수ㆍ도 가액
}

// StockExchangeTransferItem 은 주식교환·이전 결정 (stkExtrDecsn) 한 건.
type StockExchangeTransferItem struct {
	RceptNo            string `json:"rcept_no"`               // 접수번호
	CorpCls            string `json:"corp_cls"`               // 법인구분 (Y/K/N/E)
	CorpCode           string `json:"corp_code"`              // 고유번호
	CorpName           string `json:"corp_name"`              // 회사명
	ExtrSen            string `json:"extr_sen"`               // 구분
	ExtrStn            string `json:"extr_stn"`               // 교환ㆍ이전 형태
	ExtrTgcmpCmpnm     string `json:"extr_tgcmp_cmpnm"`       // 교환ㆍ이전 대상법인(회사명)
	ExtrTgcmpRp        string `json:"extr_tgcmp_rp"`          // 교환ㆍ이전 대상법인(대표자)
	ExtrTgcmpMbsn      string `json:"extr_tgcmp_mbsn"`        // 교환ㆍ이전 대상법인(주요사업)
	ExtrTgcmpRlCmpn    string `json:"extr_tgcmp_rl_cmpn"`     // 교환ㆍ이전 대상법인(회사와의 관계)
	ExtrTgcmpTisstkOstk string `json:"extr_tgcmp_tisstk_ostk"` // 교환ㆍ이전 대상법인(발행주식총수(주)(보통주식))
	ExtrTgcmpTisstkCstk string `json:"extr_tgcmp_tisstk_cstk"` // 교환ㆍ이전 대상법인(발행주식총수(주)(종류주식))
	RbsnfdtlTast       string `json:"rbsnfdtl_tast"`          // 대상법인 최근 사업연도 요약재무(원)(자산총계)
	RbsnfdtlTdbt       string `json:"rbsnfdtl_tdbt"`          // 대상법인 최근 사업연도 요약재무(원)(부채총계)
	RbsnfdtlTeqt       string `json:"rbsnfdtl_teqt"`          // 대상법인 최근 사업연도 요약재무(원)(자본총계)
	RbsnfdtlCpt        string `json:"rbsnfdtl_cpt"`           // 대상법인 최근 사업연도 요약재무(원)(자본금)
	ExtrRt             string `json:"extr_rt"`                // 교환ㆍ이전 비율
	ExtrRtBs           string `json:"extr_rt_bs"`             // 교환ㆍ이전 비율 산출근거
	ExevlAtn           string `json:"exevl_atn"`              // 외부평가에 관한 사항(외부평가 여부)
	ExevlBsRs          string `json:"exevl_bs_rs"`            // 외부평가에 관한 사항(근거 및 사유)
	ExevlIntn          string `json:"exevl_intn"`             // 외부평가에 관한 사항(외부평가기관의 명칭)
	ExevlPd            string `json:"exevl_pd"`               // 외부평가에 관한 사항(외부평가 기간)
	ExevlOp            string `json:"exevl_op"`               // 외부평가에 관한 사항(외부평가 의견)
	ExtrPp             string `json:"extr_pp"`                // 교환ㆍ이전 목적
	ExtrscExtrctrd     string `json:"extrsc_extrctrd"`        // 교환ㆍ이전일정(교환ㆍ이전계약일)
	ExtrscShddstd      string `json:"extrsc_shddstd"`         // 교환ㆍ이전일정(주주확정기준일)
	ExtrscShclspdBgd   string `json:"extrsc_shclspd_bgd"`     // 교환ㆍ이전일정(주주명부 폐쇄기간(시작일))
	ExtrscShclspdEdd   string `json:"extrsc_shclspd_edd"`     // 교환ㆍ이전일정(주주명부 폐쇄기간(종료일))
	ExtrscExtropRcpdBgd string `json:"extrsc_extrop_rcpd_bgd"` // 교환ㆍ이전일정(반대의사 통지접수기간(시작일))
	ExtrscExtropRcpdEdd string `json:"extrsc_extrop_rcpd_edd"` // 교환ㆍ이전일정(반대의사 통지접수기간(종료일))
	ExtrscGmtsckPrd    string `json:"extrsc_gmtsck_prd"`      // 교환ㆍ이전일정(주주총회 예정일자)
	ExtrscAprskhExpdBgd string `json:"extrsc_aprskh_expd_bgd"` // 교환ㆍ이전일정(주식매수청구권 행사기간(시작일))
	ExtrscAprskhExpdEdd string `json:"extrsc_aprskh_expd_edd"` // 교환ㆍ이전일정(주식매수청구권 행사기간(종료일))
	ExtrscOsprpdBgd    string `json:"extrsc_osprpd_bgd"`      // 교환ㆍ이전일정(구주권제출기간(시작일))
	ExtrscOsprpdEdd    string `json:"extrsc_osprpd_edd"`      // 교환ㆍ이전일정(구주권제출기간(종료일))
	ExtrscTrspprpd     string `json:"extrsc_trspprpd"`        // 교환ㆍ이전일정(매매거래정지예정기간)
	ExtrscTrspprpdBgd  string `json:"extrsc_trspprpd_bgd"`    // 교환ㆍ이전일정(매매거래정지예정기간(시작일))
	ExtrscTrspprpdEdd  string `json:"extrsc_trspprpd_edd"`    // 교환ㆍ이전일정(매매거래정지예정기간(종료일))
	ExtrscExtrdt       string `json:"extrsc_extrdt"`          // 교환ㆍ이전일정(교환ㆍ이전일자)
	ExtrscNstkdlprd    string `json:"extrsc_nstkdlprd"`       // 교환ㆍ이전일정(신주권교부예정일)
	ExtrscNstklstprd   string `json:"extrsc_nstklstprd"`      // 교환ㆍ이전일정(신주의 상장예정일)
	AtextrCpcmpnm      string `json:"atextr_cpcmpnm"`         // 교환ㆍ이전 후 완전모회사명
	AprskhPlnprc       string `json:"aprskh_plnprc"`          // 주식매수청구권(매수예정가격)
	AprskhPymPlpdMth   string `json:"aprskh_pym_plpd_mth"`    // 주식매수청구권(지급예정시기, 지급방법)
	AprskhLmt          string `json:"aprskh_lmt"`             // 주식매수청구권(제한 관련 내용)
	AprskhCtref        string `json:"aprskh_ctref"`           // 주식매수청구권(계약에 미치는 효력)
	BdlstAtn           string `json:"bdlst_atn"`              // 우회상장 해당 여부
	OtcprBdlstSfAtn    string `json:"otcpr_bdlst_sf_atn"`     // 타법인의 우회상장 요건 충족 여부
	Bddd               string `json:"bddd"`                   // 이사회결의일(결정일)
	OdAAtT             string `json:"od_a_at_t"`              // 사외이사 참석여부(참석(명))
	OdAAtB             string `json:"od_a_at_b"`              // 사외이사 참석여부(불참(명))
	AdtAAtn            string `json:"adt_a_atn"`              // 감사(사외이사가 아닌 감사위원) 참석여부
	PoptCtrAtn         string `json:"popt_ctr_atn"`           // 풋옵션 등 계약 체결여부
	PoptCtrCn          string `json:"popt_ctr_cn"`            // 계약내용
	RsSmAtn            string `json:"rs_sm_atn"`              // 증권신고서 제출대상 여부
	ExSmR              string `json:"ex_sm_r"`                // 제출을 면제받은 경우 그 사유
}
```

각 메서드는 위 패턴으로 작성한다.

## 에러 처리

기존 재사용: 데이터 없음 → `opendart.ErrNoData`, 그 외 status → `*opendart.APIError`.

## 테스트 전략

- `material/transfer_etc_test.go`: 기존 `material/client_test.go` 의 `newTestClient` 재사용
  (route map 값은 bare 파일명).
- 2개 메서드 각각 fixture 디코딩 → 대표 필드 검증.
- fixture 는 실 API 캡처 권장(불가 시 docs 스키마 일치 샘플).
- `integration_test.go` 에 통합 케이스 1~2개(`//go:build integration`, ErrNoData skip 허용).

## 컨벤션 (기존 유지)

- 모든 item struct 필드에 한글 코멘트, 도메인 주석 한국어.
- 표준 net/http(httpclient 재사용), 응답 캐싱 없음, string 유지, UTF-8.
- README "커버리지" DS005 줄에 "자산양수도(기타)·풋백옵션 · 주식교환·이전" 추가.

## 비범위 (후속 plan)

- DS005 합병·분할(회사합병/회사분할/회사분할합병 3) / 해외상장(상장 결정·상장·상장폐지 결정·상장폐지 4).
- DS006 증권신고서(6). DS002 개인별 보수 Ver 2.0 2종.
