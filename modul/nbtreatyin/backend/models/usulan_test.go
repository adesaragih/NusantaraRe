package models

// Uji murni pemetaan SuggestList <-> POOLDATA.HISTORYAKSEPTASIPRODUCTION
// (Activity/SaveViewSuggest langkah 2 "UNTUK TREATY", RDBList/InsertViewSuggest_SQL;
// spec-penyimpanan ID-31, AC 39-44). Nilai harapan dari XML; fixture UJI-.

import (
	"strings"
	"testing"
)

func TestUsulanBelumTersimpanMenurutSaveViewSuggest(t *testing.T) { // AC 40, 41, 43, 44
	h := HalamanBaru()
	h.Setel(HalamanQuotation+".BusinessCode", "UJI-B1")
	h.Setel(HalamanQuotation+".BusinessFac", BisnisTreaty)
	panjang := strings.Repeat("é", 3995) // 3995 KARAKTER, bukan byte
	h.SetelDaftar(DaftarUsulan, []Baris{
		{"Suggest": "UJI-lama", "IsApproved": "1", "Date": "2026-10-01 08:00:00", "OperatorName": "Uji Lama", "OperatorID": "UJI-L", "IsSave": "Yes"},
		{"Suggest": panjang, "IsApproved": "0", "Date": "2026-10-03 15:30:00", "OperatorName": "Uji Admin", "OperatorID": "UJI-ADMIN"},
		{"Suggest": "UJI-tanpa-putusan", "Date": "2026-10-03 15:31:00", "OperatorName": "Uji Sec", "OperatorID": "UJI-SH"},
	})
	u := UsulanBelumTersimpan(h)
	if len(u) != 2 {
		t.Fatalf("hanya baris ber-IsSave kosong yang ditulis (langkah 2.1), dapat %d", len(u))
	}
	a := u[0]
	// 2.1.2: CARI1 @replaceAll(pyWorkIDPrefix "NB-","-","") · CARI3 "Policy" · CARI4 .OperatorName ·
	// CARI5 .Date · CARI7 Quotation.BusinessFac · CARI8 "2" · CARI9 "0" -> "Reject" ·
	// CARI10 @substring(.Suggest,0,3990) · AKSES_LOGIN OperatorID.pyUserIdentifier ·
	// BUSINESS_CODE Quotation.BusinessCode
	for nama, got := range map[string][2]string{
		"TYPE_POLIS":    {a.TypePolis, "NB"},
		"POSISI":        {a.Posisi, "Policy"},
		"PIC":           {a.PIC, "Uji Admin"},
		"TGL_INP":       {a.TglInp, "2026-10-03 15:30:00"},
		"TYPE":          {a.Type, "T"},
		"PUTARAN":       {a.Putaran, "2"},
		"APPROVAL":      {a.Approval, "Reject"},
		"AKSES_LOGIN":   {a.AksesLogin, "UJI-ADMIN"},
		"BUSINESS_CODE": {a.BusinessCode, "UJI-B1"},
		"DIV":           {a.Div, ""}, // OperatorID.pyOrgDivision - tanpa sumber di inti (butir terbuka)
		"B2B":           {a.B2B, ""}, // OfferFacIn - halaman Fac, kosong di NB
		"PERCENT_RNM":   {a.PercentRNM, ""},
	} {
		if got[0] != got[1] {
			t.Errorf("%s = %q, harap %q", nama, got[0], got[1])
		}
	}
	if n := len([]rune(a.Keterangan)); n != 3990 || a.Keterangan != strings.Repeat("é", 3990) {
		t.Errorf("KETERANGAN dipotong pada 3990 KARAKTER (AC 40), dapat %d", n)
	}
	if u[1].Approval != "" || u[1].Keterangan != "UJI-tanpa-putusan" || u[1].PIC != "Uji Sec" {
		t.Errorf("IsApproved kosong -> APPROVAL kosong: %+v", u[1])
	}
	// 2.1.4: .IsSave = "Yes" - baris yang sudah ditulis tidak ditulis lagi
	for i, b := range h.AmbilDaftar(DaftarUsulan) {
		if b["IsSave"] != "Yes" {
			t.Errorf("baris %d belum bertanda IsSave", i+1)
		}
	}
	if len(UsulanBelumTersimpan(h)) != 0 {
		t.Fatal("panggilan kedua menulis ulang baris yang sudah tersimpan")
	}
}

