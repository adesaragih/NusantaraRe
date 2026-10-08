package services

// Rumus tab Limits proporsional — `Activity/LimitCalculation.xml`.
//
// ⛔ DISALIN dari Activity, langkah demi langkah, dengan SARANG langkahnya
// utuh (`pySteps` bersarang, objek perulangan `pyStepsObjectName`). Bukan
// diturunkan ulang dari ringkasan laporan.
//
//	[1]  keluar bila param.autocalculate == false
//	[2]  bila kindoftreaty=="qs": buang RetentionList, CessionList
//	     bila bukan: lompat ke label `sps`
//	[3]  RetentionPct := 100 - QSPct ; CessionPct := 100 - RetentionPct
//	[4]  untuk tiap .IOOLimitList:
//	       Retention += Value * @divide(RetentionPct,100,4)
//	       Cession   += Value * @divide(CessionPct,100,4)   (Currency, CurrencyID disalin)
//	[5]  keluar
//	[6]  sps: bila kindoftreaty=="surplus": buang IOOLimitList, CessionList
//	[7]  bila add!="man": buang RetentionList
//	[8]  currentcob := .TreatyGroup
//	[9]  untuk tiap TreatyIn.Limits -> tiap .Detail yang TreatyType=="QUOTA SHARE":
//	       cob := .TreatyGroup ; bila cob==currentcob:
//	         TempSurplus += .IOOLimitList ; TempCob += .COBList
//	[10] IOOPct := Surplus * 100 ; CessionPct := Surplus * 100
//	[11] bila add=="man": untuk tiap .RetentionList:
//	       IOOLimit += Value * Surplus ; Cession += Value * Surplus
//	[12] untuk tiap TempSurplus, bila add!="man" && .Value>0:
//	       IOOLimit += Value * Surplus ; Retention += Value ; Cession += Value * Surplus
//	[13–14] bila add!="man" && TempCob(1).ClassOfBusiness != "":
//	       COBList := TempCob
//
// ⚠️ DUA BACAAN YANG DIPILIH, dan keduanya dicatat:
//
//   - Langkah 6 menyimpan DUA ungkapan: `pyStepsPreCondParamsWhen` =
//     `param.kindoftreaty=="surplus"` dan `pyExpression` =
//     `… && .IOOLimit>0`. Yang dijalankan Pega adalah yang pertama; yang
//     kedua sisa pembangun ekspresi. Yang dipakai di sini yang pertama.
//   - `@divide(x,100,4)` dibulatkan SETENGAH-KE-ATAS ke empat tempat.
//     Ekspor tidak menyebut mode pembulatannya.
//
// ⭐ DIUKUR ULANG 6 Oktober 2026 terhadap SELURUH 2.868 detail di 1.080
// kontrak Treaty In (`hitung_limit_db_test.go`):
//
//	Quota Share 1.468   persen 100% · RetentionList 99,5% · CessionList 99,5%
//	Surplus     1.345   persen 100% · ketiga larik 99,6% (519 otomatis + 821 manual)
//
// ⚠️ TETAPI PEMBULATANNYA BELUM TERBUKTI OLEH DATA. Hanya 6 dari 1.468
// detail QS yang persennya lewat empat desimal sesudah dibagi 100 — satu-
// satunya kasus di mana `@divide(…,4)` berpengaruh — dan keenamnya terbelah:
// 1 cocok hanya dengan pembulatan, 1 hanya tanpa, 4 tidak cocok keduanya.
// 99,5% itu datang dari persen bulat. Pembulatan di sini berdiri di atas
// TEKS ekspor, dan mode setengah-ke-atasnya tebakan yang dinyatakan.
//
// ⛔ YANG TIDAK DIBANGUN, dan sebabnya bukan kelalaian: pemicu grid
// (`DetailLimits` sel 75–77 dan 117–119) mengirim `autocalculate = '.Layer'`
// — properti baris yang nilainya bergantung data, sehingga apakah langkah 1
// keluar atau tidak bergantung tafsiran Pega atas teks kosong — dan sel 75
// mengirim `kindoftreaty = 'QS'` HURUF BESAR ke perbandingan `=="qs"`.
// Menebak perilaku keduanya berarti mengarang. Pemicu yang dibangun hanya
// yang tak ambigu: sel 37 `.QSPct` dan sel 42 `.Surplus`, keduanya
// `autocalculate = true`.

import (
	"strings"

	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatyin/backend/models"
)

// NilaiMataUang - satu baris larik `IOOLimitList` / `RetentionList` /
// `CessionList` (kelas `ASM-FW-GISFW-Data-TreatyInTotal`).
//
// Alias tipe models — tab Share Non-Prop memuatnya lewat `KontrakWarisan`,
// dan satu bentuk untuk dua lapisan mencegah keduanya berbeda diam-diam.
type NilaiMataUang = models.NilaiMataUang

