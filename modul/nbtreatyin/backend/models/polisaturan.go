package models

// Untuk apa berkas ini: ATURAN MURNI LAIN di sekitar layar realisasi - pra-proses
// flow action (Data Transform), penanda penempatan keluar,
// status riwayat, dan bentuk nomor polis. Seluruhnya port satu rule; nama rule
// dan langkahnya ditulis di atas tiap fungsi (INVENTARIS-XML.md bab 6-7).

import (
	"strconv"
	"strings"
	"time"

	"nusantarare/inti/backend/utils"
)

// ---------------------------------------------------------------- pra-proses

// PraprosesAdmin = `DataTransform/InputPolicyTreatyIn_preDT` (pra-proses flow
// action `InboxPolicyTreatyIn`, layar admin).
//
//	1-2  StartDate / EndDate kosong -> hari ini (AC 34, 35; P35: EndDate kosong
//	     diisi HARI INI, bukan +1 tahun - ditiru apa adanya, spec §10.2)
//	3    StatementDate kosong -> sekarang
//	4-5  IsNewPolicyNonProp "1" bila Quotation.ProportionalType NonProportional,
//	     selain itu "0" bila belum "1" (AC 75)
//	6-10 MasterID "", MarketingOfficer <- Quotation.MarketingName, Suggest "",
//	     SuggestDate sekarang, OperatorName <- NAMA TAMPILAN (P33, AC 40, 42)
//	11   FlagOnGoingPolicy "1" (T_WORK_POLIS.FLAG_ONGOING_POLICY)
//	14   PolicyTreatyIn.QuotationData <- Quotation (salinan halaman)
//
// ⛔ Langkah 12-13 (`TempEmail.CARI28` / `PositionNote` dari
// `OperatorID.pyWorkBasketList(2)`) TIDAK diport: penunjukan antrean menurut
// NOMOR URUT diganti pemeriksaan keanggotaan (P25, AC 14) di services.
// ⚠️ Dua format tanggal rule ini (`dd/MM/yyyy` dan `MM/dd/yyyy hh:mm a`)
// diganti SATU format (P32, AC 33) - penyimpangan sadar spec §5.8.
func PraprosesAdmin(h *Halaman, sekarang time.Time, namaTampilan string) {
	tgl := utils.FormatTanggal(sekarang)
	if h.Ambil(HalamanPolis+".StartDate") == "" {
		h.Setel(HalamanPolis+".StartDate", tgl)
	}
	if h.Ambil(HalamanPolis+".EndDate") == "" {
		h.Setel(HalamanPolis+".EndDate", tgl)
	}
	if h.Ambil(HalamanPolis+".StatementDate") == "" {
		h.Setel(HalamanPolis+".StatementDate", utils.FormatTanggalWaktu(sekarang))
	}
	if h.Ambil(HalamanQuotation+".ProportionalType") == JenisNonProporsional {
		h.Setel(HalamanPolis+".IsNewPolicyNonProp", "1")
	}
	if h.Ambil(HalamanPolis+".IsNewPolicyNonProp") != "1" {
		h.Setel(HalamanPolis+".IsNewPolicyNonProp", "0")
	}
	h.Setel(HalamanPolis+".MasterID", "")
	h.Setel(HalamanPolis+".MarketingOfficer", h.Ambil(HalamanQuotation+".MarketingName"))
	h.Setel(HalamanPolis+".Suggest", "")
	h.Setel(HalamanPolis+".SuggestDate", utils.FormatTanggalWaktu(sekarang))
	h.Setel(HalamanPolis+".OperatorName", namaTampilan)
	h.Setel("FlagOnGoingPolicy", "1")
	SalinQuotation(h)
}

// SystemSetOneYear = `DataTransform/SystemSetOneYear_DT` langkah 1:
// `.EndDate = @DateTime.addCalendar(.StartDate,1,0,0,0,0,0,0)`.
//
// Pemicunya (XML): pra-transform (`pyPreDataTransform`) aksi refresh `change`
// medan `.StartDate` di `Section/DetailPolicyTreatyIn` - layar admin, medan
// wajib dan dapat disunting. Di `Section/DetailDeptHeadTreatyIn_UW` medan yang
// sama `pyReadOnly` (`1==1`), jadi aturan ini hanya berjalan di layar admin.
// Ia BERBEDA dari pengisian tanggal-akhir-kosong `PraprosesAdmin` (P35/AC 34):
// itu berjalan saat layar dibuka, ini saat tanggal mulai diubah, dan menimpa
// `.EndDate` tanpa syarat.
//
// `addCalendar` Pega = `java.util.Calendar.add(YEAR, 1)`: 29 Februari dijepit
// ke 28 Februari pada tahun bukan kabisat (bukan digulirkan ke 1 Maret).
// ⚠️ `[dugaan]` tanggal mulai kosong/tak terbaca: `.EndDate` dibiarkan.
func SystemSetOneYear(h *Halaman) {
	mulai, err := utils.ParseTanggal(strings.TrimSpace(h.Ambil(HalamanPolis + ".StartDate")))
	if err != nil {
		return
	}
	y, m, d := mulai.Date()
	if akhirBulan := time.Date(y+1, m+1, 0, 0, 0, 0, 0, time.UTC).Day(); d > akhirBulan {
		d = akhirBulan
	}
	h.Setel(HalamanPolis+".EndDate", utils.FormatTanggal(time.Date(y+1, m, d, 0, 0, 0, 0, time.UTC)))
}

