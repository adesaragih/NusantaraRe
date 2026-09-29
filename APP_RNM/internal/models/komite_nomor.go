package models

// Nomor akseptasi Komite - tiket 04a Komite Claim Life. MURNI.
//
// `[terverifikasi]` `Komite Claim Life/Activity/KomitePostAdjustment.xml`
// langkah 4 "Approve Last Komite" (gerbang b5695 `AcceptStatus = 1 &&
// KomiteCount == KomiteLoop`), tiap sub-langkah bergerbang
// `AdjustmentList(…).ACCEPTEDNO == ""`:
//
//	4.7  `GetKodeProdLife_SQL`       -> ParamSeq.HASIL3 (awalan `KODE_PRODUKSI`)
//	4.8  CARI1 = pyWorkPage.pxObjClass, CARI2 = HASIL3 + "A"
//	4.9  `GetSequenceNumber_SQL`     -> HASIL1 `MM.YYYY`, HASIL2 urut
//	4.10 HASIL1 = substring(0,2) + "." + substring(5,7)   (MM.YY)
//	4.11 QR,QP: ACCEPTEDNO = HASIL3 + "A"  + BusinessCode + "." + HASIL1 + "." + HASIL2
//	4.12 TR,TP: ACCEPTEDNO = HASIL3 + "AR" + BusinessCode + "." + HASIL1 + "." + HASIL2
//
// ⛔ SATU PENGHITUNG untuk kedua cabang: `JENIS = HASIL3 + "A"` disusun di
// 4.8, SEBELUM percabangan - retro tidak punya penghitung `…AR` sendiri.
//
// ⛔ `Generate_NoAccept_KMT_Life`/`…Retro` (4.4/4.5) ter-remark (b1725, b1943)
// - tidak dimigrasikan (AC 28 spec).
//
// Dibaca sesudah: polis_nomor.go (`RakitNomorPL`, `PeriodeNomorPL` dipakai ulang).

import (
	"fmt"
	"strings"

	"nusantarare/inti/penomor"
)

// ClassPenghitungKomiteLife adalah `CLASS` penghitung - `pyWorkPage.pxObjClass`.
//
// `[data DBA]` pasangan `(ASM-FW-GCNMFW-Work-KomiteLife, RNML-A)` ada di
// `GENERATE_SEQUENCE_NUMBER` (SUMBER-PENOMORAN-DBA.md).
const ClassPenghitungKomiteLife = "ASM-FW-GCNMFW-Work-KomiteLife"

// JenisPenghitungKomite adalah `CARI2 = HASIL3 + "A"` (4.8).
func JenisPenghitungKomite(awalan string) string { return awalan + "A" }

// KodeCabangAkseptasiKomite memilih "A" (QR/QP) atau "AR" (TR/TP).
func KodeCabangAkseptasiKomite(tipe string) (string, error) {
	switch strings.TrimSpace(tipe) {
	case penomor.TipePLQuotationRealisasi, penomor.TipePLQuotationProposal:
		return "A", nil
	case penomor.TipePLTreatyProposal, penomor.TipePLTreatyRealisasi:
		return "AR", nil
	}
	return "", fmt.Errorf("%w: %q", penomor.ErrTipePLTanpaCabang, tipe)
}

// NomorAkseptasiKomite merakit nomor akseptasi tingkat akhir komite.
//
// ⚠️ BENTUKNYA SAMA dengan nomor akseptasi jalur Claim Life
// (`services.RakitNomorAkseptasi`: `RNML-A` + kode bisnis + `.MM.YY.` + urut)
// tetapi PENGHITUNGNYA BERBEDA (`GENERATE_SEQUENCE_NUMBER` lawan
// `ACCEPTATIONNOLIFE_SEQ`). Dua jalur, satu ruang nomor, dua penghitung -
// tabrakan mungkin, dan korpus Komite TIDAK memeriksanya (nol
// `GetAcceptedNoCL` di `KomitePostAdjustment`). Pemeriksanya di layanan;
// kebijakannya `[terbuka — work owner]` OQ-K-04a.
func NomorAkseptasiKomite(awalan, tipe, kodeBisnis, periodeMMYYYY string, urut int) (string, error) {
	if strings.TrimSpace(awalan) == "" {
		return "", penomor.ErrAwalanProduksiKosong
	}
	cabang, err := KodeCabangAkseptasiKomite(tipe)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(kodeBisnis) == "" {
		return "", penomor.ErrKodeBisnisKosong
	}
	mmYY, err := penomor.PeriodeNomorPL(periodeMMYYYY)
	if err != nil {
		return "", err
	}
	if urut < 1 {
		return "", fmt.Errorf("models: urut penghitung %d tidak masuk akal", urut)
	}
	return penomor.RakitNomorPL(awalan, cabang, strings.TrimSpace(kodeBisnis), mmYY, urut), nil
}
