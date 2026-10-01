package services_test

import (
	"context"
	"errors"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/endorsementlife/backend/models"
	"nusantarare/modul/endorsementlife/backend/services"
	"nusantarare/modul/endorsementlife/backend/tiruan"
)

// gudangKasusSimpan - kasus `EDMLF-1` bertipe `QR`: tiga peserta `Old` (dua
// IDR, satu USD) dan satu baris `New` hasil CSV.
func gudangKasusSimpan(edmType string) *tiruan.Gudang {
	g := tiruan.Baru()
	g.Polis["EDMLF-1"] = &tiruan.Polis{ID: "EDMLF-1", OldPolicyNo: "UJI-PL-1", EdmType: edmType, ProdKe: 2,
		Kepala: map[string]string{"TYPE": models.TypeQR}}
	for _, p := range []struct{ id, uang, gp, ded, status string }{
		{"UJI-D1", "IDR", "100.25", "10", models.StatusOld},
		{"UJI-D2", "IDR", "200", "", models.StatusOld},
		{"UJI-D3", "USD", "7.5", "0.5", models.StatusOld},
		{"UJI-D4", "IDR", "50", "5", models.StatusNew},
	} {
		g.Peserta = append(g.Peserta, &tiruan.Peserta{ID: p.id, PolisID: "EDMLF-1", PLNumber: "UJI-PL-1", EdmStatus: p.status,
			Nilai: map[string]string{"CURRENCY": p.uang, "GROSS_PREMIUM": p.gp, "DEDUCTION": p.ded, "COMM": "1"}})
	}
	return g
}

func peserta(g *tiruan.Gudang, id string) *tiruan.Peserta {
	for _, p := range g.Peserta {
		if p.ID == id {
			return p
		}
	}
	return nil
}

func rekapUang(r []map[string]string, uang string) map[string]string {
	for _, b := range r {
		if b["CURRENCY"] == uang {
			return b
		}
	}
	return nil
}

