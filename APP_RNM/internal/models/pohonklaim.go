package models

// Bentuk pohon klaim yang baru - tiket 14.
//
// Untuk apa berkas ini: struct untuk tingkat-tingkat pohon klaim yang belum ada
// di klaimlife.go, yaitu akar work object dan kedua tingkat spreading.
//
// Dibaca sesudah: klaimlife.go (Klaim, Peserta, BarisAdjustment) dan money.go.
//
// Pohonnya, dari akar ke daun:
//
//	T_WORK_CLAIM                            <- WorkClaim
//	  T_GENERAL_CLAIM                       <- Klaim      (shared primary key)
//	    T_CLAIMLF_PREMIUMLIST_DETAIL        <- Peserta
//	      T_CLAIMLF_ADJUSTMENT              <- BarisAdjustment
//	        T_CLAIMLF_ADJUSTMENT_SPREADING  <- Spreading
//	          ..._SPREADING_RETRO           <- SpreadingRetro
//
// "Shared primary key" berarti Klaim.ID sama persis dengan WorkClaim.ID; tidak
// ada kolom penyambung terpisah.
//
// ⛔ Nol aturan dagang di berkas ini. Ia hanya bentuk data.

import (
	"time"

	"nusantarare/inti/uang"
)

// Lini adalah penanda lini usaha pada work object. Untuk Claim Life isinya
// tetap LiniLife.
//
// [terbuka - Non-Life] Daftar nilai enum lintas-lini ditetapkan saat konteks
// Non-Life digarap. Yang sudah pasti hanyalah nilai untuk Life.
const LiniLife = "LIFE"

// WorkClaim adalah akar pohon: satu baris work object, lintas-lini.
//
// ID-nya TEKS BERFORMAT - CLM-xxxxxx untuk baris klaim, KMTLF-xxxxxx untuk baris
// komite - bukan angka sequence. Itu penyimpangan sadar dari ADR-U-0006 yang
// dicatat di tiket 14.
//
// [terbuka] Pembangkit nomornya belum ditetapkan; pemilik DBA atau work owner.
type WorkClaim struct {
	ID string
	// CoverKey menunjuk WorkClaim induknya. Kosong bila baris ini tidak punya
	// induk - misalnya baris klaim, yang memang akar.
	CoverKey   string
	Lini       string
	PyPosition string
	// AcceptStatus sengaja TIDAK ada di sini. Hasil akseptasi milik kasus
	// komite, dan T_GENERAL_KOMITE.ACCEPT_STATUS sudah menyimpannya - keputusan
	// work owner 2026-09-18, lihat STRUKTUR-TABEL-CLAIM-LIFE.md bab
	// "ACCEPT_STATUS - DIBUANG dari tabel ini". Diagram di tiket 14 masih
	// mencantumkannya; diagram itu yang tertinggal, bukan berkas ini.
	SendtoAdmin   string
	SendtoMedical string
	Type          string
	CaseID        string
	CreateOp      string
	CreateOpName  string
	TglUpdate     time.Time
	// Tahap adalah nama assignment VERBATIM `pyTaskName` - butir **at**.
	//
	// ⛔ Ada karena `PyPosition` TIDAK dapat membedakan Input Register
	// dari Outstanding Claim: keduanya dipegang `ReasLifeAdmin`,
	// sedangkan `Send Back to Register` membuktikan keduanya keadaan
	// yang berbeda (`Section/InputOSClaimLife.xml:21404`).
	Tahap string
	// TglCreate adalah waktu kasus LAHIR - butir **au**, padanan
	// `pxCreateDateTime`.
	//
	// ⚠️ Terpisah dari `TglUpdate`, yang DITIMPA tiap perpindahan.
	// Kotak masuk diurutkan dengan yang ini
	// (`ReportDefinition/InboxPremiumList.xml:736`); mengurutkannya
	// dengan waktu ubah membuat daftar melompat-lompat setiap kali
	// seseorang memindah kasus lain.
	TglCreate time.Time
}

