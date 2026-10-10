package services_test

import (
	"testing"

	"nusantarare/modul/treatyin/backend/models"
	"nusantarare/modul/treatyin/backend/services"
)

// sumberKosong - `SumberSpreading` yang nol baris; dipakai uji yang tidak
// menyentuh Report Definition spreading.
type sumberKosong struct{}

func (sumberKosong) Induk(string, string) []models.SusunanSpreading { return nil }
func (sumberKosong) Anak(string, string) []models.SusunanSpreading  { return nil }

func shareKosong() models.ShareNP {
	return models.ShareNP{Total: map[string][]models.NilaiMataUang{}}
}

// ⭐ Pembantu `nilai` dipakai ULANG dari `hitung_limit_test.go` — satu
// bentuk, satu tempat; `services.NilaiMataUang` alias `models`-nya.

func cari(daftar []models.NilaiMataUang, cur string) string {
	for _, b := range daftar {
		if b.Currency == cur {
			return b.Value
		}
	}
	return ""
}

// ⛔ LANGKAH [23] DIKOMENTARI DI PEGA — DAN AKU SEMPAT SALAH MEMBACANYA.
//
// Ronde ini sempat MENAMBAHKAN blok [23] (menjumlah larik baris Share
// sendiri, berikut `GrossPremiumMinList`) atas dasar `praAktif=false` yang
// kukira berarti "tanpa pra-syarat, jadi selalu jalan". Rekan sesi lain
// menunjukkan tandanya yang benar, dan pengukuran membenarkannya:
//
//	[22] blockName `Share`     praAktif true   → jalan
//	[23] blockName `//`        praAktif false  → DIKOMENTARI
//	[24] blockName ``          praAktif true   → jalan
//	[25] blockName `FShare`    praAktif true   → jalan
//	[26] blockName ``          praAktif false  → JALAN (mengisi TotalFacShare*)
//
// [26] itu buktinya: `praAktif=false` saja tidak mematikan apa pun.
// `pyStepsBlockName == "//"` yang mematikan.
//
// ⛔ RALAT 8 Oktober 2026: tangkapan layar Pega kontrak 1001855 menunjukkan
// total = larik SETIAP baris Share, baris bernama ikut. Uji di bawah
// memakai angka tangkapan itu.
func TestNPSetTotalShareSamaDenganPega1001855(t *testing.T) {
	s := shareKosong()
	baris := func(limit, gross string) models.BarisShareNP {
		return models.BarisShareNP{
			SpreadingTypeXOL:    "2022 QS 145M TRT",
			RnmLimitList:        []models.NilaiMataUang{nilai("IDR", limit)},
			GrossPremiumList:    []models.NilaiMataUang{nilai("IDR", gross)},
			GrossPremiumMinList: []models.NilaiMataUang{nilai("IDR", "7")},
			DeductionTotalList:  []models.NilaiMataUang{nilai("IDR", "0")},
			NetPremiumList:      []models.NilaiMataUang{nilai("IDR", gross)},
			SpreadingListXOL: []models.BarisSpreadingNP{{
				GrossPremiumList: []models.NilaiMataUang{nilai("IDR", "5")},
			}},
		}
	}
	s.Share = []models.BarisShareNP{
		baris("4000000000", "517646250"),
		baris("6250000000", "380036250"),
		baris("25000000000", "373353750"),
	}
	services.NPSetTotalShare(&s)
	for _, u := range []struct{ kunci, mau string }{
		{"TotalShareRnmNP", "35250000000"},
		{"TotalShareGrossNP", "1271036250"},
		{"TotalShareDeductionNP", "0"},
		{"TotalShareNetNP", "1271036250"},
	} {
		if got := cari(s.Total[u.kunci], "IDR"); got != u.mau {
			t.Errorf("%s = %q, mau %q (tangkapan Pega)", u.kunci, got, u.mau)
		}
	}
	// ⛔ `Total Gross Min Premium` TETAP kosong — "No items" di tangkapan yang sama.
	if len(s.Total["TotalShareGrossMinNP"]) != 0 {
		t.Errorf("TotalShareGrossMinNP = %+v, mau kosong", s.Total["TotalShareGrossMinNP"])
	}
}

