import { readFileSync } from 'node:fs'

import { afterEach, describe, expect, it, vi } from 'vitest'

import { ambil, ambilDaftar, ambilPilihan, tambah, ubah } from './api'
import { isianDari, kodeNama, opsiTreaty, periksaIsian, rapikanIsian, syariah } from './aturan'
import { BG, MENU_BG } from './labels'

describe('label Business Group', () => {
  it('nama menu = M_NAV_MENU.LABEL migrasi inti 918, golongan MASTER', () => {
    const sql = readFileSync(`${__dirname}/../../../inti/backend/migrations/918_m_nav_menu_businessgroup.sql`, 'utf8')
    expect(sql).toContain(`'businessgroup', '${MENU_BG.kelompok}', 'MASTER', 'businessgroup'`)
  })
})

describe('aturan Business Group', () => {
  it('isian; spasi tepi dibuang; Treaty Group dan Name wajib; SYARIAH ditolak', () => {
    expect(isianDari(null)).toEqual({ topId: '', name: '', alias: '' })
    const b = { id: '10013', name: 'UJI FIRE', alias: 'UJI F', topId: '10007', treatyName: 'UJI PROPERTY' }
    expect(isianDari(b)).toEqual({ topId: '10007', name: 'UJI FIRE', alias: 'UJI F' })
    expect(rapikanIsian({ topId: ' 10007 ', name: ' a ', alias: ' b ' })).toEqual({
      topId: '10007',
      name: 'a',
      alias: 'b',
    })
    expect(periksaIsian({ topId: '', name: 'X', alias: '' })).toBe(BG.galatTreaty)
    expect(periksaIsian({ topId: '10007', name: ' ', alias: '' })).toBe(BG.galatNama)
    expect(periksaIsian({ topId: '10007', name: 'uji motor syariah ', alias: '' })).toBe(BG.galatSyariah)
    expect(periksaIsian({ topId: '10007', name: 'UJI SYARIAH MOTOR', alias: '' })).toBeNull()
    expect(syariah('fire Syariah')).toBe(true)
  })

  it('opsi Treaty Group; TOPID yatim tetap terbaca sebagai nilai lama', () => {
    const t = [{ id: '10007', name: 'UJI PROPERTY' }]
    expect(opsiTreaty(t, '10007')).toEqual([{ value: '10007', label: '10007 - UJI PROPERTY' }])
    expect(opsiTreaty(t, '10019', 'UJI HILANG')).toContainEqual({
      value: '10019',
      label: BG.nilaiLama('10019 - UJI HILANG'),
    })
    expect(kodeNama('10007', '')).toBe('10007')
  })
})

describe('klien Business Group', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('rute dan metode sama dengan handlers/rute.go; tanpa hapus', async () => {
    const panggil: string[] = []
    vi.stubGlobal(
      'fetch',
      vi.fn((url: string, init: RequestInit) => {
        panggil.push(`${String(init.method ?? 'GET')} ${url.replace(/^https?:\/\/[^/]+/, '')}`)
        return Promise.resolve(new Response('{}', { status: 200, headers: { 'Content-Type': 'application/json' } }))
      }),
    )
    await ambilDaftar(' uji ', '10007')
    await ambilDaftar('', '')
    await ambilPilihan()
    await ambil('10013')
    await tambah({ topId: '10007', name: 'UJI', alias: '' })
    await ubah('10013', { topId: '10007', name: 'UJI', alias: '' })
    expect(panggil).toEqual([
      'GET /api/business-group?q=uji&treatyGroup=10007',
      'GET /api/business-group',
      'GET /api/business-group/pilihan',
      'GET /api/business-group/10013',
      'POST /api/business-group',
      'PUT /api/business-group/10013',
    ])
  })
})
