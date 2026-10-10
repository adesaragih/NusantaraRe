package services

// Untuk apa berkas ini: PILIHAN POP-UP DAN AUTOCOMPLETE - isi grid harness (ChooseMasterTNonProp, ViewListPolicyCNP,
// CauseofLoss_Harness, CatastrofeList, Hitung_Test, ViewOldAllocation, ViewAttachmentNP, ViewHistoryMasterID_NP) dan
// sumber daftar sel (polis, adjuster, provinsi, rekening, klien, treaty spreading). Baca-saja; pilihan yang dipilih
// dikirim balik sebagai aksi dan DIBACA ULANG di server.

import (
	"context"
	"fmt"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/claimnonprop/backend/models"
)

// Jenis pilihan pop-up.
const (
	PilihanMaster        = "master"        // ChooseMasterTNonProp (BrowseDtlMasterTNP_Act)
	PilihanDaftarPolis   = "daftarPolis"   // ViewListPolicyCNP (Polis.pxResults)
	PilihanSebab         = "sebab"         // CauseofLoss_Harness
	PilihanKatastrofe    = "katastrofe"    // CatastrofeList
	PilihanSelisihAktual = "selisihAktual" // Hitung_Test (GetSelisihActual_Act)
	PilihanAlokasiLama   = "alokasiLama"   // ViewOldAllocation (GetDataOldAllocation)
	PilihanLampiranBayar = "lampiranBayar" // ViewAttachmentNP (GetPayAttachmentNP_Act, bagian baca)
	PilihanRiwayatMaster = "riwayatMaster" // ViewHistoryMasterID_NP (GetHistoryMasterID_NP)
)

// RiwayatMaster - isi harness ViewHistoryMasterID_NP.
type RiwayatMaster struct {
	Judul string             `json:"judul"`
	Baris []models.BarisXOL2 `json:"baris"`
	Total []models.Baris     `json:"total"`
}

// BerkasPolis - tombol View: berkas NB / EDM Treaty In terbaru untuk nomor polis (bawaan OQ-CNP-13).
func (l *Layanan) BerkasPolis(ctx context.Context, p inti.Pelaku, nopolis string) (models.BerkasPolis, error) {
	if err := l.periksaPelaku(p); err != nil {
		return models.BerkasPolis{}, err
	}
	nopolis = strings.TrimSpace(nopolis)
	if nopolis == "" {
		return models.BerkasPolis{}, fmt.Errorf("%w: nomor polis kosong", ErrPermintaanTidakSah)
	}
	b, ada, err := l.a.BerkasPolis(ctx, nopolis)
	if err != nil {
		return models.BerkasPolis{}, err
	}
	if !ada {
		return models.BerkasPolis{}, fmt.Errorf("%w: polis belum punya berkas NB / EDM Treaty In", ErrKasusTidakAda)
	}
	return b, nil
}

// Pilihan membaca satu daftar pilihan kasus.
func (l *Layanan) Pilihan(ctx context.Context, p inti.Pelaku, id, jenis string, n int, cari string,
	sm models.SaringanMaster) (any, error) {
	if err := l.periksaPelaku(p); err != nil {
		return nil, err
	}
	k, h, err := l.muat(ctx, nil, id)
	if err != nil {
		return nil, err
	}
	if err := l.turunkan(ctx, h); err != nil {
		return nil, err
	}
	switch jenis {
	case PilihanMaster:
		return l.a.DaftarMaster(ctx, sm)
	case PilihanDaftarPolis, models.SumberPolis:
		rows, err := l.a.DaftarPolis(ctx, h.Ambil(models.CD+"IDMaster"))
		if err != nil || jenis == PilihanDaftarPolis {
			return rows, err
		}
		out := []models.Pilihan{}
		for _, r := range rows { // autocomplete: nilai .CARI2, tampil .CARI2
			if cari == "" || strings.Contains(strings.ToUpper(r.PolicyNo), strings.ToUpper(cari)) {
				out = append(out, models.Pilihan{Nilai: r.PolicyNo, Label: r.PolicyNo,
					Tambahan: map[string]string{"TreatyGroup": r.TreatyGroup}})
			}
		}
		return out, nil
	case PilihanSebab:
		return l.a.DaftarSebab(ctx, cari)
	case PilihanKatastrofe:
		return l.a.DaftarKatastrofe(ctx, cari)
	case models.SumberAdjuster:
		return l.a.DaftarAdjuster(ctx, cari)
	case models.SumberProvinsi:
		return l.a.DaftarProvinsi(ctx, cari)
	case models.SumberKlien:
		return l.a.DaftarKlien(ctx, cari)
	case models.SumberTreaty:
		nama, err := l.a.NamaTreatySpreading(ctx, h.Ambil(models.OQ+"BusinessCode"), h.Ambil(models.TM+"Commencement"))
		if err != nil {
			return nil, err
		}
		out := []models.Pilihan{}
		for _, v := range models.SaringTreatySpreading(nama) {
			out = append(out, models.Pilihan{Nilai: v, Label: v})
		}
		return out, nil
	case models.SumberRekening:
		j := &jalanAksi{l: l, ctx: ctx, h: h}
		rek, err := j.rekeningAkseptasi(n)
		if err != nil {
			return nil, err
		}
		out := []models.Pilihan{}
		for _, r := range rek {
			out = append(out, models.Pilihan{Nilai: r.AccountNo + "|" + r.NameOfBank, Label: r.NameOfBank,
				Tambahan: map[string]string{"NoAccount": r.AccountNo, "BranchOfBank": r.BranchOfBank, "Currency": r.Currency,
					"SwiftCode": r.SwiftCode}})
		}
		return out, nil
	case PilihanSelisihAktual:
		var os []models.NilaiOS
		for _, ly := range models.KelompokLayerOS(h) {
			v, err := l.g.JumlahOS(ctx, nil, k.ID, ly["TreatyName"], ly["Currency"])
			if err != nil {
				return nil, err
			}
			os = append(os, v)
		}
		return models.HitungSelisihAktual(h, os)
	case PilihanAlokasiLama: // `AdjustmentList(<LAST>).LossAllocation`
		d := h.AmbilDaftar(models.DaftarAdjustment)
		if len(d) == 0 {
			return []models.Baris{}, nil
		}
		out := h.AmbilDaftar(models.JalurAdj(len(d), models.AnakXOLLama))
		if out == nil {
			out = []models.Baris{}
		}
		return out, nil
	case PilihanLampiranBayar:
		return l.a.LampiranBayar(ctx, k.ID)
	case PilihanRiwayatMaster:
		ids, judul := models.IDMasterRiwayat(h.Ambil(models.CD + "IDMaster"))
		rows, err := l.a.RiwayatMaster(ctx, ids)
		if err != nil {
			return nil, err
		}
		tot, err := models.TotalRiwayatMaster(rows)
		if err != nil {
			return nil, err
		}
		if rows == nil {
			rows = []models.BarisXOL2{}
		}
		return RiwayatMaster{Judul: judul, Baris: rows, Total: tot}, nil
	}
	return nil, fmt.Errorf("%w: pilihan %q", ErrPermintaanTidakSah, jenis)
}
