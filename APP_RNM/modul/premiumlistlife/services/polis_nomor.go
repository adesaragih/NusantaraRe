package services

// Penerbitan `PL_NUMBER` - tiket 03 PremiumList Life.
//
// Untuk apa berkas ini: menjalankan rantai penomoran `SubmitPremiumList_Act`
// langkah 10-13 dalam satu transaksi pendek, lalu menyimpan nomornya.
//
// ⛔ NOL PEMBENTUK BENTUK NOMOR DI SINI (AC tiket 03). Yang merakit
// `models.NomorPL`; berkas ini hanya mengurutkan pengambilan bahannya. Bentuk
// nomor yang tersebar di lapisan layanan adalah bentuk yang berubah tanpa
// satu pun uji murni menyadarinya.
//
// ⛔ TRANSAKSINYA PENDEK, DAN SENGAJA BERBEDA DARI KLAIM. `UrutNomorBerikut`
// memegang `SELECT … FOR UPDATE` atas baris penghitung; selama transaksi ini
// terbuka, tidak seorang pun dapat menerbitkan nomor premium list. Karena itu
// yang masuk ke dalamnya HANYA pengambilan nomor dan penyimpanannya - bukan
// seluruh penyimpanan polis. Pendaftaran klaim memilih sebaliknya dengan
// alasannya sendiri (lihat pendaftaran.go); keduanya dinyatakan, bukan
// kebetulan.
//
// ⚠️ Sejak tiket 05a bagian 2 rantai yang SAMA (`terbitkanDalam`) juga
// berjalan di transaksi yang lebih panjang - `SummaryPremiumList.Submit`,
// bersama rekap dan salinan peserta warisan (polis_summary.go). Kalimat
// "pendek" di atas berlaku untuk `Terbitkan`; harga transaksi panjang itu
// dinyatakan di sana.
//
// ⛔ NOMOR LAHIR SEKALI. Gerbangnya dibaca lebih dahulu, sebelum penghitung
// disentuh sama sekali: submit kedua atas polis yang sudah bernomor
// mengembalikan nomor yang sama dan TIDAK menggerakkan penghitung. Penghitung
// yang bergerak tanpa nomor tersimpan adalah nomor yang hilang selamanya.
//
// Dibaca sesudah: models/polis_nomor.go (bentuknya),
// repository/penomor.go (penghitungnya).

import (
	"context"
	"fmt"
	"time"

	"nusantarare/inti"
	"nusantarare/inti/db"
	"nusantarare/inti/galat"
	"nusantarare/inti/penomor"
	"nusantarare/modul/premiumlistlife/models"
	"nusantarare/modul/premiumlistlife/repository"
)

// HasilNomorPL adalah jawaban penerbitan atau pembacaan nomor.
//
// ⚠️ Tag JSON ditulis eksplisit. Nama medan Go adalah urusan Go; nama kunci
// JSON adalah kontrak dengan layar - dan kedua sisi pernah berselisih di repo
// ini lebih dari sekali.
type HasilNomorPL struct {
	// Nomor adalah `PL_NUMBER`; kosong berarti polis belum bernomor.
	Nomor string `json:"nomor"`
	// BaruTerbit membedakan nomor yang baru lahir dari nomor yang sudah ada.
	//
	// ⚠️ Layar memerlukannya: "nomor Anda RNML-QRLF.09.26.00007" dan "nomor
	// polis ini sudah RNML-QRLF.09.26.00007" adalah dua kalimat berbeda, dan
	// yang kedua memberi tahu orang bahwa tombolnya tidak perlu ditekan lagi.
	BaruTerbit bool `json:"baruTerbit"`
	// Periode adalah `MM.YYYY` dari penghitung - kosong bila tidak menerbitkan.
	Periode string `json:"periode"`
	// BarisPeserta adalah cacah baris peserta yang memuat nomor ini.
	BarisPeserta int `json:"barisPeserta"`
}

// NomorPremiumList melayani penerbitan dan pembacaan `PL_NUMBER`.
type NomorPremiumList struct {
	svc *Service
	// jam diserahkan supaya periode nomor dapat diuji tanpa menunggu tanggal
	// nyata - sejalan dengan `Periode` di polis_periode.go.
	jam func() time.Time
}

// NomorPremiumList menyusun layanannya.
func (s *Service) NomorPremiumList() *NomorPremiumList {
	return &NomorPremiumList{svc: s, jam: time.Now}
}

// Terbitkan menerbitkan satu `PL_NUMBER` untuk polis, sekali seumur polis.
//
// Urutannya: gerbang lahir-sekali → identitas → awalan → hari tutup buku →
// periode → kunci baris penghitung → naikkan → rakit → simpan. Gerbangnya di
// paling depan, sebelum penghitung disentuh.
func (n *NomorPremiumList) Terbitkan(ctx context.Context, pelaku inti.Pelaku, polisID string) (
	HasilNomorPL, error) {

	if err := inti.WajibIdentitas(pelaku); err != nil {
		return HasilNomorPL{}, err
	}
	if n == nil || n.svc == nil || !n.svc.PunyaDatabase() {
		return HasilNomorPL{}, db.ErrTanpaOracle
	}
	if polisID == "" {
		return HasilNomorPL{}, fmt.Errorf("%w: id polis kosong", galat.ErrPermintaanTidakSah)
	}

	// ⛔ GERBANG KASUS TERTUTUP - butir bb. Ia membaca `T_WORK_POLIS` lewat
	// `models.KasusPolisTertutup`, BUKAN `PastikanKasusTerbuka` yang membaca
	// tabel kerja klaim: layanan polis yang menanyai tabel klaim akan selalu
	// menjawab "tidak ada", yaitu penjaga yang tidak pernah menolak apa pun.
	//
	// ⚠️ Di LUAR transaksi penomoran, dan sengaja: menolak lebih dahulu
	// berarti kasus tertutup tidak pernah sempat menyentuh baris penghitung.
	keadaanKerja, err := repository.NewWorkPolis(n.svc.DB()).Keadaan(ctx, polisID)
	if err != nil {
		return HasilNomorPL{}, err
	}
	if models.KasusPolisTertutup(keadaanKerja.Status) {
		return HasilNomorPL{}, fmt.Errorf("%w: polis %q berstatus %q",
			ErrKasusPolisTertutup, polisID, keadaanKerja.Status)
	}

	saat := n.jam()
	var hasil HasilNomorPL
	err = n.svc.DalamTransaksi(ctx, func(tx *db.Tx) error {
		var err error
		hasil, err = n.terbitkanDalam(ctx, tx, polisID, saat)
		return err
	})
	if err != nil {
		return HasilNomorPL{}, err
	}
	return hasil, nil
}

