package handlers_test

// Item ganda Input Estimasi (temuan work owner 11-10-2026 "INI DI INPUT DOUBLE. TADI ADA WARNING! TAPI SAYA LANJUT INPUT
// TETEP BISA! KENAPA GITU?"): pilihan yang membuat item kembar ditolak (tidak disimpan, pesan tampil). Item kembar yang
// terlanjur tersimpan TIDAK menahan Send to PIC Claim ("SAAT DI KLIK SEND PIC DI BUANG AJA DEH"), tetapi menonaktifkan
// Download Claim Face Sheet objeknya ("JIKA ADA KEMBAR JANGAN BISA KLIK PRINT CLAIM FACE SHEET").

import (
	"net/http"
	"strings"
	"testing"

	"nusantarare/modul/claimfacin/backend/handlers"
	"nusantarare/modul/claimfacin/backend/models"
	"nusantarare/modul/claimfacin/backend/services"
	"nusantarare/modul/claimfacin/backend/tiruan"
)

func TestItemKembarDitolak(t *testing.T) {
	u := baruUji(t)
	id := u.sampaiEstimasi()
	est := models.KunciPanel(models.PanelObjekEst, models.DaftarObjek, 1)
	kode, out := u.aksiP(id, admin, "", services.PermintaanAksi{Aksi: "SetObjectItem", Konteks: est, Indeks: 1, Param: "1"})
	u.wajib(kode, http.StatusOK, out, "pilih property item baris 1")
	kode, out = u.aksiP(id, admin, "", services.PermintaanAksi{Aksi: "TambahItem", Konteks: est})
	u.wajib(kode, http.StatusOK, out, "tambah item")
	if n := len(u.g.Halaman(id).AmbilDaftar(models.DaftarItem(1))); n != 2 {
		t.Fatalf("item %d, mau 2", n)
	}
	// baris 2 memilih property item yang sama -> ditolak 422 berpesan, baris 2 tetap kosong
	kode, out = u.aksiP(id, admin, "", services.PermintaanAksi{Aksi: "SetObjectItem", Konteks: est, Indeks: 2, Param: "1"})
	if kode != http.StatusUnprocessableEntity {
		t.Fatalf("pilih item kembar: HTTP %d, mau 422: %v", kode, out)
	}
	pesan, _ := out["pesan"].([]any)
	if len(pesan) == 0 || pesan[0] != models.PesanItemKembar {
		t.Fatalf("pesan %v, mau %q", pesan, models.PesanItemKembar)
	}
	items := u.g.Halaman(id).AmbilDaftar(models.DaftarItem(1))
	if len(items) != 2 || items[1]["ObjectItemName"] != "" {
		t.Fatalf("item kembar tersimpan: %+v", items)
	}
	// data lama yang terlanjur kembar: Send to PIC Claim TIDAK ditahan pesan kembar
	h := u.g.Halaman(id)
	baris := h.AmbilDaftar(models.DaftarItem(1))
	baris[1]["ObjectItemName"] = baris[0]["ObjectItemName"]
	h.SetelDaftar(models.DaftarItem(1), baris)
	h.Setel(models.JalurIsPicTransfer, "1") // tombol Send to PIC aktif (selain item kembar, siap dikirim)
	u.g.SetelHalaman(id, h)
	kode, out = u.aksiP(id, admin, "", services.PermintaanAksi{Aksi: "SetDisable"})
	if kode != http.StatusOK {
		t.Fatalf("Send to PIC dengan item kembar: HTTP %d, mau 200 (tidak ditahan): %v", kode, out)
	}
}

