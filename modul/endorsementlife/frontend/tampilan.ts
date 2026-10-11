// Pembantu tampilan MURNI modul Endorsement Life.
//
// ⛔ Angka tidak pernah menjadi `Number`: sel angka lewat `formatNumber` (digit, bukan float).

import { DESIMAL_TAK_DIBATASI, formatNumber } from '../../../inti/frontend/lib/format'
import type { KasusEDM } from './api'
import type { KolomKorpus } from './kolomKorpus'
import { BUAT_EDM, GRID_EDM, INBOX_EDM, KASUS_EDM, OPSI_EDM_TYPE, OPSI_METODE_PREMI, OPSI_TYPE_CEDING, PERBANDINGAN_EDM, PROSES_EDM, RINCI_EDM } from './labels'

/** Baris per halaman grid - SAMA dengan `services.UkuranHalaman` (`pyGridPaginator`). */
export const UKURAN_HALAMAN_EDM = 20

/** Sel kosong ditandai, bukan dibiarkan kosong (ADR-U-0027). */
export function sel(v: string | undefined): string {
  const t = (v ?? '').trim()
  return t === '' ? '—' : t
}

/** Tanggal (dan jam bila ada) terurai dari teks server - olah teks, tanpa `Date` (nol geser zona waktu). */
interface TanggalTerurai {
  yyyy: string
  mm: string
  dd: string
  jam: string
}

/**
 * Bentuk yang dikenal: ISO `YYYY-MM-DD` (+ `[ T]HH:MM:SS`), stempel Pega `YYYYMMDD[THHMMSS]`, dan
 * `DD-MM-YYYY` / `DD/MM/YYYY`. Tahun < 1900 (waktu nol Go `0001-01-01`) = tidak ada nilai.
 */
function uraiTanggal(v: string | undefined): TanggalTerurai | null {
  const t = (v ?? '').trim()
  let m = /^(\d{4})-(\d{2})-(\d{2})(?:[ T](\d{2}:\d{2}:\d{2}))?/.exec(t)
  let hasil: TanggalTerurai | null = null
  if (m) {
    hasil = { yyyy: m[1] ?? '', mm: m[2] ?? '', dd: m[3] ?? '', jam: m[4] ?? '' }
  } else if ((m = /^(\d{4})(\d{2})(\d{2})(?:T(\d{2})(\d{2})(\d{2}))?/.exec(t))) {
    hasil = { yyyy: m[1] ?? '', mm: m[2] ?? '', dd: m[3] ?? '', jam: m[4] ? `${m[4]}:${m[5] ?? ''}:${m[6] ?? ''}` : '' }
  } else if ((m = /^(\d{2})[-/](\d{2})[-/](\d{4})$/.exec(t))) {
    hasil = { yyyy: m[3] ?? '', mm: m[2] ?? '', dd: m[1] ?? '', jam: '' }
  }
  if (hasil === null || hasil.yyyy < '1900') return null
  return hasil
}

/**
 * Sel tanggal → `dd/mm/yyyy` TANPA jam (keputusan work owner 07-10-2026; sama dengan kotak masuk PremiumList
 * Life). Kosong / tak terurai → tanda kosong.
 */
export function selTanggal(v: string | undefined): string {
  const t = uraiTanggal(v)
  return t === null ? sel('') : `${t.dd}/${t.mm}/${t.yyyy}`
}

/**
 * Sel STEMPEL WAKTU (riwayat keputusan) → `dd/mm/yyyy HH:MM:SS`; tanpa jam → `dd/mm/yyyy`.
 * Kosong / tak terurai → tanda kosong.
 */
export function selWaktu(v: string | undefined): string {
  const t = uraiTanggal(v)
  if (t === null) return sel('')
  return t.jam === '' ? `${t.dd}/${t.mm}/${t.yyyy}` : `${t.dd}/${t.mm}/${t.yyyy} ${t.jam}`
}

/** Sel angka - digit apa adanya, pemisah ribuan, nol ekor dipangkas; negatif tetap bertanda. */
export function selAngka(v: string | undefined): string {
  const t = (v ?? '').trim()
  return t === '' ? '—' : formatNumber(t, DESIMAL_TAK_DIBATASI)
}

/** Label `EDM Type` (`1` Perubahan Data, `3` Batal); nilai lain apa adanya. */
export function labelEdmType(v: string): string {
  return OPSI_EDM_TYPE[v] ?? sel(v)
}

