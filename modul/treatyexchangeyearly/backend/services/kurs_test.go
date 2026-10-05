package services_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"nusantarare/modul/treatyexchangeyearly/backend/models"
	"nusantarare/modul/treatyexchangeyearly/backend/services"
	"nusantarare/modul/treatyexchangeyearly/backend/tiruan"
)

var ctx = context.Background()

var (
	penuh = services.Aktor{AkunID: "UJI-ADMIN", Penuh: true}
	lihat = services.Aktor{AkunID: "UJI-LIHAT"}
	jam   = time.Date(2026, 10, 5, 3, 4, 5, 600e6, time.UTC)
)

func layanan(g *tiruan.Gudang) *services.Layanan {
	return services.BaruLayanan(g, tiruan.Transaksi).DenganJam(func() time.Time { return jam })
}

func isian(tahun, cur, quarter string) models.Isian {
	return models.Isian{TreatyYear: tahun, IDCurrency: cur, StartDate: tahun + "-07-01", EndDate: "2027-06-30", ToIDR: "16500.50",
		ToUSD: "1", Quarter: quarter}
}

// Add: tanggal format Pega yang benar, kode mata uang dari view CURRENCY, ID situs + sequence (nomor terpakai
// dilompati), USERID akun, DATEIU = DATEIN = waktu simpan GMT.
func TestAdd(t *testing.T) {
	g := tiruan.Contoh()
	k, err := layanan(g).Simpan(ctx, penuh, isian("2026", "10002", ""))
	if err != nil {
		t.Fatal(err)
	}
	b := g.Baris[k.Kunci]
	if b.ID != "10115" || b.Currency != "SGD" || b.StartDate != "20260701T000000.000 GMT" || b.EndDate != "20270630T000000.000 GMT" ||
		b.Quarter != "0" || b.ToIDR != "16500.50" || b.UserID != "UJI-ADMIN" || b.DateIU != "20261005T030405.600 GMT" ||
		b.DateIn != b.DateIU {
		t.Errorf("tersimpan %+v", b)
	}
	if k.Mulai != "01-07-2026" || k.Akhir != "30-06-2027" || k.Diubah != "05-10-2026 10:04" || k.CurrencyName != "UJI SGD" {
		t.Errorf("tampilan %+v", k)
	}
}

// Edit lewat Kunci (ROWID): baris lain ber-ID sama tidak ikut berubah; tanggal yang hari-nya tidak diubah dibiarkan
// apa adanya (juga yang rusak / ber-jam); tanggal yang diubah ditulis format benar; DATEIN dan kode mata uang tetap.
func TestEditLewatKunci(t *testing.T) {
	g := tiruan.Contoh()
	l := layanan(g)
	if _, err := l.Simpan(ctx, penuh, models.Isian{Kunci: "AAA2", TreatyYear: "2025", IDCurrency: "10001", StartDate: "2025-07-01",
		EndDate: "2026-06-29", ToIDR: "16400", Quarter: "0"}); err != nil {
		t.Fatal(err)
	}
	b := g.Baris["AAA2"]
	if b.StartDate != "20250701T075400.000 GMT" || b.EndDate != "20260629T000000.000 GMT" || b.ToIDR != "16400" || b.ToUSD != "" {
		t.Errorf("edit %+v", b)
	}
	if g.Baris["AAA3"].ToIDR != "18900" {
		t.Error("baris lain ber-ID sama ikut berubah")
	}
	if _, err := l.Simpan(ctx, penuh, models.Isian{Kunci: "AAA1", TreatyYear: "2019", IDCurrency: "10001", StartDate: "2019-08-01",
		EndDate: "2020-06-30", ToIDR: "14500.00", ToUSD: "0", Quarter: "0"}); err != nil {
		t.Fatal(err)
	}
	if b := g.Baris["AAA1"]; b.StartDate != "20190801T00000.000 GMT" || b.DateIn != "20230713T103923.330 GMT" || b.UserID != "UJI-ADMIN" {
		t.Errorf("tanggal warisan tidak diubah harus dibiarkan %+v", b)
	}
	if _, err := l.Simpan(ctx, penuh, models.Isian{Kunci: "TIDAK-ADA", TreatyYear: "2019", IDCurrency: "10001", StartDate: "2019-08-01",
		EndDate: "2020-06-30", ToIDR: "1", Quarter: "0"}); !errors.Is(err, services.ErrTidakAda) {
		t.Errorf("kunci tak ada: %v", err)
	}
}

