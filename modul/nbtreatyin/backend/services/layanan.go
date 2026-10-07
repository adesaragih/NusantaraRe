package services

// Untuk apa berkas ini: LAYANAN dan bantuan bersamanya - galat yang dipetakan
// handlers, keanggotaan antrean, pra-proses flow action, dan pemuatan kasus.

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/nbtreatyin/backend/models"
	"nusantarare/modul/nbtreatyin/backend/repository"
)

// Galat yang dipetakan handlers ke kode HTTP. Galat repository diteruskan
// lewat nama di sini supaya handlers tidak mengimpor repository.
var (
	ErrTanpaOracle         = db.ErrTanpaOracle
	ErrKasusTidakAda       = repository.ErrKasusTidakAda
	ErrTahapBerubah        = repository.ErrTahapBerubah
	ErrGenerasiTertutup    = repository.ErrGenerasiTertutup
	ErrNomorPolisSudahAda  = repository.ErrNomorPolisSudahAda
	ErrDataKontrakTidakAda = repository.ErrDataKontrakTidakAda
	ErrTipeNomorKosong     = repository.ErrTipeNomorKosong
	ErrOJKKosong           = repository.ErrOJKKosong
	ErrPermintaanTidakSah  = galat.ErrPermintaanTidakSah

	// ErrKasusTertutup - kasus sudah diselesaikan; tidak ada tindakan lagi.
	ErrKasusTertutup = errors.New("services: kasus sudah diselesaikan")
	// ErrBukanAnggotaAntrean - pelaku bukan anggota antrean tempat kasus
	// menunggu (AC 14: menurut NAMA antrean, bukan nomor urut).
	ErrBukanAnggotaAntrean = fmt.Errorf("%w: bukan anggota antrean tempat berkas menunggu", inti.ErrTanpaWewenang)
	// ErrTindakanTakAdaDiPosisi - tindakan itu tidak punya tombol di layar
	// posisi kasus (mis. Save di layar atasan, nomor polis di layar admin).
	ErrTindakanTakAdaDiPosisi = errors.New("services: tindakan ini tidak tersedia di posisi berkas")
)

// GalatValidasi - pesan validasi layar (medan wajib, pesan Property-Set-
// Messages). Handlers menjawabnya 422 beserta daftar pesannya.
type GalatValidasi struct{ Pesan []string }

func (e *GalatValidasi) Error() string {
	return "services: validasi layar gagal: " + strings.Join(e.Pesan, "; ")
}

// Layanan adalah pintu aturan dagang NB Treaty In.
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

// tempatPelaku - tempat berperan tiket 05 bagi pelaku: pemetaan konstanta
// `models.PemetaanPeranTempat` (K16; kosong sampai IAM menjawab = semua
// tertunda) atas peran `inti.Pelaku.Peran` (workbasket akun).
func tempatPelaku(p inti.Pelaku) map[string]bool {
	return models.TempatTampil(models.PemetaanPeranTempat, p.PunyaPeran)
}

// ------------------------------------------------------------------ daftar dan buat

