package models

// Uji seam fungsi murni pemuat dokumen lama ENDORSEMEN (tiket EDM 10): pemecah dokumen, penggolong medan, urutan
// generasi. ⛔ Fixture fiktif berawalan UJI-, dibentuk dari STRUKTUR data guide json_polis (bukan salinan isi
// dokumen produksi); berkasnya di testdata/ - nama properti nomor baris Pega tidak boleh muncul di kode Go
// (penjaga lintas modul claimlife `polaIndeksPosisi`).

import (
	_ "embed"
	"errors"
	"strings"
	"testing"
)

//go:embed testdata/dokumenlama_edm_prop.json
var dokumenUjiEDMProp string

//go:embed testdata/dokumenlama_edm_nonprop.json
var dokumenUjiEDMNonProp string

func barisUjiEDM(dokumen string) BarisJSONPolis {
	return BarisJSONPolis{
		IDPega: "ASM-FW-GISFW-WORK-ENDORSEMENTTREATY EDMT-990002", NoPolis: "UJI-QP.T1.10.2017.00001",
		NoEndors: "UJI-QP.T1.10.2017.00001/E02", ProdKe: "2",
		TglInput: "2017-10-02 08:00:00", Username: "UJI-AKUN", DataJSON: []byte(dokumen),
	}
}

func TestPecahDokumenEDMProporsional(t *testing.T) { // AC 39 (ID-33), ID-8, AC 31/33
	h, err := PecahDokumenEDM(barisUjiEDM(dokumenUjiEDMProp))
	if err != nil {
		t.Fatal(err)
	}
	if len(h.Galat) > 0 {
		t.Fatalf("galat tak terduga: %+v", h.Galat)
	}
	if h.ID != "EDMT-990002" || h.ProdKe != 2 || h.EDMNo != "UJI-QP.T1.10.2017.00001/E02" || h.EDMType != "3" ||
		h.NoPolis != "UJI-QP.T1.10.2017.00001" || h.OldDataEDMNo != "UJI-QP.T1.10.2017.00001/E01" {
		t.Errorf("kunci generasi %+v", h)
	}
	if h.Datar != (KolomDatarLama{TglInput: "2017-10-02 08:00:00", Username: "UJI-AKUN"}) {
		t.Errorf("kolom datar %+v", h.Datar)
	}
	// ⛔ AC 39: selisih dokumen APA ADANYA - termasuk nilai varian rumus berlapis Pega (180 - 50 = 130, bukan 30)
	// dan digit galat.
	for jalur, harap := range map[string]string{
		"PolicyTreatyIn.TreatyDifference.NetPremium": "-130463146.760000276",
		"PolicyTreatyIn.TreatyDifference.PremiOgp":   "130.000000276",
		"PolicyTreatyIn.TreatyDifference.RiCommOgp":  "25",
		"PolicyTreatyIn.PremiOgp":                    "180.000000276",
		"PolicyTreatyIn.StartDate":                   "2017-10-01",
		"PolicyTreatyIn.StatementDate":               "2017-10-01 00:00:00",
		"PolicyTreatyIn.EndDate":                     "", // tidak diisi (bukan keputusan EDM)
		"PolicyTreatyIn.EDMNo":                       "UJI-QP.T1.10.2017.00001/E02",
		"PolicyTreatyIn.ProdKe":                      "2",
	} {
		if got := h.Halaman.Ambil(jalur); got != harap {
			t.Errorf("%s = %q, harap %q", jalur, got, harap)
		}
	}
	sp := h.Halaman.AmbilDaftar(DaftarSelisihSpreading)
	if len(sp) != 2 || sp[0]["PremiumSpreaded"] != "78.000000166" || sp[1]["ClaimSpreaded"] != "-0.000000001" ||
		sp[1]["TreatyType"] != "UJI-SP" {
		t.Errorf("selisih spreading %+v", sp)
	}
	an := h.Halaman.AmbilDaftar(DaftarSelisihAngsuran)
	if len(an) != 1 || an[0]["DueDate"] != "2017-11-01" || an[0]["Premium"] != "-119961119.76" {
		t.Errorf("selisih angsuran %+v", an)
	}
	// ID-8: OldData tidak ikut halaman generasi; Proportional: hasil antara XOL dan rincian angsuran dibuang.
	for j := range h.Halaman.Nilai {
		if strings.Contains(j, ".OldData.") {
			t.Errorf("OldData tertinggal di halaman: %s", j)
		}
	}
	if len(h.Halaman.AmbilDaftar(TabelXOL.Daftar)) != 0 || len(h.Halaman.AmbilDaftar(DaftarSelisihXOL)) != 0 ||
		len(h.Halaman.AmbilDaftar(JalurAnak(DaftarAngsuran, 1, TabelAngsuranRinci.Daftar))) != 0 {
		t.Error("hasil antara Proportional tertinggal")
	}
	if err := PeriksaBentukSimpan(h.Halaman); err != nil {
		t.Errorf("halaman tidak lolos penjaga bentuk jalur biasa: %v", err)
	}
	harap := map[string]int{
		pmKunciPropHasilAntara: 12, "old_data": 7, "internal_pega": 7, "nourut": 1, "f3_keadaan_baris": 1,
		"keadaan_layar": 1, "turunan": 1, "selisih_turunan": 2, "selisih_tanpa_penulis": 1, "f3_tanpa_pembaca": 1,
		"": 1,
	}
	dapat := map[string]int{}
	for _, m := range h.Arsip {
		dapat[m.Kunci]++
	}
	for k, n := range harap {
		if dapat[k] != n {
			t.Errorf("arsip %q: %d medan, harap %d (seluruh arsip %v)", k, dapat[k], n, dapat)
		}
	}
	if len(dapat) != len(harap) {
		t.Errorf("arsip %v, harap %v", dapat, harap)
	}
	if len(h.BelumDiputuskan) != 1 || h.BelumDiputuskan[0].Jalur != "PolicyTreatyIn.UJIMedanFiktif" ||
		h.BelumDiputuskan[0].Nilai != "UJI-nilai-lewat" {
		t.Errorf("belum diputuskan %+v", h.BelumDiputuskan)
	}
	if h.Diabaikan[pmTeksAlasan("old_data")] != 7 {
		t.Errorf("diabaikan %v", h.Diabaikan)
	}
}

