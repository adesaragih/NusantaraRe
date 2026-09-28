import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { KEPUTUSAN_POLIS, PENGGOLONG_POLIS } from '../../assets/labels.premiumlist'
import { bolehRejectDiTahap, TAHAP_POLIS } from '../../services/api'
import { ringkasanAkibat } from './InputOffer'

// Uji layar keputusan penawaran — tiket 01 bagian 2.

const SUMBER = readFileSync(join(__dirname, 'InputOffer.tsx'), 'utf8')

describe('Reject hanya di tahap yang punya konektornya', () => {
  it('tahap penawaran TIDAK menawarkan Reject', () => {
    // ⛔ `Reject` muncul TEPAT SEKALI di seluruh flow — Transition9 b2306
    // pada Decision2, sesudah Input Premium Detail. Tombol yang pasti
    // dijawab 409 adalah tombol yang mengajari orang mengabaikan galat.
    expect(bolehRejectDiTahap(TAHAP_POLIS.penawaran)).toBe(false)
  })

  it('tahap detail menawarkannya', () => {
    expect(bolehRejectDiTahap(TAHAP_POLIS.detail)).toBe(true)
  })

  it('tahap summary tidak — ia nol konektor keputusan', () => {
    expect(bolehRejectDiTahap(TAHAP_POLIS.summary)).toBe(false)
  })

  it('layar memagarinya, bukan hanya backend', () => {
    expect(SUMBER).toContain('bolehRejectDiTahap(tahap)')
  })
})

describe('ringkasanAkibat', () => {
  it('menunggu penggolong TIDAK berkata tertutup maupun berpindah', () => {
    // ⛔ `Confirm` di tahap penawaran hanya menyerahkan kendali ke
    // Decision3. Layar yang berkata "tersimpan" menyembunyikan langkah yang
    // di sistem lama ditanyakan.
    const s = ringkasanAkibat({
      tahapTujuan: '',
      statusWork: '',
      menungguPenggolong: true,
    })
    expect(s).toContain('kelanjutan')
    expect(s).not.toContain('ditutup')
  })

  it('tertutup menyebut status kerjanya', () => {
    const s = ringkasanAkibat({
      tahapTujuan: '',
      statusWork: 'Resolved-Rejected',
      menungguPenggolong: false,
    })
    expect(s).toContain('Resolved-Rejected')
  })

  it('berpindah menyebut tahap tujuannya', () => {
    const s = ringkasanAkibat({
      tahapTujuan: TAHAP_POLIS.penawaran,
      statusWork: '',
      menungguPenggolong: false,
    })
    expect(s).toContain(TAHAP_POLIS.penawaran)
  })
})

describe('label VERBATIM', () => {
  it('ketiga keputusan dan kedua penggolong', () => {
    expect(KEPUTUSAN_POLIS.confirm).toBe('Confirm')
    expect(KEPUTUSAN_POLIS.reject).toBe('Reject')
    expect(KEPUTUSAN_POLIS.decline).toBe('Decline')
    expect(PENGGOLONG_POLIS.offer).toBe('Offer')
    expect(PENGGOLONG_POLIS.premium).toBe('Premium')
  })
})
