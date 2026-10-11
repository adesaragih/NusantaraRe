// Uji aturan `DataTransform/TreatyInSetTreatyYear.xml` beserta bentuk kabel
// tanggal `FieldTanggal`.
//
// ⛔ Angka dan bentuk di berkas ini dari ekspor Pega 2026-09 dan pengukuran
// POOLDATA 5 Oktober 2026 — bukan dari dugaan.

import { describe, expect, it } from 'vitest'

import { akhirSetahunSesudah, keKabel, keSimpan, tahunDariMulai } from './components/tanggalIso'

describe('bentuk kabel FieldTanggal', () => {
  it('`YYYYMMDD` tersimpan menjadi `DD-MM-YYYY`', () => {
    expect(keKabel('20250101')).toBe('01-01-2025')
    expect(keKabel('20261231')).toBe('31-12-2026')
  })

  // ⛔ Cacat pertama: bentuk BACA disuapkan ke medan tanggal, ditolak tanpa
  // bersuara, dan medannya terlihat kosong pada kontrak yang tanggalnya ada.
  it('⛔ bentuk BACA ditolak — ia memang bukan nilai', () => {
    expect(keKabel('01/01/2025')).toBe('')
    expect(keKabel('01/01/25')).toBe('')
  })

  it('teks aneh menjadi kosong, bukan diteruskan', () => {
    for (const v of ['', '2025', '20250101T075400.000 GMT', 'abcdefgh']) {
      expect(keKabel(v)).toBe('')
    }
  })
})

describe('TreatyInSetTreatyYear — disalin dari ekspor', () => {
  // ⛔ UJI YANG MENUTUP KELUHAN 6 Oktober 2026.
  //
  // `FieldTanggal` mengembalikan `DD-MM-YYYY`; bentuk sebelumnya hanya
  // menerima `YYYY-MM-DD`, sehingga memilih Commencement membuat `Treaty Year`
  // DAN `Termination` KOSONG — kedua turunan menolak bentuk yang mereka terima.
  it('⛔ bentuk KABEL `DD-MM-YYYY` diterima — inilah yang dikembalikan medannya', () => {
    expect(tahunDariMulai('01-01-2026')).toBe('2026')
    expect(akhirSetahunSesudah('01-01-2026')).toBe('01-01-2027')
  })

  // Bentuk ISO tetap diterima: keadaan awal pernah berbentuk itu, dan fungsi
  // yang benar hanya bila pemanggilnya ingat bentuknya akan salah suatu hari.
  it('bentuk ISO juga diterima, dan jawabannya tetap bentuk kabel', () => {
    expect(tahunDariMulai('2026-01-01')).toBe('2026')
    expect(akhirSetahunSesudah('2026-07-01')).toBe('01-07-2027')
  })

  it('⛔ masukan belum lengkap tidak melahirkan nilai', () => {
    for (const v of ['', '2025', '2025-01', '01-01', 'bukan-tanggal']) {
      expect(tahunDariMulai(v)).toBe('')
      expect(akhirSetahunSesudah(v)).toBe('')
    }
  })

  // Rantai penuh, persis seperti di layar: tersimpan → kabel → kedua turunan.
  it('⭐ rantai penuh: tersimpan 20260101 → Treaty Year 2026, Termination 2027', () => {
    const kabel = keKabel('20260101')
    expect(kabel).toBe('01-01-2026')
    expect(tahunDariMulai(kabel)).toBe('2026')
    expect(akhirSetahunSesudah(kabel)).toBe('01-01-2027')
  })
})

describe('keSimpan — bentuk tersimpan YYYYMMDD', () => {
  it('kabel dan ISO menjadi YYYYMMDD; sisanya kosong', () => {
    expect(keSimpan('01-04-2023')).toBe('20230401')
    expect(keSimpan('2023-04-01')).toBe('20230401')
    expect(keSimpan('01/04/2023')).toBe('')
    expect(keSimpan('')).toBe('')
  })
})