// terbitkanDalam adalah rantai penomoran `SubmitPremiumList_Act` langkah
// 10-14 (bahan, urut, rakit, lalu tulis nomor) di dalam transaksi MILIK
// PEMANGGIL.
//
// ⚠️ Diekstrak di tiket 05a bagian 2: penomoran kini berjalan di dua
// transaksi yang berbeda - sendirian (`Terbitkan`, tiket 03) dan bersama rekap
// serta salinan warisan (`SummaryPremiumList.Submit`). Satu fungsi untuk
// keduanya, supaya gerbang lahir-sekali tidak pernah berlaku di satu jalur
// dan terlupa di jalur lain.
func (n *NomorPremiumList) terbitkanDalam(ctx context.Context, tx *db.Tx,
	polisID string, saat time.Time) (HasilNomorPL, error) {

	nomorPolis := repository.NewNomorPolis(n.svc.DB())
	penghitung := penomor.NewPenomor(n.svc.DB())
	keadaan, err := nomorPolis.Keadaan(ctx, tx, polisID)
	if err != nil {
		return HasilNomorPL{}, err
	}
	// ⛔ GERBANG LAHIR-SEKALI. Kembali di sini berarti penghitung sama
	// sekali tidak tersentuh - bukan sekadar nomornya tidak berubah.
	if keadaan.Bernomor() {
		// ⛔ TETAPI BELUM TENTU SELESAI. `MIN` dan `MAX` melewati NULL,
		// jadi polis yang separuh barisnya bernomor terbaca "sudah
		// bernomor" di sini. Kembali begitu saja meninggalkan baris yang
		// kosong TIDAK PERNAH terisi - dan itu keadaan yang akan lahir
		// sendiri begitu unggahan CSV (tiket 04) menambah peserta SESUDAH
		// polis bernomor.
		//
		// Yang benar bukan nomor baru dan bukan galat: baris baru diberi
		// nomor yang SUDAH ada. Penghitung tetap tidak tersentuh, jadi
		// "lahir sekali" tetap utuh.
		tersentuh := 0
		if !keadaan.Utuh() {
			n, err := nomorPolis.TulisNomor(ctx, tx, polisID, keadaan.Nomor)
			if err != nil {
				return HasilNomorPL{}, err
			}
			tersentuh = n
		}
		return HasilNomorPL{
			Nomor:        keadaan.Nomor,
			BaruTerbit:   false,
			BarisPeserta: keadaan.CacahBernomor + tersentuh,
		}, nil
	}
	if keadaan.CacahPeserta == 0 {
		return HasilNomorPL{}, repository.ErrPolisTanpaPeserta
	}

	identitas, err := nomorPolis.Identitas(ctx, tx, polisID)
	if err != nil {
		return HasilNomorPL{}, err
	}
	// ⛔ Awalannya DI-LOOKUP, bukan konstanta (AC tiket 03). Awalan itu
	// milik basis data, dan lingkungan yang berbeda dapat memakai yang
	// berbeda.
	awalan, err := penghitung.AwalanProduksi(ctx, tx, penomor.TipeKodeProduksiLife)
	if err != nil {
		return HasilNomorPL{}, err
	}
	hariClosing, err := penghitung.HariClosing(ctx, tx)
	if err != nil {
		return HasilNomorPL{}, err
	}
	// ⛔ Periodenya aturan tiket 02, bukan `time.Now()` mentah - dan
	// aturannya satu, dipakai penomoran klaim maupun premium list.
	periode, err := penomor.HitungPeriodeNomor(saat, hariClosing)
	if err != nil {
		return HasilNomorPL{}, err
	}
	// ⛔ Kunci penghitungnya SATU untuk keempat tipe - lihat
	// `models.JenisPenghitungPL`.
	urut, err := penghitung.UrutNomorBerikut(ctx, tx,
		models.ClassPenghitungPL, models.JenisPenghitungPL(awalan), periode, saat)
	if err != nil {
		return HasilNomorPL{}, err
	}
	nomor, err := models.NomorPL(models.BahanNomorPL{
		Awalan:        awalan,
		Tipe:          identitas.Tipe,
		KodeBisnis:    identitas.KodeBisnis,
		PeriodeMMYYYY: periode.MMYYYY,
		Urut:          urut,
	})
	if err != nil {
		return HasilNomorPL{}, err
	}
	tersentuh, err := nomorPolis.TulisNomor(ctx, tx, polisID, nomor)
	if err != nil {
		return HasilNomorPL{}, err
	}
	return HasilNomorPL{
		Nomor:        nomor,
		BaruTerbit:   true,
		Periode:      periode.MMYYYY,
		BarisPeserta: tersentuh,
	}, nil
}
