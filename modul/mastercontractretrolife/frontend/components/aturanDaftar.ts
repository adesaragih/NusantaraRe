// Aturan bersama Reinsurer List dan Security Reinsurer - dipisah supaya kedua
// panel memakainya tanpa saling mengimpor (PanelReinsurer memuat PanelSecurity).

import type { OpsiSaring } from '../../../../inti/frontend/components/ui/pilihSaring'

/**
 * Pilihan nama reinsurer yang BELUM dipakai baris lain di daftar yang sama -
 * satu nama sekali per kontrak (Reinsurer List) dan per reinsurer induk
 * (Security Reinsurer), keputusan work owner 04-10-2026. Milik baris yang
 * sedang diubah (`idSendiri`) tetap tersedia. Server menolak hal yang sama
 * (`PesanReinsurerGanda`, `PesanSecurityGanda`).
 */
export function reinsurerTersedia(
  pilihan: readonly OpsiSaring[],
  daftar: readonly { id: string; reinsurerId: string }[],
  idSendiri: string,
): OpsiSaring[] {
  const terpakai = new Set(daftar.filter((r) => r.id !== idSendiri).map((r) => r.reinsurerId.trim()))
  return pilihan.filter((o) => !terpakai.has(o.value.trim()))
}

/** Kelas baris total - mencolok bila total ≠ 100 (tiket 05 AC 19). */
export function kelasTotal(totalBukan100: boolean): string {
  return totalBukan100 ? 'mcrl-total mcrl-total--bukan100' : 'mcrl-total'
}
