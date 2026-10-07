// Paritas popup pencari AGENT (Change SOB tiket 33, Add / Select Ceding tiket 34) - harness `SOB` / section
// `CedingCoHierarki` + tangkapan layar work owner 03-10-2026.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { POPUP_SOB, TEKS_INWARD } from '../labels'
import PopupPilihAgent, { JEDA_CARI_MS } from './PopupPilihAgent'

const HTML = renderToStaticMarkup(<PopupPilihAgent judul={POPUP_SOB.judul} onTutup={() => {}} onPilih={() => {}} />)
const SUMBER = readFileSync(join(__dirname, 'PopupPilihAgent.tsx'), 'utf8').replace(/\r\n/g, '\n')

describe('PopupPilihAgent', () => {
  it('judul dari pemanggil, kotak Search dapat diisi, tanpa tombol Search terpisah', () => {
    expect(HTML).toContain('role="dialog"')
    expect(HTML).toContain(POPUP_SOB.judul)
    expect(HTML).toMatch(new RegExp(`class="field__label">${POPUP_SOB.cari.label}<`))
    expect(HTML).not.toMatch(/<input[^>]*readonly/)
    expect(HTML).not.toMatch(/>Search<\/button>/)
  })

  it('kolom grid: ID, Client ID, Name, lalu kolom Choose', () => {
    const kolom = [...HTML.matchAll(/<th[^>]*>([^<]*)<\/th>|<th[^>]*\/>/g)].map((m) => m[1] ?? '')
    expect(kolom).toEqual([...POPUP_SOB.kolom, ''])
  })

  it('belum ada baris sebelum jawaban datang', () => {
    expect(HTML).not.toContain('<tbody')
    expect(HTML).not.toContain(TEKS_INWARD.tanpaSob)
  })

  it('baris: id, clientId, name; Choose memanggil onPilih', () => {
    expect(SUMBER).toMatch(/<td>\{b\.id\}<\/td>\s*<td>\{b\.clientId\}<\/td>\s*<td>\{b\.name\}<\/td>/)
    expect(SUMBER).toContain('onClick={() => onPilih(b)}')
    expect(SUMBER).toContain('{POPUP_SOB.pilih.label}')
  })

  it('dimuat saat dibuka; mengetik mencari sesudah jeda (jadwal lama dibatalkan); Enter mencari dari halaman 1; paging memakai kata cari terakhir; jawaban lama dibuang', () => {
    expect(SUMBER).toContain("window.setTimeout(() => void muat(kotak, 1), kotak === kunciCari.current ? 0 : JEDA_CARI_MS)")
    expect(SUMBER).toContain('return () => window.clearTimeout(jadwal)')
    expect(SUMBER).toContain('}, [kotak])')
    expect(JEDA_CARI_MS).toBeGreaterThan(0)
    expect(SUMBER).toContain('onKirim={() => void muat(kotak, 1)}')
    expect(SUMBER).toContain('onPindah={(h) => void muat(kunciCari.current, h)}')
    expect(SUMBER).toContain('if (nomor === nomorPermintaan.current) setHasil(h)')
    expect(SUMBER).toContain('<Kosong pesan={TEKS_INWARD.tanpaSob} />')
  })
})
