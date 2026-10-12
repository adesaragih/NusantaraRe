// Menu View Policy (tiket 71).

import { existsSync, readFileSync } from 'node:fs'
import { join } from 'node:path'
import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import type { RingkasanPolis } from '../api'
import { FOLLOWING as F, VIEW_POLICY as V } from '../labels'
import ViewPolicy, { labelPolicyType, RingkasanPolisLama } from './ViewPolicy'

const DDL = 'D:\\migrasi\\RNM\\DDL\\'
const baca = (...p: string[]) => readFileSync(join(__dirname, ...p), 'utf8').replace(/\r\n/g, '\n')
const SUMBER = baca('ViewPolicy.tsx')

describe('View Policy - label korpus', () => {
  it.skipIf(!existsSync(DDL + 'ViewOpenPolicyMarineInbox.xml'))('"List Polis" + tab Facultative = ViewOpenPolicyMarineInbox', () => {
    const xml = readFileSync(DDL + 'ViewOpenPolicyMarineInbox.xml', 'utf-8')
    expect(xml).toContain(`<pyValue>${V.daftar}</pyValue>`)
    expect(xml).toContain('<pyTitle>Facultative</pyTitle>')
    expect(xml).toContain(`>${V.viewData}<`)
  })
  it.skipIf(!existsSync(DDL + 'GridOpenPolicyMarine.xml'))('kolom grid + View Details = GridOpenPolicyMarine (sama dengan Following)', () => {
    const xml = readFileSync(DDL + 'GridOpenPolicyMarine.xml', 'utf-8')
    for (const k of F.kolom) expect(xml, k).toContain(`>${k}<`)
    expect(xml).toContain(`<pyLabel>${F.lihat}</pyLabel>`)
  })
  it.skipIf(!existsSync(DDL + 'PolicyType.xml'))('label Policy Type = DDL PolicyType.xml', () => {
    const xml = readFileSync(DDL + 'PolicyType.xml', 'utf-8')
    for (const k of ['0', '1', '2']) expect(xml).toMatch(new RegExp(`<pyLocalizedValue>${labelPolicyType(k)}</pyLocalizedValue>\\s*<pyStandardValue>${k}</pyStandardValue>`))
    expect(labelPolicyType('9')).toBe('9')
  })
})

describe('View Policy tab Treaty / Life - label korpus (DDL kiriman work owner 11-10-2026)', () => {
  const grid: [string, readonly string[]][] = [
    ['GridOpenoffer', V.kolom.offerFacultative],
    ['GridOpenPolicyMarineTreaty', V.kolom.polisTreaty],
    ['GridOpenPolicyLife', V.kolom.polisLife],
    ['GridOpenOfferLife', V.kolom.offerLife],
  ]
  for (const [nama, kolom] of grid) {
    it.skipIf(!existsSync(DDL + nama + '.xml'))(`kolom ${nama} verbatim`, () => {
      const xml = readFileSync(DDL + nama + '.xml', 'utf-8')
      for (const k of kolom) expect(xml, k).toContain(`<pyValue>${k}</pyValue>`)
    })
  }
  it.skipIf(!existsSync(DDL + 'FilterForViewPolicy.xml'))('View Option Treaty = FilterForViewPolicy (nilai standar + label)', () => {
    const xml = readFileSync(DDL + 'FilterForViewPolicy.xml', 'utf-8')
    for (const o of V.opsiTreaty) expect(xml).toMatch(new RegExp(`<pyLocalizedValue>${o.label}</pyLocalizedValue>\\s*<pyStandardValue>${o.value}</pyStandardValue>`))
  })
  it('View Option Offer Facultative = polis Facultative tanpa QQ Name', () => {
    expect(V.opsiOfferFacultative).toEqual(F.opsi.filter((o) => o.value !== '7'))
  })
})