// ⭐ Mata uang KOSONG tidak pernah masuk — `.Currency != ""`.
// Beda dengan cabang EGNPI dan retensi, yang menerimanya.
func TestMataUangKosongTidakMasukTotalShare(t *testing.T) {
	s := shareKosong()
	s.Share = []models.BarisShareNP{{
		GrossPremiumList: []models.NilaiMataUang{nilai("", "99"), nilai("IDR", "1")},
	}}
	services.NPSetTotalShare(&s)
	if len(s.Total["TotalShareGrossNP"]) != 1 {
		t.Fatalf("baris total = %+v, mau hanya IDR", s.Total["TotalShareGrossNP"])
	}
}

func TestNPSetTotalShareMengosongkanLebihDulu(t *testing.T) {
	s := shareKosong()
	s.Total["TotalShareGrossNP"] = []models.NilaiMataUang{nilai("IDR", "999")}
	s.Share = []models.BarisShareNP{{
		GrossPremiumList: []models.NilaiMataUang{nilai("IDR", "1")},
	}}
	services.NPSetTotalShare(&s)
	if got := cari(s.Total["TotalShareGrossNP"], "IDR"); got != "1" {
		t.Fatalf("nilai lama tidak dibuang: %q", got)
	}
}

// --- TreatyInSummaryLimitShare ----------------------------------------------

func barisRingkas(layer string, limit, ded, net, gross []models.NilaiMataUang) models.BarisShareNP {
	return models.BarisShareNP{
		LayerType: "Layer ", Layer: layer, LayerPartType: "Part ", LayerPart: "1",
		RnmLimitList: limit, NetPremiumList: net, GrossPremiumList: gross,
		DeductionList: func() []models.BarisDeduksiShare {
			out := []models.BarisDeduksiShare{}
			for _, d := range ded {
				out = append(out, models.BarisDeduksiShare{Currency: d.Currency, Deduction: d.Value})
			}
			return out
		}(),
	}
}

// ⛔ Mata uangnya DIKODE KERAS: IDR ke kolom 1, USD ke kolom 2. Yang lain
// hilang tanpa jejak — itu bunyi `@If(local.currency=="IDR",…)` ekspor.
func TestRingkasanShareHanyaIDRdanUSD(t *testing.T) {
	rows := []models.BarisShareNP{barisRingkas("1",
		[]models.NilaiMataUang{nilai("IDR", "10"), nilai("USD", "2"), nilai("SGD", "7")},
		nil, nil, nil)}
	out := services.SummaryLimitShare(rows)
	if len(out) != 1 {
		t.Fatalf("baris ringkasan = %d, mau 1", len(out))
	}
	if out[0].Limit != "10" || out[0].Limit2 != "2" {
		t.Fatalf("Limit/Limit2 = %q/%q, mau 10/2", out[0].Limit, out[0].Limit2)
	}
}

// ⛔ Baris yang 100% Limit-nya KOSONG nol masuk ringkasan — hanya langkah
// 4.3.3 yang menambah baris, dan ia di dalam loop `RnmLimitList`.
func TestBarisTanpaLimitTidakMasukRingkasan(t *testing.T) {
	rows := []models.BarisShareNP{barisRingkas("1", nil, nil,
		[]models.NilaiMataUang{nilai("IDR", "50")}, nil)}
	if out := services.SummaryLimitShare(rows); len(out) != 0 {
		t.Fatalf("baris ringkasan = %+v, mau nol", out)
	}
}

func TestRingkasanDigabungPerLayer(t *testing.T) {
	rows := []models.BarisShareNP{
		barisRingkas("1", []models.NilaiMataUang{nilai("IDR", "10")}, nil, nil, nil),
		barisRingkas("1", []models.NilaiMataUang{nilai("IDR", "5")}, nil, nil, nil),
		barisRingkas("2", []models.NilaiMataUang{nilai("IDR", "3")}, nil, nil, nil),
	}
	out := services.SummaryLimitShare(rows)
	if len(out) != 2 {
		t.Fatalf("baris ringkasan = %d, mau 2 (layer 1 dan 2)", len(out))
	}
	if out[0].Limit != "15" {
		t.Fatalf("layer 1 Limit = %q, mau 15", out[0].Limit)
	}
}

