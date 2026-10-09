package models

// Untuk apa berkas ini: PORT PRINT DLA dan ACCEPTATION layar Adjustment & Acceptation.
//
//   - Print DLA (tombol baris objek ShowObjectAdj): `CheckLimitSpreadingTreaty_Act` + `ProtectPrint` +
//     `CreateRemarksNote_Act` (pra-proses pop-up PrintDLA_dtl), lalu Submit `ChooseDla_Act` -> `GenerateDLAFacin_Act`
//     (fac retro, nomor "P") / `DLAFacintoTreaty_Act` (treaty, nomor "S").
//   - Acceptation (tombol panel InputAdjustment): `SaveAcceptation` + `HitServiceToKasir_Act` (jalur `IsCLM` 13.2),
//     status kasir `GetStatusKasir_Act` (defer load).
//
// Berkas PDF (Property-Set-HTML + HTMLToPDF + InsertDocument_Act) = OQ-CFI-20: aliran HTML tidak diekspor. Yang ditiru
// hanya akibat datanya: nomor DLA, DLA_No adjustment, FacRetroList adjustment, TotalEstimasiReas item, UPDATE
// OS_AKSEPTASI_KLAIM (InsertDLA_OS_SQL), kronologi.

import (
	"regexp"
	"strconv"

	"github.com/cockroachdb/apd/v3"
)

// Teks dan pesan VERBATIM.
const (
	TeksPrintDLA          = "Print DLA"                             // GenerateDLAFacin_Act 2, DLAFacintoTreaty_Act 2
	TeksAcceptation       = "Acceptation"                           // SaveAcceptation 3
	PesanPLASebelumDLA    = "Please print the PLA before Print DLA" // ProtectPrint 2
	JenisNomorDLAFac      = "P"                                     // GenerateDLAFacin_Act 8
	JenisNomorDLATreaty   = "S"                                     // DLAFacintoTreaty_Act 6
	JenisServiceAkseptasi = "AKSEPATSI"                             // ChooseDla_Act 7 (ejaan Pega)
	JenisServiceDLA       = "DLA"                                   // ChooseDla_Act 11
	StatusKasirSukses     = "Akseptasi Sudah Masuk ke Kasir"        // GetStatusKasir_Act 5, HitServiceToKasir_Act 13.5
	KetKasirSukses        = "Success"                               // GetStatusKasir_Act 5
)

// PLADulu = ProtectPrint 3-4: objek retro yang belum Print PLA - `pyWorkPage.Messages` berisi pesan, Submit
// PrintDLA_dtl tersembunyi (`VIS pyWorkPage.Messages = ”`). Pesan itu dihitung dari objek, bukan disimpan.
func PLADulu(h *Halaman, o int) bool {
	ob, err := Objek(h, o)
	return err == nil && ob["IsFacretro"] == "1" && ob["PlaStatus"] != "1"
}

// BukaDLA = tombol "Print DLA" baris objek (`VIS .IsFacretro = 1`, `NA .DLAStatus != 0`):
//
//   - CheckLimitSpreadingTreaty_Act - kelas ObjectItem dijalankan atas baris OBJEK: perulangan `.SpreadingList` /
//     `.Adjustment` halaman primer kosong (objek tidak memilikinya), sehingga satu-satunya akibat adalah langkah 1
//     membuang `CedingCoRetroList` polis - halaman polis dimuat ulang setiap aksi dan pembacanya hanya
//     DLAFacintoTreaty_Act. `CekLimit.CARI1` tidak pernah ditulis (OQ-CFI-24); prioritas operator langkah 20-21 tetap
//     OQ (prompt §5).
//   - ProtectPrint - `PLADulu`.
//   - CreateRemarksNote_Act - RemarksDLA = "- Refer to PLA : " + NoPLA baris PLAList ber-TipePLA 10015 terakhir
//     (langkah 2 pre=false: selalu ditimpa).
func BukaDLA(h *Halaman, o int) error {
	ob, err := Objek(h, o)
	if err != nil {
		return err
	}
	pla := ""
	for _, p := range h.AmbilDaftar(JalurAnak(DaftarObjek, o, DaftarPLA)) { // 1
		if p["TipePLA"] == TreatyFacRetro {
			pla = p["NoPLA"]
		}
	}
	ob["RemarksDLA"] = AwalanRemarksPLA + pla // 2
	if ob["ShareRetro"] == "" {               // PrintDLA_dtl: .ShareRetro DEF=1
		ob["ShareRetro"] = "1"
	}
	return nil
}

