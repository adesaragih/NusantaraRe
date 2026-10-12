// Inbox → layar kasus per tahap (11-10-2026): satu layar untuk Outstanding, Medical Check, dan Claim Analis; pita kuning
// "This case is at the … stage" dan tombol "Open claim detail" dihapus.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { BahasaUI } from '../../../inti/frontend/components/ui/bahasaUI'
import { TAHAP } from './labels'
import { TAHAP_JALUR, type Diagnosa, type Klaim, type Peserta } from './api'
import {
  bolehTolakOutstanding,
  bolehUbahAnalis,
  bolehUbahMedis,
  catatanAnalisBerubah,
  kabarPindah,
  kabarTutup,
  kalimatPenghalangTutup,
  KONFIRMASI_TUTUP,
  KonfirmasiPindah,
  konfirmasiPindah,
  KonfirmasiTutup,
  rekomendasiBerubah,
  TombolBawahTahap,
} from './pages/OutstandingClaimLife'
import { LABEL_TAB_INBOX, tabUntuk } from './pages/InboxClaimLife'
import { ApiFailure } from '../../../inti/frontend/klien'
import { GridDiagnosa } from './components/GridDiagnosa'
import { RincianPeserta } from './components/RincianPeserta'
import { PERAN } from '../../../inti/frontend/labels'

const SUMBER = readFileSync(join(__dirname, 'pages', 'OutstandingClaimLife.tsx'), 'utf8')
const RUTE = readFileSync(join(__dirname, 'rute.tsx'), 'utf8')

const en = (n: React.ReactNode) => renderToStaticMarkup(<BahasaUI.Provider value="en">{n}</BahasaUI.Provider>)

function klaimTahap(tahap: string, statusWork = ''): Klaim {
  return {
    id: 'CLMLF-UJI1',
    nomorKlaim: 'UJI-CLM-1',
    peserta: [],
    tahap,
    statusWork,
  } as unknown as Klaim
}

const tombol = (html: string) => [...html.matchAll(/<button[^>]*>([^<]*)<\/button>/g)].map((m) => m[1])

describe('Inbox membuka layar kasus yang benar per tahap', () => {
  it('pita kuning dan "Open claim detail" tidak ada lagi; rute tidak membuka layar Detail dari Inbox', () => {
    expect(SUMBER).not.toContain('This case is at the')
    expect(SUMBER).not.toContain('Open claim detail')
    expect(SUMBER).not.toContain('onDetail')
    expect(RUTE).not.toContain('onDetail')
    expect(RUTE).toContain("onPindah('outstanding')")
  })

  it('Medical Check: Send Back to Admin · Send to Claim Analyst; Save TIDAK di baris ini (12-10-2026); tanpa Close Claim', () => {
    const html = en(
      <TombolBawahTahap onGalat={() => {}} klaim={klaimTahap(TAHAP.medicalCheck)} sesudahPindah={() => {}} sesudahTutup={() => {}} />,
    )
    expect(tombol(html)).toEqual(['Send Back to Admin', 'Send to Claim Analyst'])
    expect(html).not.toContain('Close Claim')
    expect(html).not.toContain('Save to RNM')
  })

  it('Save Medical Check di kaki grid Participant Details (posisi Save Outstanding), pesan Progress saved.', () => {
    expect(SUMBER).toContain(') : bolehUbahMedis(klaim, peran) ? (')
    expect(SUMBER).toContain('onClick={() => void simpanRekomendasi()}')
    const bawah = SUMBER.slice(SUMBER.indexOf('export function TombolBawahTahap'), SUMBER.indexOf('export function KonfirmasiKirimMedis'))
    expect(bawah).not.toContain('TOMBOL_MEDIS.simpan')
    expect(bawah).not.toContain('onSimpan')
  })

  it('Claim Analis: dua jalur balik + Close Claim; tanpa Save Medical', () => {
    const html = en(
      <TombolBawahTahap onGalat={() => {}} klaim={klaimTahap(TAHAP.claimAnalis)} sesudahPindah={() => {}} sesudahTutup={() => {}} />,
    )
    const t = tombol(html)
    expect(t).toContain('Send Back to Admin')
    expect(t).toContain('Send Back to Medical')
    expect(html).toContain('Close Claim')
    expect(t).not.toContain('Save')
  })

  it('kasus tertutup: tanpa tombol', () => {
    expect(
      en(<TombolBawahTahap onGalat={() => {}} klaim={klaimTahap(TAHAP.claimAnalis, 'Resolved-Completed')} sesudahPindah={() => {}} sesudahTutup={() => {}} />),
    ).toBe('')
  })

  it('kabar Inbox per tujuan', () => {
    expect(kabarPindah('CLMLF-1', TAHAP_JALUR.outstanding)).toBe('Claim CLMLF-1 sent back to Claim Admin.')
    expect(kabarPindah('CLMLF-1', TAHAP_JALUR.claimAnalis)).toBe('Claim CLMLF-1 sent to Claim Analyst.')
    expect(kabarPindah('CLMLF-1', TAHAP_JALUR.medicalCheck)).toBe('Claim CLMLF-1 sent to Medical Advisor.')
    expect(kabarTutup('UJI-CLM-1')).toBe('Claim UJI-CLM-1 closed.')
  })
})

