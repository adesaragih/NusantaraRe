package services

// Selisih pembayaran per mata uang - tiket E13.
//
// Untuk apa berkas ini: bagian PEMBAYARAN `Activity/CountPaymentEdm_Act.xml`
// (`ASM-FW-GISFW-WORK / CountPaymentEdm_Act`) blok 13.2, 14.4, 15.5 dan
// `Activity/CountPaymentEdmTSIObj_Act.xml` blok 1.4-1.5, 2.2, langkah 3. Bagian
// porsi periode activity yang sama ada di porsipembayaran.go, yang memanggil
// fungsi di sini pada urutan korpus.
//
// Dibaca sesudah: porsipembayaran.go.
//
// ⛔ Batas port (A21): nilai baris yang DIBACA di sini - `Premium`,
// `Commision`, `BrokerageFee`, `Deduction2`, `PPh`, `PPN`, dan `EDM*` yang belum
// ditulis activity ini - adalah keadaan SESUDAH langkah 2 (`SetRIComm_Act`:
// komisi/brokerage/pajak dari query master pajak, E17 dilewati, belum diport)
// dan langkah 10 (`CountPremiEDM_DT`, premiedmdt.go).
//
// ⛔ A19 - `[terverifikasi]` `pyLocalParameters` CountPaymentEdm_Act
// mendeklarasikan 7 lokal `double` (comm, newpremi, oldpayment, premiins,
// preminew, prorateend, ricommins); 4 di antaranya dibaca di sini (preminew,
// oldpayment, premiins, ricommins). Lokal TSIObj tidak ada yang `double`. Di
// sini SEMUA aritmetika desimal eksak (CLAUDE.md §7), sehingga hasil dapat
// berbeda dari Pega di digit belakang - terukur di
// `TestPembayaranKasusNyataBlok13`.
//
// ⛔ A20 - lokal Pega bertahan sepanjang activity, termasuk antar-iterasi loop
// `[dugaan]`. Baris mata uang BARU (tanpa pasangan di OldData) membaca nilai
// lokal yang terbawa dari baris sebelumnya - mata uang lain. Diport APA ADANYA
// (K-046) dan ditandai di HasilPorsiPembayaran.BarisTerbawa; lokal yang belum
// pernah diisi → ErrLokalBelumTerisi (nilai awalnya belum terverifikasi) -
// termasuk di 14.4.3 / TSIObj 1.4.3, walau 14.4.4 / 1.4.4 lalu menolkan
// sebagian medannya.

import (
	"errors"
	"fmt"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/uang"
)

var (
	// ErrLokalBelumTerisi - rumus membaca lokal Pega yang belum pernah diisi
	// di activity ini; nilai awal lokal belum terverifikasi.
	ErrLokalBelumTerisi = errors.New("endorsement: lokal/Param Pega dibaca sebelum diisi; nilai awalnya belum terverifikasi")
	// ErrNilaiPembayaranKosong - medan pembayaran yang dibaca rumus kosong.
	ErrNilaiPembayaranKosong = errors.New("endorsement: nilai pembayaran kosong dalam rumus; perilaku Pega atas aritmetika nilai kosong belum terverifikasi")
	// ErrMataUangBaris - nilai uang baris bermata uang lain dari `.Name`
	// barisnya. Penjaga tipe Money; korpus tidak memeriksanya.
	ErrMataUangBaris = errors.New("endorsement: mata uang nilai berbeda dari Name baris CurrencyList")
)

// Cicilan - satu baris `.Policy.Payment.ListInstallment`
// (`ASM-FW-GISFW-Data-Installment`).
type Cicilan struct {
	InstallmentPercentage uang.Ratio
	// Keluaran.
	Premium, RICommision, BrokerageFee, PPN, PPh, Deduction2, PaymentTotal uang.Money
}