// LayarDLA - Section PrintDLA_dtl (harness PrintDLA_dtl_Harness) untuk objek o. Penerima surel tidak dibangun (pola
// LayarPLA, OQ-CFI-22).
func LayarDLA(o int) []Unsur {
	p := JalurObjek(o) + "."
	pesanKosong := sama(JalurPesanSurel, "")
	kirim := sama(JalurKirimSurel, "true")
	dulu := func(h *Halaman) bool { return PLADulu(h, o) }
	return []Unsur{
		tampil(label(PesanPLASebelumDLA), dulu),
		medan(JalurKirimSurel, "Send Email", KCentang),
		tampil(roJika(medan(JalurPesanSurel, "Message", KArea),
			func(h *Halaman) bool { return AmbilJalur(h, p+"IsEditable") != "1" }), kirim),
		sumber(medan(p+"ShareRetro", "Share Retro", KRadio), kode("ShareRetro")),
		wajibU(medan(p+"RemarksDLA", "RemarksDLA", KArea)),
		tampil(naJika(tombol("ChooseDla", "Submit", "ChooseDla"), dan(kirim, pesanKosong)),
			func(h *Halaman) bool { return !PLADulu(h, o) }),
	}
}

// RencanaDLA - hasil murni ChooseDla_Act langkah 1-3.
type RencanaDLA struct {
	// Keluar - langkah 1 (`.IsFacretro==1` lalu `.PlaStatus != 1` -> T=6).
	Keluar bool
	// Urutan - jenis nomor yang diterbitkan ("P" GenerateDLAFacin_Act, "S" DLAFacintoTreaty_Act) menurut kemunculan
	// pertamanya di SpreadingList item.
	//
	// `[penyimpangan sadar]` (PARITAS): Pega memanggil GenerateDLAFacin_Act sekali untuk SETIAP baris spreading 10015
	// setiap item - objek ber-dua item retro menerbitkan dua nomor DLA dalam satu klik (yang kedua tak dipakai
	// adjustment mana pun karena DLA_No sudah terisi, tetapi tetap menimpa NoDLA objek). Di sini setiap jenis sekali.
	Urutan []string
}

// ChooseDLA = ChooseDla_Act langkah 1-3.1.2 (Param.Index = o): DLAStatus 1, IsTreatyOut, jenis DLA yang diterbitkan.
// Jalur treaty (3.1.4) menuntut `CekLimit.CARI1 == 1` yang tak pernah tercapai (OQ-CFI-24) - ditiru apa adanya.
func ChooseDLA(h *Halaman, o int) (RencanaDLA, error) {
	var r RencanaDLA
	ob, err := Objek(h, o)
	if err != nil {
		return r, err
	}
	if ob["IsFacretro"] == "1" && ob["PlaStatus"] != "1" { // 1
		r.Keluar = true
		return r, nil
	}
	ob["DLAStatus"] = "1" // 2
	cek := h.Ambil(JalurCekLimit) == "1"
	lihat := map[string]bool{}
	treatyOut := ""
	for i := range h.AmbilDaftar(DaftarItem(o)) { // 3
		for _, s := range h.AmbilDaftar(DaftarDiItem(o, i+1, AnakSpreadPolis)) { // 3.1
			if spreadTreaty(s) { // 3.1.1
				treatyOut = "1"
			}
			if treatyOut != "" { // 3.1.2 (Local.SpreadingTreaty tak pernah dinolkan)
				h.Setel(JalurIsTreatyOut, treatyOut)
			}
			jenis := ""
			switch {
			case s["TreatyType"] == TreatyFacRetro: // 3.1.3
				jenis = JenisNomorDLAFac
			case spreadTreaty(s) && cek: // 3.1.4
				jenis = JenisNomorDLATreaty
			}
			if jenis != "" && !lihat[jenis] {
				lihat[jenis] = true
				r.Urutan = append(r.Urutan, jenis)
			}
		}
	}
	return r, nil
}

