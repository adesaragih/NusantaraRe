package services

// Lapis B - nilai lama per baris (properti `*Old`). Tiket E06, E08, E09, E11.
//
// Untuk apa berkas ini: `Endorsment Fac In/Activity/SetOldData.xml`, dipanggil
// `InputAddendumFacIn_PreAct` langkah 28 - SETIAP kali layar endorsement
// dibuka, bukan sekali saat kasus lahir.
//
// Dibaca sesudah: beforeimage.go, lapisc.go.
//
// ⛔ Sumber nilainya BARIS LAPIS A, bukan baris kerja: langkah 4 berjalan di
// halaman `pyWorkPage.OfferFacIn.OldData`, mengulang daftar OldData, dan
// menulis ke `pyWorkPage.OfferFacIn.<daftar>(nomor baris yang sama)`.
// Pasangan baris ditentukan NOMOR URUT (`.pxListSubscript`), bukan kunci.
//
// [terverifikasi] 52 penugasan `*Old` di 14 sub-langkah (7 lini × objek +
// cedant). Dua cara, sepakat: urai pohon langkah (7+2+6+2+4+2+5+2+6+2+8+2+2+2)
// dan cacah bahan spec `08-bahan-spec-before-image.md` §3.
//
// Pola nilai tunggal, 52 dari 52: `@If(<sumber>!="", <sumber>, 0)` - kosong
// menjadi NOL, bukan kosong. Dua penyimpangan guard dipertahankan (K-046),
// satu perbaikan sadar (A.5). Lihat `nilaiAtauNol`, `guardTujuan`.

import (
	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/uang"
	"nusantarare/modul/endorsmentfacin/backend/models"
)

// AmbangLokasiNilaiLama - gerbang 3 `SetOldData`:
// `@Utilities.SizeOfPropertyList(.OfferFacIn.LocationList)>100` → keluar.
// Alasan angka 100: pertanyaan terbuka milik work owner.
const AmbangLokasiNilaiLama = 100

// IsiNilaiLama - lapis B. Mengisi ulang seluruh properti `*Old` dari lapis A.
//
// Tiga gerbang keluar di depan (langkah 1-3) menghasilkan lapis B tidak
// tersentuh tanpa galat - perilaku yang benar, bukan kegagalan (E11).
//
// ⚠️ `[dugaan]` Bila daftar kerja lebih pendek dari daftar OldData, baris
// tujuan dibuat baru dan berisi nilai lama saja. Ini dugaan atas perilaku
// `Property-Set` Pega pada nomor baris yang belum ada; karena loop berjalan
// urut naik, baris baru selalu tepat di ujung daftar. Kasus yang memicunya:
// lini Travel, yang ditangani lapis B tetapi TIDAK punya varian lapis C yang
// menyalin `PersonList`.
func IsiNilaiLama(k *models.KasusEndorsement, p Predikat) error {
	o := &k.OfferFacIn
	// 1 - IsLife → keluar (6). Life memakai jalurnya sendiri (E18).
	if p.IsLife {
		return nil
	}
	// 2 - IsEDM → lanjut (2), selain itu keluar (6).
	if !p.IsEDM {
		return nil
	}
	// 3 - lebih dari 100 lokasi (daftar KERJA) → keluar.
	if len(o.LocationList) > AmbangLokasiNilaiLama {
		return nil
	}
	if o.OldData == nil {
		return nil
	}
	// 4 - `pyStepsPreCondition=false`: tetap jalan tanpa syarat (P-11).
	// Gerbang yang mengikat ada di tiap sub-langkah.
	lama := o.OldData
	if p.IsFire {
		nilaiLamaFire(o, lama)   // 4.1
		nilaiLamaCedant(o, lama) // 4.2
	}
	if p.IsGolfInsurance {
		nilaiLamaGolf(o, lama)   // 4.3
		nilaiLamaCedant(o, lama) // 4.4
	}
	if p.IsAneka {
		nilaiLamaAneka(o, lama)  // 4.5
		nilaiLamaCedant(o, lama) // 4.6
	}
	if p.IsPA {
		nilaiLamaPA(o, lama)     // 4.7
		nilaiLamaCedant(o, lama) // 4.8
	}
	if p.IsMarineCargo {
		nilaiLamaMarineCargo(o, lama) // 4.9
		nilaiLamaCedant(o, lama)      // 4.10
	}
	if p.IsMBU {
		nilaiLamaMBU(o, lama)    // 4.11
		nilaiLamaCedant(o, lama) // 4.12
	}
	if p.IsTravel {
		nilaiLamaTravel(o, lama) // 4.13
		nilaiLamaCedant(o, lama) // 4.14
	}
	return nil
}

