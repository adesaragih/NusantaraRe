// Tab utama Loss Record kasus FIRE (tiket 55) - data uji sintetis.

import { existsSync, readFileSync } from 'node:fs'
import { join } from 'node:path'
import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import type { TampilanLossRecord } from '../api'
import { LOSS_RECORD as L, TEKS_LOSS_RECORD as T, TEMPLATE_CSV_LOSS_RECORD } from '../labels'
import { angkaIndonesia, bacaCsvLossRecord, cocokkanKelasBisnis, pecahBarisCsv, tanggalCsv } from './csvLossRecord'
import { saranSama } from './IsianSaran'
import TabLossRecord, { persenRasio } from './TabLossRecord'

const NB = 'D:\\migrasi\\RNM\\NB FacIn\\Section\\'
const TEMPLATE = 'D:\\migrasi\\RNM\\DDL\\templatecsvlossrecord.csv'
const SUMBER = readFileSync(join(__dirname, 'TabLossRecord.tsx'), 'utf8').replace(/\r\n/g, '\n')
const HALAMAN = readFileSync(join(__dirname, '../pages/InwardFacultative.tsx'), 'utf8').replace(/\r\n/g, '\n')

describe.skipIf(!existsSync(NB + 'ObjectList_LossRecord.xml'))('label = korpus', () => {
  const baca = (f: string) => readFileSync(NB + f, 'utf-8')
  it('bar unggah, premi, grid (ObjectList_LossRecord)', () => {
    const xml = baca('ObjectList_LossRecord.xml')
    for (const t of [L.judulUnggah, L.unduhTemplate, L.unggah, L.simpan, L.premiPar, L.premiMb, L.premiPl, L.premiOther,
      L.uwYear, L.classOfBusiness, L.dateOfLoss, L.causeOfLoss, L.currency, L.lossAmount, L.lossAmountIdr]) {
      expect(xml, t).toContain(`>${t}<`)
    }
  })
  it('total dan loss ratio (SummaryLossRecord_Section)', () => {
    const xml = baca('SummaryLossRecord_Section.xml')
    for (const t of [L.classOfBusiness, L.currency, L.totalLoss, L.totalLossIdr, L.gross, L.net, ...L.rasio]) expect(xml, t).toContain(`>${t}<`)
  })
  it('tombol Save tab (InputInwardFacultativeDtl, kontainer Loss Record)', () => {
    const xml = baca('InputInwardFacultativeDtl.xml')
    const i = xml.indexOf('<pyTitle>Loss Record</pyTitle>')
    expect(i).toBeGreaterThan(0)
    const kontainer = xml.slice(i, i + 80000)
    const tombol = kontainer.indexOf(`<pyLabel>${L.simpanTab}</pyLabel>`)
    expect(tombol).toBeGreaterThan(0)
    expect(kontainer.indexOf('SaveLossRecord', tombol)).toBeGreaterThan(tombol)
  })
  it('ubah baris (InputLossRecord_Sec)', () => {
    const xml = baca('InputLossRecord_Sec.xml')
    for (const t of [L.insuredName, L.dateOfLoss, L.currency, L.lossAmount, L.causeOfLossUbah, L.lossAmountIdr, L.classOfBusiness, L.uwYear]) {
      expect(xml, t).toContain(`>${t}<`)
    }
  })
})

describe.skipIf(!existsSync(TEMPLATE))('template CSV = berkas work owner', () => {
  it('isi persis (CRLF, pemisah ;)', () => {
    expect(TEMPLATE_CSV_LOSS_RECORD).toBe(readFileSync(TEMPLATE, 'utf-8'))
  })
})

