package models

// Klausul - SATU tabel, 25 jenis, aturan per jenis - tiket 08 Treaty Contract Out.
//
// Untuk apa berkas ini: tabel aturan wajib-isi per jenis klausul di KODE
// (AC 34), gerbangnya (AC 33, 35, 36), dan dua perilaku turunan yang nyata di
// Pega: Rp/Usd anak (`HitungRpUsd`) dan total Pct anak.
//
// Sensus 29-09-2026 atas 25 aktivitas `Activity/SaveTreatyArr*` - HANYA langkah
// hidup (bukan `//`) dengan prasyarat AKTIF (`pyStepsPreCondition = true`):
//
//	`Property-Set-Messages` berprasyarat `<medan>==""` dengan WhenTrue=2
//
// ⛔ Ralat atas tabel tiket: `TerritorialLimit` pada CoinsPanel, MaxCoinsPanel,
// MinLOL, MinLOLMB, dan ExclutionTreaty berasal dari prasyarat NONAKTIF yang
// merujuk halaman lain (`InputTreatyTerr`) - residu salin-tempel, bukan gerbang.
//
// Pemetaan `TreatyDescID` dari `Activity/SetKirimIDDesc.xml` langkah 3-27:
// induk dan anak BERBAGI DescID, dibedakan `ParentReinsTypeID` ("00" = induk).
//
// Dibaca sesudah: tco_arrangement.go, tco_reinsurer.go.

import (
	"errors"
	"fmt"
	"strings"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/utils"
)

// Nama medan klausul - properti Pega VERBATIM (AC 32).
const (
	MedanReinsTypeID      = "ReinsTypeID"
	MedanLine             = "Line"
	MedanRp               = "Rp"
	MedanUsd              = "Usd"
	MedanPct              = "Pct"
	MedanPctMe            = "PctMe"
	MedanYdcf             = "Ydcf"
	MedanMethod           = "Method"
	MedanTerritorialLimit = "TerritorialLimit"
	MedanCoInsMin         = "CoIns_Min"
	MedanCoInsMax         = "CoIns_Max"
	MedanTreatyLimit      = "TreatyLimit"
	MedanIDOccupation     = "ID_Occupation"
	MedanOccupation       = "Occupation"
	MedanIDClause         = "ID_Clause"
	MedanClause           = "Clause"
	MedanLayer            = "Layer"
)

// DescID master `TREATYDESC` - `SetKirimIDDesc.xml` (`InputData.CARIDESC == ...`).
const (
	DescTreatyLimit     = "10001"
	DescPortfolio       = "10002"
	DescPLA             = "10003"
	DescCashLossLimit   = "10004"
	DescBordereAux      = "10005"
	DescFacIn           = "10006"
	DescExGratia        = "10007"
	DescTerrLimit       = "10008"
	DescEPI             = "10009"
	DescRicomm          = "10010"
	DescProfitComm      = "10011"
	DescClaimCoorp      = "10012"
	DescExclutionTreaty = "10013"
	DescCoinsPanel      = "10014"
	DescMinLOL          = "10015"
	DescMaxCoinsPanel   = "10016"
	DescLimitMB         = "10017"
	DescMinLOLMB        = "10018"
)

// Subjenis ExclutionTreaty - `Param.Type` yang dikirim tombol Save tiap
// sub-bagian (`GridTreatyArrangementExclutionTreaty*`).
const (
	SubjenisOccupation = "Occupation"
	SubjenisClause     = "Clause"
	SubjenisObject     = "Object"
	SubjenisPeriode    = "Periode"
)

