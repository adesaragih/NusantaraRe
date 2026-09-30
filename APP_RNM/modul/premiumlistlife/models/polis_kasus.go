package models

// Kasus polis baru - tombol portal `Input Offer` / `Input Premium`
// (GILIRAN-13 paket 1, butir bn). MURNI.
//
// `[terverifikasi]` dibaca sebagai pohon 29-09-2026 (nol `pyStepsBlockName //`):
//
//	Section/PremiumList.xml     `Input Offer`   -> CreateInputLife, FlagPolicy "0" (b3310, b3597)
//	                            `Input Premium` -> CreateInputLife, FlagPolicy "1" (b3958, b4233)
//	Activity/CreateInputLife    b363 `FlowType = "pyStartCase"`, b444 `Call svcAddWorkObject`,
//	                            b618 `curWorkPage.FlagOnGoingPolicy = Param.FlagPolicy`,
//	                            b726 `Obj-Save`, b858 `Commit`
//	InputPolicyHolder.xml       Start1 (b1968) -> [Always] -> Assignment2 (b1938),
//	                            `pyWorkStatus` "Input Offer Life" (b1340)
//	Activity/InputOfferLife_preAct  b271-272 `pyWorkPage.Position = "Offer"`
//	                            (pre-activity flow action `InputDataOfferLife`)
//
// ⛔ KEDUA TOMBOL MULAI DI TAHAP YANG SAMA. Konektor pertama flow `[Always]`
// ke Assignment2, apa pun benderanya. Bendera baru bekerja di `Decision3`
// (`IsFlagOnGoingPolicy`) SESUDAH `Confirm`: decision table-nya memetakan
// "0" -> Offer (b293 -> b328) dan "1" -> Premium (b294 -> b329). Brief
// GILIRAN-13 menduga bendera memilih tahap AWAL; pohon flow membantahnya.
//
// ⚠️ Posisi "Offer" disetel Pega saat flow action tahap pertama DIBUKA; di
// sini ia disetel saat kasus lahir, sebab kotak masuk menyaring posisi dan
// kasus tanpa posisi tidak tampil di tab mana pun.
//
// Dibaca sesudah: polis_penawaran.go.

import (
	"errors"
	"fmt"

	inti "nusantarare/inti/backend"
)

// Nilai `FlagOnGoingPolicy` - VERBATIM `FlagPolicy` tombol portal (butir bn).
const (
	// FlagPolisPenawaran - `Input Offer`.
	FlagPolisPenawaran = "0"
	// FlagPolisPremium - `Input Premium`.
	FlagPolisPremium = "1"
)

// ErrFlagPolisTidakSah - bendera di luar kedua nilai tombol.
var ErrFlagPolisTidakSah = errors.New(
	`models: FlagOnGoingPolicy hanya "0" (Input Offer) atau "1" (Input Premium)`)

// KasusPolisBaru adalah keadaan awal sebuah work object polis.
type KasusPolisBaru struct {
	Lini   string
	Posisi string
	Status string
	Flag   string
}

// SusunKasusPolisBaru menyusun keadaan awal dari bendera tombolnya.
//
// ⛔ Nilai dibandingkan apa adanya, tanpa `TrimSpace`: bendera datang dari
// tombol, bukan ketikan, dan nilai yang tidak persis sama adalah cacat klien.
func SusunKasusPolisBaru(flag string) (KasusPolisBaru, error) {
	if flag != FlagPolisPenawaran && flag != FlagPolisPremium {
		return KasusPolisBaru{}, fmt.Errorf("%w: %q", ErrFlagPolisTidakSah, flag)
	}
	return KasusPolisBaru{
		Lini:   inti.LiniLife,
		Posisi: PosisiOffer,
		Status: TahapPolisPenawaran,
		Flag:   flag,
	}, nil
}
