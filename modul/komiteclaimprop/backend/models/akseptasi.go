package models

// Untuk apa berkas ini: LANGKAH TINGKAT AKHIR yang membentuk data - baris OS akseptasi (`SaveAcceptation_Act`), retro
// dan batas DLA (`SaveAcceptationTreaty_TKMT`), muatan Kasir (`HitServiceToKasirKMT_Act`), dan email
// (`SendEmailKlaim_KMT`). Murni; layanan yang membaca acuan dan menulis.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"nusantarare/inti/backend/kontrak"
)

// KunciInstans - padanan `pzInsKey` kasus sistem baru di tabel warisan (OS_AKSEPTASI_KLAIM.CASEID, JSON_KLAIM.IDPEGA,
// HISTORYAKSEPTASIPEGA.ID_PEGA / ID_KOMITE, MONITORING_KLAIM_LOG.IDPEGA, DIRECTTOKASIR_LOG.IDPEGA): ID T_WORK_CLAIM
// apa adanya - sama dengan Claim Prop `models.KunciInstans` (keputusan work owner 06-10-2026, NB Treaty In).
func KunciInstans(id string) string { return id }

// DaftarAdjustment - jalur daftar baris adjustment halaman klaim induk (`ClaimData.AdjustmentList`).
const DaftarAdjustment = "ClaimData.AdjustmentList"

// jalurAdj - jalur properti baris adjustment ke-n halaman klaim induk.
func jalurAdj(n int, p string) string { return fmt.Sprintf("ClaimData.AdjustmentList(%d).%s", n, p) }

// AdjustmentKlaim - baris adjustment klaim induk yang sedang diputus (kosong bila tidak ada).
func AdjustmentKlaim(kl kontrak.KlaimTreaty) map[string]string {
	rows := kl.Daftar[DaftarAdjustment]
	if kl.Adjustment < 1 || kl.Adjustment > len(rows) {
		return map[string]string{}
	}
	return rows[kl.Adjustment-1]
}

// BarisOSAkseptasi - satu baris OS_AKSEPTASI_KLAIM Type 1 (`SaveAcceptation_Act` S1-S5, `SaveOSClaim_SQL` ->
// `PEGA_JSON_OS_AKSEP_KLAIM`), procedure tidak dipanggil. Kolom datar = kunci JSON `TempOSAkseptasi` yang bernama SAMA
// dengan kolomnya. DATA_JSON = halaman `TempOSAkseptasi` utuh (S4 `@ASM.GetPageJSONString()`) - keputusan work owner
// 08-10-2026 (diteruskan sesi ASIS CLAIM PROP): "isi json nya khusus os_akseptasi_klaim", termasuk baris akseptasi
// status 1 milik komite; meralat "DATA_JSON kosong" (07-10 dan prompt komite §3) HANYA untuk tabel ini. JSON_KLAIM tetap
// tanpa DATA_JSON.
type BarisOSAkseptasi struct {
	// Parameter procedure: CASEID (InputData.CARI2), NOCLAIM (CARI29), NOPOLIS (CARI3), STS_REJECT (CARI10 =
	// TempOSAkseptasi.Type = 1), MASTERID (CARI18); status konversi, STS_DLA, CLAIMOLD (CARI16/17/20) tidak pernah
	// diisi -> NULL.
	CaseID, NoClaim, NoPolis, StsReject, MasterID string
	// Kunci JSON = kolom.
	CauseOfLoss, CauseOfLossID, CurrencyID, Currency string
	GrossValue, PersenRNM, Value, KursValue          string
	Type, AcceptedNo, PaymentType                    string
	EstimationDate                                   time.Time
	// DataJSON - kolom DATA_JSON (CHECK `DATA_JSON IS JSON`).
	DataJSON string
}

// StsOSAkseptasi - `TempOSAkseptasi.Type := 1` (S2 ber-remark: Type 4 "final" tidak pernah jalan).
const StsOSAkseptasi = "1"

