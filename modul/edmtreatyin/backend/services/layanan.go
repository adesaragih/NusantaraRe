package services

// Untuk apa berkas ini: LAYANAN dan bantuan bersamanya - galat yang dipetakan handlers, keanggotaan antrean, daftar
// portal, layar Create (`TreatyCreateEdm` -> `CreateEDMT`), pembukaan kasus dan pra-proses flow action.

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/dokumenpolis"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/edmtreatyin/backend/models"
	"nusantarare/modul/edmtreatyin/backend/repository"
)

// Galat yang dipetakan handlers ke kode HTTP. Galat repository diteruskan lewat nama di sini supaya handlers tidak
// mengimpor repository.
var (
	ErrTanpaOracle            = db.ErrTanpaOracle
	ErrKasusTidakAda          = repository.ErrKasusTidakAda
	ErrTahapBerubah           = repository.ErrTahapBerubah
	ErrGenerasiTertutup       = repository.ErrGenerasiTertutup
	ErrGenerasiSudahDiendorse = repository.ErrGenerasiSudahDiendorse
	ErrNomorPolisSudahAda     = repository.ErrNomorPolisSudahAda
	ErrPermintaanTidakSah     = galat.ErrPermintaanTidakSah

	// ErrKasusTertutup - kasus sudah diselesaikan; tidak ada tindakan lagi.
	ErrKasusTertutup = errors.New("services: kasus sudah diselesaikan")
	// ErrBukanAnggotaAntrean - pelaku bukan anggota antrean tempat kasus menunggu (menurut NAMA antrean).
	ErrBukanAnggotaAntrean = fmt.Errorf("%w: bukan anggota antrean tempat berkas menunggu", inti.ErrTanpaWewenang)
	// ErrTindakanTakAdaDiPosisi - tindakan itu tidak punya tombol di layar posisi kasus.
	ErrTindakanTakAdaDiPosisi = errors.New("services: tindakan ini tidak tersedia di posisi berkas")
)

// GalatValidasi - pesan validasi layar (medan wajib, Property-Set-Messages, TrtERR.CARI1). Handlers menjawabnya
// 422 beserta daftar pesannya.
type GalatValidasi struct{ Pesan []string }

func (e *GalatValidasi) Error() string {
	return "services: validasi layar gagal: " + strings.Join(e.Pesan, "; ")
}

// Layanan adalah pintu aturan dagang EDM Treaty In.
type Layanan struct {
	g   Gudang
	jam func() time.Time
	// konversi dan produksi - efek keluar sesudah selesai (konversi.go).
	konversi PengirimKonversi
	produksi bool
	// pemuat - Copy Old (copyold.go); nil bila gudang tidak membaca JSON_POLIS (tombol tidak tampil).
	pemuat *Pemuat
}

// Baru menyusun layanan; `g` nil = tanpa Oracle (setiap tindakan 503). Gudang yang juga `GudangPemuat` (Oracle)
// menyalakan Copy Old.
func Baru(g Gudang, jam func() time.Time) *Layanan {
	if jam == nil {
		jam = time.Now
	}
	l := &Layanan{g: g, jam: jam}
	if gp, ok := g.(GudangPemuat); ok {
		l.pemuat = PemuatBaru(gp)
	}
	return l
}

// AdaGudang - layanan tersambung ke penyimpanan.
func (l *Layanan) AdaGudang() bool { return l != nil && l.g != nil }

// periksaPelaku - layanan tersambung dan permintaan membawa identitas.
func (l *Layanan) periksaPelaku(p inti.Pelaku) error {
	if !l.AdaGudang() {
		return ErrTanpaOracle
	}
	return inti.WajibIdentitas(p)
}

// anggota menjawab: pelaku anggota antrean (workbasket) `posisi`.
func anggota(p inti.Pelaku, posisi string) bool { return posisi != "" && p.PunyaPeran(posisi) }