/**
 * Sub-jenis EDM Type baca-saja di kepala kasus (07-10-2026, migrasi 483): hanya yang sesuai EDM Type kasus.
 * Kasus sebelum 483 (kolom kosong) tampil bertanda kosong, bukan disembunyikan.
 */
export function subJenisKasus(k: {
  edmType: string
  edmTypePerubahanData?: string
  edmTypeBatal?: string
}): ReadonlyArray<readonly [string, string]> {
  if (k.edmType === '1') return [[BUAT_EDM.edmTypePerubahanData, sel(k.edmTypePerubahanData)]]
  if (k.edmType === '3') return [[BUAT_EDM.edmTypeBatal, sel(k.edmTypeBatal)]]
  return []
}

/** Sub-jenis EDM Type kasus (Perubahan Data / Batal sesuai EDM Type); '' bila kosong atau EDM Type lain. */
export function subJenisTerisi(k: { edmType: string; edmTypePerubahanData?: string; edmTypeBatal?: string }): string {
  const v = k.edmType === '1' ? k.edmTypePerubahanData : k.edmType === '3' ? k.edmTypeBatal : ''
  return (v ?? '').trim()
}

/** Nama `Premium Method` dari kode `PRO_RATE_TYPE` (1 Single, 2 Annually, 3 Others); lain apa adanya. */
export function labelMetodePremi(v: string): string {
  return OPSI_METODE_PREMI[v.trim()] ?? sel(v)
}

/** Satu medan kepala kasus: label, nilai tampil, dan apakah selebar bagiannya. */
export type MedanKepala = readonly [label: string, nilai: string, lebar: boolean]

/**
 * Kepala layar kasus dikelompokkan (keputusan work owner 07-10-2026, pola Premium List Detail): Policy Data,
 * Parties, Endorsement. Sub-jenis hanya yang sesuai EDM Type (`subJenisKasus`); Description selebar bagian.
 */
export function bagianKepalaKasus(
  kasus: Pick<KasusEDM, 'policyNo' | 'edmType' | 'edmTypePerubahanData' | 'edmTypeBatal' | 'description' | 'edmDate' | 'kepala'>,
): ReadonlyArray<{ judul: string; medan: readonly MedanKepala[] }> {
  const k = kasus.kepala
  return [
    {
      judul: KASUS_EDM.bagianPolis,
      medan: [
        [BUAT_EDM.policyNo, sel(kasus.policyNo), false],
        [KASUS_EDM.productName, sel(k.PRODUCT_NAME), true],
        [KASUS_EDM.productNameId, sel(k.PRODUCT_NAME_ID), false],
        [KASUS_EDM.type, sel(k.TYPE), false],
        [KASUS_EDM.reinsuranceSystem, sel(labelTypeCeding(k.TYPE_CEDING ?? '')), false],
        [KASUS_EDM.classOfBusiness, sel(k.BUSINESS_NAME), false],
        [KASUS_EDM.premiumMethod, labelMetodePremi(k.PRO_RATE_TYPE ?? ''), false],
      ],
    },
    {
      judul: KASUS_EDM.bagianPihak,
      medan: [
        [KASUS_EDM.ceding, sel(k.CEDING_CO_NAME), true],
        [KASUS_EDM.sob, sel(k.SOB_NAME), false],
        [KASUS_EDM.policyHolder, sel(k.POLICY_HOLDER_NAME), true],
        [KASUS_EDM.marketingOfficer, sel(k.MARKETING_NAME), true],
      ],
    },
    {
      judul: KASUS_EDM.bagianEndorsement,
      medan: [
        [KASUS_EDM.edmType, sel(labelEdmType(kasus.edmType)), false],
        ...subJenisKasus(kasus).map(([l, n]): MedanKepala => [l, n, false]),
        [BUAT_EDM.edmDate, selTanggal(kasus.edmDate), false],
        [KASUS_EDM.description, sel(kasus.description), true],
      ],
    },
  ]
}

/**
 * Sel Work Status kotak masuk: STATUSS kosong (kasus terbuka) → `InputEDMLife`; terisi → apa adanya
 * (keputusan work owner 08-10-2026, tampilan saja - database dan filter tetap memakai kosong = terbuka).
 */
/** Kasus kotak masuk dapat dibuka? Hanya yang terbuka (`STATUSS` kosong) - tertutup tetap terdaftar (09-10-2026). */
export function kasusEDMBisaDibuka(status: string | undefined): boolean {
  return (status ?? '').trim() === ''
}