// DaftarKasus - daftar portal (`Section/SFAPortal_OpportunitiesList`). ⛔ Bab GERBANG di bawah DIGANTI RALAT
// 06-10-2026 di badan fungsi (hanya filter pembuat) - dibiarkan sebagai catatan bunyi lama.
//
// GERBANG (putaran 2, P8). Grid satu-satunya di section itu (badan REPEATING,
// `pyGridProps/pyRDName = GetListOpportunity`) bersarang di wadah
// `pyContainerVisibleWhen = OperatorID.pyWorkGroup!='ReasLife' &&
// OperatorID.pyWorkBasketList(2).pyWorkBasketName=='ReasTreatyInAdmin'`, lalu
// wadah `!IsOperatorLife` (When: `OperatorID.pyWorkGroup = "ReasLife"`). Di
// luar wadah hanya baris saringan. Maka:
//
//   - anggota `ReasTreatyInAdmin` (menurut NAMA, bukan urutan ke-2 - AC 14)
//     melihat grid itu: semua kasus terbuka, opsional per posisi (filter A
//     `pxCreateOperator = Param.UserIdentifier` tidak dibangun - antrean
//     bersama AC 11, penyimpangan putaran 1);
//   - pelaku lain tidak melihat grid itu. Tugasnya dirutekan
//     `Flow/InputRealizationTreatyIn` ke workbasket (Assignment4/6
//     `ReasTreatyInSecHead`, Assignment3 `ReasTreatyInDeptHead`,
//     `ToWorkBasket`); daftar kerja workbasket Pega tidak ada di korpus -
//     padanannya di portal tunggal ini (bab 0 butir 7): HANYA kasus yang
//     menunggu di posisi tangga yang ia pegang (AC 11, 92);
//   - tanpa satu pun posisi tangga: grid tersembunyi dan tidak ada antrean
//     -> ErrBukanAnggotaAntrean (403).
//
// ⛔ `pyWorkGroup != 'ReasLife'` (dan `!IsOperatorLife`, syarat yang sama)
// TIDAK dibangun: `inti.Pelaku` hanya membawa AkunID + workbasket, dan
// pemetaan work group Pega ke data akun tidak ada (K12 kosong; tiket 05).
func (l *Layanan) DaftarKasus(ctx context.Context, p inti.Pelaku, s models.SaringanKasus) ([]models.RingkasanKasus, error) {
	if err := l.periksaPelaku(p); err != nil {
		return nil, err
	}
	if s.Posisi != "" && !models.AdalahPosisiTangga(s.Posisi) {
		return nil, fmt.Errorf("%w: posisi %q", ErrPermintaanTidakSah, s.Posisi)
	}
	// RALAT 06-10-2026 (keputusan work owner "menu nb treaty in buat hanya filter berdasarkan create operator aja";
	// sebelumnya "isi inbox ini muncul hanya untuk akun dia saja"): HANYA filter A RD `GetListOpportunity`
	// (`A.pxCreateOperator = Param.UserIdentifier`) - berkas BUATAN akun ini di posisi mana pun, untuk siapa pun.
	// Antrean atasan tidak lagi tampil di portal (berkas yang menunggu atasan dibuka dari kotak masuk Beranda);
	// berkas tanpa CREATE_OP tidak cocok dengan akun mana pun. Saringan LINI non-life tetap (WO: "filter nonlife-nya
	// tetap"). `s.Selesai` = switch Proses / Resolved: Resolved menampilkan SEMUA berkas selesai, siapa pun
	// pembuatnya (WO 06-10-2026: "yang resolve nampilin semua yang resolve"); berkas selesai hanya-baca.
	s.Antrean, s.PembuatPosisi, s.Pembuat = nil, "", p.AkunID
	if s.Selesai {
		s.Pembuat = ""
	}
	return l.g.DaftarKasus(ctx, s)
}

// KotakMasuk - kotak masuk Beranda (keputusan work owner 06-10-2026): satu baris per workbasket tangga yang dipegang
// pelaku, urut tangga, dengan jumlah berkas yang MENUNGGU dia - Admin: buatannya yang masih di Admin; Sec Head /
// Dept Head: antrean workbasket itu. Bukan pemegang = daftar kosong (Beranda tidak gagal karenanya).
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

// DaftarMenunggu - daftar berkas kotak masuk Beranda (keputusan work owner 06-10-2026): berkas yang MENUNGGU pelaku
// di `workbasket` (kosong = semua workbasket tangga yang ia pegang) - Admin: buatannya yang masih di Admin;
// Sec Head / Dept Head: antrean itu. Workbasket di luar tangga atau tidak dipegang = ErrBukanAnggotaAntrean.
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

