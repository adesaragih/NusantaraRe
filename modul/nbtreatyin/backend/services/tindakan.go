package services

// Untuk apa berkas ini: TINDAKAN LAYAR - refresh berhitung, simpan draf,
// pilih bisnis, submit (finishAssignment), dan terbitkan nomor polis.
//
// Setiap tindakan bekerja atas halaman TERSIMPAN yang dipra-proses ulang
// (`siapkan`), lalu nilai kiriman layar DIGABUNG hanya untuk medan yang boleh
// diubah di posisi itu (`models.GabungMasukanLayar`; AC 49-52) - layar tidak
// pernah dapat menulis medan terkunci.

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/nbtreatyin/backend/models"
)

// kerjakan memuat kasus untuk ditindak: terbuka, pelaku anggota antreannya,
// halaman dipra-proses dan digabung dengan kiriman layar (hanya sel / wadah /
// grid yang terbuka - W2-W4 audit silang P3). Di layar admin, pilihan Source Of
// Business yang dipegang layar diterima sesudah penggabungan
// (`terimaSumberBisnis`, F4) dan sebelum medan turunan dihitung. Terakhir,
// action set sel yang dipicu kiriman ditulis ulang ke daftar
// (`models.TerapkanPemicu`, W5): kolom hanya-baca grid tidak pernah dari layar.
func (l *Layanan) kerjakan(ctx context.Context, p inti.Pelaku, id string, masuk *models.Halaman) (models.Kasus, *models.Halaman, error) {
	if err := l.periksaPelaku(p); err != nil {
		return models.Kasus{}, nil, err
	}
	k, h, err := l.muat(ctx, id)
	if err != nil {
		return models.Kasus{}, nil, err
	}
	if k.Tertutup() {
		return models.Kasus{}, nil, ErrKasusTertutup
	}
	if k.GenerasiTertutup {
		return models.Kasus{}, nil, ErrGenerasiTertutup
	}
	if !anggota(p, k.PositionNote) {
		return models.Kasus{}, nil, ErrBukanAnggotaAntrean
	}
	if err := l.siapkan(ctx, p, k, h); err != nil {
		return models.Kasus{}, nil, err
	}
	// `ListSuggest.ProductionDate` hanya diterima bila tampil bagi pelaku
	// (tempat berperan tiket 05; AC 52, 81).
	// Pemetaan kosong (K12, K16) = tempat tertunda = medan tidak diterima.
	pemicu, err := models.GabungMasukanLayar(h, masuk, k.PositionNote, tempatPelaku(p))
	if err != nil {
		return models.Kasus{}, nil, jawabKiriman(err)
	}
	if k.PositionNote == models.PosisiAdmin {
		// F4: pilihan Source Of Business yang dipegang layar (sumberbisnis.go).
		if err := l.terimaSumberBisnis(ctx, h, masuk); err != nil {
			return models.Kasus{}, nil, jawabKiriman(err)
		}
		if err := l.turunkan(ctx, h); err != nil {
			return models.Kasus{}, nil, err
		}
		if pemicu.Pajak {
			if err := l.pajakNonProp(ctx, h); err != nil {
				return models.Kasus{}, nil, err
			}
		}
	}
	// W5: kolom hanya-baca daftar dari server - action set sel yang dipicu
	// kiriman diputar ulang SESUDAH NetPremium / BalanceDueTo dihitung.
	if err := models.TerapkanPemicu(h, pemicu, l.jam()); err != nil {
		return models.Kasus{}, nil, err
	}
	return k, h, nil
}

// jawabKiriman - SATU jalur galat pola kiriman terkunci (F4): nilai yang bukan
// hasil tombol/popup yang dihitung ulang di server (`models.GalatKiriman` -
// Enable / Disable, Source Of Business, Choose popup bisnis) dijawab 422
// sebagai pesan validasi layar. Galat lain diteruskan apa adanya.
func jawabKiriman(err error) error {
	var gk *models.GalatKiriman
	if errors.As(err, &gk) {
		return &GalatValidasi{Pesan: []string{gk.Pesan}}
	}
	return err
}

