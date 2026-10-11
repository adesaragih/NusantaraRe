import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { POLIS_LAMA_PESERTA, POLIS_LAMA_REKAP, type KolomKorpus } from '../kolomKorpus'
import { UMUM_EDM } from '../labels'
import { UKURAN_HALAMAN_EDM } from '../tampilan'
import {
  barisRekapPolisLama,
  KOLOM_POLIS_LAMA_TERSEMBUNYI,
  kolomPolisLamaTampil,
  kolomRekapPolisLama,
  nomorAwalPeserta,
  teksEndorsementNo,
} from './PolisLama'
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
    const rekap = POLIS_LAMA.slice(POLIS_LAMA.indexOf('kolom={kolomRekapPolisLama(POLIS_LAMA_REKAP[t])}'))
    expect(rekap.slice(0, rekap.indexOf('/>'))).not.toContain('nomorAwal')
  })
})

describe('CLAIM_AMOUNT disembunyikan dari View Old Policy (keputusan work owner 08-10-2026)', () => {
  it('peserta dan rekap keempat tipe tanpa CLAIM_AMOUNT; kolom lain tetap dan berurutan sama', () => {
    expect(KOLOM_POLIS_LAMA_TERSEMBUNYI).toEqual(['CLAIM_AMOUNT'])
    for (const t of ['QR', 'QP', 'TP', 'TR'] as const) {
      for (const sumber of [POLIS_LAMA_PESERTA[t], POLIS_LAMA_REKAP[t]]) {
        const tampil = kolomPolisLamaTampil(sumber).map((k) => k.kolom)
        expect(tampil, t).not.toContain('CLAIM_AMOUNT')
        expect(tampil, t).toEqual(sumber.map((k) => k.kolom).filter((k) => k !== 'CLAIM_AMOUNT'))
      }
    }
    // Daftar korpus (bukti) tetap memuat CLAIM_AMOUNT di peserta QR/QP/TR dan rekap QP.
    expect(POLIS_LAMA_PESERTA.QR.some((k) => k.kolom === 'CLAIM_AMOUNT')).toBe(true)
    expect(POLIS_LAMA_REKAP.QP.some((k) => k.kolom === 'CLAIM_AMOUNT')).toBe(true)
  })

  it('kedua tabel popup memakai penyaring; kolom No tetap', () => {
    expect(POLIS_LAMA).toContain('kolom={kolomPolisLamaTampil(POLIS_LAMA_PESERTA[t])}')
    expect(POLIS_LAMA).toContain('kolom={kolomRekapPolisLama(POLIS_LAMA_REKAP[t])}')
    expect(POLIS_LAMA).toContain('nomorAwal={nomorAwalPeserta(halaman)}')
  })
})

describe('kolom Endorsement No rekap View Old Policy (keputusan work owner 10-10-2026)', () => {
  it('tepat sesudah PL NUMBER untuk keempat tipe; daftar korpus tidak berubah', () => {
    for (const t of ['QR', 'QP', 'TP', 'TR'] as const) {
      const k = kolomRekapPolisLama(POLIS_LAMA_REKAP[t]).map((x) => x.kolom)
      expect(k[k.indexOf('PL_NUMBER') + 1], t).toBe('ENDORSEMENT_NO_TAMPIL')
      expect(k.filter((x) => x !== 'ENDORSEMENT_NO_TAMPIL'), t).toEqual(kolomPolisLamaTampil(POLIS_LAMA_REKAP[t]).map((x) => x.kolom))
      expect(POLIS_LAMA_REKAP[t].some((x) => x.kolom === 'ENDORSEMENT_NO_TAMPIL'), t).toBe(false)
    }
  })

  it('PL_NUMBER_EDM berakhiran /<angka> apa adanya; kosong / sama PL NUMBER / tanpa akhiran → New Business', () => {
    expect(teksEndorsementNo({ PL_NUMBER: 'UJI-PL-1', PL_NUMBER_EDM: 'UJI-PL-1/01' })).toBe('UJI-PL-1/01')
    expect(teksEndorsementNo({ PL_NUMBER: 'UJI-PL-1', PL_NUMBER_EDM: '' })).toBe('New Business')
    expect(teksEndorsementNo({ PL_NUMBER: 'UJI-PL-1' })).toBe('New Business')
    expect(teksEndorsementNo({ PL_NUMBER: 'UJI-PL-1', PL_NUMBER_EDM: 'UJI-PL-1' })).toBe('New Business')
    expect(teksEndorsementNo({ PL_NUMBER: 'UJI-PL-1', PL_NUMBER_EDM: 'UJI-LAIN' })).toBe('New Business')
  })

  it('tabel rekap merender kolom itu sesudah PL NUMBER', () => {
    const baris = barisRekapPolisLama([
      { COB: 'UJI-COB', PL_NUMBER: 'UJI-PL-1', PL_NUMBER_EDM: 'UJI-PL-1/01', CURRENCY: 'IDR' },
      { COB: 'UJI-COB', PL_NUMBER: 'UJI-PL-1', PL_NUMBER_EDM: '', CURRENCY: 'IDR' },
    ])
    const html = renderToStaticMarkup(<TabelKorpus kolom={kolomRekapPolisLama(POLIS_LAMA_REKAP.QR)} baris={baris} kunci={(_, i) => String(i)} />)
    const kepala = sel(html, 'th')
    expect(kepala.slice(0, 4)).toEqual(['COB', 'PL NUMBER', 'Endorsement No', 'CURRENCY'])
    const isi = sel(html, 'td')
    const lebar = kolomRekapPolisLama(POLIS_LAMA_REKAP.QR).length
    expect(isi.slice(0, 4)).toEqual(['UJI-COB', 'UJI-PL-1', 'UJI-PL-1/01', 'IDR'])
    expect(isi.slice(lebar, lebar + 4)).toEqual(['UJI-COB', 'UJI-PL-1', 'New Business', 'IDR'])
    expect(POLIS_LAMA).toContain('baris={barisRekapPolisLama(isi.rekap)}')
  })
})
