package handlers_test

// Permintaan work owner 06-10-2026: medan UANG wajib (bagian Gross/OGP/ONP/klaim/potongan) yang dibiarkan kosong
// diisi 0 otomatis saat Save / Submit, bukan ditolak. Medan wajib lain (tanggal, pilihan, Quartal, Suggest)
// tetap ditolak bila kosong.

import (
	"net/http"
	"strings"
	"testing"
)

var uangWajibAdmin = []string{"PremiOgp", "RiCommOgp", "OveriddingCommOgp", "PremiOnp", "RiCommOnp", "OveriddingCommOnp",
	"Claim", "OutstandingClaim", "SalvageValue", "ExcessLoss", "Deduction1", "Deduction2"}

func TestSaveMengisiNolMedanUangWajibYangKosong(t *testing.T) {
	u := baru(t)
	id := u.buat()
	h := halamanLengkap("1")
	for _, m := range uangWajibAdmin {
		h.Setel("PolicyTreatyIn."+m, "")
	}
	if kode, isi := u.panggil("PUT", "/kasus/"+id, admin, map[string]any{"halaman": h}); kode != http.StatusOK {
		t.Fatalf("save dengan uang wajib kosong: %d %s", kode, isi)
	}
	for _, m := range uangWajibAdmin {
		if got := u.g.Halaman[id].Ambil("PolicyTreatyIn." + m); got != "0" {
			t.Errorf("%s tersimpan %q, harap 0", m, got)
		}
	}
}

// Permintaan work owner 06-10-2026: Approval dan Suggest bukan halangan Save - wajib hanya saat Submit.
func TestSaveTanpaApprovalDanSuggestBoleh(t *testing.T) {
	u := baru(t)
	id := u.buat()
	h := halamanLengkap("1")
	h.Setel("PolicyTreatyIn.IsApproved", "")
	h.Setel("PolicyTreatyIn.Suggest", "")
	if kode, isi := u.panggil("PUT", "/kasus/"+id, admin, map[string]any{"halaman": h}); kode != http.StatusOK {
		t.Fatalf("save tanpa Approval dan Suggest: %d %s", kode, isi)
	}
	// tetap wajib saat Submit: Approval lebih dulu, lalu Suggest
	if kode, isi := u.kirim(id, admin, h); kode != http.StatusUnprocessableEntity || !strings.Contains(isi, "Approval") {
		t.Fatalf("submit tanpa Approval harus ditolak: %d %s", kode, isi)
	}
	h.Setel("PolicyTreatyIn.IsApproved", "1")
	if kode, isi := u.kirim(id, admin, h); kode != http.StatusUnprocessableEntity || !strings.Contains(isi, "Suggest") {
		t.Fatalf("submit tanpa Suggest harus ditolak: %d %s", kode, isi)
	}
}

func TestSaveTetapMenolakMedanWajibLain(t *testing.T) {
	u := baru(t)
	id := u.buat()
	h := halamanLengkap("1")
	h.Setel("PolicyTreatyIn.EndDate", "")
	kode, isi := u.panggil("PUT", "/kasus/"+id, admin, map[string]any{"halaman": h})
	if kode != http.StatusUnprocessableEntity || !strings.Contains(isi, "To") {
		t.Fatalf("To (EndDate) kosong harus ditolak saat Save: %d %s", kode, isi)
	}
}

func TestSubmitMengisiNolMedanUangWajibYangKosong(t *testing.T) {
	u := baru(t)
	id := u.buat()
	h := halamanLengkap("1")
	for _, m := range uangWajibAdmin {
		h.Setel("PolicyTreatyIn."+m, "")
	}
	if kode, isi := u.kirim(id, admin, h); kode != http.StatusOK {
		t.Fatalf("submit dengan uang wajib kosong: %d %s", kode, isi)
	}
	for _, m := range uangWajibAdmin {
		if got := u.g.Halaman[id].Ambil("PolicyTreatyIn." + m); got != "0" {
			t.Errorf("%s tersimpan %q, harap 0", m, got)
		}
	}
}
