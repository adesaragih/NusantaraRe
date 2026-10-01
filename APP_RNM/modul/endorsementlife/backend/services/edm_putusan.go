package services

// Tiket 08 (alur keputusan), 04 (nomor endorsement), 09 (satu transaksi +
// anti-dobel), 11 (rekap warisan) - dua tombol `Submit` `InputEDMLife.xml`
// (b37494 `EmailTypePL =1`, b38109 `= 2 ||… = 7`), keduanya `VIS .IsJsonPolis=1`.
//
// ⛔ URUTAN ADALAH BAGIAN KEBENARAN (AC 42), satu transaksi (AC 70):
//
//	1 kunci kasus, terbuka, sudah `Save`          (wadah b35518 `.IsJsonPolis=1`)
//	2 riwayat `AddHistorySuggest`                 (b37776 / b38314 - kedua tombol)
//	3 `IsLifeAccepted` b272: `1` → Confirm, selain itu Decline (tanpa `Reject`, AC 30)
//	Decline: 4 `STATUSS = Resolved-Rejected` (b642), 5 jejak
//	Confirm: 4 versi berjalan + nomor `<polis>/NN` (`GenerateNoEDM_Life`, lahir sekali b1216/b1390)
//	         5 anti-dobel `(NO_POLIS, PROD_KE)`  6 resmikan kepala + peserta (11.2-11.5)
//	         7 rekap warisan (12, `InsertPLSummary`)  8 jejak `Resolved-Completed` (b686)
//	         sesudah commit: efek keluar 16 (`edm_efekkeluar.go`)
//
// ⛔ Penulisan peserta warisan `M_LIFE_PREMIUM_DETAIL` (11.6 `SaveMasterLPDet`)
// BELUM: menunggu OQ-EDM-016 (RALAT R29). `JSON_POLIS` dan `LIFEINPRODUCTION`
// tidak ditulis (R17, R18).

import (
	"context"
	"errors"
	"fmt"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/jejak"
	"nusantarare/inti/backend/outbox"
	"nusantarare/inti/backend/penomor"
	"nusantarare/modul/endorsementlife/backend/models"
	"nusantarare/modul/endorsementlife/backend/repository"
)

var (
	// ErrBelumDisimpan - `ConfirmSection` dan kedua `Submit` tampil hanya bila `.IsJsonPolis=1`.
	ErrBelumDisimpan = errors.New("services: press Save before submitting the decision")
	// ErrVersiBerubah - penjaga anti-dobel: versi polis yang dituju sudah ada, atau
	// versi berjalan berubah sejak kasus dibuat (salinannya basi).
	ErrVersiBerubah = errors.New("services: the policy already has this endorsement version, or a newer version was confirmed after this case was created; decline this case and create a new one")
)

// MasukanPutusan - radio `Status` (`EmailTypePL`, `ConfirmSection` b496) dan `Comment` b829.
type MasukanPutusan struct {
	Status   string `json:"status"`
	Komentar string `json:"comment"`
}

// HasilPutusan - status akhir dan `No. Endorsement` (`ConfirmSubmitEDM` b654).
type HasilPutusan struct {
	Status        string `json:"status"`
	NoEndorsement string `json:"noEndorsement"`
	Peserta       int    `json:"peserta"`
	RekapWarisan  int    `json:"rekapWarisan"`
	// EfekKeluar - Arasapas sesudah Confirm (tiket 10); kosong pada Decline.
	EfekKeluar *RingkasEfek `json:"efekKeluar,omitempty"`
}

