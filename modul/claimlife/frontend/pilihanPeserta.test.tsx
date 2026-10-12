import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { REGISTER_BANYAK } from './labels'
import {
  alihPilihan,
  kalimatJumlahCase,
  kelompokkanPilihan,
  kelompokPermintaan,
  sudahDipilih,
  type PesertaTerpilih,
  bacaPolisPeserta,
} from './pilihanPeserta'
import { DaftarTerpilih, idTerdaftar } from './pages/RegisterKlaim'
import { PemberitahuanTerdaftar } from './pages/InboxClaimLife'
import { PanelDataPolis } from './components/PanelDataPolis'
import KonfirmasiDaftar from './components/KonfirmasiDaftar'

// Register banyak PL (keputusan work owner 10-10-2026): pilihan bertahan, kelompok PL + mata uang, banyak nomor.

const SUMBER = readFileSync(join(__dirname, 'pages', 'RegisterKlaim.tsx'), 'utf8').replace(/\r\n/g, '\n')

const p = (pl: string, sertifikat: string, uang = 'IDR', polis = `POL-${pl}`): PesertaTerpilih => ({
  nomorPremiList: pl,
  nomorPolis: polis,
  nomorSertifikat: sertifikat,
  namaTertanggung: `UJI-NAMA-${sertifikat}`,
  mataUang: uang,
})

describe('pilihan peserta bertahan antar pencarian', () => {
  it('pilihan PL A tetap ada setelah memilih dari PL B; hasil pencarian menandai yang sudah terpilih', () => {
    let d: PesertaTerpilih[] = []
    d = alihPilihan(d, p('UJI-PL-A', '001'))
    d = alihPilihan(d, p('UJI-PL-A', '002'))
    // "mencari" PL B: daftar pilihan tidak disentuh, lalu memilih 4 peserta
    for (const s of ['101', '102', '103', '104']) d = alihPilihan(d, p('UJI-PL-B', s))
    expect(d).toHaveLength(6)
    expect(sudahDipilih(d, p('UJI-PL-A', '001'))).toBe(true)
    expect(sudahDipilih(d, p('UJI-PL-A', '003'))).toBe(false)
    // centang kedua = lepas
    expect(alihPilihan(d, p('UJI-PL-A', '001'))).toHaveLength(5)
    // sertifikat sama, polis berbeda = peserta berbeda
    expect(sudahDipilih(d, p('UJI-PL-A', '001', 'IDR', 'POL-LAIN'))).toBe(false)
  })

  it('cari() tidak lagi mengosongkan pilihan; validasi "berbeda polis atau mata uang" dibuang', () => {
    const cari = SUMBER.slice(SUMBER.indexOf('async function cari()'), SUMBER.indexOf('function pilih('))
    expect(cari).not.toContain('setDipilih([])')
    expect(SUMBER).not.toContain('berbeda polis atau mata uang')
    // sesudah sukses pilihan dikosongkan
    const kirim = SUMBER.slice(SUMBER.indexOf('async function kirim()'), SUMBER.indexOf('return (\n    <section>'))
    expect(kirim).toContain('setDipilih([])\n      onTerdaftar?.(idTerdaftar(h))')
  })
})

