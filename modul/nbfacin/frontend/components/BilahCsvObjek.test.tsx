// Bilah CSV objek General (tiket 66) - data uji sintetis.

import { existsSync, readFileSync } from 'node:fs'
import { join } from 'node:path'
import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { PERIODE as P, TEKS_CSV_OBJEK as T } from '../labels'
import BilahCsvObjek, { pesanTugas, sedangBerjalan } from './BilahCsvObjek'
import { BATAS_BARIS_CSV_OBJEK, bacaCsvObjek, hitungBarisBerisi, KOLOM_CSV_OBJEK, periksaJudulBerkas, KOLOM_CSV_PEGA, PESAN_CSV_OBJEK, TEMPLATE_CSV_OBJEK } from './csvObjek'

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
  it('Save Data (tiket 75): hanya dari unggahan valid; case berobjek (adaObjek) -> konfirmasi ganti; selesai latar -> tab detail dimuat ulang', () => {
    expect(SUMBER).toContain("if (tugas === null || tugas.status !== 'valid') {")
    expect(SUMBER).toContain('if (tugas.adaObjek) setKonfirmasi(true)')
    expect(SUMBER).toContain('terima(await simpanTugasCsvObjek(caseId, tugas.idUnggah))')
    // onTersimpan sekali, saat status berpindah dari 'menyimpan' ke 'selesai'.
    expect(SUMBER).toContain("if (t !== null && t.status === 'selesai' && statusLama.current === 'menyimpan') onTersimpan()")
    const HALAMAN = readFileSync(join(__dirname, '../pages/InwardFacultative.tsx'), 'utf8').replace(/\r\n/g, '\n')
    expect(HALAMAN).toContain('onTersimpan={() => setVersiDetail((n) => n + 1)}')
  })
})

describe('hitungBarisBerisi - berkas besar tidak dibaca utuh (bug DEV 10-10-2026: 1 juta baris tampil "kosong")', () => {
  it('menghitung baris berisi; baris kosong / hanya pemisah tidak dihitung; berhenti begitu melewati batas', async () => {
    expect(await hitungBarisBerisi(new Blob(['A;B\r\n1;2\r\n\r\n;;\r\n3;4']), 100)).toBe(3)
    expect(await hitungBarisBerisi(new Blob(['']), 100)).toBe(0)
    // 10.000 baris, batas 5: berhenti di 6 (tidak menghitung sampai habis).
    expect(await hitungBarisBerisi(new Blob(['x\n'.repeat(10000)]), 5)).toBe(6)
  })

  it('tiket 75: berkas TIDAK dibaca utuh - judul dari potongan awal, lalu di-stream ke server (tanpa batas baris)', async () => {
    const s = readFileSync(join(__dirname, 'BilahCsvObjek.tsx'), 'utf8')
    expect(s).not.toContain('f.text()')
    expect(s).not.toContain('BATAS_BARIS_CSV_OBJEK')
    expect(s.indexOf('await periksaJudulBerkas(f)')).toBeLessThan(s.indexOf('await unggahCsvObjek(caseId, f)'))
    const judul = KOLOM_CSV_OBJEK.join(';')
    expect(await periksaJudulBerkas(new Blob([judul + '\n' + 'x;'.repeat(10) + '\n']))).toBe('')
    expect(await periksaJudulBerkas(new Blob(['SALAH;JUDUL\n1;2']))).toBe(PESAN_CSV_OBJEK.judul)
    expect(await periksaJudulBerkas(new Blob([]))).toBe(PESAN_CSV_OBJEK.kosong)
  })

  it('tiket 75: pesan status latar (memvalidasi / menyimpan / gagal + Ulangi)', () => {
    const t = { idUnggah: '1', totalBaris: 1000000, barisDiproses: 350000, barisTersimpan: 120000, galat: ['baris 7: x'], jumlahGalat: 3, pesan: [], namaBerkas: 'a.csv', dibuat: '', adaObjek: false }
    expect(pesanTugas({ ...t, status: 'memvalidasi' })[0]).toContain('350.000')
    expect(pesanTugas({ ...t, status: 'menyimpan' })[0]).toContain('120.000 / 1.000.000')
    expect(pesanTugas({ ...t, status: 'gagal_validasi' })).toEqual([T.gagalValidasi(3), 'baris 7: x', T.galatLain(2)])
    expect(sedangBerjalan({ ...t, status: 'menyimpan' })).toBe(true)
    expect(sedangBerjalan({ ...t, status: 'gagal' })).toBe(false)
    expect(readFileSync(join(__dirname, 'BilahCsvObjek.tsx'), 'utf8')).toContain("{tugas?.status === 'gagal' && (")
  })
})

