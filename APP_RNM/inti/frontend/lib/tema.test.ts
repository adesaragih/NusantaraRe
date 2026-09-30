import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import type { GudangMini } from './lipatMenu'
import { atributTema, bacaTema, KUNCI_TEMA, simpanTema, temaBerlaku } from './tema'

function gudang(isi: Record<string, string> = {}): GudangMini {
  return {
    getItem: (k) => isi[k] ?? null,
    setItem: (k, v) => {
      isi[k] = v
    },
  }
}

const rusak: GudangMini = {
  getItem: () => {
    throw new Error('diblokir')
  },
  setItem: () => {
    throw new Error('kuota penuh')
  },
}

describe('tema tampilan', () => {
  it('pilihan tersimpan MENANG atas tema sistem', () => {
    expect(temaBerlaku('terang', true)).toBe('terang')
    expect(temaBerlaku('gelap', false)).toBe('gelap')
  })

  it('tanpa pilihan, tema mengikuti sistem operasi', () => {
    expect(temaBerlaku(null, true)).toBe('gelap')
    expect(temaBerlaku(null, false)).toBe('terang')
  })

  it('disimpan lalu dibaca kembali', () => {
    const g = gudang()
    simpanTema(g, 'gelap')
    expect(bacaTema(g)).toBe('gelap')
  })

  it('isi sampah dibaca sebagai "belum memilih", bukan tema', () => {
    // ⛔ Nilai asing yang diterima apa adanya menjadi `data-theme` yang tak
    // dikenali styles.css — layar terang dengan tombol yang mengaku gelap.
    expect(bacaTema(gudang({ [KUNCI_TEMA]: 'dark' }))).toBeNull()
    expect(bacaTema(gudang({ [KUNCI_TEMA]: '' }))).toBeNull()
  })

  it('penyimpanan yang gagal tidak melempar', () => {
    expect(bacaTema(rusak)).toBeNull()
    expect(bacaTema(null)).toBeNull()
    expect(() => {
      simpanTema(rusak, 'gelap')
    }).not.toThrow()
  })

  it('nilai data-theme sama dengan yang dibaca styles.css', () => {
    expect(atributTema('gelap')).toBe('dark')
    expect(atributTema('terang')).toBe('light')
    const css = readFileSync(join(__dirname, '..', 'styles.css'), 'utf8')
    expect(css).toContain(':root[data-theme="dark"]')
  })

  it('skrip sebaris index.html memakai kunci dan nilai yang SAMA', () => {
    // Skrip itu tidak dapat mengimpor tema.ts (ia jalan sebelum modul
    // dimuat), jadi kesamaannya hanya dapat dijaga di sini.
    const html = readFileSync(join(__dirname, '..', '..', '..', 'frontend', 'index.html'), 'utf8')
    expect(html).toContain(`'${KUNCI_TEMA}'`)
    expect(html).toContain("'gelap'")
    expect(html).toContain("'terang'")
    expect(html).toContain('(prefers-color-scheme: dark)')
    expect(html).toContain('dataset.theme')
  })
})
