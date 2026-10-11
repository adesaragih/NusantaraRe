package repository

// RALAT 10-10-2026 (keputusan work owner "saat register claim, data peserta belum di insert ke tabel
// os_akseptasi_klaim_life; itu hanya terjadi saat save to rnm bersamaan dengan no klaim dan no akseptasi"): baris datar
// OS_AKSEPTASI_KLAIM_LIFE lahir di Save to RNM (TulisCerminSaveRNM), tidak lagi di Register (Simpan). TANPA Oracle.

import (
	"os"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
)

func TestSQLPerbaruiCerminIdempotenDanTidakMenyentuhStatus(t *testing.T) {
	q := sqlPerbaruiCermin("S.L")
	if err := db.PeriksaSQL(q); err != nil {
		t.Fatal(err)
	}
	satu := strings.Join(strings.Fields(q), " ")
	// Kunci baris = ID adjustment + CASEID klaim; baris berkeputusan tidak disentuh.
	if !strings.HasSuffix(satu, "WHERE ID = :45 AND CASEID = :46 AND (STS_REJECT IS NULL OR STS_REJECT = :47)") {
		t.Errorf("kunci UPDATE baris datar berubah:\n%s", satu)
	}
	set := satu[strings.Index(satu, " SET ")+5 : strings.Index(satu, " WHERE ")]
	// STS_REJECT tidak disentuh UPDATE; ACCEPTATION_DATE kini DITULIS (tanggal outstanding, keputusan WO 11-10-2026).
	// Lima kolom Medical Check (STATUS, DISEASE, ICD_CODE, NOTES, KETERANGAN) TIDAK disentuh: nilai tahap lain tidak hilang.
	for _, larang := range []string{"STS_REJECT", "STATUS", "DISEASE", "ICD_CODE", "NOTES", "KETERANGAN", "NAME_OF_INSURED", "DOB ="} {
		if strings.Contains(" "+set, " "+strings.TrimSuffix(larang, " =")+" =") {
			t.Errorf("UPDATE baris datar menyentuh %s:\n%s", larang, set)
		}
	}
	if !strings.Contains(set, "ACCEPTATION_DATE = TO_DATE(:44, 'YYYY-MM-DD HH24:MI:SS')") {
		t.Errorf("ACCEPTATION_DATE tidak ditulis sebagai DATE:\n%s", set)
	}
	// Ditimpa, BUKAN NVL: Save to RNM membekukan nilai polis/peserta/baris.
	if strings.Contains(set, "NVL(") {
		t.Errorf("UPDATE masih memakai NVL - keputusan work owner 11-10-2026: timpa:\n%s", set)
	}
	for i, k := range KolomCerminSaveRNM {
		if !strings.Contains(set, k+" = ") {
			t.Errorf("kolom %s (#%d) tidak ditimpa", k, i+1)
		}
	}
	for _, mau := range []string{"NO_ACCEPTATION = :2,", "LAPSE_DATE = TO_DATE(:10, 'YYYY-MM-DD HH24:MI:SS')", "WPC = TO_DATE(:19,"} {
		if !strings.Contains(set, mau) {
			t.Errorf("SET tanpa %q:\n%s", mau, set)
		}
	}
	if strings.Contains(strings.ToUpper(q), "COMMIT") {
		t.Error("UPDATE baris datar memuat COMMIT")
	}
}

func TestSQLSisipCerminSeluruhKolomTerpetakan(t *testing.T) {
	q := sqlSisipCermin("S.L")
	if err := db.PeriksaSQL(q); err != nil {
		t.Fatal(err)
	}
	satu := strings.Join(strings.Fields(q), " ")
	kolom := strings.Split(satu[strings.Index(satu, "(")+1:strings.Index(satu, ")")], ", ")
	if len(kolom) != 47 || kolom[0] != "ID" || kolom[1] != "CASEID" || kolom[45] != "ACCEPTATION_DATE" || kolom[46] != "STS_REJECT" {
		t.Errorf("kolom INSERT: %v", kolom)
	}
	// Tanggal lewat TO_DATE berformat tetap; tanda bind berurut 1..47.
	if !strings.Contains(satu, "TO_DATE(:13, 'YYYY-MM-DD HH24:MI:SS')") || !strings.HasSuffix(satu, "TO_DATE(:46, 'YYYY-MM-DD HH24:MI:SS'), :47)") {
		t.Errorf("tanda VALUES: %s", satu)
	}
}

func TestSQLSOBPolisKunciIdPega(t *testing.T) {
	q := strings.Join(strings.Fields(sqlSOBPolis("S.PL")), " ")
	if err := db.PeriksaSQL(q); err != nil {
		t.Fatal(err)
	}
	if q != "SELECT p.SOB FROM S.PL p WHERE p.ID = :1 AND p.ID_PEGA = p.ID ORDER BY NVL(p.PROD_KE, 0) DESC FETCH FIRST 1 ROWS ONLY" {
		t.Errorf("SOB polis: %s", q)
	}
}

// Lebar kolom katalog dihormati: teks lebih panjang = galat jelas (bukan dipotong), menyebut kolom tanpa nilainya.
func TestLebarWarisanTidakDipotong(t *testing.T) {
	b := BarisLama{ID: "UJI-A1", BUSINESSNAME: strings.Repeat("x", 100), TYPE: "QR"}
	if err := PeriksaLebarWarisan(b); err != nil {
		t.Errorf("tepat 100 byte ditolak: %v", err)
	}
	b.BUSINESSNAME += "y"
	err := PeriksaLebarWarisan(b)
	if err == nil || !strings.Contains(err.Error(), "BUSINESSNAME 101 byte melebihi lebar VARCHAR2(100)") || strings.Contains(err.Error(), "xxx") {
		t.Errorf("galat lebar: %v", err)
	}
}

