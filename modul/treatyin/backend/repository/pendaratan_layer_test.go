package repository

// Bukti bahwa keempat tab yang kini dibaca dari TABEL PENDARATAN merangkai
// bentuk yang sama dengan yang dahulu dirangkai dari dokumen.
//
// ⛔ Tanpa Oracle, dan itu disengaja: yang diuji di sini PERANGKAIANNYA —
// tempat kesalahan bentuk hidup — bukan kuerinya. Uji yang menuntut basis
// data untuk membuktikan bentuk tidak akan dijalankan orang.

import (
	"testing"
)

func baris(id, induk int64, nilai map[string]string) barisPendaratan {
	return barisPendaratan{ID: id, Induk: induk, Nilai: nilai}
}

// Satu layer tanpa detail tetap menghasilkan SATU baris — aturan yang sama
// seperti `LayerDariDokumen` dahulu, dan ia nyata: layer non-proporsional
// tanpa `Detail[]` hilang dari layar bila aturan ini dilanggar.
func TestRangkaiLayerTanpaDetailTetapSatuBaris(t *testing.T) {
	lim := []barisPendaratan{baris(1, 0, map[string]string{
		"LAYER": "L-1", "LAYERTYPE": "XOL", "CURRENCY": "IDR", "LIMIT": "1500000000",
	})}
	got := RangkaiLayer(lim, nil, nil, nil)
	if len(got) != 1 {
		t.Fatalf("mau 1 baris, dapat %d", len(got))
	}
	if got[0].Layer != "L-1" || got[0].Limit100 != "1500000000" {
		t.Fatalf("medan layer tidak terbawa: %+v", got[0])
	}
	// ⛔ Irisan KOSONG, bukan nil — layar membaca `.length` di atasnya.
	if got[0].KelasBisnis == nil {
		t.Fatal("KelasBisnis nil; layar akan berhenti membaca .length")
	}
}

// Satu layer dengan DUA detail menghasilkan DUA baris, masing-masing
// membawa medan layer yang sama.
func TestRangkaiLayerSatuBarisPerDetail(t *testing.T) {
	lim := []barisPendaratan{baris(7, 0, map[string]string{"LAYER": "L-1", "CURRENCY": "IDR"})}
	det := []barisPendaratan{
		baris(11, 7, map[string]string{"TREATYGROUP": "PROPERTY", "CESSIONPCT": "25"}),
		baris(12, 7, map[string]string{"TREATYGROUP": "MARINE", "CESSIONPCT": "30"}),
	}
	got := RangkaiLayer(lim, det, nil, nil)
	if len(got) != 2 {
		t.Fatalf("mau 2 baris, dapat %d", len(got))
	}
	for _, b := range got {
		if b.Layer != "L-1" || b.MataUang != "IDR" {
			t.Fatalf("medan layer tidak diwariskan ke baris detail: %+v", b)
		}
	}
	if got[0].KelompokTreaty != "PROPERTY" || got[1].PersenCession != "30" {
		t.Fatalf("medan detail tertukar: %+v", got)
	}
}

// Kelas bisnis menempel pada detail MILIKNYA, bukan pada detail saudaranya.
// Inilah yang `IDINDUK` beli, dan satu-satunya cara membuktikannya adalah
// memberi dua detail dengan kelas bisnis berbeda.
func TestRangkaiLayerKelasBisnisTidakBocorAntarDetail(t *testing.T) {
	lim := []barisPendaratan{baris(7, 0, map[string]string{"LAYER": "L-1"})}
	det := []barisPendaratan{
		baris(11, 7, map[string]string{"TREATYGROUP": "PROPERTY"}),
		baris(12, 7, map[string]string{"TREATYGROUP": "MARINE"}),
	}
	cob := []barisPendaratan{
		baris(21, 11, map[string]string{"CLASSOFBUSINESS": "FIRE"}),
		baris(22, 11, map[string]string{"CLASSOFBUSINESS": "ENGINEERING"}),
		baris(23, 12, map[string]string{"CLASSOFBUSINESS": "CARGO"}),
	}
	got := RangkaiLayer(lim, det, cob, nil)
	if len(got[0].KelasBisnis) != 2 || got[0].KelasBisnis[0] != "FIRE" {
		t.Fatalf("kelas bisnis detail pertama salah: %v", got[0].KelasBisnis)
	}
	if len(got[1].KelasBisnis) != 1 || got[1].KelasBisnis[0] != "CARGO" {
		t.Fatalf("kelas bisnis detail kedua salah: %v", got[1].KelasBisnis)
	}
}

// `TreatyType` ada di DUA tingkat, dan yang di detail menang bila terisi —
// aturan yang dibawa apa adanya dari `LayerDariDokumen`.
func TestRangkaiLayerTreatyTypeDetailMenang(t *testing.T) {
	lim := []barisPendaratan{baris(7, 0, map[string]string{"TREATYTYPE": "SURPLUS"})}
	det := []barisPendaratan{
		baris(11, 7, map[string]string{"TREATYTYPE": "QUOTA SHARE"}),
		baris(12, 7, map[string]string{"TREATYTYPE": ""}),
	}
	got := RangkaiLayer(lim, det, nil, nil)
	if got[0].JenisTreaty != "QUOTA SHARE" {
		t.Fatalf("detail yang terisi harus menang, dapat %q", got[0].JenisTreaty)
	}
	if got[1].JenisTreaty != "SURPLUS" {
		t.Fatalf("detail KOSONG tidak boleh menimpa tingkat layer, dapat %q", got[1].JenisTreaty)
	}
}

