// Layar Outstanding mengikuti Pega (10-10-2026) — `kelvin/InboxClaimLife (Done)/Section/InputOSClaimLife.xml`.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { BahasaUI } from '../../../../inti/frontend/components/ui/bahasaUI'
import { ApiFailure } from '../../../../inti/frontend/klien'
import { TAHAP } from '../labels'
import { aturanTombolOutstanding, type BarisAdjustment, type Dokumen, type Klaim, type Peserta, type SpreadingBaris } from '../api'
import {
  GridPeserta,
  isianAwal,
  KOLOM_PESERTA,
  labelTabOS,
  pilihSemuaAccept,
  TAB_OS,
  TabRegisterClaim,
  TombolBawahOutstanding,
} from './OutstandingClaimLife'
import { GridPesertaPolis } from '../components/PanelPesertaPolis'
import { DialogTolakOutstanding } from '../components/DialogTolakOutstanding'
import { kategoriTerlaluPanjang, KONFIRMASI_HAPUS_DOKUMEN, PanelDokumenPeserta, tanggalDokumen } from '../components/PanelDokumenPeserta'
import { medanPolis, medanPolisOutstanding, PanelPolisOutstanding } from '../components/PanelDataPolis'
import {
  BagianSpreading,
  galatClaimRNM,
  galatDOL,
  IsianAdjustment,
  RincianPeserta,
  tanggalAwal,
  turunanClaimRNM,
} from '../components/RincianPeserta'
import {
  adaIsianTidakSah,
  adaPerubahan,
  bolehUbahCatatanAdmin,
  bagianBerubah,
  kabarTerkirimMedis,
  kekuranganDokumenMedis,
  KonfirmasiKirimMedis,
  KonfirmasiSimpanRNM,
  kekuranganWajib,
  pesanSimpanDulu,
  saringPeserta,
} from './OutstandingClaimLife'
import { barisKekurangan, JUDUL_KEKURANGAN_MEDIS, PitaGalat } from '../components/PitaGalat'

const SUMBER = readFileSync(join(__dirname, 'OutstandingClaimLife.tsx'), 'utf8').replace(/\r\n/g, '\n')
const API = readFileSync(join(__dirname, '..', 'api.ts'), 'utf8').replace(/\r\n/g, '\n')

function pesertaUji(id: string, isAccept = ''): Peserta {
  return {
    id,
    nomorPremiList: 'UJI-PL-1',
    nomorPolis: 'UJI-POL-1',
    nomorSertifikat: `00${id}`,
    mataUang: 'IDR',
    isCheck: 'true',
    isAccept,
    namaTertanggung: `UJI NAMA ${id}`,
    tanggalLahir: '1980-02-03',
    tanggalKejadian: '',
    tanggalTerimaKlaim: '',
    tanggalDokumenLengkap: '',
    tanggalKonfirmasi: '',
    penandaTerimaKlaim: '',
    penandaTerimaKlaimAlasan: '',
    baris: [{ id: `UJI-A${id}`, kodeStatus: '', status: 'Outstanding' }],
    dokumen: [],
    diagnosa: [],
    kodeStatus: '',
  } as unknown as Peserta
}

function klaimUji(o: { nomorKlaim?: string; type?: string; tersimpan?: boolean; katastrofe?: string } = {}): Klaim {
  const peserta = [pesertaUji('1', 'true'), pesertaUji('2')]
  if (o.tersimpan === true) for (const p of peserta) p.baris = p.baris.map((b) => ({ ...b, kodeStatus: '0' }))
  return {
    id: 'CLMLF-UJI1',
    nomorKlaim: o.nomorKlaim ?? '',
    nomorPolis: '',
    namaBisnis: '',
    kodeStatus: '',
    statusTurunan: '',
    peserta,
    cacahBaris: 2,
    tahap: TAHAP.outstandingClaim,
    statusWork: '',
    type: o.type ?? 'QP',
    stsKatastrofe: o.katastrofe ?? '',
    isKpr: '',
    stncClaim: '',
    katastrofeNote: '',
  } as unknown as Klaim
}

const en = (n: React.ReactNode) => renderToStaticMarkup(<BahasaUI.Provider value="en">{n}</BahasaUI.Provider>)
const nonaktif = (html: string, label: string) => new RegExp(`<button[^>]*disabled=""[^>]*>${label}</button>`).test(html)
const ada = (html: string, label: string) => new RegExp(`<button[^>]*>${label}</button>`).test(html)

describe('kerangka: kepala dan tab (b1607)', () => {
  it('tiga tab urut Pega: Register Claim, Policy Details, Claim Details', () => {
    expect(TAB_OS.map(labelTabOS)).toEqual(['Register Claim', 'Policy Details', 'Claim Details'])
  })
  it('judul "Claim Life (<Case ID>)" dengan Claim No atau tanda pisah di bawahnya; tombol Detail klaim lama hilang', () => {
    expect(SUMBER).toContain('{OS_PEGA.judulKasus} ({klaimID})')
    expect(SUMBER).toContain('<span className="os__label">Claim No</span>')
    expect(SUMBER).toContain("<strong>{klaim !== null && bernomor(klaim) ? klaim.nomorKlaim : '—'}</strong>")
    expect(SUMBER).not.toContain('Detail klaim\n')
    expect(SUMBER).not.toContain('className="os__detail"')
  })
})

describe('grid Participant Details (b15924)', () => {
  it('kolom urut Pega; nama dan tanggal lahir dd/mm/yyyy tampil', () => {
    expect(KOLOM_PESERTA).toEqual(['No.', 'PL Number', 'Name of Insured', 'Date of Birth', 'Accept', 'Delete'])
    const html = en(
      <GridPeserta
        klaim={klaimUji()}
        accept={{ '1': 'true', '2': 'false' }}
        onAccept={() => {}}
        onPilihSemua={() => {}}
        terbuka={null}
        onBuka={() => {}}
        onHapus={() => {}}
        boleh
        rincian={() => null}
      />,
    )
    for (const k of ['<th>PL Number</th>', '<th>Name of Insured</th>', '<th>Date of Birth</th>', '>Select All</button>', '<th>Delete</th>'])
      expect(html).toContain(k)
    expect(html).toContain('<td>UJI NAMA 1</td>')
    expect(html).toContain('<td>03/02/1980</td>')
    // Accept peserta 1 tercentang, peserta 2 tidak.
    expect(html.match(/type="checkbox" checked=""/g) ?? []).toHaveLength(1)
  })

  it('SELECT ALL menjungkit SEMUA peserta (SelectAllClaimLife_act b247/b429)', () => {
    const a = pilihSemuaAccept('', { '1': 'false', '2': 'true' })
    expect(a).toEqual({ keadaan: 'true', accept: { '1': 'true', '2': 'true' } })
    const b = pilihSemuaAccept(a.keadaan, a.accept)
    expect(b).toEqual({ keadaan: 'false', accept: { '1': 'false', '2': 'false' } })
  })

  it('DELETE hanya saat CLAIM_NO kosong (b18082), dan memanggil cabut sesudah konfirmasi', () => {
    const bernomor = en(
      <GridPeserta
        klaim={klaimUji({ nomorKlaim: 'UJI-CLM-1', tersimpan: true })}
        accept={{}}
        onAccept={() => {}}
        onPilihSemua={() => {}}
        terbuka={null}
        onBuka={() => {}}
        onHapus={() => {}}
        boleh
        rincian={() => null}
      />,
    )
    expect(bernomor).not.toContain('Delete')
    expect(SUMBER).toContain('await cabutPeserta(klaimID, pesertaID)')
    // Klik DELETE hanya membuka konfirmasi; cabut hanya dari tombol modal.
    expect(SUMBER).toContain('onHapus={setHapus}')
    expect(SUMBER).toContain('onClick={() => void jalankanHapus(hapus)}')
    expect(SUMBER).toContain('judul={OS_PEGA.konfirmasiHapus}')
  })

  it('Save menyimpan IsAccept: isian awal dari isAccept, PUT /outstanding/isian', () => {
    expect(isianAwal(klaimUji()).accept).toEqual({ '1': 'true', '2': 'false' })
    expect(isianAwal(klaimUji()).stsKatastrofe).toBe('Non-Catastrophe')
    expect(SUMBER).toContain('await simpanIsianOutstanding(klaimID, isian)')
    expect(SUMBER).toContain('ubahIsian({ ...isian, accept: { ...isian.accept, [id]: v } })')
    expect(API).toContain("/outstanding/isian`, { metode: 'PUT', badan: isian })")
  })

  it('baris dibuka menampilkan rincian peserta (ViewClaimDetailLifeGCNM b18252) lewat komponen bersama', () => {
    const html = en(
      <GridPeserta
        klaim={klaimUji()}
        accept={{}}
        onAccept={() => {}}
        onPilihSemua={() => {}}
        terbuka="2"
        onBuka={() => {}}
        onHapus={() => {}}
        boleh
        rincian={(id) => <div>UJI-RINCIAN-{id}</div>}
      />,
    )
    expect(html).toContain('UJI-RINCIAN-2')
    expect(html).not.toContain('UJI-RINCIAN-1')
    expect(html).toContain('aria-expanded="true"')
    expect(SUMBER).toContain('<RincianPeserta\n                      klaim={klaim}\n                      p={p}\n                      onMuatUlang={muatUlang}')
  })
})

describe('tombol bawah (b21102, b21404, b21839, b22808)', () => {
  const render = (k: Klaim) =>
    en(<TombolBawahOutstanding klaim={k} sibukRNM={false} onSimpanRNM={() => {}} sesudahPindah={() => {}} />)

  it('belum tersimpan, QP: Save to RNM hidup, Send to Medical tampil tapi nonaktif; tiga tombol lain TIDAK ada (11-10-2026)', () => {
    const html = render(klaimUji())
    expect(ada(html, 'Save to RNM') && !nonaktif(html, 'Save to RNM')).toBe(true)
    expect(nonaktif(html, 'Send to Medical Check')).toBe(true)
    expect(html).not.toContain('Send Back to Register')
    expect(html).not.toContain('Close Claim')
    expect(html).not.toContain('Check readiness')
    // Close Claim hanya di tombol tahap Claim Analis (TombolBawahTahap, 11-10-2026), tidak di tombol Outstanding.
    const bawahOS = SUMBER.slice(SUMBER.indexOf('export function TombolBawahOutstanding'), SUMBER.indexOf('export const KONFIRMASI_RNM'))
    expect(bawahOS).not.toContain('<PanelTutupKlaim')
  })

  it('sesudah Save to RNM: Save to RNM nonaktif, Send Back hilang, Send to Medical hidup', () => {
    const html = render(klaimUji({ nomorKlaim: 'UJI-CLM-1', tersimpan: true }))
    expect(nonaktif(html, 'Save to RNM')).toBe(true)
    expect(ada(html, 'Send Back to Register')).toBe(false)
    expect(ada(html, 'Send to Medical Check') && !nonaktif(html, 'Send to Medical Check')).toBe(true)
  })

  it('Type TP/TR: Send to Medical Check tidak tampil (b22123)', () => {
    expect(ada(render(klaimUji({ type: 'TP', nomorKlaim: 'UJI-CLM-1' })), 'Send to Medical Check')).toBe(false)
    expect(aturanTombolOutstanding(klaimUji({ type: 'TR' })).kirimMedisTampil).toBe(false)
  })
})

