package repository

// Uji pelengkap migrasi data A4 - butir at/au/bb, dokumen warisan, diagnosa.
//
// Seluruhnya MURNI: nol Oracle, nol jam. Yang diuji adalah apa yang
// DILAPORKAN, bukan hanya apa yang dipindahkan - sebab inti A4 adalah bahwa
// yang tidak dapat dipindahkan dengan yakin TIDAK ditebak.

import (
	"os"
	"strings"
	"testing"

	"nusantarare/internal/models"
)

// punyaTemuan mencari satu jenis temuan di dalam daftar.
func punyaTemuan(temuan []Temuan, jenis string) *Temuan {
	for i := range temuan {
		if temuan[i].Jenis == jenis {
			return &temuan[i]
		}
	}
	return nil
}

func TestTahapSelaluDilaporkanTidakAdaDiSumber(t *testing.T) {
	// ⛔ Inti butir at. Tabel datar 55 kolom tidak punya satu pun kolom
	// tahap; mengisinya dengan tebakan berarti kasus warisan mendarat di
	// antrean yang salah, dan orang di antrean itu tidak akan tahu kenapa.
	lengkap, temuan := LengkapiWork([]BarisLama{{CASEID: "UJI-CASE-1"}})
	if lengkap.Tahap != "" {
		t.Errorf("TAHAP = %q, mau kosong; sumbernya tidak punya kolom itu", lengkap.Tahap)
	}
	tm := punyaTemuan(temuan, TemuanTahapTidakAdaDiSumber)
	if tm == nil {
		t.Fatal("ketiadaan tahap tidak dilaporkan; ia akan lolos diam-diam")
	}
	if !strings.Contains(tm.Catatan, "55 kolom") {
		t.Errorf("catatan tidak menyebut sumbernya: %q", tm.Catatan)
	}
}

func TestWaktuLahirDiisiKolomTerdekatDanDilaporkan(t *testing.T) {
	// ⛔ `CLAIM_RECEIVED_DATE` BUKAN `pxCreateDateTime`. Ia dipakai karena ia
	// satu-satunya kolom waktu yang ada dan masuk akal - dan penggantian itu
	// harus terbaca di laporan, bukan tersembunyi di dalam kode.
	lengkap, temuan := LengkapiWork([]BarisLama{{
		CASEID: "UJI-CASE-1", ID: "UJI-1", CLAIM_RECEIVED_DATE: "2026-01-15",
	}})
	if lengkap.TglCreate.IsZero() {
		t.Fatal("TGL_CREATE kosong padahal CLAIM_RECEIVED_DATE terisi")
	}
	tm := punyaTemuan(temuan, TemuanWaktuLahirDiganti)
	if tm == nil {
		t.Fatal("penggantian kolom tidak dilaporkan")
	}
	if !strings.Contains(tm.Catatan, "pxCreateDateTime") {
		t.Errorf("catatan tidak menyebut apa yang SEHARUSNYA: %q", tm.Catatan)
	}
	// Tanggal yang tidak terurai DILAPORKAN, bukan membuat proses berhenti.
	_, temuan = LengkapiWork([]BarisLama{{
		CASEID: "UJI-CASE-2", ID: "UJI-2", CLAIM_RECEIVED_DATE: "bukan tanggal",
	}})
	if punyaTemuan(temuan, TemuanTanggalTakTerurai) == nil {
		t.Error("tanggal rusak tidak dilaporkan")
	}
}

