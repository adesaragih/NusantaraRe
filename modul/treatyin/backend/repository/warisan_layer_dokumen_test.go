package repository

// Bukti `LayerDariDokumen` — penerus `warisan_in2_test.go`.
//
// ⛔ KETIGA uji berkas lama DICABUT, dan sebabnya tertulis satu per satu —
// bukan hilang diam-diam:
//
//	TestPindaiLayerSejajarDenganDaftarKolom
//	  menjaga urutan `kolomLayer` sama dengan urutan medan di `pindaiLayer`.
//	  DICABUT: nol daftar kolom dan nol pemindai tersisa — baris layer
//	  terurai dari JSON bernama, dan nama tidak punya urutan untuk melenceng.
//
//	TestKolomLayerSamaDenganKatalogOracle
//	  mengadu ke-41 nama kolom dengan `ALL_TAB_COLUMNS`.
//	  DICABUT: nol kueri ke `M_TREATY_IN2` tersisa. Penggantinya
//	  `TestNolKueriMTreatyIn2`, yang menjaga hal yang lebih keras — bukan
//	  "namanya cocok" melainkan "tabelnya tidak disentuh sama sekali".
//
//	TestKueriLayerBerurutAngka
//	  menjaga `ORDER BY TO_NUMBER(LAYER …)` ada di kuerinya.
//	  DIUBAH ARTINYA menjadi `TestUrutanDokumenDipertahankan` di bawah:
//	  urutan kini milik dokumen, dan yang perlu dijaga adalah bahwa kita
//	  TIDAK mengurutkannya ulang.

import (
	"encoding/json"
	"testing"
)

func uraiLimits(t *testing.T, s string) []jsonLimit {
	t.Helper()
	var j struct {
		Limits []jsonLimit `json:"Limits"`
	}
	if err := json.Unmarshal([]byte(s), &j); err != nil {
		t.Fatalf("dokumen uji tidak terurai: %v", err)
	}
	return j.Limits
}

// ⭐ SATU baris per (LAYER × TREATY GROUP) — bentuk yang keempat tab
// harapkan, dan bentuk yang tabel lama pipihkan.
func TestSatuBarisPerLayerKaliTreatyGroup(t *testing.T) {
	baris := LayerDariDokumen(uraiLimits(t, `{"Limits":[
		{"Layer":"1","Currency":"IDR","Limit":"750000000",
		 "Detail":[{"TreatyGroup":"PROPERTY"},{"TreatyGroup":"MARINE"}]},
		{"Layer":"2","Currency":"USD","Limit":"1095000",
		 "Detail":[{"TreatyGroup":"ENGINEERING"}]}
	]}`))
	if len(baris) != 3 {
		t.Fatalf("%d baris, mau 3 (2 treaty group + 1): %+v", len(baris), baris)
	}
	mau := []struct{ layer, grup, mataUang string }{
		{"1", "PROPERTY", "IDR"},
		{"1", "MARINE", "IDR"},
		{"2", "ENGINEERING", "USD"},
	}
	for i, m := range mau {
		b := baris[i]
		if b.Layer != m.layer || b.KelompokTreaty != m.grup || b.MataUang != m.mataUang {
			t.Errorf("baris %d = {%s %s %s}, mau {%s %s %s}",
				i, b.Layer, b.KelompokTreaty, b.MataUang, m.layer, m.grup, m.mataUang)
		}
	}
}

// ⛔ Layer TANPA `Detail[]` tetap memberi SATU baris.
//
// Menjatuhkannya akan menyembunyikan layer yang di Pega terlihat — medan
// layernya ada, hanya treaty group-nya tidak.
func TestLayerTanpaDetailTetapSatuBaris(t *testing.T) {
	baris := LayerDariDokumen(uraiLimits(t, `{"Limits":[
		{"Layer":"1","Currency":"IDR","Limit":"750000000","Detail":[]},
		{"Layer":"2","Currency":"IDR"}
	]}`))
	if len(baris) != 2 {
		t.Fatalf("%d baris, mau 2: %+v", len(baris), baris)
	}
	for i, b := range baris {
		if b.KelompokTreaty != "" {
			t.Errorf("baris %d punya kelompok treaty %q padahal nol Detail", i, b.KelompokTreaty)
		}
	}
}

