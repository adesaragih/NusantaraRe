// Layar Premium List Detail — tiket 03.
//
// Meniru `Section/ShowLifePremiumDetail.xml` (kepala polis + nomor PL) dan
// `Section/PL_Detail_Sec.xml` (grid peserta).
//
// ⛔ DUA MEDAN GRID LAMA TIDAK DITAMPILKAN, dan ketiadaannya DINYATAKAN:
// `REINSTYPENAME` dan `RetrocadedShare` tidak punya kolom di migrasi 050–056
// mana pun. Sel kosong di layar terbaca "memang kosong"; kolom yang tidak
// pernah ada terbaca "datanya hilang". Keduanya salah, dan yang kedua membuat
// orang mencari data yang tidak pernah kami punya.
//
// ⛔ TANPA TOMBOL GENERATE PL NUMBER (keputusan work owner 01-10-2026). Di
// `ShowLifePremiumDetail` sel PL_NUMBER pun ber-`pyVisible never`: nomor PL
// terbit saat polis DISIMPAN (Confirm tahap ini / Submit summary, tiket 05a),
// bukan lewat tombol. Layar hanya menampilkan nomornya bila sudah ada.
//
// ⛔ UANG TETAP TEKS. Tidak satu pun sel melewati `Number(...)`: premi
// delapan angka desimal dibulatkan diam-diam olehnya (ADR-U-0003).
//
// ⚠️ Daftar kolomnya DATANG DARI SERVER, tidak diketik ulang di sini.

import { useCallback, useEffect, useState } from 'react'