// Putuskan menjalankan keputusan kasus dalam satu transaksi.
func (l *Layanan) Putuskan(ctx context.Context, p inti.Pelaku, id string, m MasukanPutusan) (HasilPutusan, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return HasilPutusan{}, err
	}
	if !models.KasusEDM(id) {
		return HasilPutusan{}, fmt.Errorf("%w: %q", ErrKasusTidakAda, id)
	}
	if err := models.PeriksaKeputusan(m.Status); err != nil {
		return HasilPutusan{}, fmt.Errorf("%w: %v", ErrMasukanTidakSah, err)
	}
	komentar := strings.TrimSpace(m.Komentar)
	if len(komentar) > models.BatasKomentar {
		return HasilPutusan{}, fmt.Errorf("%w: comment is longer than %d bytes", ErrMasukanTidakSah, models.BatasKomentar)
	}
	var hasil HasilPutusan
	err := l.tx(ctx, func(tx *db.Tx) error {
		k, err := l.gudang.AmbilKasus(ctx, tx, id, true)
		if errors.Is(err, repository.ErrTidakAda) {
			return fmt.Errorf("%w: %q", ErrKasusTidakAda, id)
		}
		if err != nil {
			return err
		}
		if !k.Terbuka() {
			return ErrKasusTertutup
		}
		sudah, err := l.gudang.AdaRekap(ctx, tx, id)
		if err != nil {
			return err
		}
		if !sudah {
			return ErrBelumDisimpan
		}
		waktu := l.jam()
		if err := l.gudang.SisipRiwayat(ctx, tx, models.RiwayatTulis{
			KasusID: id, Waktu: penomor.DiJakarta(waktu).Format("2006-01-02 15:04:05"), PIC: p.AkunID,
			Status: models.LabelKeputusanRiwayat(m.Status), Komentar: komentar,
		}); err != nil {
			return err
		}
		catat := func(ke string) error {
			return l.jejak.Rekam(ctx, tx, jejak.CatatanJejak{
				KlaimID: id, Dari: models.TahapInputEDMLife, Ke: ke, AkunID: p.AkunID, Waktu: waktu, Komentar: komentar,
			})
		}
		if !models.Diterima(m.Status) {
			if err := l.gudang.Tolak(ctx, tx, id); err != nil {
				return err
			}
			hasil.Status = models.StatusKasusDitolak
			return catat(models.StatusKasusDitolak)
		}

		v, ada, err := l.gudang.VersiBerjalan(ctx, tx, k.NomorPolis, 0)
		if err != nil {
			return err
		}
		if !ada {
			return fmt.Errorf("%w: no current version of %q", ErrVersiBerubah, k.NomorPolis)
		}
		nomor, prodKe, err := models.NomorEndorsement(k.NomorPolis, v)
		if err != nil {
			return err
		}
		if prodKe != k.ProdKe {
			return fmt.Errorf("%w: case version %d, next version is %d", ErrVersiBerubah, k.ProdKe, prodKe)
		}
		if k.PLNumberEDM != "" {
			nomor = k.PLNumberEDM // nomor lahir sekali (prakondisi `PL_NUMBER_EDM==""`)
		}
		// Pertahanan berlapis: di jalur modul ini dobel sudah dicegah index unik
		// kasus terbuka + kunci baris; ini menangkap versi yang ditulis jalur lain.
		dobel, err := l.gudang.AdaVersiResmi(ctx, tx, k.NomorPolis, prodKe)
		if err != nil {
			return err
		}
		if dobel {
			return fmt.Errorf("%w: %s version %d exists", ErrVersiBerubah, k.NomorPolis, prodKe)
		}
		if hasil.Peserta, err = l.gudang.Resmikan(ctx, tx, models.ResmiKasus{
			ID: id, NomorPolis: k.NomorPolis, ProdKe: prodKe, Nomor: nomor, StatusJenis: models.StatusJenis(k.Kepala["TYPE"]),
		}); err != nil {
			return err
		}
		if hasil.RekapWarisan, err = l.gudang.TulisRekapWarisan(ctx, tx, repository.RekapWarisanTulis{
			KasusID: id, NomorPolis: k.NomorPolis, Nomor: nomor, COB: k.Kepala["BUSINESS_NAME"],
		}); err != nil {
			return err
		}
		hasil.Status, hasil.NoEndorsement = models.StatusKasusSelesai, nomor
		return catat(models.StatusKasusSelesai)
	})
	if err != nil {
		return HasilPutusan{}, err
	}
	if hasil.Status == models.StatusKasusSelesai {
		// Sesudah commit: kegagalan tercatat di outbox dan tampil, versi resmi tidak dibatalkan.
		r := ringkas(l.penyalur.Salurkan(ctx, outbox.MuatanEfek{KlaimID: id, AkunID: p.AkunID, Waktu: l.jam()}))
		for _, g := range append(append([]string{}, r.Gagal...), r.TidakDiantre...) {
			l.catat("endorsement life: alarm efek keluar kasus " + id + ": " + g)
		}
		hasil.EfekKeluar = &r
	}
	return hasil, nil
}
