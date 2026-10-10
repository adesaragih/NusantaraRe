// Medan ber-sumber dengan `tampilan` (Consultant / Adjuster, work owner 08-10-2026: "yang dropdown hanya dari namanya
// aja; untuk ID dihapus dari tampilan, tapi tetap simpan ID"): yang dipilih dan tampil adalah nama, yang disimpan tetap
// nilai medan (ID). Medan tanpa `tampilan` tidak berubah: nilai sebagai label, nama sebagai keterangan.

import type { Pilihan, Tata } from '../api'

/** Butir dropdown PilihSaring. */
export function butirSaring(
  t: Pick<Tata, 'tampilan'>,
  opsi: readonly Pilihan[],
): { value: string; label: string; keterangan?: string }[] {
  if (t.tampilan) return opsi.map((o) => ({ value: o.nilai, label: o.label }))
  return opsi.map((o) => ({ value: o.nilai, label: o.nilai, keterangan: o.label !== o.nilai ? o.label : undefined }))
}

/** Teks nilai terpilih: jalur tampilan, lalu label daftar pilihan; tanpa `tampilan` = nilainya sendiri. */
export function teksTerpilih(
  t: Pick<Tata, 'tampilan'>,
  nilai: string,
  ambil: (jalur: string) => string,
  opsi: readonly Pilihan[],
): string {
  if (!t.tampilan) return nilai
  return ambil(t.tampilan) || opsi.find((o) => o.nilai === nilai)?.label || ''
}