// KelasOSAkseptasi - pxObjClass halaman TempOSAkseptasi (Pages & Classes `SaveAcceptation_Act`; seluruh DATA_JSON baris
// CLMP di DEV berkelas ini, data DEV 08-10-2026 sesi ASIS CLAIM PROP).
const KelasOSAkseptasi = "ASM-FW-GCNMFW-Data-osAkseptasi"

// SusunOSAkseptasi = SaveAcceptation_Act S1 + S3-S4 atas klaim induk (sesudah S16.9 menulis AcceptedNo). `komiteID` =
// `pyWorkPage.pyID` (KomiteNo).
func SusunOSAkseptasi(kl kontrak.KlaimTreaty, acceptedNo, komiteID string, saat time.Time) BarisOSAkseptasi {
	b := AdjustmentKlaim(kl)
	v := kl.Nilai
	// S1: 17 properti TempOSAkseptasi (+ pxObjClass); EstimationDate = @getCurrentDateStamp() (yyyyMMdd, hari yang
	// sama dengan kolom TANGGAL - DEV 233/233), TypeID = PaymentType = adj.Type, TotalGross =
	// adj.ProposeAdjustmentValue, pxCreateOperator = PEMBUAT baris adjustment (bukan penyetuju).
	dataJSON := JSONHalamanPega(map[string]string{
		"CauseOfLoss": v["ClaimData.CauseOfLoss"], "CauseOfLossID": v["ClaimData.CauseOfLossID"],
		"EstimationDate": saat.In(Jakarta).Format("20060102"), "NoClaim": v["ClaimData.NoClaim"],
		"CurrencyID": b["CurrencyID"], "Currency": b["Currency"], "GrossValue": b["GrossAdjustment"],
		"PersenRNM": v["TreatyInMaster.RNMShareP"], "Value": b["AdjustmentValue"], "Type": StsOSAkseptasi,
		"TypeID": b["Type"], "KursValue": b["KursIDR"], "AcceptedNo": acceptedNo, "TotalGross": b["ProposeAdjustmentValue"],
		"PaymentType": b["Type"], "KomiteNo": komiteID, "pxCreateOperator": b["pxCreateOperator"],
		"pxObjClass": KelasOSAkseptasi,
	})
	return BarisOSAkseptasi{
		CaseID: KunciInstans(v["pyID"]), NoClaim: v["ClaimData.NoClaim"], NoPolis: v["ClaimData.PolicyData.PolicyNo"],
		StsReject: StsOSAkseptasi, MasterID: v["ClaimData.IDMaster"],
		CauseOfLoss: v["ClaimData.CauseOfLoss"], CauseOfLossID: v["ClaimData.CauseOfLossID"],
		CurrencyID: b["CurrencyID"], Currency: b["Currency"], GrossValue: b["GrossAdjustment"],
		PersenRNM: v["TreatyInMaster.RNMShareP"], Value: b["AdjustmentValue"], KursValue: b["KursIDR"],
		Type: StsOSAkseptasi, AcceptedNo: acceptedNo, PaymentType: b["Type"], EstimationDate: saat, DataJSON: dataJSON,
	}
}

// JSONHalamanPega = `@ASM.GetPageJSONString()` atas satu halaman (fungsi tidak diekspor; bentuk dibaca dari DATA_JSON
// OS_AKSEPTASI_KLAIM warisan - 6.596 baris CLMP di DEV, sesi ASIS CLAIM PROP 08-10-2026). DISALIN dari pola Claim Prop,
// tidak diimpor (batas modul):
//   - "{" LF, pasangan pertama, lalu setiap pasangan berikutnya diawali LF ",", ditutup LF "}" LF;
//   - kunci urut tanpa membedakan huruf besar (tie-break urutan byte);
//   - semua nilai teks; properti kosong tidak ditulis;
//   - tanpa spasi; garis miring tidak di-escape; escape JSON baku tanpa escape HTML.
func JSONHalamanPega(p map[string]string) string {
	kunci := make([]string, 0, len(p))
	for k, v := range p {
		if v != "" {
			kunci = append(kunci, k)
		}
	}
	sort.Slice(kunci, func(i, j int) bool {
		a, b := strings.ToLower(kunci[i]), strings.ToLower(kunci[j])
		if a != b {
			return a < b
		}
		return kunci[i] < kunci[j]
	})
	var b strings.Builder
	b.WriteString("{\n")
	for i, k := range kunci {
		if i > 0 {
			b.WriteString("\n,")
		}
		b.WriteString(teksJSON(k))
		b.WriteByte(':')
		b.WriteString(teksJSON(p[k]))
	}
	b.WriteString("\n}\n")
	return b.String()
}

