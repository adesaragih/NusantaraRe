package services

// Rekap premium list dan submit-nya - tiket 05a bagian 2 PremiumList Life.
//
// Untuk apa berkas ini: dua layanan layar `ShowLifePremiumSummary`.
//
//	Lihat  - rekap per mata uang dihitung dari peserta, TANPA menyimpan apa
//	         pun (bacaan `SavePremiumList_Act` langkah 12, yang menerapkan
//	         `AppendCurrencySummary_DT` untuk ditampilkan).
//	Submit - SATU transaksi: penomoran (`SubmitPremiumList_Act` langkah
//	         10-14) → rekap `T_PREMIUM_LIST_SUMMARY` hapus-lalu-sisip →
//	         PL-09 rekap yang sama ke `M_LIFE_PREMIUM_SUMMARY` →
//	         pl2 salinan peserta ke `M_LIFE_PREMIUM_DETAIL` → tiket 05b
//	         `finishAssignment`: kasus ditutup Resolved-Completed + jejak →
//	         commit.
//
// Bagian simpannya (`simpanDalam`) juga dijalankan `Utility1` sesudah
// `Confirm` di tahap detail (`Transition7`) - lewat `Penawaran.terapkan`.
//
// ⛔ SATU TRANSAKSI, dan urutannya bagian dari kebenaran (AC 24 spec).
// Nomor lebih dahulu, karena salinan warisan berkunci `PL_NUMBER`; rekap
// sesudah nomor, supaya rekap tidak pernah tersimpan untuk polis yang
// penomorannya gagal. Kegagalan di tingkat mana pun membatalkan SELURUHNYA:
// tidak ada nomor tanpa rekap, tidak ada rekap tanpa nomor (AC 35, 36).
//
// ⚠️ Transaksi ini LEBIH PANJANG daripada `Terbitkan` (tiket 03), dan itu
// harga yang dibayar dengan sadar: selama ia terbuka, baris penghitung
// terkunci `FOR UPDATE`. Brief giliran 10 menuntut penomoran dan rekap
// "satu transaksi"; memisahkannya lagi menghidupkan kembali titik potong
// Pega (nomor yatim yang sudah ter-commit).
//
// ⛔ LANGKAH 15 TIDAK DITIRU - pl7. `@replaceAll(.PremiumListSummary.PL_NUMBER,
// Local.CurrentMMYY, Local.NextMMYY)` mencari `.MM.YYYY.` (b1142) di nomor
// yang berbentuk `.MM.YY.` (b3021/b3075/b3126): tidak pernah cocok, jadi di
// produksi ia no-op. Menirunya berarti menulis kode yang tidak pernah
// melakukan apa-apa - atau, bila "diperbaiki", menggeser periode nomor yang
// Pega sendiri tidak pernah geser.
//
// ⚠️ Langkah 16 (`pyMemo = 1`, `IsJsonPolis = 0`) tidak ditiru: keduanya
// properti halaman kerja Pega tanpa kolom di skema kami, dan `IsJsonPolis`
// milik jalur JSON yang dibuang (pl1, spec §6).
//
// Dibaca sesudah: models/polis_summary.go (rumusnya), polis_nomor.go,
// repository/polis_summary.go, repository/polis_warisan.go.

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"nusantarare/inti"
	"nusantarare/inti/db"
	"nusantarare/inti/galat"
	"nusantarare/inti/jejak"
	"nusantarare/inti/outbox"
	"nusantarare/inti/utils"
	"nusantarare/modul/premiumlistlife/models"
	"nusantarare/modul/premiumlistlife/repository"
)

// RekapTampil adalah satu baris rekap untuk layar - uang sebagai TEKS.
//
// ⛔ Teks, bukan angka JSON. `json.Number` atau float membuat `975.0000`
// tiba di peramban sebagai `975`, dan ekor nol itu justru bukti bahwa
// pembulatan empat angka sudah terjadi.
type RekapTampil struct {
	Currency   string            `json:"currency"`
	Premium    string            `json:"premium"`
	Commission string            `json:"commission"`
	Balance    string            `json:"balance"`
	Jumlah     map[string]string `json:"jumlah"`
	CacahBaris int               `json:"cacahBaris"`
}

