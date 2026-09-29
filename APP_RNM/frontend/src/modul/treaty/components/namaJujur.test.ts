// Penjaga NAMA JUJUR sisi React — Treaty Contract Out (penyimpangan sadar 8;
// AC 13, AC 50 spec).
//
// Nama korpus yang berbohong (`_Old`, `testingKurs`) tidak boleh merambat
// ke pengenal sistem baru. Yang dipindai: seluruh berkas sumber modul ini
// (folder `modul/treaty/components`, `modul/treaty/pages`, dan
// `modul/treaty/labels.ts`), sesudah komentar dan string
// literal dibuang — nama rule sumber boleh disebut sebagai bukti.

import { existsSync, readdirSync, readFileSync, statSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

// Refactor bentuk B (30-09-2026): folder modul kini `modul/treaty/`.
// Paket 8: akar modulnya IKUT dipindai - `api.ts`, `menu.ts`, `rute.tsx`,
// `labels.ts`. Dulu klien backend modul ini hidup di `services/api.ts`
// bersama semua modul, dan karena itu tidak pernah terbaca penjaga ini.
const AKAR = join(__dirname, '..')
const FOLDER = [join(AKAR, 'components'), join(AKAR, 'pages'), AKAR]
const BERKAS_TUNGGAL: string[] = []

function sumberModul(): string[] {
  const hasil: string[] = []
  for (const f of FOLDER) {
    if (!existsSync(f)) continue
    for (const nama of readdirSync(f)) {
      const jalur = join(f, nama)
      if (statSync(jalur).isFile() && /\.tsx?$/.test(nama) && !/\.test\.tsx?$/.test(nama)) hasil.push(jalur)
    }
  }
  for (const b of BERKAS_TUNGGAL) if (existsSync(b)) hasil.push(b)
  return hasil
}

/** Buang komentar dan string literal, sisakan pengenal. */
export function hanyaPengenal(teks: string): string {
  return teks
    .replace(/\/\*[\s\S]*?\*\//g, '')
    .split('\n')
    .filter((b) => !b.trimStart().startsWith('//') && !b.trimStart().startsWith('*'))
    .join('\n')
    .replace(/`[^`]*`|'(?:[^'\\]|\\.)*'|"(?:[^"\\]|\\.)*"/g, '""')
}

const POLA = /(?:\b|[a-z_])Old(?:[^a-z]|$)|_OLD_|\bOLD\b|[Tt]esting|JSON_KLAIM/

describe('nol nama bohong di pengenal modul', () => {
  it('ada berkas sumber modul yang terbaca', () => {
    expect(sumberModul().length).toBeGreaterThan(0)
  })
  it.each(sumberModul())('%s', (jalur) => {
    const kode = hanyaPengenal(readFileSync(jalur, 'utf8'))
    const kena = kode.match(POLA)
    expect(kena, kena ? `pengenal ${kena[0]}` : '').toBeNull()
  })
  it('polanya menggigit bentuk yang dilarang, tidak menuduh kata lain', () => {
    for (const buruk of ['browseOld()', 'kursTesting', 'const Old = 1', 'x_OLD_y']) expect(POLA.test(buruk)).toBe(true)
    for (const aman of ['Golden', 'hold', 'Threshold', 'holder']) expect(POLA.test(aman)).toBe(false)
    expect(POLA.test(hanyaPengenal("const x = 'BrowseReinsuranceType_RD_Old'"))).toBe(false)
  })
})
