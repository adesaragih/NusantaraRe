package services

// Premi endorsement per mata uang - `DataTransform/CountPremiEDM_DT.xml`
// (`ASM-FW-GISFW-WORK / CountPremiEDM_DT`), langkah 10 `CountPaymentEdm_Act`.
// Tiket E14 (sebagian: salah satu dari enam berkas rumus "premi menjadi") dan
// E13 (masukan pembayaran).
//
// Dibaca sesudah: pembayaran.go.
//
// ⛔ Urutan: CountPaymentEdm_Act langkah 4 (normalisasi tanggal polis) dan 9
// (FIX RATE) → langkah 10 (di sini) → 11 (TSIObj). HitungPorsiPembayaran
// menjalankannya di tempat itu bila `LangkahSepuluh` diisi; tanpa itu baris
// masukannya dianggap sudah sesudah langkah 10 (dipakai rekonsiliasi blok 13,
// karena langkah 10 menolak kasus nyata - lihat TestPremiEDMDTKasusNyataDitolak).
//
// `[terverifikasi]` 23 aksi `pyDisabled=true` (`grep -c
// '<pyDisabled>true</pyDisabled>' DataTransform/CountPremiEDM_DT.xml` → 23;
// cara kedua `langkah.py --dt … | grep -c pyDisabled` → 23), termasuk seluruh
// blok EdmType 4 pertama, DILEWATI (K-048 §8.2). `[terverifikasi]` langkah 10
// `PassParameterPage=false`: Param mulai kosong. Param bertahan sepanjang
// transformasi, termasuk antar-iterasi FOR_EACH `[dugaan]` - sama dengan lokal
// activity (A20): baris tanpa pasangan OldData membaca Param terbawa →
// BarisTerbawa; bila baris itu yang PERTAMA, Param lama belum pernah diisi →
// ErrLokalBelumTerisi (OldPPH/OldPPN/OldBrokerage/OldDeduction2 bahkan tidak
// dideklarasikan di `pyParameters`). `Param.x` dan `param.x` satu variabel
// `[dugaan]`. Nilai `= Param.OldPremiProrate` pada FOR_EACH cabang EdmType 1
// adalah sisa `pyPropertiesValue` yang tidak dibaca aksi FOR_EACH `[dugaan]`
// (sisa editor) - diabaikan.
//
// ⛔ `@divide` (BUKAN `@Math.divide`; A37 hanya menyebut yang kedua): mode
// pembulatannya belum terverifikasi → hasil diterima hanya bila EKSAK pada
// skalanya (semua mode setuju), selain itu ErrModeDivideBelumTerverifikasi -
// sikap yang sama dengan paket medis.

import (
	"errors"
	"fmt"
	"time"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/uang"
	"nusantarare/modul/endorsmentfacin/backend/models"
)

// ErrModeDivideBelumTerverifikasi - `@divide(a,b,n)` yang hasilnya tidak eksak
// pada n desimal: hasilnya bergantung mode pembulatan yang belum terverifikasi.
var ErrModeDivideBelumTerverifikasi = errors.New("endorsement: @divide tidak eksak; mode pembulatannya belum terverifikasi")

// MasukanPremiEDMDT - halaman yang dibaca DataTransform.
type MasukanPremiEDMDT struct {
	// StartDateTime, EndDateTime - `OfferFacIn.PolicyData` SESUDAH langkah 4;
	// EdmDate - `QuotationData.EdmDate`.
	StartDateTime, EndDateTime, EdmDate time.Time
	// LamaStart, LamaEnd - `OfferFacIn.OldData.PolicyData` (TIDAK
	// dinormalisasi langkah 4); LamaEdmDate - `OldData.QuotationData.EdmDate`.
	LamaStart, LamaEnd, LamaEdmDate time.Time
	// FixRate - `ProtectSpreading.CARI2=="FIX RATE"` sesudah langkah 9.
	FixRate bool
	EdmType models.JenisEndorsemen
	// RICommision - `.Policy.Payment.RICommision` (persen), dibaca cabang
	// EdmType 1 dan 2 DI LUAR FOR_EACH, jadi titiknya halaman primer
	// DataTransform - `[dugaan]` `pyWorkPage` (Apply-DataTransform langkah 10
	// tanpa obj).
	RICommision uang.Ratio
	// HasilProdKe - lihat MasukanPorsiPembayaran (CheckData.CARI24); dibaca
	// cabang EdmType 4 saja.
	HasilProdKe            string
	MataUang, MataUangLama []BarisBayar
}

// HasilPremiEDMDT - yang ditulis DataTransform.
type HasilPremiEDMDT struct {
	MataUang []BarisBayar
	// ProrateEDMEnd - `pyWorkPage.OfferFacIn.ProrateEDMEnd` (selalu ditulis).
	ProrateEDMEnd uang.Ratio
	BarisTerbawa  []string
}

