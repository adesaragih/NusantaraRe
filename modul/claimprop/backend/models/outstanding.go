package models

// Untuk apa berkas ini: PORT TOMBOL "Save to issue RNM" DAN "Send to Acceptation" (tiket 01, 07; Section
// OutstandingClaim / InputAcceptation_Est): `SaveOutstanding_Act`, `SetOutstanding_Act`, `CheckNopolicy_Act`,
// `UpdateTableOS`.
//
// Bagian murni (halaman, pesan, rencana nomor, baris OS) ada di sini; nomor (sequence, GENERATE_SEQUENCE_NUMBER) dan
// tulisan OS_AKSEPTASI_KLAIM / JSON_KLAIM / MONITORING_KLAIM_LOG dijalankan services lewat repository, satu transaksi.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"
)

// Pesan VERBATIM.
const (
	PesanOldIDKosong = "Bussines Old ID Not Found, Please Contact IT!" // SaveOutstanding_Act langkah 4
	PesanCetakPLA    = "Please Print PLA"                              // langkah 36 (`pyWorkPage.Message`)
	PesanPolisNull   = "Policy No is null "                            // CheckNopolicy_Act langkah 2
)

// Teks riwayat VERBATIM.
const (
	TeksSaveOutstanding = "Save Outstanding"    // SaveOutstanding_Act langkah 23.1.2
	TeksKirimAkseptasi  = "Send To Acceptation" // CheckNopolicy_Act langkah 2
)

// Status OS_AKSEPTASI_KLAIM.STS_REJECT (= `TempOSAkseptasi.Type`; kolom TYPE berisi nilai yang sama `[data DEV
// 07-10-2026]`).
const (
	StsOSEstimasi    = "0" // SaveOutstanding_Act langkah 23.1.2
	StsOSAkseptasi   = "2" // SaveAcceptationTreaty_Act
	StsOSTutupBerkas = "4" // CloseClaimProp langkah 6.3
)

// Awalan nomor (procedure DBA; Go menulis ulang tanpa procedure).
const (
	KodeNomorSementara = "KT" // SaveOutstanding_Act langkah 7 `TempCFS.CARI30`
	JenisNomorKlaim    = "K"  // langkah 16 `ParamSeq.HASIL3+"K"`
	JenisNomorPLA      = "H"  // TryMakePLA_Act
	JenisNomorDLA      = "P"  // PrintDLATreatyIn
	// AwalanNomorCFS - GENERATE_NOCLMTREATYIN `'RNM-' || KODE || KODE_BIS ...` (UpdateTableOS langkah 2-3).
	AwalanNomorCFS = "RNM-"
	// TipeKodeProduksiNonLife - `GetKodeProdNonLife_SQL` `WHERE TYPE = 'NONLIFE'`.
	TipeKodeProduksiNonLife = "NONLIFE"
)

// RencanaNomor - nomor yang diterbitkan SaveOutstanding_Act.
type RencanaNomor struct {
	// Sementara - langkah 7-8: Policy No DAN NoClaim kosong -> GENERATE_NOCLMTRTYINTEMP ("KT.<yyyy>-<urut>").
	// ⚠️ Diterbitkan setiap kali tombol ditekan selama polis kosong, walau ClaimNo sudah berisi nomor sementara -
	// langkah 10 hanya mengisi ClaimNo yang kosong (nomor yang terbit sia-sia; ditiru apa adanya).
	Sementara bool
	// Klaim - langkah 16-18: NoClaim kosong DAN Policy No terisi -> PROC_GENERATE_SEQUENCE_NUMBER.
	Klaim bool
}

// RencanakanNomor membaca kebutuhan nomor SaveOutstanding_Act.
func RencanakanNomor(h *Halaman) RencanaNomor {
	pol, no := h.Ambil(CD+"PolicyData.PolicyNo"), h.Ambil(CD+"NoClaim")
	return RencanaNomor{Sementara: pol == "" && no == "", Klaim: no == "" && pol != ""}
}

// RakitNomorSementara = GENERATE_NOCLMTRTYINTEMP `KODE || '.' || TAHUN || '-' || lpad(seq,5,'0')`.
func RakitNomorSementara(tahun string, urut int) string {
	return fmt.Sprintf("%s.%s-%05d", KodeNomorSementara, tahun, urut)
}

// RakitNomorKlaim = SaveOutstanding_Act langkah 18: `ParamSeq.CARI2 + BusinessOldId + "." + HASIL1 + ".T" + HASIL2`
// (CARI2 = kode NONLIFE + "K", HASIL1 = MM.YYYY, HASIL2 = LPAD(urut,5,'0')).
func RakitNomorKlaim(jenis, oldID, mmYYYY string, urut int) string {
	return fmt.Sprintf("%s%s.%s.T%05d", jenis, oldID, mmYYYY, urut)
}