// BarisKlaim menyatakan work object ini adalah baris klaim, bukan baris komite.
// Pembedanya CoverKey: baris klaim tidak punya induk.
func (w WorkClaim) BarisKlaim() bool { return w.CoverKey == "" }

// Spreading adalah hasil spreading satu baris adjustment, per treaty-year.
//
// Nilainya DIBEKUKAN saat adjustment disimpan: perubahan master treaty
// sesudahnya tidak mengubah angka yang sudah tersimpan.
type Spreading struct {
	ID             string
	AdjustmentID   string
	TreatyTypeID   string
	TreatyTypeName string
	TreatyYearLife string
	// RetrocadedShare adalah UANG, bukan rasio - diralat 26-09-2026 menurut
	// SpreadingClaimLife_Act langkah 8.2.1.4-7, yang mengisinya dengan sisa
	// CLAIM_GROSS atau dengan kapasitas IDR/USD treaty-year. Keduanya nilai
	// uang. Nama berakhiran _SHARE di korpus ini memang menipu - lihat
	// ADR-0003 dan bab Catatan tiket 03.
	RetrocadedShare uang.Money
	// Rate tetap RASIO, dan karena itu tidak pernah terjumlahkan dengan Money
	// (ADR-F-0004). Yang disimpan adalah rate MENTAH per mil, bukan hasil
	// baginya seribu.
	Rate uang.Ratio
	// IDR dan USD adalah nilai uang; nama kolomnya sekaligus mata uangnya.
	IDR      uang.Money
	USD      uang.Money
	Currency string
	Retro    []SpreadingRetro
}

// SpreadingRetro adalah pecahan spreading per reinsurer - tingkat terdalam.
//
// Ejaan Commision mengikuti nama kolom apa adanya, termasuk kekeliruan
// ejaannya, supaya tidak lahir dua nama untuk satu kolom.
type SpreadingRetro struct {
	ID            string
	SpreadingID   string
	ReinsurerName string
	PercentShare  uang.Ratio
	Amount        uang.Money
	Rate          uang.Ratio
	// Kedua cabang PremiumSpreadedNet ternyata dipilih oleh TAHUN POLIS -
	// lihat bab "Pembacaan ulang XML" tiket 03. Di sini ia hanya disimpan apa
	// adanya, tidak dihitung ulang.
	PremiumSpreadedGross uang.Money
	PremiumSpreadedNet   uang.Money
	// Commision dan OvrComm adalah PERSEN, bukan uang - diralat 26-09-2026:
	// SpreadingClaimLife_Act membagi keduanya seratus sebelum memakainya
	// (@divide(.COMMISION,100,5) dan @divide(.OVR_COMM,100,5)). Menyimpannya
	// sebagai Money membuat persen dapat dijumlahkan dengan uang, yang justru
	// dilarang ADR-F-0004.
	Commision      uang.Ratio
	OvrComm        uang.Ratio
	TreatyTypeID   string
	TreatyTypeName string
}

