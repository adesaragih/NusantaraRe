package services

// Reject Outstanding oleh Admin - tiket 05.
//
// Untuk apa berkas ini: jalur penolakan PERTAMA di sistem ini - Admin
// membatalkan baris yang ia input sendiri, tanpa melalui Komite. Jalur kedua
// (lewat Komite) datang di tiket 11.
//
// Dibaca sesudah: statusbaris.go. Mesin transisinya dipakai ulang apa adanya;
// berkas ini hanya menambahkan gerbang-gerbang yang khas penolakan Admin.
//
// Istilah:
//   - Reject Outstanding : membatalkan satu baris adjustment yang masih
//     menunggu keputusan. Ia membatalkan BARIS, bukan klaimnya.

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/claimlife/backend/models"
	"nusantarare/modul/claimlife/backend/repository"
)

// PeranRejectOutstanding adalah peran yang boleh menolak baris Outstanding.
//
// `[terverifikasi]` `Section/AdjustmentDetail_Section.xml` baris 15399 berkas
// pecahan: `pyWorkPage.pyPosition =='ReasLifeAdmin' &&
// pyWorkPage.ClaimData.PremiumListSummary.CLAIM_NO !=” && .STS_REJECT == 0`.
//
// ⚠️ Gerbang itu menguji `.STS_REJECT` TINGKAT BARIS - bukti bahwa wewenang
// pun diukur per baris, bukan per klaim.
const PeranRejectOutstanding = inti.PeranAdmin

var (
	// ErrKlaimBelumBernomor - padanan `CLAIM_NO != ''` pada gerbang XML.
	ErrKlaimBelumBernomor = errors.New("services: klaim belum bernomor")
)

// PeriksaKlaimBernomor menolak klaim yang belum punya nomor.
//
// Nomor klaim adalah bukti klaim itu sudah benar-benar terdaftar. Menolak baris
// pada klaim yang belum bernomor berarti membatalkan sesuatu yang belum ada.
//
// ⚠️ Sampai butir **o** diputuskan, `PenomorBelumDiputuskan` membuat TIDAK ADA
// klaim yang bernomor. Akibatnya jalur HTTP nyata menjawab 422 untuk setiap
// permintaan tolak - dinyatakan, bukan disembunyikan.
func PeriksaKlaimBernomor(k models.Klaim) error {
	if strings.TrimSpace(k.NomorKlaim) == "" {
		return fmt.Errorf("%w: %q", ErrKlaimBelumBernomor, k.ID)
	}
	return nil
}