// KelasBisnisLimit - satu baris `COBList`.
type KelasBisnisLimit struct {
	ClassOfBusiness   string `json:"ClassOfBusiness"`
	ClassOfBusinessID string `json:"ClassOfBusinessID"`
	TreatyGroup       string `json:"TreatyGroup"`
	TreatyGroupID     string `json:"TreatyGroupID"`
}

// DetailLimit - medan `Limits[].Detail[]` yang `LimitCalculation` baca atau
// tulis. Ejaan JSON-nya ejaan dokumen, supaya layar dapat mengirim simpulnya
// apa adanya.
type DetailLimit struct {
	TreatyGroup   string             `json:"TreatyGroup"`
	TreatyType    string             `json:"TreatyType"`
	QSPct         string             `json:"QSPct"`
	Surplus       string             `json:"Surplus"`
	RetentionPct  string             `json:"RetentionPct"`
	CessionPct    string             `json:"CessionPct"`
	IOOPct        string             `json:"IOOPct"`
	IOOLimitList  []NilaiMataUang    `json:"IOOLimitList"`
	RetentionList []NilaiMataUang    `json:"RetentionList"`
	CessionList   []NilaiMataUang    `json:"CessionList"`
	COBList       []KelasBisnisLimit `json:"COBList"`
}

// MasukanLimit - parameter Activity + simpul yang dihitung + pohonnya.
type MasukanLimit struct {
	// `param.kindoftreaty` — `qs` atau `surplus`.
	Jenis string `json:"jenis"`
	// `param.add` — `man` atau kosong.
	Tambah string `json:"tambah"`
	// `param.autocalculate`.
	Otomatis bool        `json:"otomatis"`
	Detail   DetailLimit `json:"detail"`
	// Seluruh `Detail` di seluruh `TreatyIn.Limits`, URUT pohon — langkah 9
	// menapakinya untuk menemukan Quota Share bergrup sama.
	Pohon []DetailLimit `json:"pohon"`
}

// HitungLimit menjalankan `LimitCalculation` atas satu simpul `Detail`.
//
// ⛔ Pure: nol akses basis data, nol jam, nol acak. Yang dikembalikan simpul
// BARU; masukannya tidak diubah.
func (l *Layanan) HitungLimit(p inti.Pelaku, m MasukanLimit) (DetailLimit, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return DetailLimit{}, err
	}
	return HitungLimit(m), nil
}

// HitungLimit — bentuk tanpa pelaku, untuk uji dan pengukuran ulang.
func HitungLimit(m MasukanLimit) DetailLimit {
	d := salinDetail(m.Detail)
	// [1]
	if !m.Otomatis {
		return d
	}
	// [2]–[5] Quota Share
	if m.Jenis == "qs" {
		d.RetentionList = []NilaiMataUang{}
		d.CessionList = []NilaiMataUang{}
		ret := kurang(seratus, angka(d.QSPct))
		ces := kurang(seratus, ret)
		d.RetentionPct = teks(ret)
		d.CessionPct = teks(ces)
		fRet, fCes := bagiBulat(ret, 100, 4), bagiBulat(ces, 100, 4)
		for _, x := range m.Detail.IOOLimitList {
			v := angka(x.Value)
			d.RetentionList = append(d.RetentionList, NilaiMataUang{Currency: x.Currency, CurrencyID: x.CurrencyID, Value: teks(kali(v, fRet))})
			d.CessionList = append(d.CessionList, NilaiMataUang{Currency: x.Currency, CurrencyID: x.CurrencyID, Value: teks(kali(v, fCes))})
		}
		return d
	}
	// [6] sps
	if m.Jenis == "surplus" {
		d.IOOLimitList = []NilaiMataUang{}
		d.CessionList = []NilaiMataUang{}
	}
	// [7]
	manual := m.Tambah == "man"
	retensiAwal := d.RetentionList
	if !manual {
		d.RetentionList = []NilaiMataUang{}
	}
	// [8]–[9]
	var tempSurplus []NilaiMataUang
	var tempCob []KelasBisnisLimit
	for _, q := range m.Pohon {
		if q.TreatyType != "QUOTA SHARE" || q.TreatyGroup != d.TreatyGroup {
			continue
		}
		tempSurplus = append(tempSurplus, q.IOOLimitList...)
		tempCob = append(tempCob, q.COBList...)
	}
	// [10]
	s := angka(d.Surplus)
	persen := teks(kali(s, seratus))
	d.IOOPct = persen
	d.CessionPct = persen
	// [11]
	if manual {
		for _, x := range retensiAwal {
			v := kali(angka(x.Value), s)
			d.IOOLimitList = append(d.IOOLimitList, NilaiMataUang{Currency: x.Currency, CurrencyID: x.CurrencyID, Value: teks(v)})
			d.CessionList = append(d.CessionList, NilaiMataUang{Currency: x.Currency, CurrencyID: x.CurrencyID, Value: teks(v)})
		}
	}
	// [12]
	if !manual {
		for _, x := range tempSurplus {
			v := angka(x.Value)
			if v.Sign() <= 0 { // `.Value>0`
				continue
			}
			hasil := teks(kali(v, s))
			d.IOOLimitList = append(d.IOOLimitList, NilaiMataUang{Currency: x.Currency, CurrencyID: x.CurrencyID, Value: hasil})
			d.RetentionList = append(d.RetentionList, NilaiMataUang{Currency: x.Currency, CurrencyID: x.CurrencyID, Value: teks(v)})
			d.CessionList = append(d.CessionList, NilaiMataUang{Currency: x.Currency, CurrencyID: x.CurrencyID, Value: hasil})
		}
	}
	// [13]–[14]
	if !manual && len(tempCob) > 0 && tempCob[0].ClassOfBusiness != "" {
		d.COBList = append([]KelasBisnisLimit{}, tempCob...)
	}
	return d
}