describe('tab Register Claim (b3475, b3781, b4726, b4905)', () => {
  it('Catastrophe terkunci sesudah bernomor; Note tampil hanya bila Catastrophe', () => {
    const k = klaimUji({ nomorKlaim: 'UJI-CLM-1', katastrofe: 'Catastrophe' })
    const html = en(<TabRegisterClaim klaim={k} isian={isianAwal(k)} onUbah={() => {}} boleh />)
    // 11-10-2026: radiogroup (bukan fieldset); terkunci = setiap radio disabled.
    expect(html).toMatch(/role="radiogroup" aria-label="Catastrophe"><dt>Catastrophe<\/dt><dd class="os__radio-datar">/)
    expect(html).toMatch(/name="stsKatastrofe"[^>]*disabled=""/)
    expect(html).toContain('Catastrophe Note')
    expect(html).toContain('readonly=""')
    const k2 = klaimUji()
    const html2 = en(<TabRegisterClaim klaim={k2} isian={isianAwal(k2)} onUbah={() => {}} boleh />)
    expect(html2).not.toContain('Catastrophe Note')
    expect(html2).not.toMatch(/name="stsKatastrofe"[^>]*disabled=""/)
    expect(html2).toContain('Is KPR')
    // Keputusan work owner 11-10-2026: pilihan KPR / Non-KPR.
    expect(html2).toContain('value="KPR"')
    expect(html2).toContain('value="Non-KPR"')
    expect(html2).not.toContain('value="Yes"')
    expect(html2).toContain('STNC Treaty')
  })
})

describe('Is KPR warisan (11-10-2026)', () => {
  it('nilai di luar KPR / Non-KPR ditampilkan apa adanya, tidak ada radio yang dicentang paksa', () => {
    const k = { ...klaimUji(), isKpr: 'Y' } as Klaim
    const html = en(<TabRegisterClaim klaim={k} isian={isianAwal(k)} onUbah={() => {}} boleh />)
    expect(html).toContain('Stored value: Y')
    expect(html).not.toMatch(/name="isKpr" value="[^"]*" checked=""/)
  })
})

describe('View Detail peserta polis (DetailPolisLife b13267, 11-10-2026)', () => {
  const hal = {
    halaman: 2,
    ukuran: 500,
    total: 1200,
    baris: [
      {
        urut: 501, nomorPolis: 'UJI-POL-1', nomorSertifikat: 'UJI-001', namaTertanggung: 'UJI NAMA POLIS', jenisKelamin: 'M',
        tanggalLahir: '1980-02-03', umur: '45', plan: 'UJI-PLAN', tanggalMulai: '2025-01-01', tanggalEfektif: '2025-01-02',
        tanggalExpired: '2026-01-01',
      },
    ],
  }
  it('kolom VERBATIM XML, tanggal dd/mm/yyyy, paginasi Halaman inti 500 per halaman', () => {
    const html = en(<GridPesertaPolis hal={hal} onPindah={() => {}} />)
    for (const k of ['No.', 'Policy No', 'Certificate No', 'Name of Insured', 'Sex', 'Date of Birth', 'Age', 'Plan',
      'Begin Date', 'Effective Date', 'Expired Date'])
      expect(html).toContain(`<th>${k}</th>`)
    expect(html).toContain('<td>501</td><td>UJI-POL-1</td><td>UJI-001</td><td>UJI NAMA POLIS</td><td>M</td><td>03/02/1980</td><td>45</td><td>UJI-PLAN</td><td>01/01/2025</td><td>02/01/2025</td><td>01/01/2026</td>')
    expect(html).toMatch(/Page\s*(<!-- -->)?2(<!-- -->)?\s*of\s*(<!-- -->)?3/)
  })
  it('halaman kosong dinyatakan; tab Policy Details TIDAK lagi memasang panelnya; api tetap ada', () => {
    expect(en(<GridPesertaPolis hal={{ ...hal, baris: [], total: 0 }} onPindah={() => {}} />)).toContain('No insured in this policy version.')
    // 11-10-2026: View Detail dicabut dari tab Policy Details (keputusan work owner).
    expect(SUMBER).not.toContain('<PanelPesertaPolis')
    expect(API).toContain('/peserta-polis?halaman=')
  })
})

describe('tab Policy Details Outstanding (screenshot 11-10-2026)', () => {
  const polis = {
    nomorPolis: 'UJI-POL-1', type: 'QR', marketingName: 'UJI-MO', cedingCoName: 'UJI-CED', policyHolderName: 'UJI-PH',
    businessName: 'UJI-COB', dateReceived: '', status: 'Resolved-Completed', statusUpdate: '', productNameId: 'P1',
    productName: 'UJI-PROD', prodKe: 0, medanTanpaSumber: [], typeCeding: 'TC1', typeCedingName: 'UJI-SYSREAS',
    proRateType: 'UJI-PRORATE', wpc: '', retroName: '', securityReinsurer: '', sobName: 'UJI-SOB',
  } as unknown as Parameters<typeof medanPolisOutstanding>[0]
  it('tiga kolom urut Pega; System Reinsurance dan Premium Method TAMPIL di Outstanding', () => {
    const k = medanPolisOutstanding(polis, 'UJI-PL-1')
    expect(k.map((c) => c.map((m) => m.label))).toEqual([
      ['PL Number', 'Type', 'System Reinsurance', 'Premium Method', 'Marketing Officer', 'WPC'],
      ['Ceding', 'Policy Holder', 'SOB', 'Class of Business', 'Product Name ID', 'Product Name'],
      ['Date Received Email', 'Response Date', 'Confirmation Date', 'Status', 'Confirmation Reserved', 'Updated Status',
        'Realization Date', 'Underwriter Note'],
    ])
    expect(k[0]?.[2]?.nilai).toBe('UJI-SYSREAS')
    expect(k[0]?.[3]?.nilai).toBe('UJI-PRORATE')
  })
  it('panel Register TETAP tanpa System Reinsurance / Premium Method (keputusan 10-10-2026)', () => {
    const label = medanPolis({ polis }).map((m) => m.label)
    expect(label).not.toContain('System Reinsurance')
    expect(label).not.toContain('Premium Method')
  })
})

describe('grid bernomor, baris terpilih, saringan (11-10-2026)', () => {
  it('nomor 1..n, baris terbuka tersorot, rincian DI BAWAH baris itu saja', () => {
    const html = en(
      <GridPeserta
        klaim={klaimUji()}
        accept={{}}
        onAccept={() => {}}
        onPilihSemua={() => {}}
        terbuka="2"
        onBuka={() => {}}
        onHapus={() => {}}
        boleh
        rincian={(id) => <div>UJI-RINCIAN-{id}</div>}
        kakiKanan={<button type="button">Save</button>}
      />,
    )
    expect(html).toContain('<th>No.</th>')
    expect(html).toMatch(/aria-expanded="false">1<\/button>/)
    expect(html).toMatch(/aria-expanded="true">2<\/button>/)
    expect(html).toMatch(/<tr class="inbox__baris os__baris--pilih">[\s\S]*UJI NAMA 2[\s\S]*<\/tr><tr class="os__rincian"><td colSpan="6"><div>UJI-RINCIAN-2/)
    expect(html.match(/os__baris--pilih/g) ?? []).toHaveLength(1)
    expect(html).toContain('class="os__simpan os__simpan--kanan"')
    expect(html).toContain('class="os__tautan"')
  })
  it('saringan nama tidak peka huruf', () => {
    const k = klaimUji()
    expect(saringPeserta(k.peserta, 'nama 2').map((p) => p.id)).toEqual(['2'])
    expect(saringPeserta(k.peserta, '').map((p) => p.id)).toEqual(['1', '2'])
  })
})

describe('rincian peserta tiga kolom + sub-tab (screenshot ke-4, 11-10-2026)', () => {
  const p = {
    ...pesertaUji('1'),
    plan: 'UJI-PLAN', umur: '45', jenisKelamin: 'M', tanggalMulai: '2025-01-01', tanggalExpired: '2026-01-01',
    tanggalEfektif: '2025-01-02', valuasiMulai: '2025-01-03', valuasiSelesai: '2025-12-31', stnc: '',
    maksKlaimKadaluarsa: '90', catatanAnalis: 'UJI ANALIS', total: {},
    // 12-10-2026: tiga MAX disembunyikan bila kosong - fixture urutan diberi nilai supaya urutannya tetap teruji.
    tanggalMaksData: '2025-06-30', penandaTerimaKlaim: '05/07/2025',
  } as unknown as Peserta
  it('medan per kolom urut screenshot; LOCATION OF DEATH baca saja; sub-tab Adjustment | Document', () => {
    const html = en(<RincianPeserta klaim={klaimUji()} p={p} onMuatUlang={async () => {}} catatanAdmin="UJI CATATAN" onUbahCatatan={() => {}} />)
    const urut = (labels: string[]) => {
      let i = -1
      for (const l of labels) {
        const j = html.indexOf(`<dt>${l}</dt>`)
        expect(j, l).toBeGreaterThan(i)
        i = j
      }
    }
    // Tata letak 11-10-2026: kelompok Insured | Coverage | Claim Dates (bukan urutan Pega 100%).
    urut(['Certificate No', 'Policy No', 'Plan', 'Age', 'Sex', 'DOB'])
    urut(['DOB', 'Begin Date', 'Expired Date', 'Effective Date', 'Valuation Begin Date', 'Valuation Expired Date', 'Received Date', 'STNC Treaty'])
    urut(['STNC Treaty', 'Date of Loss', 'Claim Received Date', 'Confirmation Date', 'Document Complete Date', 'Max Data Received', 'Max Expired Claim', 'Max Claim Received', 'STNC Data'])
    expect(html).toContain('Location of Death')
    expect(html).toContain('<dt>Place of Death</dt><dd>—</dd>')
    expect(html).toContain('>UJI CATATAN</textarea>')
    // CLAIM ANALYST hanya di tahap Claim Analis (11-10-2026).
    expect(html).not.toContain('Claim Analyst')
    const analis = en(<RincianPeserta klaim={{ ...klaimUji(), tahap: TAHAP.claimAnalis } as Klaim} p={p} onMuatUlang={async () => {}} />)
    expect(analis).toContain('UJI ANALIS')
    expect(html).toContain('>Edit Date</button>')
    expect(html).toMatch(/role="tab"[^>]*>Adjustment<\/button>[\s\S]*role="tab"[^>]*>Document<\/button>/)
    for (const k of ['Status', 'Outstanding Date', 'Acceptation No', 'Acceptation Date']) expect(html).toContain(`<th>${k}</th>`)
    // Diagnosa hanya di tahap Medical Check (b4908).
    expect(html).not.toContain('Find Disease')
  })
})

