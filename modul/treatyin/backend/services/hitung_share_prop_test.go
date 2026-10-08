package services

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"

	"nusantarare/modul/treatyin/backend/models"
)

// sumberPropUji - kedua RD spreading tiruan; mencatat parameter RD anak.
type sumberPropUji struct {
	induk   []models.SusunanSpreading
	anak    map[string][]models.SusunanSpreading
	panggil [][]string
}

func (s *sumberPropUji) Induk(_, _ string) []models.SusunanSpreading { return s.induk }

func (s *sumberPropUji) AnakProp(tahun, grup, desc, induk, tahunID string) []models.SusunanSpreading {
	s.panggil = append(s.panggil, []string{tahun, grup, desc, induk, tahunID})
	return s.anak[induk]
}

// pohonUji - Limits dari JSON, seperti yang layar kirim.
func pohonUji(t *testing.T, j string) []map[string]any {
	t.Helper()
	var out []map[string]any
	if err := json.Unmarshal([]byte(j), &out); err != nil {
		t.Fatal(err)
	}
	return out
}

func shareProp(t *testing.T, m MasukanShareProp, src SumberSpreadingProp) HasilShareProp {
	t.Helper()
	h, err := HitungShareProp(m, src)
	if err != nil {
		t.Fatalf("HitungShareProp: %v", err)
	}
	return h
}

func detailUji(h HasilShareProp, i, j int) map[string]any {
	return larikSimpul(h.Limits[i], "Detail")[j]
}

func nilaiUji(xs []map[string]any) []string {
	out := []string{}
	for _, x := range xs {
		out = append(out, teksSimpul(x, "Currency")+" "+teksSimpul(x, "Value"))
	}
	return out
}

func totalUji(xs []NilaiMataUang) []string {
	out := []string{}
	for _, x := range xs {
		out = append(out, x.Currency+" "+x.Value)
	}
	return out
}

const limitsGambar17 = `[{"ID":"1","TreatyType":"SPECIAL SURPLUS","Detail":[{"TreatyGroup":"ENGINEERING","TreatyGroupID":"10002","TreatyType":"SPECIAL SURPLUS","QSPct":"tetap","SpreadingTypeID":"P24","SpreadingType":"2024 QS 155M TRT","CessionList":[{"Currency":"IDR","Value":"10000000000"}],"IOOLimitList":[{"Currency":"IDR","Value":"20000000000"}]}]}]`

func sumberGambar17() *sumberPropUji {
	return &sumberPropUji{anak: map[string][]models.SusunanSpreading{"P24": {
		{ReinsTypeID: "R1", ReinsTypeName: "QS (R/I)", ParentReinsTypeID: "P24", Pct: "60"},
		{ReinsTypeID: "R2", ReinsTypeName: "QS (OR)", ParentReinsTypeID: "P24", Pct: "40"},
	}}}
}

