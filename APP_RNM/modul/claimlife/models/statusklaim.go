package models

import (
	"nusantarare/inti/backend/kontrak"
)

// Status klaim sebagai TURUNAN - tiket 04.
//
// Untuk apa berkas ini: status klaim DIHITUNG dari status baris-barisnya,
// tidak pernah disimpan sebagai keadaan mandiri. Unit keputusan adalah baris
// `AdjustmentList` (ADR-U-0011), dan klaim hanyalah wadahnya.
//
// Dibaca sesudah: klaimlife.go.
//
// ⚠️ Kolom `STS_REJECT` dan `ACCEPTED_NO` di `T_GENERAL_CLAIM` memang ada,
// tetapi ia CERMIN baris terakhir - bukan sumber kebenaran, dan bukan yang
// dibaca fungsi di berkas ini. Sistem lama menyimpannya; sistem ini
// menghitungnya dan menulis cerminnya untuk pembaca hilir.

// StatusKlaim adalah keadaan sebuah klaim, diturunkan dari baris-barisnya.
type StatusKlaim int

const (
	// KlaimBelumBerbaris: belum ada satu pun baris keputusan.
	KlaimBelumBerbaris StatusKlaim = iota
	// KlaimBerjalan: masih ada baris yang menunggu diputuskan.
	KlaimBerjalan
	// KlaimSelesai: tidak ada lagi yang menunggu, dan ada yang diaksep.
	KlaimSelesai
	// KlaimDitolakSeluruhnya: tidak ada lagi yang menunggu, dan tidak satu
	// pun diaksep. ⛔ Ini BUKAN keadaan akhir: klaim tetap dapat menerima
	// baris baru, dan barisnya kembali berjalan.
	KlaimDitolakSeluruhnya
	// KlaimTidakDapatDipastikan: ada baris berkode yang TIDAK DIKENAL.
	//
	// ⛔ Kode "4" ada di data warisan dan artinya belum diputuskan siapa pun.
	// Menghitungnya sebagai "ditolak seluruhnya" - yang dilakukan ronde
	// pertama berkas ini - membuat klaim tampak SUDAH diputus padahal justru
	// keputusannya yang tidak terbaca. Arah kekeliruan yang paling berbahaya.
	KlaimTidakDapatDipastikan
)

// String menulis status klaim sebagai kata yang dibaca pengguna.
func (s StatusKlaim) String() string {
	switch s {
	case KlaimBerjalan:
		return "Berjalan"
	case KlaimSelesai:
		return "Selesai"
	case KlaimDitolakSeluruhnya:
		return "Ditolak seluruhnya"
	case KlaimTidakDapatDipastikan:
		return "Tidak dapat dipastikan"
	default:
		return "Belum berbaris"
	}
}

// StatusTurunan menghitung status klaim dari seluruh baris adjustment-nya.
//
// Aturannya `[keputusan work owner 2026-09-14, CONTEXT.md Lampiran butir 1]`:
//
//	ada baris Outstanding                 -> Berjalan
//	nol Outstanding dan ada Aksep         -> Selesai
//	nol Outstanding dan nol Aksep         -> Ditolak seluruhnya
//	nol baris                             -> Belum berbaris
//
// ⛔ Kode yang TIDAK DIKENAL - misalnya "4" di data warisan - tidak dihitung
// sebagai Outstanding maupun sebagai Aksep. Menebaknya ke salah satu sisi
// berarti melaporkan keputusan yang tidak pernah diambil siapa pun.
func (k Klaim) StatusTurunan() StatusKlaim {
	var cacah, outstanding, aksep, takDikenal int
	for _, p := range k.Peserta {
		for _, b := range p.Baris {
			cacah++
			switch b.Status() {
			case kontrak.StatusOutstanding:
				outstanding++
			case kontrak.StatusAksep:
				aksep++
			case kontrak.StatusDitolak:
				// ditolak: tidak menambah apa pun, tetapi ia DIKENAL
			default:
				takDikenal++
			}
		}
	}
	switch {
	case cacah == 0:
		return KlaimBelumBerbaris
	case outstanding > 0:
		return KlaimBerjalan
	case aksep > 0:
		return KlaimSelesai
	case takDikenal > 0:
		// ⛔ Diperiksa SESUDAH Outstanding dan Aksep: baris berkode asing tidak
		// menghapus fakta bahwa ada yang masih menunggu atau sudah diaksep.
		return KlaimTidakDapatDipastikan
	default:
		return KlaimDitolakSeluruhnya
	}
}