// AturanKlausul adalah aturan satu jenis klausul.
type AturanKlausul struct {
	// Jenis - nama aktivitas Pega tanpa `SaveTreatyArr`/`_Act` (AC 32).
	Jenis    string
	DescID   string
	Anak     bool
	Subjenis string
	// Medan - yang TAMPIL dan dapat diisi di form jenis itu.
	Medan []string
	// Wajib - wajib-isi, urutan langkah Pega.
	Wajib []string
	// Turunan - dihitung server, tidak diterima dari klien (Rp/Usd anak).
	Turunan []string
	// KunciDobel - medan yang menjadikan baris "sudah pernah diinput" (AC 30).
	KunciDobel []string
	// Ditahan - alasan jenis ini BELUM dapat disimpan (AC 36).
	Ditahan string
	// BatasTotalAnak - total Pct anak <= 100 (langkah "Menjumlahkan Total PCT").
	BatasTotalAnak bool
	// PeringatanSpreading - `TreatyTestChildTotal_Act` sesudah simpan.
	PeringatanSpreading bool
	// Berkurs - form jenis ini menuntut kurs berlaku (`NewTreatyArr*`: "Tidak ada
	// Nilai Kurs di Tahun : "), tiket 11.
	Berkurs bool
	// Konversi - arah konversi kurs di form: `KonversiRpKeUsd` (Usd turunan,
	// `HitungRpUsd_depan`) atau `KonversiDuaArah` (`CalculateTSIExcludeTreaty`).
	Konversi string
	// Sumber - aktivitas Pega VERBATIM.
	Sumber string
}

var (
	medanReinsRpUsd = []string{MedanReinsTypeID, MedanLine, MedanRp, MedanUsd}
	wajibReinsRpUsd = []string{MedanReinsTypeID, MedanRp, MedanUsd}
	medanAnak       = []string{MedanReinsTypeID, MedanPct}
	wajibAnak       = []string{MedanReinsTypeID, MedanPct, MedanRp, MedanUsd}
	turunanAnak     = []string{MedanRp, MedanUsd}
	// turunanIndukKurs - `Usd` form induk hanya dibaca (mis. `GridTreatyArrangementEpi.xml`
	// b3775) dan diisi `HitungRpUsd_depan`: `Usd = Rp / Kurs` (tiket 11).
	turunanIndukKurs = []string{MedanUsd}
	kunciReins       = []string{MedanReinsTypeID}
)

// Arah konversi kurs di form (tiket 11).
const (
	KonversiRpKeUsd = "RpKeUsd"
	KonversiDuaArah = "DuaArah"
)

const alasanDitahan = "aturan wajib-isi jenis %s belum ditetapkan Product + UW (AC 36)"

func induk(jenis, desc, sumber string) AturanKlausul {
	return AturanKlausul{Jenis: jenis, DescID: desc, Medan: medanReinsRpUsd, Wajib: wajibReinsRpUsd,
		Turunan: turunanIndukKurs, KunciDobel: kunciReins, Berkurs: true, Konversi: KonversiRpKeUsd, Sumber: sumber}
}

func anak(jenis, desc, sumber string, peringatan bool) AturanKlausul {
	return AturanKlausul{Jenis: jenis, DescID: desc, Anak: true, Medan: medanAnak, Wajib: wajibAnak,
		Turunan: turunanAnak, KunciDobel: kunciReins, BatasTotalAnak: true, PeringatanSpreading: peringatan,
		Berkurs: true, Sumber: sumber}
}

