package medis

// Seam 6 - skoring medis. Tiket E20.
//
// ⛔ Tidak satu pun ambang klinis disalin ke berkas ini (E20). Yang diuji:
// keutuhan tabel bangkitan, jalannya SELURUH aturan korpus atas nilai
// sintetis, perilaku yang tidak bergantung ambang (serologi, nilai kosong,
// lab tak dikenal), dan semantik transisi atas pohon sintetis.
//
// Dibaca sesudah: medis.go.

import (
	"errors"
	"os"
	"regexp"
	"strings"
	"testing"
)

func jelajah(ls []langkah, f func(langkah)) {
	for _, l := range ls {
		f(l)
		jelajah(l.anak, f)
	}
}

// TestTabelBangkitanUtuh - cacah tabel = cacah korpus, dan SETIAP ekspresi
// korpus terurai oleh penafsir. Ekspresi di luar subset akan gagal di sini,
// bukan diam-diam di produksi.
//
// ⚠️ Satu kejanggalan korpus dikenali dengan NAMA, bukan dilonggarkan:
// `[terverifikasi]` 56 penugasan `SetParamLab_Act` bernilai `Negatif/0` TANPA
// kutip (dua cara: urai pohon dan `grep '<PropertiesValue>[^"<]'`). Pega
// membacanya sebagai ekspresi, dan artinya belum terverifikasi; di sini ia
// galat sintaks per properti, tidak ditebak sebagai teks.
//
// ⛔ Cacah pembanding di bawah BUKAN keluaran pembangkit: ia hasil `grep`
// mentah atas berkas korpus (cara kedua, CLAUDE.md §4a), dicatat tangan
// 01-10-2026:
//
//	grep -c '<pxObjClass>Embed-ActivitySteps</pxObjClass>' <berkas>   # langkah
//	grep -c '<PropertiesName' <berkas>                                # penugasan (bernama + kosong)
//	grep -c '<pyStepsBlockName>//</pyStepsBlockName>' <berkas>         # langkah di-remark
//
// Pembangkit yang salah urai memberi angka lain di sini.
func TestTabelBangkitanUtuh(t *testing.T) {
	for _, u := range []struct {
		a                                      aktivitas
		grepLangkah, grepPenugasan, grepRemark int
		tanpaKutipSah                          int
	}{
		{aturanSkor, 210, 338, 16, 0},
		{aturanFisik, 3, 7, 0, 0},
		{aturanParamLab, 9, 438, 0, 56},
	} {
		if u.a.cacahLangkah != u.grepLangkah || u.a.cacahPenugasan != u.grepPenugasan || u.a.cacahRemark != u.grepRemark {
			t.Errorf("%s: pembangkit %d/%d/%d, grep %d/%d/%d", u.a.asal, u.a.cacahLangkah, u.a.cacahPenugasan, u.a.cacahRemark,
				u.grepLangkah, u.grepPenugasan, u.grepRemark)
		}
		remark := 0
		jelajah(u.a.langkah, func(l langkah) {
			if l.diRemark {
				remark++
			}
		})
		if remark != u.grepRemark {
			t.Errorf("%s: %d langkah ditandai di-remark di tabel, grep %d", u.a.asal, remark, u.grepRemark)
		}
		if n := periksaAktivitas(t, u.a); n != u.tanpaKutipSah {
			t.Errorf("%s: %d nilai %q, mau %d", u.a.asal, n, nilaiTanpaKutip, u.tanpaKutipSah)
		}
	}
	if n := len(namaSkor()); n < 20 {
		t.Errorf("hanya %d properti skor; pembangkitnya rusak?", n)
	}
}

// nilaiTanpaKutip - satu-satunya nilai tak terurai yang dikenal di korpus.
const nilaiTanpaKutip = "Negatif/0"

