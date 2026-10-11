package models

// Untuk apa berkas ini: PORT TOMBOL "Download Claim Face Sheet" baris objek (`CLaimFaceSheet_Act`, kelas Data-Object)
// beserta activity yang dipanggilnya: `ProtectDownloadFaceClaim` (langkah 1), `SetTryeatyName_ACT` (27), `SaveCFS_ACT`
// (43.7.2 - baris OS STS 0), `CheckLimit_Act1` (50).
//
// Bagian murni ada di sini; services menjalankan nomor klaim (penomor, langkah 12-17), tulisan OS_AKSEPTASI_KLAIM,
// JSON_KLAIM (43.9 `InsertJsonClaimNonMBU_act`), MONITORING_KLAIM_LOG (48-49) dan outbox konversi (47, hanya produksi)
// dalam SATU transaksi aksi.
//
// Yang tidak dibawa (PARITAS):
//   - 4-6 pesan "click save spreading first !": gerbangnya `ClaimData.ExGratia == 1`, sedangkan ExGratia tingkat klaim
//     selalu 0 (InsertObjects_dt langkah 11) - jalur mati.
//   - 8-11 ber-remark (`//`).
//   - 18-26, 28-41: halaman sementara TempData / FacoutRetro / MataUang dan HTML -> PDF. Aliran HTML CFS tidak diekspor
//     (OQ-CFI-20) sehingga tidak ada berkas dan 42 `InsertDocument_Act` tidak menulis DOCUMENT_CLAIM.
//   - SaveCFS_ACT 2-8 menambah baris `ClaimData.osAkseptasi` yang tidak dibaca rule mana pun (sensus korpus).

import (
	"strconv"
	"strings"

	"github.com/cockroachdb/apd/v3"
)

// Pesan dan teks VERBATIM.
const (
	TeksCFS                 = "Download Claim Face Sheet (OutStanding)"                     // CLaimFaceSheet_Act 2
	PesanEstimasiTakSelaras = "Unsynchronized estimation data detected. Please contact IT." // ProtectDownloadFaceClaim 3
)

// Penanda kasus (`pyWorkPage`).
const (
	JalurIsCFS         = "IsCFS"         // CLaimFaceSheet_Act 46 - tombol Back InputEstimasiAdmin tersembunyi
	JalurIsPicTransfer = "IsPicTransfer" // CLaimFaceSheet_Act 43.5, GeneratePLATreaty_Act 19 - Send to PIC Claim
	JalurIsTreatyOut   = "IsTreatyOut"   // GeneratePLA 6
)

// Nomor klaim (CLaimFaceSheet_Act 12-15).
const (
	JenisNomorKlaim         = "K"
	JenisNomorPLATreaty     = "G" // GeneratePLATreaty_Act 6
	TipeKodeProduksiNonLife = "NONLIFE"
)

// TreatyFacRetro - jenis treaty retro fakultatif (`.TreatyType == "10015"`).
const TreatyFacRetro = "10015"

// JenisServiceEstimasi - `Param.JenisService` CLaimFaceSheet_Act 48.
const JenisServiceEstimasi = "ESTIMASI"

// Daftar halaman polis yang ditulis CheckLimit_Act1.
const (
	// DaftarCedingRetro - `OfferFacIn.QuotationData.CedingCoRetroList` (langkah 5 / 8; dibaca DLAFacintoTreaty_Act dan
	// GeneratePLATreaty_Act). Halaman polis tidak disimpan: daftar ini disusun ulang (`SusunCedingRetro`).
	DaftarCedingRetro = OQ + "CedingCoRetroList"
	AnakSpreadCeding  = "SpreadingCedingList"
	// DaftarFacRetroTreaty - pengganti `OfferFacIn.FacRetroList` yang ditulis CheckLimit_Act1 9.1.1.4 (reasuradur treaty
	// saat estimasi melampaui limit). Pega menyimpannya di halaman polis milik kasus; halaman polis di sini dibaca ulang
	// dari JSON_POLIS, maka penggantinya disimpan sebagai baris retro tingkat klaim (T_CLAIM_FAC_RETRO tanpa
	// ADJUSTMENT_ID) dan dipasang ulang sesudah polis dimuat (`TerapkanRetroTreaty`).
	DaftarFacRetroTreaty = CD + "FacRetroTreaty"
)

