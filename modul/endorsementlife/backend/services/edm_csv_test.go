package services_test

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"strings"
	"testing"

	"nusantarare/modul/endorsementlife/backend/models"
	"nusantarare/modul/endorsementlife/backend/services"
	"nusantarare/modul/endorsementlife/backend/tiruan"
)

// gudangKasusCSV - kasus Perubahan Data `EDML-1` (QR) dengan dua peserta Old.
func gudangKasusCSV(edmType string) *tiruan.Gudang {
	g := tiruan.Baru()
	g.Polis["EDML-1"] = &tiruan.Polis{ID: "EDML-1", OldPolicyNo: "UJI-PL-1", EdmType: edmType, ProdKe: 1,
		Kepala: map[string]string{"TYPE": models.TypeQR}}
	for _, id := range []string{"UJI-D1", "UJI-D2"} {
		g.Peserta = append(g.Peserta, &tiruan.Peserta{ID: id, PolisID: "EDML-1", PLNumber: "UJI-PL-1", EdmStatus: models.StatusOld,
			Nilai: map[string]string{"CERTIFICATE_NO": id, "PLAN": "UJI-PLAN", "POLICY_HOLDER": "UJI-PH", "CURRENCY": "IDR", "GROSS_PREMIUM": "100"}})
	}
	return g
}

// Judul memenuhi kolom wajib QR (keputusan work owner 10-10-2026); enam kolom pertama yang divariasikan uji,
// sisanya diisi `akhirBaris` oleh `csvUji` / `barisUji`.
const judulCSV = "CERTIFICATE_NO,PLAN,POLICY_HOLDER,CURRENCY,GROSS_PREMIUM,DOB," +
	"NAME_OF_INSURED,SEX,BEGIN_DATE,EXPIRED_DATE,SUM_INSURED,AGE,DEDUCTION,RI_ADMIN_FEE,BROKERAGE_FEE,TAX,PROF_COMM,CLAIM,NET_PREMIUM\n"

const akhirBaris = ",UJI-NAMA,M,01/01/2026,31/12/2026,1000,30,0,0,0,0,0,0,0"

// barisUji - enam kolom awal + kolom wajib lain yang sah.
func barisUji(awal string) string { return awal + akhirBaris }

func csvUji(baris ...string) *strings.Reader {
	var b strings.Builder
	b.WriteString(judulCSV)
	for _, x := range baris {
		if x != "" {
			x = barisUji(x)
		}
		b.WriteString(x + "\n")
	}
	return strings.NewReader(b.String())
}

func pesertaBaru(g *tiruan.Gudang) []*tiruan.Peserta {
	var hasil []*tiruan.Peserta
	for _, p := range g.Peserta {
		if p.EdmStatus == models.StatusNew {
			hasil = append(hasil, p)
		}
	}
	return hasil
}

func TestPeriksaCSVTanpaMenulis(t *testing.T) {
	g := gudangKasusCSV(models.EdmTypePerubahanData)
	l, j := layananJejak(g)
	h, err := l.PeriksaCSV(context.Background(), pelakuUji, "EDML-1", csvUji(
		"UJI-N1,UJI-PLAN,UJI-PH,IDR,\"10,5\",02/01/1990",
		"UJI-N2,UJI-LAIN,UJI-PH,IDR,1.234.567,02/01/1990",
		"",
		"UJI-N3,UJI-PLAN,UJI-PH,IDR,7,1990-01-02",
	))
	if err != nil {
		t.Fatal(err)
	}
	if h.Total != 3 || h.Ditolak != 2 || len(pesertaBaru(g)) != 0 || len(j.catatan) != 0 {
		t.Fatalf("tinjauan %+v; baris New %d", h, len(pesertaBaru(g)))
	}
	kolom := map[string]int{}
	for _, p := range h.Pesan {
		kolom[p.Kolom] = p.Baris
	}
	// Nomor = baris BERKAS (judul = 1): baris kosong ke-4 tidak menggeser laporan. PLAN beda tidak lagi ditolak
	// (nilai polis, 10-10-2026).
	if _, ada := kolom["PLAN"]; ada || kolom["GROSS_PREMIUM"] != 3 || kolom["DOB"] != 5 {
		t.Errorf("baris/kolom penolakan %v", h.Pesan)
	}
}