// teksJSON - literal teks JSON tanpa escape HTML.
func teksJSON(s string) string {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(s) // menyandi string tidak pernah gagal
	return strings.TrimSuffix(buf.String(), "\n")
}

// BarisRetro - satu baris `ListInsurerDla` (`GetListRetro_Sql`: TREATYREINSURER reinsurerid / name / ricomm / pctshare).
type BarisRetro struct {
	ReinsurerID, Name, RiComm, PctShare string
}

// HasilRetro - tulisan `SaveAcceptationTreaty_TKMT`.
type HasilRetro struct {
	// FacRetro - S6.3: baris baru FacRetroList (hanya bila daftar itu kosong).
	FacRetro []map[string]string
	// IsFacRetro - S7: AdjustmentValue > batas RP (`GetLimitPLATreatyin` CARI7) -> 1 + pesan "Please Print DLA".
	IsFacRetro bool
	Pesan      string
}

// PesanCetakDLA - SaveAcceptationTreaty_TKMT S7 `TempOpenPage.Message` (VERBATIM).
const PesanCetakDLA = "Please Print DLA"

// SusunRetro = SaveAcceptationTreaty_TKMT S5-S7. `batas` = `LimitDla.pxResults(1).CARI7` putaran TERAKHIR S4
// (`@toDecimal("")` = 0 bila tak ada baris), `retro` = `ListInsurerDla` putaran terakhir, `adaRetro` = SizeRetro > 0.
func SusunRetro(adjValue, batas string, retro []BarisRetro, adaRetro bool) (HasilRetro, error) {
	var h HasilRetro
	if !adaRetro { // S6.3 (S6.2 / S6.2.1 hanya menyetel FindData.CARI22 yang tak dibaca)
		for _, r := range retro {
			h.FacRetro = append(h.FacRetro, map[string]string{"ReinsurerID": r.ReinsurerID, "ReinsurerName": r.Name,
				"PctShareAllObj": r.PctShare, "RiCommAllObj": r.RiComm})
		}
	}
	v, err := desimal("AdjustmentValue", adjValue)
	if err != nil {
		return HasilRetro{}, err
	}
	b, err := desimal("CARI7", batas)
	if err != nil {
		return HasilRetro{}, err
	}
	if v.Cmp(b) > 0 { // S7
		h.IsFacRetro, h.Pesan = true, PesanCetakDLA
	}
	return h, nil
}

var bukanAngka = regexp.MustCompile(`[^0-9]`)

// YMD - "yyyyMMdd" (8 digit pertama; SaveAcceptationTreaty_TKMT S1 `FindData.CARI55`) dari teks tanggal halaman ("2006-01-02[ 15:04:05]" atau "20060102...").
func YMD(s string) string {
	d := bukanAngka.ReplaceAllString(s, "")
	if len(d) > 8 {
		return d[:8]
	}
	return d
}

// KonfigurasiKasir - kode tetap muatan `SendAcceptationToKasir` (HitServiceToKasirKMT_Act S14.1.3: CompanyName,
// LjtdId, LdcId; S14.1.4 IsPEGASyariah LdcId). `[penyimpangan sadar]` (CLAUDE.md §10): dari konfigurasi
// berdokumen (MODUL.md), bukan literal kode.
type KonfigurasiKasir struct {
	CompanyName, LjtdID, LdcID, LdcIDSyariah string
}

