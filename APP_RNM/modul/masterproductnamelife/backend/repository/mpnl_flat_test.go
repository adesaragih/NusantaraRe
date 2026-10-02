package repository

// Pemetaan flat dan rekonsiliasi pindah - TANPA Oracle (uji `db`: `mpnl_flat_db_test.go`).

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"nusantarare/modul/masterproductnamelife/backend/models"
)

// polaKolomDDLFlat - satu baris definisi kolom DDL 140–147.
var polaKolomDDLFlat = regexp.MustCompile(`^\s+([A-Z][A-Z0-9_]*)\s+(VARCHAR2\((\d+)\)|NUMBER\(38,8\)|NUMBER\(5\)|DATE|TIMESTAMP)`)

type kolomDDL struct {
	nama, tipe string
	lebar      int
}

// kolomDDLFlat - kolom setiap tabel flat menurut berkas migrasi maju modul ini.
func kolomDDLFlat(t *testing.T) map[string][]kolomDDL {
	t.Helper()
	berkas, err := filepath.Glob(filepath.Join(akarModul, "backend", "migrations", "14*.sql"))
	if err != nil {
		t.Fatal(err)
	}
	hasil := map[string][]kolomDDL{}
	for _, b := range berkas {
		if strings.HasSuffix(b, "_down.sql") {
			continue
		}
		isi, err := os.ReadFile(b)
		if err != nil {
			t.Fatal(err)
		}
		m := polaTabelDibuat.FindStringSubmatch(string(isi))
		if m == nil {
			continue
		}
		for _, baris := range strings.Split(string(isi), "\n") {
			if k := polaKolomDDLFlat.FindStringSubmatch(baris); k != nil {
				lebar, _ := strconv.Atoi(k[3])
				hasil[m[1]] = append(hasil[m[1]], kolomDDL{nama: k[1], tipe: k[2], lebar: lebar})
			}
		}
	}
	if len(hasil) != len(DaftarTabelFlat) {
		t.Fatalf("DDL terbaca %d tabel, mau %d; pembacanya yang rusak", len(hasil), len(DaftarTabelFlat))
	}
	return hasil
}

func tipeDDLMenurut(j JenisKolomFlat, lebar int) string {
	switch j {
	case FlatDesimal:
		return "NUMBER(38,8)"
	case FlatBulat:
		return "NUMBER(5)"
	case FlatTanggal:
		return "DATE"
	case FlatStempel:
		return "TIMESTAMP"
	default:
		return "VARCHAR2(" + strconv.Itoa(lebar) + ")"
	}
}

func cocokkanKolom[T any](t *testing.T, tabel string, kunci []kolomDDL, kolom []KolomFlat[T], ddl []kolomDDL) {
	t.Helper()
	mau := append([]kolomDDL{}, kunci...)
	for _, k := range kolom {
		mau = append(mau, kolomDDL{nama: k.Nama, tipe: tipeDDLMenurut(k.Jenis, k.Lebar)})
	}
	if len(mau) != len(ddl) {
		t.Errorf("%s: spesifikasi Go %d kolom, DDL %d", tabel, len(mau), len(ddl))
		return
	}
	for i := range mau {
		if mau[i].nama != ddl[i].nama || mau[i].tipe != ddl[i].tipe {
			t.Errorf("%s kolom ke-%d: Go %s %s, DDL %s %s", tabel, i+1, mau[i].nama, mau[i].tipe, ddl[i].nama, ddl[i].tipe)
		}
	}
}

