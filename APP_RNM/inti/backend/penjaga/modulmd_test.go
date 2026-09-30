package penjaga

// Pembaca `modul/<nama>/MODUL.md` - struktur tim satu folder per modul
// (30-09-2026). TANPA Oracle.
//
// Untuk apa berkas ini: penjaga di paket ini berlaku untuk SETIAP modul dan
// tidak memuat satu pun nama modul. Yang khusus satu modul - rentang migrasi
// dan slot menunya (R2, R3), tabel warisan yang ia baca tanpa membuat,
// penyuntikan wajib di handler-nya, dan sebagainya - DINYATAKAN modul itu di
// `MODUL.md`-nya sendiri, dan dibaca dari sini. Mengubahnya karena itu tidak
// pernah menyunting berkas di luar folder modul itu.
//
// Dua bentuk yang dibaca:
//
//  1. tabel kunci-nilai pertama (`| Kunci | Nilai |`) - mis. `Rentang migrasi`;
//  2. bab `## Pernyataan untuk penjaga`: satu judul `### <jenis>` per jenis
//     pernyataan, masing-masing satu tabel. Judul tanpa baris tetap berarti
//     "dinyatakan, kosong"; judul yang tidak ada berarti "tidak dinyatakan".
//
// ⛔ Setiap pembacaan gagal TERANG bila bentuknya salah: pernyataan yang tidak
// terbaca adalah penjaga yang diam-diam mati.

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
)

// folderTemplat - templat folder modul (`modul/_templat`): nilai-nilainya
// penanda, bukan pernyataan sebuah modul.
const folderTemplat = "_templat"

// Judul jenis pernyataan di bab `## Pernyataan untuk penjaga`.
const (
	pernyataanTabelWarisan  = "Tabel warisan: dibaca, tidak dibuat"
	pernyataanNamaBeda      = "Nama STRUKTUR yang berbeda di DDL"
	pernyataanNamaTerlarang = "Nama terlarang di migrasi"
	pernyataanKaskade       = "Kaskade ON DELETE CASCADE"
	pernyataanSuntikan      = "Penyuntikan wajib di handler"
	pernyataanPesanVerbatim = "Pesan verbatim yang bukan nama orang"
)

// jenisPernyataan - setiap judul yang dibaca penjaga. Judul lain di bab
// pernyataan ditolak, bukan diabaikan (`TestPernyataanBerjudulAsingDitolak`).
var jenisPernyataan = map[string]bool{
	pernyataanTabelWarisan:  true,
	pernyataanNamaBeda:      true,
	pernyataanNamaTerlarang: true,
	pernyataanKaskade:       true,
	pernyataanSuntikan:      true,
	pernyataanPesanVerbatim: true,
}

// modulMD adalah isi satu `MODUL.md` yang dibaca penjaga.
type modulMD struct {
	// folder - nama folder `modul/<folder>`.
	folder string
	jalur  string
	// kunci - tabel kunci-nilai pertama, nilai apa adanya (backtick utuh).
	kunci map[string]string
	// pernyataan - judul `###` di bab pernyataan -> baris sel tabelnya.
	pernyataan map[string][][]string
}

// bacaModulMD membaca `MODUL.md` setiap folder `modul/<nama>` selain templat.
//
// ⛔ Folder modul TANPA `MODUL.md` menggagalkan uji: ia modul yang rentang
// migrasinya tidak dinyatakan siapa pun.
func bacaModulMD(t *testing.T) []modulMD {
	t.Helper()
	entri, err := os.ReadDir(filepath.Join(akarAplikasi, "modul"))
	if err != nil {
		t.Fatal(err)
	}
	var hasil []modulMD
	for _, e := range entri {
		if !e.IsDir() || e.Name() == folderTemplat {
			continue
		}
		jalur := filepath.Join(akarAplikasi, "modul", e.Name(), "MODUL.md")
		isi, err := os.ReadFile(jalur)
		if err != nil {
			t.Errorf("modul/%s tanpa MODUL.md yang terbaca: %v", e.Name(), err)
			continue
		}
		m, err := uraiModulMD(string(isi))
		if err != nil {
			t.Errorf("modul/%s/MODUL.md: %v", e.Name(), err)
			continue
		}
		m.folder, m.jalur = e.Name(), jalur
		hasil = append(hasil, m)
	}
	if len(hasil) < 4 {
		t.Fatalf("hanya %d MODUL.md terbaca; pembacanya yang rusak", len(hasil))
	}
	return hasil
}