// RakitNomorKlaim = CLaimFaceSheet_Act 15: `ParamSeq.CARI2 + BusinessOldId + "." + HASIL1 + "." + HASIL2`
// (CARI2 = kode NONLIFE + "K", HASIL1 = MM.YYYY, HASIL2 = LPAD(urut,5,'0')). Sama untuk PLA treaty ("G",
// GeneratePLATreaty_Act 8).
func RakitNomorKlaim(jenis, oldID, mmYYYY string, urut int) string {
	return jenis + oldID + "." + mmYYYY + "." + lpad5(urut)
}

func lpad5(n int) string {
	s := strconv.Itoa(n)
	for len(s) < 5 {
		s = "0" + s
	}
	return s
}

// ---------------------------------------------------------------- ProtectDownloadFaceClaim

// EstimasiTercetak = ProtectDownloadFaceClaim langkah 4: Σ EstimationValue baris estimasi ber-PrintFaceClaim 1 seluruh
// objek. (`SET CountEstimasi + .TotalEstimasi` di langkah 4.1 adalah parameter langkah iterasi tanpa metode - tidak
// dijalankan.)
func EstimasiTercetak(h *Halaman) (*apd.Decimal, error) {
	var kk Kalkulator
	total := apd.New(0, 0)
	for o := range h.AmbilDaftar(DaftarObjek) {
		for i := range h.AmbilDaftar(DaftarItem(o + 1)) {
			for _, e := range h.AmbilDaftar(DaftarDiItem(o+1, i+1, AnakEstimasi)) {
				if e["PrintFaceClaim"] == "1" {
					total = kk.Tambah(total, kk.B(e, "EstimationValue"))
				}
			}
		}
	}
	return total, kk.Galat()
}

// ProtectDownloadFaceClaim - langkah 5-6: Σ CLAIM_VALUE REINSURANCE.TRLOSS_DETAIL_T kasus (NO_AKSEP kosong) dibanding
// langkah 4 pada dua desimal; selisih = pesan di `.PrintFaceClaim` objek.
//
// Bacaan langkah 2 hanya di produksi (`IsPEGAPROD`). Di luar produksi Pega membandingkan nilai kosong dengan Σ halaman,
// sehingga CFS objek kedua selalu tertahan; di sini pemeriksaan dilewati di luar produksi (pola Claim Prop OQ-CP-11;
// PARITAS `[penyimpangan sadar]`).
func ProtectDownloadFaceClaim(k *Konteks, h *Halaman, o int, kasusID string) error {
	if !k.Produksi {
		return nil
	}
	tersimpan, err := k.Acuan.EstimasiKasusTerbuka(k.Ctxt(), kasusID)
	if err != nil {
		return err
	}
	halaman, err := EstimasiTercetak(h)
	if err != nil {
		return err
	}
	var kk Kalkulator
	a := bulatDua(kk.Teks("TotalEst.HASILD1", tersimpan))
	if err := kk.Galat(); err != nil {
		return err
	}
	if Banding(a, bulatDua(halaman)) != 0 {
		h.TambahPesan(JalurAnak(DaftarObjek, o, "PrintFaceClaim"), PesanEstimasiTakSelaras)
	}
	return nil
}

// bulatDua = `@divide(x, 1, 2)`: dua desimal, setengah ke atas (`[inferensi]` BigDecimal ROUND_HALF_UP).
func bulatDua(x *apd.Decimal) *apd.Decimal {
	c := apd.BaseContext.WithPrecision(60)
	c.Rounding = apd.RoundHalfUp
	r := new(apd.Decimal)
	_, _ = c.Quantize(r, x, -2)
	return r
}

