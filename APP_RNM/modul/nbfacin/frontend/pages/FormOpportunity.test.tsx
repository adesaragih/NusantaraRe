// Paritas form "Opportunity" dengan gambar Pega 02-10-2026 - tiket 26.
// Halaman DIRENDER (react-dom/server). Sumbernya hanya gambar, jadi yang dijaga di sini: urutan,
// kelompok, tanda wajib, keadaan awal, dan TIDAK adanya hal yang tidak ada di gambar.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import {
  FORM_OPPORTUNITY as F,
  NILAI_AWAL_OPPORTUNITY as AWAL,
  OPSI_OPPORTUNITY_SOURCE,
  POPUP_CHOOSE_ACCOUNT as POPUP,
  TEKS_FORM_OPPORTUNITY as TEKS,
  TOMBOL_FORM_OPPORTUNITY as TOMBOL,
} from '../labels'
import FormOpportunity, { PopupChooseAccount } from './FormOpportunity'

const PEMILIK = 'UJI-PEMILIK'
const HTML = renderToStaticMarkup(<FormOpportunity pemilik={PEMILIK} />)
const SUMBER = readFileSync(join(__dirname, 'FormOpportunity.tsx'), 'utf8').replace(/\r\n/g, '\n')

/** Label medan berurutan, dengan penanda wajib. */
const label = [...HTML.matchAll(/<(?:label|span) class="field__label">([^<]*)(<span class="field__req">\*<\/span>)?/g)].map((m) => ({
  teks: m[1],
  wajib: m[2] !== undefined,
}))

/** Isi setiap `<select>` berurutan. */
const dropdown = [...HTML.matchAll(/<select[^>]*>([\s\S]*?)<\/select>/g)].map((m) =>
  [...(m[1] ?? '').matchAll(/<option value="([^"]*)"( selected="")?>([^<]*)<\/option>/g)].map((o) => ({
    nilai: o[1],
    terpilih: o[2] !== undefined,
    teks: o[3],
  })),
)

