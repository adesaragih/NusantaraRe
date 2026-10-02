package deidentifikasi

import (
	"errors"
	"strings"
	"testing"
)

// Data sintetis - BUKAN data pelanggan. Nama dan alamat di bawah karangan uji.
const contoh = `{"CaseID":"NB-0001","InsuredName":"Budi Contoh","TotalSumInsured":1234567890.123456789012,
"CedingRetention":"25.5","Loc":[{"ASMAddress":"Jl. Uji No. 1","RiskZipCode":"12345","TSI":10,"pxCreateOpName":"Operator Uji"},
{"ASMAddress":"","TSI":0.00000001}],"Note":"untuk Budi Contoh saja","ShareCeding":7,"Flag":true,"Kosong":null}`

// TestBersihkanMembuangDaftarEksplisit - nilai kedua belas kunci BUANG kosong;
// kunci lain, termasuk TotalSumInsured / CedingRetention / ShareCeding yang
// tertangkap filter pola palsu, utuh. Urutan kunci dan teks angka dipertahankan.
func TestBersihkanMembuangDaftarEksplisit(t *testing.T) {
	bersih, lap, err := Bersihkan([]byte(contoh))
	if err != nil {
		t.Fatal(err)
	}
	s := string(bersih)
	for _, mau := range []string{
		`"InsuredName":""`, `"ASMAddress":""`, `"RiskZipCode":""`, `"pxCreateOpName":""`,
		`"TotalSumInsured":1234567890.123456789012`, `"CedingRetention":"25.5"`, `"ShareCeding":7`,
		`"TSI":0.00000001`, `"Flag":true`, `"Kosong":null`,
	} {
		if !strings.Contains(s, mau) {
			t.Errorf("hasil tidak memuat %s:\n%s", mau, s)
		}
	}
	if strings.Index(s, `"CaseID"`) > strings.Index(s, `"InsuredName"`) {
		t.Error("urutan kunci berubah")
	}
	if lap.Dibuang["InsuredName"] != 1 || lap.Dibuang["ASMAddress"] != 2 || lap.Dibuang["pxCreateOpName"] != 1 {
		t.Errorf("laporan dibuang %v", lap.Dibuang)
	}
}

// TestPeriksaMenemukanKebocoran - nilai yang dibuang tetapi muncul lagi di medan
// lain (batas kata) dilaporkan sebagai kebocoran; ini cara menemukan medan
// identitas di luar daftar BUANG.
func TestPeriksaMenemukanKebocoran(t *testing.T) {
	bersih, _, err := Bersihkan([]byte(contoh))
	if err != nil {
		t.Fatal(err)
	}
	bocor, err := Periksa([]byte(contoh), bersih)
	if err != nil {
		t.Fatal(err)
	}
	if len(bocor) != 1 || !strings.Contains(bocor[0], "Note") {
		t.Errorf("kebocoran %v, mau satu di medan Note", bocor)
	}
}

// TestPeriksaMenolakAngkaBerubah - setiap nilai di luar kunci BUANG wajib
// identik; pembanding sendiri diuji dengan hasil yang sengaja dirusak.
func TestPeriksaMenolakAngkaBerubah(t *testing.T) {
	bersih, _, err := Bersihkan([]byte(contoh))
	if err != nil {
		t.Fatal(err)
	}
	rusak := strings.Replace(string(bersih), "1234567890.123456789012", "1234567890.123456789", 1)
	if _, err := Periksa([]byte(contoh), []byte(rusak)); err == nil {
		t.Error("angka yang berubah tidak terdeteksi")
	}
	rusak = strings.Replace(string(bersih), `"InsuredName":""`, `"InsuredName":"X"`, 1)
	if _, err := Periksa([]byte(contoh), []byte(rusak)); err == nil {
		t.Error("kunci BUANG yang masih berisi tidak terdeteksi")
	}
}

// TestKandidatIdentitas - kunci bernilai teks non-angka di luar daftar BUANG
// dilaporkan NAMANYA saja untuk ditinjau manusia; tidak ada yang dibuang otomatis.
func TestKandidatIdentitas(t *testing.T) {
	k, err := KandidatIdentitas([]byte(contoh))
	if err != nil {
		t.Fatal(err)
	}
	gabung := strings.Join(k, ",")
	if !strings.Contains(gabung, "Note") || !strings.Contains(gabung, "CaseID") || strings.Contains(gabung, "InsuredName") {
		t.Errorf("kandidat %v", k)
	}
}

// TestDaftarBuangDisetujui - isi daftar BUANG persis yang disetujui work owner:
// dua belas kunci tiket 15, sembilan medan bocor (keputusan 01-10-2026, butir 32),
// dua medan bocor sisa (butir 36), dan nomor polis/kasus (butir 38). Harapan disalin
// dari keputusan, bukan dari kode.
func TestDaftarBuangDisetujui(t *testing.T) {
	mau := []string{
		"InsuredName", "InsuredID", "ASMAddress", "SelectedLocationAddress", "RoadName", "ASMZipCode",
		"RiskZipCode", "Email", "pxCreateOpName", "MarketingName", "CedingCoName", "CoinsName",
		"pxCreateOperator", "AccumulationDescription", "AccumulationCode", "PIC", "PICSuggest",
		"CommentSuggest", "PropertiItemNote", "SobName", "TopRiskLocation",
		"Comment", "OperatorID",
		"PolicyNo", "OldPolicyNo", "EndorsementNo", "InvoiceNumber", "NoOfferSlip", "PolicyMasterNumber",
		"PolicyMasterIDPega", "IDNewBisnis", "IDFollowingNB", "Following",
	}
	if strings.Join(DaftarBuang, ",") != strings.Join(mau, ",") {
		t.Fatalf("DaftarBuang %v\nmau %v", DaftarBuang, mau)
	}
}

