// Tab Object per halaman untuk case besar (tiket 75 F3, K75-3) - pemilihan mode, badan PATCH, aturan Save sama dengan TabObject.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { AMBANG_HALAMAN_OBJEK, type ObjekBerId } from '../api'
import { objekBaru } from './TabObject'
import { badanPatch, BATAS_TAMBAH_PATCH, BATAS_UBAH_PATCH, objekSalah } from './TabObjectHalaman'

const baca = (...p: string[]) => readFileSync(join(__dirname, ...p), 'utf8').replace(/\r\n/g, '\n')
const INWARD = baca('..', 'pages', 'InwardFacultative.tsx')
const SUMBER = baca('TabObjectHalaman.tsx')

const lokasi = (id: string, objectType = 'Gedung'): ObjekBerId => ({ ...objekBaru(id), objectType, id })

describe('TabObjectHalaman', () => {
  it('mode per halaman hanya untuk case > ambang (500, = backend); case biasa tetap TabObject', () => {
    expect(AMBANG_HALAMAN_OBJEK).toBe(500)
    expect(INWARD).toMatch(/\(dariServer\?\.jumlahLokasi \?\? 0\) > AMBANG_HALAMAN_OBJEK \? \(\s*<TabObjectHalaman caseId=\{kasus\.caseId\} insuredName=\{insured\} \/>\s*\) : \(\s*<TabObject caseId=\{kasus\.caseId\} insuredName=\{insured\} \/>/)
  })

  it('badan PATCH: hanya perubahan; lokasi yang ditandai hapus tidak ikut "ubah"', () => {
    const ubahan = new Map([
      ['1', lokasi('1')],
      ['2', lokasi('2')],
    ])
    const b = badanPatch(ubahan, [objekBaru('9')], new Set(['2']))
    expect(b.ubah.map((o) => o.id)).toEqual(['1'])
    expect(b.tambah).toHaveLength(1)
    expect(b.hapus).toEqual(['2'])
  })

  it('aturan Save = TabObject (Object Type wajib, dst.); batas per Save = backend', () => {
    expect(objekSalah(lokasi('1', ''))).toBe(true)
    expect(objekSalah(lokasi('1'))).toBe(false)
    expect([BATAS_UBAH_PATCH, BATAS_TAMBAH_PATCH]).toEqual([500, 500])
    expect(SUMBER).toContain('ambilObjekHalaman(caseId, halaman, GRID_OBJEK.ukuran, cari)')
    expect(SUMBER).toContain('const o = ubahan.get(asal.id) ?? asal')
  })
})
