// Uji Inbox Claim Life — F0.4.
//
// Yang diuji: tab mana yang tampil untuk peran mana, label kolom yang
// VERBATIM korpus, dan bahwa nomor tahap sejalan dengan `models.Tahap` di Go.

import { existsSync, readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { PERAN, TAHAP } from '../assets/labels'
import { TAHAP_NOMOR } from '../services/api'
import { tabUntuk } from './InboxClaimLife'

const KORPUS = 'D:\\XML\\RNM_BRD\\Claim Life'
const adaKorpus = existsSync(KORPUS)
const SUMBER = readFileSync(join(__dirname, 'InboxClaimLife.tsx'), 'utf8')

describe('tab per peran', () => {
  it('Admin melihat KEDUA antrian Admin', () => {
    // `Register_Flow.xml`: Input Register 1508 `Current operator`,
    // Outstanding 1351 `ToWorklist` — keduanya ReasLifeAdmin.
    const t = tabUntuk([PERAN.admin])
    expect(t.map((x) => x.judul)).toEqual([TAHAP.inputRegister, TAHAP.outstandingClaim])
  })

  it('Medical Advisor hanya Medical Check', () => {
    expect(tabUntuk([PERAN.medis]).map((x) => x.judul)).toEqual([TAHAP.medicalCheck])
  })

  it('SPV hanya Claim Analis', () => {
    expect(tabUntuk([PERAN.spv]).map((x) => x.judul)).toEqual([TAHAP.claimAnalis])
  })

  it('peran GANDA melihat gabungan, berurut tangga', () => {
    // ⛔ Bukan salah satu. Layar yang memaksa memilih satu peran
    // menyembunyikan separuh pekerjaan orang itu.
    const t = tabUntuk([PERAN.spv, PERAN.admin])
    expect(t.map((x) => x.judul)).toEqual([
      TAHAP.inputRegister,
      TAHAP.outstandingClaim,
      TAHAP.claimAnalis,
    ])
  })

  it('tanpa peran: nol tab, bukan tab kosong', () => {
    // Tab kosong yang tidak pernah dapat berisi hanyalah tempat orang
    // menunggu sia-sia.
    expect(tabUntuk([])).toHaveLength(0)
  })
})

describe('nomor tahap sejalan dengan models.Tahap di Go', () => {
  it('keempatnya 1..4 berurut tangga', () => {
    // ⛔ Angka di kabel. Bila urutan `models.Tahap` di Go berubah, tab akan
    // meminta antrian yang salah - dan layarnya tetap tampil normal.
    expect(TAHAP_NOMOR).toEqual({
      inputRegister: 1,
      outstanding: 2,
      medicalCheck: 3,
      claimAnalis: 4,
    })
  })

  it.runIf(existsSync('D:\\XML\\RNM_BRD\\OUTPUT_HASIL_RNM\\APP_RNM\\internal\\models\\tahap.go'))(
    'urutannya sama dengan konstanta Go',
    () => {
      const go = readFileSync(
        'D:\\XML\\RNM_BRD\\OUTPUT_HASIL_RNM\\APP_RNM\\internal\\models\\tahap.go',
        'utf8',
      )
      // `TahapTidakDikenal = iota` (0), lalu keempatnya berurut.
      const urut = ['TahapInputRegister', 'TahapOutstanding', 'TahapMedicalCheck', 'TahapClaimAnalis']
      const letak = urut.map((n) => go.indexOf(`\t${n}\n`))
      expect(letak.every((i) => i > 0)).toBe(true)
      for (let i = 1; i < letak.length; i++) {
        expect(letak[i]!).toBeGreaterThan(letak[i - 1]!)
      }
    },
  )
})

describe('label kolom VERBATIM korpus', () => {
  const kolom: Array<[string, number, string]> = [
    ['Case ID', 721, 'caseId'],
    ['Create Date/Time', 736, 'tglCreate'],
    ['Create Operator Name', 751, 'createOpName'],
    ['Work Status', 765, 'status'],
  ]

  it.runIf(adaKorpus).each(kolom)(
    'label %s ada di InboxPremiumList.xml baris %i',
    (label, baris) => {
      const rd = readFileSync(join(KORPUS, 'ReportDefinition', 'InboxPremiumList.xml'), 'utf8')
      expect(rd.split('\n')[baris - 1] ?? '').toContain(`<pyFieldLabel>${label}</pyFieldLabel>`)
    },
  )

  it.each(kolom)('layar memakai label %s apa adanya', (label) => {
    expect(SUMBER).toContain(`'${label}'`)
  })

  it.runIf(adaKorpus)('Claim No VERBATIM InputOSClaimLife.xml:951', () => {
    const sec = readFileSync(join(KORPUS, 'Section', 'InputOSClaimLife.xml'), 'utf8')
    expect(sec.split('\n')[950] ?? '').toContain('Claim No')
    expect(SUMBER).toContain("'Claim No'")
  })
})

describe('halaman dan urutan dari XML', () => {
  it.runIf(adaKorpus)('50 per halaman, maksimum 500, urut DESC', () => {
    const rd = readFileSync(join(KORPUS, 'ReportDefinition', 'InboxPremiumList.xml'), 'utf8')
    const b = rd.split('\n')
    expect(b[591] ?? '').toContain('<pyPageSize>50</pyPageSize>')
    expect(b[941] ?? '').toContain('<pyMaxRecords>500</pyMaxRecords>')
    expect(b[732] ?? '').toContain('<pySortType>DESC</pySortType>')
  })
})
