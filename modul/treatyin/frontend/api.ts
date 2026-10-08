// modul/treatyin/frontend/api.ts - panggilan backend modul Treaty In, satu
// fungsi per rute (`backend/handlers/rute_treaty_in.go`). Klien HTTP-nya
// `inti/frontend/klien.ts`.
//
// ⛔ SATU rute, dan itu disengaja. Papan tiket modul ini menyatakan `L-4`:
// "tidak ada spesifikasi layar di mana pun". Rute kontrak, versi, dan
// persetujuan lahir bersama spesifikasinya.

import { minta, mintaFormulir } from '../../../inti/frontend/klien'

/** Prefix rute API modul ini - SAMA dengan `handlers.Prefix`. */
export const PREFIX_TREATYIN = '/api/treaty-in'

/** LIMA himpunan acuan tiket 15 - SAMA dengan `models.Himpunan`.
 *
 * `mata-uang` dicabut 4 Oktober 2026: kurs dan daftarnya dari
 * `TREATYEXCHANGEYEARLY`. */
export const HIMPUNAN_ACUAN = [
  'jenis-potongan',
  'kelas-bisnis',
  'kelompok-treaty',
  'bahaya',
  'jenis-reasuransi',
] as const

export type HimpunanAcuan = (typeof HIMPUNAN_ACUAN)[number]

/**
 * Satu baris tabel acuan - SAMA dengan `models.Acuan`.
 *
 * `idInduk` hanya terisi pada `jenis-reasuransi`, yang bersusun. Pada kelima
 * himpunan lain ia SELALU kosong, dan kosong di sana berarti "tabel ini memang
 * tidak bersusun" - bukan "induknya belum diisi".
 */
export interface Acuan {
  id: number
  kode: string
  nama: string
  aktif: string
  idInduk?: number
}

/** `GET /api/treaty-in/acuan/{himpunan}` - isi satu tabel acuan, urut KODE. */
export async function ambilAcuan(himpunan: HimpunanAcuan): Promise<Acuan[]> {
  return minta<Acuan[]>(`${PREFIX_TREATYIN}/acuan/${himpunan}`)
}

// ===========================================================================
// Layar daftar kontrak — ronde layar 1
// ===========================================================================

/**
 * Satu baris layar daftar - SAMA dengan `models.BarisDaftarKontrak`.
 *
 * ⚠️ `posisiKe` SELALU kosong hari ini, dan itu bukan cacat: nol kolom di
 * model baru menyimpannya. Jejak per kolom ada di berkas model Go-nya.
 */
export interface BarisDaftarKontrak {
  id: number
  namaKontrak: string
  sifatProporsi: string
  idAsalBisnis: number
  idCedant: number
  tanggalMulai: string
  tanggalBerakhir: string
  keadaanSiklusHidup: string
  posisiKe: string
}

/** `GET /api/treaty-in/kontrak` - seluruh kontrak, urut pengenal. */
export async function ambilDaftarKontrak(): Promise<BarisDaftarKontrak[]> {
  return minta<BarisDaftarKontrak[]>(`${PREFIX_TREATYIN}/kontrak`)
}

// ===========================================================================
// Layar daftar dari tabel WARISAN `POOLDATA.TREATY_IN`
// ===========================================================================

/**
 * Satu baris daftar warisan - SAMA dengan `models.BarisDaftarWarisan`.
 *
 * ⛔ `id` TEKS, bukan angka. Kolom `TREATY_IN.ID` adalah `VARCHAR2(100)`.
 * Ke-1.854 nilainya kebetulan berupa angka tujuh digit; kebetulan bukan
 * jaminan, dan kolomnya tidak mengikat baris ke-1.855 pada apa pun.
 *
 * Tiap medan terjemahan berdampingan dengan `…Asli`-nya: layar membaca yang
 * diterjemahkan, yang menyelidiki selisih pemindahan membaca yang asli.
 */
export interface BarisDaftarWarisan {
  id: string
  namaKontrak: string
  sifatProporsiAsli: string
  sifatProporsi: string
  asalBisnis: string
  cedant: string
  tanggalMulaiAsli: string
  tanggalBerakhirAsli: string
  tanggalMulai: string
  tanggalBerakhir: string
  posisiKe: string
  statusAkseptasi: string
  /** `POSITION` — syarat tampil tombol `Revision` (`.Position = ''`). */
  posisi: string
}

/** Satu halaman beserta cacah seluruhnya - SAMA dengan `models.HalamanDaftarWarisan`. */
export interface HalamanDaftarWarisan {
  baris: BarisDaftarWarisan[]
  halaman: number
  ukuran: number
  total: number
}

/**
 * `GET /api/treaty-in/kontrak-warisan?halaman=n` - satu halaman tabel warisan.
 *
 * ⚠️ BERBEDA dari `ambilDaftarKontrak`, dan bedanya bukan teknis: yang itu
 * mengembalikan kontrak yang sudah DIPINDAHKAN ke model baru (hari ini nol),
 * yang ini kontrak yang ADA DI SISTEM LAMA (1.854).
 */
export async function ambilDaftarWarisan(halaman: number): Promise<HalamanDaftarWarisan> {
  return minta<HalamanDaftarWarisan>(`${PREFIX_TREATYIN}/kontrak-warisan?halaman=${String(halaman)}`)
}

/**
 * Satu kontrak warisan - SAMA dengan `models.KontrakWarisan`.
 *
 * ⛔ `adaDiJson` memisahkan "kuncinya TIDAK ADA di dokumen" dari "kuncinya
 * ada bernilai kosong". Keduanya berarti hal yang berbeda di layar: yang
 * pertama "tidak ada di sistem lama", yang kedua "belum diisi" — dan kotak
 * kosong tidak dapat menyatakan keduanya.
 */