func TestItemKembarMenahanFaceSheetDanPLA(t *testing.T) {
	// Work owner 11-10-2026 "TAMBAHKAN": item kembar yang terlanjur tersimpan menahan Download Claim Face Sheet dan Print PLA
	// objeknya (XML hanya menonaktifkan bila IsError > 1) - outstanding / PLA tidak tercatat dobel.
	u := baruUji(t)
	id := u.sampaiEstimasi()
	est := models.KunciPanel(models.PanelObjekEst, models.DaftarObjek, 1)
	item := models.KunciPanel(models.PanelItemEst, models.DaftarItem(1), 1)
	baris := models.JalurBaris(models.DaftarDiItem(1, 1, models.AnakEstimasi), 1)
	for _, p := range []services.PermintaanAksi{
		{Aksi: "SetObjectItem", Konteks: est, Indeks: 1, Param: "1"},
		{Aksi: "PilihCoverageFire", Konteks: est, Indeks: 1, Param: "UJI-C1"},
		{Aksi: "ValidateInputEstimate", Konteks: item},
		{Aksi: "CheckEstimateValue", Konteks: item, Indeks: 1,
			Masukan: map[string]string{baris + ".CurrencyID": idr, baris + ".GrossEstimationPct": "40000000"}},
		{Aksi: "TambahItem", Konteks: est},
	} {
		kode, out := u.aksiP(id, admin, "", p)
		u.wajib(kode, http.StatusOK, out, p.Aksi)
	}
	h := u.g.Halaman(id)
	b := h.AmbilDaftar(models.DaftarItem(1))
	b[1]["ObjectItemName"] = b[0]["ObjectItemName"] // data lama: kembar tersimpan sebelum 11-10-2026
	h.SetelDaftar(models.DaftarItem(1), b)
	h.Setel(models.JalurIsError, "1")
	objek := h.AmbilDaftar(models.DaftarObjek)
	objek[0]["IsFacretro"] = "1" // tombol Print PLA aktif
	h.SetelDaftar(models.DaftarObjek, objek)
	u.g.SetelHalaman(id, h)
	// tombol Download Claim Face Sheet baris objek 1 NONAKTIF di layar, dan server menolak aksinya
	kode, out := u.minta(http.MethodGet, handlers.Prefix+"/kasus/"+id, admin, "", nil)
	u.wajib(kode, http.StatusOK, out, "buka kasus")
	if !tombolBarisNonaktif(t, out, "CLaimFaceSheet", 0) {
		t.Fatal("Download Claim Face Sheet masih aktif padahal objek punya item kembar")
	}
	kode, out = u.aksiP(id, admin, "", services.PermintaanAksi{Aksi: "CLaimFaceSheet", Indeks: 1})
	if kode == http.StatusOK {
		t.Fatalf("Download Claim Face Sheet dengan item kembar diterima: %v", out)
	}
	// Print PLA: ditolak berpesan
	kode, out = u.aksiP(id, admin, "", services.PermintaanAksi{Aksi: "BukaPLA", Indeks: 1})
	pesan, _ := out["pesan"].([]any)
	if kode != http.StatusUnprocessableEntity || len(pesan) == 0 || pesan[0] != models.PesanItemKembar {
		t.Fatalf("Print PLA dengan item kembar: HTTP %d pesan %v, mau 422 %q", kode, pesan, models.PesanItemKembar)
	}
	if no := u.g.Halaman(id).Ambil(models.CD + "NoClaim"); no != "" || len(u.g.OS) != 0 {
		t.Fatalf("face sheet tertolak tetap mencatat: NoClaim %q, OS %d", no, len(u.g.OS))
	}
	// baris kembar dihapus -> tombol aktif, face sheet berjalan
	kode, out = u.aksiP(id, admin, "", services.PermintaanAksi{Aksi: "HapusItem", Konteks: est, Indeks: 2})
	u.wajib(kode, http.StatusOK, out, "hapus item kembar")
	kode, out = u.aksiP(id, admin, "", services.PermintaanAksi{Aksi: "CLaimFaceSheet", Indeks: 1})
	u.wajib(kode, http.StatusOK, out, "face sheet sesudah kembar dihapus")
	if len(u.g.OS) != 1 {
		t.Fatalf("OS %d, mau 1", len(u.g.OS))
	}
}

// tombolBarisNonaktif - sel tombol `id` grid layar utama pada baris ke-n (0-based) berstatus nonaktif.
func tombolBarisNonaktif(t *testing.T, out map[string]any, id string, n int) bool {
	t.Helper()
	var cari func(any) (bool, bool)
	cari = func(x any) (bool, bool) {
		switch v := x.(type) {
		case map[string]any:
			if v["jenis"] == "grid" {
				kolom, _ := v["kolom"].([]any)
				for j, k := range kolom {
					if km, _ := k.(map[string]any); km["id"] == id {
						baris, _ := v["baris"].([]any)
						if n >= len(baris) {
							return false, false
						}
						sel, _ := baris[n].([]any)
						s, _ := sel[j].(map[string]any)
						na, _ := s["nonaktif"].(bool)
						return na, true
					}
				}
			}
			for _, w := range v {
				if na, ok := cari(w); ok {
					return na, true
				}
			}
		case []any:
			for _, w := range v {
				if na, ok := cari(w); ok {
					return na, true
				}
			}
		}
		return false, false
	}
	na, ok := cari(out["tata"])
	if !ok {
		t.Fatalf("tombol %q tidak ada di grid layar", id)
	}
	return na
}

