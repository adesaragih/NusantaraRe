package models

// Hitung kolom peserta Type QR - Calculate CSV (keputusan work owner 05-10-2026).
//
// Untuk apa berkas ini: aturan MURNI (tanpa Oracle) yang mengisi kolom hasil
// peserta polis Type "QR" dari SUM_INSURED, parameter produk Master Product
// Name Life, master rate (`RATE_LIFE`) dan master risk (`RIRISK_LIFE`).
// Type lain TIDAK melewati berkas ini - nilainya tetap dari CSV.
//
// Acuan: `SetRateLIfePremium_Act` langkah 1-10.4 (Find RIRATE / RIRISK, set
// Rate, set Sum at Risk), DITULIS ULANG menurut spesifikasi work owner
// 05-10-2026 - spesifikasi yang berlaku bila berbeda dari XML:
//
//   - AGE rate = ENTRY_AGE (CURRENT_AGE tidak dipakai);
//   - YEAR risk = tahun(GROSS_VALUATION_BEGIN_DATE) - tahun(BEGIN_DATE), TANPA +1;
//   - EM_PERCENT berupa PECAHAN (0.5 = 50%);
//   - rate / risk tidak ketemu atau GANDA DITOLAK (Pega diam-diam 0 / memakai
//     nilai peserta sebelumnya);
//   - master rate campuran unisex + gender DITOLAK; Ceding's Limit WAJIB.
//
// Spreading retro (langkah 10.5 dst.) dan kolom *_REFUND TIDAK dikerjakan.
//
// ⛔ NOL FLOAT (ADR-U-0003). Seluruh hitungan apd.Decimal presisi penuh; nilai
// antara TIDAK dibulatkan. Kolom hasil dibulatkan half-up 4 desimal saat
// ditulis ke baris (RISK 7 desimal).

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/utils"
)

// TipeHitungQR - satu-satunya Type yang kolom pesertanya dihitung.
const TipeHitungQR = "QR"

// HitungPesertaPerTipe menjawab apakah Calculate CSV menghitung kolom peserta
// untuk Type ini.
func HitungPesertaPerTipe(tipe string) bool { return strings.TrimSpace(tipe) == TipeHitungQR }

// KolomHasilHitungQR - kolom CSV yang untuk Type QR DIHITUNG; nilainya di CSV
// diabaikan dan ditimpa (bentuknya pun tidak divalidasi).
var KolomHasilHitungQR = []string{
	"CEDING_RETENTION", "SUM_REASURED", KolomShareNusantaraRe, KolomShareNusantaraReGross,
	"SUM_AT_RISK_GROSS", "RISK", "RATE", "GROSS_PREMIUM", "DEDUCTION", "BROKERAGE_FEE", "NET_PREMIUM",
}

// kolomHasilQR - himpunan KolomHasilHitungQR.
var kolomHasilQR = func() map[string]bool {
	m := map[string]bool{}
	for _, k := range KolomHasilHitungQR {
		m[k] = true
	}
	return m
}()

// ErrSyaratHitungQR - polis/produk/master tidak memenuhi syarat hitung QR
// (P1-P4). Seluruh unggahan ditolak, nol tulisan.
var ErrSyaratHitungQR = errors.New("models: syarat hitung Type QR tidak terpenuhi")

// ParamProdukQR - parameter produk `M_PRODUCTNAME_LIFE` (teks desimal dari
// Oracle) dan polis yang dipakai hitungan.
type ParamProdukQR struct {
	ProductName        string
	CedingRetentionNum string // Ceding Retention (%)
	CedingLimit        string // Ceding's Limit
	RNMShare           string // Nusantara Re Share (%)
	RIComm             string // Deduction (%)
	Brokerage          string // Brokerage Fee (%)
	RIRiskID           string
	// BusinessName - Class of Business polis (`T_PREMIUM_LIST.BUSINESS_NAME`).
	BusinessName string
	// ProRateType - Premium Payment Method polis: 1 Single, 2 Annually, 3 Others.
	ProRateType string
}

// PlanProdukQR - satu baris `M_PRODUCTNAME_LIFE_PLAN`.
type PlanProdukQR struct {
	Name     string // kolom "Bussines"
	RIRateID string
}

// BarisRateQR - satu baris view `RATE_LIFE` (semua kolom teks di view).
type BarisRateQR struct{ ID, Gender, Age, Contract, Rate string }

