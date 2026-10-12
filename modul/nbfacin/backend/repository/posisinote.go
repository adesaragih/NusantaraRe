package repository

// T_WORK_POLIS.POSITION_NOTE (tiket 77, keputusan work owner 12-10-2026: "kalo tambahin kolom baru untuk positionnote gimana? jadi
// di position itu tetap offer, policy, binding" -> "Ya, kerjakan"). Pega punya DUA properti: PositionNote (workbasket pemegang
// case = kotak masuk) dan IsCedingConfirm (tahap Offer / Binding / Policy); sistem baru semula memakai satu kolom POSITION untuk
// keduanya ('Offer' = antrean marketing tahap Offer, butir 77).
//
// Penyimpanan baru (case NB FacIn, LINI FAC):
//   - POSITION_NOTE = workbasket: PosisiBaru services, kecuali 'Offer' -> 'ReasFacInMarketing' (PositionNote awal Pega) dan
//     'Decline' -> NULL (keputusan work owner: Decline mengosongkan antrean);
//   - POSITION = tahap dari FLAG_ONGOING_POLICY (0 Offer, 2 Binding, lainnya Policy), disinkronkan di AKHIR setiap transaksi tulis
//     posisi (sesudah semua perubahan flag - Transition38 men-set flag 2 di pernyataan terpisah).
//
// Kosakata services TIDAK berubah `[keputusan agent]`: repository membaca "posisi efektif" (sqlPosisiEfektif) dalam kosakata lama -
// 'ReasFacInMarketing' + flag 0 dibaca 'Offer' (tahap Binding tetap 'ReasFacInMarketing' + flag 2, seperti submitbinding
// tahapSubmitBinding); NOTE NULL pada case Resolved bertahap baru -> 'Decline' bila Resolved-Decline; baris yang belum dikonversi
// (NOTE NULL, aktif) -> POSITION lama apa adanya. Dengan begitu seluruh logika routing services (kotak masuk, penjaga 409, tahap
// baca-saja, Banding / Ask, scoring) tetap sama, hanya sumber kolomnya yang pindah.
//
// Peralihan ("jangan ubah yang sudah berhasil"): migrasi 254 belum dijalankan -> kolom tidak ada -> SQL PERSIS seperti sebelum
// tiket 77 (deteksi sekali per tabel, ORA-00904 di-cache dan dicoba ulang tiap jedaUlangCatatan).

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"nusantarare/inti/backend/db"
)

// KolomCatatanPosisi - kolom workbasket T_WORK_POLIS (migrasi 254).
const KolomCatatanPosisi = "POSITION_NOTE"

// Tahap T_WORK_POLIS.POSITION sesudah tiket 77 (= IsCedingConfirm Pega).
const (
	TahapOffer   = "Offer"
	TahapBinding = "Binding"
	TahapPolicy  = "Policy"
)

// TahapDariFlag - tahap dari FLAG_ONGOING_POLICY: "0" Offer, "2" Binding, lainnya Policy (kontrak tiket 77; sama dengan SQL
// sqlTahapDariFlag). Dipakai juga POSISI riwayat Submit (services posisiRiwayatFlag).
func TahapDariFlag(flag string) string {
	switch flag {
	case "0":
		return TahapOffer
	case "2":
		return TahapBinding
	}
	return TahapPolicy
}

// sqlTahapDariFlag - ekspresi SQL TahapDariFlag atas kolom flag.
func sqlTahapDariFlag(kolom string) string {
	return "CASE " + kolom + " WHEN '0' THEN '" + TahapOffer + "' WHEN '2' THEN '" + TahapBinding + "' ELSE '" + TahapPolicy + "' END"
}

// CatatanPosisiDari - POSITION_NOTE dari PosisiBaru kosakata services: 'Offer' -> WorkbasketMarketing, 'Decline' / "" -> NULL.
func CatatanPosisiDari(posisi string) any {
	switch posisi {
	case "", posisiDeclineSimpan:
		return nil
	case PosisiAwalCaseNB:
		return WorkbasketMarketing
	}
	return posisi
}

// posisiDeclineSimpan - PosisiBaru Decline services (posisiDecline) yang dulu ditulis apa adanya ke POSITION.
const posisiDeclineSimpan = "Decline"

