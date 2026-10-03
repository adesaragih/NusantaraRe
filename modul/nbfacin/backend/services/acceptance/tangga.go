package acceptance

// Tangga akseptasi bentuk A (tiket NB-11): SATU langkah - dari jabatan pengguna
// yang baru memutuskan ke jabatan tujuan berikutnya, lalu ke antrean. Keputusan
// work owner yang mengikat: `docs/KEPUTUSAN-30-09-2026.md` butir 29, 33, 39-47.
// Bentuk B (financial): `tangga_financial.go` - jalur terpisah (butir 46).
//
// Asal (semua `D:\migrasi\RNM\NB FacIn\`):
//   - Activity\GetLimitAkseptasi_ActFlow.xml (ASM-FW-GISFW-WORK!GETLIMITAKSEPTASI_ACTFLOW,
//     pyRuleSetVersion 01-01-87) - langkah yang diport disebut per fungsi.
//   - RDBList\<nama>.xml, semuanya ASM-FW-GISFW-INT-POLICYJSON!ASM!<NAMA>: SQL tiap
//     daftar disebut di daftarBiasa / daftarBanding.
//   - Flow\InputInwardFacultativeOffer.xml (ASM-FW-GISFW-WORK!INPUTINWARDFACULTATIVEOFFER),
//     bentuk Decision23 - antrean.
//
// Tidak diport:
//   - langkah 1 `Call CountTotalTSIPremiNusaRe_Act`: menghitung ulang TotalTSINusaRe /
//     TotalTSITopRisk. Tanggung jawab PEMANGGIL - kasus wajib sudah berisi total
//     itu (butir 45);
//   - langkah 15, 16, 19, 21: label `//` = di-remark (butir 43);
//   - langkah 12-14 dan 22 (bentuk B, Bond/Kredit): NextFinancial, ditolak di sini;
//   - langkah 17: membuang "DIREKTUR TEKNIK" hanya untuk empat nomor polis literal -
//     dibuang di sistem baru, tidak dipakai lagi (butir 47);
//   - langkah 23-24 GetLimitAkseptasi_ActFlow: menyalin LetterNo, tidak mengubah keputusan.
//
// ⚠️ TSI dan limit dibandingkan sebagai desimal tanpa mata uang: tabel limit
// menyimpan MAX_LIMIT_IDR / MAX_LIMIT_USD terpisah, tetapi LIMIT_BOTTOM tidak
// bermata uang, dan mata uang TotalTSINusaRe `[pertanyaan terbuka]` (OQ-037,
// OQ-040, OQ-046). Karena itu bukan uang.Money.

import (
	"errors"
	"fmt"
	"sort"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/utils"
	"nusantarare/modul/nbfacin/backend/services/rules"
)

// Jabatan - kode jabatan: kolom JABATAN tabel limit, dan isi properti warisan
// `pyWorkPage.LetterNo` (yang TIDAK berisi nomor surat). Tipenya berbeda dari
// Antrean supaya menukar keduanya gagal saat kompilasi.
type Jabatan string

// Jabatan yang dibuang langkah 18.
const (
	jabatanManagerTeknik    Jabatan = "MANAGERTEKNIK"
	jabatanKadivFacultative Jabatan = "KADIVFACULTATIVE"
)

// Antrean - workbasket yang memegang kasus.
type Antrean string

// Antrean yang dibaca sebagai syarat (langkah 11, 22.2) dan ditulis Decision23.
const (
	antreanKadivFacultative      Antrean = "ReasFacInFacultativeDivHead"
	antreanUnderwritingFinancial Antrean = "ReasFacInUnderwritingFinancial"
	antreanFinDivHead            Antrean = "ReasFacInFinDivHead"
	antreanMarketingDirector     Antrean = "ReasFacInMarketingDirector"
	antreanTechnicalDirector     Antrean = "ReasFacInTechnicalDirector"
)

// Pengguna - pengguna yang baru memutuskan, dari model peran (butir 33):
// pengganti `LOGIN = {OperatorID.pyUserIdentifier}` di SQL dan pengganti daftar
// login literal predikat IsGroup (A14, A20).
type Pengguna struct {
	Jabatan Jabatan
	// AnggotaGrup - pengganti predikat IsGroup (langkah 3).
	AnggotaGrup bool
}

// NamaTabel - tabel limit bentuk A (skema POOLDATA).
type NamaTabel string