// TestKolomFlatCocokDenganDDL - kontrak dua sisi: spesifikasi kolom Go (nama, URUTAN, tipe, lebar) = DDL 140–147.
func TestKolomFlatCocokDenganDDL(t *testing.T) {
	ddl := kolomDDLFlat(t)
	induk := ddl[TabelFlatInduk]
	// IS_ORS ditangani tersendiri (bendera): letaknya di DDL sesudah CAUSE.
	var tanpaBendera []kolomDDL
	for _, k := range induk {
		if k.nama != KolomIsORS {
			tanpaBendera = append(tanpaBendera, k)
		}
	}
	cocokkanKolom(t, TabelFlatInduk, []kolomDDL{{nama: KolomIDFlat, tipe: "VARCHAR2(6)"}}, KolomFlatInduk, tanpaBendera)
	kunci := []kolomDDL{{nama: KolomProductID, tipe: "VARCHAR2(6)"}, {nama: KolomUrut, tipe: "NUMBER(5)"}}
	cocokkanKolom(t, AnakLien.Tabel, kunci, AnakLien.Kolom, ddl[AnakLien.Tabel])
	cocokkanKolom(t, AnakDokumen.Tabel, kunci, AnakDokumen.Kolom, ddl[AnakDokumen.Tabel])
	cocokkanKolom(t, AnakPlan.Tabel, kunci, AnakPlan.Kolom, ddl[AnakPlan.Tabel])
	cocokkanKolom(t, AnakFinUW.Tabel, kunci, AnakFinUW.Kolom, ddl[AnakFinUW.Tabel])
	cocokkanKolom(t, AnakUWLimit.Tabel, kunci, AnakUWLimit.Kolom, ddl[AnakUWLimit.Tabel])
	cocokkanKolom(t, AnakOutward.Tabel, kunci, AnakOutward.Kolom, ddl[AnakOutward.Tabel])
	cocokkanKolom(t, AnakKomentar.Tabel, kunci, AnakKomentar.Kolom, ddl[AnakKomentar.Tabel])
	ada := false
	for _, k := range induk {
		ada = ada || (k.nama == KolomIsORS && k.tipe == "NUMBER(5)")
	}
	if !ada {
		t.Error("IS_ORS NUMBER(5) tidak ada di DDL induk")
	}
}

func TestNilaiKanonik(t *testing.T) {
	for _, k := range []struct {
		jenis         JenisKolomFlat
		lebar         int
		masuk, keluar string
		masalah       string
	}{
		{FlatDesimal, 0, "150000000.50", "150000000.5", ""},
		{FlatDesimal, 0, "0,5", "0.5", ""},
		{FlatDesimal, 0, ".5", "0.5", ""},
		{FlatDesimal, 0, "-0", "0", ""},
		{FlatDesimal, 0, "1000", "1000", ""},
		{FlatDesimal, 0, " 12 ", "12", ""},
		{FlatDesimal, 0, "0.12345678", "0.12345678", ""},
		{FlatDesimal, 0, "0.123456789", "", MasalahSkala},
		{FlatDesimal, 0, "1234567890123456789012345678901", "", MasalahDigitBulat},
		{FlatDesimal, 0, "1.000,5", "", MasalahBukanAngka},
		{FlatDesimal, 0, "1e5", "", MasalahBukanAngka},
		{FlatDesimal, 0, "UJI", "", MasalahBukanAngka},
		{FlatBulat, 0, "17", "17", ""},
		{FlatBulat, 0, "17.0", "17", ""},
		{FlatBulat, 0, "100", "100", ""},
		{FlatBulat, 0, "17.5", "", MasalahBukanBulat},
		{FlatBulat, 0, "100000", "", MasalahDigitKecil},
		{FlatBulat, 0, "99999", "99999", ""},
		{FlatTanggal, 0, "2026-03-01", "2026-03-01", ""},
		{FlatTanggal, 0, "01/03/2026", "", MasalahBukanTanggal},
		{FlatTanggal, 0, "UJI", "", MasalahBukanTanggal},
		{FlatStempel, 0, "20241202T065450.847 GMT", "20241202T065450.847 GMT", ""},
		{FlatStempel, 0, "2024-12-02", "", MasalahBukanStempel},
		{FlatTeks, 3, "abc", "abc", ""},
		{FlatTeks, 3, " a ", " a ", ""},
		{FlatTeks, 3, "abcd", "", MasalahTerlaluPanjang},
		{FlatTeks, 3, "é€", "", MasalahTerlaluPanjang}, // 5 byte: lebar dalam BYTE
		{FlatDesimal, 0, "", "", ""},
	} {
		keluar, m := NilaiKanonik(k.jenis, k.lebar, k.masuk)
		if keluar != k.keluar || m != k.masalah {
			t.Errorf("NilaiKanonik(%d, %q) = %q %q, mau %q %q", k.jenis, k.masuk, keluar, m, k.keluar, k.masalah)
		}
	}
}