// HasilRekap adalah jawaban `Lihat`.
type HasilRekap struct {
	Tipe  string        `json:"tipe"`
	Rekap []RekapTampil `json:"rekap"`
}

// HasilSubmitSummary adalah jawaban `Submit`.
type HasilSubmitSummary struct {
	Nomor HasilNomorPL  `json:"nomor"`
	Rekap []RekapTampil `json:"rekap"`
	// RekapDihapus adalah cacah baris rekap lama yang diganti.
	RekapDihapus int `json:"rekapDihapus"`
	// PesertaWarisan adalah cacah baris yang tersalin ke tabel warisan.
	PesertaWarisan int `json:"pesertaWarisan"`
	// EfekKeluar - tiket 06. Ringkasan, bukan galat: simpannya sudah commit.
	EfekKeluar RingkasEfek `json:"efekKeluar"`
}

// RingkasEfek adalah hasil efek keluar yang dapat dibaca layar.
//
// ⛔ Nama efek dan sebabnya saja - tanpa alamat. Sebab galat resolver dapat
// menyebut kunci kategori, tidak pernah URL-nya (`AlamatLayanan`).
type RingkasEfek struct {
	// Dilewati - lingkungan bukan produksi; BUKAN gagal.
	Dilewati bool     `json:"dilewati"`
	Gagal    []string `json:"gagal"`
	// TidakTerantre - kegagalan yang bahkan tidak dapat diantre.
	TidakTerantre int `json:"tidakTerantre"`
}

// ringkasEfek menerjemahkan hasil penyalur menjadi ringkasan layar.
func ringkasEfek(h outbox.HasilSalur) RingkasEfek {
	r := RingkasEfek{Dilewati: h.Dilewati, Gagal: []string{}, TidakTerantre: len(h.GagalDiantre)}
	for _, g := range h.Gagal {
		r.Gagal = append(r.Gagal, g.Nama)
	}
	return r
}

// SummaryPremiumList melayani layar `ShowLifePremiumSummary`.
type SummaryPremiumList struct {
	svc      *Service
	nomor    *NomorPremiumList
	jejak    jejak.Jejak
	penyalur *PenyalurPolis
}

// ErrSubmitBukanTahapSummary - `Submit` summary di luar tahapnya.
var ErrSubmitBukanTahapSummary = errors.New(
	"services: Submit summary hanya dari tahap Input Premium Summary")

// SummaryPremiumList menyusun layanannya dengan jejak bawaan yang gagal terang.
func (s *Service) SummaryPremiumList() *SummaryPremiumList {
	return &SummaryPremiumList{svc: s, nomor: s.NomorPremiumList(),
		jejak: jejak.JejakBelumDiputuskan{}, penyalur: penyalurPolisBawaan(s)}
}

// DenganJejak mengganti perekamnya.
func (s *SummaryPremiumList) DenganJejak(j jejak.Jejak) *SummaryPremiumList {
	salin := *s
	salin.jejak = j
	return &salin
}

// DenganPenyalur mengganti penyalur efek keluarnya - tiket 06.
func (s *SummaryPremiumList) DenganPenyalur(p *PenyalurPolis) *SummaryPremiumList {
	salin := *s
	salin.penyalur = p
	return &salin
}

// keTampil mengubah rekap desimal menjadi teks layar.
func keTampil(rekap []models.RekapMataUang) []RekapTampil {
	keluar := make([]RekapTampil, 0, len(rekap))
	for _, r := range rekap {
		t := RekapTampil{
			Currency:   r.Currency,
			Premium:    utils.FormatDecimal(r.Premium),
			Commission: utils.FormatDecimal(r.Commission),
			Balance:    utils.FormatDecimal(r.Balance),
			Jumlah:     map[string]string{},
			CacahBaris: r.CacahBaris,
		}
		for _, k := range models.KolomJumlahSummary {
			t.Jumlah[k] = utils.FormatDecimal(r.Jumlah[k])
		}
		keluar = append(keluar, t)
	}
	return keluar
}