export function selWorkStatus(v: string | undefined): string {
  return (v ?? '').trim() === '' ? INBOX_EDM.statusTerbuka : (v ?? '').trim()
}

/**
 * Bagian popup Details peserta (keputusan work owner 08-10-2026, tampilan saja): urutan bagian tetap untuk keempat
 * tipe; medan yang tidak dimiliki tipe itu tidak tampil. Varian tipe lain dimasukkan ke bagian paling sesuai:
 * `SHARE_RETRO` → Coverage; `*_REFUND`, `*_RETRO`, `*_REFUND_RETRO` premi/potongan dan `RI_ADMIN_FEE_*` → Premium.
 */
export const BAGIAN_RINCI: ReadonlyArray<{ judul: string; kolom: readonly string[] }> = [
  { judul: RINCI_EDM.bagianPeserta, kolom: ['CERTIFICATE_NO', 'POLICY_HOLDER', 'PLAN', 'SEX', 'ENTRY_AGE', 'CURRENT_AGE', 'CURRENCY'] },
  { judul: RINCI_EDM.bagianPeriode, kolom: ['BEGIN_DATE', 'EXPIRED_DATE', 'EFFECTIVE_DATE', 'STNC', 'WPC', 'PERIOD_MM'] },
  {
    judul: RINCI_EDM.bagianPertanggungan,
    kolom: ['SUM_INSURED', 'CEDING_RETENTION', 'SUM_REASURED', 'SHARE_NUSANTARA_RE_GROSS', 'SUM_AT_RISK_GROSS', 'RISK',
      'SUM_AT_RISK_RETRO', 'RETROCEDED_SHARE', 'SHARE_RETRO'],
  },
  {
    judul: RINCI_EDM.bagianPremi,
    kolom: ['RATE', 'EM_PERCENT', 'FACTOR', 'GROSS_PREMIUM', 'GROSS_PREMIUM_REFUND', 'GROSS_PREMIUM_RETRO',
      'GROSS_PREMIUM_REFUND_RETRO', 'DEDUCTION', 'DEDUCTION_REFUND', 'BROKERAGE_FEE', 'BROKERAGE_FEE_REFUND',
      'RI_ADMIN_FEE_RETRO', 'RI_ADMIN_FEE_REFUND_RETRO', 'NET_PREMIUM', 'NET_PREMIUM_REFUND', 'NET_PREMIUM_RETRO',
      'NET_PREMIUM_REFUND_RETRO'],
  },
]

/**
 * Medan rincian yang nilainya biasanya panjang (nomor sertifikat, nama pemegang polis, plan) - satu baris SELEBAR
 * bagiannya, bukan setengah kisi (perbaikan 08-10-2026: Certificate No panjang menimpa Policy Holder).
 */
export const MEDAN_RINCI_LEBAR: readonly string[] = ['CERTIFICATE_NO', 'POLICY_HOLDER', 'PLAN']

/** Medan rincian satu tipe dikelompokkan `BAGIAN_RINCI`; bagian tanpa medan tidak dikembalikan. */
export function kelompokRinci(kolom: readonly KolomKorpus[]): ReadonlyArray<{ judul: string; kolom: readonly KolomKorpus[] }> {
  return BAGIAN_RINCI.map((b) => ({
    judul: b.judul,
    kolom: b.kolom.flatMap((nama) => kolom.filter((k) => k.kolom === nama)),
  })).filter((b) => b.kolom.length > 0)
}

/** Label `Reinsurance System` dari kode `TYPE_CEDING`; kode tak dikenal apa adanya. */
export function labelTypeCeding(v: string): string {
  return OPSI_TYPE_CEDING[v] ?? sel(v)
}

/** Pilihan hapus yang dikirim bersama `Save` - `models.PilihanHapus`. */
export interface PilihanHapusEDM {
  pilih: readonly string[]
  semua: boolean
  kecuali: readonly string[]
}

export const PILIHAN_KOSONG: PilihanHapusEDM = { pilih: [], semua: false, kecuali: [] }

/**
 * `DELETE ALL` b13607 → `SelectAllEdmLife_act` b259: sakelar `CARI1` - seluruh baris ikut nilai
 * barunya, centang satu per satu sebelumnya dibuang.
 */
