package services

// Porsi periode tahap pembayaran - tiket E13 (lanjutan E06).
//
// Untuk apa berkas ini: `Endorsment Fac In/Activity/CountPaymentEdm_Act.xml`
// (`ASM-FW-GISFW-WORK / CountPaymentEdm_Act`) langkah 1, 3, 4, 9, 11 dan blok
// 13 ("Endors Ext Periode"), 14 ("Endors Adj Rate"), 15 ("Endors Adj Periode")
// - bagian yang MENULIS `OfferFacIn.ProrateEDMEnd` / `ProrateStartEDM`, tanggal
// polis ternormalisasi, dan faktor periode (`Local.datedif`, `Local.day`,
// `Local.datedifbefore`) yang dipakai rumus pembayaran blok itu. Rumus
// pembayarannya sendiri (13.2, 14.4, 15.5) belum diport (E13).
//
// Dibaca sesudah: porsiperiode.go.
//
// ⛔ Inilah PENULIS KEDUA porsi periode yang K-048 sebut: activity ini
// dipanggil dari layar (`Section/PaymentCurrencyList.xml`, …), SESUDAH kasus
// lahir, sehingga menimpa nilai langkah 15 `SetValueToEDMWork`. Di dalamnya,
// `CountPremiEDM_DT` (langkah 10) dan `CountPaymentEdmTSIObj_Act` (langkah 11)
// berjalan SEBELUM blok 13-15 - untuk ketiga jenis endorsement di sini, blok ini
// penulis terakhir. `CountPremiEDM_DT` hanya MEMBACA `ProtectSpreading.CARI2`
// (`grep -n ProtectSpreading DataTransform/CountPremiEDM_DT.xml` → satu
// prakondisi, nol penugasan); porsinya belum diport.
//
// `[terverifikasi]` Rekonsiliasi kasus nyata NB-15 (`edm-fire-1.json`:
// StatusBusiness 3, EdmType 4, Type 1 → `IsEdmExtendPeriod`, dinilai lewat
// registry predikat EDM): blok 13 memberi `ProrateEDMEnd = 1` = nilai
// tersimpan (`TestPorsiPembayaranKasusNyata`). ⚠️ Rekonsiliasi ini lemah:
// blok 13 menulis konstanta. `ProrateStartEDM` tersimpan
// (1.00666666666666666667) tidak ditulis blok mana pun di sini; dari tiga
// penulisnya di korpus (`grep -rn -i
// '<PropertiesName>[^<]*ProrateStartEDM\|<pyPropertiesName>[^<]*ProrateStartEDM'`)
// - `SetValueToEDMWork` L5340, `CountPaymentEdm_Act` L6235 (blok 14 saja),
// `CountPremiEDMFacOut_DT` L469 (konstanta 0) - hanya langkah 15
// `SetValueToEDMWork` yang dapat menghasilkannya. ⚠️ `[dugaan]` Lokal `int`
// `EdmToStart` lalu bernilai 151 untuk selisih 150,5 hari: Pega TIDAK memotong
// pecahan hari - cocok dengan pembulatan pecahan hari maupun selisih tanggal
// kalender GMT. Satu kasus tidak membedakan keduanya; bahan pertanyaan
// terbuka, BUKAN dasar kode.
//
// Lokal: `[terverifikasi]` satu deklarasi `day` (`int`) di `pyLocalParameters`;
// korpus mengejanya `Local.day` (3×) dan `local.day` (31×) - `[dugaan]`
// variabel yang sama (kata kunci `Local`/`local` seperti `Param`/`param`).
// `datedif`, `datedifbefore`, `startdate`, `edmdate` juga `int`.

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/uang"
	"nusantarare/modul/endorsmentfacin/backend/models"
	"nusantarare/modul/endorsmentfacin/backend/services/predikat"
)

// ErrEDMDayBukanBulat - `QuotationData.EDMDay` terisi tetapi bukan bilangan
// bulat; penugasannya ke lokal `int` belum terverifikasi. (Kosong sah: 15.2
// mengujinya `local.day==""`.)
var ErrEDMDayBukanBulat = errors.New("endorsement: EDMDay bukan bilangan bulat; penugasan Pega ke lokal int belum terverifikasi")

// PredikatPembayaran - gerbang langkah 1, 9, 11, blok 13-15 (rule `When`
// varian EDM). Isi lewat PredikatPembayaranDari.
type PredikatPembayaran struct {
	IsEDM               bool // langkah 1 (F=6: keluar activity)
	IsEdmExtendPeriod   bool // StatusBusiness=3 ∧ EdmType=4 ∧ Type=1
	IsEdmAdjRate        bool // (StatusBusiness=3 ∧ EdmType=4 ∧ Type=3) ∨ IsEdmAdjTSI
	IsEdmAdjShareCedant bool // StatusBusiness=3 ∧ EdmType=4 ∧ Type=11
	IsEdmAdjPeriod      bool // StatusBusiness=3 ∧ EdmType=4 ∧ Type=6
	IsMarineCargo       bool // jenis bisnis POLIS LAMA (hasil query)
	// Gerbang `CountPaymentEdmTSIObj_Act` (langkah 11).
	IsEdmAdjTSI, IsEdmAddObject, IsEdmAdjRIC          bool // blok 1
	IsEdmAdjInsured, IsEdmAdjRefNo, IsEdmAdjSpreading bool // blok 2
}

