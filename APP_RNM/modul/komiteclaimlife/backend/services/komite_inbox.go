package services

// Inbox Komite dan pembacaan satu kasus - tiket 01 Komite Claim Life.
//
// Untuk apa berkas ini: menjawab "kasus komite mana yang menunggu SAYA", dan
// "seperti apa kasus ini". Nol tulisan.
//
// ⛔ PERAN KOMITE = ROSTER (km2), bukan `X-Peran` Claim Life. Anggota komite
// adalah akun yang tertulis di tangga kasus itu - `T_KOMITE_KOMITELIST`
// yang lahir dari roster `EMAILKOMITE` saat penyerahan (A2). Tidak ada peran
// "anggota komite" umum yang membuka seluruh antrean.
//
// Dibaca sesudah: repository/komite_inbox.go, komite.go (penyerahannya).

import (
	"context"
	"fmt"
	"strings"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/komiteclaimlife/backend/repository"
)

// ErrKasusKomiteTakDitemukan - tidak ada kasus komite dengan id itu.
var ErrKasusKomiteTakDitemukan = repository.ErrKasusKomiteTakDitemukan

// BarisInboxKomiteTampil adalah satu baris Inbox Komite untuk layar.
//
// ⚠️ Empat kolom pertama DIPINJAM dari kolom Pega-standar `InboxPremiumList`
// (`Case ID`, `Create Date/Time`, …, `Work Status`) - worklist `KomiteRouter`
// tidak punya section inbox sendiri di korpus. Ditandai di label layar.
type BarisInboxKomiteTampil struct {
	KasusID    string `json:"kasusId"`
	TglUpdate  string `json:"tglUpdate"`
	StatusWork string `json:"statusWork"`
	KlaimID    string `json:"klaimId"`
	NomorKlaim string `json:"nomorKlaim"`
	// Tingkat berjalan dari seluruh tingkat - `KomiteCount`/`KomiteLoop`.
	TingkatBerjalan int `json:"tingkatBerjalan"`
	KomiteLoop      int `json:"komiteLoop"`
	// NilaiKlaim TEKS (ADR-U-0003).
	NilaiKlaim string `json:"nilaiKlaim"`
	MataUang   string `json:"mataUang"`
	// StatusBaris KATA - Outstanding / Aksep / Ditolak (AC 29 spec).
	StatusBaris string `json:"statusBaris"`
}

// HalamanInboxKomite adalah satu halaman inbox.
type HalamanInboxKomite struct {
	Baris   []BarisInboxKomiteTampil `json:"baris"`
	Total   int                      `json:"total"`
	Halaman int                      `json:"halaman"`
	Ukuran  int                      `json:"ukuran"`
}

// AnggotaKasusTampil adalah satu anak tangga untuk layar.
type AnggotaKasusTampil struct {
	Urut     int    `json:"urut"`
	Jabatan  string `json:"jabatan"`
	Approval string `json:"approval"`
	// KataApproval - kata untuk kode approval; kode tak dikenal DISEBUT.
	KataApproval string `json:"kataApproval"`
	Komentar     string `json:"komentar"`
	TglApprove   string `json:"tglApprove"`
	// Saya - baris ini milik pelaku yang membaca.
	Saya bool `json:"saya"`
}

// KasusKomiteTampil adalah satu kasus beserta tangganya.
type KasusKomiteTampil struct {
	Kasus  BarisInboxKomiteTampil `json:"kasus"`
	AdjID  string                 `json:"adjustmentId"`
	Tangga []AnggotaKasusTampil   `json:"tangga"`
	// GiliranSaya - pelaku adalah anggota BERJALAN kasus ini.
	GiliranSaya bool `json:"giliranSaya"`
	// Efek - efek keluar kasus ini dan keadaannya (tiket 08): "perlu
	// intervensi" menempel pada kasusnya, bukan di halaman terpisah.
	Efek RingkasEfekKasus `json:"efek"`
}

// InboxKomite melayani Inbox Komite.
type InboxKomite struct{ svc *Service }

// InboxKomite menyusun layanannya.
func (s *Service) InboxKomite() *InboxKomite { return &InboxKomite{svc: s} }

// ukuranInboxKomiteMaks menjepit ukuran halaman - sama dengan inbox lain.
const ukuranInboxKomiteMaks = 100

