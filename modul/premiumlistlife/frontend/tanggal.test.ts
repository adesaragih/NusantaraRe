import { readdirSync, readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { isoDariTampil, tanggalJamTampil, tanggalTampil, topengTanggal } from './tanggal'

// Semua tanggal tampil dd/mm/yyyy — keputusan work owner 02-10-2026.

describe('tanggalTampil', () => {
  it('ISO dari server menjadi dd/mm/yyyy, dengan atau tanpa jam', () => {
    expect(tanggalTampil('2026-10-01')).toBe('01/10/2026')
    expect(tanggalTampil('2026-10-01T00:00:00Z')).toBe('01/10/2026')
    expect(tanggalTampil('2024-12-13 00:00:00')).toBe('13/12/2024')
  })

  it('kosong, null, dan tanggal nol Go menjadi kosong', () => {
    expect(tanggalTampil('')).toBe('')
    expect(tanggalTampil(null)).toBe('')
    expect(tanggalTampil(undefined)).toBe('')
    expect(tanggalTampil('0001-01-01T00:00:00Z')).toBe('')
  })

  it('teks yang bukan tanggal ISO dibiarkan apa adanya', () => {
    expect(tanggalTampil('31/10/2026')).toBe('31/10/2026')
    expect(tanggalTampil('UJI-PLAN')).toBe('UJI-PLAN')
    expect(tanggalTampil('1234.56')).toBe('1234.56')
  })

  it('jam ikut bila diminta', () => {
    expect(tanggalJamTampil('2026-09-30T10:15:00Z')).toBe('30/09/2026 10:15')
    expect(tanggalJamTampil('2026-09-30')).toBe('30/09/2026')
  })
})

describe('isian dd/mm/yyyy', () => {
  it('dd/mm/yyyy sah menjadi YYYY-MM-DD', () => {
    expect(isoDariTampil('13/12/2024')).toBe('2024-12-13')
    expect(isoDariTampil('29/02/2024')).toBe('2024-02-29')
  })

  it('bentuk salah atau tanggal yang tidak ada ditolak', () => {
    for (const t of ['31/02/2026', '29/02/2025', '00/01/2026', '01/13/2026', '1/2/2026', '2026-01-02', '13/12/24', '']) {
      expect(isoDariTampil(t)).toBeNull()
    }
  })

  it('topeng ketik menyisipkan garis miring dan hanya menerima angka', () => {
    expect(topengTanggal('13122024')).toBe('13/12/2024')
    expect(topengTanggal('131')).toBe('13/1')
    expect(topengTanggal('13/12/2024999')).toBe('13/12/2024')
    expect(topengTanggal('ab12')).toBe('12')
  })
})

describe('layar modul tanpa masukan tanggal bawaan yang terlihat', () => {
  it('type="date" hanya ada di IsianTanggal (sebagai kalender tersembunyi)', () => {
    const berkas = (d: string): string[] =>
      readdirSync(d, { withFileTypes: true }).flatMap((e) =>
        e.isDirectory() ? berkas(join(d, e.name)) : e.name.endsWith('.tsx') ? [join(d, e.name)] : [],
      )
    const pelanggar = berkas(__dirname).filter(
      (f) => !f.endsWith('IsianTanggal.tsx') && readFileSync(f, 'utf8').includes('type="date"'),
    )
    expect(pelanggar).toEqual([])
  })
})