func TestStatusKerjaHanyaDariCompleteDate(t *testing.T) {
	// ⛔ Kesimpulan, dan dinyatakan sebagai kesimpulan. Tanpa COMPLETE_DATE
	// kolomnya KOSONG - dan kosong berarti "belum ditutup" (migrasi 017),
	// bukan "tidak diketahui".
	kosong, temuan := LengkapiWork([]BarisLama{{CASEID: "UJI-CASE-1"}})
	if kosong.StatusWork != "" {
		t.Errorf("STATUS_WORK = %q tanpa COMPLETE_DATE, mau kosong", kosong.StatusWork)
	}
	if punyaTemuan(temuan, TemuanStatusKerjaDisimpulkan) != nil {
		t.Error("penyimpulan dilaporkan padahal tidak terjadi")
	}

	tutup, temuan := LengkapiWork([]BarisLama{{
		CASEID: "UJI-CASE-1", COMPLETE_DATE: "2026-02-01",
	}})
	if tutup.StatusWork != models.StatusWorkSelesai {
		t.Errorf("STATUS_WORK = %q, mau %q", tutup.StatusWork, models.StatusWorkSelesai)
	}
	tm := punyaTemuan(temuan, TemuanStatusKerjaDisimpulkan)
	if tm == nil {
		t.Fatal("penyimpulan TIDAK dilaporkan; ia akan terbaca sebagai bacaan")
	}
	// ⛔ Dan ia menandai dirinya terbuka. Kesimpulan yang tidak menyebut
	// dirinya kesimpulan adalah kesimpulan yang tidak akan pernah dibantah.
	if !strings.Contains(tm.Catatan, "[terbuka") {
		t.Errorf("penyimpulan tidak ditandai terbuka: %q", tm.Catatan)
	}
}

func TestDokumenWarisanDigantungLewatNoAkseptasi(t *testing.T) {
	baris := []BarisLama{
		{CASEID: "CASE-1", CERTIFICATE_NO: "006", NO_ACCEPTATION: "AKS-1"},
		{CASEID: "CASE-1", CERTIFICATE_NO: "007", NO_ACCEPTATION: "AKS-2"},
	}
	dok := []DokumenLama{
		{IDPEGA: "W-CASE-1", NOAKSEP: "AKS-1", NAMAFILE: "UJI-b.pdf", KATEGORI_1: "DL-1"},
		{IDPEGA: "W-CASE-1", NOAKSEP: "AKS-1", NAMAFILE: "UJI-a.pdf", KATEGORI_1: "DL-1"},
		{IDPEGA: "W-CASE-1", NOAKSEP: "AKS-2", NAMAFILE: "UJI-c.pdf", KATEGORI_1: "DL-2"},
	}
	hasil, temuan := PetakanDokumenLama(dok, baris)
	if len(hasil["AKS-1"]) != 2 || len(hasil["AKS-2"]) != 1 {
		t.Fatalf("pemetaan = %d/%d, mau 2/1", len(hasil["AKS-1"]), len(hasil["AKS-2"]))
	}
	// ⛔ Urutan STABIL. Peta di Go tidak berurutan; laporan yang berganti
	// susunan tiap jalan tidak dapat dibandingkan dengan jalannya kemarin.
	if hasil["AKS-1"][0].NamaFile != "UJI-a.pdf" {
		t.Errorf("urutan tidak stabil: %q lebih dulu", hasil["AKS-1"][0].NamaFile)
	}
	if len(temuan) != 0 {
		t.Errorf("jalur bersih menghasilkan temuan: %v", temuan)
	}
}

func TestDokumenTanpaPesertaDilaporkanBukanDipaksa(t *testing.T) {
	// ⛔ INTI OQ-J. Saringan Pega yang sebenarnya memakai kolom `DOCUMENT`
	// pada peserta, dan kolom itu tidak ada di mana pun yang kami terima.
	// Yang tidak cocok karena itu DILAPORKAN - bukan ditempelkan ke peserta
	// mana pun supaya "tidak ada yang hilang". Dokumen milik orang lain yang
	// menempel pada peserta yang salah jauh lebih buruk daripada dokumen
	// yang dilaporkan hilang.
	baris := []BarisLama{{CASEID: "CASE-1", NO_ACCEPTATION: "AKS-1"}}
	hasil, temuan := PetakanDokumenLama([]DokumenLama{
		{IDPEGA: "W-CASE-1", NOAKSEP: "", NAMAFILE: "UJI-yatim.pdf"},
		{IDPEGA: "W-CASE-1", NOAKSEP: "AKS-TIDAK-ADA", NAMAFILE: "UJI-asing.pdf"},
	}, baris)
	if len(hasil) != 0 {
		t.Errorf("dokumen tanpa pemilik tetap dipetakan: %v", hasil)
	}
	n := 0
	for _, tm := range temuan {
		if tm.Jenis == TemuanDokumenTanpaPeserta {
			n++
		}
	}
	if n != 2 {
		t.Errorf("temuan dokumen yatim = %d, mau 2", n)
	}
}