export interface KontrakWarisan {
  id: string
  namaKontrak: string
  lingkupWilayah: string
  tahunTreaty: string
  cedant: string
  idCedant: string
  asalBisnis: string
  idAsalBisnis: string
  sifatProporsiAsli: string
  sifatProporsi: string
  tanggalMulaiAsli: string
  tanggalBerakhirAsli: string
  tanggalMulai: string
  tanggalBerakhir: string
  bordereaux: string
  bordereauxCatatan: string
  caraPembukuan: string
  /** Nilai TERSIMPAN di balik label — jalur tulis memerlukannya. */
  caraPembukuanAsli: string
  bordereauxAsli: string
  caraPembukuanNonPropAsli: string
  /**
   * ⛔ `AccountingModeNonProp` — properti KEDUA, bukan medan yang sama.
   * Cabang non-proporsional membaca yang ini; nilainya `loss` / `risk`,
   * bukan `underwriting` / `accounting`.
   */
  caraPembukuanNonProp: string
  /**
   * ⭐ Tiga nilai yang memutuskan TAB MANA yang dirender, apa adanya dari
   * dokumen. Kosong berarti kuncinya tidak ada, dan itu BUKAN `false`.
   */
  retroBerganda: string
  edmState: string
  edmJenisMaterial: string
  /**
   * `TreatyIn.StatusAkseptasi` — tab `Information & Submit` menyembunyikan
   * `Submit` ketika nilainya `Resolve Complete`.
   */
  statusAkseptasi: string
  /** `TreatyIn.Position` — workbasket tempat berkas menunggu (tombol Actions). */
  posisi: string
  /** `TreatyIn.PositionUsername` — "Position To". */
  pemegangPosisi: string
  /**
   * Properti yang layar pegang HANYA di penampung halaman, dibaca kembali
   * dari tabelnya (migrasi `448`): `RNMShareP`, `BrokeragePercentP`,
   * `OptionLimit`, `InstallmentNo`, `RevisionState`, `ViewState`. Disemai ke
   * penampung saat kontrak dimuat. Opsional: backend lama tidak mengirimnya.
   */
  penampung?: Record<string, string>
  /**
   * Larik TOTAL yang tab pegang di penampung halaman (`TotalShareRnmProp`,
   * `TotalSpreadedRnmProp`, `TotalSpreadedRnmRIProp`, `TotalEgnpiAmountNP`,
   * `TotalInstallmentNP`), dibaca dari `T_TREATY_TOTAL`. Larik tanpa baris
   * tidak dikirim. Opsional: backend lama tidak mengirimnya.
   */
  penampungLarik?: Record<string, Record<string, unknown>[]>
  nomorRujukan: string
  pemimpinTreaty: string
  adaDiJson: Record<string, boolean>

  /**
   * ⭐ EMPAT LARIK ISI — grid dan tab, dari dokumen warisan yang sama.
   *
   * Seluruhnya `string`: nilai di `M_TREATY_IN.JSONDATA` memang string
   * (terukur atas 150 dokumen), dan mengubahnya jadi angka atau tanggal
   * di jalur baca berarti menafsirkan sebelum ada yang memintanya.
   *
   * Larik kosong berarti dokumennya memang tidak punya — bukan galat, dan
   * bukan pula "belum dibangun".
   */
  kurs: BarisKursWarisan[]
  periodePelaporan: BarisPeriodeWarisan[]
  portofolio: BarisPortofolioWarisan[]
  akumulasi: BarisAkumulasiWarisan[]

  /** Empat tab berikutnya — dari tabel pendaratan yang sama. */
  egnpi: BarisEgnpiWarisan[]
  retensi: BarisRetensiWarisan[]
  /**
   * ⭐ Panel `Total Retention Amount` — DITURUNKAN dari `retensi` oleh
   * services, menurut `Activity/TreatyInNPSetTotal.xml` cabang
   * `param.type=retention`: jumlah `Amount` per mata uang, berurut
   * kemunculan pertama.
   */
  totalRetensi: BarisTotalRetensiWarisan[]
  /**
   * ⭐ Panel `Existing Policy for Master ID` — kanan atas, KEDUA cabang
   * (gambar 01 dan 26 dokumen desain).
   */
  polisProduksi: BarisPolisProduksi[]
  angsuran: BarisAngsuranWarisan[]
  catatan: BarisCatatanWarisan[]

  /**
   * ⭐ Baris LAYER dari `M_TREATY_IN2` — melayani EMPAT tab sekaligus:
   * Limits, Share, Event Limits, RNM Share.
   *
   * ⚠️ Kosong di sini punya DUA arti. Terukur 3 Oktober 2026: 1.850 kontrak
   * punya `Limits[]` berisi di dokumennya, tetapi `M_TREATY_IN2` hanya
   * mencakup 1.340 — jadi pada 510 kontrak (1.210 elemen limit) grid kosong
   * berarti *"sumbernya tidak mencakup kontrak ini"*, bukan *"kontrak ini
   * tidak punya limit"*. Layar menyatakan bedanya lewat petunjuk kosong.
   */
  layer: BarisLayerWarisan[]
  /** POHON tab Limits proporsional — `Limits[] → Detail[] → daftar`, apa adanya. */
  limitsPohon: SimpulLimit[]
  /** Larik akar tab Limits Non-Prop — `Summary of Limit`, `Total All Layers`. */
  limitsAkar?: LimitsAkar
  /** Tab Share Non-Prop — pendaratan `T_TREATY_SHARE*` + turunan (`share_np.go`). */
  shareNP?: ShareNP
  /** Pilihan dropdown kepala: nilai TERSIMPAN ↔ label tampil. */
  opsiKepala: OpsiKepala

  /** Tab Co-Ins Scale — 186 dari 1.854 kontrak, 702 baris. */
  skalaKoasuransi: BarisSkalaKoasuransiWarisan[]

  /**
   * Dua tab TEKS — Jalan B, keputusan §15: isinya tetap di `JSONDATA`.
   * Ejaannya sudah dipilih menurut cabang di services.
   */
  pengecualian: TabTeksWarisan
  syaratKhusus: TabTeksWarisan

  /**
   * Panel Attachment — dari `M_ATTACHMENTTREATY_2`, tabel WARISAN 43 baris.
   * Nol tabel baru.
   */
  lampiran: BarisLampiranWarisan[]
  kategoriLampiran: BarisKategoriLampiran[]
}

/** Satu berkas yang terlampir pada kontrak ini. */
export interface BarisLampiranWarisan {
  id: string
  kodeKategori: string
  namaKategori: string
  namaBerkas: string
  jenisMime: string
  idSimpanan: string
  diunggah: string
  pengunggah: string
}

/**
 * Satu baris panel Attachment: Category + Count.
 *
 * ⛔ `dipastikan` menyatakan apakah pasangan kode↔nama ini TERBUKTI. Empat
 * kode (`00003` `00004` `00008` `00009`) punya nol baris dan namanya tidak
 * ada di korpus; layar menandainya alih-alih menebak, sebab berkas yang
 * mendarat di kategori yang salah baru ketahuan bertahun kemudian.
 */
export interface BarisKategoriLampiran {
  kode: string
  nama: string
  cacah: number
  dipastikan: boolean
}

/** Isi panel Attachment satu kontrak — `GET /kontrak/{id}/lampiran`. */
export interface PanelLampiranAPI {
  lampiran: BarisLampiranWarisan[]
  kategoriLampiran: BarisKategoriLampiran[]
}

/** Nasib SATU berkas yang di-Attach. */
export interface HasilBerkasUnggah {
  nama: string
  berhasil: boolean
  pesan: string
}

/** Jawaban Attach: hasil per berkas, lalu panel yang dibaca ulang. */
export interface HasilUnggahLampiran extends PanelLampiranAPI {
  berkas: HasilBerkasUnggah[]
}

/**
 * Refresh panel Attachment (`GetMasterTreatyCategory_Act`) — panel saja;
 * isian form yang belum di-Save tidak ikut dibaca ulang.
 */