// nilaiAtauNol - `@If(.X!="", .X, 0)`: sumber terisi → sumber, kosong → nol
// bermata uang sama dengan sumbernya.
//
// ⚠️ `[pertanyaan terbuka]` Mata uang nol cadangan = mata uang yang dibawa
// nilai sumber, yang bisa kosong. Literal `0` Pega tidak bermata uang, dan
// padanannya belum diputuskan siapa pun (bertalian OQ-037/040/046). Pilihan
// ini tidak menambah mata uang yang tidak ada di data; ia tidak menolak,
// karena menolak berarti setiap baris kosong menggagalkan layar.
func nilaiAtauNol(sumber uang.Money) uang.Money {
	if !sumber.Kosong() {
		return sumber
	}
	return uang.Money{Amount: apd.New(0, 0), Currency: sumber.Currency}
}

// guardTujuanRate - `RateOld = @If(.RateOld!="", .Rate, 0)`.
//
// ⛔ K-046 - DIPORT APA ADANYA, perilaku yang BENAR menurut work owner: guard
// menguji `RateOld` (properti bernama TUJUAN, dibaca di baris OldData), lalu
// mengambil `Rate`. Akibatnya polis lama yang belum pernah di-endorse - baris
// OldData-nya tanpa `RateOld` - selalu menghasilkan `RateOld = 0`. Enam dari
// enam kemunculan (4.1.2, 4.3.2, 4.5.2, 4.7.2, 4.9.2, 4.11.2). Jangan
// "diluruskan": angka selisih rate dan rekonsiliasi paralel run bergantung
// padanya.
func guardTujuanRate(barisLama models.BarisMataUang) uang.Ratio {
	if !barisLama.RateOld.Kosong() {
		return barisLama.Rate
	}
	return uang.Ratio{Value: apd.New(0, 0), Scale: barisLama.Rate.Scale}
}

// guardTujuanPremiumFire - `PremiumOld = @If(.PremiumOld!="", .Premium, 0)`.
//
// ⛔ K-046 - DIPORT APA ADANYA. Hanya di FIRE 4.1.2; 17 kemunculan
// `PremiumOld` lain menguji `.Premium` (`nilaiAtauNol`).
func guardTujuanPremiumFire(barisLama models.BarisMataUang) uang.Money {
	if !barisLama.PremiumOld.Kosong() {
		return barisLama.Premium
	}
	return uang.Money{Amount: apd.New(0, 0), Currency: barisLama.Premium.Currency}
}

// premiNusantaraReAtauNol - `PremiNusantaraReOld` di MC (4.9.1.1) dan MBU
// (4.11.1.1).
//
// ⚠️ A.5 - SATU-SATUNYA PERBAIKAN SADAR di modul before-image, BUKAN port.
// Korpus: `@If(.PremiNusantaraRe!="",.PremiNusantaraRe,"0")` - cadangan berupa
// STRING "0" (50 penugasan lain memakai angka 0). Di sini angka nol, karena
// uang bertipe `uang.Money` dan Money tidak menampung string. Pembanding
// rekonsiliasi wajib menormalkan "0" dan 0 sebelum membandingkan (E09, E22).
func premiNusantaraReAtauNol(sumber uang.Money) uang.Money { return nilaiAtauNol(sumber) }

// nilaiLamaBarisMU - pola tiga properti `TotalTSIPremiGrossList`
// (TSIOld, PremiumOld, RateOld) untuk lini selain FIRE.
func nilaiLamaBarisMU(tujuan *models.BarisMataUang, lama models.BarisMataUang) {
	tujuan.TSIOld = nilaiAtauNol(lama.TSI)
	tujuan.PremiumOld = nilaiAtauNol(lama.Premium)
	tujuan.RateOld = guardTujuanRate(lama)
}

// 4.1 - FIRE.
func nilaiLamaFire(o, lama *models.OfferFacIn) {
	for i, ll := range lama.LocationList {
		lok := barisKe(&o.LocationList, i)
		for j, it := range ll.Property.PropertyItemList {
			tj := barisKe(&lok.Property.PropertyItemList, j)
			tj.TSIObjectItemOld = nilaiAtauNol(it.TSIObjectItem)
			tj.TotalGrossPremiOld = nilaiAtauNol(it.TotalGrossPremi)
			tj.TotalPremiumNusantaraReOld = nilaiAtauNol(it.TotalPremiumNusantaraRe)
		}
		for j, b := range ll.Property.TotalTSIPremiGrossList {
			tj := barisKe(&lok.Property.TotalTSIPremiGrossList, j)
			tj.TSIOld = nilaiAtauNol(b.TSI)
			tj.PremiumOld = guardTujuanPremiumFire(b)
			tj.RateOld = guardTujuanRate(b)
		}
		for j, b := range ll.Property.TotalTSIList {
			barisKe(&lok.Property.TotalTSIList, j).TSIOld = nilaiAtauNol(b.TSI)
		}
	}
}

// 4.2/4.4/4.6/4.8/4.10/4.12/4.14 - cedant, identik di ketujuh lini.
func nilaiLamaCedant(o, lama *models.OfferFacIn) {
	for i, c := range lama.CedingCedantList {
		cd := barisKe(&o.CedingCedantList, i)
		for j, b := range c.CurrencyList {
			tj := barisKe(&cd.CurrencyList, j)
			tj.TSIOld = nilaiAtauNol(b.TSI)
			tj.PremiumOld = nilaiAtauNol(b.Premium)
		}
	}
}

