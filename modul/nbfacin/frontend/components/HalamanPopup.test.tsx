// Pager popup ber-onKirim (work owner 08-10-2026: Change SOB tidak bisa ke halaman berikutnya).

import { readFileSync } from 'node:fs'
import { join } from 'node:path'
import { describe, expect, it } from 'vitest'

import { cegahSubmit } from './HalamanPopup'

const baca = (...p: string[]) => readFileSync(join(__dirname, ...p), 'utf8').replace(/\r\n/g, '\n')

describe('HalamanPopup', () => {
  it('klik tombol pager membatalkan submit form; klik di luar tombol tidak', () => {
    let batal = 0
    const tombol = { closest: (s: string) => (s === 'button' ? {} : null) }
    const teks = { closest: () => null }
    // Lingkungan uji tanpa DOM: Element dipalsukan agar `instanceof Element` berlaku pada objek uji.
    const asli = (globalThis as { Element?: unknown }).Element
    class Palsu {}
    ;(globalThis as { Element?: unknown }).Element = Palsu
    try {
      cegahSubmit({ target: Object.assign(new Palsu(), tombol), preventDefault: () => batal++ } as never)
      cegahSubmit({ target: Object.assign(new Palsu(), teks), preventDefault: () => batal++ } as never)
    } finally {
      ;(globalThis as { Element?: unknown }).Element = asli
    }
    expect(batal).toBe(1)
  })
  it('popup Modal ber-onKirim memakai HalamanPopup, bukan Halaman inti langsung', () => {
    for (const f of [baca('PopupPilihAgent.tsx'), baca('PopupRiskAddress.tsx'), baca('..', 'pages', 'FormOpportunity.tsx')]) {
      expect(f).toContain('onKirim=')
      expect(f).toContain('<HalamanPopup halaman=')
      expect(f).not.toMatch(/<Halaman\s/)
    }
  })
})