// dlaLayak - adjustment yang menerima nomor DLA: diterima komite, sudah Print Acceptation, DLA_No kosong
// (GenerateDLAFacin_Act 15.3.11.1.1 / 20.1.1, DLAFacintoTreaty_Act 13.2.1).
func dlaLayak(b Baris) bool {
	return b["AcceptanceStatus"] == "1" && b["DLA_No"] == "" && b["IsPrintAccept"] == "1"
}

// salinRetroAdj - FacRetroList polis ke `Adjustment.FacRetroList(<CURRENT>)` (GenerateDLAFacin_Act 15.2.1.2.1:
// AdditionalInfo, PctShareAllObj, ReinsurerID, ReinsurerName; DLAFacintoTreaty_Act 13.2.3.2 Page-Copy utuh).
// EndPeriod / StartPeriod / OurRef / TFAllObj tidak disimpan: pembacanya hanya HTML DLA (OQ-CFI-20); OurRef = NoDLA
// objek.
func salinRetroAdj(h *Halaman, daftar string, retro []Baris, utuh bool) {
	ada := h.AmbilDaftar(daftar)
	out := make([]Baris, 0, len(retro))
	for n, r := range retro {
		b := Baris{}
		if n < len(ada) {
			b = ada[n]
		}
		b["AdditionalInfo"], b["PctShareAllObj"] = r["AdditionalInfo"], r["PctShareAllObj"]
		b["ReinsurerID"], b["ReinsurerName"] = r["ReinsurerID"], r["ReinsurerName"]
		if utuh {
			b["RiCommAllObj"] = r["RiCommAllObj"]
		}
		out = append(out, b)
	}
	// Baris lebih yang tidak tersentuh <CURRENT> tetap ada.
	if len(ada) > len(retro) {
		out = append(out, ada[len(retro):]...)
	}
	h.SetelDaftar(daftar, out)
}

// HasilDLA - akibat data penerbitan satu nomor DLA.
type HasilDLA struct {
	// Akseptasi - AcceptedNo adjustment yang menerima DLA_No (InsertDLA_OS_SQL `UPDATE ... WHERE AcceptedNo`).
	Akseptasi []string
	// Keluar - GenerateDLAFacin_Act 15.3.12 (tak ada adjustment layak) mengakhiri activity sebelum langkah 20.
	Keluar bool
}

// TerapkanDLAFac = GenerateDLAFacin_Act (Param.IndexObject = o) tanpa PDF:
//
//	11      NoDLA objek = nomor (pre=false: selalu, juga bila NoDLA sudah terisi)
//	2, 14   kronologi "Print DLA"
//	15.2    setiap item beradjustment: RemarksDLA objek + FacRetroList polis ke SETIAP adjustment-nya
//	15.3    per reasuradur: adjustment layak (dlaLayak) seluruh item; nol -> keluar (15.3.12, pada item pertama)
//	15.3.15 TotalEstimasiReas item = hasilperreas / hasilfacout x hasilfacoutclaim (reasuradur terakhir)
//	20      adjustment layak: DLA_No = nomor; AcceptedNo dikembalikan untuk UPDATE OS
func TerapkanDLAFac(k *Konteks, h *Halaman, o int, nomor string) (HasilDLA, error) {
	var r HasilDLA
	ob, err := Objek(h, o)
	if err != nil {
		return r, err
	}
	ob["NoDLA"] = nomor          // 11
	k.Kronologi(h, TeksPrintDLA) // 2, 14
	retro := h.AmbilDaftar(DaftarFacRetro)
	items := h.AmbilDaftar(DaftarItem(o))
	layak := false
	for i := range items {
		for _, b := range h.AmbilDaftar(DaftarAdj(o, i+1)) {
			if dlaLayak(b) {
				layak = true
			}
		}
	}
	var akhir *apd.Decimal
	if len(retro) > 0 && layak {
		if akhir, err = hasilAkhirDLA(h, o, retro); err != nil {
			return r, err
		}
	}
	for i := range items { // 15
		adj := h.AmbilDaftar(DaftarAdj(o, i+1))
		for a := range adj { // 15.2
			adj[a]["RemarksDLA"] = ob["RemarksDLA"]
			salinRetroAdj(h, DaftarDiAdj(o, i+1, a+1, AnakFacRetro), retro, false)
		}
		if len(retro) == 0 {
			continue
		}
		if !layak { // 15.3.12
			r.Keluar = true
			return r, nil
		}
		items[i]["TotalEstimasiReas"] = Teks(akhir) // 15.3.15
	}
	for i := range items { // 20 (EMBEDDED: semua item)
		for _, b := range h.AmbilDaftar(DaftarAdj(o, i+1)) {
			if dlaLayak(b) {
				r.Akseptasi = append(r.Akseptasi, b["AcceptedNo"])
				b["DLA_No"] = nomor
			}
		}
	}
	return r, nil
}

