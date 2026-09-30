package services

// Gerbang `Close Claim` - kelompok Detail & Tutup.
//
// Untuk apa berkas ini: tombol `Close Claim` di Pega tidak punya aksi lain
// selain gerbangnya. `Section/CloseClaim_Section.xml` b1081 `Close Claim`
// memanggil `pyActivity` b1101 `ProtectCloseClaim_act`, dan activity itulah
// yang memeriksa seluruh peserta lalu - hanya bila bersih - memanggil
// `FinishAssignment`. Jadi gerbangnya ADALAH aksinya.
//
// Aturannya sendiri berumah di `models/tutupklaim.go`; berkas ini yang
// membacakan klaimnya dan menyusun urutan barisnya.
//
// Dibaca sesudah: klaimlife.go.

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/jejak"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/claimlife/models"
	"nusantarare/modul/claimlife/repository"
)

// TutupKlaim membungkus pemeriksaan gerbang tutup dan penutupannya.
type TutupKlaim struct {
	svc   *Service
	jejak jejak.Jejak
}

// Tutup menyusun layanannya dengan jejak bawaan yang gagal terang.
//
// ⛔ ADR-U-0007 menuntut setiap transisi kasus terekam. Penutupan adalah
// transisi yang paling tidak dapat dibatalkan dari semuanya, jadi ia tidak
// boleh terjadi tanpa jejak - dan jejak bawaannya sengaja GAGAL, bukan diam.
func (s *Service) Tutup() *TutupKlaim {
	return &TutupKlaim{svc: s, jejak: jejak.JejakBelumDiputuskan{}}
}

// DenganJejak mengganti perekamnya - dipakai test, dan kelak oleh tiket 09.
func (t *TutupKlaim) DenganJejak(j jejak.Jejak) *TutupKlaim {
	return &TutupKlaim{svc: t.svc, jejak: j}
}

// HasilPeriksaTutup adalah jawaban gerbang.
type HasilPeriksaTutup struct {
	// Boleh menjawab padanan prasyarat `ProtectLife.CARI1==""` (b992).
	Boleh bool `json:"boleh"`
	// Penghalang memuat SELURUH peserta yang menahan, bukan yang pertama.
	Penghalang []PenghalangTampil `json:"penghalang"`
}

// PenghalangTampil adalah satu penghalang beserta kalimat yang dilihat orang.
//
// Kalimatnya disusun di server, bukan di layar: ia VERBATIM dari rule Pega
// (kecuali penunjuk pesertanya), dan kalimat verbatim yang disusun di dua
// tempat akan bergeser di salah satunya.
type PenghalangTampil struct {
	Urutan          int    `json:"urutan"`
	NomorSertifikat string `json:"nomorSertifikat"`
	Pesan           string `json:"pesan"`
}

// BarisTutupDari menyusun masukan gerbang dari sebuah klaim.
//
// ⛔ `Urutan` adalah POSISI peserta di daftar, mulai 1 - padanan
// `.pxListSubscript` (b398). Ia muncul di kalimat yang dibaca orang, jadi ia
// harus nomor yang mereka lihat di layar; memakai pengenal baris di situ
// membuat pesannya menunjuk sesuatu yang tidak ada di layar mana pun.
func BarisTutupDari(klaim *models.Klaim) []models.BarisTutup {
	if klaim == nil {
		return nil
	}
	out := make([]models.BarisTutup, 0, len(klaim.Peserta))
	for i, p := range klaim.Peserta {
		out = append(out, models.BarisTutup{
			Urutan:          i + 1,
			NomorSertifikat: p.NomorSertifikat,
			KodeStatus:      p.KodeStatus,
		})
	}
	return out
}

// Periksa menjawab apakah klaim itu boleh ditutup, dan bila tidak, siapa saja
// yang menahannya.
func (t *TutupKlaim) Periksa(ctx context.Context, id string) (HasilPeriksaTutup, error) {
	if t == nil || t.svc == nil || !t.svc.PunyaDatabase() {
		return HasilPeriksaTutup{}, db.ErrTanpaOracle
	}
	klaim, err := t.svc.KlaimLife().Ambil(ctx, id)
	if err != nil {
		return HasilPeriksaTutup{}, err
	}
	baris := BarisTutupDari(klaim)

	// ⛔ Daftar KOSONG, bukan nil: `encoding/json` menulis nil sebagai
	// `null`, dan layar yang menerima `null` harus menjaganya sendiri. Daftar
	// kosong berarti "diperiksa, tidak ada penghalang" - dan itu yang benar.
	tampil := []PenghalangTampil{}
	for _, p := range models.PenghalangTutupKlaim(baris) {
		tampil = append(tampil, PenghalangTampil{
			Urutan:          p.Urutan,
			NomorSertifikat: p.NomorSertifikat,
			Pesan:           p.Pesan(),
		})
	}
	return HasilPeriksaTutup{Boleh: len(tampil) == 0, Penghalang: tampil}, nil
}

// ErrMasihAdaPenghalang - penutupan ditolak karena ada peserta yang belum
// diaksep. Ia membawa penghalangnya, sebab pesan tanpa daftar memaksa
// pemakai menutup berulang kali untuk menemukan satu penghalang tiap kali.
var ErrMasihAdaPenghalang = errors.New("services: klaim belum boleh ditutup")

// ErrTahapTidakMenutup - penutupan dari tahap yang layarnya tidak
// menawarkannya.
var ErrTahapTidakMenutup = errors.New("services: tahap ini tidak menawarkan Close Claim")

// GalatPenghalang membawa daftar penghalang menyeberang lapisan.
//
// ⚠️ Daftarnya ikut di dalam galat, bukan dikembalikan terpisah: handler yang
// harus memanggil dua fungsi untuk menyusun satu jawaban adalah handler yang
// suatu hari memanggil satu saja.
type GalatPenghalang struct{ Penghalang []PenghalangTampil }

