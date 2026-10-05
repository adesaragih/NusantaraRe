package services_test

// Aturan Aggregate TANPA Oracle (gudang tiruan). Nilai harapan dihitung terpisah (Python decimal, presisi 38,
// setengah-ke-atas), bukan dengan rumus yang sama.

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/aggregate/backend/models"
	"nusantarare/modul/aggregate/backend/services"
	"nusantarare/modul/aggregate/backend/tiruan"
)

var (
	ctx   = context.Background()
	admin = inti.Pelaku{AkunID: "UJI-ADMIN"}
)

const kepalaCSV = "ASSESMENT ZONE;TREATY TYPE;COVERAGE;AS AT;UW YEAR;CEDING CODE;CEDING NAME;CURRENCY;TO USD;" +
	"NOR BUILDINGS;BUILDINGS;NOR STOCKS;STOCKS;NOR MACHINERY;MACHINERY;NOR OTHER CONTENTS;OTHER CONTENTS;" +
	"NOR CONSEQUENTIAL LOSS;CONSEQUENTIAL LOSS;NOR RESIDENTIAL;RESIDENTIAL;NOR COMMERCIAL;COMMERCIAL;NOR INDUSTRIAL;" +
	"INDUSTRIAL;NOR AGRICULTURE;AGRICULTURE;NOR MISCELLANEOUS;MISCELLANEOUS;NOR UTILITIES;UTILITIES;TOTAL NO OF RISK;" +
	"TOTAL IN AMOUNT;TOTAL IN AMOUNT IN USD;RNM SHARE;RNM VALUE;RNM VALUE IN USD;REMARK\n"

// barisCSV - satu baris CSV 38 kolom; isi = nomor kolom (1-38) -> nilai, kolom lain kosong.
func barisCSV(isi map[int]string) string {
	sel := make([]string, 38)
	for i, v := range isi {
		sel[i-1] = v
	}
	return strings.Join(sel, ";") + "\n"
}

// isiCSV - baris contoh CSV work owner (zona 1.1 diganti zona uji), baris EUR, baris kosong, dan baris zona tak
// dikenal ber-As At 31-08-2024.
var isiCSV = kepalaCSV +
	barisCSV(map[int]string{1: "1.1 UJI Banda", 2: "or", 3: "EQVET", 4: "30/09/2024", 5: "2023", 6: "UJI-ABAI",
		7: "UJI ABAI", 8: "IDR", 9: "9", 10: "5", 11: "5.881.261", 12: "6", 13: "319.674.758", 14: "3", 15: "234.749",
		16: "19", 17: "31.686.450", 18: "0", 32: "33", 33: "357.477.219"}) +
	barisCSV(map[int]string{1: "3.1 UJI Jakarta", 2: "QS", 3: "TSFWD", 4: "30/09/2024", 5: "2024", 8: "EUR",
		32: "78", 33: "18.438.686.215", 38: "catatan"}) +
	barisCSV(map[int]string{}) +
	barisCSV(map[int]string{1: "9.9 UJI Asing", 2: "SURPLUS", 3: "EQVET", 4: "31/08/2024", 5: "2023", 8: "IDR",
		32: "1", 33: "357.477.219"})

func layanan(g *tiruan.Gudang) *services.Layanan { return services.BaruLayanan(g, tiruan.Transaksi) }

func pratinjau(t *testing.T, g *tiruan.Gudang, treaty ...string) []models.Baris {
	t.Helper()
	p, err := layanan(g).Pratinjau(ctx, services.PermintaanPratinjau{CSV: isiCSV, MasterTreaty: treaty})
	if err != nil {
		t.Fatal(err)
	}
	return p.Baris
}

