package repository

// A4 - pelengkap migrasi data: tahap, waktu lahir, status kerja, dan dokumen
// warisan.
//
// Untuk apa berkas ini: `migrasidata.go` membongkar 55 kolom tabel datar
// menjadi pohon klaim. Yang TIDAK ada di 55 kolom itu ditangani di sini -
// termasuk hal yang ternyata tidak ada di mana pun, yang DILAPORKAN alih-alih
// dikarang.
//
// Fungsi di berkas ini MURNI: nol sentuhan basis data, nol pembacaan jam.
//
// Dibaca sesudah: migrasidata.go.

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"nusantarare/internal/models"
	"nusantarare/pkg/utils"
)

// Jenis temuan tambahan A4.
const (
	// TemuanTahapTidakAdaDiSumber - kolom tahap tidak ada di tabel datar.
	TemuanTahapTidakAdaDiSumber = "tahap tidak ada di sumber"
	// TemuanWaktuLahirDiganti - TGL_CREATE diisi kolom yang BUKAN padanannya.
	TemuanWaktuLahirDiganti = "waktu lahir diganti kolom terdekat"
	// TemuanStatusKerjaDisimpulkan - STATUS_WORK disimpulkan, bukan dibaca.
	TemuanStatusKerjaDisimpulkan = "status kerja disimpulkan dari kolom lain"
	// TemuanDokumenTanpaPeserta - dokumen warisan tidak menemukan pemiliknya.
	TemuanDokumenTanpaPeserta = "dokumen warisan tanpa peserta yang cocok"
	// TemuanKunciKelompokAsing - KATEGORI_1 tidak berbentuk DL-...
	TemuanKunciKelompokAsing = "kunci kelompok dokumen bukan DL-"
	// TemuanDokumenKlaimBerbeda - IDPEGA dan NOAKSEP menunjuk klaim berbeda.
	TemuanDokumenKlaimBerbeda = "IDPEGA dan NOAKSEP menunjuk klaim berbeda"
	// TemuanDiagnosaTakBersumber - diagnosa warisan tidak punya tabel sumber.
	TemuanDiagnosaTakBersumber = "diagnosa warisan tidak punya tabel sumber"
)

// DokumenLama adalah satu baris tabel warisan POOLDATA.DOCUMENT_CLAIM.
//
// `[data DBA - katalog DEV 26-09-2026]` 14 kolom milik kelas Pega
// `ASM-FW-GCNMFW-Int-DOCUMENT_CLAIM`, 295 baris. Yang dipakai migrasi ini
// tujuh; sisanya dibaca tetapi tidak dipindahkan.
//
// ⛔ `BASE64` TIDAK ada di sini, dan itu disengaja: isi berkas tidak pernah
// masuk artefak mana pun. Migrasi memindahkan CATATANnya; berkasnya sendiri
// sudah ada di penyimpanan dan dirujuk `T_STORAGE_ID`.
type DokumenLama struct {
	// IDPEGA adalah `pyWorkPage.pzInsKey` saat unggah
	// (`InsertDocument_Act.xml` b185, dipanggil `SaveAttachLife.xml` b1467).
	// Ia menunjuk KLAIM, bukan peserta.
	IDPEGA string
	// NOAKSEP adalah satu-satunya tautan ke PESERTA yang tabel datar punya:
	// ia sepadan dengan `OS_AKSEPTASI_KLAIM_LIFE.NO_ACCEPTATION`.
	NOAKSEP    string
	NAMAFILE   string
	MIME       string
	KATEGORI_1 string
	KATEGORI_2 string
	TStorageID string
}

// PelengkapWork adalah nilai yang A4 isikan ke baris work.
//
// ⚠️ Dipisah dari `models.WorkClaim` supaya yang DISIMPULKAN tidak bercampur
// dengan yang DIBACA. Pemanggil menempelkannya sendiri, dan laporan Temuan
// menyebut mana yang mana.
type PelengkapWork struct {
	Tahap      string
	TglCreate  time.Time
	StatusWork string
}

