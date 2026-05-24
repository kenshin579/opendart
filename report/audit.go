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
	RceptNo                       string `json:"rcept_no"`                           // 접수번호
	CorpCls                       string `json:"corp_cls"`                           // 법인구분 (Y/K/N/E)
	CorpCode                      string `json:"corp_code"`                          // 고유번호
	CorpName                      string `json:"corp_name"`                          // 회사명
	SeNm                          string `json:"se_nm"`                              // 구분
	Tm                            string `json:"tm"`                                 // 회차
	PayDe                         string `json:"pay_de"`                             // 납입일
	PayAmount                     string `json:"pay_amount"`                         // 납입금액
	CptalUsePlan                  string `json:"cptal_use_plan"`                     // 자금사용 계획
	RealCptalUseSttus             string `json:"real_cptal_use_sttus"`               // 실제 자금사용 현황
	MtrptCptalUsePlanUseprps      string `json:"mtrpt_cptal_use_plan_useprps"`       // 주요사항보고서의 자금사용 계획(사용용도)
	MtrptCptalUsePlanPrcureAmount string `json:"mtrpt_cptal_use_plan_prcure_amount"` // 주요사항보고서의 자금사용 계획(조달금액)
	RealCptalUseDtlsCn            string `json:"real_cptal_use_dtls_cn"`             // 실제 자금사용 내역(내용)
	RealCptalUseDtlsAmount        string `json:"real_cptal_use_dtls_amount"`         // 실제 자금사용 내역(금액)
	DffrncOccrrncResn             string `json:"dffrnc_occrrnc_resn"`                // 차이발생 사유 등
	StlmDt                        string `json:"stlm_dt"`                            // 결산기준일
}

// PrivatePlacementFundUsage 는 사모자금의 사용내역을 조회한다.
func (c *Client) PrivatePlacementFundUsage(ctx context.Context, p ReportParams) ([]PrivatePlacementFundUsageItem, error) {
	return getList[PrivatePlacementFundUsageItem](ctx, c.http, "/api/prvsrpCptalUseDtls.json", p)
}
