// Disalin dari `modul/claimprop/frontend/pages/inbox.test.ts` (pola, bukan impor): keputusan work owner yang disebut di bawah
// adalah keputusan layar Claim Prop yang ditiru Claim Non Prop (prompt Claim Non Prop tahap 1).
// Aturan halaman awal Claim Prop (keputusan work owner 08-10-2026): dua tab Process / Resolve; tab Process bawaan =
// worklist pembuat tanpa cek workbasket (XML `ToCurrentOperator`); switch Teknik (Input Acceptation) hanya aktif bagi
// anggota ReasKlaimTeknik; Add Claim hanya saat switch Teknik mati.

import { describe, expect, it } from 'vitest'

import { bolehTambahKlaim, jenisDaftar, switchTeknikAktif, TAB_INBOX } from './inbox'

describe('halaman awal Claim Non Prop', () => {
  it('dua tab saja, berurut Process lalu Resolve', () => {
    expect(TAB_INBOX).toEqual(['proses', 'selesai'])
  })

  it('switch Teknik mati = worklist sendiri; nyala = workbasket Teknik; Resolve = selesai', () => {
    expect(jenisDaftar('proses', false)).toBe('saya')
    expect(jenisDaftar('proses', true)).toBe('workbasket')
    expect(jenisDaftar('selesai', false)).toBe('selesai')
    expect(jenisDaftar('selesai', true)).toBe('selesai')
  })

  it('switch Teknik hanya dapat dinyalakan anggota ReasKlaimTeknik', () => {
    expect(switchTeknikAktif({ workbasketTeknik: true })).toBe(true)
    expect(switchTeknikAktif({ workbasketTeknik: false })).toBe(false)
    expect(switchTeknikAktif(null)).toBe(false)
  })

  it('Add Claim hanya di tab Process saat switch Teknik mati', () => {
    expect(bolehTambahKlaim('proses', false)).toBe(true)
    expect(bolehTambahKlaim('proses', true)).toBe(false)
    expect(bolehTambahKlaim('selesai', false)).toBe(false)
  })
})