// AturanKlausulTCO - SELURUH jenis. Satu tempat, di kode (AC 34).
var AturanKlausulTCO = []AturanKlausul{
	induk("TreatyLimit", DescTreatyLimit, "SaveTreatyArrTreatyLimit_Act"),
	anak("TreatyLimitChild", DescTreatyLimit, "SaveTreatyArrTreatyLimitChild_Act", true),
	{Jenis: "Portfolio", DescID: DescPortfolio, Medan: []string{MedanMethod},
		Ditahan: fmt.Sprintf(alasanDitahan, "Portfolio") + "; bentuk penyimpanan daftar bersarangnya ikut OQ",
		Sumber:  "SaveTreatyArrPortfolio_Act"},
	induk("PLA", DescPLA, "SaveTreatyArrPLA_Act"),
	anak("PLAList", DescPLA, "SaveTreatyArrPLAList_Act", false),
	induk("CashLossLimit", DescCashLossLimit, "SaveTreatyArrCashLossLimit_Act"),
	anak("CashLossLimitList", DescCashLossLimit, "SaveTreatyArrCashLossLimitList_Act", false),
	{Jenis: "BordereAux", DescID: DescBordereAux, Medan: []string{MedanMethod}, Wajib: []string{MedanMethod},
		KunciDobel: []string{MedanMethod}, Sumber: "SaveTreatyArrBordereAux_Act"},
	induk("FacIn", DescFacIn, "SaveTreatyArrFacIn_Act"),
	anak("FacInList", DescFacIn, "SaveTreatyArrFacInList_Act", false),
	induk("ExGratia", DescExGratia, "SaveTreatyArrExGratia_Act"),
	anak("ExGratiaChildList", DescExGratia, "SaveTreatyArrExGratiaChildList_Act", false),
	{Jenis: "TerrLimit", DescID: DescTerrLimit, Medan: []string{MedanTerritorialLimit},
		Wajib: []string{MedanTerritorialLimit}, KunciDobel: []string{MedanTerritorialLimit},
		Sumber: "SaveTreatyArrTerrLimit_Act"},
	induk("EPI", DescEPI, "SaveTreatyArrEPI_Act"),
	anak("EpiList", DescEPI, "SaveTreatyArrEpiList_Act", false),
	{Jenis: "Ricomm", DescID: DescRicomm, Medan: []string{MedanReinsTypeID, MedanPct, MedanMethod},
		Wajib: []string{MedanReinsTypeID, MedanPct, MedanMethod}, KunciDobel: kunciReins,
		Sumber: "SaveTreatyArrRicomm_Act"},
	{Jenis: "ProfitComm", DescID: DescProfitComm, Medan: []string{MedanReinsTypeID, MedanPct, MedanPctMe, MedanYdcf},
		Wajib: []string{MedanPct, MedanReinsTypeID, MedanPctMe}, KunciDobel: kunciReins,
		Sumber: "SaveTreatyArrProfitComm_Act"},
	induk("ClaimCoorp", DescClaimCoorp, "SaveTreatyArrClaimCoorp_Act"),
	anak("ClaimCoorpChild", DescClaimCoorp, "SaveTreatyArrClaimCoorpChild_Act", false),
	// ExclutionTreaty - satu jenis, empat subjenis (`Param.Type`).
	{Jenis: "ExclutionTreaty", DescID: DescExclutionTreaty, Subjenis: SubjenisOccupation,
		Medan:      []string{MedanIDOccupation, MedanOccupation, MedanLine, MedanUsd, MedanRp},
		Wajib:      []string{MedanIDOccupation, MedanOccupation, MedanLine, MedanUsd, MedanRp},
		KunciDobel: []string{MedanIDOccupation}, Konversi: KonversiDuaArah, Sumber: "SaveTreatyArrExclutionTreaty_Act"},
	{Jenis: "ExclutionTreaty", DescID: DescExclutionTreaty, Subjenis: SubjenisClause,
		Medan: []string{MedanIDClause, MedanClause}, Wajib: []string{MedanIDClause, MedanClause},
		KunciDobel: []string{MedanIDClause}, Sumber: "SaveTreatyArrExclutionTreaty_Act"},
	// ⚠️ Object dan Periode: Save mengirim `Type` KOSONG, sehingga di Pega nol
	// validasi berjalan. Satu-satunya medan formnya diwajibkan
	// [keputusan work owner 29-09-2026] (OQ-TCO-14, ditutup) - baris exclusion dengan
	// satu-satunya medannya kosong tidak bermakna. Kunci dobel per jenis
	// (`KunciDobel`) dikonfirmasi keputusan yang sama.
	{Jenis: "ExclutionTreaty", DescID: DescExclutionTreaty, Subjenis: SubjenisObject,
		Medan: []string{MedanPct}, Wajib: []string{MedanPct}, KunciDobel: []string{MedanPct},
		Sumber: "SaveTreatyArrExclutionTreaty_Act"},
	{Jenis: "ExclutionTreaty", DescID: DescExclutionTreaty, Subjenis: SubjenisPeriode,
		Medan: []string{MedanLayer}, Wajib: []string{MedanLayer}, KunciDobel: []string{MedanLayer},
		Sumber: "SaveTreatyArrExclutionTreaty_Act"},
	{Jenis: "CoinsPanel", DescID: DescCoinsPanel, Medan: []string{MedanCoInsMin, MedanCoInsMax, MedanTreatyLimit},
		Wajib: []string{MedanCoInsMin, MedanCoInsMax, MedanTreatyLimit}, KunciDobel: []string{MedanCoInsMin, MedanCoInsMax},
		Sumber: "SaveTreatyArrCoinsPanel_Act"},
	{Jenis: "MinLOL", DescID: DescMinLOL, Medan: []string{MedanPct}, Wajib: []string{MedanPct},
		KunciDobel: []string{MedanPct}, Sumber: "SaveTreatyArrMinLOL"},
	{Jenis: "MaxCoinsPanel", DescID: DescMaxCoinsPanel, Medan: []string{MedanCoInsMax}, Wajib: []string{MedanCoInsMax},
		KunciDobel: []string{MedanCoInsMax}, Sumber: "SaveTreatyArrMaxCoinsPanel"},
	{Jenis: "LimitMB", DescID: DescLimitMB, Ditahan: fmt.Sprintf(alasanDitahan, "LimitMB"),
		Medan:  []string{MedanIDOccupation, MedanPct, MedanPctMe, MedanRp, MedanUsd, MedanTerritorialLimit},
		Sumber: "SaveTreatyArrLimitMB_Act"},
	{Jenis: "MinLOLMB", DescID: DescMinLOLMB, Medan: []string{MedanPct}, Wajib: []string{MedanPct},
		KunciDobel: []string{MedanPct}, Sumber: "SaveTreatyArrMinLOLMB"},
}