// ⭐ Medan tingkat LAYER menurun ke setiap barisnya.
func TestMedanLayerMenurunKeTiapBaris(t *testing.T) {
	baris := LayerDariDokumen(uraiLimits(t, `{"Limits":[{
		"Layer":"3","LayerType":"layer","Cover":"risk","TreatyType":"QUOTA SHARE",
		"Currency":"IDR","Currency2":"USD","Limit":"5000000000","Deductible":"3500000000",
		"AdjRate":"6.5","ROLPct":"31.0600","MDPPct":"90","CurrencyRelation":"OR",
		"MDPList":[{"Currency":"IDR","Value":"413156049999"}],
		"PremiumEarnedList":[{"Currency":"IDR","Value":"170966.2500"}],
		"Detail":[{"TreatyGroup":"A"},{"TreatyGroup":"B"}]}]}`))
	if len(baris) != 2 {
		t.Fatalf("%d baris, mau 2", len(baris))
	}
	for i, b := range baris {
		for _, p := range []struct{ nama, dapat, mau string }{
			{"Layer", b.Layer, "3"},
			{"JenisLayer", b.JenisLayer, "layer"},
			{"DasarCover", b.DasarCover, "risk"},
			{"MataUang", b.MataUang, "IDR"},
			{"MataUangLimit", b.MataUangLimit, "USD"},
			{"Limit100", b.Limit100, "5000000000"},
			// ⭐ `.Deductible` — padanan yang pengaduan nilai buktikan 70,6%
			// dan gambar 30 dokumen desain beri namanya di layar.
			{"RetensiCedant", b.RetensiCedant, "3500000000"},
			{"AdjRate", b.AdjRate, "6.5"},
			{"ROL", b.ROL, "31.0600"},
			{"RasioMDP", b.RasioMDP, "90"},
			{"RelasiMataUang", b.RelasiMataUang, "OR"},
			{"MDP", b.MDP, "413156049999"},
			{"PremiEarned", b.PremiEarned, "170966.2500"},
		} {
			if p.dapat != p.mau {
				t.Errorf("baris %d %s = %q, mau %q", i, p.nama, p.dapat, p.mau)
			}
		}
	}
}

// ⭐ Medan tingkat `Detail[]` BERBEDA per baris — di situlah tabel lama
// kehilangan jejak, dan di situlah bukti bahwa pemipihan kita benar.
func TestMedanDetailBerbedaPerBaris(t *testing.T) {
	baris := LayerDariDokumen(uraiLimits(t, `{"Limits":[{"Layer":"1","Detail":[
		{"TreatyGroup":"PROPERTY","CessionPct":"50","Brokerage":"10","RNMShare":"25",
		 "RIOGR":"32,50","Earthquake":"384000000000","SpreadingType":"QS HR 60M TRT",
		 "CessionList":[{"Value":"4500000000.0000"}],
		 "EPIList":[{"Value":"3000000000"}],
		 "RNMShareList":[{"Value":"1125000000"}],
		 "RNMSpreadedList":[{"Value":"450000000"}],
		 "RNMSpreadedListRI":[{"Value":"675000000"}]},
		{"TreatyGroup":"MARINE","CessionPct":"40","RNMShare":"15"}]}]}`))
	if len(baris) != 2 {
		t.Fatalf("%d baris, mau 2", len(baris))
	}
	a, b := baris[0], baris[1]
	if a.PersenCession != "50" || b.PersenCession != "40" {
		t.Errorf("CessionPct %q lawan %q, mau 50 lawan 40", a.PersenCession, b.PersenCession)
	}
	if a.RNMShare != "25" || b.RNMShare != "15" {
		t.Errorf("RNMShare %q lawan %q", a.RNMShare, b.RNMShare)
	}
	for _, p := range []struct{ nama, dapat, mau string }{
		{"PersenBrokerage", a.PersenBrokerage, "10"},
		{"RIOGR", a.RIOGR, "32,50"},
		{"Gempa", a.Gempa, "384000000000"},
		{"JenisPenyebaran", a.JenisPenyebaran, "QS HR 60M TRT"},
		{"CessionKeRI", a.CessionKeRI, "4500000000.0000"},
		{"EPI100", a.EPI100, "3000000000"},
		{"LiabilityRNM", a.LiabilityRNM, "1125000000"},
		{"LiabilityQSOR", a.LiabilityQSOR, "450000000"},
		{"LiabilityQSRI", a.LiabilityQSRI, "675000000"},
	} {
		if p.dapat != p.mau {
			t.Errorf("%s = %q, mau %q", p.nama, p.dapat, p.mau)
		}
	}
	// Baris kedua nol — dan nol itu bukan nilai baris pertama yang bocor.
	if b.CessionKeRI != "" || b.EPI100 != "" {
		t.Errorf("nilai baris pertama bocor ke baris kedua: %+v", b)
	}
}

