import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { KEPUTUSAN_KOMITE } from '../../assets/labels.komite'
import { PERAN } from '../../assets/labels.claimlife'
import type { KasusKomite } from '../../services/api'
import { bolehEskalasi, kalimatHasilKeputusan, PILIHAN_KEPUTUSAN } from './KasusKomite'

// Uji layar keputusan Komite — tiket 02.

const SUMBER = readFileSync(join(__dirname, 'KasusKomite.tsx'), 'utf8')

describe('keputusan komite', () => {
  it('enum tertutup {1 Setuju, 2 Tolak}, tanpa pilihan ketiga', () => {
    expect(PILIHAN_KEPUTUSAN.map((p) => p.kode)).toEqual(['1', '2'])
  })

  it('label VERBATIM dari ShowTransfer', () => {
    expect(KEPUTUSAN_KOMITE.konfirmasi).toBe('Are you sure to accept this document?')
    expect(KEPUTUSAN_KOMITE.submit).toBe('Submit')
    expect(KEPUTUSAN_KOMITE.cancel).toBe('Cancel')
  })

  it('kalimat hasil menyebut naik atau berhenti', () => {
    expect(
      kalimatHasilKeputusan({ kataKeputusan: 'Setuju', tingkatDiputus: 1, berlanjut: true, tingkatBerikut: 2 }),
    ).toBe('Setuju tercatat di tingkat 1. Kasus naik ke tingkat 2.')
    expect(
      kalimatHasilKeputusan({ kataKeputusan: 'Tolak', tingkatDiputus: 2, berlanjut: false, tingkatBerikut: 0 }),
    ).toBe('Tolak tercatat di tingkat 2. Tangga berhenti.')
    expect(
      kalimatHasilKeputusan({
        kataKeputusan: 'Setuju', tingkatDiputus: 3, berlanjut: false, tingkatBerikut: 0,
        nomorAkseptasi: 'RNML-AL01.09.26.00007',
      }),
    ).toBe('Setuju tercatat di tingkat 3. Nomor akseptasi RNML-AL01.09.26.00007.')
  })

  it('eskalasi hanya untuk admin dan hanya bila ada tingkat di atas (tiket 03)', () => {
    const k = (tingkat: number): KasusKomite => ({
      kasus: {
        kasusId: 'KMTLF-UJI', tglUpdate: '', statusWork: '', klaimId: '', nomorKlaim: '',
        tingkatBerjalan: tingkat, komiteLoop: 3, nilaiKlaim: '', mataUang: '', statusBaris: '',
      },
      adjustmentId: '',
      tangga: [],
      giliranSaya: false,
    })
    expect(bolehEskalasi([PERAN.admin], k(1))).toBe(true)
    expect(bolehEskalasi([PERAN.admin], k(3))).toBe(false)
    expect(bolehEskalasi([PERAN.spv], k(1))).toBe(false)
    expect(bolehEskalasi([PERAN.admin], null)).toBe(false)
  })

  it('formulir hanya pada giliran pelaku; dropdown wajib', () => {
    expect(SUMBER).toContain('k.giliranSaya && (')
    expect(SUMBER).toMatch(/<select\s+required/)
  })
})
