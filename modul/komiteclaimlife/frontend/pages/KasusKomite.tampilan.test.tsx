import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { BahasaUI } from '../../../../inti/frontend/components/ui/bahasaUI'
import { PERAN } from '../../../../inti/frontend/labels'
import type { KasusKomite, KlaimInduk, PesertaKlaim, RiwayatKomite } from '../api'
import { EFEK_KOMITE, ESKALASI_KOMITE, FORM_KOMITE, KASUS_KOMITE, KEPUTUSAN_KOMITE, KOLOM_INBOX_KOMITE } from '../labels'
import {
  kataEfekKomite,
  kataStatusKomite,
  nadaStatusKomite,
  nilaiKlaimKomite,
  tanggalKomite,
  uangKomite,
  waktuKomite,
} from '../nilai'
import {
  adaKeputusanAsli,
  bolehKirimKeputusan,
  KonfirmasiKeputusan,
  pesertaDiserahkan,
  TampilanKasusKomite,
  tampilRetro,
  teksGiliran,
  type IndukKomite,
} from './KasusKomite'

// Layar kasus komite — work owner 12-10-2026 ("perbaiki tampilannya", lalu "perbaiki tampilannya lagi, berikut referensi
// dari pega"): urutan kotak `Section/ShowTransfer.xml` di bawah "Form of Internal Claim Requisition", kulit Claim Life,
// teks Inggris, format angka & tanggal, chip status, Note wajib, modal konfirmasi, data klaim induk baca saja.

const en = (n: React.ReactNode) =>
  renderToStaticMarkup(<BahasaUI.Provider value="en">{n}</BahasaUI.Provider>).replace(/<!-- -->/g, '')

const SUMBER = readFileSync(join(__dirname, 'KasusKomite.tsx'), 'utf8')
const API = readFileSync(join(__dirname, '..', 'api.ts'), 'utf8')
const CSS = readFileSync(join(__dirname, '..', 'komiteclaimlife.css'), 'utf8')

function kasus(giliranSaya: boolean): KasusKomite {
  return {
    kasus: {
      kasusId: 'KMTLF-UJI1', tglUpdate: '', statusWork: '', klaimId: 'UJI-K', nomorKlaim: 'UJI-CLM-01',
      tingkatBerjalan: 1, komiteLoop: 1, nilaiKlaim: '2342342', mataUang: 'IDR', statusBaris: 'Outstanding',
    },
    adjustmentId: 'UJI-ADJ-2',
    tangga: [],
    giliranSaya,
    efek: { keadaan: '', efek: [] },
  }
}

function riwayat(asli = false): RiwayatKomite {
  return {
    kasusId: 'KMTLF-UJI1',
    adjustmentId: 'UJI-ADJ-2',
    tangga: [
      {
        urut: 1, committee: 'UJI-JABATAN-1', anggota: 'UJI-AKUN-1', status: 'Setuju',
        dateApprove: '2026-10-12 09:05:44', comment: 'UJI komentar panjang',
        ...(asli ? { asli: { status: 'Setuju', comment: 'UJI awal' } } : {}),
      },
      { urut: 2, committee: 'UJI-JABATAN-2', anggota: 'UJI-AKUN-2', status: 'Tolak', dateApprove: '2026-10-12 10:00:00', comment: '' },
      { urut: 3, committee: 'UJI-JABATAN-3', anggota: 'UJI-AKUN-3', status: 'Menunggu', dateApprove: '', comment: '' },
    ],
    eskalasi: [],
    penyerahan: { tanggal: '2026-10-09 15:21:07', pic: 'UJI-PIC', remarks: 'UJI baris satu\nUJI baris dua' },
  }
}

