package repository

// Penjaga berkas data uji sintetis - GILIRAN-13 paket 3.
//
// ⛔ Berkasnya TIDAK PERNAH dijalankan executor (brief GILIRAN-13 §2 paket 3):
// work owner memuatnya ke skema uji. Karena itu seluruh kebenarannya yang
// dapat diperiksa tanpa Oracle diperiksa DI SINI - nama tabel dan kolom
// terhadap DDL migrasi, cacah nilai, pagar skema, dan nol kebocoran.

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"nusantarare/internal/config"
	"nusantarare/internal/models"
)

const letakDataUji = "skemauji/data_uji_tiga_modul.sql"

// sisipDataUji adalah satu INSERT di berkas data uji.
type sisipDataUji struct {
	tabel string
	kolom []string
	nilai []string
}

// ambil mengembalikan nilai mentah sebuah kolom - "" bila tidak disebut
// ATAU bila bernilai NULL.
//
// ⛔ NULL = kosong. Ronde pertama mengembalikan teks "NULL" apa adanya, dan
// uji mutasi membuktikannya: `ACCOUNT_NO` yang diganti NULL tetap terhitung
// "rekening lengkap".
func (s sisipDataUji) ambil(kolom string) string {
	for i, k := range s.kolom {
		if k == kolom {
			if strings.EqualFold(s.nilai[i], "NULL") {
				return ""
			}
			return s.nilai[i]
		}
	}
	return ""
}

// teks mengembalikan isi literal teks sebuah kolom, tanpa kutip.
func (s sisipDataUji) teks(kolom string) string {
	return strings.TrimSuffix(strings.TrimPrefix(s.ambil(kolom), "'"), "'")
}

// bacaDataUji membaca berkasnya TANPA komentar `--`.
func bacaDataUji(t *testing.T) string {
	t.Helper()
	isi, err := os.ReadFile(filepath.FromSlash(letakDataUji))
	if err != nil {
		t.Fatalf("membaca %s: %v", letakDataUji, err)
	}
	var keluar []string
	for _, b := range strings.Split(strings.ReplaceAll(string(isi), "\r\n", "\n"), "\n") {
		if i := indeksKomentar(b); i >= 0 {
			b = b[:i]
		}
		keluar = append(keluar, b)
	}
	return strings.Join(keluar, "\n")
}

// indeksKomentar mencari `--` di luar literal teks.
func indeksKomentar(baris string) int {
	dalamKutip := false
	for i := 0; i < len(baris); i++ {
		switch {
		case baris[i] == '\'':
			dalamKutip = !dalamKutip
		case !dalamKutip && strings.HasPrefix(baris[i:], "--"):
			return i
		}
	}
	return -1
}

// pecahNilai memecah daftar VALUES pada koma di luar kutip dan kurung.
func pecahNilai(daftar string) []string {
	var hasil []string
	var kini strings.Builder
	dalamKutip, kedalaman := false, 0
	for _, r := range daftar {
		switch {
		case r == '\'':
			dalamKutip = !dalamKutip
		case !dalamKutip && r == '(':
			kedalaman++
		case !dalamKutip && r == ')':
			kedalaman--
		case !dalamKutip && kedalaman == 0 && r == ',':
			hasil = append(hasil, strings.TrimSpace(kini.String()))
			kini.Reset()
			continue
		}
		kini.WriteRune(r)
	}
	return append(hasil, strings.TrimSpace(kini.String()))
}

var polaSisipDataUji = regexp.MustCompile(
	`(?s)INSERT INTO ([A-Z_]+)\s*\(([^)]*)\)\s*VALUES\s*\((.*?)\);`)

func sisipanDataUji(t *testing.T, sql string) []sisipDataUji {
	t.Helper()
	var hasil []sisipDataUji
	for _, m := range polaSisipDataUji.FindAllStringSubmatch(sql, -1) {
		s := sisipDataUji{tabel: m[1], nilai: pecahNilai(m[3])}
		for _, k := range strings.Split(m[2], ",") {
			s.kolom = append(s.kolom, strings.TrimSpace(k))
		}
		hasil = append(hasil, s)
	}
	if len(hasil) < 10 {
		t.Fatalf("hanya %d INSERT terbaca; pembacanya yang rusak, bukan berkasnya", len(hasil))
	}
	return hasil
}

