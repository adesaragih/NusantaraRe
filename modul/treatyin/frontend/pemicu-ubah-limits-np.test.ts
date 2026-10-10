// Peristiwa `change` Pega di tab-tab Treaty In — 8 Oktober 2026.
//
// Activity/DT Pega berjalan pada `change`: hanya bila nilainya BERUBAH.
// Bentuk sebelumnya menjalankannya pada SETIAP `blur` — melewati medan
// `Reinstatement` membangun ulang grid Reinstatement (persen yang disunting
// hilang), melewati `Adjustment Rate`/`MDP %` menimpa Premium Earned/MDP yang
// diketik tangan, melewati `Deduction %` menghitung ulang `Deduction`.
// RUMUSNYA tidak berubah — hanya kapan ia dijalankan.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

const baca = (f: string) => readFileSync(join(__dirname, 'components', f), 'utf8')
const PEMICU = baca('pemicuUbah.tsx')

describe('Treaty In — Activity hanya pada perubahan nilai', () => {
  it('pemicu bersama: patokan saat masuk, aksi hanya bila nilainya berbeda', () => {
    expect(PEMICU).toContain('export function usePemicuUbah(')
    expect(PEMICU).toContain('if (nilai === terakhir.current) return')
    expect(PEMICU).toContain('export function PemicuUbah(')
  })

  it('setiap tab berumus memakai pemicu bersama', () => {
    for (const f of ['TabLimitsNonProp.tsx', 'TabLimitsProp.tsx', 'TabShareNonProp.tsx', 'TabShareProp.tsx', 'TabEgnpi.tsx', 'TabAngsuran.tsx', 'TabAkumulasi.tsx']) {
      expect(baca(f), f).toMatch(/from '\.\/pemicuUbah'/)
    }
  })

  it('nol `onBlur` yang langsung menjalankan rumus di tab-tab itu', () => {
    for (const f of ['TabLimitsNonProp.tsx', 'TabLimitsProp.tsx', 'TabShareNonProp.tsx', 'TabShareProp.tsx', 'TabEgnpi.tsx', 'TabAkumulasi.tsx']) {
      expect(baca(f), f).not.toMatch(/onBlur=\{(\(\) =>|bisaUbah \? onLepas|hitungCadanganPremi)/)
    }
    // Installment: kotak jumlah memakai `pemicuNo.keluar`; Enter tetap memicu langsung.
    const angsuran = baca('TabAngsuran.tsx')
    expect(angsuran).toContain('onBlur={pemicuNo.keluar}')
    expect(angsuran).not.toMatch(/onBlur=\{\(e\) =>/)
  })

  it('pemetaan sel grid Reinstatement → DT tetap persis ekspor', () => {
    expect(baca('TabLimitsNonProp.tsx')).toMatch(/ReinstatementPct: 'reinst-tambahan',\s*AdditionalPct: 'reinst-jumlah',\s*ReinstatementAmount1: 'reinst-persen',/)
  })

  it('MDP berubah → grid Reinstatement ikut diterapkan dari hasil rumus', () => {
    expect(baca('TabLimitsNonProp.tsx')).toContain("mdp: ['MDPList', 'MDPMinList', 'ROLPct', 'Reinstatement_List'],")
  })
})