func TestSimpanPerubahanDataMenandaiYangDicentang(t *testing.T) {
	g := gudangKasusSimpan(models.EdmTypePerubahanData)
	l, j := layananJejak(g)
	h, err := l.Simpan(context.Background(), pelakuUji, "EDMLF-1", models.PilihanHapus{Pilih: []string{"UJI-D1", "UJI-D4"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(j.catatan) != 1 || j.catatan[0].Komentar != "Save" || j.catatan[0].AkunID != pelakuUji.AkunID || !j.catatan[0].Waktu.Equal(jamUji) {
		t.Errorf("jejak Save %+v", j.catatan)
	}
	// UJI-D4 `New`: tidak pernah dibalik walau dicentang (kotak mati di layar).
	if h.Ditandai != 1 || h.Status != models.StatusDelete {
		t.Fatalf("ditandai %d status %q", h.Ditandai, h.Status)
	}
	d1, d2, d4 := peserta(g, "UJI-D1"), peserta(g, "UJI-D2"), peserta(g, "UJI-D4")
	if d1.EdmStatus != models.StatusDelete || d1.Nilai["GROSS_PREMIUM"] != "-100.25" || d1.Nilai["DEDUCTION"] != "-10" {
		t.Errorf("UJI-D1 tidak dijurnal balik: %+v", d1)
	}
	// `COMM` bukan satu dari 32 kolom `SetPremi_EDM` 2.1 (R28): tetap.
	if d1.Nilai["COMM"] != "1" {
		t.Errorf("COMM ikut dibalik: %+v", d1)
	}
	if d2.EdmStatus != models.StatusOld || d2.Nilai["GROSS_PREMIUM"] != "200" || d4.EdmStatus != models.StatusNew || d4.Nilai["GROSS_PREMIUM"] != "50" {
		t.Errorf("peserta tak dicentang berubah: %+v %+v", d2, d4)
	}
	// IDR: PREMIUM -100.25 + 200 + 50; DEDUCTION -10 + 0 + 5; BALANCE QR = GP - DED - (…).
	idr := rekapUang(h.Rekap, "IDR")
	if len(h.Rekap) != 2 || idr["PREMIUM"] != "149.75" || idr["BALANCE"] != "154.75" || idr["COMMISSION"] != "3" {
		t.Errorf("rekap IDR %v (dari %d baris)", idr, len(h.Rekap))
	}
	if usd := rekapUang(h.Rekap, "USD"); usd["PREMIUM"] != "7.5" || usd["BALANCE"] != "7" {
		t.Errorf("rekap USD %v", usd)
	}
	k, err := l.BacaKasus(context.Background(), pelakuUji, "EDMLF-1")
	if err != nil || !k.SudahSimpan || len(k.Rekap) != 2 {
		t.Fatalf("sesudah simpan: %+v %v", k, err)
	}
	// ⛔ Simpan kedua ditolak: pembalikan kedua mustahil (R06).
	if _, err := l.Simpan(context.Background(), pelakuUji, "EDMLF-1", models.PilihanHapus{Pilih: []string{"UJI-D2"}}); !errors.Is(err, services.ErrSudahDisimpan) {
		t.Fatalf("simpan kedua: %v", err)
	}
	if d1.Nilai["GROSS_PREMIUM"] != "-100.25" || d2.EdmStatus != models.StatusOld {
		t.Error("simpan kedua mengubah peserta")
	}
}

func TestSimpanDeleteAllDenganPengecualian(t *testing.T) {
	g := gudangKasusSimpan(models.EdmTypePerubahanData)
	l, _ := layananJejak(g)
	h, err := l.Simpan(context.Background(), pelakuUji, "EDMLF-1", models.PilihanHapus{Semua: true, Kecuali: []string{"UJI-D2"}})
	if err != nil || h.Ditandai != 2 {
		t.Fatalf("%+v %v", h, err)
	}
	if peserta(g, "UJI-D2").EdmStatus != models.StatusOld || peserta(g, "UJI-D3").EdmStatus != models.StatusDelete {
		t.Error("pengecualian DELETE ALL tidak dihormati")
	}
}

func TestSimpanTanpaCentangTetapMenghitungRekap(t *testing.T) {
	g := gudangKasusSimpan(models.EdmTypePerubahanData)
	l, _ := layananJejak(g)
	h, err := l.Simpan(context.Background(), pelakuUji, "EDMLF-1", models.PilihanHapus{})
	if err != nil || h.Ditandai != 0 || len(h.Rekap) != 2 {
		t.Fatalf("%+v %v", h, err)
	}
	if rekapUang(h.Rekap, "IDR")["PREMIUM"] != "350.25" {
		t.Errorf("rekap %v", h.Rekap)
	}
}

func TestSimpanBatalMembalikSeluruhPesertaOld(t *testing.T) {
	g := gudangKasusSimpan(models.EdmTypeBatal)
	l, _ := layananJejak(g)
	// Pilihan layar diabaikan: grid Batal tanpa kotak centang.
	h, err := l.Simpan(context.Background(), pelakuUji, "EDMLF-1", models.PilihanHapus{Pilih: []string{"UJI-D1"}})
	if err != nil || h.Ditandai != 3 || h.Status != models.StatusBatal {
		t.Fatalf("%+v %v", h, err)
	}
	for _, id := range []string{"UJI-D1", "UJI-D2", "UJI-D3"} {
		if peserta(g, id).EdmStatus != models.StatusBatal {
			t.Errorf("%s tidak Batal", id)
		}
	}
	if rekapUang(h.Rekap, "USD")["PREMIUM"] != "-7.5" {
		t.Errorf("rekap Batal %v", h.Rekap)
	}
}

func TestSimpanDitolak(t *testing.T) {
	ctx := context.Background()
	tertutup := gudangKasusSimpan(models.EdmTypePerubahanData)
	tertutup.Polis["EDMLF-1"].Status = models.StatusKasusSelesai
	kosong := gudangKasusSimpan(models.EdmTypePerubahanData)
	kosong.Peserta = nil
	rusak := gudangKasusSimpan(models.EdmTypePerubahanData)
	rusak.Peserta[0].Nilai["GROSS_PREMIUM"] = "1,5"
	for _, k := range []struct {
		nama  string
		g     *tiruan.Gudang
		id    string
		galat error
	}{
		{"tertutup", tertutup, "EDMLF-1", services.ErrKasusTertutup},
		{"tanpa peserta", kosong, "EDMLF-1", services.ErrTanpaPeserta},
		{"bukan kasus EDM", gudangKasusSimpan("1"), "NBLF-1", services.ErrKasusTidakAda},
		{"tidak ada", gudangKasusSimpan("1"), "EDMLF-9", services.ErrKasusTidakAda},
	} {
		l, _ := layananJejak(k.g)
		if _, err := l.Simpan(ctx, pelakuUji, k.id, models.PilihanHapus{Semua: true}); !errors.Is(err, k.galat) {
			t.Errorf("%s: %v, mau %v", k.nama, err, k.galat)
		}
	}
	l, _ := layananJejak(rusak)
	if _, err := l.Simpan(ctx, pelakuUji, "EDMLF-1", models.PilihanHapus{Semua: true}); err == nil {
		t.Error("nilai uang rusak diterima")
	}
	if _, err := l.Simpan(ctx, inti.Pelaku{}, "EDMLF-1", models.PilihanHapus{}); err == nil {
		t.Error("tanpa identitas diterima")
	}
}