// BarisRiskQR - satu baris view `RIRISK_LIFE`.
type BarisRiskQR struct{ ID, Year, Contract, Risk string }

// rateQR / riskQR - baris master yang sudah diurai. nil = sel kosong.
type rateQR struct {
	id            string
	gender        string
	age, contract *apd.Decimal
	rate          *apd.Decimal
}

type riskQR struct {
	id             string
	year, contract *apd.Decimal
	risk           *apd.Decimal
}

// MasterQR - parameter dan master yang siap dipakai HitungPesertaQR.
type MasterQR struct {
	cedingRet, cedingLimit, rnmShare, riComm, brokerage *apd.Decimal
	riRateID, riRiskID                                  string
	proRateType                                         string
	// gender - master rate ber-GENDER F/M; false = seluruhnya U (unisex).
	gender bool
	rate   []rateQR
	risk   []riskQR
}

// ctxQR - konteks aritmetika: presisi lebar, half-up.
var ctxQR = func() *apd.Context {
	c := apd.BaseContext.WithPrecision(60)
	c.Rounding = apd.RoundHalfUp
	return c
}()

var (
	seratus  = apd.New(100, 0)
	seribu   = apd.New(1000, 0)
	satu     = apd.New(1, 0)
	nol      = apd.New(0, 0)
	duabelas = apd.New(12, 0)
)

// angkaMaster mengurai teks angka master/produk: SATU koma desimal menjadi
// titik (`DesimalKomaKeTitik`); kosong → nil.
func angkaMaster(teks string) (*apd.Decimal, error) {
	s := strings.TrimSpace(DesimalKomaKeTitik(teks))
	if s == "" {
		return nil, nil
	}
	return utils.ParseDecimal(s)
}

// PlanCocokQR - P2: tepat SATU plan produk ber-NAME (kolom "Bussines") =
// Class of Business polis (TrimSpace, tanpa peka huruf besar); P4: RIRATEID
// plan itu dan RIRISKID produk terisi. Dipakai pemanggil untuk tahu master
// rate mana yang dibaca, dan oleh SiapkanMasterQR.
func PlanCocokQR(p ParamProdukQR, plan []PlanProdukQR) (PlanProdukQR, error) {
	gagal := func(format string, a ...any) (PlanProdukQR, error) {
		return PlanProdukQR{}, fmt.Errorf("%w: %s", ErrSyaratHitungQR, fmt.Sprintf(format, a...))
	}
	bisnis := strings.TrimSpace(p.BusinessName)
	var cocok []PlanProdukQR
	for _, pl := range plan {
		if bisnis != "" && strings.EqualFold(strings.TrimSpace(pl.Name), bisnis) {
			cocok = append(cocok, pl)
		}
	}
	switch len(cocok) {
	case 0:
		return gagal("No plan in product %s matches Class of Business %s", p.ProductName, bisnis)
	case 1:
	default:
		return gagal("More than one plan in product %s matches Class of Business %s", p.ProductName, bisnis)
	}
	if strings.TrimSpace(cocok[0].RIRateID) == "" {
		return gagal("Plan %s in product %s has no R/I Rate", bisnis, p.ProductName)
	}
	if strings.TrimSpace(p.RIRiskID) == "" {
		return gagal("Product %s has no R/I Risk", p.ProductName)
	}
	return cocok[0], nil
}