// ---------------------------------------------------------------------------
// Kotak MEDICAL ADVISOR (11-10-2026)
// ---------------------------------------------------------------------------

function pesertaUji(diagnosa: Diagnosa[] = [], kodeStatus = ''): Peserta {
  return {
    id: 'P1',
    nomorPremiList: 'UJI-PL-1',
    nomorPolis: 'UJI-POL-1',
    nomorSertifikat: '001',
    mataUang: 'IDR',
    isCheck: 'true',
    isAccept: 'true',
    namaTertanggung: 'UJI NAMA 1',
    tanggalLahir: '1980-02-03',
    tanggalKejadian: '',
    tanggalTerimaKlaim: '',
    tanggalDokumenLengkap: '',
    tanggalKonfirmasi: '',
    penandaTerimaKlaim: '',
    penandaTerimaKlaimAlasan: '',
    baris: [{ id: 'UJI-A1', kodeStatus, status: 'Outstanding' }],
    dokumen: [],
    diagnosa,
    kodeStatus,
  } as unknown as Peserta
}

const dg = (id: number, nama = '', kodeIcd = '', groupDiagnose = ''): Diagnosa => ({
  id,
  pesertaId: 'P1',
  urutan: id,
  kodeIcd,
  nama,
  groupDiagnose,
  kodeStatus: '',
})

const kepala = (html: string) => [...html.matchAll(/<th scope="col">([^<]*)<\/th>/g)].map((m) => m[1])

describe('grid diagnosa MEDICAL ADVISOR', () => {
  it('dapat diubah: Find Disease | DIAGNOSE | ICD CODE | GROUP DIAGNOSE | Delete; SATU + Add kanan atas; subjudul DISEASE', () => {
    const html = en(<GridDiagnosa klaimID="K1" peserta={pesertaUji([dg(1, 'UJI PENYAKIT', 'A00')])} onBerubah={() => {}} boleh />)
    expect(kepala(html)).toEqual(['Find Disease', 'Diagnose', 'ICD Code', 'Group Diagnose', 'Delete'])
    // 12-10-2026: tombol Add hanya SATU (yang atas), sebaris dengan subjudul DISEASE (gaya RECOMMENDATION).
    expect(html.match(/>\+ Add<\/button>/g)).toHaveLength(1)
    expect(html).toMatch(/<div class="diagnosa-grid__kepala"><span class="os__label">Disease<\/span><button[^>]*>\+ Add<\/button><\/div>/)
    expect(html.indexOf('+ Add')).toBeLessThan(html.indexOf('<table'))
    expect(html).toContain('>Find Disease</button>')
    expect(html).toContain('>Delete</button>')
    // Catatan "not configured" dibuang; dropdown aktif.
    expect(html).not.toContain('not configured')
    expect(html).not.toMatch(/<select[^>]*disabled/)
    expect(html).not.toContain('tidak dapat dimuat')
  })

  it('dropdown GROUP DIAGNOSE dengan daftar (mock): opsi tampil, nilai tersimpan terpilih', () => {
    const html = en(
      <GridDiagnosa
        klaimID="K1"
        peserta={pesertaUji([dg(1, 'UJI PENYAKIT', 'A00', 'UJI-G2')])}
        onBerubah={() => {}}
        boleh
        pilihan={['UJI-G1', 'UJI-G2']}
      />,
    )
    expect(html).not.toMatch(/<select[^>]*disabled/)
    expect(html).toContain('<option value="UJI-G1">UJI-G1</option>')
    expect(html).toContain('<option value="UJI-G2" selected="">UJI-G2</option>')
    expect(html).not.toContain('not configured')
  })

  it('baca saja: tiga kolom teks, tanpa tombol', () => {
    const html = en(
      <GridDiagnosa
        klaimID="K1"
        peserta={pesertaUji([dg(1, 'UJI PENYAKIT', 'A00', 'UJI-G1')])}
        onBerubah={() => {}}
        boleh={false}
      />,
    )
    expect(kepala(html)).toEqual(['Diagnose', 'ICD Code', 'Group Diagnose'])
    expect(html).not.toContain('<button')
    expect(html).not.toContain('<select')
    expect(html).toContain('<td class="diagnosa-grid__kelompok">UJI-G1</td>')
    expect(html).toContain('<span class="os__label">Disease</span>')
  })

  it('GROUP DIAGNOSE: nilai warisan di luar daftar tetap tampil terpilih; teks panjang utuh di title (12-10-2026)', () => {
    const panjang = 'Diseases of the blood and blood-forming organs and certain disorders involving the immune mechanism'
    const html = en(
      <GridDiagnosa
        klaimID="K1"
        peserta={pesertaUji([dg(1, 'UJI PENYAKIT', 'A00', 'UJI-WARISAN'), dg(2, 'UJI LAIN', 'B00', panjang)])}
        onBerubah={() => {}}
        boleh
        pilihan={['Neoplasms', panjang]}
      />,
    )
    expect(html).toContain('<option value="UJI-WARISAN" selected="">UJI-WARISAN</option>')
    expect(html).toContain(`title="${panjang}"`)
    expect(html).toContain(`<option value="${panjang}" selected="">${panjang}</option>`)
  })

  it('kosong: + Add sekali di kepala (boleh); kalimat English', () => {
    const boleh = en(<GridDiagnosa klaimID="K1" peserta={pesertaUji()} onBerubah={() => {}} boleh />)
    expect(boleh.match(/>\+ Add<\/button>/g)).toHaveLength(1)
    expect(boleh).toContain('No diagnosis yet. Click + Add to add a row.')
    expect(en(<GridDiagnosa klaimID="K1" peserta={pesertaUji()} onBerubah={() => {}} boleh={false} />)).toContain(
      'No diagnosis yet.',
    )
  })

  it('teks Indonesia lama tidak ada lagi di grid dan pencarian', () => {
    const grid = readFileSync(join(__dirname, 'components', 'GridDiagnosa.tsx'), 'utf8')
    const cari = readFileSync(join(__dirname, 'components', 'CariDiagnosa.tsx'), 'utf8')
    for (const t of ['>Aksi<', 'Terkunci', 'Belum ada diagnosa', 'Perubahan diagnosa gagal', 'BelumTersedia'])
      expect(grid).not.toContain(t)
    for (const t of ["'Mencari…'", "'Cari'", 'Tutup\r\n', 'Tidak ada diagnosa', 'TERPOTONG']) expect(cari).not.toContain(t)
  })
})