// ------------------------------------------------------------------ portal dan kotak masuk

// DaftarKasus - grid portal `Section/SFAPortal_Endorsement_Treaty` (RD `InboxEDM_RD2`): filter A berkas BUATAN
// akun (`.pxCreateOperator = Param.UserNameID`), B kotak saring "Policy Number" (`OldPolicyNo Contains`), C belum
// selesai. ⛔ Tidak ada tab status di portal EDM (berbeda dari portal NB yang diberi switch In Progress / Resolved
// atas keputusan WO untuk NB - pertanyaan untuk WO apakah berlaku juga di EDM).
func (l *Layanan) DaftarKasus(ctx context.Context, p inti.Pelaku, cari string, selesai bool) ([]models.RingkasanKasus, error) {
	if err := l.periksaPelaku(p); err != nil {
		return nil, err
	}
	// Aturan portal NB Treaty In berlaku untuk EDM (keputusan work owner 07-10-2026 "YA"; NB 06-10-2026): In
	// Progress = berkas BUATAN akun ini (filter A) yang masih proses; Resolved = berkas selesai BUATAN akun ini (RALAT
	// 07-10-2026 "TAMBAHKAN KAN UNTUK PEMBUAT. MENU ITU HANYA UNTUK SI PEMBUAT, NB DAN EDM TREATY"; dulu semua berkas
	// selesai); berkas selesai hanya-baca (`Layar.BolehKerja`).
	s := models.SaringanKasus{Cari: cari, Pembuat: p.AkunID, Selesai: selesai}
	return l.g.DaftarKasus(ctx, s)
}

// KotakMasuk - kotak masuk Beranda (kontrak `MenuModul.antreanBeranda`, pola NB): satu baris per workbasket tangga
// yang dipegang pelaku, dengan jumlah berkas yang MENUNGGU dia. ⚠️ Alur XML merutekan berkas ke workbasket
// ReasTreatyInSecHead / ReasTreatyInDeptHead (`ToWorkBasket`); daftar kerja workbasket Pega tidak ada di korpus -
// padanannya di aplikasi ini kotak masuk Beranda (keputusan WO untuk NB 06-10-2026).
func (l *Layanan) KotakMasuk(ctx context.Context, p inti.Pelaku) ([]models.AntreanKotakMasuk, error) {
	if err := l.periksaPelaku(p); err != nil {
		return nil, err
	}
	admin := anggota(p, models.PosisiAdmin)
	var atasan []string
	for _, pos := range models.PosisiTangga {
		if pos != models.PosisiAdmin && anggota(p, pos) {
			atasan = append(atasan, pos)
		}
	}
	cacah, err := l.g.HitungKotakMasuk(ctx, p.AkunID, admin, atasan)
	if err != nil {
		return nil, err
	}
	out := []models.AntreanKotakMasuk{}
	for _, pos := range models.PosisiTangga {
		if !anggota(p, pos) {
			continue
		}
		pk, err := l.g.PemegangKotakMasuk(ctx, pos)
		if err != nil {
			return nil, err
		}
		nama := pk.NamaWorkbasket
		if strings.TrimSpace(nama) == "" {
			nama = pos
		}
		out = append(out, models.AntreanKotakMasuk{Workbasket: pos, Nama: nama, Jumlah: cacah[pos]})
	}
	return out, nil
}