func TestApprovalPerBarisHanyaAcceptReject(t *testing.T) { // AC 41
	for masuk, harap := range map[string]string{"1": "Accept", "0": "Reject", "": "", "2": ""} {
		h := HalamanBaru()
		h.SetelDaftar(DaftarUsulan, []Baris{{"IsApproved": masuk}})
		if got := UsulanBelumTersimpan(h)[0].Approval; got != harap {
			t.Errorf("IsApproved %q -> APPROVAL %q, harap %q", masuk, got, harap)
		}
	}
}

func TestBarisCatatanDariRiwayatProduksi(t *testing.T) { // AC 42, 71
	b := BarisCatatan(UsulanProduksi{NoUrut: 3, PIC: "Uji Dept", TglInp: "2026-10-03 16:00:00",
		Approval: "Accept", Keterangan: "UJI-setuju", AksesLogin: "UJI-DH"})
	for m, harap := range map[string]string{"Suggest": "UJI-setuju", "IsApproved": "1", "Date": "2026-10-03 16:00:00",
		"OperatorName": "Uji Dept", "OperatorID": "UJI-DH", "IsSave": "Yes"} {
		if b[m] != harap {
			t.Errorf("%s = %q, harap %q", m, b[m], harap)
		}
	}
	if got := BarisCatatan(UsulanProduksi{Approval: "Reject"})["IsApproved"]; got != "0" {
		t.Errorf("Reject -> IsApproved %q, harap \"0\"", got)
	}
	if _, ada := BarisCatatan(UsulanProduksi{})["IsApproved"]; ada {
		t.Error("APPROVAL kosong -> IsApproved tidak dipasang (sama dengan AddToListComments langkah 1.2)")
	}
}

// K4 [penyimpangan sadar - menunggu konfirmasi WO]: NOURUT diberikan
// repository (MAX+1 per IDPEGA, di bawah kunci kasus), XML memakai
// `.pxListSubscript` (SaveViewSuggest 2.1.2 CARI2). Keduanya SAMA: SuggestList
// dibangun ulang dari tabel berurut NOURUT 1..n (BarisCatatan, ber-IsSave),
// catatan baru ditambahkan di UJUNG (`TambahCatatan`, APPEND), sehingga baris
// yang ditulis berposisi n+1 = MAX(NOURUT)+1. Uji alur ujung-ke-ujung:
// handlers/alur_test.go TestNourutUsulanSamaDenganSubskripSuggestList.
func TestCatatanBaruBerposisiNourutBerikut(t *testing.T) {
	for n := 0; n <= 3; n++ {
		h := HalamanBaru()
		var tabel []Baris // dibaca balik repository: ORDER BY NOURUT
		for i := 1; i <= n; i++ {
			tabel = append(tabel, BarisCatatan(UsulanProduksi{NoUrut: i, PIC: "UJI-PIC", Keterangan: "UJI-lama"}))
		}
		h.SetelDaftar(DaftarUsulan, tabel)
		h.Setel(HalamanPolis+".Suggest", "UJI-baru")
		TambahCatatan(h, "UJI-AKUN")
		daftar := h.AmbilDaftar(DaftarUsulan)
		u := UsulanBelumTersimpan(h)
		if len(u) != 1 || u[0].Keterangan != "UJI-baru" {
			t.Fatalf("n=%d: satu catatan baru, dapat %+v", n, u)
		}
		// .pxListSubscript (berbasis 1) baris yang ditulis = n+1 = MAX(NOURUT)+1
		if posisi := len(daftar); posisi != n+1 || daftar[posisi-1]["Suggest"] != "UJI-baru" {
			t.Fatalf("n=%d: catatan baru di pxListSubscript %d, harap %d", n, posisi, n+1)
		}
	}
}
