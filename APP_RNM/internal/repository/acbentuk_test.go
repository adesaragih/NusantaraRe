package repository

// Test per acceptance criterion tiket 14 - TANPA Oracle.
//
// Untuk apa berkas ini: ronde 1 menulis kolom-kolomnya ke berkas migrasi tetapi
// tidak menguncinya dengan satu pun pernyataan test, sehingga empat belas AC
// berstatus "tertulis, belum teruji". Satu fungsi di sini = satu AC, dinamai
// dengan nomornya, supaya kaitannya tidak perlu ditebak siapa pun.
//
// Dibaca sesudah: migrasi_test.go dan strukturkolom_test.go.
//
// ⚠️ Beberapa AC mengeja nama kolom berbeda dari STRUKTUR-TABEL-CLAIM-LIFE.md.
// Yang dipakai adalah ejaan STRUKTUR, dan selisihnya disebut di komentar AC
// yang bersangkutan - bukan didiamkan.

import (
	"regexp"
	"strings"
	"testing"

	"nusantarare/inti"
)

// sqlTabel mengembalikan teks SQL berkas migrasi yang berawalan nomor tertentu.
func sqlTabel(t *testing.T, awalan string) string {
	t.Helper()
	for nama, teks := range seluruhSQL(t, false) {
		if strings.HasPrefix(nama, awalan) {
			return strings.ToUpper(teks)
		}
	}
	t.Fatalf("berkas migrasi berawalan %s tidak ketemu", awalan)
	return ""
}

// wajibMemuat memeriksa bahwa teks SQL menyebut seluruh nama yang diminta.
func wajibMemuat(t *testing.T, sql, tabel string, nama []string) {
	t.Helper()
	if len(nama) == 0 {
		t.Fatal("daftar kolom kosong; test-nya yang rusak")
	}
	for _, n := range nama {
		if !strings.Contains(sql, n) {
			t.Errorf("%s: kolom %s tidak ada", tabel, n)
		}
	}
}

// AC 5: kunci tamu T_CLAIMLF_DOCUMENT menunjuk PESERTA, bukan header klaim.
func TestAC05DokumenMenunjukPeserta(t *testing.T) {
	sql := sqlTabel(t, "007_")
	if !strings.Contains(sql, "REFERENCES {SKEMA}.T_CLAIMLF_PREMIUMLIST_DETAIL (ID)") {
		t.Error("FK T_CLAIMLF_DOCUMENT tidak menunjuk T_CLAIMLF_PREMIUMLIST_DETAIL (ID)")
	}
	if strings.Contains(sql, "REFERENCES {SKEMA}.T_GENERAL_CLAIM") {
		t.Error("T_CLAIMLF_DOCUMENT menggantung pada header klaim - seharusnya pada peserta")
	}
}

// AC 7: header memuat keempat field PremiumListSummary beserta penunjuk polis.
//
// ⚠️ Tiga nama dieja berbeda di AC dan di STRUKTUR, dan STRUKTUR yang dipakai:
// PL_NUMBER menjadi POLICY_NO (STRUKTUR menyebutnya ganti nama secara
// eksplisit), CASEID menjadi CASEID_POLICY (CASEID polos pindah ke
// T_WORK_CLAIM.CASE_ID), dan RISLIPRNM menjadi RI_SLIP_RNM.
func TestAC07HeaderMemuatFieldSummary(t *testing.T) {
	wajibMemuat(t, sqlTabel(t, "002_"), "T_GENERAL_CLAIM", []string{
		"CLAIM_NO", "POLICY_NO", "RI_SLIP_RNM", "BUSINESS_NAME",
		"CASEID_POLICY", "CLAIM_RETRO",
	})
}

// AC 9: peserta memuat kesembilan kolom tambahan hasil audit.
func TestAC09PesertaMemuatSembilanKolomAudit(t *testing.T) {
	wajibMemuat(t, sqlTabel(t, "003_"), "T_CLAIMLF_PREMIUMLIST_DETAIL", []string{
		"IS_CHECK", "STATUS", "RECOMMENDATION", "STS_REJECT", "SOURCE_ID",
		"CONFIRMATION_DATE", "COMPLETE_DATE", "CLAIM_RECEIVED_DATE",
		"CEDING_RETENTION",
	})
}

