// Letak kode frontend, untuk uji yang MEMBACA sumber - struktur tim satu folder
// per modul (30-09-2026).
//
// Untuk apa berkas ini: sampai 30-09-2026 seluruh kode frontend tinggal di satu
// pohon (`frontend/src`), dan uji yang menelusuri sumber cukup naik ke sana.
// Kini kode tersebar di tiga tempat - perakit `frontend/`, kerangka bersama
// `inti/frontend/`, dan setiap `modul/<nama>/frontend/`. Penelusur yang masih
// naik ke satu akar diam-diam MENYEMPIT: ia tetap hijau, hanya membaca lebih
// sedikit. Satu daftar di sini dipakai setiap penelusur.
//
// ⛔ Hanya untuk berkas uji. Kode aplikasi tidak mengimpornya (ia memakai
// `node:fs`), jadi berkas ini tidak pernah masuk bundel.

import { existsSync, readdirSync, statSync } from 'node:fs'
import { join, relative, sep } from 'node:path'

/** Folder APP_RNM - dari letak berkas ini (`inti/frontend/uji`). */
export const AKAR_APLIKASI = join(__dirname, '..', '..', '..')

/** Jalur `j` relatif APP_RNM, bergaris-miring. */
export function relatifAplikasi(j: string): string {
  return relative(AKAR_APLIKASI, j).split(sep).join('/')
}

/**
 * Setiap akar kode frontend: perakit, kerangka bersama, lalu setiap folder
 * `modul/<nama>/frontend` yang ADA - modul baru ikut tanpa menyunting berkas ini.
 */
export function akarSumberFrontend(): string[] {
  const modul = join(AKAR_APLIKASI, 'modul')
  const perModul = existsSync(modul)
    ? readdirSync(modul, { withFileTypes: true })
        .filter((d) => d.isDirectory())
        .map((d) => join(modul, d.name, 'frontend'))
        .filter((d) => existsSync(d))
    : []
  return [join(AKAR_APLIKASI, 'frontend', 'src'), join(AKAR_APLIKASI, 'inti', 'frontend'), ...perModul]
}

/** Setiap berkas .ts/.tsx di bawah akar-akar itu, termasuk berkas uji. */
export function berkasTS(akar: readonly string[] = akarSumberFrontend()): string[] {
  const hasil: string[] = []
  const telusur = (dir: string): void => {
    for (const nama of readdirSync(dir)) {
      const j = join(dir, nama)
      if (statSync(j).isDirectory()) {
        if (nama !== 'node_modules' && nama !== 'dist') telusur(j)
      } else if (/\.tsx?$/.test(nama)) {
        hasil.push(j)
      }
    }
  }
  for (const a of akar) telusur(a)
  return hasil.sort()
}

/** Seperti `berkasTS`, tanpa berkas uji. */
export function berkasSumberTS(akar?: readonly string[]): string[] {
  return berkasTS(akar).filter((j) => !/\.test\.tsx?$/.test(j))
}