// LengkapiWork mengisi TAHAP, TGL_CREATE, dan STATUS_WORK baris lama.
//
// ⛔ KETIGANYA TIDAK ADA SEBAGAI KOLOM di 55 kolom tabel datar - daftar itu
// dibaca VERBATIM dari `UpdateOsAkseptasiClaimLife_sql.xml` dan tidak memuat
// satu pun kolom tahap, waktu lahir kasus, atau status kerja. Karena itu:
//
//	TAHAP        dibiarkan KOSONG, dan dilaporkan. Tahap kosong tidak
//	             menawarkan tombol apa pun (gagal TERTUTUP) - perilaku yang
//	             benar untuk kasus yang tahapnya memang tidak diketahui.
//	TGL_CREATE   diisi CLAIM_RECEIVED_DATE, kolom terdekat yang SUNGGUH ada,
//	             dan penggantian itu dilaporkan. Ia BUKAN pxCreateDateTime.
//	STATUS_WORK  Resolved-Completed HANYA bila COMPLETE_DATE terisi; selain
//	             itu kosong (= belum ditutup, migrasi 017).
//
// ⚠️ `[terbuka - work owner]` Penyimpulan STATUS_WORK dari COMPLETE_DATE
// adalah KESIMPULAN, bukan bacaan. Ia dilaporkan setiap kali dipakai supaya
// work owner dapat membantahnya sebelum migrasi dijalankan. Alternatifnya -
// mengosongkan seluruhnya - membuat setiap kasus warisan yang sudah selesai
// tampak masih berjalan di kotak masuk, dan itu selisih yang lebih besar.
func LengkapiWork(rows []BarisLama) (PelengkapWork, []Temuan) {
	var out PelengkapWork
	var temuan []Temuan
	if len(rows) == 0 {
		return out, nil
	}
	sumber := rows[0].CASEID

	temuan = append(temuan, Temuan{
		Jenis: TemuanTahapTidakAdaDiSumber, Sumber: sumber, Medan: "TAHAP",
		Catatan: "nol kolom tahap di 55 kolom OS_AKSEPTASI_KLAIM_LIFE " +
			"(UpdateOsAkseptasiClaimLife_sql.xml); dibiarkan kosong, dan tahap " +
			"kosong tidak menawarkan tombol apa pun",
	})

	if teks := strings.TrimSpace(rows[0].CLAIM_RECEIVED_DATE); teks != "" {
		t, err := utils.ParseTanggal(teks)
		if err != nil {
			temuan = append(temuan, Temuan{
				Jenis: TemuanTanggalTakTerurai, Sumber: rows[0].ID,
				Medan: "CLAIM_RECEIVED_DATE", Nilai: teks,
				Catatan: fmt.Sprintf("%v; TGL_CREATE dibiarkan kosong", err),
			})
		} else {
			out.TglCreate = t
			temuan = append(temuan, Temuan{
				Jenis: TemuanWaktuLahirDiganti, Sumber: sumber, Medan: "TGL_CREATE",
				Nilai: teks,
				Catatan: "diisi CLAIM_RECEIVED_DATE - kolom terdekat yang ADA. " +
					"Ia BUKAN pxCreateDateTime: tanggal klaim DITERIMA tidak sama " +
					"dengan waktu kasusnya lahir di Pega",
			})
		}
	}

	if teks := strings.TrimSpace(rows[0].COMPLETE_DATE); teks != "" {
		out.StatusWork = models.StatusWorkSelesai
		temuan = append(temuan, Temuan{
			Jenis: TemuanStatusKerjaDisimpulkan, Sumber: sumber, Medan: "STATUS_WORK",
			Nilai: teks,
			Catatan: "COMPLETE_DATE terisi -> STATUS_WORK = " + models.StatusWorkSelesai +
				". KESIMPULAN, bukan bacaan: nol kolom status kerja di sumber. " +
				"[terbuka - work owner]",
		})
	}
	return out, temuan
}