// PembayaranBaris - `.Policy.Payment` satu baris mata uang.
type PembayaranBaris struct {
	Premium, Commision, BrokerageFee, Deduction2, PPh, PPN                    uang.Money
	EDMNewPremi, EDMPremiMenjadi, EDMOldPremi, EDMOldPayment, EDMOldCommision uang.Money
	EDMOldBrokerage, EDMOldPPN, EDMOldPPh, EDMOldDeduction2, NetPremium       uang.Money
	ListInstallment                                                           []Cicilan
}

// BarisBayar - satu baris `OfferFacIn.CurrencyList` (atau
// `OfferFacIn.OldData.CurrencyList`). `Name` = kode mata uang `[dugaan]`;
// seluruh nilai uang baris bermata uang `Name`.
type BarisBayar struct {
	Name            string
	Payment         PembayaranBaris
	SumTotalPayment uang.Money
}

// salinBaris - salinan baris dan cicilannya. Nilai desimal dibagi pakai:
// kode modul ini tidak pernah memutasi *apd.Decimal di tempat - setiap
// penulisan mengganti pointer (konvensi yang sama dengan premifire.go).
func salinBaris(src []BarisBayar) []BarisBayar {
	out := make([]BarisBayar, len(src))
	for i, b := range src {
		out[i] = b
		out[i].Payment.ListInstallment = append([]Cicilan(nil), b.Payment.ListInstallment...)
	}
	return out
}

// lokal - variabel `Local.*` satu activity (atau `Param.*` satu
// DataTransform). Kunci = ejaan korpus; Pega mengeja sebagian `Local.x` dan
// `local.x` - satu variabel. nil = belum pernah diisi.
type lokal map[string]*apd.Decimal

func (l lokal) baca(nama string) (*apd.Decimal, error) {
	if d := l[nama]; d != nil {
		return d, nil
	}
	return nil, fmt.Errorf("%w: %s", ErrLokalBelumTerisi, nama)
}

// pembaca - nilai baris dan lokal ke desimal, galat lengket lewat rumus.
type pembaca struct {
	*rumus
	l lokal
}

func (p pembaca) lok(nama string) *apd.Decimal {
	return p.jalan(func() (*apd.Decimal, error) { return p.l.baca(nama) })
}

func (p pembaca) uang(m uang.Money, medan, name string) *apd.Decimal {
	return p.jalan(func() (*apd.Decimal, error) {
		if m.Kosong() {
			return nil, fmt.Errorf("%w: %s", ErrNilaiPembayaranKosong, medan)
		}
		if m.Currency != name {
			return nil, fmt.Errorf("%w: %s", ErrMataUangBaris, medan)
		}
		return m.Amount, nil
	})
}

func (p pembaca) rasio(r uang.Ratio, medan string) *apd.Decimal {
	return p.jalan(func() (*apd.Decimal, error) {
		if r.Kosong() {
			return nil, fmt.Errorf("%w: %s", ErrNilaiPembayaranKosong, medan)
		}
		return r.Value, nil
	})
}

// komponen - enam komponen `.Policy.Payment` yang selalu berjalan bersama.
type komponen struct {
	premi, comm, brok, ded2, pph, ppn *apd.Decimal
}

// komponenDari - komponen satu baris; awalan menamai sumbernya di pesan galat.
func (p pembaca) komponenDari(b BarisBayar, awalan string) komponen {
	py, n := b.Payment, b.Name
	return komponen{
		premi: p.uang(py.Premium, awalan+"Premium", n), comm: p.uang(py.Commision, awalan+"Commision", n),
		brok: p.uang(py.BrokerageFee, awalan+"BrokerageFee", n), ded2: p.uang(py.Deduction2, awalan+"Deduction2", n),
		pph: p.uang(py.PPh, awalan+"PPh", n), ppn: p.uang(py.PPN, awalan+"PPN", n),
	}
}

