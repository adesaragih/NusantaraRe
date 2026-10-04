// Package tiruan adalah Gudang Aggregate di memori - untuk uji services dan handlers TANPA Oracle.
//
// Perilakunya meniru SQL repository: daftar dikelompokkan menurut kunci GridDasbordAgg, cari di Ceding Name dan
// Ceding Code, Master ID menurut CEDING memuat kata (huruf besar) dan saringan jenis, periode TREATYYEAR dan kurs
// baris pertama. Seluruh data uji berawalan `UJI`.
package tiruan

import (
	"context"
	"sort"
	"strconv"
	"strings"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/aggregate/backend/models"
	"nusantarare/modul/aggregate/backend/repository"
)

// Periode - satu baris TREATYYEAR (`yyyyMMdd`).
type Periode struct{ Tahun, Mulai, Akhir string }

// Gudang menyimpan AGGREGATE dan tabel acuan di memori.
type Gudang struct {
	Baris  []models.Baris
	Treaty []models.MasterTreaty
	Zona   map[string]string
	Tahun  []Periode
	// Kurs - "tahun|mata uang" -> {TOUSD, TOIDR}.
	TabelKurs map[string][2]string
	// Nomor - nilai SEQ_AGGREGATE berikutnya.
	Nomor int64
	// Waktu - TANGGAL_INPUT unggahan (`YYYY-MM-DD HH24:MI:SS`).
	Waktu string
}

// Transaksi tiruan - fn(nil).
func Transaksi(_ context.Context, fn func(tx *db.Tx) error) error { return fn(nil) }

func hariInput(b models.Baris) string {
	t, err := time.Parse("02-01-2006 15:04:05", b[models.KolomTanggalInput])
	if err != nil {
		return ""
	}
	return t.Format("02-01-2006")
}

func kunciBaris(b models.Baris) models.Kunci {
	return models.Kunci{TanggalInput: hariInput(b), CedingCode: b["CEDING_CODE"], CedingName: b["CEDING_NAME"],
		TreatyType: b["TREATY_TYPE"], AsAt: b["AS_AT"], UwYear: b["UW_YEAR"]}
}

// Daftar - kelompok menurut kunci, terbaru dulu.
func (g *Gudang) Daftar(_ context.Context, kueri string, offset, ukuran int) ([]models.Kelompok, int, error) {
	k := strings.ToUpper(kueri)
	per := map[models.Kunci]*models.Kelompok{}
	var urut []models.Kunci
	for _, b := range g.Baris {
		if k != "" && !strings.Contains(strings.ToUpper(b["CEDING_NAME"]), k) && !strings.Contains(strings.ToUpper(b["CEDING_CODE"]), k) {
			continue
		}
		kb := kunciBaris(b)
		if per[kb] == nil {
			per[kb] = &models.Kelompok{Kunci: kb}
			urut = append(urut, kb)
		}
		per[kb].JumlahBaris++
		if t := b[models.KolomTanggalInput]; len(t) >= 16 && t[:16] > per[kb].InputTerakhir {
			per[kb].InputTerakhir = t[:16]
		}
	}
	sort.SliceStable(urut, func(i, j int) bool { return per[urut[i]].InputTerakhir > per[urut[j]].InputTerakhir })
	out := []models.Kelompok{}
	for i := offset; i < len(urut) && i < offset+ukuran; i++ {
		out = append(out, *per[urut[i]])
	}
	return out, len(urut), nil
}

// Rincian - baris berkunci k.
func (g *Gudang) Rincian(_ context.Context, k models.Kunci) ([]models.Baris, error) {
	out := []models.Baris{}
	for _, b := range g.Baris {
		if kunciBaris(b) == k {
			out = append(out, b)
		}
	}
	return out, nil
}

// Hapus - buang baris berkunci k.
func (g *Gudang) Hapus(_ context.Context, _ *db.Tx, k models.Kunci) (int64, error) {
	var sisa []models.Baris
	var n int64
	for _, b := range g.Baris {
		if kunciBaris(b) == k {
			n++
			continue
		}
		sisa = append(sisa, b)
	}
	g.Baris = sisa
	return n, nil
}

// Ringkasan - jumlah RNM_VALUE_IN_USD per CEDING, TREATY_TYPE, COVERAGE; asAt kosong = seluruh As At (tiruan: bilangan bulat).
func (g *Gudang) Ringkasan(_ context.Context, asAt string) ([]models.IrisanRingkasan, error) {
	jumlah := map[models.IrisanRingkasan]int64{}
	for _, b := range g.Baris {
		if asAt == "" || b["AS_AT"] == asAt {
			n, _ := strconv.ParseInt(b["RNM_VALUE_IN_USD"], 10, 64)
			k := models.IrisanRingkasan{CedingCode: b["CEDING_CODE"], CedingName: b["CEDING_NAME"],
				TreatyType: b["TREATY_TYPE"], Coverage: b["COVERAGE"]}
			jumlah[k] += n
		}
	}
	out := []models.IrisanRingkasan{}
	for k, n := range jumlah {
		k.RnmValueInUSD = strconv.FormatInt(n, 10)
		out = append(out, k)
	}
	sort.Slice(out, func(i, j int) bool {
		a, b := out[i], out[j]
		return a.CedingCode+"|"+a.CedingName+"|"+a.TreatyType+"|"+a.Coverage < b.CedingCode+"|"+b.CedingName+"|"+b.TreatyType+"|"+b.Coverage
	})
	return out, nil
}