// Kolom tiruan EMAILKOMITE - dibaca dari DDL di berkas itu sendiri.
var polaTiruanRoster = regexp.MustCompile(`(?s)CREATE TABLE EMAILKOMITE \((.*?)\)'`)

func kolomTiruanRoster(t *testing.T, sql string) []string {
	t.Helper()
	m := polaTiruanRoster.FindStringSubmatch(sql)
	if m == nil {
		t.Fatal("DDL tiruan EMAILKOMITE tidak ditemukan")
	}
	var kolom []string
	for _, b := range strings.Split(m[1], ",") {
		b = strings.TrimSpace(b)
		if f := strings.Fields(b); len(f) > 0 && f[0] != "CONSTRAINT" {
			kolom = append(kolom, f[0])
		}
	}
	return kolom
}

// TestDataUjiCocokDenganDDL - setiap tabel dan kolom benar-benar ada.
//
// ⛔ Inilah pengganti "sudah dicoba di Oracle": executor tidak menjalankannya,
// jadi kolom yang salah eja hanya akan ketahuan di tangan work owner - kecuali
// dibandingkan di sini dengan DDL yang sama dengan yang `-migrate` jalankan.
func TestDataUjiCocokDenganDDL(t *testing.T) {
	sql := bacaDataUji(t)
	ddl := kolomMenurutDDL(t)
	ddl["EMAILKOMITE"] = kolomTiruanRoster(t, sql)
	for i, s := range sisipanDataUji(t, sql) {
		kolom, ada := ddl[s.tabel]
		if !ada {
			t.Errorf("INSERT #%d: tabel %s tidak dibuat migrasi mana pun", i+1, s.tabel)
			continue
		}
		dikenal := map[string]bool{}
		for _, k := range kolom {
			dikenal[k] = true
		}
		for _, k := range s.kolom {
			if !dikenal[k] {
				t.Errorf("INSERT #%d %s: kolom %s tidak ada di DDL", i+1, s.tabel, k)
			}
		}
		if len(s.kolom) != len(s.nilai) {
			t.Errorf("INSERT #%d %s: %d kolom, %d nilai", i+1, s.tabel, len(s.kolom), len(s.nilai))
		}
		// Satu-satunya kolom NOT NULL di ketujuh tabel yang diisi adalah ID.
		if s.ambil("ID") == "" {
			t.Errorf("INSERT #%d %s: tanpa ID", i+1, s.tabel)
		}
	}
}

// TestDataUjiNolKebocoran - seluruh pengenal `UJI-*`, surel `@contoh.invalid`.
//
// ⛔ Setiap literal teks yang BUKAN kode domain harus berawalan `UJI-`/`uji-`.
// Daftar kode di bawah TERTUTUP: menambah satu berarti menyatakan bahwa ia
// kode VERBATIM, bukan data.
func TestDataUjiNolKebocoran(t *testing.T) {
	// Irisan, bukan literal peta: beberapa kode berbagi nilai ("0" bendera
	// Input Offer sekaligus STS_REJECT Outstanding).
	kode := map[string]bool{}
	for _, k := range []string{
		models.LiniLife, models.PosisiOffer, models.PosisiPremium,
		models.TahapPolisPenawaran, models.TahapPolisDetail,
		models.TahapPolisSummary, models.StatusPolisSelesai,
		models.FlagPolisPenawaran, models.FlagPolisPremium,
		models.KodeOutstanding, models.KodeDitolak,
		models.PeranAdminLife, models.PeranMedicalLife, models.PeranSPVLife,
		models.TahapInputRegister.String(), models.TahapOutstanding.String(),
		models.TahapMedicalCheck.String(), models.TahapClaimAnalis.String(),
		models.PenandaDipilih, "false",
		"QP", "L1", "IDR", "M", "F",
	} {
		kode[k] = true
	}
	sql := bacaDataUji(t)
	for _, s := range sisipanDataUji(t, sql) {
		for i, v := range s.nilai {
			if !strings.HasPrefix(v, "'") {
				continue
			}
			isi := strings.Trim(v, "'")
			switch {
			case strings.Contains(isi, "@"):
				if !strings.HasPrefix(isi, "uji-") || !strings.HasSuffix(isi, "@contoh.invalid") {
					t.Errorf("%s.%s: surel %q bukan uji-…@contoh.invalid", s.tabel, s.kolom[i], isi)
				}
			case kode[isi]:
			case strings.HasPrefix(isi, "UJI-"):
			default:
				t.Errorf("%s.%s: %q bukan pengenal UJI-* dan bukan kode domain yang terdaftar",
					s.tabel, s.kolom[i], isi)
			}
		}
	}
}