// PredikatPembayaranDari - menilai seluruh gerbang lewat registry predikat EDM
// (E01). `k` wajib sudah membawa hasil query jenis bisnis polis lama
// (IsMarineCargo) - bila belum, registry panic (urutan mengikat).
func PredikatPembayaranDari(k predikat.Kasus) (PredikatPembayaran, error) {
	var p PredikatPembayaran
	for nama, tuju := range map[string]*bool{
		"IsEDM": &p.IsEDM, "IsEdmExtendPeriod": &p.IsEdmExtendPeriod, "IsEdmAdjRate": &p.IsEdmAdjRate,
		"IsEdmAdjShareCedant": &p.IsEdmAdjShareCedant, "IsEdmAdjPeriod": &p.IsEdmAdjPeriod,
		"IsMarineCargo": &p.IsMarineCargo, "IsEdmAdjTSI": &p.IsEdmAdjTSI, "IsEdmAddObject": &p.IsEdmAddObject,
		"IsEdmAdjRIC": &p.IsEdmAdjRIC, "IsEdmAdjInsured": &p.IsEdmAdjInsured, "IsEdmAdjRefNo": &p.IsEdmAdjRefNo,
		"IsEdmAdjSpreading": &p.IsEdmAdjSpreading,
	} {
		v, err := predikat.Eval(nama, k)
		if err != nil {
			return PredikatPembayaran{}, err
		}
		*tuju = v
	}
	return p, nil
}

// MasukanPorsiPembayaran - masukan langkah 1-15.
type MasukanPorsiPembayaran struct {
	Predikat PredikatPembayaran
	// FixRate - `ProtectSpreading.CARI2=="FIX RATE"` SEBELUM activity (nilai
	// halaman yang sudah ada; langkah 9 dapat menyalakannya).
	FixRate bool
	// EDMDay - `QuotationData.EDMDay`, dibaca langkah 3 dan 15.1.
	EDMDay string
	// StartDateTime, EndDateTime - `OfferFacIn.PolicyData` (kerja, BUKAN
	// OldData), SEBELUM normalisasi langkah 4.
	StartDateTime, EndDateTime time.Time
	// EdmDate - `QuotationData.EdmDate`; TIDAK dinormalisasi langkah 4.
	EdmDate time.Time
	// CacahMataUang - baris `OfferFacIn.CurrencyList` yang diulang 14.4 dan
	// langkah 3 `CountPaymentEdmTSIObj_Act`.
	CacahMataUang int
	// IsProRate - `OfferFacIn.IsProRate` (models.NilaiShortPeriod membuka 1.3
	// TSIObj).
	IsProRate string
	// EdmType - `QuotationData.EdmType` (K-029).
	EdmType models.JenisEndorsemen
}

// HasilPorsiPembayaran - yang DITULIS activity. Porsi kosong = tidak ditulis
// di sini (nilainya tetap dari penulis sebelumnya).
type HasilPorsiPembayaran struct {
	Ditulis                        bool
	ProrateEDMEnd, ProrateStartEDM uang.Ratio
	// StartDateTime, EndDateTime - `OfferFacIn.PolicyData` sesudah langkah 4
	// (ditulis balik ke halaman kerja).
	StartDateTime, EndDateTime time.Time
	// FixRate - `ProtectSpreading.CARI2=="FIX RATE"` sesudah langkah 9.
	FixRate bool
	// Datedif, DatedifBefore, Day - lokal `int` sesudah blok, untuk rumus
	// pembayaran (E13). DayKosong - `Local.day` bernilai "" (EDMDay kosong,
	// langkah 3) dan tidak ditimpa blok mana pun. Blok 15 tidak menyerahkan
	// `Local.startdate`/`edmdate`: 15.3 (tetap jalan) menyetel keduanya 1.
	Datedif, DatedifBefore, Day int64
	DayKosong                   bool
}

// satu - literal Pega `= 1`, skala 0.
func satu() uang.Ratio { return uang.Ratio{Value: apd.New(1, 0), Scale: 0} }

