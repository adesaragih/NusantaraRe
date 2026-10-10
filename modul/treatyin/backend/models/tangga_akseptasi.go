package models

import "errors"

// TANGGA AKSEPTASI `Submit` - logika murni, nol basis data.
//
// ---------------------------------------------------------------------
// ⭐ KEPUTUSAN PEMILIK PROSES MENGALAHKAN EKSPOR
// ---------------------------------------------------------------------
// 6 Oktober 2026:
//
//	"untuk SPVTREATY1 · TREATY1 jangan digunakan dulu. alur akseptasinya itu
//	 dari ReasTreatyInAdmin > ReasTreatyInSecHead > ReasTreatyInDeptHead >
//	 ReasTreatyInDirector seperti ini"
//
// `DataTransform/Akseptasi_DT.xml` menentukan langkah PERTAMA dari
// `OperatorID.pyTelephone` berisi kode regu (`SPVTREATY1`, `TREATY1`, ...).
// Cabang itu TIDAK dibangun; langkah pertama ditentukan POSISI.
// Seluruh bandingan ekspor-lawan-keputusan ada di
// `docs/ALUR-AKSEPTASI-DARI-EKSPOR.md` §2.0 dan §2.4.
//
// ---------------------------------------------------------------------
// ⛔ MENGAPA `PositionUsername` TIDAK ADA DI SINI PADA JALUR NAIK
// ---------------------------------------------------------------------
// Ekspor mengisinya dengan nama orang tertanam (`IRVANDY`, `NANDINA`, ...)
// yang HANYA dicapai lewat cabang kode regu yang dibuang di atas. Tanpa
// cabang itu, nol sumber mengisinya, dan menebaknya berarti berkas menunggu
// orang yang salah. Ia pertanyaan terbuka, bukan hal yang lupa dibangun.
//
// ⭐ Pada jalur TURUN ia tetap diketahui - keputusan pemilik proses menyebut
// jalur naiknya saja, jadi `Reject` tetap seperti ekspor. Lihat `AsalNama`.

// Posisi tangga - VERBATIM nama peran di daftar peran aplikasi.
const (
	PosisiKosong   = ""
	PosisiAdmin    = "ReasTreatyInAdmin"
	PosisiSecHead  = "ReasTreatyInSecHead"
	PosisiDeptHead = "ReasTreatyInDeptHead"
	PosisiDirector = "ReasTreatyInDirector"
)

// TanggaAkseptasi - keempat anak tangga, URUT jenjang.
//
// ⚠️ `ReasTreatyInGroupLeader` dan `ReasTreatyInUnderwriting` ada di daftar
// peran aplikasi dan SENGAJA bukan anggota di sini: nol cabang menuju ke sana.
// Itu tidak membuat keduanya boleh dihapus - peran dapat memberi hak lihat
// tanpa pernah menjadi tempat berkas menunggu.
var TanggaAkseptasi = []string{PosisiAdmin, PosisiSecHead, PosisiDeptHead, PosisiDirector}

// Pilihan `TreatyIn.ChooseStatusAkseptasi`.
const (
	PilihAccept  = "Accept"
	PilihReject  = "Reject"
	PilihDecline = "Decline"
)

// Nilai `TreatyIn.StatusAkseptasi`.
const (
	StatusAccept  = "Accept"
	StatusReject  = "Reject"
	StatusDecline = "Decline"
	// StatusTuntas menutup tangga - ⛔ pagar luar seluruh alur: kontrak
	// ber-status ini tidak dapat dikirim ulang, dan itu pula syarat yang
	// menyembunyikan tombol `Submit` di `Section/TreatyInfoSubmit.xml`.
	StatusTuntas = "Resolve Complete"
)

// AsalNama menyatakan DARI MANA `PositionUsername` diisi, bukan isinya.
type AsalNama string

const (
	// AsalNamaKosong - dikosongkan, atau belum diputuskan (jalur naik).
	AsalNamaKosong AsalNama = ""
	// AsalNamaKomentarPertama - `CommentList(1).OperatorName`, jalur biasa.
	AsalNamaKomentarPertama AsalNama = "komentar-pertama"
	// AsalNamaKomentarTerakhir - `CommentList(<LAST>).OperatorName`, revisi.
	AsalNamaKomentarTerakhir AsalNama = "komentar-terakhir"
)

var (
	// ErrPosisiTakDikenal - posisi di luar keempat anak tangga.
	ErrPosisiTakDikenal = errors.New("models: posisi bukan anak tangga akseptasi")
	// ErrPilihanTakBerlaku - pilihan itu tidak punya cabang dari posisi ini.
	ErrPilihanTakBerlaku = errors.New("models: pilihan tak punya cabang dari posisi ini")
	// ErrSudahTuntas - kontrak ber-`Resolve Complete` tidak dapat dikirim ulang.
	ErrSudahTuntas = errors.New("models: kontrak sudah Resolve Complete")
)