describe('pengelompokan dan jumlah case', () => {
  it('2 dari PL A + 4 dari PL B = 2 kelompok; satu PL dua mata uang = 2 kelompok', () => {
    const d = [p('UJI-PL-A', '001'), p('UJI-PL-B', '101'), p('UJI-PL-A', '002'), p('UJI-PL-B', '102'), p('UJI-PL-B', '103'), p('UJI-PL-B', '104')]
    const k = kelompokkanPilihan(d)
    expect(k.map((x) => [x.nomorPremiList, x.mataUang, x.peserta.length])).toEqual([
      ['UJI-PL-A', 'IDR', 2],
      ['UJI-PL-B', 'IDR', 4],
    ])
    expect(kelompokkanPilihan([p('UJI-PL-C', '201', 'IDR'), p('UJI-PL-C', '202', 'USD')])).toHaveLength(2)
    expect(kalimatJumlahCase(2, REGISTER_BANYAK.jumlahCase)).toBe('2 claim case(s) will be created')
  })

  it('badan permintaan: satu kelompok per PL + polis, sertifikat teks apa adanya', () => {
    expect(kelompokPermintaan([p('UJI-PL-A', '001'), p('UJI-PL-B', '101'), p('UJI-PL-A', '006')])).toEqual([
      { nomorPremiList: 'UJI-PL-A', nomorPolis: 'POL-UJI-PL-A', mataUang: 'IDR', sertifikat: ['001', '006'] },
      { nomorPremiList: 'UJI-PL-B', nomorPolis: 'POL-UJI-PL-B', mataUang: 'IDR', sertifikat: ['101'] },
    ])
  })

  it('Selected insured (tampilan 10-10-2026): TABEL per kelompok dengan kolom No/Certificate/Name/Policy, Remove per baris, ringkasan + Clear + Register Claim', () => {
    const html = renderToStaticMarkup(
      <DaftarTerpilih kelompok={kelompokkanPilihan([p('UJI-PL-A', '001'), { ...p('UJI-PL-B', '101'), tanggalLahir: '1990-11-22' }, p('UJI-PL-B', '102')])} onHapus={() => {}} onKosongkan={() => {}} onDaftar={() => {}} />,
    ).replace(/<!-- -->/g, '')
    expect(html).toContain('<h3>Selected insured</h3>')
    expect(html).toContain('<h4>PL UJI-PL-A<span class="register__chip">IDR</span><span class="register__jumlah">1 insured</span></h4>')
    expect(html).toContain('<span class="register__jumlah">2 insured</span>')
    expect(html.split('<thead><tr><th>No</th><th>Certificate No</th><th>Name of Insured</th><th>Date of Birth</th><th aria-label="Remove"></th></tr></thead>')).toHaveLength(3)
    expect(html).toContain('<td>1</td><td>101</td><td>UJI-NAMA-101</td><td>22/11/1990</td>')
    expect(html).toContain('<td>2</td><td>102</td><td>UJI-NAMA-102</td><td>—</td>')
    expect(html).not.toContain('Policy No')
    expect(html).not.toContain('POL-UJI-PL-B')
    expect(html).toContain('<td>2</td><td>102</td>')
    expect(html.split('class="btn btn--ghost btn--sm">Remove</button>')).toHaveLength(4)
    expect(html).not.toContain('<ul>')
    expect(html).toContain('<p class="register__ringkas">2 claim case(s) will be created</p>')
    expect(html).toContain('class="btn btn--ghost">Clear selection</button>')
    expect(html).toContain('class="btn btn--primary">Register Claim</button>')
    const kosong = renderToStaticMarkup(<DaftarTerpilih kelompok={[]} onHapus={() => {}} onKosongkan={() => {}} />)
    expect(kosong).toContain('No insured selected yet.')
    expect(kosong).not.toContain('<table')
  })

  it('halaman Register (tampilan 10-10-2026): judul, label pencarian Inggris, tabel hasil berbaris-klik, tanpa Daftarkan klaim', () => {
    expect(SUMBER).toContain('<h2>{REGISTER_BANYAK.judulHalaman}</h2>')
    expect(REGISTER_BANYAK.judulHalaman).toBe('Register Claim Life')
    expect(REGISTER_BANYAK.labelPL).toBe('Premium List No')
    expect(SUMBER).toContain(`className="btn btn--primary" disabled={sibuk || pl.trim() === ''}`)
    expect(SUMBER).toContain("className={aktif ? 'register__baris register__baris--pilih' : 'register__baris'}")
    expect(SUMBER).toContain('onClick={() => pilih(p)}')
    expect(SUMBER).not.toContain('Daftarkan klaim')
    expect(SUMBER).not.toContain('Nomor premium list')
    expect(SUMBER).not.toContain('<th>Sertifikat</th>')
  })
})