describe('FormOpportunity = gambar Pega (keadaan awal)', () => {
  it('akar halaman memasang kelas modul dan judul Opportunity', () => {
    expect(HTML.startsWith('<div class="nbfacin">')).toBe(true)
    expect(HTML).toContain(`<h4 class="panel__title">${F.judul}</h4>`)
  })

  it('urutan label: atas (tanggal, Owner), kolom kiri, kolom kanan, bawah - Type Of Facultative belum tampil', () => {
    expect(label.map((l) => l.teks)).toEqual([
      F.tanggalTutup,
      F.owner,
      F.namaProspek,
      F.grupBisnis,
      F.classOfBusiness,
      F.typeOfInward,
      F.phase,
      F.stage,
      F.sumber,
      F.statusBisnis,
      F.deskripsi,
    ])
  })

  it('Owner menampilkan nama pengguna yang login, sebagai teks (bukan kotak isian)', () => {
    expect(HTML).toContain(`<div class="nbf-opp__owner">${PEMILIK}</div>`)
  })

  it('tanda wajib persis seperti gambar', () => {
    const wajib = label.filter((l) => l.wajib).map((l) => l.teks)
    expect(wajib).toEqual([F.tanggalTutup, F.namaProspek, F.classOfBusiness, F.typeOfInward, F.phase, F.statusBisnis])
  })

  it('Opportunity Source = 13 pilihan tangkapan layar, verbatim dan berurutan', () => {
    expect([...OPSI_OPPORTUNITY_SOURCE]).toEqual([
      'Iklan', 'Analisa Referral', 'Rujukan Pelanggan', 'Direct Mail', 'Email', 'Rujukan Karyawan',
      'Telepon Masuk', 'Partner', 'Seminar', 'Sosial Media', 'Pameran', 'Web', 'Whatsapp',
    ])
  })

  it('dropdown keadaan awal: Type Of Inward kosong, Phase/Business Status bernilai gambar, Source lengkap (C-1, C-7)', () => {
    // Type Of Inward, Phase, Opportunity Source, Business Status.
    expect(dropdown).toHaveLength(4)
    const [inward, phase, sumber, status] = dropdown
    expect(inward).toEqual([
      { nilai: '', terpilih: true, teks: AWAL.inwardKosong },
      { nilai: AWAL.typeOfInward, terpilih: false, teks: AWAL.typeOfInward },
    ])
    expect(phase!.filter((o) => o.nilai !== '')).toEqual([{ nilai: AWAL.phase, terpilih: true, teks: AWAL.phase }])
    expect(sumber![0]).toEqual({ nilai: '', terpilih: true, teks: AWAL.sumberKosong })
    expect(sumber!.slice(1).map((o) => o.teks)).toEqual([...OPSI_OPPORTUNITY_SOURCE])
    expect(sumber!.slice(1).every((o) => !o.terpilih && o.nilai === o.teks)).toBe(true)
    expect(status!.filter((o) => o.nilai !== '')).toEqual([{ nilai: AWAL.statusBisnis, terpilih: true, teks: AWAL.statusBisnis }])
  })

  it('Type Of Facultative hanya dirender bila Type Of Inward = Facultative (C-7)', () => {
    expect(HTML).not.toContain(F.typeOfFacultative)
    const syarat = SUMBER.indexOf('{typeOfInward === AWAL.typeOfInward && (')
    expect(syarat).toBeGreaterThan(-1)
    expect(SUMBER.indexOf('label={F.typeOfFacultative}')).toBeGreaterThan(syarat)
  })

  it('Stage read-only berisi Opportunity (C-3)', () => {
    // react-dom/server menulis `readonly` huruf kecil, sebelum `value`.
    expect(HTML).toMatch(new RegExp(`<input[^>]*readonly=""[^>]*value="${AWAL.stage}"`))
  })

  it('tiga tombol Group Business berurutan; Search Group Business hidup, dua lainnya nonaktif (C-2, C-8)', () => {
    const tombol = [...HTML.matchAll(/<button([^>]*)>([^<]*)<\/button>/g)].map((m) => ({ teks: m[2], nonaktif: (m[1] ?? '').includes('disabled') }))
    expect(tombol).toEqual([
      { teks: TOMBOL.cariGrup, nonaktif: false },
      { teks: TOMBOL.perusahaanBaru, nonaktif: true },
      { teks: TOMBOL.grupBaru, nonaktif: true },
    ])
    expect(SUMBER).toContain('onClick={() => setCariGrup(true)}')
  })

  it('tanpa tombol simpan - tidak ada di gambar (C-5); popup tertutup pada keadaan awal', () => {
    expect(HTML).not.toMatch(/type="submit"/)
    expect(HTML).not.toContain(POPUP.judul)
  })

  it('nol catatan pengembang di layar (C-6)', () => {
    expect(HTML).not.toContain('panel__note')
    expect(HTML).not.toMatch(/belum terverifikasi|tiket|dugaan/i)
  })
})

describe('PopupChooseAccount = tangkapan layar Search Group Business (C-8)', () => {
  const P = renderToStaticMarkup(<PopupChooseAccount onTutup={() => {}} />)

  it('judul ChooseAccount, kotak Search, tombol Search', () => {
    expect(P).toContain(POPUP.judul)
    expect(P).toMatch(new RegExp(`class="field__label">${POPUP.cari}<`))
    expect(P).toMatch(new RegExp(`<button[^>]*disabled=""[^>]*>${POPUP.tombolCari}</button>`))
  })

  it('kolom grid berurutan: nomor, Insured ID, Insured Name, Group Business, kolom Choose', () => {
    const kolom = [...P.matchAll(/<th[^>]*>([^<]*)<\/th>|<th[^>]*\/>/g)].map((m) => m[1] ?? '')
    expect(kolom).toEqual(['', ...POPUP.kolom, ''])
  })

  it('daftar belum punya sumber data = BelumTersedia, bukan tabel kosong', () => {
    expect(P).toContain(`${TEKS.daftarAccount} tidak dapat dimuat saat ini.`)
    expect(P).not.toContain('<tbody')
  })
})