describe('bilah CSV per case (bug DEV 10-10-2026: unggahan FIRE tampil di case MARINE CARGO)', () => {
  it('status dikosongkan saat case berganti; bilah dipasang dengan key caseId', () => {
    const s = readFileSync(join(__dirname, 'BilahCsvObjek.tsx'), 'utf8').replace(/\r\n/g, '\n')
    expect(s).toMatch(/statusLama\.current = null\s*setTugas\(null\)\s*setPesan\(\[\]\)\s*setGalat\(null\)\s*ambilTugasCsvAktif\(caseId\)/)
    const h = readFileSync(join(__dirname, '../pages/InwardFacultative.tsx'), 'utf8').replace(/\r\n/g, '\n')
    expect(h).toMatch(/<BilahCsvObjek[\s\S]*?key=\{kasus\.caseId\}\s*caseId=\{kasus\.caseId\}/)
  })
})

describe('CSV FIRE hanya untuk case FIRE (bug DEV 11-10-2026: CSV FIRE masuk ke case MARINE CARGO)', () => {
  it('bilah CSV FIRE tidak dipasang untuk case non-FIRE (Pega ChooseActivityMappingCSV_Act: InsertUploadFire_act hanya IsFire)', () => {
    const h = readFileSync(join(__dirname, '../pages/InwardFacultative.tsx'), 'utf8').replace(/\r\n/g, '\n')
    expect(h).toContain("const jenisCsv: JenisCsv | null = kasusFire ? 'fire' : kasusCargo ? 'cargo' : dariServer?.jenisCsv === 'marinehull' ? 'marinehull' : null")
    expect(h).toMatch(/\{!bacaSaja && jenisCsv === null && !tanpaCsvPega \? \(\s*<div className="alert alert--info">\{TEKS_CSV_OBJEK\.belumUntukCob\(op\.groupBusiness \?\? ''\)\}<\/div>\s*\) : !bacaSaja && jenisCsv !== null \? \(\s*<BilahCsvObjek/)
  })
})


describe('CSV Marine Cargo (tiket 74 P6, K74-14..17)', () => {
  it('kolom = InsertUploadMarineCargo_act (18); judul dicek per NAMA, urutan bebas; template Pega ditunggu', async () => {
    const { KOLOM_CSV_CARGO, judulCsvCargoSah, periksaJudulBerkasCargo, PESAN_CSV_CARGO } = await import('./csvCargo')
    // Judul template Pega VERBATIM (D:\migrasi\RNM\DDL\templatecsvmarinecargo.csv baris 1).
    expect(KOLOM_CSV_CARGO.join(';')).toBe(
      'TRADING;ETD;DEPARTURE;DESTINATION;INVOICE;INVOICEDATE;BLNUMBER;CONVEYENCE;CONVEYENCE_NAME;DATASHIP;GOODS;GOODS_NAME;PACKING;AGENT;COVERAGE_IDX;CURRENCY;TSI;TOP_RISK',
    )
    expect(judulCsvCargoSah([...KOLOM_CSV_CARGO].reverse().join(';'))).toBe(true)
    // INVOICEDATE tidak dibaca Pega -> boleh tidak ada; kolom wajib lain tidak.
    expect(judulCsvCargoSah(KOLOM_CSV_CARGO.filter((k) => k !== 'INVOICEDATE').join(';'))).toBe(true)
    expect(judulCsvCargoSah(KOLOM_CSV_CARGO.slice(1).join(';'))).toBe(false)
    expect(await periksaJudulBerkasCargo(new Blob(['trading;' + KOLOM_CSV_CARGO.slice(1).join(';') + ';EKSTRA\n1']))).toBe('')
    expect(await periksaJudulBerkasCargo(new Blob(['A;B\n1;2']))).toBe(PESAN_CSV_CARGO.judul)
    const { CSV_PER_JENIS } = await import('./BilahCsvObjek')
    expect(CSV_PER_JENIS.cargo.berkas).toBe('templatecsvmarinecargo.csv')
    expect(CSV_PER_JENIS.fire.berkas).toBe('templatecsvfire.csv')
    expect(CSV_PER_JENIS.cargo.periksa).toBe(periksaJudulBerkasCargo)
    const s = readFileSync(join(__dirname, 'BilahCsvObjek.tsx'), 'utf8')
    expect(s).toContain('a.download = CSV_PER_JENIS[jenis].berkas')
    expect(s).toContain('const salahJudul = await CSV_PER_JENIS[jenis].periksa(f)')
    const h = readFileSync(join(__dirname, '../pages/InwardFacultative.tsx'), 'utf8')
    expect(h).toContain('jenis={jenisCsv}')
  })
})

describe('COB tanpa CSV di Pega (work owner 11-10-2026: GA (OTHERS) memang tidak ada)', () => {
  it('jenisCsv server "" -> tanpa bilah CSV dan tanpa pesan "belum tersedia"; Insert Accumulation IT tetap', () => {
    const h = readFileSync(join(__dirname, '../pages/InwardFacultative.tsx'), 'utf8').replace(/\r\n/g, '\n')
    expect(h).toContain("const tanpaCsvPega = jenisCsv === null && dariServer?.jenisCsv === ''")
    expect(h).toContain('{!hanyaLihat && (!bacaSaja || akunIT) && (!tanpaCsvPega || akunIT) && (')
    expect(h).toContain('{!tanpaCsvPega && <strong>{P.judulCsv.label}</strong>}')
    expect(T.belumUntukCob('ANEKA')).not.toContain('FIRE')
  })
})

describe('CSV Marine Hull (tiket 73 P6, K73-60)', () => {
  const TEMPLATE_HULL = 'D:/migrasi/RNM/DDL/templatecsvmarinehull.csv'
  it('kolom = judul template Pega VERBATIM (31); semua wajib, urutan bebas; template unduhan = judul saja', async () => {
    const { KOLOM_CSV_MARINE_HULL, judulCsvMarineHullSah, periksaJudulBerkasMarineHull, PESAN_CSV_MARINE_HULL, TEMPLATE_CSV_MARINE_HULL } =
      await import('./csvMarineHull')
    expect(KOLOM_CSV_MARINE_HULL).toHaveLength(31)
    if (existsSync(TEMPLATE_HULL)) {
      expect(KOLOM_CSV_MARINE_HULL.join(';')).toBe(readFileSync(TEMPLATE_HULL, 'utf8').split(/\r?\n/)[0])
    }
    expect(TEMPLATE_CSV_MARINE_HULL).toBe(KOLOM_CSV_MARINE_HULL.join(';') + '\r\n')
    expect(judulCsvMarineHullSah([...KOLOM_CSV_MARINE_HULL].reverse().join(';'))).toBe(true)
    expect(judulCsvMarineHullSah(KOLOM_CSV_MARINE_HULL.filter((k) => k !== 'LOL').join(';'))).toBe(false)
    expect(await periksaJudulBerkasMarineHull(new Blob(['no;' + KOLOM_CSV_MARINE_HULL.slice(1).join(';') + ';EKSTRA\n1']))).toBe('')
    expect(await periksaJudulBerkasMarineHull(new Blob(['A;B\n1;2']))).toBe(PESAN_CSV_MARINE_HULL.judul)
    const { CSV_PER_JENIS } = await import('./BilahCsvObjek')
    expect(CSV_PER_JENIS.marinehull.berkas).toBe('templatecsvmarinehull.csv')
    expect(CSV_PER_JENIS.marinehull.kolom).toBe(KOLOM_CSV_MARINE_HULL)
    expect(CSV_PER_JENIS.marinehull.periksa).toBe(periksaJudulBerkasMarineHull)
  })
})