describe('isian adjustment tiga kolom + bank + Spreading (11-10-2026)', () => {
  const b = {
    id: 'UJI-A1', kodeStatus: '', status: 'Tidak diketahui', jumlahKlaim: { amount: '94695.705', currency: 'IDR' },
    shareNusantaraRe: { amount: '100000', currency: 'IDR' }, cedingRetention: { amount: '1234567.5', currency: 'IDR' },
    sumInsured: { amount: '200000000', currency: 'IDR' }, sumReasured: { amount: '150000000', currency: 'IDR' },
    persenKlaim: '94.695705', klaimDibayar: { amount: '94695.705', currency: 'IDR' },
    payableTo: 'UJI PAYABLE', namaBank: 'UJI BANK', cabangBank: 'UJI CABANG', nomorRekening: '000123',
    spreading: { lapisan: [], rekap: [] }, tanggalOutstanding: '',
  } as unknown as Parameters<typeof IsianAdjustment>[0]['b']
  it('sebelum Save to RNM: tiga kolom, angka berformat, Claim RNM dan bank dapat diisi, TANPA tombol Save to Outstanding', () => {
    const html = en(<IsianAdjustment klaim={klaimUji()} b={b} onUbah={() => {}} />)
    const urut = (labels: string[]) => {
      let i = -1
      for (const l of labels) {
        const j = html.indexOf(l)
        expect(j, l).toBeGreaterThan(i)
        i = j
      }
    }
    urut(['Currency', 'Percent Claim (%)', 'Claim Paid', 'Share Nusantara Re', 'Ceding Retention', 'Claim RNM', 'Sum Reinsured', 'Sum Insured'])
    expect(html).toContain('<dt>Claim Paid</dt><dd class="os__angka">94.695,705</dd>')
    expect(html).toContain('<dt>Ceding Retention</dt><dd class="os__angka">1.234.567,5</dd>')
    // Dihitung langsung dari isian (aturan backend): 94.695,705 / 100.000 -> pecahan 0,94696 -> 94,696%.
    expect(html).toContain('<dt>Percent Claim (%)</dt><dd class="os__angka">94,696%</dd>')
    expect(html).toContain('<dt>Payable To</dt><dd>UJI PAYABLE</dd>')
    for (const l of ['Claim RNM', 'Name of Bank']) expect(html).toContain(`aria-label="${l}"`)
    // "semua read only kecuali name of bank" (11-10-2026): Branch of Bank / Account No baca saja.
    for (const l of ['Branch of Bank', 'Account No']) expect(html).not.toContain(`aria-label="${l}"`)
    expect(html).toContain('class="field__input os__isian-angka"')
    expect(html).toContain('class="field__input os__isian-teks"')
    expect(html).not.toContain('Save to Outstanding')
    // Spreading menggantikan grid Treaty Type dan grid komite.
    // Spreading belum termuat (rute PremiumList) -> Memuat; pesan kosong diuji di BagianSpreading.
    expect(html).not.toContain('It is calculated when Claim RNM is saved')
    expect(html).not.toContain('Committee Name')
    expect(html).not.toContain('No items')
  })
  it('sesudah Save to RNM / tanpa onUbah: baca saja dengan angka berformat', () => {
    const html = en(<IsianAdjustment klaim={klaimUji({ nomorKlaim: 'UJI-CLM-1' })} b={{ ...b, kodeStatus: '0' }} onUbah={() => {}} />)
    expect(html).not.toContain('aria-label="Claim RNM"')
    expect(html).toContain('<span class="os__angka">94.695,705</span>')
    expect(html).toContain('<span>UJI BANK</span>')
    expect(en(<IsianAdjustment klaim={klaimUji()} b={b} />)).not.toContain('aria-label="Claim RNM"')
  })
  it('Spreading per lapisan (Reins Type sekali per lapisan) dan rekap per reinsurer, angka rata kanan', () => {
    const m = (v: string) => ({ amount: v, currency: 'IDR' })
    const r = (v: string) => ({ value: v })
    const sp = {
      lapisan: [
        { treatyTypeName: 'QS', treatyYearLife: '', retrocededShare: m('1000.5'), reinsurer: [
          { nama: 'UJI RE A', pctShare: r('60'), amount: m('600.3'), rate: r('1.5'), komisi: r('10'), ovrComm: r('2'), gross: m('100'), net: m('90') },
          { nama: 'UJI RE B', pctShare: r('40'), amount: m('400.2'), rate: r('1.5'), komisi: r('10'), ovrComm: r('2'), gross: m('50'), net: m('45') },
        ] },
      ],
      rekap: [{ reinsurer: 'UJI RE A', currency: 'IDR', gross: m('1200.5'), net: m('108') }],
    }
    const html = en(<IsianAdjustment klaim={klaimUji()} b={b} spreading={sp as unknown as SpreadingBaris} spreadingTerbuka />)
    for (const k of ['Reins Type', 'Retroceded Share', 'Reinsurer', '% Share', 'Amount', 'Rate', 'Commission %', 'Ovr Comm %', 'Gross', 'Net',
      'Premium Spread Gross', 'Premium Spread Net']) expect(html).toContain(k)
    expect(html).toContain('<tr class="os__lapisan-awal"><td>QS</td><td class="os__angka">1.000,5</td><td>UJI RE A</td><td class="os__angka">60</td><td class="os__angka">600,3</td>')
    expect(html).toContain('<tr><td></td><td class="os__angka"></td><td>UJI RE B</td>')
    expect(html).toContain('<td>UJI RE A</td><td>IDR</td><td class="os__angka">1.200,5</td><td class="os__angka">108</td>')
    expect(html).not.toContain('No spreading yet')
  })
  it('Save grid mengirim Claim RNM + bank; total peserta tidak di Outstanding', () => {
    expect(SUMBER).toContain('onUbahAdjustment={')
    expect(SUMBER).toContain('ubahIsian({ ...isian, adjustment: { ...isian.adjustment, [adjID]: isi } })')
    expect(SUMBER).toContain('await simpanIsianOutstanding(klaimID, isian)')
    expect(SUMBER).not.toContain('tanpaTotal')
    expect(API).not.toContain('/adjustment/${encodeURIComponent(adjID)}/isian')
    const pT = { ...pesertaUji('1'), total: {}, baris: [] } as unknown as Peserta
    expect(en(<RincianPeserta klaim={klaimUji()} p={pT} onMuatUlang={async () => {}} />)).not.toContain('Total Claim Amount')
  })
})

describe('grid adjustment, Reject, tanggal isian (11-10-2026)', () => {
  const pAdj = (kode: string, tgl: string) =>
    ({ ...pesertaUji('1'), total: {}, baris: [{ id: 'UJI-A1', kodeStatus: kode, status: '', nomorAkseptasi: '', tanggalAkseptasi: '', tanggalOutstanding: tgl, jumlahKlaim: {} }] }) as unknown as Peserta
  it('Status teks dan OutStanding Date dari jejak; "—" sebelum Save to RNM', () => {
    const sebelum = en(<RincianPeserta klaim={klaimUji()} p={pAdj('', '')} onMuatUlang={async () => {}} />)
    expect(sebelum).toMatch(/<tr class="inbox__baris"><td>—<\/td><td>—<\/td><td>—<\/td><td>—<\/td><\/tr>/)
    const sesudah = en(<RincianPeserta klaim={klaimUji({ nomorKlaim: 'UJI-CLM-1' })} p={pAdj('0', '2026-10-11 09:00:00')} onMuatUlang={async () => {}} />)
    expect(sesudah).toMatch(/<tr class="inbox__baris"><td>Outstanding<\/td><td>11\/10\/2026<\/td>/)
    for (const k of ['<th>Status</th>', '<th>Outstanding Date</th>', '<th>Acceptation No</th>', '<th>Acceptation Date</th>']) expect(sesudah).toContain(k)
    expect(sesudah).not.toContain('<th>Action</th>')
  })
  it('Reject Outstanding di kaki isian, hanya sesudah bernomor dan baris berstatus 0 (b15399)', () => {
    const sumber = readFileSync(join(__dirname, '..', 'components', 'RincianPeserta.tsx'), 'utf8')
    expect(sumber).toContain('{bolehTolakOutstanding && b.status === STATUS_OUTSTANDING && klaim.nomorKlaim ? (')
    expect(sumber).toContain('aksi={')
  })
  it('lima tanggal: isian sebelum Save to RNM, baca saja sesudahnya; Edit Date tidak ada di Outstanding', () => {
    const t = { kejadian: '2026-01-02', terimaKlaim: '2026-01-03', dokumenLengkap: '', maksData: '', konfirmasi: '' }
    const sebelum = en(<RincianPeserta klaim={klaimUji()} p={{ ...pesertaUji('1'), total: {} } as unknown as Peserta} onMuatUlang={async () => {}} tanggal={t} onUbahTanggal={() => {}} />)
    for (const l of ['Date of Loss', 'Claim Received Date', 'Document Complete Date', 'Max Data Received', 'Confirmation Date'])
      expect(sebelum).toContain(`aria-label="${l}"`)
    expect(sebelum).toContain('value="2026-01-02"')
    // Wajib di Save to RNM: DATE OF LOSS, CLAIM RECEIVED DATE, CONFIRMATION DATE.
    expect(sebelum.match(/class="os__wajib"/g) ?? []).toHaveLength(3)
    expect(sebelum).not.toContain('>Edit Date</button>')
    const sesudah = en(<RincianPeserta klaim={klaimUji({ nomorKlaim: 'UJI-CLM-1' })} p={{ ...pesertaUji('1'), total: {} } as unknown as Peserta} onMuatUlang={async () => {}} tanggal={t} onUbahTanggal={() => {}} />)
    expect(sesudah).not.toContain('aria-label="Date of Loss"')
    expect(sesudah).toContain('<span>02/01/2026</span>')
  })
})

describe('spreading dari PremiumList Life; total peserta dibuang dari semua layar (11-10-2026)', () => {
  it('kosong -> pesan baru; null -> memuat; galat dinyatakan', () => {
    expect(en(<BagianSpreading s={{ lapisan: [], rekap: [] }} />)).toContain('No spreading data for this insured in the policy.')
    expect(en(<BagianSpreading s={null} />)).not.toContain('No spreading data')
    expect(en(<BagianSpreading s={null} galat="UJI galat" />)).toContain('UJI galat')
  })
  it('rincian memuat spreading lewat rute baca PremiumList; tidak ada PanelTotalPeserta di modul', () => {
    const rincian = readFileSync(join(__dirname, '..', 'components', 'RincianPeserta.tsx'), 'utf8')
    expect(rincian).toContain('ambilSpreadingPolis(klaim.id, p.id)')
    expect(rincian).not.toContain('<PanelTotalPeserta')
    expect(API).toContain('/spreading-polis`')
    for (const berkas of ['KlaimLife.tsx', 'OutstandingClaimLife.tsx'])
      expect(readFileSync(join(__dirname, berkas), 'utf8')).not.toContain('PanelTotalPeserta')
  })
})