// UploadCSVAggregate_Act: zona lewat kode, ceding dari Master ID pertama, Treaty Year dari periode As At, kurs IDR
// 1/TOIDR baris USD (20 desimal), RNM Value 4 desimal, angka bertitik ribuan, Treaty Type huruf besar.
func TestPratinjauSatuMasterID(t *testing.T) {
	b := pratinjau(t, tiruan.Contoh(), "UJI-V1")
	if len(b) != 5 {
		t.Fatalf("baris %d, mau 3 data + 2 total (baris kosong dilewati)", len(b))
	}
	idr := b[0]
	mau := map[string]string{
		"ASSESMENT_ZONE": "1.1 UJI Zona Satu", "TREATY_TYPE": "OR", "AS_AT": "30-09-2024", "TREATYYEAR": "2024",
		"CEDING_CODE": "UJI-C1", "CEDING_NAME": "UJI CEDING SATU", "TO_USD": "0.00006060606060606061",
		"BUILDINGS": "5881261", "STOCKS": "319674758", "TOTAL_IN_AMOUNT": "357477219", "RNM_SHARE": "10",
		"M_TREATY_ID": "1000001", "RNM_VALUE": "35747721.9", "TOTAL_IN_AMOUNT_IN_USD": "21665.28600000000140824359",
		"RNM_VALUE_IN_USD": "2166.528600000000140824359", "CONSEQUENTIAL_LOSS": "",
	}
	for k, v := range mau {
		if idr[k] != v {
			t.Errorf("baris IDR %s = %q, mau %q", k, idr[k], v)
		}
	}
	eur := b[1]
	for k, v := range map[string]string{"TO_USD": "1.10235", "TOTAL_IN_AMOUNT_IN_USD": "20325885749.10525",
		"RNM_VALUE": "1843868621.5", "RNM_VALUE_IN_USD": "2032588574.910525", "REMARK": "catatan"} {
		if eur[k] != v {
			t.Errorf("baris EUR %s = %q, mau %q", k, eur[k], v)
		}
	}
	// Zona tak dikenal -> kosong; 31-08-2024 jatuh di dua periode -> periode pertama (2023), kurs IDR 2023.
	asing := b[2]
	if asing["ASSESMENT_ZONE"] != "" || asing["TREATYYEAR"] != "2023" || asing["TO_USD"] != "0.00006666666666666667" {
		t.Errorf("baris zona asing %+v", asing)
	}
}

// Beberapa Master ID: share dijumlah, ID digabung `;`, RNM Value dihitung ulang, RNM Value (USD) tetap dari share
// Master ID pertama (langkah 7.25 Pega tidak menyentuhnya).
func TestPratinjauBeberapaMasterID(t *testing.T) {
	b := pratinjau(t, tiruan.Contoh(), "UJI-V1", "UJI-V2")[0]
	for k, v := range map[string]string{"RNM_SHARE": "10.5", "M_TREATY_ID": "1000001;1000002", "CEDING_CODE": "UJI-C1",
		"RNM_VALUE": "37535107.995", "RNM_VALUE_IN_USD": "2166.528600000000140824359"} {
		if b[k] != v {
			t.Errorf("%s = %q, mau %q", k, b[k], v)
		}
	}
}

// Tanpa Master ID: ceding, share, dan Master Treaty kosong - Save menolaknya.
func TestPratinjauTanpaMasterID(t *testing.T) {
	b := pratinjau(t, tiruan.Contoh())[0]
	if b["CEDING_CODE"] != "" || b["CEDING_NAME"] != "" || b["RNM_SHARE"] != "" || b["M_TREATY_ID"] != "" {
		t.Errorf("baris tanpa Master ID %+v", b)
	}
}

// Baris "Total :" per mata uang, urutan kemunculan terakhir mata uang (IDR muncul lagi sesudah EUR).
func TestBarisTotalPerMataUang(t *testing.T) {
	b := pratinjau(t, tiruan.Contoh(), "UJI-V1")
	eur, idr := b[3], b[4]
	if eur["ASSESMENT_ZONE"] != models.ZonaTotal || eur["CURRENCY"] != "EUR" || idr["CURRENCY"] != "IDR" {
		t.Fatalf("urutan total %q %q", eur["CURRENCY"], idr["CURRENCY"])
	}
	if idr["TOTAL_IN_AMOUNT"] != "714954438" || idr["TOTAL_NO_OF_RISK"] != "34" || eur["TOTAL_NO_OF_RISK"] != "78" {
		t.Errorf("jumlah total IDR %+v EUR %+v", idr, eur)
	}
}

