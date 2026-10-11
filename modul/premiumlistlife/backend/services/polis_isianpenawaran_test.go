package services_test

// Form penawaran - tiket 01 bagian 3. TANPA Oracle: yang diuji urutan pagarnya.

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/premiumlistlife/backend/models"
	"nusantarare/modul/premiumlistlife/backend/services"
)

func TestFormPenawaranMenjagaPagarnyaBerurutan(t *testing.T) {
	f := services.New(nil).FormPenawaran()
	ctx := context.Background()
	w := time.Date(2026, 9, 30, 0, 0, 0, 0, time.UTC)
	p := inti.Pelaku{AkunID: "UJI-INPUTOR"}
	// Identitas lebih dulu.
	if _, err := f.Baca(ctx, inti.Pelaku{}, "NBLF-1"); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("baca anonim: %v", err)
	}
	if _, err := f.Simpan(ctx, inti.Pelaku{}, "NBLF-1", models.IsianPenawaran{}, w); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("simpan anonim: %v", err)
	}
	if _, err := f.CariCeding(ctx, inti.Pelaku{}, "a"); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("cari anonim: %v", err)
	}
	// Pengenal kosong dijawab "wajib diisi", bukan "ORACLE_DSN".
	if _, err := f.Simpan(ctx, p, "  ", models.IsianPenawaran{}, w); !errors.Is(err, galat.ErrPermintaanTidakSah) {
		t.Errorf("pengenal kosong: %v", err)
	}
	// Bentuk sah, tanpa Oracle: gagal terang.
	if _, err := f.Baca(ctx, p, "NBLF-1"); !errors.Is(err, db.ErrTanpaOracle) {
		t.Errorf("baca tanpa Oracle: %v", err)
	}
	if _, err := f.CariPemegangPolis(ctx, p, "a"); !errors.Is(err, db.ErrTanpaOracle) {
		t.Errorf("cari tanpa Oracle: %v", err)
	}
}

// Confirm di tahap penawaran digerbangi kelengkapan isian; Decline tidak.
func TestConfirmPenawaranDigerbangiKelengkapan(t *testing.T) {
	isi, err := os.ReadFile("polis_penawaran.go")
	if err != nil {
		t.Fatal(err)
	}
	teks := string(isi)
	if !strings.Contains(teks, "keadaan.Status == models.TahapPolisPenawaran && keputusan == models.KeputusanConfirm") ||
		!strings.Contains(teks, "p.periksaPenawaranLengkap(ctx, keadaan.ID)") {
		t.Error("Confirm tahap penawaran tidak memeriksa isian Input Offer")
	}
}

// TestSaveOfferBindClosedMenutupLewatJalurKeputusan - keputusan work owner
// 08-10-2026. Tanpa Oracle (paket ini tidak punya tiruan basis data): yang diuji
// bentuk sumbernya - SATU jalur tutup (`terapkanDalam`) di transaksi Save Offer,
// sesudah penawaran tersimpan; gerbang Confirm untuk Bind; tanpa SimpanPolis.
func TestSaveOfferBindClosedMenutupLewatJalurKeputusan(t *testing.T) {
	isi, err := os.ReadFile("polis_isianpenawaran.go")
	if err != nil {
		t.Fatal(err)
	}
	teks := string(isi)
	i := strings.Index(teks, "func (f *FormPenawaran) Simpan(")
	j := strings.Index(teks[i:], "\n}\n")
	badan := teks[i : i+j]
	urut := []string{
		"models.AkibatStatusOffer(flag, siap.Status)",
		"if tutup && akibat.SimpanPolis {",
		"if tutup && keputusan == models.KeputusanConfirm {",
		"models.KekuranganPenawaran(",
		"f.svc.DalamTransaksi(ctx",
		"repo.Simpan(ctx, tx, polisID, siap)",
		"NewJSONOfferLife(f.svc.DB()).Ganti(ctx, tx, polisID)",
		"repo.SisipSuggest(ctx, tx, polisID,",
		"if !tutup {",
		".terapkanDalam(ctx, tx, pelaku, keadaan, akibat, sebab, saat)",
	}
	n := 0
	for _, u := range urut {
		k := strings.Index(badan[n:], u)
		if k < 0 {
			t.Fatalf("Simpan tidak memuat (atau urutannya salah) %q", u)
		}
		n += k + len(u)
	}
	if strings.Count(badan, "DalamTransaksi(") != 1 {
		t.Error("Save Offer harus SATU transaksi (simpan + tutup)")
	}
	// Satu jalur tutup: TutupKasus hanya dipanggil terapkanDalam.
	if strings.Contains(badan, "TutupKasus(") {
		t.Error("Save Offer menutup kasus sendiri - pakai terapkanDalam")
	}
	// Handler memasang perekam jejak Oracle.
	h, err := os.ReadFile("../handlers/rute_penawaran.go")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(h), "svc.FormPenawaran().DenganJejak(jejak.PerekamJejakOracle(svc))") ||
		!strings.Contains(string(h), "hasil.StatusTutup = akibat.StatusWork") {
		t.Error("handler Save Offer tanpa jejak Oracle atau tanpa statusTutup")
	}
}

// TestConfirmPenawaranBendera0Ditolak - pengaman server (keputusan work owner
// 08-10-2026): Putuskan memeriksa ConfirmPenawaranTersedia SEBELUM gerbang lain
// dan sebelum terapkan; Save Offer (jalur Bind) tidak lewat Putuskan.
func TestConfirmPenawaranBendera0Ditolak(t *testing.T) {
	isi, err := os.ReadFile("polis_penawaran.go")
	if err != nil {
		t.Fatal(err)
	}
	teks := string(isi)
	i := strings.Index(teks, "func (p *Penawaran) PutuskanDenganNomor(")
	j := strings.Index(teks[i:], "\n}\n")
	badan := teks[i : i+j]
	iCek := strings.Index(badan, "models.ConfirmPenawaranTersedia(keadaan.Status, keputusan, flag)")
	iLengkap := strings.Index(badan, "p.periksaPenawaranLengkap(")
	iTerapkan := strings.Index(badan, "p.terapkan(")
	if iCek < 0 || iLengkap < 0 || iTerapkan < 0 || !(iCek < iLengkap && iLengkap < iTerapkan) {
		t.Fatalf("pengaman Confirm bendera 0 harus sebelum gerbang lain dan terapkan (%d %d %d)", iCek, iLengkap, iTerapkan)
	}
	f, err := os.ReadFile("polis_isianpenawaran.go")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(f), "PutuskanDenganNomor(") || strings.Contains(string(f), ".Putuskan(") {
		t.Error("Save Offer Bind tidak boleh lewat Putuskan (akan terkena pengaman)")
	}
	h, err := os.ReadFile("../handlers/rute_penawaran.go")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(h), "galat.Tulis(w, http.StatusConflict, models.PesanConfirmPenawaranTidakTersedia)") {
		t.Error("handler tidak menjawab 409 berpesan")
	}
	if services.ErrConfirmPenawaranTidakTersedia.Error() != "services: "+models.PesanConfirmPenawaranTidakTersedia {
		t.Error("galat")
	}
}