func TestPecahDokumenEDMNonProporsional(t *testing.T) { // ID-7, AC 29; ID-33
	b := barisUjiEDM(dokumenUjiEDMNonProp)
	b.IDPega, b.NoPolis, b.NoEndors, b.ProdKe = "asm-fw-gisfw-work-endorsementtreaty EDMT-990011", "UJI-QR.T1.01.2018.00002", "", "1"
	h, err := PecahDokumenEDM(b)
	if err != nil || len(h.Galat) > 0 {
		t.Fatalf("%v %+v", err, h.Galat)
	}
	if h.EDMNo != "UJI-QR.T1.01.2018.00002/E01" {
		t.Errorf("NOENDORS kosong -> EDMNo dokumen, dapat %q", h.EDMNo)
	}
	lapisan := DatarSelisihLapisan(h.Halaman)
	if len(lapisan) != 2 || lapisan[0]["GrossPremi"] != "-6.990000001" || lapisan[1]["DueToValue"] != "-3.000000002" ||
		lapisan[1]["DueTo"] != "DUE TO YOU" {
		t.Errorf("lapisan selisih %+v", lapisan)
	}
	if len(h.Halaman.AmbilDaftar(JalurAnak(TabelXOL.Daftar, 1, TabelLayerXOL.Daftar))) != 2 {
		t.Error("TreatyXOLList NonProp hilang")
	}
	if r := h.Halaman.AmbilDaftar(JalurAnak(DaftarAngsuran, 1, TabelAngsuranRinci.Daftar)); len(r) != 1 || r[0]["DueDate"] != "2018-02-01" {
		t.Errorf("rincian angsuran NonProp %+v", r)
	}
	dapat := map[string]int{}
	for _, m := range h.Arsip {
		dapat[m.Kunci]++
	}
	// Induk XOL selisih: Currency, IDCurrency, DueTo, NetPremi, GrossPremi; pyExpanded induk XOL; nomor baris;
	// PPN rincian; OldData.EDMNo; TreatyDifference.DueTo TERISI = belum diputuskan (hanya_bila_kosong).
	for k, n := range map[string]int{"induk_selisih_xol": 5, "f3_keadaan_baris": 1, "nourut": 1, "f3_tanpa_pembaca": 1,
		"old_data": 1, "": 1} {
		if dapat[k] != n {
			t.Errorf("arsip %q: %d, harap %d (%v)", k, dapat[k], n, dapat)
		}
	}
	if len(h.BelumDiputuskan) != 1 || h.BelumDiputuskan[0].Jalur != "PolicyTreatyIn.TreatyDifference.DueTo" {
		t.Errorf("TreatyDifference.DueTo terisi wajib belum diputuskan: %+v", h.BelumDiputuskan)
	}
}