func TestPratinjauDitolak(t *testing.T) {
	banyak := make([]string, services.MaksMasterTreaty+1)
	for _, k := range []struct {
		req   services.PermintaanPratinjau
		pesan string
	}{
		{services.PermintaanPratinjau{CSV: kepalaCSV}, "no data rows"},
		{services.PermintaanPratinjau{CSV: isiCSV, MasterTreaty: banyak}, "at most"},
	} {
		_, err := layanan(tiruan.Contoh()).Pratinjau(ctx, k.req)
		if !errors.Is(err, services.ErrMasukanTidakSah) || !strings.Contains(err.Error(), k.pesan) {
			t.Errorf("%v, mau %q", err, k.pesan)
		}
	}
	_, err := layanan(tiruan.Contoh()).Pratinjau(ctx, services.PermintaanPratinjau{CSV: isiCSV, MasterTreaty: []string{"UJI-TIDAK"}})
	if !errors.Is(err, services.ErrMasukanTidakSah) {
		t.Errorf("Master ID asing: %v", err)
	}
}

// SaveAggregate_Act: pesan Pega dikumpulkan dengan nomor urut grid; baris Total dilewati.
func TestSimpanDitolakPesanPega(t *testing.T) {
	g := tiruan.Contoh()
	b := pratinjau(t, g)
	_, err := layanan(g).Simpan(ctx, admin, services.PermintaanSimpan{Baris: b})
	var gs services.GalatSimpan
	if !errors.As(err, &gs) {
		t.Fatalf("galat %v", err)
	}
	mau := []string{
		"RNM Share in List 1 Can't Null, Please Check Your Data",
		"Ceding Name in List 1 Not Found, Please Check CedingID",
		"RNM Share in List 2 Can't Null, Please Check Your Data",
		"Ceding Name in List 2 Not Found, Please Check CedingID",
		"Assesment Zone in list 3 Not Found",
		"RNM Share in List 3 Can't Null, Please Check Your Data",
		"Ceding Name in List 3 Not Found, Please Check CedingID",
	}
	if !reflect.DeepEqual(gs.Pesan, mau) {
		t.Errorf("pesan\n%q\nmau\n%q", gs.Pesan, mau)
	}
	b[0]["TREATY_TYPE"], b[0]["TO_USD"] = "XL", ""
	_, err = layanan(g).Simpan(ctx, admin, services.PermintaanSimpan{Baris: b[:1]})
	if !errors.As(err, &gs) || !strings.Contains(err.Error(), "To USD in list 1 cannot be empty") ||
		!strings.Contains(err.Error(), "Treaty Type in List 1 Can Only be Filled With OR, QS or SURPLUS") {
		t.Errorf("galat %v", err)
	}
	if len(g.Baris) != 1 {
		t.Error("baris ditulis padahal ditolak")
	}
}

// Save berhasil: baris Total dilewati, ID AGG-n melewati nomor terpakai, satu TANGGAL_INPUT, USER_INPUT = pelaku.
func TestSimpanBerhasil(t *testing.T) {
	g := tiruan.Contoh()
	b := pratinjau(t, g, "UJI-V1")
	data := []models.Baris{b[0], b[1], b[3], b[4]}
	h, err := layanan(g).Simpan(ctx, admin, services.PermintaanSimpan{Baris: data})
	if err != nil {
		t.Fatal(err)
	}
	if h.Disimpan != 2 || h.Pesan != "Successfully Saved" || len(g.Baris) != 3 {
		t.Fatalf("hasil %+v baris %d", h, len(g.Baris))
	}
	baru := g.Baris[1:]
	if baru[0][models.KolomID] != "AGG-501" || baru[1][models.KolomID] != "AGG-502" {
		t.Errorf("ID %s %s, mau AGG-501 AGG-502 (AGG-500 terpakai)", baru[0][models.KolomID], baru[1][models.KolomID])
	}
	if baru[0][models.KolomTanggalInput] != baru[1][models.KolomTanggalInput] || baru[0][models.KolomUserInput] != "UJI-ADMIN" {
		t.Errorf("jejak %+v", baru[0])
	}
}