// RakitNomorCFS = GENERATE_NOCLMTREATYIN `'RNM-' || KODE || KODE_BIS || '.' || BULAN || '.' || TAHUN || '.' || 'T' ||
// lpad(seq,5,'0')` (UpdateTableOS langkah 2-3; KODE "K", KODE_BIS = `OutOldID.pxResults(1).BRANCH_CODE`).
// ⚠️ GetOldIDBusiness hanya memilih OLDID (alias CARI1) - BRANCH_CODE tidak pernah terisi, KODE_BIS selalu kosong.
func RakitNomorCFS(kodeBis, mm, yyyy string, urut int) string {
	return fmt.Sprintf("%sK%s.%s.%s.T%05d", AwalanNomorCFS, kodeBis, mm, yyyy, urut)
}

// TerapkanNomorSementara = langkah 9-10: ClaimNo kosong diisi nomor sementara.
func TerapkanNomorSementara(h *Halaman, nomor string) {
	if h.Ambil(CD+"ClaimNo") == "" {
		h.Setel(CD+"ClaimNo", nomor)
	}
}

// TerapkanNomorKlaim = langkah 19-20 (dan UpdateTableOS langkah 4-5): NoClaim kosong diisi nomor baru.
func TerapkanNomorKlaim(h *Halaman, nomor string) {
	if h.Ambil(CD+"NoClaim") == "" && nomor != "" {
		h.Setel(CD+"NoClaim", nomor)
	}
}

// PeriksaOldID = langkah 3-5: BusinessOldId dari BUSINESS.OLDID bila kosong; tetap kosong = pesan.
func PeriksaOldID(h *Halaman, oldID string) {
	if h.Ambil(OQ+"BusinessOldId") == "" && oldID != "" {
		h.Setel(OQ+"BusinessOldId", oldID)
	}
	if h.Ambil(OQ+"BusinessOldId") == "" {
		h.TambahPesan("", PesanOldIDKosong)
	}
}

// BarisOS - satu baris OS_AKSEPTASI_KLAIM, tanpa procedure; CASEID = pzInsKey (`KunciPegaLama`), NOCLAIM = NoClaim
// atau ClaimNo (langkah 23.1.5). `[keputusan work owner 08-10-2026]` DATA_JSON DIISI khusus tabel ini (meralat
// keputusan 07-10-2026 "DATA_JSON kosong"; JSON_KLAIM tetap tanpa DATA_JSON); kolom datar tetap diisi.
type BarisOS struct {
	CaseID, NoClaim, NoPolis, MasterID string
	StsReject                          string
	CurrencyID, Currency               string
	GrossValue, Value, KursValue       string
	PersenRNM                          string
	TypeLossID, TypeLoss               string
	CauseOfLossID, CauseOfLoss         string
	AcceptedNo                         string
	EstimationDate                     time.Time
	InsertOp                           string
	// DataJSON - `InputData.CARI1` (DataPega procedure): JSON halaman yang dikirim activity (JSONHalamanPega).
	DataJSON string
}

// KelasOSAkseptasi - pxObjClass halaman TempOSAkseptasi (SaveOutstanding_Act) dan InputParamOs (CloseClaimProp),
// Pages & Classes kedua activity; seluruh 6.596 DATA_JSON baris CLMP di DEV berkelas ini `[data DEV 08-10-2026]`.
const KelasOSAkseptasi = "ASM-FW-GCNMFW-Data-osAkseptasi"

// JSONHalamanPega = `@ASM.GetPageJSONString()` atas satu halaman. Fungsinya tidak diekspor; bentuknya dibaca dari
// DATA_JSON OS_AKSEPTASI_KLAIM warisan `[data DEV 08-10-2026, 6.596 baris CLMP]`:
//   - "{" LF, pasangan pertama, lalu setiap pasangan berikutnya diawali LF ",", ditutup LF "}" LF;
//   - kunci urut tanpa membedakan huruf besar (PersenRNM < PolicyNo < pxCreateOperator < pxObjClass < pzInsKey <
//     Type);
//   - semua nilai teks; properti kosong tidak ditulis (nol nilai "" di DEV, padahal CauseOfLoss selalu di-set);
//   - tanpa spasi; garis miring tidak di-escape (119 baris); escape lain mengikuti JSON baku tanpa escape HTML.
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

// teksJSON - literal teks JSON tanpa escape HTML (`<`, `>`, `&` apa adanya).
func teksJSON(s string) string {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(s) // menyandi string tidak pernah gagal
	return strings.TrimSuffix(buf.String(), "\n")
}