export function alihSemua(p: PilihanHapusEDM): PilihanHapusEDM {
  return p.semua ? PILIHAN_KOSONG : { pilih: [], semua: true, kecuali: [] }
}

function alih(daftar: readonly string[], id: string): string[] {
  return daftar.includes(id) ? daftar.filter((x) => x !== id) : [...daftar, id]
}

/** Kotak centang `.EdmBatal` b15753 satu baris; sesudah `DELETE ALL` ia mengisi pengecualian. */
export function alihBaris(p: PilihanHapusEDM, id: string): PilihanHapusEDM {
  return p.semua ? { ...p, kecuali: alih(p.kecuali, id) } : { ...p, pilih: alih(p.pilih, id) }
}

export function tercentang(p: PilihanHapusEDM, id: string): boolean {
  return p.semua ? !p.kecuali.includes(id) : p.pilih.includes(id)
}

/**
 * Kotak centang hidup: kasus `Perubahan Data` (grid b11899; grid Batal b17500 tanpa kotak) yang
 * terbuka dan belum disimpan (`Save` mati sesudahnya, b37200), atas peserta `Old` saja (R30).
 */
export function bolehCentang(
  p: { edmStatus: string },
  k: { status: string; sudahSimpan: boolean; edmType: string },
): boolean {
  return k.status === '' && !k.sudahSimpan && k.edmType === '1' && p.edmStatus === 'Old'
}

/** Jumlah peserta `Old` yang akan ditandai Delete: centang satu per satu, atau `DELETE ALL` dikurangi kecualiannya. */
export function jumlahDicentang(p: PilihanHapusEDM, cacah: Record<string, number>): number {
  return p.semua ? Math.max(0, (cacah.Old ?? 0) - p.kecuali.length) : p.pilih.length
}

/**
 * Proteksi tanpa perubahan (keputusan work owner 08-10-2026; server `models.AdaPerubahan` sumber kebenaran):
 * Perubahan Data boleh diproses bila ada centang (`jumlahDicentang` > 0) ATAU peserta `New` (`cacah.New`);
 * Batal selalu boleh.
 */
export function adaPerubahan(k: { edmType: string; cacah: Record<string, number> }, p: PilihanHapusEDM): boolean {
  if (k.edmType !== '1') return true
  return jumlahDicentang(p, k.cacah) > 0 || (k.cacah.New ?? 0) > 0
}

/**
 * Isi konfirmasi Process Endorsement (keputusan work owner 08-10-2026). Perubahan Data: jumlah centang (+ peserta
 * baru CSV bila ada); Batal: seluruh peserta `Old` (tanpa angka bila cacahnya tidak ada). Kalimat rekap selalu.
 * Angka dari `kasus.cacah` (jumlah peserta per `EDM_STATUS`) yang sudah dimuat - tidak ada hitungan baru.
 */
export function kalimatKonfirmasiProses(
  k: { edmType: string; cacah: Record<string, number> },
  p: PilihanHapusEDM,
): string[] {
  const isi = (t: string, n: number) => t.replace('{n}', String(n))
  const hasil: string[] = []
  if (k.edmType === '1') {
    const n = jumlahDicentang(p, k.cacah)
    hasil.push(n > 0 ? isi(PROSES_EDM.hapusN, n) : PROSES_EDM.hapusNol)
    const baru = k.cacah.New ?? 0
    if (baru > 0) hasil.push(isi(PROSES_EDM.baruN, baru))
  } else if (k.edmType === '3') {
    const n = k.cacah.Old ?? 0
    hasil.push(n > 0 ? isi(PROSES_EDM.batalN, n) : PROSES_EDM.batalSemua)
  }
  hasil.push(PROSES_EDM.rekap)
  return hasil
}

/** Urutan baris Summary Comparison per mata uang (08-10-2026). */
export const URUTAN_PERBANDINGAN: ReadonlyArray<{ kunci: 'before' | 'deleted' | 'added' | 'after' | 'difference'; label: string }> = [
  { kunci: 'before', label: PERBANDINGAN_EDM.before },
  { kunci: 'deleted', label: PERBANDINGAN_EDM.deleted },
  { kunci: 'added', label: PERBANDINGAN_EDM.added },
  { kunci: 'after', label: PERBANDINGAN_EDM.after },
  { kunci: 'difference', label: PERBANDINGAN_EDM.difference },
]