func TestTambahCSVMenyisipBaruLaluMengunci(t *testing.T) {
	g := gudangKasusCSV(models.EdmTypePerubahanData)
	l, j := layananJejak(g)
	ctx := context.Background()
	h, err := l.TambahCSV(ctx, pelakuUji, "EDML-1", csvUji(
		"UJI-N1,UJI-PLAN,UJI-PH,IDR,\"10,5\",02/01/1990",
		"UJI-N2,UJI-PLAN,UJI-PH,IDR,0,03/04/1991",
	))
	if err != nil {
		t.Fatal(err)
	}
	baru := pesertaBaru(g)
	if h.Disimpan != 2 || len(baru) != 2 || baru[0].Nilai["GROSS_PREMIUM"] != "10.5" || baru[0].Nilai["DOB"] != "1990-01-02" ||
		baru[1].Nilai["GROSS_PREMIUM"] != "0" || baru[1].Nilai["DOB"] != "1991-04-03" || baru[0].PLNumber != "UJI-PL-1" {
		t.Fatalf("%+v; baris New %+v %+v", h, baru[0], baru[1])
	}
	if len(j.catatan) != 1 || j.catatan[0].Komentar != "Add CSV Data" {
		t.Errorf("jejak %+v", j.catatan)
	}
	k, _ := l.BacaKasus(ctx, pelakuUji, "EDML-1")
	if !k.CSVTerkunci {
		t.Error("Upload CSV tidak terkunci sesudah Add CSV Data (.EditInput1 = 1)")
	}
	// Unggahan kedua ditolak: tidak menumpuk.
	if _, err := l.TambahCSV(ctx, pelakuUji, "EDML-1", csvUji("UJI-N9,UJI-PLAN,UJI-PH,IDR,1,")); !errors.Is(err, services.ErrCSVTerkunci) {
		t.Fatalf("unggahan kedua: %v", err)
	}
	if len(pesertaBaru(g)) != 2 {
		t.Errorf("baris New %d sesudah unggahan kedua, mau tetap 2", len(pesertaBaru(g)))
	}
}

func TestTambahCSVSatuPenolakanMenolakSeluruhBerkas(t *testing.T) {
	g := gudangKasusCSV(models.EdmTypePerubahanData)
	l, j := layananJejak(g)
	_, err := l.TambahCSV(context.Background(), pelakuUji, "EDML-1", csvUji(
		"UJI-N1,UJI-PLAN,UJI-PH,IDR,1,02/01/1990",
		"UJI-N2,UJI-PLAN,UJI-PH,EUR,1,02/01/1990",
	))
	var gc services.GalatCSV
	if !errors.As(err, &gc) || gc.Ditolak != 1 || gc.Total != 2 || gc.Pesan[0].Kolom != "CURRENCY" || gc.Pesan[0].Baris != 3 {
		t.Fatalf("%v %+v", err, gc)
	}
	if len(pesertaBaru(g)) != 0 || len(j.catatan) != 0 {
		t.Error("baris sebelum pelanggaran ikut tersimpan (R24: Pega berhenti SESUDAH menambah)")
	}
}

// TestTambahCSVTanpaBatasBaris - AC 36: di atas 50.000 baris tetap diterima.
func TestTambahCSVTanpaBatasBaris(t *testing.T) {
	g := gudangKasusCSV(models.EdmTypePerubahanData)
	l, _ := layananJejak(g)
	var b strings.Builder
	b.WriteString(judulCSV)
	const n = 50_001
	for i := 1; i <= n; i++ {
		fmt.Fprintf(&b, "%s\n", barisUji(fmt.Sprintf("UJI-N%d,UJI-PLAN,UJI-PH,IDR,1,02/01/1990", i)))
	}
	h, err := l.TambahCSV(context.Background(), pelakuUji, "EDML-1", strings.NewReader(b.String()))
	if err != nil || h.Disimpan != n || len(pesertaBaru(g)) != n {
		t.Fatalf("%d baris: %+v %v", n, h, err)
	}
}