// PraprosesAtasan = `DataTransform/DeptHeadTreatyInUW_preDT` (pra-proses flow
// action `DeptHeadTreatyIn_UW`, layar Sec Head dan Dept Head).
//
// ⛔ Langkah 2 rule mengisi OperatorName dari `OperatorID.pyUserIdentifier`
// (pengenal akun) - BUG, bukan maksud (P33, AC 42); di sini NAMA TAMPILAN.
// ⛔ Langkah 5 (`isApprovedtoDeptHead`) tidak dibangun - P36, AC 64.
func PraprosesAtasan(h *Halaman, sekarang time.Time, namaTampilan string) {
	h.Setel(HalamanPolis+".MasterID", "")
	h.Setel(HalamanPolis+".OperatorName", namaTampilan)
	h.Setel(HalamanPolis+".SuggestDate", utils.FormatTanggalWaktu(sekarang))
	h.Setel(HalamanPolis+".IsApproved", "")
	h.Setel(HalamanPolis+".Suggest", "")
	SalinQuotation(h)
}

// SalinQuotation meniru `pyWorkPage.PolicyTreatyIn.QuotationData = pyWorkPage.Quotation`
// (salinan seluruh halaman, termasuk properti yang kosong di sumber).
func SalinQuotation(h *Halaman) {
	h.pastikan()
	for k := range h.Nilai {
		if strings.HasPrefix(k, HalamanPolis+".QuotationData.") {
			delete(h.Nilai, k)
		}
	}
	for k, v := range h.Nilai {
		if strings.HasPrefix(k, HalamanQuotation+".") {
			h.Nilai[HalamanPolis+".QuotationData."+strings.TrimPrefix(k, HalamanQuotation+".")] = v
		}
	}
}

// GeserTanggalProduksi = `InputPolicyTreatyInPre_Act` langkah 9 ("Set
// Production Date kalau diatas tanggal 25"):
//
//	syarat  `@substring(pyWorkPage.PolicyTreatyIn.StatementDate,6,2)>25`
//	        benar -> jalankan; salah -> lewati (hari = batas TIDAK digeser)
//	ProductionDate = @addCalendar(StatementDate,'0','1','0','0','0','0','0')
//	ProductionDate = @substring(ProductionDate,0,6) + "01" + @substring(ProductionDate,8)
//
// Hasilnya tanggal 1 bulan berikut; `@substring(..,8)` membawa BAGIAN WAKTU
// StatementDate (`THHmmss.SSS GMT`) apa adanya - jamnya dipertahankan.
// ⛔ RALAT putaran 2: port pertama menulis pukul 00:00:00 - tidak ada di rule.
//
// ⛔ PENYIMPANGAN SADAR. Rule menanam `>25`; `GeneratePolicyNoTreaty_Act`
// langkah 5.3 membaca hari yang sama dari `POOLDATA.TANGGAL_CLOSING`
// (`GETTanggalClosing_SQL` -> `Local.TglProd`). Dua aturan hidup berdampingan
// - pola yang SAMA dengan `InsertJsonPolisLife_Act` b1170 lawan
// `PROC_GENERATE_SEQUENCE_NUMBER` di PremiumList Life, tempat `[keputusan work
// owner]` menetapkan "ikuti yang dari DB" (penjaga
// `TestNolAmbangTutupBukuTertanam`). Preseden itu diterapkan di sini:
// `hariClosing` dibaca pemanggil dari tabel (`penomor.HariClosing`), dan
// pembandingnya tetap `>` seperti rule.
//
// ⚠️ Catatan: `@substring(StatementDate,6,2)` membaca hari dari teks internal
// DateTime Pega (GMT); di sini hari dibaca di zona waktu `statement` itu sendiri
// (jam aplikasi). Keduanya sama kecuali pukul 00:00-06:59 WIB bila jam
// aplikasi berzona WIB - dicatat, tidak ditiru (bergantung zona JVM Pega, yang
// tidak ada di korpus).
func GeserTanggalProduksi(statement time.Time, hariClosing int) time.Time {
	if statement.Day() > hariClosing {
		return time.Date(statement.Year(), statement.Month()+1, 1,
			statement.Hour(), statement.Minute(), statement.Second(), statement.Nanosecond(), statement.Location())
	}
	return statement
}

// ---------------------------------------------------------------- penempatan keluar

// Kode `TreatyType` penanda punya penempatan keluar - VERBATIM
// `DataTransform/TestTreatyToFacStatus` (AC 76).
var KodePenempatanKeluar = []string{"10015", "10218"}

