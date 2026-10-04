package services_test

// Uji seam services tiket 16, 17, 18, 19 — identitas kontrak.

import (
	"context"
	"errors"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatyin/backend/models"
	"nusantarare/modul/treatyin/backend/services"
)

func kontrakBerlaku() models.KontrakDenganVersi {
	return models.KontrakDenganVersi{Kontrak: models.Kontrak{
		ID: 7, IDCedant: 1, IDAsalBisnis: 2,
		SifatProporsi: models.SifatNonProporsional,
		TanggalMulai:  "2026-01-01", TanggalBerakhir: "2026-12-31",
	}}
}

// ── Tiket 16 ──────────────────────────────────────────────────────────────

// Uji POSITIF, dan ia POKOK tiket 16: kunci alami kembar TETAP TERSIMPAN.
// Menolaknya berarti menolak data yang sah — pembaruan tahunan yang dicatat
// ulang, atau dua perjanjian terpisah dengan cedant yang sama.
func TestKunciAlamiGandaTetapTersimpan(t *testing.T) {
	g := &gudangTiruan{serupa: []int64{41, 42}}
	l := services.LayananDengan(g)

	h, err := l.BuatKontrak(context.Background(), pelakuAda, masukanSah())

	if err != nil {
		t.Fatalf("kontrak berkunci alami kembar DITOLAK; ADR-0040 §2 melarang itu: %v", err)
	}
	if h.IDKontrak == 0 {
		t.Error("tidak tersimpan")
	}
	if len(h.Peringatan) != 1 {
		t.Fatalf("mau 1 peringatan, dapat %d", len(h.Peringatan))
	}
}

// Peringatan WAJIB menyebut pembandingnya — tanpa itu ia tidak dapat
// ditindaklanjuti, dan daftar periksa tiket 16 menolaknya.
func TestPeringatanMenyebutKontrakPembandingnya(t *testing.T) {
	l := services.LayananDengan(&gudangTiruan{serupa: []int64{41, 42}})

	h, err := l.BuatKontrak(context.Background(), pelakuAda, masukanSah())
	if err != nil {
		t.Fatal(err)
	}
	p := h.Peringatan[0]
	if p.Kode != services.KodePeringatanKunciAlami {
		t.Errorf("kode peringatan %q", p.Kode)
	}
	if len(p.Pembanding) != 2 || p.Pembanding[0] != 41 || p.Pembanding[1] != 42 {
		t.Errorf("pembanding tidak lengkap: %v", p.Pembanding)
	}
	for _, mau := range []string{"41", "42"} {
		if !strings.Contains(p.Pesan, mau) {
			t.Errorf("pesan tidak menyebut kontrak %s: %q", mau, p.Pesan)
		}
	}
}

// Kontrak yang TIDAK kembar tidak memperoleh peringatan — peringatan yang
// selalu muncul adalah peringatan yang berhenti dibaca.
func TestKontrakTidakKembarTanpaPeringatan(t *testing.T) {
	l := services.LayananDengan(&gudangTiruan{})

	h, err := l.BuatKontrak(context.Background(), pelakuAda, masukanSah())
	if err != nil {
		t.Fatal(err)
	}
	if len(h.Peringatan) != 0 {
		t.Errorf("mau nol peringatan, dapat %d", len(h.Peringatan))
	}
}

// ── Tiket 17 ──────────────────────────────────────────────────────────────

func TestCariNomorWarisanMengembalikanKontraknya(t *testing.T) {
	g := &gudangTiruan{warisan: []models.Kontrak{{ID: 7, NomorKontrakWarisan: "TRI-2019-0007"}}}
	l := services.LayananDengan(g)

	k, err := l.CariLewatNomorWarisan(context.Background(), pelakuAda, "TRI-2019-0007")

	if err != nil {
		t.Fatalf("pencarian gagal: %v", err)
	}
	if len(k) != 1 || k[0].ID != 7 {
		t.Errorf("hasil tidak sesuai: %+v", k)
	}
	if len(g.cariWarisan) != 1 || g.cariWarisan[0] != "TRI-2019-0007" {
		t.Errorf("nomor yang diteruskan ke gudang: %v", g.cariWarisan)
	}
}

// Nomor yang tidak ada → hasil KOSONG yang dinyatakan, bukan galat dan bukan
// daftar penuh. Orang yang memegang selembar kertas harus tahu nomornya tidak
// ada — bukan menerima semua kontrak.
func TestCariNomorWarisanTidakAdaMengembalikanKosong(t *testing.T) {
	l := services.LayananDengan(&gudangTiruan{warisan: []models.Kontrak{}})

	k, err := l.CariLewatNomorWarisan(context.Background(), pelakuAda, "TIDAK-ADA")

	if err != nil {
		t.Fatalf("nomor tidak ada menghasilkan galat: %v", err)
	}
	if len(k) != 0 {
		t.Errorf("mau kosong, dapat %d", len(k))
	}
}

