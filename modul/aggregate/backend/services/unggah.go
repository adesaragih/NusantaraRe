package services

// Upload CSV - versi Go `UploadCSVAggregate_Act` (local action `UploadCSV_Aggregate`, perintah work owner
// 04-10-2026: "cuman local action kamu buat versi go"). Langkah Pega yang ditandai `//` (Agent, Obj-Browse Master
// Treaty otomatis, buang ganda) tidak pernah jalan dan tidak dipindahkan; langkah 7.11 (COMMENCEMENT 20230101)
// mati karena COMMENCEMENT tidak pernah diisi.

import (
	"context"
	"encoding/csv"
	"errors"
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/utils"
	"nusantarare/modul/aggregate/backend/models"
	"nusantarare/modul/aggregate/backend/repository"
)

// Batas unggahan.
const (
	// MaksBarisCSV - baris data CSV paling banyak sekali unggah.
	MaksBarisCSV = 5000
	// MaksMasterTreaty - Master ID paling banyak sekali unggah.
	MaksMasterTreaty = 50
	// SkalaKursIDR - `@divide(1, .TO_USD, 20)` langkah 7.16.
	SkalaKursIDR = 20
	// SkalaRNMValue - `@divide((.TOTAL_IN_AMOUNT*.RNM_SHARE), 100, 4)` langkah 7.24.
	SkalaRNMValue = 4
)

// PermintaanPratinjau - isi berkas CSV (teks) dan Master ID terpilih (ID baris view, urutan pilihan).
type PermintaanPratinjau struct {
	CSV          string   `json:"csv"`
	MasterTreaty []string `json:"masterTreaty"`
}

// Pratinjau - grid `TempCSV`: baris data lalu satu baris "Total :" per mata uang.
type Pratinjau struct {
	Baris []models.Baris `json:"baris"`
	// MasterTreaty - Master ID terpilih yang dipakai, dibaca ulang dari view.
	MasterTreaty []models.MasterTreaty `json:"masterTreaty"`
}

// AngkaCSV - `@replaceAll(@replaceAll(.COLUMNn, ".", ""), ",", ".")`: titik ribuan dibuang, koma menjadi titik.
func AngkaCSV(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, ".", ""), ",", ".")
}

// AsAtCSV - `@substring(.COLUMN4,6,10)+@substring(.COLUMN4,3,5)+@substring(.COLUMN4,0,2)`: `dd/MM/yyyy` ->
// tanggal; dikembalikan dalam bentuk kabel `DD-MM-YYYY`. Bentuk lain dikembalikan apa adanya (bukan tanggal sah).
func AsAtCSV(s string) string {
	if len(s) < 10 {
		return s
	}
	t, err := time.Parse("20060102", s[6:10]+s[3:5]+s[0:2])
	if err != nil {
		return s
	}
	return t.Format("02-01-2006")
}

// asAtPega - `DD-MM-YYYY` -> `yyyyMMdd` (bentuk yang dibandingkan dengan STARTDATE/ENDDATE TREATYYEAR); "" = bukan
// tanggal sah.
func asAtPega(s string) string {
	t, err := time.Parse("02-01-2006", s)
	if err != nil {
		return ""
	}
	return t.Format("20060102")
}

// BacaCSV - baris data CSV (pemisah `;`). Baris pertama = kepala kolom (dilewati, seperti `MODUploadCSVResults`);
// baris yang seluruh selnya kosong dilewati. Kolom dibaca menurut nomor urut, bukan nama kepala.
func BacaCSV(teks string) ([][]string, error) {
	teks = strings.TrimPrefix(teks, string(rune(0xFEFF)))
	r := csv.NewReader(strings.NewReader(teks))
	r.Comma = ';'
	r.LazyQuotes = true
	r.FieldsPerRecord = -1
	semua, err := r.ReadAll()
	if err != nil {
		var pe *csv.ParseError
		if errors.As(err, &pe) {
			return nil, tolak("the CSV file cannot be read at line %d", pe.Line)
		}
		return nil, tolak("the CSV file cannot be read")
	}
	var out [][]string
	for i, rek := range semua {
		if i == 0 {
			continue
		}
		kosong := true
		for _, sel := range rek {
			kosong = kosong && strings.TrimSpace(sel) == ""
		}
		if !kosong {
			out = append(out, rek)
		}
	}
	if len(out) == 0 {
		return nil, tolak("the CSV file has no data rows")
	}
	if len(out) > MaksBarisCSV {
		return nil, tolak("the CSV file has %d data rows; at most %d rows can be uploaded at once", len(out), MaksBarisCSV)
	}
	return out, nil
}

