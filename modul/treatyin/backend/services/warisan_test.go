package services_test

// Uji tiket 40 dan 41 — jalur baca warisan.
//
// ⛔ Uji negatif saja TIDAK memisahkan rancangan yang menyalin dari rancangan
// yang merujuk: tiket 40 menyebutnya sendiri, *"salinan dan rujukan memberi
// jawaban yang sama sampai salah satunya berubah"*. Karena itu uji positif di
// bawah MENGUBAH versi dasarnya di tengah jalan, lalu membaca ulang.

import (
	"context"
	"errors"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatyin/backend/models"
	"nusantarare/modul/treatyin/backend/services"
)

// ------------------------------------------------------------------ tiket 40

func TestVersiSebelumnyaMenolakTanpaIdentitasSebelumGudang(t *testing.T) {
	g := &gudangTiruan{}
	l := services.LayananDengan(g)

	if _, err := l.NilaiVersiSebelumnya(context.Background(), inti.Pelaku{}, 5); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Fatalf("mau ErrTanpaIdentitas, dapat %v", err)
	}
	if len(g.dibaca) != 0 {
		t.Errorf("gudang dibaca %d kali walau identitas tidak ada", len(g.dibaca))
	}
}

func TestVersiSebelumnyaMenolakPengenalTidakSah(t *testing.T) {
	g := &gudangTiruan{}
	l := services.LayananDengan(g)

	for _, id := range []int64{0, -1, -99} {
		_, err := l.NilaiVersiSebelumnya(context.Background(), pelakuAda, id)
		if !errors.Is(err, services.ErrMasukanTidakSah) {
			t.Errorf("id %d: mau ErrMasukanTidakSah, dapat %v", id, err)
		}
	}
	if len(g.dibaca) != 0 {
		t.Errorf("pengenal tidak sah tetap sampai ke gudang")
	}
}

// Versi PERTAMA: jawabannya pernyataan ketiadaan, bukan galat dan bukan nol.
//
// Ini jalur gagal tiket 40 yang paling mudah salah: mengembalikan galat
// memaksa pemanggil membedakan galat yang wajar dari yang tidak, dan
// mengembalikan struct kosong tanpa penanda membuat "tidak ada versi
// sebelumnya" tidak dapat dibedakan dari "ada, tapi namanya kosong".
func TestVersiPertamaMenjawabKetiadaanBukanGalat(t *testing.T) {
	g := &gudangTiruan{versiDasar: nil}
	l := services.LayananDengan(g)

	hasil, err := l.NilaiVersiSebelumnya(context.Background(), pelakuAda, 5)
	if err != nil {
		t.Fatalf("versi pertama menghasilkan galat %v; ia keadaan yang sah", err)
	}
	if hasil.Ada {
		t.Errorf("Ada bernilai true padahal tidak ada versi dasar")
	}
	if hasil.Versi != nil {
		t.Errorf("Versi terisi %v padahal tidak ada versi dasar", hasil.Versi)
	}
}

// Uji POSITIF tiket 40, dan ia yang membuktikan SELECT benar-benar SELECT.
//
// Gudang tiruan berperan sebagai basis data: nilainya DIUBAH di antara dua
// pembacaan. Rancangan yang menyalin akan menjawab nilai lama pada pembacaan
// kedua; rancangan yang merujuk menjawab nilai yang sudah dibetulkan.
func TestPembetulanPadaVersiDasarTerlihatPadaPembacaanBerikutnya(t *testing.T) {
	g := &gudangTiruan{versiDasar: &models.VersiKontrak{ID: 11, IDKontrak: 3, NamaKontrak: "Nama salah ketik"}}
	l := services.LayananDengan(g)

	pertama, err := l.NilaiVersiSebelumnya(context.Background(), pelakuAda, 12)
	if err != nil {
		t.Fatalf("pembacaan pertama gagal: %v", err)
	}
	if !pertama.Ada || pertama.Versi.NamaKontrak != "Nama salah ketik" {
		t.Fatalf("pembacaan pertama: %+v", pertama)
	}

	// Versi DASAR dibetulkan - bukan versi anaknya.
	g.versiDasar = &models.VersiKontrak{ID: 11, IDKontrak: 3, NamaKontrak: "Nama yang benar"}

	kedua, err := l.NilaiVersiSebelumnya(context.Background(), pelakuAda, 12)
	if err != nil {
		t.Fatalf("pembacaan kedua gagal: %v", err)
	}
	if kedua.Versi.NamaKontrak != "Nama yang benar" {
		t.Errorf("pembacaan kedua menjawab %q - nilainya tersalin, bukan dirujuk", kedua.Versi.NamaKontrak)
	}
}

// ------------------------------------------------------------------ tiket 41

func TestBentukLamaMenolakTanpaIdentitasSebelumGudang(t *testing.T) {
	g := &gudangTiruan{}
	l := services.LayananDengan(g)

	if _, err := l.IdentitasBentukLama(context.Background(), inti.Pelaku{}, 5); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Fatalf("mau ErrTanpaIdentitas, dapat %v", err)
	}
	if len(g.dibaca) != 0 {
		t.Errorf("gudang dibaca %d kali walau identitas tidak ada", len(g.dibaca))
	}
}