// BuatKasus = tombol "Create opportunity" portal (`createWork` ->
// `Flow/InputRealizationTreatyIn`, connector Start1 -> Assignment2).
//
// ⛔ Tanpa gerbang peran: syarat tombolnya (`crmCreateOpportunity`,
// `isSellingMode*`, `IsNotAdmin`) milik CRM dan tidak memuat antrean NB
// Treaty In - mengarang gerbang berarti memutuskan siapa boleh membuat
// realisasi (pola premiumlistlife).
func (l *Layanan) BuatKasus(ctx context.Context, p inti.Pelaku) (models.Kasus, error) {
	if err := l.periksaPelaku(p); err != nil {
		return models.Kasus{}, err
	}
	nama, err := l.g.NamaTampilan(ctx, p.AkunID)
	if err != nil {
		return models.Kasus{}, err
	}
	var id string
	err = l.g.Transaksi(ctx, func(tx *db.Tx) error {
		var err error
		if id, err = l.g.IDKasusBerikut(ctx, tx); err != nil {
			return err
		}
		if err := l.g.SisipKasus(ctx, tx, id, p.AkunID, nama); err != nil {
			return err
		}
		h := models.HalamanBaru()
		h.Setel("PositionNote", models.PosisiAdmin)
		h.Setel("NBStatus", models.NBStatusBaru)
		h.Setel(models.HalamanQuotation+".BusinessFac", models.BisnisTreaty)
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
	// BolehKerja - pelaku anggota antrean posisi kasus dan kasus terbuka;
	// selain itu layar hanya-baca.
	BolehKerja bool               `json:"bolehKerja"`
	Tombol     models.TombolKirim `json:"tombol"`
	// MedanWajib - jalur medan wajib layar posisi ini (berlaku saat ini,
	// `models.MedanWajibBerlaku`).
	MedanWajib []string `json:"medanWajib"`
	// Tempat - tempat berperan (tiket 05) -> tampil atau tidak.
	Tempat map[string]bool `json:"tempat"`
	Pesan  []string        `json:"pesan,omitempty"`
}

// BukaKasus membuka satu kasus. Bila pelaku boleh bekerja, pra-proses flow
// action posisinya dijalankan atas halaman (tidak disimpan - di Pega
// pra-proses hanya mengubah clipboard).
func (l *Layanan) BukaKasus(ctx context.Context, p inti.Pelaku, id string) (Layar, error) {
	if err := l.periksaPelaku(p); err != nil {
		return Layar{}, err
	}
	k, h, err := l.muat(ctx, id)
	if err != nil {
		return Layar{}, err
	}
	boleh := !k.Tertutup() && anggota(p, k.PositionNote)
	if boleh {
		if err := l.siapkan(ctx, p, k, h); err != nil {
			return Layar{}, err
		}
		if k.PositionNote == models.PosisiAdmin {
			// W5: `pyDefaultValue` sel terbuka layar admin saat dirender.
			models.TerapkanNilaiBawaanSel(h)
		}
	} else if err := l.tampilan(ctx, h); err != nil {
		return Layar{}, err
	}
	return l.layar(ctx, p, k, h, boleh)
}

func (l *Layanan) layar(ctx context.Context, p inti.Pelaku, k models.Kasus, h *models.Halaman, boleh bool) (Layar, error) {
	tempat := tempatPelaku(p)
	var wajib []string
	for _, m := range models.MedanWajibBerlaku(h, k.PositionNote, tempat) {
		wajib = append(wajib, m.Jalur)
	}
	ly := Layar{Kasus: k, Halaman: h, BolehKerja: boleh, MedanWajib: wajib, Tempat: tempat, Pesan: h.SemuaPesan()}
	if boleh {
		ly.Tombol = models.TombolUntuk(h, k.PositionNote)
	}
	return ly, nil
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
	// Total spreading TURUNAN baris, tidak disimpan (models/katalog.go).
	if err := models.HitungTotalSpreading(h); err != nil {
		return models.Kasus{}, nil, err
	}
	return k, h, nil
}

// siapkan = pra-proses flow action posisi kasus:
//
//	Admin   `InputPolicyTreatyIn_preDT`, lalu `InputPolicyTreatyInPre_Act`
//	Atasan  `DeptHeadTreatyInUW_preDT`,  lalu `InputPolicyTreatyInPre_Act`
//
// Urutan Pega: data transform pra-proses lebih dulu, aktivitas pra-proses
// sesudahnya. Bagian `InputPolicyTreatyInPre_Act` yang dibangun: langkah 2
// (bisnis bila BizCode kosong), 3-4 dan 9 (tanggal; hari tutup buku dari
// TANGGAL_CLOSING - lihat `models.GeserTanggalProduksi`), dan 10
// (`TreatyRealizationCheckXOLList` - RALAT K8, `siapkanNonProp`). Langkah 1
// (`SetCategoryAttach` - lampiran, tidak ada di layar realisasi) tidak
// dibangun; langkah 5-8 menyiapkan daftar pilihan spreading (`Acuan`).
//
// Halaman master `TreatyIn` dimuat ulang dari view lewat `TreatyIn.ID` (P29).
func (l *Layanan) siapkan(ctx context.Context, p inti.Pelaku, k models.Kasus, h *models.Halaman) error {
	nama, err := l.g.NamaTampilan(ctx, p.AkunID)
	if err != nil {
		return err
	}
	sekarang := l.jam()
	if k.PositionNote == models.PosisiAdmin {
		models.PraprosesAdmin(h, sekarang, nama)
	} else {
		models.PraprosesAtasan(h, sekarang, nama)
	}
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
	if err := l.muatMaster(ctx, h); err != nil {
		return err
	}
	// langkah 10 dan master XOL (K8) - nonprop.go
	return l.siapkanNonProp(ctx, h)
}

// tampilan - buka berkas HANYA-LIHAT (bukan pemegang posisinya, atau berkas tertutup): pra-proses tidak
// dijalankan karena tidak ada pekerjaan, tetapi halaman master `TreatyIn` yang TIDAK disimpan tetap dibaca
// ulang - Commencement, Termination, deret layer, dan subsection NonProp. Di Pega halaman itu ikut tersimpan
// di clipboard berkas, sehingga tampil juga saat berkas hanya dilihat. Nol tulisan.
func (l *Layanan) tampilan(ctx context.Context, h *models.Halaman) error {
	if err := l.muatMaster(ctx, h); err != nil {
		return err
	}
	return l.tampilanNonProp(ctx, h)
}

// muatMaster mengisi halaman TreatyIn dari baris view kontrak terpilih.
func (l *Layanan) muatMaster(ctx context.Context, h *models.Halaman) error {
	id := h.Ambil(models.HalamanMaster + ".ID")
	if id == "" {
		return nil
	}
	b, err := l.g.DetailKontrak(ctx, id)
	if err != nil {
		return err
	}
	models.TerapkanMasterKontrak(h, b)
	return nil
}

// RiwayatKasus - riwayat akseptasi satu kasus, berurut waktu (AC 72).
func (l *Layanan) RiwayatKasus(ctx context.Context, p inti.Pelaku, id string) ([]models.Riwayat, error) {
	if err := l.periksaPelaku(p); err != nil {
		return nil, err
	}
	if _, err := l.g.Keadaan(ctx, nil, id); err != nil {
		return nil, err
	}
	return l.g.DaftarRiwayat(ctx, models.KunciInstans(id))
}

// Acuan adalah daftar pilihan layar.
type Acuan struct {
	MataUang  []models.Pilihan `json:"mataUang"`
	MO        []models.Pilihan `json:"mo"`
	Spreading []models.Pilihan `json:"spreading"`
	JenisReas []models.Pilihan `json:"jenisReas"`
}

// DaftarAcuan membaca keempat daftar pilihan layar.
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
	if a.Spreading, err = l.g.DaftarJenisSpreading(ctx); err != nil {
		return Acuan{}, err
	}
	if a.JenisReas, err = l.g.DaftarJenisReas(ctx); err != nil {
		return Acuan{}, err
	}
	return a, nil
}

// DaftarBisnis - isi grid AKTIF popup `BusinessAndSOBList` (RD
// `BrowseTreatyJoinEDM`, view TREATYINDETAILJOINEDM - `repository.DaftarBisnis`).
//
// Tombol `Choose Business` (`Section/DetailPolicyTreatyIn`) = showHarness
// `pySubmitData=Yes`: isian layar dikirim lebih dulu, lalu grid membaca
// `.QuotationData.ProportionalType` halaman kerja. Maka halaman kiriman
// digabung seperti tindakan lain (`kerjakan`: keanggotaan, kasus terbuka,
// medan admin), TANPA simpan - showHarness tidak ber-Obj-Save.
//
// Tombol itu hanya ada di layar admin dan hanya bila wadahnya tampil
// (`.ClaimType != 'XOL Retro'`); selain itu 409, sama dengan `PilihBisnis`.
func (l *Layanan) DaftarBisnis(ctx context.Context, p inti.Pelaku, id string, masuk *models.Halaman, saringan map[string]string) ([]models.BarisKontrak, error) {
	k, h, err := l.kerjakan(ctx, p, id, masuk)
	if err != nil {
		return nil, err
	}
	if err := bolehPilihBisnis(k, h); err != nil {
		return nil, err
	}
	s := models.SaringanPopupBisnis(h)
	s.Kolom = models.SaringanKolomSah(saringan) // keputusan WO 06-10-2026: saringan dicari di server
	return l.g.DaftarBisnis(ctx, s)
}

// bolehPilihBisnis - tombol `Choose Business` ada di posisi kasus dan tampil
// menurut isian layar.
func bolehPilihBisnis(k models.Kasus, h *models.Halaman) error {
	if k.PositionNote != models.PosisiAdmin {
		return ErrTindakanTakAdaDiPosisi
	}
	if !models.TampilPilihBisnis(h) {
		return fmt.Errorf("%w: tombol Choose Business tidak tampil bila ClaimType '%s'",
			ErrTindakanTakAdaDiPosisi, models.KlaimXOLRetro)
	}
	return nil
}
