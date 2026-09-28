package repository

// Penjaga bentuk migrasi 300-319 Treaty Contract Out - TANPA Oracle (tiket 01).
//
// Yang dijaga di sini adalah keputusan yang tercatat: awalan T_ (tco1),
// tujuh tabel, uang dan persen desimal, tanggal DATE, kunci gabungan TANPA
// FK ke kontrak, satu-satunya FK berkaskade pada security, nol tabel dokumen
// M_*, dan tabel warisan yang HANYA dibaca.

import (
	"regexp"
	"strings"
	"testing"
)

// milikTreatyContractOut menjawab apakah berkas migrasi itu milik modul ini.
//
// Rentangnya 300-319 (PROMPT-EKSEKUSI-HULU-HILIR.md §4).
func milikTreatyContractOut(nama string) bool {
	return nama >= "300_" && nama < "320_"
}

func sqlTreatyContractOut(t *testing.T) map[string]string {
	t.Helper()
	out := map[string]string{}
	for nama, teks := range seluruhSQL(t, false) {
		if milikTreatyContractOut(nama) {
			out[nama] = strings.ToUpper(teks)
		}
	}
	if len(out) == 0 {
		t.Fatal("nol migrasi 300-319 terbaca; pembacanya yang rusak")
	}
	return out
}

func TestTCOSeluruhBerkasTigaRatusanDiRentangModul(t *testing.T) {
	for nama := range seluruhSQL(t, false) {
		if strings.HasPrefix(nama, "3") && !milikTreatyContractOut(nama) {
			t.Errorf("%s di luar rentang 300-319 Treaty Contract Out", nama)
		}
	}
}

// tco1: tujuh tabel, seluruhnya berawalan T_, nama warisan TIDAK dibuat.
func TestTCOTujuhTabelBerawalanT(t *testing.T) {
	mau := []string{TabelTahunTCO, TabelKontrakTCO, TabelReinsurerTCO,
		TabelSecurityTCO, TabelBusinessTCO, TabelKlausulTCO, TabelJejakTCO}
	var gabung strings.Builder
	for _, teks := range sqlTreatyContractOut(t) {
		gabung.WriteString(teks)
	}
	sql := gabung.String()
	for _, tb := range mau {
		if !strings.Contains(sql, "CREATE TABLE {SKEMA}."+tb+" (") {
			t.Errorf("tabel %s tidak dibuat", tb)
		}
	}
	pola := regexp.MustCompile(`CREATE TABLE \{SKEMA\}\.(\w+)`)
	dibuat := pola.FindAllStringSubmatch(sql, -1)
	if len(dibuat) != len(mau) {
		t.Errorf("CREATE TABLE di 300-319: %d, mau %d", len(dibuat), len(mau))
	}
	for _, m := range dibuat {
		if !strings.HasPrefix(m[1], "T_") {
			t.Errorf("tabel %s tanpa awalan T_ (keputusan tco1)", m[1])
		}
		for _, warisan := range urutanWarisanTCO {
			if m[1] == warisan {
				t.Errorf("nama warisan %s dibuat ulang - tabel warisan sudah ada di POOLDATA", warisan)
			}
		}
	}
}