// bulat4p = `@Math.divide(x, 100, 4)`.
func bulat4p(kk *Kalkulator, x *apd.Decimal) *apd.Decimal { return bulat4(kk.Bagi(x, apd.New(100, 0))) }

// hasilAkhirDLA - GenerateDLAFacin_Act 15.3.10-15.3.15 untuk reasuradur TERAKHIR (lokal dinolkan per reasuradur di
// 15.3.10, jadi hanya putaran terakhir yang tersisa): Local.ValueADJ menjumlah adjustment layak retro seluruh item
// (1/2/5 AdjustmentValue, 3 SalvageValue, 4/6 AdjusterFeeValue); per baris spreading 10015 item ber-adjustment layak
// hasilfacout = TSINusare x share / 100, hasilfacoutclaim = ValueADJ x share / 100 (4 desimal, item terakhir menang).
// Pembagi hasilfacout nol (Pega galat langkah) -> 0, `[inferensi]`.
func hasilAkhirDLA(h *Halaman, o int, retro []Baris) (*apd.Decimal, error) {
	var kk Kalkulator
	nilai := apd.New(0, 0)
	facout, facoutKlaim := apd.New(0, 0), apd.New(0, 0)
	for i, it := range h.AmbilDaftar(DaftarItem(o)) {
		ada := false
		for _, b := range h.AmbilDaftar(DaftarAdj(o, i+1)) {
			if !dlaLayak(b) {
				continue
			}
			ada = true
			if b["IsFacRetro"] != "1" {
				continue
			}
			switch b["PaymentType"] { // 15.3.13.4.4-15.3.13.4.9
			case BayarFinal, BayarInterim, BayarAdjust:
				nilai = kk.Tambah(nilai, kk.B(b, "AdjustmentValue"))
			case BayarSalvage:
				nilai = kk.Tambah(nilai, kk.B(b, "SalvageValue"))
			case BayarFee, BayarExpense:
				nilai = kk.Tambah(nilai, kk.B(b, "AdjusterFeeValue"))
			}
		}
		if !ada { // 15.3.11.2 item tanpa adjustment layak dibuang dari TempData
			continue
		}
		for _, s := range h.AmbilDaftar(DaftarDiItem(o, i+1, AnakSpreadPolis)) { // 15.3.13.5.1
			if s["TreatyType"] != TreatyFacRetro {
				continue
			}
			share := kk.B(s, "SharePercentage")
			facout = bulat4p(&kk, kk.Kali(kk.B(it, "TSINusare"), share))
			facoutKlaim = bulat4p(&kk, kk.Kali(nilai, share))
		}
	}
	perReas := bulat4p(&kk, kk.Kali(facout, kk.B(retro[len(retro)-1], "PctShareAllObj"))) // 15.3.15
	akhir := kk.Kali(kk.BagiPega(perReas, facout), facoutKlaim)
	return akhir, kk.Galat()
}

