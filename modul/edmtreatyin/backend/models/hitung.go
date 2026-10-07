package models

// Untuk apa berkas ini: RANTAI PERHITUNGAN UANG layar realisasi - port aktivitas
// `Count*`, `CalculatePremi_Act`, `SetPPNPPH`, `SetDueTo_act` (tiket 07, 13; spec
// §5.6, §8.2; AC 27, 28, 79).
//
// ⭐ AC 79: setiap rumus di berkas ini BERASAL dari `PropertiesValue` satu
// langkah - nama aktivitas dan nomor langkahnya ditulis di atas tiap baris,
// berdampingan dengan `docs/INVENTARIS-XML.md` bab 6. Nol rumus karangan.
//
// Halaman langkah semua aktivitas ini adalah `pyWorkPage.PolicyTreatyIn`
// (kelas `ASM-FW-GISFW-Data-PolicyTreatyIn`), jadi `.X` di rule = jalur
// `PolicyTreatyIn.X` di sini; `pyWorkPage.TreatyIn.Y` = `TreatyIn.Y`.
//
// Semantik syarat langkah Pega yang ditiru (kode aksi di INVENTARIS bab 6):
// syarat dievaluasi berurutan; `lanjut` = periksa syarat berikut, `lewati` =
// langkah tidak dijalankan, `keluar` = aktivitas berhenti TANPA galat.
//
// ⚠️ PESAN `ErrorMsg1`. `Property-Set-Messages` di rule memasang pesan rule
// bernama `ErrorMsg1`; isi rule pesan itu TIDAK ada di korpus. Nama pesannya
// dipasang VERBATIM - tidak dikarang kalimatnya.

import (
	"strings"

	"github.com/cockroachdb/apd/v3"
)

// PesanErrorMsg1 adalah nama rule pesan yang dipasang `Property-Set-Messages`
// rantai Count* - verbatim, isinya tidak terekspor.
const PesanErrorMsg1 = "ErrorMsg1"

// Nilai `ProportionalType` yang dibandingkan rule - VERBATIM.
const (
	JenisProporsional    = "Proportional"
	JenisNonProporsional = "NonProportional"
)

// TypeTaxInclusive adalah satu-satunya nilai `TypeTax` yang ditulis di rule.
// ⛔ Nilai lawannya ("Exclusive") tidak pernah tertulis di mana pun; cabang
// kedua menangkap SEMUA nilai lain, termasuk kosong dan huruf kecil (AC 28).
const TypeTaxInclusive = "Inclusive"

// polis membungkus halaman untuk port aktivitas berhalaman `PolicyTreatyIn`.
type polis struct {
	h *Halaman
	k *kalkulator
}

func (p polis) teks(m string) string { return p.h.Ambil(HalamanPolis + "." + m) }

func (p polis) n(m string) *apd.Decimal { return p.nj(HalamanPolis + "." + m) }

// nj membaca jalur penuh sebagai angka (kosong = 0); galat urai dicatat.
func (p polis) nj(jalur string) *apd.Decimal {
	d, err := p.h.Angka(jalur)
	if err != nil {
		p.k.catat(err)
		return apd.New(0, 0)
	}
	return d
}

func (p polis) setel(m string, d *apd.Decimal) {
	if p.k.err == nil {
		p.h.SetelAngka(HalamanPolis+"."+m, d)
	}
}

func (p polis) pesan(m, isi string) { p.h.TambahPesan(HalamanPolis+"."+m, isi) }

// jenisProporsi membaca `pyWorkPage.PolicyTreatyIn.QuotationData.ProportionalType`.
func (p polis) jenisProporsi() string {
	return p.h.Ambil(HalamanPolis + ".QuotationData.ProportionalType")
}

func (p polis) nonProp() bool { return p.jenisProporsi() == JenisNonProporsional }

var (
	seratus     = apd.New(100, 0)
	seratus5des = apd.New(10000000, -5) // @String.toDecimal("100.00000")
)

// ---------------------------------------------------------------- SetPPNPPH

// JalurStsPKP adalah `ListAgent.pxResults(1).STS_PKP` - status PKP agen
// sumber bisnis (`BrowseClientName_RD`, `Param.ID` =
// `PolicyTreatyIn.QuotationData.SourceOfBusiness`, SetPPNPPH langkah 1-3).
// Pembacanya di repository; services menaruh hasilnya di jalur ini SEBELUM
// rantai hitung berjalan, sehingga port ini tetap murni.
const JalurStsPKP = "ListAgent.STS_PKP"

