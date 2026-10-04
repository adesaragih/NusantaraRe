// Kasus Komite — tiket 01 (baca) + tiket 02 (keputusan).
//
// Tangga persetujuan kasus dan baris yang diputuskan. Formulir keputusan
// meniru bagian keputusan `Section/ShowTransfer.xml` (lewat
// `FlowAction/ViewTransferDtl.xml` b91): dropdown `.AcceptStatus` WAJIB
// (b32607), `.KomiteComment` (b31001), tombol `Submit` b34722 / `Cancel`
// b33880. Formulir hanya tampil pada GILIRAN pelaku.
//
// ⚠️ Syarat tampil lain `ShowTransfer` (`IsTreatyIn==0` ×12, `Type TP/TR` ×3,
// `SwiftCode`, `RetrocadedShare`) menggerbangi blok RINCIAN, belum dibawa.

import { useCallback, useEffect, useState } from 'react'

import { PERAN, type KodePeran } from '../../../../inti/frontend/labels'
import {
  EFEK_KOMITE,
  ESKALASI_KOMITE,
  KASUS_KOMITE,
  KEPUTUSAN_KOMITE,
  KOLOM_INBOX_KOMITE,
} from '../labels'
import { Gagal, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import {
  ambilKasusKomite,
  ambilRiwayatKomite,
  eskalasiKomite,
  putuskanKomite,
  type KasusKomite as Kasus,
  type KeputusanAsliKomite,
  type RiwayatKomite,
} from '../api'
import { selKomite, tingkatKomite } from './InboxKomite'

/** Enum tertutup `{1 Setuju, 2 Tolak}` — urut dropdown. */
export const PILIHAN_KEPUTUSAN = [
  { kode: '1', label: KEPUTUSAN_KOMITE.setuju },
  { kode: '2', label: KEPUTUSAN_KOMITE.tolak },
] as const

/**
 * Keputusan asli tingkat yang tertimpa langkah 5.1 — MURNI (OQ-K-05,
 * GILIRAN-17): `Setuju — komentar`, atau tanda kosong bila tidak tertimpa.
 */
export function teksKeputusanAsli(asli: KeputusanAsliKomite | undefined): string {
  if (asli === undefined) return '—'
  return asli.comment.trim() === '' ? asli.status : `${asli.status} — ${asli.comment}`
}

/**
 * Kontrol eskalasi hanya untuk admin, dan hanya bila masih ada tingkat di atas.
 *
 * ⚠️ `[asumsi — OQ-007/OQ-021]` admin komite = `ReasLifeAdmin`. Server tetap
 * yang menegakkan; ini hanya menyembunyikan tombol yang pasti ditolak.
 */
export function bolehEskalasi(peran: readonly KodePeran[], k: Kasus | null): boolean {
  if (k === null) return false
  return peran.includes(PERAN.admin) && k.kasus.tingkatBerjalan > 0 &&
    k.kasus.tingkatBerjalan < k.kasus.komiteLoop
}

/** Kalimat hasil — satu tempat. */
export function kalimatHasilKeputusan(h: {
  kataKeputusan: string
  tingkatDiputus: number
  berlanjut: boolean
  tingkatBerikut: number
  nomorAkseptasi?: string
}): string {
  const dasar = `${h.kataKeputusan} tercatat di tingkat ${String(h.tingkatDiputus)}.`
  if (h.berlanjut) return `${dasar} Kasus naik ke tingkat ${String(h.tingkatBerikut)}.`
  const nomor = h.nomorAkseptasi ?? ''
  return nomor === '' ? `${dasar} Tangga berhenti.` : `${dasar} Nomor akseptasi ${nomor}.`
}

export default function KasusKomite({
  kasusID,
  peran,
  onKembali,
}: {
  kasusID: string
  /** Peran pelaku — hanya untuk menampilkan kontrol eskalasi. */
  peran: readonly KodePeran[]
  onKembali: () => void
}) {
  const [k, setK] = useState<Kasus | null>(null)
  const [riwayat, setRiwayat] = useState<RiwayatKomite | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [sibuk, setSibuk] = useState(true)
  const [keputusan, setKeputusan] = useState('')
  const [komentar, setKomentar] = useState('')
  const [kirim, setKirim] = useState(false)
  const [kabar, setKabar] = useState('')

  const muat = useCallback(async () => {
    setSibuk(true)
    setGalat(null)
    try {
      const [kasus, rw] = await Promise.all([ambilKasusKomite(kasusID), ambilRiwayatKomite(kasusID)])
      setK(kasus)
      setRiwayat(rw)
    } catch (e) {
      // ⚠️ 403 bila pelaku bukan anggota tangga — pesan server tampil apa adanya.
      setGalat(e)
      setK(null)
    } finally {
      setSibuk(false)
    }
  }, [kasusID])

  useEffect(() => {
    void muat()
  }, [muat])

  async function eskalasi(): Promise<void> {
    if (kirim) return
    setKirim(true)
    setGalat(null)
    setKabar('')
    try {
      const h = await eskalasiKomite(kasusID)
      setKabar(`Eskalasi dari tingkat ${String(h.dariTingkat)} ke tingkat ${String(h.keTingkat)}.`)
      await muat()
    } catch (e) {
      setGalat(e)
    } finally {
      setKirim(false)
    }
  }

  async function submit(): Promise<void> {
    if (kirim) return
    if (keputusan === '') {
      setKabar(KEPUTUSAN_KOMITE.wajibPilih)
      return
    }
    setKirim(true)
    setGalat(null)
    setKabar('')
    try {
      const h = await putuskanKomite(kasusID, keputusan, komentar)
      setKabar(
        kalimatHasilKeputusan(h) +
          (h.efekTertunda.length > 0
            ? ` Efek keluar diantre, belum tuntas: ${h.efekTertunda.join(', ')}.`
            : ''),
      )
      setKeputusan('')
      setKomentar('')
      await muat()
    } catch (e) {
      setGalat(e)
    } finally {
      setKirim(false)
    }
  }

  return (
    <section className="komite-kasus">
      <button type="button" onClick={onKembali}>
        {KASUS_KOMITE.kembali}
      </button>
      <h2>
        {KASUS_KOMITE.judul} {kasusID}
      </h2>
      {sibuk && <Memuat />}
      {galat !== null && <Gagal galat={galat} />}
      {kabar !== '' && <p role="status">{kabar}</p>}
      {k !== null && (
        <>
          <dl className="komite-kasus__kepala">
            <dt>{KOLOM_INBOX_KOMITE.nomorKlaim}</dt>
            <dd>{selKomite(k.kasus.nomorKlaim)}</dd>
            <dt>{KOLOM_INBOX_KOMITE.tingkat}</dt>
            <dd>{tingkatKomite(k.kasus)}</dd>
            <dt>{KOLOM_INBOX_KOMITE.nilaiKlaim}</dt>
            <dd>
              {selKomite(k.kasus.nilaiKlaim)} {k.kasus.mataUang}
            </dd>
            <dt>{KOLOM_INBOX_KOMITE.statusBaris}</dt>
            <dd>{k.kasus.statusBaris}</dd>
          </dl>
          <p role="status">
            {k.giliranSaya ? KASUS_KOMITE.giliranAnda : KASUS_KOMITE.bukanGiliran}
          </p>
          <h3>{KASUS_KOMITE.tangga}</h3>
          {riwayat !== null && (
            <table className="komite-kasus__tangga">
              <thead>
                <tr>
                  <th>{KASUS_KOMITE.urut}</th>
                  <th>{KASUS_KOMITE.committee}</th>
                  <th>{KASUS_KOMITE.anggota}</th>
                  <th>{KASUS_KOMITE.status}</th>
                  <th>{KASUS_KOMITE.dateApprove}</th>
                  <th>{KASUS_KOMITE.comment}</th>
                  <th>{KASUS_KOMITE.asli}</th>
                </tr>
              </thead>
              <tbody>
                {riwayat.tangga.map((a) => (
                  <tr key={a.urut}>
                    <td>{a.urut}</td>
                    <td>{selKomite(a.committee)}</td>
                    <td>{selKomite(a.anggota)}</td>
                    <td>{a.status}</td>
                    <td>{selKomite(a.dateApprove)}</td>
                    <td>{selKomite(a.comment)}</td>
                    <td>{teksKeputusanAsli(a.asli)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
          {riwayat !== null && riwayat.eskalasi.length > 0 && (
            <>
              <h3>{KASUS_KOMITE.eskalasi}</h3>
              <ul className="komite-kasus__eskalasi-riwayat">
                {riwayat.eskalasi.map((e, i) => (
                  <li key={`${e.waktu}-${String(i)}`}>
                    Tingkat {e.dariTingkat} → {e.keTingkat} oleh {e.oleh} pada {e.waktu}
                  </li>
                ))}
              </ul>
            </>
          )}
          {k.efek.efek.length > 0 && (
            <>
              <h3>
                {EFEK_KOMITE.judul} — {k.efek.keadaan}
              </h3>
              <table className="komite-kasus__efek">
                <thead>
                  <tr>
                    <th>{EFEK_KOMITE.jenis}</th>
                    <th>{EFEK_KOMITE.keadaan}</th>
                    <th>{EFEK_KOMITE.percobaan}</th>
                    <th>{EFEK_KOMITE.sejak}</th>
                  </tr>
                </thead>
                <tbody>
                  {k.efek.efek.map((e, i) => (
                    <tr key={`${e.jenis}-${String(i)}`}>
                      <td>{e.jenis}</td>
                      <td>{e.keadaan}</td>
                      <td>{e.percobaan}</td>
                      <td>{e.sejak}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </>
          )}
          {bolehEskalasi(peran, k) && (
            <p className="komite-kasus__eskalasi">
              <button
                type="button"
                disabled={kirim}
                onClick={() => {
                  void eskalasi()
                }}
              >
                {ESKALASI_KOMITE.tombol}
              </button>{' '}
              {ESKALASI_KOMITE.keterangan}
            </p>
          )}
          {k.giliranSaya && (
            <form
              className="komite-kasus__keputusan"
              onSubmit={(e) => {
                e.preventDefault()
                void submit()
              }}
            >
              <label>
                {KEPUTUSAN_KOMITE.konfirmasi}
                <select
                  required
                  value={keputusan}
                  onChange={(e) => {
                    setKeputusan(e.target.value)
                  }}
                >
                  <option value="">{KEPUTUSAN_KOMITE.pilih}</option>
                  {PILIHAN_KEPUTUSAN.map((p) => (
                    <option key={p.kode} value={p.kode}>
                      {p.label}
                    </option>
                  ))}
                </select>
              </label>
              <label>
                {KEPUTUSAN_KOMITE.komentar}
                <textarea
                  value={komentar}
                  onChange={(e) => {
                    setKomentar(e.target.value)
                  }}
                />
              </label>
              <button type="submit" disabled={kirim}>
                {KEPUTUSAN_KOMITE.submit}
              </button>
              <button type="button" onClick={onKembali}>
                {KEPUTUSAN_KOMITE.cancel}
              </button>
            </form>
          )}
        </>
      )}
    </section>
  )
}