describe('konsistensi isian, Claim RNM, Register Claim, Spreading (11-10-2026)', () => {
  const b = {
    id: 'UJI-A1', kodeStatus: '', status: '', jumlahKlaim: { amount: '6000000', currency: 'IDR' },
    shareNusantaraRe: { amount: '5077429.5', currency: 'IDR' }, persenKlaim: '', klaimDibayar: { amount: '6000000', currency: 'IDR' },
  } as unknown as Parameters<typeof IsianAdjustment>[0]['b']
  it('Claim RNM > Share: pesan merah dengan share berformat, isian bertanda error', () => {
    const html = en(<IsianAdjustment klaim={klaimUji()} b={b} isi={{ jumlah: '6000000', namaBank: '' }} onUbah={() => {}} />)
    expect(html).toContain('Claim RNM cannot be greater than Share Nusantara Re (5.077.429,5).')
    expect(html).toContain('field__input--error')
    expect(html).toContain('aria-invalid="true"')
  })
  it('bukan angka / sah', () => {
    // 11-10-2026: kosong / 0 = wajib (tidak menghalangi Save grid; ditegakkan Save to RNM).
    expect(galatClaimRNM('', b.shareNusantaraRe)).toBe('Claim RNM is required.')
    expect(galatClaimRNM('0', b.shareNusantaraRe)).toBe('Claim RNM is required.')
    expect(galatClaimRNM('0,00', b.shareNusantaraRe)).toBe('Claim RNM is required.')
    expect(galatClaimRNM('.', b.shareNusantaraRe)).toBe('Claim RNM must be a number.')
    expect(galatClaimRNM('5077429,5', b.shareNusantaraRe)).toBeNull()
    expect(galatClaimRNM('5077429.51', b.shareNusantaraRe)).toContain('cannot be greater')
    const sumber = readFileSync(join(__dirname, '..', 'components', 'RincianPeserta.tsx'), 'utf8')
    expect(sumber).toContain('saringKetikanAngka(e.target.value)')
    expect(sumber).toContain('inputMode={mode}')
  })
  it('Save grid: Claim RNM > Share TIDAK memblokir (simpan progres, 11-10-2026); hanya format rusak yang memblokir', () => {
    const k = { ...klaimUji(), peserta: [{ ...pesertaUji('1'), baris: [b] }] } as unknown as Klaim
    const isian = { ...isianAwal(k), adjustment: { 'UJI-A1': { jumlah: '6000000', namaBank: '' } } }
    expect(adaIsianTidakSah(k, isian)).toBe(false)
    expect(adaIsianTidakSah(k, { ...isian, adjustment: { 'UJI-A1': { jumlah: '.', namaBank: '' } } })).toBe(true)
    expect(adaIsianTidakSah(k, { ...isian, adjustment: { 'UJI-A1': { ...isian.adjustment['UJI-A1'], jumlah: '100' } } })).toBe(false)
    expect(SUMBER).toContain('disabled={sibuk || adaIsianTidakSah(klaim, isian)}')
  })
  it('Register Claim: tata letak dl/dt/dd seperti Policy Details, nilai tidak nowrap, tanpa teks catatan', () => {
    const k = { ...klaimUji(), stncClaim: '2026-01-01' } as Klaim
    const html = en(<TabRegisterClaim klaim={k} isian={isianAwal(k)} onUbah={() => {}} boleh />)
    expect(html).toContain('<section class="os__form os__form--dua" aria-label="Register Claim"><dl class="os__kolom">')
    expect(html).toContain('<dt>PL Number</dt><dd class="os__nilai-kode">UJI-PL-1</dd>')
    expect(html).toContain('<dt>STNC Treaty</dt>')
    expect(html).not.toContain('Saved with the Save button')
  })
  it('semua isian Outstanding memakai field__input (tanggal, Admin Notes, saringan, Catastrophe Note)', () => {
    const rincian = readFileSync(join(__dirname, '..', 'components', 'RincianPeserta.tsx'), 'utf8')
    expect(rincian).toContain("className={'field__input os__isian-tanggal' + (merah ? ' field__input--error' : '')}")
    expect(rincian).toContain('className="field__input os__isian-area"')
    expect(SUMBER).toContain('className="field__input os__isian-teks"')
    expect(SUMBER).toContain('className="field__input os__isian-area"')
  })
  it('Spreading Detail / Spreading Summary + baris TOTAL per mata uang', () => {
    const m = (v: string) => ({ amount: v, currency: 'IDR' })
    const html = en(
      <BagianSpreading
        s={{
          lapisan: [{ treatyTypeName: 'QS', treatyYearLife: '', retrocededShare: m('10'), reinsurer: [] }],
          rekap: [{ reinsurer: 'UJI RE A', currency: 'IDR', gross: m('100'), net: m('90') }],
          total: [{ currency: 'IDR', gross: m('1000.5'), net: m('900') }],
        } as unknown as SpreadingBaris}
      />,
    )
    expect(html.indexOf('Spreading Detail')).toBeLessThan(html.indexOf('Spreading Summary'))
    expect(html).toContain('<tr class="os__total"><th scope="row">Total</th><td>IDR</td><td class="os__angka">1.000,5</td><td class="os__angka">900</td></tr>')
  })
})

describe('bawaan Register, Catastrophe Note, Percent/Claim Paid langsung (11-10-2026)', () => {
  it('kolom kosong: bawaan Non-Catastrophe dan Non-KPR terpilih dan ikut dikirim; warisan asing tetap', () => {
    const awal = isianAwal(klaimUji())
    expect(awal.stsKatastrofe).toBe('Non-Catastrophe')
    expect(awal.isKpr).toBe('Non-KPR')
    expect(isianAwal({ ...klaimUji(), isKpr: 'Y' } as Klaim).isKpr).toBe('Y')
    const html = en(<TabRegisterClaim klaim={klaimUji()} isian={awal} onUbah={() => {}} boleh />)
    expect(html).toMatch(/name="stsKatastrofe"[^>]*checked=""[^>]*value="Non-Catastrophe"/)
    expect(html).toMatch(/name="isKpr"[^>]*checked=""[^>]*value="Non-KPR"/)
  })
  it('Catastrophe Note di kolom kanan, tepat di bawah STNC Treaty', () => {
    const k = klaimUji({ katastrofe: 'Catastrophe' })
    const html = en(<TabRegisterClaim klaim={k} isian={isianAwal(k)} onUbah={() => {}} boleh />)
    const kanan = html.slice(html.lastIndexOf('<dl class="os__kolom">'))
    expect(kanan.indexOf('STNC Treaty')).toBeGreaterThan(-1)
    expect(kanan.indexOf('Catastrophe Note')).toBeGreaterThan(kanan.indexOf('STNC Treaty'))
    const kiri = html.slice(0, html.lastIndexOf('<dl class="os__kolom">'))
    expect(kiri).not.toContain('Catastrophe Note')
  })
  it('Percent / Claim Paid berubah saat mengetik; tidak sah -> "—"', () => {
    const share = { amount: '5850038.75', currency: 'IDR' }
    expect(turunanClaimRNM('665', share)).toEqual({ persen: '0,011%', dibayar: '665' })
    expect(turunanClaimRNM('100000', share)).toEqual({ persen: '1,709%', dibayar: '100.000' })
    expect(turunanClaimRNM('', share)).toEqual({ persen: '—', dibayar: '—' })
    expect(turunanClaimRNM('6000000', share)).toEqual({ persen: '—', dibayar: '—' })
  })
  it('kolom kiri ikut kolom nilai 20ch yang sama (os__form--angka): Percent/Claim Paid sejajar dengan Share/Sum', () => {
    const sumber = readFileSync(join(__dirname, '..', 'components', 'RincianPeserta.tsx'), 'utf8')
    expect(sumber).toContain('<div className="os__form os__form--tiga os__form--angka">')
    expect(sumber).toContain('<Baris label={A.persen} nilai={turunan.persen} angka />')
  })
})

