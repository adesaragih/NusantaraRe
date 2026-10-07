import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

// Popup ShowPolicyNo sesudah Submit Dept Head (Accept): nomor polis SUDAH terbit, jadi tidak ada jalan mundur -
// tanpa Cancel, tanpa X, Escape dan klik luar tidak menutup; OK = kirim (perintah work owner 06-10-2026: "tidak mau
// ada cancel, tidak ada close, jadi selalu maju tidak bisa mundur kalo udah klik aksep").

describe('popup nomor polis', () => {
  const layar = readFileSync(join(__dirname, 'pages', 'LayarKasus.tsx'), 'utf8')
  const mulai = layar.indexOf('{nomor && (')
  const popup = layar.slice(mulai, layar.indexOf('</Modal>', mulai))

  it('modal tanpaTutup; satu-satunya aksi OK mengirim berkas', () => {
    expect(mulai).toBeGreaterThan(0)
    expect(popup).toContain('tanpaTutup')
    expect(popup).toContain('{TOMBOL.ok}')
    expect(popup).toContain('kirim()')
  })
})

// perintah work owner 06-10-2026 (RALAT "HIDE AJA"): "jangan di hide, tapi di disable aja; muncul untuk XOL, tapi di disable"
describe('tombol Enable / Disable Input Type tampil untuk XOL tetapi nonaktif', () => {
  it('LayarKasus merender tombol bersyarat XOL, disabled, tanpa aksi', () => {
    const layar = readFileSync(join(__dirname, 'pages', 'LayarKasus.tsx'), 'utf8')
    const i = layar.indexOf('{TOMBOL.enableDisable}')
    expect(i).toBeGreaterThan(0)
    const blok = layar.slice(layar.lastIndexOf("nilai(h, POLIS + 'TreatyType') === 'XOL'", i), i)
    expect(blok).toContain('<button type="button" className="btn" disabled>')
    expect(blok).not.toContain('onClick')
    expect(layar).not.toContain("aksi: 'TreatyEnableDisableInput'")
  })
})
