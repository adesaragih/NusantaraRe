package services

// Sub-tab Object Item - tiket 39: pemeriksaan item objek dan pilihan Object Item Type / Currency.

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// ErrPilihanItemTanpaDatabase - V_JN_OBJ_ITEM / CURRENCY tidak terbaca. 503.
var ErrPilihanItemTanpaDatabase = errors.New("services: basis data tidak dikonfigurasi, pilihan item objek tidak terbaca")

// polaUnit - SetErrorMessageUnit_Act `.Unit<=0` dan kontrol pxNumber pyDecimalPlaces 0: bilangan bulat > 0.
var polaUnit = regexp.MustCompile(`^[0-9]+$`)

// batasAdjust - ValidateAdjustPct `[terverifikasi]`: `Local.adjpct<60||Local.adjpct>100` atas .PctAdjustOther,
// "%Adjustment can't be less than 60% or more than 100%".
var batasAdjustBawah, batasAdjustAtas = apd.New(60, 0), apd.New(100, 0)

// batasItem - T_PROPERTYITEMLIST.SEQ_NO NUMBER(5).
const batasItem = 99999

// lebarItem - lebar kolom (BYTE) migrasi 188 tiap medan teks item.
var lebarItem = []struct {
	nama  string
	nilai func(models.ItemObjek) string
	n     int
}{
	{"itemTypeId", func(i models.ItemObjek) string { return i.ItemTypeID }, 50},
	{"itemType", func(i models.ItemObjek) string { return i.ItemType }, 50},
	{"note", func(i models.ItemObjek) string { return i.Note }, 500},
	{"propertyYear", func(i models.ItemObjek) string { return i.PropertyYear }, 50},
	{"unit", func(i models.ItemObjek) string { return i.Unit }, 50},
	{"condition", func(i models.ItemObjek) string { return i.Condition }, 500},
	{"currency", func(i models.ItemObjek) string { return i.Currency }, 50},
	{"yearOfPlanting", func(i models.ItemObjek) string { return i.YearOfPlanting }, 50},
	{"noOfTree", func(i models.ItemObjek) string { return i.NoOfTree }, 50},
	{"areaHectar", func(i models.ItemObjek) string { return i.AreaHectar }, 50},
	{"remark", func(i models.ItemObjek) string { return i.Remark }, 500},
}

// DenganPilihanItem memasang pembaca V_JN_OBJ_ITEM dan CURRENCY (tiket 39).
func (s *Service) DenganPilihanItem(p interface {
	repository.PembacaJenisItem
	repository.PembacaMataUang
}) *Service {
	s.jenisItem, s.mataUang = p, p
	return s
}

// DaftarJenisItem - pilihan Object Item Type. Tanpa identitas (pola lookup).
func (s *Service) DaftarJenisItem(ctx context.Context) ([]models.JenisItem, error) {
	if s.jenisItem == nil {
		return nil, ErrPilihanItemTanpaDatabase
	}
	return s.jenisItem.DaftarJenisItem(ctx)
}

// DaftarMataUang - pilihan Currency (tanpa ITL). Tanpa identitas (pola lookup).
func (s *Service) DaftarMataUang(ctx context.Context) ([]string, error) {
	if s.mataUang == nil {
		return nil, ErrPilihanItemTanpaDatabase
	}
	return s.mataUang.DaftarMataUang(ctx)
}

// dalamRentang - desimal sah (desimalSah) di antara bawah..atas, dibandingkan apd (tanpa float); nil = tidak sah.
func dalamRentang(d, bawah, atas *apd.Decimal) bool {
	return d != nil && desimalSah(d) && d.Cmp(bawah) >= 0 && d.Cmp(atas) <= 0
}

