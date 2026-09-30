package services

// Perpindahan tahap dan jalur balik - tiket 08.
//
// Untuk apa berkas ini: memindahkan kasus di tangga kerja, dan mengembalikannya
// ke peran sebelumnya. Keduanya TIDAK menyentuh status baris adjustment.
//
// Dibaca sesudah: wewenang.go.
//
// ⛔ Tahap BUKAN status. Memindahkan kasus ke Medical Check tidak memutuskan
// apa pun tentang barisnya (ADR-U-0011), dan berkas ini karena itu tidak
// memuat satu pun tulisan ke KodeStatus - dijaga penjaga statik tiket 07.

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/jejak"
	"nusantarare/modul/claimlife/models"
	"nusantarare/modul/claimlife/repository"
)

var (
	// ErrTahapTidakDikenal - tahap tujuan di luar keempat yang dikenal.
	ErrTahapTidakDikenal = errors.New("services: tahap tidak dikenal")
	// ErrPeranTahapBelumDiputuskan - tahap yang pemegangnya belum ditetapkan.
	ErrPeranTahapBelumDiputuskan = errors.New(
		"services: peran pemegang tahap belum diputuskan work owner")
)

// ErrPerpindahanTidakSah - langkah tahap yang tidak ada di tangga kerja.
var ErrPerpindahanTidakSah = errors.New("services: perpindahan tahap tidak sah")

// JalurBalik adalah penanda pengembalian kasus ke peran sebelumnya.
//
// `[terverifikasi]` `When/IsSendtoAdmin.xml` `pyConditionString`:
// `pyWorkPage.SendtoAdmin = 1`; `When/IsSendtoMedical.xml`:
// `pyWorkPage.SendtoMedical = 1`.
//
// ⭐ Tiket ini semula menandai `IsSendtoMedical` *"tidak terbaca dari tag"*.
// Ia TERBACA - kondisinya memang tidak ada di label rule, melainkan di
// `pyConditionString` dan `pyConditionValue1`
// (`compareTwoValues(pyWorkPage.SendtoMedical, "=", 1)`). Simetris dengan
// pasangannya, dan butir itu `[ditutup oleh XML - 26-09-2026]`.
type JalurBalik struct {
	KeAdmin   bool
	KeMedical bool
}

// NilaiSendto menulis kedua penanda sebagai teks kolom.
//
// Nilainya "1" atau kosong, seperti sistem lama - teks, bukan bilangan
// (ADR-U-0022), dan kosong berarti tidak sedang dikembalikan.
func (j JalurBalik) NilaiSendto() (admin, medical string) {
	if j.KeAdmin {
		admin = "1"
	}
	if j.KeMedical {
		medical = "1"
	}
	return admin, medical
}

// TahapLayanan memindahkan kasus di tangga kerja.
type TahapLayanan struct {
	svc   *Service
	jejak jejak.Jejak
}

// Tahap menyusun layanan itu dengan jejak bawaan yang gagal terang.
//
// ⛔ ADR-U-0007 menyebut jalur balik SECARA KHUSUS: "merekam siapa dan kapan
// untuk setiap transisi status klaim DAN setiap jalur balik (SendtoAdmin,
// SendtoMedical)". Perpindahan tahap karena itu tidak boleh terjadi tanpa
// terekam, dan tabelnya belum ada (butir am).
func (s *Service) Tahap() *TahapLayanan {
	return &TahapLayanan{svc: s, jejak: jejak.JejakBelumDiputuskan{}}
}

// DenganJejak mengganti perekamnya - dipakai test, dan kelak oleh tiket 09.
func (tl *TahapLayanan) DenganJejak(j jejak.Jejak) *TahapLayanan {
	return &TahapLayanan{svc: tl.svc, jejak: j}
}

// tahapKasus membaca tahap BERLAKU sebuah kasus; tak dikenal adalah GALAT.
//
// ⛔ Satu pintu (GILIRAN-11 paket 4): diagnosa, dialog Edit Date, perpindahan,
// dan penutupan dahulu masing-masing menyalin empat belas baris yang sama.
func tahapKasus(ctx context.Context, baca *repository.KlaimLife, klaimID string) (models.Tahap, error) {
	kolomTahap, peranPemegang, err := baca.TahapDanPeran(ctx, klaimID)
	if err != nil {
		return models.TahapTidakDikenal, err
	}
	t := models.TahapBerlaku(kolomTahap, peranPemegang)
	if !t.Diketahui() {
		return t, fmt.Errorf("%w: tahap %q, peran pemegang %q",
			ErrTahapTidakDikenal, kolomTahap, peranPemegang)
	}
	return t, nil
}