// Treaty Year + Currency + Quarter tidak boleh kembar; Quarter lain boleh; baris sendiri tidak dihitung.
func TestKembarDitolak(t *testing.T) {
	g := tiruan.Contoh()
	l := layanan(g)
	if _, err := l.Simpan(ctx, penuh, isian("2024", "10002", "0")); err == nil || !strings.Contains(err.Error(), "Treaty Year 2024, SGD, Quarter 0 already exists (ID 10050)") {
		t.Errorf("kembar: %v", err)
	}
	if _, err := l.Simpan(ctx, penuh, isian("2024", "10002", "2")); err != nil {
		t.Errorf("quarter lain: %v", err)
	}
	if _, err := l.Simpan(ctx, penuh, models.Isian{Kunci: "AAA4", TreatyYear: "2024", IDCurrency: "10002", StartDate: "2024-07-01",
		EndDate: "2025-06-30", ToIDR: "12100", Quarter: "0"}); err != nil {
		t.Errorf("baris sendiri dihitung kembar: %v", err)
	}
	if _, err := l.Simpan(ctx, penuh, models.Isian{Kunci: "AAA5", TreatyYear: "2024", IDCurrency: "10002", StartDate: "2024-07-01",
		EndDate: "2024-09-30", ToIDR: "11900", Quarter: "0"}); err == nil {
		t.Error("edit ke kombinasi baris lain harus ditolak")
	}
}

func TestValidasi(t *testing.T) {
	g := tiruan.Contoh()
	l := layanan(g)
	_, err := l.Simpan(ctx, penuh, models.Isian{TreatyYear: "25", StartDate: "2025-07-01", EndDate: "2025-06-30", ToIDR: "16.500,00",
		ToUSD: "x", Quarter: "5"})
	for _, mau := range []string{"Treaty Year must be 4 digits", "Currency is required", "End Date must not be before Start Date",
		"To IDR must be a number", "To USD must be a number", "Quarter must be 0"} {
		if !errors.Is(err, services.ErrMasukanTidakSah) || !strings.Contains(err.Error(), mau) {
			t.Errorf("tanpa %q: %v", mau, err)
		}
	}
	if _, err := l.Simpan(ctx, penuh, models.Isian{TreatyYear: "2026", IDCurrency: "10001", StartDate: "2026-02-30", EndDate: "",
		ToIDR: "1"}); err == nil || !strings.Contains(err.Error(), "Start Date is required") || !strings.Contains(err.Error(), "End Date is required") {
		t.Errorf("tanggal: %v", err)
	}
	if _, err := l.Simpan(ctx, penuh, isian("2026", "99999", "0")); err == nil || !strings.Contains(err.Error(), "Currency 99999 is not in") {
		t.Errorf("mata uang asing: %v", err)
	}
	if _, err := l.Simpan(ctx, lihat, isian("2026", "10001", "0")); !errors.Is(err, services.ErrDilarang) {
		t.Errorf("View only: %v", err)
	}
	if len(g.Baris) != 5 {
		t.Error("ditolak tetapi tertulis")
	}
}

func TestDaftarDanTampilan(t *testing.T) {
	l := layanan(tiruan.Contoh())
	d, err := l.Daftar(ctx, "", "")
	if err != nil || len(d) != 5 || d[0].TreatyYear != "2025" || d[4].TreatyYear != "2019" {
		t.Errorf("urutan %+v %v", d, err)
	}
	if d[4].Mulai != "01-08-2019" || d[4].Diubah != "13-07-2023 17:39" {
		t.Errorf("tampilan tanggal warisan %+v", d[4])
	}
	if d, _ := l.Daftar(ctx, "euro", "2025"); len(d) != 1 || d[0].Kunci != "AAA3" {
		t.Errorf("saring %+v", d)
	}
	m, th, err := l.Pilihan(ctx)
	if err != nil || len(m) != 3 || strings.Join(th, ",") != "2025,2024,2019" {
		t.Errorf("pilihan %v %v %v", m, th, err)
	}
	if s := services.TampilTanggal("lain"); s != "lain" {
		t.Errorf("bentuk lain %q", s)
	}
}
