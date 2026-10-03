// Program bangkit menulis `skema_gen.go` paket loader dari workbook rancangan tabel
// flat (tiket 22). Ia dijalankan tangan, BUKAN lewat `go generate`: workbook tidak
// ada di repositori, dan jalurnya konfigurasi, bukan literal.
//
//	go run ./modul/nbfacin/backend/services/loader/bangkit \
//	    -xlsx "D:\migrasi\RNM\OUTPUT\08-flat\Tabel-Flat-Lintas-Siklus.xlsx" \
//	    -kandidat "D:\migrasi\RNM\Claude outputs\Tabel-Flat-per-Grup-Bisnis.xlsx" \
//	    -keluar modul/nbfacin/backend/services/loader/skema_gen.go
//
// Dari workbook `-xlsx` yang dibaca HANYA dua lembar - `Kolom` (tabel, kolom, tipe
// usulan, jenis, FIELD ASLI, siklus, NULL, catatan pewarisan mata uang) dan `Jalur
// Sumber` (jalur sesudah keputusan V-22/V-24b/V-30, tabel, tabel induk) - dan isinya
// nama tabel, nama kolom, nama medan, dan tipe. Lembar berisi nilai contoh (`Kode Enumerasi`) tidak
// dibaca. Dari workbook `-kandidat` HANYA lembar `Kandidat Hapus`, dan hanya baris
// medan penunjuk `Idx*`/`Index*` (tabel, medan, kategori, status) - catatan rancangan.
// P7 (butir 70) dulu membuang sebagian penunjuk menurut lembar ini; butir 71 mencabut
// penerapannya (data membantah premis V-47), jadi mesin tidak membacanya.
// Penjaganya `TestSkemaSamaDenganKeluaranBangkit`: dengan env `FLAT_RNM_XLSX` dan
// `FLAT_RNM_KANDIDAT_XLSX` terisi, ia membangkitkan ulang di memori dan menolak
// `skema_gen.go` yang berbeda.
//
// [terverifikasi] 02-10-2026, dua cara: workbook ini = `DDL-tabel-flat-draf.sql`
// persis (78 tabel, 1.329 kolom; selisih nol dua arah, pengurai xlsx vs pencacah
// baris SQL). `Claude outputs\Tabel-Flat-per-Grup-Bisnis.xlsx` (75 / 1.290) adalah
// himpunan bagiannya: selisih 39 kolom = 3 tabel wadah V-27 (20 kolom sistem) + 18
// `CURRENCY_CODE` (K-063/K-069) + `OLD_POLIS_ID` (K-068), nol kolom ke arah
// sebaliknya - `docs/issues/22-loader-flatten.md`.
package main