// BarisCSV - langkah 6: COLUMN1-COLUMN38 ke kolom AGGREGATE.
func BarisCSV(rek []string) models.Baris {
	b := models.Baris{}
	for _, k := range models.KolomGrid {
		v := ""
		if k.CSV > 0 && k.CSV <= len(rek) {
			v = strings.TrimSpace(rek[k.CSV-1])
		}
		switch {
		case k.CSV == 0:
		case k.Nama == "TREATY_TYPE":
			v = strings.ToUpper(v)
		case k.Jenis == models.Tanggal:
			v = AsAtCSV(v)
		case k.Jenis == models.Angka && k.Nama != "TO_USD":
			v = AngkaCSV(v)
		}
		b[k.Nama] = v
	}
	return b
}

// desimal - teks -> desimal; kosong atau tak terurai = nil.
func desimal(s string) *apd.Decimal {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	d, err := utils.ParseDecimal(strings.TrimSpace(s))
	if err != nil {
		return nil
	}
	return d
}

// nolBila - nilai kosong dihitung nol, seperti aritmetika Pega atas properti desimal kosong.
func nolBila(d *apd.Decimal) *apd.Decimal {
	if d == nil {
		return apd.New(0, 0)
	}
	return d
}

func konteks() *apd.Context {
	c := utils.DecimalContext()
	c.Rounding = apd.RoundHalfUp
	return c
}

// teksAngka - desimal -> teks kabel tanpa nol buntut.
func teksAngka(d *apd.Decimal) string {
	if d == nil {
		return ""
	}
	r := new(apd.Decimal)
	r.Reduce(d)
	return utils.FormatDecimal(r)
}

// normalAngka - teks angka dalam bentuk kabel (mis. `.5` dari TO_CHAR Oracle -> `0.5`); tak terurai = apa adanya.
func normalAngka(s string) string {
	if d := desimal(s); d != nil {
		return teksAngka(d)
	}
	return s
}

func kali(a, b *apd.Decimal) *apd.Decimal {
	r := new(apd.Decimal)
	_, _ = konteks().Mul(r, nolBila(a), nolBila(b))
	return r
}

// bagi - a / b dibulatkan setengah-ke-atas pada skala; b nol = nil.
func bagi(a, b *apd.Decimal, skala int32) *apd.Decimal {
	b = nolBila(b)
	if b.IsZero() {
		return nil
	}
	c := konteks()
	r := new(apd.Decimal)
	if _, err := c.Quo(r, nolBila(a), b); err != nil {
		return nil
	}
	if _, err := c.Quantize(r, r, -skala); err != nil {
		return nil
	}
	return r
}

// rnmValue - langkah 7.24/7.25: (TOTAL_IN_AMOUNT x share) / 100, 4 desimal.
func rnmValue(tia, share *apd.Decimal) *apd.Decimal {
	return bagi(kali(tia, share), apd.New(100, 0), SkalaRNMValue)
}

// rujukan - isi acuan satu unggahan: zona, Master ID, dan cache tahun treaty serta kurs.
type rujukan struct {
	zona    map[string]string
	treaty  []models.MasterTreaty
	tahun   map[string]string
	kurs    map[string][3]string // tahun|mata uang -> TOUSD, TOIDR, ada
	gudang  Gudang
	konteks context.Context
}

