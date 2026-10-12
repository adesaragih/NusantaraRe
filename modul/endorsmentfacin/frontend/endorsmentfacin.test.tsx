// Uji frontend Endorsment Fac In - irisan I0-I2 (keputusan work owner 11-10-2026, W-EDM-01).
//
// ⛔ Komentar bukti yang tidak pernah diperiksa adalah HIASAN: bagian pertama membuka korpus (READ-ONLY) dan memastikan
// setiap teks `labels.ts` memang ada di rule yang disebut. Korpus tidak terjangkau = uji DILEWATI, bukan gagal.

import { existsSync, readFileSync } from 'node:fs'
import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { buatKasus, periksaAddendum } from './api'
import { DETAIL_EDM, FORM, JENIS_EDM, PERIKSA, PORTAL } from './labels'
import { HALAMAN_AWAL_EDMFACIN, HALAMAN_EDMFACIN, KELOMPOK_EDMFACIN, NAMA_EDMFACIN } from './menu'
import FormAddendum from './pages/FormAddendum'
import PortalEndorsement from './pages/PortalEndorsement'

const EDM = 'D:\\migrasi\\RNM\\Endorsment Fac In\\'
const DDL_EDMTYPE = 'D:\\migrasi\\RNM\\DDL\\EdmType.xml'
const ada = (f: string) => existsSync(f)
const baca = (f: string) => readFileSync(f, 'utf-8')

describe('label = korpus', () => {
  it.skipIf(!ada(EDM + 'Section\\SFAPortal_EndorsementHeader.xml'))('portal = SFAPortal_EndorsementHeader', () => {
    const xml = baca(EDM + 'Section\\SFAPortal_EndorsementHeader.xml')
    expect(xml).toContain(`<pyValue>${PORTAL.judul}</pyValue>`)
    expect(xml).toContain(`<pyLabel>${PORTAL.buat}</pyLabel>`)
  })

  it.skipIf(!ada(EDM + 'Section\\WorkPrimaryDetails.xml'))('medan = WorkPrimaryDetails', () => {
    const xml = baca(EDM + 'Section\\WorkPrimaryDetails.xml')
    for (const l of [FORM.policyNo, FORM.sumber, FORM.catatan, FORM.jenis, FORM.tanggal, FORM.status, FORM.retro]) {
      expect(xml, l).toContain(`<pyLabelFieldValue>${l}</pyLabelFieldValue>`)
    }
    for (const p of [FORM.pilihJenis, FORM.jenisBaru]) expect(xml, p).toContain(p)
    // Biaya / Survey tidak dirender: kontainer terdekat sebelum keduanya `pyContainerVisibleWhen = never`.
    for (const c of ['Biaya', 'Survey']) {
      const i = xml.indexOf(`<pyCheckboxCaption>${c}</pyCheckboxCaption>`)
      expect(i, c).toBeGreaterThan(0)
      const sebelum = xml.slice(0, i)
      const m = [...sebelum.matchAll(/<pyContainerVisibleWhen>([^<]*)<\/pyContainerVisibleWhen>/g)].pop()
      expect(m?.[1], c).toBe('never')
    }
  })

  it.skipIf(!ada(EDM + 'Section\\EdmType1.xml'))('Detail EDM Type = EdmType1', () => {
    expect(baca(EDM + 'Section\\EdmType1.xml')).toContain(`<pyLabelFieldValue>${FORM.detail}</pyLabelFieldValue>`)
  })

  it.skipIf(!ada(EDM + 'Section\\crmNewHarnessButtons.xml'))('tombol = crmNewHarnessButtons', () => {
    const xml = baca(EDM + 'Section\\crmNewHarnessButtons.xml')
    for (const t of [FORM.ok, FORM.batal]) expect(xml, t).toContain(`<pyLabel>${t}</pyLabel>`)
  })

  it.skipIf(!ada(DDL_EDMTYPE))('Endorsement Type = DDL EdmType.xml (nilai berpasangan teks)', () => {
    const xml = baca(DDL_EDMTYPE)
    for (const j of JENIS_EDM) {
      expect(xml, j.label).toMatch(
        new RegExp(`<pyLocalizedValue>${j.label.replace(/[/]/g, '\\/')}</pyLocalizedValue>\\s*<pyStandardValue>${j.nilai}</pyStandardValue>`),
      )
    }
  })

  // Pasangan nilai-teks dicocokkan lewat KUNCI (`<APPEND>` CARI1 lalu `<LAST>` CARI2 pada daftar yang sama), bukan lewat
  // urutan dua daftar terpisah (CLAUDE.md §4a jebakan 6).
  it.skipIf(!ada(EDM + 'Activity\\SetEdmType.xml'))('Detail EDM Type = SetEdmType (pasangan per daftar)', () => {
    const xml = baca(EDM + 'Activity\\SetEdmType.xml')
    const isi = [...xml.matchAll(/<(PropertiesName|PropertiesValue)>([^<]*)<\/\1>/g)].map((m) => m[2] ?? '')
    const pasangan = new Map<string, string[]>()
    for (let i = 0; i + 3 < isi.length; i++) {
      const m = /^Type(\w+)\.pxResults\(&lt;APPEND&gt;\)\.CARI1$/.exec(isi[i] ?? '')
      if (!m || isi[i + 2] !== `Type${m[1]}.pxResults(&lt;LAST&gt;).CARI2`) continue
      const daftar = pasangan.get(m[1] ?? '') ?? []
      daftar.push(`${isi[i + 1]}=${(isi[i + 3] ?? '').replace(/^"|"$/g, '')}`)
      pasangan.set(m[1] ?? '', daftar)
    }
    const dariKorpus = (nama: string) => [...new Set(pasangan.get(nama) ?? [])]
    const dariLabel = (k: string) => (DETAIL_EDM[k] ?? []).map((p) => `${p.nilai}=${p.label}`)
    // Instrumen diuji pada jawaban yang diketahui: Penambahan memuat 12 DUA kali (sekali dari langkah Pengurangan, K-046).
    expect(pasangan.get('Penambahan')?.filter((p) => p.startsWith('12=')).length).toBe(2)
    expect(dariLabel('1')).toEqual(dariKorpus('Penambahan'))
    expect(dariLabel('2')).toEqual(dariKorpus('Pengurangan'))
    expect(dariLabel('3')).toEqual(dariKorpus('Perubahan'))
  })
})

