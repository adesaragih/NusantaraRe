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

import { DETAIL_POLIS, JUDUL_KOLOM_PESERTA, KOLOM_PALING_KANAN, URUTAN_KOLOM_PESERTA } from '../labels'
import { Gagal, Kosong, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import FormDataPolis from './FormDataPolis'
import UnggahCSVPeserta from './UnggahCSVPeserta'
import { tanggalTampil } from '../tanggal'
import { periodeTampil } from './InputOffer'
import '../premiumlistlife.css'
import {
  ambilKepalaPolis,
  ambilPeriodeProduksi,
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

/**
 * Kolom grid yang BERISI di setidaknya satu baris yang tampil — permintaan
 * work owner 02-10-2026 ("tampilkan saja kolom yang ada isinya").
 *
 * ⚠️ Urutan server dipertahankan. Kolom yang SELURUHNYA nol juga
 * disembunyikan (keputusan work owner 02-10-2026) — kolom uang yang tidak ada
 * di CSV diisi 0 (`IsiNolUangKosong`). Hanya penyembunyian KOLOM: sel nol di
 * kolom yang tampil tetap ditulis `0` (`selPeserta`). Jumlah kolom yang
 * disembunyikan DINYATAKAN di layar, supaya kolom yang hilang bukan data yang
 * hilang tanpa tanda.
 */
export function kolomBerisi(
  kolom: readonly string[],
  baris: readonly { nilai: Record<string, string | undefined> }[],
): string[] {
  return kolom.filter((k) => baris.some((b) => adaIsi(b.nilai[k])))
}

/**
 * Menyusun kolom grid: `URUTAN_KOLOM_PESERTA` lebih dahulu, lalu kolom lain
 * dengan urutan server, lalu `KOLOM_PALING_KANAN` (STNC, WPC) — permintaan
 * work owner 02-10-2026.
 *
 * ⛔ Daftar urutnya hanya KUNCI URUT, bukan daftar kolom: kolom yang ada tetap
 * ditentukan server (`hal.kolom`), yang tersembunyi karena kosong tidak
 * dimunculkan kembali, dan kolom yang tidak dikenal daftar urut tetap tampil.
 */
export function susunKolom(kolom: readonly string[]): string[] {
  const kanan: readonly string[] = KOLOM_PALING_KANAN
  const depan = URUTAN_KOLOM_PESERTA.filter((k) => kolom.includes(k))
  const tengah = kolom.filter((k) => !depan.includes(k) && !kanan.includes(k))
  return [...depan, ...tengah, ...kanan.filter((k) => kolom.includes(k))]
}

/** Teks angka nol dalam bentuk apa pun: `0`, `0.00`, `-0`, `.0`. */
const POLA_NOL = /^[+-]?(0+\.?0*|\.0+)$/

/** Isi yang layak memunculkan kolom: tidak kosong dan bukan nol. */
function adaIsi(nilai: string | undefined): boolean {
  const s = (nilai ?? '').trim()
  return s !== '' && !POLA_NOL.test(s)
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
  // Periode produksi di kepala, seperti Input Offer Life (02-10-2026). Galatnya
  // DINYATAKAN: 503 berarti POOLDATA.TANGGAL_CLOSING kosong.
  const [periode, setPeriode] = useState('')
  const [galatPeriode, setGalatPeriode] = useState<unknown>(null)
  useEffect(() => {
    let hidup = true
    ambilPeriodeProduksi().then(
      (p) => {
        if (hidup) setPeriode(p)
      },
      (e: unknown) => {
        if (hidup) setGalatPeriode(e)
      },
    )
    return () => {
      hidup = false
    }
  }, [])

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
  // Hanya kolom yang berisi (02-10-2026); sisanya dihitung dan dinyatakan.
  const tampilKolom = hal === null ? [] : susunKolom(kolomBerisi(hal.kolom, hal.baris))
  const tersembunyi = hal === null ? 0 : hal.kolom.length - tampilKolom.length

  return (
    <section className="pl-detail">
      {/*
        Kepala ringkas — pola layar Input Offer Life: judul, lalu HANYA nomor
        kasus (NBLF) dan PL Number (permintaan work owner 02-10-2026; Type dan
        COB sudah tampil di form di bawahnya). PL Number yang belum ada
        DIKATAKAN, dengan gaya redup.
      */}
      <header className="pl-kepala">
        <h2 className="pl-kepala__judul">{DETAIL_POLIS.judul}</h2>
        {kepala !== null && (
          <div className="pl-kepala__meta">
            <span className="pl-kepala__chip pl-kepala__chip--utama">{kepala.polisId}</span>
            {periode !== '' && (
              <span className="pl-kepala__chip" role="status">
                <span className="pl-kepala__label">Period</span> {periodeTampil(periode)}
              </span>
            )}
            <span className="pl-kepala__chip">
              <span className="pl-kepala__label">{DETAIL_POLIS.nomor}</span>
              <span className={bernomor ? undefined : 'pl-kepala__belum'}>{kalimatNomor(kepala)}</span>
            </span>
          </div>
        )}
      </header>
      {galatPeriode !== null && <Gagal galat={galatPeriode} />}

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
                  {tampilKolom.map((k) => (
                    <th key={k}>{judulKolom(k)}</th>
                  ))}
                </tr>
              </thead>
              <tbody>
                {hal.baris.map((b) => (
                  <tr key={b.id}>
                    {tampilKolom.map((k) => (
                      <td key={k}>{selPeserta(tanggalTampil(b.nilai[k]))}</td>
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
            {tersembunyi > 0 && ` · ${tersembunyi} empty columns hidden`}
          </p>
        )}

        {/*
          ⛔ Selisih kolom grid vs `PL_Detail_Sec` (REINSTYPENAME,
          RetrocadedShare) SENGAJA TIDAK ditampilkan lagi — keputusan work owner
          02-10-2026: catatan pengembang, bukan informasi pemakai. Jawabannya
          tetap tercatat di `models.MedanGridTanpaKolom` dan tetap dikirim
          server (`kepala.medanTanpaKolom`).
        */}
      </section>
    </section>
  )
}
