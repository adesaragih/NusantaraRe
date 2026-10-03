// Paritas form "Opportunity" dengan gambar Pega 02-10-2026 - tiket 26.
// Halaman DIRENDER (react-dom/server). Sumbernya hanya gambar, jadi yang dijaga di sini: urutan,
// kelompok, tanda wajib, keadaan awal, dan TIDAK adanya hal yang tidak ada di gambar.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import {
  FORM_OPPORTUNITY as F,
  KEPALA_PORTAL,
  NILAI_AWAL_OPPORTUNITY as AWAL,
  OPSI_OPPORTUNITY_SOURCE,
  POPUP_CHOOSE_ACCOUNT as POPUP,
  TOMBOL_FORM_OPPORTUNITY as TOMBOL,
} from '../labels'
import FormOpportunity, { PopupChooseAccount, cobSah } from './FormOpportunity'

const PEMILIK = 'UJI-PEMILIK'
const HTML = renderToStaticMarkup(<FormOpportunity pemilik={PEMILIK} />)
const SUMBER = readFileSync(join(__dirname, 'FormOpportunity.tsx'), 'utf8').replace(/\r\n/g, '\n')

/** Label medan berurutan, dengan penanda wajib. */
const label = [...HTML.matchAll(/<(?:label|span) class="field__label"[^>]*>([^<]*)(<span class="field__req">\*<\/span>)?/g)].map((m) => ({
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
    expect(SUMBER).toContain('const facultative = typeOfInward === AWAL.typeOfInward')
    const syarat = SUMBER.indexOf('{facultative && (')
    expect(syarat).toBeGreaterThan(-1)
    expect(SUMBER.indexOf('label={F.typeOfFacultative}')).toBeGreaterThan(syarat)
  })

  it('Stage read-only berisi Opportunity (C-3)', () => {
    // react-dom/server menulis `readonly` huruf kecil, sebelum `value`.
    expect(HTML).toMatch(new RegExp(`<input[^>]*readonly=""[^>]*value="${AWAL.stage}"`))
  })

  it('sesudah Choose: tiga tombol diganti teks Group Business terpilih + roda gigi yang membuka popup lagi (C-10)', () => {
    const blok = SUMBER.slice(SUMBER.indexOf('{grup ? ('), SUMBER.indexOf(') : (', SUMBER.indexOf('{grup ? (')))
    expect(blok).toContain('<span>{grup.groupBusiness}</span>')
    expect(blok).toContain('aria-label={TEKS_GRUP_BISNIS.ganti}')
    expect(blok).toContain('onClick={() => setCariGrup(true)}')
    expect(blok).not.toMatch(/TOMBOL\./)
    expect(SUMBER).toMatch(/onPilih=\{\(b\) => \{[\s\S]{0,120}setGrup\(b\)\s*setCariGrup\(false\)/)
  })

  it('Class Of Business = kotak isian dengan saran BUSINESS.NOTE untuk group terpilih (C-11)', () => {
    expect(HTML).toMatch(/<input[^>]*id="nbfacin-class-of-business"[^>]*list="nbfacin-saran-cob"/)
    expect(HTML).toContain('<datalist id="nbfacin-saran-cob"></datalist>')
    // Saran dimuat ulang setiap Group Business berganti, jawaban untuk group lama dibuang.
    expect(SUMBER).toContain('daftarClassOfBusiness(idGrup)')
    expect(SUMBER).toMatch(/\}, \[idGrup\]\)/)
    expect(SUMBER).toContain('if (!batal) setSaranCOB(h.baris)')
    expect(SUMBER).toContain('<option key={s.id} value={s.note} />')
    // Group berganti = isian Class Of Business lama dikosongkan.
    expect(SUMBER).toContain("if (b.groupBusinessId !== grup?.groupBusinessId) setClassOfBusiness('')")
  })

  it('Class Of Business di luar daftar dikosongkan: saat keluar kotak dan sebelum dikirim (work owner 03-10-2026)', () => {
    const saran = [
      { id: 'UJI-1', note: 'UJI FIRE' },
      { id: 'UJI-2', note: 'UJI EARTHQUAKE' },
    ]
    expect(cobSah('UJI FIRE', saran)).toBe('UJI FIRE')
    expect(cobSah('uji fire', saran)).toBe('')
    expect(cobSah('UJI FIR', saran)).toBe('')
    expect(cobSah('UJI FIRE', [])).toBe('')
    expect(cobSah('', saran)).toBe('')
    expect(SUMBER).toContain('onBlur={() => setClassOfBusiness((v) => cobSah(v, saranCOB))}')
    const buat = SUMBER.slice(SUMBER.indexOf('async function buat()'))
    expect(buat.indexOf('cobSah(classOfBusiness, saranCOB)')).toBeLessThan(buat.indexOf('medanKosong(cob)'))
    expect(buat).toContain('classOfBusiness: cob,')
  })

  it('tiga tombol Group Business berurutan; Search Group Business hidup, dua lainnya nonaktif (C-2, C-8)', () => {
    const grupHTML = HTML.slice(HTML.indexOf('class="nbf-opp__tombol"'))
    const tombol = [...grupHTML.matchAll(/<button([^>]*)>([^<]*)<\/button>/g)]
      .slice(0, 3)
      .map((m) => ({ teks: m[2], nonaktif: (m[1] ?? '').includes('disabled') }))
    expect(tombol).toEqual([
      { teks: TOMBOL.cariGrup, nonaktif: false },
      { teks: TOMBOL.perusahaanBaru, nonaktif: true },
      { teks: TOMBOL.grupBaru, nonaktif: true },
    ])
    expect(SUMBER).toContain('onClick={() => setCariGrup(true)}')
  })

  it('tombol Create opportunity di kanan atas, sebaris dengan judul (permintaan work owner 03-10-2026); popup tertutup', () => {
    const kepala = HTML.slice(HTML.indexOf('<div class="nbf-kepala">'), HTML.indexOf('<div class="nbf-opp__atas">'))
    expect(kepala).toContain(`<h4 class="panel__title">${F.judul}</h4>`)
    expect(kepala).toContain(`<button type="button" class="btn btn--primary">${KEPALA_PORTAL.buat.label}</button>`)
    expect(kepala.indexOf('panel__title')).toBeLessThan(kepala.indexOf('btn--primary'))
    expect(HTML).not.toMatch(/type="submit"/)
    expect(HTML).not.toContain(POPUP.judul)
  })

  it('Estimated Closing Date = kotak dd/mm/yyyy (bukan input date bawaan)', () => {
    expect(HTML).toContain('placeholder="dd/mm/yyyy"')
    expect(HTML.match(/type="date"/g)).toHaveLength(1)
    expect(HTML).toContain('class="nbf-tanggal__pemilih"')
  })

  it('Create opportunity: medan wajib diperiksa dulu, urut layar; tanpa permintaan bila ada yang kosong', () => {
    const awal = SUMBER.indexOf('function medanKosong(')
    const daftar = SUMBER.slice(awal, SUMBER.indexOf('return wajib.filter', awal))
    const urut = ['F.tanggalTutup', 'F.namaProspek', 'F.classOfBusiness', 'F.typeOfInward', 'F.typeOfFacultative', 'F.phase', 'F.statusBisnis']
    const letak = urut.map((u) => daftar.indexOf(u + ','))
    expect(letak.every((l) => l >= 0)).toBe(true)
    expect([...letak].sort((a, b) => a - b)).toEqual(letak)
    expect(daftar).toContain('...(facultative ?')
    const buat = SUMBER.slice(SUMBER.indexOf('async function buat()'))
    expect(buat.indexOf('if (k.length > 0) return')).toBeLessThan(buat.indexOf('buatOpportunity(isian)'))
    expect(buat).toContain("typeOfFacultative: facultative ? typeOfFacultative : ''")
  })

  it('nol catatan pengembang di layar (C-6)', () => {
    expect(HTML).not.toContain('panel__note')
    expect(HTML).not.toMatch(/belum terverifikasi|tiket|dugaan/i)
  })
})

