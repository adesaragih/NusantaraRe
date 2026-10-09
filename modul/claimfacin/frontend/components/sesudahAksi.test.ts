// Keputusan layar sesudah aksi berhasil: modal yang dibuka / ditutup, pop-up layar sesudah aksi server tanpa akibat
// (BukaSebab, BukaKatastrofe, BukaOutstanding, BukaRetro), grid hasil Choose Polis, dan judul modal.

import { describe, expect, it } from 'vitest'

import type { Tata } from '../api'
import { CFI } from '../labels'
import { judulModal, lanjutanAksi, MODAL_POLIS, modalSesudah } from './sesudahAksi'

const ada: Tata[] = []

describe('lanjutanAksi', () => {
  it('aksi no-op pembuka harness membuka pop-up layar; Outstanding membawa baris objeknya', () => {
    expect(lanjutanAksi('BukaSebab', 0)).toEqual({ popup: 'sebab', indeks: 0 })
    expect(lanjutanAksi('BukaKatastrofe', 0)).toEqual({ popup: 'katastrofe', indeks: 0 })
    expect(lanjutanAksi('BukaOutstanding', 2)).toEqual({ popup: 'outstanding', indeks: 2 })
    expect(lanjutanAksi('BukaRetro', 1)).toEqual({ popup: 'retro', indeks: 1 })
    expect(lanjutanAksi('BukaRetroList', 0)).toEqual({ popup: 'retro', indeks: 0 })
  })

  it('Save katastrofe baru: daftar katastrofe terbuka lagi', () => {
    expect(lanjutanAksi('SaveCatasrtope', 0)).toEqual({ popup: 'katastrofe', indeks: 0 })
  })

  it('Search pop-up Choose Polis: grid hasil dibaca; aksi lain menutup pop-up', () => {
    expect(lanjutanAksi('SearchPolis', 0)).toEqual({ cariPolis: true })
    expect(lanjutanAksi('GetNameCauseofLoss', 0)).toBeNull()
    expect(lanjutanAksi('SetCatastrope', 0)).toBeNull()
  })
})

describe('modalSesudah', () => {
  it('bukaModal server membuka modalnya (Print PLA objek 1, Choose Polis sesudah CopyNB)', () => {
    expect(modalSesudah(null, { bukaModal: 'pla:1' }, false)).toBe('pla:1')
    expect(modalSesudah(MODAL_POLIS, { bukaModal: MODAL_POLIS, modal: { [MODAL_POLIS]: ada } }, true)).toBe(MODAL_POLIS)
  })

  it('tombol kaki modal berhasil tanpa pesan menutup modal; dengan pesan tetap terbuka', () => {
    expect(modalSesudah('pla:1', { modal: { 'pla:1': ada } }, true)).toBeNull()
    expect(modalSesudah('pla:1', { modal: { 'pla:1': ada }, pesan: ['UJI-PESAN'] }, true)).toBe('pla:1')
  })

  it('aksi di isi modal (Search) membiarkannya terbuka; modal yang tidak dikirim lagi ditutup', () => {
    expect(modalSesudah(MODAL_POLIS, { modal: { [MODAL_POLIS]: ada } }, false)).toBe(MODAL_POLIS)
    expect(modalSesudah('protectDOL', { modal: { [MODAL_POLIS]: ada } }, false)).toBeNull()
    expect(modalSesudah(null, { modal: { [MODAL_POLIS]: ada } }, false)).toBeNull()
  })
})

describe('judulModal', () => {
  it('kunci modal ke judul pop-up Pega', () => {
    expect(judulModal(MODAL_POLIS)).toBe(CFI.popPolis)
    expect(judulModal('protectDOL')).toBe(CFI.popProtectDOL)
    expect(judulModal('pla:2')).toBe(CFI.popPLA)
    expect(judulModal('dla:1')).toBe(CFI.popDLA)
    expect(judulModal('cedant:ClaimData.ObjectList(1).ObjectItemList(1).Adjustment(1)')).toBe(CFI.popCedant)
    expect(judulModal('komite:ClaimData.ObjectList(1).ObjectItemList(1).Adjustment(1)')).toBe(CFI.popKomite)
    expect(judulModal('tutup')).toBe(CFI.popTutup)
    expect(judulModal('tolak')).toBe(CFI.popTolak)
  })
})