export async function ambilPanelLampiran(idKontrak: string): Promise<PanelLampiranAPI> {
  return minta<PanelLampiranAPI>(`${PREFIX_TREATYIN}/kontrak/${encodeURIComponent(idKontrak)}/lampiran`)
}

/**
 * Tombol Attach modal `ASM Attach Content` (`TreatySaveAttachment`): tiap
 * berkas dikirim ke Google Storage (`InsertGoogleStorage_Act`) lalu dicatat
 * di `T_STORAGE_IMAGE` + `M_ATTACHMENTTREATY_2`.
 */
/** `DownloadAttachmentTreaty` — URL bertanda tangan (`office` = View Office Online). */
export async function ambilTautanLampiran(idKontrak: string, idLampiran: string, office: boolean): Promise<{ url: string }> {
  return minta<{ url: string }>(
    `${PREFIX_TREATYIN}/kontrak/${encodeURIComponent(idKontrak)}/lampiran/${encodeURIComponent(idLampiran)}/tautan`,
    { kueri: office ? { office: '1' } : undefined },
  )
}

/** `Delete_act` — objek di Google Storage, lalu baris `M_ATTACHMENTTREATY_2`. */
export async function hapusLampiran(idKontrak: string, idLampiran: string): Promise<PanelLampiranAPI> {
  return minta<PanelLampiranAPI>(
    `${PREFIX_TREATYIN}/kontrak/${encodeURIComponent(idKontrak)}/lampiran/${encodeURIComponent(idLampiran)}/hapus`,
    { metode: 'POST', badan: {} },
  )
}

/** `ChangeDokument_Act("Save")` — `ChangeKateAttachment2_Sql` per baris. */
export async function ubahKategoriLampiran(
  idKontrak: string,
  perubahan: readonly { id: string; kategori: string }[],
): Promise<PanelLampiranAPI> {
  return minta<PanelLampiranAPI>(`${PREFIX_TREATYIN}/kontrak/${encodeURIComponent(idKontrak)}/lampiran/kategori`, {
    metode: 'POST',
    badan: { perubahan },
  })
}

export async function unggahLampiran(idKontrak: string, kodeKategori: string, berkas: readonly File[]): Promise<HasilUnggahLampiran> {
  const isi = new FormData()
  isi.append('kategori', kodeKategori)
  for (const f of berkas) isi.append('berkas', f, f.name)
  return mintaFormulir<HasilUnggahLampiran>(`${PREFIX_TREATYIN}/kontrak/${encodeURIComponent(idKontrak)}/lampiran`, isi)
}

/**
 * Satu tab yang isinya SATU medan teks panjang — sampai 23.453 aksara.
 *
 * ⛔ `ejaanLain` bukan hiasan. Ejaan `SpecialConditions*` BUKAN sinonim:
 * dari 303 dokumen yang punya lebih dari satu, NOL yang isinya identik.
 * Pembacanya berhak tahu ada teks lain yang tidak ia lihat.
 */
export interface TabTeksWarisan {
  /** Kosong bila ejaan cabangnya tidak ada — TIDAK diisi dari ejaan lain. */
  isi: string
  /** Ejaan yang dipakai; kosong bila nol teks terpilih. */
  ejaan: string
  /** Ejaan lain yang juga berisi. */
  ejaanLain: string[]
}

/**
 * Satu baris `M_TREATY_IN2` — satu LAYER.
 *
 * Seluruhnya `string`, termasuk yang di Oracle `NUMBER`: `NUMBER(22)`
 * membawa presisi yang `number` JavaScript tidak sanggup bawa, dan yang
 * memformat untuk layar adalah `formatNumber`/`formatPersen`.
 */
export interface BarisLayerWarisan {
  masterID: string
  sifatProporsi: string
  idCedant: string
  cedant: string
  idAsalBisnis: string
  asalBisnis: string
  kelompokTreaty: string
  /**
   * ⭐ Tingkat KETIGA pohon tab Limits — `Detail[].COBList[]`.
   * Larik kosong berarti treaty group itu memang nol kelas bisnis.
   */
  kelasBisnis: string[]
  namaKontrak: string
  tanggalMulai: string
  tanggalBerakhir: string
  jenisTreaty: string
  persenCession: string
  retensiCedant: string
  dasarCover: string
  jenisLayer: string
  layer: string
  jenisPenyebaran: string
  mataUang: string
  limit100: string
  adjRate: string
  premiEarned: string
  rasioMDP: string
  mdp: string
  rol: string
  relasiMataUang: string
  cessionKeRI: string
  epi100: string
  riogr: string
  persenBrokerage: string
  gempa: string
  /**
   * ⭐ Tiga batas Event Limits yang `M_TREATY_IN2` tidak punya, beserta
   * keempat mata uangnya — gambar 28 dokumen desain.
   */
  batasRSMD: string
  batasBanjirJab: string
  batasBanjirNas: string
  mataUangRSMD: string
  mataUangGempa: string
  mataUangBanjirJab: string
  mataUangBanjirNas: string
  rnmShare: string
  mataUangLimit: string
  liabilityRNM: string
  mdpRNM100: string
  qsor: string
  qsri: string
  liabilityQSRI: string
  liabilityQSOR: string
  epiRNMQS100: string
  rnmRetainedPremi: string
  rnmQSPremi: string
}

/** Satu baris tab Co-Ins Scale — `CoInShare` adalah PITA, bukan angka. */
export interface BarisSkalaKoasuransiWarisan {
  bagianKoasuransi: string
  persenLimit: string
  penyusun: string
  disusunPada: string
}

/** Satu baris tab EGNPI — larik `EGNPI` berisi di 846 dari 1.854 kontrak. */
export interface BarisEgnpiWarisan {
  jumlah: string
  jumlahIDR: string
  perTanggal: string
  kelasBisnis: string
  mataUang: string
  keterangan: string
  proporsi: string
  kelompokTreaty: string
}

/**
 * Satu baris panel `Existing Policy for Master ID`.
 *
 * ⛔ Dari `TREATYINPRODUCTION` lewat SQL rule
 * `FetchTreatyInProductionUsingNooffer`, dicocokkan atas TUJUH aksara
 * pertama `NOOFFER`.
 */
export interface BarisPolisProduksi {
  nomorPolis: string
  pegaID: string
  kuartal: string
  tahunKuartal: string
}

/**
 * Satu baris panel `Total Retention Amount`.
 *
 * ⛔ TURUNAN — nol kolom, nol kunci JSON, nol tabel memuatnya.
 */
export interface BarisTotalRetensiWarisan {
  /** Kolom pertama panel, yang judulnya adalah nama panelnya. */
  mataUang: string
  /** Kolom `Value` — jumlah, sebagai teks. */
  nilai: string
}

/** Satu baris tab Maximum Retention. */
export interface BarisRetensiWarisan {
  jumlah: string
  kelasBisnis: string
  mataUang: string
  keterangan: string
  kelompokTreaty: string
}