var (
	// ErrKlausulJenisTakDikenal - DescID/anak/subjenis tanpa aturan di kode.
	ErrKlausulJenisTakDikenal = errors.New("models: jenis klausul tanpa aturan wajib-isi di kode")
	// ErrKlausulDitahan - jenis yang aturannya belum ditetapkan (AC 36).
	ErrKlausulDitahan = errors.New("models: jenis klausul belum dapat disimpan")
	// ErrKlausulMedanWajib - medan wajib kosong; pesan menyebut medannya (AC 35).
	ErrKlausulMedanWajib = errors.New("models: medan wajib klausul kosong")
	// ErrMedanBukanMilikJenis - medan dikirim untuk jenis yang tidak memilikinya.
	ErrMedanBukanMilikJenis = errors.New("models: medan bukan milik jenis klausul ini")
	// ErrTotalPctAnakMelebihi100 - `ASMMessageTotalPct` (teks Pega tidak diekspor).
	ErrTotalPctAnakMelebihi100 = errors.New("models: total Pct baris anak tidak boleh lebih dari 100")
)

// PesanSpreadingTCO - VERBATIM `TreatyTestChildTotal_Act.xml` langkah 5.
const PesanSpreadingTCO = "Please make sure spreading is 100%"

// CariAturanKlausul mencari aturan menurut DescID, anak/induk, dan subjenis.
func CariAturanKlausul(descID string, anak bool, subjenis string) (AturanKlausul, bool) {
	for _, a := range AturanKlausulTCO {
		if a.DescID == strings.TrimSpace(descID) && a.Anak == anak && a.Subjenis == strings.TrimSpace(subjenis) {
			return a, true
		}
	}
	return AturanKlausul{}, false
}

// AturanJenisKlausul mencari aturan menurut nama jenis VERBATIM.
func AturanJenisKlausul(jenis, subjenis string) (AturanKlausul, bool) {
	for _, a := range AturanKlausulTCO {
		if a.Jenis == jenis && (a.Subjenis == subjenis || (subjenis == "" && a.Subjenis == SubjenisOccupation)) {
			return a, true
		}
	}
	return AturanKlausul{}, false
}

// medanDesimal - medan berkolom NUMBER(38,8) (AC 51, 52).
var medanDesimal = map[string]bool{MedanRp: true, MedanUsd: true, MedanPct: true, MedanPctMe: true,
	MedanTreatyLimit: true, MedanCoInsMin: true, MedanCoInsMax: true}

func (a AturanKlausul) punya(medan string) bool {
	for _, m := range a.Medan {
		if m == medan {
			return true
		}
	}
	return false
}