// SetPPNPPH = `Activity/SetPPNPPH` (dipanggil `CountNetPremi_act` langkah 1).
//
// Langkah 4 BERSYARAT (`pyStepsPreCondition=true`, dua baris syarat):
//
//	.FlagPPH=="true"                     benar -> lewati syarat berikut (5); salah -> lanjut (2)
//	ListAgent.pxResults(1).STS_PKP == 1  salah -> lewati langkah (3)
//
// Jadi langkah 4 berjalan bila FlagPPH "true" ATAU agen sumber bisnis PKP.
// ⛔ RALAT ronde ini: port pertama menganggap langkah 4 tanpa syarat dan
// hasil RD tidak terbaca - keliru; syarat kedua membaca `STS_PKP` hasil RD.
// Bila langkah 4 lewat, BrokerageFee/BrokerageFeeSebenarnya TIDAK disentuh
// (nilai lama tetap), sedangkan PPH/PPN sudah dinolkan langkah 2.
func SetPPNPPH(h *Halaman) error {
	k := &kalkulator{}
	p := polis{h, k}
	// langkah 2
	p.setel("PPHValue", apd.New(0, 0))
	p.setel("PPNValue", apd.New(0, 0))
	// langkah 4 - syarat
	if p.teks("FlagPPH") != "true" && !samaDenganSatu(h.Ambil(JalurStsPKP)) {
		return k.err
	}
	// .BrokerageFee = @divide(2.5,100,8) * (.PremiOgp + .PremiOnp)
	p.setel("BrokerageFee", k.Kali(k.BagiBulat(k.d("2.5"), seratus, 8), k.Tambah(p.n("PremiOgp"), p.n("PremiOnp"))))
	// .BrokerageFeeSebenarnya = @if(.TypeTax=="Inclusive",@divide(.Deduction1,@divide(102.2,100,8),8),.Deduction1)
	sebenarnya := BrokerageSebenarnya(k, p.teks("TypeTax"), p.n("Deduction1"))
	p.setel("BrokerageFeeSebenarnya", sebenarnya)
	// .PPHValue = .BrokerageFeeSebenarnya* @divide(2,100,8)
	p.setel("PPHValue", k.Kali(sebenarnya, k.BagiBulat(apd.New(2, 0), seratus, 8)))
	// .PPNValue = .BrokerageFeeSebenarnya* @divide(2.2,100,8)
	p.setel("PPNValue", k.Kali(sebenarnya, k.BagiBulat(k.d("2.2"), seratus, 8)))
	return k.err
}

// samaDenganSatu meniru `x == 1` Pega atas properti teks: dibandingkan
// sebagai bilangan bila terurai ("1", "1.0"), kosong = tidak sama.
func samaDenganSatu(teks string) bool {
	if !AdalahDesimal(teks) {
		return false
	}
	d, err := AngkaTeks("STS_PKP", teks)
	return err == nil && d.Cmp(apd.New(1, 0)) == 0
}

// BrokerageSebenarnya = `@if(TypeTax=="Inclusive", @divide(x,@divide(102.2,100,8),8), x)`.
//
// ⛔ Pembandingnya TEKS PERSIS (AC 27, 28, spec §10.2 butir 2): "inclusive",
// " Inclusive", dan "" jatuh ke cabang kedua - potongan dipakai apa adanya.
func BrokerageSebenarnya(k *kalkulator, typeTax string, x *apd.Decimal) *apd.Decimal {
	if typeTax == TypeTaxInclusive {
		return k.BagiBulat(x, k.BagiBulat(k.d("102.2"), seratus, 8), 8)
	}
	return x
}

// ---------------------------------------------------------------- CountNetPremi_act