/** Satu baris JADWAL tab Installment — tabel anak, bukan induknya. */
export interface BarisAngsuranWarisan {
  angsuran: string
  mataUang: string
  jumlah: string
  persen: string
  jatuhTempo: string
  tanggalBayar: string
  wpc: string
  /** Bentuk TERSIMPAN (`YYYYMMDD`, tanggal WIB) — untuk penampung halaman. */
  jatuhTempoAsli: string
  tanggalBayarAsli: string
}

/** Satu baris tab Information & Submit. */
export interface BarisCatatanWarisan {
  tanggal: string
  operator: string
  disetujui: string
  catatan: string
}

/** Satu baris grid Rate of Exchange — dari `CurrencyList`, berisi di 297/300. */
export interface BarisKursWarisan {
  /** Pengenal baris `TREATYEXCHANGEYEARLY` — kunci Save; kosong = baris baru. */
  id?: string
  mataUang: string
  /** `.CurrencyID` — nilai dropdown sel Currency (`BrowseCurrency_RD` `.ID`). */
  mataUangID: string
  nilaiKeIDR: string
  /** Bentuk TAMPIL (`dd/mm/yy`) — untuk dibaca. */
  berlakuDari: string
  berlakuSampai: string
  /**
   * Bentuk TERSIMPAN (`YYYYMMDD`) — untuk kotak tanggal yang dapat diisi.
   *
   * ⛔ Kotak tanggal menerima bentuk kabel `DD-MM-YYYY`; `berlakuDari`
   * berupa `dd/mm/yy` (tahun DUA digit, garis miring) dan tidak terbaca
   * olehnya. Memakainya membuat kotaknya kosong tanpa satu pun galat — dan
   * Save lalu menulis kosong itu ke `STARTDATE`.
   */
  berlakuDariAsli: string
  berlakuSampaiAsli: string
}

// ---------------------------------------------------------------------------
// ⭐ TOMBOL TULIS — Save, Submit, Actions, Decline offer.
//
// Keputusan pemilik proses 6–7 Oktober 2026: isian masuk basis data HANYA
// lewat kedua fungsi ini, dan hanya dari tombolnya — tidak pernah dari
// `onChange` (dijaga `aturan-simpan.test.ts`). Sasarannya tabel masing-masing:
// `T_TREATY_*`, kepala `TREATY_IN`, kurs `TREATYEXCHANGEYEARLY`.
// ---------------------------------------------------------------------------

/** Grid Rate of Exchange yang Save kirim — baris yang BERUBAH atau BARU saja. */
export interface KursSimpan {
  tahun: string
  baris: BarisKursWarisan[]
}

/** Isi layar untuk tombol Save. */
export interface MasukanSimpan {
  /** Kosong = kontrak BARU (`Add`). */
  idKontrak: string
  /** Properti `TreatyIn` ejaan Pega — kepala, penampung halaman, Limits/Share Non-Prop. */
  dokumen: Record<string, unknown>
  kurs?: KursSimpan | null
}

/** Submit (`submit`), Actions (`akseptasi` + pilihan), Decline offer (`decline`). */
export interface MasukanKirim extends MasukanSimpan {
  aksi: 'submit' | 'akseptasi' | 'decline'
  /** `ChooseStatusAkseptasi` — Accept / Reject / Decline (Actions). */
  pilihan?: string
}

export interface HasilSimpan {
  id: string
  /** `ErrMsg` prosedur Pega: "Data Sudah Disimpan Dengan ID : …". */
  pesan: string
  posisi: string
  status: string
  pemegangPosisi: string
  /** Properti yang BELUM punya kolom di tabel pendaratan — dilaporkan, tidak ditelan. */
  kunciTakTersimpan: string[]
}

/** `POST /api/treaty-in/kontrak/simpan` — tombol Save. */
export async function simpanKontrak(m: MasukanSimpan): Promise<HasilSimpan> {
  return minta<HasilSimpan>(`${PREFIX_TREATYIN}/kontrak/simpan`, { metode: 'POST', badan: m })
}

/** `POST /api/treaty-in/kontrak/kirim` — Submit, Actions, Decline offer. */
export async function kirimKontrak(m: MasukanKirim): Promise<HasilSimpan> {
  return minta<HasilSimpan>(`${PREFIX_TREATYIN}/kontrak/kirim`, { metode: 'POST', badan: m })
}

/**
 * Tombol `Revision` daftar kontrak — `SetTreatyIn_Act(viewstate=1,
 * revisionstate=1)`: kontrak tuntas dibuka untuk direvisi dan DISIMPAN
 * seketika (RevisionState, ViewState, komentar "Create Revision").
 */
export async function mulaiRevisi(idKontrak: string): Promise<HasilSimpan> {
  return minta<HasilSimpan>(`${PREFIX_TREATYIN}/kontrak/revisi`, { metode: 'POST', badan: { idKontrak } })
}

/** Satu baris grid tab Reporting Period — dari `ReportingPeriodList`, 225/300. */
export interface BarisPeriodeWarisan {
  periode: string
  hitungOtomatis: string
  tanggalAwal: string
  jatuhTempoKirim: string
  jatuhTempoKonfirmasi: string
  jatuhTempoBayar: string
  /**
   * Bentuk TERSIMPAN (`YYYYMMDD`, tanggal WIB) — untuk KOTAK tanggal. Medan
   * di atas bentuk tampil, untuk dibaca. Terjemahan tampil yang dimasukkan
   * ke kotak tanggal membuat kotaknya kosong.
   */
  tanggalAwalAsli: string
  jatuhTempoKirimAsli: string
  jatuhTempoKonfirmasiAsli: string
  jatuhTempoBayarAsli: string
}

/** Isian kepala tab Reporting Period — dikirim tombol Apply. */
export interface MasukanPeriodePelaporan {
  mulai: string
  akhir: string
  periode: string
  interval: string
  penyerahan: string
  konfirmasi: string
  pelunasan: string
  /**
   * `param.startdate` — tanggal awal baris pertama. Kosong = `mulai`
   * (Apply); sel Initial Date mengirim `.InitialDate` barisnya.
   */
  awal?: string
}

/** Jawaban Apply: baris baru, atau pesan per medan (`Property-Set-Messages`). */
export interface HasilPeriodePelaporan {
  baris: BarisPeriodeWarisan[]
  galat: Record<string, string>
}

/**
 * `POST /api/treaty-in/hitung/periode-pelaporan` — tombol Apply,
 * `Activity/TreatyInSetReport.xml`. Rumusnya di services, bukan di sini.
 */
export async function hitungPeriodePelaporan(m: MasukanPeriodePelaporan): Promise<HasilPeriodePelaporan> {
  return minta<HasilPeriodePelaporan>(`${PREFIX_TREATYIN}/hitung/periode-pelaporan`, { metode: 'POST', badan: m })
}

