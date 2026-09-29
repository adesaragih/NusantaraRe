package models

// Reinsurer pada kombinasi (tahun, grup, jenis) - tiket 05 Treaty Contract Out.
//
// Untuk apa berkas ini: gerbang simpan reinsurer dan aturan UANG-nya - share
// dan komisi sebagai desimal presisi arbitrer (ADR-0003, AC 16), total share.
//
// Pohon yang ditiru, dibaca 29-09-2026 (nomor baris mentah):
//
//	`Activity/SetErrorMessageReinsurer.xml`     langkah 1 koma -> titik; langkah 2-3
//	                                           0 <= PctShare, Ricomm <= 100
//	                                           (`SetErrorMessageBetween` b452/b631)
//	`Activity/SaveTreatyReinsurerDetail1_Act.xml` langkah 3 b694 ReinsurerID wajib;
//	                                           langkah 7 b1357 + 9: TotalShare +
//	                                           (PctShare - PctShare1) <= 100.000 (b1452)
//
// ⛔ Uji aturan ini DIBUKTIKAN MERAH lebih dulu terhadap implementasi float64
// (0.1 + 0.2 = 0.30000000000000004; tiga sepertiga kehilangan skalanya).
//
// Dibaca sesudah: tco_kontrak.go.

import (
	"errors"
	"fmt"
	"strings"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/pkg/utils"
)

var (
	// ErrReinsurerKosong - VERBATIM `SaveTreatyReinsurerDetail1_Act.xml` b694.
	ErrReinsurerKosong = errors.New("models: Data tidak boleh kosong...!!! - ReinsurerID wajib dipilih")
	// ErrPersenKosong - share dan komisi wajib (AC 14 "masing-masing dengan
	// share dan komisi"); Pega tidak menggerbanginya.
	ErrPersenKosong = errors.New("models: nilai persen wajib diisi")
	// ErrPersenBukanDesimal - bukan desimal, atau melampaui NUMBER(38,8).
	ErrPersenBukanDesimal = errors.New("models: nilai persen bukan desimal yang sah")
	// ErrPersenDiLuarRentang - `SetErrorMessageBetween` (teks Pega tidak diekspor).
	ErrPersenDiLuarRentang = errors.New("models: nilai persen harus di antara 0 dan 100")
	// ErrTotalShareMelebihi100 - VERBATIM `SaveTreatyReinsurerDetail1_Act.xml` b1357.
	ErrTotalShareMelebihi100 = errors.New("models: Persentase tidak boleh lebih dari 100!")
)

const (
	skalaPersenTCO     = 8
	digitBulatPersen   = 30
	batasAtasPersenTCO = "100"
)

var (
	seratus = apd.New(100, 0)
)

// UraiPersenMasukTCO menormalkan satu nilai persen di BATAS MASUKAN.
//
// ⛔ Sekali, di sini - bukan `@replaceAll` yang ditempel di tengah alur simpan
// (`SetErrorMessageReinsurer` langkah 1). Koma DAN titik sekaligus ditolak:
// `1.234,5` bukan persen yang dapat ditebak maksudnya.
func UraiPersenMasukTCO(nama, teks string) (*apd.Decimal, error) {
	t := strings.TrimSpace(teks)
	if t == "" {
		return nil, fmt.Errorf("%w: %s", ErrPersenKosong, nama)
	}
	if strings.Contains(t, ",") {
		if strings.Contains(t, ".") {
			return nil, fmt.Errorf("%w: %s %q memuat koma dan titik sekaligus", ErrPersenBukanDesimal, nama, t)
		}
		t = strings.ReplaceAll(t, ",", ".")
	}
	d, err := utils.ParseDecimal(t)
	if err != nil {
		return nil, fmt.Errorf("%w: %s %q", ErrPersenBukanDesimal, nama, teks)
	}
	ringkas := new(apd.Decimal).Set(d)
	ringkas.Reduce(ringkas)
	if ringkas.Exponent < -skalaPersenTCO {
		return nil, fmt.Errorf("%w: %s %q lebih dari %d angka di belakang koma", ErrPersenBukanDesimal,
			nama, teks, skalaPersenTCO)
	}
	if ringkas.NumDigits()+int64(ringkas.Exponent) > digitBulatPersen {
		return nil, fmt.Errorf("%w: %s %q terlalu besar", ErrPersenBukanDesimal, nama, teks)
	}
	if d.Negative && !d.IsZero() || d.Cmp(seratus) > 0 {
		return nil, fmt.Errorf("%w: %s = %s", ErrPersenDiLuarRentang, nama, d.Text('f'))
	}
	return d, nil
}

// TotalShareTCO menjumlahkan share - `BrowseTreatyReinsurerList_Act` langkah 3.
//
// Sel NULL warisan dihitung nol: baris tanpa share tidak menggagalkan
// tampilnya total seluruh kombinasi.
func TotalShareTCO(shares []*apd.Decimal) (*apd.Decimal, error) {
	jumlah := new(apd.Decimal)
	for _, s := range shares {
		if s == nil {
			continue
		}
		if _, err := utils.DecimalContext().Add(jumlah, jumlah, s); err != nil {
			return nil, fmt.Errorf("models: menjumlahkan share: %w", err)
		}
	}
	return jumlah, nil
}

// PeriksaTotalShareTCO menegakkan batas 100 pada total BARU.
//
// `lain` = share seluruh reinsurer LAIN pada kombinasi itu (baris yang sedang
// diubah tidak termasuk) - padanan `TotalShare - PctShare1` Pega.
//
// ⚠️ Hanya LEBIH dari 100 yang ditolak. Kurang dari 100 bukan gerbang - di Pega
// pun tidak (`<= 100.000` b1452), dan tiket 05 menyatakannya `[terbuka]`.
func PeriksaTotalShareTCO(lain []*apd.Decimal, baru *apd.Decimal) (*apd.Decimal, error) {
	total, err := TotalShareTCO(append(append([]*apd.Decimal{}, lain...), baru))
	if err != nil {
		return nil, err
	}
	if total.Cmp(seratus) > 0 {
		return total, fmt.Errorf("%w (total share %s, batas %s)", ErrTotalShareMelebihi100,
			total.Text('f'), batasAtasPersenTCO)
	}
	return total, nil
}

// PeriksaReinsurerTCO menjalankan gerbang yang bukan angka.
func PeriksaReinsurerTCO(r ReinsurerTreaty) error {
	if strings.TrimSpace(r.ReinsurerID) == "" {
		return ErrReinsurerKosong
	}
	return nil
}