// NilaiMedanKlausul membaca satu medan sebagai teks (desimal dalam bentuk 'f').
func NilaiMedanKlausul(k KlausulTreaty, medan string) string {
	d := func(x *apd.Decimal) string {
		if x == nil {
			return ""
		}
		return x.Text('f')
	}
	switch medan {
	case MedanReinsTypeID:
		return k.ReinsTypeID
	case MedanLine:
		return k.Line
	case MedanRp:
		return d(k.Rp)
	case MedanUsd:
		return d(k.Usd)
	case MedanPct:
		return d(k.Pct)
	case MedanPctMe:
		return d(k.PctMe)
	case MedanYdcf:
		return k.Ydcf
	case MedanMethod:
		return k.Method
	case MedanTerritorialLimit:
		return k.TerritorialLimit
	case MedanCoInsMin:
		return d(k.CoinsMin)
	case MedanCoInsMax:
		return d(k.CoinsMax)
	case MedanTreatyLimit:
		return d(k.TreatyLimit)
	case MedanIDOccupation:
		return k.IDOccupation
	case MedanOccupation:
		return k.Occupation
	case MedanIDClause:
		return k.IDClause
	case MedanClause:
		return k.Clause
	case MedanLayer:
		return k.Layer
	}
	return ""
}

// IsiMedanKlausulTCO menuliskan masukan klien ke baris - HANYA medan milik
// jenis itu, desimal dinormalkan sekali di batas masukan.
func IsiMedanKlausulTCO(a AturanKlausul, k *KlausulTreaty, masuk map[string]string) error {
	for medan, teks := range masuk {
		if !a.punya(medan) {
			return fmt.Errorf("%w: %s untuk jenis %s", ErrMedanBukanMilikJenis, medan, a.Jenis)
		}
		t := strings.TrimSpace(teks)
		if medanDesimal[medan] {
			var d *apd.Decimal
			if t != "" {
				var err error
				if d, err = UraiDesimalMasukTCO(medan, t); err != nil {
					return err
				}
			}
			switch medan {
			case MedanRp:
				k.Rp = d
			case MedanUsd:
				k.Usd = d
			case MedanPct:
				k.Pct = d
			case MedanPctMe:
				k.PctMe = d
			case MedanTreatyLimit:
				k.TreatyLimit = d
			case MedanCoInsMin:
				k.CoinsMin = d
			case MedanCoInsMax:
				k.CoinsMax = d
			}
			continue
		}
		switch medan {
		case MedanReinsTypeID:
			k.ReinsTypeID = t
		case MedanLine:
			k.Line = t
		case MedanYdcf:
			k.Ydcf = t
		case MedanMethod:
			k.Method = t
		case MedanTerritorialLimit:
			k.TerritorialLimit = t
		case MedanIDOccupation:
			k.IDOccupation = t
		case MedanOccupation:
			k.Occupation = t
		case MedanIDClause:
			k.IDClause = t
		case MedanClause:
			k.Clause = t
		case MedanLayer:
			k.Layer = t
		}
	}
	return nil
}

// UraiDesimalMasukTCO menormalkan satu desimal uang/persen klausul di batas
// masukan: koma -> titik, tolak koma+titik, <= 8 desimal, <= 30 digit bulat.
// Tanpa batas rentang - Pega tidak menegakkannya untuk klausul.
func UraiDesimalMasukTCO(nama, teks string) (*apd.Decimal, error) {
	t := strings.TrimSpace(teks)
	if strings.Contains(t, ",") {
		if strings.Contains(t, ".") {
			return nil, fmt.Errorf("%w: %s %q memuat koma dan titik sekaligus", ErrPersenBukanDesimal, nama, t)
		}
		t = strings.ReplaceAll(t, ",", ".")
	}
	d, _, err := apd.NewFromString(t)
	if err != nil || d.Form != apd.Finite {
		return nil, fmt.Errorf("%w: %s %q", ErrPersenBukanDesimal, nama, teks)
	}
	r := new(apd.Decimal).Set(d)
	r.Reduce(r)
	if r.Exponent < -skalaPersenTCO || r.NumDigits()+int64(r.Exponent) > digitBulatPersen {
		return nil, fmt.Errorf("%w: %s %q melampaui NUMBER(38,8)", ErrPersenBukanDesimal, nama, teks)
	}
	return d, nil
}

// PeriksaKlausulTCO menjalankan gerbang wajib-isi jenis itu (AC 33, 35, 36).
func PeriksaKlausulTCO(a AturanKlausul, k KlausulTreaty) error {
	if a.Ditahan != "" {
		return fmt.Errorf("%w: %s", ErrKlausulDitahan, a.Ditahan)
	}
	var kosong []string
	for _, m := range a.Wajib {
		if strings.TrimSpace(NilaiMedanKlausul(k, m)) == "" {
			kosong = append(kosong, m)
		}
	}
	if len(kosong) > 0 {
		return fmt.Errorf("%w - jenis %s: %s", ErrKlausulMedanWajib, a.Jenis, strings.Join(kosong, ", "))
	}
	return nil
}

