package services

// Keputusan penawaran polis - tiket 01 PremiumList Life.
//
// Untuk apa berkas ini: ketiga tombol keputusan `Confirm` / `Reject` /
// `Decline`, dan penggolong `Offer` / `Premium` yang menyusul `Confirm` di
// tahap penawaran.
//
// Aturannya MURNI dan hidup di `models.TransisiPenawaran`; yang di sini
// gerbang, transaksi, dan jejak.
//
// ⛔ NOL aturan otomatis yang menetapkan keputusan - AC 6 tiket 01.
// Kedua decision table Pega (`IsLifeAccepted`, `IsFlagOnGoingPolicy`)
// mengekspor NOL baris keputusan; yang ditiru AKIBAT keputusan, bukan
// formula yang memilihnya. Karena itu keputusan SELALU datang sebagai
// parameter dari tindakan pengguna, tidak pernah dihitung di sini.
//
// Dibaca sesudah: models/polis_penawaran.go (peta konektornya).

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"nusantarare/internal/models"
	"nusantarare/internal/repository"
)

var (
	// ErrKasusPolisTertutup - kasus sudah selesai, tidak dapat diputus ulang.
	//
	// AC 3 tiket 01: "case tertutup tidak dapat dilanjutkan maupun diputuskan
	// ulang".
	ErrKasusPolisTertutup = errors.New("services: kasus polis sudah ditutup")
	// ErrPenggolongBelumSaatnya - `Offer`/`Premium` di luar tempatnya.
	ErrPenggolongBelumSaatnya = errors.New(
		"services: penggolong Offer/Premium hanya sesudah Confirm di tahap penawaran")
)

// Penawaran melayani keputusan atas penawaran polis.
type Penawaran struct {
	svc      *Service
	jejak    Jejak
	penyalur *PenyalurPolis
}

// Penawaran menyusun layanannya dengan jejak dan penyalur bawaan yang gagal
// terang.
func (s *Service) Penawaran() *Penawaran {
	return &Penawaran{svc: s, jejak: JejakBelumDiputuskan{}, penyalur: penyalurPolisBawaan(s)}
}

// DenganJejak mengganti perekamnya.
func (p *Penawaran) DenganJejak(j Jejak) *Penawaran {
	salin := *p
	salin.jejak = j
	return &salin
}

// DenganPenyalur mengganti penyalur efek keluarnya - tiket 06.
func (p *Penawaran) DenganPenyalur(s *PenyalurPolis) *Penawaran {
	salin := *p
	salin.penyalur = s
	return &salin
}

// pagari menjalankan gerbang bersama ketiga keputusan.
//
// Urutannya: identitas, bentuk permintaan, database, lalu keadaan kasus.
// Permintaan tanpa pengenal harus dijawab "pengenal wajib diisi", bukan
// "ORACLE_DSN belum dikonfigurasi" - cacat yang pernah nyata di Claim Life
// dan ditangkap `TestUbahStatusMenjagaPagarnya`.
func (p *Penawaran) pagari(ctx context.Context, pelaku Pelaku, polisID string) (
	repository.KeadaanPolis, error) {

	var k repository.KeadaanPolis
	if err := WajibIdentitas(pelaku); err != nil {
		return k, err
	}
	if strings.TrimSpace(polisID) == "" {
		return k, fmt.Errorf("%w: pengenal polis wajib diisi", ErrPermintaanTidakSah)
	}
	if p == nil || p.svc == nil || !p.svc.PunyaDatabase() {
		return k, repository.ErrTanpaOracle
	}
	k, err := repository.NewWorkPolis(p.svc.db).Keadaan(ctx, polisID)
	if err != nil {
		return k, err
	}
	// ⛔ Kasus tertutup tidak dapat diputus ulang. Diperiksa DI SINI, satu
	// tempat untuk ketiga keputusan - penjaga yang dipasang di tiga pintu
	// menuju satu ruang adalah tiga tempat untuk lupa.
	if models.KasusPolisTertutup(k.Status) {
		return k, fmt.Errorf("%w: polis %q berstatus %q",
			ErrKasusPolisTertutup, polisID, k.Status)
	}
	return k, nil
}

// ⛔ RALAT 28-09-2026 - TAHAP DIBACA DARI `STATUS`, BUKAN `POSITION`.
//
// Ronde pertama membaca tahap dari kolom `POSITION`. `ProtectAccept.xml`
// membantahnya: ia membandingkan `pyWorkPage.Position` dengan `"Offer"`
// b1207 dan `"Premium"` b2288 - posisi LAYAR, bukan tahap. Tahap adalah
// `pyWorkStatus`, dan itulah kolom `STATUS`.
//
// ⚠️ Bila tidak diralat: setiap keputusan akan dibandingkan dengan
// "Offer"/"Premium", nol di antaranya cocok dengan ketiga nama tahap, dan
// SELURUH permintaan dijawab "tahap polis tidak dikenal" - fitur yang
// hijau di setiap uji murni dan mati pada baris nyata pertama.