func TestPecahDokumenEDMLewatDanGalat(t *testing.T) {
	ganti := func(lama, baru string) string { return strings.Replace(dokumenUjiEDMProp, lama, baru, 1) }
	for nama, tt := range map[string]struct {
		ubah func(*BarisJSONPolis)
		err  error
	}{
		"kelas lain dilewati":    {func(b *BarisJSONPolis) { b.DataJSON = []byte(`{"pxObjClass":"UJI-Lain"}`) }, ErrBukanTreatyIn},
		"PRODKE 0 milik NB":      {func(b *BarisJSONPolis) { b.ProdKe = "0" }, ErrBukanGenerasiEndorsemen},
		"PRODKE rusak":           {func(b *BarisJSONPolis) { b.ProdKe = "x" }, ErrProdKe},
		"IDPEGA kasus NB":        {func(b *BarisJSONPolis) { b.IDPega = "ASM-FW-GISFW-WORK-NB NB-77" }, ErrIDPega},
		"IDPEGA kelas lain EDMT": {func(b *BarisJSONPolis) { b.IDPega = "ASM-FW-GISFW-WORK-LAIN EDMT-77" }, ErrIDPega},
		"dokumen rusak":          {func(b *BarisJSONPolis) { b.DataJSON = []byte(`{`) }, ErrDokumenRusak},
	} {
		b := barisUjiEDM(dokumenUjiEDMProp)
		tt.ubah(&b)
		if _, err := PecahDokumenEDM(b); !errors.Is(err, tt.err) {
			t.Errorf("%s: harap %v, dapat %v", nama, tt.err, err)
		}
	}
	for nama, tt := range map[string]struct {
		dok  string
		ubah func(*BarisJSONPolis)
		err  error
	}{
		"ProdKe dokumen beda": {ganti(`"ProdKe": "2"`, `"ProdKe": "1"`), nil, ErrProdKeBeda},
		"EDMNo beda":          {ganti(`"EDMNo": "UJI-QP.T1.10.2017.00001/E02"`, `"EDMNo": "UJI-X/E02"`), nil, ErrEDMNoBeda},
		"EDMNo kosong": {ganti(`"EDMNo": "UJI-QP.T1.10.2017.00001/E02",`, ``),
			func(b *BarisJSONPolis) { b.NoEndors = "" }, ErrEDMNoKosong},
		"PolicyNo beda":      {ganti(`"PolicyNo": "UJI-QP.T1.10.2017.00001"`, `"PolicyNo": "UJI-LAIN"`), nil, ErrNoPolisBeda},
		"NOPOLIS kosong":     {dokumenUjiEDMProp, func(b *BarisJSONPolis) { b.NoPolis = "" }, ErrNoPolisKosong},
		"tanggal ambigu":     {ganti(`"StartDate": "20171001"`, `"StartDate": "05/06/2017"`), nil, ErrTanggalAmbigu},
		"uang selisih rusak": {ganti(`"NetPremium": "-130463146.760000276"`, `"NetPremium": "UJI-bukan-angka"`), nil, ErrNilaiKolom},
		"tanggal selisih rusak": {ganti(`"InstallmentNo": "1", "DueDate": "20171101", "InstallmentPercentage"`,
			`"InstallmentNo": "1", "DueDate": "2017-11-01", "InstallmentPercentage"`), nil, ErrFormatTanggal},
	} {
		b := barisUjiEDM(tt.dok)
		if tt.ubah != nil {
			tt.ubah(&b)
		}
		h, err := PecahDokumenEDM(b)
		if err != nil {
			t.Errorf("%s: galat struktural %v", nama, err)
			continue
		}
		ada := false
		for _, g := range h.Galat {
			ada = ada || errors.Is(g.Err, tt.err)
		}
		if !ada {
			t.Errorf("%s: harap %v di %+v", nama, tt.err, h.Galat)
		}
	}
}