func TestNilaiTulisLamaTanggalDanKosong(t *testing.T) {
	b := BarisLama{ID: "UJI-A1", LAPSE_DATE: "2025-06-01", NO_ACCEPTATION: "11/10/2026", SEX: " "}
	n, err := NilaiTulisLama(b, []string{"LAPSE_DATE", "NO_ACCEPTATION", "SEX"})
	if err != nil || n[0] != "2025-06-01 00:00:00" || n[1] != "11/10/2026" || n[2] != nil {
		t.Errorf("%v %v", n, err)
	}
	if _, err := NilaiTulisLama(BarisLama{ID: "X", BEGIN_DATE: "bukan tanggal"}, []string{"BEGIN_DATE"}); err == nil {
		t.Error("tanggal tak terbaca harus galat, bukan NULL diam-diam")
	}
}

func TestCerminHanyaDilahirkanUntukBarisBelumDisimpan(t *testing.T) {
	for kode, mau := range map[string]bool{"": true, "  ": true, "0": false, "1": false, "2": false} {
		if got := CerminBolehDisisip(kode); got != mau {
			t.Errorf("CerminBolehDisisip(%q) = %v, mau %v", kode, got, mau)
		}
	}
}

// Urutan TulisCerminSaveRNM: UPDATE dulu, INSERT hanya bila nol baris diperbarui DAN boleh disisip - Save to RNM yang
// ditekan ulang tidak menggandakan baris datar.
func TestTulisCerminSaveRNMPerbaruiDuluLaluSisip(t *testing.T) {
	isi, err := os.ReadFile("pohonklaim.go")
	if err != nil {
		t.Fatal(err)
	}
	fn := potongFungsi(buangKomentarGo(string(isi)), "func (r *PohonKlaim) TulisCerminSaveRNM(")
	if fn == "" {
		t.Fatal("TulisCerminSaveRNM tidak ditemukan; pembacanya yang rusak")
	}
	perbarui := strings.Index(fn, "tx.ExecContext(ctx, perbarui,")
	jikaNol := strings.Index(fn, "if n == 0 {")
	tolak := strings.Index(fn, "if !sisip {")
	sisip := strings.Index(fn, "tx.ExecContext(ctx, sisipQ,")
	if perbarui < 0 || jikaNol < 0 || tolak < 0 || sisip < 0 || !(perbarui < jikaNol && jikaNol < tolak && tolak < sisip) {
		t.Errorf("urutan UPDATE -> (n == 0) -> sisip? -> INSERT berubah: %d %d %d %d", perbarui, jikaNol, tolak, sisip)
	}
	if strings.Contains(fn, "DalamTransaksi(") || strings.Contains(fn, "Commit(") {
		t.Error("TulisCerminSaveRNM membuka/menutup transaksi sendiri; ia wajib ikut transaksi utama Save to RNM")
	}
}

// Hapus klaim / cabut peserta SEBELUM Save to RNM: baris datar belum ada, dan itu bukan galat.
func TestHapusSebelumSaveRNMTidakMenuntutBarisDatar(t *testing.T) {
	isi, err := os.ReadFile("pohonklaim.go")
	if err != nil {
		t.Fatal(err)
	}
	kode := buangKomentarGo(string(isi))
	for _, kepala := range []string{"func (r *KlaimLife) HapusCerminBelumDisimpan(", "func (r *PohonKlaim) HapusFisik("} {
		fn := potongFungsi(kode, kepala)
		if fn == "" {
			t.Fatalf("%s tidak ditemukan", kepala)
		}
		if strings.Contains(fn, "RowsAffected") {
			t.Errorf("%s memeriksa cacah baris datar - klaim yang belum Save to RNM tidak punya baris datar", kepala)
		}
	}
}

// 11-10-2026: ambang hari produk dari tabel DATAR Master Product Name Life (dua kolom yang sama, berkunci ID).
func TestSQLAmbangProdukTabelDatar(t *testing.T) {
	q := strings.Join(strings.Fields(sqlAmbangProdukDatar("S.MPNL")), " ")
	if err := db.PeriksaSQL(q); err != nil {
		t.Fatal(err)
	}
	if q != "SELECT t.MAXEXPIREDCLAIM, t.MAXDATARECEIVE FROM S.MPNL t WHERE t.ID = :1" || TabelProdukNameLife != "M_PRODUCTNAME_LIFE" {
		t.Errorf("ambang datar: %s", q)
	}
	isi, err := os.ReadFile("ambangproduk.go")
	if err != nil {
		t.Fatal(err)
	}
	fn := potongFungsi(buangKomentarGo(string(isi)), "func (r *ProdukLife) Ambang(")
	datar, view := strings.Index(fn, "sqlAmbangProdukDatar("), strings.Index(fn, "sqlAmbangProduk(produk)")
	if datar < 0 || view < 0 || datar > view || !strings.Contains(fn, "errors.Is(err, sql.ErrNoRows)") {
		t.Errorf("urutan baca: tabel datar dulu (%d), view hanya bila tidak ada baris (%d)", datar, view)
	}
}
