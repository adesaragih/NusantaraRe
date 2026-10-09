package services

// Baca satu kontrak WARISAN untuk form — keputusan pemilik proses.
//
// ⛔ Terjemahan tampilnya MEMAKAI ULANG `TanggalTampil` dan
// `SifatProporsiTampil` yang sudah hidup di `warisan_daftar.go`. Dua
// penerjemah untuk satu aturan adalah dua tempat untuk salah, dan yang kedua
// selalu yang basi.

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatyin/backend/models"
	"nusantarare/modul/treatyin/backend/repository"
)

var (
	// ErrWarisanTidakAda - pengenalnya tidak menunjuk kontrak mana pun.
	ErrWarisanTidakAda = repository.ErrWarisanTidakAda
	// ErrJSONWarisanRusak - dokumennya ada tetapi tidak dapat diurai.
	ErrJSONWarisanRusak = repository.ErrJSONWarisanRusak
)

// KunciPenampungRevisi - properti `T_TREATY_REVISION` (migrasi `448`) yang
// layar pegang hanya di penampung halaman; lihat `KontrakWarisan.Penampung`.
var KunciPenampungRevisi = []string{
	"RNMShareP", "BrokeragePercentP", "OptionLimit", "InstallmentNo", "RevisionState", "ViewState",
	// ⭐ 8 Oktober 2026 — kedua total skalar tab EGNPI (`TabEgnpi`): kolomnya
	// sudah ada di `T_TREATY_REVISION` dan Save menulisnya, tetapi tab
	// menyemai kosong saat kontrak dibuka.
	"TotalEgnpiAmount", "TotalEgnpiProportion",
	// ⭐ 8 Oktober 2026 — ketujuh medan kepala tab Reporting Period
	// (`TabReportingPeriod` memegangnya di penampung dan menyemai kosong).
	// Kolomnya ada sejak migrasi `444` dan Save menulisnya; tanpa disemai,
	// Save berikutnya menimpanya kosong (laporan pemakai).
	"ReportingStart", "ReportingEnd", "ReportingPeriod", "ReportingInterval",
	"ReportingSubmission", "ReportingConfirmation", "ReportingSettlement",
}