// HitungPremiEDMDT - aksi aktif DataTransform berurutan.
func HitungPremiEDMDT(m MasukanPremiEDMDT) (HasilPremiEDMDT, error) {
	h := HasilPremiEDMDT{MataUang: salinBaris(m.MataUang)}
	param := lokal{}

	// WHEN A: tanggal mulai SAMA atau tanggal akhir SAMA dengan polis lama.
	// WHEN B: mulai BEDA atau akhir BEDA. `[terverifikasi]` dua WHEN berurutan
	// tanpa OTHERWISE: keduanya dapat benar sekaligus dan B lalu menimpa A -
	// diport apa adanya.
	beda := func(a, b time.Time) (int64, error) {
		if a.IsZero() || b.IsZero() {
			return 0, ErrTanggalPorsiPeriodeKosong
		}
		return hariBulat(b.Sub(a))
	}
	dMulai, err := beda(m.StartDateTime, m.LamaStart)
	if err != nil {
		return HasilPremiEDMDT{}, err
	}
	dAkhir, err := beda(m.EndDateTime, m.LamaEnd)
	if err != nil {
		return HasilPremiEDMDT{}, err
	}
	var startdate, edmdate int64
	if dMulai == 0 || dAkhir == 0 {
		if startdate, err = beda(m.StartDateTime, m.EndDateTime); err != nil {
			return HasilPremiEDMDT{}, err
		}
		// Param.enddate (Start→EdmDate) tidak dibaca aksi aktif mana pun.
		if edmdate, err = beda(m.EdmDate, m.EndDateTime); err != nil {
			return HasilPremiEDMDT{}, err
		}
		edmdate++
	}
	if dMulai != 0 || dAkhir != 0 {
		if startdate, err = beda(m.LamaStart, m.LamaEnd); err != nil {
			return HasilPremiEDMDT{}, err
		}
		if edmdate, err = beda(m.LamaEdmDate, m.LamaEnd); err != nil {
			return HasilPremiEDMDT{}, err
		}
		edmdate++
	}
	if startdate == 0 {
		startdate = 1
	}
	awal, err := bagiPorsi(startdate-edmdate, startdate) // prorateedmstart
	if err != nil {
		return HasilPremiEDMDT{}, err
	}
	akhir, err := bagiPorsi(edmdate, startdate) // prorateedmend
	if err != nil {
		return HasilPremiEDMDT{}, err
	}
	if m.FixRate {
		awal, akhir = uang.Ratio{Value: apd.New(0, 0), Scale: 0}, satu()
	}
	h.ProrateEDMEnd = akhir

	jenis := func(kode models.JenisEndorsemen) (bool, error) { return samaKode(string(m.EdmType), string(kode)) }
	perubahan, err := jenis(models.EdmPerubahan)
	if err != nil {
		return HasilPremiEDMDT{}, err
	}
	sejakSemula, err := jenis(models.EdmBatalSejakSemula)
	if err != nil {
		return HasilPremiEDMDT{}, err
	}
	prorata, err := jenis(models.EdmBatalProrata)
	if err != nil {
		return HasilPremiEDMDT{}, err
	}

	if perubahan {
		if err := dtPerubahan(param, h.MataUang, m.MataUangLama, akhir.Value, m.HasilProdKe, &h.BarisTerbawa); err != nil {
			return HasilPremiEDMDT{}, err
		}
	}
	if sejakSemula {
		h.ProrateEDMEnd = satu()
		if err := dtBatal(param, h.MataUang, m.MataUangLama, m.RICommision, nil, &h.BarisTerbawa); err != nil {
			return HasilPremiEDMDT{}, err
		}
	}
	if prorata {
		h.ProrateEDMEnd = awal
		if err := dtBatal(param, h.MataUang, m.MataUangLama, m.RICommision, awal.Value, &h.BarisTerbawa); err != nil {
			return HasilPremiEDMDT{}, err
		}
	}
	return h, nil
}