// AC 10: adjustment memuat keempat kolom intinya.
//
// ⚠️ AC mengeja ACCEPTEDNO; STRUKTUR mengeja ACCEPTED_NO, dan seluruh kolom
// baru proyek ini memakai garis bawah. Ejaan STRUKTUR yang dipakai.
func TestAC10AdjustmentMemuatKolomInti(t *testing.T) {
	wajibMemuat(t, sqlTabel(t, "004_"), "T_CLAIMLF_ADJUSTMENT", []string{
		"CLAIM_AMOUNT", "STS_REJECT", "ACCEPTED_NO", "ACCEPTATION_DATE",
	})
}

// AC 23: T_WORK_CLAIM memuat kolom pindahan dari header.
//
// ⚠️ CASEID di AC menjadi CASE_ID di STRUKTUR.
func TestAC23WorkClaimMemuatKolomPindahan(t *testing.T) {
	wajibMemuat(t, sqlTabel(t, "001_"), "T_WORK_CLAIM", []string{
		"CREATE_OP", "CREATE_OP_NAME", "TGL_UPDATE", "CASE_ID", "LINI",
	})
}

// AC 28: kolom LINI ada di T_WORK_CLAIM, dan untuk Life isinya konstanta.
//
// Yang `[terbuka]` hanya DAFTAR NILAI enum lintas-lini, bukan keberadaan
// kolomnya - karena itu constraint CHECK sengaja tidak dipasang.
func TestAC28KolomLiniAdaDanKonstantanyaTunggal(t *testing.T) {
	wajibMemuat(t, sqlTabel(t, "001_"), "T_WORK_CLAIM", []string{"LINI"})
	if inti.LiniLife != "LIFE" {
		t.Errorf("konstanta lini Life = %q, mau %q", inti.LiniLife, "LIFE")
	}
}

// AC 33: keempat kolom identitas bertipe SAMA dengan T_WORK_CLAIM.ID.
//
// ⚠️ T_GENERAL_KOMITE.ID tidak diperiksa di sini: tabelnya milik konteks
// Komite dan TIDAK dibuat tiket ini. Bagian AC itu tetap terbuka.
func TestAC33IdentitasBertipeSama(t *testing.T) {
	const tipe = "VARCHAR2(32)"
	kasus := []struct{ awalan, kolom string }{
		{"001_", "ID"},
		{"001_", "COVER_KEY"},
		{"002_", "ID"},
		{"004_", "KOMITE_ID"},
	}
	for _, k := range kasus {
		sql := sqlTabel(t, k.awalan)
		var ketemu bool
		for _, baris := range strings.Split(sql, "\n") {
			b := strings.TrimSpace(baris)
			if strings.HasPrefix(b, k.kolom+" ") && strings.Contains(b, tipe) {
				ketemu = true
			}
		}
		if !ketemu {
			t.Errorf("%s kolom %s tidak bertipe %s", k.awalan, k.kolom, tipe)
		}
	}
}

// AC 39: spreading memuat kedelapan kolomnya.
func TestAC39SpreadingMemuatKolomnya(t *testing.T) {
	wajibMemuat(t, sqlTabel(t, "005_"), "T_CLAIMLF_ADJUSTMENT_SPREADING", []string{
		"TREATY_TYPE_ID", "TREATY_TYPE_NAME", "TREATY_YEAR_LIFE",
		"RETROCADED_SHARE", "RATE", "IDR", "USD", "CURRENCY",
	})
}

// AC 40: spreading retro memuat kesepuluh kolomnya.
//
// ⚠️ Ejaan COMMISION kurang satu huruf S, dan itu memang ejaan korpus - 2.691
// kemunculan di 369 berkas. Ia disalin apa adanya supaya tidak lahir dua nama
// untuk satu kolom.
func TestAC40SpreadingRetroMemuatKolomnya(t *testing.T) {
	wajibMemuat(t, sqlTabel(t, "006_"), "T_CLAIMLF_ADJ_SPREADING_RETRO", []string{
		"REINSURER_NAME", "PERCENT_SHARE", "AMOUNT", "RATE",
		"PREMIUM_SPREADED_GROSS", "PREMIUM_SPREADED_NET", "COMMISION",
		"OVR_COMM", "TREATY_TYPE_ID", "TREATY_TYPE_NAME",
	})
}