// periksaDasar memeriksa pelaku, basis data, dan id - dipakai kedua layanan.
func (s *SummaryPremiumList) periksaDasar(pelaku inti.Pelaku, polisID string) error {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return err
	}
	if s == nil || s.svc == nil || !s.svc.PunyaDatabase() {
		return db.ErrTanpaOracle
	}
	if polisID == "" {
		return fmt.Errorf("%w: id polis kosong", galat.ErrPermintaanTidakSah)
	}
	return nil
}

// siapkan = periksaDasar + gerbang kasus tertutup, untuk layanan yang MENULIS.
//
// Mengembalikan keadaan kerja - tahapnya dibutuhkan `Submit`.
func (s *SummaryPremiumList) siapkan(ctx context.Context, pelaku inti.Pelaku, polisID string) (
	repository.KeadaanPolis, error) {

	if err := s.periksaDasar(pelaku, polisID); err != nil {
		return repository.KeadaanPolis{}, err
	}
	// ⛔ GERBANG KASUS TERTUTUP - butir bb, lewat `T_WORK_POLIS`.
	keadaan, err := repository.NewWorkPolis(s.svc.DB()).Keadaan(ctx, polisID)
	if err != nil {
		return repository.KeadaanPolis{}, err
	}
	if models.KasusPolisTertutup(keadaan.Status) {
		return repository.KeadaanPolis{}, fmt.Errorf("%w: polis %q berstatus %q",
			ErrKasusPolisTertutup, polisID, keadaan.Status)
	}
	return keadaan, nil
}

// rekapDalam membaca peserta dan menghitung rekap di dalam transaksi.
func (s *SummaryPremiumList) rekapDalam(ctx context.Context, tx *db.Tx,
	polisID string) (string, []models.RekapMataUang, error) {

	identitas, err := repository.NewNomorPolis(s.svc.DB()).Identitas(ctx, tx, polisID)
	if err != nil {
		return "", nil, err
	}
	baris, mataUang, err := repository.NewSummaryPolis(s.svc.DB()).BarisUang(ctx, tx, polisID)
	if err != nil {
		return "", nil, err
	}
	if len(baris) == 0 {
		return "", nil, repository.ErrPolisTanpaPeserta
	}
	// ⛔ `Type` di luar QR/QP/TP/TR DITOLAK di sini - `ErrTipePLTanpaCabang`
	// - bukan direkap dengan cabang kosong yang menghasilkan nol.
	rekap, err := models.RekapPerMataUang(identitas.Tipe, baris, mataUang)
	if err != nil {
		return "", nil, err
	}
	return identitas.Tipe, rekap, nil
}

// Lihat menghitung rekap untuk layar tanpa menyimpan apa pun.
//
// ⚠️ Ia tetap memakai transaksi, supaya identitas polis dan baris pesertanya
// terbaca dari SATU keadaan. Transaksinya tidak menulis apa-apa.
//
// ⚠️ TANPA gerbang kasus tertutup, dengan sengaja: ia membaca saja. Melarang
// orang MELIHAT rekap polis yang sudah selesai tidak melindungi apa pun -
// alasan yang sama dengan tinjauan unggahan (tutupkontrak_test.go).
func (s *SummaryPremiumList) Lihat(ctx context.Context, pelaku inti.Pelaku, polisID string) (
	HasilRekap, error) {

	if err := s.periksaDasar(pelaku, polisID); err != nil {
		return HasilRekap{}, err
	}
	var hasil HasilRekap
	err := s.svc.DalamTransaksi(ctx, func(tx *db.Tx) error {
		tipe, rekap, err := s.rekapDalam(ctx, tx, polisID)
		if err != nil {
			return err
		}
		hasil = HasilRekap{Tipe: tipe, Rekap: keTampil(rekap)}
		return nil
	})
	if err != nil {
		return HasilRekap{}, err
	}
	return hasil, nil
}

