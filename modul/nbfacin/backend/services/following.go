package services

// Following Policy Number (tiket 70; kontrak sesi 9d, K70-1..3 / A70-1..4 di docs/issues/70-following-policy.md). Bukti
// `[terverifikasi]` `D:\migrasi\RNM\NB FacIn\`:
//   - Activity\GetListViewFollowing_act.xml langkah 3 (filter kosong -> keluar) dan 4.1-4.9 (klausa ditimpa langkah sesudahnya,
//     BERURUTAN): 4.1 opsi 2 NOPOLIS; 4.2 opsi 1 + "NB-" / "RNW-" IDPEGA; 4.3 opsi 3 ARASAPAS (ditunda, K70-3); 4.4 filter memuat
//     "RNM-E" (opsi apa pun) NOENDORS; 4.5 opsi 1 + "EDM-"; 4.6-4.9 opsi 4-7. Akibat urutan: "RNM-E" menimpa opsi 1 / 2, tetapi
//     opsi 4-7 (langkah 4.6-4.9, sesudah 4.4) menimpanya balik. Klausa kosong -> tanpa kueri (4.10);
//   - Activity\CopyNopolisFollowing_Act.xml (Following / FollowingNB / IDFollowingNB) dan CopyLocationFollowing_Act.xml
//     (langkah 2 LocationList / CargoList dihapus; 9 ketiga medan; 11 PercentShare + Policy.Payment.RICommision; 12-15 bisnis
//     sama -> lokasi utuh, berbeda -> tanpa FEAList / BuildingConstruction / ListCauseOfLoss / OccupationList / PropertyItemList /
//     SurroundingRisk / TotalTSIList).
// Sistem baru (K70-1 / A70-2): sumber = TABEL case sumber; nbfacin hanya memodelkan objek FIRE (kasusTangga "Fire"), jadi bisnis
// selalu sama -> salinan lokasi utuh (lewat gantiObjek seperti Save tab Object: premi dihitung ulang dengan periode case
// penerima). FollowingNB / IDFollowingNB tanpa kolom - hanya FOLLOWING ditulis.

import (
	"context"
	"errors"
	"log"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

const ukuranHalamanFollowing = 20 // GridViewFollowingNB 20 per halaman

var (
	// ErrFollowingTanpaDatabase - penyimpan Following belum dipasang (503).
	ErrFollowingTanpaDatabase = errors.New("services: Following Policy tidak tersedia tanpa basis data")
	// ErrMasukanFollowing - View Option / halaman / badan tidak sah (400).
	ErrMasukanFollowing = errors.New("isian Following Policy tidak sah")
	// ErrFollowingSesudahBinding - Section\Periode.xml L19826 `.ConfirmBinding = 1` (409).
	ErrFollowingSesudahBinding = errors.New("Following Policy tidak dapat diubah setelah Binding.")
	// ErrPolisFollowingTidakAda - polis (case + nomor) tidak ada di FACINPRODUCTION (404).
	ErrPolisFollowingTidakAda = errors.New("Polis Following tidak ditemukan.")
	// ErrCaseSumberTidakAda - case sumber tidak ada di tabel sistem baru (Copy All, A70-2; 404).
	ErrCaseSumberTidakAda = errors.New("Case sumber tidak ditemukan di sistem baru; Copy All hanya untuk case yang dibuat di sistem baru.")
)

// DenganFollowing - memasang penyimpan Following Policy (tiket 70).
func (s *Service) DenganFollowing(p repository.PenyimpanFollowing) *Service {
	s.following = p
	return s
}

// HalamanFollowing - satu halaman hasil cari.
type HalamanFollowing struct {
	Baris                  []models.BarisFollowing
	Total, Halaman, Ukuran int
}

// jenisCariFollowing - GetListViewFollowing_act langkah 4.1-4.9 berurutan -> (jenis klausa, nilai bind); jenis "" = tanpa kueri.
// Opsi 3 (No Slip) ditunda (K70-3) dan opsi di luar 1-7 -> 400.
func jenisCariFollowing(opsi, filter string) (string, string, error) {
	switch opsi {
	case "1", "2", "4", "5", "6", "7":
	default:
		return "", "", errors.Join(ErrMasukanFollowing, errors.New("View Option wajib 1, 2, 4, 5, 6 atau 7"))
	}
	if filter == "" {
		return "", "", nil // langkah 3
	}
	pola := "%" + filter + "%"
	var jenis, nilai string
	if opsi == "2" {
		jenis, nilai = repository.CariFollowingNoPolis, filter
	}
	if opsi == "1" && (strings.Contains(filter, "NB-") || strings.Contains(filter, "RNW-")) {
		jenis, nilai = repository.CariFollowingIDPega, filter // IDPEGA dua bentuk dicocokkan repository (bindFollowing)
	}
	if strings.Contains(filter, "RNM-E") {
		jenis, nilai = repository.CariFollowingNoEndors, filter
	}
	if opsi == "1" && strings.Contains(filter, "EDM-") {
		jenis, nilai = repository.CariFollowingEDM, filter
	}
	switch opsi {
	case "4":
		jenis, nilai = repository.CariFollowingReff, pola
	case "5":
		jenis, nilai = repository.CariFollowingInsured, pola
	case "6":
		jenis, nilai = repository.CariFollowingSOB, pola
	case "7":
		jenis, nilai = repository.CariFollowingQQ, pola
	}
	return jenis, nilai, nil
}

// CariFollowing - GET …/following (popup Search; baca saja, tanpa penjaga tahap).
func (s *Service) CariFollowing(ctx context.Context, pelaku inti.Pelaku, id, opsi, filter string, halaman int) (HalamanFollowing, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return HalamanFollowing{}, err
	}
	if !idKasusSah(id) {
		return HalamanFollowing{}, ErrKasusTidakAda
	}
	if halaman < 1 {
		return HalamanFollowing{}, errors.Join(ErrMasukanFollowing, errors.New("halaman wajib >= 1"))
	}
	jenis, nilai, err := jenisCariFollowing(opsi, strings.TrimSpace(filter))
	if err != nil {
		return HalamanFollowing{}, err
	}
	h := HalamanFollowing{Baris: []models.BarisFollowing{}, Halaman: halaman, Ukuran: ukuranHalamanFollowing}
	if jenis == "" {
		return h, nil
	}
	if s.following == nil {
		return HalamanFollowing{}, ErrFollowingTanpaDatabase
	}
	h.Baris, h.Total, err = s.following.CariFollowing(ctx, jenis, nilai, (halaman-1)*ukuranHalamanFollowing, ukuranHalamanFollowing)
	return h, err
}

