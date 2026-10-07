// Aturan halaman awal Claim Prop (keputusan work owner 07-10-2026): dua tab Process / Resolve; tab Process memilih
// workbasket Admin (Outstanding Claim) atau Teknik (Input Acceptation); Add Claim hanya di workbasket Admin.

import { describe, expect, it } from 'vitest'

import { bolehTambahKlaim, jenisDaftar, TAB_INBOX, WORKBASKET } from './inbox'

describe('halaman awal Claim Prop', () => {
  it('dua tab saja, berurut Process lalu Resolve', () => {
    expect(TAB_INBOX).toEqual(['proses', 'selesai'])
  })

  it('workbasket Admin lalu Teknik', () => {
    expect(WORKBASKET).toEqual(['admin', 'teknik'])
  })

  it('tab dan workbasket menentukan daftar yang diminta ke server', () => {
    expect(jenisDaftar('proses', 'admin')).toBe('saya')
    expect(jenisDaftar('proses', 'teknik')).toBe('workbasket')
    expect(jenisDaftar('selesai', 'admin')).toBe('selesai')
    expect(jenisDaftar('selesai', 'teknik')).toBe('selesai')
  })

  it('Add Claim hanya di tab Process, workbasket Admin', () => {
    expect(bolehTambahKlaim('proses', 'admin')).toBe(true)
    expect(bolehTambahKlaim('proses', 'teknik')).toBe(false)
    expect(bolehTambahKlaim('selesai', 'admin')).toBe(false)
    expect(bolehTambahKlaim('selesai', 'teknik')).toBe(false)
  })
})