// DaftarMenunggu - daftar berkas kotak masuk Beranda: berkas yang MENUNGGU pelaku di `workbasket` (kosong = semua
// workbasket tangga yang ia pegang) - Admin: buatannya yang masih di Admin; Sec Head / Dept Head: antrean itu.
func (l *Layanan) DaftarMenunggu(ctx context.Context, p inti.Pelaku, workbasket string) ([]models.RingkasanKasus, error) {
	if err := l.periksaPelaku(p); err != nil {
		return nil, err
	}
	if workbasket != "" && (!models.AdalahPosisiTangga(workbasket) || !anggota(p, workbasket)) {
		return nil, fmt.Errorf("%w (kotak masuk %q)", ErrBukanAnggotaAntrean, workbasket)
	}
	s := models.SaringanKasus{Posisi: workbasket}
	if anggota(p, models.PosisiAdmin) && (workbasket == "" || workbasket == models.PosisiAdmin) {
		s.Pembuat, s.PembuatPosisi = p.AkunID, models.PosisiAdmin
	}
	for _, pos := range models.PosisiTangga {
		if pos != models.PosisiAdmin && anggota(p, pos) && (workbasket == "" || workbasket == pos) {
			s.Antrean = append(s.Antrean, pos)
		}
	}
	if s.Pembuat == "" && len(s.Antrean) == 0 {
		return []models.RingkasanKasus{}, nil
	}
	return l.g.DaftarKasus(ctx, s)
}

// ------------------------------------------------------------------ layar Create

// HasilPeriksaPolis - jawaban sel "No Polis Treaty" layar `TreatyCreateEdm` (change -> TrtEdmCheckPolicyError lalu
// CheckNopolisAvailability).
type HasilPeriksaPolis struct {
	// NoMaster - "No Master Treaty" (`DisplayData.CARI2`, TrtEdmCheckPolicyError langkah 6-7).
	NoMaster string `json:"noMaster"`
	// Galat - `TrtERR.CARI1` (kosong = tombol Create boleh tampil bila No Polis dan No Master terisi).
	Galat string `json:"galat"`
}

// PeriksaPolis = `Activity/TrtEdmCheckPolicyError` (2-4 RD GetListEdmTreaty: EDM polis itu yang belum selesai ->
// pesan; 6-7 FetchNoOfferFromNoPolis -> No Master) lalu `Activity/CheckNopolisAvailability` (FetchNopolisCount:
// tidak ada di json_polis -> pesan). Urutan sel: CheckNopolisAvailability berjalan SESUDAH dan menimpa TrtERR.CARI1
// bila nomor tidak ada.
func (l *Layanan) PeriksaPolis(ctx context.Context, p inti.Pelaku, nopolis string) (HasilPeriksaPolis, error) {
	if err := l.periksaPelaku(p); err != nil {
		return HasilPeriksaPolis{}, err
	}
	var r HasilPeriksaPolis
	nopolis = strings.TrimSpace(nopolis)
	if nopolis == "" {
		return r, nil
	}
	berjalan, err := l.g.AdaEDMBerjalan(ctx, nopolis)
	if err != nil {
		return r, err
	}
	if berjalan {
		r.Galat = models.PesanEDMBelumSelesai
	}
	if r.NoMaster, err = l.g.NoMasterDariNoPolis(ctx, nopolis); err != nil {
		return r, err
	}
	cacah, err := l.g.CacahGenerasiJSONPolis(ctx, nopolis)
	if err != nil {
		return r, err
	}
	if cacah == 0 {
		r.Galat = models.PesanNopolisSalah
	}
	return r, nil
}

// PermintaanBuat - isian layar `TreatyCreateEdm`.
type PermintaanBuat struct {
	// NoPolis - "No Polis Treaty" (`DisplayData.CARI1`).
	NoPolis string `json:"noPolis"`
	// EDMType - "Source of Change" (`DisplayData.CARI3` -> Param.edmtype).
	EDMType string `json:"edmType"`
}

