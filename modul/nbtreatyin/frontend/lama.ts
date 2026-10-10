// Copy Old (perintah work owner 07-10-2026) - fungsi murni popup `components/DialogCopyOld.tsx`. Pola sama dengan Copy
// Old EDM Treaty In (`modul/edmtreatyin/frontend/lama.ts`) dan Product Name Life.

import type { DokumenLama, HakPortal, HasilSalinLama, StatusSalinLama } from './api'

/** Tombol Copy Old tampil hanya bila server menyatakan akun superadmin ber-hak penuh. */
export const tampilCopyOld = (h: HakPortal | null | undefined): boolean => h?.copyOld === true

/** `Search` menyaring NB Number, Master ID, Policy Number, Insured, Group Business, SOB, Ceding (tanpa beda huruf). */
export function saringLama(daftar: readonly DokumenLama[], kata: string): DokumenLama[] {
  const k = kata.trim().toUpperCase()
  if (k === '') return [...daftar]
  return daftar.filter((d) =>
    [d.id, d.noOffer, d.noPolis, d.insuredName, d.businessName, d.sobName, d.cedingCoName].some((v) =>
      v.toUpperCase().includes(k),
    ),
  )
}

/** ID yang dikirim `Process Copy`: terpilih DAN boleh disalin, urutan daftar. */
export function idBolehDisalin(daftar: readonly DokumenLama[], terpilih: ReadonlySet<string>): string[] {
  return daftar.filter((d) => d.bolehDisalin && terpilih.has(d.id)).map((d) => d.id)
}

/** Cacah hasil `Process Copy` per status. */
export function ringkasSalin(hasil: readonly HasilSalinLama[]): Record<StatusSalinLama, number> {
  const r: Record<StatusSalinLama, number> = { disalin: 0, sudahAda: 0, ditolak: 0, gagal: 0 }
  for (const h of hasil) r[h.status]++
  return r
}