describe('sesudah Register langsung ke Inbox (keputusan work owner 10-10-2026)', () => {
  const REG = readFileSync(join(__dirname, 'pages', 'RegisterKlaim.tsx'), 'utf8').replace(/\r\n/g, '\n')
  const RUTE = readFileSync(join(__dirname, 'rute.tsx'), 'utf8').replace(/\r\n/g, '\n')
  const kirim = REG.slice(REG.indexOf('async function kirim()'), REG.indexOf('  return (', REG.indexOf('async function kirim()')))
  const coba = kirim.slice(kirim.indexOf('try {'), kirim.indexOf('} catch'))
  const tangkap = kirim.slice(kirim.indexOf('} catch'), kirim.indexOf('} finally'))

  it('idTerdaftar: Case ID urut jawaban server', () => {
    const h = {
      kasus: [
        { id: 'UJI-W1', nomorKlaim: '', nomorPremiList: 'UJI-PL-A', mataUang: 'IDR', jumlahPeserta: 2 },
        { id: 'UJI-W2', nomorKlaim: '', nomorPremiList: 'UJI-PL-B', mataUang: 'IDR', jumlahPeserta: 4 },
      ],
    }
    expect(idTerdaftar(h)).toEqual(['UJI-W1', 'UJI-W2'])
  })

  it('berhasil: onTerdaftar dipanggil dengan Case ID sesudah daftarKlaimLife; gagal: tidak pindah, pilihan utuh', () => {
    expect(coba).toContain('onTerdaftar?.(idTerdaftar(h))')
    expect(coba.indexOf('daftarKlaimLife(')).toBeLessThan(coba.indexOf('onTerdaftar?.('))
    expect(tangkap).toContain('setGalat(kalimatGalat(e))')
    expect(tangkap).not.toContain('onTerdaftar')
    expect(tangkap).not.toContain('setDipilih')
    // Modal tertutup pada kedua jalur.
    expect(kirim.slice(kirim.indexOf('} finally'))).toContain('setKonfirmasi(false)')
  })

  it('tabel hasil dihapus dari layar Register', () => {
    expect(REG).not.toContain('HasilKasus')
    expect(REG).not.toContain('setHasil')
    expect(REG).not.toContain('REGISTER_BANYAK.hasil')
  })

  it('rute: Case ID lewat state modul (bukan localStorage), pindah ke inbox, notifikasi bisa ditutup', () => {
    expect(RUTE).toContain('const [terdaftar, setTerdaftar] = useState<string[]>([])')
    expect(RUTE).toContain('setTerdaftar(ids)')
    expect(RUTE.indexOf('setTerdaftar(ids)')).toBeLessThan(RUTE.indexOf("onPindah('inbox')", RUTE.indexOf('setTerdaftar(ids)')))
    expect(RUTE).toContain('terdaftar={terdaftar}')
    expect(RUTE).toContain('onTutupTerdaftar={() => setTerdaftar([])}')
    expect(RUTE).not.toContain('localStorage')
    expect(RUTE).not.toContain('sessionStorage')
  })

  it('inbox menampilkan pemberitahuan sukses berisi Case ID; kosong = tidak tampil', () => {
    const html = renderToStaticMarkup(<PemberitahuanTerdaftar ids={['UJI-W1', 'UJI-W2']} onTutup={() => {}} />)
    expect(html).toContain('alert alert--ok')
    expect(html).toContain('2 claim case(s) registered: UJI-W1, UJI-W2')
    expect(html).toContain('>Close</button>')
    expect(renderToStaticMarkup(<PemberitahuanTerdaftar ids={[]} />)).toBe('')
  })
})