func (r *rujukan) tahunTreaty(asAt string) (string, error) {
	if t, ada := r.tahun[asAt]; ada {
		return t, nil
	}
	t, err := r.gudang.TahunTreaty(r.konteks, asAt)
	if err != nil {
		return "", err
	}
	r.tahun[asAt] = t
	return t, nil
}

func (r *rujukan) ambilKurs(tahun, mataUang string) (string, string, bool, error) {
	k := tahun + "|" + mataUang
	if v, ada := r.kurs[k]; ada {
		return v[0], v[1], v[2] == "1", nil
	}
	usd, idr, ada, err := r.gudang.Kurs(r.konteks, tahun, mataUang)
	if err != nil {
		return "", "", false, err
	}
	tanda := "0"
	if ada {
		tanda = "1"
	}
	r.kurs[k] = [3]string{usd, idr, tanda}
	return usd, idr, ada, nil
}

// kodeZona - `@substring(.ASSESMENT_ZONE,0,3)`.
func kodeZona(z string) string {
	r := []rune(z)
	if len(r) > 3 {
		r = r[:3]
	}
	return string(r)
}

// lengkapi - langkah 7 untuk satu baris data.
func (r *rujukan) lengkapi(b models.Baris) error {
	// 7.2-7.4 Assessment Zone: kode 3 huruf pertama -> ASSESSMENT_NOTE; tidak ada = kosong.
	if z := b["ASSESMENT_ZONE"]; z != "" {
		b["ASSESMENT_ZONE"] = r.zona[kodeZona(z)]
	}
	// 7.7 Ceding Code dan Ceding Name dari Master ID pertama (isi CSV ditimpa).
	pertama := models.MasterTreaty{}
	if len(r.treaty) > 0 {
		pertama = r.treaty[0]
	}
	b["CEDING_CODE"], b["CEDING_NAME"] = pertama.CedingID, pertama.Ceding
	// 7.8-7.10 Treaty Year: periode TREATYYEAR yang memuat As At.
	b["TREATYYEAR"] = ""
	if a := asAtPega(b["AS_AT"]); a != "" {
		t, err := r.tahunTreaty(a)
		if err != nil {
			return err
		}
		b["TREATYYEAR"] = t
	}
	// 7.12-7.16 To USD: kurs tahun treaty; IDR memakai TOIDR baris USD lalu 1/x; tanpa kurs = kosong.
	b["TO_USD"] = ""
	if tahun := b["TREATYYEAR"]; tahun != "" {
		mataUang := b["CURRENCY"]
		cari := mataUang
		if mataUang == models.MataUangIDR {
			cari = models.MataUangUSD
		}
		usd, idr, ada, err := r.ambilKurs(tahun, cari)
		if err != nil {
			return err
		}
		if ada {
			v := usd
			if mataUang == models.MataUangIDR {
				v = idr
			}
			if d := desimal(v); d != nil {
				v = teksAngka(d)
			}
			b["TO_USD"] = strings.TrimSpace(v)
		}
		if d := desimal(b["TO_USD"]); mataUang == models.MataUangIDR && d != nil && !d.IsZero() {
			b["TO_USD"] = teksAngka(bagi(apd.New(1, 0), d, SkalaKursIDR))
		}
	}
	// 7.24 RNM Share, RNM Value, nilai USD, Master Treaty dari Master ID pertama.
	toUSD, tia := desimal(b["TO_USD"]), desimal(b["TOTAL_IN_AMOUNT"])
	b["RNM_SHARE"], b["M_TREATY_ID"] = normalAngka(pertama.RnmShare), pertama.TreatyID
	nilai := rnmValue(tia, desimal(pertama.RnmShare))
	b["RNM_VALUE"] = teksAngka(nilai)
	b["TOTAL_IN_AMOUNT_IN_USD"] = teksAngka(kali(toUSD, tia))
	b["RNM_VALUE_IN_USD"] = teksAngka(kali(toUSD, nilai))
	// 7.23 + 7.25 lebih dari satu Master ID: share dijumlah, ID digabung `;`, RNM Value dihitung ulang. RNM Value
	// (USD) TIDAK dihitung ulang - langkah 7.25 Pega tidak menyentuhnya (keputusan work owner: ikuti XML).
	if len(r.treaty) > 1 {
		jumlah := apd.New(0, 0)
		ids := make([]string, 0, len(r.treaty))
		for _, t := range r.treaty {
			_, _ = konteks().Add(jumlah, jumlah, nolBila(desimal(t.RnmShare)))
			ids = append(ids, t.TreatyID)
		}
		b["RNM_SHARE"], b["M_TREATY_ID"] = teksAngka(jumlah), strings.Join(ids, ";")
		b["RNM_VALUE"] = teksAngka(rnmValue(tia, jumlah))
	}
	return nil
}