// periksaAktivitas - cacah tabel = cacah korpus, setiap ekspresi terurai
// kecuali `nilaiTanpaKutip`; mengembalikan cacah kemunculan yang terakhir itu.
func periksaAktivitas(t *testing.T, a aktivitas) (tanpaKutip int) {
	t.Helper()
	nLangkah, nTugas := 0, 0
	jelajah(a.langkah, func(l langkah) {
		nLangkah++
		nTugas += len(l.penugasan)
		for _, b := range l.prakondisi {
			if _, err := urai(b.ekspresi); err != nil {
				t.Errorf("langkah %s prakondisi: %v", l.nomor, err)
			}
		}
		for _, p := range l.penugasan {
			if p.nama == "" && p.ekspresi == "" {
				continue
			}
			if _, err := urai(p.ekspresi); err != nil {
				if p.ekspresi == nilaiTanpaKutip && errors.Is(err, ErrSintaks) {
					tanpaKutip++
					continue
				}
				t.Errorf("langkah %s %s: %v", l.nomor, p.nama, err)
			}
		}
	})
	if nLangkah != a.cacahLangkah || nTugas != a.cacahPenugasan {
		t.Errorf("%s: tabel %d langkah / %d penugasan, korpus %d / %d", a.asal, nLangkah, nTugas, a.cacahLangkah, a.cacahPenugasan)
	}
	return tanpaKutip
}

// TestNilaiRujukanLab - `SetParamLab_Act`: nilai rujukan terisi per lab;
// `Negatif/0` menjadi galat, bukan teks; lab kosong menghapus halaman.
func TestNilaiRujukanLab(t *testing.T) {
	for _, sex := range []string{"1", "2"} {
		for _, lab := range []string{"Prodia", "Pramita", "Biotest"} {
			nilai, galat, err := NilaiRujukanLab(sex, lab)
			if err != nil {
				t.Fatal(err)
			}
			if len(nilai) == 0 {
				t.Errorf("%s/%s: nol nilai rujukan", sex, lab)
			}
			for _, g := range galat {
				if !errors.Is(g, ErrSintaks) {
					t.Errorf("%s/%s: galat tak terduga %v", sex, lab, g)
				}
			}
		}
	}
	nilai, galat, err := NilaiRujukanLab("1", "")
	if err != nil || len(nilai) != 0 || len(galat) != 0 {
		t.Errorf("lab kosong: nilai %d, galat %v, err %v", len(nilai), galat, err)
	}
}

// TestPeriksaFisik - BMI eksak dinilai; tak eksak tidak ditebak (skala
// `@divide` tanpa argumen belum terverifikasi). Nilai sintetis.
func TestPeriksaFisik(t *testing.T) {
	// 100 ÷ (200/100)² = 25, eksak.
	h, err := PeriksaFisik(MasukanFisik{JenisKelamin: "1", Pemeriksaan: map[string]string{"Weight": "100", "Height": "200"}})
	if err != nil {
		t.Fatal(err)
	}
	if h.Nilai["BM_Ratio"] != "25" || h.Nilai["BMI_Ratio_Note"] == "" || h.Nilai["EM"] == "" {
		t.Errorf("BMI eksak: %v / galat %v", h.Nilai, h.Galat)
	}
	// 100 ÷ 3² tidak eksak.
	h, err = PeriksaFisik(MasukanFisik{JenisKelamin: "1", Pemeriksaan: map[string]string{"Weight": "100", "Height": "300"}})
	if err != nil {
		t.Fatal(err)
	}
	if _, ada := h.Nilai["BM_Ratio"]; ada || len(h.Galat) == 0 || !errors.Is(h.Galat[0], ErrSkalaBawaanBelumTerverifikasi) {
		t.Errorf("BMI tak eksak: %v / galat %v", h.Nilai, h.Galat)
	}
}