/**
 * Masukan `LimitCalculation` — parameter Activity + simpul + pohonnya.
 *
 * ⛔ Simpul dikirim APA ADANYA: kuncinya ejaan dokumen (`QSPct`,
 * `IOOLimitList`, …), dan services hanya membaca yang ia perlukan.
 */
export interface MasukanLimit {
  /** `param.kindoftreaty`. */
  jenis: 'qs' | 'surplus'
  /** `param.add` — `man` atau kosong. */
  tambah: '' | 'man'
  /** `param.autocalculate`. */
  otomatis: boolean
  detail: SimpulLimit
  /** Seluruh `Detail` di seluruh `Limits`, URUT pohon — langkah 9. */
  pohon: readonly SimpulLimit[]
}

/** `POST /api/treaty-in/hitung/limit` — `Activity/LimitCalculation.xml`. */
export async function hitungLimit(m: MasukanLimit): Promise<SimpulLimit> {
  return minta<SimpulLimit>(`${PREFIX_TREATYIN}/hitung/limit`, { metode: 'POST', badan: m })
}

/** Satu baris tab Portfolio — dari `Portfolio`, 152/300. */
export interface BarisPortofolioWarisan {
  jenis: string
  jenisPortfolio: string
  keterangan: string
}

/** Satu baris tab Accumulation — dari `AccumulationList`, 12/300. */
export interface BarisAkumulasiWarisan {
  periode: string
  tanggalLapor: string
  hariKirim: string
  jatuhTempoKirim: string
  /** Bentuk TERSIMPAN (`YYYYMMDD`, tanggal WIB) — untuk penampung halaman. */
  tanggalLaporAsli: string
  jatuhTempoKirimAsli: string
}

/** `GET /api/treaty-in/kontrak-warisan/{id}` - satu kontrak warisan. */
export async function ambilKontrakWarisan(id: string): Promise<KontrakWarisan> {
  return minta<KontrakWarisan>(`${PREFIX_TREATYIN}/kontrak-warisan/${encodeURIComponent(id)}`)
}

// ===========================================================================
// Isi pemilih "Choose Ceding" dan "Choose Source of Business"
//
// ⛔ Keduanya KATALOG — nol parameter kontrak. Isinya sama untuk setiap
// kontrak, jadi ia diambil sekali ketika pemilihnya dibuka, bukan ikut
// menumpang pada pembacaan kontrak.
// ===========================================================================

/** Satu baris pemilih: pengenal yang DITULIS, beserta namanya. */
export interface PilihanWarisan {
  id: string
  nama: string
  /**
   * ⚠️ Benar bila nama ini dipakai LEBIH DARI SATU pengenal. Terukur
   * 4 Oktober 2026: `REASURANSI MAIPARK INDONESIA` punya 3, tujuh nama lain
   * punya 2. Layar WAJIB menandainya — pemilih yang menampilkan dua baris
   * bernama sama tanpa penjelasan terbaca sebagai kesalahan layar, dan yang
   * memilih akan menebak.
   */
  kembar: boolean
  /** `SOANOTE` — kolom SOA Name menu Reinsurance Type (pilihan Treaty Type). */
  namaSoa?: string
}

/** `GET /api/treaty-in/warisan/cedant` - 131 baris, urut nama lalu pengenal. */
export async function ambilDaftarCedant(): Promise<PilihanWarisan[]> {
  return minta<PilihanWarisan[]>(`${PREFIX_TREATYIN}/warisan/cedant`)
}

/** `GET /api/treaty-in/warisan/asal-bisnis` - 96 baris, urutan sama. */
export async function ambilDaftarAsalBisnis(): Promise<PilihanWarisan[]> {
  return minta<PilihanWarisan[]>(`${PREFIX_TREATYIN}/warisan/asal-bisnis`)
}

/**
 * `GET /api/treaty-in/warisan/jenis-treaty` — isi dropdown `Treaty Type` tab
 * Limits: `REINSURANCETYPE` lewat `BrowseReinsuranceType_RD` (Flag `active`,
 * ID menurun). `id` = `.TreatyTypeID`, `nama` = `.Note`.
 */
export async function ambilDaftarJenisTreaty(): Promise<PilihanWarisan[]> {
  return minta<PilihanWarisan[]>(`${PREFIX_TREATYIN}/warisan/jenis-treaty`)
}

/**
 * Isi SELURUH dropdown tab Limits proporsional — tiga RD Pega:
 * `jenisTreaty` (`BrowseReinsuranceType_RD`, REINSURANCETYPE),
 * `kelompokTreaty` (`BrowseTreatyGroup_RD`, TREATYGROUP: `id` = `.ID`,
 * `nama` = `.TreatyGroupName`), `mataUang` (`BrowseCurrencyTreatyIn_RD` /
 * `BrowseCurrency_RD`, CURRENCY tanpa ITL: `id` = `.ID`, `nama` = `.Currency`).
 */
export interface OpsiLimits {
  jenisTreaty: PilihanWarisan[]
  kelompokTreaty: PilihanWarisan[]
  mataUang: PilihanWarisan[]
  /**
   * Dropdown `associated` tab Limits Non-Prop — daftarnya di rule Property
   * yang tidak diekspor; nilainya domain TERUKUR di data, label apa adanya.
   */
  jenisLayer?: OpsiPilihan[]
  cover?: OpsiPilihan[]
  relasiMataUang?: OpsiPilihan[]
  catatanReinstatement?: OpsiPilihan[]
}

/** Bagian AKAR dokumen yang tab Limits Non-Prop tampilkan. */
export interface LimitsAkar {
  LimitSummaryList: SimpulLimit[]
  /** `TotalLimitIOONP`, `TotalLimitDeductblNP`, `TotalLimitPremiEarnNP`, `TotalLimitMDPNP`. */
  Total: Record<string, SimpulLimit[]>
  TotalLimitsROL: string
}

/** Aksi rumus tab Limits Non-Prop — `services/hitung_limit_np.go`. */
export type AksiLimitNP =
  | 'egnpi'
  | 'reinstatement'
  | 'reinst-jumlah'
  | 'reinst-persen'
  | 'reinst-tambahan'
  | 'adj'
  | 'mdp'
  | 'total'
  | 'nilai-list'

/** Satu baris EGNPI yang `TotalEgnpi` baca. */
export interface EgnpiLimitNP {
  TreatyGroup: string
  Currency: string
  CurrencyID: string
  Amount: string
}

/** Satu baris Rate of Exchange yang `DetailCalculationROL` baca. */
export interface KursLimitNP {
  Currency: string
  Conversion: string
}

export interface MasukanLimitNP {
  aksi: AksiLimitNP
  /** Seluruh `Limits[]`, simpul APA ADANYA. */
  layers: readonly SimpulLimit[]
  egnpi: readonly EgnpiLimitNP[]
  kurs: readonly KursLimitNP[]
  /** Indeks layer dan baris reinstatement — mulai 0. */
  indeks: number
  baris: number
}