// polisPA - dokumen JSON_POLIS fixture lini PA: satu peserta ber-dua coverage (ASMCoverage), spreading QS. ASMIDCard
// "1": item PA mencari peserta polis lewat IndexObject (= ObjectID objek = ASMIDCard) yang dibandingkan dengan nomor
// baris PersonList (`OpsiCoverageObjek`, warisan XML).
var polisPA = strings.NewReplacer(`"BusinessType": "Fire"`, `"BusinessType": "PA"`).Replace(polisFire[:strings.Index(polisFire,
	`"LocationList"`)]) + `"PersonList": [{"pyFullName": "UJI PESERTA", "ASMIDCard": "1", "ASMJobName": "UJI KERJA",
  "ASMCoverage": [
   {"Coverage": "UJI-PA1", "CoverageNote": "UJI MENINGGAL", "TSILiability": "100000000", "TSINusantaraRe": "25000000",
    "Currency": {"Name": "IDR", "ID": "10026"},
    "SpreadingList": [{"TreatyType": "10003", "TreatyName": "UJI QS", "SharePercentage": "100", "TSISpreaded": "25000000"}]},
   {"Coverage": "UJI-PA2", "CoverageNote": "UJI CACAT", "TSILiability": "50000000", "TSINusantaraRe": "12500000",
    "Currency": {"Name": "IDR", "ID": "10026"},
    "SpreadingList": [{"TreatyType": "10003", "TreatyName": "UJI QS", "SharePercentage": "100", "TSISpreaded": "12500000"}]}]}]
}`

