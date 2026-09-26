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

import "time"

// Lini adalah penanda lini usaha pada work object. Untuk Claim Life isinya
// tetap LiniLife.
//
// [terbuka - Non-Life] Daftar nilai enum lintas-lini ditetapkan saat konteks
// Non-Life digarap. Yang sudah pasti hanyalah nilai untuk Life.
const LiniLife = "LIFE"

// WorkClaim adalah akar pohon: satu baris work object, lintas-lini.
//
// ID-nya TEKS BERFORMAT - CLM-xxxxxx untuk baris klaim, KMT-xxxxxx untuk baris
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
	// RetrocadedShare dan Rate adalah RASIO, bukan uang, dan karena itu tidak
	// dapat dijumlahkan dengan Money (ADR-F-0004).
	RetrocadedShare Ratio
	Rate            Ratio
	// IDR dan USD adalah nilai uang; nama kolomnya sekaligus mata uangnya.
	IDR      Money
	USD      Money
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
	PercentShare  Ratio
	Amount        Money
	Rate          Ratio
	// [terbuka] Rumus PremiumSpreadedNet punya dua cabang di rule yang sama;
	// pemiliknya Product dan Underwriting. Di sini ia hanya disimpan apa
	// adanya, tidak dihitung ulang.
	PremiumSpreadedGross Money
	PremiumSpreadedNet   Money
	Commision            Money
	OvrComm              Money
	TreatyTypeID         string
	TreatyTypeName       string
}

// Dokumen adalah satu dokumen pendukung milik seorang peserta.
//
// [terbuka] [data DBA] Kolom isinya TIDAK DAPAT DITURUNKAN dari korpus - kelas
// Pega-nya ASM-FW-GCNMFW-Int-DOCUMENT_CLAIM dan SQL-nya dibuat Pega sendiri.
// Yang pasti hanya identitas dan induknya. Menambah medan lain berarti
// mengarang.
type Dokumen struct {
	ID        int64
	PesertaID string
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
