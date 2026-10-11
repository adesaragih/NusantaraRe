// Tab Object MARINE CARGO (tiket 74 P1) - label VERBATIM korpus, gerbang layar, fungsi murni, dan render form cargo.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import type { BarisKapal } from '../api'
import { FORM_CARGO as F, GRID_CARGO, GRID_CARGO_UW, POPUP_AGENT_CARGO, POPUP_KAPAL_CARGO } from '../labels'
import { cargoBaru, rapikanCargo, tampilTanggal, terapkanAgent, terapkanKapal, terapkanKapalTambahan, kapalTambahanBaru } from './cargo'
import { FormCargo } from './TabObjectCargo'

const INWARD = readFileSync(join(__dirname, '..', 'pages', 'InwardFacultative.tsx'), 'utf8').replace(/\r\n/g, '\n')

const KAPAL: BarisKapal = {
  id: 'K1',
  nama: 'UJI KAPAL',
  flag: 'UJI',
  tahun: '2001',
  nrt: '1,5',
  grt: '2,25',
  dwt: '3',
  klasifikasi: 'KL',
  konstruksi: 'BAJA',
  shipTypeId: '7',
  shipTypeName: 'TIPE',
  register: 'R1',
}

describe('TabObjectCargo - label korpus', () => {
  it('kolom grid CargoList non-UW / UW = pyValue InputDtlObject_FacIn(_IsUW)', () => {
    expect(GRID_CARGO).toEqual(['No.', 'Top Risk', 'SailDate', 'Goods', 'Ship Name', 'Conveyance'])
    expect(GRID_CARGO_UW).toEqual(['No.', 'Top Risk', 'Trading', 'Goods', 'Packing', 'Conveyance'])
  })

  it('label form = pyLabelFieldValue (bukan pyLabelFor "Barang" / "Pengakutan")', () => {
    expect([F.trading, F.goods, F.conveyance, F.packing, F.pilih]).toEqual(['Trading', 'Goods', 'Conveyance', 'Packing', 'Choose'])
    expect([F.grupAgent, F.grupLc, F.grupKapal, F.grupKapalTambahan]).toEqual(['Survey Agent', 'LC', 'Ship Data', 'Additional Ship'])
    expect([F.sailDate, F.fromRute, F.toRute, F.blNumber, F.invoiceNumber, F.invoiceDate, F.maxAge, F.maxAgeDetail]).toEqual([
      'Sail Date', 'From Rute', 'To Rute', 'No BL', 'Invoice Number', 'Invoice Date', 'Max Age ID', 'Max Age Detail',
    ])
    expect(F.kolomKapalTambahan).toEqual(['Register ID', 'Ship Name', 'Voyage No'])
    expect(POPUP_AGENT_CARGO.kolom).toEqual(['ID', 'Name', 'Address', 'City', 'Country'])
    expect(POPUP_KAPAL_CARGO.kolom).toHaveLength(16)
  })
})

describe('TabObjectCargo - gerbang layar', () => {
  it('Marine Cargo hanya menambah cabang: FIRE dan Aneka tetap menang', () => {
    expect(INWARD).toContain('const kasusCargo = !kasusFire && !kasusAneka && lini?.marineCargo === true')
    expect(INWARD).toContain("tab === 'Object' && kasusCargo ? (\n              <TabObjectCargo caseId={kasus.caseId} />")
  })
})

describe('cargo.ts', () => {
  it('SelectShip_act: GRT / NRT / DWT koma desimal → titik; Age dikosongkan (dihitung server)', () => {
    const c = terapkanKapal({ ...cargoBaru(), ship: { ...cargoBaru().ship, age: '9' } }, KAPAL)
    expect(c.ship).toMatchObject({ id: 'K1', nmShip: 'UJI KAPAL', grt: '2.25', nrt: '1.5', dwt: '3', yMake1: '2001', shipTypeName: 'TIPE', age: '' })
  })

  it('Survey Agent per cargo (K74-6); Additional Ship mengambil medan autocomplete, Voyage No tetap', () => {
    const c = terapkanAgent(cargoBaru(), { id: 'A1', name: 'AGEN', address: 'ALAMAT', city: 'KOTA', country: 'ID', contact: '' })
    expect(c.surveyAgent).toEqual({ agentId: 'A1', fullName: 'AGEN', address: 'ALAMAT', city: 'KOTA', country: 'ID' })
    const t = terapkanKapalTambahan({ ...kapalTambahanBaru(), voyageNo: 'V9' }, KAPAL)
    expect(t).toMatchObject({ id: 'K1', register: 'R1', nmShip: 'UJI KAPAL', voyageNo: 'V9' })
  })

  it('rapikanCargo mengisi halaman anak yang absen; tanggal kabel DD-MM-YYYY tampil dd/mm/yyyy', () => {
    const r = rapikanCargo({ id: 'X', goodNote: 'G' })
    expect(r.ship.nmShip).toBe('')
    expect(r.lc.lc).toBe(false)
    expect(r.additionalShip).toEqual([])
    expect(tampilTanggal('05-11-2026')).toBe('05/11/2026')
  })
})

describe('FormCargo - render', () => {
  it('urutan kartu: kepala, lalu kiri Survey Agent / LC / Ship Data, kanan pelayaran / Additional Ship', () => {
    const h = renderToStaticMarkup(
      <FormCargo c={cargoBaru()} idx={0} ubah={() => {}} master={{ trading: [{ id: '1', description: 'IMPOR' }], goods: [], conveyance: [], packing: [] }} />,
    )
    const urut = [F.trading, F.conveyance, F.goods, F.packing, F.grupAgent, F.grupLc, F.grupKapal, F.sailDate, F.maxAgeDetail, F.grupKapalTambahan]
    const posisi = urut.map((t) => h.indexOf(`>${t}<`))
    expect(posisi.every((p) => p >= 0)).toBe(true)
    expect([...posisi].sort((a, b) => a - b)).toEqual(posisi)
    expect(h).toContain('<option value="1">IMPOR</option>')
    expect(h).toContain(`>${F.selectShip}<`)
    expect(h).toContain(`>${F.selectAgent}<`)
  })
})