// salinDetail menyalin simpul beserta lariknya — masukan tidak boleh berubah.
func salinDetail(d DetailLimit) DetailLimit {
	k := d
	k.IOOLimitList = append([]NilaiMataUang{}, d.IOOLimitList...)
	k.RetentionList = append([]NilaiMataUang{}, d.RetentionList...)
	k.CessionList = append([]NilaiMataUang{}, d.CessionList...)
	k.COBList = append([]KelasBisnisLimit{}, d.COBList...)
	return k
}

// konteksLimit — presisi lebar dan `RoundHalfUp`, idiom yang sama dengan
// pemakai `apd` lain di repo ini. Presisi 34 = desimal128: jauh melampaui
// angka uang mana pun di korpus, jadi perkalian di sini tidak pernah
// membulatkan; satu-satunya pembulatan adalah `Quantize` milik `@divide`.
var konteksLimit = func() *apd.Context {
	c := apd.BaseContext.WithPrecision(34)
	c.Rounding = apd.RoundHalfUp
	return c
}()

// `seratus` dipakai ULANG dari `pemulihan.go` — nilai yang sama, satu tempat.

// angka membaca teks desimal Pega. Kosong atau tak terbaca = 0 — Pega
// memperlakukan properti kosong sebagai nol di dalam aritmetika.
//
// ⭐ Bentuk ketikan Indonesia ikut terbaca (8 Oktober 2026). Kontrak
// 1001855 tersimpan dengan `Limit` "16.000.000.000,00", `AdjRate` "4,183",
// dan `ROLPct` "14,379" — diketik di medan teks biasa. Dulu ketiganya
// terbaca NOL tanpa suara: `Total ROL` tetap 0 sesudah Update Total, dan
// `DetailCalculationROL` berhenti di "limit kosong".
func angka(s string) *apd.Decimal {
	s = strings.TrimSpace(s)
	if s == "" {
		return apd.New(0, 0)
	}
	d, _, err := apd.NewFromString(s)
	if err != nil {
		d, _, err = apd.NewFromString(kabelDariKetikan(s))
		if err != nil {
			return apd.New(0, 0)
		}
	}
	return d
}

// kabelDariKetikan mengubah ketikan Indonesia ke bentuk kabel:
// "16.000.000.000,00" → "16000000000.00", "14,379" → "14.379",
// "16.000.000" → "16000000".
//
// ⛔ HANYA dipanggil sesudah bentuk kabel GAGAL dibaca. Nilai kabel
// ("1.659") tidak pernah ditafsir ulang sebagai ribuan.
func kabelDariKetikan(s string) string {
	if strings.Contains(s, ",") {
		return strings.ReplaceAll(strings.ReplaceAll(s, ".", ""), ",", ".")
	}
	if strings.Count(s, ".") > 1 {
		return strings.ReplaceAll(s, ".", "")
	}
	return s
}

func kali(a, b *apd.Decimal) *apd.Decimal {
	r := new(apd.Decimal)
	_, _ = konteksLimit.Mul(r, a, b)
	return r
}

func kurang(a, b *apd.Decimal) *apd.Decimal {
	r := new(apd.Decimal)
	_, _ = konteksLimit.Sub(r, a, b)
	return r
}

// bagiBulat = `@divide(x, pembagi, skala)`: x/pembagi lalu dibulatkan
// setengah-ke-atas ke `skala` tempat desimal (`Quantize`).
func bagiBulat(x *apd.Decimal, pembagi int64, skala int32) *apd.Decimal {
	r := new(apd.Decimal)
	_, _ = konteksLimit.Quo(r, x, apd.New(pembagi, 0))
	_, _ = konteksLimit.Quantize(r, r, -skala)
	return r
}

// teks menulis bilangan sebagai desimal tanpa nol di ekor dan tanpa
// notasi ilmiah — bentuk yang layar dan pemformat `format.ts` terima.
func teks(d *apd.Decimal) string {
	r := new(apd.Decimal)
	r.Reduce(d)
	s := r.Text('f')
	if s == "-0" {
		return "0"
	}
	return s
}
