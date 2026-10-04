package services

// Tiket 10 (E5) - efek keluar sesudah `Confirm`: `InsertJsonPolisLife_Act`
// langkah 16 b7162 (prakondisi b7230 `IsPEGAPROD`) → `serviceInsertArasapasLife_act`.
//
// ⛔ Di LUAR transaksi keputusan, dan kegagalannya tidak membatalkan versi
// yang sudah resmi (ADR-U-0008): ia tercatat di outbox `T_LOG_SERVICE_RNM`
// (`MODUL = ENDORSEMENTLIFE`), terlihat di layar, dan dapat diulang.
//
// ⛔ ALARM DIHIDUPKAN (spec penyimpangan 3), pemicunya bergeser seperti
// PremiumList: simpan di sini atomik - keadaan separuh yang dicari langkah 13
// `//` (`GetNopolisByIDPega`) mustahil - jadi alarm menyala pada kegagalan
// efek keluar itu sendiri. Alarm = pesan di layar + log TANPA alamat (E5).
//
// ⛔ Kedua efek STUB yang gagal terang (OQ-EDM-013): alamat di-resolve
// sungguhan dari `M_LINK_SERVICE` (ADR-U-0013), panggilan nyata dan email
// nyata tidak pernah dilakukan. Langkah 15 `//` `SendEmailNotification` memuat
// data rahasia dan TIDAK dikutip ke mana pun.

import (
	"context"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/layanan"
	"nusantarare/inti/backend/outbox"
)

// ModulEndorsementLife mengisi kolom `MODUL` outbox `T_LOG_SERVICE_RNM`.
const ModulEndorsementLife = "ENDORSEMENTLIFE"

// KunciArasapasEndorsement - `serviceInsertArasapasLife_act` langkah 4 b748
// `GetLinkService`, `Kategori_1` b793 `"Production"`, `Kategori_2` b795
// `"convertJsonNusareToProduction"`. Kunci, bukan alamat.
var KunciArasapasEndorsement = layanan.KunciLayanan{
	Kategori1: "Production",
	Kategori2: "convertJsonNusareToProduction",
}

// Nama kedua efek - kolom `JENIS_EFEK` outbox.
const (
	NamaEfekArasapasEDM = "arasapas-endorsement"
	NamaEfekAlarmEDM    = "email-alarm-endorsement"
)

// EfekArasapasEDM me-resolve alamatnya sungguhan, lalu berhenti terang.
type EfekArasapasEDM struct {
	Resolver layanan.ResolverEndpoint
}

// Nama menyebut efek ini di catatan kegagalan.
func (EfekArasapasEDM) Nama() string { return NamaEfekArasapasEDM }

// Jalankan - kunci yang tidak ada di `M_LINK_SERVICE` terlihat sebagai kegagalan konfigurasi.
func (e EfekArasapasEDM) Jalankan(ctx context.Context, _ outbox.MuatanEfek) error {
	if e.Resolver == nil {
		return layanan.ErrResolverBelumDiputuskan
	}
	if _, err := layanan.AlamatLayanan(ctx, e.Resolver, KunciArasapasEndorsement); err != nil {
		return err
	}
	return outbox.ErrArasapasBelumDisetujui
}

// EfekAlarmEmailEDM - email alarm; penerima dan isinya menunggu OQ-EDM-013.
type EfekAlarmEmailEDM struct{}

// Nama menyebut efek ini di catatan kegagalan.
func (EfekAlarmEmailEDM) Nama() string { return NamaEfekAlarmEDM }

// Jalankan selalu gagal, dan menyebut apa yang ditunggu.
func (EfekAlarmEmailEDM) Jalankan(context.Context, outbox.MuatanEfek) error {
	return outbox.ErrEmailBelumDisetujui
}

// Penyalur - yang `Putuskan` panggil sesudah commit (uji: penyalur tiruan).
type Penyalur interface {
	Salurkan(ctx context.Context, m outbox.MuatanEfek) outbox.HasilSalur
}

