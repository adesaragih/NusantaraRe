// Tab Scoring Risk tampilan polis NB (FIRE) - NILAI TERSIMPAN (keputusan work owner 11-10-2026): pilihan, skor, total,
// hasil dibaca dari halaman polis (tabel datar nbfacin 205 / bentuk halaman Pega), tidak dihitung ulang.

import { describe, expect, it } from 'vitest'

import type { Halaman } from './api'
import { FAKTOR_SCORING_NB, HALAMAN_FAKTOR_NB, nadaKelasNB, skorPolisNB, subtotalFaktorNB } from './polisNBScoring'

const SR = 'OfferFacIn.ScoringRisk'
const D = SR + '.DataScoringRiskList'

describe('skorPolisNB', () => {
  it('bentuk tabel datar: faktor (FactorGroup / FactorName) + opsi CheckBox "true" -> pilihan, ScorePerFactor -> skor', () => {
    const h = {
      nilai: {
        [SR + '.FinalScore']: '80.5',
        [SR + '.NoteFinalScore']: 'Above Average',
        [SR + '.Status']: 'New Business',
        [SR + '.ScoringResult.AcceptableWithApproval']: 'true',
        [SR + '.ScoringResult.OperationalDirector']: 'true',
        [SR + '.ScoringResult.LossRatio']: 'false',
        'OfferFacIn.QuotationData.InsuredName': 'UJI TERTANGGUNG',
        'OfferFacIn.PolicyData.StartDateTime': '20260101T000000.000 GMT',
      },
      daftar: {
        [D]: [{ OccupationNote: 'UJI OKUPASI', TopRiskLocation: 'UJI LOKASI', Score: '80.5' }, { IndexLocation: '2' }],
        [D + '(1).ScoringFactor']: [
          { FactorName: 'Occupation', ScorePerFactor: '15' },
          { FactorGroup: 'LossRatio', FactorName: 'LossRatio', PercentLossRatio: '12.5', ScorePerFactor: '' },
          { FactorGroup: 'RiskLocation', FactorName: 'Riot', ScorePerFactor: '0.5' },
        ],
        [D + '(1).ScoringFactor(1).ScoringOption']: [
          { CheckBox: 'false', OptionNo: '1' },
          { CheckBox: 'true', OptionNo: '2' },
        ],
        // tercentang ganda -> nomor terkecil (nbfacin BacaScoring); ScorePerFactor kosong -> Score opsi
        [D + '(1).ScoringFactor(2).ScoringOption']: [
          { CheckBox: 'true', OptionNo: '6', Score: '9' },
          { CheckBox: 'true', OptionNo: '5', Score: '11' },
        ],
      },
    } as unknown as Halaman
    const v = skorPolisNB(h)
    expect(v.finalScore).toBe('80.5')
    expect(v.noteFinalScore).toBe('Above Average')
    expect(v.status).toBe('New Business')
    expect(v.insuredName).toBe('UJI TERTANGGUNG')
    expect(v.hasil.acceptableWithApproval).toBe(true)
    expect(v.hasil.opDirector).toBe(true)
    expect(v.hasil.dueLossRatio).toBe(false)
    expect(v.kolom).toHaveLength(2)
    const k = v.kolom[0]!
    expect(k.occupationNote).toBe('UJI OKUPASI')
    expect(k.total).toBe('80.5')
    expect(k.pilihan['occupation']).toBe(2)
    expect(k.skor['occupation']).toBe('15')
    expect(k.pilihan['lossRatio']).toBe(5)
    expect(k.skor['lossRatio']).toBe('11')
    expect(k.persenLossRatio).toBe('12.5')
    expect(k.pilihan['riot'] ?? 0).toBe(0)
    expect(v.kolom[1]!.indexLocation).toBe('2')
    expect(v.kolom[1]!.occupationNote).toBe('')
  })

  it('bentuk halaman Pega: `.<Grup>.<Nama>.ChechBoxN` / `.ScorePerFactor` / `.LossRatio.LossRatio.PercentLossRatio`', () => {
    const h = {
      nilai: {},
      daftar: {
        [D]: [
          {
            'Occupation.ChechBox3': 'true',
            'Occupation.ScorePerFactor': '7.5',
            'ObjectConditions.ClassConstruction.ChechBox1': 'false',
            'ObjectConditions.ClassConstruction.ChechBox4': 'true',
            'ObjectConditions.ClassConstruction.Score4': '1',
            'LossRatio.LossRatio.PercentLossRatio': '30',
          },
        ],
      },
    } as unknown as Halaman
    const k = skorPolisNB(h).kolom[0]!
    expect(k.pilihan['occupation']).toBe(3)
    expect(k.skor['occupation']).toBe('7.5')
    expect(k.pilihan['classConstruction']).toBe(4)
    expect(k.skor['classConstruction']).toBe('1')
    expect(k.persenLossRatio).toBe('30')
  })

  it('tanpa Scoring Risk: dua kolom kosong, tanpa skor', () => {
    const v = skorPolisNB({ nilai: {}, daftar: {} } as unknown as Halaman)
    expect(v.kolom).toHaveLength(2)
    expect(v.finalScore).toBe('')
    expect(Object.values(v.hasil).some(Boolean)).toBe(false)
  })
})

describe('faktor', () => {
  it('24 sub-faktor label nbfacin = 24 halaman faktor `DaftarFaktorScoring`', () => {
    const kunci = FAKTOR_SCORING_NB.flatMap((f) => f.sub.map((s) => s.kunci))
    expect(kunci).toHaveLength(24)
    expect([...kunci].sort()).toEqual(Object.keys(HALAMAN_FAKTOR_NB).sort())
  })

  it('subtotal faktor: kosong bila belum ada sub dicentang; jumlah eksak skor sub', () => {
    const f = FAKTOR_SCORING_NB.find((x) => x.no === '2')!
    const kolom = { pilihan: { riot: 2, flood: 0 }, skor: { riot: '0.5', flood: '', earthquake: '0.25' } }
    expect(subtotalFaktorNB(kolom, f)).toBe('0.75')
    expect(subtotalFaktorNB({ pilihan: {}, skor: { riot: '1' } }, f)).toBe('')
    expect(subtotalFaktorNB({ pilihan: { riot: 1 }, skor: { riot: 'x' } }, f)).toBe('')
  })

  it('nada kelas skor', () => {
    expect(nadaKelasNB('Good')).toBe('baik')
    expect(nadaKelasNB('Average')).toBe('sedang')
    expect(nadaKelasNB('Poor')).toBe('buruk')
    expect(nadaKelasNB('')).toBe('')
  })
})