func TestKunciKelompokAsingDilaporkanTapiTetapDibawa(t *testing.T) {
	// ⚠️ Arah selisihnya disengaja. `KATEGORI_1` adalah pembukuan Pega,
	// bukan identitas: kunci yang bentuknya aneh tetap menunjuk kelompok
	// unggahan yang nyata. Ia dilaporkan supaya terlihat, dan dibawa supaya
	// tidak hilang.
	baris := []BarisLama{{CASEID: "CASE-1", NO_ACCEPTATION: "AKS-1"}}
	hasil, temuan := PetakanDokumenLama([]DokumenLama{
		{IDPEGA: "W-CASE-1", NOAKSEP: "AKS-1", NAMAFILE: "UJI-x.pdf",
			KATEGORI_1: "LAMA-99"},
	}, baris)
	if len(hasil["AKS-1"]) != 1 {
		t.Fatalf("dokumen berkunci asing dibuang; mau tetap dibawa")
	}
	if hasil["AKS-1"][0].Kategori1 != "LAMA-99" {
		t.Errorf("kunci diperbaiki diam-diam menjadi %q", hasil["AKS-1"][0].Kategori1)
	}
	tm := punyaTemuan(temuan, TemuanKunciKelompokAsing)
	if tm == nil {
		t.Fatal("kunci kelompok asing tidak dilaporkan")
	}
	if !strings.Contains(tm.Catatan, models.AwalanKunciDokumen) {
		t.Errorf("catatan tidak menyebut bentuk yang diharap: %q", tm.Catatan)
	}
}

func TestIdpegaBertentanganDilaporkan(t *testing.T) {
	baris := []BarisLama{{CASEID: "CASE-1", NO_ACCEPTATION: "AKS-1"}}
	_, temuan := PetakanDokumenLama([]DokumenLama{
		{IDPEGA: "W-CASE-9", NOAKSEP: "AKS-1", NAMAFILE: "UJI-x.pdf", KATEGORI_1: "DL-1"},
	}, baris)
	if punyaTemuan(temuan, TemuanDokumenKlaimBerbeda) == nil {
		t.Error("pertentangan IDPEGA lawan NOAKSEP tidak dilaporkan")
	}
}

func TestDiagnosaWarisanDinyatakanTidakBersumber(t *testing.T) {
	// ⛔ Yang PALING mudah dikarang di seluruh A4, dan karena itu dijaga
	// paling keras: `.DiagnoseList` adalah daftar yang hidup di BLOB halaman
	// kerja Pega. Tabel datar hanya punya SATU pasang DISEASE/ICD_CODE.
	// Migrasi tidak membuat satu pun baris T_CLAIMLF_DIAGNOSE dari tebakan.
	tm := TemuanDiagnosaWarisan(12)
	if tm.Jenis != TemuanDiagnosaTakBersumber {
		t.Errorf("jenis = %q", tm.Jenis)
	}
	if !strings.Contains(tm.Catatan, "BLOB") {
		t.Errorf("catatan tidak menyebut di mana daftarnya hidup: %q", tm.Catatan)
	}
	if !strings.Contains(tm.Catatan, "[terbuka") {
		t.Errorf("tidak ditandai terbuka: %q", tm.Catatan)
	}
	if !strings.Contains(tm.Nilai, "12") {
		t.Errorf("cacah peserta tidak terbawa: %q", tm.Nilai)
	}
}

func TestA4TidakMenyentuhBasisData(t *testing.T) {
	// ⛔ Seluruh A4 MURNI, dan penjaga ini yang menahannya tetap begitu.
	// Fungsi migrasi data yang diam-diam membuka koneksi akan menjalankan
	// dirinya sendiri terhadap DEV saat seseorang menjalankan uji - dan
	// brief melarang eksekusi A4 sampai work owner menyetujuinya.
	b, err := os.ReadFile("migrasidokumen.go")
	if err != nil {
		t.Fatal(err)
	}
	isi := string(b)
	for _, terlarang := range []string{
		"db.sql", "QueryContext", "ExecContext", "tx.tx", "Qualify(",
		"time.Now()",
	} {
		if strings.Contains(isi, terlarang) {
			t.Errorf("migrasidokumen.go memuat %q; A4 harus murni", terlarang)
		}
	}
}