func TestBentukLamaMenolakPengenalTidakSahDanKontrakTidakAda(t *testing.T) {
	g := &gudangTiruan{}
	l := services.LayananDengan(g)
	if _, err := l.IdentitasBentukLama(context.Background(), pelakuAda, 0); !errors.Is(err, services.ErrMasukanTidakSah) {
		t.Errorf("pengenal 0: mau ErrMasukanTidakSah, dapat %v", err)
	}

	// Gudang menjawab kontrak kosong - pengenalnya tidak menunjuk apa pun.
	kosong := &gudangTiruan{}
	lk := services.LayananDengan(kosong)
	if _, err := lk.IdentitasBentukLama(context.Background(), pelakuAda, 404); !errors.Is(err, services.ErrKontrakTidakAda) {
		t.Errorf("kontrak tidak ada: mau ErrKontrakTidakAda, dapat %v", err)
	}
}

// Jalur gagal tiket 41: kontrak yang LAHIR DI SISTEM BARU tidak punya nomor
// lama, dan jawabannya pernyataan - bukan nomor karangan.
func TestKontrakBaruTidakDiberiNomorLamaKarangan(t *testing.T) {
	g := &gudangTiruan{kontrak: models.KontrakDenganVersi{Kontrak: models.Kontrak{
		ID: 7, NomorKontrakWarisan: "", IDCedant: 2, IDAsalBisnis: 3,
		SifatProporsi: models.SifatProporsional, TanggalMulai: "2026-01-01", TanggalBerakhir: "2026-12-31",
	}}}
	l := services.LayananDengan(g)

	hasil, err := l.IdentitasBentukLama(context.Background(), pelakuAda, 7)
	if err != nil {
		t.Fatalf("mau diterima, dapat %v", err)
	}
	if hasil.NomorLamaAda {
		t.Errorf("NomorLamaAda true untuk kontrak yang lahir di sistem baru")
	}
	if hasil.NomorLama != "" {
		t.Errorf("nomor lama dikarang: %q", hasil.NomorLama)
	}
	// Sisa identitasnya tetap terjawab - ketiadaan nomor lama bukan alasan
	// menolak seluruh pertanyaannya.
	if hasil.IDCedant != 2 || hasil.SifatProporsi != models.SifatProporsional {
		t.Errorf("identitas selain nomor lama ikut hilang: %+v", hasil)
	}
}

// Uji POSITIF tiket 41: kontrak warisan menjawab nomor lamanya apa adanya,
// beserta keempat ruas identitas lainnya.
func TestKontrakWarisanMenjawabIdentitasBentukLama(t *testing.T) {
	g := &gudangTiruan{kontrak: models.KontrakDenganVersi{Kontrak: models.Kontrak{
		ID: 7, NomorKontrakWarisan: "  TR-2019-0042  ", IDCedant: 2, IDAsalBisnis: 3,
		SifatProporsi: models.SifatNonProporsional, TanggalMulai: "2019-01-01", TanggalBerakhir: "2019-12-31",
	}}}
	l := services.LayananDengan(g)

	hasil, err := l.IdentitasBentukLama(context.Background(), pelakuAda, 7)
	if err != nil {
		t.Fatalf("mau diterima, dapat %v", err)
	}
	if !hasil.NomorLamaAda {
		t.Fatalf("NomorLamaAda false untuk kontrak warisan")
	}
	// Spasi di tepi dipangkas: nomor berspasi dan nomor tanpa spasi adalah
	// nomor yang sama bagi hilir, dan hilir tidak punya cara menyamakannya.
	if hasil.NomorLama != "TR-2019-0042" {
		t.Errorf("nomor lama %q, mau TR-2019-0042", hasil.NomorLama)
	}
	for _, p := range []struct {
		nama, dapat, mau string
	}{
		{"sifat proporsi", hasil.SifatProporsi, models.SifatNonProporsional},
		{"tanggal mulai", hasil.TanggalMulai, "2019-01-01"},
		{"tanggal berakhir", hasil.TanggalBerakhir, "2019-12-31"},
	} {
		if p.dapat != p.mau {
			t.Errorf("%s: %q, mau %q", p.nama, p.dapat, p.mau)
		}
	}
}

// ⛔ Penjaga INV-58/INV-59 yang dapat GAGAL: bentuk baca tiket 41 tidak boleh
// punya jalur tulis. Bila seseorang kelak menambahkan `SimpanIdentitasLama`
// atau sejenisnya ke seam gudang, uji ini menangkapnya.
func TestBentukLamaTidakPunyaJalurTulis(t *testing.T) {
	g := &gudangTiruan{kontrak: models.KontrakDenganVersi{Kontrak: models.Kontrak{
		ID: 7, IDCedant: 2, IDAsalBisnis: 3, SifatProporsi: models.SifatProporsional,
		TanggalMulai: "2026-01-01", TanggalBerakhir: "2026-12-31",
	}}}
	l := services.LayananDengan(g)

	if _, err := l.IdentitasBentukLama(context.Background(), pelakuAda, 7); err != nil {
		t.Fatalf("mau diterima, dapat %v", err)
	}
	// Nol penulisan: tidak ada kontrak dibuat, diperbarui, versi ditambah,
	// maupun pemulihan dicatat.
	if len(g.dibuat) != 0 || len(g.diperbarui) != 0 || len(g.versiBaru) != 0 || len(g.pemulihan) != 0 {
		t.Errorf("bentuk baca menulis ke gudang: dibuat=%d diperbarui=%d versiBaru=%d pemulihan=%d",
			len(g.dibuat), len(g.diperbarui), len(g.versiBaru), len(g.pemulihan))
	}
}