import (
	"archive/zip"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"go/format"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func main() {
	jalurXlsx := flag.String("xlsx", "", "workbook Tabel-Flat-Lintas-Siklus.xlsx (baca saja)")
	jalurKandidat := flag.String("kandidat", "", "workbook Tabel-Flat-per-Grup-Bisnis.xlsx, lembar Kandidat Hapus (baca saja)")
	keluar := flag.String("keluar", "", "berkas skema_gen.go")
	flag.Parse()
	if *jalurXlsx == "" || *jalurKandidat == "" || *keluar == "" {
		fmt.Fprintln(os.Stderr, "bangkit: -xlsx, -kandidat dan -keluar wajib")
		os.Exit(2)
	}
	src, ringkas, err := bangkitkan(*jalurXlsx, *jalurKandidat)
	if err != nil {
		fmt.Fprintln(os.Stderr, "bangkit:", err)
		os.Exit(1)
	}
	if err := os.WriteFile(*keluar, src, 0o644); err != nil {
		fmt.Fprintln(os.Stderr, "bangkit:", err)
		os.Exit(1)
	}
	fmt.Println(ringkas)
}

// Jumlah yang wajib dicapai - angka rancangan yang sudah diukur dua cara.
const (
	jumlahTabel = 78
	jumlahKolom = 1329
)

type kolom struct {
	tabel, nama, tipe, jenis, medan, siklus, null, catatan string
}

type jalur struct{ jalur, tabel, induk string }

var polaWaris = regexp.MustCompile(`diwarisi dari (T_[A-Z_]+)`)

var polaPenunjuk = regexp.MustCompile(`^(Idx|Index)[A-Z]`)

// bangkitkan membaca workbook dan mengembalikan sumber Go yang sudah diformat.
func bangkitkan(jalurXlsx, jalurKandidat string) ([]byte, string, error) {
	kandidat, sumKandidat, err := bacaKandidat(jalurKandidat)
	if err != nil {
		return nil, "", fmt.Errorf("kandidat: %w", err)
	}
	mentah, err := os.ReadFile(jalurXlsx)
	if err != nil {
		return nil, "", err
	}
	wb, err := bukaWorkbook(mentah)
	if err != nil {
		return nil, "", err
	}
	barisKolom, err := wb.lembar("Kolom", 14)
	if err != nil {
		return nil, "", err
	}
	barisJalur, err := wb.lembar("Jalur Sumber", 4)
	if err != nil {
		return nil, "", err
	}
	if len(barisKolom) == 0 || strings.Join(barisKolom[0][:6], "|") != "TABEL|KOLOM|PJG|TIPE USULAN|JENIS|FIELD ASLI" ||
		barisKolom[0][9] != "SIKLUS" || barisKolom[0][13] != "NULL" || barisKolom[0][8] != "CATATAN" {
		return nil, "", errors.New("kepala lembar Kolom tidak dikenal")
	}
	var semua []kolom
	tabel := map[string][]kolom{}
	for _, r := range barisKolom[1:] {
		if r[0] == "" {
			continue
		}
		k := kolom{tabel: r[0], nama: r[1], tipe: r[3], jenis: r[4], medan: r[5], siklus: r[9], null: r[13], catatan: r[8]}
		switch k.jenis {
		case "data", "sistem", "PII", "kunci":
		default:
			return nil, "", fmt.Errorf("%s.%s: jenis %q tidak dikenal", k.tabel, k.nama, k.jenis)
		}
		if k.null != "NULL" && k.null != "NOT NULL" {
			return nil, "", fmt.Errorf("%s.%s: NULL %q tidak dikenal", k.tabel, k.nama, k.null)
		}
		semua = append(semua, k)
		tabel[k.tabel] = append(tabel[k.tabel], k)
	}
	if len(tabel) != jumlahTabel || len(semua) != jumlahKolom {
		return nil, "", fmt.Errorf("lembar Kolom %d tabel / %d kolom, mau %d / %d", len(tabel), len(semua), jumlahTabel, jumlahKolom)
	}
	// Lembar Jalur Sumber: baris 0 catatan, baris 1 kepala.
	if len(barisJalur) < 2 || strings.Join(barisJalur[1], "|") != "NAMA TABEL|JALUR SUMBER (JSON)|TABEL INDUK|GRUP" {
		return nil, "", errors.New("kepala lembar Jalur Sumber tidak dikenal")
	}
	var daftarJalur []jalur
	for _, r := range barisJalur[2:] {
		if r[0] == "" {
			continue
		}
		j := jalur{tabel: r[0], induk: r[2]}
		switch r[1] {
		case "(objek kerja Pega)":
			continue // T_WORK_POLIS: lahir dari IDPEGA, bukan dari jalur dokumen
		case "(akar pagedata)":
			j.jalur = ""
		default:
			// BAHAN §0: jalur lembar ini diturunkan dari XML; di JSON ruas kode mata uang tidak ada.
			j.jalur = strings.ReplaceAll(r[1], "/<kode mata uang>", "")
		}
		if _, ada := tabel[j.tabel]; !ada {
			return nil, "", fmt.Errorf("jalur %q: tabel %s tidak ada di lembar Kolom", r[1], j.tabel)
		}
		if _, ada := tabel[j.induk]; !ada {
			return nil, "", fmt.Errorf("jalur %q: induk %s tidak ada di lembar Kolom", r[1], j.induk)
		}
		daftarJalur = append(daftarJalur, j)
	}
	waris := map[string]string{}
	for _, k := range semua {
		if k.siklus == "turunan rancangan" && k.nama == "CURRENCY_CODE" && k.jenis == "sistem" {
			m := polaWaris.FindStringSubmatch(k.catatan)
			if m == nil {
				return nil, "", fmt.Errorf("%s.CURRENCY_CODE turunan tanpa sumber pewarisan", k.tabel)
			}
			waris[k.tabel] = m[1]
		}
	}

	var b bytes.Buffer
	sum := md5.Sum(mentah)
	fmt.Fprintf(&b, "// Code generated by ./bangkit; DO NOT EDIT.\n\n")
	fmt.Fprintf(&b, "// Sumber: %s (md5 %s) - lembar Kolom dan Jalur Sumber.\n", filepath.Base(jalurXlsx), hex.EncodeToString(sum[:]))
	fmt.Fprintf(&b, "// Sumber: %s (md5 %s) - lembar Kandidat Hapus, baris penunjuk Idx*/Index*.\n", filepath.Base(jalurKandidat), sumKandidat)
	fmt.Fprintf(&b, "// %d tabel, %d kolom, %d jalur, %d pewarisan mata uang, %d penunjuk.\n\npackage loader\n\n", len(tabel), len(semua), len(daftarJalur), len(waris), len(kandidat))
	nama := make([]string, 0, len(tabel))
	for t := range tabel {
		nama = append(nama, t)
	}
	sort.Strings(nama)
	b.WriteString("var skemaTabel = map[string][]kolomSkema{\n")
	for _, t := range nama {
		fmt.Fprintf(&b, "%s: {\n", strconv.Quote(t))
		for _, k := range tabel[t] {
			medan := k.medan
			if medan == "-" {
				medan = ""
			}
			fmt.Fprintf(&b, "{nama: %s, tipe: %s, medan: %s, wajib: %t, turunan: %t},\n",
				strconv.Quote(k.nama), strconv.Quote(k.tipe), strconv.Quote(medan),
				k.null == "NOT NULL", k.siklus == "turunan rancangan")
		}
		b.WriteString("},\n")
	}
	b.WriteString("}\n\nvar jalurSumber = []jalurSkema{\n")
	for _, j := range daftarJalur {
		fmt.Fprintf(&b, "{jalur: %s, tabel: %s, induk: %s},\n", strconv.Quote(j.jalur), strconv.Quote(j.tabel), strconv.Quote(j.induk))
	}
	b.WriteString("}\n\n// warisMataUang - K-063 kelompok (c): tabel -> tabel leluhur pembawa mata uangnya.\nvar warisMataUang = map[string]string{\n")
	wk := make([]string, 0, len(waris))
	for t := range waris {
		wk = append(wk, t)
	}
	sort.Strings(wk)
	for _, t := range wk {
		fmt.Fprintf(&b, "%s: %s,\n", strconv.Quote(t), strconv.Quote(waris[t]))
	}
	b.WriteString("}\n\n// penunjukKandidat - lembar Kandidat Hapus: \"TABEL.Medan\" -> {kategori, status}.\nvar penunjukKandidat = map[string]kandidatPenunjuk{\n")
	kk := make([]string, 0, len(kandidat))
	for k := range kandidat {
		kk = append(kk, k)
	}
	sort.Strings(kk)
	for _, k := range kk {
		fmt.Fprintf(&b, "%s: {kategori: %s, status: %s},\n", strconv.Quote(k), strconv.Quote(kandidat[k][0]), strconv.Quote(kandidat[k][1]))
	}
	b.WriteString("}\n")
	src, err := format.Source(b.Bytes())
	if err != nil {
		return nil, "", err
	}
	return src, fmt.Sprintf("%d tabel, %d kolom, %d jalur, %d pewarisan, %d penunjuk", len(tabel), len(semua), len(daftarJalur), len(waris), len(kandidat)), nil
}

// bacaKandidat - baris penunjuk Idx*/Index* lembar `Kandidat Hapus`: kunci
// "TABEL.Medan" -> {KATEGORI, STATUS DI SKEMA INI}. Satu kunci dengan dua baris
// disimpan KEDUANYA, dipisah " + " (kategori) dan " | " (status) - tidak dipilih:
// [terverifikasi] 02-10-2026 lembar itu sendiri berselisih untuk DUA kunci -
// T_FR_ANEKALIST.IdxOccupation dan T_FR_DEDUCTIBLELIST.IndexProperty (R2 "SUDAH
// DIHAPUS" dan R3 "SUDAH JADI FK (V-47)"); keduanya disimpan apa adanya. Salinan polaPenunjuk paket loader: paket main ini
// tidak boleh mengimpor paket yang ia bangkitkan.
func bacaKandidat(jalur string) (map[string][2]string, string, error) {
	mentah, err := os.ReadFile(jalur)
	if err != nil {
		return nil, "", err
	}
	wb, err := bukaWorkbook(mentah)
	if err != nil {
		return nil, "", err
	}
	baris, err := wb.lembar("Kandidat Hapus", 8)
	if err != nil {
		return nil, "", err
	}
	if len(baris) == 0 || baris[0][0] != "TABEL" || baris[0][1] != "KOLOM" || baris[0][2] != "KATEGORI" || baris[0][7] != "STATUS DI SKEMA INI" {
		return nil, "", errors.New("kepala lembar Kandidat Hapus tidak dikenal")
	}
	hasil := map[string][2]string{}
	for _, r := range baris[1:] {
		if r[0] == "" || !polaPenunjuk.MatchString(r[1]) {
			continue
		}
		k := r[0] + "." + r[1]
		if lama, ada := hasil[k]; ada {
			status := lama[1]
			if status != r[7] {
				status += " | " + r[7]
			}
			hasil[k] = [2]string{lama[0] + " + " + r[2], status}
			continue
		}
		hasil[k] = [2]string{r[2], r[7]}
	}
	sum := md5.Sum(mentah)
	return hasil, hex.EncodeToString(sum[:]), nil
}

// ---- pembaca xlsx minimal (zip + XML), tanpa dependensi ----

type workbook struct {
	berkas  map[string]*zip.File
	bersama []string
	lembarK map[string]string // nama lembar -> jalur XML-nya di dalam zip
}

func bukaWorkbook(mentah []byte) (*workbook, error) {
	z, err := zip.NewReader(bytes.NewReader(mentah), int64(len(mentah)))
	if err != nil {
		return nil, err
	}
	wb := &workbook{berkas: map[string]*zip.File{}, lembarK: map[string]string{}}
	for _, f := range z.File {
		wb.berkas[f.Name] = f
	}
	if f, ada := wb.berkas["xl/sharedStrings.xml"]; ada {
		var sst struct {
			SI []struct {
				T string `xml:"t"`
				R []struct {
					T string `xml:"t"`
				} `xml:"r"`
			} `xml:"si"`
		}
		if err := bacaXML(f, &sst); err != nil {
			return nil, err
		}
		for _, si := range sst.SI {
			s := si.T
			for _, r := range si.R {
				s += r.T
			}
			wb.bersama = append(wb.bersama, s)
		}
	}
	var buku struct {
		Sheets []struct {
			Name string `xml:"name,attr"`
			// Atribut `r:id`: tag tanpa ruang nama mencocokkan nama lokal `id` dari ruang
			// nama mana pun - satu-satunya atribut `id` pada <sheet>.
			RID string `xml:"id,attr"`
		} `xml:"sheets>sheet"`
	}
	var rel struct {
		R []struct {
			ID     string `xml:"Id,attr"`
			Target string `xml:"Target,attr"`
		} `xml:"Relationship"`
	}
	if err := bacaXML(wb.berkas["xl/workbook.xml"], &buku); err != nil {
		return nil, err
	}
	if err := bacaXML(wb.berkas["xl/_rels/workbook.xml.rels"], &rel); err != nil {
		return nil, err
	}
	target := map[string]string{}
	for _, r := range rel.R {
		target[r.ID] = "xl/" + strings.TrimPrefix(strings.TrimPrefix(r.Target, "/"), "xl/")
	}
	for _, s := range buku.Sheets {
		wb.lembarK[s.Name] = target[s.RID]
	}
	return wb, nil
}

func bacaXML(f *zip.File, v any) error {
	if f == nil {
		return errors.New("bagian workbook tidak ada")
	}
	r, err := f.Open()
	if err != nil {
		return err
	}
	defer r.Close()
	isi, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	return xml.Unmarshal(isi, v)
}

// lembar mengembalikan baris dengan sel DITEMPATKAN menurut huruf kolomnya - sel
// kosong menjadi "" dan tidak menggeser sel sesudahnya (cacat instrumen yang
// tertangkap 02-10-2026 pada pembaca pertama).
func (wb *workbook) lembar(nama string, lebar int) ([][]string, error) {
	var sh struct {
		Rows []struct {
			C []struct {
				R  string `xml:"r,attr"`
				T  string `xml:"t,attr"`
				V  string `xml:"v"`
				IS struct {
					T string `xml:"t"`
				} `xml:"is"`
			} `xml:"c"`
		} `xml:"sheetData>row"`
	}
	if err := bacaXML(wb.berkas[wb.lembarK[nama]], &sh); err != nil {
		return nil, fmt.Errorf("lembar %s: %w", nama, err)
	}
	var hasil [][]string
	for _, row := range sh.Rows {
		r := make([]string, lebar)
		for _, c := range row.C {
			i := indeksKolom(c.R)
			if i < 0 || i >= lebar {
				continue
			}
			switch c.T {
			case "s":
				n, err := strconv.Atoi(c.V)
				if err != nil || n < 0 || n >= len(wb.bersama) {
					return nil, fmt.Errorf("lembar %s sel %s: rujukan teks bersama %q", nama, c.R, c.V)
				}
				r[i] = wb.bersama[n]
			case "inlineStr":
				r[i] = c.IS.T
			default:
				r[i] = c.V
			}
			r[i] = strings.TrimSpace(r[i])
		}
		hasil = append(hasil, r)
	}
	return hasil, nil
}

func indeksKolom(ref string) int {
	n := 0
	for _, ch := range ref {
		if ch < 'A' || ch > 'Z' {
			break
		}
		n = n*26 + int(ch-'A'+1)
	}
	return n - 1
}
