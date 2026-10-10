// Disalin dari `modul/claimnonprop/frontend/ketikTanggal.test.ts` (pola, bukan impor; asal Claim Prop): keputusan work
// owner yang disebut di bawah adalah keputusan layar Claim Prop yang ditiru Claim Fac In (prompt Claim Fac In tahap 1).
// Isian tanggal Claim Prop (work owner 08-10-2026: "format tanggalnya kenapa susah banget di ketik ya? perbaiki"): cukup
// ketik angka, pemisah otomatis, tampil dd-mm-yyyy; nilai halaman tetap YYYY-MM-DD / YYYY-MM-DD HH:MM:SS.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { halamanTanggal, rapikanTanggal, tampilTanggal } from './ketikTanggal'

describe('tampilTanggal - nilai halaman ke tampilan', () => {
  it('tanggal: dd-mm-yyyy dari bentuk halaman maupun bentuk lama', () => {
    expect(tampilTanggal('2026-05-01', 'tanggal')).toBe('01-05-2026')
    expect(tampilTanggal('2026-05-01 00:00:00', 'tanggal')).toBe('01-05-2026')
    expect(tampilTanggal('20260501', 'tanggal')).toBe('01-05-2026')
    expect(tampilTanggal('', 'tanggal')).toBe('')
  })
  it('tanggal-waktu: dd-mm-yyyy hh:mm', () => {
    expect(tampilTanggal('2026-05-01 08:30:00', 'tanggal-waktu')).toBe('01-05-2026 08:30')
    expect(tampilTanggal('2026-05-01', 'tanggal-waktu')).toBe('01-05-2026')
  })
  it('bukan tanggal apa adanya', () => {
    expect(tampilTanggal('UJI', 'tanggal')).toBe('UJI')
  })
})

describe('rapikanTanggal - mengetik', () => {
  it('cukup angka; pemisah disisipkan otomatis', () => {
    expect(rapikanTanggal('0105', 'tanggal')).toBe('01-05')
    expect(rapikanTanggal('01052026', 'tanggal')).toBe('01-05-2026')
    expect(rapikanTanggal('010520260830', 'tanggal-waktu')).toBe('01-05-2026 08:30')
  })
  it('huruf dan pemisah ketikan sendiri tidak mengganggu; kelebihan digit dibuang', () => {
    expect(rapikanTanggal('01/05/2026', 'tanggal')).toBe('01-05-2026')
    expect(rapikanTanggal('01a05b2026999', 'tanggal')).toBe('01-05-2026')
  })
  it('menghapus pemisah terakhir tidak menyangkut', () => {
    expect(rapikanTanggal('01-', 'tanggal')).toBe('01')
  })
})

describe('halamanTanggal - isian ke nilai halaman', () => {
  it('lengkap dan ada', () => {
    expect(halamanTanggal('01-05-2026', 'tanggal')).toBe('2026-05-01')
    expect(halamanTanggal('29-02-2024', 'tanggal')).toBe('2024-02-29')
    expect(halamanTanggal('01-05-2026 08:30', 'tanggal-waktu')).toBe('2026-05-01 08:30:00')
  })
  it('setengah jadi atau tanggal yang tidak ada = kosong', () => {
    expect(halamanTanggal('01-05', 'tanggal')).toBe('')
    expect(halamanTanggal('31-02-2026', 'tanggal')).toBe('')
    expect(halamanTanggal('01-13-2026', 'tanggal')).toBe('')
    expect(halamanTanggal('01-05-2026', 'tanggal-waktu')).toBe('')
    expect(halamanTanggal('01-05-2026 24:00', 'tanggal-waktu')).toBe('')
  })
})

describe('semua isian tanggal Claim Fac In memakai InputTanggal', () => {
  it('kendali "tanggal" dan "tanggal-waktu" di TataView dirender InputTanggal, bukan <input type="date"> polos', () => {
    const src = readFileSync(join(__dirname, 'components', 'TataView.tsx'), 'utf8')
    const medan = src.slice(src.indexOf('function Medan('), src.indexOf('function Tombol('))
    expect(medan).toContain('<InputTanggal')
    expect(medan).not.toContain('type="date"')
    expect(medan).not.toContain('type="datetime-local"')
  })
})