// 4.3 - Golf.
func nilaiLamaGolf(o, lama *models.OfferFacIn) {
	for i, ll := range lama.LocationList {
		lok := barisKe(&o.LocationList, i)
		for j, a := range ll.Property.RiskLocation.AnekaList {
			an := barisKe(&lok.Property.RiskLocation.AnekaList, j)
			an.TSIOld = nilaiAtauNol(a.TSI)
			for c, cv := range a.CoverageList {
				cov := barisKe(&an.CoverageList, c)
				cov.TSIOld = nilaiAtauNol(cv.TSI)
				cov.PremiumOld = nilaiAtauNol(cv.Premium)
			}
		}
		for j, b := range ll.Property.TotalTSIPremiGrossList {
			nilaiLamaBarisMU(barisKe(&lok.Property.TotalTSIPremiGrossList, j), b)
		}
	}
}

// 4.5 - Aneka.
func nilaiLamaAneka(o, lama *models.OfferFacIn) {
	for i, ll := range lama.LocationList {
		lok := barisKe(&o.LocationList, i)
		for j, okp := range ll.Property.RiskLocation.OccupationList {
			tujuanOkp := barisKe(&lok.Property.RiskLocation.OccupationList, j)
			for a, an := range okp.AnekaList {
				barisKe(&tujuanOkp.AnekaList, a).TSIOld = nilaiAtauNol(an.TSI)
			}
		}
		for j, b := range ll.Property.TotalTSIPremiGrossList {
			nilaiLamaBarisMU(barisKe(&lok.Property.TotalTSIPremiGrossList, j), b)
		}
	}
}

// 4.7 - PA.
func nilaiLamaPA(o, lama *models.OfferFacIn) {
	for i, pl := range lama.PersonList {
		p := barisKe(&o.PersonList, i)
		for c, cv := range pl.ASMCoverage {
			cov := barisKe(&p.ASMCoverage, c)
			cov.PremiumOld = nilaiAtauNol(cv.Premium)
			cov.TSIOld = nilaiAtauNol(cv.TSI)
		}
		for j, b := range pl.TotalTSIPremiGrossList {
			nilaiLamaBarisMU(barisKe(&p.TotalTSIPremiGrossList, j), b)
		}
	}
}

// 4.9 - Marine Cargo.
func nilaiLamaMarineCargo(o, lama *models.OfferFacIn) {
	for i, kl := range lama.CargoList {
		kg := barisKe(&o.CargoList, i)
		for c, cv := range kl.CoverageList {
			cov := barisKe(&kg.CoverageList, c)
			cov.PremiumOld = nilaiAtauNol(cv.Premium)
			cov.PremiNusantaraReOld = premiNusantaraReAtauNol(cv.PremiNusantaraRe) // A.5
			cov.TSIOld = nilaiAtauNol(cv.TSI)
		}
		for j, b := range kl.TotalTSIPremiGrossList {
			nilaiLamaBarisMU(barisKe(&kg.TotalTSIPremiGrossList, j), b)
		}
	}
}

// 4.11 - MBU.
func nilaiLamaMBU(o, lama *models.OfferFacIn) {
	for i, vl := range lama.VehicleList {
		v := barisKe(&o.VehicleList, i)
		for c, cv := range vl.CoverageList {
			cov := barisKe(&v.CoverageList, c)
			cov.TSIOld = nilaiAtauNol(cv.TSI)
			cov.PremiumGrossDiscountFleetOld = nilaiAtauNol(cv.PremiumGrossDiscountFleet)
			cov.PremiRpOld = nilaiAtauNol(cv.PremiRp)
			cov.PremiNusantaraReOld = premiNusantaraReAtauNol(cv.PremiNusantaraRe) // A.5
			cov.PremiumOld = nilaiAtauNol(cv.Premium)
		}
		for j, b := range vl.TotalTSIPremiGrossList {
			nilaiLamaBarisMU(barisKe(&v.TotalTSIPremiGrossList, j), b)
		}
	}
}

// 4.13 - Travel. Hanya coverage peserta; tanpa TotalTSIPremiGrossList.
func nilaiLamaTravel(o, lama *models.OfferFacIn) {
	for i, pl := range lama.PersonList {
		p := barisKe(&o.PersonList, i)
		for c, cv := range pl.ASMCoverage {
			cov := barisKe(&p.ASMCoverage, c)
			cov.TSIOld = nilaiAtauNol(cv.TSI)
			cov.PremiumOld = nilaiAtauNol(cv.Premium)
		}
	}
}

// barisKe - baris ke-i (0-based) daftar kerja; dibuat di ujung bila belum
// ada (lihat `[dugaan]` di IsiNilaiLama). Loop pemanggil selalu naik urut,
// sehingga i tidak pernah melompati ujung daftar.
func barisKe[T any](daftar *[]T, i int) *T {
	for len(*daftar) <= i {
		var nol T
		*daftar = append(*daftar, nol)
	}
	return &(*daftar)[i]
}
