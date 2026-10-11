// Tombol di panel rinci baris adjustment (Send to Committe, Acceptation, Generate DLA, View Komite No) adalah aksi baris:
// server menyaring tata baris lewat Indeks (`services/aksi.go` aksiBarisAdj). Tanpa nomor baris, "Send to Committe"
// terkirim ber-indeks 0 dan ditolak "aksi ini tidak tersedia di layar kasus saat ini" (laporan work owner 09-10-2026).

import { createElement, useContext } from 'react'
import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import type { Halaman, Tata } from '../api'
import { barisModal, DAFTAR_RINCI } from './rincian'
import TataView, { BarisAdjustment, indeksTombol, type KonteksTata } from './TataView'

const dasar: Omit<KonteksTata, 'h'> = { ubah: () => {}, aksi: () => {}, opsi: () => [], pesanMedan: {}, sibuk: false }

function Penguji() {
  return createElement('i', null, `BARIS-${useContext(BarisAdjustment)}`)
}

describe('panel rinci membawa nomor baris adjustment', () => {
  const grid: Tata = {
    jenis: 'grid',
    jalur: DAFTAR_RINCI,
    kolom: [{ jenis: 'medan', jalur: 'Type', label: 'Type', kendali: 'tampil' }],
    baris: [[{ tampil: true, hanyaBaca: true }], [{ tampil: true, hanyaBaca: true }]],
  }
  const h: Halaman = { nilai: {}, daftar: { [DAFTAR_RINCI]: [{ Type: 'UJI-1' }, { Type: 'UJI-2' }] } }
  const html = renderToStaticMarkup(
    createElement(TataView, {
      tata: [grid],
      k: { ...dasar, h, rincian: { daftar: DAFTAR_RINCI, isi: () => createElement(Penguji) } },
    }),
  )

  it('isi panel baris ke-2 membaca nomor baris 2', () => {
    expect(html).toContain('BARIS-2')
  })

  it('di luar panel rinci nomor barisnya 0', () => {
    expect(renderToStaticMarkup(createElement(Penguji))).toContain('BARIS-0')
  })
})

// Modal baris adjustment (Send Claim to Committee, Submit DLA) dibuka berkunci `komite:n` / `dla:n`.
describe('barisModal', () => {
  it('modal baris membawa nomor barisnya', () => {
    expect(barisModal('komite:2')).toBe(2)
    expect(barisModal('dla:1')).toBe(1)
  })
  it('modal halaman (PLA, tutup klaim) bernomor 0', () => {
    expect(barisModal('pla')).toBe(0)
    expect(barisModal('tutupKlaim')).toBe(0)
  })
})

describe('indeksTombol', () => {
  it('tombol tanpa indeks sendiri memakai nomor baris panel', () => {
    expect(indeksTombol(0, 2)).toBe(2)
  })
  it('indeks baris grid sendiri tetap menang', () => {
    expect(indeksTombol(3, 2)).toBe(3)
  })
  it('di luar panel tetap 0', () => {
    expect(indeksTombol(0, 0)).toBe(0)
  })
})