// edmLama - lima komponen `EDMOld*` baris (premi, komisi, brokerage, pajak,
// potongan) sebagaimana tertulis saat dibaca.
func (p pembaca) edmLama(b BarisBayar) komponen {
	py, n := b.Payment, b.Name
	return komponen{
		premi: p.uang(py.EDMOldPremi, "EDMOldPremi", n), comm: p.uang(py.EDMOldCommision, "EDMOldCommision", n),
		brok: p.uang(py.EDMOldBrokerage, "EDMOldBrokerage", n), ded2: p.uang(py.EDMOldDeduction2, "EDMOldDeduction2", n),
		pph: p.uang(py.EDMOldPPh, "EDMOldPPh", n), ppn: p.uang(py.EDMOldPPN, "EDMOldPPN", n),
	}
}

// netto - `premi − Commision − BrokerageFee − Deduction2 + PPh + PPN`.
// ⛔ K-046 `K046_NettoPremi_PPhPPNDitambahkan`: PPh dan PPN DITAMBAHKAN,
// diport apa adanya. premi terpisah karena sebagian rumus memakai lokal
// `preminew`, bukan medan baris.
func (p pembaca) netto(premi *apd.Decimal, k komponen) *apd.Decimal {
	return p.tambah(p.tambah(p.kurang(p.kurang(p.kurang(premi, k.comm), k.brok), k.ded2), k.pph), k.ppn)
}

// cariLama - pasangan `OldData.CurrencyList` ber-Name sama; loop korpus
// menulis lokal di SETIAP baris yang cocok, jadi yang terakhir menang.
func cariLama(lama []BarisBayar, name string) (BarisBayar, bool) {
	var hasil BarisBayar
	ada := false
	for _, b := range lama {
		if b.Name == name {
			hasil, ada = b, true
		}
	}
	return hasil, ada
}

// isiLokalLama - `Set Old Premi & Comm` (13.2.2.1 / 14.4.2.1 / 15.5.2.1 /
// TSIObj 1.4.2.1): enam komponen lama + `oldpayment` dari baris lama.
func isiLokalLama(p pembaca, lm BarisBayar) {
	k := p.komponenDari(lm, "OldData ")
	p.l["oldpremi"], p.l["oldcomm"], p.l["OldBrokeragefee"] = k.premi, k.comm, k.brok
	p.l["OldDeduction2"], p.l["Oldpph"], p.l["Oldppn"] = k.ded2, k.pph, k.ppn
	p.l["oldpayment"] = p.netto(k.premi, k)
}

// prologBaris - awal setiap iterasi blok: `Local.preminew = .Premium`, lalu
// pasangan OldData (isi lokal lama) atau penanda terbawa. Mengembalikan
// komponen baris baru dan apakah pasangan ada (`Local.check`).
func prologBaris(p pembaca, b BarisBayar, lama []BarisBayar, terbawa *[]string) (komponen, bool) {
	k := p.komponenDari(b, "")
	p.l["preminew"] = k.premi
	lm, ada := cariLama(lama, b.Name)
	if ada {
		isiLokalLama(p, lm)
	} else {
		*terbawa = append(*terbawa, b.Name)
	}
	return k, ada
}

// uangBaris - desimal → Money bermata uang baris.
func uangBaris(d *apd.Decimal, name string) uang.Money { return uang.Money{Amount: d, Currency: name} }

// tulisLama - `Set Edm Premi dkk` bagian komponen lama (13.2.3 / 14.4.3 /
// 15.5.3 / TSIObj 1.4.3). denganDed2 - `[terverifikasi]` 14.4.3 dan TSIObj
// 1.4.3 TIDAK menulis `EDMOldDeduction2` (diport apa adanya).
func tulisLama(p pembaca, b *BarisBayar, denganDed2 bool) {
	n, py := b.Name, &b.Payment
	py.EDMOldPremi = uangBaris(p.lok("oldpremi"), n)
	py.EDMOldPayment = uangBaris(p.lok("oldpayment"), n)
	py.EDMOldCommision = uangBaris(p.lok("oldcomm"), n)
	py.EDMOldBrokerage = uangBaris(p.lok("OldBrokeragefee"), n)
	py.EDMOldPPN = uangBaris(p.lok("Oldppn"), n)
	py.EDMOldPPh = uangBaris(p.lok("Oldpph"), n)
	if denganDed2 {
		py.EDMOldDeduction2 = uangBaris(p.lok("OldDeduction2"), n)
	}
}