// wajibFollowing - A70-4: ConfirmBinding != "1" (409 ErrFollowingSesudahBinding) dan case boleh diubah akun ini (aturan tab).
// !IsLife selalu benar (nbfacin = FIRE).
func (s *Service) wajibFollowing(ctx context.Context, pelaku inti.Pelaku, id string) error {
	t, ada, err := s.tahapKasus(ctx, id)
	if err != nil || !ada {
		return err
	}
	if t.ConfirmBinding == confirmBindingYa {
		return ErrFollowingSesudahBinding
	}
	_, err = alasanUbah(t, pelaku, true)
	return err
}

// AksesFollowing - `bolehFollowing` GET kasus (LUNAK: galat -> false).
func (s *Service) AksesFollowing(ctx context.Context, pelaku inti.Pelaku, id string) bool {
	return s.wajibFollowing(ctx, pelaku, id) == nil
}

// siapFollowing - identitas, penjaga, badan, polis ada di FACINPRODUCTION (Pega tidak memeriksa - nilai grid; di sini dicek
// supaya nomor sembarang tidak tertulis).
func (s *Service) siapFollowing(ctx context.Context, pelaku inti.Pelaku, id, sumber, noPolis string) error {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return err
	}
	if !idKasusSah(id) {
		return ErrKasusTidakAda
	}
	if s.following == nil || s.transaksi == nil {
		return ErrFollowingTanpaDatabase
	}
	if err := s.wajibFollowing(ctx, pelaku, id); err != nil {
		return err
	}
	if !idKasusSah(sumber) || strings.TrimSpace(noPolis) == "" {
		return errors.Join(ErrMasukanFollowing, errors.New("caseIdSumber dan policyNo wajib diisi"))
	}
	ada, err := s.following.AdaPolisFollowing(ctx, sumber, noPolis)
	if err != nil {
		return err
	}
	if !ada {
		return ErrPolisFollowingTidakAda
	}
	return nil
}

// SalinNopolisFollowing - POST …/following/nopolis (CopyNopolisFollowing_Act): Old Policy Number = nomor polis.
func (s *Service) SalinNopolisFollowing(ctx context.Context, pelaku inti.Pelaku, id, sumber, noPolis string) error {
	if err := s.siapFollowing(ctx, pelaku, id, sumber, noPolis); err != nil {
		return err
	}
	err := s.transaksi(ctx, func(tx *db.Tx) error { return s.following.TulisFollowing(ctx, tx, id, noPolis, nil) })
	if errors.Is(err, repository.ErrKasusTidakAda) {
		return ErrKasusTidakAda
	}
	if err == nil {
		log.Printf("nbfacin: Copy Nopolis Following kasus %s oleh %s: %s (case %s)", id, pelaku.AkunID, noPolis, sumber)
	}
	return err
}

// SalinSemuaFollowing - POST …/following/salin-semua (CopyLocationFollowing_Act, A70-2): FOLLOWING + PercentShare + RI
// Commission case sumber, seluruh lokasi case diganti lokasi case sumber (objek, item, coverage, deductible, okupasi, FEA,
// konstruksi, sebab kerugian, catatan kerugian, surrounding risk) beserta spreading-nya - satu transaksi.
func (s *Service) SalinSemuaFollowing(ctx context.Context, pelaku inti.Pelaku, id, sumber, noPolis string) error {
	if err := s.siapFollowing(ctx, pelaku, id, sumber, noPolis); err != nil {
		return err
	}
	if s.objek == nil || s.deduction == nil {
		return ErrFollowingTanpaDatabase
	}
	objek, err := s.BacaObjek(ctx, sumber)
	if errors.Is(err, ErrKasusTidakAda) {
		return ErrCaseSumberTidakAda
	}
	if err != nil {
		return err
	}
	sp, err := s.bacaSpreading(ctx, nil, sumber)
	if errors.Is(err, ErrKasusTidakAda) {
		return ErrCaseSumberTidakAda
	}
	if err != nil {
		return err
	}
	ded, err := s.deduction.BacaDeduction(ctx, sumber)
	if errors.Is(err, repository.ErrKasusTidakAda) {
		return ErrCaseSumberTidakAda
	}
	if err != nil {
		return err
	}
	if err := periksaObjek(objek); err != nil {
		return err
	}
	nilai := &models.NilaiFollowing{PercentShare: sp.PercentShare, RICommision: ded.RICommision}
	err = s.gantiObjek(ctx, id, objek, opsiGantiObjek{spreading: &sp, sebelum: func(tx *db.Tx) error {
		return s.following.TulisFollowing(ctx, tx, id, noPolis, nilai)
	}})
	if err == nil {
		log.Printf("nbfacin: Copy All Following kasus %s oleh %s: %s (case %s, %d lokasi)", id, pelaku.AkunID, noPolis, sumber, len(objek))
	}
	return err
}
