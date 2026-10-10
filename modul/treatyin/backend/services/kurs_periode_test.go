package services_test

import (
	"context"
	"testing"

	"nusantarare/modul/treatyin/backend/models"
	"nusantarare/modul/treatyin/backend/services"
)

// Kesembilan baris `TREATYEXCHANGEYEARLY` tahun 2026 di DEV_NUSARE2,
// 10 Oktober 2026 — `…Asli` sudah berbentuk tanggal WIB.
func kurs2026() []models.BarisKursWarisan {
	b := func(id, cur, dari, sampai string) models.BarisKursWarisan {
		return models.BarisKursWarisan{ID: id, MataUang: cur, BerlakuDariAsli: dari, BerlakuSampaiAsli: sampai}
	}
	return []models.BarisKursWarisan{
		b("10116", "USD", "20260701", "20270630"),
		b("10124", "IDR", "20261008", "20271008"),
		b("10125", "USD", "20261008", "20271008"),
		b("10130", "IDR", "20260401", "20270331"),
		b("10131", "USD", "20260401", "20270331"),
		b("10132", "IDR", "20260701", "20270630"),
		b("10133", "USD", "20260701", "20270630"),
		b("10141", "IDR", "20260101", "20261231"),
		b("10142", "USD", "20260101", "20261231"),
	}
}

func idKurs(bs []models.BarisKursWarisan) []string {
	out := make([]string, 0, len(bs))
	for _, b := range bs {
		out = append(out, b.ID)
	}
	return out
}

