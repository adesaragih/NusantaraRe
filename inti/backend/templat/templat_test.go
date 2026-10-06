package templat_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"nusantarare/inti/backend/templat"
)

var ctx = context.Background()

func slotUji() templat.Slot {
	return templat.Slot{
		Kode: "uji.premi.fire", Menu: "uji", Grup: "UJI Menu", Nama: "PREMIUM · FIRE",
		DipakaiDi: "UJI › Template", Ekstensi: ".csv", Pemisah: ';', JumlahKolom: 3,
		NamaUnduhan: "UJI - PREMIUM - FIRE.csv",
		Bawaan:      []byte("\xef\xbb\xbfCOB;\"POLICY\nNUMBER\";NOTE\r\nUJI-A;UJI-1;x\r\n"),
	}
}

func layanan(t *testing.T, g templat.Gudang) *templat.Layanan {
	t.Helper()
	k, err := templat.NewKatalog([]templat.Slot{slotUji()})
	if err != nil {
		t.Fatal(err)
	}
	return templat.NewLayanan(k, g)
}

func TestSlotDiperiksaSaatKatalogDirakit(t *testing.T) {
	salah := slotUji()
	salah.JumlahKolom = 4
	if _, err := templat.NewKatalog([]templat.Slot{salah}); err == nil || !strings.Contains(err.Error(), "3 kolom, slot 4") {
		t.Errorf("berkas bawaan beda jumlah kolom harus ditolak: %v", err)
	}
	if _, err := templat.NewKatalog([]templat.Slot{slotUji(), slotUji()}); err == nil {
		t.Error("kode ganda harus ditolak")
	}
	for _, ubah := range []func(*templat.Slot){
		func(s *templat.Slot) { s.Kode = "Uji" },
		func(s *templat.Slot) { s.Menu = "" },
		func(s *templat.Slot) { s.Ekstensi = "CSV" },
		func(s *templat.Slot) { s.NamaUnduhan = "uji.xlsx" },
		func(s *templat.Slot) { s.Bawaan = nil },
	} {
		s := slotUji()
		ubah(&s)
		if s.Periksa() == nil {
			t.Errorf("slot cacat diterima: %+v", s)
		}
	}
}

func TestKepalaJudulBerbarisBaruTetapSatuKolom(t *testing.T) {
	k, err := templat.Kepala(slotUji().Bawaan, ';')
	if err != nil || len(k) != 3 || k[1] != "POLICY\nNUMBER" {
		t.Errorf("kepala %q %v", k, err)
	}
}

func TestPeriksaMenolakDanMemperingatkan(t *testing.T) {
	l := layanan(t, templat.NewGudangMemori())
	h, err := l.Periksa(ctx, "uji.premi.fire", "baru.csv", []byte("COB;POLICY NUMBER;CATATAN\n"))
	if err != nil || len(h.Galat) != 0 || h.JumlahKolom != 3 {
		t.Fatalf("periksa %+v %v", h, err)
	}
	if len(h.Perbedaan) != 1 || h.Perbedaan[0] != (templat.Perbedaan{Kolom: 3, Lama: "NOTE", Baru: "CATATAN"}) {
		t.Errorf("judul berbaris baru dirapikan, hanya NOTE yang berbeda: %+v", h.Perbedaan)
	}
	for nama, isi := range map[string]string{
		"baru.xlsx":  "COB;POLICY;NOTE\n",
		"kolom.csv":  "COB;POLICY\n",
		"koma.csv":   "COB,POLICY,NOTE\n",
		"kosong.csv": "",
	} {
		h, _ := l.Periksa(ctx, "uji.premi.fire", nama, []byte(isi))
		if len(h.Galat) == 0 {
			t.Errorf("%s harus ditolak", nama)
		}
	}
	if _, err := l.Periksa(ctx, "uji.tidak.ada", "a.csv", nil); !errors.Is(err, templat.ErrSlotTidakAda) {
		t.Errorf("slot tak terdaftar: %v", err)
	}
}

func TestUnggahBerversiAktifkanDanUnduh(t *testing.T) {
	g := templat.NewGudangMemori()
	l := layanan(t, g)
	b, err := l.Unduh(ctx, "uji.premi.fire", -1)
	if err != nil || b.Nama != "UJI - PREMIUM - FIRE.csv" || string(b.Isi) != string(slotUji().Bawaan) || !strings.HasPrefix(b.Mime, "text/csv") {
		t.Fatalf("tanpa versi = bawaan: %+v %v", b, err)
	}
	if _, err := l.Unggah(ctx, "uji.premi.fire", "a.csv", []byte("A;B\n"), "", "UJI-IT"); !errors.Is(err, templat.ErrDitolak) {
		t.Errorf("jumlah kolom salah harus ditolak: %v", err)
	}
	v1, err := l.Unggah(ctx, "uji.premi.fire", `C:\fakepath\v1.csv`, []byte("COB;POLICY;NOTE\nUJI;1;x\n"), " awal ", "UJI-IT")
	if err != nil || v1 != 1 {
		t.Fatalf("unggah v1 %d %v", v1, err)
	}
	v2, _ := l.Unggah(ctx, "uji.premi.fire", "v2.csv", []byte("COB;POLICY;NOTE\nUJI;2;x\n"), "", "UJI-IT")
	b, _ = l.Unduh(ctx, "uji.premi.fire", -1)
	if v2 != 2 || !strings.Contains(string(b.Isi), "UJI;2") || b.Nama != "UJI - PREMIUM - FIRE.csv" {
		t.Errorf("versi aktif = v2, nama unduhan tetap nama slot: %+v", b)
	}
	r, _ := l.Riwayat(ctx, "uji.premi.fire")
	if len(r) != 2 || r[0].Versi != 2 || !r[0].Aktif || r[1].Aktif || r[1].NamaBerkas != "v1.csv" || r[1].Catatan != "awal" {
		t.Errorf("riwayat %+v", r)
	}
	if err := l.Aktifkan(ctx, "uji.premi.fire", 1); err != nil {
		t.Fatal(err)
	}
	b, _ = l.Unduh(ctx, "uji.premi.fire", -1)
	if !strings.Contains(string(b.Isi), "UJI;1") {
		t.Errorf("aktifkan v1: %s", b.Isi)
	}
	if err := l.Aktifkan(ctx, "uji.premi.fire", 0); err != nil {
		t.Fatal(err)
	}
	d, _ := l.Daftar(ctx)
	if len(d) != 1 || d[0].Aktif != nil || d[0].Pemisah != ";" {
		t.Errorf("kembali ke bawaan: %+v", d)
	}
	b, _ = l.Unduh(ctx, "uji.premi.fire", 2)
	if b.Nama != "UJI - PREMIUM - FIRE v2.csv" || !strings.Contains(string(b.Isi), "UJI;2") {
		t.Errorf("unduh versi tertentu: %+v", b)
	}
	if err := l.Aktifkan(ctx, "uji.premi.fire", 9); !errors.Is(err, templat.ErrVersiTidakAda) {
		t.Errorf("versi tak ada: %v", err)
	}
}

func TestTanpaDatabaseHanyaBawaan(t *testing.T) {
	l := layanan(t, nil)
	if b, err := l.Unduh(ctx, "uji.premi.fire", -1); err != nil || len(b.Isi) == 0 {
		t.Errorf("bawaan tetap bisa diunduh: %v", err)
	}
	if _, err := l.Unggah(ctx, "uji.premi.fire", "a.csv", []byte("COB;P;N\n"), "", "UJI-IT"); !errors.Is(err, templat.ErrTanpaDatabase) {
		t.Errorf("unggah tanpa database: %v", err)
	}
}
