// Aturan layar Reinsurance Type - fungsi murni, diuji tanpa DOM. Aturan sebenarnya dijaga backend (`backend/services`):
// huruf besar, nama tidak kembar, nilai warisan dibiarkan selama tidak diubah, hak View only.

import type { Isian, Jenis } from './api'
import { RT } from './labels'

/** Pilihan untuk nilai baru (backend `models.Type`, `models.Flag`, `models.GroupType`). */
export const TYPE = ['1', '2', '3', '4'] as const
export const FLAG = ['active', 'inactive'] as const
export const GROUP_TYPE = ['OR', 'QS', 'RI', 'SPL'] as const

/** Isian form dari satu baris (Edit) atau kosong (Add: Flag active, Code 00). */
export function isianDari(j: Jenis | null): Isian {
  if (j === null) return { name: '', type: '', soaName: '', code: '00', flag: 'active', noUrut: '', groupType: '' }
  return {
    name: j.name,
    type: j.type,
    soaName: j.soaName,
    code: j.code,
    flag: j.flag,
    noUrut: j.noUrut,
    groupType: j.groupType,
  }
}

/** Isian yang dikirim: spasi tepi dibuang (huruf besar dan Code kosong = 00 dikerjakan backend). */
export function rapikanIsian(i: Isian): Isian {
  return {
    name: i.name.trim(),
    type: i.type.trim(),
    soaName: i.soaName.trim(),
    code: i.code.trim(),
    flag: i.flag.trim(),
    noUrut: i.noUrut.trim(),
    groupType: i.groupType.trim(),
  }
}

/**
 * Pemeriksaan awal sebelum dikirim; `null` = boleh dikirim. Type dan Flag wajib untuk baris baru; baris lama boleh
 * membawa nilai warisannya (kosong) selama tidak diubah. Backend tetap memeriksa seluruhnya.
 */
export function periksaIsian(i: Isian, lama: Jenis | null): string | null {
  if (i.name.trim() === '') return RT.galatNama
  if (i.type.trim() === '' && (lama === null || lama.type !== '')) return RT.galatType
  if (i.flag.trim() === '' && (lama === null || lama.flag !== '')) return RT.galatFlag
  if (i.code.trim() !== '' && !/^[0-9]+$/.test(i.code.trim())) return RT.galatCode
  if (i.noUrut.trim() !== '' && !/^[0-9]+$/.test(i.noUrut.trim())) return RT.galatNoUrut
  return null
}

/**
 * Opsi dropdown: pilihan sah + nilai warisan terpasang yang di luar pilihan (Flag `1`, Type kosong) - ditandai "(old
 * value)" supaya baris lama dapat disimpan tanpa mengubahnya.
 */
export function opsiDengan(
  pilihan: readonly string[],
  terpasang: string,
  label: (v: string) => string = (v) => v,
): { value: string; label: string }[] {
  const opsi = pilihan.map((v) => ({ value: v, label: label(v) }))
  if (terpasang !== '' && !pilihan.includes(terpasang)) opsi.push({ value: terpasang, label: RT.nilaiLama(terpasang) })
  return opsi
}

/** Kelas lencana Flag. */
export function kelasFlag(flag: string): string {
  if (flag === 'active') return 'reinsurancetype__lencana reinsurancetype__lencana--aktif'
  if (flag === 'inactive') return 'reinsurancetype__lencana reinsurancetype__lencana--nonaktif'
  return 'reinsurancetype__lencana'
}
