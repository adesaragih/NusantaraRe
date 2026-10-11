// Isian tanggal NB Treaty In bisa diketik dan ditempel (work owner 10-10-2026, Statement Period / To): helper
// `ketikTanggal.ts` dan pemasangannya di `KotakMedan` - bukan lagi `<input type="date">` bawaan.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { halamanTanggal, rapikanTanggal, tampilTanggal, tanggalUtuh } from './ketikTanggal'

describe('mengetik', () => {
  it('cukup angka; pemisah disisipkan, paling panjang dd-mm-yyyy', () => {
    expect(['0', '01', '010', '0107', '01072', '01072024', '010720245'].map(rapikanTanggal)).toEqual([
      '0',
      '01',
      '01-0',
      '01-07',
      '01-07-2',
      '01-07-2024',
      '01-07-2024',
    ])
    expect(rapikanTanggal('01-07-')).toBe('01-07')
  })

  it('nilai halaman hanya untuk tanggal lengkap yang ada', () => {
    expect(halamanTanggal('01-07-2024')).toBe('2024-07-01')
    expect(halamanTanggal('29-02-2024')).toBe('2024-02-29')
    expect(['29-02-2025', '31-04-2024', '01-13-2024', '01-07-202', ''].map(halamanTanggal)).toEqual([
      '',
      '',
      '',
      '',
      '',
    ])
  })
})

describe('menempel', () => {
  it('format umum dibaca utuh; hari dulu untuk d/m/yyyy (format aplikasi)', () => {
    expect(
      ['07/01/2024', '7/1/2024', '1.7.2024', '2024-07-01', '2024-07-01 00:00:00', '2024/7/1', ' 01-07-2024\n'].map(
        rapikanTanggal,
      ),
    ).toEqual(['07-01-2024', '07-01-2024', '01-07-2024', '01-07-2024', '01-07-2024', '01-07-2024', '01-07-2024'])
    expect(['1 Jul 2024', '01-Jul-2024', '1 Agustus 2024', '15 Okt 2026'].map(rapikanTanggal)).toEqual([
      '01-07-2024',
      '01-07-2024',
      '01-08-2024',
      '15-10-2026',
    ])
    expect(tanggalUtuh('teks bebas')).toBeNull()
    expect(rapikanTanggal('tgl 01072024')).toBe('01-07-2024')
  })

  it('nilai halaman tampil dd-mm-yyyy (yang disalin dari kotak)', () => {
    expect(tampilTanggal('2024-07-01')).toBe('01-07-2024')
    expect(tampilTanggal('2024-07-01 00:00:00')).toBe('01-07-2024')
    expect(tampilTanggal('')).toBe('')
  })
})

describe('KotakMedan memakai InputTanggal', () => {
  it('bukan input date bawaan; tempelan mengganti isi; aksi saat fokus meninggalkan kotak + tombol kalender', () => {
    const kotak = readFileSync(join(__dirname, 'components', 'KotakMedan.tsx'), 'utf8')
    const isian = readFileSync(join(__dirname, 'components', 'InputTanggal.tsx'), 'utf8')
    expect(kotak).not.toContain('type="date"')
    expect(kotak).toContain(
      '<InputTanggal id={idIsian} value={v} galat={!!pesan} onChange={(x) => onUbah(medan.jalur, x)} />',
    )
    expect(kotak).toContain(
      'if (!e.currentTarget.contains(e.relatedTarget as Node | null)) onSelesai(medan, nilai(halaman, medan.jalur))',
    )
    expect(isian).toMatch(/type="text"\s+inputMode="numeric"/)
    expect(isian).toContain("ketik(e.clipboardData.getData('text'))")
  })
})