describe('baca CSV Loss Record', () => {
  it('template terbaca: 9 baris, kolom menurut nama kepala, catatan kolom ke-10 diabaikan', () => {
    const h = bacaCsvLossRecord(TEMPLATE_CSV_LOSS_RECORD)
    expect(h.galat).toEqual([])
    expect(h.baris).toHaveLength(9)
    expect(h.baris[0]).toEqual({
      uwYear: '', classOfBusiness: 'PROPERTY ALL RISK', dateOfLoss: '', causeOfLoss: 'WINDSTORM', currency: 'IDR',
      lossAmount: '100000000', lossAmountIdr: '',
    })
    expect(h.baris[1]).toMatchObject({ dateOfLoss: '01-05-2024', causeOfLoss: 'FIRE', currency: 'USD', lossAmount: '10000' })
    expect(h.baris[4]).toMatchObject({ classOfBusiness: 'MACHINERY BREAK DOWN', lossAmount: '1000' })
  })
  it('angka format Indonesia, tanggal DD/MM/YYYY, sel berkutip', () => {
    expect(angkaIndonesia('100.000.000')).toBe('100000000')
    expect(angkaIndonesia('10.000,50')).toBe('10000.50')
    expect(angkaIndonesia('abc')).toBeNull()
    expect(tanggalCsv('03/01/2025')).toBe('03-01-2025')
    expect(tanggalCsv('')).toBe('')
    expect(tanggalCsv('31/02/2025')).toBeNull()
    expect(tanggalCsv('2025-01-03')).toBeNull()
    expect(pecahBarisCsv('"A;B";"x""y";3', ';')).toEqual(['A;B', 'x"y', '3'])
  })
  it('galat: kolom hilang, tanggal / jumlah salah, berkas tanpa baris - grid tidak diganti (L55-3)', () => {
    expect(bacaCsvLossRecord('CLASS_OF_BUSINESS;DATE_OF_LOSS\nA;01/01/2025').galat).toEqual([`${T.csvKolom}CAUSE_OF_LOSS, CURRENCY, LOSS_AMOUNT`])
    const h = bacaCsvLossRecord('CLASS_OF_BUSINESS,DATE_OF_LOSS,CAUSE_OF_LOSS,CURRENCY,LOSS_AMOUNT\nA,32/01/2025,x,IDR,1\nB,,y,IDR,abc')
    expect(h.baris).toEqual([])
    expect(h.galat).toEqual([`${T.csvBaris}2${T.csvTanggal}`, `${T.csvBaris}3${T.csvJumlah}`])
    expect(bacaCsvLossRecord('CLASS_OF_BUSINESS;DATE_OF_LOSS;CAUSE_OF_LOSS;CURRENCY;LOSS_AMOUNT\n;;;;\n').galat).toEqual([T.csvKosong])
  })
})

describe('CSV: Class Of Business wajib dari daftar BUSINESS (L55-7)', () => {
  const DAFTAR = ['PROPERTY ALL RISK', 'MACHINERY BREAK DOWN', 'PUBLIC LIABILITY']
  const cari = (q: string) =>
    Promise.resolve({ baris: DAFTAR.filter((d) => d.includes(q.toUpperCase())).map((d, i) => ({ id: String(i), label: d })) })

  it('template: semua ada di daftar -> lolos, nomor baris berkas ikut', async () => {
    const h = await cocokkanKelasBisnis(bacaCsvLossRecord(TEMPLATE_CSV_LOSS_RECORD), cari)
    expect(h.galat).toEqual([])
    expect(h.baris).toHaveLength(9)
    expect(h.nomor).toEqual([2, 3, 4, 5, 6, 7, 8, 9, 10])
  })
  it('ejaan disamakan ke daftar; tidak ada / kosong -> galat per baris, grid tidak diganti', async () => {
    const kepala = 'CLASS_OF_BUSINESS;DATE_OF_LOSS;CAUSE_OF_LOSS;CURRENCY;LOSS_AMOUNT'
    const ok = await cocokkanKelasBisnis(bacaCsvLossRecord(`${kepala}\nproperty all risk ;;x;IDR;1`), cari)
    expect(ok.baris[0]?.classOfBusiness).toBe('PROPERTY ALL RISK')
    const salah = await cocokkanKelasBisnis(bacaCsvLossRecord(`${kepala}\nPROPERTY;;x;IDR;1\n;;y;IDR;2\nPUBLIC LIABILITY;;z;IDR;3`), cari)
    expect(salah.baris).toEqual([])
    expect(salah.galat).toEqual([`${T.csvBaris}2${T.csvKelasTidakAda}PROPERTY`, `${T.csvBaris}3${T.csvKelasKosong}`])
  })
})

