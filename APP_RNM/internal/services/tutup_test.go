package services

import (
	"context"
	"testing"

	"nusantarare/internal/models"
)

// Uji penyusun masukan gerbang tutup.
//
// Yang dijaga: urutan yang dipakai pesan adalah POSISI DI LAYAR, dan status
// dibawa apa adanya dari peserta.

func TestBarisTutupDariMemakaiPosisiLayar(t *testing.T) {
	klaim := &models.Klaim{Peserta: []models.Peserta{
		{ID: "P-ZZ", NomorSertifikat: "UJI-009", KodeStatus: models.KodeAksep},
		{ID: "P-AA", NomorSertifikat: "UJI-001", KodeStatus: models.KodeOutstanding},
	}}
	got := BarisTutupDari(klaim)
	if len(got) != 2 {
		t.Fatalf("cacah = %d, mau 2", len(got))
	}
	// ⛔ Urutan 1 dan 2 mengikuti POSISI, bukan pengenal baris dan bukan
	// urutan abjad sertifikat. Pesannya dibaca orang yang sedang melihat
	// daftar itu; nomor yang tidak ada di layar tidak menolong siapa pun.
	if got[0].Urutan != 1 || got[1].Urutan != 2 {
		t.Errorf("urutan = %d,%d; mau 1,2", got[0].Urutan, got[1].Urutan)
	}
	if got[0].NomorSertifikat != "UJI-009" || got[1].NomorSertifikat != "UJI-001" {
		t.Errorf("sertifikat tidak mengikuti urutan daftar: %v", got)
	}
	if got[1].KodeStatus != models.KodeOutstanding {
		t.Errorf("status = %q, mau dibawa apa adanya", got[1].KodeStatus)
	}
}

func TestBarisTutupDariKlaimKosong(t *testing.T) {
	if got := BarisTutupDari(nil); got != nil {
		t.Errorf("klaim nil menghasilkan %v, mau nil", got)
	}
	if got := BarisTutupDari(&models.Klaim{}); len(got) != 0 {
		t.Errorf("klaim tanpa peserta menghasilkan %d baris", len(got))
	}
}

func TestPeriksaTutupTanpaDatabase(t *testing.T) {
	// Gerbang yang tidak dapat membaca pesertanya TIDAK menjawab "boleh".
	// Menjawab boleh saat data tidak terbaca adalah cara paling mudah
	// menutup klaim yang belum selesai.
	var s Service
	h, err := s.Tutup().Periksa(context.Background(), "K-1")
	if err == nil {
		t.Fatal("tanpa database seharusnya galat")
	}
	if h.Boleh {
		t.Error("menjawab boleh padahal pesertanya tidak terbaca")
	}
}
