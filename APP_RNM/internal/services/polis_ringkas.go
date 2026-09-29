package services

// `PolicyDataLife` untuk Claim Life - butir pl4, dipakai butir av.
//
// Untuk apa berkas ini: menyerahkan sepuluh medan polis kepada Claim Life,
// dalam bentuk yang layar Register dan Outstanding pakai.
//
// ⛔ CLAIM LIFE TIDAK MEMBACA CERMIN JSON-NYA `[keputusan work owner, av]`.
// Sumbernya `T_PREMIUM_LIST`, tabel relasional modul PremiumList Life.
//
// ⛔ TIGA MEDAN TETAP TANPA SUMBER, dan itu DINYATAKAN, bukan diisi teks
// kosong: `TanggalRespon`, `TanggalKonfirmasi`, `TanggalRealisasi` tidak punya
// kolom di migrasi 050-056 mana pun. Layar tetap menampilkannya - himpunan
// medannya harus tetap sama dengan `InputRegisterClaimLife.xml` - dengan
// penanda bahwa sumbernya belum ada. Menghilangkannya membuat paritas tampak
// lengkap padahal tidak, dan tidak ada yang akan mencarinya lagi.
//
// Dibaca sesudah: polis_inbox.go.

import (
	"context"
	"fmt"

	"nusantarare/internal/repository"
	"nusantarare/inti"
	"nusantarare/inti/db"
	"nusantarare/inti/galat"
)

// ErrPolisNomorTakDitemukan dirujuk ulang supaya handler tidak perlu
// mengimpor repository.
var ErrPolisNomorTakDitemukan = repository.ErrPolisNomorTakDitemukan

// MedanTanpaSumberPolicyData menyebut medan layar Claim Life yang belum punya
// kolom di mana pun.
//
// ⛔ DIKIRIM KE LAYAR, bukan disimpan di komentar Go. Orang yang membandingkan
// layar baru dengan layar lama akan menghitung medannya.
var MedanTanpaSumberPolicyData = []string{
	"TanggalRespon", "TanggalKonfirmasi", "TanggalRealisasi",
	// av-2 (GILIRAN-11 paket 2): `Confirmation Reserved` b14198 dan
	// `Underwriter Note` b14806 - nol kolom di migrasi 050-056.
	"TanggalKonfirmasiBalik", "KetentuanUnderwriting",
}

// PolicyDataLife adalah bentuk yang layar Claim Life pakai.
//
// ⚠️ Nama kuncinya sengaja sama dengan nama properti Pega
// (`.PolicyDataLife.*`), supaya siapa pun yang membandingkan layar dengan
// rule tidak perlu menerjemahkan.
type PolicyDataLife struct {
	NomorPolis       string `json:"nomorPolis"`
	Type             string `json:"type"`
	MarketingName    string `json:"marketingName"`
	CedingCoName     string `json:"cedingCoName"`
	PolicyHolderName string `json:"policyHolderName"`
	BusinessName     string `json:"businessName"`
	// DateReceived `YYYY-MM-DD`; kosong berarti kolomnya belum diisi.
	DateReceived string `json:"dateReceived"`
	Status       string `json:"status"`
	StatusUpdate string `json:"statusUpdate"`
	// ProductNameID adalah KUNCI ambang batas hari (butir ba).
	ProductNameID string `json:"productNameId"`
	ProductName   string `json:"productName"`
	// ProdKe adalah versi polis yang terbaca.
	//
	// ⛔ IKUT DIKIRIM. Satu nomor polis punya banyak versi, dan yang dibaca
	// adalah `PROD_KE` TERBESAR. Layar yang tidak dapat menyebut versi mana
	// yang ditampilkannya membuat selisih angka mustahil ditelusuri.
	ProdKe int `json:"prodKe"`
	// MedanTanpaSumber menyebut medan layar yang belum punya kolom.
	MedanTanpaSumber []string `json:"medanTanpaSumber"`

	// ⭐ av-2 (GILIRAN-11 paket 2) - tujuh medan lagi dari `T_PREMIUM_LIST`.
	// `TypeCeding` kode, `TypeCedingName` katanya; layar menampilkan kata.
	TypeCeding     string `json:"typeCeding"`
	TypeCedingName string `json:"typeCedingName"`
	ProRateType    string `json:"proRateType"`
	// WPC `YYYY-MM-DD`; kosong berarti kolomnya belum diisi.
	WPC string `json:"wpc"`
	// RetroName dan SecurityReinsurer hanya TAMPIL bagi Type TP/TR
	// (`pyCondition` b10541/b10824); keputusan tampilnya milik layar.
	RetroName         string `json:"retroName"`
	SecurityReinsurer string `json:"securityReinsurer"`
	SobName           string `json:"sobName"`
}

// RingkasPolis melayani pembacaan `PolicyDataLife`.
type RingkasPolis struct{ svc *Service }

// RingkasPolis menyusun layanannya.
func (s *Service) RingkasPolis() *RingkasPolis { return &RingkasPolis{svc: s} }

// Ambil membaca data polis versi berjalan untuk sebuah nomor polis.
func (r *RingkasPolis) Ambil(ctx context.Context, pelaku inti.Pelaku, nomorPolis string) (
	PolicyDataLife, error) {

	if err := inti.WajibIdentitas(pelaku); err != nil {
		return PolicyDataLife{}, err
	}
	if r == nil || r.svc == nil || !r.svc.PunyaDatabase() {
		return PolicyDataLife{}, db.ErrTanpaOracle
	}
	if nomorPolis == "" {
		return PolicyDataLife{}, fmt.Errorf("%w: nomor polis kosong", galat.ErrPermintaanTidakSah)
	}
	p, err := repository.NewRingkasPolisLife(r.svc.DB()).Ringkas(ctx, nomorPolis)
	if err != nil {
		return PolicyDataLife{}, err
	}
	hasil := PolicyDataLife{
		NomorPolis:       p.NomorPolis,
		Type:             p.Type,
		MarketingName:    p.MarketingName,
		CedingCoName:     p.CedingCoName,
		PolicyHolderName: p.PolicyHolderName,
		BusinessName:     p.BusinessName,
		Status:           p.Status,
		StatusUpdate:     p.StatusUpdate,
		ProductNameID:    p.ProductNameID,
		ProductName:      p.ProductName,
		ProdKe:           p.ProdKe,
		MedanTanpaSumber: MedanTanpaSumberPolicyData,

		TypeCeding:        p.TypeCeding,
		TypeCedingName:    p.TypeCedingName,
		ProRateType:       p.ProRateType,
		RetroName:         p.RetroName,
		SecurityReinsurer: p.SecurityReinsurer,
		SobName:           p.SobName,
	}
	if p.DateReceived != nil {
		hasil.DateReceived = *p.DateReceived
	}
	if p.WPC != nil {
		hasil.WPC = *p.WPC
	}
	return hasil, nil
}