// sumTotal - `.SumTotalPayment = EDMPremiMenjadi − EDMOldPayment`.
func sumTotal(p pembaca, b *BarisBayar) {
	n := b.Name
	b.SumTotalPayment = uangBaris(p.kurang(p.uang(b.Payment.EDMPremiMenjadi, "EDMPremiMenjadi", n),
		p.uang(b.Payment.EDMOldPayment, "EDMOldPayment", n)), n)
}

// nolkanLama - "kalau currency baru" (14.4.4 / TSIObj 1.4.4): EDMOldPremi,
// EDMOldPayment, EDMOldCommision = 0, dan EDMOldDeduction2 = 0 bila denganDed2
// (TSIObj saja).
func nolkanLama(b *BarisBayar, denganDed2 bool) {
	n, py, nol := b.Name, &b.Payment, apd.New(0, 0)
	py.EDMOldPremi, py.EDMOldPayment, py.EDMOldCommision = uangBaris(nol, n), uangBaris(nol, n), uangBaris(nol, n)
	if denganDed2 {
		py.EDMOldDeduction2 = uangBaris(nol, n)
	}
}

// selisihKomponen - `(X × rBaru) [+ (lama × rLama)] − lama`; rLama nil = tanpa
// bagian lama.
func selisihKomponen(p pembaca, baru, lama, rBaru, rLama *apd.Decimal) *apd.Decimal {
	v := p.kali(baru, rBaru)
	if rLama != nil {
		v = p.tambah(v, p.kali(lama, rLama))
	}
	return p.kurang(v, lama)
}

// isiSelisih - enam lokal `Xins` yang dibaca cicilan. rLamaDed2 terpisah:
// `[terverifikasi]` 14.4.6 `Deduction2ins` tanpa bagian lama (K-046).
func isiSelisih(p pembaca, baru, lama komponen, rBaru, rLama, rLamaDed2 *apd.Decimal) {
	p.l["premiins"] = selisihKomponen(p, baru.premi, lama.premi, rBaru, rLama)
	p.l["ricommins"] = selisihKomponen(p, baru.comm, lama.comm, rBaru, rLama)
	p.l["Brokeragefeeins"] = selisihKomponen(p, baru.brok, lama.brok, rBaru, rLama)
	p.l["ppnins"] = selisihKomponen(p, baru.ppn, lama.ppn, rBaru, rLama)
	p.l["pphins"] = selisihKomponen(p, baru.pph, lama.pph, rBaru, rLama)
	p.l["Deduction2ins"] = selisihKomponen(p, baru.ded2, lama.ded2, rBaru, rLamaDed2)
}

// cicilan - `@Math.divide((Local.Xins * .InstallmentPercentage), 100, 20)` per
// komponen, lalu `PaymentTotal`. denganDed2 - CountPaymentEdm_Act mengurangkan
// `Deduction2` dari `PaymentTotal`; ⛔ K-046 `[terverifikasi]` TSIObj 3.5.1
// TIDAK (diport apa adanya).
func cicilan(p pembaca, b *BarisBayar, denganDed2 bool) {
	n := b.Name
	seratus := apd.New(100, 0)
	for i := range b.Payment.ListInstallment {
		c := &b.Payment.ListInstallment[i]
		pct := p.rasio(c.InstallmentPercentage, "InstallmentPercentage")
		bagian := func(lok string) *apd.Decimal { return p.bagi(p.kali(p.lok(lok), pct), seratus, 20) }
		premi, comm, brok := bagian("premiins"), bagian("ricommins"), bagian("Brokeragefeeins")
		ppn, pph, ded2 := bagian("ppnins"), bagian("pphins"), bagian("Deduction2ins")
		total := p.kurang(p.kurang(premi, comm), brok)
		if denganDed2 {
			total = p.kurang(total, ded2)
		}
		total = p.tambah(p.tambah(total, pph), ppn)
		c.Premium, c.RICommision, c.BrokerageFee = uangBaris(premi, n), uangBaris(comm, n), uangBaris(brok, n)
		c.PPN, c.PPh, c.Deduction2, c.PaymentTotal = uangBaris(ppn, n), uangBaris(pph, n), uangBaris(ded2, n), uangBaris(total, n)
	}
}