// AC 51/52/53: tipe kolom uang, persen, dan tanggal di tabel T_*.
func TestTCOUangDesimalTanggalDATE(t *testing.T) {
	tipe := tipeMenurutDDL(t)
	for _, warisan := range urutanWarisanTCO {
		tabelBaru := tabelBaruDariWarisanTCO[warisan]
		kolomDDL, ada := tipe[tabelBaru]
		if !ada {
			t.Fatalf("%s tidak terbaca dari DDL", tabelBaru)
		}
		for _, k := range KolomWarisanTCO(warisan) {
			dapat, punya := kolomDDL[k]
			if !punya {
				t.Errorf("%s: kolom warisan %s tidak ada di DDL - nama harus VERBATIM", tabelBaru, k)
				continue
			}
			switch TipeTujuanTCO(warisan, k) {
			case TujuanDesimal:
				if dapat != "NUMBER(38,8)" {
					t.Errorf("%s.%s bertipe %s, mau NUMBER(38,8) (ADR-0003)", tabelBaru, k, dapat)
				}
			case TujuanTanggal:
				if dapat != "DATE" {
					t.Errorf("%s.%s bertipe %s, mau DATE (AC 53)", tabelBaru, k, dapat)
				}
			default:
				if strings.HasPrefix(dapat, "NUMBER") || dapat == "DATE" {
					t.Errorf("%s.%s bertipe %s padahal STRUKTUR menyebut teks", tabelBaru, k, dapat)
				}
			}
		}
	}
	// Daftar uang AC 51 dan persen AC 52 dinamai satu per satu.
	uang := map[string][]string{
		TabelKlausulTCO:   {"RP", "USD", "MORERP", "MOREUSD", "TREATYLIMIT", "COINS_MIN", "COINS_MAX", "PCT", "PCTME", "KURS"},
		TabelSecurityTCO:  {"PCT_SHARE"},
		TabelReinsurerTCO: {"RICOMM", "PCTSHARE"},
	}
	for tabel, kolom := range uang {
		for _, k := range kolom {
			if tipe[tabel][k] != "NUMBER(38,8)" {
				t.Errorf("%s.%s = %q, mau NUMBER(38,8)", tabel, k, tipe[tabel][k])
			}
		}
	}
	// Security PK surrogate: ID VARCHAR2(32) NOT NULL (AC 68).
	if tipe[TabelSecurityTCO]["ID"] != "VARCHAR2(32)" {
		t.Errorf("T_MTREATYSECURITY.ID = %q, mau VARCHAR2(32)", tipe[TabelSecurityTCO]["ID"])
	}
}

// Kebijakan FK modul ini - berbeda dari Claim Life dan PremiumList, dan
// disengaja (spec §2): anak menggantung pada kunci gabungan.
func TestTCOKebijakanKunciTamu(t *testing.T) {
	berkas := sqlTreatyContractOut(t)
	tanpaFK := []string{"300_", "302_", "304_", "305_", "306_"}
	for _, awalan := range tanpaFK {
		for nama, isi := range berkas {
			if strings.HasPrefix(nama, awalan) && strings.Contains(isi, "REFERENCES") {
				t.Errorf("%s memuat REFERENCES - reinsurer, business, klausul menggantung pada "+
					"KUNCI GABUNGAN, bukan FK (fakta bisnis work owner)", nama)
			}
		}
	}
	kontrak := cariMigrasi(t, berkas, "301_")
	if !strings.Contains(kontrak, "REFERENCES {SKEMA}.T_TREATYYEAR (ID)") {
		t.Error("T_TREATYCONTRACT tidak menunjuk T_TREATYYEAR")
	}
	if strings.Contains(kontrak, "ON DELETE CASCADE") {
		t.Error("kontrak berkaskade dari tahun - tidak ada jalur hapus tahun; harus ditolak terang (ORA-02292)")
	}
	security := cariMigrasi(t, berkas, "303_")
	if !strings.Contains(security, "REFERENCES {SKEMA}.T_TREATYREINSURER (ID) ON DELETE CASCADE") {
		t.Error("T_MTREATYSECURITY harus menunjuk T_TREATYREINSURER dengan ON DELETE CASCADE (tiket 06)")
	}
	// Setiap FK ber-index pada kolomnya.
	polaFK := regexp.MustCompile(`FOREIGN KEY \(([A-Z_]+)\)`)
	diperiksa := 0
	for nama, isi := range berkas {
		for _, m := range polaFK.FindAllStringSubmatch(isi, -1) {
			diperiksa++
			if !strings.Contains(isi, "("+m[1]+")\n") {
				t.Errorf("%s: FK %s tanpa index pada kolomnya", nama, m[1])
			}
		}
	}
	if diperiksa != 2 {
		t.Errorf("FK di 300-319: %d, mau tepat 2 (kontrak->tahun, security->reinsurer)", diperiksa)
	}
}

