package models

// Uji seam fungsi murni pemuat dokumen lama (tiket 22): pembaca tanggal lama,
// pemecah dokumen, dan penggolong medan. ⛔ Fixture fiktif berawalan UJI-,
// dibentuk dari STRUKTUR contoh DATA_JSON (PERTANYAAN-untuk-DBA P29 "lima
// sifat", rancangan 4quinque) - bukan salinan isinya; berkasnya di testdata/.

import (
	_ "embed"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBacaTanggalLama(t *testing.T) { // spec-penyimpanan AC 21, 22; K15
	for _, tt := range []struct {
		masuk string
		gol   Golongan
		harap string
	}{
		{"", GolTanggal, ""},
		{" 20171130 ", GolTanggal, "2017-11-30"},                            // spasi tepi bukan bagian tanggal
		{"20171130", GolTanggal, "2017-11-30"},                              // AC 21 - YYYYMMDD
		{"20171130", GolTanggalWaktu, "2017-11-30"},                         // tanggal saja, jam tidak dikarang
		{"20170930T170000.000 GMT", GolTanggalWaktu, "2017-10-01 00:00:00"}, // AC 22 - 17:00 GMT = 00:00 WIB
		{"20170930T170000.000 GMT", GolTanggal, "2017-10-01"},
		{"20260102T050000.000 GMT", GolTanggalWaktu, "2026-01-02 12:00:00"},
	} {
		got, err := BacaTanggalLama(tt.masuk, tt.gol)
		if err != nil || got != tt.harap {
			t.Errorf("BacaTanggalLama(%q, %s) = %q, %v; harap %q", tt.masuk, tt.gol, got, err, tt.harap)
		}
	}
}

func TestTanggalAmbiguTidakDitebak(t *testing.T) { // K15
	for _, s := range []string{"05/06/2017", "05/06/2017 10:30 AM", "1/2/2018"} {
		if _, err := BacaTanggalLama(s, GolTanggal); !errors.Is(err, ErrTanggalAmbigu) {
			t.Errorf("%q: harap ErrTanggalAmbigu, dapat %v", s, err)
		}
	}
	// Bentuk lain - termasuk garis miring yang TIDAK ambigu - tidak diterjemahkan:
	// cara menentukan susunan per baris belum diputuskan (P32 butir 1).
	for _, s := range []string{"13/06/2017", "2017-11-30", "20171331", "20170230", "30-09-2017", "UJI-bukan-tanggal"} {
		if _, err := BacaTanggalLama(s, GolTanggal); !errors.Is(err, ErrFormatTanggal) {
			t.Errorf("%q: harap ErrFormatTanggal, dapat %v", s, err)
		}
	}
}

// dokumenUjiProp - bentuk contoh 2 (proporsional): ListInstallment DATAR
// (baris pertama membawa properti nomor baris Pega), SpreadingRiskList,
// QuotationData{CedingCoList}, pxObjClass di setiap simpul. Fixture dokumen
// lama ada di testdata/ (JSON): nama properti nomor baris Pega tidak boleh
// muncul di kode Go (penjaga lintas modul claimlife, PERMINTAAN-TIM-INTI.md A).
//
//go:embed testdata/dokumen_uji_prop.json
var dokumenUjiProp string

func barisUji(dokumen string) BarisJSONPolis {
	return BarisJSONPolis{
		IDPega: "ASM-FW-GISFW-WORK-NB UJI-77", NoPolis: "UJI-QP.T1.10.2017.00001", ProdKe: "0",
		TglInput: "2017-10-02 08:00:00", Username: "UJI-AKUN", DataJSON: []byte(dokumen),
	}
}

func TestPecahDokumenProporsionalDatar(t *testing.T) { // AC 52, 55, 29, 69; ID-19
	h, err := PecahDokumenLama(barisUji(dokumenUjiProp))
	if err != nil {
		t.Fatal(err)
	}
	if len(h.Galat) > 0 {
		t.Fatalf("galat tak terduga: %+v", h.Galat)
	}
	if h.ID != "UJI-77" {
		t.Errorf("ID kasus %q, harap UJI-77 (pyID dari IDPEGA)", h.ID)
	}
	for jalur, harap := range map[string]string{
		"PolicyTreatyIn.PremiOgp":                       "592629512.880000276",          // AC 55: tidak dibulatkan ke presisi mata uang
		"PolicyTreatyIn.StartDate":                      "2017-10-01",                   // AC 21
		"PolicyTreatyIn.EndDate":                        "2017-10-01",                   // AC 69: kosong -> = StartDate
		"PolicyTreatyIn.StatementDate":                  "2017-10-01 00:00:00",          // AC 22
		"PolicyTreatyIn.CedingCoName":                   "UJI-CEDING A; UJI-CEDING B; ", // AC 29: apa adanya
		"PolicyTreatyIn.QuotationData.GroupPanel":       "006",
		"PolicyTreatyIn.QuotationData.BusinessOldId":    "01",
		"PolicyTreatyIn.QuotationData.ProportionalType": "Proportional",
	} {
		if got := h.Halaman.Ambil(jalur); got != harap {
			t.Errorf("%s = %q, harap %q", jalur, got, harap)
		}
	}
	ceding := h.Halaman.AmbilDaftar(HalamanPolis + ".QuotationData.CedingCoList")
	if len(ceding) != 2 || ceding[1]["CedingCo"] != "UJI-C2" || ceding[1]["CedingCoName"] != "UJI-CEDING B" {
		t.Errorf("ceding %+v", ceding)
	}
	angsuran := h.Halaman.AmbilDaftar(DaftarAngsuran)
	if len(angsuran) != 2 || angsuran[0]["DueDate"] != "2017-11-01" || angsuran[1]["Premium"] != "148157378.220000069" {
		t.Errorf("angsuran datar %+v", angsuran)
	}
	if rinci := h.Halaman.AmbilDaftar(JalurAnak(DaftarAngsuran, 1, "InstallmentList")); len(rinci) != 0 {
		t.Errorf("bentuk datar tidak punya sarang, dapat %+v", rinci)
	}
	// Medan tak dikenal TERSIMPAN beserta nilainya (AC 57) - termasuk yang kosong.
	harapTak := map[string]string{
		"PolicyTreatyIn.UJIMedanFiktif":                   "UJI-nilai-lewat",
		"PolicyTreatyIn.QuotationData.UJIFiktifQuotation": "",
	}
	if len(h.BelumDiputuskan) != len(harapTak) {
		t.Errorf("medan tak dikenal %+v, harap %v", h.BelumDiputuskan, harapTak)
	}
	for _, m := range h.BelumDiputuskan {
		if v, ada := harapTak[m.Jalur]; !ada || v != m.Nilai {
			t.Errorf("medan tak dikenal tak terduga %+v", m)
		}
	}
	// Yang diabaikan menurut keputusan tertulis dihitung, tidak hilang diam-diam.
	if h.Diabaikan[AlasanKeadaanLayar] != 1 || h.Diabaikan[AlasanTurunan] != 1 ||
		h.Diabaikan[AlasanInternalPega] != 4 || h.Diabaikan[AlasanNourut] != 1 {
		t.Errorf("diabaikan %v", h.Diabaikan)
	}
}

// dokumenUjiNonProp - bentuk contoh 3 (non-proporsional XOL): ListInstallment
// BERSARANG (InstallmentList), TreatyXOLList > ValueList, OldData berisi
// TreatyXOLList kosong, TreatyDifference.
//
//go:embed testdata/dokumen_uji_nonprop.json
var dokumenUjiNonProp string

func TestPecahDokumenNonProporsionalBersarang(t *testing.T) { // AC 53; ID-26, ID-29
	b := barisUji(dokumenUjiNonProp)
	b.NoPolis = "UJI-QR.T1.01.2018.00002"
	h, err := PecahDokumenLama(b)
	if err != nil || len(h.Galat) > 0 {
		t.Fatalf("%v %+v", err, h.Galat)
	}
	rinci := h.Halaman.AmbilDaftar(JalurAnak(DaftarAngsuran, 1, "InstallmentList"))
	if len(rinci) != 2 || rinci[0]["DueDate"] != "2018-01-31" || rinci[1]["Premium"] != "1500.25" {
		t.Errorf("rincian angsuran bersarang %+v", rinci)
	}
	xol := h.Halaman.AmbilDaftar(HalamanPolis + ".TreatyXOLList")
	if len(xol) != 2 || xol[0]["Deduction"] != "200.5" {
		t.Errorf("XOL %+v", xol)
	}
	layer := h.Halaman.AmbilDaftar(JalurAnak(HalamanPolis+".TreatyXOLList", 1, "ValueList"))
	if len(layer) != 2 || layer[1]["Layer"] != "2" || layer[1]["LayerType"] != "UJI-LT" || layer[0]["Deduction"] != "100.25" {
		t.Errorf("layer %+v", layer)
	}
	if len(h.Halaman.AmbilDaftar(JalurAnak(HalamanPolis+".TreatyXOLList", 2, "ValueList"))) != 0 {
		t.Error("ValueList kosong baris kedua harus tetap kosong")
	}
	// pyExpanded: keputusan F3 (InputPolicyTreatyInDetail_NonProp 12.1 menulis,
	// nol pembaca) -> dibuang berbukti, bukan lagi "belum diputuskan".
	if len(h.BelumDiputuskan) != 0 {
		t.Errorf("medan belum diputuskan %+v", h.BelumDiputuskan)
	}
	if h.Diabaikan[AlasanOldData] != 1 || h.Diabaikan[AlasanSelisih] != 1 || h.Diabaikan[AlasanF3KeadaanBaris] != 1 {
		t.Errorf("diabaikan %v", h.Diabaikan)
	}
}

// TestUjiPemecahMencakupDuaBentuk - AC 54: uji yang hanya mencakup satu bentuk
// ListInstallment tidak memadai. Fixture testdata/ uji ini wajib memuat keduanya.
func TestUjiPemecahMencakupDuaBentuk(t *testing.T) {
	bentuk := map[string]bool{}
	for _, d := range []string{dokumenUjiProp, dokumenUjiNonProp} {
		b := barisUji(d)
		h, err := PecahDokumenLama(b)
		if err != nil {
			t.Fatal(err)
		}
		bersarang := false
		for i := range h.Halaman.AmbilDaftar(DaftarAngsuran) {
			if len(h.Halaman.AmbilDaftar(JalurAnak(DaftarAngsuran, i+1, "InstallmentList"))) > 0 {
				bersarang = true
			}
		}
		bentuk[map[bool]string{true: "bersarang", false: "datar"}[bersarang]] = true
	}
	if !bentuk["datar"] || !bentuk["bersarang"] {
		t.Fatalf("cakupan bentuk ListInstallment tidak memadai: %v", bentuk)
	}
}

// kasusPenggolong - pola medan -> alasan yang diharapkan; JSON karena nama
// properti nomor baris Pega tidak boleh muncul di kode Go (penjaga lintas
// modul claimlife, PERMINTAAN-TIM-INTI.md A).
//
//go:embed testdata/penggolong_kasus.json
var kasusPenggolong []byte

// TestPenggolongMedanDiabaikan mengunci penggolong medan yang SENGAJA tidak
// disimpan: setiap cabang (simpul utuh lebih dulu, ruas di mana pun, nomor
// baris hanya di baris daftar, skalar tingkat polis saja) beserta teks
// alasannya - teks itu tercetak di ringkasan pemuat (AC 59).
func TestPenggolongMedanDiabaikan(t *testing.T) {
	var uji struct {
		Alasan map[string]string `json:"alasan"`
		Kasus  []struct {
			Pola  string `json:"pola"`
			Harap string `json:"harap"`
		} `json:"kasus"`
	}
	if err := json.Unmarshal(kasusPenggolong, &uji); err != nil {
		t.Fatal(err)
	}
	if len(uji.Kasus) < 30 {
		t.Fatalf("hanya %d kasus terbaca; fixture-nya yang rusak", len(uji.Kasus))
	}
	for kunci, v := range map[string]string{
		"internal_pega": AlasanInternalPega, "nourut": AlasanNourut, "keadaan_layar": AlasanKeadaanLayar,
		"turunan": AlasanTurunan, "pantulan_layer": AlasanPantulanLayer, "breakdown": AlasanBreakdown,
		"persetujuan_dh": AlasanPersetujuanDH, "old_data": AlasanOldData, "selisih": AlasanSelisih,
		"f3_tanpa_pembaca": AlasanF3TanpaPembaca, "f3_keadaan_baris": AlasanF3KeadaanBaris,
		"f3_salinan_generasi": AlasanF3SalinanGenerasi,
	} {
		if v != uji.Alasan[kunci] {
			t.Errorf("alasan %s = %q, harap %q", kunci, v, uji.Alasan[kunci])
		}
	}
	for _, k := range uji.Kasus {
		harap := ""
		if k.Harap != "" {
			harap = uji.Alasan[k.Harap]
			if harap == "" {
				t.Fatalf("kasus %s: kunci alasan %q tidak ada di fixture", k.Pola, k.Harap)
			}
		}
		if got := alasanDiabaikan(k.Pola); got != harap {
			t.Errorf("alasanDiabaikan(%s) = %q, harap %q", k.Pola, got, harap)
		}
	}
}

// TestBerkasPenggolongCacatDitolak - berkas data tidak diperiksa penyusun Go:
// rujukan alasan yang salah ketik dan medan JSON tak dikenal wajib ditolak,
// bukan diam-diam menjadikan medannya "tak dikenal".
func TestBerkasPenggolongCacatDitolak(t *testing.T) {
	for nama, isi := range map[string]string{
		"alasan simpul tak terdefinisi": `{"alasan": {"a": "UJI-a"}, "simpul": {"UJI": "b"}}`,
		"alasan ruas tak terdefinisi":   `{"alasan": {"a": "UJI-a"}, "ruas": {"UJI": {"alasan": "b"}}}`,
		"alasan skalar tak terdefinisi": `{"alasan": {"a": "UJI-a"}, "skalar_polis": {"UJI": ""}}`,
		"medan JSON tak dikenal":        `{"alasan": {"a": "UJI-a"}, "simpull": {}}`,
		"alasan pola tak terdefinisi":   `{"alasan": {"a": "UJI-a"}, "pola": {"UJI": {"alasan": "b", "bukti": "UJI-bukti"}}}`,
		// F3: setiap keputusan per medan wajib berbukti XML.
		"pola tanpa bukti": `{"alasan": {"a": "UJI-a"}, "pola": {"UJI": {"alasan": "a", "bukti": " "}}}`,
	} {
		if _, err := muatPenggolongAbaikan([]byte(isi)); err == nil {
			t.Errorf("%s: harap ditolak", nama)
		}
	}
	if _, err := muatPenggolongAbaikan([]byte(`{"alasan": {"a": "UJI-a"}, "simpul": {"UJI": "a"}}`)); err != nil {
		t.Errorf("berkas sah ditolak: %v", err)
	}
}

func TestPetaKatalogDokumenDigerakkanKatalog(t *testing.T) {
	peta, err := petaKatalogDokumen()
	if err != nil {
		t.Fatalf("setiap daftar bersarang katalog wajib punya induk: %v", err)
	}
	for pola, gol := range map[string]Golongan{
		"PolicyTreatyIn.StartDate":                                   GolTanggal,
		"PolicyTreatyIn.QuotationData.GroupPanel":                    GolKode,
		"PolicyTreatyIn.QuotationData.CedingCoList().CedingCo":       GolKode,
		"PolicyTreatyIn.ListInstallment().InstallmentList().DueDate": GolTanggal,
		"PolicyTreatyIn.TreatyXOLList().ValueList().Deduction":       GolUang, // ID-30: uang
		"PolicyTreatyIn.SpreadingRiskList().SharePercentage":         GolPersen,
		// F3 (a): dibaca syarat `InputPolicyTreatyInPre_Act` langkah 10
		// (`.PolicyTreatyIn.EDMType=="3"` -> lewati langkah) - wajib berkolom.
		"PolicyTreatyIn.EDMType": GolKode,
	} {
		if k, ada := peta[pola]; !ada || k.Golongan != gol {
			t.Errorf("%s: %+v (ada %v), harap golongan %s", pola, k, ada, gol)
		}
	}
	for _, pola := range []string{"PositionNote", "NBStatus", "PolicyTreatyIn.Layer", "PolicyTreatyIn.TotalPremium"} {
		if _, ada := peta[pola]; ada {
			t.Errorf("%s bukan kolom dokumen", pola)
		}
	}
}

func TestDokumenBergalatTidakDimuatDanSebabnyaDisebut(t *testing.T) { // AC 58, K15
	ganti := func(dari, ke string) string { return strings.Replace(dokumenUjiProp, dari, ke, 1) }
	for _, tt := range []struct {
		nama    string
		dokumen string
		nopol   string
		jalur   string
		harap   error
	}{
		{"tanggal ambigu", ganti(`"StartDate": "20171001"`, `"StartDate": "05/06/2017"`), "", "PolicyTreatyIn.StartDate", ErrTanggalAmbigu},
		{"tanggal baris", ganti(`"DueDate": "20171201"`, `"DueDate": "20171331"`), "", "PolicyTreatyIn.ListInstallment(2).DueDate", ErrFormatTanggal},
		{"uang berkoma", ganti(`"PremiOgp": "592629512.880000276"`, `"PremiOgp": "12,5"`), "", "PolicyTreatyIn.PremiOgp", ErrNilaiKolom},
		{"nomor polis beda", dokumenUjiProp, "UJI-LAIN", "PolicyTreatyIn.PolicyNo", ErrNoPolisBeda},
	} {
		b := barisUji(tt.dokumen)
		if tt.nopol != "" {
			b.NoPolis = tt.nopol
		}
		h, err := PecahDokumenLama(b)
		if err != nil {
			t.Fatalf("%s: %v", tt.nama, err)
		}
		if len(h.Galat) != 1 || h.Galat[0].Jalur != tt.jalur || !errors.Is(h.Galat[0].Err, tt.harap) {
			t.Errorf("%s: galat %+v, harap %s di %s", tt.nama, h.Galat, tt.harap, tt.jalur)
		}
	}
	b := barisUji(dokumenUjiProp)
	b.NoPolis = " "
	if h, _ := PecahDokumenLama(b); len(h.Galat) == 0 || !errors.Is(h.Galat[len(h.Galat)-1].Err, ErrNoPolisKosong) {
		t.Errorf("NOPOLIS kosong: %+v", h.Galat)
	}
}

func TestDokumenDiLuarLingkupAtauRusak(t *testing.T) {
	for _, tt := range []struct {
		nama         string
		dokumen      string
		idpega, prod string
		harap        error
	}{
		{"JSON rusak", `{"pxObjClass": `, "", "0", ErrDokumenRusak},
		{"akar daftar", `[]`, "", "0", ErrDokumenRusak},
		{"lini lain", `{"pxObjClass": "UJI-Data-PolicyFacIn"}`, "", "0", ErrBukanTreatyIn},
		{"generasi endorsemen", dokumenUjiProp, "", "1", ErrGenerasiEndorsemen},
		{"PRODKE kosong", dokumenUjiProp, "", "", ErrProdKe},
		{"IDPEGA tanpa kelas", dokumenUjiProp, "UJI-77", "0", ErrIDPega},
	} {
		b := barisUji(tt.dokumen)
		b.ProdKe = tt.prod
		if tt.idpega != "" {
			b.IDPega = tt.idpega
		}
		if _, err := PecahDokumenLama(b); !errors.Is(err, tt.harap) {
			t.Errorf("%s: %v, harap %v", tt.nama, err, tt.harap)
		}
	}
}

func TestAngkaJSONTidakLewatFloat(t *testing.T) { // ADR-0003
	d := strings.Replace(dokumenUjiProp, `"PremiOgp": "592629512.880000276"`, `"PremiOgp": 592629512.880000276`, 1)
	h, err := PecahDokumenLama(barisUji(d))
	if err != nil || len(h.Galat) > 0 {
		t.Fatal(err, h.Galat)
	}
	if got := h.Halaman.Ambil("PolicyTreatyIn.PremiOgp"); got != "592629512.880000276" {
		t.Errorf("angka JSON = %q, harap literal utuh", got)
	}
}

// TestPanduanBentukDokumenNolMedanBelumDiputuskan - AC 59 (RALAT F3, WO
// 04-10-2026: "nol medan yang BELUM DIPUTUSKAN"). Setiap jalur daun panduan
// bentuk dokumen `docs/dataguide-json-polis.json` (`JSON_DATAGUIDE` atas
// POOLDATA.JSON_POLIS, 378 jalur) dirakit menjadi satu dokumen fiktif lalu
// dipecah: tiap medannya wajib berkolom, disalin ke riwayat produksi, atau
// dibuang dengan keputusan tertulis. ⚠️ Panduan itu basi (diagram R45) -
// medan di luar panduan baru terlihat saat uji-kering atas data nyata (F7).
func TestPanduanBentukDokumenNolMedanBelumDiputuskan(t *testing.T) {
	isi, err := os.ReadFile(filepath.Join("..", "..", "docs", "dataguide-json-polis.json"))
	if err != nil {
		t.Fatal(err)
	}
	var panduan []struct {
		Jalur string `json:"o:path"`
		Jenis string `json:"type"`
	}
	if err := json.Unmarshal(isi, &panduan); err != nil {
		t.Fatal(err)
	}
	daftar := map[string]bool{}
	for _, p := range panduan {
		if p.Jenis == "array" {
			daftar[p.Jalur] = true
		}
	}
	akar := map[string]any{"pxObjClass": KelasDokumenTreatyIn}
	daun := 0
	for _, p := range panduan {
		if p.Jenis == "object" || p.Jenis == "array" || p.Jalur == "$.pxObjClass" {
			continue
		}
		daun++
		ruas := strings.Split(strings.TrimPrefix(p.Jalur, "$."), ".")
		simpul, jalur := akar, "$"
		for i, r := range ruas {
			jalur += "." + r
			if i == len(ruas)-1 {
				simpul[r] = "" // kosong sah untuk setiap golongan kolom
				break
			}
			if daftar[jalur] {
				xs, _ := simpul[r].([]any)
				if len(xs) == 0 {
					xs = []any{map[string]any{}}
					simpul[r] = xs
				}
				simpul = xs[0].(map[string]any)
				continue
			}
			anak, _ := simpul[r].(map[string]any)
			if anak == nil {
				anak = map[string]any{}
				simpul[r] = anak
			}
			simpul = anak
		}
	}
	if daun < 340 {
		t.Fatalf("hanya %d jalur daun terbaca; panduannya yang rusak", daun)
	}
	dok, err := json.Marshal(akar)
	if err != nil {
		t.Fatal(err)
	}
	h, err := PecahDokumenLama(barisUji(string(dok)))
	if err != nil || len(h.Galat) > 0 {
		t.Fatalf("%v %+v", err, h.Galat)
	}
	for _, m := range h.BelumDiputuskan {
		t.Errorf("medan belum diputuskan: %s", m.Pola)
	}
}
