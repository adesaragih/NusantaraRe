import { describe, expect, it } from 'vitest'

import type { Keputusan } from './api'
import { isianKirim, periksa, PESAN_KOSONG, tampil, tampilCatatanSubjectivity, tampilSubjectivity } from './nilai'

const kosong: Keputusan = {
  acceptStatus: '',
  comment: '',
  isSubjectivity: false,
  subjectivityNote: '',
  usulTutup: false,
  usulCadang: false,
}

describe('tampil - angka 4 desimal pemisah Indonesia, tanggal dd-mm-yyyy', () => {
  it('angka', () => {
    expect(tampil('1234567.5', 'angka')).toBe('1.234.567,5')
    expect(tampil('0.123456', 'angka')).toBe('0,1235')
    expect(tampil('', 'angka')).toBe('')
  })
  it('tanggal dan tanggal-jam', () => {
    expect(tampil('2026-09-15', 'tanggal')).toBe('15-09-2026')
    expect(tampil('2026-10-08 10:05:00', 'tanggalJam')).toBe('08-10-2026 10:05')
    expect(tampil('2026-10-08T10:05:00+07:00', 'tanggalJam')).toBe('08-10-2026 10:05')
    expect(tampil('UJI', 'tanggal')).toBe('UJI')
  })
  it('teks apa adanya', () => {
    expect(tampil(' UJI ', 'teks')).toBe('UJI')
  })
})

describe('isian ShowTransfer', () => {
  it('Subjectivity tampil hanya saat AcceptStatus 1; catatannya hanya bila dicentang', () => {
    expect(tampilSubjectivity({ ...kosong, acceptStatus: '1' })).toBe(true)
    expect(tampilSubjectivity({ ...kosong, acceptStatus: '2' })).toBe(false)
    expect(tampilCatatanSubjectivity({ ...kosong, acceptStatus: '1', isSubjectivity: true })).toBe(true)
    expect(tampilCatatanSubjectivity({ ...kosong, acceptStatus: '2', isSubjectivity: true })).toBe(false)
  })
  it('wajib isi: AcceptStatus, Note, dan Subjectivity Note bila dicentang di tingkat 1', () => {
    expect(periksa(kosong, true)).toEqual({ acceptStatus: PESAN_KOSONG, comment: PESAN_KOSONG })
    const subj = { ...kosong, acceptStatus: '1', comment: 'UJI', isSubjectivity: true }
    expect(periksa(subj, true)).toEqual({ subjectivityNote: PESAN_KOSONG })
    expect(periksa(subj, false)).toEqual({})
    expect(periksa({ ...subj, subjectivityNote: 'UJI' }, true)).toEqual({})
  })
  it('isian bertingkat-1 tidak terkirim di tingkat lain atau saat tersembunyi', () => {
    const k = {
      ...kosong,
      acceptStatus: '2',
      comment: 'UJI',
      isSubjectivity: true,
      subjectivityNote: 'X',
      usulTutup: true,
    }
    expect(isianKirim(k, true)).toEqual({ ...k, isSubjectivity: false, subjectivityNote: '' })
    expect(isianKirim({ ...k, acceptStatus: '1' }, false)).toEqual({
      ...k,
      acceptStatus: '1',
      isSubjectivity: false,
      subjectivityNote: '',
      usulTutup: false,
    })
  })
})
