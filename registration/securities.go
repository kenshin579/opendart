package registration

import (
	"context"
	"encoding/json"

	"github.com/kenshin579/opendart/internal/httpclient"
)

// RsGeneralItem 은 증권신고서 일반사항 그룹 항목(지분증권/증권예탁증권 공통).
type RsGeneralItem struct {
	RceptNo  string `json:"rcept_no"`  // 접수번호
	CorpCls  string `json:"corp_cls"`  // 법인구분 (Y/K/N/E)
	CorpCode string `json:"corp_code"` // 고유번호
	CorpName string `json:"corp_name"` // 회사명
	Sbd      string `json:"sbd"`       // 청약기일
	Pymd     string `json:"pymd"`      // 납입기일
	Sband    string `json:"sband"`     // 청약공고일
	Asand    string `json:"asand"`     // 배정공고일
	Asstd    string `json:"asstd"`     // 배정기준일
	Exstk    string `json:"exstk"`     // 신주인수권에 관한 사항(행사대상증권)
	Exprc    string `json:"exprc"`     // 신주인수권에 관한 사항(행사가격)
	Expd     string `json:"expd"`      // 신주인수권에 관한 사항(행사기간)
	RptRcpn  string `json:"rpt_rcpn"`  // 주요사항보고서(접수번호)
}

// RsSecurityTypeItem 은 증권신고서 증권의종류 그룹 항목(지분증권/증권예탁증권 공통).
type RsSecurityTypeItem struct {
	RceptNo  string `json:"rcept_no"`  // 접수번호
	CorpCls  string `json:"corp_cls"`  // 법인구분 (Y/K/N/E)
	CorpCode string `json:"corp_code"` // 고유번호
	CorpName string `json:"corp_name"` // 회사명
	Stksen   string `json:"stksen"`    // 증권의종류
	Stkcnt   string `json:"stkcnt"`    // 증권수량
	Fv       string `json:"fv"`        // 액면가액
	Slprc    string `json:"slprc"`     // 모집(매출)가액
	Slta     string `json:"slta"`      // 모집(매출)총액
	Slmthn   string `json:"slmthn"`    // 모집(매출)방법
}

// RsUnderwriterItem 은 증권신고서 인수인정보 그룹 항목(지분증권/증권예탁증권 공통).
type RsUnderwriterItem struct {
	RceptNo  string `json:"rcept_no"`  // 접수번호
	CorpCls  string `json:"corp_cls"`  // 법인구분 (Y/K/N/E)
	CorpCode string `json:"corp_code"` // 고유번호
	CorpName string `json:"corp_name"` // 회사명
	Actsen   string `json:"actsen"`    // 인수인구분
	Actnmn   string `json:"actnmn"`    // 인수인명
	Stksen   string `json:"stksen"`    // 증권의종류
	Udtcnt   string `json:"udtcnt"`    // 인수수량
	Udtamt   string `json:"udtamt"`    // 인수금액
	Udtprc   string `json:"udtprc"`    // 인수대가
	Udtmth   string `json:"udtmth"`    // 인수방법
}

// RsFundUsageItem 은 증권신고서 자금의사용목적 그룹 항목(지분증권/증권예탁증권 공통).
type RsFundUsageItem struct {
	RceptNo  string `json:"rcept_no"`  // 접수번호
	CorpCls  string `json:"corp_cls"`  // 법인구분 (Y/K/N/E)
	CorpCode string `json:"corp_code"` // 고유번호
	CorpName string `json:"corp_name"` // 회사명
	Se       string `json:"se"`        // 구분
	Amt      string `json:"amt"`       // 금액
}

// RsSellerItem 은 증권신고서 매출인에관한사항 그룹 항목(지분증권/증권예탁증권 공통).
type RsSellerItem struct {
	RceptNo   string `json:"rcept_no"`   // 접수번호
	CorpCls   string `json:"corp_cls"`   // 법인구분 (Y/K/N/E)
	CorpCode  string `json:"corp_code"`  // 고유번호
	CorpName  string `json:"corp_name"`  // 회사명
	Hdr       string `json:"hdr"`        // 보유자
	RlCmp     string `json:"rl_cmp"`     // 회사와의관계
	BfslHdstk string `json:"bfsl_hdstk"` // 매출전보유증권수
	Slstk     string `json:"slstk"`      // 매출증권수
	AtslHdstk string `json:"atsl_hdstk"` // 매출후보유증권수
}

// EquityRetailPutbackOptionItem 은 지분증권 일반청약자환매청구권 그룹 항목(지분증권 전용).
type EquityRetailPutbackOptionItem struct {
	RceptNo  string `json:"rcept_no"`  // 접수번호
	CorpCls  string `json:"corp_cls"`  // 법인구분 (Y/K/N/E)
	CorpCode string `json:"corp_code"` // 고유번호
	CorpName string `json:"corp_name"` // 회사명
	Grtrs    string `json:"grtrs"`     // 부여사유
	Exavivr  string `json:"exavivr"`   // 행사가능 투자자
	Grtcnt   string `json:"grtcnt"`    // 부여수량
	Expd     string `json:"expd"`      // 행사기간
	Exprc    string `json:"exprc"`     // 행사가격
}

// EquitySecuritiesRegistration 은 지분증권 증권신고서(estkRs)의 그룹별 항목.
type EquitySecuritiesRegistration struct {
	General             []RsGeneralItem                 // 일반사항
	SecurityTypes       []RsSecurityTypeItem            // 증권의종류
	Underwriters        []RsUnderwriterItem             // 인수인정보
	FundUsage           []RsFundUsageItem               // 자금의사용목적
	Sellers             []RsSellerItem                  // 매출인에관한사항
	RetailPutbackOption []EquityRetailPutbackOptionItem // 일반청약자환매청구권
}

// EquitySecurities 는 지분증권 증권신고서(DS006)를 조회한다.
func (c *Client) EquitySecurities(ctx context.Context, p Params) (*EquitySecuritiesRegistration, error) {
	groups, err := httpclient.GetGroups(ctx, c.http, "/api/estkRs.json", p.toMap())
	if err != nil {
		return nil, err
	}
	out := &EquitySecuritiesRegistration{}
	for _, g := range groups {
		var derr error
		switch g.Title {
		case "일반사항":
			derr = json.Unmarshal(g.List, &out.General)
		case "증권의종류":
			derr = json.Unmarshal(g.List, &out.SecurityTypes)
		case "인수인정보":
			derr = json.Unmarshal(g.List, &out.Underwriters)
		case "자금의사용목적":
			derr = json.Unmarshal(g.List, &out.FundUsage)
		case "매출인에관한사항":
			derr = json.Unmarshal(g.List, &out.Sellers)
		case "일반청약자환매청구권":
			derr = json.Unmarshal(g.List, &out.RetailPutbackOption)
		}
		if derr != nil {
			return nil, derr
		}
	}
	return out, nil
}