describe('PopupChooseAccount = tangkapan layar Search Group Business (C-8, C-9)', () => {
  const P = renderToStaticMarkup(<PopupChooseAccount onTutup={() => {}} onPilih={() => {}} />)
  const POP = SUMBER

  it('judul ChooseAccount, kotak Search dapat diisi, tombol Search mengirim form', () => {
    expect(P).toContain(POPUP.judul)
    expect(P).toMatch(new RegExp(`class="field__label">${POPUP.cari}<`))
    expect(P).toMatch(new RegExp(`<button type="submit"[^>]*>${POPUP.tombolCari}</button>`))
    expect(P).not.toMatch(/<input[^>]*readonly/)
  })

  it('kolom grid berurutan: nomor, Insured ID, Insured Name, Group Business, kolom Choose', () => {
    const kolom = [...P.matchAll(/<th[^>]*>([^<]*)<\/th>|<th[^>]*\/>/g)].map((m) => m[1] ?? '')
    expect(kolom).toEqual(['', ...POPUP.kolom, ''])
  })

  it('baris: nomor berlanjut antar halaman, tiga kolom data, tombol Choose memanggil onPilih', () => {
    expect(POP).toContain('<td>{awal + i + 1}</td>')
    expect(POP).toMatch(/<td>\{b\.insuredId\}<\/td>\s*<td>\{b\.insuredName\}<\/td>\s*<td>\{b\.groupBusiness\}<\/td>/)
    expect(POP).toContain('onClick={() => onPilih(b)}')
    expect(POP).toContain('{POPUP.pilih}')
  })

  it('daftar dimuat saat dibuka dengan kotak kosong; Search memuat ulang dari halaman 1; paging memakai kata cari terakhir', () => {
    expect(POP).toMatch(/useEffect\(\(\) => \{\s*void muat\('', 1\)/)
    expect(POP).toContain('onKirim={() => void muat(kotak, 1)}')
    expect(POP).toContain('onPindah={(h) => void muat(kunciCari.current, h)}')
  })

  it('jawaban lama dibuang (nomor permintaan); galat ditampilkan apa adanya; tabel tanpa tbody sebelum data datang', () => {
    expect(POP).toContain('if (nomor === nomorPermintaan.current) setHasil(h)')
    expect(P).not.toContain('<tbody')
    expect(POP).toContain('<Gagal galat={galat} />')
  })
})