const (
	TabelProperty                    NamaTabel = "M_LIMIT_PROPERTYY"
	TabelPropertyNonPreferred        NamaTabel = "M_LIMIT_PROPERTY_NON_PREFERREDD"
	TabelPropertyPreferredCommercial NamaTabel = "M_LIMIT_PROPERTY_PREFERRED_COMMERCIALL"
	TabelEngineering                 NamaTabel = "M_LIMIT_ENGINEERINGG"
	TabelNonPropEng                  NamaTabel = "M_LIMIT_NONPROPANDENGG"
)

// BarisLimit - satu baris tabel limit, hanya kolom yang menentukan tangga.
// `MAX_LIMIT_*` dan `BATAS_WAKTU` (CARI2-CARI4) tidak dipakai keputusan ini.
type BarisLimit struct {
	Jabatan      Jabatan
	TeamGroup    string
	LimitBottom  *apd.Decimal
	LimitBottom2 *apd.Decimal
}

// TabelLimit - isi tabel limit, disuntikkan pemanggil (kelak pemuat repository).
type TabelLimit map[NamaTabel][]BarisLimit

// Transisi - hasil satu langkah tangga. Hasil keputusan underwriting (Keputusan,
// NB-13) bukan bagian dari sini: ketiga field state terpisah lewat tipe (butir 45).
type Transisi struct {
	// Selesai - Decision23 jatuh ke Else: tangga berakhir. Penyelesaian normal,
	// bukan galat. JabatanTujuan bisa tetap terisi (mis. SENIORUW tidak punya
	// konektor).
	Selesai bool
	// JabatanTujuan - LetterNo yang ditulis langkah 20 (bentuk A) atau 22.2 (bentuk B);
	// kosong bila tidak ada kandidat.
	JabatanTujuan Jabatan
	// Antrean - workbasket assignment tujuan (parameter `Workbasket` bentuk
	// Assignment*); kosong bila Selesai.
	Antrean Antrean
	// PositionNoteDitulis - konektor menulis `pyWorkPage.PositionNote = Antrean`.
	// ToKadivFacultative TIDAK menulisnya: kasus berpindah ke antrean Kadiv
	// Facultative, tetapi PositionNote tetap berisi nilai lama.
	PositionNoteDitulis bool
}

var (
	// ErrDaftarGanda - lebih dari satu daftar limit dimuat ke halaman yang sama
	// (bentuk A: IsFire dan IsEngineering sama-sama benar; bentuk B: lebih dari
	// satu predikat IsLimit*); daftar mana yang tersisa di Pega `[pertanyaan
	// terbuka]` (butir 44).
	ErrDaftarGanda = errors.New("acceptance: lebih dari satu tabel limit terpilih")
	// ErrTabelTakTerpilih - tidak ada daftar limit yang dimuat (mis. IsFire tanpa
	// IsPreferredRisk yang dikenal); isi halaman LimitAkseptasi `[pertanyaan
	// terbuka]` (butir 44).
	ErrTabelTakTerpilih = errors.New("acceptance: tidak ada tabel limit terpilih")
	// ErrUrutanTakPasti - bentuk A: kandidat pertama seri LIMIT_BOTTOM antar-jabatan berbeda.
	// `[terverifikasi]` ketujuh SQL bentuk A hanya `ORDER BY LIMIT_BOTTOM ASC`,
	// jadi urutan Oracle antar-baris seri tidak ditentukan (butir 44).
	ErrUrutanTakPasti = errors.New("acceptance: kandidat pertama seri, urutan tidak pasti")
	// ErrLimitKosong - kolom limit sebuah baris kosong. Apakah kolomnya nullable
	// `[pertanyaan terbuka]`; NULL di SQL akan menyaring baris itu diam-diam, jadi
	// ditolak di sini, tidak ditebak (A28).
	ErrLimitKosong = errors.New("acceptance: kolom limit kosong")
	// ErrBarisPenggunaGanda - jabatan pengguna punya lebih dari satu limit
	// berbeda di team group yang sama; subquery Pega akan gagal (ORA-01427).
	ErrBarisPenggunaGanda = errors.New("acceptance: limit jabatan pengguna ganda dan berbeda")
	// ErrBentukB - kasus Bond/Kredit/Trade (langkah 12-14, 22) dibawa ke Next;
	// jalurnya NextFinancial.
	ErrBentukB = errors.New("acceptance: kasus bentuk B (Bond/Kredit), pakai NextFinancial")
)

