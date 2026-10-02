package services

// `Generate` b60122 → `GenerateUpload_Act` b60228 (paket 9, PARITAS §3.3):
//
//	1 b226 `·` Page-Remove `TempWorkpage1`
//	2 b332 `·` satu baris `pxResults` - `CARI1..CARI40` dari halaman `ProductName`
//	           dan `ProductNameInward` di clipboard (isi form saat itu, tanpa simpan)
//	3 b1285 `·` `pxConvertResultsToCSV` `FileName=SeeDetail` b1335,
//	           `CSVPropHeaders` b1340 (39 judul), `CSVProperties` b1343 (40 properti)
//
// ⚠️ OQ-MPNL-06: 39 judul untuk 40 properti, dan properti ke-12 tertulis `CARI2`
// (bukan `CARI12`) - pengisian menurut posisi menggeser `BENEFIT`…`TREATYNUMBER`.
// Bawaan: tiap kolom berisi nilai yang NAMANYA disebut judulnya; `CARI36`
// (`LIENCLAUSE`) tanpa judul tidak ikut. Ekspresi `@if` ditiru persis, termasuk
// ejaan `SUPRLUS`. Nilai tidak diubah (tanpa pembersih rumus spreadsheet) - sama
// dengan Pega.

import (
	"context"
	"encoding/csv"
	"io"
	"strconv"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/masterproductnamelife/backend/models"
	"nusantarare/modul/masterproductnamelife/backend/repository"
)

// NamaBerkasGenerate - `FileName=SeeDetail` b1335.
const NamaBerkasGenerate = "SeeDetail.csv"

type kolomCSV struct {
	judul  string // `CSVPropHeaders` b1340 VERBATIM
	inward bool   // halaman `ProductNameInward`, bukan `ProductName`
	kunci  string
	ubah   func(string) string
}

// kodeKeTeks - `@if(X==1,"A",@if(X==2,"B",…,"lainnya"))`.
func kodeKeTeks(nilai []string, lainnya string) func(string) string {
	return func(v string) string {
		for i, n := range nilai {
			if strings.TrimSpace(v) == strconv.Itoa(i+1) {
				return n
			}
		}
		return lainnya
	}
}

func umumCSV(k string) kolomCSV   { return kolomCSV{judul: k, kunci: k} }
func inwardCSV(k string) kolomCSV { return kolomCSV{judul: k, inward: true, kunci: k} }

// kolomGenerate - urutan judul b1340; sumber tiap judul dari langkah 2 b332.
var kolomGenerate = []kolomCSV{
	{judul: "TYPE", kunci: "TYPE", ubah: kodeKeTeks([]string{"Basic"}, "Rider")},              // CARI1
	{judul: "TYPE_CEDING", kunci: "TYPE_CEDING", ubah: kodeKeTeks([]string{"QS"}, "SUPRLUS")}, // CARI2
	umumCSV("CEDING"), // CARI3
	{judul: "GRUP", kunci: "GRUP", ubah: kodeKeTeks([]string{"KREDIT LIFE", "NON KREDIT LIFE", "PA"}, "HEALTH")}, // CARI4
	umumCSV("PRODUCTNAME"), umumCSV("PRODUCTCODE"), umumCSV("PRODUCTTYPE"), umumCSV("RIRATE"), // CARI5-8
	umumCSV("RICOMM"), umumCSV("RIRISK"), umumCSV("INWARDNAME"), umumCSV("BENEFIT"), umumCSV("CAUSE"), // CARI9-13
	inwardCSV("POLICYHODERNAME"), inwardCSV("INSURED"), inwardCSV("ADDENDUMWORD"), inwardCSV("AMANDEMENTSCHD"), // CARI14-17
	inwardCSV("BEGIN"), inwardCSV("MATURE"), inwardCSV("BIRTHDAY"), inwardCSV("CURRENCY"), // CARI18-21
	inwardCSV("EXTRAMORTALITY"), inwardCSV("MAXCONTRACT"), inwardCSV("CEDINGRETENTIONNUM"), inwardCSV("CEDINGLIMIT"), // CARI22-25
	inwardCSV("BROKERAGE"), inwardCSV("MINAGE"), inwardCSV("MAXAGE"), inwardCSV("EXTRAPREMI"), inwardCSV("MONTHS"), // CARI26-30
	inwardCSV("MINSUMINSURED"), inwardCSV("MAXSUMINSURED"), inwardCSV("MAXSUMREASURED"), inwardCSV("RNMSHARE"), // CARI31-34
	inwardCSV("RNMLIMITNUM"), // CARI35 (CARI36 LIENCLAUSE tanpa judul)
	{judul: "PAYMENT", inward: true, kunci: "PAYMENT",
		ubah: kodeKeTeks([]string{"Annual", "Semi Annual", "Quarterly", "Monthly"}, "Single")}, // CARI37
	inwardCSV("PROPORTIONALTABLE"), inwardCSV("SUBJECTTO"), inwardCSV("TREATYNUMBER"), // CARI38-40
}

// Generate menulis `SeeDetail.csv` dari isi form yang dikirim (tidak disimpan,
// tidak divalidasi - langkah 2 b332 tanpa prakondisi).
func (l *Layanan) Generate(_ context.Context, p inti.Pelaku, m models.Produk, ke io.Writer) error {
	if err := inti.WajibIdentitas(p); err != nil {
		return err
	}
	umum, inward := repository.HalamanPega(m)
	judul := make([]string, len(kolomGenerate))
	baris := make([]string, len(kolomGenerate))
	for i, k := range kolomGenerate {
		judul[i] = k.judul
		v := umum[k.kunci]
		if k.inward {
			v = inward[k.kunci]
		}
		if k.ubah != nil {
			v = k.ubah(v)
		}
		baris[i] = v
	}
	w := csv.NewWriter(ke)
	w.UseCRLF = true
	_ = w.Write(judul)
	_ = w.Write(baris)
	w.Flush()
	return w.Error()
}