describe('Save to RNM: kolom wajib semua peserta (keputusan work owner 11-10-2026)', () => {
  const uang = (v: string) => ({ amount: v, currency: 'IDR' })
  const barisUji = (id: string, jumlah: string) =>
    ({ id: 'UJI-A' + id, kodeStatus: '', status: '', jumlahKlaim: uang(jumlah), shareNusantaraRe: uang('1000') }) as unknown as BarisAdjustment
  // Peserta LENGKAP: lolos pre-check.
  const lengkap = (id: string): Peserta =>
    ({
      ...pesertaUji(id, 'true'),
      tanggalKejadian: '2025-06-01',
      tanggalTerimaKlaim: '2025-06-02',
      tanggalKonfirmasi: '2025-06-03',
      tanggalMulai: '2025-01-01',
      tanggalExpired: '2026-01-01',
      dokumen: [{ id: 'UJI-D' + id }],
      baris: [barisUji(id, '100')],
    }) as unknown as Peserta
  const klaimLengkap = (...p: Peserta[]) => ({ ...klaimUji(), isKpr: 'Non-KPR', stsKatastrofe: 'Non-Catastrophe', peserta: p }) as Klaim

  it('pre-check mengumpulkan SEMUA kekurangan semua peserta; data lengkap = kosong', () => {
    expect(kekuranganWajib(klaimLengkap(lengkap('1'), lengkap('2')), isianAwal(klaimLengkap(lengkap('1'), lengkap('2'))))).toEqual([])
    const p1 = { ...lengkap('1'), tanggalKejadian: '', tanggalTerimaKlaim: '' } as Peserta
    p1.baris = [barisUji('1', '0')]
    const p2 = { ...lengkap('2'), tanggalKonfirmasi: '', dokumen: [] } as unknown as Peserta
    const k = { ...klaimLengkap(p1, p2), isKpr: '' } as Klaim
    const isian = { ...isianAwal(k), isKpr: '' }
    const kurang = kekuranganWajib(k, isian)
    expect(kurang.map((x) => [x.nomorPeserta, x.kolom])).toEqual([
      [0, 'Is KPR'],
      [1, 'Date of Loss'],
      [1, 'Claim Received Date'],
      [1, 'Claim RNM'],
      [2, 'Document'],
      [2, 'Confirmation Date'],
    ])
    expect(barisKekurangan(kurang)).toEqual([
      'Register Claim: Is KPR',
      'No 1 · UJI NAMA 1: Date of Loss, Claim Received Date, Claim RNM',
      'No 2 · UJI NAMA 2: Document, Confirmation Date',
    ])
    // Isian di layar yang belum tersimpan ikut dibaca (Claim RNM diketik).
    const terisi = { ...isian, adjustment: { 'UJI-A1': { jumlah: '50', namaBank: '' } } }
    expect(kekuranganWajib(k, terisi).some((x) => x.kolom === 'Claim RNM')).toBe(false)
  })

  it('SATU kotak galat berjudul, satu baris per peserta', () => {
    const html = en(
      <PitaGalat
        kekurangan={[
          { nomorPeserta: 1, pesertaId: '1', sertifikat: '001', nama: 'UJI NAMA 1', kolom: 'Date of Loss', pesan: '' },
          { nomorPeserta: 1, pesertaId: '1', sertifikat: '001', nama: 'UJI NAMA 1', kolom: 'Claim RNM', pesan: '' },
          { nomorPeserta: 2, pesertaId: '2', sertifikat: '002', nama: '', kolom: 'Claim RNM', pesan: '' },
        ]}
      />,
    )
    expect(html.match(/role="alert"/g)).toHaveLength(1)
    expect(html).toContain('<strong>Please complete the required fields before Save to RNM:</strong>')
    expect(html).toContain('<li>No 1 · UJI NAMA 1: Date of Loss, Claim RNM</li><li>No 2 · 002: Claim RNM</li>')
  })

  it('kolom yang kurang bertanda merah; Claim RNM bertanda * dengan pesan wajib di bawahnya', () => {
    const p = lengkap('1')
    const t = { ...tanggalAwal(p), kejadian: '' }
    const html = en(
      <RincianPeserta klaim={klaimLengkap(p)} p={p} onMuatUlang={async () => {}} tanggal={t} onUbahTanggal={() => {}}
        kolomKurang={new Set(['Date of Loss'])} />,
    )
    expect(html).toMatch(/class="field__input os__isian-tanggal field__input--error"[^>]*aria-label="Date of Loss"[^>]*aria-invalid="true"/)
    expect(html.match(/field__input--error/g)).toHaveLength(1)
    const isi = en(<IsianAdjustment klaim={klaimLengkap(p)} b={barisUji('1', '100')} isi={{ jumlah: '', namaBank: '' }} onUbah={() => {}} />)
    expect(isi).toContain('Claim RNM<span class="os__wajib" aria-hidden="true">*</span>')
    expect(isi).toContain('<span role="alert" class="os__galat-isian">Claim RNM is required.</span>')
  })

  it('Claim RNM kosong TIDAK menghalangi Save grid', () => {
    const k = klaimLengkap(lengkap('1'))
    const isian = { ...isianAwal(k), adjustment: { 'UJI-A1': { jumlah: '', namaBank: '' } } }
    expect(adaIsianTidakSah(k, isian)).toBe(false)
  })

  it('perubahan belum disimpan: Save to RNM nonaktif dengan "Save your changes first."', () => {
    const p = lengkap('1')
    const k = klaimLengkap(p)
    expect(adaPerubahan(k, isianAwal(k))).toBe(false)
    const ubahTgl = { ...isianAwal(k), tanggal: { '1': { ...tanggalAwal(p), konfirmasi: '2025-07-01' } } }
    expect(adaPerubahan(k, ubahTgl)).toBe(true)
    // Diubah lalu dikembalikan = tidak ada perubahan.
    expect(adaPerubahan(k, { ...isianAwal(k), tanggal: { '1': tanggalAwal(p) } })).toBe(false)
    expect(adaPerubahan(k, { ...isianAwal(k), adjustment: { 'UJI-A1': { jumlah: '7', namaBank: '' } } })).toBe(true)
    // Bawaan Non-KPR klaim lama (kolom tersimpan kosong) belum tersimpan.
    const lama = { ...k, isKpr: '' } as Klaim
    expect(adaPerubahan(lama, isianAwal(lama))).toBe(true)
    const html = en(<TombolBawahOutstanding klaim={k} sibukRNM={false} onSimpanRNM={() => {}} sesudahPindah={() => {}} adaPerubahan />)
    expect(nonaktif(html, 'Save to RNM')).toBe(true)
    // Teks bantuan DI ATAS deretan tombol, bukan menempel di tombol (11-10-2026).
    expect(html).toContain('<p role="note" class="alert alert--info os__bantuan">Save your changes first.</p><div class="os__aksi">')
    const bersih = en(<TombolBawahOutstanding klaim={k} sibukRNM={false} onSimpanRNM={() => {}} sesudahPindah={() => {}} adaPerubahan={false} />)
    expect(nonaktif(bersih, 'Save to RNM')).toBe(false)
    expect(bersih).not.toContain('Save your changes first.')
    expect(SUMBER).toContain('adaPerubahan={adaPerubahan(klaim, isian)}')
  })

  it('pre-check menolak SEBELUM permintaan; 422 berkekurangan ke kotak yang sama; tanpa alert ganda', () => {
    const fungsi = SUMBER.slice(SUMBER.indexOf('async function simpanRNM()'))
    const cek = fungsi.indexOf('const kurang = kekuranganWajib(klaim, isian)')
    const keluar = fungsi.indexOf('tunjukkanKekurangan(kurang)\n      return')
    const panggil = fungsi.indexOf('simpanKeRNM(klaimID)')
    expect(cek).toBeGreaterThan(-1)
    expect(keluar).toBeGreaterThan(cek)
    expect(panggil).toBeGreaterThan(keluar)
    expect(fungsi).toContain('if (e instanceof GagalSimpanRNM) tunjukkanKekurangan(e.kekurangan)')
    // Peserta bermasalah pertama dibuka di tab Claim Details.
    expect(SUMBER).toContain("setTab('klaim')\n      setTerbuka(pertama.pesertaId)")
    // Gagal (di dalam PitaGalat) sudah menampilkan pesan server: tidak ada <p> salinan.
    expect(SUMBER).not.toContain('pesanGalat(')
    expect(SUMBER).not.toContain('<Gagal')
  })
})