// PenyalurEDM menjalankan Arasapas, lalu alarm HANYA bila Arasapas gagal -
// dua `outbox.Penyalur` bersarang (satu penyalur berisi dua efek akan
// mengirim alarm pada setiap Confirm yang berhasil).
type PenyalurEDM struct {
	kirim *outbox.Penyalur
	alarm *outbox.Penyalur
}

// NewPenyalurEDM menyusunnya. Gerbang lingkungan milik `outbox.Penyalur` - SATU tempat.
func NewPenyalurEDM(l inti.Lingkungan, a outbox.Antrean, kirim, alarm outbox.EfekKeluar) *PenyalurEDM {
	return &PenyalurEDM{kirim: outbox.NewPenyalur(l, a, kirim), alarm: outbox.NewPenyalur(l, a, alarm)}
}

// Salurkan menjalankan kiriman, lalu alarm bila kiriman gagal.
func (p *PenyalurEDM) Salurkan(ctx context.Context, m outbox.MuatanEfek) outbox.HasilSalur {
	h := p.kirim.Salurkan(ctx, m)
	if h.Dilewati || len(h.Gagal) == 0 {
		return h
	}
	a := p.alarm.Salurkan(ctx, m)
	h.Gagal = append(h.Gagal, a.Gagal...)
	h.GagalDiantre = append(h.GagalDiantre, a.GagalDiantre...)
	return h
}

// penyalurBawaan - bukan produksi: efek keluar dilewati, keputusan tetap tersimpan penuh.
func penyalurBawaan() Penyalur {
	return NewPenyalurEDM(inti.BukanProduksi, outbox.AntreanBelumDiputuskan{},
		EfekArasapasEDM{Resolver: layanan.ResolverBelumDiputuskan{}}, EfekAlarmEmailEDM{})
}

// PenyalurOracle - outbox dan resolver Oracle; kedua efek tetap stub.
func PenyalurOracle(s *Service) Penyalur {
	return NewPenyalurEDM(s.Lingkungan(), outbox.AntreanEfekOracleModul(s, ModulEndorsementLife),
		EfekArasapasEDM{Resolver: layanan.ResolverLinkServiceOracle(s)}, EfekAlarmEmailEDM{})
}

// RingkasEfek - keadaan efek keluar yang dibaca layar (alarm di layar, E5).
type RingkasEfek struct {
	Dilewati bool     `json:"dilewati"`
	Gagal    []string `json:"gagal"`
	// TidakDiantre - kegagalan yang bahkan tidak tercatat di outbox.
	TidakDiantre []string `json:"tidakDiantre"`
}

// sebabTetap - sebab kegagalan sebagai kalimat TETAP: teks galat apa adanya
// dapat membawa alamat layanan, dan alarm tidak pernah menampilkan atau
// mencatat alamat (E5). Sebab lengkapnya tinggal di outbox.
func sebabTetap(sebab string) string {
	switch {
	case strings.Contains(sebab, outbox.ErrArasapasBelumDisetujui.Error()):
		return "the Arasapas call is not approved yet (OQ-EDM-013); queued in the outbox"
	case strings.Contains(sebab, outbox.ErrEmailBelumDisetujui.Error()):
		return "the alarm email is not approved yet (OQ-EDM-013); queued in the outbox"
	case strings.Contains(sebab, layanan.ErrResolverBelumDiputuskan.Error()),
		strings.Contains(sebab, layanan.ErrEndpointTidakDitemukan.Error()):
		return "the service address could not be resolved from M_LINK_SERVICE"
	}
	return "failed; the cause is recorded in the outbox T_LOG_SERVICE_RNM"
}

// ringkas menerjemahkan hasil penyalur menjadi kalimat tetap per efek.
func ringkas(h outbox.HasilSalur) RingkasEfek {
	r := RingkasEfek{Dilewati: h.Dilewati, Gagal: []string{}, TidakDiantre: []string{}}
	for _, g := range h.Gagal {
		r.Gagal = append(r.Gagal, g.Nama+": "+sebabTetap(g.Sebab))
	}
	for range h.GagalDiantre {
		r.TidakDiantre = append(r.TidakDiantre, "a failed effect could not be queued in the outbox")
	}
	return r
}
