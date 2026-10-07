// Pembantu tampilan MURNI modul Endorsement Life.
//
// ⛔ Angka tidak pernah menjadi `Number`: sel angka lewat `formatNumber` (digit, bukan float).

import { DESIMAL_TAK_DIBATASI, formatNumber } from '../../../inti/frontend/lib/format'
import type { KasusEDM } from './api'
import { BUAT_EDM, KASUS_EDM, OPSI_EDM_TYPE, OPSI_METODE_PREMI, OPSI_TYPE_CEDING } from './labels'

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