describe('DOL dalam periode valuation, alur Save to RNM, galat tunggal, rincian simetris (11-10-2026)', () => {
  const pDol = (o: Partial<Record<string, string>> = {}) =>
    ({
      ...pesertaUji('1', 'true'),
      valuasiMulai: '2025-03-18 00:00:00',
      valuasiSelesai: '2026-03-18 00:00:00',
      valuasiRetroMulai: '2025-04-01',
      valuasiRetroSelesai: '2025-10-01',
      ...o,
    }) as unknown as Peserta
  const RENTANG = 'Date of Loss must be within the valuation period (18/03/2025 – 18/03/2026).'

  it('galatDOL = backend DOLDalamJendela: [BEGIN, EXPIRED] inklusif tanpa geser, gross dan retro; Type lain = server', () => {
    // Kasus batas SAMA dengan Go TestBatasJendelaTiapType / TestSimpanRNMDOLJendelaTanpaGeserRetro (keputusan WO 11-10-2026).
    const RETRO = 'Date of Loss must be within the valuation period (01/04/2025 – 01/10/2025).'
    for (const tipe of ['QR', 'QP']) {
      const k = klaimUji({ type: tipe })
      expect(galatDOL(k, pDol(), '2025-03-18'), tipe + ' BEGIN').toBeNull()
      expect(galatDOL(k, pDol(), '2026-03-18'), tipe + ' EXPIRED').toBeNull()
      expect(galatDOL(k, pDol(), '2025-03-17'), tipe + ' BEGIN-1').toBe(RENTANG)
      expect(galatDOL(k, pDol(), '2026-03-19'), tipe + ' EXPIRED+1').toBe(RENTANG)
    }
    for (const tipe of ['TR', 'TP']) {
      const k = klaimUji({ type: tipe })
      expect(galatDOL(k, pDol(), '2025-04-01'), tipe + ' BEGIN').toBeNull()
      expect(galatDOL(k, pDol(), '2025-10-01'), tipe + ' EXPIRED').toBeNull()
      expect(galatDOL(k, pDol(), '2025-03-31'), tipe + ' BEGIN-1').toBe(RETRO)
      expect(galatDOL(k, pDol(), '2025-10-02'), tipe + ' EXPIRED+1').toBe(RETRO)
    }
    expect(galatDOL(klaimUji({ type: 'QR' }), pDol(), '2026-10-11')).toBe(RENTANG)
    expect(galatDOL(klaimUji({ type: 'QR' }), pDol(), '')).toBeNull() // kosong = kewajiban Save to RNM
    expect(galatDOL(klaimUji({ type: '' }), pDol(), '2030-01-01')).toBeNull()
  })

  it('DOL di luar: pesan berentang di bawah isian, isian merah, Save grid TETAP aktif (simpan progres)', () => {
    const p = pDol()
    const k = { ...klaimUji({ type: 'QR' }), peserta: [p] } as Klaim
    const html = en(
      <RincianPeserta klaim={k} p={p} onMuatUlang={async () => {}} tanggal={{ ...tanggalAwal(p), kejadian: '2026-10-11' }} onUbahTanggal={() => {}} />,
    )
    expect(html).toMatch(/class="field__input os__isian-tanggal field__input--error"[^>]*aria-label="Date of Loss"[^>]*aria-invalid="true"/)
    expect(html).toContain('<span role="alert" class="os__galat-isian">' + RENTANG + '</span>')
    const isian = { ...isianAwal(k), tanggal: { '1': { ...tanggalAwal(p), kejadian: '2026-10-11' } } }
    // 11-10-2026 (terbaru): peringatan tetap tampil, Save tidak diblokir; Save to RNM yang menolak.
    expect(adaIsianTidakSah(k, isian)).toBe(false)
    expect(kekuranganWajib(k, { ...isian, isKpr: 'KPR' }).find((x) => x.kolom === 'Date of Loss')?.pesan).toBe('Insured number 1: ' + RENTANG)
    expect(adaIsianTidakSah(k, { ...isian, tanggal: { '1': { ...tanggalAwal(p), kejadian: '2025-06-01' } } })).toBe(false)
  })

  it('Save to RNM tidak menjalankan pre-check selama ada perubahan; isian berubah / Save grid membersihkan kotak lama', () => {
    const fungsi = SUMBER.slice(SUMBER.indexOf('async function simpanRNM()'))
    expect(fungsi.indexOf('adaPerubahan(klaim, isian)) return')).toBeGreaterThan(-1)
    expect(fungsi.indexOf('adaPerubahan(klaim, isian)) return')).toBeLessThan(fungsi.indexOf('kekuranganWajib(klaim, isian)'))
    const ubah = SUMBER.slice(SUMBER.indexOf('function ubahIsian('), SUMBER.indexOf('const muatUlang'))
    expect(ubah).toContain('setKekurangan(null)')
    expect(SUMBER).not.toContain('setIsian({ ...isian,')
    const simpan = SUMBER.slice(SUMBER.indexOf('async function simpan()'), SUMBER.indexOf('function tunjukkanKekurangan'))
    expect(simpan.indexOf('setKekurangan(null)')).toBeLessThan(simpan.indexOf('simpanIsianOutstanding('))
    expect(simpan).toContain('setGalat(e)')
  })

  it('SATU tempat galat di atas isi (sebelum tab), tidak ada alert di bawah tombol', () => {
    const tempat = SUMBER.indexOf('<div className="cl-tempat-galat">')
    expect(tempat).toBeGreaterThan(-1)
    expect(tempat).toBeLessThan(SUMBER.indexOf('<StripTab tab={TAB_OS}'))
    expect(SUMBER.match(/<PitaGalat/g)).toHaveLength(2) // galat muat + galat aksi/kekurangan, satu tempat
    const bawah = SUMBER.slice(SUMBER.indexOf('<TombolBawahOutstanding\n'), SUMBER.indexOf('{hapus !== null && ('))
    expect(bawah).not.toMatch(/alert|PitaGalat/)
    for (const berkas of ['RegisterKlaim.tsx', 'InboxClaimLife.tsx', 'KlaimLife.tsx']) {
      const isi = readFileSync(join(__dirname, berkas), 'utf8')
      expect(isi, berkas).toContain('<PitaGalat galat={galat')
      expect(isi, berkas).not.toMatch(/<p[^>]*role="alert"/)
    }
  })

  it('PitaGalat: dapat ditutup; kalimat; galat server SATU kali', () => {
    const teks = en(<PitaGalat galat="UJI GALAT" onTutup={() => {}} />)
    expect(teks).toContain('<div role="alert" class="alert alert--error cl-pita__isi">UJI GALAT</div>')
    expect(teks).toContain('aria-label="Dismiss"')
    expect(en(<PitaGalat galat="" />)).toBe('')
    expect(en(<PitaGalat galat={null} kekurangan={[]} />)).toBe('')
    const server = en(<PitaGalat galat={new ApiFailure(422, { code: 'DITOLAK_BACKEND', message: 'UJI SERVER' })} onTutup={() => {}} />)
    expect(server.match(/UJI SERVER/g)).toHaveLength(1)
  })

  it('rincian: tiga kelompok berjudul Insured | Coverage | Claim Dates; Location of Death lalu Admin Notes di kolom Insured', () => {
    const p = pDol()
    const html = en(<RincianPeserta klaim={klaimUji()} p={p} onMuatUlang={async () => {}} catatanAdmin="UJI" onUbahCatatan={() => {}} />)
    const i = html.indexOf('<section class="rincian__kelompok" aria-label="Insured"><h4 class="rincian__judul">Insured</h4>')
    const c = html.indexOf('<section class="rincian__kelompok" aria-label="Coverage">')
    const d = html.indexOf('<section class="rincian__kelompok" aria-label="Claim Dates">')
    expect(i).toBeGreaterThan(-1)
    expect(c).toBeGreaterThan(i)
    expect(d).toBeGreaterThan(c)
    // Susunan lama (11-10-2026): di kolom Insured, sesudah DOB, LOCATION OF DEATH lalu ADMIN NOTES di bawahnya.
    const insured = html.slice(i, c)
    expect(insured.indexOf('<dt>DOB</dt>')).toBeLessThan(insured.indexOf('Location of Death'))
    expect(insured.indexOf('Location of Death')).toBeLessThan(insured.indexOf('Admin Notes'))
    expect(html).not.toContain('rincian__bawah')
    expect(html).toContain('<dt>Certificate No</dt><dd class="rincian__nilai-kode">')
    const css = readFileSync(join(__dirname, '..', 'claimlife.css'), 'utf8').replace(/\s+/g, ' ')
    expect(css).toContain('.claimlife .rincian__grid { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr));')
    expect(css).toMatch(/@media \(max-width: 1280px\) \{ \.claimlife \.rincian__grid \{ grid-template-columns: repeat\(2, minmax\(0, 1fr\)\); \}/)
    expect(css).toContain('.claimlife .rincian__grid .os__baris-form { grid-template-columns: 24ch minmax(0, 1fr); }')
    expect(css).toContain('.claimlife .rincian__nilai-kode { white-space: normal; overflow-wrap: anywhere; }')
  })
})

describe('rincian tanpa gulir; Admin Notes hanya Admin + Outstanding (11-10-2026)', () => {
  it('penjaga gaya: tidak ada overflow auto/scroll pada aturan rincian', () => {
    const css = readFileSync(join(__dirname, '..', 'claimlife.css'), 'utf8')
    const aturan = [...css.matchAll(/([^{}]+)\{([^}]*)\}/g)].filter((m) => (m[1] ?? '').includes('rincian'))
    expect(aturan.length).toBeGreaterThan(5)
    for (const m of aturan) {
      expect(m[2], (m[1] ?? '').trim()).not.toMatch(/overflow(-x|-y)?\s*:\s*(auto|scroll)/)
    }
  })

  it('Admin Notes dapat diisi hanya Admin di tahap Outstanding; selain itu teks baca saja tanpa kotak isian', () => {
    const k = klaimUji()
    expect(bolehUbahCatatanAdmin(k, ['ReasLifeAdmin'])).toBe(true)
    expect(bolehUbahCatatanAdmin(k, ['ReasLifeSPV'])).toBe(false)
    expect(bolehUbahCatatanAdmin(k, undefined)).toBe(false)
    expect(bolehUbahCatatanAdmin({ ...k, tahap: TAHAP.claimAnalis } as Klaim, ['ReasLifeAdmin'])).toBe(false)
    expect(bolehUbahCatatanAdmin({ ...k, tahap: TAHAP.medicalCheck } as Klaim, ['ReasLifeAdmin', 'ReasLifeMedicalAdvisor'])).toBe(false)
    // Sesudah Save to RNM (bernomor) TETAP dapat diisi: XML b2588/b2628 hanya mengunci peserta diputus.
    expect(bolehUbahCatatanAdmin({ ...k, nomorKlaim: 'UJI-CLM-1' } as Klaim, ['ReasLifeAdmin'])).toBe(true)
    expect(SUMBER).toContain('bolehUbahCatatanAdmin(klaim, peran)')
    const rute = readFileSync(join(__dirname, '..', 'rute.tsx'), 'utf8')
    expect(rute).toMatch(/<OutstandingClaimLife\s+klaimID=\{kasus\}\s+peran=\{masuk\.peran\}/)

    const p = { ...pesertaUji('1'), catatanAdmin: 'UJI CATATAN TERSIMPAN' } as unknown as Peserta
    const baca = en(<RincianPeserta klaim={{ ...k, tahap: TAHAP.claimAnalis } as Klaim} p={p} onMuatUlang={async () => {}} />)
    expect(baca).not.toContain('<textarea')
    expect(baca).toContain('<p class="rincian__catatan-teks">UJI CATATAN TERSIMPAN</p>')
    const kosong = en(<RincianPeserta klaim={k} p={{ ...p, catatanAdmin: '' } as Peserta} onMuatUlang={async () => {}} />)
    expect(kosong).toContain('<p class="rincian__catatan-teks">—</p>')
    const isi = en(<RincianPeserta klaim={k} p={p} onMuatUlang={async () => {}} onUbahCatatan={() => {}} />)
    expect(isi).toContain('>UJI CATATAN TERSIMPAN</textarea>')
    // Peserta diputus (b2628): baca saja walau Admin.
    const putus = en(<RincianPeserta klaim={k} p={{ ...p, kodeStatus: '1' } as Peserta} onMuatUlang={async () => {}} onUbahCatatan={() => {}} />)
    expect(putus).not.toContain('<textarea')
  })
})

describe('Save = simpan progres, bank, spreading terlipat, tombol bawah, Register (11-10-2026)', () => {
  const uang = (v: string) => ({ amount: v, currency: 'IDR' })
  const bUji = {
    id: 'UJI-A1', kodeStatus: '', status: '', jumlahKlaim: uang('100'), shareNusantaraRe: uang('1000'),
    namaBank: 'UJI BANK', cabangBank: 'UJI CABANG', nomorRekening: '000123', payableTo: 'UJI PAYABLE',
  } as unknown as BarisAdjustment

  it('bank: hanya Name of Bank isian; Payable To / Branch / Account baca saja ("—" bila kosong)', () => {
    const html = en(<IsianAdjustment klaim={klaimUji()} b={bUji} isi={{ jumlah: '100', namaBank: 'UJI BANK' }} onUbah={() => {}} />)
    expect(html).toContain('aria-label="Name of Bank"')
    expect(html).toContain('<dt>Branch of Bank</dt><dd>UJI CABANG</dd>')
    expect(html).toContain('<dt>Account No</dt><dd>000123</dd>')
    expect(html).toContain('<dt>Payable To</dt><dd>UJI PAYABLE</dd>')
    const kosong = en(<IsianAdjustment klaim={klaimUji()} b={{ ...bUji, cabangBank: '', nomorRekening: '' }} isi={{ jumlah: '1', namaBank: '' }} onUbah={() => {}} />)
    expect(kosong).toContain('<dt>Branch of Bank</dt><dd>—</dd>')
    expect(kosong).toContain('<dt>Account No</dt><dd>—</dd>')
    expect(API).toMatch(/export interface IsianBarisAdjustment \{\s+jumlah: string[\s\S]*?namaBank: string\s+\}/)
  })

  it('spreading TERLIPAT bawaan; dimuat SEKALI saat pertama dibuka', () => {
    const html = en(<IsianAdjustment klaim={klaimUji()} b={bUji} />)
    expect(html).toContain('<button type="button" class="btn btn--ghost btn--sm os__lipat" aria-expanded="false">View Spreading ▸</button>')
    expect(html).not.toContain('Spreading Detail')
    const buka = en(<IsianAdjustment klaim={klaimUji()} b={bUji} spreading={null} spreadingTerbuka />)
    expect(buka).toContain('Hide Spreading ▾')
    const rincian = readFileSync(join(__dirname, '..', 'components', 'RincianPeserta.tsx'), 'utf8').replace(/\r\n/g, '\n')
    expect(rincian).toContain('if (!mintaSpreading) return')
    expect(rincian).toContain('onMuatSpreading={() => setMintaSpreading(true)}')
    expect(rincian).toContain('if (!lihatSpreading) onMuatSpreading?.()')
  })

  it('Save to RNM pre-check mengumpulkan Claim RNM > Share dengan kalimat server', () => {
    const p = { ...pesertaUji('1', 'true'), tanggalKejadian: '2025-06-01', tanggalTerimaKlaim: '2025-06-02', tanggalKonfirmasi: '2025-06-03',
      tanggalMulai: '2025-01-01', tanggalExpired: '2026-01-01', dokumen: [{ id: 'D' }],
      baris: [{ ...bUji, jumlahKlaim: uang('5000') }] } as unknown as Peserta
    const k = { ...klaimUji({ type: '' }), isKpr: 'KPR', peserta: [p] } as Klaim
    const kurang = kekuranganWajib(k, isianAwal(k))
    expect(kurang.map((x) => x.pesan)).toEqual(['Insured number 1: Claim RNM cannot be greater than Share Nusantara Re (1.000).'])
  })

  it('Save hanya menyimpan progres: pita "Progress saved."; tidak memanggil Save to RNM / spreading', () => {
    const simpan = SUMBER.slice(SUMBER.indexOf('async function simpan()'), SUMBER.indexOf('function tunjukkanKekurangan'))
    expect(simpan).toContain("setPesan('Progress saved.')")
    for (const dilarang of ['simpanKeRNM', 'ambilSpreadingPolis', 'kekuranganWajib']) expect(simpan).not.toContain(dilarang)
    expect(SUMBER).toContain('disabled={sibuk || adaIsianTidakSah(klaim, isian)}')
  })

  it('tombol bawah HANYA di tab Claim Details; petunjuk menyebut bagian yang berubah', () => {
    const klaimTab = SUMBER.slice(SUMBER.indexOf("{tab === 'klaim' && ("), SUMBER.indexOf('{hapus !== null && ('))
    expect(klaimTab).toContain('<TombolBawahOutstanding')
    expect(SUMBER.match(/<TombolBawahOutstanding\n/g)).toHaveLength(1)
    const luar = SUMBER.slice(SUMBER.indexOf('<StripTab tab={TAB_OS}'), SUMBER.indexOf("{tab === 'klaim' && ("))
    expect(luar).not.toContain('TombolBawahOutstanding')
    const k = { ...klaimUji(), isKpr: 'Non-KPR', stsKatastrofe: 'Non-Catastrophe' } as Klaim
    expect(bagianBerubah(k, isianAwal(k))).toEqual([])
    expect(bagianBerubah(k, { ...isianAwal(k), isKpr: 'KPR' })).toEqual(['Register Claim'])
    const p1 = k.peserta[0]
    if (p1 === undefined) throw new Error('fixture')
    expect(bagianBerubah(k, { ...isianAwal(k), isKpr: 'KPR', tanggal: { [p1.id]: { ...tanggalAwal(p1), konfirmasi: '2025-07-01' } } }))
      .toEqual(['Register Claim', 'Participant Details'])
    expect(pesanSimpanDulu(['Register Claim'])).toBe('You have unsaved changes (Register Claim). Click Save first.')
    const html = en(<TombolBawahOutstanding klaim={k} sibukRNM={false} onSimpanRNM={() => {}} sesudahPindah={() => {}} adaPerubahan bagian={['Register Claim']} />)
    expect(html).toContain('<p role="note" class="alert alert--info os__bantuan">You have unsaved changes (Register Claim). Click Save first.</p>')
  })

  it('Register Claim: radio MENDATAR dan bulat; Catastrophe Note berlabel dengan placeholder, lebar penuh', () => {
    const k = klaimUji({ katastrofe: 'Catastrophe' })
    const html = en(<TabRegisterClaim klaim={k} isian={isianAwal(k)} onUbah={() => {}} boleh />)
    expect(html.match(/<dd class="os__radio-datar">/g)).toHaveLength(2)
    expect(html).not.toContain('os__radio-tegak')
    expect(html).toContain('placeholder="Describe the catastrophe…"')
    expect(html).toContain('<div class="os__baris-form os__catatan-katastrofe"><dt>Catastrophe Note</dt>')
    const css = readFileSync(join(__dirname, '..', 'claimlife.css'), 'utf8').replace(/\s+/g, ' ')
    expect(css).toContain('.claimlife .os__radio-datar { display: flex; flex-wrap: wrap; align-items: center; gap: 8px 24px;')
    expect(css).toContain('.claimlife .os__radio input[type="radio"] { appearance: auto;')
    expect(css).toContain('border-radius: 50%;')
  })
})