// Next menjalankan SATU langkah tangga bentuk A untuk kasus `k`.
func Next(k rules.Kasus, tabel TabelLimit, p Pengguna) (Transisi, error) {
	tsi, keluar, err := langkah2Sampai5(k, p)
	if err != nil || keluar {
		return Transisi{Selesai: keluar}, err
	}
	// Langkah 12-14 memuat daftar FINANCIALINS ke halaman yang sama.
	fin, err := pilihDaftarFinancial(k)
	if err != nil {
		return Transisi{}, err
	}
	if fin != nil {
		return Transisi{}, fmt.Errorf("%w: %s", ErrBentukB, fin.sql)
	}
	d, err := pilihDaftar(k)
	if err != nil {
		return Transisi{}, err
	}
	kandidat, err := d.saring(tabel[d.tabel], teks(k, "pyWorkPage.OfferFacIn.QuotationData.TeamGroup"), p.Jabatan, tsi.carid2)
	if err != nil {
		return Transisi{}, err
	}
	tujuan, err := kandidatPertama(buangJabatanLangkah18(kandidat, bandingReject(k)))
	if err != nil {
		return Transisi{}, err
	}
	return decision23(k, tujuan)
}

func teks(k rules.Kasus, jalur string) string {
	v, _ := k.Nilai(jalur)
	return v
}

// angka - properti angka; kosong atau tak terbaca = galat (tidak ditebak 0).
func angka(k rules.Kasus, jalur string) (*apd.Decimal, error) {
	d, err := utils.ParseDecimal(teks(k, jalur))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", jalur, err)
	}
	return d, nil
}

// bulatNol - `@Math.divide(x,1,0)`, setengah-ke-atas: mode yang sama dengan
// premium.bagiBulat. `[dugaan]` berdasar analogi fungsi `@Math.divide` yang sama -
// A37 dikonfirmasi untuk premi saja (butir 56); seri pada nilai dasar akseptasi
// belum pernah teramati. Kembaran premium.bulatkan; satu
// tempat kelak di inti (USULAN-PR-TIM-INTI-NB01).
func bulatNol(k rules.Kasus, jalur string) (*apd.Decimal, error) {
	x, err := angka(k, jalur)
	if err != nil {
		return nil, err
	}
	ctx := utils.DecimalContext()
	ctx.Rounding = apd.RoundHalfUp
	hasil := new(apd.Decimal)
	if _, err := ctx.Quantize(hasil, x, 0); err != nil {
		return nil, err
	}
	return hasil, nil
}

// langkah2Sampai5 - urutan korpus: langkah 2 (TSI dasar), langkah 3 "exit kalau
// group" (IsGroup → Exit Activity; LetterNo tetap "" dari langkah 2, Decision23
// jatuh ke Else), baru langkah 4-5 yang membaca top risk dan OldData. Dipakai
// kedua bentuk.
func langkah2Sampai5(k rules.Kasus, p Pengguna) (tsi nilaiTSI, keluar bool, err error) {
	dasar, err := bulatNol(k, "pyWorkPage.OfferFacIn.TotalTSINusaRe")
	if err != nil {
		return nilaiTSI{}, false, err
	}
	if p.AnggotaGrup {
		return nilaiTSI{}, true, nil
	}
	tsi, err = basisTSI(k, dasar)
	return tsi, false, err
}

// nilaiTSI - dua nilai TSI yang ditulis langkah 2, 4, 5. Keduanya sama, KECUALI
// endorsement dengan selisih top risk > 0.
type nilaiTSI struct {
	// carid2 - DataSearch.CARID2, nilai {DataSearch.CARID2} di SQL.
	carid2 *apd.Decimal
	// totalTSI - Local.TotalTSI, dibaca langkah 22 (bentuk B).
	totalTSI *apd.Decimal
}