describe('RincianPeserta: kotak MEDICAL ADVISOR per tahap dan peran', () => {
  const r = (tahap: string, bolehMedis: boolean, kodeStatus = '') => {
    const p = pesertaUji([dg(1, 'UJI PENYAKIT', 'A00')], kodeStatus)
    const k = { ...klaimTahap(tahap), peserta: [p] } as Klaim
    return en(<RincianPeserta klaim={k} p={p} onMuatUlang={async () => {}} bolehMedis={bolehMedis} />)
  }

  it('Medical Check + Medical Advisor: kotak berjudul MEDICAL ADVISOR, grid dapat diubah', () => {
    const html = r(TAHAP.medicalCheck, true)
    expect(html).toContain('<legend class="os__label">Medical Advisor</legend>')
    expect(html).toContain('>Find Disease</button>')
  })

  it('Medical Check bukan Medical Advisor: baca saja', () => {
    const html = r(TAHAP.medicalCheck, false)
    expect(html).toContain('Medical Advisor')
    expect(html).not.toContain('>Find Disease</button>')
  })

  it('peserta sudah diputus (STS_REJECT): baca saja walau Medical Advisor', () => {
    const html = r(TAHAP.medicalCheck, true, '2')
    expect(html).not.toContain('>Find Disease</button>')
    expect(html).toContain('This insured has been decided; the medical data is read-only.')
  })

  it('Claim Analis dan Outstanding: tampil baca saja (12-10-2026: semua tahap)', () => {
    expect(r(TAHAP.claimAnalis, false)).toContain('Medical Advisor')
    expect(r(TAHAP.claimAnalis, false)).not.toContain('>Find Disease</button>')
    const os = r(TAHAP.outstandingClaim, false)
    expect(os).toContain('Medical Advisor')
    expect(os).not.toContain('>Find Disease</button>')
    expect(os).not.toContain('+ Add')
    expect(os).not.toContain('<textarea')
  })

  it('bolehUbahMedis: hanya Medical Check + ReasLifeMedicalAdvisor + kasus terbuka', () => {
    expect(bolehUbahMedis(klaimTahap(TAHAP.medicalCheck), [PERAN.medis])).toBe(true)
    expect(bolehUbahMedis(klaimTahap(TAHAP.medicalCheck), [PERAN.admin])).toBe(false)
    expect(bolehUbahMedis(klaimTahap(TAHAP.claimAnalis), [PERAN.medis])).toBe(false)
    expect(bolehUbahMedis(klaimTahap(TAHAP.medicalCheck, 'Resolved-Completed'), [PERAN.medis])).toBe(false)
    expect(bolehUbahMedis(klaimTahap(TAHAP.medicalCheck), undefined)).toBe(false)
  })
})