// ---------------------------------------------------------------- SetTryeatyName_ACT

// NamaTreatySpreading = `SetTryeatyName_ACT` (CLaimFaceSheet_Act 27, setiap item objek o): TreatyName SpreadingList dan
// SpreadingClaim yang kosong diisi `REINSURANCETYPE.NOTE` menurut TreatyType.
func NamaTreatySpreading(k *Konteks, h *Halaman, o int) error {
	for i := range h.AmbilDaftar(DaftarItem(o)) {
		for _, anak := range []string{AnakSpreadPolis, AnakSpreadKlaim} {
			for _, s := range h.AmbilDaftar(DaftarDiItem(o, i+1, anak)) {
				if s["TreatyName"] != "" {
					continue
				}
				nama, err := k.Acuan.NamaJenisReas(k.Ctxt(), s["TreatyType"])
				if err != nil {
					return err
				}
				s["TreatyName"] = nama
			}
		}
	}
	return nil
}

// ---------------------------------------------------------------- langkah 43-46 + SaveCFS_ACT

// SetelCFS = CLaimFaceSheet_Act langkah 43-46 untuk objek o (nomor klaim sudah terisi, langkah 12-17):
//
//	43.1-43.5  item ber-spreading 10015 -> IsFacretro 1, selainnya 0 + `IsPicTransfer` kasus 1; item ber-TotalEstimasi ->
//	           PrintFaceClaim 1
//	43.7       setiap estimasi bernilai yang belum tercetak -> satu baris OS STS 0 (`SaveCFS_ACT`), lalu PrintFaceClaim 1
//	44         objek IsFacretro 1 bila item TERAKHIR ber-10015 (nilai lokal loop terakhir)
//	45-46      objek PrintFaceClaim 1, PlaStatus 0, CFS 0; kasus IsCFS 1 (tombol Back tersembunyi)
func SetelCFS(k *Konteks, h *Halaman, o int, kasusID string) ([]BarisOS, error) {
	ob, err := Objek(h, o)
	if err != nil {
		return nil, err
	}
	retro := "2"
	var rows []BarisOS
	for i, it := range h.AmbilDaftar(DaftarItem(o)) {
		n := i + 1
		retro = "2"                                                            // 43.1
		for _, s := range h.AmbilDaftar(DaftarDiItem(o, n, AnakSpreadPolis)) { // 43.2
			if s["TreatyType"] == TreatyFacRetro {
				retro = "1"
			}
		}
		if it["TotalEstimasi"] != "" { // 43.3
			it["PrintFaceClaim"] = "1"
		}
		if retro == "1" { // 43.4
			it["IsFacretro"] = "1"
		} else { // 43.5
			it["IsFacretro"] = "0"
			h.Setel(JalurIsPicTransfer, "1")
		}
		for e, est := range h.AmbilDaftar(DaftarDiItem(o, n, AnakEstimasi)) { // 43.7
			if est["EstimationValue"] != "" && est["PrintFaceClaim"] != "1" { // 43.7.2
				rows = append(rows, BarisOSEstimasi(k, h, o, n, e+1, kasusID))
			}
			if est["EstimationValue"] != "" { // 43.7.3
				est["PrintFaceClaim"] = "1"
			}
		}
	}
	if retro == "1" { // 44
		ob["IsFacretro"] = "1"
	}
	ob["PrintFaceClaim"], ob["PlaStatus"], ob["CFS"] = "1", "0", "0" // 45
	h.Setel(JalurIsCFS, "1")                                         // 46
	return rows, nil
}

