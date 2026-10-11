// Blok Email Type + Submit tahap marketing + popup Cause of Decline (tiket 57) - data uji sintetis.

import { existsSync, readFileSync } from 'node:fs'
import { join } from 'node:path'
import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import type { PilihanSebabTc, TampilanSubmit } from '../api'
import { EMAIL_TYPE as E, EMAIL_TYPE_BINDING as EB, EMAIL_TYPE_TL as ETL, EMAIL_TYPE_UW as EUW, SEBAB_DECLINE as S, TEKS_EMAIL_TYPE as T, TEKS_SEBAB_DECLINE as TS } from '../labels'
import { ApiFailure } from '../../../../inti/frontend/klien'
import PortalOpportunity from '../pages/PortalOpportunity'
import { navigasi, PILIH_UBAH_AKUMULASI, terapkanBacaSaja } from './BacaSajaKasus'
import BlokEmailType, {
  diproteksi,
  galatIsianSubmit,
  isianSubmitKosong,
  opsiEmailType,
  pesanProteksiDari,
  potongRiwayat,
  tanggalRiwayat,
  UKURAN_RIWAYAT,
  type SubmitMarketing,
} from './BlokEmailType'
import { centangPrioritas, komentarDecline, optionAktif, pctRnmAktif, pctSobAktif, pilihanKosong, remarksAktif, teksSebabTc, validasiSebab } from './sebabDecline'

const NB = 'D:\\migrasi\\RNM\\NB FacIn\\'
const DDL_TL = 'D:\\migrasi\\RNM\\DDL\\EmailTypeTL.xml'
const HALAMAN = readFileSync(join(__dirname, '../pages/InwardFacultative.tsx'), 'utf8').replace(/\r\n/g, '\n')
const ada = existsSync(NB + 'Section\\EmailSection.xml')
const baca = (f: string) => readFileSync(NB + f, 'utf-8')

describe.skipIf(!ada)('label dan aturan = korpus', () => {
  it('EmailSection: label, grid riwayat; konfirmasi RejectNotificationSection', () => {
    const xml = baca('Section\\EmailSection.xml')
    for (const t of [E.label, E.comment, E.bandingTo, E.pilihSebab, E.bandingKosong, ...E.kolomRiwayat]) expect(xml, t).toContain(`>${t}<`)
    expect(xml).toContain('<pyRequiredWhen>.EmailType = 4</pyRequiredWhen>')
    expect(xml).toContain('<pyDisabledWhen>.EmailType = 7</pyDisabledWhen>')
    const rn = baca('Section\\RejectNotificationSection.xml')
    for (const t of [E.yakinReject, E.yakinDecline]) expect(rn, t).toContain(`>${t}<`)
  })
  it('kode Email Type -> label (SaveViewSuggest CARI9)', () => {
    const xml = baca('Activity\\SaveViewSuggest.xml').replace(/&quot;/g, '"')
    for (const o of E.opsi) expect(xml, o.label).toContain(`.Approval="${o.kode}","${o.label}"`)
  })
  it('popup Cause of Decline: label (InputCauseOfDecline), pesan verbatim (CekCauseOfDecline_ACT), sumber daftar', () => {
    const xml = baca('Section\\InputCauseOfDecline.xml')
    for (const t of [S.judul, S.nonTc, ...S.kolom]) expect(xml, t).toContain(`>${t}<`)
    for (const t of [S.cek, S.submit]) expect(xml, t).toContain(`<pyLabel>${t}</pyLabel>`)
    const cek = baca('Activity\\CekCauseOfDecline_ACT.xml').replace(/&quot;/g, '"')
    for (const p of Object.values(S.pesan)) expect(cek, p).toContain(`"${p}"`)
    expect(baca('RDBList\\GetCauseOfDecline.xml')).toContain('select TERMCONDITION from POOLDATA.CAUSEOFDECLINE WHERE DESCRIPTION IS NULL')
    expect(baca('RDBList\\GetCauseOfDecline2.xml')).toContain('select TERMCONDITION from POOLDATA.CAUSEOFDECLINE WHERE TERMCONDITION2 IS NOT NULL')
  })
  it('syarat kolom = pyDisabledWhen InputCauseOfDecline (nama T/C verbatim)', () => {
    const xml = baca('Section\\InputCauseOfDecline.xml').replace(/&quot;/g, '"').replace(/&amp;/g, '&')
    for (const t of [...S.term.option, ...S.term.pctRnm, ...S.term.pctSob]) expect(xml, t).toContain(`"${t}"`)
    for (const t of S.term.remarks) expect(xml, t).toContain(`'${t}'`)
  })
})