// turunkan menghitung ULANG di server medan turunan layar admin dari
// isiannya - layar tidak pernah menjadi sumber nilai turunan (AC 49; temuan
// tinjauan 2026-10-03):
//
//	SetPPNPPH langkah 1-3   STS_PKP agen sumber bisnis
//	SetCurrency_act         Currency <- nama mata uang IDCurrency
//	CountNetPremi_act       NetPremium, Balance*, PPN/PPH, BrokerageFee*, DueTo
//	CountSpreading_Act 4.2-5  total spreading (jumlah baris)
//
// Baris daftar (PremiumSpreaded, ListInstallment) yang dipicu kiriman ditulis
// SESUDAHNYA oleh `models.TerapkanPemicu` (memakai NetPremium / BalanceDueTo ini).
//
// Ketiganya rumus yang sama yang dijalankan refresh layar, dan hasilnya
// bergantung pada isian saja - menjalankannya lagi tidak mengubah nilai yang
// layar sudah tampilkan.
func (l *Layanan) turunkan(ctx context.Context, h *models.Halaman) error {
	if err := l.pasangStsPKP(ctx, h); err != nil {
		return err
	}
	if id := h.Ambil(models.HalamanPolis + ".IDCurrency"); id != "" {
		nama, err := l.g.NamaMataUang(ctx, id)
		if err != nil {
			return err
		}
		models.SetelNamaMataUang(h, nama)
	}
	// ⛔ K8: polis NonProp baru - medan uang berada di kontainer
	// `.IsNewPolicyNonProp != 1` yang tersembunyi, tidak satu pun refresh memicu
	// CountNetPremi_act; nilainya milik InputPolicyTreatyInDetail_NonProp 18-19.
	if !models.PolisNonPropBaru(h) {
		if err := models.CountNetPremi(h); err != nil {
			return err
		}
	}
	return models.HitungTotalSpreading(h)
}

// tulis menjalankan fn di satu transaksi SESUDAH mengunci baris kasus dan
// memastikan tahapnya masih tahap yang dibaca `kerjakan`.
func (l *Layanan) tulis(ctx context.Context, k models.Kasus, fn func(tx *db.Tx) error) error {
	return l.g.Transaksi(ctx, func(tx *db.Tx) error {
		if err := l.g.KunciKasus(ctx, tx, k.ID, k.StatusWork); err != nil {
			return err
		}
		return fn(tx)
	})
}

// ------------------------------------------------------------------ hitung

// PermintaanHitung - action set satu sel layar: SATU bentuk, selalu `Urutan`
// (temuan tinjauan P9). Sel yang memuat lebih dari satu refresh berhitung
// (mis. `.RiCommOgp`: `CountResult1_Act(Data="Pct")` lalu `CountOGPONP_Act`)
// menjalankannya berurutan atas halaman yang sama, seperti clipboard Pega;
// sel satu refresh = urutan satu langkah.
type PermintaanHitung struct {
	Urutan []LangkahHitung `json:"urutan"`
	// Indeks - Param.Index / Param.idx, berbasis 1.
	Indeks  int             `json:"indeks"`
	Halaman *models.Halaman `json:"halaman"`
}

// LangkahHitung - satu refresh berhitung.
type LangkahHitung struct {
	// Aksi - nama aktivitas Pega tanpa akhiran `_Act` (lihat `aksiHitung`).
	Aksi string `json:"aksi"`
	// Param - parameter aktivitas: Data ("Pct"/"Amount"), DiscountType,
	// Result, Overidding, Action ("PREMIUM"/"CLAIM").
	Param string `json:"param"`
}

// aksiHalaman, aksiParam, aksiIndeks - tiga bentuk tanda tangan aktivitas Pega
// yang dipanggil refresh layar: tanpa parameter, dengan Param.<nama> teks,
// dan dengan Param.Index/idx.
type aksiFn = func(l *Layanan, ctx context.Context, h *models.Halaman, param string, indeks int) error

func aksiHalaman(f func(*models.Halaman) error) aksiFn {
	return func(_ *Layanan, _ context.Context, h *models.Halaman, _ string, _ int) error { return f(h) }
}

func aksiParam(f func(*models.Halaman, string) error) aksiFn {
	return func(_ *Layanan, _ context.Context, h *models.Halaman, param string, _ int) error { return f(h, param) }
}

func aksiIndeks(f func(*models.Halaman, int) error) aksiFn {
	return func(_ *Layanan, _ context.Context, h *models.Halaman, _ string, indeks int) error {
		return f(h, indeks)
	}
}

func tanpaGalat(f func(*models.Halaman)) func(*models.Halaman) error {
	return func(h *models.Halaman) error { f(h); return nil }
}