// BacaKontrakWarisan membaca satu kontrak warisan, siap tampil.
func (l *Layanan) BacaKontrakWarisan(ctx context.Context, p inti.Pelaku, id string) (models.KontrakWarisan, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return models.KontrakWarisan{}, err
	}
	// ⛔ Pengenal KOSONG ditolak di sini, bukan diteruskan: kueri dengan
	// pengenal kosong mengembalikan nol baris, dan "tidak ada" adalah
	// jawaban yang berbeda dari "tidak ditanyakan".
	if strings.TrimSpace(id) == "" {
		return models.KontrakWarisan{}, fmt.Errorf("%w: pengenal kontrak kosong", ErrMasukanTidakSah)
	}
	k, err := l.gudang.BacaKontrakWarisan(ctx, id)
	if err != nil {
		return models.KontrakWarisan{}, err
	}
	// ⭐ TIGA TAB dari tabel pendaratan, bukan lagi dari CLOB. Satu
	// mekanisme untuk prop dan non-prop: nol cabang menurut
	// `PROPORTIONTYPE` di sini, dan tab yang kosong mengembalikan nol baris
	// alih-alih galat.
	if k.PeriodePelaporan, err = l.gudang.BacaPeriodePelaporan(ctx, id); err != nil {
		return models.KontrakWarisan{}, err
	}
	if k.Portofolio, err = l.gudang.BacaPortofolio(ctx, id); err != nil {
		return models.KontrakWarisan{}, err
	}
	if k.Akumulasi, err = l.gudang.BacaAkumulasi(ctx, id); err != nil {
		return models.KontrakWarisan{}, err
	}
	// Empat tab berikutnya - dimuat migrasi 430/431, dan sampai 3 Oktober
	// 2026 nol layar membacanya. Tabel terisi yang tidak dibaca siapa pun
	// adalah pekerjaan yang terlihat selesai dan tidak sampai ke pemakai.
	if k.Egnpi, err = l.gudang.BacaEgnpi(ctx, id); err != nil {
		return models.KontrakWarisan{}, err
	}
	if k.Retensi, err = l.gudang.BacaRetensi(ctx, id); err != nil {
		return models.KontrakWarisan{}, err
	}
	if k.Angsuran, err = l.gudang.BacaAngsuran(ctx, id); err != nil {
		return models.KontrakWarisan{}, err
	}
	if k.Catatan, err = l.gudang.BacaCatatan(ctx, id); err != nil {
		return models.KontrakWarisan{}, err
	}

	// ⭐⭐ KEEMPAT TAB DIBACA DARI TABEL PENDARATAN, 6 Oktober 2026.
	//
	// Keputusan pemilik proses: *"nilai yg ditarik dari JSON data itu
	// dilarang keras, gunakan table baru"*, lalu dipersempit ke keempat
	// belas tabel yang sudah berdiri. Sebelum ini keempat tab terurai dari
	// `M_TREATY_IN.JSONDATA` di dalam `BacaKontrakWarisan`.
	//
	// Dua pembacaan, sebab bentuknya dua: PIPIH untuk Limits non-prop,
	// Share, Event Limits, dan RNM Share; POHON untuk Limits proporsional.
	if k.Layer, err = l.gudang.BacaLayerPendaratan(ctx, id); err != nil {
		return models.KontrakWarisan{}, err
	}
	if k.LimitsPohon, err = l.gudang.BacaPohonLimitsPendaratan(ctx, id); err != nil {
		return models.KontrakWarisan{}, err
	}
	if k.LimitsAkar, err = l.gudang.BacaLimitsAkarPendaratan(ctx, id); err != nil {
		return models.KontrakWarisan{}, err
	}

	// ⭐⭐ MEDAN KEPALA dan grid Rate of Exchange — dari pendaratan pula.
	//
	// ⛔ `AdaDiJSON` tetap diisi, dan namanya tetap tepat: yang ia nyatakan
	// adalah "kunci ini ADA di dokumen sistem lama". Pembedaannya kini
	// datang dari `NULL` lawan teks kosong di tabel, bukan dari penunjuk
	// `*string` — arti yang sama, sumber yang berbeda.
	rev, err := l.gudang.BacaRevisiPendaratan(ctx, id)
	if err != nil {
		return models.KontrakWarisan{}, err
	}
	for nama, ke := range map[string]*string{
		"Bordeaux":              &k.Bordereaux,
		"BordereauxNote":        &k.BordereauxCatatan,
		"AccountingMode":        &k.CaraPembukuan,
		"AccountingModeNonProp": &k.CaraPembukuanNonProp,
		"ContractRefNo":         &k.NomorRujukan,
		"TreatyLeader":          &k.PemimpinTreaty,
		"IsMultipleRetro":       &k.RetroBerganda,
		"EDMState":              &k.EDMState,
		"EDMMaterialType":       &k.EDMJenisMaterial,
		"StatusAkseptasi":       &k.StatusAkseptasi,
		// ⭐ Posisi tangga akseptasi — pendaratan menimpa cadangan `TREATY_IN`.
		"Position":         &k.Posisi,
		"PositionUsername": &k.PemegangPosisi,
		// ⭐ Ketujuh medan kepala Reporting Period — migrasi `444`.
		// Ejaannya ejaan DOKUMEN; pemetaan kolomnya ada di
		// `repository.kolomRevisi`.
		"ReportingStart":        &k.PeriodeMulai,
		"ReportingEnd":          &k.PeriodeAkhir,
		"ReportingPeriod":       &k.PeriodeJenis,
		"ReportingInterval":     &k.PeriodeInterval,
		"ReportingSubmission":   &k.PeriodePenyerahan,
		"ReportingConfirmation": &k.PeriodeKonfirmasi,
		"ReportingSettlement":   &k.PeriodePelunasan,
	} {
		if v, ada := rev.Medan[nama]; ada {
			k.AdaDiJSON[nama] = true
			*ke = v
		}
	}
	k.LimitsAkar.TotalLimitsROL = rev.Medan["TotalLimitsROL"]
	// ⭐ Properti penampung halaman — migrasi `448`. Disemai ke penampung
	// saat form memuat, sehingga isian yang di-Save tampil kembali.
	k.Penampung = map[string]string{}
	for _, nama := range KunciPenampungRevisi {
		if v, ada := rev.Medan[nama]; ada {
			k.Penampung[nama] = v
		}
	}
	// ⭐ Medan akar `T_TREATY_HAZARD_LIMIT` — tab Event Limits (Non-Prop)
	// dan kedua batas Co-Ins. Tabelnya ditulis Save sejak migrasi `446`,
	// tetapi NOL yang pernah membacanya: tab itu selalu kosong.
	//
	// ⚠️ Ditaruh SESUDAH gelung di atas supaya nilai dari tabel bahaya
	// menang bila suatu hari ejaan yang sama muncul di kedua tabel —
	// `T_TREATY_HAZARD_LIMIT` yang memilikinya menurut peta pendaratan.
	bahaya, err := l.gudang.BacaBatasBahaya(ctx, id)
	if err != nil {
		return models.KontrakWarisan{}, err
	}
	for nama, v := range bahaya {
		k.Penampung[nama] = v
	}
	// ⭐ Larik total penampung (`T_TREATY_TOTAL`) — laporan pemakai
	// 8 Oktober 2026: total Share Prop "No items" sesudah Save sampai
	// Refresh ditekan.
	if k.PenampungLarik, err = l.gudang.BacaTotalPenampung(ctx, id); err != nil {
		return models.KontrakWarisan{}, err
	}
	for nama, v := range rev.Teks {
		k.AdaDiJSON[nama] = true
		k.TeksMentah[nama] = v
	}
	// ⛔ Grid Rate of Exchange dari `TREATYEXCHANGEYEARLY`.
	// ⭐ 455 (8 Oktober 2026): HANYA baris milik kontrak ini (`T_TREATY_KURS`);
	// kontrak tanpa catatan membaca kurs tahun treaty-nya seperti dulu.
	// Dibaca SESUDAH `TahunTreaty` terisi dari kolom `TREATY_IN`.
	if k.Kurs, err = l.gudang.BacaKursKontrak(ctx, id, k.TahunTreaty); err != nil {
		return models.KontrakWarisan{}, err
	}
	if k.SkalaKoasuransi, err = l.gudang.BacaSkalaKoasuransi(ctx, id); err != nil {
		return models.KontrakWarisan{}, err
	}

	// ⭐ Dua tab TEKS, dipilih menurut cabang — Jalan B, keputusan §15.
	// Dijalankan SESUDAH `SifatProporsiAsli` terisi, sebab cabangnya yang
	// menentukan ejaan mana yang dipakai.
	isiTabTeks(&k)

	// ⭐ Panel Attachment. Dua pembacaan, dan keduanya diperlukan: daftar
	// berkasnya, dan katalog kategori yang juga memuat kategori BERNOL
	// berkas - panel lama menampilkan seluruh kategori beserta cacahnya.
	if k.Lampiran, err = l.gudang.BacaLampiranKontrak(ctx, id); err != nil {
		return models.KontrakWarisan{}, err
	}
	katalog, err := l.gudang.BacaKatalogKategoriLampiran(ctx)
	if err != nil {
		return models.KontrakWarisan{}, err
	}
	// ⭐ Daftar kategori DICABANGKAN menurut sifat proporsinya — butir
	// kesepuluh berbeda, dan hanya itu. Panel yang memakai daftar cabang
	// seberang menampilkan kategori yang di cabang ini tidak pernah ada.
	namaKategori := models.NamaKategoriLampiranProp
	if !SifatProporsional(k.SifatProporsiAsli) {
		namaKategori = models.NamaKategoriLampiranNonProp
	}
	k.KategoriLampiran = SusunKategoriLampiran(namaKategori, katalog, k.Lampiran)
	for i := range k.Lampiran {
		// Terjemahan tanggal MEMAKAI ULANG penerjemah yang sama.
		k.Lampiran[i].Diunggah = TanggalTampil(k.Lampiran[i].Diunggah)
	}

	// ⭐ Panel `Existing Policy for Master ID` — pembacaan kelima belas, dan
	// satu-satunya yang menyentuh `TREATYINPRODUCTION`.
	//
	// ⚠️ Nol baris BUKAN galat: kontrak 1001846 di gambar 01 memang berbunyi
	// "No items", dan kontrak yang belum punya polis produksi adalah keadaan
	// yang lazim, bukan kegagalan.
	if k.PolisProduksi, err = l.gudang.BacaPolisProduksi(ctx, id); err != nil {
		return models.KontrakWarisan{}, err
	}

	// ⭐ Panel `Total Retention Amount` — dihitung SESUDAH `Retensi` terbaca.
	k.TotalRetensi = TotalRetensiPerMataUang(k.Retensi)

	k.SifatProporsi = SifatProporsiTampil(k.SifatProporsiAsli)
	// ⛔ MEDAN KEPALA — `dd/mm/yyyy`, bukan `dd/mm/yy`.
	//
	// Gambar 01 dokumen desain: `Commencement 01/01/2025`, sementara grid
	// Rate of Exchange di layar yang sama berbunyi `01/01/25`. Keduanya ada
	// dengan sengaja; lihat `TanggalTampilPanjang`.
	k.TanggalMulai = TanggalTampilPanjang(k.TanggalMulaiAsli)
	k.TanggalBerakhir = TanggalTampilPanjang(k.TanggalBerakhirAsli)

	// ⭐ LABEL, bukan nilai tersimpan — gambar 01 dan 26.
	//
	// ⚠️ Nilai ASLINYA tetap dibawa di medan `…Asli`: layar menampilkan
	// labelnya, dan apa pun yang kelak menulis kembali memerlukan yang
	// tersimpan. Menerjemahkan di tempat akan membuang yang asli.
	k.CaraPembukuanAsli = k.CaraPembukuan
	k.CaraPembukuanNonPropAsli = k.CaraPembukuanNonProp
	k.BordereauxAsli = k.Bordereaux
	k.CaraPembukuan = CaraPembukuanTampil(k.CaraPembukuan)
	k.CaraPembukuanNonProp = CaraPembukuanTampil(k.CaraPembukuanNonProp)
	k.Bordereaux = BordereauxTampil(k.Bordereaux)
	k.OpsiKepala = OpsiKepalaKontrak()
	// ⭐ `SetTreatyIn_Act` langkah 13: kontrak Non-Prop yang dibuka membangun
	// grid Reinstatement tiap layer (`TreatySetReinstatement`) — grid itu
	// tidak punya tabel pendaratan, persis seperti ia belum ada sebelum
	// Pega membangunnya saat kontrak dibuka.
	if !SifatProporsional(k.SifatProporsiAsli) {
		SiapkanReinstatementPohon(k.LimitsPohon)
	}

	// ⭐ Tab Share Non-Prop — pohon pendaratan, lalu larik turunannya
	// (100% Limit RNM, bagian OR / R/I, Summary, Total) dengan rumus
	// Activity yang sama. Lihat `share_np_muat.go`.
	sp, err := l.gudang.BacaSharePendaratan(ctx, id)
	if err != nil {
		return models.KontrakWarisan{}, err
	}
	k.ShareNP = ShareDariPendaratan(sp, rev.Medan)
	akarShare, err := l.gudang.BacaShareAkarRevisi(ctx, id)
	if err != nil {
		return models.KontrakWarisan{}, err
	}
	detailShare, err := l.gudang.BacaShareDetailWarisan(ctx, id)
	if err != nil {
		return models.KontrakWarisan{}, err
	}
	TerapkanAkarShare(&k.ShareNP, akarShare, detailShare)
	idTreaty := rev.Medan["ID"]
	if idTreaty == "" {
		idTreaty = id
	}
	SiapkanShareNP(&k.ShareNP, LayerDariPohon(k.LimitsPohon), idTreaty)

	// ⭐ TANGGAL DI DALAM LARIK ikut diterjemahkan.
	//
	// ⛔ Terlewat sampai 3 Oktober 2026: ketiga medan kepala sudah memakai
	// `TanggalTampil`, sementara tanggal di dalam grid tidak - di layar
	// terbaca `20250118`, bentuk simpanan, bukan bentuk baca. Satu aturan
	// yang berlaku separuh lebih buruk daripada yang tidak berlaku: yang
	// melihatnya mengira dua medan itu memang bertipe berbeda.
	//
	// Nilai TERSIMPAN tidak berubah - tabel pendaratan tetap memegang
	// `YYYYMMDD`, dan yang diterjemahkan hanya jawaban API. `TanggalTampil`
	// mengembalikan apa adanya untuk yang bukan delapan angka, jadi nilai
	// kosong dan nilai aneh lewat tanpa dikarang.
	for i := range k.Kurs {
		b := &k.Kurs[i]
		// ⛔ URUTANNYA MENGIKAT: `…Asli` diisi SEBELUM medan tampil
		// ditimpa. `TanggalTampil` menghasilkan `dd/mm/yy`, dan dari bentuk
		// itu tahun empat digitnya tidak dapat dipulihkan.
		b.BerlakuDariAsli, b.BerlakuSampaiAsli = TanggalWIB(b.BerlakuDari), TanggalWIB(b.BerlakuSampai)
		b.BerlakuDari = TanggalTampil(b.BerlakuDari)
		b.BerlakuSampai = TanggalTampil(b.BerlakuSampai)
	}
	for i := range k.PeriodePelaporan {
		p := &k.PeriodePelaporan[i]
		// `InitialDate` disimpan Pega sebagai stempel GMT — tanggal WIB-nya
		// yang layar tampilkan. Jatuh tempo sudah Date.
		p.TanggalAwal = TanggalWIB(p.TanggalAwal)
		*p = tampilkanPeriode(*p)
	}
	for i := range k.Angsuran {
		g := &k.Angsuran[i]
		g.JatuhTempoAsli, g.TanggalBayarAsli = TanggalWIB(g.JatuhTempo), TanggalWIB(g.TanggalBayar)
		g.JatuhTempo = TanggalTampil(g.JatuhTempo)
		g.TanggalBayar = TanggalTampil(g.TanggalBayar)
	}
	for i := range k.Egnpi {
		k.Egnpi[i].PerTanggal = TanggalTampil(k.Egnpi[i].PerTanggal)
	}
	for i := range k.Catatan {
		k.Catatan[i].Tanggal = TanggalTampil(k.Catatan[i].Tanggal)
	}
	for i := range k.Akumulasi {
		a := &k.Akumulasi[i]
		// ReportDate/SubDueDate tersimpan sebagai stempel DateTime GMT
		// (panjang 23) - bentuk simpannya tanggal WIB.
		a.TanggalLaporAsli, a.JatuhTempoKirimAsli = TanggalWIB(a.TanggalLapor), TanggalWIB(a.JatuhTempoKirim)
		a.TanggalLapor = TanggalTampil(a.TanggalLapor)
		a.JatuhTempoKirim = TanggalTampil(a.JatuhTempoKirim)
	}
	// ⛔ PENJAGA TERAKHIR, dan ia dijalankan SESUDAH segalanya terisi:
	// tidak ada larik `nil` yang boleh meninggalkan fungsi ini. Alasannya,
	// dan mengapa ia satu fungsi alih-alih tambalan di tempat, ada di
	// `warisan_nol.go`.
	// Pecahan spreading manual yang tak tersimpan (`hitung_share_prop.go`).
	l.lengkapiPecahanKontrak(ctx, &k)
	// Total yang tak tersimpan dihitung dari rinciannya (`total_cadangan.go`).
	lengkapiTotalPenampung(&k)
	nolkanLarik(&k)
	return k, nil
}