describe('data polis otomatis dari IDPEGA peserta (keputusan work owner 10-10-2026)', () => {
  it('bacaPolisPeserta: hanya IDPEGA, tanpa fallback; kosong = tidak ditemukan; tanpa peserta = tidak membaca', async () => {
    const dipanggil: string[] = []
    const ambil = (ada: string[]) => async (k: string) => {
      dipanggil.push(k)
      if (!ada.includes(k)) throw new Error('404')
      return { nomorPolis: k }
    }
    expect(await bacaPolisPeserta(ambil(['NBLF-UJI-1']), 'NBLF-UJI-1')).toEqual({ polis: { nomorPolis: 'NBLF-UJI-1' }, tidakDitemukan: false, idPega: 'NBLF-UJI-1' })
    expect(await bacaPolisPeserta(ambil([]), 'EDML-UJI-2')).toEqual({ polis: null, tidakDitemukan: true, idPega: 'EDML-UJI-2' })
    expect(dipanggil).toEqual(['NBLF-UJI-1', 'EDML-UJI-2'])
    expect(await bacaPolisPeserta(ambil([]), ' ')).toEqual({ polis: null, tidakDitemukan: true, idPega: '' })
    expect(await bacaPolisPeserta(ambil([]), undefined)).toEqual({ polis: null, tidakDitemukan: false, idPega: '' })
    expect(dipanggil).toHaveLength(2)
  })

  it('panel: tanpa teks teknis .PolicyDataLife; header tak ditemukan = satu kalimat menyebut IDPEGA; tanpa centang = tanpa catatan', () => {
    const kosong = renderToStaticMarkup(<PanelDataPolis polis={null} />)
    expect(kosong).not.toContain('PolicyDataLife')
    expect(kosong).not.toContain('polis__catatan')
    const gagal = renderToStaticMarkup(<PanelDataPolis polis={null} polisTidakDitemukan={{ idPega: 'NBLF-UJI-1' }} />).replace(/<!-- -->/g, '')
    expect(gagal).toContain('Policy data not found for this insured (NBLF-UJI-1).')
    expect(gagal).not.toContain('PolicyDataLife')
  })

  it('Register dan Outstanding memanggil rute IDPEGA; Type/Kode bisnis hilang; blok Selected insured dirender', () => {
    expect(SUMBER).not.toContain('setType')
    expect(SUMBER).not.toContain('Kode bisnis')
    expect(SUMBER).not.toContain('kodeBisnis')
    expect(SUMBER).toContain('dipilih[dipilih.length - 1]')
    expect(SUMBER).toContain("const idPegaTerpilih = pertamaDipilih?.idPega ?? ''")
    expect(SUMBER).toContain('bacaPolisPeserta(ambilDataPolisIdPega, idPegaTerpilih)')
    expect(SUMBER).not.toContain('ambilDataPolis(')
    expect(SUMBER).toContain('<DaftarTerpilih')
    expect(SUMBER).toContain('daftarKlaimLife({ kelompok: kelompokPermintaan(dipilih) })')
    const os = readFileSync(join(__dirname, 'pages', 'OutstandingClaimLife.tsx'), 'utf8')
    expect(os).toContain('bacaPolisPeserta(ambilDataPolisIdPega, idPega)')
    expect(os).not.toContain('ambilDataPolis(')
  })

  it('pilihan membawa IDPEGA peserta', () => {
    const d = alihPilihan([], { ...p('UJI-PL-A', '001'), idPega: 'NBLF-UJI-1' })
    expect(d[0]?.idPega).toBe('NBLF-UJI-1')
  })
})

describe('polis wajib Resolved-Completed (keputusan work owner 10-10-2026)', () => {
  it('panel memperingatkan polis yang belum Resolved-Completed; yang sudah selesai tanpa peringatan', () => {
    const dasar = { nomorPolis: '', type: 'QR', marketingName: '', cedingCoName: '', policyHolderName: '', businessName: '', dateReceived: '', status: 'Pending', statusUpdate: '', productNameId: '', productName: '', prodKe: 1, medanTanpaSumber: [], typeCeding: '', typeCedingName: '', proRateType: '', wpc: '', retroName: '', securityReinsurer: '', sobName: '' }
    const belum = renderToStaticMarkup(<PanelDataPolis polis={dasar} />).replace(/<!-- -->/g, '')
    expect(belum).toContain('Policy not Resolved-Completed (status: Pending). A claim cannot be registered yet.')
    const selesai = renderToStaticMarkup(<PanelDataPolis polis={{ ...dasar, status: 'Resolved-Completed' }} />)
    expect(selesai).not.toContain('not Resolved-Completed')
  })
})

describe('panel Data Polis rapi (tampilan 10-10-2026)', () => {
  it('tanpa System Reinsurance / Premium Method, tanpa baris teknis versi; chip versi di judul; tanggal dd/mm/yyyy', () => {
    const polis = { nomorPolis: 'UJI-PL-A', type: 'QR', marketingName: '', cedingCoName: '', policyHolderName: '', businessName: '', dateReceived: '2026-01-15', status: 'Resolved-Completed', statusUpdate: '', productNameId: '', productName: 'UJI PRODUK', prodKe: 0, medanTanpaSumber: [], typeCeding: '1', typeCedingName: 'UJI QS', proRateType: 'UJI-PRORATA', wpc: '2026-02-01', retroName: '', securityReinsurer: '', sobName: '' }
    const html = renderToStaticMarkup(<PanelDataPolis polis={polis} />).replace(/<!-- -->/g, '')
    expect(html).not.toContain('System Reinsurance')
    expect(html).not.toContain('Premium Method')
    expect(html).not.toContain('versi')
    expect(html).toContain('<span class="polis__chip">New business</span>')
    expect(html).toContain('<dd>01/02/2026</dd>')
    expect(html).toContain('<dd>15/01/2026</dd>')
    expect(renderToStaticMarkup(<PanelDataPolis polis={{ ...polis, prodKe: 3 }} />)).toContain('Endorsement /03')
    expect(renderToStaticMarkup(<PanelDataPolis polis={{ ...polis, status: 'Pending' }} />)).toContain('class="alert alert--warn"')
  })
})