// produkUjiFlat - satu produk dengan setiap jenis kolom terisi (fixture UJI-, nol data orang).
func produkUjiFlat() models.Produk {
	p := models.Produk{ID: "100044"}
	p.Umum.ProductName, p.Umum.Ceding, p.Umum.SOBName, p.Umum.RIComm = "UJI PRODUK", "UJI CEDING", "UJI SOB", "12.50"
	p.Umum.IsORS, p.Umum.CreateOp, p.Umum.UpdateOp, p.Umum.Comment = true, "UJI-A", "UJI-B", "komentar"
	p.Umum.TypeBasicRider = "" // medan mati kosong
	p.Inward.ID, p.Inward.ProductID = "100099", "100044"
	p.Inward.PolicyHolder, p.Inward.PolicyHolderName = "UJI-ORG-1", "UJI PEMEGANG"
	p.Inward.Begin, p.Inward.MinAge, p.Inward.CedingLimit = "2026-03-01", "17", "150000000.50"
	p.LienClause = []models.BarisLien{{Usia: "60", Manfaat: "50", Asli: `{"pxObjClass":"UJI"}`}}
	p.UnderwritingLimit = []models.BarisUWLimit{{MinInsured: "1", MaxInsured: "1,5", MinAge: "18", MaxAge: "65"}}
	p.CommentList = []models.BarisKomentar{{Date: "20241202T065450.847 GMT", OperatorName: "UJI-A", Suggest: "ok"}}
	return p
}

func TestNormalkanFlatKanonikDanMasalah(t *testing.T) {
	p := produkUjiFlat()
	n, masalah := NormalkanFlat(p)
	if len(masalah) != 0 {
		t.Fatalf("masalah: %+v", masalah)
	}
	if n.Umum.RIComm != "12.5" || n.Inward.CedingLimit != "150000000.5" || n.UnderwritingLimit[0].MaxInsured != "1.5" {
		t.Errorf("angka kanonik: %q %q %q", n.Umum.RIComm, n.Inward.CedingLimit, n.UnderwritingLimit[0].MaxInsured)
	}
	if n.Inward.ID != "100044" || n.Inward.ProductID != "100044" || n.Umum.PolicyHolder != "UJI-ORG-1" ||
		n.Umum.PolicyHolderName != "UJI PEMEGANG" {
		t.Errorf("identitas inward dan pemegang polis: %+v", n.Inward)
	}
	if n.Umum.Comment != "" || n.LienClause[0].Asli != "" || !n.Umum.IsORS {
		t.Errorf("Comment/Asli dibuang, IsORS tetap: %q %q %v", n.Umum.Comment, n.LienClause[0].Asli, n.Umum.IsORS)
	}
	if p.Umum.RIComm != "12.50" || p.LienClause[0].Asli == "" || p.UnderwritingLimit[0].MaxInsured != "1,5" {
		t.Error("NormalkanFlat mengubah produk pemanggil")
	}
	// Idempoten: bentuk kanonik dinormalkan lagi tidak berubah.
	if n2, m2 := NormalkanFlat(n); len(BedaProduk(n, n2)) != 0 || len(m2) != 0 {
		t.Errorf("tidak idempoten: %v %v", BedaProduk(n, n2), m2)
	}
	p.Inward.Mature, p.Inward.MinAge, p.UnderwritingLimit[0].MaxInsured = "31/02/2026x", "17.5", "UJI"
	_, masalah = NormalkanFlat(p)
	jenis := map[string]string{}
	for _, m := range masalah {
		jenis[m.Tabel+"."+m.Kolom] = m.Jenis
	}
	if jenis["M_PRODUCTNAME_LIFE.MATURE"] != MasalahBukanTanggal || jenis["M_PRODUCTNAME_LIFE.MINAGE"] != MasalahBukanBulat ||
		jenis["M_PRODUCTNAME_LIFE_UWLIMIT.MAXINSURED"] != MasalahBukanAngka || len(masalah) != 3 {
		t.Errorf("masalah: %+v", masalah)
	}
}