// BuatKasus = tombol "Create" layar `TreatyCreateEdm` -> `Activity/CreateEDMT(edmtype)` -> `openWorkByHandle`.
// Tombol tampil bila `TrtERR.CARI1 = ” && DisplayData.CARI1 != ” && DisplayData.CARI2 != ”` - syarat itu
// diperiksa ulang di server (`PeriksaPolis`). Source of Change TIDAK disyaratkan tombol (XML) - tidak diwajibkan.
//
// Satu transaksi: generasi terakhir polis (bernomor, PRODKE terbesar) -> ID kasus `EDMT-<n>` -> baris T_WORK_POLIS
// + generasi baru (PRODKE + 1, OLD_POLIS_ID = generasi itu, NOENDORS = EDMNo, EDM_TYPE) -> halaman lahir
// (`models.RakitHalamanBaru`). Endorsemen serentak atas generasi yang sama: yang kedua ditolak basis data (UNIQUE
// OLD_POLIS_ID, AC 4) -> 409.
func (l *Layanan) BuatKasus(ctx context.Context, p inti.Pelaku, r PermintaanBuat) (models.Kasus, error) {
	periksa, err := l.PeriksaPolis(ctx, p, r.NoPolis)
	if err != nil {
		return models.Kasus{}, err
	}
	nopolis := strings.TrimSpace(r.NoPolis)
	if nopolis == "" || periksa.NoMaster == "" {
		return models.Kasus{}, fmt.Errorf("%w: tombol Create tidak tampil (No Polis / No Master kosong)", ErrPermintaanTidakSah)
	}
	if periksa.Galat != "" {
		return models.Kasus{}, &GalatValidasi{Pesan: []string{periksa.Galat}}
	}
	nama, err := l.g.NamaTampilan(ctx, p.AkunID)
	if err != nil {
		return models.Kasus{}, err
	}
	var id string
	err = l.g.Transaksi(ctx, func(tx *db.Tx) error {
		terakhir, err := l.g.GenerasiTerakhir(ctx, tx, nopolis)
		if errors.Is(err, repository.ErrGenerasiPolisTidakAda) {
			return &GalatValidasi{Pesan: []string{models.PesanGenerasiBelumDimuat}}
		}
		if err != nil {
			return err
		}
		// Tinjauan kode 06-10-2026: json_polis lebih maju dari rantai relasional (generasi Pega belum dimuat) ->
		// tolak. Utility1 menomori json_polis dengan cacah barisnya (`TreatyInSearchProdKe`), jadi generasi baru
		// PRODKE terakhir+1 hanya sah bila json_polis memuat tepat terakhir+1 generasi.
		cacah, err := l.g.CacahGenerasiJSONPolis(ctx, nopolis)
		if err != nil {
			return err
		}
		if cacah > terakhir.ProdKe+1 {
			return &GalatValidasi{Pesan: []string{models.PesanGenerasiBelumDimuat}}
		}
		lama, err := l.g.BacaGenerasi(ctx, tx, terakhir.ID)
		if err != nil {
			return err
		}
		if id, err = l.g.IDKasusBerikut(ctx, tx); err != nil {
			return err
		}
		prodKe := terakhir.ProdKe + 1
		h := models.RakitHalamanBaru(lama, nopolis, strings.TrimSpace(r.EDMType), p.AkunID, prodKe)
		gb := repository.GenerasiBaru{
			ProdKe: prodKe, EDMNo: h.Ambil(models.HalamanPolis + ".EDMNo"),
			OldPolisID: terakhir.ID, EDMType: h.Ambil(models.HalamanPolis + ".EDMType"),
		}
		if err := l.g.SisipKasus(ctx, tx, id, p.AkunID, nama, gb); err != nil {
			return err
		}
		return l.g.SimpanHalaman(ctx, tx, id, h)
	})
	if err != nil {
		return models.Kasus{}, err
	}
	return l.g.Keadaan(ctx, nil, id)
}

// ------------------------------------------------------------------ buka

// Layar adalah satu kasus siap ditampilkan.
type Layar struct {
	Kasus   models.Kasus    `json:"kasus"`
	Halaman *models.Halaman `json:"halaman"`
	// BolehKerja - pelaku anggota antrean posisi kasus dan kasus terbuka; selain itu layar hanya-baca.
	BolehKerja bool               `json:"bolehKerja"`
	Tombol     models.TombolKirim `json:"tombol"`
	// MedanWajib - jalur medan wajib layar posisi ini (berlaku saat ini).
	MedanWajib []string `json:"medanWajib"`
	Pesan      []string `json:"pesan,omitempty"`
}

