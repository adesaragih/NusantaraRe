// Copy Old (perintah work owner 07-10-2026 "SAMA SEPERTI MASTER PRODUCTNAME LIFE, KHUSUS BUAT SUPERUSER") - fungsi murni
// popup `components/DialogCopyOld.tsx`. Pola: `modul/masterproductnamelife/frontend/bentuk.ts` (03-10-2026).

import type { DokumenLama, HakPortal, HasilSalinLama, StatusSalinLama } from './api'

/** Tombol Copy Old tampil hanya bila server menyatakan akun superadmin ber-hak penuh. */
export const tampilCopyOld = (h: HakPortal | null | undefined): boolean => h?.copyOld === true

/** `Search` menyaring EDM Number, Policy Number, EDM No, SOB, Ceding (tanpa membedakan huruf). */
export function saringLama(daftar: readonly DokumenLama[], kata: string): DokumenLama[] {
  const k = kata.trim().toUpperCase()
  if (k === '') return [...daftar]
  return daftar.filter((d) =>
    [d.id, d.noPolis, d.edmNo, d.sobName, d.cedingCoName].some((v) => v.toUpperCase().includes(k)),
  )
}

/** ID yang dikirim `Process Copy`: terpilih DAN boleh disalin, urutan daftar (urutan generasi). */
export function idBolehDisalin(daftar: readonly DokumenLama[], terpilih: ReadonlySet<string>): string[] {
  return daftar.filter((d) => d.bolehDisalin && terpilih.has(d.id)).map((d) => d.id)
}

/** Cacah hasil `Process Copy` per status. */
export function ringkasSalin(hasil: readonly HasilSalinLama[]): Record<StatusSalinLama, number> {
  const r: Record<StatusSalinLama, number> = { disalin: 0, sudahAda: 0, ditolak: 0, gagal: 0 }
  for (const h of hasil) r[h.status]++
  return r
}