// TetapkanHasFacOut = `InboxPolicyTreatyIn_postDT` langkah 2-3:
// `.PolicyTreatyIn.HasFacOut = "0"`, lalu `TestTreatyToFacStatus` memasang "1"
// bila salah satu baris SpreadingRiskList ber-TreatyType salah satu kode itu.
func TetapkanHasFacOut(h *Halaman) {
	h.Setel(HalamanPolis+".HasFacOut", "0")
	for _, b := range h.AmbilDaftar(DaftarSpreading) {
		for _, kode := range KodePenempatanKeluar {
			if b["TreatyType"] == kode {
				h.Setel(HalamanPolis+".HasFacOut", "1")
			}
		}
	}
}

// ---------------------------------------------------------------- riwayat

// StatusRiwayat = `Activity/InsertHistoryAkseptasiPega` langkah 2-3 untuk
// kasus treaty: IsApproved "1" -> ACCEPT, "0" -> REJECT, selain itu kosong
// (keenam penanda EmailType* milik modul lain, nol di halaman treaty).
func StatusRiwayat(isApproved string) string {
	switch isApproved {
	case "1":
		return "ACCEPT"
	case "0":
		return "REJECT"
	}
	return ""
}

// ---------------------------------------------------------------- nomor polis

// Kode tipe nomor polis - VERBATIM `GeneratePolicyNoTreaty_Act` langkah 7-9.
const (
	TipeNomorQR = "QR" // DueTo == 1 (premi)
	TipeNomorQP = "QP" // DueTo == 0 (bayar)
	TipeNomorTP = "TP" // ClaimType == "XOL Retro"
)

// TipeNomorPolis = langkah 7-9; langkah 9 (XOL Retro) dijalankan TERAKHIR,
// sehingga menimpa QR/QP.
func TipeNomorPolis(h *Halaman) string {
	t := ""
	switch h.Ambil(HalamanPolis + ".DueTo") {
	case "1":
		t = TipeNomorQR
	case "0":
		t = TipeNomorQP
	}
	if h.Ambil(HalamanPolis+".ClaimType") == KlaimXOLRetro {
		t = TipeNomorTP
	}
	return t
}

// KelasDeret dan JenisDeret adalah kunci `GENERATE_SEQUENCE_NUMBER` yang
// dikirim langkah 26: `ParamSeq.CARI1 = pyWorkPage.pxObjClass`,
// `ParamSeq.CARI2 = ParamSeq.HASIL3+"QR/QP/TP"` - teks "QR/QP/TP" HARFIAH,
// sehingga ketiga tipe berbagi SATU deret.
//
// `pyWorkPage.pxObjClass` kasus NB Treaty In = `ASM-FW-GISFW-Work-NB` (kelas
// halaman `pyWorkPage` di `pyPagesAndClasses` DecisionTable/isApproved dan
// DataTransform/DeptHeadTreatyInUW_preDT).
const KelasDeret = "ASM-FW-GISFW-Work-NB"

// JenisDeret = awalan produksi + "QR/QP/TP".
func JenisDeret(awalan string) string { return awalan + "QR/QP/TP" }

// TanggalProduksiNomor = `GeneratePolicyNoTreaty_Act` langkah 5.1-5.3.
//
//	5.1 ProductionDate = @CurrentDateTime()        (kotak When tak dicentang)
//	5.2 StatementDate di depan hari ini -> ProductionDate = StatementDate
//	5.3 hari ProductionDate (Asia/Jakarta) > tanggal closing ->
//	    tanggal 1 bulan berikut, pukul 05:00 GMT
func TanggalProduksiNomor(sekarang, statement time.Time, hariClosing int) time.Time {
	prod := sekarang
	// 5.2 `@toInt(@DateTime.DateTimeDifference(StatementDate,@CurrentDateTime(),D))<0`
	// - selisih (sekarang - statement) dalam HARI, dipotong ke nol oleh @toInt:
	// benar hanya bila StatementDate sedikitnya satu hari penuh di depan.
	if !statement.IsZero() && sekarang.Sub(statement) <= -24*time.Hour {
		prod = statement
	}
	lokal := prod.In(zonaJakarta())
	if lokal.Day() > hariClosing {
		bulan := time.Date(lokal.Year(), lokal.Month(), 1, 5, 0, 0, 0, time.UTC)
		return bulan.AddDate(0, 1, 0)
	}
	return prod
}

// RakitNomorPolis = langkah 28:
// `ParamSeq.CARI4 + ".T" + OJKBusinessID + "." + ParamSeq.HASIL1 + "." + ParamSeq.HASIL2`
// dengan CARI4 = awalan + tipe, HASIL1 = `MM.YYYY`, HASIL2 = urut lima angka.
func RakitNomorPolis(awalan, tipe, ojkBusinessID, mmYYYY string, urut int) string {
	u := strconv.Itoa(urut)
	for len(u) < 5 {
		u = "0" + u
	}
	return awalan + tipe + ".T" + ojkBusinessID + "." + mmYYYY + "." + u
}

func zonaJakarta() *time.Location {
	if l, err := time.LoadLocation("Asia/Jakarta"); err == nil {
		return l
	}
	return time.FixedZone("WIB", 7*3600)
}