// TestSharePropGambar17 - Refresh atas Limits yang tab Limits isi: RNM Share
// = Cession × % RNM Share, spreading 40/60 (FetchQSfromMaster), total
// TIDAK berlipat - gambar Pega 17: 3.000.000.000 · 1.200.000.000 · 1.800.000.000.
func TestSharePropGambar17(t *testing.T) {
	src := sumberGambar17()
	masukan := pohonUji(t, limitsGambar17)
	h := shareProp(t, MasukanShareProp{Aksi: AksiSharePropShare, Limits: masukan, RNMShareP: "30", BrokeragePercentP: "0", OptionLimit: "1", Commencement: "20250101"}, src)
	d := detailUji(h, 0, 0)
	if got := nilaiUji(larikSimpul(d, "RNMShareList")); !reflect.DeepEqual(got, []string{"IDR 3000000000"}) {
		t.Errorf("RNMShareList %v", got)
	}
	if teksSimpul(d, "RNMShare") != "30" || teksSimpul(d, "Brokerage") != "0" || teksSimpul(d, "ShareNote") != " of 100%" {
		t.Errorf("RNMShare/Brokerage/ShareNote %v %v %q", d["RNMShare"], d["Brokerage"], d["ShareNote"])
	}
	if got := nilaiUji(larikSimpul(d, "RNMSpreadedList")); !reflect.DeepEqual(got, []string{"IDR 1200000000"}) {
		t.Errorf("OR %v", got)
	}
	if got := nilaiUji(larikSimpul(d, "RNMSpreadedListRI")); !reflect.DeepEqual(got, []string{"IDR 1800000000"}) {
		t.Errorf("R/I %v", got)
	}
	if teksSimpul(d, "SpreadingTotalPct") != "100" || len(larikSimpul(d, "SpreadingList")) != 2 {
		t.Errorf("spreading %v %v", d["SpreadingTotalPct"], d["SpreadingList"])
	}
	if !reflect.DeepEqual(totalUji(h.TotalShareRnmProp), []string{"IDR 3000000000"}) ||
		!reflect.DeepEqual(totalUji(h.TotalSpreadedRnmProp), []string{"IDR 1200000000"}) ||
		!reflect.DeepEqual(totalUji(h.TotalSpreadedRnmRIProp), []string{"IDR 1800000000"}) {
		t.Errorf("total %v %v %v", h.TotalShareRnmProp, h.TotalSpreadedRnmProp, h.TotalSpreadedRnmRIProp)
	}
	if len(h.Pesan) != 0 {
		t.Errorf("pesan %v", h.Pesan)
	}
	// Kunci yang rumus tidak kenal TIDAK hilang; masukan tidak berubah.
	if teksSimpul(d, "QSPct") != "tetap" || teksSimpul(h.Limits[0], "ID") != "1" {
		t.Errorf("kunci lain hilang: %v", d)
	}
	if _, ada := larikSimpul(masukan[0], "Detail")[0]["RNMShareList"]; ada {
		t.Error("masukan berubah")
	}
	// Refresh TANPA ParentReinsTypeID: RD anak dipanggil tanpa tahun (A, E dilewati).
	if !reflect.DeepEqual(src.panggil, [][]string{{"", "10002", "10001", "P24", ""}}) {
		t.Errorf("parameter RD anak %v", src.panggil)
	}
}

// TestSharePropOpsiDanQuotaShare - Option 2 memakai 100% Limit; QUOTA SHARE
// menulis " of <CessionPct>% of 100%".
func TestSharePropOpsiDanQuotaShare(t *testing.T) {
	src := sumberGambar17()
	h := shareProp(t, MasukanShareProp{Aksi: AksiSharePropShare, Limits: pohonUji(t, limitsGambar17), RNMShareP: "10", OptionLimit: "2"}, src)
	if got := nilaiUji(larikSimpul(detailUji(h, 0, 0), "RNMShareList")); !reflect.DeepEqual(got, []string{"IDR 2000000000"}) {
		t.Errorf("Option 2 %v", got)
	}
	qs := pohonUji(t, `[{"TreatyType":"QUOTA SHARE","Detail":[{"TreatyType":"QUOTA SHARE","CessionPct":"35","SpreadingTypeID":"P24","CessionList":[{"Currency":"USD","Value":"100"}]}]}]`)
	h = shareProp(t, MasukanShareProp{Aksi: AksiSharePropShare, Limits: qs, RNMShareP: "50", OptionLimit: "1"}, src)
	if got := teksSimpul(detailUji(h, 0, 0), "ShareNote"); got != " of 35% of 100%" {
		t.Errorf("ShareNote %q", got)
	}
}

// TestSharePropTanpaSpreadingDisalinApaAdanya - Detail tanpa Spreading Type
// → SetSpreadName: totalnya DITAMBAHKAN dua kali (SetSpreadName [5] lalu
// TreatyInPropshare [6]) dan pesan total spreading terpasang - bunyi Activity.
func TestSharePropTanpaSpreadingDisalinApaAdanya(t *testing.T) {
	limits := pohonUji(t, `[{"TreatyType":"SURPLUS","Detail":[{"TreatyGroup":"FIRE","TreatyType":"SURPLUS","CessionList":[{"Currency":"IDR","Value":"1000"}]}]}]`)
	h := shareProp(t, MasukanShareProp{Aksi: AksiSharePropShare, Limits: limits, RNMShareP: "10", OptionLimit: "1"}, &sumberPropUji{})
	if got := nilaiUji(larikSimpul(detailUji(h, 0, 0), "RNMShareList")); !reflect.DeepEqual(got, []string{"IDR 100"}) {
		t.Errorf("RNMShareList %v", got)
	}
	if !reflect.DeepEqual(totalUji(h.TotalShareRnmProp), []string{"IDR 200"}) {
		t.Errorf("total berlipat (disalin) %v", h.TotalShareRnmProp)
	}
	if !reflect.DeepEqual(h.Pesan, []string{pesanTotalShareSpreading}) {
		t.Errorf("pesan %v", h.Pesan)
	}
}

