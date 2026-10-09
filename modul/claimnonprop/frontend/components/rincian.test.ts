// Expand pane Claim Non Prop: hanya tiga grid ber-`pyEditingMode=expandPane` di XML yang dapat dibuka (Acceptation List,
// Insured Interests Outstanding, XOL Allocation panel akseptasi). Pola panel Claim Prop (bukan impor).

import { describe, expect, it } from 'vitest'

import {
  barisModal,
  barisTerbuka,
  bukaAwal,
  DAFTAR_ADJ,
  DAFTAR_INTEREST,
  jalurXOL,
  pecahJalurXOL,
  rincianUntuk,
  type RincianGrid,
} from './rincian'

const adj: RincianGrid = { daftar: DAFTAR_ADJ, isi: () => null, bukaAwal: true, nomorAkseptasi: true }
const interest: RincianGrid = { daftar: DAFTAR_INTEREST, isi: () => null }

describe('grid ber-expand pane', () => {
  it('Acceptation List, Insured Interests, dan XOL Allocation panel akseptasi', () => {
    expect(DAFTAR_ADJ).toBe('ClaimData.AdjustmentList')
    expect(DAFTAR_INTEREST).toBe('ClaimData.InterestList')
    expect(jalurXOL(2)).toBe('ClaimData.AdjustmentList(2).SpreadingRisk')
    expect(rincianUntuk([adj, interest], DAFTAR_INTEREST)).toBe(interest)
    expect(rincianUntuk([adj, interest], DAFTAR_ADJ)).toBe(adj)
  })

  it('grid lain tidak dapat dibuka', () => {
    expect(rincianUntuk([adj, interest], 'ClaimData.TotalInterestInsured')).toBeUndefined()
    expect(rincianUntuk([adj, interest], 'ClaimData.SpreadingRisk')).toBeUndefined()
    expect(rincianUntuk([adj], undefined)).toBeUndefined()
    expect(rincianUntuk(undefined, DAFTAR_ADJ)).toBeUndefined()
  })

  it('jalur grid XOL panel membawa nomor akseptasinya', () => {
    expect(pecahJalurXOL('ClaimData.AdjustmentList(3).SpreadingRisk')).toBe(3)
    expect(pecahJalurXOL('ClaimData.SpreadingRisk')).toBeNull()
    expect(pecahJalurXOL('ClaimData.AdjustmentList(3).AlokasiXOLPaid')).toBeNull()
  })
})

describe('baris ber-panel terbuka', () => {
  it('Acceptation List: baris terbaru terbuka sejak awal; pane lain tertutup sampai diklik', () => {
    expect(bukaAwal(3)).toBe(3)
    expect(bukaAwal(0)).toBeNull()
    expect(bukaAwal(3, false)).toBeNull()
  })

  it('Add membuka baris baru; hapus menutup baris yang hilang; lainnya tetap', () => {
    expect(barisTerbuka(1, 1, 2)).toBe(2)
    expect(barisTerbuka(null, 0, 1, false)).toBe(1)
    expect(barisTerbuka(2, 2, 1)).toBe(1)
    expect(barisTerbuka(2, 2, 1, false)).toBeNull()
    expect(barisTerbuka(1, 1, 0)).toBeNull()
    expect(barisTerbuka(1, 3, 2)).toBe(1)
    expect(barisTerbuka(null, 2, 2)).toBeNull()
  })
})

describe('barisModal', () => {
  it('modal komite membawa nomor akseptasinya', () => {
    expect(barisModal('komite:2')).toBe(2)
  })
  it('modal halaman bernomor 0', () => {
    expect(barisModal('pla')).toBe(0)
    expect(barisModal('tutupKlaim')).toBe(0)
    expect(barisModal('cwp')).toBe(0)
  })
})
