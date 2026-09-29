import { describe, expect, it } from 'vitest'

import { BATAS_REMARKS_BYTE, remarksSah } from './DialogTolakOutstanding'

// Dialog Reject Outstanding — OQ-M5 (GILIRAN-17). `Remarks` WAJIB
// (`RejectOSClaimLife_Sec` b1653, b1695 `always`) dan disimpan di
// `T_CLAIMLF_JEJAK.KOMENTAR` VARCHAR2(4000) BYTE (migrasi 021).

describe('remarksSah', () => {
  it('kosong dan spasi saja ditolak — Remarks wajib', () => {
    expect(remarksSah('')).toBe(false)
    expect(remarksSah('   ')).toBe(false)
    expect(remarksSah('UJI alasan')).toBe(true)
  })

  it('batasnya BYTE, bukan karakter — huruf beraksen dihitung dua', () => {
    expect(remarksSah('x'.repeat(BATAS_REMARKS_BYTE))).toBe(true)
    expect(remarksSah('x'.repeat(BATAS_REMARKS_BYTE + 1))).toBe(false)
    expect(remarksSah('é'.repeat(BATAS_REMARKS_BYTE / 2 + 1))).toBe(false)
  })
})
