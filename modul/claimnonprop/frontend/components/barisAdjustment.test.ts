// Disalin dari `modul/claimprop/frontend/components/barisAdjustment.test.ts` (pola, bukan impor): keputusan work owner yang disebut di bawah
// adalah keputusan layar Claim Prop yang ditiru Claim Non Prop (prompt Claim Non Prop tahap 1).
// Tombol di panel rinci baris adjustment (Send to Committe, Acceptation, Generate DLA, View Komite No) adalah aksi baris:
// server menyaring tata baris lewat Indeks (`services/aksi.go` aksiBarisAdj). Tanpa nomor baris, "Send to Committe"
// terkirim ber-indeks 0 dan ditolak "aksi ini tidak tersedia di layar kasus saat ini" (laporan work owner 09-10-2026).

import { createElement, useContext } from 'react'
import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import type { Halaman, Tata } from '../api'
import { DAFTAR_ADJ, DAFTAR_INTEREST } from './rincian'
import TataView, { alamatAksi, BarisAdjustment, TutupPanel, type KonteksTata } from './TataView'

const dasar: Omit<KonteksTata, 'h'> = { ubah: () => {}, aksi: () => {}, opsi: () => [], pesanMedan: {}, sibuk: false }

function Penguji() {
  return createElement('i', null, `BARIS-${useContext(BarisAdjustment)}`)
}

describe('panel rinci membawa nomor baris adjustment', () => {
  const grid: Tata = {
    jenis: 'grid',
    jalur: DAFTAR_ADJ,
    kolom: [{ jenis: 'medan', jalur: 'Type', label: 'Type', kendali: 'tampil' }],
    baris: [[{ tampil: true, hanyaBaca: true }], [{ tampil: true, hanyaBaca: true }]],
  }
  const h: Halaman = { nilai: {}, daftar: { [DAFTAR_ADJ]: [{ Type: 'UJI-1' }, { Type: 'UJI-2' }] } }
  const html = renderToStaticMarkup(
    createElement(TataView, {
      tata: [grid],
      k: {
        ...dasar,
        h,
        rincian: [{ daftar: DAFTAR_ADJ, isi: () => createElement(Penguji), bukaAwal: true, nomorAkseptasi: true }],
      },
    }),
  )

  it('isi panel baris ke-2 membaca nomor baris 2', () => {
    expect(html).toContain('BARIS-2')
  })

  it('di luar panel rinci nomor barisnya 0', () => {
    expect(renderToStaticMarkup(createElement(Penguji))).toContain('BARIS-0')
  })
})

describe('alamatAksi', () => {
  it('sel grid di panel akseptasi: indeks = nomor akseptasi, baris = baris grid panel', () => {
    expect(alamatAksi(3, 2, true)).toEqual({ indeks: 2, baris: 3 })
  })
  it('medan / tombol panel di luar grid: indeks = nomor akseptasi, baris 0', () => {
    expect(alamatAksi(0, 2, false)).toEqual({ indeks: 2, baris: 0 })
  })
  it('di luar panel: indeks baris grid-nya sendiri, baris 0', () => {
    expect(alamatAksi(3, 0, true)).toEqual({ indeks: 3, baris: 0 })
    expect(alamatAksi(0, 0, false)).toEqual({ indeks: 0, baris: 0 })
  })
})

// Expand pane InputDtlInterest (bukan panel akseptasi): tanpa nomor akseptasi, tombol Cancel / Submit menutup pane.
describe('expand pane Insured Interests', () => {
  function PengujiTutup() {
    return createElement('i', null, `TUTUP-${useContext(TutupPanel) === null ? 'tidak' : 'ada'}`)
  }
  const grid: Tata = {
    jenis: 'grid',
    jalur: DAFTAR_INTEREST,
    kolom: [{ jenis: 'medan', jalur: 'ObjectName', label: 'Insured Interest', kendali: 'teks' }],
    baris: [[{ tampil: true, hanyaBaca: true }]],
  }
  const h: Halaman = { nilai: {}, daftar: { [DAFTAR_INTEREST]: [{ ObjectName: 'UJI-1' }] } }
  const tutup = renderToStaticMarkup(
    createElement(TataView, {
      tata: [grid],
      k: { ...dasar, h, rincian: [{ daftar: DAFTAR_INTEREST, isi: () => createElement(PengujiTutup) }] },
    }),
  )

  it('pane tertutup sampai baris diklik (tanpa bukaAwal)', () => {
    expect(tutup).not.toContain('TUTUP-')
    expect(tutup).toContain('aria-expanded="false"')
  })
})