export interface HasilLimitNP {
  layers: SimpulLimit[]
  Total?: Record<string, SimpulLimit[]>
  TotalLimitsROL?: string
  LimitSummaryList?: SimpulLimit[]
  pesan: string[]
}

/** `POST /api/treaty-in/hitung/limit-np` — Activity tab Limits Non-Prop. */
export async function hitungLimitNP(m: MasukanLimitNP): Promise<HasilLimitNP> {
  return minta<HasilLimitNP>(`${PREFIX_TREATYIN}/hitung/limit-np`, { metode: 'POST', badan: m })
}

// ---------------------------------------------------------------------
// Tab EGNPI (non-prop) — `SetAmountConversion`, `TreatyInEGNPIListValue`,
// `TreatyInNPSetTotal(egnpi)`, `TreatyInNonAddItem(egnpi)`.
// ---------------------------------------------------------------------

/**
 * Satu baris tab EGNPI.
 *
 * ⭐ MELUASKAN `EgnpiLimitNP`, bukan menggantikannya: tab Limits mengirim
 * empat medan yang sama ke rumusnya sendiri, dan satu baris layar yang
 * dijelaskan dua tipe adalah cara termudah keduanya berbeda diam-diam.
 */
export interface BarisEgnpi extends EgnpiLimitNP {
  ID: string
  TreatyGroupID: string
  AsDate: string
  Proportion: string
  AmountIDR: string
  ClassOfBusiness: string
  ClassOfBusinessID: string
  Note: string
}

export type AksiEgnpi = 'konversi' | 'nilai' | 'total' | 'tambah' | 'hapus'

export interface MasukanEgnpi {
  aksi: AksiEgnpi
  egnpi: readonly BarisEgnpi[]
  kurs: readonly KursLimitNP[]
  /** `TreatyIn.Retention` — hanya baris PERTAMA yang dipakai, oleh `tambah`. */
  retensi: readonly { Currency: string; CurrencyID: string }[]
  /** Indeks baris (mulai 0) untuk `konversi` dan `hapus`. */
  indeks: number
}

export interface HasilEgnpi {
  egnpi: BarisEgnpi[]
  TotalEgnpiAmount: string
  TotalEgnpiProportion: string
  TotalEgnpiAmountNP: { Currency: string; CurrencyID: string; Value: string }[]
  pesan: string[]
}

/** `POST /api/treaty-in/hitung/egnpi` — rumus tab EGNPI, nol tulis basis data. */
export async function hitungEgnpi(m: MasukanEgnpi): Promise<HasilEgnpi> {
  return minta<HasilEgnpi>(`${PREFIX_TREATYIN}/hitung/egnpi`, { metode: 'POST', badan: m })
}

// ---------------------------------------------------------------------
// Tab Maximum Retention (non-prop) — `TreatyInNPSetTotal(retention)` dan
// `TreatyInNonAddItem(retention)`.
// ---------------------------------------------------------------------

/**
 * Satu baris tab Maximum Retention.
 *
 * ⚠️ `Note` hanya terlihat di RINCIAN baris dan `ClassOfBusiness` nol
 * terlihat di mana pun — keduanya tetap dibawa: medan yang dibuang di
 * perjalanan tidak dapat dikembalikan saat Save.
 */
export interface BarisRetensi {
  ID: string
  TreatyGroup: string
  TreatyGroupID: string
  Currency: string
  CurrencyID: string
  Amount: string
  ClassOfBusiness: string
  ClassOfBusinessID: string
  Note: string
}

export type AksiRetensi = 'total' | 'tambah' | 'hapus'

export interface MasukanRetensi {
  aksi: AksiRetensi
  retensi: readonly BarisRetensi[]
  /** Indeks baris (mulai 0) untuk `hapus`. */
  indeks: number
}

export interface HasilRetensi {
  retensi: BarisRetensi[]
  TotalRetentionAmountNP: { Currency: string; CurrencyID: string; Value: string }[]
  pesan: string[]
}

/** `POST /api/treaty-in/hitung/retensi` — rumus tab Maximum Retention. */
export async function hitungRetensi(m: MasukanRetensi): Promise<HasilRetensi> {
  return minta<HasilRetensi>(`${PREFIX_TREATYIN}/hitung/retensi`, { metode: 'POST', badan: m })
}

export interface HasilDeduksi {
  DeductionList: SimpulLimit[]
  DeductionTotalList: SimpulLimit[]
  pesan: string[]
}

/** `POST /api/treaty-in/hitung/limit-deduksi` — `Activity/CalculateDeduction.xml`. */
export async function hitungDeduksi(m: {
  sts: '' | 'val' | 'pct'
  indeks: number
  DeductionList: readonly SimpulLimit[]
  GrossPremiumList: readonly SimpulLimit[]
}): Promise<HasilDeduksi> {
  return minta<HasilDeduksi>(`${PREFIX_TREATYIN}/hitung/limit-deduksi`, { metode: 'POST', badan: m })
}

export interface HasilCadangan {
  ReserveList: SimpulLimit[]
  pesan: string[]
}

/** `POST /api/treaty-in/hitung/limit-cadangan` — `Activity/PremiumReserveCalculate.xml`. */
export async function hitungCadangan(m: {
  PremiumReservePct: string
  CessionList: readonly SimpulLimit[]
}): Promise<HasilCadangan> {
  return minta<HasilCadangan>(`${PREFIX_TREATYIN}/hitung/limit-cadangan`, { metode: 'POST', badan: m })
}

/** Masukan `GetAchievement` — tombol Refresh / pilihan Quarter Year. */
export interface MasukanAchievement {
  idKontrak: string
  /** `param.search == "search"` — Quarter Year dipilih. */
  cari: boolean
  asAt: string
  tahun: string
  /** `TreatyIn.RNMShareP` — kosong bila tidak diketahui. */
  rnmShareP: string
  limits: readonly { TreatyType: string; Detail: readonly { TreatyGroup: string; EPIList: readonly SimpulLimit[] }[] }[]
}

/** Satu baris grid tab `Achievement In IDR` (`AchievementCombine.xml`). */
export interface BarisRingkasAchievement {
  TreatyGroup: string
  TreatyType: string
  TotalAchPremium: string
  TotalAchNetPremium: string
  TotalAchIncured: string
  LossRatio: string
}

/** Isi tab `Achievement In IDR`: grid + TIGA sel kaki. */
export interface RingkasanAchievement {
  baris: BarisRingkasAchievement[]
  SumTotalAchievNetPremium: string
  SumTotalAchievIncured: string
  SumLossRatio: string
}

export interface HasilAchievement {
  /** Per Kind of Treaty, per Treaty Group — medan yang Activity tulis. */
  limits: SimpulLimit[][]
  kuartal: string[]
  tahunKuartal: string[]
  flagExcel: boolean
  /** Tab `Achievement In IDR` — proyeksi atas `limits` di atas. */
  ringkasan: RingkasanAchievement
}