// Pindah memindahkan kasus ke tahap lain, dengan atau tanpa jalur balik.
//
// ⛔ Status baris TIDAK disentuh. Itu bukan kelalaian melainkan aturannya:
// perpindahan tahap memindahkan PEKERJAAN, bukan memutuskan klaim
// (ADR-U-0011). Penjaga statik tiket 07 memastikan berkas ini tidak pernah
// menulis KodeStatus.
//
// ⚠️ TANPA argumen jalur balik. Ronde pertama menerimanya sebagai bendera
// bebas, sehingga `SENDTO_ADMIN=1` dapat ditulis pada perpindahan menuju
// Medical Check, dan keduanya dapat menyala sekaligus - padahal keduanya
// menunjuk tujuan yang berbeda. Kini ia DITURUNKAN dari pasangan tahapnya.
func (tl *TahapLayanan) Pindah(ctx context.Context, pelaku inti.Pelaku,
	klaimID string, ke models.Tahap, saat time.Time) error {

	if err := inti.WajibIdentitas(pelaku); err != nil {
		return err
	}
	if !ke.Diketahui() {
		return fmt.Errorf("%w: tujuan %v", ErrTahapTidakDikenal, ke)
	}
	if strings.TrimSpace(klaimID) == "" {
		return fmt.Errorf("%w: pengenal klaim wajib diisi", galat.ErrPermintaanTidakSah)
	}
	if !tl.svc.PunyaDatabase() {
		return db.ErrTanpaOracle
	}

	// ⛔ BUTIR bb: kasus yang sudah ditutup tidak dapat diubah lagi.
	// Satu pintu untuk seluruh rute pengubah - lihat
	// services.PastikanKasusTerbuka, yang pemanggilannya ditagih penjaga
	// statik. Diperiksa SESUDAH wewenang: pemanggil yang tidak berhak tidak
	// berhak pula tahu keadaan kasusnya.
	if err := tl.svc.PastikanKasusTerbuka(ctx, klaimID); err != nil {
		return err
	}

	// ⛔ Yang diperiksa adalah peran tahap ASAL, bukan tahap tujuan. Ronde
	// pertama memeriksa tujuan - dan itu membalik seluruh jalur balik:
	// pengembalian ke Admin oleh Medical Advisor akan menuntut pelakunya
	// berperan Admin, yang justru bukan dia. Orang memindahkan pekerjaan yang
	// SEDANG IA PEGANG.
	baca := repository.NewKlaimLife(tl.svc.DB())
	// ⛔ BUTIR at: TAHAP dan PY_POSITION dibaca BERSAMA. Yang tersimpan
	// di PY_POSITION adalah NAMA PERAN, bukan pengenal shape - ronde
	// pertama menganggapnya `"Assignment<n>"` dan setiap pembacaan baris
	// nyata berakhir "tidak dikenal".
	// ⚠️ Kolom TAHAP menang; PY_POSITION hanya CADANGAN untuk baris lama.
	// Baris lama yang sebenarnya di Input Register akan tampak Outstanding
	// sampai kolomnya terisi - diterima, dan dicatat di models.TahapDariPeran.
	asal, err := tahapKasus(ctx, baca, klaimID)
	if err != nil {
		return err
	}
	peranAsalTahap, ada := models.PeranPemegangTahap(asal)
	if !ada {
		return fmt.Errorf("%w: tahap %q", ErrPeranTahapBelumDiputuskan, asal)
	}
	// Orang memindahkan pekerjaan yang SEDANG IA PEGANG.
	if err := inti.WajibPeran(pelaku, peranAsalTahap); err != nil {
		return err
	}
	peranTujuan, ada := models.PeranPemegangTahap(ke)
	if !ada {
		return fmt.Errorf("%w: tahap %q", ErrPeranTahapBelumDiputuskan, ke)
	}
	// ⛔ Perpindahan diperiksa antar-TAHAP, bukan antar-PERAN: peta peran
	// tidak dapat menyatakan Input Register ⇄ Outstanding Claim, yang
	// keduanya dipegang ReasLifeAdmin.
	if !models.SerahTerimaSah(asal, ke) {
		return fmt.Errorf("%w: %s -> %s", ErrPerpindahanTidakSah, asal, ke)
	}
	keAdmin, keMedical := models.JalurBalikTahap(asal, ke)
	admin, medical := JalurBalik{KeAdmin: keAdmin, KeMedical: keMedical}.NilaiSendto()
	return tl.svc.DalamTransaksi(ctx, func(tx *db.Tx) error {
		if err := baca.PerbaruiTahap(ctx, tx, klaimID, asal.String(), ke.String(),
			peranTujuan, admin, medical, saat); err != nil {
			return err
		}
		return tl.jejak.Rekam(ctx, tx, jejak.CatatanJejak{
			// ⛔ KlaimID, bukan AdjustmentID: yang berpindah KASUSNYA.
			KlaimID: klaimID,
			// ⚠️ Jejak mencatat TAHAP, bukan peran: perpindahan Admin→Admin
			// punya peran asal dan tujuan yang SAMA, dan jejak yang
			// mencatat peran akan berbunyi "dari Admin ke Admin".
			Dari:   asal.String(),
			Ke:     ke.String(),
			AkunID: pelaku.AkunID,
			Waktu:  saat,
		})
	})
}