// aksiHitung - aksi refresh yang tertulis di sel layar (INVENTARIS bab 5).
var aksiHitung = map[string]aksiFn{
	"CountOGPONP":              aksiHalaman(models.CountOGPONP),
	"CountResult1":             aksiParam(models.CountResult1),
	"CountResult1Onp":          aksiParam(models.CountResult1Onp),
	"CountResult2Ogp":          aksiParam(models.CountResult2Ogp),
	"CountResult2Onp":          aksiParam(models.CountResult2Onp),
	"CountRiCommOgp":           aksiParam(models.CountRiCommOgp),
	"CountRiCommOnp":           aksiParam(models.CountRiCommOnp),
	"CountOverridingCommOgp":   aksiParam(models.CountOverridingCommOgp),
	"CountOverridingCommOnp":   aksiParam(models.CountOverridingCommOnp),
	"CountNetPremi":            aksiHalaman(models.CountNetPremi),
	"CalculatePremi":           aksiParam(models.CalculatePremi),
	"SetValidateInstallment":   aksiHalaman(models.SetValidateInstallment),
	"CountPctInstallment":      aksiIndeks(models.CountPctInstallment),
	"CountSpreading":           aksiIndeks(models.CountSpreading),
	"SetDueTo":                 aksiHalaman(models.SetDueTo),
	"ProtectDate":              aksiHalaman(tanpaGalat(models.ProtectDate)),
	"SystemSetOneYear":         aksiHalaman(tanpaGalat(models.SystemSetOneYear)),
	"RemoveTypeTax":            aksiHalaman(tanpaGalat(models.RemoveTypeTax)),
	"TreatyEnableDisableInput": aksiHalaman(tanpaGalat(models.TreatyEnableDisableInput)),
	"FillPaymentInstallment": func(l *Layanan, _ context.Context, h *models.Halaman, _ string, _ int) error {
		return models.FillPaymentInstallment(h, l.jam())
	},
	// SetCurrency_act(CURR=.IDCurrency): RDB GetCurrency -> .Currency
	"SetCurrency": func(l *Layanan, ctx context.Context, h *models.Halaman, _ string, _ int) error {
		nama, err := l.g.NamaMataUang(ctx, h.Ambil(models.HalamanPolis+".IDCurrency"))
		if err != nil {
			return err
		}
		models.SetelNamaMataUang(h, nama)
		return nil
	},
	// CheckDataMkt: Obj-Browse marketing officer (dilewati bila MOID kosong)
	"CheckDataMkt": func(l *Layanan, ctx context.Context, h *models.Halaman, _ string, _ int) error {
		var mo models.BarisMO
		if id := h.Ambil(models.HalamanPolis + ".QuotationData.MOID"); id != "" {
			var err error
			if mo, err = l.g.MO(ctx, id); err != nil {
				return err
			}
		}
		models.TerapkanMO(h, mo)
		return nil
	},
	// sel `.TypeTax` - `[keputusan work owner 06-10-2026]` pajak dihitung ulang saat diubah (XML hanya
	// postValue). Hitungannya dijalankan `kerjakan` (pemicu Pajak); aksi ini hanya putaran ke server.
	"HitungPajak": aksiHalaman(func(*models.Halaman) error { return nil }),
}

// Hitung menjalankan action set sel (`Urutan`) dan mengembalikan layarnya.
// Hasil TIDAK disimpan - kecuali `CheckDataMkt`, yang di XML menutup dirinya
// dengan `Obj-Save pyWorkPage` (langkah 5).
func (l *Layanan) Hitung(ctx context.Context, p inti.Pelaku, id string, r PermintaanHitung) (Layar, error) {
	if len(r.Urutan) == 0 {
		return Layar{}, fmt.Errorf("%w: urutan hitung kosong", ErrPermintaanTidakSah)
	}
	fs := make([]aksiFn, len(r.Urutan))
	simpan := false
	for i, s := range r.Urutan {
		f, ada := aksiHitung[s.Aksi]
		if !ada {
			return Layar{}, fmt.Errorf("%w: aksi hitung %q", ErrPermintaanTidakSah, s.Aksi)
		}
		fs[i] = f
		simpan = simpan || s.Aksi == "CheckDataMkt"
	}
	k, h, err := l.kerjakan(ctx, p, id, r.Halaman)
	if err != nil {
		return Layar{}, err
	}
	for _, s := range r.Urutan { // aksiposisi.go
		if !aksiTerbuka(k.PositionNote, s.Aksi, h) {
			return Layar{}, fmt.Errorf("%w: aksi %q", ErrTindakanTakAdaDiPosisi, s.Aksi)
		}
	}
	if err := l.pasangStsPKP(ctx, h); err != nil {
		return Layar{}, err
	}
	for i, s := range r.Urutan {
		if err := fs[i](l, ctx, h, s.Param, r.Indeks); err != nil {
			return Layar{}, err
		}
	}
	if simpan {
		if err := l.tulis(ctx, k, func(tx *db.Tx) error { return l.g.SimpanHalaman(ctx, tx, id, h) }); err != nil {
			return Layar{}, err
		}
	}
	return l.layar(ctx, p, k, h, true)
}

