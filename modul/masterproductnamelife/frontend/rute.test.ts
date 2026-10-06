import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

// Menu Product Name Life ditekan → layar kembali ke daftar awal (04-10-2026).

describe('rute Master Product Name Life', () => {
  it('halaman dipasang ulang setiap kali menunya dipilih (key = ketukMenu)', () => {
    const rute = readFileSync(join(__dirname, 'rute.tsx'), 'utf8')
    expect(rute).toContain('<MasterProductNameLife key={ketukMenu ?? 0} />')
    const app = readFileSync(join(__dirname, '..', '..', '..', 'frontend', 'App.tsx'), 'utf8')
    expect(app).toContain('ketukMenu={ketukMenu}')
  })
})