// tanggalStempel = `@getCurrentDateStamp()` (yyyyMMdd); di DEV selalu sama dengan tanggal kolom TANGGAL (519/519
// baris CLMP sejak 2025), maka dihitung dari saat yang sama dengan TANGGAL.
func tanggalStempel(t time.Time) string {
	return t.In(Jakarta).Format("20060102")
}

// nomorOS = InputData.CARI29 (langkah 23.1.4-5): NoClaim, atau ClaimNo bila keduanya kosong.
func nomorOS(h *Halaman) string {
	if no := h.Ambil(CD + "NoClaim"); no != "" {
		return no
	}
	return h.Ambil(CD + "ClaimNo")
}

// KirimEstimasi = SaveOutstanding_Act langkah 23: SETIAP estimasi belum terkirim (`PrintFaceClaim != 1`) ditandai
// terkirim (23.1.1, prakondisi nonaktif) dan melahirkan satu baris OS ber-STS_REJECT 0 (23.1.2-23.1.6, prakondisi
// nonaktif). ESTIMATIONDATE = `@getCurrentDateStamp()` (jam simpan, bukan tanggal estimasi baris). DATA_JSON =
// halaman TempOSAkseptasi (23.1.2 "set jsondata untuk table os_akseptasi_klaim", 23.1.3 GetPageJSONString).
func KirimEstimasi(k *Konteks, h *Halaman, kunciKasus string) []BarisOS {
	var out []BarisOS
	for _, e := range h.AmbilDaftar(DaftarEstimasi) {
		if e["PrintFaceClaim"] == "1" {
			continue
		}
		e["PrintFaceClaim"] = "1"
		out = append(out, BarisOS{
			CaseID: kunciKasus, NoClaim: nomorOS(h), NoPolis: h.Ambil(CD + "PolicyData.PolicyNo"),
			MasterID: h.Ambil(CD + "IDMaster"), StsReject: StsOSEstimasi,
			CurrencyID: e["CurrencyID"], Currency: e["Currency"], GrossValue: e["GrossEstimationPct"],
			Value: e["EstimationValue"], KursValue: e["KursValue"], PersenRNM: h.Ambil(TM + "RNMShareP"),
			TypeLossID: e["TypeLossID"], TypeLoss: e["TypeLoss"],
			CauseOfLossID: h.Ambil(CD + "CauseOfLossID"), CauseOfLoss: h.Ambil(CD + "CauseOfLoss"),
			EstimationDate: k.Sekarang, InsertOp: k.Pelaku,
			DataJSON: JSONHalamanPega(map[string]string{
				"CauseOfLoss":      h.Ambil(CD + "CauseOfLoss"),
				"CauseOfLossID":    h.Ambil(CD + "CauseOfLossID"),
				"EstimationDate":   tanggalStempel(k.Sekarang),
				"PersenRNM":        h.Ambil(TM + "RNMShareP"),
				"NoClaim":          h.Ambil(CD + "NoClaim"),
				"CurrencyID":       e["CurrencyID"],
				"Currency":         e["Currency"],
				"GrossValue":       e["GrossEstimationPct"],
				"Value":            e["EstimationValue"],
				"Type":             StsOSEstimasi,
				"TypeID":           e["Type"],
				"KursValue":        e["KursValue"],
				"TypeLossID":       e["TypeLossID"],
				"TypeLoss":         e["TypeLoss"],
				"pxCreateOperator": k.Pelaku,
				"pxObjClass":       KelasOSAkseptasi,
			}),
		})
	}
	return out
}

// BekukanDataLama = SaveOutstanding_Act langkah 24-28 (prakondisi nonaktif - selalu): claim amount `Note=Yes`, loss
// allocation / spreading / break QS `IsOldData=Yes`, interest `IsAdjVal=Yes`. Baris beku tidak dihitung ulang.
func BekukanDataLama(h *Halaman) {
	for _, b := range h.AmbilDaftar(DaftarClaimAmount) {
		b["Note"] = "Yes"
	}
	for _, j := range []string{DaftarLossAlloc, DaftarSpreading, DaftarBreakQS} {
		for _, b := range h.AmbilDaftar(j) {
			b["IsOldData"] = "Yes"
		}
	}
	for _, b := range h.AmbilDaftar(DaftarInterest) {
		b["IsAdjVal"] = "Yes"
	}
}

// SelesaiOutstanding = langkah 29-30.
func SelesaiOutstanding(k *Konteks, h *Halaman) {
	h.Setel("IsOutstanding", "1")
	h.Setel("IsCFS", "")
	h.Setel("ReCFS", "")
	k.Riwayat(h, TeksSaveOutstanding)
}