describe('ViewPolicy', () => {
  it('render awal: judul, View Option + Filter, grid 8 kolom + View Details', () => {
    const html = renderToStaticMarkup(<ViewPolicy />)
    expect(html).toContain(`<h2 class="inbox__judul nbf-tema__judul">${V.judul}</h2>`)
    expect(html).toContain(`placeholder="${F.placeholder}"`)
    // Polis Facultative 6 + Offer Facultative 5 + Treaty 5 (Life tanpa View Option, work owner 11-10-2026).
    expect(html.match(/<option value="\d">/g)).toHaveLength(16)
    for (const k of F.kolom) expect(html).toContain(`<th scope="col">${k}</th>`)
  })
  it('tiga tab Facultative / Treaty / Life; Facultative = List Polis + List Offer, Treaty = List Polis saja, Life = keduanya', () => {
    const html = renderToStaticMarkup(<ViewPolicy />)
    for (const t of V.tab) expect(html).toContain(`>${t}<`)
    expect(html.match(new RegExp(`<h3 class="nbf-vp__subjudul">${V.daftar}</h3>`, 'g'))).toHaveLength(3)
    expect(html.match(new RegExp(`<h3 class="nbf-vp__subjudul">${V.daftarOffer}</h3>`, 'g'))).toHaveLength(2)
    for (const kolom of Object.values(V.kolom)) for (const k of kolom) expect(html).toContain(`<th scope="col">${k}</th>`)
    // Detail Treaty / Life belum diport: tanpa kolom View Details; Offer Facultative membuka case baca-saja.
    expect(SUMBER).toContain("{detail?.jenis === 'kasus' && <InwardFacultative kasus={{ caseId: detail.caseId }} onBatal={() => setDetail(null)} hanyaLihat />}")
    expect(SUMBER).toContain('cari={(_o, f, h) => cariPolisLife(f, h)}')
    expect(SUMBER).not.toContain('cariOfferTreaty')
  })
  it('detail: dropdown View Data hanya bila versi > 1; versi sistem baru = layar case, lainnya ringkasan (DetailPolis bersama)', () => {
    const DETAIL = readFileSync(join(__dirname, '..', 'components', 'DetailPolis.tsx'), 'utf8').replace(/\r\n/g, '\n')
    expect(DETAIL).toContain('{versi !== null && versi.length > 1 && (')
    expect(DETAIL).toContain('{kini !== null && kini.sistemBaru && <div key={kini.caseId}>{tampilKasus(kini.caseId)}</div>}')
    expect(DETAIL).toContain('if (kini === null || kini.sistemBaru) return')
    expect(SUMBER).toContain('tampilKasus={(caseId) => <InwardFacultative kasus={{ caseId }} onBatal={() => setDetail(null)} hanyaLihat />}')
  })
  it('ringkasan polis lama: kepala + baris + ceding, angka gaya Indonesia, Policy Type berlabel', () => {
    const v: RingkasanPolis = {
      kepala: {
        caseId: 'NB-UJI1', policyNo: 'UJI-POLIS-1', noEndors: '', insuredName: 'UJI TERTANGGUNG', qqName: '', sourceOfBusiness: 'UJI SOB',
        cedingCo: 'UJI CEDING', startPeriod: '01-01-2026', endPeriod: '01-01-2027', policyType: '1', reffNumber: 'UJI-REFF', groupName: '',
        statusBusiness: '', noMasterPolis: '',
      },
      baris: [{ objectNo: '1', riskLocation: 'UJI LOKASI', objectItem: 'BUILDING', coverage: 'FLEXAS', accumulationCode: 'UJI-AK', currency: 'IDR',
        tsi100: '1000000', tsiRnm: '100000', rate: '0.5', premiRnm: '50', treaty: 'QS', pctShare: '10', riComm: '5' }],
      ceding: [{ cedingName: 'UJI CEDING', pctShare: '100' }],
    }
    const html = renderToStaticMarkup(<RingkasanPolisLama v={v} />)
    expect(html).toContain('<dd>Master Policy</dd>')
    expect(html).toContain('<dd>01/01/2026 – 01/01/2027</dd>')
    expect(html).toContain('<td class="nbf-angka">1.000.000,0000</td>')
    for (const c of V.kolomRingkasan) expect(html).toContain(`<th scope="col">${c}</th>`)
  })
  it('menu View Policy di bawah Home (menuAtas, work owner 09-10-2026); rute merender halamannya', () => {
    expect(baca('..', 'menu.ts')).toContain("menuAtas: [{ label: VIEW_POLICY.judul, halaman: 'nbfacin-view-policy', semua: true }]")
    const rute = baca('..', 'rute.tsx')
    expect(rute).toContain("<div hidden={halaman !== 'nbfacin-view-policy'}>")
    expect(rute).toContain('<ViewPolicy onLihatBerkas={onLihatBerkas} onBukaModul={onBukaModul} />')
  })
})

describe('View Policy - pilihan kosong View Option', () => {
  it('"Choose" saja (work owner 11-10-2026), di List Polis / List Offer Facultative dan List Polis Treaty', () => {
    const html = renderToStaticMarkup(<ViewPolicy />)
    expect(html.match(new RegExp(`<option value="" selected="">${F.pilihKosong}</option>`, 'g'))).toHaveLength(3)
    expect(html).not.toContain(`${F.viewOption}: ${F.pilihKosong}`)
  })
})

describe('View Policy - View Details Treaty (work owner 11-10-2026)', () => {
  it('tombol View Details di List Polis Treaty bila onLihatBerkas ada; berkas dicari lalu dibuka di jendela tampilan-saja', () => {
    const tanpa = renderToStaticMarkup(<ViewPolicy />)
    const dengan = renderToStaticMarkup(<ViewPolicy onLihatBerkas={() => true} />)
    const kolomLihat = (h: string) => (h.match(/<th scope="col"><\/th>/g) ?? []).length
    expect(kolomLihat(dengan)).toBe(kolomLihat(tanpa) + 1)
    // Di tempat (onBukaModul) seperti detail NB FacIn; jendela onLihatBerkas hanya cadangan.
    expect(SUMBER).toContain('const bukaBerkas = onBukaModul ?? onLihatBerkas')
    expect(SUMBER).toContain('onLihat={bukaBerkas ? lihatTreaty : undefined}')
    expect(SUMBER).toContain('if (bukaBerkas?.(x.modul, x.kasus) !== true) setPesanTreaty(V.teks.berkasTakTerpasang)')
    expect(SUMBER).toContain(
      "if (g instanceof ApiFailure && g.status === 404 && g.detail.code === 'DITOLAK_BACKEND') setPesanTreaty(V.teks.polisTanpaBerkas(b.policyNo))",
    )
  })
})