// ⭐ `Note` disusun ekspor langkah 6: LayerType + Layer + " of " +
// LayerPartType + LayerPart.
func TestNoteRingkasanDisusunDariEmpatMedanLayer(t *testing.T) {
	rows := []models.BarisShareNP{barisRingkas("1", []models.NilaiMataUang{nilai("IDR", "1")}, nil, nil, nil)}
	out := services.SummaryLimitShare(rows)
	if out[0].Note != "Layer 1 of Part 1" {
		t.Fatalf("Note = %q", out[0].Note)
	}
}

// --- TreatyInNonAddItem(share) ----------------------------------------------

// ⛔ `RNMShare == 0` menghentikan `Update Summary` dengan pesan ekspor, dan
// kedua larik tetap kosong.
func TestUpdateSummaryMenolakRNMShareNol(t *testing.T) {
	s := shareKosong()
	s.RNMShare = "0"
	pesan := services.NonAddItemShare(&s, []services.LayerNP{{Layer: "1"}})
	if len(pesan) != 1 || pesan[0] != services.PesanShareNilaiKosong {
		t.Fatalf("pesan = %v, mau %q", pesan, services.PesanShareNilaiKosong)
	}
	if len(s.Share) != 0 {
		t.Fatalf("baris Share lahir walau RNMShare nol: %d", len(s.Share))
	}
}

// ⭐ `Share to Other Retro` > 0 memotong RNM Share dan melahirkan larik
// facultative — `RnmShareDeducted` yang tampil di `Share to RNM :`.
func TestShareKeRetroMemotongRnmShare(t *testing.T) {
	s := shareKosong()
	s.RNMShare = "30"
	s.FacultativeShare = "10"
	layers := []services.LayerNP{{Layer: "1", LayerType: "Layer "}}
	if pesan := services.NonAddItemShare(&s, layers); len(pesan) != 0 {
		t.Fatalf("pesan tak diduga: %v", pesan)
	}
	if s.RnmShareDeducted != "20" {
		t.Fatalf("RnmShareDeducted = %q, mau 20", s.RnmShareDeducted)
	}
	if len(s.FacultativeShareList) != 1 || len(s.Share) != 1 {
		t.Fatalf("larik = %d fac / %d share, mau 1/1", len(s.FacultativeShareList), len(s.Share))
	}
}

// ⚠️ Nol `Share to Other Retro` berarti nol baris facultative — dan
// `RnmShareDeducted` TIDAK disentuh.
func TestTanpaShareKeRetroNolBarisFacultative(t *testing.T) {
	s := shareKosong()
	s.RNMShare = "30"
	services.NonAddItemShare(&s, []services.LayerNP{{Layer: "1"}})
	if len(s.FacultativeShareList) != 0 {
		t.Fatalf("baris facultative lahir tanpa Share to Other Retro")
	}
}

// --- rantai HitungShareNP ---------------------------------------------------

// ⭐ `set-brokerage` = `TreatyInSetBrokerage` SAJA, dan `brokerage` =
// SetBrokerage LALU XOLAddSpreading — `set-brokerage` disusul `rnm`
// harus sampai di keadaan yang SAMA dengan `brokerage`.
func TestSetBrokerageSajaLaluRnmSamaDenganBrokerage(t *testing.T) {
	awal := shareKosong()
	awal.RNMShare = "50"
	awal.BrokeragePercent = "10"
	layers := []services.LayerNP{{Layer: "1", MDPList: []models.NilaiMataUang{{Currency: "IDR", Value: "1000"}}}}
	// Gross Premium baris ADA — deduksi < 1 dibuang langkah [6].
	awal.Share = []models.BarisShareNP{{Layer: "1", RNMShare: "50", GrossPremiumList: []models.NilaiMataUang{{Currency: "IDR", Value: "1000"}}}}
	sekali := services.HitungShareNP(services.MasukanShareNP{Aksi: services.AksiShareBrokerage, Share: awal, Layers: layers}, sumberKosong{})
	dua := services.HitungShareNP(services.MasukanShareNP{Aksi: services.AksiShareSetBrokerage, Share: awal, Layers: layers}, sumberKosong{})
	if len(dua.Share.Share) != 1 || len(dua.Share.Share[0].DeductionList) == 0 || dua.Share.Share[0].DeductionList[0].DeductionPct != "10" {
		t.Fatalf("set-brokerage tidak menulis Brokerage fee: %+v", dua.Share.Share)
	}
	dua = services.HitungShareNP(services.MasukanShareNP{Aksi: services.AksiShareRNM, Share: dua.Share, Layers: layers}, sumberKosong{})
	if cari(sekali.Share.Total["TotalShareNetNP"], "IDR") != cari(dua.Share.Total["TotalShareNetNP"], "IDR") ||
		len(sekali.Share.Share[0].DeductionList) != len(dua.Share.Share[0].DeductionList) {
		t.Fatalf("set-brokerage+rnm berbeda dari brokerage:\n%+v\n%+v", sekali.Share.Share[0], dua.Share.Share[0])
	}
}