describe('RECOMMENDATION per peserta (11-10-2026)', () => {
  const r = (tahap: string, o: { boleh?: boolean; kodeStatus?: string; isian?: string; ubah?: boolean } = {}) => {
    const p = { ...pesertaUji([], o.kodeStatus ?? ''), rekomendasi: 'UJI TERSIMPAN' } as Peserta
    const k = { ...klaimTahap(tahap), peserta: [p] } as Klaim
    return en(
      <RincianPeserta
        klaim={k}
        p={p}
        onMuatUlang={async () => {}}
        bolehMedis={o.boleh ?? false}
        rekomendasi={o.isian}
        onUbahRekomendasi={o.ubah === false ? undefined : () => {}}
      />,
    )
  }

  it('Medical Advisor di Medical Check: textarea di bawah grid diagnosa, nilai tersimpan, batas 1000', () => {
    const html = r(TAHAP.medicalCheck, { boleh: true })
    expect(html).toMatch(/<span class="os__label">Recommendation<\/span><textarea[^>]*maxLength="1000"[^>]*>UJI TERSIMPAN<\/textarea>/)
    expect(html.indexOf('diagnosa-grid')).toBeLessThan(html.indexOf('Recommendation'))
    expect(r(TAHAP.medicalCheck, { boleh: true, isian: 'UJI BARU' })).toContain('>UJI BARU</textarea>')
  })

  it('baca saja: bukan Medical Advisor, Claim Analis, atau peserta diputus — teks tanpa kotak isian', () => {
    for (const html of [
      r(TAHAP.medicalCheck, { boleh: false }),
      r(TAHAP.claimAnalis, { boleh: false, ubah: false }),
      r(TAHAP.medicalCheck, { boleh: true, kodeStatus: '1' }),
    ]) {
      expect(html).toContain('<span class="os__label">Recommendation</span><p class="rincian__catatan-teks">UJI TERSIMPAN</p>')
      expect(html).not.toContain('<textarea')
    }
  })

  it('rekomendasiBerubah: hanya yang berbeda dari tersimpan', () => {
    const p1 = { ...pesertaUji(), id: 'P1', rekomendasi: 'A' } as Peserta
    const p2 = { ...pesertaUji(), id: 'P2' } as Peserta
    const k = { ...klaimTahap(TAHAP.medicalCheck), peserta: [p1, p2] } as Klaim
    expect(rekomendasiBerubah(k, { P1: 'A', P2: '' })).toEqual({})
    expect(rekomendasiBerubah(k, { P1: 'B', P2: 'C', P9: 'asing' })).toEqual({ P1: 'B', P2: 'C' })
  })

  it('ada perubahan belum disimpan: perpindahan nonaktif + petunjuk', () => {
    const html = en(
      <TombolBawahTahap onGalat={() => {}} klaim={klaimTahap(TAHAP.medicalCheck)} adaPerubahan sesudahPindah={() => {}} sesudahTutup={() => {}} />,
    )
    expect(html).toContain('<p role="note" class="alert alert--info os__bantuan">Save your changes first.</p>')
    expect(html).toMatch(/<button[^>]*disabled=""[^>]*>Send Back to Admin<\/button>/)
    expect(html).toMatch(/<button[^>]*disabled=""[^>]*>Send to Claim Analyst<\/button>/)
  })

  it('halaman: Save Medical Check memanggil simpanMedis dengan yang berubah saja', () => {
    expect(SUMBER).toContain('await simpanMedis(klaimID, { rekomendasi: berubah })')
    expect(SUMBER.split("setPesan('Progress saved.')").length - 1).toBe(5)
  })
})

// ---------------------------------------------------------------------------
// 12-10-2026: Save Adjustment / Reject Outstanding per tahap; konfirmasi perpindahan
// ---------------------------------------------------------------------------

describe('Save Adjustment dan Add tidak di Medical Check (12-10-2026)', () => {
  // Prasyarat SaveAdjustment_Act terpenuhi: peserta dipilih, baris terakhir Outstanding tanpa nomor akseptasi.
  const p = {
    ...pesertaUji(),
    baris: [{ id: 'UJI-A1', kodeStatus: '0', status: 'Outstanding', nomorAkseptasi: '' }],
  } as unknown as Peserta

  it('Save Adjustment TIDAK dirender di tahap mana pun (12-10-2026, XML 1=2) walau prasyarat terpenuhi', () => {
    for (const tahap of [TAHAP.outstandingClaim, TAHAP.medicalCheck, TAHAP.claimAnalis]) {
      const html = en(<RincianPeserta klaim={{ ...klaimTahap(tahap), peserta: [p] } as Klaim} p={p} onMuatUlang={async () => {}} />)
      expect(html, tahap).not.toContain('Save Adjustment')
    }
    expect(readFileSync(join(__dirname, 'components', 'RincianPeserta.tsx'), 'utf8')).not.toContain('simpanAdjustment(')
  })

  it('rincian di Medical Check: tanpa Save Adjustment dan tanpa Add putaran', () => {
    const html = en(
      <RincianPeserta klaim={{ ...klaimTahap(TAHAP.medicalCheck), peserta: [p] } as Klaim} p={p} onMuatUlang={async () => {}} />,
    )
    expect(html).not.toContain('Save Adjustment')
    expect(html).not.toMatch(/>Add<\/button>/)
  })
})

describe('Reject Outstanding hanya di inbox Admin (12-10-2026)', () => {
  it('bolehTolakOutstanding: tahap Outstanding + ReasLifeAdmin + kasus terbuka', () => {
    expect(bolehTolakOutstanding(klaimTahap(TAHAP.outstandingClaim), [PERAN.admin])).toBe(true)
    expect(bolehTolakOutstanding(klaimTahap(TAHAP.outstandingClaim), [PERAN.medis])).toBe(false)
    expect(bolehTolakOutstanding(klaimTahap(TAHAP.medicalCheck), [PERAN.admin])).toBe(false)
    expect(bolehTolakOutstanding(klaimTahap(TAHAP.claimAnalis), [PERAN.admin, PERAN.spv])).toBe(false)
    expect(bolehTolakOutstanding(klaimTahap(TAHAP.outstandingClaim, 'Resolved-Completed'), [PERAN.admin])).toBe(false)
    expect(bolehTolakOutstanding(klaimTahap(TAHAP.outstandingClaim), undefined)).toBe(false)
  })

  it('rincian memakai gerbang pemanggil; halaman mengisinya dari bolehTolakOutstanding', () => {
    const rincian = readFileSync(join(__dirname, 'components', 'RincianPeserta.tsx'), 'utf8')
    expect(rincian).toContain('{bolehTolakOutstanding && b.status === STATUS_OUTSTANDING && klaim.nomorKlaim ? (')
    expect(rincian).toContain('bolehTolakOutstanding = false,')
    expect(SUMBER).toContain('bolehTolakOutstanding={bolehTolakOutstanding(klaim, peran)}')
  })
})