describe('Cause of Decline - aturan', () => {
  const daftar = ['Rate', 'Occupancy', 'Period Of Insurance (POI)', 'Loss Ratio (net)'].map(pilihanKosong)
  it('kolom aktif per T/C', () => {
    expect([optionAktif('Risk Exposure'), optionAktif('Rate')]).toEqual([true, false])
    expect([remarksAktif('Clauses', ''), remarksAktif('Extension of coverage', 'Others'), remarksAktif('Rate', '')]).toEqual([true, true, false])
    expect([pctRnmAktif('Loss Ratio (net)'), pctSobAktif('Loss Ratio (net)'), pctSobAktif('OJK Discount')]).toEqual([true, false, true])
  })
  it('satu baris per prioritas; klik lagi = lepas; % RNM otomatis prioritas I', () => {
    let d = centangPrioritas(daftar, 0, 'I', '1.2345')
    expect(d[0]).toMatchObject({ prioritas: 'I', pctRnm: '1.2345' })
    d = centangPrioritas(d, 1, 'I')
    expect([d[0]?.prioritas, d[1]?.prioritas]).toEqual(['', 'I'])
    d = centangPrioritas(d, 1, 'II')
    expect(d[1]?.prioritas).toBe('II')
    expect(centangPrioritas(d, 1, 'II')[1]).toEqual(pilihanKosong('Occupancy'))
  })
  it('validasi verbatim; wajib prioritas I atau Non - T/C', () => {
    const p = (isi: Partial<PilihanSebabTc>): PilihanSebabTc => ({ ...pilihanKosong('Rate'), prioritas: 'I', ...isi })
    expect(validasiSebab([p({ pctRnm: '1' })], '')).toEqual([S.pesan.pctRnmSob])
    expect(validasiSebab([p({ termCondition: 'Loss Ratio (net)' })], '')).toEqual([S.pesan.pctRnm])
    expect(validasiSebab([p({ termCondition: 'Occupancy' })], '')).toEqual([S.pesan.remaks])
    expect(validasiSebab([p({ termCondition: 'Risk Exposure', description1: 'Others' })], '')).toEqual([S.pesan.remarks])
    expect(validasiSebab([p({ termCondition: 'Period Of Insurance (POI)' })], '')).toEqual([S.pesan.option])
    expect(validasiSebab([], '')).toEqual([TS.kosong])
    expect(validasiSebab([], 'UJI-NON-TC')).toEqual([])
  })
  it('teks alasan + Comment = Reason1 \\n Reason2; Non - T/C menimpa Reason1', () => {
    const a: PilihanSebabTc = { termCondition: 'Rate', prioritas: 'I', description1: '', description2: '', pctRnm: '1.5', pctSob: '2' }
    const b: PilihanSebabTc = { termCondition: 'Occupancy', prioritas: 'II', description1: '', description2: 'UJI', pctRnm: '', pctSob: '' }
    expect(teksSebabTc(a)).toBe('Rate   1.5 % RNM 2 % SOB ')
    expect(komentarDecline([a, b], '')).toBe('Rate   1.5 % RNM 2 % SOB \nOccupancy  UJI ')
    expect(komentarDecline([a], 'UJI-NON-TC')).toBe('UJI-NON-TC\n')
  })
})

