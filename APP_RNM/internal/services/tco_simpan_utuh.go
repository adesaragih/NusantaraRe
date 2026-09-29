package services

// Simpan atomik satu kontrak lintas enam tabel - tiket 09 Treaty Contract Out.
//
// Untuk apa berkas ini: kontrak, reinsurer (+ security-nya), business, dan
// klausul dari SATU permintaan ditulis dalam SATU transaksi; kegagalan di mana
// pun membatalkan seluruhnya (AC 37). Penulis per baris yang sudah ada (tiket
// 04-08) dipakai APA ADANYA - aturan tidak disalin - dengan dua penyesuaian
// yang diberikan repository:
//
//   - `repository.DenganTransaksiUtuhTCO`: pembaca modul melihat tulisan
//     permintaan ini (reinsurer baru terbaca oleh security-nya, induk klausul
//     baru oleh anaknya);
//   - identitas SEMENTARA sampai seluruh baris lolos; baru sesudah itu nomor
//     sequence diambil (`TetapkanIdentitasTCO`). Kegagalan tidak menghabiskan
//     satu nomor pun (AC tiket 09).
//
// `[data DBA]` keenam prosedur penulis tidak COMMIT sendiri; COMMIT yang
// terlihat di rule Connect-SQL Pega ditulis Pega. Di sini commit SEKALI, oleh
// `DalamTransaksi`, di akhir.
//
// ⚠️ Urutan dalam permintaan dipakai apa adanya: kontrak -> reinsurer (tiap
// reinsurer disusul security-nya) -> business -> klausul (induk sebelum anak).
//
// Dibaca sesudah: tco_kontrak.go ... tco_klausul.go, repository/tco_transaksi_utuh.go.

import (
	"context"
	"errors"
	"fmt"
	"time"

	"nusantarare/internal/models"
	"nusantarare/internal/repository"
)

// batasBarisUtuhTCO - batas jumlah baris satu permintaan (transaksi tidak dibiarkan tak berbatas).
const batasBarisUtuhTCO = 500

var (
	// ErrSimpanUtuhTidakSah - bentuk permintaan simpan utuh tidak sah (400).
	ErrSimpanUtuhTidakSah = errors.New("services: permintaan simpan utuh tidak sah")
	// ErrPenetapIdentitasBelumDisuntik - penetap identitas belum dipasang.
	ErrPenetapIdentitasBelumDisuntik = errors.New("services: penetap identitas transaksi utuh belum disuntik")
)

// GalatSimpanUtuhTCO menyebut baris MANA yang menggagalkan simpan (AC 38).
type GalatSimpanUtuhTCO struct {
	Bagian string
	Galat  error
}

func (g GalatSimpanUtuhTCO) Error() string {
	// Identitas sementara tidak berarti bagi pemakai: disamarkan.
	pesan := repository.PolaIdentitasSementaraTCO().ReplaceAllString(g.Galat.Error(), "(baris baru dalam permintaan ini)")
	return "simpan utuh dibatalkan seluruhnya - gagal pada " + g.Bagian + ": " + pesan
}

// Unwrap menjaga pemetaan HTTP galat aslinya (409, 422, ...).
func (g GalatSimpanUtuhTCO) Unwrap() error { return g.Galat }

// ReinsurerUtuhMasuk adalah satu reinsurer beserta security-nya.
type ReinsurerUtuhMasuk struct {
	ReinsurerMasuk
	Security []SecurityMasuk `json:"security"`
}

// KontrakUtuhMasuk adalah seluruh perubahan satu kontrak.
type KontrakUtuhMasuk struct {
	Kontrak   KontrakMasuk         `json:"kontrak"`
	Reinsurer []ReinsurerUtuhMasuk `json:"reinsurer"`
	Business  []BusinessMasuk      `json:"business"`
	Klausul   []KlausulMasuk       `json:"klausul"`
}

func (m KontrakUtuhMasuk) jumlahBaris() int {
	n := 1 + len(m.Business) + len(m.Klausul)
	for _, r := range m.Reinsurer {
		n += 1 + len(r.Security)
	}
	return n
}

// ReinsurerUtuhTampil adalah satu reinsurer tersimpan beserta security-nya.
type ReinsurerUtuhTampil struct {
	Reinsurer ReinsurerTampil  `json:"reinsurer"`
	Security  []SecurityTampil `json:"security"`
}

// HasilSimpanUtuhTampil adalah jawaban simpan utuh.
//
// `Status` "1" = sukses (AC 39): klien membaca SELAIN "1" - kosong, NULL,
// apa pun - sebagai kegagalan (`models.StatusSimpanSuksesTCO`).
type HasilSimpanUtuhTampil struct {
	Status     string                `json:"status"`
	Kontrak    KontrakTampil         `json:"kontrak"`
	Reinsurer  []ReinsurerUtuhTampil `json:"reinsurer"`
	Business   []BusinessTampil      `json:"business"`
	Klausul    []KlausulTampil       `json:"klausul"`
	Peringatan []string              `json:"peringatan"`
}