// WarisanTidakAda menjawab apakah galatnya "kontraknya tidak ada".
//
// Dipisah supaya handler tidak perlu mengimpor `repository` - arah
// ketergantungan `handlers -> services -> repository` dijaga penjaga inti.
func WarisanTidakAda(err error) bool { return errors.Is(err, ErrWarisanTidakAda) }

// WarisanJSONRusak menjawab apakah galatnya "dokumennya tidak dapat diurai".
func WarisanJSONRusak(err error) bool { return errors.Is(err, ErrJSONWarisanRusak) }

// TotalRetensiPerMataUang menjumlahkan `Amount` tab Maximum Retention,
// dikelompokkan per mata uang.
//
// ⛔ RUMUSNYA DISALIN dari `Activity/TreatyInNPSetTotal.xml` cabang
// `param.type=retention` — bukan diturunkan dari nama panelnya. Tiga hal
// yang Activity itu nyatakan dan yang ditiru di sini:
//
//  1. Pengelompokannya per `.Currency`, dan baris bermata-uang sama
//     DIJUMLAHKAN ke baris yang sudah ada (`Appendflag`), bukan ditambahkan
//     sebagai baris kedua.
//  2. Urutannya URUTAN KEMUNCULAN mata uang pertama kali — `<APPEND>`, bukan
//     urutan abjad. Mengurutkan abjad di sini akan membuat layar ini berbeda
//     dari layar lama tanpa satu pun sumber yang memintanya.
//  3. Nilainya dijumlahkan sebagai ANGKA. Yang tersimpan teks, jadi yang
//     tidak dapat diurai diperlakukan nol — Activity-nya punya pesan
//     "Error Amount is empty" untuk keadaan itu, dan pesan itu milik jalur
//     TULIS yang belum dibangun; layar baca-saja tidak boleh menolak
//     menampilkan apa pun karena satu baris kotor.
//
// ⚠️ Nol baris masuk -> nol baris keluar (irisan kosong, bukan nil), dan
// panelnya yang menyatakan "No items".
func TotalRetensiPerMataUang(baris []models.BarisRetensiWarisan) []models.BarisTotalRetensiWarisan {
	hasil := []models.BarisTotalRetensiWarisan{}
	di := map[string]int{}
	for _, b := range baris {
		mu := strings.TrimSpace(b.MataUang)
		i, ada := di[mu]
		if !ada {
			di[mu] = len(hasil)
			hasil = append(hasil, models.BarisTotalRetensiWarisan{MataUang: mu, Nilai: "0"})
			i = len(hasil) - 1
		}
		hasil[i].Nilai = tambahTeksAngka(hasil[i].Nilai, b.Jumlah)
	}
	return hasil
}

// tambahTeksAngka menjumlahkan dua angka yang tersimpan sebagai teks.
//
// ⚠️ `big.Float` dan bukan `float64`: nilai retensi adalah UANG, dan
// penjumlahan biner ganda memunculkan sisa seperti `0,30000000000000004`
// yang lalu tampil di layar sebagai angka yang tidak pernah diketik siapa
// pun. Presisi 200 bit melampaui lebar nilai mana pun di POOLDATA.
func tambahTeksAngka(a, b string) string {
	x, _, err := big.ParseFloat(strings.TrimSpace(a), 10, 200, big.ToNearestEven)
	if err != nil {
		x = new(big.Float).SetPrec(200)
	}
	y, _, err := big.ParseFloat(strings.TrimSpace(b), 10, 200, big.ToNearestEven)
	if err != nil {
		// Baris kotor dihitung NOL, dan itu disengaja - lihat butir 3.
		y = new(big.Float).SetPrec(200)
	}
	return new(big.Float).SetPrec(200).Add(x, y).Text('f', -1)
}