// IsiRetroDanPLA = langkah 31-36: tahun treaty (TREATYYEAR menurut treaty group dan tanggal mulai polis), lalu untuk
// SETIAP baris SpreadingClaim batas PLA (PROPORTIONALARRG kolom RP) dan daftar retro (TREATYREINSURER). FacRetroList
// yang kosong saat baris spreading itu diperiksa diisi seluruh retro baris itu (langkah 35.3; sesudah terisi baris
// spreading berikutnya tidak menambah). Cash call = batas baris spreading TERAKHIR (langkah 34 menimpa).
// Estimasi IDR melampaui cash call -> IsPLA 1, pesan cetak PLA, IsRealisation 2 (langkah 36).
func IsiRetroDanPLA(k *Konteks, h *Halaman) error {
	grup := h.Ambil(CD + "TreatyGroupID")
	tahun, err := k.Acuan.TahunTreaty(k.Ctxt(), grup, YMD(h.Ambil(CD+"PolicyData.StartDateTime")))
	if err != nil {
		return err
	}
	cash := ""
	for _, s := range h.AmbilDaftar(DaftarSpreading) {
		lim, _, err := k.Acuan.LimitPLA(k.Ctxt(), tahun, grup, s["TreatyType"])
		if err != nil {
			return err
		}
		cash = lim
		retro, err := k.Acuan.DaftarRetro(k.Ctxt(), s["TreatyType"], tahun, grup)
		if err != nil {
			return err
		}
		if len(h.AmbilDaftar(DaftarFacRetro)) > 0 {
			continue // 35.2 hanya menyetel FindData.CARI22 yang tidak dibaca
		}
		for _, r := range retro {
			h.TambahBaris(DaftarFacRetro, Baris{"ReinsurerID": r.ReinsurerID, "ReinsurerName": r.ReinsurerName,
				"PctShareAllObj": r.PctShare, "RiCommAllObj": r.RiComm})
		}
	}
	var kal Kalkulator
	tot := kal.H(h, CD+"TotalEstimasiIDR")
	cc := kal.Teks("CashCall", cash)
	if err := kal.Galat(); err != nil {
		return err
	}
	if Lebih(tot, cc) {
		h.Setel(CD+"IsPLA", "1")
		h.Setel("Message", PesanCetakPLA)
		h.Setel("IsRealisation", "2")
	}
	return nil
}

// BolehKirimAkseptasi - tombol "Send to Acceptation" tampil (`IsOutstanding = 1`) dan tidak nonaktif
// (`IsAcceptation==1`).
func BolehKirimAkseptasi(h *Halaman) bool {
	return h.Ambil("IsOutstanding") == "1" && h.Ambil("IsAcceptation") != "1"
}

// PeriksaPolisAkseptasi = CheckNopolicy_Act langkah 2-3: pesan bila Policy No kosong. Mengembalikan true bila
// UpdateTableOS boleh dijalankan (langkah 4).
func PeriksaPolisAkseptasi(h *Halaman) bool {
	if h.Ambil(CD+"PolicyData.PolicyNo") == "" {
		h.TambahPesan(CD+"PolicyData.PolicyNo", PesanPolisNull)
		return false
	}
	return true
}

// KlaimSementaraBerpolis = SetOutstanding_Act langkah 4-5 / UpdateTableOS langkah 7-8: Policy No terisi DAN ClaimNo
// memuat "KT".
// ⚠️ Kedua langkah memanggil `SaveOSKlaimTreaty_SQL` (procedure PEGA_JSON_OS_AKSEP_KLAIMTRT) yang TIDAK ADA di DEV
// dan sumbernya tidak terekspor - isinya tidak diketahui (OQ-CP-08). Yang dijalankan hanya salinan JSON_KLAIM.
func KlaimSementaraBerpolis(h *Halaman) bool {
	return h.Ambil(CD+"PolicyData.PolicyNo") != "" && strings.Contains(h.Ambil(CD+"ClaimNo"), KodeNomorSementara)
}

// SelesaiKirimAkseptasi = UpdateTableOS langkah 9 + CheckNopolicy_Act langkah 5.
func SelesaiKirimAkseptasi(k *Konteks, h *Halaman) {
	if h.Ambil(CD+"PolicyData.PolicyNo") != "" {
		h.Setel("IsAcceptation", "1")
		k.Riwayat(h, TeksKirimAkseptasi)
	}
}

// BolehSubmitOutstanding - tombol "Submit" Section OutstandingClaim: tampil bila `IsAcceptation==1`, menjalankan
// finishAssignment hanya bila Policy No dan NoClaim terisi (`pyActionConditions`).
func BolehSubmitOutstanding(h *Halaman) bool {
	return h.Ambil("IsAcceptation") == "1" && h.Ambil(CD+"PolicyData.PolicyNo") != "" && h.Ambil(CD+"NoClaim") != ""
}