// TestCSVSesudahProsesDitolak - Process Endorsement langkah FINAL (keputusan
// work owner 08-10-2026, mengganti R31): sesudah Save, Upload CSV (periksa) dan
// Add CSV Data ditolak 409 `ErrSudahDiproses` - peserta, rekap, jejak tetap.
func TestCSVSesudahProsesDitolak(t *testing.T) {
	g := gudangKasusCSV(models.EdmTypePerubahanData)
	l, j := layananJejak(g)
	ctx := context.Background()
	if _, err := l.Simpan(ctx, pelakuUji, "EDML-1", models.PilihanHapus{Pilih: []string{"UJI-D1"}}); err != nil {
		t.Fatal(err)
	}
	sebelum := len(g.Peserta)
	rekap, jejakSebelum := g.Rekap["EDML-1"], len(j.catatan)
	if _, err := l.PeriksaCSV(ctx, pelakuUji, "EDML-1", csvUji("UJI-N1,UJI-PLAN,UJI-PH,IDR,50,")); !errors.Is(err, services.ErrSudahDiproses) {
		t.Fatalf("periksa: %v", err)
	}
	_, err := l.TambahCSV(ctx, pelakuUji, "EDML-1", csvUji("UJI-N1,UJI-PLAN,UJI-PH,IDR,50,"))
	if !errors.Is(err, services.ErrSudahDiproses) || services.Pesan(err) != models.PesanSudahDiproses {
		t.Fatalf("tambah: %v", err)
	}
	if len(g.Peserta) != sebelum || len(pesertaBaru(g)) != 0 || g.Rekap["EDML-1"] != rekap || len(j.catatan) != jejakSebelum {
		t.Errorf("tulisan sesudah ditolak: peserta %d→%d rekap %d→%d jejak %d→%d",
			sebelum, len(g.Peserta), rekap, g.Rekap["EDML-1"], jejakSebelum, len(j.catatan))
	}
}

func TestCSVDitolak(t *testing.T) {
	ctx := context.Background()
	tertutup := gudangKasusCSV(models.EdmTypePerubahanData)
	tertutup.Polis["EDML-1"].Status = models.StatusKasusDitolak
	kosong := gudangKasusCSV(models.EdmTypePerubahanData)
	kosong.Peserta = nil
	gagal := gudangKasusCSV(models.EdmTypePerubahanData)
	gagal.GagalSisipKe = 2 // baris berkas pertama sesudah judul
	for _, k := range []struct {
		nama  string
		g     *tiruan.Gudang
		isi   string
		galat error
	}{
		{"Batal", gudangKasusCSV(models.EdmTypeBatal), judulCSV + "UJI-N1,UJI-PLAN,UJI-PH,IDR,1,\n", services.ErrCSVBukanPerubahanData},
		{"tertutup", tertutup, judulCSV + "UJI-N1,UJI-PLAN,UJI-PH,IDR,1,\n", services.ErrKasusTertutup},
		{"tanpa peserta lama", kosong, judulCSV + "UJI-N1,UJI-PLAN,UJI-PH,IDR,1,\n", services.ErrCSVTanpaAcuan},
		{"kosong", gudangKasusCSV("1"), judulCSV, services.ErrCSVKosong},
		{"tanpa judul", gudangKasusCSV("1"), "", services.ErrCSVKosong},
		{"tanpa POLICY_HOLDER", gudangKasusCSV("1"), "PLAN,CURRENCY\nUJI-PLAN,IDR\n", services.ErrMasukanTidakSah},
		{"kutip rusak", gudangKasusCSV("1"), judulCSV + "UJI-N1,\"UJI-PLAN,UJI-PH,IDR,1,\n", services.ErrCSVRusak},
	} {
		l, _ := layananJejak(k.g)
		if _, err := l.TambahCSV(ctx, pelakuUji, "EDML-1", strings.NewReader(k.isi)); !errors.Is(err, k.galat) {
			t.Errorf("%s: %v, mau %v", k.nama, err, k.galat)
		}
	}
	l, j := layananJejak(gagal)
	if _, err := l.TambahCSV(ctx, pelakuUji, "EDML-1", csvUji("UJI-N1,UJI-PLAN,UJI-PH,IDR,1,02/01/1990")); err == nil || len(j.catatan) != 0 {
		t.Errorf("galat sisip tidak diteruskan: %v", err)
	}
}