var (
	seratusKlausul = apd.New(100, 0)
)

// RpUsdAnakTCO - `Activity/HitungRpUsd.xml`: `Usd = Pct * Usd_induk / 100`,
// `Rp = Pct * Rp_induk / 100`, dikuantisasi ke 8 desimal (NUMBER(38,8)).
//
// ⛔ Angka TURUNAN inilah yang tersimpan - klien tidak mengirimnya.
func RpUsdAnakTCO(pct, rpInduk, usdInduk *apd.Decimal) (rp, usd *apd.Decimal, err error) {
	hitung := func(x *apd.Decimal) (*apd.Decimal, error) {
		if pct == nil || x == nil {
			return nil, nil
		}
		hasil := new(apd.Decimal)
		if _, err := utils.DecimalContext().Mul(hasil, pct, x); err != nil {
			return nil, err
		}
		if _, err := utils.DecimalContext().Quo(hasil, hasil, seratusKlausul); err != nil {
			return nil, err
		}
		kuantum := new(apd.Decimal)
		// SALINAN per panggilan (temuan /code-review): menulis Rounding ke konteks
		// bersama adalah data race.
		k := utils.DecimalContext()
		k.Rounding = apd.RoundHalfUp
		if _, err := k.Quantize(kuantum, hasil, -skalaPersenTCO); err != nil {
			return nil, err
		}
		return kuantum, nil
	}
	if rp, err = hitung(rpInduk); err != nil {
		return nil, nil, fmt.Errorf("models: menghitung Rp anak: %w", err)
	}
	if usd, err = hitung(usdInduk); err != nil {
		return nil, nil, fmt.Errorf("models: menghitung Usd anak: %w", err)
	}
	return rp, usd, nil
}

// PeriksaTotalAnakTCO - langkah "Menjumlahkan Total PCT" + "Valdiasi Total agar
// tidak boleh lebih dari 100": total Pct anak BARU > 100 ditolak.
//
// ⚠️ Pega menjumlahkan lewat teks berkoma lalu membagi 100 - kompensasi atas
// `toDecimal("100,00")` yang membaca koma sebagai pemisah ribuan. Maksudnya
// jelas (total <= 100); yang ditiru MAKSUD itu, dengan desimal persis.
func PeriksaTotalAnakTCO(lain []*apd.Decimal, baru *apd.Decimal) (*apd.Decimal, error) {
	total, err := TotalShareTCO(append(append([]*apd.Decimal{}, lain...), baru))
	if err != nil {
		return nil, err
	}
	if total.Cmp(seratusKlausul) > 0 {
		return total, fmt.Errorf("%w (total %s)", ErrTotalPctAnakMelebihi100, total.Text('f'))
	}
	return total, nil
}

// PeringatanSpreadingTCO - `TreatyTestChildTotal_Act`: total != 100 -> pesan
// peringatan (baris SUDAH tersimpan; bukan penolakan).
func PeringatanSpreadingTCO(total *apd.Decimal) string {
	if total == nil || total.Cmp(seratusKlausul) != 0 {
		return PesanSpreadingTCO
	}
	return ""
}

// SubjenisExclusionTCO menurunkan subjenis baris ExclutionTreaty dari medannya.
func SubjenisExclusionTCO(k KlausulTreaty) string {
	switch {
	case strings.TrimSpace(k.IDOccupation) != "":
		return SubjenisOccupation
	case strings.TrimSpace(k.IDClause) != "":
		return SubjenisClause
	case strings.TrimSpace(k.Layer) != "":
		return SubjenisPeriode
	case k.Pct != nil:
		return SubjenisObject
	}
	return ""
}

// KlausulAnakTCO - baris anak adalah baris dengan induk (bukan sentinel "00").
func KlausulAnakTCO(k KlausulTreaty) bool {
	p := strings.TrimSpace(k.ParentReinsTypeID)
	return p != "" && p != ParentReinsTypeTanpaInduk
}
