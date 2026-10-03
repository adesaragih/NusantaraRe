package acceptance

// Tangga akseptasi bentuk B - lini financial (tiket NB-12). Jalur TERPISAH dari
// bentuk A: tabel, saringan, dan cara naiknya berbeda, dan keduanya tidak
// disatukan di balik satu abstraksi (tiket 12).
//
// Asal: langkah 12-14 dan 22 `D:\migrasi\RNM\NB FacIn\Activity\GetLimitAkseptasi_ActFlow.xml`
// (ASM-FW-GISFW-WORK!GETLIMITAKSEPTASI_ACTFLOW, pyRuleSetVersion 01-01-87).
// `[terverifikasi]` langkah 23 `D:\migrasi\RNM\NB FacIn\Activity\GetLimitAkseptasi_Act.xml`
// (ASM-FW-GISFW-WORK!GETLIMITAKSEPTASI_ACT) berisi blok yang sama dengan langkah 22
// (pembongkaran per langkah, 01-10-2026).
//
// ⚠️ `[keputusan work owner]` Mengoreksi K-026 dan tiket 12 (butir 46):
// `[terverifikasi]` SQL bentuk B tanpa ORDER BY / JABATAN_ATASAN / WORKBASKET,
// tetapi tangganya ditulis di ACTIVITY - langkah 22 naik satu tingkat menurut
// antrean saat ini, lalu Decision23 merutekannya.
//
// Tidak diport:
//   - langkah 17 (butir 47);
//   - langkah 10 dan 20 yang ikut berjalan bila IsNonPropertyandNonEngineering:
//     `[dugaan]` daftar langkah 10 diganti RDB-List langkah 12-14 (A21), dan
//     `[terverifikasi]` LetterNo langkah 20 ditimpa langkah 22.1 (`LetterNo = ""`);
//   - salinan `Data.LetterNo`, dan `pyWorkPage.NBStatus` (teks status berisi CARI5
//     = kolom NAMA, yang tidak diekspor).

import (
	"errors"
	"fmt"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/modul/nbfacin/backend/services/rules"
)

// BarisFinancial - satu baris M_LIMIT_FINANCIALINS. `[terverifikasi]` JABATAN-nya
// berspasi (`DIREKTUR TEKNIK`, fixture butir 42) dan dicocokkan persis, tanpa
// normalisasi. `[terverifikasi]` kolom LIMITTRADE_BOTTOM tidak dibaca rule mana pun:
// `grep -ril "LIMITTRADE_BOTTOM" "D:\migrasi\RNM\NB FacIn"` = 0 berkas.
type BarisFinancial struct {
	Jabatan        Jabatan
	LimitBond      *apd.Decimal
	LimitCreditCL  *apd.Decimal
	LimitCreditNCL *apd.Decimal
}

var (
	// ErrBukanBentukB - kasus bentuk A dibawa ke NextFinancial; jalurnya Next.
	ErrBukanBentukB = errors.New("acceptance: bukan kasus bentuk B, pakai Next")
	// ErrTafsirLimit - langkah 22.2.1 `Local.Limit = .CARID2` membaca kolom yang
	// tidak dipilih SQL (`[terverifikasi]` limit dipetakan ke CARI2), jadi
	// `[dugaan]` Limit kosong = 0. Bila tafsir itu dan tafsir "Limit = kolom
	// limit" memberi hasil berbeda, kasus ditolak (A26).
	ErrTafsirLimit = errors.New("acceptance: syarat TotalTSI >= Limit langkah 22 bergantung tafsir Local.Limit")
)

// daftarFinancial - satu RDB-List bentuk B: SQL, predikat pembukanya, dan kolom limitnya.
type daftarFinancial struct {
	sql      string   // RequestType = rule ASM-FW-GISFW-INT-POLICYJSON!ASM!<NAMA>
	predikat []string // langkah membuka bila salah satu benar
	limit    func(BarisFinancial) *apd.Decimal
}

// daftarBentukB - langkah 12-14. `[terverifikasi]` WHERE ketiga SQL
// (`D:\migrasi\RNM\NB FacIn\RDBList\<sql>.xml`, `<pyBrowseSQL>`):
// `<kolom> <= {DataSearch.CARID2}`, tanpa ORDER BY. Trade Credit memakai SQL
// Kredit Cash Loan (langkah 13 when2).
var daftarBentukB = []daftarFinancial{
	{"GetLimitAkseptasiBond_SQL", []string{"IsLimitSBondKBG", "IsLimitCustomBond"}, func(b BarisFinancial) *apd.Decimal { return b.LimitBond }},
	{"GetLimitAkseptasiKreditCL_SQL", []string{"IsLimitCreditCL", "IsLimitTradeCredit"}, func(b BarisFinancial) *apd.Decimal { return b.LimitCreditCL }},
	{"GetLimitAkseptasiKreditNCL_SQL", []string{"IsLimitCreditNCL"}, func(b BarisFinancial) *apd.Decimal { return b.LimitCreditNCL }},
}

