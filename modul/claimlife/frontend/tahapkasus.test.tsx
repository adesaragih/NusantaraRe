// Inbox → layar kasus per tahap (11-10-2026): satu layar untuk Outstanding, Medical Check, dan Claim Analis; pita kuning
// "This case is at the … stage" dan tombol "Open claim detail" dihapus.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { BahasaUI } from '../../../inti/frontend/components/ui/bahasaUI'
import { TAHAP } from './labels'
import { TAHAP_JALUR, type Diagnosa, type Klaim, type Peserta } from './api'
import { bolehUbahMedis, kabarPindah, kabarTutup, rekomendasiBerubah, TombolBawahTahap } from './pages/OutstandingClaimLife'
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

  it('Medical Check: Save · Send Back to Admin · Send to Claim Analyst (MedicalCheckClaimLife); tanpa Close Claim', () => {
    const html = en(
      <TombolBawahTahap
        klaim={klaimTahap(TAHAP.medicalCheck)}
        onSimpan={() => {}}
        bolehSimpan
        sesudahPindah={() => {}}
        sesudahTutup={() => {}}
      />,
    )
    expect(tombol(html)).toEqual(['Save', 'Send Back to Admin', 'Send to Claim Analyst'])
    expect(html).not.toContain('Close Claim')
    expect(html).not.toContain('Save to RNM')
  })

  it('Medical Check tanpa izin simpan: Save nonaktif', () => {
    const html = en(
      <TombolBawahTahap klaim={klaimTahap(TAHAP.medicalCheck)} onSimpan={() => {}} sesudahPindah={() => {}} sesudahTutup={() => {}} />,
    )
    expect(html).toMatch(/<button[^>]*disabled=""[^>]*>Save<\/button>/)
  })

  it('Claim Analis: dua jalur balik + Close Claim; tanpa Save Medical', () => {
    const html = en(
      <TombolBawahTahap klaim={klaimTahap(TAHAP.claimAnalis)} onSimpan={() => {}} sesudahPindah={() => {}} sesudahTutup={() => {}} />,
    )
    const t = tombol(html)
    expect(t).toContain('Send Back to Admin')
    expect(t).toContain('Send Back to Medical')
    expect(html).toContain('Close Claim')
    expect(t).not.toContain('Save')
  })

  it('kasus tertutup: tanpa tombol', () => {
    expect(
      en(<TombolBawahTahap klaim={klaimTahap(TAHAP.claimAnalis, 'Resolved-Completed')} sesudahPindah={() => {}} sesudahTutup={() => {}} />),
    ).toBe('')
  })

  it('kabar Inbox per tujuan', () => {
    expect(kabarPindah('CLMLF-1', TAHAP_JALUR.outstanding)).toBe('Claim CLMLF-1 sent back to Admin.')
    expect(kabarPindah('CLMLF-1', TAHAP_JALUR.claimAnalis)).toBe('Claim CLMLF-1 sent to Claim Analyst.')
    expect(kabarPindah('CLMLF-1', TAHAP_JALUR.medicalCheck)).toBe('Claim CLMLF-1 sent to Medical Check.')
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
  it('dapat diubah: Find Disease | DIAGNOSE | ICD CODE | GROUP DIAGNOSE | Delete; + Add di kepala dan kaki', () => {
    const html = en(<GridDiagnosa klaimID="K1" peserta={pesertaUji([dg(1, 'UJI PENYAKIT', 'A00')])} onBerubah={() => {}} boleh />)
    expect(kepala(html)).toEqual(['Find Disease', 'DIAGNOSE', 'ICD CODE', 'GROUP DIAGNOSE', 'Delete'])
    expect(html.match(/>\+ Add<\/button>/g)).toHaveLength(2)
    expect(html).toContain('>Find Disease</button>')
    expect(html).toContain('>Delete</button>')
    // Daftar kelompok belum ada: dropdown nonaktif + keterangan English.
    expect(html).toMatch(/<select[^>]*disabled=""[^>]*title="Group Diagnose list is not configured yet."/)
    expect(html).toContain('<p class="diagnosa-grid__catatan">Group Diagnose list is not configured yet.</p>')
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
    expect(kepala(html)).toEqual(['DIAGNOSE', 'ICD CODE', 'GROUP DIAGNOSE'])
    expect(html).not.toContain('<button')
    expect(html).not.toContain('<select')
    expect(html).toContain('<td>UJI-G1</td>')
  })

  it('kosong: + Add sekali di kepala (boleh); kalimat English', () => {
    const boleh = en(<GridDiagnosa klaimID="K1" peserta={pesertaUji()} onBerubah={() => {}} boleh />)
    expect(boleh.match(/>\+ Add<\/button>/g)).toHaveLength(1)
    expect(boleh).toContain('No diagnosis yet. Click + Add to add a row.')
    expect(en(<GridDiagnosa klaimID="K1" peserta={pesertaUji()} onBerubah={() => {}} boleh={false} />)).toContain(
      'No diagnosis recorded.',
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
    expect(html).toContain('<legend class="os__label">MEDICAL ADVISOR</legend>')
    expect(html).toContain('>Find Disease</button>')
  })

  it('Medical Check bukan Medical Advisor: baca saja', () => {
    const html = r(TAHAP.medicalCheck, false)
    expect(html).toContain('MEDICAL ADVISOR')
    expect(html).not.toContain('>Find Disease</button>')
  })

  it('peserta sudah diputus (STS_REJECT): baca saja walau Medical Advisor', () => {
    const html = r(TAHAP.medicalCheck, true, '2')
    expect(html).not.toContain('>Find Disease</button>')
    expect(html).toContain('This insured has been decided; the medical data is read-only.')
  })

  it('Claim Analis: tampil baca saja; Outstanding: tidak tampil', () => {
    expect(r(TAHAP.claimAnalis, false)).toContain('MEDICAL ADVISOR')
    expect(r(TAHAP.claimAnalis, false)).not.toContain('>Find Disease</button>')
    expect(r(TAHAP.outstandingClaim, false)).not.toContain('MEDICAL ADVISOR')
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
    expect(html).toMatch(/<span class="os__label">RECOMMENDATION<\/span><textarea[^>]*maxLength="1000"[^>]*>UJI TERSIMPAN<\/textarea>/)
    expect(html.indexOf('diagnosa-grid')).toBeLessThan(html.indexOf('RECOMMENDATION'))
    expect(r(TAHAP.medicalCheck, { boleh: true, isian: 'UJI BARU' })).toContain('>UJI BARU</textarea>')
  })

  it('baca saja: bukan Medical Advisor, Claim Analis, atau peserta diputus — teks tanpa kotak isian', () => {
    for (const html of [
      r(TAHAP.medicalCheck, { boleh: false }),
      r(TAHAP.claimAnalis, { boleh: false, ubah: false }),
      r(TAHAP.medicalCheck, { boleh: true, kodeStatus: '1' }),
    ]) {
      expect(html).toContain('<span class="os__label">RECOMMENDATION</span><p class="rincian__catatan-teks">UJI TERSIMPAN</p>')
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

  it('ada perubahan belum disimpan: Save hidup, perpindahan nonaktif + petunjuk', () => {
    const html = en(
      <TombolBawahTahap
        klaim={klaimTahap(TAHAP.medicalCheck)}
        onSimpan={() => {}}
        bolehSimpan
        adaPerubahan
        sesudahPindah={() => {}}
        sesudahTutup={() => {}}
      />,
    )
    expect(html).toContain('<p role="note" class="alert alert--info os__bantuan">Save your changes first.</p>')
    expect(html).toMatch(/<button[^>]*class="btn btn--primary"[^>]*>Save<\/button>/)
    expect(html).not.toMatch(/<button[^>]*class="btn btn--primary"[^>]*disabled/)
    expect(html).toMatch(/<button[^>]*disabled=""[^>]*>Send Back to Admin<\/button>/)
    expect(html).toMatch(/<button[^>]*disabled=""[^>]*>Send to Claim Analyst<\/button>/)
  })

  it('halaman: Save Medical Check memanggil simpanMedis dengan yang berubah saja', () => {
    expect(SUMBER).toContain('await simpanMedis(klaimID, { rekomendasi: berubah })')
    expect(SUMBER).toContain("setPesan('Medical Check saved.')")
  })
})
