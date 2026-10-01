package services

// Simpan produk (paket 3: sisi umum, tiket 02) - padanan tombol `Save` b58998
// → `SaveProductName_Confirm` (submit `Save` b34) → `SaveProductName_Act`
// (`pyLocalActionActivity` b101) TANPA prosedur `PEGA_M_PRODUCT_LIFE`:
//
//	1  b359  `·` PRE=false  UPDATEOP ← operator; POLICYHODER/POLICYHODERNAME ← inward
//	6  b1368 `·` PRE=true `@PropertyHasValue(CREATEOP)` T=3 F=2  CREATEOP ← operator bila kosong
//	8  b1623 `·` PRE=false  JSON halaman ProductName
//	9  b1831 `·` PRE=false  simpan JSON - ID baru '1' ‖ LPAD(M_PRODUCT_LIFE_SEQ, 5) (P6)
//	10 b2019 `·` PRE=false  kolom datar RIRISKID, RIRISK (R8)
//	11 b2207 `·` PRE=false  ID dikembalikan, form ditutup
//
// ⛔ Baru = POST, ubah = PUT: identitas tidak pernah dari klien (ADR-0006).
// ⛔ Medan milik server (pembuat, pengubah, medan layar mati, riwayat
// komentar, daftar outward) TIDAK PERNAH dari klien - diambil dari baris
// tersimpan di transaksi yang sama.
// ⛔ Satu transaksi; galat di mana pun = nol tulisan (gagal terang-terangan).

import (
	"context"
	"fmt"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/masterproductnamelife/backend/models"
	"nusantarare/modul/masterproductnamelife/backend/repository"
)

// GudangTulis - penulis produk (bagian Gudang).
type GudangTulis interface {
	// KunciProduk membaca produk utuh dan mengunci barisnya (FOR UPDATE).
	KunciProduk(ctx context.Context, tx *db.Tx, id string) (models.Produk, error)
	// SisipProduk menerbitkan ID baru dari sequence dan menulis produk baru.
	SisipProduk(ctx context.Context, tx *db.Tx, p models.Produk) (string, error)
	// PerbaruiProduk menulis ulang produk; kunci JSON tak dikelola dipertahankan.
	PerbaruiProduk(ctx context.Context, tx *db.Tx, p models.Produk) error
}

// Label VERBATIM medan form (PARITAS §3) - dipakai kalimat penolakan.
const (
	labelProductName = "Product Name"  // b3652
	labelCeding      = "Ceding"        // b4107
	labelSOB         = "SOB"           // b4495
	labelDeduction   = "Deduction (%)" // b7137
	labelRIRisk      = "R/I Risk Name" // b7430
	labelCause       = "Cause Of Loss" // b10693
)

// periksaUmum - gerbang murni sisi umum (angka, panjang kolom datar).
func periksaUmum(pk *periksa, u *models.ProdukUmum) {
	pk.desimal(labelDeduction, &u.RIComm)
	pk.panjang(labelProductName, u.ProductName, repository.LebarProductName)
	pk.panjang(labelRIRisk, u.RIRisk, repository.LebarRIRisk)
	pk.panjang(labelRIRisk+" ID", u.RIRiskID, repository.LebarRIRiskID)
}

// pilihan - satu medan yang diisi pemilih master (`Choose*` → `set*_DT`).
type pilihan struct {
	jenis            models.JenisMaster
	label            string
	id, nama         *string
	idLama, namaLama string
}

// periksaPilihan - tiket 04 AC: pilihan yang tidak ada di master ditolak; nama
// tersimpan dibangun ulang dari ID master. Hanya nilai yang BERUBAH dari yang
// tersimpan yang diperiksa (data lama yang tidak disentuh tidak ditolak).
func (l *Layanan) periksaPilihan(ctx context.Context, pk *periksa, daftar []pilihan) error {
	for _, c := range daftar {
		if *c.id == c.idLama && *c.nama == c.namaLama {
			continue
		}
		if *c.id == "" {
			if *c.nama != "" {
				pk.tolak("%s %q must be chosen from the master list", c.label, *c.nama)
			}
			continue
		}
		v, ada, err := l.gudang.AmbilMaster(ctx, c.jenis, *c.id)
		if err != nil {
			return err
		}
		if !ada {
			pk.tolak("%s %q is not in the master list", c.label, *c.id)
			continue
		}
		*c.nama = v.Nama
	}
	return nil
}