// pilihDaftarFinancial - daftar bentuk B yang dimuat, atau nil untuk kasus bentuk A.
// `[terverifikasi]` kode BusinessOldId kelima predikat (`NB FacIn\When\IsLimit*.xml`)
// saling lepas - sensus literal `kanan` per predikat di `registry_gen.go`, irisan 0
// (01-10-2026) - tetapi lebih dari satu daftar tetap ditolak (butir 44 d).
func pilihDaftarFinancial(k rules.Kasus) (*daftarFinancial, error) {
	var terpilih []*daftarFinancial
	for i := range daftarBentukB {
		for _, pr := range daftarBentukB[i].predikat {
			buka, err := rules.Eval(pr, k)
			if err != nil {
				return nil, err
			}
			if buka {
				terpilih = append(terpilih, &daftarBentukB[i])
				break
			}
		}
	}
	switch len(terpilih) {
	case 0:
		return nil, nil
	case 1:
		return terpilih[0], nil
	}
	return nil, fmt.Errorf("%w: %s dan %s", ErrDaftarGanda, terpilih[0].sql, terpilih[1].sql)
}

// gerbangLangkah22 - `[terverifikasi]` 22.2.2-22.2.4: antrean saat ini, jabatan di
// daftar, dan LetterNo yang ditulis. Jabatan berspasi (tabel); LetterNo tanpa spasi
// (predikat To*). Tiap antrean membuka paling banyak satu gerbang.
var gerbangLangkah22 = map[Antrean]struct {
	jabatan Jabatan
	tujuan  Jabatan
}{
	antreanUnderwritingFinancial: {"KADIV KEUANGAN", "KADIVFINANCIAL"},
	antreanFinDivHead:            {"DIREKTUR MARKETING", "DIREKTURMARKETING"},
	antreanMarketingDirector:     {"DIREKTUR TEKNIK", "DIREKTURTEKNIK"},
}

// NextFinancial menjalankan SATU langkah tangga bentuk B untuk kasus `k`.
// `p.Jabatan` tidak dipakai: SQL bentuk B tanpa subquery pengguna. Anggota grup
// keluar di langkah 3, sebelum jenis daftarnya diperiksa - seperti di activity.
func NextFinancial(k rules.Kasus, tabel []BarisFinancial, p Pengguna) (Transisi, error) {
	tsi, keluar, err := langkah2Sampai5(k, p)
	if err != nil || keluar {
		return Transisi{Selesai: keluar}, err
	}
	daf, err := pilihDaftarFinancial(k)
	if err != nil {
		return Transisi{}, err
	}
	if daf == nil {
		return Transisi{}, ErrBukanBentukB
	}
	// 22.1 LetterNo = ""; 22.2 loop atas hasil SQL. Hanya gerbang antrean saat ini
	// yang dapat terbuka.
	var letterNo Jabatan
	g, ada := gerbangLangkah22[Antrean(teks(k, "pyWorkPage.PositionNote"))]
	for _, b := range tabel {
		limit := daf.limit(b)
		if limit == nil {
			return Transisi{}, fmt.Errorf("%w: %s (%s)", ErrLimitKosong, b.Jabatan, daf.sql)
		}
		if !ada || b.Jabatan != g.jabatan || limit.Cmp(tsi.carid2) > 0 { // WHERE <kolom> <= CARID2
			continue
		}
		// 22.2.2 when1 `Local.TotalTSI >= @toDecimal(Local.Limit)` dengan dua tafsir
		// Local.Limit: kosong = 0, atau kolom limit baris ini.
		lolosNol := tsi.totalTSI.Sign() >= 0
		lolosKolom := tsi.totalTSI.Cmp(limit) >= 0
		if lolosNol != lolosKolom {
			return Transisi{}, fmt.Errorf("%w: %s", ErrTafsirLimit, b.Jabatan)
		}
		if lolosNol {
			letterNo = g.tujuan
		}
	}
	return decision23(k, letterNo)
}
