package models_test

// Uji gabungan kondisi pembungkus layar (temuan cek layar 10-10-2026): pembungkus tampil / nonaktif yang ditumpuk
// MENGGABUNG kondisi, tidak menimpanya - tombol OQ tetap nonaktif, container ber-visible-when tetap menyaring unsurnya.

import (
	"testing"

	"nusantarare/modul/claimfacin/backend/models"
)

// cariTata - unsur pertama (rekursif) yang memenuhi `cocok`.
func cariTata(ts []models.Tata, cocok func(models.Tata) bool) (models.Tata, bool) {
	for _, t := range ts {
		if cocok(t) {
			return t, true
		}
		if t, ok := cariTata(t.Anak, cocok); ok {
			return t, true
		}
	}
	return models.Tata{}, false
}

func TestTombolOQTetapNonaktifWalauDibungkusNaJika(t *testing.T) {
	// "View KMT NO" = tombolOQ dibungkus naJika(`.KomiteNo != ''`) + tampil: KomiteNo kosong -> naJika salah, tombol
	// tetap nonaktif ber-OQ (pembungkus MENGGABUNG kondisi).
	h := models.HalamanBaru()
	h.SetelDaftar(models.DaftarObjek, []models.Baris{{}})
	h.SetelDaftar(models.DaftarItem(1), []models.Baris{{}})
	h.SetelDaftar(models.DaftarAdj(1, 1), []models.Baris{{"IsKomite": "1"}})
	ts := models.Evaluasi(h, models.LayarAdjustment(1, 1, 1), false)
	lihat, ok := cariTata(ts, func(t models.Tata) bool { return t.ID == "ViewKMTNo" })
	if !ok || !lihat.Nonaktif || lihat.Catatan == "" {
		t.Fatalf("View KMT NO harus nonaktif ber-OQ: %+v", lihat)
	}
	if models.AksiTerbuka(ts, "ViewKMTNo", 0) {
		t.Fatal("aksi View KMT NO terbuka")
	}
}

func TestYesRejectTT3MenggabungKondisi(t *testing.T) {
	// SureRejectClaim "Yes" (KCF-03): naJika(IsReject = 1) + tampil(Komite.CARI1 kosong) - keduanya berlaku bersama.
	cari := func(h *models.Halaman) (models.Tata, bool, []models.Tata) {
		ts := models.Evaluasi(h, models.LayarTolak(), false)
		y, ok := cariTata(ts, func(t models.Tata) bool { return t.ID == "SendRejectClaimToKomite2" })
		return y, ok, ts
	}
	h := models.HalamanBaru()
	if y, ok, ts := cari(h); !ok || y.Nonaktif || !models.AksiTerbuka(ts, "SendRejectClaimToKomite2", 0) {
		t.Fatalf("Yes Reject aktif: %+v", y)
	}
	h.Setel("IsReject", "1")
	if y, _, ts := cari(h); !y.Nonaktif || models.AksiTerbuka(ts, "SendRejectClaimToKomite2", 0) {
		t.Fatalf("IsReject 1 -> Yes nonaktif: %+v", y)
	}
	h.Setel("IsReject", "")
	h.Setel(models.JalurKomiteBaru, "KMT-UJI1")
	y, ok, ts := cari(h)
	if ok || models.AksiTerbuka(ts, "SendRejectClaimToKomite2", 0) {
		t.Fatalf("sesudah KMT lahir Yes tersembunyi: %+v", y)
	}
	if _, ok := cariTata(ts, func(t models.Tata) bool { return t.Label == models.TeksSuksesKomite }); !ok {
		t.Fatal("label Success Create Request to Committee tidak tampil")
	}
}

func TestContainerBertumpukMenyaringUnsur(t *testing.T) {
	grid := func(t models.Tata) bool {
		return t.Jenis == models.JenisGrid && t.Jalur == models.DaftarDiItem(1, 1, "CoverageList")
	}
	for _, c := range []struct {
		lini string
		mau  bool
	}{{"Fire", false}, {"MarineCargo", true}} {
		h := halamanAdj(models.Baris{})
		h.Setel(models.OQ+"BusinessType", c.lini)
		_, ada := cariTata(models.Evaluasi(h, models.LayarItemAdj(1, 1), false), grid)
		if ada != c.mau {
			t.Errorf("%s: grid coverage tampil=%v, mau %v", c.lini, ada, c.mau)
		}
	}
	// Send Claim to Committee: VIS isian wajib (LS7) DAN proteksi SendPICProtect_Act
	selalu := func(*models.Halaman) bool { return true }
	h := halamanAdj(models.Baris{})
	kirim := func(t models.Tata) bool { return t.ID == "KirimKomite" }
	if _, ada := cariTata(models.Evaluasi(h, models.LayarKomite(1, 1, 1, selalu), false), kirim); ada {
		t.Fatal("Send Claim to Committee tampil tanpa Remarks")
	}
	models.SetelJalur(h, models.JalurAdj(1, 1, 1)+".DataCommitteFacin.Remarks", "UJI")
	if _, ada := cariTata(models.Evaluasi(h, models.LayarKomite(1, 1, 1, selalu), false), kirim); !ada {
		t.Fatal("Send Claim to Committee tidak tampil walau Remarks terisi dan proteksi lolos")
	}
}

