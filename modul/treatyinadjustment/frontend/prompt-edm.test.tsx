// Uji teks pilihan `EDMState` / `EDMMaterialType` — *(8 Okt, E)*.
//
// Bukti: `D:\XML_NURE\_migration-docs\treaty-in-adjustment\ekspor-tambahan\`
// `EDMState.xml` (1 Internal @6122, 2 External @6296) dan
// `EDMMaterialType.xml` (1 Material @6154, 2 Non Material @6328), kelas
// `ASM-FW-GISFW-Int-TREATY_IN` — sama dengan peta `promptValue` modul Treaty In.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'
import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import type { BarisPenyesuaian } from './api'
import PilihMaster from './komponen/PilihMaster'
import { PROMPT_EDM, teksPromptEDM } from './labelsPromptEDM'
import { cocokCari, selDaftar } from './pages/PenyesuaianKontrak'

describe('prompt value EDMState / EDMMaterialType', () => {
  it('peta persis rule Property ekspor-tambahan', () => {
    expect(PROMPT_EDM).toEqual({
      EDMState: { '1': 'Internal', '2': 'External' },
      EDMMaterialType: { '1': 'Material', '2': 'Non Material' },
    })
  })

  it('kode di luar peta tampil apa adanya — `EDMState = 3` tidak ditebak', () => {
    expect(teksPromptEDM('EDMState', '1')).toBe('Internal')
    expect(teksPromptEDM('EDMMaterialType', '2')).toBe('Non Material')
    expect(teksPromptEDM('EDMState', '3')).toBe('3')
    expect(teksPromptEDM('EDMState', '')).toBe('')
    expect(teksPromptEDM('Lain', '1')).toBe('1')
  })

  it('radio Material Type picker revisi menampilkan teks, nilainya tetap kode', () => {
    const html = renderToStaticMarkup(
      <PilihMaster jenis="revisi" onTutup={() => undefined} onDraf={() => undefined} />,
    )
    expect(html).toMatch(/name="tria-picker-material"[^>]*value="1"|value="1"[^>]*name="tria-picker-material"/)
    expect(html).toContain('Material')
    expect(html).toContain('Non Material')
  })

  it('grid daftar: kolom Type / Material Type menampilkan teks, bukan kode', () => {
    const baris: BarisPenyesuaian = {
      id: '1002307/R01', idAsal: '1002307', jenisPenyesuaian: '2', jenisMaterial: '1',
      namaKontrak: 'X', sifatProporsi: '', asalBisnis: '', cedant: '',
      tanggalMulai: '', tanggalBerakhir: '', posisi: '', statusAkseptasi: '',
    }
    const sel = selDaftar(baris)
    expect([sel[2], sel[3]]).toEqual(['External', 'Material'])
    expect(selDaftar({ ...baris, jenisPenyesuaian: '1', jenisMaterial: '2' }).slice(2, 4)).toEqual(['Internal', 'Non Material'])
    // EDMState 3 = label kepala "Adjustment Premium" @82379; Material kosong tetap kosong
    expect(selDaftar({ ...baris, jenisPenyesuaian: '3', jenisMaterial: '' }).slice(2, 4)).toEqual(['Adjustment Premium', ''])
    // kode lain di luar peta apa adanya
    expect(selDaftar({ ...baris, jenisPenyesuaian: '9' })[2]).toBe('9')
    // pencarian mencocokkan teks yang tampil
    expect(cocokCari(baris, 'external')).toBe(true)
  })

  it('kepala mode detail memakai teks prompt untuk kedua dropdown', () => {
    const src = readFileSync(join(__dirname, 'pages', 'PenyesuaianKontrak.tsx'), 'utf8')
    // ⭐ 9 Oktober 2026 — kepala bentuk Pega: nilai TEKS (bukan kotak isian).
    expect(src).toContain("butir(PENYESUAIAN.jenisPenyesuaian, teksPromptEDM('EDMState', m.EDMState ?? ''))")
    expect(src).toContain("butir(PENYESUAIAN.jenisMaterial, teksPromptEDM('EDMMaterialType', m.EDMMaterialType ?? ''))")
  })
})