// TestTambahCSVGagalDiTengahMembatalkanSemua - kelompok pertama (1000 baris)
// sudah tersisip saat baris ke-1500 gagal: seluruhnya batal bersama transaksi.
func TestTambahCSVGagalDiTengahMembatalkanSemua(t *testing.T) {
	g := gudangKasusCSV(models.EdmTypePerubahanData)
	g.GagalSisipKe = 1501 // baris berkas rekaman ke-1500 (judul = baris 1)
	l := services.BaruLayanan(g, tiruan.TransaksiPulih(g), nil).DenganJejak(&jejakTiruan{})
	var b strings.Builder
	b.WriteString(judulCSV)
	for i := 1; i <= 2500; i++ {
		fmt.Fprintf(&b, "%s\n", barisUji(fmt.Sprintf("UJI-N%d,UJI-PLAN,UJI-PH,IDR,1,02/01/1990", i)))
	}
	if _, err := l.TambahCSV(context.Background(), pelakuUji, "EDML-1", strings.NewReader(b.String())); err == nil {
		t.Fatal("galat sisip di kelompok kedua tidak diteruskan")
	}
	if n := len(pesertaBaru(g)); n != 0 {
		t.Fatalf("%d baris New tersisa sesudah pembatalan, mau 0", n)
	}
}

// TestCSVWajibMataUangSertifikat - validasi CSV keputusan work owner 10-10-2026 di layanan: kolom wajib di judul,
// mata uang polis, CERTIFICATE_NO ganda dalam berkas dan bentrok peserta Old (peserta Delete tidak dihitung).
func TestCSVWajibMataUangSertifikat(t *testing.T) {
	ctx := context.Background()
	g := gudangKasusCSV(models.EdmTypePerubahanData)
	g.Peserta = append(g.Peserta, &tiruan.Peserta{ID: "UJI-DEL", PolisID: "EDML-1", PLNumber: "UJI-PL-1", EdmStatus: models.StatusDelete,
		Nilai: map[string]string{"CERTIFICATE_NO": "UJI-HAPUS", "PLAN": "UJI-PLAN", "POLICY_HOLDER": "UJI-PH", "CURRENCY": "USD"}})
	l, _ := layananJejak(g)
	h, err := l.PeriksaCSV(ctx, pelakuUji, "EDML-1", csvUji(
		"UJI-N1,UJI-PLAN,UJI-PH,IDR,1,02/01/1990",
		"uji-n1 ,UJI-PLAN,UJI-PH,IDR,1,02/01/1990",   // ganda dalam berkas (tidak peka huruf)
		"UJI-D2,UJI-PLAN,UJI-PH,IDR,1,02/01/1990",    // sertifikat peserta Old
		"UJI-HAPUS,UJI-PLAN,UJI-PH,IDR,1,02/01/1990", // milik peserta Delete: boleh
		"UJI-N5,UJI-PLAN,UJI-PH,USD,1,02/01/1990",    // USD hanya milik peserta Delete
		"UJI-N6,UJI-PLAN,UJI-PH,EUR,1,02/01/1990",
	))
	if err != nil {
		t.Fatal(err)
	}
	mau := map[int]string{
		3: `Certificate No "uji-n1" is used more than once in this file (also in row 2).`,
		4: `Certificate No "UJI-D2" is already registered in this policy.`,
		6: "CURRENCY USD is not used in this policy (use IDR).",
		7: models.PesanCSVMataUang,
	}
	if h.Total != 6 || h.Ditolak != 4 || len(h.Pesan) != 4 {
		t.Fatalf("tinjauan %+v", h)
	}
	for _, p := range h.Pesan {
		if mau[p.Baris] != p.Pesan {
			t.Errorf("baris %d: %q, mau %q", p.Baris, p.Pesan, mau[p.Baris])
		}
	}
	// Judul tanpa kolom wajib: berkas ditolak di tahap judul, seluruh kolom yang kurang disebut.
	_, err = l.PeriksaCSV(ctx, pelakuUji, "EDML-1", strings.NewReader("PLAN,POLICY_HOLDER,CERTIFICATE_NO\nUJI-PLAN,UJI-PH,UJI-N1\n"))
	if !errors.Is(err, services.ErrMasukanTidakSah) || services.Pesan(err) != "Your file is missing these columns: NAME_OF_INSURED, SEX, DOB, CURRENCY, BEGIN_DATE, EXPIRED_DATE, SUM_INSURED, NET_PREMIUM, AGE or ENTRY_AGE, GROSS_PREMIUM, DEDUCTION, RI_ADMIN_FEE, BROKERAGE_FEE." {
		t.Errorf("judul kurang: %v", err)
	}
	// Kolom wajib kasus untuk layar (satu sumber).
	k, _ := l.BacaKasus(ctx, pelakuUji, "EDML-1")
	if strings.Join(k.KolomWajibCSV, ",") != strings.Join(models.DaftarKolomWajibCSV(models.TypeQR), ",") {
		t.Errorf("kolom wajib kasus %v", k.KolomWajibCSV)
	}
}