// Tolak membatalkan satu baris adjustment yang masih Outstanding.
//
// Meniru `Claim Life/Activity/RejectOSClaimLife_Act.xml` langkah 2 (pecahan
// 443-518), yang menulis BERSAMAAN:
//
//	.STS_REJECT                          = 2   (baris adjustment)
//	PremiumListDetail(idx).STS_REJECT    = 2   (peserta pemilik baris)
//	PremiumListDetail(idx).IsCheck       = "false"
//
// ⭐ Pencabutan `IsCheck` adalah AC BARU menurut XML - tiket ini semula tidak
// menyebutnya. Artinya: peserta yang barisnya dibatalkan berhenti terhitung
// sebagai "dipilih untuk diklaim", sehingga ia dapat dipilih ulang dengan
// baris pengganti.
//
// ⛔ Yang TIDAK dilakukan di sini: menulis baris datar warisan. Langkah 5 rule
// itu memanggil `UpdateOsAkseptasiClaimLife_sql`, yang meng-INSERT 55 kolom
// termasuk `POLICY_HOLDER` dan `NAME_OF_INSURED` - nama orang, yang sengaja
// TIDAK dibawa model ini. Memetakan baris adjustment ke baris datar juga belum
// ada; itu pekerjaan tiket 13. `[terbuka - tiket 13]`
//
// ⛔ OQ-M5 DITUTUP (GILIRAN-17): `komentar` = isian `Remarks` dialog Reject
// Outstanding (`RejectOSClaimLife_Sec` b1687, WAJIB b1695), disimpan di
// `T_CLAIMLF_JEJAK.KOMENTAR` pada baris jejak transisi penolakan ini - di
// transaksi yang sama. Pega menyalinnya ke `KomiteList.KomiteComment`
// (`RejectOSClaimLife_Act` b2262-b2263).
func (st *Status) Tolak(ctx context.Context, pelaku inti.Pelaku,
	klaimID, adjID, komentar string, saat time.Time) error {

	// ⛔ Identitas DULU, lalu peran, lalu barulah apa pun dibaca. Urutan ini
	// bukan gaya: permintaan yang tidak berwenang tidak berhak tahu apakah
	// klaimnya ada, dan permintaan tanpa identitas tidak berhak diberi tahu
	// bahwa yang kurang adalah perannya.
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return err
	}
	if err := inti.WajibPeran(pelaku, PeranRejectOutstanding); err != nil {
		return err
	}
	if strings.TrimSpace(klaimID) == "" || strings.TrimSpace(adjID) == "" {
		return fmt.Errorf("%w: pengenal klaim dan baris wajib diisi", galat.ErrPermintaanTidakSah)
	}
	komentar = strings.TrimSpace(komentar)
	if komentar == "" {
		return fmt.Errorf("%w: Remarks wajib diisi", galat.ErrPermintaanTidakSah)
	}
	if len(komentar) > BatasKomentarJejak {
		return fmt.Errorf("%w: Remarks melebihi %d byte", galat.ErrPermintaanTidakSah, BatasKomentarJejak)
	}
	if !st.svc.PunyaDatabase() {
		return db.ErrTanpaOracle
	}

	baca := repository.NewKlaimLife(st.svc.DB())
	hdr, err := baca.AmbilHeader(ctx, klaimID)
	if err != nil {
		return err
	}
	if hdr == nil {
		return fmt.Errorf("%w: klaim %q tidak ada", galat.ErrPermintaanTidakSah, klaimID)
	}
	if err := PeriksaKlaimBernomor(*hdr); err != nil {
		return err
	}

	pesertaID, err := st.pemilikBaris(ctx, baca, klaimID, adjID)
	if err != nil {
		return err
	}

	// Mesin transisi tiket 04 dipakai APA ADANYA - termasuk kefinalannya,
	// pencerminan tiga tingkat, dan jejak auditnya. Nol aturan status ditulis
	// ulang di sini: dua salinan aturan kefinalan berarti dua kesempatan untuk
	// berbeda.
	//
	// ⛔ Pencabutan IS_CHECK disisipkan ke DALAM transaksi yang sama, bukan
	// menyusul sesudahnya. Rule Pega menulis ketiganya dalam satu
	// Property-Set, dan dua transaksi berarti peserta dapat tertinggal masih
	// "dipilih" padahal barisnya sudah batal.
	return st.ubah(ctx, pelaku, klaimID, pesertaID, adjID,
		kontrak.StatusDitolak, saat, true, komentar)
}

// BatasKomentarJejak - lebar `T_CLAIMLF_JEJAK.KOMENTAR` VARCHAR2(4000) dalam
// BYTE (migrasi 021). Diperiksa di pintu supaya kelebihan menjadi 400, bukan
// ORA-12899 di tengah transaksi.
const BatasKomentarJejak = 4000

// pemilikBaris mencari peserta yang memiliki sebuah baris adjustment.
func (st *Status) pemilikBaris(ctx context.Context, baca *repository.KlaimLife,
	klaimID, adjID string) (string, error) {
	perBaris, err := baca.AmbilBaris(ctx, klaimID)
	if err != nil {
		return "", err
	}
	for pesertaID, daftar := range perBaris {
		for _, b := range daftar {
			if b.ID == adjID {
				return pesertaID, nil
			}
		}
	}
	return "", fmt.Errorf("%w: baris %q bukan milik klaim %q",
		galat.ErrPermintaanTidakSah, adjID, klaimID)
}