// pasangStsPKP = `SetPPNPPH` langkah 1-3: RD `BrowseClientName_RD` atas
// `QuotationData.SourceOfBusiness`; hasilnya dibaca syarat langkah 4.
func (l *Layanan) pasangStsPKP(ctx context.Context, h *models.Halaman) error {
	sts, err := l.g.StsPKPAgen(ctx, h.Ambil(models.HalamanPolis+".QuotationData.SourceOfBusiness"))
	if err != nil {
		return err
	}
	h.Setel(models.JalurStsPKP, sts)
	return nil
}

// ------------------------------------------------------------------ simpan dan pilih bisnis

// SimpanDraf = tombol "Save" layar admin (`click→save`). Layar atasan tidak
// punya tombol Save.
func (l *Layanan) SimpanDraf(ctx context.Context, p inti.Pelaku, id string, masuk *models.Halaman) (Layar, error) {
	k, h, err := l.kerjakan(ctx, p, id, masuk)
	if err != nil {
		return Layar{}, err
	}
	if k.PositionNote != models.PosisiAdmin {
		return Layar{}, ErrTindakanTakAdaDiPosisi
	}
	// `[keputusan work owner]` AC 48: "Berkas tidak dapat disimpan bila medan
	// wajib pada tingkat itu kosong" - juga tombol Save, bukan hanya Submit. Medan UANG wajib yang kosong
	// diisi 0 lebih dulu (permintaan work owner 06-10-2026).
	models.IsiNolWajibUang(h, k.PositionNote, tempatPelaku(p))
	// Approval dan Suggest wajib saat Submit saja, bukan halangan Save (permintaan work owner 06-10-2026).
	if kosong := models.MedanWajibKosongSimpan(h, k.PositionNote, tempatPelaku(p)); len(kosong) > 0 {
		return Layar{}, &GalatValidasi{Pesan: kosong}
	}
	if err := l.tulis(ctx, k, func(tx *db.Tx) error { return l.g.SimpanHalaman(ctx, tx, id, h) }); err != nil {
		return Layar{}, err
	}
	return l.layar(ctx, p, k, h, true)
}

