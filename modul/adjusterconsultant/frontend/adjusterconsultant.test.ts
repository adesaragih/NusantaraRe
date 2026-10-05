import { readFileSync } from 'node:fs'

import { afterEach, describe, expect, it, vi } from 'vitest'

import { ambil, ambilDaftar, setelAktif, tambah, ubah } from './api'
import { isianDari, namaTampil, periksaIsian, rapikanIsian, STATUS_AWAL, URUTAN_STATUS } from './aturan'
import { ADJ, MENU_ADJ } from './labels'

describe('label Adjuster Consultant', () => {
  it('nama menu = M_NAV_MENU.LABEL migrasi inti 915, golongan MASTER', () => {
    const sql = readFileSync(`${__dirname}/../../../inti/backend/migrations/915_m_nav_menu_adjusterconsultant.sql`, 'utf8')
    expect(sql).toContain(`'adjusterconsultant', '${MENU_ADJ.kelompok}', 'MASTER', 'adjusterconsultant'`)
  })

  it('caption kolom layar Pega MstAdjusterConsultant', () => {
    expect([ADJ.id, ADJ.nama, ADJ.alamat, ADJ.telp, ADJ.tglUbah]).toEqual(['ID', 'Name', 'Address', 'Telp No', 'Edit Date'])
  })
})

describe('aturan Adjuster Consultant', () => {
  it('saringan status: bawaan Active, lalu Inactive dan All', () => {
    expect(STATUS_AWAL).toBe('active')
    expect(URUTAN_STATUS).toEqual(['active', 'inactive', ''])
  })

  it('isian: kosong untuk Add, terisi untuk Edit; spasi tepi dibuang; Name wajib', () => {
    expect(isianDari(null)).toEqual({ name: '', address: '', telpNo: '' })
    const a = { id: '10001', name: 'UJI', address: 'JL', telpNo: '1', username: 'U', editDate: '', active: true }
    expect(isianDari(a)).toEqual({ name: 'UJI', address: 'JL', telpNo: '1' })
    expect(rapikanIsian({ name: ' uji ', address: ' jl ', telpNo: ' 1 ' })).toEqual({ name: 'uji', address: 'jl', telpNo: '1' })
    expect(periksaIsian({ name: '  ', address: '', telpNo: '' }, ADJ.galatNama)).toBe(ADJ.galatNama)
    expect(periksaIsian({ name: 'UJI', address: '', telpNo: '' }, ADJ.galatNama)).toBeNull()
    expect(namaTampil(a)).toBe('UJI (10001)')
  })
})

describe('klien Adjuster Consultant', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('rute dan metode sama dengan handlers/rute.go', async () => {
    const panggil: string[] = []
    vi.stubGlobal(
      'fetch',
      vi.fn((url: string, init: RequestInit) => {
        panggil.push(`${String(init.method ?? 'GET')} ${url.replace(/^https?:\/\/[^/]+/, '')}`)
        return Promise.resolve(new Response('{}', { status: 200, headers: { 'Content-Type': 'application/json' } }))
      }),
    )
    await ambilDaftar(' uji ', 'active')
    await ambilDaftar('', '')
    await ambil('10001')
    await tambah({ name: 'UJI', address: '', telpNo: '' })
    await ubah('10001', { name: 'UJI', address: '', telpNo: '' })
    await setelAktif('10001', false)
    expect(panggil).toEqual([
      'GET /api/adjuster-consultant?q=uji&status=active',
      'GET /api/adjuster-consultant',
      'GET /api/adjuster-consultant/10001',
      'POST /api/adjuster-consultant',
      'PUT /api/adjuster-consultant/10001',
      'POST /api/adjuster-consultant/10001/aktif',
    ])
  })
})