// TerapkanDLATreaty = DLAFacintoTreaty_Act (Param.IndexObject = o) tanpa PDF: NoDLA objek (9, pre=false), kronologi
// (12), per item / adjustment: DLA_No bagi yang layak (13.2.1), RemarksDLA (13.2.2), FacRetroList polis utuh
// (13.2.3.2). 13.3.13 (TotalEstimasiReas) ber-remark; 13.8 berulang atas `.Adjustment` baris FacRetroList (kosong).
// Tanpa InsertDLA_OS_SQL (tidak ada di activity ini).
func TerapkanDLATreaty(k *Konteks, h *Halaman, o int, nomor string) error {
	ob, err := Objek(h, o)
	if err != nil {
		return err
	}
	ob["NoDLA"] = nomor          // 9
	k.Kronologi(h, TeksPrintDLA) // 2, 12
	retro := h.AmbilDaftar(DaftarFacRetro)
	for i := range h.AmbilDaftar(DaftarItem(o)) { // 13
		for a, b := range h.AmbilDaftar(DaftarAdj(o, i+1)) {
			if dlaLayak(b) { // 13.2.1
				b["DLA_No"] = nomor
			}
			b["RemarksDLA"] = ob["RemarksDLA"] // 13.2.2
			salinRetroAdj(h, DaftarDiAdj(o, i+1, a+1, AnakFacRetro), retro, true)
		}
	}
	return nil
}

// ---------------------------------------------------------------- Acceptation

// BolehAkseptasi - tombol "Acceptation": tampil `.AcceptanceStatus = 1`, nonaktif `.IsPrintAccept = 1 &&
// (.DirectToKasir = 'false' || .StatusKasir != ”)`.
func BolehAkseptasi(b Baris) bool {
	if b["AcceptanceStatus"] != "1" {
		return false
	}
	return !(b["IsPrintAccept"] == "1" && (b["DirectToKasir"] == "false" || b["StatusKasir"] != ""))
}

// SaveAcceptation = SaveAcceptation langkah 1-7 (adjustment a item i objek o): false bila sudah Print Acceptation
// (1, T=6). Kronologi "Acceptation"; IsFacretro objek / item / adjustment menurut spreading polis 10015 item (5-7).
// Langkah 8 (nama treaty spreading adjustment kosong, SetTreatNameAdjustment) dijalankan services lewat
// `NamaJenisReas`.
func SaveAcceptation(k *Konteks, h *Halaman, o, i, a int) (bool, error) {
	b, err := Adj(h, o, i, a)
	if err != nil {
		return false, err
	}
	if b["IsPrintAccept"] == "1" { // 1
		return false, nil
	}
	ob, err := Objek(h, o)
	if err != nil {
		return false, err
	}
	it, err := Item(h, o, i)
	if err != nil {
		return false, err
	}
	k.Kronologi(h, TeksAcceptation) // 3-4
	retro := "0"
	for _, s := range h.AmbilDaftar(DaftarDiItem(o, i, AnakSpreadPolis)) { // 5
		if s["TreatyType"] == TreatyFacRetro {
			retro = "1"
		}
	}
	ob["IsFacretro"], b["IsFacRetro"], it["IsFacretro"] = retro, retro, retro // 6-7
	return true, nil
}

// SelesaiAkseptasi = SaveAcceptation langkah 13: IsPrintAccept adjustment + objek 1, DLAStatus objek 0 (tombol Print
// DLA terbuka). Langkah 12 PrintPDFAccep_MultiAksep = OQ-CFI-20.
func SelesaiAkseptasi(h *Halaman, o, i, a int) error {
	b, err := Adj(h, o, i, a)
	if err != nil {
		return err
	}
	ob, err := Objek(h, o)
	if err != nil {
		return err
	}
	b["IsPrintAccept"] = "1"
	ob["IsPrintAccept"], ob["DLAStatus"] = "1", "0"
	return nil
}

// TerapkanStatusKasir = GetStatusKasir_Act (defer load blok kasir): hanya baris diterima ber-DirectToKasir "true" yang
// belum sukses (1); KET log "Success" -> teks sukses, selainnya KET apa adanya (5). Tanpa baris log: tidak berubah.
func TerapkanStatusKasir(b Baris, ket string, ada bool) {
	if !(b["AcceptanceStatus"] == "1" && b["DirectToKasir"] == "true") || b["StatusKasir"] == StatusKasirSukses || !ada {
		return
	}
	if ket == KetKasirSukses {
		b["StatusKasir"] = StatusKasirSukses
		return
	}
	b["StatusKasir"] = ket
}