describe('DOB dan konfirmasi Register Claim (keputusan work owner 10-10-2026)', () => {
  it('tabel hasil: kolom Date of Birth dd/mm/yyyy menggantikan Policy No', () => {
    expect(SUMBER).toContain('<th>{REGISTER_BANYAK.kolomTanggalLahir}</th>')
    expect(SUMBER).toContain('<td>{tanggalTampil(p.tanggalLahir)}</td>')
    expect(SUMBER).not.toContain('kolomPolisCeding')
    expect(SUMBER).not.toContain('<td>{p.nomorPolis}</td>')
    expect(REGISTER_BANYAK.kolomTanggalLahir).toBe('Date of Birth')
    // nomor polis ceding tetap di data pilihan (tidak dihapus), hanya tidak ditampilkan
    expect(alihPilihan([], { ...p('UJI-PL-A', '001'), tanggalLahir: '1990-11-22' })[0]).toMatchObject({ nomorPolis: 'POL-UJI-PL-A', tanggalLahir: '1990-11-22' })
  })

  it('modal: judul, kalimat jumlah case, tabel per case, Continue?, Cancel + Register Claim', () => {
    const kel = kelompokkanPilihan([p('UJI-PL-A', '001'), p('UJI-PL-B', '101'), p('UJI-PL-B', '102')])
    const el = KonfirmasiDaftar({ kelompok: kel, sibuk: false, onLanjut: () => {}, onTutup: () => {} })
    expect(el.props.judul).toBe('Register Claim')
    expect(el.props.labelBatal).toBe('Cancel')
    const isi = renderToStaticMarkup(<>{el.props.children}</>)
    expect(isi).toContain('<p>2 claim case(s) will be created:</p>')
    expect(isi).toContain('<td>UJI-PL-A</td><td>IDR</td><td>1</td>')
    expect(isi).toContain('<td>UJI-PL-B</td><td>IDR</td><td>2</td>')
    expect(isi).toContain('<p>Continue?</p>')
    expect(renderToStaticMarkup(el.props.aksi)).toBe('<button type="button" class="btn btn--primary">Register Claim</button>')
  })

  it('Cancel tidak mendaftar; Register Claim memanggil sekali; saat proses tombol nonaktif "Registering…" dan modal tak dapat ditutup', () => {
    let lanjut = 0
    let tutup = 0
    const kel = kelompokkanPilihan([p('UJI-PL-A', '001')])
    const el = KonfirmasiDaftar({ kelompok: kel, sibuk: false, onLanjut: () => lanjut++, onTutup: () => tutup++ })
    el.props.onTutup()
    expect([lanjut, tutup]).toEqual([0, 1])
    el.props.aksi.props.onClick()
    expect([lanjut, tutup]).toEqual([1, 1])
    const proses = KonfirmasiDaftar({ kelompok: kel, sibuk: true, onLanjut: () => lanjut++, onTutup: () => tutup++ })
    expect(proses.props.tanpaTutup).toBe(true)
    expect(renderToStaticMarkup(proses.props.aksi)).toBe('<button type="button" class="btn btn--primary" disabled="">Registering…</button>')
    // Halaman: Register Claim di layar hanya MEMBUKA modal; daftarKlaimLife hanya di kirim(), dan kirim() hanya dari modal; penjaga klik ganda.
    expect(SUMBER).toContain('onDaftar={() => setKonfirmasi(true)}')
    expect(SUMBER).toContain('onLanjut={() => void kirim()}')
    expect(SUMBER.split('void kirim()').length - 1).toBe(1) // satu-satunya pemanggil: tombol modal
    expect(SUMBER.split('daftarKlaimLife(').length - 1).toBe(1)
    expect(SUMBER).toContain('if (sedangKirim.current) return')
    expect(SUMBER).toContain('if (!sibuk) setKonfirmasi(false)')
  })
})