// perBaris - kerangka loop baris dengan galat lengket per baris.
func perBaris(baris []BarisBayar, l lokal, f func(p pembaca, b *BarisBayar)) error {
	p := pembaca{rumus: &rumus{}, l: l}
	for i := range baris {
		f(p, &baris[i])
		if p.err != nil {
			return fmt.Errorf("baris %s: %w", baris[i].Name, p.err)
		}
	}
	return nil
}

// bayarBlok13 - `CountPaymentEdm_Act` 13.2 (Extend Period). 13.2 dan 13.2.4
// ber-`pyStepsPreCondition=false` → tetap jalan (P-11): semua baris, semua
// lini. r = @Math.divide(datedif, day, 20) = 1/1 (13.1). Selisih terhadap
// medan `EDMOld*` yang baru ditulis 13.2.3.
func bayarBlok13(baris, lama []BarisBayar, l lokal, r *apd.Decimal, terbawa *[]string) error {
	return perBaris(baris, l, func(p pembaca, b *BarisBayar) {
		k := satuPorsiAwal(p, b, lama, r, terbawa)
		isiSelisih(p, k, p.edmLama(*b), r, nil, nil)
		cicilan(p, b, true)
	})
}

// bayarBlok15 - `CountPaymentEdm_Act` 15.5 (Adj Periode), r =
// @Math.divide(edmdate, startdate, 20) = 1 (15.3). Sama dengan 13.2 kecuali
// 15.5.4/15.5.5 memakai LOKAL lama (bukan medan `EDMOld*`) - nilainya sama
// karena 15.5.3 baru saja menyalinnya - dan premi baru dari `Local.preminew`.
func bayarBlok15(baris, lama []BarisBayar, l lokal, r *apd.Decimal, terbawa *[]string) error {
	return perBaris(baris, l, func(p pembaca, b *BarisBayar) {
		k := satuPorsiAwal(p, b, lama, r, terbawa)
		k.premi = p.lok("preminew")
		lamaLok := komponen{
			premi: p.lok("oldpremi"), comm: p.lok("oldcomm"), brok: p.lok("OldBrokeragefee"),
			ded2: p.lok("OldDeduction2"), pph: p.lok("Oldpph"), ppn: p.lok("Oldppn"),
		}
		isiSelisih(p, k, lamaLok, r, nil, nil)
		cicilan(p, b, true)
	})
}

// satuPorsiAwal - bagian bersama 13.2 dan 15.5: prolog, lalu `Set Edm Premi
// dkk` (EDMNewPremi = netto, EDMPremiMenjadi = netto × r, komponen lama
// TERMASUK EDMOldDeduction2, SumTotalPayment).
func satuPorsiAwal(p pembaca, b *BarisBayar, lama []BarisBayar, r *apd.Decimal, terbawa *[]string) komponen {
	k, _ := prologBaris(p, *b, lama, terbawa)
	net := p.netto(p.lok("preminew"), k)
	b.Payment.EDMNewPremi = uangBaris(net, b.Name)
	b.Payment.EDMPremiMenjadi = uangBaris(p.kali(net, r), b.Name)
	tulisLama(p, b, true)
	sumTotal(p, b)
	return k
}