// PilihBisnis = tombol "Choose" popup `BusinessAndSOBList` -> `SetValue_Act`
// (ID=.ID) -> `InputPolicyTreatyInDetail_preACT`, lalu `Obj-Save`. Bagian
// yang dibangun: lihat `models/pilihbisnis.go` dan `models/komisi.go`
// (langkah 17).
//
// ⛔ Pembacaan kontrak yang gagal MENGHENTIKAN proses dan tidak menyimpan
// apa pun (AC 36-38) - di Pega `pxRetrieveReportData` yang kosong mengisi
// halaman dengan nilai kosong lalu menyimpannya.
func (l *Layanan) PilihBisnis(ctx context.Context, p inti.Pelaku, id, idDetail string, masuk *models.Halaman) (Layar, error) {
	if strings.TrimSpace(idDetail) == "" {
		return Layar{}, fmt.Errorf("%w: ID kontrak kosong", ErrPermintaanTidakSah)
	}
	k, h, err := l.kerjakan(ctx, p, id, masuk)
	if err != nil {
		return Layar{}, err
	}
	// `Choose` hidup hanya di dalam popup tombol `Choose Business` (layanan.go).
	if err := bolehPilihBisnis(k, h); err != nil {
		return Layar{}, err
	}
	// Pola F4: `Choose` hanya ada di baris grid popup - RD `BrowseTreatyJoinEDM`
	// tersaring kasus ini (`DaftarBisnis`) dijalankan ulang; ID di luarnya 422.
	// ID dicari langsung (saringan kolom popup bekerja di server, keputusan WO 06-10-2026): baris yang sah =
	// baris view yang dapat dikembalikan RD dengan filter H kasus ini, bukan hanya 500 baris pertama.
	s := models.SaringanPopupBisnis(h)
	s.ID = idDetail
	daftar, err := l.g.DaftarBisnis(ctx, s)
	if err != nil {
		return Layar{}, err
	}
	if err := models.PeriksaPilihanBisnis(daftar, idDetail); err != nil {
		return Layar{}, jawabKiriman(err)
	}
	b, err := l.g.DetailKontrak(ctx, idDetail)
	if err != nil {
		return Layar{}, err
	}
	models.TerapkanDetailKontrak(h, b) // langkah 3, 5, 6, 11
	models.TerapkanMasterKontrak(h, b)
	if h.Ambil(models.HalamanPolis+".IDCurrency") == "" { // langkah 4
		idCur, err := l.g.IDMataUangDariNama(ctx, h.Ambil(models.HalamanPolis+".Currency"))
		if err != nil {
			return Layar{}, err
		}
		models.SetelIDMataUang(h, idCur)
	}
	ojk, err := l.g.OJKGrupTreaty(ctx, h.Ambil(models.HalamanPolis+".TreatyGroupID")) // langkah 8
	if err != nil {
		return Layar{}, err
	}
	models.SetelOJK(h, ojk)
	klien, err := l.g.KlienDariNama(ctx, h.Ambil(models.HalamanQuotation+".InsuredName")) // 14.1-14.3
	if err != nil {
		return Layar{}, err
	}
	models.TerapkanKlien(h, klien)
	bis, err := l.g.BisnisDariKunci(ctx, models.KunciCariBisnis(h.Ambil(models.HalamanQuotation+".BusinessName"), false)) // 14.4-14.6
	if err != nil {
		return Layar{}, err
	}
	models.TerapkanBisnisPilih(h, bis) // 14.7-14.9
	// 16 (NonProportional) dan 18 - K8, nonprop.go (16 dan 17 saling meniadakan;
	// 18 hanya berbuat bila master XOL termuat).
	if err := l.pilihBisnisNonProp(ctx, h); err != nil {
		return Layar{}, err
	}
	// langkah 17 (bukan NonProportional) -> TreatyInputPctCommSpreading:
	// RiCommOgp dari baris view kontrak NoOffer (models/komisi.go).
	if models.LangkahKomisiProporsional(h) {
		baris, err := l.g.KomisiKontrak(ctx, h.Ambil(models.HalamanPolis+".NoOffer"))
		if err != nil {
			return Layar{}, err
		}
		models.TreatyInputPctCommSpreading(h, baris)
	}
	if err := l.tulis(ctx, k, func(tx *db.Tx) error { return l.g.SimpanHalaman(ctx, tx, id, h) }); err != nil {
		return Layar{}, err
	}
	return l.layar(ctx, p, k, h, true)
}

// ------------------------------------------------------------------ kirim

// Kirim = `finishAssignment` flow action posisi kasus.
//
//	Admin  `InboxPolicyTreatyIn`: tombol Submit (IsApproved 1: SetDueTo_act
//	       lalu finish; 0: `PolicyTreatyInDeclineConfirm`, local action
//	       SetDueTo_act, Yes = finish) -> pasca DT `InboxPolicyTreatyIn_postDT`
//	       -> `InputPolicyTreatyInPost_Act` (riwayat, cek duplikat bila
//	       IsApproved 1, SaveViewSuggest)
//	Atasan `DeptHeadTreatyIn_UW`: pasca DT `DeptHeadTreatyIn_UW_postDT` ->
//	       `InsertHistoryAkseptasiPega`
//
// `SaveViewSuggest` (InputPolicyTreatyInPost_Act langkah 4) menulis catatan
// yang baru ditambahkan pasca DT ke `POOLDATA.HISTORYAKSEPTASIPRODUCTION`
// (`models.UsulanBelumTersimpan`, K4) - tanpa syarat `BusinessFac == "F"`
// (`[penyimpangan sadar]` K4) dan di KETIGA jenjang, TGL_INP 24 jam, NOURUT
// repository (`[penyimpangan sadar — disetujui WO 04-10-2026]`); rinciannya di
// models/usulan.go.
//
// lalu connector flow (`models.Langkah`). Bila realisasi selesai (Decision8 Else),
// Utility1 `SaveJsonPolisTreatyIn_Act` ditulis di transaksi yang sama: medan
// halaman (`models.PrasimpanPolis`), lalu json_polis TANPA DATA_JSON, ACHIEVEMENT,
// TREATYINPRODUCTION (`[keputusan work owner 06-10-2026]` "JSON-nya tidak
// disimpan, tapi tetap insert kolom lainnya"; models/produksi.go). Semuanya SATU
// transaksi (AC 29, 83). Sesudah transaksi: Utility2 `serviceInsertArasapas_act`
// (`konversikan`, KEPUTUSAN-RONDE-12 butir 7) - gagalnya tidak membatalkan apa pun.
func (l *Layanan) Kirim(ctx context.Context, p inti.Pelaku, id string, masuk *models.Halaman) (HasilKirim, error) {
	k, err := l.kirim(ctx, p, id, masuk)
	if err != nil {
		return HasilKirim{}, err
	}
	hasil := HasilKirim{Kasus: k}
	if k.StatusWork == models.StatusSelesai {
		h, err := l.g.BacaHalaman(ctx, nil, id)
		if err != nil {
			// Kasus sudah selesai dan tersimpan - submit TIDAK digagalkan; konversi
			// yang tidak dapat disusun dicatat dan ditandai di layar.
			log.Printf("nbtreatyin: konversi %s tidak disusun: %v", id, err)
			hasil.PesanKonversi = models.PesanGagalKonversi
			return hasil, nil
		}
		hasil.PesanKonversi = l.konversikan(ctx, id, h)
	}
	return hasil, nil
}