// SiapkanMasterQR menegakkan syarat polis P1-P4 dan mengurai master rate /
// risk. Galatnya ErrSyaratHitungQR berpesan yang dapat dibaca pemakai.
func SiapkanMasterQR(p ParamProdukQR, plan []PlanProdukQR, rate []BarisRateQR, risk []BarisRiskQR) (MasterQR, error) {
	gagal := func(format string, a ...any) (MasterQR, error) {
		return MasterQR{}, fmt.Errorf("%w: %s", ErrSyaratHitungQR, fmt.Sprintf(format, a...))
	}
	m := MasterQR{proRateType: strings.TrimSpace(p.ProRateType)}
	persen := []struct {
		teks, label string
		ke          **apd.Decimal
	}{
		{p.CedingRetentionNum, "Ceding Retention (%)", &m.cedingRet},
		{p.RNMShare, "Nusantara Re Share (%)", &m.rnmShare},
		{p.RIComm, "Deduction (%)", &m.riComm},
		{p.Brokerage, "Brokerage Fee (%)", &m.brokerage},
	}
	for _, x := range persen {
		d, err := angkaMaster(x.teks)
		if err != nil {
			return gagal("%s %q on product %s is not a number", x.label, x.teks, p.ProductName)
		}
		if d == nil {
			d = nol // kosong = 0 (kolom produk nullable)
		}
		*x.ke = d
	}
	// P1 - Ceding's Limit terisi dan > 0.
	limit, err := angkaMaster(p.CedingLimit)
	if err != nil || limit == nil || limit.Sign() <= 0 {
		return gagal("Ceding's Limit is not set on the product")
	}
	m.cedingLimit = limit
	// P2 + P4.
	pl, err := PlanCocokQR(p, plan)
	if err != nil {
		return MasterQR{}, err
	}
	m.riRateID = strings.TrimSpace(pl.RIRateID)
	m.riRiskID = strings.TrimSpace(p.RIRiskID)
	// Master rate: urai, lalu P3 - seragam U atau seragam F/M.
	var u, fm int
	for _, r := range rate {
		g := strings.ToUpper(strings.TrimSpace(r.Gender))
		baris := rateQR{id: strings.TrimSpace(r.ID), gender: g}
		for _, x := range []struct {
			teks, kolom string
			ke          **apd.Decimal
		}{{r.Age, "AGE", &baris.age}, {r.Contract, "CONTRACT", &baris.contract}, {r.Rate, "RATE", &baris.rate}} {
			d, err := angkaMaster(x.teks)
			if err != nil {
				return gagal("R/I Rate %s row ID %s: %s %q is not a number", m.riRateID, baris.id, x.kolom, x.teks)
			}
			*x.ke = d
		}
		switch g {
		case "U":
			u++
		case "F", "M":
			fm++
		default:
			return gagal("R/I Rate %s row ID %s: GENDER %q is not U, F or M", m.riRateID, baris.id, r.Gender)
		}
		m.rate = append(m.rate, baris)
	}
	if len(m.rate) == 0 {
		return gagal("R/I Rate %s has no rows", m.riRateID)
	}
	if u > 0 && fm > 0 {
		return gagal("R/I Rate %s mixes unisex and gendered rows", m.riRateID)
	}
	m.gender = fm > 0
	if m.gender && m.proRateType != "1" && m.proRateType != "2" && m.proRateType != "3" {
		return gagal("Premium Payment Method (1, 2 or 3) must be set on the policy for gendered R/I Rate %s", m.riRateID)
	}
	for _, r := range risk {
		baris := riskQR{id: strings.TrimSpace(r.ID)}
		for _, x := range []struct {
			teks, kolom string
			ke          **apd.Decimal
		}{{r.Year, "YEAR", &baris.year}, {r.Contract, "CONTRACT", &baris.contract}, {r.Risk, "RISK", &baris.risk}} {
			d, err := angkaMaster(x.teks)
			if err != nil {
				return gagal("R/I Risk %s row ID %s: %s %q is not a number", m.riRiskID, baris.id, x.kolom, x.teks)
			}
			*x.ke = d
		}
		m.risk = append(m.risk, baris)
	}
	return m, nil
}

// sama - nilai master terisi dan sama dengan v.
func sama(master, v *apd.Decimal) bool { return master != nil && master.Cmp(v) == 0 }

// bulatTeks - half-up `digit` desimal, teks tanpa notasi eksponen.
func bulatTeks(d *apd.Decimal, digit int32) string {
	var r apd.Decimal
	_, _ = ctxQR.Quantize(&r, d, -digit)
	return utils.FormatDecimal(&r)
}

// KontrakQR - CONTRACT dari PERIOD_MM: 1 bila < 12, selain itu round(PERIOD_MM/12) half-up.
func KontrakQR(periodMM int64) *apd.Decimal {
	if periodMM < 12 {
		return apd.New(1, 0)
	}
	var q, r apd.Decimal
	_, _ = ctxQR.Quo(&q, apd.New(periodMM, 0), duabelas)
	_, _ = ctxQR.Quantize(&r, &q, 0)
	return &r
}