// ⛔ TERBALIK 7 Oktober 2026 — pesan "Spreading is incomplete at : " TIDAK
// PERNAH LAHIR.
//
// Uji ini dulu MEMAKU pesan itu, dan dengan begitu memaku langkah yang
// ekspornya KOMENTARI. Langkah [8] `TreatyInPropshare`
// (`Property-Set-Messages`, pra-syarat `Local.ErrorFlag=="1"`) ber-
// `pyStepsBlockName == "//"`.
//
// ⚠️ Langkah [7] yang MENYUSUN teksnya tetap hidup — itulah yang membuat
// kekeliruan ini mudah: membaca [7] saja memberi kesan pesannya dipakai.
// Yang membacanya cuma [8], dan [8] mati.
//
// ⭐ Keadaan yang sama tetap diuji: totalnya tetap dihitung, hanya pesannya
// yang tidak ada.
func TestSharePropSpreadingBelumLengkapNolPesan(t *testing.T) {
	src := &sumberPropUji{anak: map[string][]models.SusunanSpreading{"P9": {{ReinsTypeName: "QS (R/I)", Pct: "100"}}}}
	limits := pohonUji(t, `[{"TreatyType":"SURPLUS","Detail":[{"TreatyGroup":"FIRE","SpreadingTypeID":"P9","CessionList":[{"Currency":"IDR","Value":"1000"}]}]}]`)
	h := shareProp(t, MasukanShareProp{Aksi: AksiSharePropShare, Limits: limits, RNMShareP: "10", OptionLimit: "1"}, src)
	for _, p := range h.Pesan {
		if strings.HasPrefix(p, pesanSpreadingBelumLengkap) {
			t.Errorf("pesan langkah MATI [8] lahir: %q", p)
		}
	}
	// Totalnya tetap terhitung — `jumlahkanTotal` tidak ikut dicabut.
	if len(h.TotalShareRnmProp) == 0 {
		t.Error("total akar kosong — `jumlahkanTotal` ikut tercabut")
	}
}

// TestSharePropDetailDanSpreading - % RNM Share satu Detail
// (TreatyInPropshareDetail) dan Spreading Type (FetchQSfromMaster bertahun).
func TestSharePropDetailDanSpreading(t *testing.T) {
	src := sumberGambar17()
	limits := pohonUji(t, limitsGambar17)
	larikSimpul(limits[0], "Detail")[0]["RNMShare"] = "20"
	h := shareProp(t, MasukanShareProp{Aksi: AksiSharePropDetail, Limits: limits, OptionLimit: "1"}, src)
	if got := nilaiUji(larikSimpul(detailUji(h, 0, 0), "RNMShareList")); !reflect.DeepEqual(got, []string{"IDR 2000000000"}) {
		t.Errorf("detail RNMShareList %v", got)
	}
	if !reflect.DeepEqual(totalUji(h.TotalShareRnmProp), []string{"IDR 2000000000"}) {
		t.Errorf("total %v", h.TotalShareRnmProp)
	}

	src = sumberGambar17()
	src.induk = []models.SusunanSpreading{{ReinsTypeID: "P24", ReinsTypeName: "2024 QS 155M TRT", TreatyYearID: "TY24", TreatyYear: "2024"}}
	h = shareProp(t, MasukanShareProp{Aksi: AksiSharePropSpreading, Limits: h.Limits}, src)
	if !reflect.DeepEqual(src.panggil, [][]string{{"2024", "10002", "10001", "P24", "TY24"}}) {
		t.Errorf("RD anak bertahun %v", src.panggil)
	}
	if got := nilaiUji(larikSimpul(detailUji(h, 0, 0), "RNMSpreadedList")); !reflect.DeepEqual(got, []string{"IDR 800000000"}) {
		t.Errorf("OR %v", got)
	}
}