// TestTemplatCSVBolakBalik - Download Template (10-10-2026): judul kolom wajib di depan, baris contoh dari peserta
// Old pertama LOLOS unggahan bila CERTIFICATE_NO / NAME_OF_INSURED diisi; gerbang sama dengan unggahan.
func TestTemplatCSVBolakBalik(t *testing.T) {
	ctx := context.Background()
	g := gudangKasusCSV(models.EdmTypePerubahanData)
	for k, v := range map[string]string{"NAME_OF_INSURED": "UJI-NAMA-1", "SEX": "F", "DOB": "02/01/1990", "BEGIN_DATE": "01/01/2026",
		"EXPIRED_DATE": "31/12/2026", "SUM_INSURED": "1000.5", "AGE": "36", "DEDUCTION": "1.25", "STNC": "2026-04-03"} {
		g.Peserta[0].Nilai[k] = v
	}
	// Tiruan menyimpan tanggal DATE sebagai YYYY-MM-DD; SQL nyata memberi DD/MM/YYYY (TO_CHAR) - samakan.
	g.Peserta[0].Nilai["DOB"], g.Peserta[0].Nilai["BEGIN_DATE"], g.Peserta[0].Nilai["EXPIRED_DATE"] = "02/01/1990", "01/01/2026", "31/12/2026"
	l, _ := layananJejak(g)
	b, err := l.TemplatCSV(ctx, pelakuUji, "EDML-1")
	if err != nil || b.Nama != "Template_EDM_EDML-1.csv" {
		t.Fatalf("%+v %v", b.Nama, err)
	}
	r := csv.NewReader(bytes.NewReader(b.Isi))
	rek, err := r.ReadAll()
	if err != nil || len(rek) != 2 {
		t.Fatalf("templat %v %v", rek, err)
	}
	judul, contoh := rek[0], rek[1]
	wajib := models.KolomWajibCSV(models.TypeQR)
	if len(judul) != len(models.PemetaanCSV) || strings.Join(judul[:len(wajib)], ",") != strings.Join(wajib, ",") {
		t.Fatalf("judul %v", judul)
	}
	idx := map[string]int{}
	for i, k := range judul {
		idx[k] = i
	}
	if contoh[idx["CERTIFICATE_NO"]] != "" || contoh[idx["NAME_OF_INSURED"]] != "" || contoh[idx["PLAN"]] != "UJI-PLAN" ||
		contoh[idx["STNC"]] != "03/04/2026" || contoh[idx["GROSS_PREMIUM"]] != "100" {
		t.Errorf("contoh %v", contoh)
	}
	contoh[idx["CERTIFICATE_NO"]], contoh[idx["NAME_OF_INSURED"]] = "UJI-N1", "UJI-NAMA"
	var w bytes.Buffer
	cw := csv.NewWriter(&w)
	_ = cw.Write(judul)
	_ = cw.Write(contoh)
	cw.Flush()
	h, err := l.PeriksaCSV(ctx, pelakuUji, "EDML-1", bytes.NewReader(w.Bytes()))
	if err != nil || h.Ditolak != 0 || h.Total != 1 {
		t.Fatalf("baris contoh tidak lolos unggahan: %+v %v", h, err)
	}
	// Gerbang sama dengan unggahan.
	tutup := gudangKasusCSV(models.EdmTypePerubahanData)
	tutup.Polis["EDML-1"].Status = models.StatusKasusDitolak
	lt, _ := layananJejak(tutup)
	if _, err := lt.TemplatCSV(ctx, pelakuUji, "EDML-1"); !errors.Is(err, services.ErrKasusTertutup) {
		t.Errorf("kasus tertutup: %v", err)
	}
	lb, _ := layananJejak(gudangKasusCSV(models.EdmTypeBatal))
	if _, err := lb.TemplatCSV(ctx, pelakuUji, "EDML-1"); !errors.Is(err, services.ErrCSVBukanPerubahanData) {
		t.Errorf("Batal: %v", err)
	}
}