describe('menu', () => {
  it('nama folder, label 900, halaman awal portal', () => {
    expect(NAMA_EDMFACIN).toBe('endorsmentfacin')
    expect(KELOMPOK_EDMFACIN).toBe('Endorsment Fac In')
    expect(HALAMAN_EDMFACIN).toContain(HALAMAN_AWAL_EDMFACIN)
    expect(HALAMAN_AWAL_EDMFACIN).toBe('endorsmentfacin-portal')
  })
})

describe('layar', () => {
  it('portal: judul Addendum + Create Addendum, tanpa grid (sumber tidak ada di korpus)', () => {
    const html = renderToStaticMarkup(<PortalEndorsement onBuat={() => {}} />)
    expect(html).toContain(`>${PORTAL.judul}<`)
    expect(html).toContain(`>${PORTAL.buat}<`)
    expect(html).not.toContain('<table')
  })

  it('form awal: medan WorkPrimaryDetails, EDM Type & Detail tersembunyi, Retro bawaan 0, OK mati', () => {
    const html = renderToStaticMarkup(<FormAddendum onBatal={() => {}} onDibuat={() => {}} />)
    for (const l of [FORM.policyNo, FORM.sumber, FORM.catatan, FORM.jenis, FORM.tanggal, FORM.status, FORM.retro]) {
      expect(html, l).toContain(l)
    }
    expect(html).not.toContain('Biaya')
    expect(html).not.toContain('type="checkbox"')
    // Endorsement Status bawaan Internal; Retro bawaan No (Yes / No, W-EDM-13/14).
    expect(html).toContain('<option value="0" selected="">Internal</option>')
    expect(html).toMatch(/name="edmf-retro" checked="" value="0"\/> No/)
    // EdmTypeNew hanya bila Endorsement Type 4; Detail hanya bila EdmTypeNew 1/2/3.
    expect(html).not.toContain(`<label class="field__label">${FORM.jenisBaru}`)
    expect(html).not.toContain(`<label class="field__label">${FORM.detail}`)
    expect(html).toContain('name="edmf-retro" checked="" value="0"')
    expect(html).not.toContain('name="edmf-retro" checked="" value="1"')
    expect(html).toMatch(/<button type="button" class="btn btn--primary edmf-tombol edmf-tombol--utama" disabled=""[^>]*>OK<\/button>/)
    // Templat NB FacIn: tanggal teks dd/mm/yyyy (bukan input date bawaan yang tampil mm/dd/yyyy), alasan OK mati tampil.
    expect(html).toContain('placeholder="dd/mm/yyyy"')
    expect(html).not.toMatch(/<input class="field__input" type="date"/)
    expect(html).toContain(`Required: ${FORM.policyNo}, ${FORM.sumber}, ${FORM.catatan}, ${FORM.jenis}, ${FORM.tanggal}.`)
    // Panel rincian gerbang tidak ditampilkan (work owner 12-10-2026 "ini tidak perlu ditampilkan").
    expect(html).not.toContain("edmf-periksa")
  })

  it('portal menampilkan pesan kasus lahir, dan tanpa pesan bila kosong', () => {
    const ada = renderToStaticMarkup(<PortalEndorsement onBuat={() => {}} pesan={PERIKSA.dibuat('EDM-UJI')} onTutupPesan={() => {}} />)
    expect(ada).toContain('Endorsement case EDM-UJI created.')
    expect(renderToStaticMarkup(<PortalEndorsement onBuat={() => {}} />)).not.toContain('alert')
  })
})