// TestDaftarKosongkanDisetujui - teks bebas dan bagian alamat (butir 38): dikosongkan,
// tetapi BUKAN sumber deteksi kebocoran. `ObjectName` tidak termasuk (label kategori, A16).
func TestDaftarKosongkanDisetujui(t *testing.T) {
	mau := []string{
		"Remarks", "REMARK", "EdmNote", "EdmSourceNote", "GoodNote", "OccupationNote", "CoverageNote",
		"TradingNote", "ConveyanceNote", "PackingNote", "EndorsmentReason", "Description", "NoteFinalScore",
		"Detail", "NM_SHIP", "BranchName",
		"ASMCity", "ASMDistrict", "ASMRW",
	}
	if strings.Join(DaftarKosongkan, ",") != strings.Join(mau, ",") {
		t.Fatalf("DaftarKosongkan %v\nmau %v", DaftarKosongkan, mau)
	}
}

// TestKosongkanBukanSumberKebocoran - nilai teks bebas yang muncul lagi di medan
// lain (salinan angka TSI, nama okupasi) bukan kebocoran; nilai DaftarBuang tetap.
// Data sintetis.
func TestKosongkanBukanSumberKebocoran(t *testing.T) {
	const masuk = `{"Remarks":"1000","TSI":"1000","OccupationNote":"Toko Rekaan","OccupationName":"Toko Rekaan",` +
		`"PolicyNo":"POLIS-REKAAN-1","Following":"","Catatan":"lihat POLIS-REKAAN-1"}`
	bersih, lap, err := Bersihkan([]byte(masuk))
	if err != nil {
		t.Fatal(err)
	}
	if lap.Dibuang["Remarks"] != 1 || lap.Dibuang["OccupationNote"] != 1 {
		t.Fatalf("laporan %v", lap.Dibuang)
	}
	for _, s := range []string{`"Remarks":""`, `"OccupationNote":""`, `"TSI":"1000"`, `"OccupationName":"Toko Rekaan"`} {
		if !strings.Contains(string(bersih), s) {
			t.Fatalf("keluaran tanpa %s: %s", s, bersih)
		}
	}
	bocor, err := Periksa([]byte(masuk), bersih)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(bocor, ",") != "$.Catatan" {
		t.Fatalf("bocor %v, mau hanya $.Catatan", bocor)
	}
}

// TestBersihkanJalurEksplisit - `ID` di tingkat dokumen (identitas kasus) dikosongkan
// lewat JalurBuang yang dicocokkan PERSIS; `ID` di jalur lain (kode mata uang,
// baris daftar) utuh. Data sintetis.
func TestBersihkanJalurEksplisit(t *testing.T) {
	const masuk = `{"ID":"KASUS-REKAAN 1","OldData":{"ID":"KASUS-REKAAN 0","OldData":{"ID":"9"}},` +
		`"Currency":{"ID":"IDR"},"L":[{"ID":"7"}],"OldDataX":{"ID":"tetap"}}`
	const mau = `{"ID":"","OldData":{"ID":"","OldData":{"ID":""}},` +
		`"Currency":{"ID":"IDR"},"L":[{"ID":"7"}],"OldDataX":{"ID":"tetap"}}`
	got, lap, err := Bersihkan([]byte(masuk))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != mau {
		t.Fatalf("dapat %s\nmau   %s", got, mau)
	}
	if lap.Dibuang["$.ID"] != 1 || lap.Dibuang["$.OldData.OldData.ID"] != 1 {
		t.Fatalf("laporan %v", lap.Dibuang)
	}
	if _, err := Periksa([]byte(masuk), got); err != nil {
		t.Fatalf("Periksa: %v", err)
	}
	kand, err := KandidatIdentitas([]byte(masuk))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(kand, ",") != "ID" { // dari Currency.ID dan OldDataX.ID, bukan dari jalur BUANG
		t.Fatalf("kandidat %v", kand)
	}
	// `ID` yang hanya ada di jalur BUANG bukan kandidat.
	if kand, err := KandidatIdentitas([]byte(`{"ID":"KASUS-REKAAN 1","Kode":"A"}`)); err != nil || strings.Join(kand, ",") != "Kode" {
		t.Fatalf("kandidat %v (%v), mau hanya Kode", kand, err)
	}
}

// TestJalurBuangDisetujui - butir 38: identitas kasus di tingkat dokumen.
func TestJalurBuangDisetujui(t *testing.T) {
	mau := "$.ID,$.OldData.ID,$.OldData.OldData.ID"
	if strings.Join(JalurBuang, ",") != mau {
		t.Fatalf("JalurBuang %v, mau %s", JalurBuang, mau)
	}
}

// TestSudahBersih - berkas yang semua daun BUANG-nya kosong lolos; satu nilai
// tersisa (kunci, kosongkan, atau jalur) ditolak dengan jalurnya. Data sintetis.
func TestSudahBersih(t *testing.T) {
	if err := SudahBersih([]byte(`{"InsuredName":"","ID":"","Remarks":"","TSI":"10","Currency":{"ID":"IDR"}}`)); err != nil {
		t.Fatalf("bersih: %v", err)
	}
	for _, kotor := range []string{
		`{"InsuredName":"Rekaan","TSI":"10"}`,
		`{"Remarks":"catatan rekaan"}`,
		`{"ID":"KASUS-REKAAN 1"}`,
	} {
		if err := SudahBersih([]byte(kotor)); !errors.Is(err, ErrBelumBersih) {
			t.Errorf("%s: galat %v, mau ErrBelumBersih", kotor, err)
		}
	}
}
