import { readdirSync, readFileSync, statSync } from 'node:fs'
import { dirname, join, relative, resolve, sep } from 'node:path'

import ts from 'typescript'
import { describe, expect, it } from 'vitest'

// Penjaga LAPISAN frontend - refactor bentuk B paket 8 (30-09-2026).
//
// Padanan `inti/penjaga/impor_lintas_modul_test.go` di backend. Aturannya,
// atas SETIAP berkas .ts/.tsx di bawah src/ (termasuk berkas uji):
//
//   1. `inti/**` hanya mengimpor `inti/**` - inti tidak mengenal modul.
//   2. `modul/X/**` hanya mengimpor `inti/**` dan `modul/X/**`. Yang mengenal
//      semua modul hanya lapisan aplikasi: `App.tsx`, `Beranda.tsx`,
//      `main.tsx`, dan daftar modul `modul/daftar.ts`.
//
// ⚠️ Impor dibaca pengurai TypeScript (`preProcessFile`), bukan pola teks:
// teks yang MENYEBUT sebuah impor di dalam string atau komentar tidak
// dihitung, dan impor dinamis `import('...')` ikut dihitung.

const SRC = join(__dirname, '..')

/** Lapisan sebuah berkas (jalur relatif src, bergaris-miring). */
function lapisan(rel: string): string {
  if (rel.startsWith('inti/')) return 'inti'
  const m = /^modul\/([^/]+)\//.exec(rel)
  if (m) return `modul:${m[1] ?? ''}`
  return 'aplikasi'
}

/** Alasan `dari` tidak boleh mengimpor `ke`; `null` = boleh. */
function pelanggaranLapisan(dari: string, ke: string): string | null {
  const a = lapisan(dari)
  const b = lapisan(ke)
  if (a === 'aplikasi') return null
  if (a === 'inti') return b === 'inti' ? null : 'inti tidak mengenal modul maupun lapisan aplikasi'
  if (b === 'inti' || b === a) return null
  return 'modul hanya mengimpor inti/ dan dirinya sendiri; yang merakit modul hanya modul/daftar.ts dan App.tsx'
}

function semuaBerkas(dir: string, out: string[] = []): string[] {
  for (const nama of readdirSync(dir)) {
    const jalur = join(dir, nama)
    if (statSync(jalur).isDirectory()) semuaBerkas(jalur, out)
    else if (/\.tsx?$/.test(nama)) out.push(jalur)
  }
  return out
}

/** Setiap impor RELATIF tiap berkas, sebagai jalur relatif src. */
function imporRelatif(): { dari: string; ke: string; spec: string }[] {
  const hasil: { dari: string; ke: string; spec: string }[] = []
  for (const f of semuaBerkas(SRC)) {
    const dari = relative(SRC, f).split(sep).join('/')
    for (const imp of ts.preProcessFile(readFileSync(f, 'utf8'), true, true).importedFiles) {
      if (!imp.fileName.startsWith('.')) continue // paket npm
      const ke = relative(SRC, resolve(dirname(f), imp.fileName)).split(sep).join('/')
      hasil.push({ dari, ke, spec: imp.fileName })
    }
  }
  return hasil
}

describe('lapisan frontend: inti <- modul <- aplikasi', () => {
  it('nol impor lintas lapisan di seluruh src/', () => {
    const impor = imporRelatif()
    const perLapis = new Map<string, number>()
    for (const { dari } of impor) perLapis.set(lapisan(dari), (perLapis.get(lapisan(dari)) ?? 0) + 1)
    // ⛔ Penjaga yang membaca nol impor di satu lapis lulus atas apa pun.
    expect(perLapis.get('inti') ?? 0).toBeGreaterThan(30)
    for (const m of ['claimlife', 'premiumlist', 'komite', 'treaty']) {
      expect(perLapis.get(`modul:${m}`) ?? 0, `impor modul ${m}`).toBeGreaterThan(10)
    }
    const langgar = impor
      .map((i) => ({ ...i, alasan: pelanggaranLapisan(i.dari, i.ke) }))
      .filter((i) => i.alasan !== null)
      .map((i) => `${i.dari} -> ${i.spec}: ${i.alasan}`)
    expect(langgar).toEqual([])
  })

  it('aturannya menggigit dua arah', () => {
    const kasus: [string, string, boolean][] = [
      ['modul/claimlife/pages/X.tsx', 'modul/premiumlist/api.ts', false],
      ['modul/komite/pages/X.test.ts', 'modul/claimlife/labels.ts', false],
      ['modul/komite/rute.tsx', 'modul/daftar.ts', false],
      ['modul/komite/rute.tsx', 'App.tsx', false],
      ['modul/komite/rute.tsx', 'modul/komite/pages/InboxKomite.tsx', true],
      ['modul/komite/api.ts', 'inti/klien.ts', true],
      ['inti/components/Shell.tsx', 'modul/daftar.ts', false],
      ['inti/lib/daftarMenu.ts', 'modul/treaty/labels.ts', false],
      ['inti/components/Shell.tsx', 'inti/lib/daftarMenu.ts', true],
      ['modul/daftar.ts', 'modul/treaty/rute.tsx', true],
      ['App.tsx', 'modul/claimlife/api.ts', true],
      ['modul/komitex/api.ts', 'modul/komite/api.ts', false],
    ]
    for (const [dari, ke, boleh] of kasus) {
      expect(pelanggaranLapisan(dari, ke) === null, `${dari} -> ${ke}`).toBe(boleh)
    }
  })
})
