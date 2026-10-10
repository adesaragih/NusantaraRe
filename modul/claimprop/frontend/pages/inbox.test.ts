// Aturan halaman awal Claim Prop (keputusan work owner 08-10-2026): dua tab Process / Resolve; tab Process bawaan =
// worklist pembuat tanpa cek workbasket (XML `ToCurrentOperator`); switch Teknik (Input Acceptation) hanya aktif bagi
// anggota ReasKlaimTeknik; Add Claim hanya saat switch Teknik mati.

import { describe, expect, it } from 'vitest'

import { bolehTambahKlaim, jenisDaftar, switchTeknikAktif, TAB_INBOX, tabelKomiteTampil } from './inbox'

describe('halaman awal Claim Prop', () => {
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
    expect(switchTeknikAktif({ workbasketTeknik: true, komite: false })).toBe(true)
    expect(switchTeknikAktif({ workbasketTeknik: false, komite: true })).toBe(false)
    expect(switchTeknikAktif(null)).toBe(false)
  })

  // Menu Komite Claim Prop dibuang (keputusan work owner 09-10-2026): tabel komite di bawah inbox hanya bagi pemegang
  // workbasket roster EMAILKOMITE PROP, tanpa switch dan tanpa ikut tab.
  it('tabel komite hanya bagi pemegang workbasket roster komite PROP', () => {
    expect(tabelKomiteTampil({ workbasketTeknik: false, komite: true })).toBe(true)
    expect(tabelKomiteTampil({ workbasketTeknik: true, komite: false })).toBe(false)
    expect(tabelKomiteTampil(null)).toBe(false)
  })

  it('Add Claim hanya di tab Process saat switch Teknik mati', () => {
    expect(bolehTambahKlaim('proses', false)).toBe(true)
    expect(bolehTambahKlaim('proses', true)).toBe(false)
    expect(bolehTambahKlaim('selesai', false)).toBe(false)
  })
})
