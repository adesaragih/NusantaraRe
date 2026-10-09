// Label tab **Achievement In IDR** — teks APA ADANYA dari ekspor.
//
//   Section/TreatyInTabsProportional.xml   nama tab, @1581141
//   Section/TreatyInTabsAchievement.xml    cangkang tab
//   Section/AchievementCombine.xml         grid + kaki
//   Activity/GetAchievement.xml            angkanya
//
// ⛔ `Incured` dieja begitu di ekspor (bukan `Incurred`). Dipertahankan —
// label adalah yang pemakai lama baca.

/**
 * Keenam kolom grid, urut ekspor, beserta kunci dan desimalnya.
 *
 * ⭐ `Net Loss Ratio` PERSEN, lima lainnya UANG. Dibedakan supaya rasio
 * tidak ikut diformat sebagai rupiah.
 */
export const KOLOM_ACHIEVEMENT = [
  { label: 'Treaty Group', kunci: 'TreatyGroup', jenis: 'teks' },
  { label: 'Reins Type', kunci: 'TreatyType', jenis: 'teks' },
  { label: 'Gross Premium Before Claim in IDR', kunci: 'TotalAchPremium', jenis: 'uang' },
  { label: 'Net Premium Before Claim in IDR', kunci: 'TotalAchNetPremium', jenis: 'uang' },
  { label: 'Incured Claim in IDR', kunci: 'TotalAchIncured', jenis: 'uang' },
  { label: 'Net Loss Ratio', kunci: 'LossRatio', jenis: 'persen' },
] as const

/**
 * ⛔ KAKI HANYA PUNYA TIGA SEL, dan letaknya BUKAN di bawah tiap kolom.
 *
 * Ekspor menaruh `Spacers` ber-`pyCondition 1=2` di kolom 1-3 dan 5, lalu
 * sel nilai di bawah `Net Premium`, `Incured Claim`, dan `Net Loss Ratio`.
 * Jadi `Gross Premium` NOL punya total di kaki — itu bunyi ekspornya, bukan
 * sel yang lupa dipasang.
 */
export const KAKI_ACHIEVEMENT: readonly (string | null)[] = [
  null,
  null,
  null,
  'SumTotalAchievNetPremium',
  'SumTotalAchievIncured',
  'SumLossRatio',
]

export const ACHIEVEMENT = {
  judul: 'Achievement In IDR',
  /**
   * ⭐ Label kaki. `pyValue = Total in IDR` ber-`pyCondition 1=2` TETAPI
   * `pyVisible = ALWAYS` — dan `pyCondition` hanya berlaku ketika
   * `pyVisible = OTHER`. Jadi labelnya TAMPIL. Spacers di sekitarnya
   * ber-`pyVisible = OTHER` dan karena itu memang tersembunyi.
   */
  totalIDR: 'Total in IDR',
  /** Baris terakhir grid — `.Detail(<LAST>).TreatyType` di ekspor. */
  barisTotal: ' Total In IDR',
  tanpaBaris: 'No items',
  petunjukKosong: 'No Achievement rows for this contract yet.',
  gagalMuat: 'Failed to load Achievement',
} as const

/** Desimal — `uang` dua, `persen` dua. */
export const DESIMAL_ACHIEVEMENT = { uang: 2, persen: 2 } as const
