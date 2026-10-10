// Bilah CSV objek General (tiket 66) - data uji sintetis.

import { existsSync, readFileSync } from 'node:fs'
import { join } from 'node:path'
import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { PERIODE as P, TEKS_CSV_OBJEK as T } from '../labels'
import BilahCsvObjek from './BilahCsvObjek'
import { BATAS_BARIS_CSV_OBJEK, bacaCsvObjek, KOLOM_CSV_OBJEK, KOLOM_CSV_PEGA, PESAN_CSV_OBJEK, TEMPLATE_CSV_OBJEK } from './csvObjek'

const TEMPLATE_DDL = 'D:\\migrasi\\RNM\\DDL\\templatecsvfire.csv'
const AKT = 'D:\\migrasi\\RNM\\NB FacIn\\Activity\\InsertUploadFire_act.xml'
const PERIODE_XML = 'D:\\migrasi\\RNM\\NB FacIn\\Section\\Periode.xml'
const SUMBER = readFileSync(join(__dirname, 'BilahCsvObjek.tsx'), 'utf8').replace(/\r\n/g, '\n')
const JUDUL = KOLOM_CSV_OBJEK.join(';')
const baris = (isi: Partial<Record<(typeof KOLOM_CSV_OBJEK)[number], string>>) => KOLOM_CSV_OBJEK.map((k) => isi[k] ?? '').join(';')

describe('csvObjek - template FIRE', () => {
  it.skipIf(!existsSync(TEMPLATE_DDL))('kolom = baris judul templatecsvfire.csv (pemisah ;)', () => {
    const judul = readFileSync(TEMPLATE_DDL, 'utf8').split(/\r?\n/)[0] ?? ''
    expect(judul.split(';').filter((x) => x !== '')).toEqual([...KOLOM_CSV_OBJEK])
  })
  it.skipIf(!existsSync(AKT))('pesan Please Upload CSV verbatim InsertUploadFire_act', () => {
    expect(readFileSync(AKT, 'utf8')).toContain(`<PropertiesValue>"${T.pleaseUpload}"</PropertiesValue>`)
  })
  it.skipIf(!existsSync(PERIODE_XML))('label tombol = Periode.xml sel 93-96', () => {
    const xml = readFileSync(PERIODE_XML, 'utf8')
    for (const t of [P.downloadTemplateCsv, P.uploadCsv, P.viewUpload, P.saveData]) expect(xml).toContain(`>${t.label}<`)
  })
  it('template unduhan = baris judul saja (baris contoh Pega berisi nama, tidak disalin)', () => {
    expect(TEMPLATE_CSV_OBJEK).toBe(JUDUL + '\r\n')
    expect(KOLOM_CSV_OBJEK).toHaveLength(61)
    expect(KOLOM_CSV_PEGA).toBe(55)
  })
  it('batas 3000 baris per unggahan (work owner 08-10-2026)', () => {
    expect(BATAS_BARIS_CSV_OBJEK).toBe(3000)
    const isi = (n: number) => JUDUL + '\n' + Array.from({ length: n }, (_, i) => baris({ NO: String(i + 1) })).join('\n')
    expect(bacaCsvObjek(isi(3000)).baris).toHaveLength(3000)
    expect(bacaCsvObjek(isi(3001)).galat).toBe(PESAN_CSV_OBJEK.terlaluBanyak(3001))
  })
  it('judul 55 kolom (template Pega lama) tetap diterima; 6 kolom wilayah kosong', () => {
    const lama = KOLOM_CSV_OBJEK.slice(0, KOLOM_CSV_PEGA)
    const h = bacaCsvObjek(lama.join(';') + '\r\n' + lama.map((k) => (k === 'NO' ? '1' : '')).join(';'))
    expect(h.galat).toBe('')
    expect(h.baris[0]?.KODE_POS).toBe('')
    expect(bacaCsvObjek(KOLOM_CSV_OBJEK.slice(0, 50).join(';') + '\n1').galat).toBe(PESAN_CSV_OBJEK.judul)
  })
  it('bacaCsvObjek: judul wajib sama, baris kosong dilewati, sel apa adanya (dipangkas)', () => {
    const h = bacaCsvObjek(`\uFEFF${JUDUL};;\r\n${baris({ NO: '1', ADDRESS: ' UJI-ALAMAT ', TSI_OBJECT_ITEM: '1.000.000,50' })}\r\n;;;\r\n`)
    expect(h.galat).toBe('')
    expect(h.baris).toHaveLength(1)
    expect(h.baris[0]?.ADDRESS).toBe('UJI-ALAMAT')
    expect(h.baris[0]?.TSI_OBJECT_ITEM).toBe('1.000.000,50')
    expect(h.baris[0]?.RATE_OTHERS).toBe('')
    expect(bacaCsvObjek('').galat).toBe(PESAN_CSV_OBJEK.kosong)
    expect(bacaCsvObjek('NO;ADDRESS\r\n1;X').galat).toBe(PESAN_CSV_OBJEK.judul)
    expect(bacaCsvObjek(JUDUL + '\r\n').galat).toBe(PESAN_CSV_OBJEK.tanpaBaris)
    expect(bacaCsvObjek(JUDUL.toLowerCase() + '\n' + baris({ NO: '1' })).galat).toBe('')
  })
})

describe('BilahCsvObjek', () => {
  it('empat tombol aktif + tombol akhir', () => {
    const html = renderToStaticMarkup(<BilahCsvObjek caseId="NB-UJI" akhir={<span>UJI-AKHIR</span>} onTersimpan={() => {}} />)
    for (const t of [P.downloadTemplateCsv, P.uploadCsv, P.viewUpload, P.saveData]) {
      expect(html).toMatch(new RegExp(`<button type="button" class="[^"]*">${t.label.replace(/[/.]/g, '\\$&')}</button>`))
    }
    expect(html).toContain('UJI-AKHIR')
    expect(html).toContain('type="file"')
  })
  it('Save Data: tanpa unggahan -> "Please Upload CSV"; case berobjek -> konfirmasi ganti; berhasil -> tab detail dimuat ulang', () => {
    expect(SUMBER).toMatch(/if \(unggahan === null\) \{\s*setPesan\(\[T\.pleaseUpload\]\)/)
    expect(SUMBER).toContain('const ada = (await ambilObjek(caseId)).baris.length > 0')
    expect(SUMBER).toContain('if (ada) setKonfirmasi(true)')
    expect(SUMBER).toContain('const h = await simpanCsvObjek(caseId, unggahan)')
    expect(SUMBER).toContain('onTersimpan()')
    const HALAMAN = readFileSync(join(__dirname, '../pages/InwardFacultative.tsx'), 'utf8').replace(/\r\n/g, '\n')
    expect(HALAMAN).toContain('onTersimpan={() => setVersiDetail((n) => n + 1)}')
  })
})
