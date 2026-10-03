package services

// Create opportunity - tiket 29, butir 76 (keputusan work owner 03-10-2026). Satu
// transaksi: nomor NB dari sequence, baris T_WORK_POLIS (tabel yang ada, K-064), baris
// T_NB_OPPORTUNITY. Urutan pemeriksaan sepola KasusPolis.Buat premiumlistlife:
// identitas (401) -> isian (400) -> basis data (503).

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// Transaksi - pembuka satu transaksi (inti.Dasar.DalamTransaksi): galat fn = batal.
type Transaksi func(ctx context.Context, fn func(tx *db.Tx) error) error

// DenganTransaksi memasang pembuka transaksi - SATU untuk semua penulis (case NB tiket
// 29, General tiket 31).
func (s *Service) DenganTransaksi(t Transaksi) *Service {
	s.transaksi = t
	return s
}

// DenganCaseNB memasang pembuat case NB (tiket 29); transaksinya lewat DenganTransaksi.
func (s *Service) DenganCaseNB(c repository.PenulisCaseNB) *Service {
	s.caseNB = c
	return s
}

// NilaiFacultative - Type Of Inward yang memunculkan Type Of Facultative (satu-satunya
// pilihan yang terlihat di tangkapan layar, `[dugaan]` frontend NILAI_AWAL_OPPORTUNITY).
const NilaiFacultative = "Facultative"

// BentukTanggalKabel - bentuk kabel inti DD-MM-YYYY (NFR-14, inti/frontend/lib/tanggalInput.ts).
const BentukTanggalKabel = "02-01-2006"

// ErrMasukanOpportunity - medan wajib kosong, tanggal tak sah, atau isian melebihi
// lebar kolomnya. 400; rinciannya (nama medan kabel) di pesan.
var ErrMasukanOpportunity = errors.New("services: isian opportunity tidak sah")

// ErrOpportunityTanpaDatabase - layanan tanpa basis data: case NB tidak dapat dibuat. 503.
var ErrOpportunityTanpaDatabase = errors.New("services: basis data tidak dikonfigurasi, case NB tidak dapat dibuat")

// IsianOpportunity - badan POST /api/nbfacin/opportunity apa adanya (kontrak frontend
// `IsianOpportunity`, modul/nbfacin/frontend/api.ts).
type IsianOpportunity struct {
	EstimatedClosingDate, BusinessProspectName, AccountID, InsuredID, GroupBusinessID, GroupBusiness string
	ClassOfBusiness, TypeOfInward, TypeOfFacultative, Phase, Stage, OpportunitySource                string
	BusinessStatus, Description                                                                      string
}

// batasMedan - lebar kolom T_NB_OPPORTUNITY (migrasi 180, A81). karakter = semantik
// CHAR (dihitung rune), selainnya BYTE (dihitung bita).
type batasMedan struct {
	nama     string
	nilai    func(IsianOpportunity) string
	n        int
	karakter bool
	wajib    bool
}

var medanOpportunity = []batasMedan{
	{nama: "businessProspectName", nilai: func(i IsianOpportunity) string { return i.BusinessProspectName }, n: 255, wajib: true},
	{nama: "accountId", nilai: func(i IsianOpportunity) string { return i.AccountID }, n: 255, karakter: true},
	{nama: "insuredId", nilai: func(i IsianOpportunity) string { return i.InsuredID }, n: 255, karakter: true},
	{nama: "groupBusinessId", nilai: func(i IsianOpportunity) string { return i.GroupBusinessID }, n: 32, karakter: true},
	{nama: "groupBusiness", nilai: func(i IsianOpportunity) string { return i.GroupBusiness }, n: 64, karakter: true},
	{nama: "classOfBusiness", nilai: func(i IsianOpportunity) string { return i.ClassOfBusiness }, n: 4000, wajib: true},
	{nama: "typeOfInward", nilai: func(i IsianOpportunity) string { return i.TypeOfInward }, n: 255, wajib: true},
	{nama: "typeOfFacultative", nilai: func(i IsianOpportunity) string { return i.TypeOfFacultative }, n: 255},
	{nama: "phase", nilai: func(i IsianOpportunity) string { return i.Phase }, n: 255, wajib: true},
	{nama: "stage", nilai: func(i IsianOpportunity) string { return i.Stage }, n: 255},
	{nama: "opportunitySource", nilai: func(i IsianOpportunity) string { return i.OpportunitySource }, n: 255},
	{nama: "businessStatus", nilai: func(i IsianOpportunity) string { return i.BusinessStatus }, n: 255, wajib: true},
	{nama: "description", nilai: func(i IsianOpportunity) string { return i.Description }, n: 4000},
}