// CountNetPremi = `Activity/CountNetPremi_act`.
func CountNetPremi(h *Halaman) error {
	// langkah 1 Call SetPPNPPH
	if err := SetPPNPPH(h); err != nil {
		return err
	}
	k := &kalkulator{}
	p := polis{h, k}
	// langkah 2 (NonProp -> lewati): .GrossClaim = @divide(.Claim,pyWorkPage.TreatyIn.RNMShareP,4) * 100
	// ⛔ Dilewati bila nilai master tidak tersedia - lihat `MasterTersedia`.
	if !p.nonProp() && MasterTersedia(h, "RNMShareP") {
		p.setel("GrossClaim", k.Kali(k.BagiBulat(p.n("Claim"), p.nj(HalamanMaster+".RNMShareP"), 4), seratus))
	}
	// langkah 3 (hanya NonProp): .GrossClaim = @divide(.Claim,pyWorkPage.TreatyIn.RNMShare,4) * 100
	if p.nonProp() && MasterTersedia(h, "RNMShare") {
		p.setel("GrossClaim", k.Kali(k.BagiBulat(p.n("Claim"), p.nj(HalamanMaster+".RNMShare"), 4), seratus))
	}
	if k.err != nil {
		return k.err
	}
	// langkah 4
	// .NetPremium = ((.PremiOgp-.ResultOgp1)+(.PremiOnp-.ResultOnp1))-(.OveriddingCommOgp*.PremiOgp/100)-(.OveriddingCommOnp*.PremiOnp/100)-.Deduction1-.Deduction2
	net := k.Tambah(k.Kurang(p.n("PremiOgp"), p.n("ResultOgp1")), k.Kurang(p.n("PremiOnp"), p.n("ResultOnp1")))
	net = k.Kurang(net, k.Bagi(k.Kali(p.n("OveriddingCommOgp"), p.n("PremiOgp")), seratus))
	net = k.Kurang(net, k.Bagi(k.Kali(p.n("OveriddingCommOnp"), p.n("PremiOnp")), seratus))
	net = k.Kurang(k.Kurang(net, p.n("Deduction1")), p.n("Deduction2"))
	p.setel("NetPremium", net)
	// langkah 5: .BalanceDueTo = .NetPremium-.Claim-.ExcessLoss+.SalvageValue
	p.setel("BalanceDueTo", k.Tambah(k.Kurang(k.Kurang(net, p.n("Claim")), p.n("ExcessLoss")), p.n("SalvageValue")))
	// langkah 6: syarat `.Deduction1!="" || .Deduction1!=0`
	// ⚠️ Kosong = 0 di sisi kanan, jadi syarat ini setara "Deduction1 terisi".
	if p.teks("Deduction1") != "" || p.n("Deduction1").Sign() != 0 {
		ppn, pph := p.n("PPNValue"), p.n("PPHValue")
		klaim := k.Tambah(k.Kurang(apd.New(0, 0), k.Tambah(p.n("Claim"), p.n("ExcessLoss"))), p.n("SalvageValue"))
		// .BalanceDueTo = .NetPremium +.PPNValue+.PPHValue-.Claim-.ExcessLoss+.SalvageValue
		p.setel("BalanceDueTo", k.Tambah(k.Tambah(k.Tambah(net, ppn), pph), klaim))
		// .BalanceBeforePPH = .NetPremium +.PPNValue-.Claim-.ExcessLoss+.SalvageValue
		p.setel("BalanceBeforePPH", k.Tambah(k.Tambah(net, ppn), klaim))
		// .BalanceBeforeTax = .NetPremium -.Claim-.ExcessLoss+.SalvageValue
		p.setel("BalanceBeforeTax", k.Tambah(net, klaim))
	}
	if k.err != nil {
		return k.err
	}
	// langkah 7 call SetDueTo_act
	return SetDueTo(h)
}

// SetDueTo = `Activity/SetDueTo_act`: `.DueTo` = 1 bila BalanceDueTo >= 0, 0 bila < 0.
func SetDueTo(h *Halaman) error {
	k := &kalkulator{}
	p := polis{h, k}
	b := p.n("BalanceDueTo")
	if k.err != nil {
		return k.err
	}
	if b.Sign() >= 0 {
		h.Setel(HalamanPolis+".DueTo", "1")
	} else {
		h.Setel(HalamanPolis+".DueTo", "0")
	}
	return nil
}

// ---------------------------------------------------------------- CountOGPONP_Act