func (l *Layanan) kirim(ctx context.Context, p inti.Pelaku, id string, masuk *models.Halaman) (models.Kasus, error) {
	k, h, err := l.kerjakan(ctx, p, id, masuk)
	if err != nil {
		return models.Kasus{}, err
	}
	posisi := k.PositionNote
	if models.TombolUntuk(h, posisi) == models.TombolTidakAda {
		return models.Kasus{}, &GalatValidasi{Pesan: []string{"Approval"}}
	}
	if err := models.SetDueTo(h); err != nil {
		return models.Kasus{}, err
	}
	if err := l.pasangStsPKP(ctx, h); err != nil {
		return models.Kasus{}, err
	}
	if err := l.validasiKirim(ctx, p, h, posisi); err != nil {
		return models.Kasus{}, err
	}
	nama, err := l.g.NamaTampilan(ctx, p.AkunID)
	if err != nil {
		return models.Kasus{}, err
	}
	// pasca DT
	if posisi == models.PosisiAdmin {
		models.PascaAdmin(h, p.AkunID, nama)
	} else {
		pembuat, err := l.g.NamaTampilan(ctx, k.CreateOp)
		if err != nil {
			return models.Kasus{}, err
		}
		models.PascaAtasan(h, p.AkunID, nama, pembuat)
	}
	// InputPolicyTreatyInPost_Act langkah 3 (admin, IsApproved == 1)
	if posisi == models.PosisiAdmin && h.Ambil(models.HalamanPolis+".IsApproved") == "1" {
		if err := l.cekDuplikat(ctx, h); err != nil {
			return models.Kasus{}, err
		}
	}
	tr, err := models.Langkah(posisi, h.Ambil(models.HalamanPolis+".IsApproved"), h.Ambil(models.HalamanPolis+".PolicyNo") != "")
	if err != nil {
		return models.Kasus{}, err
	}
	// NBStatus "NB IS IN <nama>'S INBOX": nama dari pemegang aktif workbasket tujuan (keputusan WO 06-10-2026)
	// - penolakan atasan: nama pembuat berkas; berkas lama tanpa pembuat -> pemegang workbasket Admin
	namaKotak := ""
	if tr.KembaliKePembuat {
		if namaKotak, err = l.g.NamaTampilan(ctx, k.CreateOp); err != nil {
			return models.Kasus{}, err
		}
	}
	if tr.NBStatusKePosisi || (tr.KembaliKePembuat && strings.TrimSpace(namaKotak) == "") {
		pk, err := l.g.PemegangKotakMasuk(ctx, tr.PosisiBaru)
		if err != nil {
			return models.Kasus{}, err
		}
		namaKotak = models.NamaKotakMasuk(tr.PosisiBaru, pk)
	}
	err = l.tulis(ctx, k, func(tx *db.Tx) error {
		// InsertHistoryAkseptasiPega - sebelum connector: WORKBASKET = posisi
		// tempat putusan diambil.
		if err := l.g.CatatRiwayat(ctx, tx, models.Riwayat{
			IDPega:     models.KunciInstans(id),
			Status:     models.StatusRiwayat(h.Ambil(models.HalamanPolis + ".IsApproved")),
			Username:   nama,
			Workbasket: posisi,
			OperatorID: p.AkunID,
		}); err != nil {
			return err
		}
		// SaveViewSuggest (K4): catatan baru -> HISTORYAKSEPTASIPRODUCTION.
		if baru := models.UsulanBelumTersimpan(h); len(baru) > 0 {
			if err := l.g.CatatUsulan(ctx, tx, models.KunciInstans(id), baru); err != nil {
				return err
			}
		}
		if tr.Ditutup() {
			// Utility1 SaveJsonPolisTreatyIn_Act - hanya jalur selesai (Decision8 Else), bukan penolakan admin
			var simpanan models.SimpananPolis
			if tr.Simpan {
				hari, err := l.g.HariClosing(ctx, tx)
				if err != nil {
					return err
				}
				models.PrasimpanPolis(h, id, l.jam(), hari)
				if simpanan, err = models.SusunSimpananPolis(h, id, p.AkunID); err != nil {
					return err
				}
			}
			if err := l.g.SimpanHalaman(ctx, tx, id, h); err != nil {
				return err
			}
			if tr.Simpan {
				if err := l.g.SimpanPolisProduksi(ctx, tx, simpanan); err != nil {
					return err
				}
			}
			return l.g.TutupKasus(ctx, tx, id, k.StatusWork, tr.StatusTutup)
		}
		// NBStatus tidak pernah kosong selama berkas berjalan (keputusan WO 06-10-2026)
		if tr.NBStatusKePosisi || tr.KembaliKePembuat {
			h.Setel("NBStatus", models.TeksNBStatusKotakMasuk(namaKotak))
		}
		h.Setel("PositionNote", tr.PosisiBaru)
		if err := l.g.SimpanHalaman(ctx, tx, id, h); err != nil {
			return err
		}
		return l.g.PindahPosisi(ctx, tx, id, k.StatusWork, tr.PosisiBaru)
	})
	if err != nil {
		return models.Kasus{}, err
	}
	return l.g.Keadaan(ctx, nil, id)
}

