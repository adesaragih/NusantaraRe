// Asal: salinan `modul/nbtreatyin/frontend/ketikAngka.test.ts` (06-10-2026), kelas `nbti__` -> `edmt__`.
// Ketikan angka (permintaan work owner 06-10-2026: "inputannya decimal, selain decimal tidak bisa terketik; format
// angka langsung kedetek decimal; tidak harus 123456.678"). Tampilan selama mengetik mengikuti pola inti
// (titik ribuan, koma desimal - `sajian.ts`); yang dikirim ke backend selalu MENTAH (`123456.678`).

import { describe, expect, it } from 'vitest'

import { tampilKetik, ubahKetikan } from './ketikAngka'

/** Ketik `teks` satu karakter demi satu karakter di ujung isian, mulai dari kosong. */
function ketik(teks: string, awal = '') {
  let tampil = awal
  let mentah = ''
  for (const c of teks) {
    const r = ubahKetikan(tampil, tampil + c, tampil.length + 1)
    tampil = r.tampil
    mentah = r.mentah
  }
  return { tampil, mentah }
}

describe('tampilKetik - nilai mentah ke tampilan isian', () => {
  it('titik ribuan, koma desimal, tanpa pembulatan', () => {
    expect(tampilKetik('1234567.5')).toBe('1.234.567,5')
    expect(tampilKetik('-1000')).toBe('-1.000')
    expect(tampilKetik('0.125')).toBe('0,125')
    expect(tampilKetik('')).toBe('')
  })
  it('teks bukan angka apa adanya', () => {
    expect(tampilKetik('UJI')).toBe('UJI')
  })
})

describe('ubahKetikan - mengetik', () => {
  it('ribuan dipasang otomatis', () => {
    expect(ketik('1234567')).toEqual({ tampil: '1.234.567', mentah: '1234567' })
  })
  it('titik ATAU koma yang diketik = pemisah desimal', () => {
    expect(ketik('12345.67')).toEqual({ tampil: '12.345,67', mentah: '12345.67' })
    expect(ketik('12345,67')).toEqual({ tampil: '12.345,67', mentah: '12345.67' })
  })
  it('koma baru diketik: tampil ber-koma, mentah tanpa titik', () => {
    expect(ketik('12,')).toEqual({ tampil: '12,', mentah: '12' })
  })
  it('pemisah desimal kedua diabaikan', () => {
    expect(ketik('1,5,3')).toEqual({ tampil: '1,53', mentah: '1.53' })
    expect(ketik('1.5.3')).toEqual({ tampil: '1,53', mentah: '1.53' })
  })
  it('selain angka tidak terketik', () => {
    expect(ketik('12a3b-')).toEqual({ tampil: '123', mentah: '123' })
  })
  it('minus hanya di depan', () => {
    expect(ketik('-12')).toEqual({ tampil: '-12', mentah: '-12' })
    expect(ketik('1-2')).toEqual({ tampil: '12', mentah: '12' })
  })
  it('desimal diawali koma menjadi 0,', () => {
    expect(ketik(',5')).toEqual({ tampil: '0,5', mentah: '0.5' })
  })
  it('nol di depan dibuang', () => {
    expect(ketik('007')).toEqual({ tampil: '7', mentah: '7' })
    expect(ketik('0,05')).toEqual({ tampil: '0,05', mentah: '0.05' })
  })
  it('menghapus satu digit menata ulang ribuan', () => {
    // "12.345" -> hapus "3" -> "12.45"
    expect(ubahKetikan('12.345', '12.45', 3)).toMatchObject({ tampil: '1.245', mentah: '1245' })
  })
  it('kosong', () => {
    expect(ubahKetikan('1', '', 0)).toMatchObject({ tampil: '', mentah: '' })
  })
})

describe('ubahKetikan - menempel teks berformat', () => {
  const tempel = (t: string) => ubahKetikan('', t, t.length)
  it('format Indonesia dan format Inggris dikenali', () => {
    expect(tempel('123.456,678')).toMatchObject({ mentah: '123456.678', tampil: '123.456,678' })
    expect(tempel('123,456.678')).toMatchObject({ mentah: '123456.678', tampil: '123.456,678' })
    expect(tempel('163.125.000')).toMatchObject({ mentah: '163125000', tampil: '163.125.000' })
    expect(tempel('1,234,567')).toMatchObject({ mentah: '1234567' })
    expect(tempel('123456.678')).toMatchObject({ mentah: '123456.678' })
    expect(tempel('1,5')).toMatchObject({ mentah: '1.5' })
    expect(tempel(' IDR 2.500,75 ')).toMatchObject({ mentah: '2500.75' })
  })
})

describe('ubahKetikan - posisi kursor', () => {
  it('kursor tetap sesudah digit yang sama ketika titik ribuan bertambah', () => {
    // "1.234" kursor di ujung, ketik "5" -> "12.345", kursor di ujung
    expect(ubahKetikan('1.234', '1.2345', 6).kursor).toBe(6)
    // "1.234" sisip "9" sesudah "1" -> "19.234", kursor sesudah "9" (indeks 2)
    expect(ubahKetikan('1.234', '19.234', 2).kursor).toBe(2)
  })
})
