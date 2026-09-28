package services

// Riwayat tangga Komite - tiket 09 Komite Claim Life.
//
// Untuk apa berkas ini: siapa pun yang teridentifikasi dapat MEMBACA riwayat
// sebuah kasus - melihat bukan memutuskan (AC tiket 09). Tangganya dari
// `T_KOMITE_KOMITELIST` urut `KOMITE_URUT` (bukan page runtime maupun JSON,
// AC 33 spec); eskalasinya dari jejak (tiket 03).
//
// `[terverifikasi]` Judul kolom VERBATIM grid `.KomiteList`
// `Section/ShowTransfer.xml` (b29727): `Committee` (`.IDKomite` = jabatan),
// `Status` (`.KomiteAproval`), `Date Approve`, `Comment` - di
// `labels.komite.ts`.

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"nusantarare/internal/models"
	"nusantarare/internal/repository"
)

// Bentuk teks jejak tangga - SATU tempat, dipakai penulis (komite_keputusan.go)
// dan pembaca (berkas ini).
const (
	awalanJejakTingkat  = "Komite tingkat "
	awalanJejakEskalasi = "Eskalasi ke tingkat "
	KataTingkatDilewati = "Dilewati (eskalasi)"
	KataTingkatMenunggu = "Menunggu"
)

// BarisRiwayatKomite adalah satu tingkat tangga.
type BarisRiwayatKomite struct {
	Urut int `json:"urut"`
	// Committee = `.IDKomite` = jabatan (roster.go: `IDKomite = .JABATAN`).
	Committee string `json:"committee"`
	// Anggota = `KOMITE_OPERATORID` (pengenal akun, bukan nama orang).
	Anggota      string `json:"anggota"`
	Status       string `json:"status"`
	TanggalPutus string `json:"dateApprove"`
	Comment      string `json:"comment"`
}

// EskalasiRiwayat adalah satu eskalasi yang terbaca dari jejak.
type EskalasiRiwayat struct {
	DariTingkat int    `json:"dariTingkat"`
	KeTingkat   int    `json:"keTingkat"`
	Oleh        string `json:"oleh"`
	Waktu       string `json:"waktu"`
}

// RiwayatKomite adalah riwayat satu kasus.
type RiwayatKomite struct {
	KasusID      string               `json:"kasusId"`
	AdjustmentID string               `json:"adjustmentId"`
	Tangga       []BarisRiwayatKomite `json:"tangga"`
	Eskalasi     []EskalasiRiwayat    `json:"eskalasi"`
}

// KataStatusTangga menerjemahkan `KOMITE_APPROVAL` - murni.
//
// ⛔ Dilewati ≠ menunggu (AC tiket 09): NULL hanya lahir dari eskalasi
// (tiket 03), `0` adalah tingkat yang belum memutuskan.
func KataStatusTangga(kode string) string {
	switch k := strings.TrimSpace(kode); k {
	case "":
		return KataTingkatDilewati
	case repository.ApprovalKomiteMenunggu:
		return KataTingkatMenunggu
	default:
		if kata := models.KataKeputusanKomite(k); kata != "" {
			return kata
		}
		return "Kode " + k
	}
}

var polaEskalasi = regexp.MustCompile(`^` + awalanJejakEskalasi + `(\d+) \(`)

// susunRiwayat - murni.
func susunRiwayat(k repository.KasusKomite, jejak []repository.JejakKomite) RiwayatKomite {
	r := RiwayatKomite{KasusID: k.Baris.KasusID, AdjustmentID: k.AdjID,
		Tangga: []BarisRiwayatKomite{}, Eskalasi: []EskalasiRiwayat{}}
	for _, a := range k.Tangga {
		r.Tangga = append(r.Tangga, BarisRiwayatKomite{
			Urut: a.Urut, Committee: a.Jabatan, Anggota: a.OperatorID,
			Status: KataStatusTangga(a.Approval), Comment: a.Komentar,
			TanggalPutus: waktuTampil(a.TglAprove.Time, a.TglAprove.Valid),
		})
	}
	for _, j := range jejak {
		m := polaEskalasi.FindStringSubmatch(j.Ke)
		if m == nil {
			continue
		}
		ke, _ := strconv.Atoi(m[1])
		dari, _ := strconv.Atoi(strings.TrimPrefix(j.Dari, awalanJejakTingkat))
		r.Eskalasi = append(r.Eskalasi, EskalasiRiwayat{DariTingkat: dari, KeTingkat: ke,
			Oleh: j.AkunID, Waktu: j.Waktu.Format("2006-01-02 15:04:05")})
	}
	return r
}

// Riwayat membaca riwayat satu kasus - untuk siapa pun yang teridentifikasi.
func (i *InboxKomite) Riwayat(ctx context.Context, pelaku Pelaku, kasusID string) (RiwayatKomite, error) {
	if err := WajibIdentitas(pelaku); err != nil {
		return RiwayatKomite{}, err
	}
	if i == nil || i.svc == nil || !i.svc.PunyaDatabase() {
		return RiwayatKomite{}, repository.ErrTanpaOracle
	}
	if strings.TrimSpace(kasusID) == "" {
		return RiwayatKomite{}, fmt.Errorf("%w: id kasus komite kosong", ErrPermintaanTidakSah)
	}
	baca := repository.NewInboxKomite(i.svc.db)
	k, err := baca.Kasus(ctx, kasusID)
	if err != nil {
		return RiwayatKomite{}, err
	}
	jejak, err := baca.JejakEskalasi(ctx, k.AdjID, awalanJejakEskalasi+"%("+kasusID+")")
	if err != nil {
		return RiwayatKomite{}, err
	}
	return susunRiwayat(k, jejak), nil
}
