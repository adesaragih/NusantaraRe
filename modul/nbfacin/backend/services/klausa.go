package services

// Tab Clauses kasus FIRE (tiket 47): ClauseList kasus (baca / ganti utuh tiap Save), pencarian klausa (popup Choose
// Clause; SearchClauseFireSQL_PreAct + RetrieveClauseSQL, dikirim work owner 05-10-2026), dan argumen master klausa.
//
// Pemeriksaan isian (K47-3, keputusan agent): clauseCode wajib dan tidak ganda (SearchClauseFireSQL_PostAct hanya
// menambah kode yang belum ada); clauseLanguage kosong atau "0" / "1" / "2" (kode PostAct); argumentCount kosong atau
// angka; setiap teks muat kolomnya (migrasi 197) - yang melebihi DITOLAK, tidak dipotong.

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

var (
	// ErrMasukanKlausa - isian ClauseList / parameter klausa tidak sah. 400.
	ErrMasukanKlausa = errors.New("services: isian klausa tidak sah")
	// ErrKlausaTanpaDatabase - T_CLAUSELIST / M_ARGCLAUSEFIRE tidak terbaca. 503.
	ErrKlausaTanpaDatabase = errors.New("services: basis data tidak dikonfigurasi, klausa tidak terbaca")
)

const (
	// lebarTeksKlausa - VARCHAR2(4000) T_CLAUSELIST / T_CLAUSEARGUMENTLIST (migrasi 197).
	lebarTeksKlausa = 4000
	// lebarArgumentCount - T_CLAUSELIST.ARGUMENT_COUNT VARCHAR2(50).
	lebarArgumentCount = 50
	// banyakKlausaMaks - batas baris ClauseList per Save (penjaga badan; Pega tanpa batas tertulis).
	banyakKlausaMaks = 500
	// UkuranHalamanKlausa - hasil Choose Clause per halaman (C-3).
	UkuranHalamanKlausa = 10
	// lebarKataKlausa - batas kata kunci / bahasa (pola A73).
	lebarKataKlausa = 255
)

// labelBahasaKlausa - PromptList `DDL\ClauseLanguageID.xml` (ASM-FW-GISFW-DATA-CLAUSE!CLAUSELANGUAGEID)
// `[terverifikasi]`: kode (pyStandardValue) -> label (pyLocalizedValue). Label = `.Language` yang PostAct petakan kembali
// ke ClauseLanguage ("Indonesia" "0", "Inggris" "1", lain "2").
var labelBahasaKlausa = map[string]string{"0": "Indonesia", "1": "Inggris", "2": "Dual Bahasa"}

// HalamanKlausa - satu halaman Choose Clause.
type HalamanKlausa struct {
	Baris                []models.HasilKlausa
	Total, Nomor, Ukuran int
}

// kodeBahasaKlausa - ClauseLanguage SearchClauseFireSQL_PostAct: Indonesia "0", Inggris "1", lain "2".
var kodeBahasaKlausa = []string{"0", "1", "2"}

// DenganKlausa memasang penyimpan ClauseList kasus dan argumen klausa (tiket 47).
func (s *Service) DenganKlausa(k repository.PenyimpanKlausa) *Service {
	s.klausa = k
	return s
}

// periksaKlausa - lihat kepala berkas; masalah dikumpulkan berurutan.
func periksaKlausa(baris []models.KlausaKasus) error {
	if len(baris) > banyakKlausaMaks {
		return fmt.Errorf("%w: paling banyak %d klausa", ErrMasukanKlausa, banyakKlausaMaks)
	}
	var masalah []string
	kode := map[string]bool{}
	lebar := func(nama, v string, maks int) {
		if len(v) > maks {
			masalah = append(masalah, fmt.Sprintf("%s paling banyak %d bita", nama, maks))
		}
	}
	for i, b := range baris {
		p := fmt.Sprintf("baris[%d].", i)
		switch {
		case strings.TrimSpace(b.Code) == "":
			masalah = append(masalah, p+"clauseCode wajib diisi")
		case kode[b.Code]:
			masalah = append(masalah, p+"clauseCode "+b.Code+" ganda")
		}
		kode[b.Code] = true
		if b.Language != "" && !slices.Contains(kodeBahasaKlausa, b.Language) {
			masalah = append(masalah, p+"clauseLanguage harus 0, 1, atau 2")
		}
		if strings.Trim(b.ArgumentCount, "0123456789") != "" {
			masalah = append(masalah, p+"argumentCount harus angka")
		}
		lebar(p+"argumentCount", b.ArgumentCount, lebarArgumentCount)
		for _, t := range []struct{ nama, nilai string }{{"clauseCode", b.Code}, {"clauseTitle", b.Title},
			{"clauseDescription", b.Description}, {"clauseLanguageId", b.LanguageID}, {"clauseContent", b.Content},
			{"clauseContentTemp", b.ContentTemp}} {
			lebar(p+t.nama, t.nilai, lebarTeksKlausa)
		}
		for j, a := range b.Arguments {
			q := fmt.Sprintf("%sargumentList[%d].", p, j)
			lebar(q+"argumentNumber", a.Number, lebarTeksKlausa)
			lebar(q+"argumentDescription", a.Description, lebarTeksKlausa)
			lebar(q+"argumentValue", a.Value, lebarTeksKlausa)
		}
	}
	if len(masalah) > 0 {
		return fmt.Errorf("%w: %s", ErrMasukanKlausa, strings.Join(masalah, "; "))
	}
	return nil
}

