// Disalin dari `modul/claimnonprop/frontend/pages/inbox.test.ts` (pola, bukan impor; asal Claim Prop): keputusan work
// owner yang disebut di bawah adalah keputusan layar Claim Prop yang ditiru Claim Fac In (prompt Claim Fac In tahap 1).
// Aturan halaman awal: dua tab Process / Resolve; tab Process bawaan = worklist pembuat (Input Register / Input
// Estimasi); switch Teknik (Choose Surveyor) hanya aktif bagi anggota ReasKlaimTeknik; Add Claim hanya saat switch
// Teknik mati.

import { describe, expect, it } from 'vitest'

import { bolehTambahKlaim, jenisDaftar, switchTeknikAktif, TAB_INBOX, tabelKomiteTampil } from './inbox'

describe('halaman awal Claim Fac In', () => {
  it('dua tab saja, berurut Process lalu Resolve', () => {
    expect(TAB_INBOX).toEqual(['proses', 'selesai'])
  })

  it('switch Teknik mati = worklist sendiri; nyala = workbasket Choose Surveyor; Resolve = selesai', () => {
    expect(jenisDaftar('proses', false)).toBe('saya')
    expect(jenisDaftar('proses', true)).toBe('workbasket')
    expect(jenisDaftar('selesai', false)).toBe('selesai')
    expect(jenisDaftar('selesai', true)).toBe('selesai')
  })

  it('switch Teknik hanya dapat dinyalakan anggota ReasKlaimTeknik', () => {
    expect(switchTeknikAktif({ workbasketSurveyor: true, komite: false })).toBe(true)
    expect(switchTeknikAktif({ workbasketSurveyor: false, komite: true })).toBe(false)
    expect(switchTeknikAktif(null)).toBe(false)
  })

  // Modul Komite Claim Fac In tanpa menu (prompt tahap 2 §2 butir 2, pola Claim Non Prop): tabel komite di bawah inbox
  // hanya bagi anggota roster komite FACIN, tanpa switch dan tanpa ikut tab.
  it('tabel komite hanya bagi anggota roster komite FACIN', () => {
    expect(tabelKomiteTampil({ workbasketSurveyor: false, komite: true })).toBe(true)
    expect(tabelKomiteTampil({ workbasketSurveyor: true, komite: false })).toBe(false)
    expect(tabelKomiteTampil(null)).toBe(false)
  })

  it('Add Claim hanya di tab Process saat switch Teknik mati', () => {
    expect(bolehTambahKlaim('proses', false)).toBe(true)
    expect(bolehTambahKlaim('proses', true)).toBe(false)
    expect(bolehTambahKlaim('selesai', false)).toBe(false)
  })
})