func TestCoverageKembarPADitolak(t *testing.T) {
	// Work owner 11-10-2026 "PERBAIKI JUGA UNTUK BISNIS LAIN!": proteksi item kembar Fire / Aneka berlaku juga untuk
	// MBU / PA / Travel (kriteria `ProtectCoverage_Act` / DeleteTest = Coverage ID sama). XML membuang item TERAKHIR
	// ber-coverage sama - bila yang diubah bukan baris terakhir, item lain yang terhapus; di sini pilihannya DIBATALKAN.
	u := baruUji(t)
	u.a.Dokumen[tiruan.KunciPolis(polisUji, "1")] = []byte(polisPA)
	id := u.sampaiEstimasi()
	est := models.KunciPanel(models.PanelObjekEst, models.DaftarObjek, 1)
	for _, p := range []services.PermintaanAksi{ // objek PA lahir berisi satu item kosong (InsertObjects)
		{Aksi: "PilihCoverageObjek", Konteks: est, Indeks: 1, Param: "UJI-PA1"},
		{Aksi: "TambahItem", Konteks: est},
		{Aksi: "PilihCoverageObjek", Konteks: est, Indeks: 2, Param: "UJI-PA2"},
	} {
		kode, out := u.aksiP(id, admin, "", p)
		u.wajib(kode, http.StatusOK, out, p.Aksi)
	}
	// baris 1 memilih coverage baris 2 -> ditolak 422 berpesan; kedua item tetap seperti semula (XML: baris 2 terhapus)
	kode, out := u.aksiP(id, admin, "", services.PermintaanAksi{Aksi: "PilihCoverageObjek", Konteks: est, Indeks: 1,
		Param: "UJI-PA2"})
	if kode != http.StatusUnprocessableEntity {
		t.Fatalf("coverage kembar: HTTP %d, mau 422: %v", kode, out)
	}
	if pesan, _ := out["pesan"].([]any); len(pesan) == 0 || pesan[0] != models.PesanItemKembar {
		t.Fatalf("pesan %v, mau %q", out["pesan"], models.PesanItemKembar)
	}
	items := u.g.Halaman(id).AmbilDaftar(models.DaftarItem(1))
	if len(items) != 2 || items[0]["CoverageOLDID"] != "UJI-PA1" || items[1]["CoverageOLDID"] != "UJI-PA2" {
		t.Fatalf("item sesudah coverage kembar ditolak: %+v", items)
	}
	// data lama yang terlanjur kembar (Coverage ID sama) menahan Download Claim Face Sheet dan Print PLA objeknya
	h := u.g.Halaman(id)
	items[1]["CoverageOLDID"] = items[0]["CoverageOLDID"]
	h.SetelDaftar(models.DaftarItem(1), items)
	objek := h.AmbilDaftar(models.DaftarObjek)
	objek[0]["IsFacretro"] = "1" // tombol Print PLA aktif
	h.SetelDaftar(models.DaftarObjek, objek)
	u.g.SetelHalaman(id, h)
	kode, out = u.minta(http.MethodGet, handlers.Prefix+"/kasus/"+id, admin, "", nil)
	u.wajib(kode, http.StatusOK, out, "buka kasus")
	if !tombolBarisNonaktif(t, out, "CLaimFaceSheet", 0) {
		t.Fatal("Download Claim Face Sheet masih aktif padahal objek PA punya coverage kembar")
	}
	if kode, out = u.aksiP(id, admin, "", services.PermintaanAksi{Aksi: "CLaimFaceSheet", Indeks: 1}); kode == http.StatusOK {
		t.Fatalf("Download Claim Face Sheet dengan coverage kembar diterima: %v", out)
	}
	kode, out = u.aksiP(id, admin, "", services.PermintaanAksi{Aksi: "BukaPLA", Indeks: 1})
	if pesan, _ := out["pesan"].([]any); kode != http.StatusUnprocessableEntity || len(pesan) == 0 ||
		pesan[0] != models.PesanItemKembar {
		t.Fatalf("Print PLA coverage kembar: HTTP %d pesan %v, mau 422 %q", kode, out["pesan"], models.PesanItemKembar)
	}
	if len(u.g.OS) != 0 {
		t.Fatalf("face sheet tertolak tetap mencatat OS %d", len(u.g.OS))
	}
}

func TestCoverageFireDitolakSesudahFaceSheet(t *testing.T) {
	// Work owner 11-10-2026 "KENAPA MASIH BISA EDIT SETELAH KLIK PRINT CLAIM FACESHEET?!": sesudah Download Claim Face
	// Sheet, Coverage Name item Fire hanya-baca dan aksinya ditolak server (outstanding sudah tercatat).
	u := baruUji(t)
	id := u.sampaiEstimasi()
	est := models.KunciPanel(models.PanelObjekEst, models.DaftarObjek, 1)
	item := models.KunciPanel(models.PanelItemEst, models.DaftarItem(1), 1)
	baris := models.JalurBaris(models.DaftarDiItem(1, 1, models.AnakEstimasi), 1)
	for _, p := range []services.PermintaanAksi{
		{Aksi: "SetObjectItem", Konteks: est, Indeks: 1, Param: "1"},
		{Aksi: "PilihCoverageFire", Konteks: est, Indeks: 1, Param: "UJI-C1"},
		{Aksi: "ValidateInputEstimate", Konteks: item},
		{Aksi: "CheckEstimateValue", Konteks: item, Indeks: 1,
			Masukan: map[string]string{baris + ".CurrencyID": idr, baris + ".GrossEstimationPct": "40000000"}},
		{Aksi: "CLaimFaceSheet", Indeks: 1},
	} {
		kode, out := u.aksiP(id, admin, "", p)
		u.wajib(kode, http.StatusOK, out, p.Aksi)
	}
	if it := u.g.Halaman(id).AmbilDaftar(models.DaftarItem(1))[0]; it["PrintFaceClaim"] != "1" {
		t.Fatalf("item belum tercetak CFS: %+v", it)
	}
	kode, out := u.aksiP(id, admin, "", services.PermintaanAksi{Aksi: "PilihCoverageFire", Konteks: est, Indeks: 1,
		Param: "UJI-C1"})
	if kode == http.StatusOK {
		t.Fatalf("Coverage Name sesudah Claim Face Sheet diterima: %v", out)
	}
}