// pilihanUmum - keempat pemilih sisi umum.
func pilihanUmum(m *models.Produk, lama models.ProdukUmum) []pilihan {
	return []pilihan{
		{models.MasterCeding, labelCeding, &m.Umum.CedingID, &m.Umum.Ceding, lama.CedingID, lama.Ceding},
		{models.MasterSOB, labelSOB, &m.Umum.SOBID, &m.Umum.SOBName, lama.SOBID, lama.SOBName},
		{models.MasterRIRisk, labelRIRisk, &m.Umum.RIRiskID, &m.Umum.RIRisk, lama.RIRiskID, lama.RIRisk},
		{models.MasterPenyebab, labelCause, &m.Umum.CauseID, &m.Umum.Cause, lama.CauseID, lama.Cause},
	}
}

// lengkapiMilikServer menimpa medan milik server dari baris tersimpan (lama
// nil = produk baru).
func lengkapiMilikServer(m *models.Produk, lama *models.Produk, p inti.Pelaku) {
	var simpan models.Produk
	if lama != nil {
		simpan = *lama
	}
	u, su := &m.Umum, simpan.Umum
	// Medan layar mati (`OTHER 1=2`) - kuncinya dibaca view, nilainya dari JSON lama.
	u.TypeBasicRider, u.TypeCeding, u.Grup, u.ProductCode = su.TypeBasicRider, su.TypeCeding, su.Grup, su.ProductCode
	u.ProductType, u.ProductTypeID, u.RIRate, u.RIRateID = su.ProductType, su.ProductTypeID, su.RIRate, su.RIRateID
	u.RICommID, u.OutwardName, u.OutwardNameID = su.RICommID, su.OutwardName, su.OutwardNameID
	u.OutwardRate, u.OutwardRateID, u.OutwardComm, u.OutwardCommID = su.OutwardRate, su.OutwardRateID, su.OutwardComm, su.OutwardCommID
	u.Benefit, u.BenefitID = su.Benefit, su.BenefitID
	// Langkah 6 b1368: CREATEOP diisi hanya bila kosong; produk baru (termasuk
	// salinan `Copy`) = pelaku (OQ-MPNL-13). Langkah 1 b359: UPDATEOP = pelaku.
	u.CreateOp = su.CreateOp
	if lama == nil || u.CreateOp == "" {
		u.CreateOp = p.AkunID
	}
	u.UpdateOp = p.AkunID
	// Langkah 1 b359: pemegang polis disalin dari halaman inward ke halaman umum.
	u.PolicyHolder, u.PolicyHolderName = m.Inward.PolicyHolder, m.Inward.PolicyHolderName
	// Riwayat komentar dan daftar outward milik server.
	m.CommentList = append([]models.BarisKomentar{}, simpan.CommentList...)
	m.OutwardList = append([]models.BarisOutward{}, simpan.OutwardList...)
}

// SimpanProduk menyimpan produk baru (`baru`) atau mengubah yang ada.
func (l *Layanan) SimpanProduk(ctx context.Context, p inti.Pelaku, m models.Produk, baru bool) (models.Produk, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return models.Produk{}, err
	}
	if baru && m.ID != "" {
		return models.Produk{}, ErrIDDariKlien
	}
	var pk periksa
	periksaUmum(&pk, &m.Umum)
	var hasil models.Produk
	err := l.tx(ctx, func(tx *db.Tx) error {
		var lama *models.Produk
		if !baru {
			s, err := l.gudang.KunciProduk(ctx, tx, m.ID)
			if err != nil {
				return err
			}
			lama = &s
		}
		var su models.ProdukUmum
		if lama != nil {
			su = lama.Umum
		}
		if err := l.periksaPilihan(ctx, &pk, pilihanUmum(&m, su)); err != nil {
			return err
		}
		if err := pk.galat(); err != nil {
			return err
		}
		lengkapiMilikServer(&m, lama, p)
		if baru {
			id, err := l.gudang.SisipProduk(ctx, tx, m)
			if err != nil {
				return err
			}
			m.ID = id
		} else if err := l.gudang.PerbaruiProduk(ctx, tx, m); err != nil {
			return err
		}
		var err error
		hasil, err = l.gudang.AmbilProduk(ctx, tx, m.ID)
		return err
	})
	if err != nil {
		return models.Produk{}, tidakAda(err, ErrProdukTidakAda, m.ID)
	}
	l.catat(fmt.Sprintf("master product name life: product %s saved (new=%v) at %s", hasil.ID, baru,
		time.Now().UTC().Format(time.RFC3339)))
	return hasil, nil
}
