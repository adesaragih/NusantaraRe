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

	"nusantarare/internal/models"
	"nusantarare/internal/repository"
)

// TutupKlaim membungkus pemeriksaan gerbang tutup.
type TutupKlaim struct{ svc *Service }

// Tutup menyusun layanannya.
func (s *Service) Tutup() *TutupKlaim { return &TutupKlaim{svc: s} }

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
		return HasilPeriksaTutup{}, repository.ErrTanpaOracle
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