// BacaKlausa - GET /api/nbfacin/kasus/{caseId}/klausa.
func (s *Service) BacaKlausa(ctx context.Context, id string) ([]models.KlausaKasus, error) {
	if s.klausa == nil {
		return nil, ErrKlausaTanpaDatabase
	}
	if !idKasusSah(id) {
		return nil, ErrKasusTidakAda
	}
	baris, err := s.klausa.BacaKlausa(ctx, id)
	if errors.Is(err, repository.ErrKasusTidakAda) {
		return nil, ErrKasusTidakAda
	}
	return baris, err
}

// GantiKlausa - PUT /api/nbfacin/kasus/{caseId}/klausa (Save tab, SaveFacIn_Act): ganti utuh, jawab baca ulang.
func (s *Service) GantiKlausa(ctx context.Context, pelaku inti.Pelaku, id string, baris []models.KlausaKasus) ([]models.KlausaKasus, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return nil, err
	}
	if err := periksaKlausa(baris); err != nil {
		return nil, err
	}
	if s.klausa == nil || s.transaksi == nil {
		return nil, ErrKlausaTanpaDatabase
	}
	if !idKasusSah(id) {
		return nil, ErrKasusTidakAda
	}
	err := s.transaksi(ctx, func(tx *db.Tx) error { return s.klausa.GantiKlausa(ctx, tx, id, baris) })
	if errors.Is(err, repository.ErrKasusTidakAda) {
		return nil, ErrKasusTidakAda
	}
	if err != nil {
		return nil, err
	}
	return s.BacaKlausa(ctx, id)
}

// CariKlausa - GET /api/nbfacin/klausa?bahasa=&q=&halaman=: bahasa kode "0" (bawaan) / "1" / "2"; q = Info LIKE
// '%q%' TIDAK peka huruf (K47-7); halaman >= 1, di luar jangkauan = baris kosong.
func (s *Service) CariKlausa(ctx context.Context, bahasa, kata string, nomor int) (HalamanKlausa, error) {
	bahasa = strings.TrimSpace(bahasa)
	if bahasa == "" {
		bahasa = "0"
	}
	label, sah := labelBahasaKlausa[bahasa]
	switch {
	case !sah:
		return HalamanKlausa{}, fmt.Errorf("%w: bahasa harus 0, 1, atau 2", ErrMasukanKlausa)
	case len(kata) > lebarKataKlausa:
		return HalamanKlausa{}, fmt.Errorf("%w: q paling banyak %d bita", ErrMasukanKlausa, lebarKataKlausa)
	case nomor < 1:
		return HalamanKlausa{}, fmt.Errorf("%w: halaman mulai dari 1", ErrMasukanKlausa)
	}
	if s.klausa == nil {
		return HalamanKlausa{}, ErrKlausaTanpaDatabase
	}
	baris, total, err := s.klausa.CariKlausa(ctx, bahasa, kata, nomor, UkuranHalamanKlausa)
	if err != nil {
		return HalamanKlausa{}, err
	}
	for i := range baris {
		baris[i].Language = label
	}
	return HalamanKlausa{Baris: baris, Total: total, Nomor: nomor, Ukuran: UkuranHalamanKlausa}, nil
}

// ArgumenKlausa - GET /api/nbfacin/klausa/{id}/argumen.
func (s *Service) ArgumenKlausa(ctx context.Context, id string) ([]models.ArgumenKlausa, error) {
	id = strings.TrimSpace(id)
	if id == "" || len(id) > lebarTeksKlausa {
		return nil, fmt.Errorf("%w: id klausa kosong atau terlalu panjang", ErrMasukanKlausa)
	}
	if s.klausa == nil {
		return nil, ErrKlausaTanpaDatabase
	}
	return s.klausa.ArgumenKlausa(ctx, id)
}