// uraiModulMD mengurai isi satu MODUL.md.
func uraiModulMD(isi string) (modulMD, error) {
	m := modulMD{kunci: map[string]string{}, pernyataan: map[string][][]string{}}
	baris := strings.Split(strings.ReplaceAll(isi, "\r\n", "\n"), "\n")
	// 1. Tabel kunci-nilai pertama.
	for i := 0; i < len(baris); i++ {
		if sel := selTabel(baris[i]); len(sel) == 2 && sel[0] == "Kunci" && sel[1] == "Nilai" {
			for j := i + 1; j < len(baris) && strings.HasPrefix(strings.TrimSpace(baris[j]), "|"); j++ {
				s := selTabel(baris[j])
				if barisPemisah(s) {
					continue
				}
				if len(s) != 2 {
					return m, fmt.Errorf("baris tabel kunci-nilai tidak berdua sel: %q", baris[j])
				}
				if _, ganda := m.kunci[s[0]]; ganda {
					return m, fmt.Errorf("kunci %q ganda", s[0])
				}
				m.kunci[s[0]] = s[1]
			}
			break
		}
	}
	if len(m.kunci) == 0 {
		return m, fmt.Errorf("tanpa tabel `| Kunci | Nilai |`")
	}
	// 2. Bab pernyataan.
	dalamBab, judul, kepala := false, "", false
	for _, b := range baris {
		t := strings.TrimSpace(b)
		switch {
		case t == "## Pernyataan untuk penjaga":
			dalamBab, judul = true, ""
			continue
		case strings.HasPrefix(t, "## "):
			dalamBab, judul = false, ""
			continue
		case dalamBab && strings.HasPrefix(t, "### "):
			judul, kepala = strings.TrimSpace(strings.TrimPrefix(t, "### ")), false
			if !jenisPernyataan[judul] {
				return m, fmt.Errorf("judul pernyataan %q tidak dikenal penjaga; jenis yang dibaca: docs/bersama/PANDUAN-TIM-PER-MODUL.md bab 6", judul)
			}
			if _, ganda := m.pernyataan[judul]; ganda {
				return m, fmt.Errorf("pernyataan %q ganda", judul)
			}
			m.pernyataan[judul] = [][]string{}
			continue
		}
		if !dalamBab || judul == "" || !strings.HasPrefix(t, "|") {
			continue
		}
		s := selTabel(t)
		if barisPemisah(s) {
			continue
		}
		if !kepala {
			kepala = true // baris pertama tabel = kepala
			continue
		}
		m.pernyataan[judul] = append(m.pernyataan[judul], s)
	}
	return m, nil
}

// selTabel memecah satu baris tabel markdown menjadi sel yang dirapikan.
func selTabel(baris string) []string {
	t := strings.TrimSpace(baris)
	if !strings.HasPrefix(t, "|") || !strings.HasSuffix(t, "|") || len(t) < 2 {
		return nil
	}
	bagian := strings.Split(t[1:len(t)-1], "|")
	for i := range bagian {
		bagian[i] = strings.TrimSpace(bagian[i])
	}
	return bagian
}

// barisPemisah - baris `| --- | --- |` di bawah kepala tabel.
func barisPemisah(sel []string) bool {
	if len(sel) == 0 {
		return false
	}
	for _, s := range sel {
		if strings.Trim(s, "-: ") != "" {
			return false
		}
	}
	return true
}

var polaKode = regexp.MustCompile("`([^`]+)`")

// kodeDi - setiap nilai di dalam backtick sebuah sel, berurutan.
func kodeDi(sel string) []string {
	var hasil []string
	for _, k := range polaKode.FindAllStringSubmatch(sel, -1) {
		hasil = append(hasil, k[1])
	}
	return hasil
}

// satuKode - nilai backtick SATU-SATUNYA di sebuah sel.
func satuKode(t *testing.T, m modulMD, judul, sel string) string {
	t.Helper()
	k := kodeDi(sel)
	if len(k) != 1 {
		t.Fatalf("%s, %q: sel %q harus memuat tepat satu nilai `...`, dapat %d", m.jalur, judul, sel, len(k))
	}
	return k[0]
}

