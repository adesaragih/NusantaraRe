import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import type { KolomKorpus } from '../kolomKorpus'
import { UMUM_EDM } from '../labels'
import { UKURAN_HALAMAN_EDM } from '../tampilan'
import { nomorAwalPeserta } from './PolisLama'
import { TabelKorpus } from './TabelKorpus'

// Kolom nomor baris "No" tabel peserta popup View Old Policy - keputusan work owner 07-10-2026.

const POLIS_LAMA = readFileSync(join(__dirname, 'PolisLama.tsx'), 'utf8')
const KOLOM: readonly KolomKorpus[] = [{ label: 'UJI-KOLOM', kolom: 'NAMA', jenis: 't', baris: 1 } as KolomKorpus]
const BARIS = [{ NAMA: 'UJI-A' }, { NAMA: 'UJI-B' }]

function sel(html: string, tag: 'th' | 'td'): string[] {
  return [...html.matchAll(new RegExp(`<${tag}[^>]*>([^<]*)</${tag}>`, 'g'))].map((m) => m[1] ?? '')
}

describe('TabelKorpus nomorAwal', () => {
  it('tanpa nomorAwal: tidak ada kolom No (pemakai lain tidak berubah)', () => {
    const html = renderToStaticMarkup(<TabelKorpus kolom={KOLOM} baris={BARIS} kunci={(_, i) => String(i)} />)
    expect(sel(html, 'th')).toEqual(['UJI-KOLOM'])
    expect(sel(html, 'td')).toEqual(['UJI-A', 'UJI-B'])
    expect(html).not.toContain('edm-nomor')
  })

  it('dengan nomorAwal: kolom No paling kiri mulai dari nilai itu', () => {
    const html = renderToStaticMarkup(<TabelKorpus kolom={KOLOM} baris={BARIS} kunci={(_, i) => String(i)} nomorAwal={21} />)
    expect(sel(html, 'th')).toEqual([UMUM_EDM.nomor, 'UJI-KOLOM'])
    expect(sel(html, 'td')).toEqual(['21', 'UJI-A', '22', 'UJI-B'])
    expect(UMUM_EDM.nomor).toBe('No')
  })
})

describe('PolisLama: nomor peserta berlanjut antar halaman', () => {
  it('nomorAwal = (halaman - 1) × ukuran + 1', () => {
    expect(nomorAwalPeserta(1)).toBe(1)
    expect(nomorAwalPeserta(2)).toBe(UKURAN_HALAMAN_EDM + 1)
    expect(nomorAwalPeserta(2, 10)).toBe(11)
    expect(nomorAwalPeserta(3, 10)).toBe(21)
    expect(nomorAwalPeserta(0, 10)).toBe(1)
  })

  it('tabel peserta bernomor; tabel rekap tanpa nomor', () => {
    expect(POLIS_LAMA).toContain('nomorAwal={nomorAwalPeserta(halaman)}')
    expect(POLIS_LAMA.split('nomorAwal=').length - 1).toBe(1)
    const rekap = POLIS_LAMA.slice(POLIS_LAMA.indexOf('kolom={POLIS_LAMA_REKAP[t]}'))
    expect(rekap.slice(0, rekap.indexOf('/>'))).not.toContain('nomorAwal')
  })
})