describe('Blok Email Type - isian dan tampilan', () => {
  it('kelengkapan isian sebelum Submit', () => {
    const k = isianSubmitKosong()
    expect(galatIsianSubmit(k)).toBe(T.wajibType)
    expect(galatIsianSubmit({ ...k, emailType: '1' })).toBe(T.wajibComment)
    expect(galatIsianSubmit({ ...k, emailType: '1', comment: 'ok' })).toBe('')
    expect(galatIsianSubmit({ ...k, emailType: '4', comment: 'ok' })).toBe(T.wajibBanding)
    expect(galatIsianSubmit({ ...k, emailType: '7', comment: 'x' })).toBe(T.wajibSebab)
    expect(galatIsianSubmit({ ...k, emailType: '7', comment: 'x', sebabNonTc: 'UJI' })).toBe('')
  })
  it('tanggal riwayat DD/MM/YYYY HH:MM', () => {
    expect(tanggalRiwayat('07-10-2026 09:05')).toBe('07/10/2026 09:05')
  })
  it('render: 4 opsi, Comment, riwayat; Banding tanpa opsi -> catatan', () => {
    const v: TampilanSubmit = {
      bolehSubmit: true, alasan: '', opsiBanding: [], sebabTc: [], sebabNonTc: [], pctRnmOtomatis: { rate: '', riComm: '' },
      tahap: 'marketing', deductibleBedaMataUang: false, deductibleAcceptance: false, commentDeductible: '', daysValidity: '', dateValidity: '', waitingBindDate: '', confirmBinding: false, receivedRiSlip: false,
      riwayat: [{ tanggal: '07-10-2026 09:05', pic: 'UJI-PIC', approval: 'Accept', status: 'Offer', comment: 'UJI-KOMENTAR' }],
    }
    const s = { v, isi: { ...isianSubmitKosong(), emailType: '4' }, proteksi: ['UJI-PROTEKSI'], galat: null, mengajukan: false, memeriksa: false, berhasil: null, konfirmasi: false } as unknown as SubmitMarketing
    const html = renderToStaticMarkup(<BlokEmailType s={s} />)
    for (const o of E.opsi) expect(html).toContain(`${o.label}</label>`)
    expect(html).toContain('>07/10/2026 09:05<')
    expect(html).toContain('>UJI-KOMENTAR<')
    expect(html).toContain(T.bandingTanpaOpsi)
    expect(html).toContain('>UJI-PROTEKSI<')
  })
  it('case sudah tidak di tahap marketing: form disembunyikan TANPA keterangan (work owner 08-10-2026), riwayat tetap', () => {
    const v: TampilanSubmit = {
      bolehSubmit: false, alasan: 'UJI-ALASAN', opsiBanding: [], sebabTc: [], sebabNonTc: [], pctRnmOtomatis: { rate: '', riComm: '' },
      tahap: 'marketing', deductibleBedaMataUang: false, deductibleAcceptance: false, commentDeductible: '', daysValidity: '', dateValidity: '', waitingBindDate: '', confirmBinding: false, receivedRiSlip: false,
      riwayat: [{ tanggal: '07-10-2026 09:25', pic: 'UJI-PIC', approval: 'Accept', status: 'Offer', comment: 'UJI-OK' }],
    }
    const s = { v, isi: isianSubmitKosong(), proteksi: [], galat: null, mengajukan: false, memeriksa: false, berhasil: null, konfirmasi: false } as unknown as SubmitMarketing
    const html = renderToStaticMarkup(<BlokEmailType s={s} />)
    expect(html).not.toContain('UJI-ALASAN')
    expect(html).not.toContain('nbf-et__info')
    expect(html).not.toContain('type="radio"')
    expect(html).not.toContain('<textarea')
    expect(html).toContain('>UJI-OK<')
  })
  it('halaman: blok FIRE di bawah detail, Submit aktif hanya bila boleh', () => {
    // Tiket 73 P4: FIRE dan keluarga Aneka.
    expect(HALAMAN).toContain('const bolehSubmit = kasusFire || kasusAneka')
    expect(HALAMAN).toContain('{bolehSubmit && <BlokEmailType s={sub} hanyaRiwayat={hanyaLihat} />}')
    expect(HALAMAN).toContain('disabled={!bolehSubmit || !sub.bisaSubmit}')
  })
})

describe('penolakan proteksi 422', () => {
  it('pesan dipulihkan dari galat "; "; galat lain bukan proteksi', () => {
    expect(pesanProteksiDari(new ApiFailure(422, { code: 'DITOLAK_BACKEND', message: 'Scoring Risk Cannot Empty !!; List Payment can\'t be null!' }))).toEqual([
      'Scoring Risk Cannot Empty !!',
      "List Payment can't be null!",
    ])
    expect(pesanProteksiDari(new ApiFailure(409, { code: 'DITOLAK_BACKEND', message: 'x' }))).toBeNull()
    expect(pesanProteksiDari(new Error('x'))).toBeNull()
  })
})