// CountOGPONP = `Activity/CountOGPONP_Act` - perhitungan ulang seluruh baris uang.
func CountOGPONP(h *Halaman) error {
	k := &kalkulator{}
	p := polis{h, k}
	nol := apd.New(0, 0)
	// langkah 1: .PremiOgp kosong atau "0" -> lima medan OGP jadi 0
	if p.teks("PremiOgp") == "" || p.teks("PremiOgp") == "0" {
		for _, m := range []string{"PremiOgp", "RiCommOgp", "ResultOgp1", "OveriddingCommOgp", "ResultOgp2"} {
			p.setel(m, nol)
		}
	}
	// langkah 2: .PremiOnp kosong atau "0" -> lima medan ONP jadi 0
	if p.teks("PremiOnp") == "" || p.teks("PremiOnp") == "0" {
		for _, m := range []string{"PremiOnp", "RiCommOnp", "ResultOnp1", "OveriddingCommOnp", "ResultOnp2"} {
			p.setel(m, nol)
		}
	}
	// langkah 3-7
	for _, f := range []func(*Halaman, string) error{CountResult1, CountResult2Ogp, CountResult1Onp, CountResult2Onp} {
		if err := f(h, ""); err != nil {
			return err
		}
	}
	if err := CountNetPremi(h); err != nil {
		return err
	}
	// langkah 8 - kotak When TIDAK dicentang: SELALU dijalankan.
	// .ClaimType = @if((.Claim!=0&&.Claim!="")||(.SalvageValue!=0&&.SalvageValue!=""),"SOA","")
	// .ClaimPaymentType = @if(...,"Claim","")
	k2 := &kalkulator{}
	p2 := polis{h, k2}
	terisi := func(m string) bool { return p2.n(m).Sign() != 0 && p2.teks(m) != "" }
	adaKlaim := terisi("Claim") || terisi("SalvageValue")
	if k2.err != nil {
		return k2.err
	}
	if adaKlaim {
		h.Setel(HalamanPolis+".ClaimType", "SOA")
		h.Setel(HalamanPolis+".ClaimPaymentType", "Claim")
	} else {
		h.Setel(HalamanPolis+".ClaimType", "")
		h.Setel(HalamanPolis+".ClaimPaymentType", "")
	}
	// langkah 9: setiap baris SpreadingRiskList -> CountSpreading_Act(Index)
	daftar := h.AmbilDaftar(HalamanPolis + ".SpreadingRiskList")
	for i := range daftar {
		if err := CountSpreading(h, i+1); err != nil {
			return err
		}
	}
	// langkah 10
	return SetValidateInstallment(h)
}

// ---------------------------------------------------------------- CountResult*

// CountResult1 = `Activity/CountResult1_Act` (OGP hasil 1). `data` = Param.Data.
func CountResult1(h *Halaman, data string) error {
	k := &kalkulator{}
	p := polis{h, k}
	// langkah 1: PremiOgp kosong/"0" -> keluar
	if p.teks("PremiOgp") == "" || p.teks("PremiOgp") == "0" {
		return nil
	}
	ri := p.teks("RiCommOgp")
	// langkah 2: isDouble||kosong -> lanjut, selain itu KELUAR; lalu >100 -> pesan
	if !(AdalahDesimal(ri) || ri == "") {
		return nil
	}
	if p.n("RiCommOgp").Cmp(seratus) > 0 {
		p.pesan("RiCommOgp", PesanErrorMsg1)
	}
	// langkah 3: bukan NonProp dan data=="" -> GrossPremium = @divide(.PremiOgp*100 ,pyWorkPage.TreatyIn.RNMShareP ,4)
	// ⛔ Langkah 3-4 membutuhkan nilai master - dilewati bila tidak tersedia
	// (`MasterTersedia`).
	if !p.nonProp() && data == "" && MasterTersedia(h, "RNMShareP") {
		p.setel("GrossPremium", k.BagiBulat(k.Kali(p.n("PremiOgp"), seratus), p.nj(HalamanMaster+".RNMShareP"), 4))
	}
	// langkah 4: NonProp dan data=="" -> pembaginya TreatyIn.RNMShare
	if p.nonProp() && data == "" && MasterTersedia(h, "RNMShare") {
		p.setel("GrossPremium", k.BagiBulat(k.Kali(p.n("PremiOgp"), seratus), p.nj(HalamanMaster+".RNMShare"), 4))
	}
	if k.err != nil {
		return k.err
	}
	// langkah 5: (RiCommOgp<100 && isDouble) || kosong -> lanjut, selain itu KELUAR
	if !((AdalahDesimal(ri) && p.n("RiCommOgp").Cmp(seratus) < 0) || ri == "") {
		return nil
	}
	if data == "Pct" {
		// .ResultOgp1 = .PremiOgp*(.RiCommOgp/@String.toDecimal("100.00000"))
		p.setel("ResultOgp1", k.Kali(p.n("PremiOgp"), k.Bagi(p.n("RiCommOgp"), seratus5des)))
	}
	// langkah 6: data=="Amount" -> .RiCommOgp = (.ResultOgp1/.PremiOgp)*@String.toDecimal("100.00000")
	if data == "Amount" {
		p.setel("RiCommOgp", k.Kali(k.Bagi(p.n("ResultOgp1"), p.n("PremiOgp")), seratus5des))
	}
	if k.err != nil {
		return k.err
	}
	// langkah 7
	return CountNetPremi(h)
}