func TestIDKasusDariIDPegaEDM(t *testing.T) {
	for _, s := range []string{"ASM-FW-GISFW-WORK-ENDORSEMENTTREATY EDMT-1", "ASM-FW-GISFW-Work-EndorsementTreaty EDMT-990045"} {
		if _, err := IDKasusDariIDPegaEDM(s); err != nil {
			t.Errorf("%q: %v", s, err)
		}
	}
	for _, s := range []string{"", "EDMT-1", "ASM-FW-GISFW-WORK-ENDORSEMENTTREATY NB-1", "ASM-FW-GISFW-WORK-ENDORSEMENTTREATY EDMT-",
		"ASM-FW-GISFW-WORK-ENDORSEMENTTREATY EDMT-1A", "ASM-FW-GISFW-WORK-NB EDMT-1"} {
		if _, err := IDKasusDariIDPegaEDM(s); !errors.Is(err, ErrIDPega) {
			t.Errorf("%q: harap ErrIDPega, dapat %v", s, err)
		}
	}
}

func TestUrutKunciGenerasiMenurutBilangan(t *testing.T) { // tiket 10: generasi dimuat menurut nomornya
	k := []KunciJSONPolis{
		{"r5", "UJI-B", "1"}, {"r1", "UJI-A", "10"}, {"r2", "UJI-A", "9"}, {"r3", "UJI-A", "x"}, {"r4", "UJI-A", "2"},
		{"r0", "UJI-A", "2"},
	}
	UrutKunciGenerasi(k)
	var dapat []string
	for _, x := range k {
		dapat = append(dapat, x.Kunci)
	}
	if strings.Join(dapat, ",") != "r0,r4,r2,r1,r3,r5" {
		t.Errorf("urutan %v", dapat)
	}
}

func TestPenggolongEDMMenolakBerkasCacat(t *testing.T) {
	for nama, isi := range map[string]string{
		"alasan tak terdefinisi": `{"alasan":{"a":"A"},"simpul":{"X":"b"}}`,
		"pola tanpa bukti":       `{"alasan":{"a":"A"},"pola":{"X":{"alasan":"a"}}}`,
		"ruas F3 tanpa bukti":    `{"alasan":{"f3_x":"A"},"ruas":{"X":{"alasan":"f3_x"}}}`,
		"medan tak dikenal":      `{"alasan":{"a":"A"},"lain":1}`,
	} {
		if _, err := pmMuatPenggolong([]byte(isi)); err == nil {
			t.Errorf("%s: harap ditolak", nama)
		}
	}
	if _, err := pmMuatPenggolong([]byte(`{"alasan":{"f3_x":"A"},"ruas":{"X":{"alasan":"f3_x"}},"bukti_ruas":{"X":"UJI bukti"}}`)); err != nil {
		t.Errorf("berkas sah ditolak: %v", err)
	}
}

func TestKunciDibuangHanyaBilaKosong(t *testing.T) {
	if k := pmKunciDibuang(HalamanPolis+".TreatyDifference.Currency", ""); k != "selisih_tanpa_penulis" {
		t.Errorf("kosong: %q", k)
	}
	if k := pmKunciDibuang(HalamanPolis+".TreatyDifference.Currency", "UJI-IDR"); k != "" {
		t.Errorf("terisi harus belum diputuskan, dapat %q", k)
	}
	if k := pmKunciDibuang(HalamanPolis+".pyExpanded", "true"); k != "" {
		t.Errorf("pyExpanded tingkat polis bukan baris daftar, dapat %q", k)
	}
}

func TestHalamanPembandingGenerasiNB(t *testing.T) {
	b := BarisJSONPolis{IDPega: "ASM-FW-GISFW-WORK-NB NB-990000", NoPolis: "UJI-QP.T1.10.2017.00001", ProdKe: "0",
		DataJSON: []byte(`{"pxObjClass":"ASM-FW-GISFW-Data-PolicyTreatyIn","QuotationData":{"ProportionalType":"Proportional"},
		"SpreadingRiskList":[{"TreatyType":"UJI-QS"}],"TreatyXOLList":[{"Currency":"UJI"}],"OldData":{"EDMNo":"UJI-X"}}`)}
	h, err := HalamanPembanding(b)
	if err != nil {
		t.Fatal(err)
	}
	if h.Ambil(pmJalurEDMNo) != "" || h.Ambil(pmJalurOldEDMNo) != "" || len(h.AmbilDaftar(TabelXOL.Daftar)) != 0 ||
		len(h.AmbilDaftar(DaftarSpreading)) != 1 {
		t.Errorf("pembanding %+v", h)
	}
}
