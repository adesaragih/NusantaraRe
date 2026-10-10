// Tombol `Add Revision` (@500688) dan `Add Adjustment Premium` (@515563) —
// hanya workbasket ReasTreatyInAdmin (diterapkan 9 Oktober 2026).

import { describe, expect, it } from 'vitest'

import { tambahTampil } from './pages/PenyesuaianKontrak'

describe('tombol tambah daftar Adjustment', () => {
  it('tampil hanya bagi pemegang ReasTreatyInAdmin', () => {
    expect(tambahTampil(['ReasTreatyInAdmin'])).toBe(true)
    expect(tambahTampil(['ReasTreatyInSecHead', 'ReasTreatyInAdmin'])).toBe(true)
    expect(tambahTampil(['ReasTreatyInSecHead'])).toBe(false)
    expect(tambahTampil([])).toBe(false)
  })
})