// Dokumen adalah satu dokumen pendukung milik seorang peserta.
//
// Kolom isinya dihitung DUA KALI dari sumber yang berbeda, 26-09-2026:
// katalog POOLDATA.DOCUMENT_CLAIM (14 kolom) dan seluruh Property-Set di
// InsertDocument_Act sebelum Obj-Save (12 properti). Keduanya cocok - lihat
// STRUKTUR-TABEL-CLAIM-LIFE.md bab kolom isi. Yang sengaja TIDAK ikut: kunci
// internal Pega (IDPEGA, INSKEY_*), pelaku (milik jejak audit tiket 09), dan
// NOAKSEP/NOPREKAS yang rumahnya baris adjustment.
type Dokumen struct {
	// ⛔ Nama JSON DITULIS. Tanpa tag, Go mengirim `ID`/`NamaFile`
	// berhuruf besar dan React membaca `undefined` tanpa satu pun galat -
	// daftar dokumen yang tampil kosong padahal barisnya ada.
	// ⛔ `,string` — DIKIRIM SEBAGAI TEKS, dan itu BUKAN gaya.
	//
	// Nilainya cap waktu `@CurrentDate("yyyyMMddhhmmssSSS")`
	// (`InsertDocument_Act.xml` b648), yaitu **17 angka** ≈ 2,0e16.
	// `Number.MAX_SAFE_INTEGER` di JavaScript 9.007.199.254.740.991 ≈ 9,0e15
	// — jadi SETIAP pengenal dokumen berada di LUAR jangkauan aman, dan
	// `JSON.parse` membulatkannya diam-diam:
	//
	//	20260927103000123  ->  20260927103000124
	//
	// Akibatnya tautan unduh menunjuk dokumen yang tidak ada, penghapusan
	// mengenai baris yang salah atau tidak ada, dan TIDAK ADA satu pun galat
	// di kedua sisi. Ini cacat lintas-lapis KETUJUH di modul ini dan yang
	// paling buruk: enam pendahulunya membuat layar DIAM, yang ini membuat
	// layar BERBOHONG.
	//
	// ⚠️ Ditemukan uji — `PanelDokumenPeserta.test.ts` membandingkan
	// tautan yang dirakit dengan tautan yang diharap, dan angkanya berbeda
	// satu. Bukan pembacaan ulang yang menemukannya.
	//
	// ⚠️ `Diagnosa.ID` TIDAK berubah, dan itu disengaja: ia dari
	// sequence yang mulai dari 1, jauh di dalam jangkauan aman. Yang
	// menentukan bukan tipe Go-nya melainkan BESAR nilainya.
	ID        int64  `json:"id,string"`
	PesertaID string `json:"pesertaId"`
	NamaFile  string `json:"namaFile"`
	Mime      string `json:"mime"`
	// Kategori1 dan Kategori2 keduanya ditulis InsertDocument_Act; yang
	// ditampilkan Section/DocumentLife.xml kepada manusia adalah Kategori2.
	Kategori1  string `json:"kategori1"`
	Kategori2  string `json:"kategori2"`
	TStorageID string `json:"tStorageId"`
	// ⛔ Kedua tanggal PENUNJUK, dan nil BUKAN tanggal nol. Kolomnya nullable
	// (ADR-U-0027), dan ADR-U-0022 Akibat 2 berbunyi "teks kosong pada kolom
	// angka atau tanggal menjadi KOSONG, bukan nol dan bukan tanggal nol".
	// time.Time biasa tidak dapat membedakan NULL dari 0001-01-01, dan yang
	// membaca kolom itu tidak akan pernah tahu mana yang ia pegang.
	//
	// ⚠️ [terbuka - tiket 14] Medan tanggal lain di model ini (mis.
	// BarisAdjustment.TanggalAkseptasi) masih time.Time biasa dan punya
	// masalah yang sama. Mengubahnya menyentuh pembaca dan penulis yang sudah
	// ada; itu keputusan sekali untuk seluruh model, bukan keputusan tiket 03.
	Tanggal     *time.Time `json:"tanggal"`
	PaymentDate *time.Time `json:"paymentDate"`
}

// PohonKlaim adalah satu klaim utuh, dari akar work object sampai daun.
type PohonKlaim struct {
	Work  WorkClaim
	Klaim Klaim
}

// CacahSpreading menghitung seluruh baris spreading di seluruh pohon.
func (p PohonKlaim) CacahSpreading() int {
	n := 0
	for _, ps := range p.Klaim.Peserta {
		for _, b := range ps.Baris {
			n += len(b.Spreading)
		}
	}
	return n
}

// CacahSpreadingRetro menghitung seluruh baris di tingkat terdalam.
func (p PohonKlaim) CacahSpreadingRetro() int {
	n := 0
	for _, ps := range p.Klaim.Peserta {
		for _, b := range ps.Baris {
			for _, s := range b.Spreading {
				n += len(s.Retro)
			}
		}
	}
	return n
}