// Submit adalah tombol `Submit` layar summary - tiket 05a + 05b.
//
// `[terverifikasi]` b27471 `Submit` → b26414 `InsertJsonPolisLife_Act` →
// b26442 `finishAssignment` → `Transition2` → `END52` Resolved-Completed.
// Yang tersisa dari activity itu sesudah JSON dibuang (`simpanDalam`)
// berjalan di transaksi yang SAMA dengan penutupan kasus dan jejaknya -
// lewat `Penawaran.terapkan`, satu tempat untuk menutup kasus polis.
//
// ⛔ HANYA dari tahap `Input Premium Summary`. `finishAssignment` menyelesaikan
// assignment yang sedang dibuka; memanggilnya dari tahap lain berarti
// menutup kasus lewat konektor yang tahap itu tidak punya.
func (s *SummaryPremiumList) Submit(ctx context.Context, pelaku inti.Pelaku,
	polisID string, saat time.Time) (HasilSubmitSummary, error) {

	keadaan, err := s.siapkan(ctx, pelaku, polisID)
	if err != nil {
		return HasilSubmitSummary{}, err
	}
	if strings.TrimSpace(keadaan.Status) != models.TahapPolisSummary {
		return HasilSubmitSummary{}, fmt.Errorf("%w: polis %q berada di %q, bukan %q",
			ErrSubmitBukanTahapSummary, polisID, keadaan.Status, models.TahapPolisSummary)
	}
	return s.svc.Penawaran().DenganJejak(s.jejak).DenganPenyalur(s.penyalur).terapkan(ctx, pelaku, keadaan,
		models.PenyelesaianSummary(), "Submit", saat)
}

// simpanDalam menomori polis, mengganti rekapnya, dan menyalin peserta
// warisan - di dalam transaksi MILIK PEMANGGIL.
//
// ⛔ URUTANNYA DIKUNCI `TestSimpanDalamUrutanTerkunci`: nomor → rekap →
// ganti rekap → kepala + ganti summary warisan → sumber warisan → ganti warisan.
func (s *SummaryPremiumList) simpanDalam(ctx context.Context, tx *db.Tx,
	polisID string, saat time.Time) (HasilSubmitSummary, error) {

	ringkas := repository.NewSummaryPolis(s.svc.DB())
	summaryWarisan := repository.NewSummaryWarisan(s.svc.DB())
	warisan := repository.NewPesertaWarisan(s.svc.DB())

	// 1. Penomoran - langkah 10-14. Lahir sekali: simpan ulang memakai nomor
	//    yang sama dan tidak menggerakkan penghitung.
	nomor, err := s.nomor.terbitkanDalam(ctx, tx, polisID, saat)
	if err != nil {
		return HasilSubmitSummary{}, err
	}
	// 2. Rekap - rumus murni atas peserta yang baru saja dinomori.
	_, rekap, err := s.rekapDalam(ctx, tx, polisID)
	if err != nil {
		return HasilSubmitSummary{}, err
	}
	// 3. Hapus lalu sisip rekap.
	dihapus, _, err := ringkas.GantiRekap(ctx, tx, polisID, rekap)
	if err != nil {
		return HasilSubmitSummary{}, err
	}
	// 3b. PL-09 - rekap yang SAMA ke tabel warisan `M_LIFE_PREMIUM_SUMMARY`,
	//     seperti `PEGA_M_LIFE_PREMIUM_SUMMARY` (langkah 8, b2228-b3252),
	//     berkunci nomor + work. Prosedurnya tidak dipanggil (keputusan o).
	kepala, err := ringkas.KepalaSummaryWarisan(ctx, tx, polisID, nomor.Nomor)
	if err != nil {
		return HasilSubmitSummary{}, err
	}
	if _, _, err := summaryWarisan.Ganti(ctx, tx, kepala, rekap); err != nil {
		return HasilSubmitSummary{}, err
	}
	// 4-5. pl2 - salinan peserta ke tabel warisan, berkunci nomor + work.
	sumber, err := ringkas.SumberWarisan(ctx, tx, polisID)
	if err != nil {
		return HasilSubmitSummary{}, err
	}
	_, disalin, err := warisan.Ganti(ctx, tx, nomor.Nomor, polisID, sumber)
	if err != nil {
		return HasilSubmitSummary{}, err
	}
	return HasilSubmitSummary{
		Nomor:          nomor,
		Rekap:          keTampil(rekap),
		RekapDihapus:   dihapus,
		PesertaWarisan: disalin,
	}, nil
}