// Langkah adalah HASIL satu penekanan - apa yang berubah, bukan bagaimana.
type Langkah struct {
	Posisi   string
	Status   string
	AsalNama AsalNama

	// RevisionState dan ViewState: nil berarti JANGAN SENTUH.
	//
	// ⚠️ Bedanya penting: jalur revisi MENGOSONGKAN keduanya saat tuntas dan
	// MENYETEL `RevisionState` ke "1" saat ditolak, sementara jalur biasa nol
	// menyentuhnya. Nilai kosong dan "tidak disentuh" bukan hal yang sama.
	RevisionState *string
	ViewState     *string
}

func teks(s string) *string { return &s }

// BolehKirim - pagar luar: kontrak tuntas tidak dapat dikirim ulang.
func BolehKirim(statusAkseptasi string) bool { return statusAkseptasi != StatusTuntas }

// AdalahAnakTangga - posisi itu salah satu dari empat.
func AdalahAnakTangga(p string) bool {
	for _, x := range TanggaAkseptasi {
		if x == p {
			return true
		}
	}
	return false
}

// LangkahBerikut menghitung langkah berikutnya.
//
// `posisi` kosong disamakan dengan `ReasTreatyInAdmin`: berkas yang belum
// pernah dikirim ada di tangan pengajunya.
func LangkahBerikut(posisi, pilihan, statusSekarang string, revisi bool) (Langkah, error) {
	if !BolehKirim(statusSekarang) {
		return Langkah{}, ErrSudahTuntas
	}
	if posisi == PosisiKosong {
		posisi = PosisiAdmin
	}
	if !AdalahAnakTangga(posisi) {
		return Langkah{}, ErrPosisiTakDikenal
	}
	if revisi {
		return langkahRevisi(posisi, pilihan)
	}
	return langkahBiasa(posisi, pilihan)
}

// ⛔ `ReasTreatyInAdmin` HANYA punya cabang `Accept`. Pengaju tidak dapat
// menolak atau menampik berkasnya sendiri - nol cabang turun dari sana di
// ekspor maupun di keputusan pemilik proses.
func langkahBiasa(posisi, pilihan string) (Langkah, error) {
	switch pilihan {
	case PilihAccept:
		switch posisi {
		case PosisiAdmin:
			return Langkah{Posisi: PosisiSecHead, Status: StatusAccept}, nil
		case PosisiSecHead:
			return Langkah{Posisi: PosisiDeptHead, Status: StatusAccept}, nil
		case PosisiDeptHead:
			return Langkah{Posisi: PosisiDirector, Status: StatusAccept}, nil
		case PosisiDirector:
			return Langkah{Posisi: PosisiKosong, Status: StatusTuntas}, nil
		}
	case PilihReject:
		if posisi == PosisiAdmin {
			return Langkah{}, ErrPilihanTakBerlaku
		}
		return Langkah{Posisi: PosisiAdmin, Status: StatusReject, AsalNama: AsalNamaKomentarPertama}, nil
	case PilihDecline:
		if posisi == PosisiAdmin {
			return Langkah{}, ErrPilihanTakBerlaku
		}
		return Langkah{Posisi: PosisiKosong, Status: StatusDecline}, nil
	}
	return Langkah{}, ErrPilihanTakBerlaku
}

// Jalur REVISI - ⛔ tangganya LEBIH PENDEK dan itu disengaja: revisi tuntas
// di Sec Head, tidak lewat Dept Head maupun Director. Dan `Reject`-nya memakai
// komentar TERAKHIR, sementara jalur biasa memakai yang PERTAMA.
func langkahRevisi(posisi, pilihan string) (Langkah, error) {
	switch posisi {
	case PosisiAdmin:
		if pilihan != PilihAccept {
			return Langkah{}, ErrPilihanTakBerlaku
		}
		return Langkah{Posisi: PosisiSecHead, Status: StatusAccept}, nil
	case PosisiSecHead:
		switch pilihan {
		case PilihAccept:
			return Langkah{
				Posisi: PosisiKosong, Status: StatusTuntas,
				RevisionState: teks(""), ViewState: teks(""),
			}, nil
		case PilihReject:
			return Langkah{
				Posisi: PosisiAdmin, Status: StatusReject,
				AsalNama: AsalNamaKomentarTerakhir, RevisionState: teks("1"),
			}, nil
		case PilihDecline:
			return Langkah{
				Posisi: PosisiKosong, Status: StatusDecline, RevisionState: teks(""),
			}, nil
		}
		return Langkah{}, ErrPilihanTakBerlaku
	}
	// Dept Head dan Director nol dicapai di jalur revisi.
	return Langkah{}, ErrPilihanTakBerlaku
}