// periksaItem - isian item baris objek ke-`n` (tanpa basis data). Mata uang WAJIB (A133: kolom CURRENCY NOT NULL
// K-069, dan aplikasi tidak menulis 'UNKNOWN' - K-012); keanggotaannya di CURRENCY diperiksa periksaMataUang.
// Condition / PctAdjust2 tidak dicocokkan ke daftar (aturan propertinya belum ada).
func periksaItem(n int, item []models.ItemObjek) []string {
	var masalah []string
	if len(item) > batasItem {
		return []string{fmt.Sprintf("baris[%d].items paling banyak %d", n, batasItem)}
	}
	for m, it := range item {
		awal := fmt.Sprintf("baris[%d].items[%d].", n, m)
		if strings.TrimSpace(it.Currency) == "" {
			masalah = append(masalah, awal+"currency wajib diisi")
		}
		if !desimalSah(it.TSI) {
			masalah = append(masalah, awal+"tsi"+pesanDesimal)
		}
		if it.Unit != "" && (!polaUnit.MatchString(it.Unit) || strings.Trim(it.Unit, "0") == "") {
			masalah = append(masalah, awal+"unit harus bilangan bulat > 0")
		}
		if !desimalSah(it.PctAdjust2) {
			masalah = append(masalah, awal+"pctAdjust2"+pesanDesimal)
		}
		if !desimalSah(it.PctAdjustOther) {
			masalah = append(masalah, awal+"pctAdjustOther"+pesanDesimal)
		} else if it.IsAdjustable && !dalamRentang(it.PctAdjustOther, batasAdjustBawah, batasAdjustAtas) {
			masalah = append(masalah, awal+"pctAdjustOther: %Adjustment can't be less than 60% or more than 100%")
		}
		for _, l := range lebarItem {
			if len(l.nilai(it)) > l.n {
				masalah = append(masalah, fmt.Sprintf("%s%s paling banyak %d byte", awal, l.nama, l.n))
			}
		}
		if !desimalSah(it.TotalNetRate) {
			masalah = append(masalah, awal+"totalNetRate"+pesanDesimal)
		}
		// CURRENCY_CODE T_COVERAGELIST (migrasi 193) = VARCHAR2(10) rancangan, diisi mata uang item.
		if len(it.Coverages) > 0 && len(it.Currency) > lebarMataUangCoverage {
			masalah = append(masalah, fmt.Sprintf("%scurrency paling banyak %d byte bila item ber-coverage", awal, lebarMataUangCoverage))
		}
		masalah = append(masalah, periksaCoverage(n, m, it.Coverages)...)
	}
	return masalah
}

// mataUangDi - kode mata uang beserta jalur medannya (untuk pesan 400).
type mataUangDi struct{ jalur, kode string }

// periksaMataUang - setiap mata uang item (tiket 39) dan catatan kerugian (tiket 42) ada di CURRENCY tanpa ITL
// (pilihan RD BrowseCurrency_RD). Butuh basis data hanya bila ada item / catatan.
func (s *Service) periksaMataUang(ctx context.Context, baris []models.ObjekFire) error {
	var masalah []string
	var sah map[string]bool
	for n, o := range baris {
		mataUang := make([]mataUangDi, 0, len(o.Items)+len(o.LossRecords))
		for m, it := range o.Items {
			mataUang = append(mataUang, mataUangDi{fmt.Sprintf("baris[%d].items[%d]", n, m), it.Currency})
		}
		for m, c := range o.LossRecords {
			mataUang = append(mataUang, mataUangDi{fmt.Sprintf("baris[%d].lossRecords[%d]", n, m), c.Currency})
		}
		for _, mu := range mataUang {
			if sah == nil {
				if s.mataUang == nil {
					return ErrPilihanItemTanpaDatabase
				}
				daftar, err := s.mataUang.DaftarMataUang(ctx)
				if err != nil {
					return err
				}
				sah = make(map[string]bool, len(daftar))
				for _, c := range daftar {
					sah[c] = true
				}
			}
			if !sah[mu.kode] {
				masalah = append(masalah, fmt.Sprintf("%s.currency %q tidak ada di daftar mata uang", mu.jalur, mu.kode))
			}
		}
	}
	if len(masalah) > 0 {
		return fmt.Errorf("%w: %s", ErrMasukanObjek, strings.Join(masalah, "; "))
	}
	return nil
}
