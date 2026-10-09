package services_test

import (
	"testing"

	"nusantarare/modul/treatyin/backend/services"
)

func TestTotalRetensiDijumlahPerMataUang(t *testing.T) {
	rows := []services.RetensiNP{
		{Currency: "IDR", CurrencyID: "26", Amount: "3500000000"},
		{Currency: "USD", CurrencyID: "1", Amount: "100"},
		{Currency: "IDR", CurrencyID: "26", Amount: "500000000"},
	}
	total := services.NPSetTotalRetensi(rows)
	if len(total) != 2 {
		t.Fatalf("baris total = %d, mau 2", len(total))
	}
	if total[0].Currency != "IDR" || total[0].Value != "4000000000" {
		t.Fatalf("baris IDR = %+v", total[0])
	}
	if total[1].Currency != "USD" || total[1].Value != "100" {
		t.Fatalf("baris USD = %+v", total[1])
	}
}

// ⛔ Urutan baris total mengikuti kemunculan PERTAMA tiap mata uang -
// pola `Appendflag` Pega, bukan urutan abjad.
func TestUrutanTotalMengikutiKemunculanPertama(t *testing.T) {
	rows := []services.RetensiNP{
		{Currency: "USD", Amount: "1"},
		{Currency: "IDR", Amount: "1"},
	}
	total := services.NPSetTotalRetensi(rows)
	if total[0].Currency != "USD" {
		t.Fatalf("urutan = %q dulu, mau USD", total[0].Currency)
	}
}

// ⚠️ Beda dengan cabang `egnpi`: baris retensi baru KOSONG, nol mewarisi.
func TestTambahBarisRetensiKosong(t *testing.T) {
	rows := services.TambahBarisRetensi(nil)
	if len(rows) != 1 {
		t.Fatalf("jumlah baris = %d", len(rows))
	}
	if rows[0] != (services.RetensiNP{}) {
		t.Fatalf("baris baru tidak kosong: %+v", rows[0])
	}
}

func TestHitungRetensiHapusBaris(t *testing.T) {
	h := services.HitungRetensi(services.MasukanRetensi{
		Aksi:   services.AksiRetensiHapus,
		Indeks: 0,
		Retensi: []services.RetensiNP{
			{Currency: "USD", Amount: "1"},
			{Currency: "IDR", Amount: "2"},
		},
	})
	if len(h.Retensi) != 1 || h.Retensi[0].Currency != "IDR" {
		t.Fatalf("sesudah hapus = %+v", h.Retensi)
	}
	// ⭐ Total ikut dihitung ulang, tanpa menekan `Update Total`.
	if len(h.TotalRetentionAmountNP) != 1 || h.TotalRetentionAmountNP[0].Value != "2" {
		t.Fatalf("total = %+v", h.TotalRetentionAmountNP)
	}
}

// ⛔ Masukan tidak boleh berubah - rute `/hitung/` murni.
func TestHitungRetensiTidakMengubahMasukan(t *testing.T) {
	asal := []services.RetensiNP{{Currency: "USD", Amount: "1"}}
	services.HitungRetensi(services.MasukanRetensi{Aksi: services.AksiRetensiTambah, Retensi: asal})
	if len(asal) != 1 {
		t.Fatalf("masukan bertambah: %d", len(asal))
	}
}

// ⭐ Nol larik `nil` - `null` di JSON membuat layar jatuh pada `.length`.
func TestHasilRetensiNolLarikNil(t *testing.T) {
	h := services.HitungRetensi(services.MasukanRetensi{Aksi: "entah"})
	if h.Retensi == nil || h.TotalRetentionAmountNP == nil || h.Pesan == nil {
		t.Fatalf("ada larik nil: %+v", h)
	}
}

// ⛔ Tab ini NOL punya pagar bagi-nol, dan itu disengaja: cabang retensi
// tidak pernah membagi. Nol baris hanya berarti nol total.
func TestRetensiKosongNolPesan(t *testing.T) {
	h := services.HitungRetensi(services.MasukanRetensi{Aksi: services.AksiRetensiTotal})
	if len(h.Pesan) != 0 {
		t.Fatalf("pesan tak diduga: %v", h.Pesan)
	}
	if len(h.TotalRetentionAmountNP) != 0 {
		t.Fatalf("total tak diduga: %+v", h.TotalRetentionAmountNP)
	}
}