// barisPernyataan - baris jenis `judul` sebuah modul, masing-masing sekurangnya
// `lebar` sel; `ada` = judulnya dinyatakan.
func barisPernyataan(t *testing.T, m modulMD, judul string, lebar int) (baris [][]string, ada bool) {
	t.Helper()
	baris, ada = m.pernyataan[judul]
	for _, b := range baris {
		if len(b) < lebar {
			t.Fatalf("%s, %q: baris %v kurang dari %d sel", m.jalur, judul, b, lebar)
		}
	}
	return baris, ada
}

var polaRentang = regexp.MustCompile(`^(\d{3})-(\d{3})$`)

// rentangNomor membaca `NNN-NNN` (backtick boleh) - nomor SELALU tiga digit.
func rentangNomor(nilai string) (awal, akhir int, err error) {
	v := strings.Trim(strings.TrimSpace(nilai), "`")
	m := polaRentang.FindStringSubmatch(v)
	if m == nil {
		return 0, 0, fmt.Errorf("%q bukan rentang tiga digit NNN-NNN", nilai)
	}
	awal, _ = strconv.Atoi(m[1])
	akhir, _ = strconv.Atoi(m[2])
	if awal > akhir {
		return 0, 0, fmt.Errorf("%q: awal lebih besar dari akhir", nilai)
	}
	return awal, akhir, nil
}

// rentangModul - rentang migrasi dan slot menu sebuah modul (R2, R3).
func rentangModul(t *testing.T, m modulMD) (migrasi, slot [2]int) {
	t.Helper()
	for kunci, tuju := range map[string]*[2]int{"Rentang migrasi": &migrasi, "Slot menu": &slot} {
		nilai, ada := m.kunci[kunci]
		if !ada {
			t.Fatalf("%s tanpa baris %q", m.jalur, kunci)
		}
		a, b, err := rentangNomor(nilai)
		if err != nil {
			t.Fatalf("%s, %q: %v", m.jalur, kunci, err)
		}
		*tuju = [2]int{a, b}
	}
	return migrasi, slot
}

// --- Agregat pernyataan, dibaca penjaga di paket ini ------------------------

// tabelBukanMilikKita - tabel yang STRUKTUR gambarkan tetapi SENGAJA tidak kita
// buat, beserta alasannya, dari pernyataan setiap modul.
func tabelBukanMilikKita(t *testing.T) map[string]string {
	t.Helper()
	hasil := map[string]string{}
	for _, m := range bacaModulMD(t) {
		baris, _ := barisPernyataan(t, m, pernyataanTabelWarisan, 2)
		for _, b := range baris {
			tabel := satuKode(t, m, pernyataanTabelWarisan, b[0])
			if lain, ganda := hasil[tabel]; ganda {
				t.Fatalf("tabel warisan %s dinyatakan dua kali (%s)", tabel, lain)
			}
			hasil[tabel] = b[1]
		}
	}
	return hasil
}

// namaStrukturBeda - nama STRUKTUR -> nama DDL, untuk `jenis` tabel atau kolom.
func namaStrukturBeda(t *testing.T, jenis string) map[string]string {
	t.Helper()
	hasil := map[string]string{}
	for _, m := range bacaModulMD(t) {
		baris, _ := barisPernyataan(t, m, pernyataanNamaBeda, 3)
		for _, b := range baris {
			if b[0] != "tabel" && b[0] != "kolom" {
				t.Fatalf("%s, %q: jenis %q, mau tabel atau kolom", m.jalur, pernyataanNamaBeda, b[0])
			}
			if b[0] != jenis {
				continue
			}
			hasil[satuKode(t, m, pernyataanNamaBeda, b[1])] = satuKode(t, m, pernyataanNamaBeda, b[2])
		}
	}
	return hasil
}

// namaTerlarangDiMigrasi - nama yang tidak boleh muncul di migrasi mana pun.
func namaTerlarangDiMigrasi(t *testing.T) map[string]string {
	t.Helper()
	hasil := map[string]string{}
	for _, m := range bacaModulMD(t) {
		baris, _ := barisPernyataan(t, m, pernyataanNamaTerlarang, 2)
		for _, b := range baris {
			hasil[satuKode(t, m, pernyataanNamaTerlarang, b[0])] = b[1]
		}
	}
	return hasil
}