describe('tombol committee, konfirmasi Save to RNM (11-10-2026)', () => {
  it('"Send Claim to Committee" hanya di tahap Claim Analis', () => {
    const rincian = readFileSync(join(__dirname, '..', 'components', 'RincianPeserta.tsx'), 'utf8')
    // 11-10-2026 (terbaru): HANYA di tahap Claim Analis.
    expect(rincian).toContain('{klaim.tahap === TAHAP.claimAnalis && bolehSerahkanDiLayar(b) ? (')
    expect(rincian.split('TOMBOL_KOMITE.serahkan').length - 1).toBe(1)
  })

  it('modal: judul, pertanyaan, peringatan, ringkas peserta + PL, Cancel + "Yes, Save to RNM"', () => {
    const k = klaimUji()
    const el = KonfirmasiSimpanRNM({ klaim: k, sibuk: false, onLanjut: () => {}, onTutup: () => {} })
    expect(el.props.judul).toBe('Save to RNM')
    expect(el.props.labelBatal).toBe('Cancel')
    expect(el.props.tanpaTutup).toBe(false)
    const isi = renderToStaticMarkup(<>{el.props.children}</>)
    expect(isi).toContain('<p>Are you sure you want to Save to RNM?</p>')
    expect(isi).toContain('<p role="note" class="alert alert--warn">This action cannot be undone. A claim number will be generated, the insured will be set to Outstanding, and the claim data will be locked.</p>')
    expect(isi).toContain('<dt>Participant Details</dt><dd>2</dd>')
    expect(isi).toContain('<dt>PL Number</dt><dd class="os__nilai-kode">UJI-PL-1</dd>')
    expect(renderToStaticMarkup(el.props.aksi)).toBe('<button type="button" class="btn btn--primary">Yes, Save to RNM</button>')
  })

  it('Cancel tidak memanggil server; Yes memanggil sekali; saat proses nonaktif "Saving…" dan tak dapat ditutup', () => {
    let lanjut = 0
    let tutup = 0
    const el = KonfirmasiSimpanRNM({ klaim: klaimUji(), sibuk: false, onLanjut: () => lanjut++, onTutup: () => tutup++ })
    el.props.onTutup()
    expect([lanjut, tutup]).toEqual([0, 1])
    el.props.aksi.props.onClick()
    expect([lanjut, tutup]).toEqual([1, 1])
    const proses = KonfirmasiSimpanRNM({ klaim: klaimUji(), sibuk: true, onLanjut: () => lanjut++, onTutup: () => tutup++ })
    expect(proses.props.tanpaTutup).toBe(true)
    expect(renderToStaticMarkup(proses.props.aksi)).toBe('<button type="button" class="btn btn--primary" disabled="">Saving…</button>')
    // Halaman: klik Save to RNM -> pre-check -> modal; simpanKeRNM HANYA dari jalankanSimpanRNM (tombol modal).
    const rnm = SUMBER.slice(SUMBER.indexOf('async function simpanRNM()'), SUMBER.indexOf('async function jalankanSimpanRNM()'))
    expect(rnm.indexOf('kekuranganWajib(klaim, isian)')).toBeLessThan(rnm.indexOf('setKonfirmasiRNM(true)'))
    expect(rnm).not.toContain('simpanKeRNM(')
    expect(SUMBER.split('simpanKeRNM(klaimID)').length - 1).toBe(1)
    expect(SUMBER.split('void jalankanSimpanRNM()').length - 1).toBe(1)
    const jalan = SUMBER.slice(SUMBER.indexOf('async function jalankanSimpanRNM()'))
    expect(jalan).toContain('if (sibukRNM) return')
    expect(SUMBER).toContain('if (!sibukRNM) setKonfirmasiRNM(false)')
  })
})

describe('MAX CLAIM RECEIVED tanpa teks galat teknis (11-10-2026; 12-10-2026 disembunyikan bila kosong)', () => {
  it('alasan tidak dirender; tidak dihitung = baris tidak ada (tanpa ⓘ); penanda terisi tampil apa adanya', () => {
    const p = { ...pesertaUji('1'), penandaTerimaKlaim: '', penandaTerimaKlaimAlasan: 'models: ambang hari produk belum diisi di Master Product Name Life' } as unknown as Peserta
    const html = en(<RincianPeserta klaim={klaimUji()} p={p} onMuatUlang={async () => {}} />)
    expect(html).not.toContain('not computed')
    expect(html).not.toContain('models:')
    expect(html).not.toContain('Max Claim Received')
    expect(html).not.toContain('ⓘ')
    const ada = en(<RincianPeserta klaim={klaimUji()} p={{ ...p, penandaTerimaKlaim: '05/07/2025', penandaTerimaKlaimAlasan: '' } as Peserta} onMuatUlang={async () => {}} />)
    expect(ada).toContain('<dt>Max Claim Received</dt><dd>05/07/2025</dd>')
  })
})