func TestSimpanIsiTakSah(t *testing.T) {
	g := tiruan.Contoh()
	b := pratinjau(t, g, "UJI-V1")
	b[1]["UW_YEAR"] = "UJI-TERLALU-PANJANG"
	_, err := layanan(g).Simpan(ctx, admin, services.PermintaanSimpan{Baris: b[:2]})
	if err == nil || err.Error() != "Failed to Save, Error in row 2" {
		t.Errorf("galat %v", err)
	}
	_, err = layanan(g).Simpan(ctx, admin, services.PermintaanSimpan{Baris: []models.Baris{{"UJI_KOLOM": "x"}}})
	if !errors.Is(err, services.ErrMasukanTidakSah) || !strings.Contains(err.Error(), "unknown column") {
		t.Errorf("kolom asing %v", err)
	}
	if _, err := layanan(g).Simpan(ctx, inti.Pelaku{}, services.PermintaanSimpan{Baris: b}); !errors.Is(err, services.ErrTanpaPelaku) {
		t.Errorf("tanpa pelaku %v", err)
	}
}

// GetMasterIDAgg_Act: kata kosong = nol baris; CEDING memuat kata; Proportional non-PROPERTY tidak; ganda dibuang
// dengan menyisakan kemunculan terakhir.
func TestCariTreaty(t *testing.T) {
	l := layanan(tiruan.Contoh())
	if d, _ := l.CariTreaty(ctx, " "); len(d) != 0 {
		t.Errorf("kata kosong %+v", d)
	}
	d, err := l.CariTreaty(ctx, "satu")
	if err != nil {
		t.Fatal(err)
	}
	if len(d) != 1 || d[0].ID != "UJI-V3" {
		t.Errorf("hasil %+v, mau UJI-V3 saja", d)
	}
}

// Daftar, rincian, hapus, ringkasan.
func TestDaftarRincianHapusRingkasan(t *testing.T) {
	g := tiruan.Contoh()
	l := layanan(g)
	h, err := l.Daftar(ctx, "lama", 0, 0)
	if err != nil || h.Total != 1 || h.Daftar[0].JumlahBaris != 1 {
		t.Fatalf("daftar %+v %v", h, err)
	}
	k := h.Daftar[0].Kunci
	if k.TanggalInput != "01-03-2025" || k.AsAt != "30-09-2024" {
		t.Errorf("kunci %+v", k)
	}
	// Chart tanpa As At (work owner 04-10-2026): seluruh As At dijumlah.
	lain := models.Baris{}
	for k, v := range g.Baris[0] {
		lain[k] = v
	}
	lain[models.KolomID], lain["AS_AT"], lain["RNM_VALUE_IN_USD"] = "AGG-501", "31-12-2024", "5"
	g.Baris = append(g.Baris, lain)
	r, err := l.Ringkasan(ctx)
	mau := models.IrisanRingkasan{CedingCode: "UJI-C9", CedingName: "UJI CEDING LAMA", TreatyType: "OR", Coverage: "EQVET", RnmValueInUSD: "12"}
	if err != nil || len(r.Irisan) != 1 || r.Irisan[0] != mau {
		t.Errorf("ringkasan %+v %v", r, err)
	}
	g.Baris = g.Baris[:1]
	if b, err := l.Rincian(ctx, k); err != nil || len(b) != 1 {
		t.Errorf("rincian %v %v", b, err)
	}
	if _, err := l.Rincian(ctx, models.Kunci{AsAt: "31/12/2024"}); !errors.Is(err, services.ErrMasukanTidakSah) {
		t.Errorf("tanggal salah bentuk %v", err)
	}
	if n, err := l.Hapus(ctx, admin, k); err != nil || n != 1 || len(g.Baris) != 0 {
		t.Errorf("hapus %d %v", n, err)
	}
	if _, err := l.Hapus(ctx, admin, k); !errors.Is(err, services.ErrTidakAda) {
		t.Errorf("hapus kedua %v", err)
	}
}

func TestAngkaDanTanggalCSV(t *testing.T) {
	for masuk, mau := range map[string]string{"5.881.261": "5881261", "52.375.772.698,01": "52375772698.01", "": ""} {
		if got := services.AngkaCSV(masuk); got != mau {
			t.Errorf("AngkaCSV(%q) = %q, mau %q", masuk, got, mau)
		}
	}
	if got := services.AsAtCSV("31/12/2024"); got != "31-12-2024" {
		t.Errorf("AsAtCSV = %q", got)
	}
	if got := services.AsAtCSV("2024-12-31"); got != "2024-12-31" {
		t.Errorf("bentuk lain = %q, mau apa adanya", got)
	}
}
