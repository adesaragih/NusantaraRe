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
	"strings"
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

// Pesan wajib-isi VERBATIM - `SaveProductName_Act` langkah 1 b359 menyetel
// `local.errMsg1..4`, langkah 2–5 (PRE=true, WHEN `== ""`, F=3) memasangnya.
const (
	PesanProductNameKosong   = "Product Name Empty"  // b431 - langkah 2 b668, WHEN b805 `ProductName.PRODUCTNAME==""`
	PesanCedingKosong        = "Ceding Empty"        // b452 - langkah 3 b843, WHEN b980 `ProductName.CEDING==""`
	PesanPemegangPolisKosong = "Policy Holder Empty" // b473 - langkah 4 b1018, WHEN b1155 `ProductNameInward.POLICYHODERNAME==""`
	PesanSOBKosong           = "SOB Empty"           // b494 - langkah 5 b1193, WHEN b1330 `ProductName.SOBNAME==""`
)

// periksaWajibIsi - SATU aturan wajib-isi untuk setiap jalur simpan (tiket 05).
//
// ⛔ Hanya keempat pemeriksaan HIDUP Pega. "Tipe & grup terisi" BUKAN
// pemeriksaan: langkah 1 b359 adalah `Property-Set` ber-PRE=false, dan medan
// `TYPE`/`GRUP` mati (`1=2`) - menegakkannya membuat setiap simpan gagal (R7).
// Pega berhenti di pesan pertama (transisi `1==1` → 6); di sini SEMUA
// dilaporkan sekaligus (keputusan tertulis tiket 05). Spasi = kosong.
func periksaWajibIsi(pk *periksa, m *models.Produk) {
	for _, w := range []struct {
		nilai, pesan string
	}{
		{m.Umum.ProductName, PesanProductNameKosong},
		{m.Umum.Ceding, PesanCedingKosong},
		{m.Inward.PolicyHolderName, PesanPemegangPolisKosong},
		{m.Umum.SOBName, PesanSOBKosong},
	} {
		if strings.TrimSpace(w.nilai) == "" {
			pk.tolak("%s", w.pesan)
		}
	}
}

// Label VERBATIM medan inward (PARITAS §3.2).
const (
	labelPolicyHolder = "Policy Holder" // b17129
	labelCurrency     = "Currency"      // b28173
	labelBegin        = "Begin Date"    // b22001
	labelSTNC         = "STNC"          // b22336
	labelMature       = "Expired Date"  // b27282
	labelMinAge       = "Minimum Age (Years)"
	labelMaxAge       = "Maximum Age (Years)"
	labelMinSI        = "Min Sum Insured"
	labelMaxSI        = "Max Sum Insured"
)

// periksaInward - gerbang murni sisi inward: setiap medan `pxNumber` desimal,
// ketiga tanggal, dan rentang (R17).
func periksaInward(pk *periksa, i *models.ProdukInward) {
	angka := []struct {
		label string
		v     *string
	}{
		{"Addendum No.", &i.AddendumNo}, {"Amandement No.", &i.AmandementNo},
		{"Max Notification Claim Expired", &i.MaxExpiredClaim}, {"Ceding Retention (%)", &i.CedingRetentionNum},
		{"Ceding's Limit", &i.CedingLimit}, {"Brokerage Fee (%)", &i.Brokerage}, {"Expiry Age (Years)", &i.ExpiryAge},
		{"Extra Premium", &i.ExtraPremi}, {"Max Sum Reasured", &i.MaxSumReasured}, {"Nusantara Re Share (%)", &i.RNMShare},
		{"Nusantara Re's Limit", &i.RNMLimitNum}, {"Premium Factor (%)", &i.PremiumFactor},
		{"Annuity Interest (%)", &i.AnnuityInterest}, {"Premium Refund Factor (%)", &i.PremiumRefundFactor},
		{"Max Production Data Receive", &i.MaxDataReceive}, {"Extra Mortality (%)", &i.ExtraMortality},
		{"Max Contract (year)", &i.MaxContract}, {"Proportional Table", &i.ProportionalTable},
	}
	for _, a := range angka {
		pk.desimal(a.label, a.v)
	}
	minUsia, maxUsia := pk.desimal(labelMinAge, &i.MinAge), pk.desimal(labelMaxAge, &i.MaxAge)
	minSI, maxSI := pk.desimal(labelMinSI, &i.MinSumInsured), pk.desimal(labelMaxSI, &i.MaxSumInsured)
	pk.tidakLebihBesar(labelMinAge, minUsia, labelMaxAge, maxUsia)
	pk.tidakLebihBesar(labelMinSI, minSI, labelMaxSI, maxSI)
	pk.tanggal(labelBegin, &i.Begin)
	pk.tanggal(labelSTNC, &i.STNC)
	pk.tanggal(labelMature, &i.Mature)
}

