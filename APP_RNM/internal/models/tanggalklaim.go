package models

// Tiga tanggal klaim dialog Edit Date - sensus Claim Life §3.1.
//
// Untuk apa berkas ini: gerbang tahap dialog `Edit Date` selain DOL. DOL
// sendiri sudah ditulis tiket 06 (`services/dol.go` Set).
//
// `[terverifikasi]` `Claim Life/Section/EditDateClaimLife_Section.xml`
// (pecahan `sed 's/></>\n</g'`):
//
//	DATE_OF_LOSS         b796   change -> ValidasiDOL_Act b837 (tiket 06)
//	CLAIM_RECEIVED_DATE  b1076  change -> ValidasiClaimReceived_Act b1120
//	COMPLETE_DATE        b1387  caption "DOCUMENT COMPLETE DATE"; postValue saja
//	CONFIRMATION_DATE    b1626  postValue saja
//	Save b1910 -> UpdateDateClaimLife_Act b1929 -> closeContainer -> refresh
//
// Keempat isian berprasyarat baca-saja YANG SAMA (b1000, b1313, b1550, b1788):
//
//	pyWorkPage.pyPosition!='ReasLifeAdmin' || pyWorkPage.ClaimData.PremiumListSummary.CLAIM_NO!=''
//
// ⚠️ `[terbuka - OQ-M1]` separuh kedua, `CLAIM_NO != ''`, TIDAK ditiru. Pega
// menomori klaim saat Save Outstanding (`GetSequenceNumber_SQL` dari
// `SaveOutStandingLife_Act` b8057); aplikasi ini menomorinya saat PENDAFTARAN
// (`services/pendaftaran.go`). Menirunya huruf demi huruf membuat keempat
// isian terkunci pada SETIAP klaim - fitur yang tidak pernah dapat dipakai.
// Maksudnya ("sesudah Save Outstanding, tanggal tidak diubah lagi") menunggu
// padanan Save Outstanding, yang belum punya rute.
//
// ⛔ BUTIR bk `[DIPUTUSKAN 28-09-2026, veto work owner]` - OQ-M9 ditutup.
// `ValidasiClaimReceived_Act` (dua langkah: RDB-List produk b248, Property-Set
// b425) menulis `.MAXCLAIM_RECEIVED` HANYA ke halaman (b582); nol `Obj-Save`,
// nol `RDB-Save`, nol SQL menyebutnya di seluruh korpus, dan satu-satunya
// pembacanya sel read-only `MAX CLAIM RECEIVED` `ClaimLifeDetailGCNM` b12131.
// Maka ia DIHITUNG SAAT BACA (PenandaTerimaKlaim), tanpa kolom dan tanpa
// migrasi. Di Pega penandanya pun tidak memblokir (OQ-G).
//
// Dibaca sesudah: tahap.go, diagnosa.go.

import "time"

// PenandaTerimaKlaim adalah `.MAXCLAIM_RECEIVED` seorang peserta - MURNI.
//
// `[terverifikasi]` b562 `@DateTimeDifference(DOL, CLAIM_RECEIVED_DATE, "D")`,
// b583 `@if(selisih <= MAXEXPIREDCLAIM, "", dd/MM/yyyy CLAIM_RECEIVED_DATE)`
// - persis `PenandaBatasHari`. Kosong berarti sah.
//
// Salah satu tanggal kosong: tidak ada yang diperiksa (activity-nya dipicu
// PERUBAHAN `CLAIM_RECEIVED_DATE`, b1120). Ambang kosong: galat, bukan nol.
func PenandaTerimaKlaim(p Peserta, maxExpiredClaim string) (string, error) {
	if p.TanggalKejadian == "" || p.TanggalTerimaKlaim == "" {
		return "", nil
	}
	return PenandaBatasHari(p.TanggalKejadian, p.TanggalTerimaKlaim, maxExpiredClaim)
}

// TanggalKlaimTeks adalah tiga tanggal dialog Edit Date selain DOL, sebagai
// TEKS yang tersimpan (`YYYY-MM-DD HH24:MI:SS`) - yang dibaca peserta.
type TanggalKlaimTeks struct {
	TanggalTerimaKlaim    string // CLAIM_RECEIVED_DATE b1076
	TanggalDokumenLengkap string // COMPLETE_DATE b1387
	TanggalKonfirmasi     string // CONFIRMATION_DATE b1626
}

// TanggalKlaim adalah tiga tanggal dialog Edit Date selain DOL.
//
// nil berarti KOSONG - kolomnya ditulis NULL, sama seperti Pega menulis
// isian yang dikosongkan (ADR-U-0027: kosong bukan tanggal nol).
type TanggalKlaim struct {
	TerimaKlaim    *time.Time // CLAIM_RECEIVED_DATE
	DokumenLengkap *time.Time // COMPLETE_DATE
	Konfirmasi     *time.Time // CONFIRMATION_DATE
}

// TahapBolehUbahTanggalKlaim menjawab gerbang tahap dialog Edit Date - MURNI.
//
// Dua syarat, keduanya dari XML: tombol `Edit Date` (b14115) berdiri di grid
// peserta `ClaimLifeDetailGCNM`, jadi tahapnya harus membuka grid itu
// (`TahapBergridPeserta`); dan `pyPosition == 'ReasLifeAdmin'` - tahapnya
// dipegang Admin. Irisannya hanya Outstanding Claim.
func TahapBolehUbahTanggalKlaim(t Tahap) bool {
	peran, ada := PeranPemegangTahap(t)
	return ada && peran == PeranAdminLife && TahapBergridPeserta(t)
}
