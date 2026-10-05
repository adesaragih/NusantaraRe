import { describe, expect, it } from 'vitest'

import { bolehUbahMenu, HAK_PENUH } from './hakMenu'

describe('hak menu View only (migrasi 914)', () => {
  it('hanya menu ber-hak LIHAT yang menolak tulis (juga superadmin, 05-10-2026); tanpa penyedia = penuh', () => {
    const h = { lihat: ['accounts'] }
    expect(bolehUbahMenu(h, 'accounts')).toBe(false)
    expect(bolehUbahMenu(h, 'aggregate')).toBe(true)
    expect(bolehUbahMenu(HAK_PENUH, 'accounts')).toBe(true)
  })
})