function peserta(id: string, adj: string): PesertaKlaim {
  return {
    id, idPega: `UJI-PEGA-${id}`, nomorPremiList: 'UJI-PL-01', nomorPolis: 'UJI-POL-01', nomorSertifikat: `UJI-CERT-${id}`,
    namaTertanggung: `UJI-NAMA-${id}`, umur: '41', jenisKelamin: 'M', tanggalLahir: '1985-01-02',
    tanggalMulai: '2025-01-01', tanggalExpired: '2026-12-31', tanggalEfektif: '2025-01-15',
    valuasiMulai: '2025-02-01', valuasiSelesai: '2026-01-31',
    tanggalTerimaKlaim: '2026-10-01 00:00:00', tanggalDokumenLengkap: '2026-10-05 00:00:00', tanggalKejadian: '2026-09-20',
    rekomendasi: 'UJI rekomendasi', catatanAnalis: 'UJI analis\nbaris dua',
    baris: [
      {
        id: adj,
        jumlahKlaim: { amount: '54000000', currency: 'IDR' },
        shareNusantaraRe: { amount: '60000000.5', currency: 'IDR' },
        cedingRetention: { amount: '0', currency: 'IDR' },
        sumInsured: { amount: '100000000', currency: 'IDR' },
        sumReasured: { amount: '60000000', currency: 'IDR' },
        shareRetro: { amount: '12000000', currency: 'IDR' },
        retrocededShare: { amount: '12000000', currency: 'IDR' },
        payableTo: 'UJI-PAYABLE', namaBank: 'UJI-BANK', cabangBank: 'UJI-CABANG', nomorRekening: 'UJI-REK-001',
      },
    ],
    dokumen: [{ id: '20261009152107123', namaFile: 'UJI-berkas.pdf', kategori2: 'UJI-KATEGORI', tanggal: '2026-10-09 15:21:07' }],
    diagnosa: [{ id: 1, nama: 'UJI-PENYAKIT' }],
  }
}

const KLAIM: KlaimInduk = {
  id: 'UJI-K',
  nomorKlaim: 'UJI-CLM-01',
  claimRetro: { amount: '10800000', currency: 'IDR' },
  peserta: [peserta('P1', 'UJI-ADJ-1'), peserta('P2', 'UJI-ADJ-2')],
}

function induk(o: Partial<IndukKomite> = {}): IndukKomite {
  return {
    klaim: KLAIM,
    diserahkan: pesertaDiserahkan(KLAIM, 'UJI-ADJ-2'),
    polis: { cedingCoName: 'UJI-CEDING', businessName: 'UJI-PLAN', policyHolderName: 'UJI-HOLDER', sobName: 'UJI-SOB', type: 'TR' },
    spreading: { lapisan: [{ treatyTypeName: 'UJI-QS', retrocededShare: { amount: '12000000', currency: 'IDR' } }] },
    ...o,
  }
}

const NOP = () => undefined

function layar(o: Partial<Parameters<typeof TampilanKasusKomite>[0]> = {}) {
  return en(
    <TampilanKasusKomite
      kasusID="KMTLF-UJI1"
      k={kasus(true)}
      riwayat={riwayat()}
      induk={induk()}
      memuatInduk={false}
      galatInduk={null}
      peran={[PERAN.spv]}
      sibuk={false}
      kirim={false}
      galat={null}
      kabar=""
      keputusan=""
      komentar=""
      konfirmasi={false}
      onKembali={NOP}
      onTutupGalat={NOP}
      onKeputusan={NOP}
      onKomentar={NOP}
      onMintaKonfirmasi={NOP}
      onBatalKonfirmasi={NOP}
      onSubmit={NOP}
      onEskalasi={NOP}
      onRefresh={NOP}
      onLihat={NOP}
      {...o}
    />,
  )
}