func TestMedanTanpaKolom(t *testing.T) {
	p := produkUjiFlat()
	if m := MedanTanpaKolom(p); len(m) != 0 {
		t.Fatalf("produk uji: %v", m)
	}
	p.Umum.Grup, p.Inward.Months = "2", "12"
	if m := MedanTanpaKolom(p); len(m) != 2 {
		t.Errorf("GRUP dan inward.MONTHS terisi: %v", m)
	}
}

var polaPenampungSQL = regexp.MustCompile(`:([0-9]+)`)

// penampungBerurutan - :1..:n tepat berurutan di teks (go-ora mengikat menurut urutan kemunculan).
func penampungBerurutan(t *testing.T, nama, q string, n int) {
	t.Helper()
	var dapat []string
	for _, m := range polaPenampungSQL.FindAllStringSubmatch(q, -1) {
		dapat = append(dapat, m[1])
	}
	if len(dapat) != n {
		t.Errorf("%s: %d penampung, mau %d", nama, len(dapat), n)
		return
	}
	for i, d := range dapat {
		if d != strconv.Itoa(i+1) {
			t.Errorf("%s: penampung ke-%d = :%s", nama, i+1, d)
			return
		}
	}
}

func TestSQLFlatPenampungBerurutan(t *testing.T) {
	n, _ := NormalkanFlat(produkUjiFlat())
	args := argInduk(n)
	penampungBerurutan(t, "sisip induk", sqlSisipInduk("S.T"), len(args)+1)
	penampungBerurutan(t, "perbarui induk", sqlPerbaruiInduk("S.T"), len(args)+1)
	if !strings.HasSuffix(sqlPerbaruiInduk("S.T"), "WHERE ID = :"+strconv.Itoa(len(args)+1)) {
		t.Errorf("WHERE ID di penampung terakhir: %s", sqlPerbaruiInduk("S.T"))
	}
	penampungBerurutan(t, "baca induk", sqlBacaInduk("S.T", true), 1)
	nAnak := 2
	for _, k := range AnakUWLimit.Kolom {
		nAnak += len(argTulis(k.Jenis, ""))
	}
	penampungBerurutan(t, "sisip anak", sqlSisipAnak(AnakUWLimit, "S.T"), nAnak)
	penampungBerurutan(t, "sisip komentar", sqlSisipAnak(AnakKomentar, "S.T"), 2+3)
	for _, q := range []string{sqlSisipInduk("S.T"), sqlPerbaruiInduk("S.T"), sqlBacaInduk("S.T", false),
		sqlSisipAnak(AnakPlan, "S.T"), sqlBacaAnak(AnakPlan, "S.T"), sqlDaftarFlat("S.T")} {
		if strings.Contains(strings.ToUpper(q), "COMMIT") {
			t.Errorf("COMMIT di SQL: %s", q)
		}
	}
}

func TestArgTulisAngkaTanpaNLS(t *testing.T) {
	for _, k := range []struct {
		v     string
		koef  any
		skala any
	}{
		{"150000000.5", "1500000005", int64(1)}, {"0.5", "5", int64(1)}, {"-12.25", "-1225", int64(2)},
		{"100", "100", int64(0)}, {"0", "0", int64(0)}, {"", nil, int64(0)},
	} {
		a := argTulis(FlatDesimal, k.v)
		if len(a) != 2 || a[0] != k.koef || a[1] != k.skala {
			t.Errorf("argTulis(%q) = %v, mau [%v %v]", k.v, a, k.koef, k.skala)
		}
	}
	if a := argTulis(FlatStempel, "20241202T065450.847 GMT"); len(a) != 1 || a[0] != "20241202065450.847" {
		t.Errorf("stempel: %v", a)
	}
}

// --- rekonsiliasi --------------------------------------------------------------------

func barisJ(id, isi string) barisJSON { return barisJSON{id: id, isi: isi} }