// validasiKirim - medan wajib layar (AC 45-48) dan pesan validasi yang
// dipasang rantai layar (Property-Set-Messages / Page-Set-Messages): halaman
// berpesan tidak dapat di-submit, sama dengan Pega.
func (l *Layanan) validasiKirim(ctx context.Context, p inti.Pelaku, h *models.Halaman, posisi string) error {
	// medan UANG wajib yang kosong diisi 0 lebih dulu (permintaan work owner 06-10-2026); tersimpan bersama halaman
	models.IsiNolWajibUang(h, posisi, tempatPelaku(p))
	pesan := models.MedanWajibKosong(h, posisi, tempatPelaku(p))
	// 7.4 audit silang P3: tombol Submit Dept Head (IsApproved 1) menjalankan
	// GeneratePolicyNoTreaty_Act lebih dulu (TerbitkanNomor); langkah 11 ->
	// FetchTreatyGroupOldID langkah 3 memasang pesan halaman bila TreatyGroupID
	// kosong - halaman berpesan tidak dapat di-submit (OK = finishAssignment).
	if posisi == models.PosisiDeptHead && models.TombolUntuk(h, posisi) == models.TombolNomorPolis &&
		models.GrupTreatyTakTerbaca(h) {
		pesan = append(pesan, models.PesanGrupTreatyTakTerbaca)
	}
	if posisi == models.PosisiAdmin {
		// Rantai layar admin dijalankan atas SALINAN: hanya pesannya yang
		// diambil, nilai yang akan disimpan tidak berubah.
		salin := h.Salin()
		salin.BersihkanPesan()
		// K8: rantai uang (CountOGPONP_Act) hanya terpicu dari medan kontainer
		// proporsional - tersembunyi bagi polis NonProp baru.
		if !models.PolisNonPropBaru(h) {
			if err := models.CountOGPONP(salin); err != nil {
				return err
			}
		}
		models.ProtectDate(salin)
		pesan = append(pesan, salin.SemuaPesan()...)
	}
	if len(pesan) > 0 {
		return &GalatValidasi{Pesan: pesan}
	}
	return nil
}

// PesanDuplikatAwal - VERBATIM `TreatyRealizationCheckDuplicate` langkah 4.
const PesanDuplikatAwal = "Protect Duplicate Policy; data is similar to "