// BukaKasus membuka satu kasus. Bila pelaku boleh bekerja, pra-proses flow action posisinya dijalankan atas
// halaman (tidak disimpan - di Pega pra-proses hanya mengubah clipboard).
func (l *Layanan) BukaKasus(ctx context.Context, p inti.Pelaku, id string) (Layar, error) {
	if err := l.periksaPelaku(p); err != nil {
		return Layar{}, err
	}
	k, h, err := l.muat(ctx, id)
	if err != nil {
		return Layar{}, err
	}
	boleh := !k.Tertutup() && !k.GenerasiTertutup && anggota(p, k.PositionNote)
	if boleh {
		if err := l.siapkan(ctx, p, k, h); err != nil {
			return Layar{}, err
		}
		models.TerapkanNilaiBawaanSel(h)
	}
	return l.layar(k, h, boleh), nil
}

func (l *Layanan) layar(k models.Kasus, h *models.Halaman, boleh bool) Layar {
	var wajib []string
	for _, m := range models.MedanWajibBerlaku(h, k.PositionNote) {
		wajib = append(wajib, m.Jalur)
	}
	ly := Layar{Kasus: k, Halaman: h, BolehKerja: boleh, MedanWajib: wajib, Pesan: h.SemuaPesan()}
	if boleh {
		ly.Tombol = models.TombolUntuk(h, k.PositionNote)
	}
	return ly
}

// muat membaca keadaan dan halaman tersimpan satu kasus.
func (l *Layanan) muat(ctx context.Context, id string) (models.Kasus, *models.Halaman, error) {
	k, err := l.g.Keadaan(ctx, nil, id)
	if err != nil {
		return models.Kasus{}, nil, err
	}
	h, err := l.g.BacaHalaman(ctx, nil, id)
	if err != nil {
		return models.Kasus{}, nil, err
	}
	h.Setel("Position", k.Position)
	if k.PositionNote != "" {
		h.Setel("PositionNote", k.PositionNote)
	}
	// pyWorkIDPrefix (EDMChooseBusiness_Act langkah 11-12)
	h.Setel("pyWorkIDPrefix", models.AwalanKasus)
	// Total spreading dan total selisih TURUNAN baris, tidak disimpan.
	if err := models.HitungTotalSpreading(h); err != nil {
		return models.Kasus{}, nil, err
	}
	if err := models.HitungTotalSpreadingOldData(h); err != nil {
		return models.Kasus{}, nil, err
	}
	if err := models.HitungTotalSelisih(h); err != nil {
		return models.Kasus{}, nil, err
	}
	return k, h, nil
}