// berkasPemisahUji - judulCSV + baris barisUji berpemisah `sep`; pada ";" GROSS_PREMIUM diberi desimal koma (berkutip
// atau tidak, seperti Excel Indonesia). `awal` ditaruh di depan berkas (mis. BOM).
func berkasPemisahUji(awal, sep string, kutip bool, baris ...string) string {
	var b strings.Builder
	b.WriteString(awal)
	b.WriteString(strings.ReplaceAll(strings.TrimSuffix(judulCSV, "\n"), ",", sep) + "\r\n")
	for _, x := range baris {
		medan := strings.Split(barisUji(x), ",")
		if sep == ";" {
			medan[4] = "14216802,6"
			if kutip {
				medan[4] = `"14216802,6"`
			}
		}
		b.WriteString(strings.Join(medan, sep) + "\r\n")
	}
	return b.String()
}

// TestPemisahKomaDanTitikKoma - opsi a (keputusan work owner 10-10-2026): pemisah dideteksi dari judul, BOM
// diabaikan, desimal koma pada berkas ";" lolos, nomor baris tetap (judul = 1), "sep=" tidak didukung.
func TestPemisahKomaDanTitikKoma(t *testing.T) {
	ctx := context.Background()
	sah := "UJI-N1,UJI-PLAN,UJI-PH,IDR,1,02/01/1990"
	salah := "UJI-N2,UJI-PLAN,UJI-PH,EUR,1,02/01/1990" // ditolak di CURRENCY, baris berkas 3
	for _, k := range []struct {
		nama, awal, sep string
		kutip           bool
	}{
		{"koma", "", ",", false},
		{"koma ber-BOM", "\uFEFF", ",", false},
		{"titik koma, desimal koma tanpa kutip", "", ";", false},
		{"titik koma ber-BOM, desimal koma berkutip", "\uFEFF", ";", true},
	} {
		l, _ := layananJejak(gudangKasusCSV(models.EdmTypePerubahanData))
		h, err := l.PeriksaCSV(ctx, pelakuUji, "EDML-1", strings.NewReader(berkasPemisahUji(k.awal, k.sep, k.kutip, sah, salah)))
		if err != nil {
			t.Fatalf("%s: %v", k.nama, err)
		}
		if h.Total != 2 || h.Ditolak != 1 || len(h.Pesan) != 1 || h.Pesan[0].Baris != 3 || h.Pesan[0].Kolom != "CURRENCY" {
			t.Errorf("%s: %+v", k.nama, h)
		}
		g := gudangKasusCSV(models.EdmTypePerubahanData)
		lt, _ := layananJejak(g)
		ht, err := lt.TambahCSV(ctx, pelakuUji, "EDML-1", strings.NewReader(berkasPemisahUji(k.awal, k.sep, k.kutip, sah)))
		if err != nil || ht.Disimpan != 1 {
			t.Fatalf("%s tambah: %+v %v", k.nama, ht, err)
		}
		mau := "1"
		if k.sep == ";" {
			mau = "14216802.6"
		}
		if gp := pesertaBaru(g)[0].Nilai["GROSS_PREMIUM"]; gp != mau {
			t.Errorf("%s: GROSS_PREMIUM %q, mau %q", k.nama, gp, mau)
		}
	}
	// "sep=" tidak didukung: dibaca sebagai judul → kolom wajib kurang.
	l, _ := layananJejak(gudangKasusCSV(models.EdmTypePerubahanData))
	if _, err := l.PeriksaCSV(ctx, pelakuUji, "EDML-1", strings.NewReader("sep=;\r\n"+berkasPemisahUji("", ";", false, sah))); !errors.Is(err, services.ErrMasukanTidakSah) {
		t.Errorf("sep=: %v", err)
	}
}

