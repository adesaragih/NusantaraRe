import { describe, expect, it } from 'vitest'

import { LABEL_CARI_DIAGNOSA, ringkasanHasil } from './CariDiagnosa'

// Uji ringkasan hasil pencarian diagnosa — kelompok Medis.

describe('ringkasanHasil', () => {
  it('daftar yang MENYENTUH batas dinyatakan terpotong', () => {
    // ⛔ Inti berkas ini. `DISEASE_LIFE` 97.586 baris, dan hasil yang
    // terpotong diam-diam terbaca sebagai hasil yang LENGKAP — pemakai akan
    // menyimpulkan diagnosanya tidak ada, lalu berhenti mencari.
    const s = ringkasanHasil(50, 50)
    expect(s).toContain('TERPOTONG')
    expect(s).toContain('50')
  })

  it('di bawah batas tidak menakut-nakuti', () => {
    const s = ringkasanHasil(7, 50)
    expect(s).not.toContain('TERPOTONG')
    expect(s).toContain('7')
  })

  it('kosong berkata kosong, bukan diam', () => {
    // Layar yang diam sesudah pencarian tidak dapat dibedakan dari layar
    // yang permintaannya tidak pernah berangkat.
    expect(ringkasanHasil(0, 50)).toContain('Tidak ada')
  })

  it('label VERBATIM dari korpus', () => {
    expect(LABEL_CARI_DIAGNOSA.buka).toBe('Find Disease')
    expect(LABEL_CARI_DIAGNOSA.pilih).toBe('Choose')
  })
})