// Pohon tiga tingkat: Limits -> Detail -> COBList/AchievementLists.
func TestRangkaiPohonLimitsTigaTingkat(t *testing.T) {
	lim := []barisPendaratan{baris(7, 0, map[string]string{
		"TREATYTYPE": "QUOTA SHARE", "TREATYTYPEID": "10042",
	})}
	det := []barisPendaratan{baris(11, 7, map[string]string{
		"TREATYGROUP": "PROPERTY", "QSPCT": "25", "RETENTIONPCT": "75",
	})}
	cob := []barisPendaratan{baris(21, 11, map[string]string{"CLASSOFBUSINESS": "FIRE"})}
	ach := []barisPendaratan{baris(31, 11, map[string]string{"QUARTER": "Q 1", "PREMIUM": "1000"})}

	got := RangkaiPohonLimits(lim, det, cob, ach)
	if len(got) != 1 {
		t.Fatalf("mau 1 simpul puncak, dapat %d", len(got))
	}
	// ⛔ Kuncinya ejaan PEGA, bukan nama kolom Oracle: layar mengikat
	// medannya dengan ejaan itu, dan menukar sumber data tidak boleh
	// menukar nama medan di layar.
	if got[0]["TreatyType"] != "QUOTA SHARE" {
		t.Fatalf("kunci puncak harus ejaan Pega: %v", got[0])
	}
	rinci, ok := got[0]["Detail"].([]map[string]any)
	if !ok || len(rinci) != 1 {
		t.Fatalf("Detail bukan larik satu simpul: %#v", got[0]["Detail"])
	}
	if rinci[0]["QSPct"] != "25" {
		t.Fatalf("medan detail tidak terbawa: %v", rinci[0])
	}
	for _, nama := range []string{"COBList", "AchievementLists"} {
		larik, ok := rinci[0][nama].([]map[string]any)
		if !ok || len(larik) != 1 {
			t.Fatalf("%s bukan larik satu simpul: %#v", nama, rinci[0][nama])
		}
	}
}

// Detail TANPA grid tetap membawa larik KOSONG, bukan nihil — `TabLimitsProp`
// membaca `.length` di atasnya, dan `null` menghentikan halaman.
func TestRangkaiPohonLimitsGridKosongBukanNihil(t *testing.T) {
	lim := []barisPendaratan{baris(7, 0, map[string]string{"TREATYTYPE": "SURPLUS"})}
	det := []barisPendaratan{baris(11, 7, map[string]string{"TREATYGROUP": "PROPERTY"})}
	got := RangkaiPohonLimits(lim, det, nil, nil)
	rinci := got[0]["Detail"].([]map[string]any)
	for _, nama := range []string{"COBList", "AchievementLists"} {
		larik, ok := rinci[0][nama].([]map[string]any)
		if !ok {
			t.Fatalf("%s bukan larik: %#v", nama, rinci[0][nama])
		}
		if larik == nil {
			t.Fatalf("%s nihil; layar akan berhenti membaca .length", nama)
		}
	}
}

// Nol layer mengembalikan larik KOSONG, bukan nihil.
func TestRangkaiNolBarisTetapIrisanKosong(t *testing.T) {
	if got := RangkaiLayer(nil, nil, nil, nil); got == nil {
		t.Fatal("RangkaiLayer mengembalikan nil")
	}
	if got := RangkaiPohonLimits(nil, nil, nil, nil); got == nil {
		t.Fatal("RangkaiPohonLimits mengembalikan nil")
	}
}