// TestPlanPolicyHolderDariPolis - PLAN / POLICY_HOLDER peserta baru dari polis (keputusan work owner 10-10-2026,
// seperti PremiumList Life): tidak wajib di judul, nilai CSV diabaikan; kepala kasus dulu, lalu peserta Old;
// keduanya kosong → unggahan ditolak.
func TestPlanPolicyHolderDariPolis(t *testing.T) {
	ctx := context.Background()
	// Judul TANPA PLAN / POLICY_HOLDER.
	judul := "CERTIFICATE_NO,CURRENCY,GROSS_PREMIUM,DOB,NAME_OF_INSURED,SEX,BEGIN_DATE,EXPIRED_DATE,SUM_INSURED,AGE,DEDUCTION,RI_ADMIN_FEE,BROKERAGE_FEE,NET_PREMIUM\n"
	baris := "UJI-N1,IDR,1,02/01/1990,UJI-NAMA,M,01/01/2026,31/12/2026,1000,30,0,0,0,0\n"

	// Kepala kasus terisi: nilai kepala dipakai (walau berbeda dari peserta Old).
	g := gudangKasusCSV(models.EdmTypePerubahanData)
	g.Polis["EDML-1"].Kepala["BUSINESS_NAME"] = "GROUP TERM LIFE"
	g.Polis["EDML-1"].Kepala["POLICY_HOLDER_NAME"] = "UJI-PEMEGANG-POLIS"
	l, _ := layananJejak(g)
	if _, err := l.TambahCSV(ctx, pelakuUji, "EDML-1", strings.NewReader(judul+baris)); err != nil {
		t.Fatal(err)
	}
	if b := pesertaBaru(g)[0].Nilai; b["PLAN"] != "GROUP TERM LIFE" || b["POLICY_HOLDER"] != "UJI-PEMEGANG-POLIS" {
		t.Errorf("nilai kepala: %q %q", b["PLAN"], b["POLICY_HOLDER"])
	}

	// Nilai CSV berbeda diabaikan; kepala kosong → peserta Old pertama.
	g = gudangKasusCSV(models.EdmTypePerubahanData)
	l, _ = layananJejak(g)
	if _, err := l.TambahCSV(ctx, pelakuUji, "EDML-1", csvUji("UJI-N1,UJI-PLAN-CSV,UJI-PH-CSV,IDR,1,02/01/1990")); err != nil {
		t.Fatal(err)
	}
	if b := pesertaBaru(g)[0].Nilai; b["PLAN"] != "UJI-PLAN" || b["POLICY_HOLDER"] != "UJI-PH" {
		t.Errorf("fallback peserta Old: %q %q", b["PLAN"], b["POLICY_HOLDER"])
	}

	// Kepala dan peserta Old kosong → ditolak seluruhnya, berpesan.
	g = gudangKasusCSV(models.EdmTypePerubahanData)
	for _, p := range g.Peserta {
		p.Nilai["PLAN"], p.Nilai["POLICY_HOLDER"] = "", ""
	}
	l, _ = layananJejak(g)
	_, err := l.PeriksaCSV(ctx, pelakuUji, "EDML-1", strings.NewReader(judul+baris))
	if !errors.Is(err, services.ErrMasukanTidakSah) || !strings.Contains(err.Error(), "POLICY_HOLDER and PLAN are empty on the policy") {
		t.Errorf("keduanya kosong: %v", err)
	}
	// Required columns tanpa PLAN / POLICY_HOLDER.
	k, _ := layananUji(gudangKasusCSV(models.EdmTypePerubahanData)).BacaKasus(ctx, pelakuUji, "EDML-1")
	for _, c := range k.KolomWajibCSV {
		if c == "PLAN" || c == "POLICY_HOLDER" {
			t.Errorf("Required columns memuat %s", c)
		}
	}
}
