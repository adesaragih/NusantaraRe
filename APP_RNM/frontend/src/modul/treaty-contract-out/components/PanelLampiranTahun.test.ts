// Uji panel lampiran tahun treaty — tiket 12 Treaty Contract Out.
//
// Aturan murni diuji langsung; kabel dan larangan diuji dari sumbernya.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { AKAR_APLIKASI } from '../../../../../inti/frontend/uji/sumber'
import { LAMPIRAN_TCO } from '../labels'
import type { LampiranTahun } from '../api'
import { jalurIsiLampiran, jalurSemuaLampiran } from '../api'
import { adaTerkirim, bolehUlangi, bolehUnduh, labelStatusLampiran } from './PanelLampiranTahun'

function tanpaKomentar(teks: string): string {
  return teks
    .split('\n')
    .filter((b) => !b.trimStart().startsWith('//') && !b.trimStart().startsWith('*'))
    .join('\n')
}

const PANEL = tanpaKomentar(readFileSync(join(__dirname, 'PanelLampiranTahun.tsx'), 'utf8'))
const API = readFileSync(join(__dirname, '..', 'api.ts'), 'utf8')
const bagianApiLampiran = API.slice(API.indexOf('Treaty Contract Out tiket 12'))

function lampiran(status: LampiranTahun['status']): LampiranTahun {
  return {
    id: '1000000001', idTreatyYear: '1000001', fileName: 'UJI-kontrak.pdf', fileMimeType: 'application/pdf',
    category: 'CLAUSES', userId: 'UJI-ADMIN', tglUpload: '2026-09-29 09:00:00', status,
    percobaan: 1, galat: '',
  }
}

describe('aturan panel lampiran', () => {
  it('label status memakai kosakata tiket 12', () => {
    expect(labelStatusLampiran('terkirim')).toBe(LAMPIRAN_TCO.statusTerkirim)
    expect(labelStatusLampiran('tertunda')).toBe(LAMPIRAN_TCO.statusTertunda)
    expect(labelStatusLampiran('gagal')).toBe(LAMPIRAN_TCO.statusGagal)
  })
  it('hanya yang terkirim dapat diunduh; yang belum dapat diulang', () => {
    expect(bolehUnduh(lampiran('terkirim'))).toBe(true)
    expect(bolehUnduh(lampiran('tertunda'))).toBe(false)
    expect(bolehUlangi(lampiran('gagal'))).toBe(true)
    expect(bolehUlangi(lampiran('tertunda'))).toBe(true)
    expect(bolehUlangi(lampiran('terkirim'))).toBe(false)
    expect(adaTerkirim([lampiran('tertunda'), lampiran('terkirim')])).toBe(true)
    expect(adaTerkirim([lampiran('gagal')])).toBe(false)
  })
  it('tco4: ukuran tidak tampil — M_ATTACHMENTTREATY_2 tidak menyimpannya (GetAllAttachment2_Sql b84)', () => {
    expect(PANEL).not.toMatch(/\.ukuran|teksUkuran/)
  })
  it('jalur unduh menyebut tahun treaty dan lampirannya', () => {
    expect(jalurIsiLampiran('1000001', '1000000001')).toBe(
      '/api/treaty-contract-out/tahun/1000001/lampiran/1000000001/isi',
    )
    expect(jalurSemuaLampiran('1000001')).toBe('/api/treaty-contract-out/tahun/1000001/lampiran/semua')
  })
})

describe('kabel dan larangan', () => {
  it('unduhan lewat fetch berheader identitas, BUKAN tautan biasa', () => {
    expect(PANEL).toContain('unduhBerkasBeridentitas(')
    expect(PANEL).not.toMatch(/href=/)
    expect(PANEL).not.toMatch(/window\.open|location\.href/)
    // Refactor bentuk B: pengunduh beridentitas kini pembantu klien bersama
    // (`inti/klien.ts`) - ia membaca amplop galat, dan amplop satu rumah.
    const klien = readFileSync(join(AKAR_APLIKASI, 'inti', 'frontend', 'klien.ts'), 'utf8')
    expect(klien).toContain('export async function unduhBerkasBeridentitas')
    const fungsi = klien.slice(klien.indexOf('export async function unduhBerkasBeridentitas'))
    expect(fungsi).toContain('...headerIdentitas()')
  })
  it('unggah multipart lewat mintaFormulir yang membawa identitas', () => {
    expect(bagianApiLampiran).toContain('mintaFormulir<HasilLampiranTahun>(')
  })
  it('seluruh label panel dari LAMPIRAN_TCO, bukan literal', () => {
    for (const kunci of ['attachmentFor', 'forTreatyContractOut', 'addAttachment', 'refresh', 'downloadAll',
      'kolomFileName', 'kolomType', 'delete'] as const) {
      expect(PANEL).toContain(`LAMPIRAN_TCO.${kunci}`)
    }
    expect(PANEL).not.toMatch(/>\s*(Add attachment|Download All|Refresh|Delete)\s*</)
  })
  it('layar hanya menyebut tahun treaty — TREATYID (TreatyYear + TreatyYearID) dirakit backend (tco4)', () => {
    for (const teks of [PANEL, bagianApiLampiran]) {
      expect(teks).not.toMatch(/treatyIn\b|treaty_in|TREATYID|M_ATTACHMENTTREATY/i)
    }
  })
  it('status tampil apa adanya dan galat terakhirnya terlihat', () => {
    expect(PANEL).toContain('labelStatusLampiran(')
    expect(PANEL).toContain('.galat')
  })
})

describe('bahasa Inggris (keputusan work owner 30-09-2026)', () => {
  it('kode perbaikan dan kepala kolom tombol tidak tampil mentah', () => {
    const kode = readFileSync(join(__dirname, 'PanelLampiranTahun.tsx'), 'utf8')
    expect(kode).toContain('LAMPIRAN_TCO.perbaikanUlangi')
    expect(kode).not.toMatch(/\(\{t\.perbaikan\}\)|aria-label="aksi"/)
  })
})