// Putuskan menerapkan satu keputusan penawaran.
//
// Mengembalikan akibatnya supaya pemanggil - dan layar - tahu apakah kasus
// berpindah, tertutup, atau menunggu penggolong.
//
// ⚠️ `MenungguPenggolong` TIDAK menulis apa pun. Di flow, `Confirm` di tahap
// penawaran hanya memindahkan kendali ke `Decision3`; yang memindahkan kasus
// adalah hasil penggolong itu. Menuliskan tahap di sini berarti kasus
// berpindah sebelum ada yang memutuskan ke mana.
func (p *Penawaran) Putuskan(ctx context.Context, pelaku Pelaku,
	polisID, keputusan string, saat time.Time) (models.AkibatKeputusan, error) {

	keadaan, err := p.pagari(ctx, pelaku, polisID)
	if err != nil {
		return models.AkibatKeputusan{}, err
	}
	akibat, err := models.TransisiPenawaran(keadaan.Status, keputusan)
	if err != nil {
		return models.AkibatKeputusan{}, err
	}
	if akibat.MenungguPenggolong {
		// Nol tulisan - lihat komentar di atas.
		return akibat, nil
	}
	if _, err := p.terapkan(ctx, pelaku, keadaan, akibat, keputusan, saat); err != nil {
		return models.AkibatKeputusan{}, err
	}
	return akibat, nil
}

// Golongkan menerapkan hasil `IsFlagOnGoingPolicy` sesudah `Confirm`.
//
// ⛔ HANYA dari tahap penawaran. Penggolong itu `Decision3`, dan satu-satunya
// jalan masuk ke sana `Transition4` `[Confirm]` dari `Decision1` - yaitu
// sesudah tahap penawaran. Menerimanya dari tahap lain berarti membuka jalur
// yang flow tidak punya.
func (p *Penawaran) Golongkan(ctx context.Context, pelaku Pelaku,
	polisID, hasil string, saat time.Time) (models.AkibatKeputusan, error) {

	keadaan, err := p.pagari(ctx, pelaku, polisID)
	if err != nil {
		return models.AkibatKeputusan{}, err
	}
	if strings.TrimSpace(keadaan.Status) != models.TahapPolisPenawaran {
		return models.AkibatKeputusan{}, fmt.Errorf("%w: polis %q berada di %q",
			ErrPenggolongBelumSaatnya, polisID, keadaan.Status)
	}
	akibat, err := models.LanjutanPenggolong(hasil)
	if err != nil {
		return models.AkibatKeputusan{}, err
	}
	if _, err := p.terapkan(ctx, pelaku, keadaan, akibat, hasil, saat); err != nil {
		return models.AkibatKeputusan{}, err
	}
	return akibat, nil
}

// terapkan menulis akibat sebuah keputusan, beserta jejaknya.
//
// ⛔ Satu transaksi: perpindahan atau penutupan BERSAMA jejaknya. Jejak yang
// ditulis terpisah dapat hilang sendirian, dan transisi tanpa jejak persis
// yang ADR-0007 larang.
//
// ⛔ TIKET 05b - `SimpanPolis`. Bila jalurnya melewati
// `InsertJsonPolisLife_Act` (`Utility1`, atau `Submit` layar summary), sisa
// activity itu sesudah JSON dibuang - nomor, rekap, salinan peserta warisan -
// berjalan DI DALAM transaksi ini, SEBELUM kasus ditutup. Kasus yang tertutup
// tanpa rekapnya, atau rekap yang tersimpan untuk kasus yang gagal ditutup,
// keduanya tidak mungkin: satu commit untuk semuanya.
func (p *Penawaran) terapkan(ctx context.Context, pelaku Pelaku,
	keadaan repository.KeadaanPolis, akibat models.AkibatKeputusan,
	sebab string, saat time.Time) (HasilSubmitSummary, error) {

	kerja := repository.NewWorkPolis(p.svc.db)
	var simpan HasilSubmitSummary
	err := p.svc.DalamTransaksi(ctx, func(tx *repository.Tx) error {
		if akibat.SimpanPolis {
			var err error
			simpan, err = p.svc.SummaryPremiumList().simpanDalam(ctx, tx, keadaan.ID, saat)
			if err != nil {
				return err
			}
		}
		switch {
		case akibat.Ditutup():
			if err := kerja.TutupKasus(ctx, tx, keadaan.ID, akibat.StatusWork); err != nil {
				return err
			}
		case akibat.TahapTujuan != "":
			if err := kerja.PindahTahap(ctx, tx, keadaan.ID,
				keadaan.Status, akibat.TahapTujuan); err != nil {
				return err
			}
		default:
			return fmt.Errorf("%w: akibat tanpa perpindahan maupun penutupan",
				ErrPermintaanTidakSah)
		}
		// ⚠️ `Dari` dan `Ke` memuat dua kosakata - tahap, dan status akhir -
		// sebagaimana jejak Claim Life. Keduanya "keadaan sebelum" dan
		// "keadaan sesudah"; `sebab` merekam keputusan yang membawanya.
		ke := akibat.TahapTujuan
		if akibat.Ditutup() {
			ke = akibat.StatusWork
		}
		return p.jejak.Rekam(ctx, tx, CatatanJejak{
			KlaimID: keadaan.ID,
			Dari:    keadaan.Status,
			Ke:      ke + " (" + sebab + ")",
			AkunID:  pelaku.AkunID,
			Waktu:   saat,
		})
	})
	if err != nil {
		return HasilSubmitSummary{}, err
	}
	// ⛔ TIKET 06 - efek keluar SESUDAH commit, dan hanya untuk jalur yang
	// menyimpan (`InsertJsonPolisLife_Act` langkah 14-15). Hasilnya TIDAK
	// menjadi galat: premium list sudah tersimpan, dan layanan luar yang
	// gagal tidak boleh membuatnya tampak gagal (ADR-U-0008).
	if akibat.SimpanPolis {
		simpan.EfekKeluar = ringkasEfek(p.penyalur.Salurkan(ctx, MuatanEfek{
			KlaimID: keadaan.ID,
			AkunID:  pelaku.AkunID,
			Waktu:   saat,
		}))
	}
	return simpan, nil
}
