package models

import (
	"fmt"
	"strings"
)

// Form Rate Detail - section Pega `InboxRIRate` (kelas `ASM-FW-GISFW-Int-M_RATE_LIFE`, judul "R/I RATE DETAIL" b367;
// View Detail.xml, work owner 06-10-2026): tambah dan ubah satu baris `M_RATE_LIFE` milik ringkasan yang sedang
// dilihat. R/I RATE NAME = `TempIDUsedBy.USEDBY` b1747 (disabled, wajib b1766) dan IDUSEDBY = `TempIDUsedBy.ID` b1526
// (tersembunyi `1=2` b1657): keduanya dari ringkasan, bukan isian pengguna.

// UkuranHalamanRate - grid Rate Detail `pyPageSize` Other b10280, `pyPageSizeOther` 20 b10206.
const UkuranHalamanRate = 20

// IsianRate - isian form Rate Detail (Save b3118 -> `AddToList_Act` b3196; Edit b9773 -> `EditList_DT` b9853).
type IsianRate struct {
	// Gender - `.GENDER` b1938 (pxRadioButtons, tidak wajib).
	Gender string `json:"gender"`
	// Contract - `.CONTRACT` b2122 (pxNumber, wajib b2141, placeholder 0 b2145).
	Contract string `json:"contract"`
	// Age - `.AGE` b2400 (pxNumber, tidak wajib, placeholder 0 b2418).
	Age string `json:"age"`
	// Rate - `.RATE` b2565 (pxNumber, wajib b2579, placeholder `0,0000` b2583).
	Rate string `json:"rate"`
}

// PeriksaIsianRate merapikan isian Rate Detail; galat = kalimat untuk pengguna. CONTRACT dan RATE wajib, GENDER
// dan AGE boleh kosong (XML). GENDER bila diisi U/M/F (pilihan radio tidak ada di XML - isi `RATE_LIFE` DEV); AGE dan
// CONTRACT bulat 0..UmurMaks; RATE disimpan berkoma desimal (`ChangeDotToPoint_DT` b2171, b2609).
func PeriksaIsianRate(isi IsianRate) (IsianRate, error) {
	var out IsianRate
	out.Gender = strings.ToUpper(strings.TrimSpace(isi.Gender))
	if out.Gender != "" && !sahGender(out.Gender) {
		return IsianRate{}, fmt.Errorf("GENDER must be U, M, or F (got %q)", isi.Gender)
	}
	var err error
	if out.Contract, err = normalBulat("CONTRACT", isi.Contract, false); err != nil {
		return IsianRate{}, err
	}
	if out.Age, err = normalBulat("AGE", isi.Age, true); err != nil {
		return IsianRate{}, err
	}
	if out.Rate, err = NormalRate(isi.Rate); err != nil {
		return IsianRate{}, err
	}
	return out, nil
}