// cekDuplikat = `Activity/TreatyRealizationCheckDuplicate`, dipanggil HANYA
// `InputPolicyTreatyInPost_Act` langkah 3 (pasca-proses flow action admin)
// bersyarat `.PolicyTreatyIn.IsApproved==1` (AC 59, RALAT putaran 2):
//
//	1    Page-Clear-Messages pyWorkPage
//	2-3  RDB `TreatyRealizationCheckDuplicate` atas TREATYINPRODUCTION
//	     (`repository.PolisSerupa`, pemetaan parameter apa adanya)
//	4    local.msg = "Protect Duplicate Policy; data is similar to "
//	5.1  setiap baris: local.msg = local.msg + .CARI1 + " "
//	6    `@SizeOfPropertyList(ResultData.pxResults) > 0 && .PolicyTreatyIn.ClaimType != "XOL"`
//	     -> Page-Set-Messages pyWorkPage: submit tertahan (422)
//
// ⛔ `CheckDuplicateOffer` (peringatan "jumlah klaim > 0") TIDAK dibangun:
// langkah 1-4 berlabel `//`, dan satu-satunya pemanggilnya (`SetTreatyIn_Act`
// langkah 14) memanggilnya tanpa parameter (`pyPassCurrentParameterPage=false`)
// sehingga `GetCountClaim` selalu menghitung `masterid = NULL` = 0 - langkah
// 7-8 (pesannya) tidak pernah benar.
func (l *Layanan) cekDuplikat(ctx context.Context, h *models.Halaman) error {
	sama, err := l.g.PolisSerupa(ctx, h)
	if err != nil {
		return err
	}
	if len(sama) == 0 || h.Ambil(models.HalamanPolis+".ClaimType") == "XOL" {
		return nil
	}
	msg := PesanDuplikatAwal
	for _, n := range sama {
		msg += n + " "
	}
	return &GalatValidasi{Pesan: []string{msg}}
}

// ------------------------------------------------------------------ nomor polis

// NomorPolis adalah hasil `ShowPolicyNoTreaty_SC` (pyID dan PolicyNo).
type NomorPolis struct {
	ID       string `json:"id"`
	PolicyNo string `json:"policyNo"`
}

// TerbitkanNomor = tombol Submit Dept Head (IsApproved 1):
// `GeneratePolicyNoTreaty_Act` -> `ShowPolicyNoTreaty` (OK = Kirim).
//
// ⛔ Nomor dibentuk SEKALI per berkas (AC 74; langkah 28 bersyarat
// `PolicyNo==""`): berkas yang sudah bernomor mengembalikan nomornya.
func (l *Layanan) TerbitkanNomor(ctx context.Context, p inti.Pelaku, id string, masuk *models.Halaman) (NomorPolis, error) {
	k, h, err := l.kerjakan(ctx, p, id, masuk)
	if err != nil {
		return NomorPolis{}, err
	}
	if models.TombolUntuk(h, k.PositionNote) != models.TombolNomorPolis {
		return NomorPolis{}, ErrTindakanTakAdaDiPosisi
	}
	if nopol := h.Ambil(models.HalamanPolis + ".PolicyNo"); nopol != "" {
		return NomorPolis{ID: id, PolicyNo: nopol}, nil
	}
	// langkah 11, 13 - hanya bila kosong
	if h.Ambil(models.HalamanPolis+".TreatyGroupOldID") == "" {
		old, err := l.g.OldIDGrupTreaty(ctx, h.Ambil(models.HalamanPolis+".TreatyGroupID"))
		if err != nil {
			return NomorPolis{}, err
		}
		models.SetelGrupLama(h, old)
	}
	if h.Ambil(models.HalamanPolis+".OJKBusinessID") == "" {
		ojk, err := l.g.OJKGrupTreaty(ctx, h.Ambil(models.HalamanPolis+".TreatyGroupID"))
		if err != nil {
			return NomorPolis{}, err
		}
		models.SetelOJK(h, ojk)
	}
	models.SalinQuotation(h) // langkah 10
	var hasil NomorPolis
	err = l.tulis(ctx, k, func(tx *db.Tx) error {
		bahan, err := l.g.TerbitkanNomorPolis(ctx, tx, h, l.jam())
		if err != nil {
			return err
		}
		if err := l.g.SetelNomorPolis(ctx, tx, id, bahan.NoPolis); err != nil {
			return err
		}
		h.Setel(models.HalamanPolis+".PolicyNo", bahan.NoPolis)
		h.Setel(models.HalamanPolis+".ProductionDate", utils.FormatTanggalWaktu(bahan.ProductionDate))
		hasil = NomorPolis{ID: id, PolicyNo: bahan.NoPolis}
		return l.g.SimpanHalaman(ctx, tx, id, h) // langkah 30 Obj-Save
	})
	if err != nil {
		return NomorPolis{}, err
	}
	return hasil, nil
}