// BarisOSEstimasi = `SaveCFS_ACT` (Data-ObjectItem; IdxObject o, IdxObejctItem i, IdxEstimation e): halaman
// TempOSAkseptasi per lini (langkah 9-15, berurutan - lini yang cocok lebih dari satu saling menimpa), pxCreateOperator
// dan Type 0 (16), STS_DLA 8 bila item IsFacretro 1 selainnya 7 (17), DATA_JSON (18), baris OS STS 0 (19).
func BarisOSEstimasi(k *Konteks, h *Halaman, o, i, e int, kasusID string) BarisOS {
	ob, _ := Objek(h, o)
	it, _ := Item(h, o, i)
	est, _ := Estimasi(h, o, i, e)
	p := map[string]string{}
	umum := func() {
		p["CauseOfLoss"] = h.Ambil(CD + "CauseOfLoss")
		p["PersenRNM"] = h.Ambil(AwalanPolis + ".PercentShare")
		p["ObjectID"] = ob["ObjectID"]
		p["Type"] = "0"
		p["Value"] = est["EstimationValue"]
		p["GrossValue"] = est["GrossEstimationPct"]
		p["DeductibleValue"] = est["Deductible"]
		p["Currency"] = est["Currency"]
		p["CurrencyID"] = est["CurrencyID"]
		p["EstimationDate"] = k.Sekarang.In(Jakarta).Format("20060102") // @getCurrentDateStamp() (DEV: 8 digit)
	}
	jaminan := func() {
		p["CauseOfLossID"] = h.Ambil(CD + "CauseOfLossID")
		p["CoverageID"] = it["CoverageID"]
		p["CoverageName"] = it["CoverageNote"]
	}
	if IsFire(h) || IsAneka(h) || IsGolfInsurance(h) { // 9-11
		umum()
		jaminan()
		p["ObjectName"] = it["ObjectItemName"]
	}
	if IsMarineCargo(h) { // 12
		umum()
		p["CoverageName"] = ""
		if c := h.AmbilDaftar(DaftarDiItem(o, i, "CoverageList")); len(c) > 0 {
			p["CoverageName"] = c[0]["CoverageNote"]
		}
	}
	if IsMBU(h) { // 13
		umum()
		jaminan()
		p["GrossValue"] = est["GrossEstimationPctMBU"]
		p["NoClaim"] = h.Ambil(CD + "NoClaim")
	}
	if IsTravel(h) || IsPA(h) { // 14-15
		umum()
		jaminan()
		p["ObjectName"] = ob["ObjectName"]
		p["ObjectIDCard"] = ob["ObjectIDCard"]
		p["ObjectParticipantStatus"] = ob["ObjectParticipantStatus"]
		p["ObjectDateOfBirth"] = ob["ObjectDateOfBirth"]
		p["NoClaim"] = h.Ambil(CD + "NoClaim")
		if IsPA(h) {
			p["ObjectJob"] = ob["ObjectJob"]
			p["Gender"] = ob["Gender"]
		}
	}
	p["pxCreateOperator"] = k.Pelaku // 16
	p["Type"] = "0"
	p["pxObjClass"] = KelasOSAkseptasi
	dla := StsDLAFac // 17
	if it["IsFacretro"] == "1" {
		dla = StsDLARetro
	}
	return BarisOS{CaseID: KunciInstans(kasusID), NoClaim: h.Ambil(CD + "NoClaim"), NoPolis: h.Ambil(JalurNoPolis),
		StsReject: StsOSOutstanding, StsDLA: dla, DataJSON: JSONHalamanPega(HalamanJSON{Nilai: p})}
}

// ---------------------------------------------------------------- CheckLimit_Act1

// BatasTreaty - satu baris `GetDataTreatyLimit_Sql` (TREATYBUSINESS x PROPORTIONALARRG, TREATYDESCID 10001, kontrak
// berlaku pada tanggal mulai polis).
type BatasTreaty struct {
	TreatyYear, TreatyGroupID, Limit string
}

// BarisQuotaShare - satu baris `GetQuotaShare` (PROPORTIONALARRG anak menurut PARENTREINSTYPEID).
type BarisQuotaShare struct {
	TreatyType, TreatyName, SharePercentage string
}

// ReasuradurTreaty - satu baris `GetListRetro_Sql` (TREATYREINSURER).
type ReasuradurTreaty struct {
	ReinsurerID, ReinsurerName, PctShareAllObj, RiCommAllObj string
}