/** Angka bertanda untuk baris Difference: positif diberi '+', negatif tetap '-', nol tanpa tanda. */
export function selBertanda(v: string | undefined): string {
  const t = (v ?? '').trim()
  if (t === '') return selAngka(t)
  const nol = /^[-+]?0*(\.0*)?$/.test(t)
  if (nol || t.startsWith('-')) return selAngka(t)
  return '+' + selAngka(t.replace(/^\+/, ''))
}

/**
 * Teks kolom Participants Summary Comparison: JUMLAH orang (nilai mutlak) di semua baris, tanpa tanda +/− dan tanpa
 * warna (permintaan work owner 09-10-2026; sebelumnya Deleted negatif, Difference bertanda).
 */
export function teksPeserta(n: number): string {
  return String(Math.abs(n))
}

/**
 * Kelas merah untuk nilai NEGATIF (08-10-2026): teks desimal diputuskan dari tanda awalan "-" sesudah trim (tanpa
 * parseFloat); "-0"/nol bukan negatif. Angka (cacah) < 0 juga negatif. Selain itu `undefined`.
 */
export function kelasNegatif(v: string | number | undefined): string | undefined {
  if (typeof v === 'number') return v < 0 ? 'edm-negatif' : undefined
  const t = (v ?? '').trim()
  if (!t.startsWith('-') || /^-0*(\.0*)?$/.test(t)) return undefined
  return 'edm-negatif'
}

/** Kelas sel angka: `edm-angka` + `edm-negatif` bila negatif. */
export function kelasAngka(v: string | number | undefined): string {
  const n = kelasNegatif(v)
  return n === undefined ? 'edm-angka' : `edm-angka ${n}`
}

/** Arah nilai Difference untuk warnanya: naik (+), turun (−), atau nol. */
export function arahSelisih(v: string | undefined): 'naik' | 'turun' | 'nol' {
  const t = (v ?? '').trim()
  if (t === '' || /^[-+]?0*(\.0*)?$/.test(t)) return 'nol'
  return t.startsWith('-') ? 'turun' : 'naik'
}

/** Kolom grid peserta, urut korpus; `holder` false = grid Batal b17500 (tanpa `POLICY HOLDER`). */
/**
 * Kolom grid korpus yang DISEMBUNYIKAN dari tabel Participants karena sudah tampil di popup Details per peserta
 * (`RincianPeserta` → `RINCI_PESERTA`) - keputusan work owner 08-10-2026, tampilan saja. `SEX` ikut sesudah
 * ditambahkan ke popup Details (`RINCI_PESERTA` keempat tipe, sesudah Entry Age) - keputusan yang sama.
 */
export const KOLOM_GRID_DI_RINCIAN: readonly string[] = [
  'POLICY_HOLDER',
  'SEX',
  'ENTRY_AGE',
  'PLAN',
  'BEGIN_DATE',
  'EFFECTIVE_DATE',
  'EXPIRED_DATE',
]

/**
 * Kolom grid peserta yang tampil, urut korpus (b11899 / Batal b17500 tanpa `POLICY HOLDER`), dikurangi
 * `KOLOM_GRID_DI_RINCIAN` (08-10-2026). `POLICY HOLDER` kini tersembunyi untuk kedua grid.
 */
export function kolomGrid(): ReadonlyArray<readonly [label: string, kolom: string, jenis: 't' | 'd' | 'n']> {
  const semua: Array<readonly [string, string, 't' | 'd' | 'n']> = [
    [GRID_EDM.policyNo, 'POLICY_NO', 't'],
    [GRID_EDM.policyHolder, 'POLICY_HOLDER', 't'],
    [GRID_EDM.certificateNo, 'CERTIFICATE_NO', 't'],
    [GRID_EDM.nameOfInsured, 'NAME_OF_INSURED', 't'],
    [GRID_EDM.sex, 'SEX', 't'],
    [GRID_EDM.dateOfBirth, 'DOB', 'd'],
    [GRID_EDM.entryAge, 'ENTRY_AGE', 'n'],
    [GRID_EDM.plan, 'PLAN', 't'],
    [GRID_EDM.beginDate, 'BEGIN_DATE', 'd'],
    [GRID_EDM.effectiveDate, 'EFFECTIVE_DATE', 'd'],
    [GRID_EDM.expiredDate, 'EXPIRED_DATE', 'd'],
  ]
  return semua.filter(([, k]) => !KOLOM_GRID_DI_RINCIAN.includes(k))
}