// siapkan = pra-proses flow action posisi kasus:
//
//	Admin   `InboxPolicyTreatyInAddendum`: pra-DT `InputPolicyTreatyInAddendum_preAddDT`
//	        (`models.PraprosesAdmin`), lalu pra-activity `InputPolicyTreatyInPre_Act` (identik NB):
//	        1  `SetCategoryAttach` (lampiran Data-OfferFacIn - tidak ada di layar EDM) - tidak dibangun
//	        2  bisnis bila BizCode kosong (RDB GetOldIDBusiness_SQL)
//	        3-4 StatementDate / ProductionDate dari tanggal SQL (WHEN tak dicentang - selalu)
//	        5-8 daftar pilihan spreading (`Acuan`)
//	        9  ProductionDate geser bila tanggal > closing
//	        10 `TreatyRealizationCheckXOLList` (IsNewPolicyNonProp == 1 && EDMType != "3")
//	Atasan  `DeptHeadTreatyIn_UWAddendum`: pra-DT `DeptHeadTreatyInAddendum_PreDT` (`models.PraprosesAtasan`),
//	        pra-activity `SetCategoryAttach` versi Work (langkah 1 lampiran - tidak dibangun; 2-4 daftar pilihan
//	        spreading `SelectSpreadingTreatyInProduction` - `Acuan`)
func (l *Layanan) siapkan(ctx context.Context, p inti.Pelaku, k models.Kasus, h *models.Halaman) error {
	nama, err := l.g.NamaTampilan(ctx, p.AkunID)
	if err != nil {
		return err
	}
	sekarang := l.jam()
	if k.PositionNote != models.PosisiAdmin {
		models.PraprosesAtasan(h, sekarang, nama)
		return nil
	}
	models.PraprosesAdmin(h, sekarang, nama)
	if h.Ambil(models.HalamanPolis+".BizCode") == "" && h.Ambil(models.HalamanQuotation+".BusinessName") != "" {
		b, err := l.g.BisnisDariKunci(ctx, models.KunciCariBisnis(h.Ambil(models.HalamanQuotation+".BusinessName"), true))
		if err != nil {
			return err
		}
		models.TerapkanBisnisPra(h, b)
	}
	var hari int
	if err := l.g.Transaksi(ctx, func(tx *db.Tx) error {
		var err error
		hari, err = l.g.HariClosing(ctx, tx)
		return err
	}); err != nil {
		return err
	}
	models.PraprosesTanggal(h, sekarang, hari)
	return l.siapkanNonProp(ctx, h)
}

// Acuan adalah daftar pilihan layar.
type Acuan struct {
	MataUang  []models.Pilihan `json:"mataUang"`
	MO        []models.Pilihan `json:"mo"`
	JenisReas []models.Pilihan `json:"jenisReas"`
	// JenisEDM - pilihan "Source of Change" (`DisplayData.CARI3`, DT `TreatyEDMListType`): kode
	// models.KodeJenisEDM, teks models.LabelJenisEDM (screenshot Pega, work owner 07-10-2026).
	JenisEDM []models.Pilihan `json:"jenisEdm"`
}

// DaftarAcuan membaca daftar pilihan layar.
func (l *Layanan) DaftarAcuan(ctx context.Context, p inti.Pelaku) (Acuan, error) {
	if err := l.periksaPelaku(p); err != nil {
		return Acuan{}, err
	}
	var a Acuan
	var err error
	if a.MataUang, err = l.g.DaftarMataUang(ctx); err != nil {
		return Acuan{}, err
	}
	if a.MO, err = l.g.DaftarMO(ctx); err != nil {
		return Acuan{}, err
	}
	if a.JenisReas, err = l.g.DaftarJenisReas(ctx); err != nil {
		return Acuan{}, err
	}
	for _, kode := range models.KodeJenisEDM {
		a.JenisEDM = append(a.JenisEDM, models.Pilihan{Nilai: kode, Label: models.LabelJenisEDM[kode]})
	}
	return a, nil
}

// KasusLampiran - kasus pemilik lampiran "Reas" (`inti/backend/dokumenpolis`, grid `AttachmentGridReas` NB FacIn yang
// dipinjam EDM Treaty In lewat `SetCategoryAttach`): lampiran baru berkunci `KunciInstans`, lampiran Pega lama berkunci
// pzInsKey `ASM-FW-GISFW-WORK <pyID>`. Upload / Delete selama kasus belum Resolve - keputusan work owner 08-10-2026:
// "semua bisa asal belum resolve".
func (l *Layanan) KasusLampiran(ctx context.Context, p inti.Pelaku, id string) (dokumenpolis.Kasus, error) {
	if err := l.periksaPelaku(p); err != nil {
		return dokumenpolis.Kasus{}, err
	}
	k, err := l.g.Keadaan(ctx, nil, id)
	if err != nil {
		return dokumenpolis.Kasus{}, err
	}
	return dokumenpolis.KasusDari(models.KunciInstans(k.ID), models.PyIDKasus(k.ID), !k.Tertutup()), nil
}
