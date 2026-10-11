package models_test

// Uji gabungan kondisi pembungkus layar (temuan cek layar 10-10-2026): pembungkus tampil / nonaktif yang ditumpuk
// MENGGABUNG kondisi, tidak menimpanya - tombol OQ tetap nonaktif, container ber-visible-when tetap menyaring unsurnya.

import (
	"strings"
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
	// Send Claim to Committee: tampil menurut proteksi SendPICProtect_Act saja - VIS isian wajib LS7 dievaluasi server
	// dan mengetik di pop-up tidak memanggil server, jadi tombol tak pernah muncul (work owner 11-10-2026 "TOMBOL SEND
	// NYA GA ADA!"). Remarks tetap wajib (ditegakkan saat diklik).
	selalu := func(*models.Halaman) bool { return true }
	tidak := func(*models.Halaman) bool { return false }
	h := halamanAdj(models.Baris{})
	kirim := func(t models.Tata) bool { return t.ID == "KirimKomite" }
	ts := models.Evaluasi(h, models.LayarKomite(1, 1, 1, selalu), false)
	if _, ada := cariTata(ts, kirim); !ada {
		t.Fatal("Send Claim to Committee tidak tampil padahal proteksi lolos (Remarks belum diketik)")
	}
	remarks, ok := cariTata(ts, func(x models.Tata) bool { return x.Label == "Remarks" })
	if !ok || !remarks.Wajib {
		t.Fatalf("Remarks harus wajib: %+v", remarks)
	}
	if _, ada := cariTata(models.Evaluasi(h, models.LayarKomite(1, 1, 1, tidak), false), kirim); ada {
		t.Fatal("Send Claim to Committee tampil padahal proteksi lampiran / rekening / premi gagal")
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

func TestDetailPolisTampilanNBDenganPenandaFire(t *testing.T) {
	// Perintah work owner 10-10-2026 "IKUTI MODUL NB FACIN. TAMPILAN READ ONLY": detail polis Choose Polis = unsur
	// polisNB tunggal (blok General Pega diganti), kendali "fire" hanya bila When IsFire benar; View Retro List dalam
	// baris Inline sehingga tidak terbawa ke kaki modal bersama Cancel / Submit.
	polis := func(jenis string) []models.Tata {
		h := models.HalamanBaru()
		h.Setel(models.JalurDetailPolis, "true")
		h.Setel(models.OQ+"BusinessType", jenis)
		h.SetelDaftar(models.DaftarFacRetro, []models.Baris{{"ReinsurerName": "UJI RETRO"}})
		return models.Evaluasi(h, models.LayarPilihPolis(), false)
	}
	for _, c := range []struct{ jenis, kendali string }{{"Fire", models.KPolisFire}, {"FireStyle1", models.KPolisFire}, {"Liability", ""}} {
		ts := polis(c.jenis)
		var nb []models.Tata
		var cari func([]models.Tata)
		cari = func(xs []models.Tata) {
			for _, x := range xs {
				if x.Jenis == models.JenisPolisNB {
					nb = append(nb, x)
				}
				cari(x.Anak)
			}
		}
		cari(ts)
		if len(nb) != 1 || nb[0].Kendali != c.kendali {
			t.Fatalf("%s: polisNB %+v, mau satu berkendali %q", c.jenis, nb, c.kendali)
		}
		if _, ok := cariTata(ts, func(x models.Tata) bool { return x.Jenis == models.JenisMedan && x.Jalur == models.JalurNoPolis }); ok {
			t.Fatalf("%s: medan blok General lama masih ada", c.jenis)
		}
		retro, ok := cariTata(ts, func(x models.Tata) bool {
			return x.Letak == "sebaris" && len(x.Anak) == 1 && x.Anak[0].ID == "ViewRetroList"
		})
		if !ok || !models.AksiTerbuka(ts, "ViewRetroList", 0) {
			t.Fatalf("%s: View Retro List harus di baris Inline sendiri dan terbuka: %+v", c.jenis, retro)
		}
	}
}

func TestInputRegisterDuaKolomLabelAtasUrutXML(t *testing.T) {
	// Screenshot Pega work owner 11-10-2026 "UBAH LAYOUTNYA BEGINI" = kontainer 650 Section InputRegisterDetail: dua kolom
	// diisi per baris, label di atas isian.
	h := models.HalamanBaru()
	h.Setel(models.JalurNoPolis, "UJI-RNM-F.001")
	ts := models.Evaluasi(h, models.LayarRegister(), false)
	grid, ok := cariTata(ts, func(x models.Tata) bool { return x.Letak == models.LetakDuaAtas })
	if !ok {
		t.Fatal("tanpa layout duaAtas")
	}
	var urut []string
	for _, c := range grid.Anak {
		switch {
		case c.Jenis == models.JenisMedan:
			urut = append(urut, c.Label)
		case len(c.Anak) > 0:
			urut = append(urut, c.Anak[0].Label)
		}
	}
	mau := []string{"Policy No", "Insured Name", "QQ Name", "Ceding Co Name", "Source Of Business", "Date of Loss",
		"Report Date", "Received Date", "Reporter Name", "Reporter Phone Number", "Choose Cause of Loss", "Catastrophe"}
	if strings.Join(urut, "|") != strings.Join(mau, "|") {
		t.Fatalf("urutan sel\n dapat %v\n mau   %v", urut, mau)
	}
}

func TestTabProgressClaimDisembunyikan(t *testing.T) {
	// Work owner 11-10-2026 "TAB PROGRESS CLAIM HIDE AJA! GA KEPAKE JUGA" - di ketiga layar kasus.
	h := models.HalamanBaru()
	h.SetelDaftar(models.DaftarObjek, []models.Baris{{}})
	h.SetelDaftar(models.DaftarItem(1), []models.Baris{{}})
	for nama, defs := range map[string][]models.Unsur{"register": models.LayarRegister(), "estimasi": models.LayarEstimasi(),
		"adjustment": models.LayarSurveyor()} {
		ts := models.Evaluasi(h, defs, false)
		if _, ok := cariTata(ts, func(x models.Tata) bool { return x.Label == "Progress Claim" }); ok {
			t.Errorf("%s: tab Progress Claim masih tampil", nama)
		}
		if _, ok := cariTata(ts, func(x models.Tata) bool { return strings.HasPrefix(x.Label, "Policy Detail") }); !ok {
			t.Errorf("%s: tab lain ikut hilang", nama)
		}
	}
}

func TestPanelItemEstimasiTotalDiKakiGrid(t *testing.T) {
	// Work owner 11-10-2026 "tolong rapihin. buat lebih user frendly dan enak dilihat!": total estimasi / spreading =
	// kaki grid-nya; Deductible Info satu baris berlabel. "PERBAIKI JUGA UNTUK BISNIS LAIN!" (11-10-2026): semua lini -
	// Estimasi (Fire / Aneka / Golf / MBU), EstimasiPA (PA / Travel), EstimasiMarine (Marine Cargo; kaki "Total" tanpa
	// label XML mengulang dua total yang sama).
	const totalEst = "Total Estimation Gross|Total Estimation Gross in IDR|Total Estimation RNM|Total Estimation RNM in IDR"
	for _, c := range []struct {
		lini       string
		estimasi   string // label kaki grid EstimationList ("" = tanpa kaki)
		spreading  bool   // grid Spreading Policy berkaki
		deductible bool
	}{
		{"Fire", totalEst, true, true}, {"Liability", totalEst, true, true}, {"GolfInsurance", totalEst, true, true},
		{"MBUCar", totalEst, true, true}, {"PA", totalEst, true, false}, {"MarineCargo", totalEst, true, false},
		{"Travel", "", false, false},
	} {
		h := models.HalamanBaru()
		h.Setel(models.OQ+"BusinessType", c.lini)
		h.Setel(models.CD+"ExGratia", "0")
		h.SetelDaftar(models.DaftarObjek, []models.Baris{{}})
		h.SetelDaftar(models.DaftarItem(1), []models.Baris{{}})
		ts := models.Evaluasi(h, models.LayarItemEstimasi(1, 1), false)
		kaki := map[string][]string{}
		var lepas []string
		var jalan func([]models.Tata)
		jalan = func(xs []models.Tata) {
			for _, x := range xs {
				if x.Jenis == models.JenisGrid {
					for _, k := range x.Kaki {
						kaki[x.Jalur] = append(kaki[x.Jalur], k.Label)
					}
				}
				if (x.Jenis == models.JenisMedan && strings.HasPrefix(x.Label, "Total ")) || x.Label == "Total" {
					lepas = append(lepas, x.Label)
				}
				jalan(x.Anak)
			}
		}
		jalan(ts)
		if len(lepas) > 0 {
			t.Errorf("%s: total di luar kaki grid / kaki tanpa label: %v", c.lini, lepas)
		}
		mau := map[string]string{"EstimationList": c.estimasi}
		if c.spreading {
			mau["SpreadingList"] = "Total Share Percentage|Total TSI Spreaded"
		}
		for akhiran, label := range mau {
			ok := false
			for j, l := range kaki {
				if strings.HasSuffix(j, akhiran) && strings.Join(l, "|") == label {
					ok = true
				}
			}
			if !ok && label != "" {
				t.Errorf("%s: kaki grid ...%s = %q tidak ada; kaki %v", c.lini, akhiran, label, kaki)
			}
		}
		ded, ok := cariTata(ts, func(x models.Tata) bool { return x.Letak == models.LetakSebaris && x.Label == "Deductible Info" })
		if ok != c.deductible {
			t.Errorf("%s: Deductible Info tampil=%v, mau %v", c.lini, ok, c.deductible)
		}
		if ok && (len(ded.Anak) == 0 || ded.Anak[0].Kendali != models.KCentang) {
			t.Errorf("%s: Deductible Info harus satu baris berlabel berisi centang: %+v", c.lini, ded)
		}
	}
}

func TestPanelAdjustmentTigaKolomIkutLayoutLama(t *testing.T) {
	// Screenshot layout lama work owner 11-10-2026 "PERBAIKI TAMPILAN ADJUSTMENT! INI SAYA LAMPIRKAN LAYOUT LAMA": LS12 dan
	// LS15 Inline grid triple berlabel di atas; sel = kontainer XML (tersembunyi tidak memakan sel); Payable / bank
	// berlabel kiri; total spreading di kaki grid.
	selPertama := func(c models.Tata) string {
		for c.Jenis == models.JenisBagian && len(c.Anak) > 0 {
			c = c.Anak[0]
		}
		return c.Label
	}
	for _, c := range []struct {
		komite string
		mau    []string
	}{
		{"", []string{"Payment Type", "Currency Estimation", "RNM Share(%)", "Choose Currency", "Gross Adjustment (100%)",
			"Deductible Type", "Deductible (100%)", "Nett Adjustment (100%)", "Payable To", "Name of Bank"}},
		{"1", []string{"Payment Type", "Currency Estimation", "RNM Share(%)", "Currency Adjustment",
			"Gross Adjustment (100%)", "Deductible Type", "Deductible (100%)", "Nett Adjustment (100%)", "Payable To",
			"Name of Bank", "Komite Number"}},
	} {
		h := models.HalamanBaru()
		h.Setel(models.OQ+"BusinessType", "Fire")
		h.SetelDaftar(models.DaftarObjek, []models.Baris{{}})
		h.SetelDaftar(models.DaftarItem(1), []models.Baris{{}})
		// komite: sudah diakseptasi seperti layout lama (View KMT NO hanya selama AcceptanceStatus kosong)
		h.SetelDaftar(models.DaftarAdj(1, 1), []models.Baris{{"PaymentType": "1", "Currency": "IDR", "IsKomite": c.komite,
			"KomiteNo": "UJI-KMT-1", "ExGratia": "0", "AcceptanceStatus": c.komite}})
		ts := models.Evaluasi(h, models.LayarAdjustment(1, 1, 1), false)
		detail, ok := cariTata(ts, func(x models.Tata) bool { return x.Label == "Adjusment Detail" })
		if !ok {
			t.Fatal("tanpa Adjusment Detail")
		}
		var tiga []models.Tata
		var jalan func([]models.Tata)
		jalan = func(xs []models.Tata) {
			for _, x := range xs {
				if x.Letak == models.LetakTiga {
					tiga = append(tiga, x)
				}
				jalan(x.Anak)
			}
		}
		jalan(detail.Anak)
		if len(tiga) != 2 {
			t.Fatalf("komite %q: layout tiga %d, mau 2 (LS12, LS15)", c.komite, len(tiga))
		}
		var atas []string
		for _, x := range tiga[0].Anak {
			atas = append(atas, selPertama(x))
		}
		if mau := []string{"DLA No Ceding", "Net For Collection", "Transfer Direct to Kasir"}; strings.Join(atas, "|") !=
			strings.Join(mau, "|") {
			t.Errorf("komite %q LS12\n dapat %v\n mau   %v", c.komite, atas, mau)
		}
		var urut []string
		for _, x := range tiga[1].Anak {
			if x.Jenis == models.JenisBagian && len(x.Anak) == 0 {
				continue // kontainer tersembunyi seluruhnya: tidak memakan sel
			}
			urut = append(urut, selPertama(x))
			if selPertama(x) == "Currency Estimation" && len(x.Anak) != 3 {
				t.Errorf("LS17 sejajar %d medan, mau 3", len(x.Anak))
			}
			if lt := x.Letak; (selPertama(x) == "Payable To" || selPertama(x) == "Name of Bank") != (lt == models.LetakKiri) {
				t.Errorf("sel %q letak %q", selPertama(x), lt)
			}
		}
		if strings.Join(urut, "|") != strings.Join(c.mau, "|") {
			t.Errorf("komite %q LS15\n dapat %v\n mau   %v", c.komite, urut, c.mau)
		}
		g, ok := cariTata(ts, func(x models.Tata) bool { return x.Jenis == models.JenisGrid && x.Label == "Spreading Adjustment" })
		if !ok || len(g.Kaki) != 2 || g.Kaki[0].Label != "Total Share Percentage(%)" || g.Kaki[1].Label != "Total Claim Spread" {
			t.Errorf("kaki Spreading Adjustment %+v", g.Kaki)
		}
	}
	// prompt values PaymentType (screenshot property rule work owner 11-10-2026) dan Ex Gratia layout lama
	for kode, mau := range map[string]string{"1": "Final", "2": "Claim", "3": "Salvage", "4": "Adjuster Fee", "5": "Adjustment",
		"6": "Consultant Fee", "7": "Cancellation"} {
		if got := models.LabelKode["PaymentType"][kode]; got != mau {
			t.Errorf("PaymentType %s = %q, mau %q", kode, got, mau)
		}
	}
	if models.LabelKode["ExGratia"]["0"] != "NO" {
		t.Errorf("ExGratia 0 = %q, mau NO", models.LabelKode["ExGratia"]["0"])
	}
}

func TestDetailChronologyKomiteTigaKolom(t *testing.T) {
	// Work owner 11-10-2026 "RAPIHKAN!": Detail Chronology to Comitee (LS3 + LS4, tiga medan sebaris masing-masing) = satu
	// Inline grid triple berlabel atas, urutan XML; hanya tampil bila adjustment sudah dikirim ke komite.
	h := models.HalamanBaru()
	h.SetelDaftar(models.DaftarObjek, []models.Baris{{}})
	h.SetelDaftar(models.DaftarItem(1), []models.Baris{{}})
	h.SetelDaftar(models.DaftarAdj(1, 1), []models.Baris{{"IsKomite": "1", "DataCommitteFacin.Remarks": "UJI"}})
	bag, ok := cariTata(models.Evaluasi(h, models.LayarAdjustment(1, 1, 1), false),
		func(x models.Tata) bool { return x.Label == "Detail Chronology to Comitee" })
	if !ok || len(bag.Anak) != 1 || bag.Anak[0].Letak != models.LetakTiga {
		t.Fatalf("Detail Chronology to Comitee: %+v", bag)
	}
	var label []string
	for _, c := range bag.Anak[0].Anak {
		label = append(label, c.Label)
		if !c.Nonaktif {
			t.Errorf("%s harus hanya-baca (NA always)", c.Label)
		}
	}
	mau := "Chronology|Extent Of Loss|Legal Liability|Salvage|Adjuster Fee|Remarks"
	if strings.Join(label, "|") != mau {
		t.Errorf("urutan %v, mau %s", label, mau)
	}
	h.SetelDaftar(models.DaftarAdj(1, 1), []models.Baris{{}})
	if _, ada := cariTata(models.Evaluasi(h, models.LayarAdjustment(1, 1, 1), false),
		func(x models.Tata) bool { return x.Label == "Detail Chronology to Comitee" }); ada {
		t.Error("Detail Chronology to Comitee tampil sebelum dikirim ke komite")
	}
}

func TestPanelAdjustmentBentukTab(t *testing.T) {
	// Work owner 11-10-2026 "INI BENTUK TAB!": InputAdjustment LS1 = layout group Tab - "Detail Chronology to Comitee"
	// (VIS IsKomite = 1) dan "Adjusment Detail"; spreading, Committee Accept Status, dan tombol di dalam tab kedua.
	// "STATUS NYA BENERIN JANGAN ANGKA": AcceptanceStatus berlabel.
	for _, c := range []struct {
		komite string
		mau    string
	}{{"1", "Detail Chronology to Comitee|Adjusment Detail"}, {"", "Adjusment Detail"}} {
		h := models.HalamanBaru()
		h.SetelDaftar(models.DaftarObjek, []models.Baris{{}})
		h.SetelDaftar(models.DaftarItem(1), []models.Baris{{}})
		h.SetelDaftar(models.DaftarAdj(1, 1), []models.Baris{{"IsKomite": c.komite, "PaymentType": "1"}})
		ts := models.Evaluasi(h, models.LayarAdjustment(1, 1, 1), false)
		if len(ts) != 1 || ts[0].Letak != models.LetakTab {
			t.Fatalf("komite %q: akar panel bukan satu layout group Tab: %d unsur", c.komite, len(ts))
		}
		var tab []string
		for _, b := range ts[0].Anak {
			tab = append(tab, b.Label)
		}
		if strings.Join(tab, "|") != c.mau {
			t.Errorf("komite %q: tab %v, mau %s", c.komite, tab, c.mau)
		}
		akhir := ts[0].Anak[len(ts[0].Anak)-1]
		for _, id := range []string{"SimpanAdjustment", "SendToCommitte"} {
			if _, ok := cariTata(akhir.Anak, func(x models.Tata) bool { return x.ID == id }); !ok {
				t.Errorf("komite %q: tombol %s tidak di tab Adjusment Detail", c.komite, id)
			}
		}
		if _, ok := cariTata(akhir.Anak, func(x models.Tata) bool { return x.Label == "Committee Accept Status" }); !ok {
			t.Errorf("komite %q: Committee Accept Status tidak di tab Adjusment Detail", c.komite)
		}
	}
	for k, v := range map[string]string{"0": "Transfer to committee", "1": "Approve", "2": "Reject"} {
		if models.LabelKode["AcceptanceStatus"][k] != v {
			t.Errorf("AcceptanceStatus %s = %q, mau %q", k, models.LabelKode["AcceptanceStatus"][k], v)
		}
	}
}

func TestBackSurveyorTersembunyiSesudahAkseptasi(t *testing.T) {
	// Work owner 11-10-2026 "TOMBOL BACK KALO UDAH ADA AKSEPTASI, HIDE AJA": tombol Back (BackToEstimasi) layar Choose
	// Surveyor hilang bila ada adjustment ber-AcceptanceStatus 1; Close Claim tetap.
	back := func(x models.Tata) bool { return x.ID == "Back" }
	tutup := func(x models.Tata) bool { return x.ID == "CloseClaim" }
	for _, c := range []struct {
		status string
		mau    bool
	}{{"", true}, {"2", true}, {"1", false}} {
		h := models.HalamanBaru()
		h.SetelDaftar(models.DaftarObjek, []models.Baris{{}})
		h.SetelDaftar(models.DaftarItem(1), []models.Baris{{}})
		h.SetelDaftar(models.DaftarAdj(1, 1), []models.Baris{{"AcceptanceStatus": c.status}})
		ts := models.Evaluasi(h, models.LayarSurveyor(), false)
		if _, ada := cariTata(ts, back); ada != c.mau {
			t.Errorf("AcceptanceStatus %q: Back tampil=%v, mau %v", c.status, ada, c.mau)
		}
		if _, ada := cariTata(ts, tutup); !ada {
			t.Errorf("AcceptanceStatus %q: Close Claim hilang", c.status)
		}
	}
}

func TestKursDanTSIRNMIDRGridItemHanyaBaca(t *testing.T) {
	// Work owner 11-10-2026 "INI READ ONLY. BUAT FIRE ANEKA DAN BISNIS LAIN NYA!": kolom Currency Value (IDR) dan TSI RNM
	// in IDR grid item (Input Estimasi dan Adjustment) hanya-baca di semua lini (XML grid Fire: Text Input).
	for _, lini := range []string{"Fire", "Liability", "GolfInsurance", "MBUCar", "PA", "Travel"} {
		h := models.HalamanBaru()
		h.Setel(models.OQ+"BusinessType", lini)
		h.SetelDaftar(models.DaftarObjek, []models.Baris{{}})
		h.SetelDaftar(models.DaftarItem(1), []models.Baris{{"KursObjectItem": "1", "ValueTSINusareIDR": "50000000"}})
		for nama, layar := range map[string][]models.Unsur{"estimasi": models.LayarObjekEstimasi(1),
			"adjustment": models.LayarObjekAdj(1)} {
			g, ok := cariTata(models.Evaluasi(h, layar, false), func(x models.Tata) bool {
				return x.Jenis == models.JenisGrid && x.Jalur == models.DaftarItem(1)
			})
			if !ok || len(g.Baris) != 1 {
				t.Fatalf("%s %s: grid item tidak ada", lini, nama)
			}
			ada := 0
			for j, k := range g.Kolom {
				if k.Jalur != "KursObjectItem" && k.Jalur != "ValueTSINusareIDR" {
					continue
				}
				ada++
				if sel := g.Baris[0][j]; sel.Tampil && !sel.HanyaBaca {
					t.Errorf("%s %s: kolom %q (%s) dapat diketik", lini, nama, k.Label, k.Jalur)
				}
			}
			if ada != 2 {
				t.Errorf("%s %s: kolom kurs / TSI RNM in IDR %d, mau 2", lini, nama, ada)
			}
		}
	}
}

func TestCoverageItemTerkunciSesudahFaceSheet(t *testing.T) {
	// Work owner 11-10-2026 "KENAPA MASIH BISA EDIT SETELAH KLIK PRINT CLAIM FACESHEET?!": Coverage Name grid item
	// terkunci begitu item tercetak CFS (`RO=.PrintFaceClaim = 1`, semua lini) dan nonaktif menurut NA XML-nya - Fire
	// `IsAdjustment = 1 || .IsFacretro = 1 || IsError > 1`, Aneka / Golf `IsError > 1`.
	sel := func(lini string, item models.Baris, isi func(*models.Halaman)) (models.SelTata, bool) {
		h := models.HalamanBaru()
		h.Setel(models.OQ+"BusinessType", lini)
		h.SetelDaftar(models.DaftarObjek, []models.Baris{{}})
		h.SetelDaftar(models.DaftarItem(1), []models.Baris{item})
		if isi != nil {
			isi(h)
		}
		g, ok := cariTata(models.Evaluasi(h, models.LayarObjekEstimasi(1), false), func(x models.Tata) bool {
			return x.Jenis == models.JenisGrid && x.Jalur == models.DaftarItem(1)
		})
		if !ok {
			t.Fatalf("%s: grid item tidak ada", lini)
		}
		for j, k := range g.Kolom {
			if k.Label == "Coverage Name" && k.Aksi != "" {
				return g.Baris[0][j], true
			}
		}
		return models.SelTata{}, false
	}
	for _, lini := range []string{"Fire", "Liability", "GolfInsurance", "MBUCar", "PA", "Travel"} {
		s, ok := sel(lini, models.Baris{}, nil)
		if lini == "MBUCar" || lini == "PA" || lini == "Travel" { // Coverage ID (bukan Coverage Name) yang dipilih
			continue
		}
		if !ok || s.HanyaBaca || s.Nonaktif {
			t.Fatalf("%s: Coverage Name item belum tercetak harus dapat dipilih: %+v %v", lini, s, ok)
		}
		if s, _ = sel(lini, models.Baris{"PrintFaceClaim": "1"}, nil); !s.HanyaBaca {
			t.Errorf("%s: Coverage Name item tercetak CFS masih dapat diubah", lini)
		}
		if s, _ = sel(lini, models.Baris{}, func(h *models.Halaman) { h.Setel(models.JalurIsError, "2") }); !s.Nonaktif {
			t.Errorf("%s: Coverage Name saat IsError > 1 tidak nonaktif", lini)
		}
	}
	if s, _ := sel("Fire", models.Baris{"IsFacretro": "1"}, nil); !s.Nonaktif {
		t.Error("Fire: Coverage Name item IsFacretro 1 tidak nonaktif")
	}
	if s, _ := sel("Fire", models.Baris{}, func(h *models.Halaman) { h.Setel(models.JalurIsAdjustment, "1") }); !s.Nonaktif {
		t.Error("Fire: Coverage Name saat IsAdjustment 1 tidak nonaktif")
	}
	// Coverage ID MBU / PA / Travel: RO=.PrintFaceClaim = 1
	for _, lini := range []string{"MBUCar", "PA", "Travel"} {
		h := models.HalamanBaru()
		h.Setel(models.OQ+"BusinessType", lini)
		h.SetelDaftar(models.DaftarObjek, []models.Baris{{}})
		h.SetelDaftar(models.DaftarItem(1), []models.Baris{{"PrintFaceClaim": "1"}})
		g, _ := cariTata(models.Evaluasi(h, models.LayarObjekEstimasi(1), false), func(x models.Tata) bool {
			return x.Jenis == models.JenisGrid && x.Jalur == models.DaftarItem(1)
		})
		for j, k := range g.Kolom {
			if k.Aksi == "PilihCoverageObjek" && !g.Baris[0][j].HanyaBaca {
				t.Errorf("%s: Coverage ID item tercetak CFS masih dapat diubah", lini)
			}
		}
	}
}

func TestKolomNomorBarisBerjudulNo(t *testing.T) {
	// Work owner 11-10-2026 "KODE INI PERBAIKI. TADIKAN SUDAH DIPERINTAHKAN PERBAIKI UNTUK SEMUA BISNIS!": nomor baris
	// grid Adjustment (`.pxListSubscript`, judul XML kosong) terbaca sebagai nilai kolom Status di sebelahnya. Kolom nomor
	// grid panel adjustment berjudul "No" di semua lini; Status baris disetujui = label, bukan kode.
	if got := models.LabelKode["AcceptanceStatus"]["1"]; got != "Approve" {
		t.Fatalf("label AcceptanceStatus 1 = %q", got)
	}
	for _, lini := range []string{"Fire", "Liability", "GolfInsurance", "MBUCar", "PA", "Travel", "MarineCargo"} {
		h := models.HalamanBaru()
		h.Setel(models.OQ+"BusinessType", lini)
		h.SetelDaftar(models.DaftarObjek, []models.Baris{{}})
		h.SetelDaftar(models.DaftarItem(1), []models.Baris{{}})
		h.SetelDaftar(models.DaftarAdj(1, 1), []models.Baris{{"IsApproved": "1", "AcceptanceStatus": "1", "IsKomite": "1"}})
		h.SetelDaftar(models.DaftarDiAdj(1, 1, 1, models.AnakAdjSpread), []models.Baris{{"TreatyType": "UJI"}})
		h.SetelDaftar(models.DaftarDiAdj(1, 1, 1, models.AnakAdjQS), []models.Baris{{"TreatyName": "UJI QS"}})
		ts := models.Evaluasi(h, append(models.LayarItemAdj(1, 1), models.LayarAdjustment(1, 1, 1)...), false)
		var grids []models.Tata
		var jalan func([]models.Tata)
		jalan = func(xs []models.Tata) {
			for _, x := range xs {
				if x.Jenis == models.JenisGrid {
					grids = append(grids, x)
				}
				jalan(x.Anak)
			}
		}
		jalan(ts)
		for _, g := range grids {
			for j, k := range g.Kolom {
				if k.Jalur == models.PropNomorBaris && k.Label != "No" {
					t.Errorf("%s: grid %s kolom nomor ke-%d berjudul %q, mau No", lini, g.Jalur, j, k.Label)
				}
				if k.Jalur == "AcceptanceStatus" && (k.Sumber != "kode:AcceptanceStatus" || !g.Baris[0][j].Tampil) {
					t.Errorf("%s: Status baris disetujui tanpa label kode (%q, tampil %v)", lini, k.Sumber, g.Baris[0][j].Tampil)
				}
			}
		}
	}
}