describe('Submit berhasil -> case ditutup, kembali ke portal (work owner 07-10-2026)', () => {
  it('halaman memanggil onSelesai; rute membuka portal dengan pesan', () => {
    expect(HALAMAN).toContain('if (onSelesai) onSelesai(pesan)')
    const RUTE = readFileSync(join(__dirname, '../rute.tsx'), 'utf8').replace(/\r\n/g, '\n')
    expect(RUTE).toContain('setPesanPortal(pesan)')
    expect(RUTE).toContain('setKasus(null)')
    // Dibuka dari kotak masuk Beranda: Submit kembali ke Beranda (work owner 07-10-2026), sesudah jendela pesan NB Status
    // ditutup OK (work owner 08-10-2026: "saat di aksep uw, tadi kenapa tidak ada info nb statusnya ya").
    expect(RUTE).toMatch(/if \(dariBeranda && onBeranda\) \{\s*setPesanBeranda\(pesan\)\s*return/)
    expect(RUTE).toMatch(/labelBatal=\{TEKS_EMAIL_TYPE\.ok\}\s*onTutup=\{\(\) => \{\s*setPesanBeranda\(''\)\s*setDariBeranda\(false\)\s*onBeranda\?\.\(\)/)
    expect(RUTE).toContain('<p>{pesanBeranda}</p>')
    // Sesudah login ulang: permintaan buka case lama tidak dijalankan, halaman case tanpa case kembali ke Beranda
    // (work owner 08-10-2026: "baru login langsung kebuka gini, harusnya kan ada di menu home terlebih dahulu").
    expect(RUTE).toContain('const ketukSaatPasang = useRef(ketukBuka)')
    expect(RUTE).toContain('ketukBuka !== ketukSaatPasang.current')
    expect(RUTE).toMatch(/if \(halaman === 'nbfacin-inward' && kasus === null && pesanBeranda === ''\) \{\s*if \(onBeranda\) onBeranda\(\)/)
    const portal = renderToStaticMarkup(<PortalOpportunity onBuat={() => {}} onBuka={() => {}} pesan="UJI-PESAN" onTutupPesan={() => {}} />)
    expect(portal).toContain('>UJI-PESAN<')
  })
})

describe('baca-saja case di luar tahap marketing (work owner 07-10-2026)', () => {
  const el = (o: { aria?: boolean; role?: string; nav?: boolean; pager?: boolean }) => ({
    hasAttribute: (n: string) => n === 'aria-expanded' && o.aria === true,
    getAttribute: (n: string) => (n === 'role' ? (o.role ?? null) : null),
    closest: (s: string) => ((s === '[data-navigasi]' && o.nav) || (s === '.pager' && o.pager) ? ({} as Element) : null),
  })
  it('navigasi tetap aktif: buka-baris, tab, paginasi, data-navigasi; lainnya dinonaktifkan', () => {
    expect([navigasi(el({ aria: true })), navigasi(el({ role: 'tab' })), navigasi(el({ pager: true })), navigasi(el({ nav: true }))]).toEqual([true, true, true, true])
    expect(navigasi(el({}))).toBe(false)
  })
  it('terapkanBacaSaja: elemen dalam selektor kecuali tetap aktif (tiket 60)', () => {
    type Palsu = { disabled: boolean; hasAttribute: (n: string) => boolean; getAttribute: (n: string) => string | null; closest: (s: string) => unknown }
    const buat = (akumulasi: boolean): Palsu => ({
      disabled: false,
      hasAttribute: () => false,
      getAttribute: () => null,
      closest: (s: string) => (s === PILIH_UBAH_AKUMULASI && akumulasi ? {} : null),
    })
    const a = buat(true), b = buat(false)
    const akar = { querySelectorAll: () => ({ forEach: (fn: (e: Palsu) => void) => [a, b].forEach(fn) }) } as unknown as ParentNode
    expect(terapkanBacaSaja(akar, PILIH_UBAH_AKUMULASI)).toBe(1)
    expect([a.disabled, b.disabled]).toEqual([false, true])
    const c = buat(true)
    const akar2 = { querySelectorAll: () => ({ forEach: (fn: (e: Palsu) => void) => [c].forEach(fn) }) } as unknown as ParentNode
    expect(terapkanBacaSaja(akar2)).toBe(1)
    expect(c.disabled).toBe(true)
  })
  it('terapkanBacaSaja sembunyi (mode hanya lihat, work owner 09-10-2026): tombol aksi disembunyikan, isian dinonaktifkan, navigasi tetap', () => {
    type Palsu = { tagName: string; hidden: boolean; disabled: boolean; hasAttribute: (n: string) => boolean; getAttribute: () => null; closest: (s: string) => unknown }
    const buat = (tagName: string, nav = false): Palsu => ({
      tagName,
      hidden: false,
      disabled: false,
      hasAttribute: () => false,
      getAttribute: () => null,
      closest: (s: string) => (s === '[data-navigasi]' && nav ? {} : null),
    })
    const tombol = buat('BUTTON'), isian = buat('INPUT'), lihat = buat('BUTTON', true)
    const akar = { querySelectorAll: () => ({ forEach: (fn: (e: Palsu) => void) => [tombol, isian, lihat].forEach(fn) }) } as unknown as ParentNode
    terapkanBacaSaja(akar, undefined, true)
    expect([tombol.hidden, isian.hidden, isian.disabled, lihat.hidden, lihat.disabled]).toEqual([true, false, true, false, false])
    expect(HALAMAN).toContain('{bacaSaja && !hanyaLihat && (')
  })
  it('halaman: General dan detail dibungkus BacaSajaKasus; Save for later nonaktif; banner', () => {
    expect(HALAMAN).toContain('const bacaSaja = hanyaLihat || dariServer?.bolehUbah === false')
    // General: bacaSaja; detail: detailBacaSaja (tahap UW membuka tab Spreading saja - tiket 59).
    expect(HALAMAN.match(/<BacaSajaKasus aktif=\{bacaSaja\} sembunyi=\{hanyaLihat\}>/g)).toHaveLength(1)
    expect(HALAMAN).toContain("<BacaSajaKasus key={`${detailBacaSaja ? `bacasaja${kecualiDetail ?? ''}` : 'ubah'}-${versiDetail}`} aktif={detailBacaSaja} kecuali={kecualiDetail} sembunyi={hanyaLihat}>")
    expect(HALAMAN).toContain('disabled={menyimpan || tanggalSalah || bacaSaja}')
  })
})

describe('tahap Team Leader (tiket 58)', () => {
  it('opsi Email Type TL = properti EmailTypeTL (DDL) - Accept / Reject / Ask / Decline', () => {
    expect(opsiEmailType('teamLeader').map((o) => `${o.kode}:${o.label}`)).toEqual(['1:Accept', '2:Reject', '3:Ask', '7:Decline'])
    expect(opsiEmailType('marketing').map((o) => o.kode)).toEqual(['1', '2', '4', '7'])
  })
  it.skipIf(!existsSync(DDL_TL))('opsi = PromptList EmailTypeTL.xml', () => {
    const xml = readFileSync(DDL_TL, 'utf-8')
    for (const o of ETL.opsi) expect(xml).toMatch(new RegExp(`<pyLocalizedValue>${o.label}</pyLocalizedValue>\\s*<pyStandardValue>${o.kode}</pyStandardValue>`))
  })
  it.skipIf(!ada)('label Deductible Acceptance = Periode.xml', () => {
    const xml = baca('Section\\Periode.xml')
    for (const t of [ETL.deductibleAcceptance, ETL.commentDeductible]) expect(xml, t).toContain(`>${t}<`)
  })
  it('render TL: 4 opsi TL, isian Deductible bila beda mata uang', () => {
    const v: TampilanSubmit = {
      bolehSubmit: true, alasan: '', opsiBanding: [], sebabTc: [], sebabNonTc: [], pctRnmOtomatis: { rate: '', riComm: '' }, riwayat: [],
      tahap: 'teamLeader', deductibleBedaMataUang: true, deductibleAcceptance: true, commentDeductible: 'UJI-KOMENTAR-DED', daysValidity: '', dateValidity: '', waitingBindDate: '', confirmBinding: false, receivedRiSlip: false,
    }
    const s = { v, isi: { ...isianSubmitKosong(), deductibleAcceptance: true, commentDeductible: 'UJI-KOMENTAR-DED' }, proteksi: [], galat: null, mengajukan: false, memeriksa: false, berhasil: null, konfirmasi: false } as unknown as SubmitMarketing
    const html = renderToStaticMarkup(<BlokEmailType s={s} />)
    expect(html).toContain('Ask</label>')
    expect(html).not.toContain('Banding</label>')
    expect(html).toContain(ETL.deductibleAcceptance)
    expect(html).toContain('>UJI-KOMENTAR-DED<')
  })
})

describe('riwayat akseptasi berhalaman 6 baris (work owner 08-10-2026)', () => {
  const r = (n: number) => ({ tanggal: '08-10-2026 15:0' + (n % 10), pic: `UJI-${n}`, approval: 'Accept', status: 'Offer', comment: 'OK' })
  it('potongRiwayat: 6 per halaman, halaman dijepit', () => {
    const baris = Array.from({ length: 10 }, (_, i) => r(i + 1))
    expect(UKURAN_RIWAYAT).toBe(6)
    expect(potongRiwayat(baris, 1).baris.map((x) => x.pic)).toEqual(['UJI-1', 'UJI-2', 'UJI-3', 'UJI-4', 'UJI-5', 'UJI-6'])
    expect(potongRiwayat(baris, 2).baris.map((x) => x.pic)).toEqual(['UJI-7', 'UJI-8', 'UJI-9', 'UJI-10'])
    expect(potongRiwayat(baris, 9).kini).toBe(2)
    expect(potongRiwayat([], 3)).toEqual({ baris: [], kini: 1 })
  })
  it('render: 6 baris + pager bila lebih dari 6; tanpa pager bila <= 6', () => {
    const buat = (n: number) => {
      const v: TampilanSubmit = {
        bolehSubmit: false, alasan: '', opsiBanding: [], sebabTc: [], sebabNonTc: [], pctRnmOtomatis: { rate: '', riComm: '' },
        riwayat: Array.from({ length: n }, (_, i) => r(i + 1)), tahap: '', deductibleBedaMataUang: false, deductibleAcceptance: false,
        commentDeductible: '', daysValidity: '', dateValidity: '', waitingBindDate: '', confirmBinding: false, receivedRiSlip: false,
      }
      const s = { v, isi: isianSubmitKosong(), proteksi: [], galat: null, mengajukan: false, memeriksa: false, berhasil: null, konfirmasi: false } as unknown as SubmitMarketing
      return renderToStaticMarkup(<BlokEmailType s={s} />)
    }
    const banyak = buat(10)
    expect(banyak).toContain('>UJI-6<')
    expect(banyak).not.toContain('>UJI-7<')
    expect(banyak).toContain('class="pager"')
    expect(buat(6)).not.toContain('class="pager"')
    // Pager di atas tabel riwayat.
    expect(banyak.indexOf('class="pager"')).toBeLessThan(banyak.indexOf('nbf-et__riwayat'))
  })
})

describe('pilihan Email Type tidak berganti sendiri (work owner 08-10-2026)', () => {
  const SUMBER = readFileSync(join(__dirname, 'BlokEmailType.tsx'), 'utf8').replace(/\r\n/g, '\n')
  const CSS = readFileSync(join(__dirname, '../nbfacin.css'), 'utf8').replace(/\r\n/g, '\n')
  it('radio melepas fokus sesudah dipilih (tombol panah tidak memindahkan pilihan)', () => {
    expect(SUMBER).toMatch(/onChange=\{\(e\) => \{[\s\S]{0,400}e\.currentTarget\.blur\(\)\s*void s\.pilihType\(o\.kode\)/)
  })
  it('area klik label pilihan selebar teksnya, bukan selebar kolom', () => {
    expect(CSS).toMatch(/\.nbfacin \.nbf-et__opsi label \{[\s\S]{0,200}display: inline-flex;\s*justify-self: start;/)
  })
})

describe('jalur Policy (tiket 69)', () => {
  it('Submit yang menyelesaikan Policy menampilkan nomor polis (69a-2, work owner 09-10-2026)', () => {
    expect(T.nomorPolis('RNM-F01.10.2026.00001')).toBe('Policy Number: RNM-F01.10.2026.00001')
    expect(HALAMAN).toContain("h.nomorPolis ? TEKS_EMAIL_TYPE.nomorPolis(h.nomorPolis) : ''")
  })
  const DDL_UWP = 'D:\\migrasi\\RNM\\DDL\\EmailTypeUWPolicy.xml'
  it.skipIf(!existsSync(DDL_UWP))('UW Policy = EmailTypeUWPolicy.xml tanpa Banding (K69-6); TL Policy = opsi TL empat (K69-2)', () => {
    const xml = readFileSync(DDL_UWP, 'utf-8')
    expect(xml).toContain('ASM-FW-GISFW-WORK!EMAILTYPEUWPOLICY')
    for (const o of opsiEmailType('uwPolicy')) expect(xml).toMatch(new RegExp(`<pyLocalizedValue>${o.label}</pyLocalizedValue>\\s*<pyStandardValue>${o.kode}</pyStandardValue>`))
    expect(xml).toContain('<pyLocalizedValue>Banding</pyLocalizedValue>')
  })
  it('opsi & isian: TL Policy = TL Offer; UW Policy Accept / Decline, Decline tanpa sebab, Validity seperti UW Offer', () => {
    expect(opsiEmailType('uwPolicy').map((o) => `${o.kode}:${o.label}`)).toEqual(['1:Accept', '7:Decline'])
    expect(opsiEmailType('teamLeaderPolicy')).toEqual(opsiEmailType('teamLeader'))
    const isi = { ...isianSubmitKosong(), emailType: '7' as const, comment: 'UJI' }
    expect(galatIsianSubmit(isi, 'uwPolicy')).toBe('')
    expect(galatIsianSubmit(isi, 'teamLeaderPolicy')).toBe(T.wajibSebab)
    // Validity Of Acceptance tampil juga di flag 1 (EmailSection_IsUW L9297-L10935 tanpa syarat; verifikasi sesi c3).
    expect(galatIsianSubmit({ ...isi, emailType: '1', daysValidity: '7x' }, 'uwPolicy')).toBe(T.hariAngka)
  })
})

describe('tahap Underwriting (tiket 61)', () => {
  const SUMBER = readFileSync(join(__dirname, 'BlokEmailType.tsx'), 'utf8').replace(/\r\n/g, '\n')
  const DDL_UW = 'D:\\migrasi\\RNM\\DDL\\EmailTypeUW.xml'
  const SEKSI_UW = 'D:\\migrasi\\RNM\\NB FacIn\\Section\\EmailSection_IsUW.xml'
  it.skipIf(!existsSync(DDL_UW))('opsi = PromptList EmailTypeUW.xml (Accept 1 / Reject 2 / Ask 3 / Decline 7)', () => {
    const xml = readFileSync(DDL_UW, 'utf-8')
    expect(xml).toContain('ASM-FW-GISFW-WORK!EMAILTYPEUW')
    for (const o of EUW.opsi) expect(xml).toMatch(new RegExp(`<pyLocalizedValue>${o.label}</pyLocalizedValue>\\s*<pyStandardValue>${o.kode}</pyStandardValue>`))
  })
  it.skipIf(!existsSync(SEKSI_UW))('label Validity / days, or up to / Waiting Bind = EmailSection_IsUW; tanpa Cause of decline', () => {
    const xml = readFileSync(SEKSI_UW, 'utf-8')
    for (const t of [EUW.validitas, EUW.hariAtau, EUW.waitingBind]) expect(xml, t).toContain(`<pyLabelFieldValue>${t}</pyLabelFieldValue>`)
    expect(xml).not.toContain('Cause of decline')
  })
  it('opsi UW; Decline UW tanpa Cause of decline; Validity Of Acceptance angka bulat', () => {
    expect(opsiEmailType('uw').map((o) => `${o.kode}:${o.label}`)).toEqual(['1:Accept', '2:Reject', '3:Ask', '7:Decline'])
    const isi = { ...isianSubmitKosong(), emailType: '7' as const, comment: 'UJI' }
    expect(galatIsianSubmit(isi, 'uw')).toBe('')
    expect(galatIsianSubmit(isi, 'teamLeader')).toBe(T.wajibSebab)
    expect(galatIsianSubmit({ ...isi, daysValidity: '7x' }, 'uw')).toBe(T.hariAngka)
    expect(galatIsianSubmit({ ...isi, daysValidity: '14' }, 'uw')).toBe('')
    expect(galatIsianSubmit({ ...isi, comment: ' ' }, 'uw')).toBe(T.wajibComment)
  })
  it('render UW: 4 opsi, baris Validity + tanggal server, tanpa tombol Cause of decline, Comment Decline tetap dapat diisi', () => {
    const v: TampilanSubmit = {
      bolehSubmit: true, alasan: '', opsiBanding: [], sebabTc: [], sebabNonTc: [], pctRnmOtomatis: { rate: '', riComm: '' }, riwayat: [],
      tahap: 'uw', deductibleBedaMataUang: false, deductibleAcceptance: false, commentDeductible: '',
      daysValidity: '14', dateValidity: '22-10-2026', waitingBindDate: '08-11-2026', confirmBinding: false, receivedRiSlip: false,
    }
    const s = {
      v, isi: { ...isianSubmitKosong(), emailType: '7', daysValidity: '14' }, proteksi: [], galat: null, mengajukan: false, memeriksa: false,
      berhasil: null, konfirmasi: false, validitas: { dateValidity: '22-10-2026', waitingBindDate: '08-11-2026' },
    } as unknown as SubmitMarketing
    const html = renderToStaticMarkup(<BlokEmailType s={s} />)
    expect(html).toContain('Ask</label>')
    expect(html).not.toContain('Banding</label>')
    expect(html).not.toContain(E.pilihSebab)
    expect(html).toContain(EUW.validitas)
    expect(html).toContain('>22/10/2026<')
    expect(html).toContain('>08/11/2026<')
    expect(html).toMatch(/<textarea id="nbf-et-comment" rows="5"><\/textarea>/)
  })
  it('konfirmasi UW hanya Decline; Validity dihitung ulang server saat lepas fokus', () => {
    expect(SUMBER).toContain("if (isi.emailType === '7' || isi.emailType === '9' || (isi.emailType === '2' && !tahapUW(v?.tahap))) setKonfirmasi(true)")
    expect(SUMBER).toContain('const h = await hitungValiditas(caseId, hari)')
    expect(SUMBER).toContain('await muatValiditas(isi.daysValidity.trim())')
    // Klik Email Type UW = GetDateValidity_ACT mode default (work owner 08-10-2026: "saat uw pilih accept ... otomatis").
    expect(SUMBER).toContain("if (tahapUW(v?.tahap)) void muatValiditas('')")
    expect(SUMBER).toContain('onBlur={() => void s.ubahHari()}')
  })
})

describe('proteksi Banding seperti Accept (Protection_Act langkah 8, work owner 08-10-2026 "iya disamakan")', () => {
  const SUMBER = readFileSync(join(__dirname, 'BlokEmailType.tsx'), 'utf8').replace(/\r\n/g, '\n')
  it('Accept dan Banding diperiksa server dan ditahan bila ada pesan; Reject / Decline / Ask / Revise tidak', () => {
    expect(['1', '2', '3', '4', '7', '9'].map((k) => diproteksi(k))).toEqual([true, false, false, true, false, false])
    expect(SUMBER).toContain('if (!diproteksi(baru.emailType)) return')
    expect(SUMBER).toContain('emailType: baru.emailType as KodeEmailType,')
    expect(SUMBER).toContain('(!diproteksi(isi.emailType) || proteksi.length === 0)')
  })
})

describe('tahap Binding (tiket 64)', () => {
  const DDL_B = 'D:\\migrasi\\RNM\\DDL\\EmailTypeBinding.xml'
  const SEKSI_B = 'D:\\migrasi\\RNM\\NB FacIn\\Section\\EmailSectionCeding.xml'
  const REVISE = 'D:\\migrasi\\RNM\\NB FacIn\\Section\\ReviseNotification.xml'
  const SUMBER = readFileSync(join(__dirname, 'BlokEmailType.tsx'), 'utf8').replace(/\r\n/g, '\n')
  const vBinding = (confirmBinding: boolean): TampilanSubmit => ({
    bolehSubmit: true, alasan: '', opsiBanding: ['UJI-JABATAN'], sebabTc: [], sebabNonTc: [], pctRnmOtomatis: { rate: '', riComm: '' }, riwayat: [],
    tahap: 'binding', deductibleBedaMataUang: false, deductibleAcceptance: false, commentDeductible: '',
    daysValidity: '', dateValidity: '', waitingBindDate: '', confirmBinding, receivedRiSlip: false,
  })
  const render = (v: TampilanSubmit, emailType: string) =>
    renderToStaticMarkup(
      <BlokEmailType
        s={{ v, isi: { ...isianSubmitKosong(), emailType }, proteksi: [], proteksiBinding: [], galat: null, mengajukan: false, memeriksa: false, membinding: false, berhasil: null, konfirmasi: false } as unknown as SubmitMarketing}
      />,
    )

  it.skipIf(!existsSync(DDL_B))('opsi = PromptList EmailTypeBinding.xml (1 Accept / 4 Banding / 7 Decline / 9 Revise)', () => {
    const xml = readFileSync(DDL_B, 'utf-8')
    expect(xml).toContain('ASM-FW-GISFW-WORK!EMAILTYPEBINDING')
    for (const o of EB.opsi) expect(xml).toMatch(new RegExp(`<pyLocalizedValue>${o.label}</pyLocalizedValue>\\s*<pyStandardValue>${o.kode}</pyStandardValue>`))
  })
  it.skipIf(!existsSync(SEKSI_B))('label Binding / Received RI Slip = EmailSectionCeding; teks Revise = ReviseNotification', () => {
    const xml = readFileSync(SEKSI_B, 'utf-8')
    expect(xml).toContain(`<pyLabel>${EB.tombolBinding}</pyLabel>`)
    expect(xml).toContain(EB.receivedRiSlip)
    if (existsSync(REVISE)) expect(readFileSync(REVISE, 'utf-8')).toContain(`>${EB.yakinRevise}<`)
  })
  it('tombol Binding hanya ditahan pesan penahan Binding dari server (ProtectSpreading / ProtectDeductible), bukan semua proteksi', () => {
    expect(SUMBER).toContain('disabled={!bisaUbah || v.confirmBinding || s.membinding || s.memeriksa || s.proteksiBinding.length > 0}')
    expect(SUMBER).toContain('setProteksiBinding(h.pesanBinding ?? [])')
    const html = renderToStaticMarkup(
      <BlokEmailType
        s={{ v: vBinding(false), isi: { ...isianSubmitKosong(), emailType: '1' }, proteksi: ['UJI-MOID'], proteksiBinding: [], galat: null, mengajukan: false, memeriksa: false, membinding: false, berhasil: null, konfirmasi: false } as unknown as SubmitMarketing}
      />,
    )
    expect(html).toMatch(new RegExp(`<button type="button" class="btn nbf-tema__tombol btn--sm">${EB.tombolBinding}</button>`))
  })
  it('opsi Binding; Accept butuh tombol Binding; Revise lewat konfirmasi', () => {
    expect(opsiEmailType('binding').map((o) => `${o.kode}:${o.label}`)).toEqual(['1:Accept', '4:Banding', '7:Decline', '9:Revise'])
    const isi = { ...isianSubmitKosong(), emailType: '1' as const, comment: 'UJI' }
    expect(galatIsianSubmit(isi, 'binding', false)).toBe(T.wajibBinding)
    expect(galatIsianSubmit(isi, 'binding', true)).toBe('')
    expect(galatIsianSubmit({ ...isi, emailType: '9' }, 'binding', false)).toBe('')
    expect(SUMBER).toContain("isi.emailType === '9' ? EB.yakinRevise : E.yakinDecline")
  })
  it('render: tombol Binding saat Accept (aktif sebelum di-bind), Received RI Slip baru aktif sesudah Binding', () => {
    const belum = render(vBinding(false), '1')
    expect(belum).toContain('Revise</label>')
    expect(belum).toMatch(new RegExp(`<button type="button" class="btn nbf-tema__tombol btn--sm">${EB.tombolBinding}</button>`))
    expect(belum).toMatch(/<input type="checkbox" disabled=""\/> Received RI Slip/)
    const sudah = render(vBinding(true), '1')
    expect(sudah).toMatch(new RegExp(`<button type="button" class="btn nbf-tema__tombol btn--sm" disabled="">${EB.tombolBinding}</button>`))
    expect(sudah).toMatch(/<input type="checkbox"\/> Received RI Slip/)
    expect(render(vBinding(false), '9')).not.toContain(EB.tombolBinding + '</button>')
  })
})

describe('BlokEmailType - Binding lalu Accept (bug DEV 10-10-2026)', () => {
  it('sesudah tombol Binding sukses, proteksi diperiksa ulang supaya pesan "Tekan tombol Binding…" tidak menahan Submit', () => {
    const s = readFileSync(join(__dirname, 'BlokEmailType.tsx'), 'utf8').replace(/\r\n/g, '\n')
    expect(s).toMatch(/const h = await tekanBinding\(caseId\)\s*setV\(h\)\s*const baru: IsianSubmit = \{ \.\.\.isi, receivedRiSlip: h\.receivedRiSlip \}\s*setIsi\(baru\)[\s\S]*?await periksa\(baru\)/)
  })
})
