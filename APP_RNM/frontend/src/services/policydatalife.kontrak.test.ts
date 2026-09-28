import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import type { PolicyDataLife } from './api'

// Kontrak `GET /api/polis-life/ringkas` — pl4, tiket 08. SISI KLIEN.
//
// ⛔ Pasangan `internal/services/polis_ringkas_kontrak_test.go`. Di sini
// TypeScript membaca tag JSON struct Go. Daftar kunci di bawah bertipe
// `Record<keyof PolicyDataLife, true>`: medan antarmuka yang terlupa atau
// kunci yang tidak ada di antarmuka membuat `tsc` gagal — jadi daftar ini
// tidak dapat menyimpang dari antarmukanya tanpa ketahuan.

const KUNCI: Record<keyof PolicyDataLife, true> = {
  nomorPolis: true,
  type: true,
  marketingName: true,
  cedingCoName: true,
  policyHolderName: true,
  businessName: true,
  dateReceived: true,
  status: true,
  statusUpdate: true,
  productNameId: true,
  productName: true,
  prodKe: true,
  medanTanpaSumber: true,
  typeCeding: true,
  typeCedingName: true,
  proRateType: true,
  wpc: true,
  retroName: true,
  securityReinsurer: true,
  sobName: true,
}

describe('kontrak PolicyDataLife dua sisi', () => {
  it('tag JSON struct Go sama dengan medan antarmuka klien', () => {
    const go = readFileSync(
      join(__dirname, '..', '..', '..', 'internal', 'services', 'polis_ringkas.go'),
      'utf8',
    )
    const awal = go.indexOf('type PolicyDataLife struct {')
    expect(awal).toBeGreaterThanOrEqual(0)
    const blok = go.slice(awal, go.indexOf('\n}', awal))
    const server = [...blok.matchAll(/json:"([^",]+)/g)].map((m) => m[1]).sort()
    expect(server.length).toBeGreaterThanOrEqual(10)
    expect(server).toEqual(Object.keys(KUNCI).sort())
  })
})
