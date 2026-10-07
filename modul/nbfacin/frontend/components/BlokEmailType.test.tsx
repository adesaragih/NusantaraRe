// Blok Email Type + Submit tahap marketing + popup Cause of Decline (tiket 57) - data uji sintetis.

import { existsSync, readFileSync } from 'node:fs'
import { join } from 'node:path'
import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import type { PilihanSebabTc, TampilanSubmit } from '../api'
import { EMAIL_TYPE as E, SEBAB_DECLINE as S, TEKS_EMAIL_TYPE as T, TEKS_SEBAB_DECLINE as TS } from '../labels'
import { ApiFailure } from '../../../../inti/frontend/klien'
import PortalOpportunity from '../pages/PortalOpportunity'
import { navigasi } from './BacaSajaKasus'
import BlokEmailType, { galatIsianSubmit, isianSubmitKosong, pesanProteksiDari, tanggalRiwayat, type SubmitMarketing } from './BlokEmailType'
import { centangPrioritas, komentarDecline, optionAktif, pctRnmAktif, pctSobAktif, pilihanKosong, remarksAktif, teksSebabTc, validasiSebab } from './sebabDecline'

const NB = 'D:\\migrasi\\RNM\\NB FacIn\\'
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
  it('case sudah tidak di tahap marketing: form disembunyikan, keterangan jelas, riwayat tetap', () => {
    const v: TampilanSubmit = {
      bolehSubmit: false, alasan: 'UJI-ALASAN', opsiBanding: [], sebabTc: [], sebabNonTc: [], pctRnmOtomatis: { rate: '', riComm: '' },
      riwayat: [{ tanggal: '07-10-2026 09:25', pic: 'UJI-PIC', approval: 'Accept', status: 'Offer', comment: 'UJI-OK' }],
    }
    const s = { v, isi: isianSubmitKosong(), proteksi: [], galat: null, mengajukan: false, memeriksa: false, berhasil: null, konfirmasi: false } as unknown as SubmitMarketing
    const html = renderToStaticMarkup(<BlokEmailType s={s} />)
    expect(html).toContain(T.sudahDiproses)
    expect(html).toContain('UJI-ALASAN')
    expect(html).not.toContain('type="radio"')
    expect(html).not.toContain('<textarea')
    expect(html).toContain('>UJI-OK<')
  })
  it('halaman: blok FIRE di bawah detail, Submit aktif hanya bila boleh', () => {
    expect(HALAMAN).toContain('{kasusFire && <BlokEmailType s={sub} />}')
    expect(HALAMAN).toContain('disabled={!kasusFire || !sub.bisaSubmit}')
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
  it('halaman: General dan detail dibungkus BacaSajaKasus; Save for later nonaktif; banner', () => {
    expect(HALAMAN).toContain('const bacaSaja = dariServer?.bolehUbah === false')
    expect(HALAMAN.match(/<BacaSajaKasus aktif=\{bacaSaja\}>/g)).toHaveLength(2)
    expect(HALAMAN).toContain('disabled={menyimpan || tanggalSalah || bacaSaja}')
  })
})
