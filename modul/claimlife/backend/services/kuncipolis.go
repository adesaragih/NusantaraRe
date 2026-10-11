package services

// Kunci data polis klaim Life = IDPEGA peserta - keputusan work owner 10-10-2026: "ambil ke tabel T_PREMIUM_LIST
// berdasarkan idpega peserta nya saja" (menggantikan kunci PL / POLICY_NO pagi itu).
//
// ⛔ M_LIFE_PREMIUM_DETAIL.IDPEGA = pengenal kasus versi polis (pzInsKey -> NBLF-/EDML-); peserta yang dipilih selalu
// dari versi TERAKHIR yang hidup (AmbilUntukKlaim / Cari, OQ-N14), dan M_LIFE hanya ditulis saat versi diresmikan
// (NB Bind, EDM Confirm - seluruh peserta versi EDM, termasuk Old, ber-IDPEGA kasus EDM). Jadi header yang dibaca
// pasti versi yang benar dan kasus EDM terbuka tidak mungkin terbawa. NO_POLIS / PL / POLICY_NO TIDAK dipakai lagi.
//
// Pembacanya `RingkasDariIdPega` milik PremiumList (`repository.RingkasPolisLife`). ⛔ kontrak.PembacaPolis (inti, milik
// tim inti) TIDAK diubah: kemampuan itu diminta lewat antarmuka sempit PembacaPolisIdPega yang dipenuhi tipe
// pembaca yang sama (penegasan tipe) - pembaca yang tidak memenuhinya gagal terang, bukan jatuh ke NO_POLIS.
//
// ⚠️ Fallback PL / POLICY_NO DIBUANG. Peserta tanpa IDPEGA (baris warisan Pega yang IDPEGA-nya kosong) tidak
// mendapat data polis - dilaporkan pemanggil, tidak dikarang.

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/claimlife/backend/models"
)

// PembacaPolisIdPega - pembaca header polis berkunci IDPEGA peserta (dipenuhi repository.RingkasPolisLife PremiumList).
type PembacaPolisIdPega interface {
	RingkasDariIdPega(ctx context.Context, idPega string) (kontrak.PolisRingkas, error)
}

// ErrPembacaTanpaIdPega - pembaca polis yang disambung tidak dapat membaca berkunci IDPEGA.
var ErrPembacaTanpaIdPega = errors.New("services: pembaca polis tidak mendukung pembacaan berkunci IDPEGA")

// IdPegaPeserta - IDPEGA unik peserta, urut kemunculan; yang kosong dilewati.
func IdPegaPeserta(peserta []models.Peserta) []string {
	var out []string
	sudah := map[string]bool{}
	for _, p := range peserta {
		if id := strings.TrimSpace(p.IdPega); id != "" && !sudah[id] {
			sudah[id] = true
			out = append(out, id)
		}
	}
	return out
}

// RingkasPolisIdPega membaca header setiap IDPEGA dan mengembalikan versi TERBARU (PROD_KE terbesar) beserta IDPEGA-nya.
// Seluruh peserta satu case semestinya ber-IDPEGA sama; bila berbeda, yang terbaru dipakai (bukan ditebak). IDPEGA
// kosong / tidak satu pun ditemukan = ErrPolisNomorTakDitemukan.
func RingkasPolisIdPega(ctx context.Context, pembaca kontrak.PembacaPolis, idPega []string) (kontrak.PolisRingkas, string, error) {
	baca, ok := pembaca.(PembacaPolisIdPega)
	if !ok {
		return kontrak.PolisRingkas{}, "", ErrPembacaTanpaIdPega
	}
	var (
		terbaik   kontrak.PolisRingkas
		kunci     string
		tidakAda  []string
		ditemukan bool
	)
	for _, id := range idPega {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		p, err := baca.RingkasDariIdPega(ctx, id)
		if errors.Is(err, kontrak.ErrPolisNomorTakDitemukan) {
			tidakAda = append(tidakAda, id)
			continue
		}
		if err != nil {
			return kontrak.PolisRingkas{}, "", err
		}
		if !ditemukan || p.ProdKe > terbaik.ProdKe {
			terbaik, kunci, ditemukan = p, id, true
		}
	}
	if !ditemukan {
		return kontrak.PolisRingkas{}, strings.Join(tidakAda, ", "),
			fmt.Errorf("%w: IDPEGA %q", kontrak.ErrPolisNomorTakDitemukan, strings.Join(tidakAda, ", "))
	}
	return terbaik, kunci, nil
}