// dtPerubahan - WHEN EdmType==4 (aktif).
//
//	EDMNewPremi = actualpremi = netto baru
//	pasangan lama: OldPremiProrate = Premium lama;
//	  CARI24=="0": OldPremi = Premium − Commision − BrokerageFee + PPN + PPh (lama; TANPA Deduction2)
//	  CARI24=="1": OldPremi = EDMPremiMenjadi lama
//	EDMOldPremi = OldPremiProrate, EDMPremiMenjadi = actualpremi, EDMOldPayment = OldPremi, EDMOld* = Param
//	SumTotalPayment = actualpremi×end − OldPremi×end
//	Check=="0": EDMPremiMenjadi = EDMNewPremi; EDMOldPayment = EDMOldPremi = EDMOldCommision = 0;
//	            SumTotalPayment = EDMPremiMenjadi×end − 0×end
//
// ⛔ K-046: `Param.Check` bertipe teks ("0"/"1"); EDMOldPremi diisi Premium
// lama (BUKAN premi prorata walau nama Param-nya "OldPremiProrate").
func dtPerubahan(param lokal, baris, lama []BarisBayar, end *apd.Decimal, hasilProdKe string, terbawa *[]string) error {
	nol := apd.New(0, 0)
	// CARI24 hanya dibaca di dalam pasangan; galatnya hanya berlaku di sana.
	cari24EDM, errCari24 := cari24(hasilProdKe)
	return perBaris(baris, param, func(p pembaca, b *BarisBayar) {
		n, py := b.Name, &b.Payment
		k := p.komponenDari(*b, "")
		net := p.netto(k.premi, k)
		py.EDMNewPremi = uangBaris(net, n)
		p.l["actualpremi"] = net
		lm, ada := cariLama(lama, n)
		if ada {
			lp := lm.Payment
			p.l["OldPremiProrate"] = p.uang(lp.Premium, "OldData Premium", n)
			if errCari24 != nil {
				p.jalan(func() (*apd.Decimal, error) { return nil, errCari24 })
			} else if cari24EDM {
				p.l["OldPremi"] = p.uang(lp.EDMPremiMenjadi, "OldData EDMPremiMenjadi", n)
			} else {
				// `[terverifikasi]` tanpa Deduction2 (K-046).
				p.l["OldPremi"] = p.tambah(p.tambah(p.kurang(p.kurang(p.uang(lp.Premium, "OldData Premium", n),
					p.uang(lp.Commision, "OldData Commision", n)), p.uang(lp.BrokerageFee, "OldData BrokerageFee", n)),
					p.uang(lp.PPN, "OldData PPN", n)), p.uang(lp.PPh, "OldData PPh", n))
			}
			isiParamLama(p, lm)
		} else {
			*terbawa = append(*terbawa, n)
		}
		py.EDMOldPremi = uangBaris(p.lok("OldPremiProrate"), n)
		py.EDMPremiMenjadi = uangBaris(p.lok("actualpremi"), n)
		py.EDMOldPayment = uangBaris(p.lok("OldPremi"), n)
		tulisParamLama(p, b)
		b.SumTotalPayment = uangBaris(p.kurang(p.kali(p.lok("actualpremi"), end), p.kali(p.lok("OldPremi"), end)), n)
		if !ada {
			// Check=="0": tiga medan dinolkan; PPN/PPh/Brokerage/Deduction2 lama
			// TERBAWA (menyimpang dari AC E13 "seluruh nilai sebelum nol").
			py.EDMPremiMenjadi = py.EDMNewPremi
			py.EDMOldPayment, py.EDMOldPremi, py.EDMOldCommision = uangBaris(nol, n), uangBaris(nol, n), uangBaris(nol, n)
			b.SumTotalPayment = uangBaris(p.kurang(p.kali(net, end), p.kali(nol, end)), n)
		}
	})
}