// SimpanUtuhTCO mengorkestrasi simpan utuh satu kontrak.
type SimpanUtuhTCO struct {
	svc       *Service
	kontrak   *KontrakTreatyTCO
	reinsurer *ReinsurerTCO
	security  *SecurityTCO
	business  *BusinessTCO
	klausul   *KlausulTCO
	tetapkan  func(ctx context.Context, tx *repository.Tx) (map[string]string, error)
	jejak     func(ctx context.Context, tx *repository.Tx, akunID, tabel, barisID, aksi, keterangan string, waktu time.Time) error
	jam       func() time.Time
	transaksi func(ctx context.Context, fn func(tx *repository.Tx) error) error
}

// SimpanUtuhTCO menyusun orkestratornya; penulis bawaannya gagal terang.
func (s *Service) SimpanUtuhTCO() *SimpanUtuhTCO {
	l := &SimpanUtuhTCO{svc: s, kontrak: s.KontrakTreatyTCO(), reinsurer: s.ReinsurerTCO(), security: s.SecurityTCO(),
		business: s.BusinessTCO(), klausul: s.KlausulTCO(), jam: time.Now, transaksi: s.DalamTransaksi,
		tetapkan: func(context.Context, *repository.Tx) (map[string]string, error) {
			return nil, ErrPenetapIdentitasBelumDisuntik
		},
		jejak: func(context.Context, *repository.Tx, string, string, string, string, string, time.Time) error {
			return ErrPenetapIdentitasBelumDisuntik
		}}
	return l
}

func (l *SimpanUtuhTCO) salin() *SimpanUtuhTCO { s := *l; return &s }

// DenganKontrak memasang penulis kontrak (tiket 04).
func (l *SimpanUtuhTCO) DenganKontrak(k *KontrakTreatyTCO) *SimpanUtuhTCO {
	s := l.salin()
	s.kontrak = k
	return s
}

// DenganReinsurer memasang penulis reinsurer (tiket 05).
func (l *SimpanUtuhTCO) DenganReinsurer(r *ReinsurerTCO) *SimpanUtuhTCO {
	s := l.salin()
	s.reinsurer = r
	return s
}

// DenganSecurity memasang penulis security (tiket 06).
func (l *SimpanUtuhTCO) DenganSecurity(x *SecurityTCO) *SimpanUtuhTCO {
	s := l.salin()
	s.security = x
	return s
}

// DenganBusiness memasang penulis business (tiket 07).
func (l *SimpanUtuhTCO) DenganBusiness(b *BusinessTCO) *SimpanUtuhTCO {
	s := l.salin()
	s.business = b
	return s
}

// DenganKlausul memasang penulis klausul (tiket 08/11).
func (l *SimpanUtuhTCO) DenganKlausul(k *KlausulTCO) *SimpanUtuhTCO {
	s := l.salin()
	s.klausul = k
	return s
}

// DenganPenetapIdentitas memasang penetap identitas + perekam jejak utuh.
func (l *SimpanUtuhTCO) DenganPenetapIdentitas(
	tetapkan func(ctx context.Context, tx *repository.Tx) (map[string]string, error),
	jejak func(ctx context.Context, tx *repository.Tx, akunID, tabel, barisID, aksi, keterangan string, waktu time.Time) error,
) *SimpanUtuhTCO {
	s := l.salin()
	s.tetapkan, s.jejak = tetapkan, jejak
	return s
}

// PenetapIdentitasOracle - `TetapkanIdentitasTCO` + `SisipJejakTCO` di atas Oracle.
func PenetapIdentitasOracle(svc *Service) (
	func(ctx context.Context, tx *repository.Tx) (map[string]string, error),
	func(ctx context.Context, tx *repository.Tx, akunID, tabel, barisID, aksi, keterangan string, waktu time.Time) error,
) {
	return svc.db.TetapkanIdentitasTCO, svc.db.SisipJejakTCO
}

// DenganJam mengganti sumber waktu - dipakai uji.
func (l *SimpanUtuhTCO) DenganJam(j func() time.Time) *SimpanUtuhTCO {
	s := l.salin()
	s.jam = j
	return s
}

// DenganTransaksi mengganti pelaksana transaksi - dipakai uji.
func (l *SimpanUtuhTCO) DenganTransaksi(f func(ctx context.Context, fn func(tx *repository.Tx) error) error) *SimpanUtuhTCO {
	s := l.salin()
	s.transaksi = f
	return s
}