// bayarBlok14 - `CountPaymentEdm_Act` 14.4 (Adj Rate / Adj Share Cedant).
// rAkhir = @Math.divide(datedif, day, 20), rAwal = @Math.divide(datedifbefore,
// day, 20). Mengembalikan IsEDMShare (14.4.7) bila ditulis 1.
//
// ⛔ Diport apa adanya (K-046), `[terverifikasi]` pohon langkah 14.4:
//   - 14.4.4 mata uang BARU (check==0) hanya menolkan EDMOldPremi/Payment/
//     Commision dan menghitung ulang EDMPremiMenjadi = (preminew − Commision) ×
//     rAkhir - TANPA brokerage/pajak; EDMOldBrokerage/PPN/PPh dari 14.4.3 tetap
//     (nilai lokal terbawa, A20). Tiket E13 menulis "seluruh nilai sebelum
//     dipaksa nol" - korpus tidak.
//   - 14.4.3 tidak menulis EDMOldDeduction2; 14.4.6 `Deduction2ins` tanpa
//     bagian lama.
func bayarBlok14(baris, lama []BarisBayar, l lokal, rAkhir, rAwal *apd.Decimal, share shareEDM, terbawa *[]string) (*bool, error) {
	var isEDMShare *bool
	err := perBaris(baris, l, func(p pembaca, b *BarisBayar) {
		n := b.Name
		// 14.4.1 (`pyStepsPreCondition=false`) - 14.4.2
		k, check := prologBaris(p, *b, lama, terbawa)
		// 14.4.3
		net := p.netto(p.lok("preminew"), k)
		b.Payment.EDMNewPremi = uangBaris(net, n)
		b.Payment.EDMPremiMenjadi = uangBaris(p.tambah(p.kali(net, rAkhir), p.kali(p.lok("oldpayment"), rAwal)), n)
		tulisLama(p, b, false)
		sumTotal(p, b)
		// 14.4.4 kalau currency baru
		if !check {
			b.Payment.EDMPremiMenjadi = uangBaris(p.kali(p.kurang(p.lok("preminew"), k.comm), rAkhir), n)
			nolkanLama(b, false)
			sumTotal(p, b)
		}
		// 14.4.5 / 14.4.6 - terhadap medan EDMOld* saat itu.
		isiSelisih(p, k, p.edmLama(*b), rAkhir, rAwal, nil)
		// 14.4.7 isedm share: `NewShare>OldShare || NewShare<OldShare`.
		if beda, ok := share.beda(p); ok && beda {
			v := true
			isEDMShare = &v
		}
		// 14.4.8 cicilan
		cicilan(p, b, true)
	})
	if err != nil {
		return nil, err
	}
	return isEDMShare, nil
}

// shareEDM - `Local.NewShare` / `Local.OldShare` (14.1; lokal `Decimal`).
type shareEDM struct{ baru, lama uang.Ratio }

func (s shareEDM) beda(p pembaca) (bool, bool) {
	baru, lama := p.rasio(s.baru, "OfferFacIn.PercentShare"), p.rasio(s.lama, "OldData.PercentShare")
	if p.err != nil {
		return false, false
	}
	return baru.Cmp(lama) != 0, true
}

// bayarTSIObjBlok1 - `CountPaymentEdmTSIObj_Act` 1.4 (Adj TSI / Add Object /
// Adj RIC), r = @Math.divide(datedif, day, 20) sesudah 1.1-1.3.
//
// ⛔ K-046 `[terverifikasi]`: 1.4.1 menolkan `oldpremi`/`oldcomm` di awal
// SETIAP baris, tetapi komponen lama lain terbawa (A20); 1.4.3 tidak menulis
// EDMOldDeduction2; 1.4.4 (mata uang baru) menolkan EDMOldPremi/Payment/
// Commision/Deduction2 tanpa menghitung ulang EDMPremiMenjadi. 1.4.3 juga
// menulis `CheckData.CARI1`/`CARI2`/`CARI3` - TIDAK diport (grep korpus
// Endorsment Fac In: hanya muncul di berkas ini).
func bayarTSIObjBlok1(baris, lama []BarisBayar, l lokal, r *apd.Decimal, terbawa *[]string) error {
	return perBaris(baris, l, func(p pembaca, b *BarisBayar) {
		nol := apd.New(0, 0)
		l["oldpremi"], l["oldcomm"] = nol, nol // 1.4.1
		k, check := prologBaris(p, *b, lama, terbawa)
		net := p.netto(p.lok("preminew"), k)
		b.Payment.EDMNewPremi = uangBaris(net, b.Name)
		b.Payment.EDMPremiMenjadi = uangBaris(p.kali(net, r), b.Name)
		tulisLama(p, b, false)
		sumTotal(p, b)
		if !check {
			nolkanLama(b, true)
			sumTotal(p, b)
		}
	})
}

