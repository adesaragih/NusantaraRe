// Kotak masuk PremiumList Life — tiket 01 bagian 2.
//
// Meniru `ReportDefinition/InboxPremiumList.xml` dan portal
// `Section/PremiumList.xml`.
//
// ⛔ SEBELAS KOLOM, dan dua yang sengaja tidak ada — lihat
// `modul/premiumlistlife/frontend/labels.ts`. `KetentuanUnderwriting` b891 tidak punya
// kolom di migrasi mana pun, dan sel kosong di layar terbaca "memang kosong"
// alih-alih "kami tidak punya datanya".
//
// ⛔ DUA TOMBOL PORTAL MEMBUAT KASUS (GILIRAN-13 butir bn). `Input Offer`
// b3273 dan `Input Premium` b3921 memanggil `CreateInputLife` b3291/b3939
// dengan `FlagPolicy` "0"/"1", yang:
//
//   b444 `Call svcAddWorkObject`           -> membuat work object baru
//   b618 `curWorkPage.FlagOnGoingPolicy = Param.FlagPolicy`
//   b726 `Obj-Save`
//   b982 "ASSIGN-WORKLIST <pzInsKey>!InputPolicyHolder" -> assignment
//        tahap pertamanya langsung dibuka
//
// ⚠️ RALAT 29-09-2026: komentar lama berkata `SEQ_WORK_POLIS` menunggu
// keputusan skema. Ia SUDAH diputuskan di pl3 (brief modul PremiumList) dan
// terlewat; migrasi 057 kini membuatnya bersama kolom bendera.

import { useCallback, useEffect, useState } from 'react'

import { KOLOM_INBOX_POLIS, TOMBOL_POLIS } from '../labels'
import { Gagal, Kosong, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { unduhXlsx, type KolomEksporXlsx } from '../../../../inti/frontend/lib/exportXlsx'
import {
  ambilKotakMasukPolis,
  buatKasusPolis,
  FLAG_POLIS,
  type BarisInboxPolis,
  type FlagPolis,
  type HalamanInboxPolis,
} from '../api'

/**
 * Kolom ekspor — TEPAT yang tampil, berurutan sama.
 *
 * ⚠️ Tanggal keluar apa adanya (`YYYY-MM-DD`), yang sudah terurut benar
 * sebagai teks — berbeda dengan kotak masuk Claim Life, yang menampilkan
 * `DD-MM-YYYY` dan karena itu harus mengekspor bentuk ISO-nya.
 */
export const KOLOM_EKSPOR_POLIS: KolomEksporXlsx<BarisInboxPolis>[] = [
  { kunci: 'caseId', label: KOLOM_INBOX_POLIS.caseId },
  { kunci: 'tglCreate', label: KOLOM_INBOX_POLIS.tglCreate },
  { kunci: 'createOpName', label: KOLOM_INBOX_POLIS.createOpName },
  { kunci: 'statusWork', label: KOLOM_INBOX_POLIS.statusWork },
  { kunci: 'cedingCoName', label: KOLOM_INBOX_POLIS.cedingCoName },
  { kunci: 'policyHolderName', label: KOLOM_INBOX_POLIS.policyHolderName },
  { kunci: 'plNumber', label: KOLOM_INBOX_POLIS.plNumber },
  { kunci: 'riSlipRnm', label: KOLOM_INBOX_POLIS.riSlipRnm },
  { kunci: 'type', label: KOLOM_INBOX_POLIS.type },
  { kunci: 'marketingName', label: KOLOM_INBOX_POLIS.marketingName },
  { kunci: 'sobName', label: KOLOM_INBOX_POLIS.sobName },
  { kunci: 'dateReceived', label: KOLOM_INBOX_POLIS.dateReceived },
]

/** Sel kosong ditandai, bukan dibiarkan kosong (ADR-U-0027). */
export function sel(nilai: string): string {
  return nilai.trim() === '' ? '—' : nilai
}

export default function InboxPremiumList({
  onBuka,
}: {
  /** Membuka satu polis. */
  onBuka: (caseID: string, tahap: string) => void
}) {
  const [hal, setHal] = useState<HalamanInboxPolis | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [sibuk, setSibuk] = useState(true)
  const [membuat, setMembuat] = useState(false)
  const [galatBuat, setGalatBuat] = useState<unknown>(null)

  /** `CreateInputLife` — lalu kasusnya langsung dibuka di tahap pertamanya. */
  async function buat(flag: FlagPolis): Promise<void> {
    if (membuat) return
    setMembuat(true)
    setGalatBuat(null)
    try {
      const hasil = await buatKasusPolis(flag)
      onBuka(hasil.caseId, hasil.statusWork)
    } catch (e) {
      setGalatBuat(e)
    } finally {
      setMembuat(false)
    }
  }

  const muat = useCallback(async () => {
    setSibuk(true)
    setGalat(null)
    try {
      setHal(await ambilKotakMasukPolis())
    } catch (e) {
      // ⛔ Galat DINYATAKAN, bukan menjadi daftar kosong. Daftar kosong
      // sesudah gagal memuat terbaca "tidak ada pekerjaan".
      setGalat(e)
      setHal(null)
    } finally {
      setSibuk(false)
    }
  }, [])

  useEffect(() => {
    void muat()
  }, [muat])

  return (
    <section className="inbox">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">PremiumList</h2>
        <button
          type="button"
          disabled={membuat}
          onClick={() => {
            void buat(FLAG_POLIS.inputOffer)
          }}
        >
          {TOMBOL_POLIS.inputOffer}
        </button>{' '}
        <button
          type="button"
          disabled={membuat}
          onClick={() => {
            void buat(FLAG_POLIS.inputPremium)
          }}
        >
          {TOMBOL_POLIS.inputPremium}
        </button>{' '}
        <button
          type="button"
          className="inbox__ekspor"
          disabled={hal === null || hal.baris.length === 0}
          onClick={() => {
            if (hal === null) return
            unduhXlsx('premiumlist.xlsx', KOLOM_EKSPOR_POLIS, hal.baris)
          }}
        >
          Export xlsx
        </button>
      </header>

      {/* Kasus yang GAGAL lahir dinyatakan — tombol yang diam tanpa sebab
          terbaca "tidak terjadi apa-apa". */}
      {galatBuat !== null && <Gagal galat={galatBuat} />}
      {sibuk && <Memuat />}
      {galat !== null && <Gagal galat={galat} />}

      {hal !== null && hal.baris.length === 0 && (
        <Kosong pesan="Antrean PremiumList kosong." />
      )}

      {hal !== null && hal.baris.length > 0 && (
        <table className="inbox__tabel">
          <thead>
            <tr>
              {KOLOM_EKSPOR_POLIS.map((k) => (
                <th key={String(k.kunci)}>{k.label}</th>
              ))}
            </tr>
          </thead>
          <tbody>
            {hal.baris.map((b) => (
              <tr
                key={b.caseId}
                className="inbox__baris"
                onClick={() => {
                  onBuka(b.caseId, b.statusWork)
                }}
              >
                {KOLOM_EKSPOR_POLIS.map((k) => (
                  <td key={String(k.kunci)}>{sel(b[k.kunci])}</td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      )}

      {hal !== null && (
        <p className="inbox__cacah" role="status">
          {hal.baris.length} dari {hal.total}
        </p>
      )}
    </section>
  )
}