// ⭐ Besaran per layer DIPISAH menurut `JENIS` dan URUT — elemen ke-0 milik
// mata uang pertama, ke-1 milik yang kedua.
//
// ⛔ Yang dijaga BUKAN "kodenya jalan" melainkan PASANGAN mata uangnya.
// `T_TREATY_LIMIT_MEASURE` menampung TIGA larik dalam satu tabel, dan
// pembaca yang lupa menyaring `JENIS` akan mengambil angka `PremiumEarnedList`
// lalu menaruhnya di kolom MDP — angka yang salah di tempat yang benar,
// yang tidak terlihat sampai seseorang membandingkan dengan sistem lama.
func TestBesaranLayerDipisahMenurutJenisDanUrutan(t *testing.T) {
	lim := []barisPendaratan{{ID: 7, Nilai: map[string]string{"LAYER": "1"}}}
	ukur := []barisPendaratan{
		{ID: 1, Induk: 7, Nilai: map[string]string{"JENIS": "MDPList", "VALUE": "mdp-IDR"}},
		{ID: 2, Induk: 7, Nilai: map[string]string{"JENIS": "PremiumEarnedList", "VALUE": "premi-IDR"}},
		{ID: 3, Induk: 7, Nilai: map[string]string{"JENIS": "MDPList", "VALUE": "mdp-USD"}},
		{ID: 4, Induk: 7, Nilai: map[string]string{"JENIS": "EgnpiTotalList", "VALUE": "egnpi"}},
		{ID: 5, Induk: 7, Nilai: map[string]string{"JENIS": "PremiumEarnedList", "VALUE": "premi-USD"}},
	}
	got := RangkaiLayer(lim, nil, nil, ukur)
	if len(got) != 1 {
		t.Fatalf("%d baris, mau 1", len(got))
	}
	b := got[0]
	for _, p := range [][2]string{
		{"MDP", b.MDP}, {"MDPKedua", b.MDPKedua},
		{"PremiEarned", b.PremiEarned}, {"PremiEarnedKedua", b.PremiEarnedKedua},
	} {
		mau := map[string]string{
			"MDP": "mdp-IDR", "MDPKedua": "mdp-USD",
			"PremiEarned": "premi-IDR", "PremiEarnedKedua": "premi-USD",
		}[p[0]]
		if p[1] != mau {
			t.Errorf("%s = %q, mau %q", p[0], p[1], mau)
		}
	}
	// ⚠️ `EgnpiTotalList` ikut mendarat di tabel yang sama dan TIDAK boleh
	// bocor ke salah satu dari keempat medan di atas.
	if b.MDP == "egnpi" || b.PremiEarned == "egnpi" {
		t.Error("nilai EgnpiTotalList bocor ke medan MDP/PremiEarned")
	}
}

// ⛔ Larik yang hanya berisi SATU elemen memberi medan kedua KOSONG, bukan
// menyalin yang pertama. Menyalinnya membuat layar memperlihatkan angka
// yang sama dua kali seolah kontraknya memang bermata uang dua.
func TestBesaranSatuElemenMemberiKeduaKosong(t *testing.T) {
	lim := []barisPendaratan{{ID: 3, Nilai: map[string]string{"LAYER": "1"}}}
	ukur := []barisPendaratan{
		{ID: 1, Induk: 3, Nilai: map[string]string{"JENIS": "MDPList", "VALUE": "satu"}},
	}
	got := RangkaiLayer(lim, nil, nil, ukur)
	if got[0].MDP != "satu" {
		t.Errorf("MDP = %q, mau \"satu\"", got[0].MDP)
	}
	if got[0].MDPKedua != "" {
		t.Errorf("MDPKedua = %q, mau kosong", got[0].MDPKedua)
	}
}

// ⭐ 454 — `BreakDownSprdList` terpasang ke baris spreading INDUKNYA
// (`IDINDUK` = ID baris `T_TREATY_LIMIT_SPREADING`), bukan ke Detail.
func TestPohonLimitsMemasangPecahanSpreading(t *testing.T) {
	pohon := RangkaiPohonLimitsPeta(map[string][]barisPendaratan{
		"T_TREATY_LIMITS":       {baris(1, 0, map[string]string{"TREATYTYPE": "QUOTA SHARE"})},
		"T_TREATY_LIMIT_DETAIL": {baris(10, 1, map[string]string{"TREATYGROUP": "HOSPITAL"})},
		"T_TREATY_LIMIT_SPREADING": {
			baris(100, 10, map[string]string{"REINSTYPEID": "10263", "REINSTYPENAME": "2025 QS 181M TRT", "PCT": "25"}),
			baris(101, 10, map[string]string{"REINSTYPEID": "10007", "REINSTYPENAME": "ORS", "PCT": "0"}),
		},
		"T_TREATY_LIMIT_SPRD_BREAKDOWN": {
			baris(1000, 100, map[string]string{"REINSNAME": "QS (OR)", "CURRENCY": "IDR", "SHAREPCT": "40", "AMOUNT": "22500000"}),
			baris(1001, 100, map[string]string{"REINSNAME": "QS (R/I)", "CURRENCY": "IDR", "SHAREPCT": "60", "AMOUNT": "33750000"}),
		},
	})
	sebar := pohon[0]["Detail"].([]map[string]any)[0]["SpreadingList"].([]map[string]any)
	if len(sebar) != 2 || sebar[0]["ReinsTypeName"] != "2025 QS 181M TRT" {
		t.Fatalf("SpreadingList %v", sebar)
	}
	pecah := sebar[0]["BreakDownSprdList"].([]map[string]any)
	if len(pecah) != 2 || pecah[0]["ReinsName"] != "QS (OR)" || pecah[1]["Amount"] != "33750000" || pecah[0]["SharePct"] != "40" {
		t.Errorf("BreakDownSprdList baris 1 %v", pecah)
	}
	// Baris tanpa pecahan — larik KOSONG, bukan nihil (layar membaca `.length`).
	if kosong, ok := sebar[1]["BreakDownSprdList"].([]map[string]any); !ok || len(kosong) != 0 {
		t.Errorf("BreakDownSprdList baris 2 %v", sebar[1]["BreakDownSprdList"])
	}
}
