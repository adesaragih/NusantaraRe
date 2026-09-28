import { describe, expect, it } from 'vitest'

import type { Peserta } from '../../services/api'
import { keIsianTanggal, tanggalAwal } from './PanelTanggalKlaim'

// Uji tiga tanggal klaim dialog Edit Date — sensus §3.1.

describe('keIsianTanggal', () => {
  it('bentuk backend berjam dipotong ke tanggalnya', () => {
    // ⛔ Inti uji ini: `<input type="date">` menolak "2026-03-02 00:00:00"
    // DIAM-DIAM — kotaknya tampak kosong walau tanggalnya tersimpan.
    expect(keIsianTanggal('2026-03-02 00:00:00')).toBe('2026-03-02')
    expect(keIsianTanggal('2026-03-02')).toBe('2026-03-02')
  })

  it('kosong dan bentuk asing menjadi kosong, bukan tebakan', () => {
    expect(keIsianTanggal('')).toBe('')
    expect(keIsianTanggal('02/03/2026')).toBe('')
  })
})

describe('tanggalAwal', () => {
  it('ketiga tanggal dibaca dari PESERTA, bukan dari klaim', () => {
    const p = {
      tanggalTerimaKlaim: '2026-03-02 00:00:00',
      tanggalDokumenLengkap: '',
      tanggalKonfirmasi: '2026-03-04 00:00:00',
    } as Peserta
    expect(tanggalAwal(p)).toEqual({
      tanggalTerimaKlaim: '2026-03-02',
      tanggalDokumenLengkap: '',
      tanggalKonfirmasi: '2026-03-04',
    })
  })
})