// TestDataUjiBerpagarSepertiMigrateDown - pagar skema di SETIAP blok.
//
// ⛔ "Pagar yang sama dengan `-migrate-down`" (brief): `config.PagarSkemaUji`
// menolak skema yang MEMUAT `POOLDATA`. SQL tidak dapat membaca env
// `ORACLE_SKEMA_UJI`, jadi yang ditiru syarat keduanya - dengan nama yang
// DIAMBIL dari konstanta itu, bukan diketik ulang.
//
// Setiap blok memeriksanya sendiri, SEBELUM pernyataan pengubah pertamanya:
// alat yang meneruskan skrip sesudah galat (bukan SQL*Plus) tetap berhenti
// di blok berikutnya.
func TestDataUjiBerpagarSepertiMigrateDown(t *testing.T) {
	sql := bacaDataUji(t)
	pagar := "INSTR(UPPER(v_skema), '" + config.NamaSkemaWarisan + "') > 0"
	var blok []string
	for _, b := range strings.Split(sql, "\n/\n") {
		if strings.Contains(b, "BEGIN") {
			blok = append(blok, b)
		}
	}
	if len(blok) < 2 {
		t.Fatalf("%d blok PL/SQL, mau sedikitnya 2 (pagar+tiruan, data)", len(blok))
	}
	for i, b := range blok {
		p := strings.Index(b, pagar)
		if p < 0 {
			t.Errorf("blok %d tanpa pagar %q", i+1, pagar)
			continue
		}
		for _, ubah := range []string{"INSERT INTO", "EXECUTE IMMEDIATE 'CREATE"} {
			if u := strings.Index(b, ubah); u >= 0 && u < p {
				t.Errorf("blok %d: %q mendahului pagarnya", i+1, ubah)
			}
		}
	}
	// Seluruh INSERT di SATU blok: gagal di mana pun = nol baris tertinggal.
	berisi := 0
	for _, b := range blok {
		if strings.Contains(b, "INSERT INTO") {
			berisi++
		}
	}
	if berisi != 1 {
		t.Errorf("INSERT tersebar di %d blok, mau 1", berisi)
	}
	// Migrasi 057 wajib sudah berjalan - kolom bendera dan sequence-nya.
	if !strings.Contains(sql, "'"+kunciLangkah("057_seq_work_polis_dan_flag_ongoing.sql")+"'") {
		t.Error("berkas tidak memeriksa T_MIGRASI untuk 057")
	}
	// ⛔ Nol pernyataan perusak dan nol COMMIT (ADR-U-0029): work owner yang
	// memutuskan kapan menetapkannya.
	atas := strings.ToUpper(sql)
	for _, dilarang := range []string{"COMMIT", "DELETE ", "UPDATE ", "DROP ", "TRUNCATE", "MERGE "} {
		if strings.Contains(atas, dilarang) {
			t.Errorf("berkas memuat %q", strings.TrimSpace(dilarang))
		}
	}
}