describe('konfirmasi perpindahan Medical Check / Claim Analis (12-10-2026)', () => {
  it('kalimat per tujuan', () => {
    expect(konfirmasiPindah(TAHAP_JALUR.claimAnalis)).toEqual({
      judul: 'Send to Claim Analyst',
      tanya: 'Are you sure you want to send this claim to Claim Analyst?',
      peringatan: 'The claim moves to the Claim Analyst queue and can no longer be edited from the Medical Advisor screen.',
      ya: 'Yes, Send to Claim Analyst',
    })
    expect(konfirmasiPindah(TAHAP_JALUR.outstanding)).toEqual({
      judul: 'Send Back to Admin',
      tanya: 'Are you sure you want to send this claim back to Claim Admin?',
      peringatan: 'The claim returns to the Claim Admin queue.',
      ya: 'Yes, Send Back to Admin',
    })
    expect(konfirmasiPindah(TAHAP_JALUR.medicalCheck).tanya).toBe('Are you sure you want to send this claim back to Medical Advisor?')
  })

  it('Modal: HANYA pertanyaan + peringatan (12-10-2026, tanpa Claim No / peserta); Cancel / Yes; sibuk = Sending…', () => {
    const html = en(<KonfirmasiPindah tujuan={TAHAP_JALUR.claimAnalis} sibuk={false} onLanjut={() => {}} onTutup={() => {}} />)
    expect(html).toContain('Are you sure you want to send this claim to Claim Analyst?')
    expect(html).toContain('<p role="note" class="alert alert--warn">The claim moves to the Claim Analyst queue')
    expect(html).not.toContain('Claim No')
    expect(html).not.toContain('Insured')
    expect(html).not.toContain('<dl')
    expect(html).toContain('>Cancel<')
    expect(html).toContain('>Yes, Send to Claim Analyst</button>')
    const sibuk = en(<KonfirmasiPindah tujuan={TAHAP_JALUR.outstanding} sibuk onLanjut={() => {}} onTutup={() => {}} />)
    expect(sibuk).toMatch(/<button[^>]*disabled=""[^>]*>Sending…<\/button>/)
    expect(sibuk).not.toContain('>Cancel<')
  })

  it('setiap tombol perpindahan membuka modal dulu; satu pemanggil pindahTahap; sukses ke Inbox dengan pita', () => {
    const tahap = SUMBER.slice(SUMBER.indexOf('export function TombolBawahTahap'), SUMBER.indexOf('export function KonfirmasiKirimMedis'))
    expect(tahap).toContain('setMenunggu(t.tujuan)\n            return false')
    expect(tahap.split('pindahTahap(').length - 1).toBe(1)
    expect(tahap).toContain('if (sibukPindah || menunggu === null) return')
    expect(tahap).toContain('if (!sibukPindah) setMenunggu(null)')
    expect(SUMBER).toContain('sesudahPindah={(tujuan) => onPindah(kabarPindah(klaimID, tujuan))}')
    expect(kabarPindah('CLMLF-1', TAHAP_JALUR.claimAnalis)).toBe('Claim CLMLF-1 sent to Claim Analyst.')
    expect(kabarPindah('CLMLF-1', TAHAP_JALUR.outstanding)).toBe('Claim CLMLF-1 sent back to Claim Admin.')
  })
})

// ---------------------------------------------------------------------------
// 12-10-2026 (lanjutan): MAX disembunyikan bila kosong; kulit layar kasus
// ---------------------------------------------------------------------------

describe('MAX DATA / MAX EXPIRED / MAX CLAIM RECEIVED disembunyikan bila kosong (12-10-2026)', () => {
  const kosong = { ...pesertaUji(), tanggalMaksData: '', maksKlaimKadaluarsa: '', penandaTerimaKlaim: '', penandaTerimaKlaimAlasan: 'x' } as unknown as Peserta
  const isi = { ...kosong, tanggalMaksData: '2025-06-30', maksKlaimKadaluarsa: '90', penandaTerimaKlaim: '05/07/2025' } as Peserta
  const baca = (tahap: string, p: Peserta) =>
    en(<RincianPeserta klaim={{ ...klaimTahap(tahap), peserta: [p] } as Klaim} p={p} onMuatUlang={async () => {}} />)

  it('baca saja + kosong: ketiganya tidak dirender, tanpa ikon ⓘ; STNC DATA tetap', () => {
    const html = baca(TAHAP.medicalCheck, kosong)
    for (const l of ['Max Data Received', 'Max Expired Claim', 'Max Claim Received']) expect(html, l).not.toContain(l)
    expect(html).not.toContain('ⓘ')
    expect(html).toContain('<dt>STNC Data</dt>')
  })

  it('baca saja + terisi: ketiganya tampil dengan nilainya', () => {
    const html = baca(TAHAP.claimAnalis, isi)
    expect(html).toContain('<dt>Max Expired Claim</dt><dd>90</dd>')
    expect(html).toContain('<dt>Max Claim Received</dt><dd>05/07/2025</dd>')
    expect(html).toContain('<dt>Max Data Received</dt>')
  })

  it('Outstanding sebelum Save to RNM: MAX DATA RECEIVED tetap ISIAN walau kosong', () => {
    const k = { ...klaimTahap(TAHAP.outstandingClaim), nomorKlaim: '', peserta: [kosong] } as Klaim
    const html = en(<RincianPeserta klaim={k} p={kosong} onMuatUlang={async () => {}} onUbahTanggal={() => {}} />)
    expect(html).toContain('Max Data Received')
    expect(html).not.toContain('Max Expired Claim')
    expect(html).not.toContain('Max Claim Received')
  })
})