/** `POST /api/treaty-in/hitung/achievement` — `Activity/GetAchievement.xml`. */
export async function hitungAchievement(m: MasukanAchievement): Promise<HasilAchievement> {
  return minta<HasilAchievement>(`${PREFIX_TREATYIN}/hitung/achievement`, { metode: 'POST', badan: m })
}

/**
 * `GET /api/treaty-in/warisan/kelas-bisnis?treatyGroupId=` — autocomplete
 * `Class of Business` (`BrowseTreatyBusinessWOType_RD`): `id` = `.BizCode`,
 * `nama` = `.BIZNAME`.
 */
export async function ambilKelasBisnis(treatyGroupId: string): Promise<PilihanWarisan[]> {
  return minta<PilihanWarisan[]>(
    `${PREFIX_TREATYIN}/warisan/kelas-bisnis?treatyGroupId=${encodeURIComponent(treatyGroupId)}`,
  )
}

/** `GET /api/treaty-in/warisan/opsi-limits`. */
export async function ambilOpsiLimits(): Promise<OpsiLimits> {
  return minta<OpsiLimits>(`${PREFIX_TREATYIN}/warisan/opsi-limits`)
}

/**
 * Satu simpul pohon Limits: medan TEKS apa adanya, atau larik simpul.
 *
 * ⛔ Kunci yang TIDAK ADA di dokumen tidak muncul — `kunci in simpul` yang
 * salah berarti "tidak ada di sistem lama", bukan "kosong".
 */
export interface SimpulLimit {
  [kunci: string]: string | SimpulLimit[]
}

/** Satu pilihan dropdown — nilai TERSIMPAN dan labelnya. */
export interface OpsiPilihan {
  value: string
  label: string
}

/** Pilihan ketiga dropdown kepala. Labelnya disusun services, bukan di sini. */
export interface OpsiKepala {
  bordereaux: OpsiPilihan[]
  caraPembukuan: OpsiPilihan[]
  caraPembukuanNonProp: OpsiPilihan[]
  periodePelaporan: OpsiPilihan[]
}

/** `GET /api/treaty-in/warisan/opsi-kepala` — juga untuk kontrak BARU. */
export async function ambilOpsiKepala(): Promise<OpsiKepala> {
  return minta<OpsiKepala>(`${PREFIX_TREATYIN}/warisan/opsi-kepala`)
}

// ---------------------------------------------------------------------
// Tab Share NON-PROPORSIONAL — `hitung_share_np.go`. Ejaan medan = ejaan
// dokumen Pega, seperti `SimpulLimit`.
// ---------------------------------------------------------------------

/** Satu baris larik bernilai per mata uang (`RnmLimitList`, `Total…`). */
export interface NilaiShare {
  Currency: string
  CurrencyID: string
  Value: string
}

export interface BarisDeduksiShare {
  Comment: string
  Currency: string
  CurrencyID: string
  Deduction: string
  DeductionPct: string
  /** `Auto Calculate %` — `'true'` bila dicentang. */
  DeductionPctCalculate?: string
}

export interface BarisSpreadingShare {
  ReinsTypeName: string
  ReinsTypeID: string
  ParentReinsTypeID: string
  Pct: string
  Rp: string
  Usd: string
  RnmLimitList: NilaiShare[]
  GrossPremiumList: NilaiShare[]
  GrossPremiumMinList: NilaiShare[]
  DeductionTotalList: NilaiShare[]
  NetPremiumList: NilaiShare[]
}

/** Satu baris `TreatyIn.Share` — SATU per layer Limits. */
export interface BarisShareNP {
  LayerType: string
  Layer: string
  LayerPartType: string
  LayerPart: string
  Cover: string
  RNMShare: string
  TreatyGroupList: { TreatyGroup: string; TreatyGroupID: string }[]
  Limit: string
  Limit2: string
  SpreadingTypeXOL: string
  SpreadingTypeIDXOL: string
  SpreadingTotalPctXOL: string
  SpreadingListXOL: BarisSpreadingShare[]
  DeductionList: BarisDeduksiShare[]
  RnmLimitList: NilaiShare[]
  GrossPremiumList: NilaiShare[]
  GrossPremiumMinList: NilaiShare[]
  DeductionTotalList: NilaiShare[]
  NetPremiumList: NilaiShare[]
  RNMSpreadedListXOL: NilaiShare[]
  RNMSpreadedListRIXOL: NilaiShare[]
  RNMSpreadedListGrossXOL: NilaiShare[]
  RNMSpreadedListGrossRIXOL: NilaiShare[]
  RNMSpreadedListGrossMinXOL: NilaiShare[]
  RNMSpreadedListGrossRIMinXOL: NilaiShare[]
  RNMSpreadedListDeductXOL: NilaiShare[]
  RNMSpreadedListDeductRIXOL: NilaiShare[]
  RNMSpreadedListNetXOL: NilaiShare[]
  RNMSpreadedListNetRIXOL: NilaiShare[]
}

/** Satu baris grid `Reinsurer Name` / `Facultative Reinsurers`. */
export interface BarisReinsShare {
  ID: string
  ReinsID: string
  ReinsName: string
  Layer: string
  SharePct: string
}

/** Satu baris "Summarry of RNM Share". */
export interface RingkasanShareNP {
  LayerType: string
  Layer: string
  LayerPartType: string
  LayerPart: string
  Note: string
  Limit: string
  Limit2: string
  MDP: string
  MDP2: string
  Deductible: string
  Deductible2: string
  NetPremi: string
  NetPremi2: string
}

export interface ShareNP {
  RNMShare: string
  BrokeragePercent: string
  RNMShareAcrossTheBoard: string
  FacultativeShare: string
  FacultativeShareBrokerage: string
  RnmShareDeducted: string
  IsProRate: string
  ShareReins: BarisReinsShare[]
  ShareFacultativeReinsurers: BarisReinsShare[]
  Share: BarisShareNP[]
  FacultativeShareList: BarisShareNP[]
  LimitShareSummaryList: RingkasanShareNP[]
  LimitFacShareSummaryList: RingkasanShareNP[]
  /** Kesembilan grid "Total All Layers RNM Share" — kunci nama properti. */
  Total: Record<string, NilaiShare[]>
}

/** Aksi = satu tombol / isian ekspor (`AksiShare…` di services). */
export type AksiShareNP =
  | 'rnm'
  | 'brokerage'
  /** `TreatyInSetBrokerage` saja — langkah bersyarat kotak centang Across The Board. */
  | 'set-brokerage'
  | 'centang'
  | 'fac'
  | 'summary'
  | 'total'
  | 'nilai-share'
  | 'spreading-type'
  | 'rnm-baris'
  | 'deduksi'
  | 'spreading-tambah'
  | 'spreading-hapus'
  | 'spreading-pct'