// sqlPosisiEfektif - posisi kosakata lama dari baris alias `w` (lihat kepala berkas):
//  1. NOTE NULL + Resolved + POSITION tahap (Decline biner baru) -> 'Decline' bila Resolved-Decline, selain itu NULL;
//  2. x = NVL(POSITION_NOTE, POSITION) (baris belum dikonversi -> POSITION lama); x 'ReasFacInMarketing' + flag 0 -> 'Offer'
//     (antrean marketing tahap Offer; services menyamakan keduanya - posisiMarketing, posisiMarketingKotakMasuk, workbasketTahap),
//     sehingga pembacaan SAMA sebelum dan sesudah konversi docs/alat/konversi_position_note_254*.sql.
func sqlPosisiEfektif(w string) string {
	x := "NVL(" + w + ".POSITION_NOTE, " + w + ".POSITION)"
	return "CASE WHEN " + w + ".POSITION_NOTE IS NULL AND " + w + ".STATUS_WORK LIKE '" + polaSelesai + "' AND " + w + ".POSITION IN ('" +
		TahapOffer + "', '" + TahapBinding + "', '" + TahapPolicy + "') THEN CASE WHEN " + w + ".STATUS_WORK = 'Resolved-Decline' THEN '" +
		posisiDeclineSimpan + "' END WHEN " + x + " = '" + WorkbasketMarketing + "' AND " + w + ".FLAG_ONGOING_POLICY = '0' THEN '" +
		PosisiAwalCaseNB + "' ELSE " + x + " END"
}

// PosisiEfektif - padanan Go sqlPosisiEfektif (uji dan dokumentasi; SQL yang dipakai pembaca).
func PosisiEfektif(note, posisi, flag, status string) string {
	if note == "" && strings.HasPrefix(status, strings.TrimSuffix(polaSelesai, "%")) &&
		(posisi == TahapOffer || posisi == TahapBinding || posisi == TahapPolicy) {
		if status == "Resolved-Decline" {
			return posisiDeclineSimpan
		}
		return ""
	}
	x := note
	if x == "" {
		x = posisi
	}
	if x == WorkbasketMarketing && flag == "0" {
		return PosisiAwalCaseNB
	}
	return x
}

// sqlSinkronTahap - POSITION = tahap dari flag tersimpan, di akhir transaksi tulis posisi (hanya bila 254 sudah jalan).
func sqlSinkronTahap(work string) string {
	return "UPDATE " + work + " SET POSITION = " + sqlTahapDariFlag("FLAG_ONGOING_POLICY") + " WHERE ID = :1 AND LINI = :2"
}

// ---- deteksi kolom ----------------------------------------------------------------------------------------------

// jedaUlangCatatan - kolom belum ada: dicoba lagi sesudah selang ini (DBA dapat menjalankan 254 saat biner hidup).
const jedaUlangCatatan = time.Minute

var catatanPosisi = struct {
	sync.Mutex
	ada  map[string]bool
	coba map[string]time.Time
}{ada: map[string]bool{}, coba: map[string]time.Time{}}

func sqlUjiCatatanPosisi(work string) string {
	return "SELECT " + KolomCatatanPosisi + " FROM " + work + " WHERE ROWNUM < 1"
}

// adaCatatanPosisi - true bila T_WORK_POLIS.POSITION_NOTE ada (migrasi 254). Galat lain (koneksi) -> false + log: kueri utama yang
// gagal akan melaporkannya; perilaku lama tetap aman.
func adaCatatanPosisi(ctx context.Context, d *db.DB, work string) bool {
	catatanPosisi.Lock()
	ada, tahu := catatanPosisi.ada[work]
	sejak := catatanPosisi.coba[work]
	catatanPosisi.Unlock()
	if ada || (tahu && time.Since(sejak) < jedaUlangCatatan) {
		return ada
	}
	baris, err := d.QueryContext(ctx, sqlUjiCatatanPosisi(work))
	if err == nil {
		baris.Close()
	}
	ada = err == nil
	if err != nil && !kolomBelumAda(err) {
		log.Printf("nbfacin: uji kolom %s.%s: %v", work, KolomCatatanPosisi, err)
	}
	catatanPosisi.Lock()
	catatanPosisi.ada[work], catatanPosisi.coba[work] = ada, time.Now()
	catatanPosisi.Unlock()
	return ada
}

// gantiSekali - teks SQL `q` dengan TEPAT satu kemunculan `lama` diganti `baru`; selain itu q apa adanya (uji mengunci teksnya,
// jadi bentuk lama tidak pernah tertukar diam-diam).
func gantiSekali(q, lama, baru string) string {
	if strings.Count(q, lama) != 1 {
		return q
	}
	return strings.Replace(q, lama, baru, 1)
}