// CariTreaty - CEDING memuat kata (huruf besar) dan (Proportional + PROPERTY atau NonProportional).
func (g *Gudang) CariTreaty(_ context.Context, kueri string, batas int) ([]models.MasterTreaty, error) {
	out := []models.MasterTreaty{}
	for _, t := range g.Treaty {
		jenisSah := (t.ProportionType == repository.ProporsionalPropertyJenis && t.TreatyGroup == repository.ProporsionalPropertyGrup) ||
			t.ProportionType == repository.NonProporsional
		if jenisSah && strings.Contains(t.Ceding, strings.ToUpper(kueri)) && len(out) < batas {
			out = append(out, t)
		}
	}
	return out, nil
}

// AmbilTreaty - satu baris view menurut ID.
func (g *Gudang) AmbilTreaty(_ context.Context, id string) (models.MasterTreaty, error) {
	for _, t := range g.Treaty {
		if t.ID == id {
			return t, nil
		}
	}
	return models.MasterTreaty{}, repository.ErrTidakAda
}

// DaftarZona - kode -> catatan.
func (g *Gudang) DaftarZona(context.Context) (map[string]string, error) { return g.Zona, nil }

// TahunTreaty - periode pertama yang memuat asAt.
func (g *Gudang) TahunTreaty(_ context.Context, asAt string) (string, error) {
	for _, p := range g.Tahun {
		if p.Mulai <= asAt && p.Akhir >= asAt {
			return p.Tahun, nil
		}
	}
	return "", nil
}

// Kurs - TOUSD dan TOIDR.
func (g *Gudang) Kurs(_ context.Context, tahun, mataUang string) (string, string, bool, error) {
	v, ada := g.TabelKurs[tahun+"|"+mataUang]
	return v[0], v[1], ada, nil
}

// Sekarang - waktu tiruan.
func (g *Gudang) Sekarang(context.Context, *db.Tx) (string, error) { return g.Waktu, nil }

// NomorBerikut - pencacah tiruan.
func (g *Gudang) NomorBerikut(context.Context, *db.Tx) (int64, error) {
	n := g.Nomor
	g.Nomor++
	return n, nil
}

// AdaID - ID sudah ada di Baris.
func (g *Gudang) AdaID(_ context.Context, _ *db.Tx, id string) (bool, error) {
	for _, b := range g.Baris {
		if b[models.KolomID] == id {
			return true, nil
		}
	}
	return false, nil
}

// Sisip - tambah baris dengan ID dan jejak.
func (g *Gudang) Sisip(_ context.Context, _ *db.Tx, id, sekarang, pelaku string, b models.Baris) error {
	salinan := models.Baris{}
	for k, v := range b {
		salinan[k] = v
	}
	t, _ := time.Parse("2006-01-02 15:04:05", sekarang)
	salinan[models.KolomID], salinan[models.KolomTanggalInput], salinan[models.KolomUserInput] = id, t.Format("02-01-2006 15:04:05"), pelaku
	g.Baris = append(g.Baris, salinan)
	return nil
}

// Contoh - gudang berisi data acuan uji: dua zona, periode 2023 dan 2024 (bertumpuk pada 31-08-2024), kurs USD dan
// EUR, tiga baris Master ID (satu kembar, satu Proportional non-PROPERTY), dan satu baris AGGREGATE lama yang ID-nya
// sama dengan nomor sequence berikutnya.
func Contoh() *Gudang {
	return &Gudang{
		Zona:  map[string]string{"1.1": "1.1 UJI Zona Satu", "3.1": "3.1 UJI Zona Tiga"},
		Tahun: []Periode{{"2023", "20230901", "20240901"}, {"2024", "20240701", "20250930"}},
		TabelKurs: map[string][2]string{
			"2024|USD": {"1", "16500.00"}, "2024|EUR": {"1.10235", "18954.96"}, "2023|USD": {"1", "15000.00"},
		},
		Treaty: []models.MasterTreaty{
			{ID: "UJI-V1", TreatyID: "1000001", CedingID: "UJI-C1", Ceding: "UJI CEDING SATU", RnmShare: "10",
				TreatyYear: "2024", ProportionType: "NonProportional", TreatyGroup: "PROPERTY"},
			{ID: "UJI-V2", TreatyID: "1000002", CedingID: "UJI-C2", Ceding: "UJI CEDING DUA", RnmShare: ".5",
				TreatyYear: "2024", ProportionType: "Proportional", TreatyGroup: "PROPERTY"},
			{ID: "UJI-V3", TreatyID: "1000001", CedingID: "UJI-C1", Ceding: "UJI CEDING SATU", RnmShare: "10",
				TreatyYear: "2024", ProportionType: "NonProportional", TreatyGroup: "PROPERTY"},
			{ID: "UJI-V4", TreatyID: "1000004", CedingID: "UJI-C1", Ceding: "UJI CEDING SATU", RnmShare: "5",
				TreatyYear: "2024", ProportionType: "Proportional", TreatyGroup: "ENGINEERING"},
		},
		Baris: []models.Baris{{models.KolomID: "AGG-500", models.KolomTanggalInput: "01-03-2025 10:00:00",
			"CEDING_CODE": "UJI-C9", "CEDING_NAME": "UJI CEDING LAMA", "TREATY_TYPE": "OR", "AS_AT": "30-09-2024",
			"UW_YEAR": "2023", "COVERAGE": "EQVET", "RNM_VALUE_IN_USD": "7"}},
		Nomor: 500,
		Waktu: "2026-10-04 09:30:00",
	}
}