// CountResult1Onp = `Activity/CountResult1Onp_Act` (ONP hasil 1).
func CountResult1Onp(h *Halaman, data string) error {
	k := &kalkulator{}
	p := polis{h, k}
	if p.teks("PremiOnp") == "" || p.teks("PremiOnp") == "0" { // langkah 1
		return nil
	}
	ri := p.teks("RiCommOnp")
	if !(AdalahDesimal(ri) || ri == "") { // langkah 2
		return nil
	}
	if p.n("RiCommOnp").Cmp(seratus) > 0 {
		p.pesan("RiCommOnp", PesanErrorMsg1)
	}
	// langkah 3: (RiCommOnp<100 && isDouble) || RiCommOnp="" -> lanjut, selain itu KELUAR
	if !((AdalahDesimal(ri) && p.n("RiCommOnp").Cmp(seratus) < 0) || ri == "") {
		return nil
	}
	if data == "Pct" {
		// .ResultOnp1 = .PremiOnp*(.RiCommOnp/@String.toDecimal("100.00000"))
		p.setel("ResultOnp1", k.Kali(p.n("PremiOnp"), k.Bagi(p.n("RiCommOnp"), seratus5des)))
	}
	if data == "Amount" { // langkah 4
		// .RiCommOnp = (.ResultOnp1/.PremiOnp)*@String.toDecimal("100.00000")
		p.setel("RiCommOnp", k.Kali(k.Bagi(p.n("ResultOnp1"), p.n("PremiOnp")), seratus5des))
	}
	if k.err != nil {
		return k.err
	}
	return CountNetPremi(h) // langkah 5
}

// CountResult2Ogp = `Activity/CountResult2Ogp_act` (OGP hasil 2 - overriding).
func CountResult2Ogp(h *Halaman, data string) error {
	k := &kalkulator{}
	p := polis{h, k}
	if p.teks("PremiOgp") == "" || p.teks("PremiOgp") == "0" { // langkah 1
		return nil
	}
	ov := p.teks("OveriddingCommOgp")
	// langkah 2: @String.isDouble(.OveriddingCommOgp) -> lanjut, selain itu KELUAR
	// ⚠️ Berbeda dari CountResult1_Act: di sini kosong TIDAK diterima.
	if !AdalahDesimal(ov) {
		return nil
	}
	if p.n("OveriddingCommOgp").Cmp(seratus) > 0 {
		p.pesan("OveriddingCommOgp", PesanErrorMsg1)
	}
	// langkah 3: (<100 && isDouble) || =="" -> lanjut, selain itu KELUAR
	if !((AdalahDesimal(ov) && p.n("OveriddingCommOgp").Cmp(seratus) < 0) || ov == "") {
		return nil
	}
	if data == "Pct" {
		// .ResultOgp2 = (.OveriddingCommOgp/@String.toDecimal("100.00000"))*.PremiOgp
		p.setel("ResultOgp2", k.Kali(k.Bagi(p.n("OveriddingCommOgp"), seratus5des), p.n("PremiOgp")))
	}
	if data == "Amount" { // langkah 4
		// .OveriddingCommOgp = (.ResultOgp2/.PremiOgp)*@String.toDecimal("100.00000")
		p.setel("OveriddingCommOgp", k.Kali(k.Bagi(p.n("ResultOgp2"), p.n("PremiOgp")), seratus5des))
	}
	if k.err != nil {
		return k.err
	}
	return CountNetPremi(h) // langkah 5
}