func TestRekonsiliasiLolosDanKanonik(t *testing.T) {
	umum := []barisJSON{
		barisJ("100045", `{"ID":"100045","PRODUCTNAME":"UJI B","CEDING":"UJI C","pxObjClass":"UJI-KELAS","IsView":"false",
			"DocumentClaim":[{"Document":"UJI DOK 1"},{"Document":"UJI DOK 2","pxObjClass":"UJI"}],
			"OutwardList":[{"REINSTYPEID":"","OVR_COMM":""},{"REINSTYPEID":"10200","TRANSACTIONYEAR":"2026"}],
			"CommentList":[{"Date":"20241202T065450.847 GMT","OperatorName":"UJI-A","Suggest":"ok","IsApproved":""}]}`),
		barisJ("100044", `{"ID":"100044","PRODUCTNAME":"UJI A","RICOMM":"10","IsORS":"true","POLICYHODERNAME":"UJI P"}`),
	}
	inward := []barisJSON{
		barisJ("100044", `{"ID":"100044","PRODUCTID":"100044","BEGIN":"01/03/2026","MINAGE":"17","POLICYHODERNAME":"UJI P"}`),
		barisJ("100099", `{"ID":"100099","PRODUCTID":"100045","CEDINGLIMIT":"1500"}`),
	}
	lap, produk := rekonsiliasi(umum, inward)
	if !lap.Lolos() || len(produk) != 2 {
		t.Fatalf("mau lolos:\n%s", lap.Teks())
	}
	if lap.K4OutwardKosong != 1 || lap.Baris[TabelFlatOutward] != 1 || lap.Baris[TabelFlatDokumen] != 2 ||
		lap.Baris[TabelFlatKomentar] != 1 || lap.Baris[TabelFlatInduk] != 2 {
		t.Errorf("cacah baris / K4:\n%s", lap.Teks())
	}
	if lap.InwardBerIDLain != 1 || lap.InwardYatim != 0 || lap.ProdukTanpaInward != 0 {
		t.Errorf("inward: %+v", lap)
	}
	if lap.DibuangD2["kunci halaman umum pxObjClass"] != 1 || lap.DibuangD2["IsView (keadaan layar)"] != 1 ||
		lap.DibuangD2["baris DocumentClaim berkunci internal Pega"] != 1 || lap.DibuangD2["baris CommentList berkunci internal Pega"] != 0 {
		t.Errorf("dibuang D2: %v", lap.DibuangD2)
	}
	if produk[0].ID != "100044" || produk[0].Inward.Begin != "2026-03-01" || !produk[0].Umum.IsORS ||
		produk[1].Inward.ID != "100045" || produk[1].Inward.CedingLimit != "1500" {
		t.Errorf("produk: %+v / %+v", produk[0].Inward, produk[1].Inward)
	}
	// ⛔ Laporan agregat: nol nilai data.
	for _, nilai := range []string{"UJI A", "UJI DOK", "UJI P", "1500", "10200"} {
		if strings.Contains(lap.Teks(), nilai) {
			t.Errorf("laporan memuat nilai data %q", nilai)
		}
	}
}