import { DETAIL_POLIS, JUDUL_KOLOM_PESERTA } from '../labels'
import { Gagal, Kosong, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import FormDataPolis from './FormDataPolis'
import UnggahCSVPeserta from './UnggahCSVPeserta'
import '../premiumlistlife.css'
import {
  ambilKepalaPolis,
  ambilPesertaPolis,
  type HalamanPesertaPolis,
  type KepalaPolis,
} from '../api'

/**
 * Judul satu kolom: yang dikenal diperindah, yang tidak tampil apa adanya.
 *
 * ⚠️ Kolom asing TIDAK disembunyikan. Kolom yang hilang dari layar karena
 * judulnya belum terdaftar adalah data yang hilang tanpa satu pun tanda.
 */
export function judulKolom(nama: string): string {
  return JUDUL_KOLOM_PESERTA[nama] ?? nama
}

/** Sel kosong ditandai, bukan dibiarkan kosong (ADR-U-0027). */
export function selPeserta(nilai: string | undefined): string {
  return nilai === undefined || nilai.trim() === '' ? '—' : nilai
}

/** Kalimat keadaan nomor — satu tempat, supaya layar tidak mengarang. */
export function kalimatNomor(kepala: KepalaPolis | null): string {
  if (kepala === null) return ''
  return kepala.plNumber.trim() === ''
    ? DETAIL_POLIS.belumBernomor
    : kepala.plNumber
}

export default function PremiumListDetail({ polisID }: { polisID: string }) {
  const [kepala, setKepala] = useState<KepalaPolis | null>(null)
  const [hal, setHal] = useState<HalamanPesertaPolis | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [sibuk, setSibuk] = useState(true)

  const muat = useCallback(async () => {
    setSibuk(true)
    setGalat(null)
    try {
      const [k, p] = await Promise.all([
        ambilKepalaPolis(polisID),
        ambilPesertaPolis(polisID),
      ])
      setKepala(k)
      setHal(p)
    } catch (e) {
      // ⛔ Galat DINYATAKAN, bukan menjadi layar kosong. Layar kosong sesudah
      // gagal memuat terbaca "polis ini memang tidak punya peserta".
      setGalat(e)
      setKepala(null)
      setHal(null)
    } finally {
      setSibuk(false)
    }
  }, [polisID])

  useEffect(() => {
    void muat()
  }, [muat])

  const bernomor = kepala !== null && kepala.plNumber.trim() !== ''

  return (
    <section className="pl-detail">
      {/*
        Kepala ringkas — pola layar Input Offer Life: judul, lalu lencana nomor
        kasus, Type, Class of Business, dan PL_NUMBER (bila belum, dikatakan).
      */}
      <header className="pl-kepala">
        <h2 className="pl-kepala__judul">{DETAIL_POLIS.judul}</h2>
        {kepala !== null && (
          <div className="pl-kepala__meta">
            <span className="pl-kepala__chip">{kepala.polisId}</span>
            <span className="pl-kepala__chip">
              <span className="pl-kepala__label">Type</span> {selPeserta(kepala.type)}
            </span>
            <span className="pl-kepala__chip">
              <span className="pl-kepala__label">COB</span> {selPeserta(kepala.businessCode)}
            </span>
            <span className="pl-kepala__chip">
              <span className="pl-kepala__label">{DETAIL_POLIS.nomor}</span> {kalimatNomor(kepala)}
            </span>
          </div>
        )}
      </header>

      {/*
        Tiket 03 bagian 2 — tiga kolom atas `ShowLifePremiumDetail` (data polis,
        pihak, tanggal). Sesudah tersimpan kepala dimuat ulang: Type adalah bahan
        penomoran PL.
      */}
      <FormDataPolis
        polisID={polisID}
        bernomor={bernomor}
        onTersimpan={() => {
          void muat()
        }}
      />

      {/*
        ⛔ UNGGAHAN BERDIRI DI LAYAR YANG SAMA dengan gridnya, dan itu bentuk
        aslinya: `ViewCSVResult_LifePremium*` menampilkan hasil unggahan di
        konteks polis yang sedang dibuka. Layar terpisah memaksa orang
        mengingat polis mana yang sedang diunggahi.
      */}
      <UnggahCSVPeserta
        polisID={polisID}
        onTersimpan={() => {
          void muat()
        }}
      />

      <section className="panel">
        <h3 className="panel__title">Participants</h3>
        {sibuk && <Memuat />}
        {galat !== null && <Gagal galat={galat} />}

        {hal !== null && hal.baris.length === 0 && (
          <Kosong pesan="This policy has no participant rows yet." />
        )}

        {hal !== null && hal.baris.length > 0 && (
          <div className="pl-offer__riwayat-gulir">
            <table className="inbox__tabel">
              <thead>
                <tr>
                  {hal.kolom.map((k) => (
                    <th key={k}>{judulKolom(k)}</th>
                  ))}
                </tr>
              </thead>
              <tbody>
                {hal.baris.map((b) => (
                  <tr key={b.id}>
                    {hal.kolom.map((k) => (
                      <td key={k}>{selPeserta(b.nilai[k])}</td>
                    ))}
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}

        {hal !== null && hal.baris.length > 0 && (
          <p className="panel__note pl-detail__cacah" role="status">
            {hal.baris.length} of {hal.total}
          </p>
        )}

        {/*
          ⛔ SELISIH KOLOM DIJAWAB DI LAYAR, bukan hanya di komentar Go.
          `PL_Detail_Sec` menampilkan empat puluh medan; grid ini tiga puluh
          delapan. Siapa pun yang menghitungnya akan bertanya, dan jawaban yang
          hanya ada di kode bukan jawaban bagi yang bertanya.
        */}
        {kepala !== null && kepala.medanTanpaKolom.length > 0 && (
          <details className="pl-detail__absen">
            <summary>
              {kepala.medanTanpaKolom.length} legacy screen fields not shown
            </summary>
            <ul>
              {kepala.medanTanpaKolom.map((m) => (
                <li key={m.medan}>
                  <code>{m.medan}</code> — {m.alasan}
                </li>
              ))}
            </ul>
          </details>
        )}
      </section>
    </section>
  )
}