// Nomor kosong BUKAN "cari semua": ia permintaan yang belum lengkap.
func TestCariNomorWarisanKosongDitolak(t *testing.T) {
	g := &gudangTiruan{}
	l := services.LayananDengan(g)

	for _, n := range []string{"", "   "} {
		if _, err := l.CariLewatNomorWarisan(context.Background(), pelakuAda, n); !errors.Is(err, services.ErrMasukanTidakSah) {
			t.Errorf("nomor %q: mau ErrMasukanTidakSah, dapat %v", n, err)
		}
	}
	if len(g.cariWarisan) != 0 {
		t.Error("gudang tersentuh untuk nomor kosong")
	}
}

func TestCariNomorWarisanMenolakTanpaIdentitas(t *testing.T) {
	g := &gudangTiruan{}
	l := services.LayananDengan(g)

	if _, err := l.CariLewatNomorWarisan(context.Background(), inti.Pelaku{}, "X"); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Fatalf("mau ErrTanpaIdentitas, dapat %v", err)
	}
	if len(g.cariWarisan) != 0 {
		t.Error("gudang tersentuh tanpa identitas")
	}
}

// ── Tiket 18 ──────────────────────────────────────────────────────────────

func ubahSah() services.MasukanUbahKontrak {
	k := kontrakBerlaku().Kontrak
	return services.MasukanUbahKontrak{
		NomorKontrakWarisan: "TRI-BARU",
		IDCedant:            k.IDCedant,
		IDAsalBisnis:        k.IDAsalBisnis,
		SifatProporsi:       k.SifatProporsi,
		TanggalMulai:        k.TanggalMulai,
		TanggalBerakhir:     k.TanggalBerakhir,
	}
}

// INV-19: KELIMA ruas beku, satu per satu — lima uji, bukan satu. Daftar
// periksa tiket 18 menuntutnya begitu, dan satu uji gabungan akan hijau walau
// hanya satu ruas yang benar-benar dijaga.
func TestRuasBekuTidakDapatDiubah(t *testing.T) {
	kasus := []struct {
		ruas  string
		rusak func(*services.MasukanUbahKontrak)
	}{
		{"idCedant", func(m *services.MasukanUbahKontrak) { m.IDCedant = 99 }},
		{"idAsalBisnis", func(m *services.MasukanUbahKontrak) { m.IDAsalBisnis = 99 }},
		{"sifatProporsi", func(m *services.MasukanUbahKontrak) { m.SifatProporsi = models.SifatProporsional }},
		{"tanggalMulai", func(m *services.MasukanUbahKontrak) { m.TanggalMulai = "2025-01-01" }},
		{"tanggalBerakhir", func(m *services.MasukanUbahKontrak) { m.TanggalBerakhir = "2027-12-31" }},
	}
	for _, k := range kasus {
		g := &gudangTiruan{kontrak: kontrakBerlaku()}
		l := services.LayananDengan(g)
		m := ubahSah()
		k.rusak(&m)

		err := l.PerbaruiKontrak(context.Background(), pelakuAda, 7, m)

		if !errors.Is(err, services.ErrRuasBekuBerubah) {
			t.Errorf("%s: mau ErrRuasBekuBerubah, dapat %v", k.ruas, err)
			continue
		}
		if !strings.Contains(err.Error(), k.ruas) {
			t.Errorf("%s: pesan tidak menyebut ruasnya: %v", k.ruas, err)
		}
		if len(g.diperbarui) != 0 {
			t.Errorf("%s: penolakan TETAP menulis ke gudang", k.ruas)
		}
	}
}

// Dua ruas beku sekaligus → pesannya menyebut KEDUANYA, bukan yang pertama.
// Pada borang berkolom puluhan, menebak berarti mencoba satu per satu.
func TestDuaRuasBekuBerubahDisebutKeduanya(t *testing.T) {
	l := services.LayananDengan(&gudangTiruan{kontrak: kontrakBerlaku()})
	m := ubahSah()
	m.IDCedant = 99
	m.TanggalMulai = "2025-01-01"

	err := l.PerbaruiKontrak(context.Background(), pelakuAda, 7, m)

	if err == nil {
		t.Fatal("diterima, mau ditolak")
	}
	for _, ruas := range []string{"idCedant", "tanggalMulai"} {
		if !strings.Contains(err.Error(), ruas) {
			t.Errorf("pesan tidak menyebut %s: %v", ruas, err)
		}
	}
}