// basisTSI menghitung kedua nilai dari `batas` hasil langkah 2:
//
//	langkah 2: keduanya @Math.divide(pyWorkPage.OfferFacIn.TotalTSINusaRe,1,0)
//	langkah 4: bila pyWorkPage.IsAdaTopRisk=="true" || pyWorkPage.OfferFacIn.TotalTSITopRisk>0
//	           → keduanya @Math.divide(pyWorkPage.OfferFacIn.TotalTSITopRisk,1,0)
//	langkah 5: bila StatusBusiness=="3" (Endorsement, butir 29) → TotalTSI = selisih
//	           mutlak TSI dengan OldData (`@if(x<0,x*-1,x)`); CARID2 = selisih mutlak
//	           top risk bila > 0, selain itu TotalTSI.
func basisTSI(k rules.Kasus, batas *apd.Decimal) (nilaiTSI, error) {
	var err error
	topRisk := teks(k, "pyWorkPage.IsAdaTopRisk") == "true"
	if !topRisk {
		tr, err := angka(k, "pyWorkPage.OfferFacIn.TotalTSITopRisk")
		if err != nil {
			return nilaiTSI{}, err
		}
		topRisk = tr.Sign() > 0
	}
	if topRisk {
		if batas, err = bulatNol(k, "pyWorkPage.OfferFacIn.TotalTSITopRisk"); err != nil {
			return nilaiTSI{}, err
		}
	}
	if teks(k, "pyWorkPage.OfferFacIn.QuotationData.StatusBusiness") != "3" {
		return nilaiTSI{carid2: batas, totalTSI: batas}, nil
	}
	selisih := func(baru, lama string) (*apd.Decimal, error) {
		b, err := bulatNol(k, baru)
		if err != nil {
			return nil, err
		}
		l, err := bulatNol(k, lama)
		if err != nil {
			return nil, err
		}
		d := new(apd.Decimal)
		if _, err := utils.DecimalContext().Sub(d, b, l); err != nil {
			return nil, err
		}
		return d.Abs(d), nil
	}
	tsi, err := selisih("pyWorkPage.OfferFacIn.TotalTSINusaRe", "pyWorkPage.OfferFacIn.OldData.TotalTSINusaRe")
	if err != nil {
		return nilaiTSI{}, err
	}
	diffTop, err := selisih("pyWorkPage.OfferFacIn.TotalTSITopRisk", "pyWorkPage.OfferFacIn.OldData.TotalTSITopRisk")
	if err != nil {
		return nilaiTSI{}, err
	}
	if diffTop.Sign() > 0 {
		return nilaiTSI{carid2: diffTop, totalTSI: tsi}, nil
	}
	return nilaiTSI{carid2: tsi, totalTSI: tsi}, nil
}

// bandingReject - `pyWorkPage.OfferFacIn.IsBanding=="true" && pyWorkPage.OfferFacIn.IsFlagReject=="true"`
// (langkah 11 dan 18.1).
func bandingReject(k rules.Kasus) bool {
	return teks(k, "pyWorkPage.OfferFacIn.IsBanding") == "true" && teks(k, "pyWorkPage.OfferFacIn.IsFlagReject") == "true"
}

// langkahDaftar - satu RDB-List: tabel, SQL-nya, dan prakondisinya (semua When
// harus benar).
type langkahDaftar struct {
	tabel     NamaTabel
	sql       string   // RequestType = nama rule RDBList
	predikat  string   // When rule
	preferred []string // nilai pyWorkPage.OfferFacIn.IsPreferredRisk yang membuka; nil = tidak dicek
}

// daftarBiasa - langkah 6-10.
var daftarBiasa = []langkahDaftar{
	{TabelProperty, "GetLimitAkseptasi_SQL", "IsFire", []string{"Preferred Risk"}},
	{TabelPropertyNonPreferred, "GetLimitAkseptasiNonPrefer_SQL", "IsFire", []string{"Non-Preferred Risk"}},
	{TabelPropertyPreferredCommercial, "GetLimitAkseptasiPreferedComm_SQL", "IsFire", []string{"Preferred Risk Commercial"}},
	{TabelEngineering, "GetLimitAccEngineeringUW_SQL", "IsEngineering", nil},
	{TabelNonPropEng, "GetLimitAkseptasiNonFire_SQL", "IsNonPropertyandNonEngineering", nil},
}

// daftarBanding - langkah 11.1-11.4. Preferred Risk Commercial memakai tabel
// PROPERTYY di sini, sesuai 11.1.
var daftarBanding = []langkahDaftar{
	{TabelProperty, "GetLimitAkseptasiBanding_SQL", "IsFire", []string{"Preferred Risk", "Preferred Risk Commercial"}},
	{TabelPropertyNonPreferred, "GetLimitAkseptasiNonPreferBanding_SQL", "IsFire", []string{"Non-Preferred Risk"}},
	{TabelEngineering, "GetLimitAccEngineeringBanding_SQL", "IsEngineering", nil},
	{TabelNonPropEng, "GetLimitAkseptasiNonFireBanding_SQL", "IsNonPropertyandNonEngineering", nil},
}