// TestSharePropSebarManual - SetSpreadName: rincian per anak × mata uang,
// QS (OR)/ORS ke OR, QS (R/I) ke R/I.
func TestSharePropSebarManual(t *testing.T) {
	src := &sumberPropUji{
		induk: []models.SusunanSpreading{{ReinsTypeID: "M1", ReinsTypeName: "2025 QS TRT", TreatyYearID: "TY25"}},
		anak: map[string][]models.SusunanSpreading{"M1": {
			{ReinsTypeID: "a", ReinsTypeName: "QS (OR)", Pct: "40"},
			{ReinsTypeID: "b", ReinsTypeName: "QS (R/I)", Pct: "60"},
		}},
	}
	limits := pohonUji(t, `[{"TreatyType":"SURPLUS","Detail":[{"TreatyGroup":"FIRE","RNMShare":"10","SpreadingList":[{"ReinsTypeID":"M1","Pct":"10"}],"RNMShareList":[{"Currency":"IDR","Value":"1000"}]}]}]`)
	h := shareProp(t, MasukanShareProp{Aksi: AksiSharePropSebarNama, Limits: limits}, src)
	d := detailUji(h, 0, 0)
	if got := nilaiUji(larikSimpul(d, "RNMSpreadedList")); !reflect.DeepEqual(got, []string{"IDR 400"}) {
		t.Errorf("OR %v", got)
	}
	if got := nilaiUji(larikSimpul(d, "RNMSpreadedListRI")); !reflect.DeepEqual(got, []string{"IDR 600"}) {
		t.Errorf("R/I %v", got)
	}
	r := larikSimpul(d, "SpreadingList")[0]
	if teksSimpul(r, "ReinsTypeName") != "2025 QS TRT" || len(larikSimpul(r, "BreakDownSprdList")) != 2 {
		t.Errorf("baris spreading %v", r)
	}
	if teksSimpul(d, "SpreadingTotalPct") != "10" || len(h.Pesan) != 0 {
		t.Errorf("total pct %v pesan %v", d["SpreadingTotalPct"], h.Pesan)
	}
	if !reflect.DeepEqual(src.panggil, [][]string{{"", "", "", "M1", "TY25"}}) {
		t.Errorf("RD anak %v", src.panggil)
	}
}

func TestSharePropDitolak(t *testing.T) {
	for _, m := range []MasukanShareProp{{Aksi: "hapus"}, {Aksi: AksiSharePropDetail, IndeksLimit: 3}} {
		if _, err := HitungShareProp(m, &sumberPropUji{}); !errors.Is(err, ErrMasukanTidakSah) {
			t.Errorf("%+v: %v", m, err)
		}
	}
}

// ⛔ LANGKAH [2] `CalculateShareList` MATI — `RNMShareAcrossTheBoard`
// TIDAK BOLEH mengubah apa pun di cabang proporsional.
//
// Langkah [2] (`RNMShareAcrossTheBoard == true` → `CessionList` ×
// `RNMShareP`) ber-`pyStepsBlockName == "//"` di `Treaty In` MAUPUN
// `Treaty In Adjustment`. Selama ia dibangun, `CessionList` masuk DUA KALI
// — sekali di [2] ber-`RNMShareP`, sekali di [3] ber-`.RNMShare` — dan
// keduanya bernilai sama sebab langkah 5 `TreatyInPropshare` menyalin
// `.RNMShare = RNMShareP`.
//
// ⚠️ Seluruh uji lama DIAM soal ini sebab nol satu pun yang menyetel
// `RNMShareAcrossTheBoard`; nilainya kosong, dan cabang itu tidak pernah
// menyala. Itulah sebabnya cacatnya hidup sampai pemakai melihatnya.
//
// ⭐ Angkanya dari laporan pemakai 7 Oktober 2026 — bukan karangan.
func TestSharePropLintasBoardTidakMenggandakan(t *testing.T) {
	const limits = `[{"ID":"1","TreatyType":"QUOTA SHARE","Detail":[{"TreatyGroup":"PROPERTY","TreatyGroupID":"10007","TreatyType":"QUOTA SHARE","QSPct":"100","CessionPct":"100","RetentionPct":"0","IOOLimitList":[{"Currency":"IDR","Value":"1757675000000"}],"RetentionList":[{"Currency":"IDR","Value":"0"}],"CessionList":[{"Currency":"IDR","Value":"1757675000000"}]}]}]`
	jalan := func(board string) []string {
		return nilaiUji(larikSimpul(detailUji(shareProp(t, MasukanShareProp{
			Aksi: AksiSharePropShare, Limits: pohonUji(t, limits),
			RNMShareP: "1.28", OptionLimit: "1", RNMShareAcrossTheBoard: board,
			Commencement: "20250101",
		}, &sumberPropUji{}), 0, 0), "RNMShareList"))
	}
	// Cession 1.757.675.000.000 × 1,28% = 22.498.240.000 — SATU baris.
	mau := []string{"IDR 22498240000"}
	for _, board := range []string{"true", "false", ""} {
		if got := jalan(board); !reflect.DeepEqual(got, mau) {
			t.Errorf("RNMShareAcrossTheBoard=%q → %v, mau %v", board, got, mau)
		}
	}
}

