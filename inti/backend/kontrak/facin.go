package kontrak

// Kontrak mesin Fac In bersama - keputusan work owner 01-10-2026 (nbfacin
// `docs/KEPUTUSAN-30-09-2026.md` butir 55).
//
// Mesin perhitungan, registry predikat, dan tangga akseptasi Fac In milik
// New Business (`modul/nbfacin`). Renewal (`modul/rnwfacin`) - kelak juga
// Endorsement - memakai ulang mesin yang sama tanpa salinan: `[terverifikasi]`
// 1.907 dari 1.927 berkas korpus RNW identik byte dengan NB (K-031). Karena modul
// tidak saling mengimpor, yang menyeberang hanyalah antarmuka dan tipe di bawah.
//
// Disediakan nbfacin (`backend/services/kontrakfacin`), dipakai rnwfacin (dan
// endorsmentfacin). Penyambungan lewat `Pendaftaran()` menyusul layar pertama
// modulnya (A32); sampai itu pemakai merakitnya sendiri, mis. `uji/lintasmodul`.

import "nusantarare/inti/backend/uang"

// KasusFacIn - halaman kerja kasus Fac In sebagai jalur properti Pega → nilai,
// mis. `pyWorkPage.OfferFacIn.QuotationData.TeamGroup`. Jalur yang tidak ada
// dibaca kosong, seperti clipboard Pega. Sama bentuknya dengan rules.Kasus nbfacin.
type KasusFacIn interface {
	Nilai(jalur string) (nilai string, ada bool)
}

// JabatanFacIn - kode jabatan tangga (`pyWorkPage.LetterNo`). Berbeda tipe dari
// AntreanFacIn supaya menukar keduanya gagal saat kompilasi (tiket NB-11).
type JabatanFacIn string

// AntreanFacIn - workbasket yang memegang kasus.
type AntreanFacIn string

// PenggunaFacIn - pengguna yang baru memutuskan, dari model peran (butir 33).
type PenggunaFacIn struct {
	Jabatan     JabatanFacIn
	AnggotaGrup bool
}

// TransisiFacIn - hasil SATU langkah tangga akseptasi (bentuk A atau B).
type TransisiFacIn struct {
	// Selesai - tangga berakhir; penyelesaian normal, bukan galat.
	Selesai             bool
	JabatanTujuan       JabatanFacIn
	Antrean             AntreanFacIn
	PositionNoteDitulis bool
}

// TanggaAkseptasiFacIn menjalankan satu langkah tangga akseptasi atas kasus;
// penyedia memilih sendiri bentuk A (tabel M_LIMIT_* bentuk A) atau bentuk B
// (M_LIMIT_FINANCIALINS) menurut predikat kasus.
type TanggaAkseptasiFacIn interface {
	Langkah(k KasusFacIn, p PenggunaFacIn) (TransisiFacIn, error)
}

// PenilaiPredikatFacIn menilai satu predikat `When` registry NB atas kasus.
type PenilaiPredikatFacIn interface {
	Eval(nama string, k KasusFacIn) (bool, error)
}

// MasukanPremiFacIn - masukan perhitungan premi satu coverage; medan dan artinya
// sama dengan `premium.Input` nbfacin (teks apa adanya dari halaman kerja).
type MasukanPremiFacIn struct {
	LiniBisnis             string // label K-018: PA, FIRE, MBU, ANEKA, BONDING, GOLF, MARINE CARGO
	CalculateMethod        string
	MataUang               string
	TSI                    string
	Rate                   string
	ProRatePercent         string
	PctShortPeriod         string
	Discount               string
	DiscountType           string
	DiscountPercentage     string
	PremiSebelumnya        string
	Loading                string
	ProRatePercentCoverage string
	// Medan FIRE/ANEKA/GOLF/MARINE CARGO (tiket nbfacin 18), aditif 01-10-2026.
	IndemnityPercentage string
	LossLimit           string // `.LostLimit` coverage (ejaan korpus)
	PctAdjustment       string
	CoverageBasis       string
	NetRate             string
	FirstScale          string
	MBD                 bool // predikat IsMBD kasus
	MasterPolicy        bool // MARINE: PolicyType == 1 && IsMOP == "MOP"
}

// LiniFacIn - satu lini bisnis yang gerbang pemilih rumusnya terbuka, beserta
// pembagi satuan rate-nya (resolver K-018: ‰ = 1000, % = 100).
type LiniFacIn struct {
	Lini         string
	PembagiRate  int64
	SimbolSatuan string
}

// HasilLiniFacIn - lini kasus dari gerbang predikat; Peringatan = jalur
// BusinessType yang berbeda nilai.
type HasilLiniFacIn struct {
	Lini       []LiniFacIn
	Peringatan []string
}

// MesinPremiFacIn - perhitungan premi dan resolver lini bisnis NB.
type MesinPremiFacIn interface {
	// Hitung - premi satu coverage, rumus lini sistem lama dan pembulatannya.
	Hitung(in MasukanPremiFacIn) (uang.Money, error)
	// AsalRumus - rule dan langkah sistem lama yang dipakai Hitung; kosong bila
	// lini/metodenya belum diport.
	AsalRumus(in MasukanPremiFacIn) string
	// LiniDariPredikat - lini bisnis kasus dari gerbang pemilih rumus.
	LiniDariPredikat(k KasusFacIn) (HasilLiniFacIn, error)
}