// daftar - daftar limit yang dimuat: tabelnya dan varian SQL-nya.
type daftar struct {
	tabel   NamaTabel
	banding bool // varian `…Banding_SQL`: saring LIMIT_BOTTOM2
}

// pilihDaftar - langkah 6-11. Langkah 11 (banding + reject di antrean Kadiv
// Facultative) memuat ulang daftar dengan varian Banding; `[dugaan]` RDB-List
// mengganti isi halaman (A21). Lebih dari satu atau tidak ada daftar dalam satu
// kelompok → galat (butir 44).
func pilihDaftar(k rules.Kasus) (daftar, error) {
	langkah, banding := daftarBiasa, false
	if bandingReject(k) && Antrean(teks(k, "pyWorkPage.PositionNote")) == antreanKadivFacultative {
		langkah, banding = daftarBanding, true
	}
	var terpilih []NamaTabel
	for _, l := range langkah {
		buka, err := rules.Eval(l.predikat, k)
		if err != nil {
			return daftar{}, err
		}
		if buka && l.preferred != nil {
			buka = false
			for _, v := range l.preferred {
				buka = buka || teks(k, "pyWorkPage.OfferFacIn.IsPreferredRisk") == v
			}
		}
		if buka {
			terpilih = append(terpilih, l.tabel)
		}
	}
	switch len(terpilih) {
	case 0:
		return daftar{}, ErrTabelTakTerpilih
	case 1:
		return daftar{terpilih[0], banding}, nil
	}
	return daftar{}, fmt.Errorf("%w: %v", ErrDaftarGanda, terpilih)
}

// limit - kolom yang disaring: LIMIT_BOTTOM, atau LIMIT_BOTTOM2 untuk varian Banding.
func (d daftar) limit(b BarisLimit) *apd.Decimal {
	if d.banding {
		return b.LimitBottom2
	}
	return b.LimitBottom
}

// saring - `[terverifikasi]` WHERE ketujuh SQL bentuk A (dan keempat varian Banding):
//
//	team_group = {TeamGroup} AND LIMIT_BOTTOM > (limit pengguna) AND LIMIT_BOTTOM < {CARID2}
//	ORDER BY LIMIT_BOTTOM ASC
//
// Varian Banding menyaring LIMIT_BOTTOM2 tetapi TETAP mengurutkan LIMIT_BOTTOM.
// Jabatan pengguna yang tidak ada → subquery kosong → tidak ada kandidat (butir 44).
// `[dugaan]` team_group dibandingkan sebagai teks: kolom VARCHAR2(20)
// (DDL\M_LIMIT_PROPERTYY.txt), dan {TeamGroup} disisipkan tanpa kutip (A24).
func (d daftar) saring(baris []BarisLimit, teamGroup string, j Jabatan, batas *apd.Decimal) ([]BarisLimit, error) {
	var milikPengguna *apd.Decimal
	for _, b := range baris {
		if b.TeamGroup == teamGroup && (d.limit(b) == nil || b.LimitBottom == nil) {
			return nil, fmt.Errorf("%w: %s team group %s", ErrLimitKosong, b.Jabatan, teamGroup)
		}
	}
	for _, b := range baris {
		if b.TeamGroup != teamGroup || b.Jabatan != j {
			continue
		}
		if milikPengguna != nil && milikPengguna.Cmp(d.limit(b)) != 0 {
			return nil, fmt.Errorf("%w: %s team group %s", ErrBarisPenggunaGanda, j, teamGroup)
		}
		milikPengguna = d.limit(b)
	}
	if milikPengguna == nil {
		return nil, nil
	}
	var hasil []BarisLimit
	for _, b := range baris {
		if b.TeamGroup == teamGroup && d.limit(b).Cmp(milikPengguna) > 0 && d.limit(b).Cmp(batas) < 0 {
			hasil = append(hasil, b)
		}
	}
	sort.SliceStable(hasil, func(i, j int) bool { return hasil[i].LimitBottom.Cmp(hasil[j].LimitBottom) < 0 })
	return hasil, nil
}

