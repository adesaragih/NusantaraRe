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

import { PERAN, type KodePeran } from '../../assets/labels.claimlife'
import {
  ESKALASI_KOMITE,
  KASUS_KOMITE,
  KEPUTUSAN_KOMITE,
  KOLOM_INBOX_KOMITE,
} from '../../assets/labels.komite'
import { Gagal, Memuat } from '../../components/ui/dasar'
import {
  ambilKasusKomite,
  eskalasiKomite,
  putuskanKomite,
  type KasusKomite as Kasus,
} from '../../services/api'
import { selKomite, tingkatKomite } from './InboxKomite'

/** Enum tertutup `{1 Setuju, 2 Tolak}` — urut dropdown. */
export const PILIHAN_KEPUTUSAN = [
  { kode: '1', label: KEPUTUSAN_KOMITE.setuju },
  { kode: '2', label: KEPUTUSAN_KOMITE.tolak },
] as const

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
      setK(await ambilKasusKomite(kasusID))
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
      setKabar(kalimatHasilKeputusan(h))
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
          <table className="komite-kasus__tangga">
            <thead>
              <tr>
                <th>{KASUS_KOMITE.urut}</th>
                <th>{KASUS_KOMITE.jabatan}</th>
                <th>{KASUS_KOMITE.approval}</th>
                <th>{KASUS_KOMITE.komentar}</th>
                <th>{KASUS_KOMITE.tanggal}</th>
              </tr>
            </thead>
            <tbody>
              {k.tangga.map((a) => (
                <tr key={a.urut} className={a.saya ? 'komite-kasus__saya' : undefined}>
                  <td>{a.urut}</td>
                  <td>{selKomite(a.jabatan)}</td>
                  <td>{a.kataApproval}</td>
                  <td>{selKomite(a.komentar)}</td>
                  <td>{selKomite(a.tglApprove)}</td>
                </tr>
              ))}
            </tbody>
          </table>
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