// HitungPorsiPembayaran - langkah berurutan, gerbang korpus.
func HitungPorsiPembayaran(m MasukanPorsiPembayaran) (HasilPorsiPembayaran, error) {
	var h HasilPorsiPembayaran
	// 1 - IsEDM T=2 F=6: bukan endorsement → keluar activity, tidak menulis.
	if !m.Predikat.IsEDM {
		return h, nil
	}
	// 3 - `Local.day = QuotationData.EDMDay` (lokal int).
	var err error
	if h.Day, h.DayKosong, err = hariEDM(m.EDMDay); err != nil {
		return HasilPorsiPembayaran{}, err
	}
	// 4 - normalisasi tanggal polis (tanpa syarat).
	if h.StartDateTime, err = normalTanggalPolis(m.StartDateTime); err != nil {
		return HasilPorsiPembayaran{}, err
	}
	if h.EndDateTime, err = normalTanggalPolis(m.EndDateTime); err != nil {
		return HasilPorsiPembayaran{}, err
	}
	m.StartDateTime, m.EndDateTime = h.StartDateTime, h.EndDateTime
	// 9 - IsMarineCargo T=5 F=2, IsEdmAdjShareCedant T=2 F=3 (ATAU) →
	// `ProtectSpreading.CARI2 = "FIX RATE"` ("biar g usah dikali prorate").
	if m.Predikat.IsMarineCargo || m.Predikat.IsEdmAdjShareCedant {
		m.FixRate = true
	}
	h.FixRate = m.FixRate

	// 11 - `Call CountPaymentEdmTSIObj_Act`: lokalnya milik activity itu, yang
	// menyeberang hanya `OfferFacIn.ProrateEDMEnd`.
	if err := porsiTSIObj(m, &h); err != nil {
		return HasilPorsiPembayaran{}, err
	}

	// 13 - IsEdmExtendPeriod. 13.1 tanpa syarat di dalam blok.
	if m.Predikat.IsEdmExtendPeriod {
		h.Datedif, h.Day, h.DayKosong = 1, 1, false
		h.ProrateEDMEnd, h.Ditulis = satu(), true
	}

	// 14 - IsEdmAdjRate ATAU IsEdmAdjShareCedant (pola (5,2)(2,3)).
	if m.Predikat.IsEdmAdjRate || m.Predikat.IsEdmAdjShareCedant {
		// 14.1 - `@DateTimeDifference(a,b,"D")` = b − a dalam hari `[dugaan]`
		// (deskripsi langkah: "set date diff enddate-effdate").
		if h.Datedif, err = hariBulat(m.EndDateTime.Sub(m.EdmDate)); err != nil {
			return HasilPorsiPembayaran{}, err
		}
		if h.DatedifBefore, err = hariBulat(m.EdmDate.Sub(m.StartDateTime)); err != nil {
			return HasilPorsiPembayaran{}, err
		}
		if h.Day, err = hariBulat(m.EndDateTime.Sub(m.StartDateTime)); err != nil {
			return HasilPorsiPembayaran{}, err
		}
		h.DayKosong = false
		h.ProrateEDMEnd, h.Ditulis = satu(), true
		// 14.2
		if h.Day == 0 {
			h.Day = 365
		}
		// 14.3 - FIX RATE ATAU marine cargo (pola (5,2)(2,3)).
		if m.FixRate || m.Predikat.IsMarineCargo {
			h.Datedif, h.Day, h.DatedifBefore = 1, 1, 0
			h.ProrateEDMEnd = satu()
		}
		// 14.4 - loop CurrencyList (`pyStepsPreCondition=false`: tanpa syarat);
		// 14.4.5 menulis kedua porsi di setiap baris, nilainya sama.
		if m.CacahMataUang > 0 {
			akhir, err := bagiPorsi(h.Datedif, h.Day)
			if err != nil {
				return HasilPorsiPembayaran{}, err
			}
			awal, err := bagiPorsi(h.DatedifBefore, h.Day)
			if err != nil {
				return HasilPorsiPembayaran{}, err
			}
			h.ProrateEDMEnd, h.ProrateStartEDM = akhir, awal
		}
	}

	// 15 - IsEdmAdjPeriod. 15.1 menghitung startdate/edmdate dari selisih
	// tanggal, tetapi 15.3 (`pyStepsPreCondition=false` → tetap jalan, P-11)
	// menimpa keduanya dengan 1; 15.4 `@Math.divide(edmdate, startdate, 20)` = 1.
	// ⛔ K-046 - porsi Adj Periode SELALU 1, diport apa adanya. Selisih 15.1
	// karena itu tidak berpengaruh dan tidak dihitung.
	if m.Predikat.IsEdmAdjPeriod {
		// 15.1 `local.day = EDMDay`; 15.2 `local.day=="" || local.day==0` → 365.
		d, kosong, err := hariEDM(m.EDMDay)
		if err != nil {
			return HasilPorsiPembayaran{}, err
		}
		if kosong || d == 0 {
			d = 365
		}
		h.Day, h.DayKosong = d, false
		akhir, err := bagiPorsi(1, 1)
		if err != nil {
			return HasilPorsiPembayaran{}, err
		}
		h.ProrateEDMEnd, h.Ditulis = akhir, true
	}
	return h, nil
}