// bayarTSIObjBlok2 - `CountPaymentEdmTSIObj_Act` 2.2 (Adj Insured / Ref No /
// Adj Spreading): `oldpayment` = Premium − Commision baris lama bila
// `CheckData.CARI24==0` ("OLD NB"), atau EDMPremiMenjadi baris lama bila
// `==1` ("OLD EDM"); lalu EDMPremiMenjadi = EDMOldPayment = oldpayment.
func bayarTSIObjBlok2(baris, lama []BarisBayar, l lokal, cari24EDM bool, terbawa *[]string) error {
	return perBaris(baris, l, func(p pembaca, b *BarisBayar) {
		n := b.Name
		if lm, ada := cariLama(lama, n); ada {
			if cari24EDM {
				l["oldpayment"] = p.uang(lm.Payment.EDMPremiMenjadi, "OldData EDMPremiMenjadi", n)
			} else {
				l["oldpayment"] = p.kurang(p.uang(lm.Payment.Premium, "OldData Premium", n), p.uang(lm.Payment.Commision, "OldData Commision", n))
			}
		} else {
			*terbawa = append(*terbawa, n)
		}
		b.Payment.EDMPremiMenjadi = uangBaris(p.lok("oldpayment"), n)
		b.Payment.EDMOldPayment = uangBaris(p.lok("oldpayment"), n)
		sumTotal(p, b)
	})
}

// bayarTSIObjLangkah3 - `CountPaymentEdmTSIObj_Act` langkah 3 ("hitung
// installment"), TANPA gerbang: semua baris, semua jenis endorsement.
// prorateend = `OfferFacIn.ProrateEDMEnd` saat itu (3.2).
//
//	3.3 ListInstallment kosong → salinan ListInstallment baris PERTAMA
//	3.4 NetPremium = EDMPremiMenjadi×p − EDMOldPayment×p; Xins = X×p − EDMOldX×p
//	3.5 cicilan; ⛔ K-046 PaymentTotal TANPA −Deduction2
//
// ⚠️ Langkah ini membaca SEMUA `EDMOld*` baris - untuk Extend Period / Adj
// Periode nilainya datang dari `CountPremiEDM_DT` (premiedmdt.go).
func bayarTSIObjLangkah3(baris []BarisBayar, l lokal, prorateend *apd.Decimal) error {
	return perBaris(baris, l, func(p pembaca, b *BarisBayar) {
		py := &b.Payment
		if len(py.ListInstallment) < 1 {
			py.ListInstallment = append([]Cicilan(nil), baris[0].Payment.ListInstallment...)
		}
		menjadi := p.uang(py.EDMPremiMenjadi, "EDMPremiMenjadi", b.Name)
		bayarLama := p.uang(py.EDMOldPayment, "EDMOldPayment", b.Name)
		py.NetPremium = uangBaris(p.kurang(p.kali(menjadi, prorateend), p.kali(bayarLama, prorateend)), b.Name)
		// 3.4 `Xins = (.X × prorateend) − (.EDMOldX × prorateend)`.
		baru, lama := p.komponenDari(*b, ""), p.edmLama(*b)
		kurangProrata := func(x, y *apd.Decimal) *apd.Decimal { return p.kurang(p.kali(x, prorateend), p.kali(y, prorateend)) }
		l["premiins"] = kurangProrata(baru.premi, lama.premi)
		l["ricommins"] = kurangProrata(baru.comm, lama.comm)
		l["ppnins"] = kurangProrata(baru.ppn, lama.ppn)
		l["pphins"] = kurangProrata(baru.pph, lama.pph)
		l["Brokeragefeeins"] = kurangProrata(baru.brok, lama.brok)
		l["Deduction2ins"] = kurangProrata(baru.ded2, lama.ded2)
		cicilan(p, b, false)
	})
}
