import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { KASUS_EDM } from '../labels'
import { JudulKasus } from './InputEDMLife'

// Chip Case ID di bawah judul Endorsement Life Detail - permintaan work owner 08-10-2026.

const SUMBER = readFileSync(join(__dirname, 'InputEDMLife.tsx'), 'utf8')

describe('chip Case ID di bawah judul', () => {
  it('judul lalu chip utama berisi id kasus', () => {
    const html = renderToStaticMarkup(<JudulKasus id="EDML-828" />)
    expect(html.indexOf(KASUS_EDM.judul)).toBeLessThan(html.indexOf('EDML-828'))
    expect(html).toContain('<span class="edm-kepala__chip edm-kepala__chip--utama">EDML-828</span>')
  })

  it('id kosong: tanpa baris chip', () => {
    expect(renderToStaticMarkup(<JudulKasus id="" />)).not.toContain('edm-kepala__chip')
  })

  it('layar kasus memakai JudulKasus dengan kasus.id; tombol kanan View Old Policy; tanpa impor premiumlistlife', () => {
    expect(SUMBER).toContain('<JudulKasus id={kasus.id} />')
    const kepala = SUMBER.slice(SUMBER.indexOf('<JudulKasus id={kasus.id} />'), SUMBER.indexOf('</header>'))
    expect(kepala).toContain('POLIS_LAMA_EDM.viewOldPolicy')
    const impor = SUMBER.split('\n').filter((l) => l.startsWith('import ') && l.includes('premiumlistlife'))
    expect(impor).toEqual([])
  })
})

describe('header kasus: judul+chip kiri, tombol rata kanan (perbaikan 08-10-2026)', () => {
  const CSS = readFileSync(join(__dirname, '..', 'endorsementlife.css'), 'utf8').replace(/\r\n/g, '\n')
  const aturan = (pemilih: string) => {
    const i = CSS.indexOf(pemilih + ' {')
    expect(i, pemilih).toBeGreaterThanOrEqual(0)
    return CSS.slice(i, CSS.indexOf('}', i))
  }

  it('header: dua anak langsung - blok judul (kiri) lalu blok tombol (kanan)', () => {
    const awal = SUMBER.indexOf('<header className="inbox__kepala edm-kasus__kepala">')
    const kepala = SUMBER.slice(awal, SUMBER.indexOf('</header>', awal))
    expect(kepala.indexOf('<JudulKasus id={kasus.id} />')).toBeGreaterThan(0)
    expect(kepala.indexOf('<JudulKasus')).toBeLessThan(kepala.indexOf('<div className="edm-aksi">'))
    const html = renderToStaticMarkup(<JudulKasus id="EDML-828" />)
    expect(html.startsWith('<div class="edm-kepala__judul"><h2 class="inbox__judul">')).toBe(true)
  })

  it('CSS: space-between, blok kiri mengisi ruang, tombol margin-left:auto dan tidak menyusut', () => {
    expect(aturan('.edm-kasus__kepala')).toContain('justify-content: space-between')
    expect(aturan('.edm-kasus__kepala')).toContain('align-items: flex-start')
    expect(aturan('.edm-kasus__kepala > .edm-kepala__judul')).toMatch(/flex: 1 1 [^;]+;/)
    expect(aturan('.edm-kasus__kepala > .edm-kepala__judul')).toContain('min-width: 0')
    expect(aturan('.edm-kasus__kepala > .edm-aksi')).toContain('margin-left: auto')
    expect(aturan('.edm-kasus__kepala > .edm-aksi')).toContain('flex-shrink: 0')
  })
})