describe('kulit Claim Fac In di layar kasus (12-10-2026) - gaya saja', () => {
  const CSS = readFileSync(join(__dirname, 'claimlife.css'), 'utf8')
  it('section kasus memakai kelas kulit; struktur (kepala, tab, grid) tetap', () => {
    expect(SUMBER).toContain('<section className="os claimlife__akar claimlife__kasus">')
    expect(SUMBER).not.toMatch(/from '[^']*claimfacin/)
  })
  it('tab bersegmen, sub-tab rincian, tabel kartu, baris rinci putih, tanda ▸/▾ pada tombol buka yang sudah ada', () => {
    for (const p of [
      '.claimlife .claimlife__kasus .tabs {',
      '.claimlife .claimlife__kasus .tabs > .tabs__item--aktif,',
      '.claimlife .claimlife__kasus .os__rincian .tabs > .tabs__item {',
      '.claimlife .claimlife__kasus :is(.inbox__tabel, .dokumen__tabel, .diagnosa-grid__tabel, .diagnosa__tabel) thead th {',
      '.claimlife .claimlife__kasus tbody tr.os__rincian > td,',
      '.claimlife .claimlife__kasus .register__kartu,',
      '.claimlife .claimlife__kasus .os__buka::before {',
      '.claimlife .claimlife__kasus .os__buka[aria-expanded="true"]::before {',
    ])
      expect(CSS, p).toContain(p)
  })
})

// ---------------------------------------------------------------------------
// 12-10-2026 (lanjutan 4): nama tab, Register, Close Claim, catatan Claim Analyst
// ---------------------------------------------------------------------------

describe('nama tab inbox dan tombol Register (12-10-2026)', () => {
  const INBOX = readFileSync(join(__dirname, 'pages', 'InboxClaimLife.tsx'), 'utf8').replace(/\r\n/g, '\n')
  it('label tampilan Claim Admin / Medical Advisor / Claim Analyst; judul korpus tetap', () => {
    const semua = tabUntuk([PERAN.admin, PERAN.medis, PERAN.spv])
    expect(semua.map((t) => t.label)).toEqual([TAHAP.inputRegister, 'Claim Admin', 'Medical Advisor', 'Claim Analyst'])
    expect(semua.map((t) => t.judul)).toEqual([TAHAP.inputRegister, TAHAP.outstandingClaim, TAHAP.medicalCheck, TAHAP.claimAnalis])
    expect(LABEL_TAB_INBOX).toEqual({ admin: 'Claim Admin', medis: 'Medical Advisor', analis: 'Claim Analyst' })
    expect(INBOX).toContain('<span className="inbox__tab-judul">{t.label}</span>')
  })
  it('Register hanya di tab Claim Admin; Export xlsx tetap di semua tab', () => {
    expect(INBOX).toContain('{aktif === TAHAP_NOMOR.outstanding && (\n        <button type="button" className="btn btn--primary" onClick={onRegister}>')
    expect(INBOX).toContain('Export xlsx')
    // Syarat Register tertutup SEBELUM tombol Export (Export tidak bersyarat tab).
    const tutupRegister = INBOX.indexOf('{MENU.register}\n        </button>\n        )}')
    expect(tutupRegister).toBeGreaterThan(0)
    expect(INBOX.indexOf('className="btn btn--ghost"', tutupRegister)).toBeGreaterThan(tutupRegister)
  })
})