// TestDataUjiMenutupSetiapTahap - satu kasus per tahap, dua modul.
func TestDataUjiMenutupSetiapTahap(t *testing.T) {
	sisip := sisipanDataUji(t, bacaDataUji(t))
	polis, klaim := map[string]string{}, map[string]string{}
	for _, s := range sisip {
		switch s.tabel {
		case "T_WORK_POLIS":
			polis[s.teks("STATUS")] = s.teks("ID")
		case "T_WORK_CLAIM":
			klaim[s.teks("TAHAP")] = s.teks("ID")
		}
	}
	for _, tahap := range []string{models.TahapPolisPenawaran, models.TahapPolisDetail,
		models.TahapPolisSummary, models.StatusPolisSelesai} {
		if polis[tahap] == "" {
			t.Errorf("tidak ada polis di tahap %q", tahap)
		}
	}
	for _, tahap := range []models.Tahap{models.TahapInputRegister, models.TahapOutstanding,
		models.TahapMedicalCheck, models.TahapClaimAnalis} {
		if klaim[tahap.String()] == "" {
			t.Errorf("tidak ada klaim di tahap %q", tahap)
		}
	}
	// PY_POSITION = pemegang tahapnya (butir at): tanpa itu `TahapBerlaku`
	// dan kolom `TAHAP` dapat berselisih.
	for _, s := range sisip {
		if s.tabel != "T_WORK_CLAIM" {
			continue
		}
		tahap := models.TahapDariNama(s.teks("TAHAP"))
		if peran, _ := models.PeranPemegangTahap(tahap); peran != s.teks("PY_POSITION") {
			t.Errorf("%s: tahap %q dipegang %q, PY_POSITION %q",
				s.teks("ID"), s.teks("TAHAP"), peran, s.teks("PY_POSITION"))
		}
	}
}

// TestDataUjiPunyaBarisSiapSerahKomite - satu baris lolos seluruh gerbang
// penyerahan yang dapat diperiksa tanpa Oracle, dan roster menutupnya.
//
// Gerbangnya `services.Penyerahan.Serahkan`: Outstanding, belum tertaut,
// rekening lengkap (`PeriksaRekening`), jumlah terisi; roster `STS_KLAIM`
// `LIFE` (`services.StatusKlaimRoster`), `STS_AKTIF` "1", `LIMIT_BOTTOM` <=
// nilai klaim.
func TestDataUjiPunyaBarisSiapSerahKomite(t *testing.T) {
	sisip := sisipanDataUji(t, bacaDataUji(t))
	var siap []sisipDataUji
	for _, s := range sisip {
		if s.tabel == "T_CLAIMLF_ADJUSTMENT" && s.teks("STS_REJECT") == models.KodeOutstanding &&
			s.ambil("KOMITE_ID") == "" && s.teks("NAME_OF_BANK") != "" &&
			s.teks("ID_BANK") != "" && s.teks("ACCOUNT_NO") != "" && s.ambil("CLAIM_AMOUNT") != "" {
			siap = append(siap, s)
		}
	}
	if len(siap) != 1 {
		t.Fatalf("%d baris siap serah Komite, mau tepat 1", len(siap))
	}
	nilai := siap[0].ambil("CLAIM_AMOUNT")
	tingkat := 0
	for _, s := range sisip {
		if s.tabel == "EMAILKOMITE" && s.teks("STS_KLAIM") == "LIFE" && s.teks("STS_AKTIF") == "1" &&
			bandingAngka(s.ambil("LIMIT_BOTTOM"), nilai) <= 0 {
			tingkat++
		}
	}
	if tingkat < 1 {
		t.Errorf("roster tidak menutup nilai klaim %s", nilai)
	}
}

// bandingAngka membandingkan dua literal bilangan bulat non-negatif tanpa
// float (ADR-U-0003): panjang dulu, lalu leksikal.
func bandingAngka(a, b string) int {
	a = strings.SplitN(strings.TrimLeft(a, "0"), ".", 2)[0]
	b = strings.SplitN(strings.TrimLeft(b, "0"), ".", 2)[0]
	switch {
	case len(a) != len(b):
		if len(a) < len(b) {
			return -1
		}
		return 1
	case a < b:
		return -1
	case a > b:
		return 1
	}
	return 0
}