// Uji POSITIF: ruas yang BUKAN kunci alami tetap dapat diubah. Gerbang yang
// menolak segalanya lulus kelima uji negatif di atas.
func TestRuasBukanKunciAlamiTetapDapatDiubah(t *testing.T) {
	g := &gudangTiruan{kontrak: kontrakBerlaku()}
	l := services.LayananDengan(g)

	if err := l.PerbaruiKontrak(context.Background(), pelakuAda, 7, ubahSah()); err != nil {
		t.Fatalf("mengubah nomor warisan DITOLAK: %v", err)
	}
	if len(g.diperbarui) != 1 {
		t.Fatalf("gudang menerima %d perubahan, mau 1", len(g.diperbarui))
	}
	if g.diperbarui[0].NomorKontrakWarisan != "TRI-BARU" {
		t.Errorf("nomor warisan tidak tersimpan: %+v", g.diperbarui[0])
	}
	// Lapisan beku ikut dikirim APA ADANYA, bukan dari masukan.
	if g.diperbarui[0].IDCedant != 1 {
		t.Errorf("lapisan beku ikut berubah: %+v", g.diperbarui[0])
	}
}

// ── Tiket 19 ──────────────────────────────────────────────────────────────

func versiSah() services.MasukanVersiTambahan {
	k := kontrakBerlaku().Kontrak
	return services.MasukanVersiTambahan{
		NomorUrutVersi: 2, KeadaanSiklusHidup: "DRAFT", NamaKontrak: "Versi kedua",
		IDCedant: k.IDCedant, IDAsalBisnis: k.IDAsalBisnis, SifatProporsi: k.SifatProporsi,
		TanggalMulai: k.TanggalMulai, TanggalBerakhir: k.TanggalBerakhir,
	}
}

// Kelima ruas beku, satu per satu — versi tidak dapat membawa lapisan beku
// yang menyimpang dari kontraknya.
func TestVersiMenyimpangDitolak(t *testing.T) {
	kasus := []struct {
		ruas  string
		rusak func(*services.MasukanVersiTambahan)
	}{
		{"idCedant", func(m *services.MasukanVersiTambahan) { m.IDCedant = 99 }},
		{"idAsalBisnis", func(m *services.MasukanVersiTambahan) { m.IDAsalBisnis = 99 }},
		{"sifatProporsi", func(m *services.MasukanVersiTambahan) { m.SifatProporsi = models.SifatProporsional }},
		{"tanggalMulai", func(m *services.MasukanVersiTambahan) { m.TanggalMulai = "2025-01-01" }},
		{"tanggalBerakhir", func(m *services.MasukanVersiTambahan) { m.TanggalBerakhir = "2027-12-31" }},
	}
	for _, k := range kasus {
		g := &gudangTiruan{kontrak: kontrakBerlaku()}
		l := services.LayananDengan(g)
		m := versiSah()
		k.rusak(&m)

		_, err := l.TambahVersi(context.Background(), pelakuAda, 7, m)

		if !errors.Is(err, services.ErrVersiMenyimpang) {
			t.Errorf("%s: mau ErrVersiMenyimpang, dapat %v", k.ruas, err)
			continue
		}
		// Pesannya menyebut ruas DAN nilai yang berlaku.
		if !strings.Contains(err.Error(), k.ruas) || !strings.Contains(err.Error(), "berlaku") {
			t.Errorf("%s: pesan tidak menyebut ruas dan nilai berlakunya: %v", k.ruas, err)
		}
		if len(g.versiBaru) != 0 {
			t.Errorf("%s: versi menyimpang TETAP ditulis", k.ruas)
		}
	}
}

// Uji POSITIF: versi yang lapisan bekunya SAMA diterima.
func TestVersiSelarasDiterima(t *testing.T) {
	g := &gudangTiruan{kontrak: kontrakBerlaku()}
	l := services.LayananDengan(g)

	id, err := l.TambahVersi(context.Background(), pelakuAda, 7, versiSah())

	if err != nil {
		t.Fatalf("versi selaras DITOLAK: %v", err)
	}
	if id == 0 {
		t.Error("pengenal versi tidak dikembalikan")
	}
	if len(g.versiBaru) != 1 {
		t.Fatalf("gudang menerima %d versi, mau 1", len(g.versiBaru))
	}
	// INV-59: lapisan beku TIDAK disalin ke VERSI_KONTRAK. Yang dikirim ke
	// gudang hanya nomor urut, keadaan, dan nama.
	v := g.versiBaru[0]
	if v.NomorUrutVersi == nil || *v.NomorUrutVersi != 2 || v.NamaKontrak != "Versi kedua" {
		t.Errorf("bentuk versi tidak sesuai: %+v", v)
	}
}

func TestTambahVersiMenolakMasukanTidakLengkap(t *testing.T) {
	l := services.LayananDengan(&gudangTiruan{kontrak: kontrakBerlaku()})
	m := versiSah()
	m.NamaKontrak = "  "
	m.NomorUrutVersi = 0

	_, err := l.TambahVersi(context.Background(), pelakuAda, 7, m)

	if !errors.Is(err, services.ErrMasukanTidakSah) {
		t.Fatalf("mau ErrMasukanTidakSah, dapat %v", err)
	}
	for _, ruas := range []string{"namaKontrak", "nomorUrutVersi"} {
		if !strings.Contains(err.Error(), ruas) {
			t.Errorf("pesan tidak menyebut %s: %v", ruas, err)
		}
	}
}