// ⛔ `QS (OR)` dan `QS (R/I)` dicocokkan menurut NAMA, bukan urutan.
//
// Sapuan menemukan keduanya dalam urutan yang berbeda antar dokumen; membaca
// elemen ke-0 sebagai OR akan menukar keduanya diam-diam, dan dua persen
// yang tertukar terbaca benar.
func TestSpreadingDicocokkanMenurutNama(t *testing.T) {
	// Sengaja TERBALIK: R/I lebih dulu.
	baris := LayerDariDokumen(uraiLimits(t, `{"Limits":[{"Layer":"1","Detail":[{
		"SpreadingList":[
			{"ReinsTypeName":"QS (R/I)","Pct":"60"},
			{"ReinsTypeName":"QS (OR)","Pct":"40"}]}]}]}`))
	if len(baris) != 1 {
		t.Fatalf("%d baris", len(baris))
	}
	if baris[0].QSOR != "40" {
		t.Errorf("QSOR = %q, mau 40 — dicocokkan menurut urutan, bukan nama", baris[0].QSOR)
	}
	if baris[0].QSRI != "60" {
		t.Errorf("QSRI = %q, mau 60", baris[0].QSRI)
	}
}

// ⭐ URUTAN DOKUMEN dipertahankan — penerus `TestKueriLayerBerurutAngka`.
//
// ⛔ Dahulu urutan datang dari `ORDER BY TO_NUMBER(LAYER)`. Kini ia milik
// dokumen, dan yang perlu dijaga BERBALIK: kita tidak boleh mengurutkannya
// ulang. Urutan di dokumen adalah urutan yang layar lama tampilkan.
func TestUrutanDokumenDipertahankan(t *testing.T) {
	baris := LayerDariDokumen(uraiLimits(t, `{"Limits":[
		{"Layer":"10"},{"Layer":"2"},{"Layer":"1"}]}`))
	var urut []string
	for _, b := range baris {
		urut = append(urut, b.Layer)
	}
	mau := []string{"10", "2", "1"}
	for i := range mau {
		if urut[i] != mau[i] {
			t.Fatalf("urutan %v, mau %v — dokumen tidak boleh diurutkan ulang", urut, mau)
		}
	}
}

// ⛔ Irisan KOSONG, bukan nil — penerus
// `TestKontrakTanpaBarisLayerMengembalikanKosong`.
func TestNolLimitsMemberiIrisanKosong(t *testing.T) {
	for nama, dok := range map[string]string{
		"nol kunci":    `{}`,
		"larik kosong": `{"Limits":[]}`,
		"null":         `{"Limits":null}`,
	} {
		baris := LayerDariDokumen(uraiLimits(t, dok))
		if baris == nil {
			t.Errorf("%s: nil, mau irisan kosong", nama)
		}
		if len(baris) != 0 {
			t.Errorf("%s: %d baris", nama, len(baris))
		}
	}
}

// ⛔ TIPE BERCAMPUR tidak menggugurkan dokumen.
//
// Sapuan menemukan `.Layer` sebagai angka JSON pada sebagian dokumen dan
// sebagai string pada sebagian lain. `*string` menolak keduanya sekaligus
// dan membuang SELURUH kontrak; `json.RawMessage` menerima apa adanya.
func TestTipeBercampurDiterima(t *testing.T) {
	baris := LayerDariDokumen(uraiLimits(t, `{"Limits":[
		{"Layer":1,"Limit":750000000,"MDPPct":90.5},
		{"Layer":"2","Limit":"750000000","MDPPct":"90.5"}]}`))
	if len(baris) != 2 {
		t.Fatalf("%d baris, mau 2", len(baris))
	}
	if baris[0].Layer != "1" || baris[1].Layer != "2" {
		t.Errorf("Layer %q dan %q", baris[0].Layer, baris[1].Layer)
	}
	if baris[0].Limit100 != "750000000" || baris[1].Limit100 != "750000000" {
		t.Errorf("Limit100 %q dan %q", baris[0].Limit100, baris[1].Limit100)
	}
}

// ⛔ ANGKA BESAR tidak kehilangan digit.
//
// `json.RawMessage` membawa teks angka apa adanya. Melewatkannya melalui
// `float64` mengembalikan `9007199254740992` untuk `9007199254740993`, dan
// digit yang hilang itu tidak pernah terlihat sebagai galat.
func TestAngkaBesarUtuh(t *testing.T) {
	baris := LayerDariDokumen(uraiLimits(t,
		`{"Limits":[{"Layer":"1","Limit":9007199254740993}]}`))
	if baris[0].Limit100 != "9007199254740993" {
		t.Errorf("Limit100 = %q, digit hilang", baris[0].Limit100)
	}
}