// namaPemeriksaan - nama properti pemeriksaan yang DIBACA aturan.
func namaPemeriksaan() []string {
	set := map[string]bool{}
	jelajah(aturanSkor.langkah, func(l langkah) {
		for _, p := range l.penugasan {
			if strings.HasPrefix(p.nama, awalanPemeriksaan) && !strings.HasPrefix(p.nama, awalanSkor) {
				set[strings.TrimPrefix(p.nama, awalanPemeriksaan)] = true
			}
		}
	})
	var h []string
	for k := range set {
		h = append(h, k)
	}
	return h
}

// TestSeluruhAturanBerjalan - setiap jenis kelamin × lab, semua pemeriksaan
// diminta, nilai sintetis "1": tidak ada galat struktur, sintaks, atau
// boolean - hanya galat yang memang dirancang (nilai/pembulatan).
func TestSeluruhAturanBerjalan(t *testing.T) {
	periksa := map[string]string{}
	for _, n := range namaPemeriksaan() {
		periksa[n] = "1"
	}
	for _, sex := range []string{"1", "2"} {
		for _, lab := range []string{"Prodia", "Pramita", "Biotest"} {
			h, err := ScoreMedical(Masukan{JenisKelamin: sex, Lab: lab, Usia: "30", JN: "All", JNSub: "All",
				Pemeriksaan: periksa, ParamLab: periksa})
			if err != nil {
				t.Fatalf("%s/%s: %v", sex, lab, err)
			}
			for _, g := range h.Galat {
				if errors.Is(g.Err, ErrSintaks) || errors.Is(g.Err, ErrBukanBoolean) {
					t.Errorf("%s/%s: %v", sex, lab, g)
				}
			}
			if len(h.Skor) == 0 {
				t.Errorf("%s/%s: nol skor", sex, lab)
			}
		}
	}
}

// TestSerologiReaktifDitolak - E20: keluaran "ditolak" dapat muncul dan ia
// KEPUTUSAN, bukan galat.
func TestSerologiReaktifDitolak(t *testing.T) {
	for _, sex := range []string{"1", "2"} {
		h, err := ScoreMedical(Masukan{JenisKelamin: sex, Lab: "Prodia", JN: "Serologi",
			Pemeriksaan: map[string]string{"HbsAg": "Reaktif", "HbeAg": "UJI-NONREAKTIF", "AntiHIV": "UJI-NONREAKTIF"}})
		if err != nil {
			t.Fatal(err)
		}
		if h.Skor["HbsAg"] != "Decline" || h.Skor["HbeAg"] != "Standard" || h.Skor["AntiHIV"] != "Standard" {
			t.Errorf("jenis kelamin %s: skor %v", sex, h.Skor)
		}
	}
}

// TestNilaiKosongTidakMenjadiStandar - pemeriksaan yang diminta tetapi
// nilainya kosong TIDAK diberi skor tebakan.
func TestNilaiKosongTidakMenjadiStandar(t *testing.T) {
	h, err := ScoreMedical(Masukan{JenisKelamin: "1", Lab: "Prodia", JN: "Tumor Marker", JNSub: "CEA",
		Pemeriksaan: map[string]string{"CEA": ""}})
	if err != nil {
		t.Fatal(err)
	}
	if v, ada := h.Skor["CEA"]; ada {
		t.Errorf("skor CEA %q padahal nilainya kosong", v)
	}
	if len(h.Galat) == 0 || !errors.Is(h.Galat[0], ErrNilaiBukanAngka) {
		t.Errorf("galat %v, mau ErrNilaiBukanAngka", h.Galat)
	}
}

