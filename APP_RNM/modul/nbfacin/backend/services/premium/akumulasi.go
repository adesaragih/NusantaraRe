package premium

// Akumulasi premi per mata uang (tiket NB-07): premi tiap coverage dibulatkan DI
// DALAM loop (rumus lininya), lalu dijumlahkan - galat pembulatan menumpuk per
// coverage, dan port yang membulatkan sekali di akhir akan meleset.
//
// Asal: `D:\migrasi\RNM\NB FacIn\Activity\FillPremiMBU_FacIn.xml`
// (ASM-FW-GISFW-DATA-COVERAGE!FILLPREMIMBU_FACIN) langkah 2.6 jalur NB:
//
//	2.6   loop CurrencyMaster.pxResults (tabel master mata uang, GetCurrencyMaster)
//	2.6.1   Local.PremiPerCurrency = 0
//	2.6.2.1.1  bila Local.CurrencyMaster==.Currency.Name && .FlagDelete!=1:
//	           Local.PremiPerCurrency += .Premium; Local.currencycoverage = .Currency.Name
//	2.6.3   bila !IsEDM && Local.CurrencyMaster==Local.currencycoverage:
//	        CurrencyList(<APPEND>).Name = .CURR; .Premium = Local.PremiPerCurrency
//
// Tidak diport: jalur EDM (2.6.4, menulis TSI, bukan premi), 2.5 salinan OldData
// untuk StatusBusiness 3. `[terverifikasi]` jalur NB tidak menulis TSI ke CurrencyList.
// `[terverifikasi]` PA tidak punya penulis CurrencyList.Premium (`grep -l` atas
// `NB FacIn\Activity\*.xml`); AddCurrencyListPA_ACT hanya menjumlahkan TSI.

import (
	"fmt"

	"nusantarare/inti/backend/uang"
	"nusantarare/inti/backend/utils"
)

// CoverageMBU - satu baris VehicleList(*).CoverageList: masukan premi (MataUang =
// `.Currency.Name`) dan `.FlagDelete`.
type CoverageMBU struct {
	Input      Input
	FlagDelete string
}

// PremiMataUang - satu baris CurrencyList hasil langkah 2.6.3.
type PremiMataUang struct {
	MataUang string
	Premi    uang.Money
}

// PremiPerMataUang - langkah 2.6 jalur NB. `urutanMaster` = kolom CURR tabel master
// mata uang, urut seperti hasil GetCurrencyMaster (isi tabelnya tidak ada di korpus).
// Mata uang tanpa coverage tidak ditambahkan; coverage bermata uang di luar master
// tidak dijumlahkan.
func PremiPerMataUang(coverage []CoverageMBU, urutanMaster []string) ([]PremiMataUang, error) {
	premi := make([]uang.Money, len(coverage))
	for i, c := range coverage {
		hapus, err := dihapus(c.FlagDelete)
		if err != nil {
			return nil, fmt.Errorf("coverage %d: %w", i+1, err)
		}
		if hapus {
			continue
		}
		// Rumus lini, termasuk pembulatannya - di dalam loop.
		if premi[i], err = Calculate(c.Input); err != nil {
			return nil, fmt.Errorf("coverage %d: %w", i+1, err)
		}
	}
	var hasil []PremiMataUang
	for _, mu := range urutanMaster {
		var jumlah uang.Money
		for i, c := range coverage {
			if premi[i].Kosong() || c.Input.MataUang != mu {
				continue
			}
			if jumlah.Kosong() {
				jumlah = premi[i]
				continue
			}
			var err error
			if jumlah, err = jumlah.Add(premi[i]); err != nil {
				return nil, err
			}
		}
		if !jumlah.Kosong() {
			hasil = append(hasil, PremiMataUang{MataUang: mu, Premi: jumlah})
		}
	}
	return hasil, nil
}

// dihapus - `.FlagDelete!=1`, dibandingkan sebagai angka; kosong = bukan 1.
// Teks yang bukan angka ditolak, tidak ditebak.
func dihapus(flag string) (bool, error) {
	if flag == "" {
		return false, nil
	}
	d, err := utils.ParseDecimal(flag)
	if err != nil {
		return false, fmt.Errorf("%w: FlagDelete %q", ErrAngkaTakTerbaca, flag)
	}
	satu, _ := utils.ParseDecimal("1")
	return d.Cmp(satu) == 0, nil
}