describe('TabLossRecord - tampilan dan aturan', () => {
  const v: TampilanLossRecord = {
    premiPar: '9391000', premiMb: '200000', premiPl: '25000', premiOther: '',
    baris: [
      { uwYear: '2024', classOfBusiness: 'UJI-COB-1', dateOfLoss: '03-01-2025', causeOfLoss: 'UJI-SEBAB', currency: 'IDR', lossAmount: '100000000', lossAmountIdr: '100000000' },
    ],
    totalCob: [{ classOfBusiness: 'UJI-COB-1', currency: 'IDR', lossAmount: '100000000', lossAmountIdr: '100000000' }],
    rasio: { gross1: '1039.9334', gross3: '346.6444', gross5: '207.9866', net1: '', net3: '', net5: '' },
    pesan: [],
  }
  it('render: bar unggah, premi, grid, total, loss ratio', () => {
    const html = renderToStaticMarkup(<TabLossRecord caseId="NB-UJI" insuredName="UJI-TERTANGGUNG" awal={v} />)
    for (const t of [L.judulUnggah, L.unduhTemplate, L.unggah, L.simpan, L.simpanTab, L.premiPar, L.totalLossIdr, ...L.rasio]) expect(html).toContain(`>${t}<`)
    expect(html).toContain('>03/01/2025<')
    expect(html).toContain('>100.000.000,0000<')
    expect(html).toContain('>1.039,9334%<')
  })
  it('persen loss ratio 4 desimal; kosong tetap kosong', () => {
    expect(persenRasio('6239.6')).toBe('6.239,6000%')
    expect(persenRasio('')).toBe('')
  })
  it('Class Of Business wajib dari daftar: sama persis dipilih, selain itu ketikan dihapus (L55-7)', () => {
    const baris = [{ id: '1', label: 'PROPERTY ALL RISK' }, { id: '2', label: 'PROPERTY ALL RISK 2' }]
    expect(saranSama('property all risk ', baris)).toEqual({ id: '1', label: 'PROPERTY ALL RISK' })
    expect(saranSama('PROPERTY', baris)).toBeNull()
    expect(saranSama('', baris)).toBeNull()
    expect(SUMBER).toMatch(/onPilih=\{\(s\) => ubahBaris\(i, \{ classOfBusiness: s\.label \}\)\}\n\s+muatSaatFokus\n\s+wajibDariDaftar/)
    const SARAN = readFileSync(join(__dirname, 'IsianSaran.tsx'), 'utf8').replace(/\r\n/g, '\n')
    expect(SARAN).toContain("if (s) onPilih(s)\n                else onKetik('')")
  })
  it('akses: tab hanya bila backend mengizinkan (POSITION IT, tiket 55)', () => {
    expect(HALAMAN).toContain("(t !== 'Loss Record' || dariServer?.aksesLossRecord === true)")
    expect(HALAMAN).toContain("tab === 'Loss Record' && kasusFire && dariServer?.aksesLossRecord === true")
  })
  it('hitungan di backend; tanpa float', () => {
    expect(SUMBER).toContain('hitungLossRecord(caseId, baru)')
    expect(SUMBER).toContain('setV(await simpanLossRecord(caseId, v))')
    expect(SUMBER).not.toMatch(/(?<![A-Za-z.])Number\(|parseFloat/)
  })
})