describe('api', () => {
  it('periksaAddendum: POST /api/endorsmentfacin/periksa dengan isian apa adanya', async () => {
    const asli = globalThis.fetch
    let url = ''
    let metode = ''
    let badan = ''
    globalThis.fetch = (async (u: RequestInfo | URL, init?: RequestInit) => {
      url = String(u)
      metode = init?.method ?? ''
      badan = String(init?.body ?? '')
      return new Response(
        JSON.stringify({ ditolak: false, lengkap: false, pesan: [], gerbang: [], bisnis: { hasil: 'lolos', jenisBisnis: 'UJI' }, tanggal: null }),
        { status: 200 },
      )
    }) as typeof fetch
    try {
      const isi = { policyNo: 'UJI-1', edmType: '4', edmTypeNew: '1', endorsementDate: '01-02-2026', endorsementInternalRetro: '0' }
      const h = await periksaAddendum(isi)
      expect(metode).toBe('POST')
      expect(url).toMatch(/\/api\/endorsmentfacin\/periksa$/)
      expect(JSON.parse(badan)).toEqual(isi)
      expect(h.lengkap).toBe(false)
      expect(h.tanggal).toBeNull()
    } finally {
      globalThis.fetch = asli
    }
  })

  it('buatKasus: POST /api/endorsmentfacin/kasus; penolakan = 200 tanpa caseId (hasil bisnis, bukan galat)', async () => {
    const asli = globalThis.fetch
    let url = ''
    let badan = ''
    globalThis.fetch = (async (u: RequestInfo | URL, init?: RequestInit) => {
      url = String(u)
      badan = String(init?.body ?? '')
      return new Response(JSON.stringify({ caseId: '', periksa: { ditolak: true, lengkap: true, pesan: [], gerbang: [] } }), {
        status: 200,
      })
    }) as typeof fetch
    try {
      const isi = {
        policyNo: 'UJI-1',
        edmType: '2',
        edmTypeNew: '',
        endorsementDate: '01-02-2026',
        endorsementInternalRetro: '0',
        note: 'UJI catatan',
        endorsementSource: '',
        endorsementSourceNote: 'UJI sumber',
        endorsementStatus: '0',
        detailType: '',
      }
      const h = await buatKasus(isi)
      expect(url).toMatch(/\/api\/endorsmentfacin\/kasus$/)
      expect(JSON.parse(badan)).toEqual(isi)
      expect(h.caseId).toBe('')
      expect(h.periksa?.ditolak).toBe(true)
    } finally {
      globalThis.fetch = asli
    }
  })
})