// MuatanKasir - badan `SendAcceptationToKasir` (`TAllPaymentData`, HitServiceToKasir_Act 13.2.2.1.1.3-13.3, jalur
// `IsCLM`). Disimpan di outbox (efek "kasir"), dikirim hanya di produksi.
type MuatanKasir struct {
	NoTrans, NoKlaim, LbuId, NoPolis, AcceptType, Kepada, AccountNo, TglAksep string
	Nett, Deductible, KaliDeduct                                              string
	StsSyariah, CompanyName, LjtdId, LdcId, StsAp, LkuId, LbgID               string
	TglBolehBayar, Email, UserInput                                           string
}

// Kode tetap muatan kasir (13.2.2.1.1.3; bukan data orang).
const (
	kasirCompany = "NUSARE"
	kasirLjtd    = "D0031"
	kasirLdc     = "100081"
)

var bukanAngka = regexp.MustCompile(`[^0-9]`)

// PanjangNoAksepCLM - HitServiceToKasir_Act 13.2: `@length(.AcceptedNo) = 21 || 22`; selainnya keluar (T=6).
func PanjangNoAksepCLM(no string) bool { return len(no) == 21 || len(no) == 22 }

// BolehKasir - HitServiceToKasir_Act langkah 2: `.DirectToKasir=="true" && .StatusKasir = ""`, selainnya keluar.
func BolehKasir(b Baris) bool { return b["DirectToKasir"] == "true" && b["StatusKasir"] == "" }

// KunciCedingKasir - 13.2.2.1.1.1 `@substring(pyWorkPage.Quotation.CedingCo, 0, 8)`. `pyWorkPage.Quotation` tidak
// pernah diisi rule Claim Fac In mana pun; `[inferensi]` (PARITAS): kode ceding polis (`QuotationData.CedingCo`).
func KunciCedingKasir(h *Halaman) string {
	c := h.Ambil(OQ + "CedingCo")
	if len(c) > 8 {
		return c[:8]
	}
	return c
}

// tglAksepTeks = `@substring(d,6,8)+"-"+@substring(d,4,6)+"-"+@substring(d,0,4)` atas "yyyyMMdd".
func tglAksepTeks(ymd string) string {
	if len(ymd) < 8 {
		return ""
	}
	return ymd[6:8] + "-" + ymd[4:6] + "-" + ymd[0:4]
}

// TglBolehBayar = 13.2.2.1.1.5: hari > 25 -> tanggal 01 dua bulan sesudahnya, selainnya hari yang sama bulan berikutnya;
// bulan 13 -> "1"; tahun naik bila bulan hasil 01 DAN bulan berjalan Desember.
func TglBolehBayar(ymd string, bulanKini int) string {
	if len(ymd) < 8 {
		return ""
	}
	hari, _ := strconv.Atoi(ymd[6:8])
	bulan, _ := strconv.Atoi(ymd[4:6])
	tahun, _ := strconv.Atoi(ymd[0:4])
	bulan++
	if hari > 25 {
		bulan++
	}
	if bulan >= 13 {
		bulan = 1
	}
	hs := strconv.Itoa(hari)
	if hari > 25 {
		hs = "01"
	}
	if len(hs) < 2 {
		hs = "0" + hs
	}
	bs := strconv.Itoa(bulan)
	if len(bs) < 2 {
		bs = "0" + bs
	}
	if bulan == 1 && bulanKini == 12 {
		tahun++
	}
	return hs + "-" + bs + "-" + strconv.Itoa(tahun)
}