func samaUrut(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// ⛔ CACAT YANG DILAPORKAN 10 Oktober 2026: kontrak Non-Prop 1002309
// (01/01/2026–31/12/2026) diisi dua baris, dan sesudah Save grid
// menampilkan sembilan — seluruh kurs 2026 dari kontrak lain.
func TestKursHanyaPeriodeKontrak(t *testing.T) {
	got := idKurs(services.KursDalamPeriode(kurs2026(), "20260101", "20261231"))
	if mau := []string{"10141", "10142"}; !samaUrut(got, mau) {
		t.Fatalf("kontrak 01/01/2026–31/12/2026: %v, mau %v", got, mau)
	}
}

// ⚠️ Akibat yang diterima keputusan "Saring periode": kontrak dengan
// periode yang PERSIS sama tetap saling melihat barisnya.
func TestKursPeriodeSamaTetapBersama(t *testing.T) {
	got := idKurs(services.KursDalamPeriode(kurs2026(), "20260701", "20270630"))
	if mau := []string{"10116", "10132", "10133"}; !samaUrut(got, mau) {
		t.Fatalf("kontrak 01/07/2026–30/06/2027: %v, mau %v", got, mau)
	}
}

// Kurs yang berlaku di dalam periode — perubahan kurs per kuartal — tetap
// tampil; yang melewati Expiry atau mendahului Commencement tidak.
func TestKursDiDalamPeriodeTetapTampil(t *testing.T) {
	bs := []models.BarisKursWarisan{
		{ID: "a", BerlakuDariAsli: "20260101", BerlakuSampaiAsli: "20260331"},
		{ID: "b", BerlakuDariAsli: "20260401", BerlakuSampaiAsli: "20261231"},
		{ID: "c", BerlakuDariAsli: "20251201", BerlakuSampaiAsli: "20261231"},
		{ID: "d", BerlakuDariAsli: "20260101", BerlakuSampaiAsli: "20270101"},
	}
	got := idKurs(services.KursDalamPeriode(bs, "20260101", "20261231"))
	if mau := []string{"a", "b"}; !samaUrut(got, mau) {
		t.Fatalf("%v, mau %v", got, mau)
	}
}

// Batas kontrak yang kosong tidak menyaring sisinya; tanpa kedua batas
// seluruh baris tahun itu kembali apa adanya.
func TestKursTanpaPeriodeKontrak(t *testing.T) {
	if got := services.KursDalamPeriode(kurs2026(), "", ""); len(got) != 9 {
		t.Fatalf("tanpa periode: %d baris, mau 9", len(got))
	}
	got := idKurs(services.KursDalamPeriode(kurs2026(), "20261008", ""))
	if mau := []string{"10124", "10125"}; !samaUrut(got, mau) {
		t.Fatalf("hanya Commencement: %v, mau %v", got, mau)
	}
}

// Baris yang tanggalnya kosong tidak dapat ditempatkan di periode mana pun.
func TestKursTanpaTanggalTidakLolos(t *testing.T) {
	bs := []models.BarisKursWarisan{
		{ID: "x", BerlakuDariAsli: "", BerlakuSampaiAsli: "20261231"},
		{ID: "y", BerlakuDariAsli: "20260101", BerlakuSampaiAsli: "20261231"},
	}
	got := idKurs(services.KursDalamPeriode(bs, "20260101", "20261231"))
	if mau := []string{"y"}; !samaUrut(got, mau) {
		t.Fatalf("%v, mau %v", got, mau)
	}
}

// ⛔ Kontrak warisan 01/01–31/12 dengan kurs tahunnya berperiode 01/07–30/06:
// nol baris di dalam periode. Yang TUNTAS melihat seluruh tahun — grid
// tidak kosong (1.302 kontrak di DEV_NUSARE2, 10 Oktober 2026).
func TestKursKontrakTuntasNolDiPeriodeMelihatSeluruhTahun(t *testing.T) {
	bs := []models.BarisKursWarisan{
		{ID: "10112", MataUang: "USD", BerlakuDariAsli: "20240701", BerlakuSampaiAsli: "20250630"},
		{ID: "10097", MataUang: "EUR", BerlakuDariAsli: "20240701", BerlakuSampaiAsli: "20250630"},
	}
	if got := services.KursDalamPeriode(bs, "20240101", "20241231"); len(got) != 0 {
		t.Fatalf("saringan periode meloloskan %v", idKurs(got))
	}
	got := idKurs(services.KursGridKontrak(bs, "20240101", "20241231", models.StatusTuntas))
	if mau := []string{"10112", "10097"}; !samaUrut(got, mau) {
		t.Fatalf("kontrak tuntas: %v, mau %v", got, mau)
	}
}

// ⛔ Permintaan pemakai 10 Oktober 2026: Add, Delete, lalu Save harus
// menghasilkan jumlah baris yang tampil. Kontrak yang masih dapat di-Save
// TIDAK mendapat cadangan seluruh tahun — kalau dapat, menghapus baris
// terakhir periodenya akan memunculkan kurs periode lain sesudah Save.
func TestKursKontrakTerbukaTanpaCadangan(t *testing.T) {
	for _, status := range []string{"", models.StatusDecline, "Pending"} {
		if got := services.KursGridKontrak(kurs2026(), "20250101", "20251231", status); len(got) != 0 {
			t.Errorf("status %q: %v, mau kosong", status, idKurs(got))
		}
	}
	got := idKurs(services.KursGridKontrak(kurs2026(), "20260101", "20261231", ""))
	if mau := []string{"10141", "10142"}; !samaUrut(got, mau) {
		t.Fatalf("%v, mau %v", got, mau)
	}
}

// ⭐ Delete grid sampai ke rencana simpan — tetapi HANYA baris tahun ini
// yang berada di periode kontrak yang disimpan.
func TestDeleteKursHanyaBarisPeriodeKontrak(t *testing.T) {
	g := gudangSimpan()
	g.kepalaTreatyIn["1001001"]["Commencement"] = "20260101"
	g.kepalaTreatyIn["1001001"]["Termination"] = "20261231"
	g.kurs = []models.BarisKursWarisan{
		{ID: "10141", MataUang: "IDR", BerlakuDari: "20260101T000000.000 GMT", BerlakuSampai: "20261231T000000.000 GMT"},
		{ID: "10143", MataUang: "AUD", BerlakuDari: "20260101T000000.000 GMT", BerlakuSampai: "20261231T000000.000 GMT"},
		{ID: "10133", MataUang: "USD", BerlakuDari: "20260701T000000.000 GMT", BerlakuSampai: "20270630T000000.000 GMT"},
	}
	_, err := services.LayananDengan(g).SimpanKontrak(context.Background(), admin, services.MasukanSimpan{
		IDKontrak: "1001001",
		Kurs: &models.KursSimpan{
			Tahun: "2026",
			// 10143 milik periode ini; 10133 periode 01/07 (kontrak lain);
			// 99999 tidak ada di tahun ini.
			Hapus: []string{"10143", "10133", "99999"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	k := g.disimpan[0].Kurs
	if k == nil {
		t.Fatal("rencana tanpa kurs")
	}
	if mau := []string{"10143"}; !samaUrut(k.Hapus, mau) {
		t.Fatalf("hapus %v, mau %v", k.Hapus, mau)
	}
}

// ⛔ Laporan 10 Oktober 2026: kontrak 01/01/2025–01/01/2026 menyimpan kurs
// 01/01/2026–31/12/2026; barisnya tersimpan tetapi tersaring, grid kosong.
// Save kini menolaknya dengan pesan yang menyebut periodenya.
func TestSaveKursDiLuarPeriodeKontrakDitolak(t *testing.T) {
	g := gudangSimpan()
	g.kepalaTreatyIn["1001001"]["Commencement"] = "20250101"
	g.kepalaTreatyIn["1001001"]["Termination"] = "20260101"
	simpan := func(dari, sampai string) error {
		_, err := services.LayananDengan(g).SimpanKontrak(context.Background(), admin, services.MasukanSimpan{
			IDKontrak: "1001001",
			Kurs: &models.KursSimpan{Tahun: "2025", Baris: []models.BarisKursWarisan{
				{MataUang: "USD", MataUangID: "10001", NilaiKeIDR: "16000", BerlakuDariAsli: dari, BerlakuSampaiAsli: sampai},
			}},
		})
		return err
	}
	if err := simpan("20260101", "20261231"); err == nil {
		t.Fatal("kurs di luar periode kontrak lolos")
	}
	if err := simpan("", ""); err == nil {
		t.Fatal("kurs tanpa tanggal lolos")
	}
	if err := simpan("20250101", "20251231"); err != nil {
		t.Fatalf("kurs di dalam periode ditolak: %v", err)
	}
}