// kaskadePerModul - modul yang MENYATAKAN kebijakan kaskadenya -> awalan
// berkas migrasinya yang boleh `ON DELETE CASCADE`.
func kaskadePerModul(t *testing.T) map[string]map[string]bool {
	t.Helper()
	hasil := map[string]map[string]bool{}
	for _, m := range bacaModulMD(t) {
		baris, ada := barisPernyataan(t, m, pernyataanKaskade, 1)
		if !ada {
			continue
		}
		awalan := map[string]bool{}
		for _, b := range baris {
			awalan[satuKode(t, m, pernyataanKaskade, b[0])] = true
		}
		hasil[m.folder] = awalan
	}
	return hasil
}

// suntikan adalah satu pasangan: penyusun layanan, dan yang wajib menyertainya.
type suntikan struct {
	// penyusun adalah pemanggilan yang mengembalikan layanan ber-stub.
	penyusun string
	// wajib adalah penyuntikan yang harus muncul di fungsi yang sama.
	wajib []string
}

// petaSuntikan - berkas handler (jalur relatif akar aplikasi) -> kewajibannya,
// dari pernyataan setiap modul.
func petaSuntikan(t *testing.T) map[string][]suntikan {
	t.Helper()
	hasil := map[string][]suntikan{}
	for _, m := range bacaModulMD(t) {
		baris, _ := barisPernyataan(t, m, pernyataanSuntikan, 3)
		for _, b := range baris {
			berkas := filepath.ToSlash(filepath.Join("modul", m.folder, "backend", "handlers",
				satuKode(t, m, pernyataanSuntikan, b[0])))
			wajib := kodeDi(b[2])
			if len(wajib) == 0 {
				t.Fatalf("%s, %q: %s tanpa satu pun penyuntikan wajib", m.jalur, pernyataanSuntikan, b[0])
			}
			hasil[berkas] = append(hasil[berkas], suntikan{penyusun: satuKode(t, m, pernyataanSuntikan, b[1]), wajib: wajib})
		}
	}
	return hasil
}

// pesanVerbatimDinyatakan - (paket relatif akar aplikasi, konstanta) yang
// cocok dengan pola nama orang tetapi bukan nama orang.
func pesanVerbatimDinyatakan(t *testing.T) map[string][]string {
	t.Helper()
	hasil := map[string][]string{}
	for _, m := range bacaModulMD(t) {
		baris, _ := barisPernyataan(t, m, pernyataanPesanVerbatim, 2)
		for _, b := range baris {
			paket := filepath.ToSlash(filepath.Join("modul", m.folder, satuKode(t, m, pernyataanPesanVerbatim, b[0])))
			hasil[paket] = append(hasil[paket], satuKode(t, m, pernyataanPesanVerbatim, b[1]))
		}
	}
	for _, k := range hasil {
		sort.Strings(k)
	}
	return hasil
}

// ⛔ Judul `###` yang bukan salah satu jenis pernyataan ditolak: judul yang
// salah ketik (mis. `Kaskade ON DELETE CASCADE` tanpa satu huruf) dulu dibaca
// sebagai jenis BARU yang tidak dibaca penjaga mana pun - pernyataan modul itu
// diam-diam mati, dan penjaganya tetap hijau (temuan code review 30-09-2026).
func TestPernyataanBerjudulAsingDitolak(t *testing.T) {
	kepala := "| Kunci | Nilai |\n| --- | --- |\n| Nama modul | `alfa` |\n\n## Pernyataan untuk penjaga\n\n"
	tabel := "\n\n| Kode | Alasan |\n| --- | --- |\n| `ALFA_A` | uji |\n"
	if _, err := uraiModulMD(kepala + "### " + pernyataanKaskade + tabel); err != nil {
		t.Fatalf("prasyarat: judul yang dikenal ditolak: %v", err)
	}
	salah := strings.Replace(pernyataanKaskade, "CASCADE", "CASCAD", 1)
	if salah == pernyataanKaskade {
		t.Fatal("prasyarat: pengganti tidak cocok, cacat tidak terpasang")
	}
	_, err := uraiModulMD(kepala + "### " + salah + tabel)
	if err == nil {
		t.Fatalf("judul %q diterima sebagai jenis pernyataan", salah)
	}
	if !strings.Contains(err.Error(), salah) {
		t.Errorf("galat %q tidak menyebut judulnya", err)
	}
}