// kolomTotal - kolom yang dijumlah baris "Total :" (langkah 9.1.1).
var kolomTotal = []string{"RNM_VALUE", "TOTAL_NO_OF_RISK", "TOTAL_IN_AMOUNT", "TOTAL_IN_AMOUNT_IN_USD", "RNM_VALUE_IN_USD"}

// BarisTotal - langkah 7.26, 8, 9: satu baris "Total :" per mata uang, urutan kemunculan TERAKHIR mata uang itu
// (langkah Java Pega membuang ganda dari belakang).
func BarisTotal(data []models.Baris) []models.Baris {
	var urut []string
	dilihat := map[string]bool{}
	for i := len(data) - 1; i >= 0; i-- {
		m := data[i]["CURRENCY"]
		if !dilihat[m] {
			dilihat[m] = true
			urut = append([]string{m}, urut...)
		}
	}
	out := make([]models.Baris, 0, len(urut))
	for _, m := range urut {
		total := models.Baris{}
		for _, k := range models.KolomGrid {
			total[k.Nama] = ""
		}
		total["ASSESMENT_ZONE"], total["CURRENCY"] = models.ZonaTotal, m
		for _, kol := range kolomTotal {
			jumlah := apd.New(0, 0)
			for _, b := range data {
				if b["CURRENCY"] == m && b["ASSESMENT_ZONE"] != models.ZonaTotal {
					_, _ = konteks().Add(jumlah, jumlah, nolBila(desimal(b[kol])))
				}
			}
			total[kol] = teksAngka(jumlah)
		}
		out = append(out, total)
	}
	return out
}

// Pratinjau - Upload CSV: baca berkas, lengkapi setiap baris, lalu tambahkan baris "Total :" per mata uang.
func (l *Layanan) Pratinjau(ctx context.Context, p PermintaanPratinjau) (Pratinjau, error) {
	if len(p.MasterTreaty) > MaksMasterTreaty {
		return Pratinjau{}, tolak("at most %d Master IDs can be selected", MaksMasterTreaty)
	}
	rekaman, err := BacaCSV(p.CSV)
	if err != nil {
		return Pratinjau{}, err
	}
	r := &rujukan{tahun: map[string]string{}, kurs: map[string][3]string{}, gudang: l.gudang, konteks: ctx,
		treaty: []models.MasterTreaty{}}
	dilihat := map[string]bool{}
	for _, id := range p.MasterTreaty {
		id = strings.TrimSpace(id)
		if id == "" || dilihat[id] {
			continue
		}
		dilihat[id] = true
		t, err := l.gudang.AmbilTreaty(ctx, id)
		if errors.Is(err, repository.ErrTidakAda) {
			return Pratinjau{}, tolak("Master ID %s is no longer in the treaty list; choose it again", id)
		}
		if err != nil {
			return Pratinjau{}, err
		}
		r.treaty = append(r.treaty, t)
	}
	if r.zona, err = l.gudang.DaftarZona(ctx); err != nil {
		return Pratinjau{}, err
	}
	data := make([]models.Baris, 0, len(rekaman))
	for _, rek := range rekaman {
		b := BarisCSV(rek)
		if err := r.lengkapi(b); err != nil {
			return Pratinjau{}, err
		}
		data = append(data, b)
	}
	return Pratinjau{Baris: append(data, BarisTotal(data)...), MasterTreaty: r.treaty}, nil
}
