package main

import (
	"archive/zip"
	"bytes"
	"os"
	"reflect"
	"testing"
)

// TestSkemaSamaDenganKeluaranBangkit - penjaga `skema_gen.go`: hasil bangkit ulang
// dari workbook harus sama persis dengan yang di repositori. Workbook tidak ada di
// repositori, jadi uji ini butuh env FLAT_RNM_XLSX (jalur
// Tabel-Flat-Lintas-Siklus.xlsx) dan FLAT_RNM_KANDIDAT_XLSX (jalur
// Tabel-Flat-per-Grup-Bisnis.xlsx), dan dilewati tanpanya.
func TestSkemaSamaDenganKeluaranBangkit(t *testing.T) {
	jalur, kandidat := os.Getenv("FLAT_RNM_XLSX"), os.Getenv("FLAT_RNM_KANDIDAT_XLSX")
	if jalur == "" || kandidat == "" {
		t.Skip("FLAT_RNM_XLSX / FLAT_RNM_KANDIDAT_XLSX kosong: penjaga skema_gen.go dilewati (workbook tidak ada di repositori)")
	}
	baru, _, err := bangkitkan(jalur, kandidat)
	if err != nil {
		t.Fatal(err)
	}
	lama, err := os.ReadFile("../skema_gen.go")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(bytes.ReplaceAll(lama, []byte("\r\n"), []byte("\n")), baru) {
		t.Error("skema_gen.go berbeda dari keluaran bangkit - jalankan ulang pembangkit (lihat main.go)")
	}
}

// TestLembarMenempatkanSelMenurutKolom - uji instrumen dengan jawaban yang diketahui:
// sel kosong di tengah baris TIDAK menggeser sel sesudahnya (cacat pembaca pertama
// 02-10-2026), teks bersama dan teks sebaris sama-sama terbaca.
func TestLembarMenempatkanSelMenurutKolom(t *testing.T) {
	var b bytes.Buffer
	z := zip.NewWriter(&b)
	tulis := func(nama, isi string) {
		w, err := z.Create(nama)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := w.Write([]byte(isi)); err != nil {
			t.Fatal(err)
		}
	}
	tulis("xl/workbook.xml", `<workbook><sheets><sheet name="Uji" sheetId="1" r:id="rId1"/></sheets></workbook>`)
	tulis("xl/_rels/workbook.xml.rels", `<Relationships><Relationship Id="rId1" Target="worksheets/sheet1.xml"/></Relationships>`)
	tulis("xl/sharedStrings.xml", `<sst><si><t>A1</t></si><si><r><t>C</t></r><r><t>1</t></r></si></sst>`)
	tulis("xl/worksheets/sheet1.xml", `<worksheet><sheetData>`+
		`<row r="1"><c r="A1" t="s"><v>0</v></c><c r="C1" t="s"><v>1</v></c><c r="D1" t="inlineStr"><is><t> D1 </t></is></c></row>`+
		`<row r="2"><c r="B2"><v>7</v></c></row></sheetData></worksheet>`)
	if err := z.Close(); err != nil {
		t.Fatal(err)
	}
	wb, err := bukaWorkbook(b.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	got, err := wb.lembar("Uji", 4)
	if err != nil {
		t.Fatal(err)
	}
	mau := [][]string{{"A1", "", "C1", "D1"}, {"", "7", "", ""}}
	if !reflect.DeepEqual(got, mau) {
		t.Errorf("lembar %q, mau %q", got, mau)
	}
}