describe('Close Claim di layar Claim Analyst (12-10-2026)', () => {
  const tombolSemua = (html: string) => [...html.matchAll(/<button[^>]*class="([^"]*)"[^>]*>([^<]*)<\/button>/g)].map((m) => [m[2], m[1]])
  it('tanpa Check readiness; Close Claim merah dan PALING KANAN; jalur balik ghost di kiri', () => {
    const html = en(<TombolBawahTahap onGalat={() => {}} klaim={klaimTahap(TAHAP.claimAnalis)} sesudahPindah={() => {}} sesudahTutup={() => {}} />)
    expect(html).not.toContain('Check readiness')
    const t = tombolSemua(html)
    expect(t.map((x) => x[0])).toEqual(['Send Back to Admin', 'Send Back to Medical', 'Close Claim'])
    expect(t[2]?.[1]).toContain('btn--danger')
    expect(t[0]?.[1]).toContain('btn--ghost')
  })
  it('perubahan belum disimpan menonaktifkan Close Claim juga', () => {
    const html = en(<TombolBawahTahap onGalat={() => {}} klaim={klaimTahap(TAHAP.claimAnalis)} adaPerubahan sesudahPindah={() => {}} sesudahTutup={() => {}} />)
    expect(html).toMatch(/<button[^>]*disabled=""[^>]*>Close Claim<\/button>/)
  })
  it('Modal: pertanyaan Pega b499 + peringatan + Cancel / Yes; sibuk = Closing… tanpa tutup', () => {
    const html = en(<KonfirmasiTutup sibuk={false} onLanjut={() => {}} onTutup={() => {}} />)
    expect(KONFIRMASI_TUTUP.tanya).toBe('Are you sure want to Close Claim?')
    expect(html).toContain('<p>Are you sure want to Close Claim?</p>')
    expect(html).toContain('<p role="note" class="alert alert--warn">')
    expect(html).toContain('>Yes, Close Claim</button>')
    expect(html).toContain('>Cancel<')
    const sibuk = en(<KonfirmasiTutup sibuk onLanjut={() => {}} onTutup={() => {}} />)
    expect(sibuk).toMatch(/<button[^>]*disabled=""[^>]*>Closing…<\/button>/)
    expect(sibuk).not.toContain('>Cancel<')
  })
  it('409 berpenghalang menjadi satu kalimat di tempat galat; galat lain apa adanya; satu pemanggil tutupKlaim', () => {
    const e = new ApiFailure(409, { penghalang: [{ urutan: 1, nomorSertifikat: '001', pesan: 'Insured 001 is not accepted.' }] } as never)
    expect(kalimatPenghalangTutup(e)).toBe('The claim cannot be closed yet: Insured 001 is not accepted.')
    const lain = new ApiFailure(500, {} as never)
    expect(kalimatPenghalangTutup(lain)).toBe(lain)
    expect(SUMBER.split('await tutupKlaim(').length - 1).toBe(1)
    expect(SUMBER).not.toContain('<PanelTutupKlaim')
    expect(SUMBER).toContain('onGalat(kalimatPenghalangTutup(e))')
  })
})

describe('catatan CLAIM ANALYST (12-10-2026)', () => {
  const r = (tahap: string, ubah: boolean, isian?: string) => {
    const p = { ...pesertaUji(), catatanAnalis: 'UJI CATATAN' } as Peserta
    return en(
      <RincianPeserta
        klaim={{ ...klaimTahap(tahap), peserta: [p] } as Klaim}
        p={p}
        onMuatUlang={async () => {}}
        catatanAnalis={isian}
        onUbahCatatanAnalis={ubah ? () => {} : undefined}
      />,
    )
  }
  it('Claim Analyst + SPV: textarea (≤ 1000); tahap lain baca saja; Claim Admin tidak tampil', () => {
    expect(r(TAHAP.claimAnalis, true)).toMatch(/<span class="os__label">Claim Analyst<\/span><textarea[^>]*maxLength="1000"[^>]*>UJI CATATAN<\/textarea>/)
    expect(r(TAHAP.claimAnalis, true, 'UJI BARU')).toContain('>UJI BARU</textarea>')
    for (const t of [TAHAP.medicalCheck, TAHAP.claimAnalis]) {
      const html = r(t, false)
      expect(html, t).toContain('<span class="os__label">Claim Analyst</span><p class="rincian__catatan-teks">UJI CATATAN</p>')
      expect(html, t).not.toMatch(/CLAIM ANALYST<\/span><textarea/)
    }
    expect(r(TAHAP.outstandingClaim, false)).not.toContain('Claim Analyst')
  })
  it('bolehUbahAnalis: Claim Analis + ReasLifeSPV + kasus terbuka; catatanAnalisBerubah', () => {
    expect(bolehUbahAnalis(klaimTahap(TAHAP.claimAnalis), [PERAN.spv])).toBe(true)
    expect(bolehUbahAnalis(klaimTahap(TAHAP.claimAnalis), [PERAN.medis])).toBe(false)
    expect(bolehUbahAnalis(klaimTahap(TAHAP.medicalCheck), [PERAN.spv])).toBe(false)
    expect(bolehUbahAnalis(klaimTahap(TAHAP.claimAnalis, 'Resolved-Completed'), [PERAN.spv])).toBe(false)
    const p1 = { ...pesertaUji(), id: 'P1', catatanAnalis: 'A' } as Peserta
    const k = { ...klaimTahap(TAHAP.claimAnalis), peserta: [p1] } as Klaim
    expect(catatanAnalisBerubah(k, { P1: 'A' })).toEqual({})
    expect(catatanAnalisBerubah(k, { P1: 'B', P9: 'x' })).toEqual({ P1: 'B' })
  })
  it('Save di kaki grid (pola Save Medical), PUT /analis, Progress saved.; perubahan menahan pindah/Close', () => {
    expect(SUMBER).toContain(') : bolehUbahAnalis(klaim, peran) ? (')
    expect(SUMBER).toContain('onClick={() => void simpanCatatanAnalis()}')
    expect(SUMBER).toContain('await simpanAnalis(klaimID, { catatan: berubah })')
    expect(SUMBER).toContain('Object.keys(catatanAnalisBerubah(klaim, catatanAnalis)).length > 0')
    const api = readFileSync(join(__dirname, 'api.ts'), 'utf8')
    expect(api).toContain("/analis`, { metode: 'PUT', badan: isian })")
  })
})