// buangJabatanLangkah18 - langkah 18 "hapus": 18.1 KADIVFACULTATIVE kecuali
// banding + reject; 18.2 MANAGERTEKNIK selalu. Langkah 18 ber-pyStepsPreCondition=false
// = tetap jalan (P-11). `[dugaan]` Property-Remove di dalam loop membuang SEMUA
// baris yang cocok (A23).
func buangJabatanLangkah18(baris []BarisLimit, bandingReject bool) []BarisLimit {
	var sisa []BarisLimit
	for _, b := range baris {
		if b.Jabatan == jabatanManagerTeknik || (b.Jabatan == jabatanKadivFacultative && !bandingReject) {
			continue
		}
		sisa = append(sisa, b)
	}
	return sisa
}

// kandidatPertama - langkah 20: `pyWorkPage.LetterNo = LimitAkseptasi.pxResults(1).CARI1`.
// Daftar kosong → "" `[dugaan]` (A25). Seri antar-jabatan berbeda → galat.
func kandidatPertama(baris []BarisLimit) (Jabatan, error) {
	if len(baris) == 0 {
		return "", nil
	}
	for _, b := range baris[1:] {
		if b.LimitBottom.Cmp(baris[0].LimitBottom) == 0 && b.Jabatan != baris[0].Jabatan {
			return "", fmt.Errorf("%w: %s dan %s", ErrUrutanTakPasti, baris[0].Jabatan, b.Jabatan)
		}
	}
	return baris[0].Jabatan, nil
}

// konektorDecision23 - bentuk Decision23 (XOR). `[terverifikasi]` workbasket tiap
// assignment tujuan (parameter `Workbasket`) dan PositionNote tiap konektor, dibaca
// dari `pyToTasks/…/pyPropertySet` dan `pyConnectors/…/pyPropertyAssigns` yang
// sepakat; pada enam konektor yang menulis PositionNote, nilainya = workbasket.
var konektorDecision23 = []struct {
	predikat          string
	workbasket        Antrean
	tulisPositionNote bool
}{
	{"ToKadivTeknik", "ReasFacInGroupLeader", true},        // Transition30 → Assignment1
	{"ToDepHeadUW", "ReasFacInDepHeadUnderwriting", true},  // Transition133 → Assignment17
	{"ToKadivFacultative", antreanKadivFacultative, false}, // Transition134 → Assignment8
	{"ToKadivFin", antreanFinDivHead, true},                // Transition94 → Assignment11
	{"ToDirMarketing", antreanMarketingDirector, true},     // Transition13 → Assignment6
	{"ToDirTeknik", antreanTechnicalDirector, true},        // Transition19 → Assignment4
	{"ToManagerTeknik", "ReasFacInManagerTeknik", true},    // Transition75 → Assignment5
}

// kasusLetterNo - kasus dengan pyWorkPage.LetterNo hasil langkah 20.
type kasusLetterNo struct {
	rules.Kasus
	letterNo Jabatan
}

func (k kasusLetterNo) Nilai(jalur string) (string, bool) {
	if jalur == "pyWorkPage.LetterNo" {
		return string(k.letterNo), true
	}
	return k.Kasus.Nilai(jalur)
}

// decision23 - konektor yang predikatnya benar; tidak ada = Else (Transition35,
// "finish") = tangga selesai.
func decision23(k rules.Kasus, tujuan Jabatan) (Transisi, error) {
	kl := kasusLetterNo{k, tujuan}
	var hasil []Transisi
	for _, c := range konektorDecision23 {
		buka, err := rules.Eval(c.predikat, kl)
		if err != nil {
			return Transisi{}, err
		}
		if buka {
			hasil = append(hasil, Transisi{JabatanTujuan: tujuan, Antrean: c.workbasket, PositionNoteDitulis: c.tulisPositionNote})
		}
	}
	switch len(hasil) {
	case 0:
		return Transisi{Selesai: true, JabatanTujuan: tujuan}, nil
	case 1:
		return hasil[0], nil
	}
	// Predikat To* membandingkan LetterNo dengan jabatan berbeda: dua benar
	// berarti registry berubah - kesalahan program. Bukan KondisiTakDikenali:
	// tangga tidak punya Fase.
	panic(fmt.Sprintf("acceptance: lebih dari satu konektor Decision23 benar untuk %q", tujuan))
}
