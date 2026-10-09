// Kontrak Resolve Complete / Decline = form BACA — laporan pemakai 9 Oktober
// 2026: sesudah approve, isian masih dapat diubah.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

const FORM = readFileSync(join(__dirname, 'pages', 'FormKontrakTreatyIn.tsx'), 'utf8')

describe('form Treaty In — status tuntas mengunci', () => {
  it('mode efektif = lihat bila Resolve Complete / Decline; bisaUbah dari mode efektif', () => {
    expect(FORM).toContain("statusMuat === 'Resolve Complete' || statusMuat === 'Decline'")
    // ⭐ 9 Oktober 2026 — Force Edit (IT) membuka kunci, termasuk Resolve Complete.
    expect(FORM).toContain("const mode: ModeForm = paksaUbah ? 'ubah' : statusTuntas ? 'lihat' : modeRute")
    expect(FORM.indexOf("const bisaUbah = mode === 'ubah'")).toBeGreaterThan(FORM.indexOf('const statusTuntas'))
  })
})
