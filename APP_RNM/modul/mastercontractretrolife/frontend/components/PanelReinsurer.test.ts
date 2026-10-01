// Akses Security Reinsurer (keputusan work owner 02-10-2026): SAMA dengan Treaty Contract Out
// (`PanelReinsurerKombinasi`) - panel Security tampil DI BAWAH daftar reinsurer, di dalam panel Reinsurer,
// tanpa menggantikannya; baris lain mengganti isinya; reinsurer yang dihapus menutupnya.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

const KODE = readFileSync(join(__dirname, 'PanelReinsurer.tsx'), 'utf8')

describe('akses Security Reinsurer seperti Treaty Contract Out', () => {
  it('panel Reinsurer TIDAK diganti selama Security terbuka - nol return dini bergantung security', () => {
    expect(KODE).not.toMatch(/if \(security !== null\) \{\s*return/)
  })

  it('panel Security dirender sesudah daftar reinsurer dan Total Share, di dalam section yang sama', () => {
    const tabel = KODE.indexOf('<div className="mcrl-tabel">')
    const total = KODE.indexOf('REINSURER_MCRL.totalShare')
    const panel = KODE.indexOf('<PanelSecurity')
    const akhir = KODE.lastIndexOf('</section>')
    expect(tabel).toBeGreaterThan(-1)
    expect(panel).toBeGreaterThan(total)
    expect(total).toBeGreaterThan(tabel)
    expect(akhir).toBeGreaterThan(panel)
  })

  it('baris lain mengganti isi panel Security (key = reinsurer)', () => {
    expect(KODE).toContain('setSecurity(r)')
    expect(KODE).toMatch(/<PanelSecurity\s+key=\{security\.id\}/)
  })

  it('reinsurer yang dihapus ikut menutup panel Security-nya', () => {
    expect(KODE).toContain('setSecurity((s) => (s?.id === id ? null : s))')
  })
})
