// Tombol `Generate Excel` sub-tab Achievement rincian Limits Prop —
// `Activity/GenerateCSVTreaty.xml` korpus Adjustment:
//
//   4    FOR_EACH `.AchievementLists`: TempData.pxResults(idx).CARI1..13 =
//        Quarter, QUARTERYEAR, Currency, PREMIUM, RICOMM, BROKERAGE,
//        NETPREMIUM, PaidClaim, CASHCALL, OutstandingClaim, IncuredClaim,
//        Total, LossRatio
//   5    MSOGenerateExcelFile (`FSFileName` = CSVAchievementTreatyIn.xlsx)
//
// ⚠️ Korpus ini mengekspor SEMUA baris — tanpa penyaring baris ber-Quarter
// kosong yang dipakai layar Treaty In. Judul kolom = judul kolom grid
// `AchievementLists` di ekspor (templat Excel Pega tidak diekspor).

import { unduhXlsx } from '../../../../inti/frontend/lib/exportXlsx'
import type { SisiPenyesuaian } from '../api'
import type { ButirKerangka } from '../ekspor/jenis'
import { KERANGKA_RINCIAN } from '../ekspor/kerangka.gen'
import { PENYESUAIAN } from '../labelsPenyesuaian'
import { teks } from './baris'

/** Urutan kolom `CARI1`..`CARI13` langkah 4. */
export const KOLOM_CSV_ACHIEVEMENT = [
  'Quarter', 'QUARTERYEAR', 'Currency', 'PREMIUM', 'RICOMM', 'BROKERAGE', 'NETPREMIUM',
  'PaidClaim', 'CASHCALL', 'OutstandingClaim', 'IncuredClaim', 'Total', 'LossRatio',
] as const

/** Judul kolom grid `AchievementLists` Section `DetailLimits`, per kunci. */
function judulKolom(): Record<string, string> {
  const out: Record<string, string> = {}
  const jalan = (bs: readonly ButirKerangka[]) => {
    for (const b of bs) {
      if (b.t === 'blok') jalan(b.anak)
      if (b.t === 'grid' && b.larik === 'AchievementLists') b.kunci.forEach((k, i) => (out[k] = b.kolom[i] ?? k))
    }
  }
  jalan(KERANGKA_RINCIAN.DetailLimits ?? [])
  return out
}

/** Baris Excel dari halaman Detail (Treaty Group) tombol itu. */
export function barisAchievement(halaman: SisiPenyesuaian): Record<string, string>[] {
  return (halaman.larik.AchievementLists ?? []).map((b) => Object.fromEntries(KOLOM_CSV_ACHIEVEMENT.map((k) => [k, teks(b[k])])))
}

export function unduhAchievement(halaman: SisiPenyesuaian): void {
  const judul = judulKolom()
  unduhXlsx(
    PENYESUAIAN.berkasAchievement,
    KOLOM_CSV_ACHIEVEMENT.map((k) => ({ kunci: k, label: judul[k] ?? k })),
    barisAchievement(halaman),
  )
}