// hariEDM - `QuotationData.EDMDay` ke lokal int: kosong sah (dibawa sebagai
// kosong), teks lain yang bukan bulat → ErrEDMDayBukanBulat.
func hariEDM(s string) (hari int64, kosong bool, err error) {
	if s == "" {
		return 0, true, nil
	}
	d, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, false, fmt.Errorf("%w: %q", ErrEDMDayBukanBulat, s)
	}
	return d, false, nil
}

// normalTanggalPolis - langkah 4:
//
//	@DateTime.FormatDateTime(t,"yyyyMMdd","Asia/Jakarta","in_ID")+"T050000.000 GMT"
//
// = tanggal kalender Jakarta (UTC+7 tetap, `jakarta` validasitanggal.go) pukul
// 05:00 GMT. Tanggal kosong → ErrTanggalPorsiPeriodeKosong (FormatDateTime atas
// kosong belum terverifikasi).
func normalTanggalPolis(t time.Time) (time.Time, error) {
	if t.IsZero() {
		return time.Time{}, ErrTanggalPorsiPeriodeKosong
	}
	j := t.In(jakarta)
	return time.Date(j.Year(), j.Month(), j.Day(), 5, 0, 0, 0, time.UTC), nil
}

// porsiTSIObj - `Activity/CountPaymentEdmTSIObj_Act.xml`, bagian porsi periode.
//
//	1  IsEdmAdjTSI ∨ IsEdmAddObject ∨ IsEdmAdjRIC (pola (5,2)(5,2)(2,3)):
//	   1.1 datedif = @DateTimeDifference(EdmDate, End, "D") + 1
//	       day     = @DateTimeDifference(Start, End, "D")
//	   1.2 day ""/0 → 1
//	   1.3 IsProRate=="ShortPeriod" ∨ FIX RATE ∨ IsMarineCargo → datedif = day = 1, ProrateEDMEnd = 1
//	   1.6 (tanpa syarat) ProrateEDMEnd = @Math.divide(datedif, day, 20)
//	2  IsEdmAdjInsured ∨ IsEdmAdjRefNo ∨ IsEdmAdjSpreading: ProrateEDMEnd = 1
//	3  per baris CurrencyList, EdmType==2 (Batal Prorata):
//	   datedif = … + 1, day = …, ProrateEDMEnd = @Math.divide(datedif, day, 20)
//
// Tanggal polis yang dibaca sudah dinormalisasi langkah 4 pemanggil.
//
// ⛔ K-046 - diport apa adanya: pembilang menghitung hari INKLUSIF (+1),
// penyebut tidak, sehingga porsi bisa melebihi 1 sebanyak 1/day; blok 14
// `CountPaymentEdm_Act` tidak memakai +1. Langkah 3 tidak punya guard day = 0.
func porsiTSIObj(m MasukanPorsiPembayaran, h *HasilPorsiPembayaran) error {
	p := m.Predikat
	selisih := func() (datedif, day int64, err error) {
		if datedif, err = hariBulat(m.EndDateTime.Sub(m.EdmDate)); err != nil {
			return 0, 0, err
		}
		day, err = hariBulat(m.EndDateTime.Sub(m.StartDateTime))
		return datedif + 1, day, err
	}
	if p.IsEdmAdjTSI || p.IsEdmAddObject || p.IsEdmAdjRIC {
		datedif, day, err := selisih()
		if err != nil {
			return err
		}
		if day == 0 {
			day = 1 // 1.2
		}
		if m.IsProRate == models.NilaiShortPeriod || m.FixRate || p.IsMarineCargo {
			datedif, day = 1, 1 // 1.3 (ProrateEDMEnd = 1 ditimpa 1.6)
		}
		porsi, err := bagiPorsi(datedif, day)
		if err != nil {
			return err
		}
		h.ProrateEDMEnd, h.Ditulis = porsi, true
	}
	if p.IsEdmAdjInsured || p.IsEdmAdjRefNo || p.IsEdmAdjSpreading {
		h.ProrateEDMEnd, h.Ditulis = satu(), true
	}
	if m.CacahMataUang > 0 {
		batalProrata, err := samaKode(string(m.EdmType), string(models.EdmBatalProrata))
		if err != nil {
			return err
		}
		if batalProrata {
			datedif, day, err := selisih()
			if err != nil {
				return err
			}
			porsi, err := bagiPorsi(datedif, day) // day 0 → ErrBagiNol
			if err != nil {
				return err
			}
			h.ProrateEDMEnd, h.Ditulis = porsi, true
		}
	}
	return nil
}