func (g *GalatPenghalang) Error() string {
	return fmt.Sprintf("%v: %d peserta menahan", ErrMasihAdaPenghalang, len(g.Penghalang))
}

func (g *GalatPenghalang) Unwrap() error { return ErrMasihAdaPenghalang }

// Tutup menutup kasus: padanan `Call FinishAssignment` (b838).
//
// ⛔ KEPUTUSAN bb, dan penyimpangannya DINYATAKAN. `Flow/Register_Flow.xml`
// tidak punya konektor bernama `CloseClaim` - ia local action, dan perilaku
// mesin Pega untuk `FinishAssignment` dari local action tanpa konektor
// senama TIDAK DAPAT diturunkan dari ekspor (OQ-I). Yang ditiru adalah niat
// nyatanya: label `Close Claim`, konfirmasi b499 "Are you sure want to Close
// Claim?", dan gerbang "is not approved yet" yang menahan seluruh peserta
// yang belum diaksep. Ketiganya hanya masuk akal bila tombolnya MENUTUP.
//
// Urutan pemeriksaannya sengaja: identitas -> kasus terbuka -> tahap
// menawarkan -> peran pemegang tahap -> gerbang peserta. Gerbang peserta
// PALING AKHIR karena ia yang paling mahal (membaca seluruh klaim) dan
// paling tidak berguna bila pemanggilnya memang tidak berhak.
func (t *TutupKlaim) Tutup(ctx context.Context, pelaku inti.Pelaku, klaimID string,
	saat time.Time) error {

	if err := inti.WajibIdentitas(pelaku); err != nil {
		return err
	}
	if strings.TrimSpace(klaimID) == "" {
		return fmt.Errorf("%w: pengenal klaim wajib diisi", ErrWajibIsi)
	}
	if t == nil || t.svc == nil || !t.svc.PunyaDatabase() {
		return db.ErrTanpaOracle
	}

	baca := repository.NewKlaimLife(t.svc.DB())
	status, err := baca.StatusWorkKlaim(ctx, klaimID)
	if err != nil {
		return err
	}
	if kontrak.KasusTertutup(status) {
		return fmt.Errorf("%w: %s", kontrak.ErrKasusSudahTertutup, klaimID)
	}

	// ⛔ Tahap ASAL, sama seperti perpindahan: orang menutup kasus yang
	// SEDANG IA PEGANG. Kolom TAHAP menang; PY_POSITION cadangan baris lama.
	asal, err := tahapKasus(ctx, baca, klaimID)
	if err != nil {
		return err
	}
	if !models.TahapMenawarkanTutup(asal) {
		return fmt.Errorf("%w: %s (yang menawarkannya %v)",
			ErrTahapTidakMenutup, asal, models.TahapPenawarTutup())
	}
	peranTahap, ada := models.PeranPemegangTahap(asal)
	if !ada {
		return fmt.Errorf("%w: tahap %q", ErrPeranTahapBelumDiputuskan, asal)
	}
	if err := inti.WajibPeran(pelaku, peranTahap); err != nil {
		return err
	}

	hasil, err := t.Periksa(ctx, klaimID)
	if err != nil {
		return err
	}
	if !hasil.Boleh {
		return &GalatPenghalang{Penghalang: hasil.Penghalang}
	}

	return t.svc.DalamTransaksi(ctx, func(tx *db.Tx) error {
		if err := baca.TutupKasus(ctx, tx, klaimID, asal.String(),
			kontrak.StatusWorkSelesai, saat); err != nil {
			return err
		}
		return t.jejak.Rekam(ctx, tx, jejak.CatatanJejak{
			// ⛔ KlaimID, bukan AdjustmentID: yang tertutup KASUSNYA.
			KlaimID: klaimID,
			Dari:    asal.String(),
			// ⚠️ Tujuannya status kerja, bukan tahap: sesudah tutup TAHAP
			// kosong, dan jejak yang mencatat "ke: (kosong)" tidak dapat
			// dibaca siapa pun setahun kemudian.
			Ke:     kontrak.StatusWorkSelesai,
			AkunID: pelaku.AkunID,
			Waktu:  saat,
		})
	})
}

// PastikanKasusTerbuka menolak perubahan atas kasus yang sudah ditutup.
//
// ⛔ SATU pintu untuk SELURUH rute pengubah, dan itu disengaja. Aturan yang
// ditulis ulang di tujuh berkas adalah aturan yang suatu hari hanya ada di
// enam - dan yang ketujuh tidak akan berbunyi, sebab tiap berkas hijau
// sendirian. `TestSetiapLayananPengubahMemeriksaKasusTerbuka` menagih
// pemanggilannya.
//
// ⚠️ Diletakkan di `Service`, bukan di `TutupKlaim`: yang memanggilnya adalah
// layanan LAIN, dan menaruhnya pada tipe gerbang tutup akan membuat setiap
// layanan pengubah menyusun gerbang tutup hanya untuk bertanya.
func (s *Service) PastikanKasusTerbuka(ctx context.Context, klaimID string) error {
	if s == nil || !s.PunyaDatabase() {
		return db.ErrTanpaOracle
	}
	if strings.TrimSpace(klaimID) == "" {
		return fmt.Errorf("%w: pengenal klaim wajib diisi", ErrWajibIsi)
	}
	status, err := repository.NewKlaimLife(s.DB()).StatusWorkKlaim(ctx, klaimID)
	if err != nil {
		return err
	}
	if kontrak.KasusTertutup(status) {
		return fmt.Errorf("%w: %s tidak dapat diubah lagi", kontrak.ErrKasusSudahTertutup, klaimID)
	}
	return nil
}