export interface MasukanShareNP {
  aksi: AksiShareNP
  share: ShareNP
  /** Layer Limits Non-Prop — sumber Update Summary. */
  layers: readonly SimpulLimit[]
  indeks?: number
  baris?: number
  sts?: string
  idKontrak?: string
  /** `TreatyIn.Commencement`, YYYYMMDD. */
  commencement?: string
}

export interface HasilShareNP {
  share: ShareNP
  /** Pesan tingkat tab. */
  pesan: string[]
  /** Pesan yang Pega tempelkan pada medan SATU baris Share (panel rinciannya). */
  pesanBaris: { indeks: number; pesan: string }[]
}

/** Satu induk / anak susunan treaty master (`PROPORTIONALARRG`). */
export interface SusunanSpreading {
  reinsTypeId: string
  reinsTypeName: string
  parentReinsTypeId: string
  treatyYearId: string
  pct: string
  rp: string
  usd: string
}

/** `POST /api/treaty-in/hitung/share-np`. */
export async function hitungShareNP(m: MasukanShareNP): Promise<HasilShareNP> {
  return minta<HasilShareNP>(`${PREFIX_TREATYIN}/hitung/share-np`, { metode: 'POST', badan: m })
}

// ---------------------------------------------------------------------
// Penampung halaman `TreatyIn` cabang PROPORSIONAL — baris ejaan Pega
// (`halaman.tsx`). Rumus tab Accumulation dan Share Prop di services.
// ---------------------------------------------------------------------

/** Satu baris `TreatyIn.AccumulationList`. Tanggal bentuk simpan `YYYYMMDD`. */
export interface BarisAkumulasi {
  Period: string
  ReportDate: string
  SubDays: string
  SubDueDate: string
}

/** `TreatyInSetAccountReport` (`periode`) / `TreatyInAccumulationSetSubDue` (`jatuh-tempo`). */
export interface MasukanAkumulasi {
  aksi: 'periode' | 'jatuh-tempo'
  AccumulationPeriod: string
  /** Milik tab Reporting Period — dibaca dari penampung halaman. */
  ReportingStart: string
  ReportingEnd: string
  AccumulationList: readonly BarisAkumulasi[]
}

/** `POST /api/treaty-in/hitung/akumulasi`. */
export async function hitungAkumulasi(m: MasukanAkumulasi): Promise<{ AccumulationList: BarisAkumulasi[] }> {
  return minta<{ AccumulationList: BarisAkumulasi[] }>(`${PREFIX_TREATYIN}/hitung/akumulasi`, { metode: 'POST', badan: m })
}

/** Satu baris total per mata uang (`TotalShareRnmProp` …). */
export interface NilaiTotalProp {
  Currency: string
  CurrencyID: string
  Value: string
}

/** Satu aksi tab Share Prop — `Limits` = page list yang tab Limits isi. */
export interface MasukanShareProp {
  aksi: 'share' | 'detail' | 'spreading' | 'sebar-nama'
  Limits: readonly SimpulLimit[]
  RNMShareP: string
  BrokeragePercentP: string
  OptionLimit: string
  RNMShareAcrossTheBoard: string
  /** `TreatyIn.Commencement` (`YYYYMMDD`). */
  Commencement: string
  TotalShareRnmProp: readonly NilaiTotalProp[]
  TotalSpreadedRnmProp: readonly NilaiTotalProp[]
  TotalSpreadedRnmRIProp: readonly NilaiTotalProp[]
  indeksLimit: number
  indeksDetail: number
}

export interface HasilShareProp {
  Limits: SimpulLimit[]
  TotalShareRnmProp: NilaiTotalProp[]
  TotalSpreadedRnmProp: NilaiTotalProp[]
  TotalSpreadedRnmRIProp: NilaiTotalProp[]
  pesan: string[]
}

/** `POST /api/treaty-in/hitung/share-prop`. */
export async function hitungShareProp(m: MasukanShareProp): Promise<HasilShareProp> {
  return minta<HasilShareProp>(`${PREFIX_TREATYIN}/hitung/share-prop`, { metode: 'POST', badan: m })
}

/** Satu baris `TreatyIn.Installment(n).InstallmentList` (ejaan Pega, tanggal `YYYYMMDD`). */
export interface BarisAngsuran {
  Installment: string
  DueDate: string
  WPC: string
  PaymentDate: string
  Currency: string
  InstallmentPct: string
  Amount: string
}

/** Satu halaman `TreatyIn.Installment` — satu mata uang. */
export interface Angsuran {
  Currency: string
  AmountTotal: string
  PctTotal: string
  InstallmentList: BarisAngsuran[]
}

export interface MasukanAngsuran {
  /** `tanggal-bayar` = TreatyInUpdatePaymentDate (satu halaman); `tanggal-bayar-semua` = `_Act`. */
  aksi: 'nilai' | 'total-baris' | 'total' | 'tanggal-bayar' | 'tanggal-bayar-semua'
  /** `update` (Update Value), `editpercentage` (sel % Installment), atau kosong. */
  status: string
  installmentNo: string
  edmState: string
  angsuran: readonly Angsuran[]
  /** `TreatyIn.OLDDATA.Installment` — kosong di Treaty In. */
  angsuranLama: readonly Angsuran[]
  /** `TreatyIn.TotalShareNetNP` — dari tab Share Non-Prop. */
  netPremium: readonly NilaiShare[]
  indeks: number
}

export interface HasilAngsuran {
  angsuran: Angsuran[]
  /** `null` bila aksinya tidak menyentuh total (`total-baris`). */
  TotalInstallmentNP: NilaiShare[] | null
  InstallmentNo: string
  pesan: string[]
}

/**
 * `POST /api/treaty-in/hitung/angsuran` — `TreatyInSetValueInstallment`,
 * `SetTotalInstallment`, `TreatyInNPSetTotal(installment)`.
 */
export async function hitungAngsuran(m: MasukanAngsuran): Promise<HasilAngsuran> {
  return minta<HasilAngsuran>(`${PREFIX_TREATYIN}/hitung/angsuran`, { metode: 'POST', badan: m })
}

/** `GET /api/treaty-in/warisan/spreading-induk` — dropdown `Spreading Type`. */
export async function ambilIndukSpreading(treatyGroupId: string, mulai: string): Promise<SusunanSpreading[]> {
  return minta<SusunanSpreading[]>(
    `${PREFIX_TREATYIN}/warisan/spreading-induk?treatyGroupId=${encodeURIComponent(treatyGroupId)}&mulai=${encodeURIComponent(mulai)}`,
  )
}

/** `GET /api/treaty-in/warisan/reasuradur-share` — autocomplete `Reinsurer Name`. */
export async function ambilReasuradurShare(): Promise<PilihanWarisan[]> {
  return minta<PilihanWarisan[]>(`${PREFIX_TREATYIN}/warisan/reasuradur-share`)
}