// qsFac - `@contains(.TreatyName,"QS") && @contains(.TreatyName,"FAC")`.
func qsFac(nama string) bool { return strings.Contains(nama, "QS") && strings.Contains(nama, "FAC") }

// CheckLimit = `CheckLimit_Act1` (CLaimFaceSheet_Act 50, setiap item objek o):
//
//	1     daftar sementara dibuang, termasuk CedingCoRetroList polis dan `.SpreadingAdjustment` item (Break QS)
//	2-3   nilai = TotalEstimationValueinIDR (kosong -> TotalEstimasi); IsMoreThanTreatyLimit item false
//	4     SpreadingList ber-"QS"+"FAC" -> GetDataTreatyLimit_Sql: TreatyGroupID / TreatyYear klaim, limit (RP)
//	5-8   CedingCoRetroList + Break QS item (`SusunCedingRetro`); 6 Σ klaim polis terbuka (hanya produksi)
//	9     limit terisi dan (nilai > limit ATAU Σ polis > limit) -> IsFacretro / IsMoreThanTreatyLimit item + objek, lalu
//	      FacRetroList polis diganti reasuradur treaty (GetListRetro_Sql)
//
// Perbandingan 9 numerik (`[inferensi]`: FindData.HASILD2 bertipe desimal).
func CheckLimit(k *Konteks, h *Halaman, o, i int) error {
	it, err := Item(h, o, i)
	if err != nil {
		return err
	}
	ob, err := Objek(h, o)
	if err != nil {
		return err
	}
	nilai := it["TotalEstimationValueinIDR"] // 2
	if nilai == "" {                         // 3
		nilai = it["TotalEstimasi"]
	}
	it["IsMoreThanTreatyLimit"] = "false"
	if ob["IsMoreThanTreatyLimit"] != "true" {
		ob["IsMoreThanTreatyLimit"] = "false"
	}
	limit := ""
	for _, s := range h.AmbilDaftar(DaftarDiItem(o, i, AnakSpreadPolis)) { // 4
		if !qsFac(s["TreatyName"]) {
			continue
		}
		b, _, err := k.Acuan.BatasTreaty(k.Ctxt(), h.Ambil(OQ+"BusinessCode"), s["TreatyType"], h.Ambil(JalurMulaiPolis))
		if err != nil {
			return err
		}
		h.Setel(CD+"TreatyGroupID", b.TreatyGroupID) // 4.1.3 (Out.pxResults(1) kosong -> "")
		h.Setel(CD+"TreatyYear", b.TreatyYear)
		limit = b.Limit
	}
	if err := SusunCedingRetro(k, h, o, i); err != nil { // 1, 5, 8
		return err
	}
	total := ""
	if k.Produksi { // 6-7
		if total, err = k.Acuan.EstimasiPolisTerbuka(k.Ctxt(), h.Ambil(JalurNoPolis)); err != nil {
			return err
		}
	}
	if limit == "" { // 9 [T=3]
		return nil
	}
	lewat, err := lebihDari(nilai, limit)
	if err != nil {
		return err
	}
	lewatPolis, err := lebihDari(total, limit)
	if err != nil {
		return err
	}
	if !lewat && !lewatPolis {
		return nil
	}
	it["IsFacretro"], ob["IsFacretro"] = "1", "1"
	it["IsMoreThanTreatyLimit"], ob["IsMoreThanTreatyLimit"] = "true", "true"
	for _, s := range h.AmbilDaftar(DaftarDiItem(o, i, AnakSpreadPolis)) { // 9.1
		if !qsFac(s["TreatyName"]) { // 9.1.1
			continue
		}
		rs, err := k.Acuan.ReasuradurTreaty(k.Ctxt(), s["TreatyType"], h.Ambil(CD+"TreatyYear"),
			h.Ambil(CD+"TreatyGroupID"))
		if err != nil {
			return err
		}
		var fr []Baris
		for _, r := range rs { // 9.1.1.4.1
			fr = append(fr, Baris{"ReinsurerID": r.ReinsurerID, "ReinsurerName": r.ReinsurerName,
				"PctShareAllObj": r.PctShareAllObj, "RiCommAllObj": r.RiCommAllObj})
		}
		h.SetelDaftar(DaftarFacRetro, fr) // 9.1.1.2 Property-Remove lalu isi
		h.SetelDaftar(DaftarFacRetroTreaty, SalinDaftar(fr))
	}
	return nil
}