// AC 41: seluruh kolom uang dan persen pada kedua tabel spreading bertipe
// desimal - nol float, nol teks.
func TestAC41UangDanPersenBertipeDesimal(t *testing.T) {
	kolom := map[string][]string{
		"005_": {"RETROCADED_SHARE", "RATE", "IDR", "USD"},
		"006_": {"PERCENT_SHARE", "AMOUNT", "RATE", "PREMIUM_SPREADED_GROSS",
			"PREMIUM_SPREADED_NET", "COMMISION", "OVR_COMM"},
	}
	for awalan, daftar := range kolom {
		sql := sqlTabel(t, awalan)
		for _, n := range daftar {
			var tipe string
			for _, baris := range strings.Split(sql, "\n") {
				b := strings.TrimSpace(baris)
				if strings.HasPrefix(b, n+" ") {
					tipe = b
				}
			}
			if tipe == "" {
				t.Errorf("%s kolom %s tidak ketemu", awalan, n)
				continue
			}
			// Diperketat ronde 3: "memuat kata NUMBER" terlalu longgar -
			// NUMBER polos lolos, dan justru itu yang terjadi sampai ronde 2.
			// Keputusan work owner c menetapkan NUMBER(38,8).
			if !strings.Contains(tipe, "NUMBER(38,8)") {
				t.Errorf("%s kolom %s bukan NUMBER(38,8): %q", awalan, n, tipe)
			}
			for _, terlarang := range []string{"FLOAT", "BINARY_DOUBLE", "BINARY_FLOAT", "VARCHAR"} {
				if strings.Contains(tipe, terlarang) {
					t.Errorf("%s kolom %s bertipe %s - uang tidak pernah begitu (ADR-U-0003)",
						awalan, n, terlarang)
				}
			}
		}
	}
}

// AC 42: nilai spreading DIBEKUKAN - ia kolom tersimpan, bukan turunan.
//
// Yang dapat dijaga dari sisi skema adalah ini: tabel spreading tidak dibuat
// sebagai VIEW, dan migrasi tidak menautkannya ke master treaty mana pun.
// Kalau angkanya dihitung ulang dari master, perubahan master akan mengubah
// angka yang sudah tersimpan - persis yang AC ini larang.
func TestAC42SpreadingDibekukanBukanTurunan(t *testing.T) {
	for _, awalan := range []string{"005_", "006_"} {
		sql := sqlTabel(t, awalan)
		if strings.Contains(sql, "CREATE VIEW") || strings.Contains(sql, "CREATE OR REPLACE VIEW") {
			t.Errorf("%s dibuat sebagai VIEW - nilainya tidak beku", awalan)
		}
		if strings.Contains(sql, "GENERATED ALWAYS AS") {
			t.Errorf("%s memuat kolom turunan - nilainya tidak beku", awalan)
		}
		for _, master := range []string{"M_TREATY", "MASTER_TREATY", "M_CONTRACT_RETRO"} {
			if strings.Contains(sql, master) {
				t.Errorf("%s menautkan diri ke master %s - nilainya tidak beku", awalan, master)
			}
		}
	}
}

// AC 45 dan 48: rujukan ke Komite memakai KOMITE_ID, bukan indeks posisi.
//
// ⛔ Test yang menemukan padanan indeks posisi sebagai kunci rujukan GAGAL.
// Pega memakai .pxListSubscript; memindahkannya apa adanya akan membuat rujukan
// rusak begitu satu baris disisipkan di tengah.
func TestAC45Dan48RujukanKomiteBukanIndeksPosisi(t *testing.T) {
	sql := sqlTabel(t, "004_")
	if !strings.Contains(sql, "KOMITE_ID") {
		t.Fatal("T_CLAIMLF_ADJUSTMENT tidak memuat KOMITE_ID")
	}
	for _, terlarang := range []string{
		"INDEXPREMIUMLIST", "INDEXADJUSTMENT", "INDEX_PREMIUM_LIST",
		"INDEX_ADJUSTMENT", "PXLISTSUBSCRIPT", "LIST_SUBSCRIPT",
	} {
		if strings.Contains(sql, terlarang) {
			t.Errorf("kunci rujukan komite memakai indeks posisi %q", terlarang)
		}
	}
	// AC 48: KOMITE_ID ber-index. Diperiksa atas pernyataan CREATE INDEX-nya
	// sendiri, bukan atas keberadaan kata CREATE di mana pun - berkas ini selalu
	// memuat CREATE, sehingga pemeriksaan begitu tidak pernah gagal.
	satuBaris := strings.Join(strings.Fields(sql), " ")
	var berIndex bool
	for _, potong := range strings.Split(satuBaris, "CREATE ") {
		if !strings.HasPrefix(potong, "INDEX") && !strings.HasPrefix(potong, "UNIQUE INDEX") {
			continue
		}
		if strings.Contains(potong, "T_CLAIMLF_ADJUSTMENT (KOMITE_ID)") {
			berIndex = true
		}
	}
	if !berIndex {
		t.Error("tidak ada CREATE INDEX atas T_CLAIMLF_ADJUSTMENT (KOMITE_ID)")
	}
}