// ⭐ `Update Total` menjalankan NPSetTotal + kedua ringkasan sekaligus.
func TestAksiTotalMengisiRingkasanDanTotal(t *testing.T) {
	m := services.MasukanShareNP{
		Aksi: services.AksiShareTotal,
		Share: models.ShareNP{
			Total: map[string][]models.NilaiMataUang{},
			Share: []models.BarisShareNP{func() models.BarisShareNP {
				b := barisRingkas("1", []models.NilaiMataUang{nilai("IDR", "10")}, nil, nil, nil)
				// ⛔ RALAT 8 Oktober 2026: larik baris Share SENDIRI yang
				// dijumlah (tangkapan Pega 1001855), bukan tingkat spreading.
				b.GrossPremiumList = []models.NilaiMataUang{nilai("IDR", "4")}
				b.SpreadingListXOL = []models.BarisSpreadingNP{{
					GrossPremiumList: []models.NilaiMataUang{nilai("IDR", "999")},
				}}
				return b
			}()},
		},
	}
	h := services.HitungShareNP(m, sumberKosong{})
	if len(h.Share.LimitShareSummaryList) != 1 {
		t.Fatalf("ringkasan = %+v", h.Share.LimitShareSummaryList)
	}
	if got := cari(h.Share.Total["TotalShareGrossNP"], "IDR"); got != "4" {
		t.Fatalf("TotalShareGrossNP = %q, mau 4", got)
	}
}

// ⛔ Masukan tidak boleh berubah.
func TestHitungShareNPTidakMengubahMasukan(t *testing.T) {
	asal := models.ShareNP{
		Total: map[string][]models.NilaiMataUang{},
		Share: []models.BarisShareNP{{GrossPremiumList: []models.NilaiMataUang{nilai("IDR", "1")}}},
	}
	services.HitungShareNP(services.MasukanShareNP{Aksi: services.AksiShareTotal, Share: asal}, sumberKosong{})
	if len(asal.Total) != 0 {
		t.Fatalf("masukan berubah: %+v", asal.Total)
	}
}

// ⭐ Kesembilan kunci total SELALU ada — grid yang kuncinya hilang akan
// jatuh pada `.length` di layar.
func TestKesembilanKunciTotalSelaluAda(t *testing.T) {
	h := services.HitungShareNP(services.MasukanShareNP{
		Aksi:  services.AksiShareTotal,
		Share: models.ShareNP{Total: map[string][]models.NilaiMataUang{}},
	}, sumberKosong{})
	if len(services.KunciTotalShareNP) != 9 {
		t.Fatalf("kunci total = %d, mau 9", len(services.KunciTotalShareNP))
	}
	for _, k := range services.KunciTotalShareNP {
		if h.Share.Total[k] == nil {
			t.Fatalf("kunci %q nil — grid akan jatuh pada .length", k)
		}
	}
}

// ⭐ Nol larik nil di seluruh hasil.
func TestHasilShareNPNolLarikNil(t *testing.T) {
	h := services.HitungShareNP(services.MasukanShareNP{
		Aksi:  "entah",
		Share: models.ShareNP{Total: map[string][]models.NilaiMataUang{}},
	}, sumberKosong{})
	if h.Share.Share == nil || h.Share.FacultativeShareList == nil ||
		h.Share.ShareReins == nil || h.Share.ShareFacultativeReinsurers == nil ||
		h.Share.LimitShareSummaryList == nil || h.Pesan == nil {
		t.Fatalf("ada larik nil: %+v", h.Share)
	}
}