// ⭐ TOTAL BERLIPAT ITU DISENGAJA — KEPUTUSAN PEMILIK PROSES 7 Oktober 2026.
//
// Pemakai melaporkan `Total Share RNM Limit` tampil dua kali lipat. Setelah
// ditelusuri ke ekspor dan ditanyakan dengan angkanya di tangan, jawabannya:
// **ikuti Pega, biarkan berlipat**.
//
// ---------------------------------------------------------------------
// ⛔ MENGAPA UJI INI ADA
// ---------------------------------------------------------------------
// Penggandaan ini TERLIHAT seperti cacat, dan siapa pun yang membuka layar
// berikutnya akan tergoda "memperbaikinya". Komentar saja tidak cukup
// menahan — komentar tidak merah. Uji ini merah.
//
// Sebabnya, dari ekspor: `SetSpreadName` langkah [5] ber-
// `pyStepsObjectName = TreatyIn.Limits`, jadi ia menapaki SELURUH pohon,
// bukan Detail yang sedang dikerjakan. Ia dipanggil di dalam gelung Detail
// (`TreatyInPropshare` [5.1.4], saat `.SpreadingTypeID == ""`), lalu langkah
// [6] menapaki pohon yang sama lagi dengan rumus yang persis sama.
// Totalnya dikosongkan HANYA SEKALI, di langkah [3].
//
// ⚠️ Angkanya dari laporan pemakai — bukan karangan.
func TestSharePropTotalBerlipatDisengaja(t *testing.T) {
	// TANPA Spreading Type → `SetSpreadName` → BERLIPAT.
	tanpa := `[{"TreatyType":"QUOTA SHARE","Detail":[{"TreatyGroup":"PROPERTY","CessionList":[{"Currency":"IDR","Value":"1757675000000"}]}]}]`
	h := shareProp(t, MasukanShareProp{
		Aksi: AksiSharePropShare, Limits: pohonUji(t, tanpa),
		RNMShareP: "1.28", OptionLimit: "1",
	}, &sumberPropUji{})

	if got := nilaiUji(larikSimpul(detailUji(h, 0, 0), "RNMShareList")); !reflect.DeepEqual(got, []string{"IDR 22498240000"}) {
		t.Errorf("daftar per Detail %v, mau [IDR 22498240000] — yang INI harus tetap benar", got)
	}
	if got := totalUji(h.TotalShareRnmProp); !reflect.DeepEqual(got, []string{"IDR 44996480000"}) {
		t.Errorf("total %v, mau [IDR 44996480000] (2x, DISENGAJA — lihat komentar di atas)", got)
	}
}

// ⭐ PASANGANNYA: cabang ber-Spreading Type TIDAK berlipat, dan selisih
// kedua cabang itulah yang membuat penggandaan di atas terbaca sebagai
// perilaku Pega, bukan sebagai kekeliruan kita.
//
// Cabang ini memanggil `FetchQSfromMaster`, bukan `SetSpreadName`. Gambar
// Pega 17 memperlihatkan totalnya tidak berlipat.
func TestSharePropBerSpreadingTypeTidakBerlipat(t *testing.T) {
	h := shareProp(t, MasukanShareProp{
		Aksi: AksiSharePropShare, Limits: pohonUji(t, limitsGambar17),
		RNMShareP: "10", OptionLimit: "2", Commencement: "20250101",
	}, sumberGambar17())
	daftar := nilaiUji(larikSimpul(detailUji(h, 0, 0), "RNMShareList"))
	total := totalUji(h.TotalShareRnmProp)
	if !reflect.DeepEqual(daftar, total) {
		t.Errorf("daftar %v lawan total %v — cabang ber-Spreading Type TIDAK boleh berlipat", daftar, total)
	}
}