func cariMigrasi(t *testing.T, berkas map[string]string, awalan string) string {
	t.Helper()
	for nama, isi := range berkas {
		if strings.HasPrefix(nama, awalan) {
			return isi
		}
	}
	t.Fatalf("migrasi berawalan %s tidak ada", awalan)
	return ""
}

// AC 63/64: nol tabel dokumen di skema dan nol pembacaan tabel M_* warisan
// di sumber produksi mana pun.
func TestTCONolTabelDokumenWarisan(t *testing.T) {
	terlarang := []string{"M_PROPORTIONALARRG", "M_TREATYCONTRACT", "M_TREATYYEAR",
		"M_TREATYBUSINESS", "PROSESCOPY", "JSON_KLAIM"}
	diperiksa := 0
	for nama, isi := range berkasSumberProduksi(t) {
		diperiksa++
		kode := strings.ToUpper(buangKomentarSumber(nama, isi))
		for _, n := range terlarang {
			if strings.Contains(kode, n) {
				t.Errorf("%s menyebut %s - tabel dokumen MATI / fitur salin dibuang / teks galat "+
					"salin-tempel tidak dibawa (spec penyimpangan sadar 1, 3, 8)", nama, n)
			}
		}
	}
	if diperiksa < 30 {
		t.Fatalf("hanya %d berkas terbaca; pembacanya yang rusak", diperiksa)
	}
}

// tco1: tabel warisan HANYA dibaca. Fungsi yang menyebut nama warisan tidak
// boleh memuat kata kerja tulis.
//
// ⚠️ Batasnya dinyatakan: diperiksa per potongan `func`, komentar sebaris
// dibuang. Nama warisan yang dirakit dari potongan teks tidak tertangkap.
func TestTCOWarisanHanyaDibaca(t *testing.T) {
	sumber := berkasSumberProduksi(t)
	diperiksa := 0
	for nama, isi := range sumber {
		if !strings.HasSuffix(nama, ".go") || !strings.Contains(nama, "/tco_") {
			continue
		}
		kode := buangKomentarSumber(nama, isi)
		for _, potongan := range strings.Split(kode, "\nfunc ") {
			atas := strings.ToUpper(potongan)
			menulis := false
			for _, tulis := range []string{"INSERT INTO", "UPDATE ", "DELETE FROM", "MERGE INTO"} {
				if strings.Contains(atas, tulis) {
					menulis = true
				}
			}
			if !menulis {
				continue
			}
			diperiksa++
			// Tabel warisan modul ini DAN master yang hanya dibaca (spec b107).
			dijaga := append(append([]string{}, urutanWarisanTCO...), masterDibacaSajaTCO...)
			for _, warisan := range dijaga {
				if strings.Contains(potongan, `"`+warisan+`"`) || strings.Contains(potongan, "."+warisan) {
					t.Errorf("%s: potongan yang menulis menyebut tabel warisan %s:\n%.120s",
						nama, warisan, strings.TrimSpace(potongan))
				}
			}
		}
	}
	if diperiksa == 0 {
		t.Fatal("nol potongan penulis terbaca; pembacanya yang rusak")
	}
}

// Instrumen penjaga di atas benar-benar membaca: DDL buatan yang salah
// harus tertangkap oleh pembanding tipe yang sama.
func TestTCOPenjagaTipeMasihMenggigit(t *testing.T) {
	const buruk = `CREATE TABLE {skema}.T_UJI (
  ID   VARCHAR2(32) NOT NULL,
  RP   VARCHAR2(1000),
  STARTDATE VARCHAR2(20)
)`
	_, kolom := KolomCreateTable(buruk)
	if len(kolom) != 3 {
		t.Fatalf("pengurai membaca %d kolom, mau 3", len(kolom))
	}
	if TipeTujuanTCO(warisanKlausulTCO, "RP") != TujuanDesimal ||
		TipeTujuanTCO(warisanTahunTCO, "STARTDATE") != TujuanTanggal {
		t.Fatal("peta tipe tujuan tidak mengenal RP/STARTDATE; penjaga di atas hampa")
	}
}
