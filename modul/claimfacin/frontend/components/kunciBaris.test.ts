// Disalin dari `modul/claimnonprop/frontend/components/kunciBaris.test.ts` (pola, bukan impor; asal Claim Prop):
// keputusan work owner yang disebut di bawah adalah keputusan layar Claim Prop yang ditiru Claim Fac In (prompt Claim Fac
// In tahap 1). Key baris grid hasil Choose Polis unik walau nomor polisnya sama - key kembar membuat baris lama
// tertinggal saat daftar berganti.

import { describe, expect, it } from 'vitest'

import type { BarisPolisCari } from '../api'
import { kunciPolis } from './kunciBaris'

const dasar: BarisPolisCari = {
  policyNo: 'UJI-POL-1',
  customerName: 'UJI TERTANGGUNG',
  sourceOfBusinessName: 'UJI-SOB',
  cedingCoName: 'UJI-CEDING',
  qq: 'UJI-QQ',
  startDateTime: '2026-01-01',
  endDateTime: '2026-12-31',
  prodke: '0',
  businessName: 'UJI-BISNIS',
}

describe('kunci baris grid hasil Choose Polis', () => {
  it('baris yang hanya berbeda Prodke / tanggal tetap ber-key berbeda', () => {
    const k = new Set([
      kunciPolis(dasar),
      kunciPolis({ ...dasar, prodke: '1' }),
      kunciPolis({ ...dasar, endDateTime: '2027-12-31' }),
    ])
    expect(k.size).toBe(3)
  })

  it('pemisah tidak membuat dua baris berbeda bertabrakan', () => {
    expect(kunciPolis({ ...dasar, sourceOfBusinessName: 'A', cedingCoName: 'BC' })).not.toBe(
      kunciPolis({ ...dasar, sourceOfBusinessName: 'AB', cedingCoName: 'C' }),
    )
  })
})