// CountResult2Onp = `Activity/CountResult2Onp_act` (ONP hasil 2 - overriding).
func CountResult2Onp(h *Halaman, data string) error {
	k := &kalkulator{}
	p := polis{h, k}
	if p.teks("PremiOnp") == "" || p.teks("PremiOnp") == "0" { // langkah 1
		return nil
	}
	ov := p.teks("OveriddingCommOnp")
	if !AdalahDesimal(ov) { // langkah 2
		return nil
	}
	if p.n("OveriddingCommOnp").Cmp(seratus) > 0 {
		p.pesan("OveriddingCommOnp", PesanErrorMsg1)
	}
	// langkah 3: .OveriddingCommOnp<100 && isDouble -> lanjut, selain itu KELUAR
	// ⚠️ Di sini kosong TIDAK diterima (berbeda dari CountResult2Ogp_act).
	if !(AdalahDesimal(ov) && p.n("OveriddingCommOnp").Cmp(seratus) < 0) {
		return nil
	}
	if data == "Pct" {
		// .ResultOnp2 = (.OveriddingCommOnp/@String.toDecimal("100.00000"))*.PremiOnp
		p.setel("ResultOnp2", k.Kali(k.Bagi(p.n("OveriddingCommOnp"), seratus5des), p.n("PremiOnp")))
	}
	if data == "Amount" { // langkah 4
		// .OveriddingCommOnp = (.ResultOnp2/.PremiOnp) *@String.toDecimal("100.00000")
		p.setel("OveriddingCommOnp", k.Kali(k.Bagi(p.n("ResultOnp2"), p.n("PremiOnp")), seratus5des))
	}
	if k.err != nil {
		return k.err
	}
	return CountNetPremi(h) // langkah 5
}

// ---------------------------------------------------------------- layar Dept Head

// ---------------------------------------------------------------- CalculatePremi_Act

// ---------------------------------------------------------------- bantu

// MasterTersedia - nilai `pyWorkPage.TreatyIn.<m>` terisi.
//
// ⛔ PENYIMPANGAN SADAR yang dipaksa keputusan work owner P29. Di Pega halaman
// `TreatyIn` diisi dari JSON master kontrak (`FetchMasterTreatyIn`,
// `InputPolicyTreatyInDetail_preACT` langkah 10 `adoptJSONObject`) - SELALU
// terisi saat langkah-langkah ini berjalan. P29 membuang JSON; view
// `TREATYINDETAILJOINEDM` tidak punya `RNMShareP`, `RNMShare`, maupun
// `BrokeragePercentP`, dan kolom `RNM_SHARE` di view TIDAK BOLEH dipakai dalam
// perhitungan sebelum maknanya dijelaskan (PERTANYAAN-untuk-DBA, `[terbuka]`
// spec §9.2 butir 8). Memetakannya berarti mengarang pemetaan.
//
// Tanpa penjaga ini setiap `@divide(x, TreatyIn.RNMShareP)` membagi dengan nol
// (kosong = 0) dan SELURUH rantai hitung gagal - layar tidak dapat dipakai sama
// sekali. Yang dipilih: langkah yang membutuhkan nilai master DILEWATI selama
// nilainya kosong, langkah lain berjalan seperti XML. Begitu sumber nilai
// master ditetapkan, penjaga ini tidak pernah aktif lagi.
func MasterTersedia(h *Halaman, m string) bool {
	return strings.TrimSpace(h.Ambil(HalamanMaster+"."+m)) != ""
}

// formatAngka menulis desimal sebagai teks bertitik tanpa notasi ilmiah.
func formatAngka(d *apd.Decimal) string { return d.Text('f') }

// KlaimXOLRetro - nilai `PolicyTreatyIn.ClaimType` jalur XOL Retro (`Activity/SetValueEDM_Act` langkah 5/7,
// `InputPolicyTreatyEDMDetail_NP` langkah 8/9, `TreatyRealizationCheckXOLListEDM` langkah 6/7). Asal: salinan
// konstanta `modul/nbtreatyin/backend/models/sumberbisnis.go`.
const KlaimXOLRetro = "XOL Retro"
