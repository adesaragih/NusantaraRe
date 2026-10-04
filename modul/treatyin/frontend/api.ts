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

  /** Tab Co-Ins Scale — 186 dari 1.854 kontrak, 702 baris. */
  skalaKoasuransi: BarisSkalaKoasuransiWarisan[]

  /**
   * Dua tab TEKS — Jalan B, keputusan §15: isinya tetap di `JSONDATA`.
   * Ejaannya sudah dipilih menurut cabang di services.
   */
  pengecualian: TabTeksWarisan
  syaratKhusus: TabTeksWarisan
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