func TestLabTakDikenalTanpaSkor(t *testing.T) {
	h, err := ScoreMedical(Masukan{JenisKelamin: "1", Lab: "UJI-LAB", JN: "All", JNSub: "All",
		Pemeriksaan: map[string]string{"HbsAg": "Reaktif"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(h.Skor) != 0 || len(h.Galat) != 0 {
		t.Errorf("lab tak dikenal memberi skor %v / galat %v", h.Skor, h.Galat)
	}
}

// TestTransisiPrakondisi - semantik kode transisi atas pohon sintetis.
func TestTransisiPrakondisi(t *testing.T) {
	atau := []barisPrakondisi{{`param.JN=="All"`, "5", "2"}, {`param.JN=="X"`, "2", "3"}}
	dan := []barisPrakondisi{{`param.JN=="All"`, "2", "3"}, {`param.JN_Sub=="Y"`, "2", "3"}}
	for _, u := range []struct {
		nama   string
		baris  []barisPrakondisi
		jn, js string
		mau    aksiLangkah
	}{
		{"atau: A", atau, "All", "", aksiJalankan},
		{"atau: B", atau, "X", "", aksiJalankan},
		{"atau: bukan keduanya", atau, "Z", "", aksiLewati},
		{"dan: keduanya", dan, "All", "Y", aksiJalankan},
		{"dan: satu", dan, "All", "Z", aksiLewati},
	} {
		h := &halamanKerja{nilai: map[string]string{"param.JN": u.jn, "param.JN_Sub": u.js}, gagal: map[string]error{}}
		got, err := putuskan(langkah{berprakondisi: true, prakondisi: u.baris}, h)
		if err != nil || got != u.mau {
			t.Errorf("%s: %v %v, mau %v", u.nama, got, err, u.mau)
		}
	}
	// Prakondisi nonaktif (P-11) → tetap jalan, walau syaratnya salah.
	h := &halamanKerja{nilai: map[string]string{}, gagal: map[string]error{}}
	if got, _ := putuskan(langkah{berprakondisi: false, prakondisi: []barisPrakondisi{{`1==2`, "2", "3"}}}, h); got != aksiJalankan {
		t.Errorf("pyStepsPreCondition=false: %v", got)
	}
	// Kode tak dikenal → galat struktur.
	if _, err := putuskan(langkah{berprakondisi: true, prakondisi: []barisPrakondisi{{`1==1`, "4", "3"}}}, h); !errors.Is(err, ErrTransisiTakDikenal) {
		t.Errorf("galat %v, mau ErrTransisiTakDikenal", err)
	}
	// "6" menghentikan seluruh activity.
	pohon := []langkah{
		{nomor: "1", berprakondisi: true, prakondisi: []barisPrakondisi{{`1==1`, "6", "2"}}},
		{nomor: "2", penugasan: []penugasan{{".X", `"jalan"`}}},
	}
	var galat []Galat
	if keluar, err := jalankan(pohon, h, &galat); err != nil || !keluar || h.nilai[".X"] != "" {
		t.Errorf("keluar=%v err=%v X=%q", keluar, err, h.nilai[".X"])
	}
}

// TestGalatMenjalar - properti yang gagal tidak dibaca sebagai kosong oleh
// ekspresi sesudahnya.
func TestGalatMenjalar(t *testing.T) {
	pohon := []langkah{{nomor: "1", penugasan: []penugasan{
		{".A", `@toDecimal(.Kosong)`},
		{".B", `@if(.A==1,"Y","Z")`},
	}}}
	h := &halamanKerja{nilai: map[string]string{}, gagal: map[string]error{}}
	var galat []Galat
	if _, err := jalankan(pohon, h, &galat); err != nil {
		t.Fatal(err)
	}
	if len(galat) != 2 || h.nilai[".B"] != "" {
		t.Errorf("galat %v, B=%q", galat, h.nilai[".B"])
	}
}

// TestUsiaMemilihCabang - E20: skor yang bergantung usia bercabang menurut
// usia. Nilai pemeriksaan "0" (sintetis) dipakai untuk dua usia; yang ditegaskan
// hanya bahwa keduanya dinilai dan hasilnya BERBEDA - tidak satu ambang pun
// disalin ke test.
func TestUsiaMemilihCabang(t *testing.T) {
	skor := func(usia string) string {
		h, err := ScoreMedical(Masukan{JenisKelamin: "1", Lab: "Prodia", Usia: usia, JN: "Blood Glucose", JNSub: "HbA1c (NGSP)",
			Pemeriksaan: map[string]string{"HbA1c_NGSP": "0"}})
		if err != nil {
			t.Fatal(err)
		}
		return h.Skor["HbA1c_NGSP"]
	}
	muda, tua := skor("18"), skor("45")
	if muda == "" || tua == "" || muda == tua {
		t.Errorf("usia 18 → %q, usia 45 → %q; mau dua skor berbeda", muda, tua)
	}
}

// TestLangkahDiRemarkTidakDijalankan - butir 43: langkah berlabel `//` beserta
// sub-langkahnya tidak dijalankan. Di korpus, SEMUA penugasan Score.SGOT,
// Score.SGPT, dan Score.GammaGT ada di langkah ber-`//` - jadi tidak pernah
// terisi, walau pemeriksaannya diminta spesifik.
func TestLangkahDiRemarkTidakDijalankan(t *testing.T) {
	for _, js := range []string{"SGOT / AST", "SGPT / ALT", "Gamma GT"} {
		h, err := ScoreMedical(Masukan{JenisKelamin: "1", Lab: "Prodia", JN: "Liver Function Tes", JNSub: js,
			Pemeriksaan: map[string]string{"SGOT": "1", "SGPT": "1", "GammaGT": "1"}})
		if err != nil {
			t.Fatal(err)
		}
		for _, skor := range []string{"SGOT", "SGPT", "GammaGT"} {
			if v, ada := h.Skor[skor]; ada {
				t.Errorf("%s: Score.%s terisi %q dari langkah di-remark", js, skor, v)
			}
		}
	}
	// Pohon sintetis: sub-langkah langkah ber-remark ikut tidak jalan.
	pohon := []langkah{{nomor: "1", diRemark: true, penugasan: []penugasan{{".A", `"x"`}},
		anak: []langkah{{nomor: "1.1", penugasan: []penugasan{{".B", `"y"`}}}}}}
	h := &halamanKerja{nilai: map[string]string{}, gagal: map[string]error{}}
	var galat []Galat
	if _, err := jalankan(pohon, h, &galat); err != nil || len(h.nilai) != 0 {
		t.Errorf("err %v nilai %v", err, h.nilai)
	}
}

// TestTransisiKorpusTigaPola - cara pertama sensus pola transisi yang dikutip
// medis.go (cara kedua: grep, di komentar itu): satu baris (2,3) ×30, dua
// baris (2,3)(2,3) ×10, dua baris (5,2)(2,3) ×170.
func TestTransisiKorpusTigaPola(t *testing.T) {
	pola := map[string]int{}
	jelajah(aturanSkor.langkah, func(l langkah) {
		k := ""
		for _, b := range l.prakondisi {
			k += "(" + b.bilaBenar + "," + b.bilaSalah + ")"
		}
		if k != "" {
			pola[k]++
		}
	})
	mau := map[string]int{"(2,3)": 30, "(2,3)(2,3)": 10, "(5,2)(2,3)": 170}
	if len(pola) != len(mau) {
		t.Errorf("pola %v, mau %v", pola, mau)
	}
	for k, n := range mau {
		if pola[k] != n {
			t.Errorf("pola %s ×%d, mau ×%d", k, pola[k], n)
		}
	}
}

// TestAturanBangkitanTanpaNomorPolis - berkas bangkitan tidak memuat pola
// nomor polis produksi (CLAUDE.md §4.10). Pembangkit juga menolaknya; test ini
// menjaga berkas yang SUDAH ada.
func TestAturanBangkitanTanpaNomorPolis(t *testing.T) {
	isi, err := os.ReadFile("aturan_gen.go")
	if err != nil {
		t.Fatal(err)
	}
	if regexp.MustCompile(`(?i)RNM-[A-Z0-9]`).Match(isi) {
		t.Error("aturan_gen.go memuat pola nomor polis")
	}
}
