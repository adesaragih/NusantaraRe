// Layar tabel acuan Treaty In — tiket 15.
//
// ⛔ INI SELURUH LAYAR MODUL INI HARI INI, dan itu disengaja. Papan tiket
// menyatakan `L-4` ("tidak ada spesifikasi layar di mana pun") dan melarang
// mengarang layar untuk memenuhi bentuk irisan tegak: layar yang dikarang
// dipakai sebagai kriteria selesai oleh orang yang tidak tahu ia karangan.
// Halaman ini TIDAK mengarang alur kontrak - ia memperlihatkan isi kelima tabel
// acuan yang tiket 15 buat, tidak lebih.
//
// Rupa (02-10-2026): isinya duduk di dalam `Panel` - kartu bertepi yang sudah
// dipakai claimlife, mastercontractretrolife, dan treatycontractout - bukan
// menempel langsung di `--bg`. Keadaan kosong memakai `Kosong`, BUKAN paragraf
// telanjang; ketiga keadaan (gagal · memuat · kosong) tetap tiga cabang yang
// terpisah, lihat `dasar.tsx` baris 344.

import { useEffect, useState } from 'react'

import { Gagal, Kosong, Memuat, Panel, StripTab } from '../../../../inti/frontend/components/ui/dasar'
import { ambilAcuan, HIMPUNAN_ACUAN, type Acuan, type HimpunanAcuan } from '../api'
import { ACUAN_TREATYIN, LABEL_HIMPUNAN, MENU_TREATYIN } from '../labels'

/**
 * Penanda aktif sebagai lencana.
 *
 * ⛔ Teksnya nilai APA ADANYA. `AKTIF` adalah `VARCHAR2(40 CHAR)` tanpa CHECK
 * (`400_tabel_acuan.sql`) dan `KAMUS-KOLOM.md` §10.22 tidak menyebut domainnya,
 * jadi memetakan nilai menjadi "Ya"/"Tidak" berarti mengarang domain yang
 * speknya tidak punya - dan nilai di luar peta itu akan hilang dari layar.
 * Yang diwarnai hanya `'1'`, satu-satunya nilai yang repo ini sendiri tulis.
 */
function LencanaAktif({ nilai }: { nilai: string }) {
  const hidup = nilai === '1'
  return (
    <span className={'trin__lencana' + (hidup ? ' trin__lencana--aktif' : '')}>{nilai}</span>
  )
}

export default function AcuanTreatyIn() {
  const [himpunan, setHimpunan] = useState<HimpunanAcuan>('jenis-potongan')
  const [baris, setBaris] = useState<Acuan[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)

  useEffect(() => {
    let dibuang = false
    setBaris(null)
    setGalat(null)
    ambilAcuan(himpunan)
      .then((h) => {
        if (!dibuang) setBaris(h)
      })
      .catch((e: unknown) => {
        if (!dibuang) setGalat(e)
      })
    return () => {
      dibuang = true
    }
  }, [himpunan])

  const bersusun = himpunan === 'jenis-reasuransi'

  return (
    <div className="inbox">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{MENU_TREATYIN.acuan}</h2>
      </header>
      <StripTab
        tab={HIMPUNAN_ACUAN}
        aktif={himpunan}
        onPilih={setHimpunan}
        label={(t) => LABEL_HIMPUNAN[t]}
      />
      {/* Judul panel = himpunan yang sedang terbuka. Tanpa itu, isi tabel dan
          tab yang dipilih hanya terhubung lewat ingatan pembacanya. */}
      {galat !== null && <Gagal galat={galat} />}

      {/* ⛔ Disamakan dengan kedua layar baru (ronde layar 2): kepala kolom
          SELALU dirender, termasuk saat nol baris. Sebelumnya seluruh
          <table> disembunyikan ketika kosong, sehingga rancangan tabelnya
          tidak terlihat sama sekali — dan satu modul yang memperlihatkan
          kolomnya di dua layar lalu menyembunyikannya di layar ketiga punya
          dua bahasa visual. */}
      <Panel judul={LABEL_HIMPUNAN[himpunan]}>
        <div className="table-wrap">
          <table className="trin__tabel trin__tabel--pega">
            <thead>
              <tr>
                <th scope="col">{ACUAN_TREATYIN.kolomKode}</th>
                <th scope="col">{ACUAN_TREATYIN.kolomNama}</th>
                <th scope="col">{ACUAN_TREATYIN.kolomAktif}</th>
                {bersusun && <th scope="col">{ACUAN_TREATYIN.kolomInduk}</th>}
              </tr>
            </thead>
            <tbody>
              {baris === null && galat === null && (
                <tr>
                  <td colSpan={bersusun ? 4 : 3}>
                    <Memuat />
                  </td>
                </tr>
              )}
              {baris !== null && baris.length === 0 && (
                <tr>
                  <td colSpan={bersusun ? 4 : 3}>
                    <Kosong pesan={ACUAN_TREATYIN.kosong} petunjuk={ACUAN_TREATYIN.kosongPetunjuk} />
                  </td>
                </tr>
              )}
              {baris?.map((b) => (
                <tr key={b.id}>
                  <td>{b.kode}</td>
                  <td>{b.nama}</td>
                  <td>
                    <LencanaAktif nilai={b.aktif} />
                  </td>
                  {bersusun && <td>{b.idInduk ?? ''}</td>}
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </Panel>
    </div>
  )
}