// MuatanKasir - badan `TAllPaymentData` (HitServiceToKasirKMT_Act S14.3 Java, jalur CLMP S14.1).
type MuatanKasir struct {
	NoTrans, NoKlaim, LbuId, NoPolis, AcceptType, Kepada, AccountNo, TglAksep string
	Nett, Deductible, KaliDeduct                                              string
	StsSyariah, CompanyName, LjtdId, LdcId, StsAp, LkuId, LbgID               string
	TglBolehBayar, Email, UserInput                                           string
}

// PanjangNoAksepCLMP - S14.1 `@length(.AcceptedNo) = 23 || 24`; selainnya keluar (kode 6).
func PanjangNoAksepCLMP(no string) bool { return len(no) == 23 || len(no) == 24 }

// tglAksepTeks = `@substring(d,6,8)+"-"+@substring(d,4,6)+"-"+@substring(d,0,4)` atas "yyyyMMdd".
func tglAksepTeks(ymd string) string {
	if len(ymd) < 8 {
		return ""
	}
	return ymd[6:8] + "-" + ymd[4:6] + "-" + ymd[0:4]
}

// TglBolehBayar = S14.1.5: hari > 25 -> tanggal 01 dua bulan sesudahnya, selainnya hari yang sama bulan berikutnya;
// bulan 13 -> "1"; tahun naik bila bulan hasil 01 DAN bulan berjalan Desember.
func TglBolehBayar(ymd string, bulanKini int) string {
	if len(ymd) < 8 {
		return ""
	}
	hari, _ := strconv.Atoi(ymd[6:8])
	bulan, _ := strconv.Atoi(ymd[4:6])
	tahun, _ := strconv.Atoi(ymd[0:4])
	bulan++
	if hari > 25 {
		bulan++
	}
	if bulan >= 13 {
		bulan = 1
	}
	hs := strconv.Itoa(hari)
	if hari > 25 {
		hs = "01"
	}
	if len(hs) < 2 {
		hs = "0" + hs
	}
	bs := strconv.Itoa(bulan)
	if len(bs) < 2 {
		bs = "0" + bs
	}
	if bulan == 1 && bulanKini == 12 {
		tahun++
	}
	return hs + "-" + bs + "-" + strconv.Itoa(tahun)
}

// SusunMuatanKasir = HitServiceToKasirKMT_Act S14.1.1-S14.1.5 (jalur CLMP). `b` = baris adjustment induk (sesudah S13
// IDOfBank), `email` = `gl.f_get_email(TreatyInMaster.CedingID)` (S14.1.2), `syariah` = IsPEGASyariah (OQ, false).
func SusunMuatanKasir(kl kontrak.KlaimTreaty, b map[string]string, email string, syariah bool, cfg KonfigurasiKasir,
	pelaku string, saat time.Time) MuatanKasir {
	ldc := cfg.LdcID
	if syariah {
		ldc = cfg.LdcIDSyariah
	}
	tgl := YMD(b["AcceptedDate"])
	user := b["pxCreateOperator"]
	if user == "" {
		user = pelaku // `pxRequestor.pxUserIdentifier`
	}
	v := kl.Nilai
	return MuatanKasir{
		NoTrans: b["AcceptedNo"], NoKlaim: v["ClaimData.NoClaim"], LbuId: v["OfferFacIn.QuotationData.BusinessOldId"],
		NoPolis: v["ClaimData.PolicyData.PolicyNo"], AcceptType: b["Type"], Kepada: b["PayableTo"],
		AccountNo: bukanAngka.ReplaceAllString(b["NoAccount"], ""), TglAksep: tglAksepTeks(tgl),
		Nett: b["AdjustmentValue"], Deductible: b["IndividualRiskRNM"], KaliDeduct: b["IndividualRiskPercentage"],
		StsSyariah: "0", CompanyName: cfg.CompanyName, LjtdId: cfg.LjtdID, LdcId: ldc, StsAp: "0",
		LkuId: b["CurrencyID"], LbgID: b["IDOfBank"], TglBolehBayar: TglBolehBayar(tgl, int(saat.In(Jakarta).Month())),
		Email: email, UserInput: user,
	}
}