describe('CLAIM ANALYST di bawah RECOMMENDATION (12-10-2026)', () => {
  const r = (tahap: string, ubah: boolean) => {
    const p = { ...pesertaUji(), catatanAnalis: 'UJI CATATAN', rekomendasi: 'UJI REKOMENDASI' } as Peserta
    return en(
      <RincianPeserta
        klaim={{ ...klaimTahap(tahap), peserta: [p] } as Klaim}
        p={p}
        onMuatUlang={async () => {}}
        onUbahCatatanAnalis={ubah ? () => {} : undefined}
      />,
    )
  }
  it('urutan: MEDICAL ADVISOR → RECOMMENDATION → CLAIM ANALYST → sub-tab Adjustment', () => {
    for (const [tahap, ubah] of [
      [TAHAP.claimAnalis, true],
      [TAHAP.medicalCheck, false],
    ] as const) {
      const html = r(tahap, ubah)
      const i = ['Medical Advisor', 'Recommendation', 'Claim Analyst', 'role="tablist"'].map((k) => html.indexOf(k))
      expect(i.every((x) => x >= 0), tahap).toBe(true)
      expect([...i].sort((a, b) => a - b), tahap).toEqual(i)
      // Kotak sendiri SESUDAH kotak MEDICAL ADVISOR (bukan di dalamnya).
      expect(html.indexOf('</fieldset>', html.indexOf('Recommendation')), tahap).toBeLessThan(html.indexOf('Claim Analyst'))
    }
  })
  it('Claim Admin: MEDICAL ADVISOR tampil, CLAIM ANALYST tetap tersembunyi', () => {
    const html = r(TAHAP.outstandingClaim, false)
    expect(html).toContain('Recommendation')
    expect(html).not.toContain('Claim Analyst')
  })
})

describe('Title Case label Claim Life (keputusan work owner 12-10-2026)', () => {
  const LABELS = readFileSync(join(__dirname, 'labels.ts'), 'utf8')
  const CSS_MODUL = readFileSync(join(__dirname, 'claimlife.css'), 'utf8')
  it('tanpa label SEMUA HURUF BESAR (kecuali singkatan PL/ICD/DOB/ID/RNM/KPR/STNC/SOB/WPC)', () => {
    const singkatan = new Set(['PL', 'ICD', 'DOB', 'ID', 'RNM', 'KPR', 'STNC', 'SOB', 'WPC', 'QP', 'QR', 'TP', 'TR', 'CSV', 'PIC'])
    const kapital = [...LABELS.matchAll(/^\s*[a-zA-Z]+: '([^']+)'/gm)]
      .map((m) => m[1] ?? '')
      .filter((v) => v.split(' ').some((w) => /^[A-Z]{3,}$/.test(w) && !singkatan.has(w)))
    expect(kapital).toEqual([])
  })
  it('contoh label rincian dan grid', () => {
    for (const l of ['Place of Death', 'Country', 'Province', 'City', 'Certificate No', 'Policy No', 'Date of Loss',
      'Claim Received Date', 'Valuation Begin Date', 'Admin Notes', 'Recommendation', 'Claim Analyst', 'Disease',
      'Medical Advisor', 'Location of Death', 'PL Number', 'Name of Insured', 'Date of Birth', 'Select All', 'Delete',
      'Group Diagnose', 'ICD Code', 'Outstanding Date'])
      expect(LABELS, l).toContain(`'${l}'`)
  })
  it('CSS modul tanpa text-transform: uppercase', () => {
    expect(CSS_MODUL.replace(/\/\*[\s\S]*?\*\//g, '')).not.toContain('text-transform: uppercase')
  })
})

describe('ejaan dan kapitalisasi teks layar (keputusan work owner 12-10-2026 "jika ada yang seperti itu lagi diperbaiki")', () => {
  const LABELS = readFileSync(join(__dirname, 'labels.ts'), 'utf8')
  it('status baris "Outstanding" (bukan OutStanding); ejaan Commission / Reinsured / Spread', async () => {
    const { kataStatusBaris } = await import('./api')
    expect(kataStatusBaris('0')).toBe('Outstanding')
    expect(LABELS).toContain("'Commission %'")
    expect(LABELS).toContain("sumReasured: 'Sum Reinsured'")
    expect(LABELS).toContain("'Premium Spread Gross', 'Premium Spread Net'")
    for (const salah of ["'Commision", "'Sum Reasured'", 'Spreaded', "'OutStanding"]) expect(LABELS, salah).not.toContain(salah)
  })
  it('layar Detail lama dan panel bersama tanpa teks Indonesia', () => {
    const teks = ['pages/KlaimLife.tsx', 'components/PanelPindahTahap.tsx', 'components/PanelTanggalKlaim.tsx']
      .map((f) => readFileSync(join(__dirname, f), 'utf8'))
      .join('\n')
      .split('\n')
      .filter((b) => !/^\s*(\/\/|\/\*|\*|\{\/\*)/.test(b))
      .join('\n')
    for (const id of ["'Gagal ", "'Memuat", "'Buka'", 'Hapus klaim', "'Memindahkan", 'Perpindahan dari', "'Menyimpan", "'Mencabut", "'Ya, hapus'", 'Pengenal klaim'])
      expect(teks, id).not.toContain(id)
  })
})