// Pra-proses EstimasiMarine_FA (GetSpreadingMarine_Act): spreading + kurs item Marine Cargo dari coverage TERAKHIR,
// hanya selama Spreading Policy item kosong (temuan review: tanpa kait ini estimasi Marine tidak pernah bisa ditambah).
func TestLengkapiMarine(t *testing.T) {
	k := konteksUji()
	h := models.HalamanBaru()
	h.Setel(models.OQ+"BusinessType", "MarineCargo")
	h.SetelDaftar(models.DaftarObjek, []models.Baris{{"ObjectName": "UJI KARGO"}})
	h.SetelDaftar(models.DaftarItem(1), []models.Baris{{"ObjectItemName": "UJI"}})
	cov := models.DaftarDiItem(1, 1, "CoverageList")
	h.SetelDaftar(cov, []models.Baris{
		{"Currency.Name": "IDR", "Currency.ID": idr, "TSINusantaraRe": "5"},
		{"Currency.Name": "USD", "Currency.ID": "UJI-USD", "TSINusantaraRe": "1000"},
	})
	h.SetelDaftar(models.JalurAnak(cov, 2, "SpreadingList"), []models.Baris{
		{"TreatyType": "10003", "SharePercentage": "60", "TSISpreaded": "600"},
		{"TreatyType": "10015", "SharePercentage": "40", "TSISpreaded": "400"},
	})
	if err := models.LengkapiMarine(k, h); err != nil {
		t.Fatal(err)
	}
	it, _ := models.Item(h, 1, 1)
	sp := h.AmbilDaftar(models.DaftarDiItem(1, 1, models.AnakSpreadPolis))
	if len(sp) != 2 || it["Currency"] != "USD" || it["KursObjectItem"] != "15000" || it["ValueTSINusareIDR"] != "15000000" {
		t.Fatalf("item %v spreading %v", it, sp)
	}
	kl := h.AmbilDaftar(models.DaftarDiItem(1, 1, models.AnakSpreadKlaim))
	kl[0]["ClaimSpreaded"] = "UJI-TETAP"
	if err := models.LengkapiMarine(k, h); err != nil {
		t.Fatal(err)
	}
	if v := h.AmbilDaftar(models.DaftarDiItem(1, 1, models.AnakSpreadKlaim))[0]["ClaimSpreaded"]; v != "UJI-TETAP" {
		t.Fatalf("pra-proses kedua menimpa Spreading Claim: %q", v)
	}
	// lini lain: tidak menyentuh
	h2 := models.HalamanBaru()
	h2.Setel(models.OQ+"BusinessType", "Fire")
	if err := models.LengkapiMarine(k, h2); err != nil {
		t.Fatal(err)
	}
}

// SetchronologyKlaimFacIn 1.1: jabatan pelaku untuk setiap baris kronologi (bukan hanya Committee); tanpa jabatan roster
// = "Claim Admin".
func TestKronologiJabatanSetiapBaris(t *testing.T) {
	for _, c := range []struct{ tingkat, mau string }{{"UJI HEAD", "UJI HEAD"}, {"", models.TingkatClaimAdmin}} {
		k := konteksUji()
		k.Langkah, k.Tingkat = "Input Estimasi", c.tingkat
		h := models.HalamanBaru()
		k.Kronologi(h, "UJI CATATAN")
		kr := h.AmbilDaftar(models.DaftarKronologi)
		if len(kr) != 1 || kr[0]["ASMUserID"] != c.mau || kr[0]["ASMNoteType"] == models.JenisCatatanKomite {
			t.Fatalf("tingkat %q: %v", c.tingkat, kr)
		}
	}
}

func TestCentangCWPTampilSelalu(t *testing.T) {
	// PreventRejectClaim LS21 (checkbox Close Without Payment) tanpa syarat tampil; hanya LS23 (isian + Yes) yang
	// tersembunyi sesudah kasus komite TT4 lahir (`Komite.CARI1 != ''`, LS28 label sukses).
	h := models.HalamanBaru()
	h.Setel(models.JalurKomiteBaru, "KMT-UJI1")
	ts := models.Evaluasi(h, models.LayarTutup(), false)
	if _, ok := cariTata(ts, func(t models.Tata) bool { return t.Jalur == models.JalurAlokasiSalvage }); !ok {
		t.Fatal("checkbox Close Without Payment hilang sesudah KMT lahir")
	}
	if _, ok := cariTata(ts, func(t models.Tata) bool { return t.Jalur == models.JalurTKRemarks }); ok {
		t.Fatal("isian LS23 masih tampil sesudah KMT lahir")
	}
}

func TestTambahEstimasiTanpaSpreadingKeluarDiLangkah6(t *testing.T) {
	// ValidateInputEstimate_act 6 (pre `countSpreading==0`): baris baru dibuang lalu transisi pasca `true` -> 6 keluar;
	// langkah 7+ (tanggal, Deductible, CopyCurrency, kurs) TIDAK berjalan atas baris lain.
	h := models.HalamanBaru()
	h.SetelDaftar(models.DaftarObjek, []models.Baris{{}})
	h.SetelDaftar(models.DaftarItem(1), []models.Baris{{"NetDeductibleValue": "5"}})
	h.SetelDaftar(models.DaftarDiItem(1, 1, models.AnakEstimasi), []models.Baris{{"EstimationDate": "2026-01-02",
		"Deductible": "999"}})
	if err := models.TambahEstimasi(konteksUji(), h, 1, 1); err != nil {
		t.Fatal(err)
	}
	rows := h.AmbilDaftar(models.DaftarDiItem(1, 1, models.AnakEstimasi))
	if len(rows) != 1 || rows[0]["Deductible"] != "999" {
		t.Fatalf("baris estimasi sesudah Add tanpa spreading: %+v", rows)
	}
	if len(h.Pesan) == 0 {
		t.Fatal("pesan tanpa spreading tidak tampil")
	}
}