// HitungPesertaQR menghitung kolom hasil setiap baris dan MENULISNYA ke
// `baris[i].Nilai`. Baris yang ditolak tidak ditulis. Mengembalikan seluruh
// penolakan (nomor baris CSV, kolom, nilai yang dicari).
//
// ⚠️ Baris yang sudah ditolak validasi (`lewati`) tidak dihitung - penolakan
// bentuknya sudah cukup, dan menghitung nilai rusak hanya menambah kalimat.
func HitungPesertaQR(m MasterQR, baris []BarisUnggah, lewati map[int]bool) []Penolakan {
	tolak := []Penolakan{}
	for i := range baris {
		b := &baris[i]
		if lewati[b.Nomor] {
			continue
		}
		if p := hitungSatuQR(m, b); p != nil {
			tolak = append(tolak, *p)
		}
	}
	return tolak
}

func hitungSatuQR(m MasterQR, b *BarisUnggah) *Penolakan {
	ambil := func(k string) string { return strings.TrimSpace(b.Nilai[k]) }
	tolak := func(kolom, pesan string) *Penolakan {
		return &Penolakan{Baris: b.Nomor, Kolom: kolom, Pesan: pesan, Sebab: pesan}
	}
	// Masukan.
	si, err := UangCSV(ambil("SUM_INSURED"))
	if err != nil {
		return tolak("SUM_INSURED", "SUM_INSURED is not a valid number")
	}
	umur, err := strconv.ParseInt(ambil("ENTRY_AGE"), 10, 64)
	if err != nil {
		return tolak("ENTRY_AGE", fmt.Sprintf("ENTRY_AGE %q is not a whole number", ambil("ENTRY_AGE")))
	}
	if ambil("GROSS_VALUATION_BEGIN_DATE") == "" {
		return tolak("GROSS_VALUATION_BEGIN_DATE", PesanGrossValMulai)
	}
	gvb, err := TanggalCSV(ambil("GROSS_VALUATION_BEGIN_DATE"))
	if err != nil {
		return tolak("GROSS_VALUATION_BEGIN_DATE", PesanGrossValMulai)
	}
	mulai, err := TanggalCSV(ambil("BEGIN_DATE"))
	if err != nil {
		return tolak("BEGIN_DATE", PesanBeginDate)
	}
	periode, err := strconv.ParseInt(ambil("PERIOD_MM"), 10, 64)
	if err != nil || periode < 1 {
		return tolak("PERIOD_MM", fmt.Sprintf("PERIOD_MM %q must be a whole number of at least 1", ambil("PERIOD_MM")))
	}
	kontrak := KontrakQR(periode)
	sex := strings.ToUpper(ambil("SEX"))
	if m.gender && sex != "F" && sex != "M" {
		return tolak("SEX", fmt.Sprintf("SEX must be F or M (R/I Rate %s is gendered), found %q", m.riRateID, ambil("SEX")))
	}
	faktor := satu
	if v := ambil("FACTOR"); v != "" {
		d, err := UangCSV(v)
		if err != nil {
			return tolak("FACTOR", fmt.Sprintf("FACTOR %q is not a valid number", v))
		}
		if d.Sign() != 0 {
			faktor = d
		}
	}
	em := nol
	if v := ambil("EM_PERCENT"); v != "" {
		d, err := UangCSV(v)
		if err != nil {
			return tolak("EM_PERCENT", fmt.Sprintf("EM_PERCENT %q is not a valid number", v))
		}
		em = d
	}

	// 10 - RATE dari master plan, AGE = ENTRY_AGE.
	age := apd.New(umur, 0)
	var rate *apd.Decimal
	var cocok int
	for _, r := range m.rate {
		ok := sama(r.age, age)
		switch {
		case !m.gender:
			ok = ok && sama(r.contract, kontrak)
		case m.proRateType == "1":
			ok = ok && sama(r.contract, kontrak) && r.gender == sex
		default: // Pro Rate Type 2/3: CONTRACT tidak dicocokkan
			ok = ok && r.gender == sex
		}
		if ok && r.rate != nil {
			cocok++
			rate = r.rate
		}
	}
	dicari := fmt.Sprintf("age %d", umur)
	if !m.gender || m.proRateType == "1" {
		dicari += ", contract " + kontrak.Text('f')
	}
	if m.gender {
		dicari += ", sex " + sex
	}
	if cocok == 0 {
		return tolak("RATE", fmt.Sprintf("No rate found for %s in R/I Rate %s", dicari, m.riRateID))
	}
	if cocok > 1 {
		return tolak("RATE", fmt.Sprintf("More than one rate found for %s in R/I Rate %s", dicari, m.riRateID))
	}

	// 5-7.
	var cr, sr, share apd.Decimal
	_, _ = ctxQR.Mul(&cr, si, m.cedingRet)
	_, _ = ctxQR.Quo(&cr, &cr, seratus)
	if cr.Cmp(m.cedingLimit) > 0 {
		cr.Set(m.cedingLimit)
	}
	_, _ = ctxQR.Sub(&sr, si, &cr)
	_, _ = ctxQR.Mul(&share, &sr, m.rnmShare)
	_, _ = ctxQR.Quo(&share, &share, seratus)

	// 8-9 - tahun pertama ⇔ BEGIN_DATE = GROSS_VALUATION_BEGIN_DATE.
	var sar apd.Decimal
	riskTeks := ""
	if tanggalSama(mulai, gvb) {
		sar.Set(&share)
	} else {
		tahun := apd.New(int64(gvb.Year()-mulai.Year()), 0)
		var risk *apd.Decimal
		var n int
		for _, r := range m.risk {
			if sama(r.year, tahun) && sama(r.contract, kontrak) && r.risk != nil {
				n++
				risk = r.risk
			}
		}
		cari := fmt.Sprintf("year %s, contract %s", tahun.Text('f'), kontrak.Text('f'))
		if n == 0 {
			return tolak("RISK", fmt.Sprintf("No risk found for %s in R/I Risk %s", cari, m.riRiskID))
		}
		if n > 1 {
			return tolak("RISK", fmt.Sprintf("More than one risk found for %s in R/I Risk %s", cari, m.riRiskID))
		}
		var riskKolom apd.Decimal
		_, _ = ctxQR.Quo(&riskKolom, risk, seribu)
		riskTeks = bulatTeks(&riskKolom, 7)
		_, _ = ctxQR.Mul(&sar, &share, risk)
		_, _ = ctxQR.Quo(&sar, &sar, seribu)
	}

	// 11-14.
	var gross, emSatu, ded, brk, net apd.Decimal
	_, _ = ctxQR.Quo(&gross, rate, seribu)
	_, _ = ctxQR.Mul(&gross, &gross, &share)
	_, _ = ctxQR.Mul(&gross, &gross, faktor)
	_, _ = ctxQR.Add(&emSatu, em, satu)
	_, _ = ctxQR.Mul(&gross, &gross, &emSatu)
	_, _ = ctxQR.Mul(&ded, &gross, m.riComm)
	_, _ = ctxQR.Quo(&ded, &ded, seratus)
	_, _ = ctxQR.Mul(&brk, &gross, m.brokerage)
	_, _ = ctxQR.Quo(&brk, &brk, seratus)
	_, _ = ctxQR.Sub(&net, &gross, &ded)
	_, _ = ctxQR.Sub(&net, &net, &brk)

	// Tulis - kolom hasil CSV ditimpa.
	b.Nilai["CEDING_RETENTION"] = bulatTeks(&cr, 4)
	b.Nilai["SUM_REASURED"] = bulatTeks(&sr, 4)
	b.Nilai[KolomShareNusantaraRe] = bulatTeks(&share, 4)
	b.Nilai[KolomShareNusantaraReGross] = bulatTeks(&share, 4)
	b.Nilai["SUM_AT_RISK_GROSS"] = bulatTeks(&sar, 4)
	b.Nilai["RISK"] = riskTeks
	b.Nilai["RATE"] = bulatTeks(rate, 4)
	b.Nilai["GROSS_PREMIUM"] = bulatTeks(&gross, 4)
	b.Nilai["DEDUCTION"] = bulatTeks(&ded, 4)
	b.Nilai["BROKERAGE_FEE"] = bulatTeks(&brk, 4)
	b.Nilai["NET_PREMIUM"] = bulatTeks(&net, 4)
	// RI_ADMIN_FEE: nilai CSV bila terisi, selain itu 0 (tidak mengurangi NET_PREMIUM).
	if ambil("RI_ADMIN_FEE") == "" {
		b.Nilai["RI_ADMIN_FEE"] = "0"
	}
	return nil
}

// tanggalSama - hari, bulan, tahun sama persis.
func tanggalSama(a, b time.Time) bool {
	return a.Year() == b.Year() && a.Month() == b.Month() && a.Day() == b.Day()
}