describe('Policy Details PL Number, Reject Outstanding, dokumen (11-10-2026)', () => {
  it('PL Number: label rapi, nilai kode satu baris, kolom cukup lebar (tanpa gulir)', () => {
    const html = en(<PanelPolisOutstanding polis={null} plNumber="RNML-QRL2.10.26.00043" />)
    expect(html).toContain('<section class="os__form os__form--tiga os__form--polis" aria-label="Policy Details">')
    expect(html).toContain('<dt>PL Number</dt><dd class="os__nilai-pl">RNML-QRL2.10.26.00043</dd>')
    const css = readFileSync(join(__dirname, '..', 'claimlife.css'), 'utf8').replace(/\s+/g, ' ')
    expect(css).toContain('.claimlife .os__form--polis { grid-template-columns: repeat(auto-fit, minmax(min(100%, 380px), 1fr)); }')
    expect(css).toContain('.claimlife .os__nilai-pl { white-space: nowrap; }')
    expect(css).not.toMatch(/os__form--polis[^{]*\{[^}]*overflow/)
  })

  it('Reject Outstanding: modal konfirmasi + peringatan, Remarks tetap, "Yes, Reject" nonaktif saat proses', () => {
    const isi = en(<DialogTolakOutstanding sibuk={false} galat={null} onKirim={() => {}} onBatal={() => {}} subjek="UJI NAMA 1" />)
    expect(isi).toContain('Reject Outstanding')
    expect(isi).toContain('<p>Are you sure you want to reject this outstanding claim for UJI NAMA 1?</p>')
    expect(isi).toContain('<p role="note" class="alert alert--warn">This action cannot be undone. The participant&#x27;s claim will be marked as Rejected.</p>')
    expect(isi).toContain('Remarks<textarea')
    // Remarks kosong: "Yes, Reject" nonaktif (Remarks tetap wajib seperti Pega).
    expect(isi).toContain('<button type="submit" class="btn btn--danger" disabled="">Yes, Reject</button>')
    const proses = en(<DialogTolakOutstanding sibuk galat={null} onKirim={() => {}} onBatal={() => {}} />)
    expect(proses).toContain('Rejecting…')
    expect(proses).toContain('for this insured?')
    const dialog = readFileSync(join(__dirname, '..', 'components', 'DialogTolakOutstanding.tsx'), 'utf8')
    expect(dialog).toContain('tanpaTutup={sibuk}')
    expect(dialog).toContain('if (sah && !sibuk) onKirim(remarks.trim())')
    const rincian = readFileSync(join(__dirname, '..', 'components', 'RincianPeserta.tsx'), 'utf8')
    expect(rincian).toContain('subjek={p.namaTertanggung || p.nomorSertifikat}')
    expect(rincian).toContain('if (menolak !== null) return')
  })

  it('Dokumen: label Inggris, Category dropdown dari server, tombol "Add attachment" sesudah memilih', () => {
    const html = en(<PanelDokumenPeserta klaimID="K-1" pesertaID="P-1" dokumen={[]} onBerubah={() => {}} />)
    expect(html).toContain('<h4 class="dokumen__judul">Claim Documents</h4>')
    expect(html).toContain('<p class="dokumen__kosong">No claim documents for this insured yet.</p>')
    expect(html).toMatch(/Category <select class="field__input"[^>]*disabled=""[^>]*><option value="" selected="">Select category…<\/option><\/select>/)
    expect(html).toContain('<button type="button" class="btn btn--primary btn--sm" disabled="">Add attachment</button>')
    expect(html).not.toContain('type="text"')
    const sumber = readFileSync(join(__dirname, '..', 'components', 'PanelDokumenPeserta.tsx'), 'utf8')
    expect(sumber).toContain('ambilKategoriDokumen(klaimID)')
    expect(sumber).toContain("disabled={sibuk || kategori === '' || berkas === null}")
  })
})

describe('hapus dan lihat dokumen (11-10-2026)', () => {
  const dok = (id: string, storage: string) =>
    ({ id, namaFile: 'UJI-' + id + '.pdf', kategori1: 'DL-1', kategori2: 'Document Claim', tanggal: '2026-10-11T16:02:41+07:00', tStorageId: storage }) as unknown as Dokumen
  it('setiap baris (termasuk yang masih diproses) punya View + Delete; tanggal dd/mm/yyyy HH:mm', () => {
    const html = en(<PanelDokumenPeserta klaimID="K-1" pesertaID="P-1" dokumen={[dok('1', ''), dok('2', 'IMG')]} onBerubah={() => {}} />)
    expect(html.match(/>View<\/button>/g)).toHaveLength(2)
    expect(html.match(/>Delete<\/button>/g)).toHaveLength(2)
    expect(html.match(/Processing — waiting for storage<\/small>/g)).toHaveLength(1)
    expect(html).not.toContain('2026-10-11T16:02:41')
    expect(tanggalDokumen('')).toBe('—')
    expect(tanggalDokumen('bukan tanggal')).toBe('bukan tanggal')
    expect(tanggalDokumen('2026-10-11T16:02:41+07:00')).toMatch(/^\d{2}\/\d{2}\/2026 \d{2}:\d{2}$/)
  })
  it('Delete lewat modal konfirmasi; View berbatas klaim', () => {
    const sumber = readFileSync(join(__dirname, '..', 'components', 'PanelDokumenPeserta.tsx'), 'utf8')
    expect(KONFIRMASI_HAPUS_DOKUMEN).toMatchObject({ judul: 'Delete this document?', peringatan: 'This cannot be undone.', ya: 'Yes, Delete' })
    expect(sumber).toContain('setHapus(b)')
    expect(sumber).toContain('await hapusDokumen(klaimID, target.id)')
    expect(sumber).toContain('tanpaTutup={sibuk}')
    expect(sumber).toContain('if (!sibuk) setHapus(null)')
    expect(sumber).toContain('ambilIsiDokumenKlaim(klaimID, b.id)')
    expect(sumber.split('hapusDokumen(').length - 1).toBe(1) // satu-satunya pemanggil: tombol modal
  })
})

describe('kategori dokumen dari produk (11-10-2026)', () => {
  it('pesan produk tanpa Document Claim; nama > 64 byte dinyatakan, bukan dipotong; rute berbatas klaim', () => {
    const sumber = readFileSync(join(__dirname, '..', 'components', 'PanelDokumenPeserta.tsx'), 'utf8')
    expect(sumber).toContain("kategoriKosong: 'No claim document types are set for this product in Master Product Name Life.'")
    expect(kategoriTerlaluPanjang('Sertifikat peserta (Participant certificate)')).toBe(false)
    expect(kategoriTerlaluPanjang('Copy identitas diri KTP/SIM/Paspor (Identity card copy KTP/SIM/Passport)')).toBe(true)
    expect(API).toContain('/api/klaim-life/${encodeURIComponent(klaimID)}/dokumen/kategori')
  })
})

describe('proteksi Send to Medical Check (11-10-2026)', () => {
  const dok = { id: '1', namaFile: 'UJI.pdf', kategori1: '', kategori2: 'Document Claim', tanggal: '', tStorageId: '' }
  const k = (peserta: Peserta[]) => ({ ...klaimUji({ nomorKlaim: 'UJI-CLM-1', tersimpan: true }), peserta }) as Klaim

  it('pra-periksa: peserta Accept tanpa dokumen didaftar SEMUA; berdokumen / tanpa Accept tidak', () => {
    const p1 = { ...pesertaUji('1', 'true'), dokumen: [] } as unknown as Peserta
    const p2 = { ...pesertaUji('2', 'true'), dokumen: [dok] } as unknown as Peserta
    const p3 = { ...pesertaUji('3', 'false'), dokumen: [] } as unknown as Peserta
    const p4 = { ...pesertaUji('4', 'true'), dokumen: [] } as unknown as Peserta
    const kurang = kekuranganDokumenMedis(k([p1, p2, p3, p4]))
    expect(kurang.map((x) => [x.nomorPeserta, x.pesertaId])).toEqual([[1, '1'], [4, '4']])
    expect(kurang[0]?.pesan).toBe('Insured number 1 (UJI NAMA 1): at least one claim document is required because the insured is accepted.')
    expect(kekuranganDokumenMedis(k([p2, p3]))).toEqual([])
  })

  it('kotak kekurangan berjudul Send to Medical; peserta pertama dibuka di sub-tab Document', () => {
    const html = en(<PitaGalat judul={JUDUL_KEKURANGAN_MEDIS} kekurangan={[{ nomorPeserta: 1, pesertaId: '1', sertifikat: '001', nama: 'UJI A', kolom: 'Document', pesan: '' }]} />)
    expect(html).toContain('<strong>Please upload at least one claim document for every accepted insured before Send to Medical Check:</strong>')
    expect(html).toContain('<li>No 1 · UJI A: Document</li>')
    const p = { ...pesertaUji('1', 'true'), dokumen: [] } as unknown as Peserta
    const rincian = en(<RincianPeserta klaim={k([p])} p={p} onMuatUlang={async () => {}} kolomKurang={new Set(['Document'])} />)
    expect(rincian).toMatch(/role="tab" aria-selected="true"[^>]*>Document<\/button>/)
  })

  it('perubahan belum disimpan: Send to Medical Check nonaktif + petunjuk, juga sesudah Save to RNM', () => {
    const kk = k([pesertaUji('1', 'true')])
    const html = en(<TombolBawahOutstanding klaim={kk} sibukRNM={false} onSimpanRNM={() => {}} sesudahPindah={() => {}} adaPerubahan />)
    expect(nonaktif(html, 'Send to Medical Check')).toBe(true)
    expect(html).toContain('<p role="note" class="alert alert--info os__bantuan">Save your changes first.</p>')
    const bersih = en(<TombolBawahOutstanding klaim={kk} sibukRNM={false} onSimpanRNM={() => {}} sesudahPindah={() => {}} adaPerubahan={false} />)
    expect(ada(bersih, 'Send to Medical Check') && !nonaktif(bersih, 'Send to Medical Check')).toBe(true)
  })

  it('halaman: pra-periksa sebelum pindah, galat server berdaftar ke tempat galat yang sama', () => {
    expect(SUMBER).toContain('const kurang = kekuranganDokumenMedis(klaim)')
    expect(SUMBER).toContain('tunjukkanKekurangan(kurang, JUDUL_KEKURANGAN_MEDIS)')
    expect(SUMBER).toContain('if (e instanceof GagalSimpanRNM) {\n                        tunjukkanKekurangan(e.kekurangan, JUDUL_KEKURANGAN_MEDIS)')
    const panel = readFileSync(join(__dirname, '..', 'components', 'PanelPindahTahap.tsx'), 'utf8')
    expect(panel).toContain('if (sebelumPindah !== undefined && !sebelumPindah(t)) return')
    expect(panel).toContain('if (tanganiGalat?.(e) !== true)')
  })
})

describe('konfirmasi Send to Medical Check (11-10-2026)', () => {
  it('modal: judul, pertanyaan, peringatan, tombol - TANPA ringkasan Claim No / peserta (12-10-2026)', () => {
    const el = KonfirmasiKirimMedis({ sibuk: false, onLanjut: () => {}, onTutup: () => {} })
    expect(el.props.judul).toBe('Send to Medical Check')
    expect(el.props.labelBatal).toBe('Cancel')
    const isi = renderToStaticMarkup(<>{el.props.children}</>)
    expect(isi).toContain('<p>Are you sure you want to send this claim to Medical Advisor?</p>')
    expect(isi).toContain('<p role="note" class="alert alert--warn">After sending, the claim moves to the Medical Advisor queue and can no longer be edited from the Claim Admin screen.</p>')
    expect(isi).not.toContain('Claim No')
    expect(isi).not.toContain('<dl')
    expect(renderToStaticMarkup(el.props.aksi)).toBe('<button type="button" class="btn btn--primary">Yes, Send to Medical Check</button>')
  })

  it('Cancel tidak memanggil; Yes sekali; saat proses nonaktif "Sending…" dan tak dapat ditutup', () => {
    let lanjut = 0
    let tutup = 0
    const el = KonfirmasiKirimMedis({ sibuk: false, onLanjut: () => lanjut++, onTutup: () => tutup++ })
    el.props.onTutup()
    expect([lanjut, tutup]).toEqual([0, 1])
    el.props.aksi.props.onClick()
    expect([lanjut, tutup]).toEqual([1, 1])
    const proses = KonfirmasiKirimMedis({ sibuk: true, onLanjut: () => {}, onTutup: () => {} })
    expect(proses.props.tanpaTutup).toBe(true)
    expect(renderToStaticMarkup(proses.props.aksi)).toBe('<button type="button" class="btn btn--primary" disabled="">Sending…</button>')
  })

  it('pra-periksa lolos membuka modal; satu pemanggil pindahTahap; sukses ke Inbox dengan pita', () => {
    expect(SUMBER).toContain('setKonfirmasiMedis(true)\n                      return false')
    // Satu pemanggil per alur (12-10-2026): halaman Outstanding (Send to Medical Check) dan TombolBawahTahap (Medical Check /
    // Claim Analis, sesudah Modal konfirmasi) - masing-masing tepat satu.
    const halaman = SUMBER.slice(SUMBER.indexOf('export default function OutstandingClaimLife'))
    const tahap = SUMBER.slice(SUMBER.indexOf('export function TombolBawahTahap'), SUMBER.indexOf('export function KonfirmasiKirimMedis'))
    expect(halaman.split('pindahTahap(').length - 1).toBe(1)
    expect(tahap.split('pindahTahap(').length - 1).toBe(1)
    expect(SUMBER.split('pindahTahap(').length - 1).toBe(2)
    expect(SUMBER.split('void jalankanKirimMedis()').length - 1).toBe(1)
    const jalan = SUMBER.slice(SUMBER.indexOf('async function jalankanKirimMedis()'))
    expect(jalan).toContain('if (sibukMedis) return')
    expect(jalan).toContain('onPindah(kabarTerkirimMedis(klaimID))')
    expect(SUMBER).toContain('if (!sibukMedis) setKonfirmasiMedis(false)')
    expect(kabarTerkirimMedis('CLMLF-000012')).toBe('Claim CLMLF-000012 sent to Medical Advisor.')
    const rute = readFileSync(join(__dirname, '..', 'rute.tsx'), 'utf8')
    expect(rute).toContain("setKabar(k ?? '')")
    expect(rute).toContain('kabar={kabar}')
    const inbox = readFileSync(join(__dirname, 'InboxClaimLife.tsx'), 'utf8')
    expect(inbox).toContain('<PemberitahuanKabar kabar={kabar} onTutup={onTutupKabar} />')
  })
})