/** Isi satu kartu (`data-kartu`) — sampai kartu berikutnya. */
function kartu(h: string, id: string): string {
  const awal = h.indexOf(`data-kartu="${id}"`)
  expect(awal).toBeGreaterThan(-1)
  const sisa = h.slice(awal + 1)
  const ujung = sisa.search(/data-kartu="/)
  return ujung < 0 ? sisa : sisa.slice(0, ujung)
}

describe('format angka dan tanggal layar kasus komite', () => {
  it('uang dua desimal seperti Pega (54.000.000,00), tanpa float; tanggal dd/mm/yyyy; waktu dd/mm/yyyy HH:mm', () => {
    expect(uangKomite('54000000')).toBe('54.000.000,00')
    expect(uangKomite('60000000.5')).toBe('60.000.000,50')
    expect(uangKomite('1500000.255')).toBe('1.500.000,26')
    expect(uangKomite('0')).toBe('0,00')
    expect(uangKomite('')).toBe('—')
    expect(uangKomite(undefined)).toBe('—')
    expect(nilaiKlaimKomite('2342342', 'IDR')).toBe('2.342.342 IDR')
    expect(tanggalKomite('2026-09-20')).toBe('20/09/2026')
    expect(tanggalKomite('2026-10-01 00:00:00')).toBe('01/10/2026')
    expect(tanggalKomite(null)).toBe('—')
    expect(waktuKomite('2026-10-12 17:02:58')).toBe('12/10/2026 17:02')
    expect(waktuKomite('')).toBe('—')
    expect(SUMBER).not.toMatch(/\bNumber\(|parseFloat|parseInt/)
  })
})

describe('teks layar kasus komite berbahasa Inggris', () => {
  it('kata status server diterjemahkan; nada chip', () => {
    expect(['Menunggu', 'Setuju', 'Tolak', 'Dilewati (eskalasi)', 'Kode 9', ''].map(kataStatusKomite)).toEqual([
      'Waiting', 'Approved', 'Rejected', 'Skipped (Escalation)', 'Code 9', '—',
    ])
    expect(['Menunggu', 'Setuju', 'Tolak', 'Dilewati (eskalasi)'].map(nadaStatusKomite)).toEqual([
      'menunggu', 'setuju', 'tolak', 'lain',
    ])
    expect(['antre', 'gagal-permanen', 'perlu intervensi', 'tuntas'].map(kataEfekKomite)).toEqual([
      'Queued', 'Failed Permanently', 'Needs Intervention', 'Complete',
    ])
  })

  it('label Title Case Inggris; nol kata Indonesia', () => {
    expect(KASUS_KOMITE.kembali).toBe('Back to Claim Life Inbox')
    expect(KASUS_KOMITE.asli).toBe('Before Final Reject')
    expect(KASUS_KOMITE.tangga).toBe('List of Committee')
    expect([KEPUTUSAN_KOMITE.setuju, KEPUTUSAN_KOMITE.tolak]).toEqual(['Approve', 'Reject'])
    expect(KEPUTUSAN_KOMITE.komentar).toBe('Note')
    expect(FORM_KOMITE.sumReinsured).toBe('Sum Reinsured')
    const semua = [KOLOM_INBOX_KOMITE, KASUS_KOMITE, KEPUTUSAN_KOMITE, EFEK_KOMITE, ESKALASI_KOMITE, FORM_KOMITE].flatMap((o) =>
      Object.values(o),
    )
    for (const v of semua) {
      expect(v).not.toMatch(/\b(Tingkat|tingkat|Kasus|kasus|Klaim|klaim|Menunggu|Setuju|Tolak|Giliran|Anda|Keputusan|keputusan|Pilih|Efek|Sejak|tidak|wajib|Tanggal|Peserta)\b/)
      expect(v).not.toMatch(/^[A-Z ]{6,}$/)
    }
    expect(SUMBER).not.toMatch(/tercatat di|naik ke tingkat|Tangga berhenti|Nomor akseptasi|diantre| oleh \{| pada \{/)
  })
})

describe('urutan kartu mengikuti ShowTransfer', () => {
  it('Form of Internal Claim Requisition, lalu Policy Detail … Decision', () => {
    const h = layar()
    const urut = [...h.matchAll(/data-kartu="([a-z-]+)"/g)].map((m) => m[1])
    expect(urut).toEqual(['polis', 'klaim', 'medis', 'adjustment', 'dokumen', 'penyerahan', 'tangga', 'keputusan'])
    expect(h.indexOf('>Form of Internal Claim Requisition</h3>')).toBeLessThan(h.indexOf('data-kartu="polis"'))
    for (const [id, judul] of [
      ['polis', 'Policy Detail'], ['klaim', 'Claim Detail'], ['medis', 'Medical Advisor and Claim Analyst'],
      ['adjustment', 'Adjustment Detail'], ['dokumen', 'Document Claim'], ['tangga', 'List of Committee'],
      ['keputusan', 'Decision'],
    ] as const) {
      expect(kartu(h, id)).toContain(`>${judul}</h3>`)
    }
  })

  it('kepala: Back ghost kecil, judul Committee Case, chip giliran, ringkasan kasus', () => {
    const h = layar()
    expect(h).toContain('class="btn btn--ghost btn--sm komiteclaimlife__kembali">Back to Claim Life Inbox</button>')
    expect(h).toContain('>Committee Case KMTLF-UJI1</h2>')
    expect(h).toContain('komiteclaimlife__chip--giliran">Your turn to decide</span>')
    expect(h).toContain('<dt>Parent Claim</dt><dd>UJI-CLM-01</dd>')
    expect(h).toContain('<dt>Level</dt><dd>1 / 1</dd>')
    expect(h).toContain('<dt>Row Status</dt><dd><span class="komiteclaimlife__chip komiteclaimlife__chip--lain">Outstanding</span></dd>')
    expect(teksGiliran(kasus(false), riwayat())).toBe('Waiting for level 3')
    expect(teksGiliran(kasus(false), null)).toBe('Waiting for level 1')
    const selesai = riwayat()
    selesai.tangga = selesai.tangga.slice(0, 2)
    expect(teksGiliran(kasus(false), selesai)).toBe('Ladder Completed')
  })
})

describe('medan per kartu', () => {
  it('Policy Detail: Ceding Company, Plan, PL Number | Policy Holder Name, SOB Name', () => {
    const c = kartu(layar(), 'polis')
    for (const s of [
      '<dt>Ceding Company</dt><dd>UJI-CEDING</dd>', '<dt>Plan</dt><dd>UJI-PLAN</dd>', '<dt>PL Number</dt><dd>UJI-PL-01</dd>',
      '<dt>Policy Holder Name</dt><dd>UJI-HOLDER</dd>', '<dt>SOB Name</dt><dd>UJI-SOB</dd>',
    ]) expect(c).toContain(s)
    expect(c.indexOf('PL Number')).toBeLessThan(c.indexOf('Policy Holder Name'))
    const tanpa = kartu(layar({ induk: induk({ polis: null }) }), 'polis')
    expect(tanpa).toContain('The policy header was not found.')
    expect(tanpa).toContain('<dt>Ceding Company</dt><dd>—</dd>')
  })

  it('Claim Detail tiga kolom; HANYA peserta yang diserahkan (P2, bukan P1); tanggal dd/mm/yyyy', () => {
    const c = kartu(layar(), 'klaim')
    for (const s of [
      '<dt>Name of Insured</dt><dd>UJI-NAMA-P2</dd>', '<dt>Certificate No</dt><dd>UJI-CERT-P2</dd>',
      '<dt>Policy No</dt><dd>UJI-POL-01</dd>', '<dt>Age</dt><dd>41</dd>', '<dt>Sex</dt><dd>M</dd>',
      '<dt>Date of Birth</dt><dd>02/01/1985</dd>', '<dt>Begin Date</dt><dd>01/01/2025</dd>',
      '<dt>Expired Date</dt><dd>31/12/2026</dd>', '<dt>Effective Date</dt><dd>15/01/2025</dd>',
      '<dt>Valuation Begin Date</dt><dd>01/02/2025</dd>', '<dt>Valuation Expired Date</dt><dd>31/01/2026</dd>',
      '<dt>Claim Received Date</dt><dd>01/10/2026</dd>', '<dt>Document Complete Date</dt><dd>05/10/2026</dd>',
      '<dt>Date of Loss</dt><dd>20/09/2026</dd>',
    ]) expect(c).toContain(s)
    expect(layar()).not.toContain('UJI-NAMA-P1')
    expect((c.match(/komiteclaimlife__grid--tumpuk/g) ?? []).length).toBe(3)
  })

  it('Medical Advisor and Claim Analyst: grid Diagnose (No items bila kosong), Recommendation, Claim Analyst', () => {
    const c = kartu(layar(), 'medis')
    expect(c).toContain('<th>Diagnose</th>')
    expect(c).toContain('<td>UJI-PENYAKIT</td>')
    expect(c).toContain('<dt>Recommendation</dt><dd class="komiteclaimlife__teks-panjang">UJI rekomendasi</dd>')
    expect(c).toContain('<dt>Claim Analyst</dt><dd class="komiteclaimlife__teks-panjang">UJI analis\nbaris dua</dd>')
    const kosong = structuredClone(KLAIM)
    const p = kosong.peserta[1]
    if (p === undefined) throw new Error('fixture')
    p.diagnosa = []
    expect(kartu(layar({ induk: induk({ klaim: kosong, diserahkan: pesertaDiserahkan(kosong, 'UJI-ADJ-2') }) }), 'medis'))
      .toContain('>No items</td>')
  })

  it('Adjustment Detail: Currency; uang berformat; bank; grid Treaty Type | Retroceded Share', () => {
    const c = kartu(layar(), 'adjustment')
    for (const s of [
      '<dt>Currency</dt><dd>IDR</dd>', '<dt>Share Nusantara Re</dt><dd>60.000.000,50</dd>',
      '<dt>Ceding Retention</dt><dd>0,00</dd>', '<dt>Claim RNM</dt><dd>54.000.000,00</dd>',
      '<dt>Claim Retro</dt><dd>10.800.000,00</dd>', '<dt>Sum Reinsured</dt><dd>60.000.000,00</dd>',
      '<dt>Sum Insured</dt><dd>100.000.000,00</dd>', '<dt>Share Retro</dt><dd>12.000.000,00</dd>',
      '<dt>Retroceded Share</dt><dd>12.000.000,00</dd>', '<dt>Payable To</dt><dd>UJI-PAYABLE</dd>',
      '<dt>Name of Bank</dt><dd>UJI-BANK</dd>', '<dt>Branch of Bank</dt><dd>UJI-CABANG</dd>',
      '<dt>Account No</dt><dd>UJI-REK-001</dd>',
      '<th>Treaty Type</th><th class="komiteclaimlife__kanan">Retroceded Share</th>',
      '<td>UJI-QS</td><td class="komiteclaimlife__kanan">12.000.000,00</td>',
    ]) expect(c).toContain(s)
    const urut = ['Currency', 'Share Nusantara Re', 'Ceding Retention', 'Claim RNM', 'Claim Retro', 'Sum Reinsured',
      'Sum Insured', 'Share Retro', '<dt>Retroceded Share', 'Payable To', 'Name of Bank', 'Branch of Bank', 'Account No',
      'Treaty Type']
    const posisi = urut.map((s) => c.indexOf(s))
    expect(posisi).toEqual([...posisi].sort((a, b) => a - b))
  })

  it('Claim Retro / Share Retro / Retroceded Share hanya polis Type TP / TR (ShowTransfer b15429)', () => {
    expect(tampilRetro(null)).toBe(true)
    expect(['TP', 'tr', 'QS', ''].map((type) => tampilRetro({ cedingCoName: '', businessName: '', policyHolderName: '', sobName: '', type })))
      .toEqual([true, true, false, false])
    const polisQS = { cedingCoName: 'UJI-C', businessName: 'UJI-B', policyHolderName: 'UJI-H', sobName: 'UJI-S', type: 'QS' }
    const c = kartu(layar({ induk: induk({ polis: polisQS }) }), 'adjustment')
    expect(c).not.toContain('Claim Retro')
    expect(c).not.toContain('Share Retro')
    expect(c).not.toContain('<dt>Retroceded Share')
  })

  it('Document Claim: Refresh, File Name dapat diklik (View), Attachment Type, Date', () => {
    const c = kartu(layar(), 'dokumen')
    expect(c).toContain('>Refresh</button>')
    expect(c).toContain('<th>File Name</th><th>Attachment Type</th><th>Date</th>')
    expect(c).toContain('class="komiteclaimlife__tautan" title="View UJI-berkas.pdf">UJI-berkas.pdf</button>')
    expect(c).toContain('<td>UJI-KATEGORI</td><td class="komiteclaimlife__waktu">09/10/2026 15:21</td>')
    expect(kartu(layar({ memuatInduk: true }), 'dokumen')).toContain('disabled="">Refresh</button>')
  })

  it('Submission Date dd/mm/yyyy HH:mm dan Remarks dibungkus', () => {
    const c = kartu(layar(), 'penyerahan')
    expect(c).toContain('<dt>Submission Date</dt><dd>09/10/2026 15:21</dd>')
    expect(c).toContain('<dd class="komite-kasus__remarks">UJI baris satu\nUJI baris dua</dd>')
    expect(CSS).toMatch(/\.komite-kasus__remarks \{[^}]*white-space: pre-wrap;/)
  })

  it('List of Committee: Committee | Status | Date Approve | Comment; chip status; tanggal terformat', () => {
    const c = kartu(layar(), 'tangga')
    expect(c).toContain('<th>Committee</th><th>Status</th><th>Date Approve</th><th>Comment</th></tr>')
    expect(c).toContain('komiteclaimlife__chip--setuju">Approved</span>')
    expect(c).toContain('komiteclaimlife__chip--tolak">Rejected</span>')
    expect(c).toContain('komiteclaimlife__chip--menunggu">Waiting</span>')
    expect(c).toContain('>12/10/2026 09:05</td>')
    expect(c).toContain('<div class="komiteclaimlife__gulir"><table class="komiteclaimlife__tabel komite-kasus__tangga">')
  })

  it('Before Final Reject di List of Committee hanya bila ada tingkat yang tertimpa', () => {
    expect(adaKeputusanAsli(riwayat())).toBe(false)
    expect(layar()).not.toContain('Before Final Reject')
    expect(adaKeputusanAsli(riwayat(true))).toBe(true)
    const c = kartu(layar({ riwayat: riwayat(true) }), 'tangga')
    expect(c).toContain('<th>Before Final Reject</th>')
    expect(c).toContain('>Approved — UJI awal</td>')
  })
})

describe('data klaim induk baca saja', () => {
  it('peserta + baris yang diserahkan dipilih dari adjustment kasus', () => {
    expect(pesertaDiserahkan(KLAIM, 'UJI-ADJ-1')?.peserta.id).toBe('P1')
    expect(pesertaDiserahkan(KLAIM, 'UJI-ADJ-2')?.peserta.id).toBe('P2')
    expect(pesertaDiserahkan(KLAIM, 'UJI-ADJ-X')).toBeNull()
    expect(layar({ induk: induk({ diserahkan: null }) })).toContain('The submitted insured row was not found in the parent claim.')
  })

  it('klaim induk gagal dibaca: dinyatakan, layar komite tetap jalan', () => {
    const h = layar({ induk: null, galatInduk: 'UJI galat' })
    expect(h).toContain('The parent claim could not be loaded.')
    expect(h).toContain('data-kartu="tangga"')
    expect(layar({ induk: null, memuatInduk: true })).toContain('Loading claim data…')
  })

  it('hanya GET ke rute Claim Life / PremiumList Life; tanpa impor modul lain', () => {
    const bagian = API.slice(API.indexOf('Form of Internal Claim Requisition'))
    expect(bagian).not.toMatch(/metode|method: '(POST|PUT|PATCH|DELETE)'/)
    expect([...bagian.matchAll(/`(\/api\/[^`?]*)/g)].map((m) => m[1])).toEqual([
      '/api/klaim-life/${encodeURIComponent(klaimID)}',
      '/api/polis-life/ringkas',
      '/api/klaim-life/${encodeURIComponent(klaimID)}/peserta/${encodeURIComponent(pesertaID)}/spreading-polis',
      '/api/klaim-life/${encodeURIComponent(klaimID)}/dokumen/${encodeURIComponent(dokID)}/isi',
    ])
    for (const src of [SUMBER, API]) expect(src).not.toMatch(/from '(\.\.\/)+(claimlife|premiumlistlife)\//)
  })
})

describe('keputusan: Note wajib (ShowTransfer b32893-b32940)', () => {
  it('label VERBATIM di atas dropdown; Note bertanda wajib; Cancel kiri, Submit kanan', () => {
    const c = kartu(layar(), 'keputusan')
    expect(c).toContain('Are you sure to accept this document? <span aria-hidden="true">*</span>')
    expect(c).toContain('<option value="1">Approve</option><option value="2">Reject</option>')
    expect(c).toContain('Note <span aria-hidden="true">*</span>')
    expect(c).toMatch(/<textarea required=""/)
    expect(c.indexOf('>Cancel</button>')).toBeLessThan(c.indexOf('>Submit</button>'))
    expect(CSS).toMatch(/\.komiteclaimlife__aksi--antara \{\s*justify-content: space-between;/)
  })

  it('Submit nonaktif sampai keputusan DAN Note terisi', () => {
    expect(bolehKirimKeputusan('', '')).toBe(false)
    expect(bolehKirimKeputusan('1', '')).toBe(false)
    expect(bolehKirimKeputusan('1', '   ')).toBe(false)
    expect(bolehKirimKeputusan('', 'UJI note')).toBe(false)
    expect(bolehKirimKeputusan('2', 'UJI note')).toBe(true)
    expect(layar({ keputusan: '1' })).toContain('<button type="submit" class="btn btn--primary" disabled="">Submit</button>')
    expect(layar({ keputusan: '1', komentar: 'UJI note' })).toContain('<button type="submit" class="btn btn--primary">Submit</button>')
    expect(layar({ keputusan: '1', komentar: 'UJI note', kirim: true })).toContain('disabled="">Submitting…</button>')
  })

  it('kartu Decision hanya pada giliran', () => {
    expect(layar({ k: kasus(false) })).not.toContain('data-kartu="keputusan"')
  })

  it('galat di pita atas yang dapat ditutup', () => {
    const h = layar({ galat: 'UJI galat' })
    expect(h).toContain('<div class="komiteclaimlife__pita"><div role="alert" class="alert alert--error">UJI galat</div>')
    expect(h).toContain('aria-label="Dismiss"')
    expect(h.indexOf('komiteclaimlife__pita')).toBeLessThan(h.indexOf('Form of Internal Claim Requisition'))
  })
})

describe('modal konfirmasi keputusan', () => {
  it('pertanyaan + keputusan + Note + peringatan + Cancel / Yes', () => {
    const h = en(<KonfirmasiKeputusan keputusan="2" note="UJI note" kirim={false} onBatal={NOP} onYa={NOP} />)
    expect(h).toContain('Confirm Decision')
    expect(h).toContain('<p class="komiteclaimlife__tanya">Submit this decision?</p>')
    expect(h).toContain('<dt>Decision</dt><dd><strong>Reject</strong></dd>')
    expect(h).toContain('<dt>Note</dt><dd class="komiteclaimlife__teks-panjang">UJI note</dd>')
    expect(h).toContain('class="alert alert--warn"')
    expect(h).toContain('>Cancel</button>')
    expect(h).toContain('<button type="button" class="btn btn--primary">Yes</button>')
  })

  it('saat mengirim: Submitting…, tanpa tutup', () => {
    const h = en(<KonfirmasiKeputusan keputusan="1" note="UJI note" kirim onBatal={NOP} onYa={NOP} />)
    expect(h).toContain('disabled="">Submitting…</button>')
    expect(h).not.toContain('>Cancel</button>')
  })

  it('modal hanya dibuka lewat Submit formulir', () => {
    expect(layar({ keputusan: '1', komentar: 'UJI note', konfirmasi: true })).toContain('Submit this decision?')
    expect(SUMBER).toContain('onMintaKonfirmasi()')
  })
})