// Simpan menulis seluruh perubahan satu kontrak - bersama, atau tidak sama sekali.
func (l *SimpanUtuhTCO) Simpan(ctx context.Context, pelaku Pelaku, tahunID string, m KontrakUtuhMasuk) (
	HasilSimpanUtuhTampil, error) {

	if err := WajibIdentitas(pelaku); err != nil {
		return HasilSimpanUtuhTampil{}, err
	}
	if n := m.jumlahBaris(); n > batasBarisUtuhTCO {
		return HasilSimpanUtuhTampil{}, fmt.Errorf("%w: %d baris, batas %d", ErrSimpanUtuhTidakSah, n, batasBarisUtuhTCO)
	}
	var h HasilSimpanUtuhTampil
	err := l.transaksi(ctx, func(tx *repository.Tx) error {
		c := repository.DenganTransaksiUtuhTCO(ctx, tx)
		// Setiap penulis per baris berjalan di transaksi YANG SAMA: tidak ada
		// commit per tabel.
		dalam := func(_ context.Context, fn func(tx *repository.Tx) error) error { return fn(tx) }
		k, err := l.kontrak.DenganTransaksi(dalam).Simpan(c, pelaku, tahunID, m.Kontrak)
		if err != nil {
			return GalatSimpanUtuhTCO{Bagian: "kontrak", Galat: err}
		}
		h.Kontrak = k
		for i, r := range m.Reinsurer {
			hr, err := l.reinsurer.DenganTransaksi(dalam).Simpan(c, pelaku, tahunID, k.ID, r.ReinsurerMasuk)
			if err != nil {
				return GalatSimpanUtuhTCO{Bagian: fmt.Sprintf("reinsurer ke-%d", i+1), Galat: err}
			}
			ru := ReinsurerUtuhTampil{Reinsurer: hr.Reinsurer, Security: []SecurityTampil{}}
			for j, s := range r.Security {
				hs, err := l.security.DenganTransaksi(dalam).Simpan(c, pelaku, tahunID, k.ID, hr.Reinsurer.ID, s)
				if err != nil {
					return GalatSimpanUtuhTCO{Bagian: fmt.Sprintf("security ke-%d pada reinsurer ke-%d", j+1, i+1), Galat: err}
				}
				ru.Security = append(ru.Security, hs)
			}
			h.Reinsurer = append(h.Reinsurer, ru)
		}
		for i, b := range m.Business {
			hb, err := l.business.DenganTransaksi(dalam).Simpan(c, pelaku, tahunID, k.ID, b)
			if err != nil {
				return GalatSimpanUtuhTCO{Bagian: fmt.Sprintf("business ke-%d", i+1), Galat: err}
			}
			h.Business = append(h.Business, hb)
		}
		for i, kl := range m.Klausul {
			hk, err := l.klausul.DenganTransaksi(dalam).Simpan(c, pelaku, tahunID, kl)
			if err != nil {
				return GalatSimpanUtuhTCO{Bagian: fmt.Sprintf("klausul ke-%d", i+1), Galat: err}
			}
			h.Klausul = append(h.Klausul, hk.Klausul)
			if hk.Peringatan != "" {
				h.Peringatan = append(h.Peringatan, fmt.Sprintf("klausul ke-%d: %s", i+1, hk.Peringatan))
			}
		}
		// Jejak permintaan utuh (AC 41) - sebelum identitas ditetapkan supaya
		// ikut diganti.
		ket := fmt.Sprintf("simpan utuh: 1 kontrak, %d reinsurer, %d security, %d business, %d klausul",
			len(m.Reinsurer), jumlahSecurity(m), len(m.Business), len(m.Klausul))
		if err := l.jejak(c, tx, pelaku.AkunID, repository.TabelKontrakTCO, k.ID, repository.AksiJejakSimpan, ket, l.jam()); err != nil {
			return GalatSimpanUtuhTCO{Bagian: "jejak", Galat: err}
		}
		peta, err := l.tetapkan(c, tx)
		if err != nil {
			return GalatSimpanUtuhTCO{Bagian: "penetapan identitas", Galat: err}
		}
		gantiIdentitasUtuh(&h, peta)
		return nil
	})
	if err != nil {
		return HasilSimpanUtuhTampil{}, err
	}
	h.Status = models.StatusSimpanSuksesTeksTCO
	if h.Reinsurer == nil {
		h.Reinsurer = []ReinsurerUtuhTampil{}
	}
	if h.Business == nil {
		h.Business = []BusinessTampil{}
	}
	if h.Klausul == nil {
		h.Klausul = []KlausulTampil{}
	}
	if h.Peringatan == nil {
		h.Peringatan = []string{}
	}
	return h, nil
}

func jumlahSecurity(m KontrakUtuhMasuk) int {
	n := 0
	for _, r := range m.Reinsurer {
		n += len(r.Security)
	}
	return n
}

// gantiIdentitasUtuh menulis identitas tetap ke jawaban.
func gantiIdentitasUtuh(h *HasilSimpanUtuhTampil, peta map[string]string) {
	ganti := func(s *string) {
		if t, ada := peta[*s]; ada {
			*s = t
		}
	}
	ganti(&h.Kontrak.ID)
	for i := range h.Reinsurer {
		ganti(&h.Reinsurer[i].Reinsurer.ID)
		for j := range h.Reinsurer[i].Security {
			ganti(&h.Reinsurer[i].Security[j].ID)
			ganti(&h.Reinsurer[i].Security[j].ReasID)
		}
	}
	for i := range h.Business {
		ganti(&h.Business[i].ID)
	}
	for i := range h.Klausul {
		ganti(&h.Klausul[i].ID)
	}
}
