import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { layarSesudahMenu, menuEDMDipilihUlang, type Layar } from './rute'

// Klik menu Endorsement Life (ulang) selalu kembali ke kotak masuk - keputusan work owner 07-10-2026.

const RUTE = readFileSync(join(__dirname, 'rute.tsx'), 'utf8')

describe('menuEDMDipilihUlang', () => {
  it('ketuk tidak berubah → false', () => {
    expect(menuEDMDipilihUlang(3, 3, 'edm-inbox')).toBe(false)
  })
  it('ketuk berubah dan halaman edm-inbox → true', () => {
    expect(menuEDMDipilihUlang(3, 4, 'edm-inbox')).toBe(true)
    expect(menuEDMDipilihUlang(undefined, 1, 'edm-inbox')).toBe(true)
  })
  it('ketukMenu undefined → false', () => {
    expect(menuEDMDipilihUlang(3, undefined, 'edm-inbox')).toBe(false)
    expect(menuEDMDipilihUlang(undefined, undefined, 'edm-inbox')).toBe(false)
  })
  it('halaman lain → false', () => {
    expect(menuEDMDipilihUlang(3, 4, 'premiumlist')).toBe(false)
    expect(menuEDMDipilihUlang(3, 4, '')).toBe(false)
  })
})

describe('layar sesudah klik menu', () => {
  const buat: Layar = { jenis: 'buat' }
  const kasus: Layar = { jenis: 'kasus', id: 'UJI-EDM-1' }
  it('klik menu ulang dari Create Addendum dan dari kasus → kotak masuk', () => {
    expect(layarSesudahMenu(buat, 1, 2, 'edm-inbox')).toEqual({ jenis: 'inbox' })
    expect(layarSesudahMenu(kasus, 1, 2, 'edm-inbox')).toEqual({ jenis: 'inbox' })
  })
  it('perpindahan tanpa klik menu → layar bertahan', () => {
    expect(layarSesudahMenu(kasus, 2, 2, 'edm-inbox')).toBe(kasus)
    expect(layarSesudahMenu(buat, undefined, undefined, 'edm-inbox')).toBe(buat)
  })
  it('rute memakai ketukMenu + useRef + useEffect, tanpa impor modul lain', () => {
    expect(RUTE).toContain('export function RuteEDM({ halaman, ketukMenu }: PropsRute<HalamanEDM>)')
    expect(RUTE).toContain('const ketukLalu = useRef(ketukMenu)')
    expect(RUTE).toContain('if (menuEDMDipilihUlang(ketukLalu.current, ketukMenu, halaman)) {')
    expect(RUTE).toContain("setLayar({ jenis: 'inbox' })")
    expect(RUTE).not.toMatch(/from '[^']*premiumlistlife/)
    expect(RUTE).toContain('keputusan work owner 07-10-2026')
    // Hook dipanggil sebelum `return null` (aturan hook React).
    expect(RUTE.indexOf('useEffect(() => {')).toBeLessThan(RUTE.indexOf('return null'))
  })
})