// SusunCedingRetro = CheckLimit_Act1 langkah 1 (sebagian), 5 dan 8: CedingCoRetroList polis dari SpreadingClaim item
// ber-"QS"+"FAC" (CedingCoName / PercentShare polis, TreatyTypeCeding, mata uang), anak SpreadingCedingList dan Break QS
// item (`.SpreadingAdjustment`) dari `GetQuotaShare` (TreatyYear / TreatyGroupID klaim, induk = TreatyTypeCeding).
//
// Langkah 8.1 mengambil TreatyYear / TreatyGroupID klaim ke FindData.CARI3 / CARI6 yang juga dipakai 9.1.1.3; di sini
// keduanya dibaca langsung dari klaim (nilai yang sama).
func SusunCedingRetro(k *Konteks, h *Halaman, o, i int) error {
	h.SetelDaftar(DaftarCedingRetro, nil)
	h.SetelDaftar(DaftarDiItem(o, i, AnakBreakQS), nil)
	var ceding []Baris
	for _, s := range h.AmbilDaftar(DaftarDiItem(o, i, AnakSpreadKlaim)) { // 5
		if qsFac(s["TreatyName"]) {
			ceding = append(ceding, Baris{"CedingCoName": h.Ambil(OQ + "CedingCoName"),
				"PersenShare": h.Ambil(AwalanPolis + ".PercentShare"), "TreatyTypeCeding": s["TreatyType"],
				"Currency": s["Currency"], "CurrencyID": s["CurrencyID"]})
		}
	}
	h.SetelDaftar(DaftarCedingRetro, ceding)
	var breakQS []Baris
	for c, cr := range ceding { // 8
		rows, err := k.Acuan.QuotaShare(k.Ctxt(), h.Ambil(CD+"TreatyYear"), h.Ambil(CD+"TreatyGroupID"),
			cr["TreatyTypeCeding"])
		if err != nil {
			return err
		}
		var anak []Baris
		for _, q := range rows { // 8.4.1
			anak = append(anak, Baris{"SharePercentage": q.SharePercentage, "TreatyName": q.TreatyName})
			breakQS = append(breakQS, Baris{"TreatyName": q.TreatyName, "SharePercentage": q.SharePercentage,
				"TreatyType": q.TreatyType, "Currency": cr["Currency"], "CurrencyID": cr["CurrencyID"]})
		}
		h.SetelDaftar(JalurAnak(DaftarCedingRetro, c+1, AnakSpreadCeding), anak)
	}
	h.SetelDaftar(DaftarDiItem(o, i, AnakBreakQS), breakQS)
	return nil
}

// lebihDari - `a > b` numerik; nilai kosong tidak pernah lebih besar.
func lebihDari(a, b string) (bool, error) {
	if strings.TrimSpace(a) == "" || strings.TrimSpace(b) == "" {
		return false, nil
	}
	var kk Kalkulator
	x, y := kk.Teks("a", a), kk.Teks("b", b)
	if err := kk.Galat(); err != nil {
		return false, err
	}
	return Lebih(x, y), nil
}

// TerapkanRetroTreaty memasang pengganti FacRetroList polis yang tersimpan di klaim (CheckLimit_Act1 9.1.1.4) sesudah
// halaman polis dimuat ulang dari JSON_POLIS.
func TerapkanRetroTreaty(h *Halaman) {
	if fr := h.AmbilDaftar(DaftarFacRetroTreaty); len(fr) > 0 {
		h.SetelDaftar(DaftarFacRetro, SalinDaftar(fr))
	}
}
