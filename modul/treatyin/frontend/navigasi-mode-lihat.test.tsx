// Laporan pemakai 8 Oktober 2026: "kenapa di treaty in saat view malah tidak
// bisa klik down colom contoh yg di limits". Mode lihat mengunci isi tab
// dengan `<fieldset disabled>`, dan HTML mematikan SETIAP `<button>` di
// dalamnya — juga yang hanya membuka rincian atau memindah sub-tab. Di Pega
// keduanya tetap hidup di mode lihat.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { StripTabNavigasi, TombolNavigasi } from './components/navigasi'

const baca = (f: string) => readFileSync(join(__dirname, 'components', f), 'utf8')

describe('navigasi tetap hidup di mode lihat', () => {
  it('kontrol navigasi BUKAN kontrol form — `fieldset disabled` tidak menjangkaunya', () => {
    const html = renderToStaticMarkup(
      <fieldset disabled>
        <TombolNavigasi className="trin__buka" terbuka={false} label="Rincian 1" onKlik={() => undefined}>
          ▸
        </TombolNavigasi>
        <StripTabNavigasi tab={['Event Limits', 'Deduction']} aktif="Event Limits" onPilih={() => undefined} />
      </fieldset>,
    )
    expect(html).not.toMatch(/<button/)
    expect(html).toContain('role="button" tabindex="0" class="trin__buka" aria-label="Rincian 1" aria-expanded="false"')
    // Markup strip SAMA dengan `StripTab` inti — kelas, peran, dan pilihan.
    expect(html).toContain('<div class="tabs" role="tablist">')
    expect(html).toContain('role="tab" tabindex="0" aria-selected="true" class="tabs__item tabs__item--aktif"')
  })

  it('tombol buka rincian dan strip sub-tab di dalam tab memakai navigasi, bukan `<button>`', () => {
    // ⭐ 8 Oktober 2026 — `gridPega.tsx` (grid bentuk Pega) memegang tombol
    // ▸/▾ rincian untuk tab yang memakainya (TabShareProp, Limits, …).
    for (const f of ['limitsUI.tsx', 'TabAngsuran.tsx', 'TabShareNonProp.tsx', 'gridPega.tsx']) {
      expect(baca(f), f).not.toMatch(/<button[^>]*aria-expanded/)
      expect(baca(f), f).toContain('<TombolNavigasi')
    }
    expect(baca('TabShareProp.tsx')).toContain('<GridPega')
    expect(baca('TabShareProp.tsx')).not.toMatch(/<button[^>]*aria-expanded/)
    for (const f of ['TabLimitsProp.tsx', 'TabShare.tsx', 'TabShareNonProp.tsx', 'TabShareProp.tsx']) {
      expect(baca(f), f).toContain('<StripTabNavigasi')
      expect(baca(f), f).not.toMatch(/<StripTab /)
    }
  })
})
