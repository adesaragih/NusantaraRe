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
// ⛔ NOMOR LAHIR SEKALI. Tombolnya mati begitu polis bernomor, dan layar
// MENGATAKAN sebabnya — tombol yang tetap hidup tetapi selalu menjawab hal
// yang sama mengajari orang mengabaikan jawabannya.
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
import {
  ambilKepalaPolis,
  ambilPesertaPolis,
  terbitkanNomorPL,
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
  const [terbit, setTerbit] = useState(false)
  const [kabar, setKabar] = useState('')

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
  const tanpaPeserta = hal !== null && hal.total === 0

  async function terbitkan(): Promise<void> {
    if (terbit || bernomor) return
    setTerbit(true)
    setGalat(null)
    setKabar('')
    try {
      const hasil = await terbitkanNomorPL(polisID)
      setKepala((lama) => (lama === null ? lama : { ...lama, plNumber: hasil.nomor }))
      setKabar(
        hasil.baruTerbit
          ? `Nomor terbit untuk periode ${hasil.periode} pada ${String(hasil.barisPeserta)} baris peserta.`
          : DETAIL_POLIS.sudahBernomor,
      )
      // Grid ikut disegarkan: setiap barisnya kini memuat nomornya.
      setHal(await ambilPesertaPolis(polisID))
    } catch (e) {
      setGalat(e)
    } finally {
      setTerbit(false)
    }
  }

  return (
    <section className="pl-detail">
      <header className="pl-detail__kepala">
        <h2 className="pl-detail__judul">{DETAIL_POLIS.judul}</h2>
        {kepala !== null && (
          <p className="pl-detail__identitas">
            {kepala.polisId} — {selPeserta(kepala.type)} /{' '}
            {selPeserta(kepala.businessCode)}
          </p>
        )}
        <p className="pl-detail__nomor">
          {DETAIL_POLIS.nomor}: <strong>{kalimatNomor(kepala)}</strong>
        </p>
        <button
          type="button"
          className="pl-detail__terbitkan"
          disabled={sibuk || terbit || bernomor || tanpaPeserta}
          onClick={() => {
            void terbitkan()
          }}
        >
          {DETAIL_POLIS.terbitkan}
        </button>
        {/* ⛔ Sebab tombolnya mati DIKATAKAN, bukan dibiarkan ditebak. */}
        {bernomor && <p role="status">{DETAIL_POLIS.sudahBernomor}</p>}
        {!bernomor && tanpaPeserta && <p role="status">{DETAIL_POLIS.perluPeserta}</p>}
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

      {sibuk && <Memuat />}
      {galat !== null && <Gagal galat={galat} />}
      {kabar !== '' && <p role="status">{kabar}</p>}

      {hal !== null && hal.baris.length === 0 && (
        <Kosong pesan="Polis ini belum punya baris peserta." />
      )}

      {hal !== null && hal.baris.length > 0 && (
        <table className="pl-detail__tabel">
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
      )}

      {hal !== null && (
        <p className="pl-detail__cacah" role="status">
          {hal.baris.length} dari {hal.total}
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
            {kepala.medanTanpaKolom.length} medan layar lama tidak ditampilkan
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
  )
}
