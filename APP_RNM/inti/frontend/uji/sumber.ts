// Letak kode frontend, untuk uji yang MEMBACA sumber - struktur tim satu folder
// per modul (30-09-2026).
//
// Untuk apa berkas ini: sampai 30-09-2026 seluruh kode frontend tinggal di satu
// pohon (dulu `frontend/src`), dan uji yang menelusuri sumber cukup naik ke sana.
// Kini kode tersebar di tiga tempat - perakit `frontend/`, kerangka bersama
// `inti/frontend/`, dan setiap `modul/<nama>/frontend/`. Penelusur yang masih
// naik ke satu akar diam-diam MENYEMPIT: ia tetap hijau, hanya membaca lebih
// sedikit. Satu daftar di sini dipakai setiap penelusur.
//
// ⛔ Hanya untuk berkas uji. Kode aplikasi tidak mengimpornya (ia memakai
// `node:fs`), jadi berkas ini tidak pernah masuk bundel.

import { existsSync, readdirSync, readFileSync, statSync } from 'node:fs'
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
  return [join(AKAR_APLIKASI, 'frontend'), join(AKAR_APLIKASI, 'inti', 'frontend'), ...perModul]
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

/**
 * Folder korpus setiap modul yang `MODUL.md`-nya berstatus `belum dimigrasi`,
 * urut abjad - dibaca dari baris `Status` dan `Folder korpus` tabel kuncinya.
 *
 * Untuk apa: uji "kelompok yang belum dimigrasi tetap BERDIRI" dulu mengunci
 * angkanya (`toHaveLength(16)`) di berkas bersama, sehingga modul yang
 * mendapat butir menu pertamanya harus menyunting uji milik tim inti - dan dua
 * modul yang memulai bersamaan berkonflik di baris yang sama. Kini angkanya
 * PERNYATAAN pemilik modul di foldernya sendiri: modul yang mendapat butir
 * menu pertamanya mengubah `Status`-nya menjadi `dimigrasi`, dan uji yang
 * memakai daftar ini merah bila pernyataan itu tidak sesuai dengan menu.
 *
 * ⛔ Folder berawalan `_` (`_templat`) bukan modul. MODUL.md tanpa kedua baris
 * itu, status di luar dua nilai, atau nol MODUL.md terbaca = galat, bukan
 * daftar kosong yang meluluskan segalanya.
 */
export function folderKorpusBelumDimigrasi(): string[] {
  const modul = join(AKAR_APLIKASI, 'modul')
  const hasil: string[] = []
  let dibaca = 0
  for (const d of readdirSync(modul, { withFileTypes: true })) {
    if (!d.isDirectory() || d.name.startsWith('_')) continue
    const isi = readFileSync(join(modul, d.name, 'MODUL.md'), 'utf8')
    const status = /^\| Status \| (.+?) \|\s*$/m.exec(isi)?.[1]
    const korpus = /^\| Folder korpus \| `(.+?)` \|\s*$/m.exec(isi)?.[1]
    if (status === undefined || korpus === undefined) {
      throw new Error(`modul/${d.name}/MODUL.md tanpa baris Status atau Folder korpus`)
    }
    if (status !== 'dimigrasi' && status !== 'belum dimigrasi') {
      throw new Error(`modul/${d.name}/MODUL.md: Status ${JSON.stringify(status)}, mau dimigrasi atau belum dimigrasi`)
    }
    dibaca++
    if (status === 'belum dimigrasi') hasil.push(korpus)
  }
  if (dibaca === 0) throw new Error(`nol MODUL.md terbaca di ${modul}; pembacanya yang rusak`)
  return hasil.sort()
}
