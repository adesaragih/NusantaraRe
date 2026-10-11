package repository

// Simulasi lintasan tulis GantiObjekAneka tanpa basis data (bug DEV 10-10-2026, ORA-01722 T_SPREADINGLIST): traversal,
// pengelompokan, periksaBindBaris, dan pilihan array-bind / per baris (jalankanKelompok) yang SAMA dengan jalur Oracle; hanya
// nomor sequence (teks digit berurutan) dan eksekusi (direkam, tidak dikirim) yang ditiru. Dipakai uji services - jalur PUT
// spreading-aneka utuh dengan bind yang benar-benar akan dikirim.

import (
	"context"
	"fmt"

	"nusantarare/modul/nbfacin/backend/models"
)

// EksekusiSimulasi - satu eksekusi INSERT simulasi: tabel, teks SQL, baris (nilai per bind) yang dikirim bersama.
type EksekusiSimulasi struct {
	Tabel, SQL string
	Baris      [][]any
}

// SimulasiTulisAneka - lihat kepala berkas; id = ID case (PARENT_ID T_LOCATIONLIST). Galat = galat yang akan dikembalikan
// sebelum Oracle (mis. ErrBindTidakSah).
func SimulasiTulisAneka(ctx context.Context, id string, lokasi []models.LokasiAneka) ([]EksekusiSimulasi, error) {
	t, err := tabelAnekaDari(func(n string) (string, error) { return n, nil })
	if err != nil {
		return nil, err
	}
	e := &eksekutorSimulasi{}
	n := 0
	uid := func() (string, error) {
		n++
		return fmt.Sprintf("00000000-0000-4000-8000-%012d", n), nil
	}
	if _, _, err := sisipAnekaBerkelompok(ctx, e, t, id, lokasi, uid); err != nil {
		return e.eksekusi, err
	}
	return e.eksekusi, nil
}

// eksekutorSimulasi - eksekutorKelompok tanpa basis data.
type eksekutorSimulasi struct {
	nomor    int64
	eksekusi []EksekusiSimulasi
}

func (e *eksekutorSimulasi) nomorBerkelompok(_ context.Context, sequence []string, n int) ([][]string, error) {
	hasil := make([][]string, len(sequence))
	for i := range sequence {
		for range n {
			e.nomor++
			hasil[i] = append(hasil[i], fmt.Sprint(e.nomor))
		}
	}
	return hasil, nil
}

func (e *eksekutorSimulasi) sisipArray(_ context.Context, q, tabel string, kolom []any, n int) error {
	baris := make([][]any, n)
	for r := range baris {
		baris[r] = make([]any, len(kolom))
		for c := range kolom {
			baris[r][c] = kolom[c].([]any)[r]
		}
	}
	e.eksekusi = append(e.eksekusi, EksekusiSimulasi{Tabel: tabel, SQL: q, Baris: baris})
	return nil
}