// ⛔⛔ MATA UANG KEDUA — cacat `nilaiPertama`, ditutup 6 Oktober 2026.
//
// ⚠️ Uji POSITIFNYA memakai bentuk layer bermata uang DUA, dan itu bukan
// pilihan gaya: layer bermata uang TUNGGAL lulus bahkan dengan cacatnya,
// sehingga uji yang memakainya tidak membuktikan apa pun.
func TestMataUangKeduaTerbaca(t *testing.T) {
	// Bentuk yang disalin dari kontrak NYATA `1000003` Limits[0] — terukur
	// 6 Oktober 2026, salah satu dari 106 layer bermata uang dua.
	baris := LayerDariDokumen(uraiLimits(t, `{"Limits":[{
		"Layer":"1","Currency":"IDR","Currency2":"USD",
		"Limit":"35629453681.71","Limit2":"2284.50",
		"Deductible":"150000000","Deductible2":"1155000",
		"MDPList":[
			{"Currency":"IDR","Value":"255030285.3504000"},
			{"Currency":"USD","Value":"11798.18156160"}],
		"PremiumEarnedList":[
			{"Currency":"IDR","Value":"318787856.688000"},
			{"Currency":"USD","Value":"14747.7269520"}]}]}`))
	if len(baris) != 1 {
		t.Fatalf("%d baris, mau 1", len(baris))
	}
	b := baris[0]
	for _, p := range []struct{ nama, dapat, mau string }{
		{"MataUang", b.MataUang, "IDR"},
		{"MataUangLimit", b.MataUangLimit, "USD"},
		{"Limit100", b.Limit100, "35629453681.71"},
		{"Limit100Kedua", b.Limit100Kedua, "2284.50"},
		{"RetensiCedant", b.RetensiCedant, "150000000"},
		{"RetensiCedantKedua", b.RetensiCedantKedua, "1155000"},
		{"MDP", b.MDP, "255030285.3504000"},
		// ⛔ INILAH yang dahulu hilang.
		{"MDPKedua", b.MDPKedua, "11798.18156160"},
		{"PremiEarned", b.PremiEarned, "318787856.688000"},
		{"PremiEarnedKedua", b.PremiEarnedKedua, "14747.7269520"},
	} {
		if p.dapat != p.mau {
			t.Errorf("%s = %q, mau %q", p.nama, p.dapat, p.mau)
		}
	}
}

// ⭐ Uji NEGATIF — layer bermata uang TUNGGAL tetap kosong di medan kedua.
//
// ⛔ Kosong, BUKAN salinan yang pertama. Menyalin akan membuat layar
// menampilkan satu nilai dua kali di bawah dua mata uang yang berbeda —
// bentuk yang terbaca benar dan salah.
func TestMataUangKeduaKosongBilaTunggal(t *testing.T) {
	baris := LayerDariDokumen(uraiLimits(t, `{"Limits":[{
		"Layer":"1","Currency":"IDR","Limit":"750000000",
		"MDPList":[{"Currency":"IDR","Value":"455112272998"}],
		"PremiumEarnedList":[{"Currency":"IDR","Value":"170966.25"}]}]}`))
	b := baris[0]
	for _, p := range []struct{ nama, dapat string }{
		{"Limit100Kedua", b.Limit100Kedua},
		{"RetensiCedantKedua", b.RetensiCedantKedua},
		{"MDPKedua", b.MDPKedua},
		{"PremiEarnedKedua", b.PremiEarnedKedua},
	} {
		if p.dapat != "" {
			t.Errorf("%s = %q, mau kosong — nilai pertama tidak boleh disalin", p.nama, p.dapat)
		}
	}
	// Yang pertama tetap terbaca.
	if b.MDP != "455112272998" || b.PremiEarned != "170966.25" {
		t.Errorf("nilai pertama rusak: MDP=%q PremiEarned=%q", b.MDP, b.PremiEarned)
	}
}

// ⛔ Larik BERPANJANG TIGA tidak membuang elemen ketiga diam-diam — ia
// memang tidak punya tempat, dan itu harus terlihat di sini bila kelak ada.
//
// ⚠️ Terukur 6 Oktober 2026: panjang maksimum di seluruh 1.854 dokumen
// adalah DUA. Uji ini merah pada hari yang ketiga muncul, dan hari itu
// bentuk dua-kolom perlu ditinjau ulang.
func TestLarikNilaiTidakPernahLebihDariDua(t *testing.T) {
	baris := LayerDariDokumen(uraiLimits(t, `{"Limits":[{"Layer":"1",
		"MDPList":[{"Currency":"IDR","Value":"1"},
		           {"Currency":"USD","Value":"2"},
		           {"Currency":"SGD","Value":"3"}]}]}`))
	b := baris[0]
	if b.MDP != "1" || b.MDPKedua != "2" {
		t.Errorf("dua pertama tidak terbaca: %q %q", b.MDP, b.MDPKedua)
	}
	t.Log("⚠️ elemen KETIGA (SGD 3) tidak punya tempat di bentuk dua-kolom; " +
		"terukur nol kemunculan di POOLDATA pada 6 Oktober 2026")
}