// AC 50: penunjuk KE ATAS dipasangi REFERENCES T_WORK_CLAIM(ID).
//
// `[keputusan work owner 26-09-2026, butir d]` sesudah empat sesi berstatus
// [USULAN]. Kedua kolomnya tetap nullable; yang dilarang hanyalah menunjuk
// baris yang tidak ada.
//
// Test ini juga mengunci bahwa keduanya TANPA "ON DELETE": kaskade di penunjuk
// ke atas akan membuat penghapusan satu baris komite ikut menghapus baris
// adjustment yang menunjuknya - kebalikan dari arah kepemilikan.
func TestAC50PenunjukKeAtasBerReferences(t *testing.T) {
	kasus := []struct {
		awalan, kolom, constraint string
	}{
		{"001_", "COVER_KEY", "FK_WORK_COVER_KEY"},
		{"004_", "KOMITE_ID", "FK_ADJ_KOMITE"},
	}
	for _, k := range kasus {
		sql := sqlTabel(t, k.awalan)
		if !strings.Contains(sql, "CONSTRAINT "+k.constraint+" FOREIGN KEY ("+k.kolom+")") {
			t.Errorf("%s: %s tidak dipasangi FOREIGN KEY", k.awalan, k.kolom)
			continue
		}
		// Hanya klausa milik constraint INI yang diperiksa. Jendela sekian
		// byte akan menembus pernyataan berikutnya - seluruh pernyataan satu
		// berkas disambung menjadi satu teks - sehingga ON DELETE milik
		// constraint lain akan dituduhkan ke constraint ini.
		klausa := klausaConstraint(sql, k.constraint)
		if !strings.Contains(klausa, "REFERENCES {SKEMA}.T_WORK_CLAIM (ID)") {
			t.Errorf("%s: %s tidak menunjuk T_WORK_CLAIM (ID); klausanya: %q",
				k.awalan, k.kolom, klausa)
		}
		if strings.Contains(klausa, "ON DELETE") {
			t.Errorf("%s: %s memakai ON DELETE; penunjuk ke atas bukan kepemilikan",
				k.awalan, k.kolom)
		}
		// Nullable: kolomnya tidak boleh dideklarasikan NOT NULL.
		//
		// DDL meratakan kolom dengan BANYAK spasi, jadi teksnya dirapatkan
		// dulu. Tanpa itu pola berspasi tunggal tidak pernah cocok dan
		// penjaganya mati - persis yang terjadi sampai tinjauan ronde 4.
		if strings.Contains(rapatkanSpasi(sql), k.kolom+" VARCHAR2(32) NOT NULL") {
			t.Errorf("%s: %s NOT NULL; keputusan d menuntutnya tetap nullable", k.awalan, k.kolom)
		}
	}
}

// spasiBeruntun dipakai merapatkan perataan kolom pada DDL.
var spasiBeruntun = regexp.MustCompile(`[ \t]+`)

// rapatkanSpasi mengubah setiap deretan spasi menjadi satu spasi.
func rapatkanSpasi(s string) string { return spasiBeruntun.ReplaceAllString(s, " ") }

// klausaConstraint memotong teks SATU constraint, dari sesudah namanya sampai
// tepat sebelum constraint berikutnya atau akhir daftar kolom.
//
// Kenapa perlu: seluruh pernyataan sebuah berkas migrasi disambung menjadi satu
// teks, sehingga memotong "sekian byte sesudah nama constraint" akan ikut
// menelan pernyataan CREATE INDEX di bawahnya dan menuduhkan isinya ke sini.
func klausaConstraint(sql, nama string) string {
	awal := strings.Index(sql, "CONSTRAINT "+nama)
	if awal < 0 {
		return ""
	}
	sisa := sql[awal+len("CONSTRAINT "+nama):]
	akhir := len(sisa)
	for _, batas := range []string{"CONSTRAINT ", "\n)"} {
		if i := strings.Index(sisa, batas); i >= 0 && i < akhir {
			akhir = i
		}
	}
	return rapatkanSpasi(sisa[:akhir])
}