// periksaOpportunity - medan wajib (bertanda * di layar, urutan tiket 29): Estimated
// Closing Date, Business Prospect Name, Class Of Business, Type Of Inward, Type Of
// Facultative bila Facultative, Phase, Business Status. Kosong = hanya spasi. Nilai
// TIDAK dipangkas: disimpan apa adanya.
func periksaOpportunity(i IsianOpportunity) (models.Opportunity, error) {
	var masalah []string
	tanggal, err := time.Parse(BentukTanggalKabel, i.EstimatedClosingDate)
	switch {
	case strings.TrimSpace(i.EstimatedClosingDate) == "":
		masalah = append(masalah, "estimatedClosingDate wajib diisi")
	case err != nil || tanggal.Format(BentukTanggalKabel) != i.EstimatedClosingDate:
		masalah = append(masalah, "estimatedClosingDate bukan tanggal DD-MM-YYYY yang sah")
	}
	for _, m := range medanOpportunity {
		v := m.nilai(i)
		wajib := m.wajib || (m.nama == "typeOfFacultative" && i.TypeOfInward == NilaiFacultative)
		panjang := len(v)
		if m.karakter {
			panjang = utf8.RuneCountInString(v)
		}
		switch {
		case wajib && strings.TrimSpace(v) == "":
			masalah = append(masalah, m.nama+" wajib diisi")
		case panjang > m.n:
			satuan := "byte"
			if m.karakter {
				satuan = "karakter"
			}
			masalah = append(masalah, fmt.Sprintf("%s paling banyak %d %s", m.nama, m.n, satuan))
		}
	}
	if len(masalah) > 0 {
		return models.Opportunity{}, fmt.Errorf("%w: %s", ErrMasukanOpportunity, strings.Join(masalah, "; "))
	}
	return models.Opportunity{EstimatedClosingDate: tanggal, BusinessProspectName: i.BusinessProspectName,
		AccountID: i.AccountID, InsuredID: i.InsuredID, GroupBusinessID: i.GroupBusinessID, GroupBusiness: i.GroupBusiness,
		ClassOfBusiness: i.ClassOfBusiness, TypeOfInward: i.TypeOfInward, TypeOfFacultative: i.TypeOfFacultative,
		Phase: i.Phase, Stage: i.Stage, OpportunitySource: i.OpportunitySource, BusinessStatus: i.BusinessStatus,
		Description: i.Description}, nil
}

// BuatOpportunity - membuat case NB beserta isian opportunity-nya; mengembalikan nomor
// case (NB-<n>). Galat apa pun di dalam transaksi membatalkan ketiganya.
func (s *Service) BuatOpportunity(ctx context.Context, pelaku inti.Pelaku, isian IsianOpportunity) (string, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return "", err
	}
	o, err := periksaOpportunity(isian)
	if err != nil {
		return "", err
	}
	if s.caseNB == nil || s.transaksi == nil {
		return "", ErrOpportunityTanpaDatabase
	}
	var id string
	err = s.transaksi(ctx, func(tx *db.Tx) error {
		var err error
		if id, err = s.caseNB.PengenalBerikut(ctx, tx); err != nil {
			return err
		}
		if err := s.caseNB.SisipCase(ctx, tx, id, pelaku.AkunID); err != nil {
			return err
		}
		return s.caseNB.SisipOpportunity(ctx, tx, id, o)
	})
	if err != nil {
		return "", err
	}
	return id, nil
}