// PetakanDokumenLama menggantung dokumen warisan ke pesertanya.
//
// ⛔ SATU-SATUNYA tautan ke peserta yang sumbernya punya adalah `NOAKSEP`
// lawan `NO_ACCEPTATION`. `IDPEGA` menunjuk KLAIM saja
// (`pyWorkPage.pzInsKey`, `InsertDocument_Act.xml` b185), dan `KATEGORI_1`
// adalah kunci KELOMPOK unggahan (`DL-` + angka, `SaveAttachLife.xml` b595) -
// bukan identitas peserta.
//
// ⛔ OQ-J: saringan Pega yang sebenarnya `.KATEGORI_1 = .DOCUMENT` pada
// halaman PESERTA (`LoadDocumentLife_ACT.xml` b495), dan kolom `DOCUMENT`
// tidak ada di mana pun yang kami terima. Karena itu pemetaan ini TIDAK dapat
// meniru saringan aslinya, dan setiap ketidakcocokan DILAPORKAN - tidak
// dipaksa menempel ke peserta mana pun.
//
// Mengembalikan peta `NO_ACCEPTATION -> dokumen`, sebab itulah kunci yang
// pemanggil punya sebelum pengenal peserta baru terbit.
func PetakanDokumenLama(dok []DokumenLama, baris []BarisLama) (
	map[string][]models.Dokumen, []Temuan) {

	klaimDari := map[string]string{}
	for _, b := range baris {
		no := strings.TrimSpace(b.NO_ACCEPTATION)
		if no == "" {
			continue
		}
		klaimDari[no] = b.CASEID
	}

	hasil := map[string][]models.Dokumen{}
	var temuan []Temuan
	for _, d := range dok {
		no := strings.TrimSpace(d.NOAKSEP)
		if no == "" {
			temuan = append(temuan, Temuan{
				Jenis: TemuanDokumenTanpaPeserta, Sumber: d.IDPEGA, Medan: "NOAKSEP",
				Nilai: d.NAMAFILE,
				Catatan: "NOAKSEP kosong; satu-satunya tautan ke peserta hilang, " +
					"dan IDPEGA hanya menunjuk klaim. TIDAK dipindahkan",
			})
			continue
		}
		caseID, cocok := klaimDari[no]
		if !cocok {
			temuan = append(temuan, Temuan{
				Jenis: TemuanDokumenTanpaPeserta, Sumber: d.IDPEGA, Medan: "NOAKSEP",
				Nilai: no,
				Catatan: "nol baris OS_AKSEPTASI_KLAIM_LIFE bernomor akseptasi ini; " +
					"TIDAK dipindahkan, dan TIDAK ditempelkan ke peserta mana pun",
			})
			continue
		}
		// ⚠️ `IDPEGA` dibandingkan hanya bila terisi. Ia menunjuk work
		// object, bukan CASEID, jadi yang diperiksa keberadaan pertentangan
		// yang JELAS - bukan kesamaan bentuk.
		if id := strings.TrimSpace(d.IDPEGA); id != "" && !strings.Contains(id, caseID) {
			temuan = append(temuan, Temuan{
				Jenis: TemuanDokumenKlaimBerbeda, Sumber: d.IDPEGA, Medan: "IDPEGA",
				Nilai: no,
				Catatan: fmt.Sprintf("NOAKSEP menunjuk klaim %q sedangkan IDPEGA "+
					"tidak memuatnya; dokumen tetap dipindahkan menurut NOAKSEP, "+
					"dan selisihnya dilaporkan", caseID),
			})
		}
		if k := strings.TrimSpace(d.KATEGORI_1); k != "" &&
			!strings.HasPrefix(k, models.AwalanKunciDokumen) {
			temuan = append(temuan, Temuan{
				Jenis: TemuanKunciKelompokAsing, Sumber: d.IDPEGA, Medan: "KATEGORI_1",
				Nilai: k,
				Catatan: "kunci kelompok tidak berawalan " + models.AwalanKunciDokumen +
					" (SaveAttachLife.xml b595); dibawa apa adanya, tidak diperbaiki",
			})
		}
		hasil[no] = append(hasil[no], models.Dokumen{
			NamaFile:   d.NAMAFILE,
			Mime:       d.MIME,
			Kategori1:  d.KATEGORI_1,
			Kategori2:  d.KATEGORI_2,
			TStorageID: d.TStorageID,
		})
	}
	// Urutan stabil: peta di Go tidak berurutan, dan laporan yang berganti
	// susunan tiap jalan tidak dapat dibandingkan dengan jalannya kemarin.
	for no := range hasil {
		sort.SliceStable(hasil[no], func(i, j int) bool {
			return hasil[no][i].NamaFile < hasil[no][j].NamaFile
		})
	}
	return hasil, temuan
}

// TemuanDiagnosaWarisan menyatakan bahwa diagnosa lama tidak punya sumber.
//
// ⛔ BUKAN kelalaian, dan bukan hal yang dapat ditutup dengan mencari lebih
// lama. Tabel datar `OS_AKSEPTASI_KLAIM_LIFE` punya `DISEASE` dan `ICD_CODE`
// TUNGGAL - satu pasang per baris peserta - sedangkan `.DiagnoseList`
// (butir bd) adalah DAFTAR. Daftar itu hidup di halaman kerja Pega, yang
// tersimpan sebagai BLOB milik mesin Pega dan TIDAK diekspor sebagai tabel.
//
// Akibatnya migrasi hanya dapat membawa SATU diagnosa per peserta - yang ada
// di kolom datar - dan daftar yang lebih panjang dari satu HILANG. Itu
// dinyatakan di sini, dilaporkan setiap jalan, dan menjadi butir yang work
// owner harus putuskan sebelum cutover.
//
// ⚠️ Yang TIDAK dilakukan: mengarang. Nol baris `T_CLAIMLF_DIAGNOSE` dibuat
// dari tebakan.
func TemuanDiagnosaWarisan(cacahPeserta int) Temuan {
	return Temuan{
		Jenis: TemuanDiagnosaTakBersumber,
		Medan: "T_CLAIMLF_DIAGNOSE",
		Nilai: fmt.Sprintf("%d peserta", cacahPeserta),
		Catatan: "OS_AKSEPTASI_KLAIM_LIFE hanya punya DISEASE dan ICD_CODE TUNGGAL; " +
			".DiagnoseList adalah DAFTAR yang hidup di BLOB halaman kerja Pega dan " +
			"tidak diekspor sebagai tabel. Migrasi membawa paling banyak SATU " +
			"diagnosa per peserta; daftar yang lebih panjang tidak dapat dipulihkan. " +
			"[terbuka - work owner]",
	}
}
