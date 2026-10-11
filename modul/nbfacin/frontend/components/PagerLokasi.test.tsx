// 10 baris per halaman tab Spreading / Deduction (work owner 11-10-2026) - indeks asli dipertahankan, pager di atas tabel.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import PagerLokasi, { halamanAman, potongHalaman, UKURAN_HALAMAN_LOKASI } from './PagerLokasi'

const baca = (f: string) => readFileSync(join(__dirname, f), 'utf8').replace(/\r\n/g, '\n')

describe('potongHalaman / halamanAman', () => {
  const daftar = Array.from({ length: 28 }, (_, i) => `L${i + 1}`)

  it('10 per halaman dengan indeks ASLI (penyunting baris memakai posisi di seluruh daftar)', () => {
    expect(UKURAN_HALAMAN_LOKASI).toBe(10)
    expect(potongHalaman(daftar, 1)).toHaveLength(10)
    expect(potongHalaman(daftar, 2)[0]).toEqual(['L11', 10])
    expect(potongHalaman(daftar, 3)).toEqual(daftar.slice(20).map((x, i) => [x, 20 + i]))
  })

  it('halaman dijepit bila daftar memendek / kosong', () => {
    expect(halamanAman(28, 9)).toBe(3)
    expect(halamanAman(28, 0)).toBe(1)
    expect(halamanAman(0, 4)).toBe(1)
    expect(potongHalaman(daftar.slice(0, 5), 3)).toEqual(daftar.slice(0, 5).map((x, i) => [x, i]))
    expect(potongHalaman([], 2)).toEqual([])
  })
})

describe('PagerLokasi', () => {
  it('tidak tampil bila muat satu halaman; tampil "1–10 dari 28" bila lebih', () => {
    expect(renderToStaticMarkup(<PagerLokasi total={10} halaman={1} onPindah={() => {}} />)).toBe('')
    const h = renderToStaticMarkup(<PagerLokasi total={28} halaman={1} onPindah={() => {}} />)
    expect(h).toContain('28')
    expect(h).toContain('aria-current="page"')
  })

  it('dipasang di Spreading FIRE / Aneka / Cargo dan tiga daftar komisi Deduction', () => {
    expect(baca('TabSpreading.tsx')).toContain('{potongHalaman(v.lokasi, halLokasi).map(([l, o]) => [')
    expect(baca('TabSpreadingAneka.tsx')).toContain('{potongHalaman(v.lokasi, halLokasi).map(([l, a]) => {')
    expect(baca('TabSpreadingCargo.tsx')).toContain('{potongHalaman(v.cargo, halCargo).map(([x, ci]) => {')
    const d = baca('TabDeduction.tsx')
    for (const t of [
      '{potongHalaman(lokasi, hal).map(([l, a]) => [',
      '{potongHalaman(lokasi, hal).map(([l, a]) => {',
      '{potongHalaman(cargo, hal).map(([x, i]) => {',
    ]) {
      expect(d).toContain(t)
    }
    expect(d.match(/<PagerLokasi total=\{(lokasi|cargo)\.length\} halaman=\{hal\} onPindah=\{setHal\} \/>/g)).toHaveLength(3)
  })

  it('dipasang juga di tab Coverage FIRE / Aneka / Cargo (work owner 11-10-2026 "DI COVERAGE JUGA")', () => {
    expect(baca('TabCoverage.tsx')).toContain('{potongHalaman(objek, halLokasi).map(([x, o]) => [')
    expect(baca('TabCoverageAneka.tsx')).toContain('{potongHalaman(lokasi, halLokasi).map(([l, a]) => {')
    expect(baca('TabCoverageCargo.tsx')).toContain('{potongHalaman(cargo, halCargo).map(([x, i]) => [')
    for (const [f, total, hal] of [
      ['TabCoverage.tsx', 'objek', 'HalLokasi'],
      ['TabCoverageAneka.tsx', 'lokasi', 'HalLokasi'],
      ['TabCoverageCargo.tsx', 'cargo', 'HalCargo'],
    ] as const) {
      const h = hal.charAt(0).toLowerCase() + hal.slice(1)
      expect(baca(f)).toContain(`<PagerLokasi total={${total}.length} halaman={${h}} onPindah={set${hal}} />`)
    }
  })
})
