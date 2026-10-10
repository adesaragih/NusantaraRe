// Paritas form objek Aneka (tiket 73 P1) - `Section\InputDtlObjectAneka_FacIn.xml` kontainer C1-C9. Data uji sintetis.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import type { TampilanAneka } from '../api'
import { FORM_MEMORANDUM, FORM_OBJEK_ANEKA as F, TEKS_ANEKA } from '../labels'
import { objekAnekaBaru } from './aneka'
import FormObjekAneka, { baca, kelompokkan, KONTAINER, medanTampil, tulis } from './FormObjekAneka'

const SUMBER = readFileSync(join(__dirname, 'FormObjekAneka.tsx'), 'utf8').replace(/\r\n/g, '\n')

const T = (p: Partial<TampilanAneka>): TampilanAneka => ({
  kontainer: [], kolomTahunJumlah: false, he: false, golf: false, tanamanHasil: false, engineering: false, bonding: false,
  briguna: false, topRisk: false, sectionLol: false, jumlahSatuan: true, currencyCoverage: false, mbd: false, ...p,
})

const label = (t: TampilanAneka) => medanTampil(t).map((m) => m.label)

describe('FormObjekAneka', () => {
  it('jumlah sel per kontainer = urai korpus (84 FIELD dikurangi Briguna 9, Others, Memorandum, pyTemplateInputBox C2)', () => {
    // Urai sel `InputDtlObjectAneka_FacIn.xml` (sesi 0f, an2a_idoa_cells): C1 11 + Briguna 9 + Others 1; C2 15 (+1 tombol 1=2);
    // C3 5; C4 8; C5 3; C6 8; C7 9; C8 6; C9 8.
    const n = Object.fromEntries(Object.entries(KONTAINER).map(([k, v]) => [k, v.length]))
    expect(n).toEqual({ C1: 11, C2: 15, C3: 5, C4: 8, C5: 3, C6: 8, C7: 9, C8: 6, C9: 8 })
    // Rekonsiliasi dua arah dengan sensus urai (84 sel FIELD): 73 sel kontainer + Briguna 9 + Others 1 + tombol C2 `1=2` 1.
    const total = Object.values(KONTAINER).reduce((s, v) => s + v.length, 0)
    expect(total + 9 + 1 + 1).toBe(84)
  })

  it('Liability (C1 saja): urutan sel; Draft Wording / Maintenance Period tersembunyi (OTHER IsBonding / IsEngineering)', () => {
    expect(label(T({ kontainer: ['C1'] }))).toEqual([
      F.objectDescription, F.section, F.noKontrak, F.startDate, F.endDate, F.noKtp, F.noNpwp, F.currency, F.tsi,
    ])
    expect(label(T({ kontainer: ['C1'], bonding: true, engineering: true }))).toContain(F.draftWordingId)
    expect(label(T({ kontainer: ['C1'], bonding: true, engineering: true }))).toContain(F.maintenancePeriod)
  })

  it('Boiler (C1 + C4 + C5): properti yang sama tampil sekali (A-1), kontainer pertama menang', () => {
    const m = medanTampil(T({ kontainer: ['C1', 'C4', 'C5'] }))
    const jalur = m.map((x) => x.jalur)
    expect(jalur.filter((j) => j === 'objectName')).toHaveLength(1)
    expect(jalur.filter((j) => j === 'section')).toHaveLength(1)
    expect(m.find((x) => x.jalur === 'objectName')?.label).toBe(F.objectDescription)
    expect(jalur).toEqual(expect.arrayContaining(['vehicleHe.chassisNumber', 'serialNo', 'year', 'quantity', 'unitQuantity']))
  })

  it('C5 Quantity / Unit bersyarat `BusinessCode!=20` (jumlahSatuan)', () => {
    expect(label(T({ kontainer: ['C5'], jumlahSatuan: false }))).toEqual([F.year])
    expect(label(T({ kontainer: ['C5'] }))).toEqual([F.year, F.quantity, F.unit])
  })

  it('C7 Aviation Hull tanpa Currency dan TSI', () => {
    const l = label(T({ kontainer: ['C7'] }))
    expect(l).not.toContain(F.currency)
    expect(l).not.toContain(F.tsi)
    expect(l[0]).toBe(F.modelTipe)
  })

  it('C2 Marine Hull: Production Year, Place of Manufacture, Flag, MarineHull = tampil-saja (pyReadOnlyCondition Always)', () => {
    const m = medanTampil(T({ kontainer: ['C2'] }))
    expect(m.filter((x) => x.jenis === 'tampil').map((x) => x.jalur)).toEqual([
      // Per baris layar Pega (work owner 09-10-2026): Type of Vessel | Classification | Construction, NRT | GRT | DWT, Length | Breadth | Depth.
      'year', 'builtIn', 'flag', 'marineHull.typeOfVessel', 'marineHull.classification', 'marineHull.construction', 'marineHull.nrtrt',
      'marineHull.grtrt', 'marineHull.dwtrt', 'marineHull.lengthm', 'marineHull.breadthm', 'marineHull.depthm',
    ])
    expect(m[0]).toMatchObject({ jalur: 'objectName', label: F.vesselName, jenis: 'kapal' })
  })

  it('render: Bonding menampilkan Akseptasi Memorandum; Briguna hanya bila IsCreditBriguna; tanpa kontainer = pesan', () => {
    const o = objekAnekaBaru()
    const bonding = renderToStaticMarkup(<FormObjekAneka o={o} ubah={() => {}} tampilan={T({ kontainer: ['C1'], bonding: true })} mataUang={['IDR']} />)
    expect(bonding).toContain(FORM_MEMORANDUM.judul)
    expect(bonding).toContain(FORM_MEMORANDUM.principalsGuarantee)
    expect(bonding).not.toContain(F.briguna.loanName)
    const kredit = renderToStaticMarkup(<FormObjekAneka o={o} ubah={() => {}} tampilan={T({ kontainer: ['C1'], briguna: true })} mataUang={[]} />)
    expect(kredit).toContain(F.briguna.loanName)
    expect(kredit).not.toContain(FORM_MEMORANDUM.judul)
    expect(renderToStaticMarkup(<FormObjekAneka o={o} ubah={() => {}} tampilan={T({})} mataUang={[]} />)).toContain(TEKS_ANEKA.tanpaKontainer)
  })

  it('kartu (perapian 09-10-2026): Object / Period / Sum Insured / Identity; HE: Object Description pertama, detail di Vehicle', () => {
    const lia = kelompokkan(medanTampil(T({ kontainer: ['C1'] })))
    expect(lia.objek.map((m) => m.jalur)).toEqual(['objectName', 'section', 'noKontrak'])
    expect(lia.periode.map((m) => m.jalur)).toEqual(['startDate', 'endDate'])
    expect(lia.nilai.map((m) => m.jalur)).toEqual(['currency', 'tsi'])
    expect(lia.identitas.map((m) => m.jalur)).toEqual(['noKtp', 'noNpwp'])
    const he = kelompokkan(medanTampil(T({ kontainer: ['C5', 'C6'] })))
    expect(he.objek.map((m) => m.jalur)).toEqual(['objectName', 'year', 'quantity', 'unitQuantity'])
    expect(he.kendaraan.map((m) => m.jalur)).toEqual([
      'vehicleHe.objectNameHe', 'vehicleHe.chassisNumber', 'vehicleHe.engineNumber', 'vehicleHe.brandName', 'vehicleHe.typeName',
    ])
    const kapal = kelompokkan(medanTampil(T({ kontainer: ['C2'] })))
    expect(kapal.objek.map((m) => m.jalur)).toEqual(['objectName'])
    expect(kapal.kapal).toHaveLength(12)
    // Tidak ada medan yang hilang saat dikelompokkan.
    for (const k of Object.keys(KONTAINER)) {
      const m = medanTampil(T({ kontainer: [k] }))
      expect(Object.values(kelompokkan(m)).flat()).toHaveLength(m.length)
    }
  })

  it('baca / tulis jalur bersarang tanpa mengubah objek asal', () => {
    const o = objekAnekaBaru()
    const b = tulis(o, 'aviationHull.pilots', '2')
    expect(baca(b, 'aviationHull.pilots')).toBe('2')
    expect(baca(o, 'aviationHull.pilots')).toBe('')
    expect(baca(tulis(o, 'others', 'X'), 'others')).toBe('X')
  })

  it('Object Name / Merk HE teks bebas (K73-10) memicu SetObjectDescHE_DT; Type dari lookup tipe-kendaraan', () => {
    expect(SUMBER).toContain("ubah(ubahHe(o, { [anak]: nilai, [idnya]: '' }))")
    expect(SUMBER).toContain("onPilih={(s) => ubah(ubahHe(o, { typeName: s.label, type: s.id }))}")
    expect(SUMBER).toContain('daftarTipeKendaraan()')
  })
})

describe('FormObjekAneka - urutan kartu kecil', () => {
  it('Sum Insured paling kiri, lalu Identity, lalu Period (work owner 09-10-2026)', () => {
    const html = renderToStaticMarkup(<FormObjekAneka o={objekAnekaBaru()} ubah={() => {}} tampilan={T({ kontainer: ['C1'] })} mataUang={['IDR']} />)
    const posisi = [TEKS_ANEKA.grup.nilai, TEKS_ANEKA.grup.identitas, TEKS_ANEKA.grup.periode].map((j) => html.indexOf(`>${j}</h5>`))
    expect(posisi.every((p) => p >= 0)).toBe(true)
    expect([...posisi].sort((a, b) => a - b)).toEqual(posisi)
  })
})