// dtBatal - WHEN EdmType==1 (batal sejak semula, r == nil) dan EdmType==2
// (batal prorata, r = prorateedmstart). Param OldPremi/OldPremiProrate/
// oldcommision DINOLKAN tiap baris; OldPPH/OldPPN/OldBrokerage/OldDeduction2
// TIDAK (terbawa, A20).
//
//	pasangan lama: OldPremi = Premium lama
//	  OldPremiProrate (EdmType 1) = @divide(Premium,1,4) − Premium×@divide(RICommision,100,4) − BrokerageFee + PPN + PPh
//	  OldPremiProrate (EdmType 2) =          Premium     − Premium×@divide(RICommision,100,4) − BrokerageFee + PPN + PPh
//	  (`[terverifikasi]` kedua cabang CARI24 identik → CARI24 tidak dibaca)
//	EdmType 1: EDMNewPremi = 0; EDMPremiMenjadi = netto baru;               SumTotal = Menjadi×1 − OldPayment×1
//	EdmType 2: EDMNewPremi = netto baru; EDMPremiMenjadi = OldPremiProrate×r; SumTotal = Menjadi − OldPayment
//	EDMOldPremi = OldPremi, EDMOldPayment = OldPremiProrate, EDMOld* = Param
//
// ⛔ K-046 `K046_EdmType1_HanyaLife…`/K-029 tidak berlaku di sini: DataTransform
// ini menghitung EdmType 1 untuk SEMUA lini. Deduction2 lama tidak masuk
// OldPremiProrate (diport apa adanya).
func dtBatal(param lokal, baris, lama []BarisBayar, riComm uang.Ratio, r *apd.Decimal, terbawa *[]string) error {
	nol, satuD, seratus := apd.New(0, 0), apd.New(1, 0), apd.New(100, 0)
	sejakSemula := r == nil
	return perBaris(baris, param, func(p pembaca, b *BarisBayar) {
		n, py := b.Name, &b.Payment
		p.l["oldcommision"], p.l["OldPremi"], p.l["OldPremiProrate"] = nol, nol, nol
		if lm, ada := cariLama(lama, n); ada {
			lp := lm.Payment
			premi := p.uang(lp.Premium, "OldData Premium", n)
			p.l["OldPremi"] = premi
			pokok := premi
			if sejakSemula {
				pokok = p.jalan(func() (*apd.Decimal, error) { return bagiDivideEksak(premi, satuD, 4) })
			}
			komisi := p.jalan(func() (*apd.Decimal, error) {
				if riComm.Kosong() {
					return nil, fmt.Errorf("%w: Policy.Payment.RICommision", ErrNilaiPembayaranKosong)
				}
				return bagiDivideEksak(riComm.Value, seratus, 4)
			})
			// `[terverifikasi]` tanpa Deduction2 lama (K-046).
			p.l["OldPremiProrate"] = p.tambah(p.tambah(p.kurang(p.kurang(pokok, p.kali(premi, komisi)),
				p.uang(lp.BrokerageFee, "OldData BrokerageFee", n)), p.uang(lp.PPN, "OldData PPN", n)), p.uang(lp.PPh, "OldData PPh", n))
			isiParamLama(p, lm)
		} else {
			*terbawa = append(*terbawa, n)
		}
		k := p.komponenDari(*b, "")
		net := p.netto(k.premi, k)
		py.EDMOldPremi = uangBaris(p.lok("OldPremi"), n)
		py.EDMOldPayment = uangBaris(p.lok("OldPremiProrate"), n)
		if sejakSemula {
			py.EDMNewPremi = uangBaris(nol, n)
			py.EDMPremiMenjadi = uangBaris(net, n)
		} else {
			py.EDMNewPremi = uangBaris(net, n)
			py.EDMPremiMenjadi = uangBaris(p.kali(p.lok("OldPremiProrate"), r), n)
		}
		tulisParamLama(p, b)
		menjadi, lamaBayar := py.EDMPremiMenjadi.Amount, py.EDMOldPayment.Amount
		if sejakSemula {
			menjadi, lamaBayar = p.kali(menjadi, satuD), p.kali(lamaBayar, satuD)
		}
		b.SumTotalPayment = uangBaris(p.kurang(menjadi, lamaBayar), n)
	})
}

// isiParamLama - Param komponen lama dari baris lama (ketiga cabang).
// `Param.oldcommision` dan `param.oldcommision` satu variabel.
func isiParamLama(p pembaca, lm BarisBayar) {
	lp, n := lm.Payment, lm.Name
	p.l["oldcommision"] = p.uang(lp.Commision, "OldData Commision", n)
	p.l["OldPPH"] = p.uang(lp.PPh, "OldData PPh", n)
	p.l["OldPPN"] = p.uang(lp.PPN, "OldData PPN", n)
	p.l["OldBrokerage"] = p.uang(lp.BrokerageFee, "OldData BrokerageFee", n)
	p.l["OldDeduction2"] = p.uang(lp.Deduction2, "OldData Deduction2", n)
}

// tulisParamLama - EDMOldCommision, PPN, PPh, Brokerage, Deduction2 dari Param.
func tulisParamLama(p pembaca, b *BarisBayar) {
	n, py := b.Name, &b.Payment
	py.EDMOldCommision = uangBaris(p.lok("oldcommision"), n)
	py.EDMOldPPN = uangBaris(p.lok("OldPPN"), n)
	py.EDMOldPPh = uangBaris(p.lok("OldPPH"), n)
	py.EDMOldBrokerage = uangBaris(p.lok("OldBrokerage"), n)
	py.EDMOldDeduction2 = uangBaris(p.lok("OldDeduction2"), n)
}

// bagiDivideEksak - `@divide(a, b, n)`: diterima hanya bila eksak pada n
// desimal (lihat kepala berkas).
func bagiDivideEksak(a, b *apd.Decimal, n int32) (*apd.Decimal, error) {
	hasil, err := bagiHalfUp(a, b, n)
	if err != nil {
		return nil, err
	}
	// Eksak pada n desimal ⇔ hasil × b = a.
	balik, err := kaliEksak(hasil, b)
	if err != nil {
		return nil, err
	}
	if balik.Cmp(a) != 0 {
		return nil, ErrModeDivideBelumTerverifikasi
	}
	return hasil, nil
}
