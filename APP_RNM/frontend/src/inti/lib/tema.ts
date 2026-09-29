/**
 * Tema tampilan — terang atau gelap (29-09-2026).
 *
 * Tiga keputusan, ketiganya teruji di tema.test.ts:
 *
 *   1. Pilihan pemakai MENANG. Tanpa pilihan tersimpan, tema mengikuti
 *      pengaturan sistem operasi (`prefers-color-scheme`).
 *   2. Tema dipasang sebagai `data-theme` pada <html> — SATU sakelar yang
 *      dibaca styles.css. Skrip sebaris di index.html memasangnya SEBELUM
 *      halaman tergambar; tanpa itu mode gelap berkedip putih setiap muat.
 *   3. Penyimpanan tidak dapat dipercaya (pola lipatMenu.ts): localStorage
 *      diblokir, penuh, atau berisi sampah → jatuh ke tema sistem tanpa
 *      melempar.
 *
 * ⚠️ KUNCI dan kedua nilainya DIULANG di skrip sebaris index.html: skrip itu
 * berjalan sebelum modul mana pun dimuat, jadi tidak dapat mengimpor berkas
 * ini. Mengubah salah satunya saja membuat halaman dimuat dengan tema yang
 * salah lalu melompat — tema.test.ts menjaga keduanya tetap sama.
 */
import type { GudangMini } from './lipatMenu'

export type Tema = 'terang' | 'gelap'

export const KUNCI_TEMA = 'tema-tampilan'

/** Tema tersimpan, atau null bila belum pernah dipilih / tak terbaca. */
export function bacaTema(gudang: GudangMini | null): Tema | null {
  if (!gudang) return null
  try {
    const v = gudang.getItem(KUNCI_TEMA)
    return v === 'terang' || v === 'gelap' ? v : null
  } catch {
    return null
  }
}

export function simpanTema(gudang: GudangMini | null, tema: Tema): void {
  if (!gudang) return
  try {
    gudang.setItem(KUNCI_TEMA, tema)
  } catch {
    // Kuota penuh / mode privat: tema tidak diingat, dan itu bukan galat.
  }
}

/** Tema yang berlaku: pilihan tersimpan, atau tema sistem. */
export function temaBerlaku(tersimpan: Tema | null, sistemGelap: boolean): Tema {
  return tersimpan ?? (sistemGelap ? 'gelap' : 'terang')
}

/** Nilai `data-theme` pada <html> — nama yang dibaca styles.css. */
export function atributTema(tema: Tema): 'light' | 'dark' {
  return tema === 'gelap' ? 'dark' : 'light'
}
