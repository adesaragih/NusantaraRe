import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { menuDipilihUlang } from './rute'

// Memilih menu PremiumList Life (lagi) kembali ke kotak masuk.

describe('menu dipilih ulang', () => {
  it('ketukan baru di halaman PremiumList — kembali ke kotak masuk', () => {
    expect(menuDipilihUlang(3, 4, 'premiumlist')).toBe(true)
  })

  it('tanpa ketukan baru — kasus yang terbuka tetap', () => {
    expect(menuDipilihUlang(4, 4, 'premiumlist')).toBe(false)
  })

  it('ketukan menu modul lain — keadaan polis tidak disentuh', () => {
    expect(menuDipilihUlang(4, 5, 'beranda')).toBe(false)
  })

  it('Shell tanpa ketukMenu — perilaku lama', () => {
    expect(menuDipilihUlang(undefined, undefined, 'premiumlist')).toBe(false)
  })

  it('App meneruskan ketukMenu dari pilihan menu', () => {
    const app = readFileSync(join(__dirname, '../../../frontend/App.tsx'), 'utf8')
    expect(app).toContain('pilihDariMenu(h)')
    expect(app).toContain('ketukMenu={ketukMenu}')
  })
})