func waktuTampil(t time.Time, ada bool) string {
	if !ada {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

func keBarisTampil(b repository.BarisInboxKomite) BarisInboxKomiteTampil {
	return BarisInboxKomiteTampil{
		KasusID:         b.KasusID,
		TglUpdate:       waktuTampil(b.TglUpdate.Time, b.TglUpdate.Valid),
		StatusWork:      b.StatusWork,
		KlaimID:         b.KlaimID,
		NomorKlaim:      b.NomorKlaim,
		TingkatBerjalan: b.TingkatBerjalan,
		KomiteLoop:      b.KomiteLoop,
		NilaiKlaim:      b.NilaiKlaim,
		MataUang:        b.MataUang,
		StatusBaris:     kontrak.StatusBarisDariKode(b.StsReject).String(),
	}
}

// Ambil membaca satu halaman inbox milik pelaku.
func (i *InboxKomite) Ambil(ctx context.Context, pelaku inti.Pelaku, halaman, ukuran int) (
	HalamanInboxKomite, error) {

	if err := inti.WajibIdentitas(pelaku); err != nil {
		return HalamanInboxKomite{}, err
	}
	if i == nil || i.svc == nil || !i.svc.PunyaDatabase() {
		return HalamanInboxKomite{}, db.ErrTanpaOracle
	}
	if halaman < 1 {
		halaman = 1
	}
	if ukuran < 1 || ukuran > ukuranInboxKomiteMaks {
		ukuran = ukuranInboxKomiteMaks
	}
	baris, total, err := repository.NewInboxKomite(i.svc.DB()).Ambil(ctx, pelaku.AkunID,
		kontrak.StatusWorkSelesai, (halaman-1)*ukuran, ukuran)
	if err != nil {
		return HalamanInboxKomite{}, err
	}
	hasil := HalamanInboxKomite{Baris: []BarisInboxKomiteTampil{}, Total: total,
		Halaman: halaman, Ukuran: ukuran}
	for _, b := range baris {
		hasil.Baris = append(hasil.Baris, keBarisTampil(b))
	}
	return hasil, nil
}

// KataApprovalKomite menerjemahkan `KOMITE_APPROVAL` menjadi kata.
//
// `[terverifikasi]` hanya `0` = belum memutuskan (`KomiteRouter` b382). Kode
// lain datang dari dropdown `.KomiteAproval` layar `ShowTransfer` (b30623) -
// artinya dibaca di tiket 02; sampai itu kodenya DISEBUT, tidak diterjemahkan
// dengan tebakan.
//
// ⚠️ Sejak tiket 02/03/09 kodenya sudah terbaca: `1`/`2` = Setuju/Tolak,
// kosong = dilewati eskalasi. Satu penerjemah - `KataStatusTangga`.
func KataApprovalKomite(kode string) string { return KataStatusTangga(kode) }

// Kasus membaca satu kasus. Hanya anggota tangga kasus itu yang boleh.
//
// ⛔ ADR-0014: melihat kasus di luar tangganya sendiri adalah melihat
// pekerjaan orang lain - 403, bukan 404, sebab kasusnya memang ada.
func (i *InboxKomite) Kasus(ctx context.Context, pelaku inti.Pelaku, kasusID string) (
	KasusKomiteTampil, error) {

	if err := inti.WajibIdentitas(pelaku); err != nil {
		return KasusKomiteTampil{}, err
	}
	if i == nil || i.svc == nil || !i.svc.PunyaDatabase() {
		return KasusKomiteTampil{}, db.ErrTanpaOracle
	}
	if strings.TrimSpace(kasusID) == "" {
		return KasusKomiteTampil{}, fmt.Errorf("%w: id kasus komite kosong", galat.ErrPermintaanTidakSah)
	}
	baca := repository.NewInboxKomite(i.svc.DB())
	k, err := baca.Kasus(ctx, kasusID)
	if err != nil {
		return KasusKomiteTampil{}, err
	}
	out, err := susunKasusTampil(k, pelaku.AkunID)
	if err != nil {
		return KasusKomiteTampil{}, err
	}
	efek, err := baca.EfekKasus(ctx, kasusID)
	if err != nil {
		return KasusKomiteTampil{}, err
	}
	out.Efek = ringkasEfekKasus(efek)
	return out, nil
}

// susunKasusTampil murni - dapat diuji tanpa Oracle.
func susunKasusTampil(k repository.KasusKomite, akunID string) (KasusKomiteTampil, error) {
	out := KasusKomiteTampil{Kasus: keBarisTampil(k.Baris), AdjID: k.AdjID,
		Tangga: []AnggotaKasusTampil{}, Efek: RingkasEfekKasus{Efek: []EfekTampil{}}}
	anggota := false
	for _, a := range k.Tangga {
		saya := strings.TrimSpace(a.OperatorID) != "" && a.OperatorID == akunID
		anggota = anggota || saya
		out.Tangga = append(out.Tangga, AnggotaKasusTampil{
			Urut: a.Urut, Jabatan: a.Jabatan, Approval: a.Approval,
			KataApproval: KataApprovalKomite(a.Approval), Komentar: a.Komentar,
			TglApprove: waktuTampil(a.TglAprove.Time, a.TglAprove.Valid), Saya: saya,
		})
		if saya && a.Urut == k.Baris.TingkatBerjalan &&
			a.Approval == repository.ApprovalKomiteMenunggu &&
			!kontrak.KasusTertutup(k.Baris.StatusWork) {
			out.GiliranSaya = true
		}
	}
	if !anggota {
		return KasusKomiteTampil{}, fmt.Errorf("%w: pelaku bukan anggota tangga kasus %q",
			inti.ErrTanpaWewenang, k.Baris.KasusID)
	}
	return out, nil
}
