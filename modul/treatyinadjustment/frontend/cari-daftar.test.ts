// Pencarian daftar Adjustment — permintaan pemakai 8 Oktober 2026: *"di
// adjustment buat fitur pencarian"*.

import { describe, expect, it } from 'vitest'

import type { BarisPenyesuaian } from './api'
import { cocokCari, kataCari, potongSorot } from './pages/PenyesuaianKontrak'

const B: BarisPenyesuaian = {
  id: '1001700',
  idAsal: '1001200',
  jenisPenyesuaian: 'Revision',
  jenisMaterial: '1',
  namaKontrak: 'Motor Risk and Catastrophe Excess of Loss Treaty 2024',
  sifatProporsi: 'NonProportional',
  asalBisnis: 'MARSH REINSURANCE BROKERS INDONESIA',
  cedant: 'MNC ASURANSI INDONESIA',
  tanggalMulai: '',
  tanggalBerakhir: '',
  posisi: 'AGUNGPUTRAANDALAS',
  statusAkseptasi: 'Accept',
}

describe('cocokCari', () => {
  it('kosong / spasi saja = semua baris', () => {
    expect(cocokCari(B, '')).toBe(true)
    expect(cocokCari(B, '   ')).toBe(true)
  })

  it('tanpa membedakan huruf besar/kecil, di kolom mana pun', () => {
    expect(cocokCari(B, 'marsh')).toBe(true)
    expect(cocokCari(B, '1001200')).toBe(true)
    expect(cocokCari(B, 'agungputra')).toBe(true)
  })

  it('beberapa kata = SEMUANYA harus ada, boleh di kolom berbeda', () => {
    expect(cocokCari(B, 'marsh 2024')).toBe(true)
    expect(cocokCari(B, 'marsh 2025')).toBe(false)
  })

  it('kataCari memecah di spasi dan mengecilkan huruf', () => {
    expect(kataCari('  Marsh   2024 ')).toEqual(['marsh', '2024'])
  })
})

describe('potongSorot', () => {
  it('memisahkan bagian cocok (huruf asli dipertahankan)', () => {
    expect(potongSorot('MARSH REINSURANCE', ['rein'])).toEqual([
      { isi: 'MARSH ', cocok: false },
      { isi: 'REIN', cocok: true },
      { isi: 'SURANCE', cocok: false },
    ])
  })

  it('beberapa kata dan beberapa kemunculan', () => {
    expect(potongSorot('ab ab', ['ab']).filter((p) => p.cocok).map((p) => p.isi)).toEqual(['ab', 'ab'])
  })

  it('tanpa kata = satu bagian polos', () => {
    expect(potongSorot('TESTS', [])).toEqual([{ isi: 'TESTS', cocok: false }])
  })
})