// pilihanInward - kedua pemilih sisi inward.
func pilihanInward(m *models.Produk, lama models.ProdukInward) []pilihan {
	return []pilihan{
		{models.MasterPemegangPolis, labelPolicyHolder, &m.Inward.PolicyHolder, &m.Inward.PolicyHolderName,
			lama.PolicyHolder, lama.PolicyHolderName},
		{models.MasterMataUang, labelCurrency, &m.Inward.CurrencyID, &m.Inward.Currency, lama.CurrencyID, lama.Currency},
	}
}

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
		// Spasi = kosong (sama dengan wajib-isi); ID master tidak berspasi tepi.
		*c.id = strings.TrimSpace(*c.id)
		if strings.TrimSpace(*c.nama) == "" {
			*c.nama = ""
		}
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
//
// `sumber` = produk tersimpan (ubah) atau produk asal `Copy` (baru, `CopyProduct`
// b138/b161 hanya mengosongkan kedua ID - medan lain ikut tersalin); nil = baru.
func lengkapiMilikServer(m *models.Produk, sumber *models.Produk, baru bool, p inti.Pelaku) {
	var simpan models.Produk
	if sumber != nil {
		simpan = *sumber
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
	if baru || u.CreateOp == "" {
		u.CreateOp = p.AkunID
	}
	u.UpdateOp = p.AkunID
	// Sisi inward: identitas dan kunci view tanpa medan form milik server.
	in, si := &m.Inward, simpan.Inward
	in.Ceding, in.TreatyNumber, in.InwardTreatyNm = si.Ceding, si.TreatyNumber, si.InwardTreatyNm
	in.CedingRetentionPct, in.CedingLimitXPN, in.RNMLimitPct = si.CedingRetentionPct, si.CedingLimitXPN, si.RNMLimitPct
	in.LienClause, in.Months = si.LienClause, si.Months
	in.ID, in.ProductID = si.ID, m.ID
	if baru {
		in.ID = "" // `CopyProduct` 2 b161; penulis memberi ID = ID produk (R14)
	}
	// Langkah 1 b359: pemegang polis disalin dari halaman inward ke halaman umum.
	u.PolicyHolder, u.PolicyHolderName = m.Inward.PolicyHolder, m.Inward.PolicyHolderName
	// Riwayat komentar dan daftar outward milik server (salinan mewarisinya, seperti Pega).
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
	if !baru && m.SalinanDari != "" {
		return models.Produk{}, ErrSalinanPadaUbah
	}
	var pk periksa
	periksaWajibIsi(&pk, &m)
	periksaUmum(&pk, &m.Umum)
	periksaInward(&pk, &m.Inward)
	periksaPlan(&pk, m.PlanList)
	periksaUWLimit(&pk, m.UnderwritingLimit)
	periksaFinUW(&pk, m.FinancialUnderwriting)
	var hasil models.Produk
	err := l.tx(ctx, func(tx *db.Tx) error {
		var lama *models.Produk
		if !baru {
			s, err := l.gudang.KunciProduk(ctx, tx, m.ID)
			if err != nil {
				return err
			}
			lama = &s
		} else if m.SalinanDari != "" {
			s, err := l.gudang.AmbilProduk(ctx, tx, m.SalinanDari)
			if err != nil {
				return tidakAda(err, ErrSalinanTidakAda, m.SalinanDari)
			}
			lama = &s
		}
		var tersimpan models.Produk
		if lama != nil {
			tersimpan = *lama
		}
		if err := l.periksaPilihan(ctx, &pk, append(pilihanUmum(&m, tersimpan.Umum),
			pilihanInward(&m, tersimpan.Inward)...)); err != nil {
			return err
		}
		if err := l.periksaPilihanPlan(ctx, &pk, m.PlanList, tersimpan.PlanList); err != nil {
			return err
		}
		if err := pk.galat(); err != nil {
			return err
		}
		lengkapiMilikServer(&m, lama, baru, p)
		if perluHitungOutward(&m, baru) {
			if err := l.hitungOutward(ctx, tx, &m); err != nil {
				return err
			}
		}
		// Langkah 7 b1513 `·` (tanpa prakondisi): `AddCommentList_Act` - SETIAP
		// simpan menambah satu baris, juga bila komentarnya kosong (OQ-MPNL-14).
		m.CommentList = append(m.CommentList, barisKomentar(l.jam(), p.AkunID, m.Umum.Comment))
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