func TestRekonsiliasiMenolak(t *testing.T) {
	for _, k := range []struct {
		nama         string
		umum, inward string
		gagal        bool
		normalisasi  string
		k3           int
	}{
		{"K3 MATURE bukan tanggal", `{"ID":"100044"}`, `{"ID":"100044","PRODUCTID":"100044","MATURE":"UJI"}`, false, "", 1},
		{"K3 MaxInsured UW bukan angka", `{"ID":"100044","UnderwritingLimitList":[{"MaxInsured":"UJI"}]}`,
			`{"ID":"100044","PRODUCTID":"100044"}`, false, "", 1},
		{"BEGIN bukan tanggal = gagal (bukan K3)", `{"ID":"100044"}`, `{"ID":"100044","PRODUCTID":"100044","BEGIN":"UJI"}`, true, "", 0},
		{"MINAGE pecahan = gagal", `{"ID":"100044"}`, `{"ID":"100044","PRODUCTID":"100044","MINAGE":"17.5"}`, true, "", 0},
		{"teks melebihi lebar = gagal", `{"ID":"100044","CURRENCY":"x"}`,
			`{"ID":"100044","PRODUCTID":"100044","CURRENCY":"UJI-MATA-UANG-PANJANG-SEKALI"}`, true, "", 0},
		{"medan mati terisi = gagal", `{"ID":"100044","GRUP":"2"}`, `{"ID":"100044","PRODUCTID":"100044"}`, true, "", 0},
		{"kunci inward warisan terisi = gagal", `{"ID":"100044"}`, `{"ID":"100044","PRODUCTID":"100044","MONTHS":"12"}`, true, "", 0},
		{"pemegang polis kedua sisi beda = gagal", `{"ID":"100044","POLICYHODER":"UJI-ORG-1"}`,
			`{"ID":"100044","PRODUCTID":"100044","POLICYHODER":"UJI-ORG-2"}`, true, "", 0},
		{"nol ekor = normalisasi", `{"ID":"100044","RICOMM":"12.50"}`, `{"ID":"100044","PRODUCTID":"100044"}`, false,
			"M_PRODUCTNAME_LIFE.RICOMM (nol ekor desimal)", 0},
		{"koma desimal = normalisasi", `{"ID":"100044"}`, `{"ID":"100044","PRODUCTID":"100044","CEDINGLIMIT":"1,5"}`, false,
			"M_PRODUCTNAME_LIFE.CEDINGLIMIT (koma desimal)", 0},
		{"normalisasi baris anak", `{"ID":"100044","UnderwritingLimitList":[{"MinAge":"018"}]}`,
			`{"ID":"100044","PRODUCTID":"100044"}`, false, "M_PRODUCTNAME_LIFE_UWLIMIT.MINAGE (nol depan)", 0},
	} {
		t.Run(k.nama, func(t *testing.T) {
			lap, _ := rekonsiliasi([]barisJSON{barisJ("100044", k.umum)}, []barisJSON{barisJ("100044", k.inward)})
			if (len(lap.Gagal) > 0) != k.gagal || len(lap.K3) != k.k3 || (k.normalisasi != "") != (lap.Normalisasi[k.normalisasi] == 1) {
				t.Errorf("gagal %v (mau %v), K3 %d (mau %d), normalisasi %v (mau %q):\n%s", lap.Gagal, k.gagal, len(lap.K3),
					k.k3, lap.Normalisasi, k.normalisasi, lap.Teks())
			}
			if lap.Lolos() == (k.gagal || k.normalisasi != "") {
				t.Errorf("Lolos() = %v", lap.Lolos())
			}
			// -terima-normalisasi meloloskan normalisasi, TIDAK PERNAH kegagalan.
			lap.TerimaNormalisasi = true
			if lap.BolehDitulis() == k.gagal {
				t.Errorf("BolehDitulis() dengan -terima-normalisasi = %v, gagal %v", lap.BolehDitulis(), k.gagal)
			}
			if strings.Contains(lap.Teks(), "UJI-MATA-UANG") || strings.Contains(lap.Teks(), "UJI-ORG") {
				t.Error("laporan memuat nilai data")
			}
		})
	}
}

func TestJenisNormalisasi(t *testing.T) {
	for asal, mau := range map[string]string{
		"12.50": "nol ekor desimal", "1,5": "koma desimal", " 12": "spasi tepi", "+5": "tanda plus", "12.": "titik tanpa pecahan",
		".5": "tanpa nol depan", "-.5": "tanpa nol depan", "018": "nol depan", "-018": "nol depan", "0.5": "lain",
	} {
		if dapat := JenisNormalisasi(asal); dapat != mau {
			t.Errorf("JenisNormalisasi(%q) = %q, mau %q", asal, dapat, mau)
		}
	}
}

func TestRekonsiliasiInwardYatimDanProdukTanpaInward(t *testing.T) {
	lap, produk := rekonsiliasi([]barisJSON{barisJ("100044", `{"ID":"100044"}`)},
		[]barisJSON{barisJ("100070", `{"ID":"100070","PRODUCTID":"100071"}`)})
	if lap.ProdukTanpaInward != 1 || lap.InwardYatim != 1 || lap.Lolos() || len(produk) != 1 {
		t.Errorf("yatim menggagalkan, produk tanpa inward dihitung:\n%s", lap.Teks())
	}
	lap, _ = rekonsiliasi([]barisJSON{barisJ("100044", `{"ID":"100044"}`), barisJ("100044", `{"ID":"100044"}`)}, nil)
	if lap.Lolos() {
		t.Error("ID produk ganda harus gagal")
	}
	lap, _ = rekonsiliasi([]barisJSON{barisJ("100044", `{"ID":`)}, nil)
	if lap.Lolos() {
		t.Error("JSON rusak harus gagal")
	}
}