// StatusKasirSukses - S14.5 `@if(.StatusServiceKasir.ReponseCode="1","Akseptasi Sudah Masuk ke Kasir",...)`.
const StatusKasirSukses = "Akseptasi Sudah Masuk ke Kasir"

// bulanEmail - SendEmailKlaim_KMT S5 (VERBATIM, termasuk ejaannya).
var bulanEmail = map[string]string{"01": "Januari", "02": "Febuari", "03": "Maret", "04": "April", "05": "Mei",
	"06": "Juni", "07": "July", "08": "Agustus", "09": "September", "10": "Oktober", "11": "November", "12": "Desember"}

// TanggalEmail = S5 `Temp.CARI23`: "dd <Bulan> yyyy" dari DateOfLoss.
func TanggalEmail(dol string) string {
	d := YMD(dol)
	if len(d) < 8 {
		return ""
	}
	return d[6:8] + " " + bulanEmail[d[4:6]] + " " + d[0:4]
}

// SubjekEmail = SendEmailKlaim_KMT S12 / S14 / S15: "<awalan>Pengajuan Akseptasi : <ID klaim>/<KomiteNo>
// <InsuredName> DOL <tanggal>". Awalan "" (penyetuju berikut), "(Approval) " (tingkat akhir), "(Reject) " (tolak).
func SubjekEmail(jenis string, kl kontrak.KlaimTreaty, komiteNo string) string {
	awal := ""
	switch jenis {
	case EmailPembuatSetuju:
		awal = "(Approval) "
	case EmailPembuatTolak:
		awal = "(Reject) "
	}
	v := kl.Nilai
	return awal + "Pengajuan Akseptasi : " + v["pyID"] + "/" + komiteNo + " " + v["ClaimData.InsuredName"] + " DOL " +
		TanggalEmail(v["ClaimData.DateOfLoss"])
}

// AdaSpreadingAdjustment - SendEmailKlaim_KMT S1: `.SpreadingAdjustment` kosong -> keluar (tanpa email).
func AdaSpreadingAdjustment(kl kontrak.KlaimTreaty) bool {
	return len(kl.Daftar[jalurAdj(kl.Adjustment, "SpreadingAdjustment")]) > 0
}

// SpreadingAdjustmentTerakhir - TreatyType baris terakhir `.SpreadingAdjustment` (SaveAcceptationTreaty_TKMT S4:
// LimitDla / ListInsurerDla ditimpa setiap putaran; S5-S6 membaca putaran terakhir). Kosong bila tak ada baris.
func SpreadingAdjustmentTerakhir(kl kontrak.KlaimTreaty) (string, bool) {
	rows := kl.Daftar[jalurAdj(kl.Adjustment, "SpreadingAdjustment")]
	if len(rows) == 0 {
		return "", false
	}
	return rows[len(rows)-1]["TreatyType"], true
}

// AdaFacRetro - S5 `@SizeOfPropertyList(TempOpenPage.ClaimData.FacRetroList) > 0`.
func AdaFacRetro(kl kontrak.KlaimTreaty) bool { return len(kl.Daftar["ClaimData.FacRetroList"]) > 0 }

// NoKomite - `.KomiteNo` adjustment (SetKomiteNo_Act: nomor kasus komite); kosong -> ID kasus komite.
func NoKomite(kl kontrak.KlaimTreaty, komiteID string) string {
	if no := strings.TrimSpace(AdjustmentKlaim(kl)["KomiteNo"]); no != "" {
		return no
	}
	return komiteID
}