// SusunMuatanKasir = HitServiceToKasir_Act 7 dan 13.2 (jalur CLM, Local.IndexObject = o): seluruh adjustment item objek
// o ber-AcceptedNo sama dijumlah - Nett menurut Payment Type (1/2/5 AdjustmentValue, 4/6 AdjusterFeeValue, 3
// SalvageValue), Deductible Σ IndividualRiskRNM, KaliDeduct Σ IndividualRiskPercentage; medan lain dari baris cocok
// TERAKHIR. `email` = GL.F_GET_EMAIL ceding. IsPEGASyariah (13.2.2.1.1.4, pemeriksaan node server) = konvensional.
func SusunMuatanKasir(k *Konteks, h *Halaman, o int, b Baris, email string) (MuatanKasir, error) {
	b["NoAccount"] = bukanAngka.ReplaceAllString(b["NoAccount"], "") // 7
	var kk Kalkulator
	nett, ded, kali := apd.New(0, 0), apd.New(0, 0), apd.New(0, 0)
	akhir := b
	for i := range h.AmbilDaftar(DaftarItem(o)) { // 13.2.2
		for _, x := range h.AmbilDaftar(DaftarAdj(o, i+1)) {
			if x["AcceptedNo"] != b["AcceptedNo"] { // 13.2.2.1.1
				continue
			}
			akhir = x
			ded = kk.Tambah(ded, kk.B(x, "IndividualRiskRNM"))
			kali = kk.Tambah(kali, kk.B(x, "IndividualRiskPercentage"))
			switch x["PaymentType"] { // 13.2.2.1.1.6-8
			case BayarFinal, BayarInterim, BayarAdjust:
				nett = kk.Tambah(nett, kk.B(x, "AdjustmentValue"))
			case BayarFee, BayarExpense:
				nett = kk.Tambah(nett, kk.B(x, "AdjusterFeeValue"))
			case BayarSalvage:
				nett = kk.Tambah(nett, kk.B(x, "SalvageValue"))
			}
		}
	}
	if err := kk.Galat(); err != nil {
		return MuatanKasir{}, err
	}
	tgl := YMD(akhir["AcceptedDate"])
	user := akhir["pxCreateOperator"]
	if user == "" {
		user = k.Pelaku
	}
	return MuatanKasir{
		NoTrans: akhir["AcceptedNo"], NoKlaim: h.Ambil(CD + "NoClaim"), LbuId: h.Ambil(OQ + "BusinessOldId"),
		NoPolis: h.Ambil(JalurNoPolis), AcceptType: akhir["PaymentType"], Kepada: akhir["PayableTo"],
		AccountNo: bukanAngka.ReplaceAllString(akhir["NoAccount"], ""), TglAksep: tglAksepTeks(tgl), Nett: Teks(nett),
		Deductible: Teks(ded), KaliDeduct: Teks(kali), StsSyariah: "0", CompanyName: kasirCompany, LjtdId: kasirLjtd,
		LdcId: kasirLdc, StsAp: "0", LkuId: akhir["CurrencyID"], LbgID: akhir["IDOfBank"],
		TglBolehBayar: TglBolehBayar(tgl, int(k.Sekarang.In(Jakarta).Month())), Email: email, UserInput: user,
	}, nil
}

// TeksDraftDLA - DraftGenerateDLAFacin_Act 3 (VERBATIM).
const TeksDraftDLA = "Print Draft DLA"

// DraftDLA = DraftGenerateDLAFacin_Act (tombol "Send Claim to Committee", sebelum CreateKMTNo_Act) tanpa PDF: hanya
// adjustment retro (1 Exit-Activity bila `.IsFacRetro != 1`); ShareRetro objek 1 dan kronologi (3, 11); setiap item
// beradjustment: RemarksDLA objek + FacRetroList polis ke setiap adjustment-nya (12.2). Penomoran (6-8) ber-remark:
// draft tanpa nomor. Berkas draft = OQ-CFI-20.
func DraftDLA(k *Konteks, h *Halaman, o, i, a int) error {
	b, err := Adj(h, o, i, a)
	if err != nil {
		return err
	}
	if b["IsFacRetro"] != "1" { // 1
		return nil
	}
	ob, err := Objek(h, o)
	if err != nil {
		return err
	}
	ob["ShareRetro"] = "1"       // 3
	k.Kronologi(h, TeksDraftDLA) // 3, 11
	retro := h.AmbilDaftar(DaftarFacRetro)
	for n := range h.AmbilDaftar(DaftarItem(o)) { // 12
		for m, x := range h.AmbilDaftar(DaftarAdj(o, n+1)) { // 12.2.1
			x["RemarksDLA"] = ob["RemarksDLA"]
			salinRetroAdj(h, DaftarDiAdj(o, n+1, m+1, AnakFacRetro), retro, false)
		}
	}
	return nil
}
