// modul/treatyin/frontend/api.ts - panggilan backend modul Treaty In, satu
// fungsi per rute (`backend/handlers/rute_treaty_in.go`). Klien HTTP-nya
// `inti/frontend/klien.ts`.
//
// ⛔ SATU rute, dan itu disengaja. Papan tiket modul ini menyatakan `L-4`:
// "tidak ada spesifikasi layar di mana pun". Rute kontrak, versi, dan
// persetujuan lahir bersama spesifikasinya.

import { minta } from '../../../inti/frontend/klien'

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
  mataUang: string
  nilaiKeIDR: string
  berlakuDari: string
  berlakuSampai: string
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
