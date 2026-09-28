// Kotak masuk PremiumList Life — tiket 01 bagian 2.
//
// Meniru `ReportDefinition/InboxPremiumList.xml` dan portal
// `Section/PremiumList.xml`.
//
// ⛔ SEBELAS KOLOM, dan dua yang sengaja tidak ada — lihat
// `assets/labels.premiumlist.ts`. `KetentuanUnderwriting` b891 tidak punya
// kolom di migrasi mana pun, dan sel kosong di layar terbaca "memang kosong"
// alih-alih "kami tidak punya datanya".
//
// ⛔ DUA TOMBOL PORTAL BELUM DAPAT MEMBUAT KASUS, dan sebabnya dinyatakan,
// bukan disembunyikan. `Input Offer` b3273 dan `Input Premium` b3921
// memanggil `CreateInputLife` b3291/b3939, yang:
//
//   b444 `Call svcAddWorkObject`           -> membuat work object baru
//   b618 `curWorkPage.FlagOnGoingPolicy = Param.FlagPolicy`
//   b726 `Obj-Save`
//
// Dua hal yang belum ada di skema kami: `SEQ_WORK_POLIS` (tidak dibuat
// migrasi 050-056 mana pun) dan kolom untuk `FlagOnGoingPolicy`. Keduanya
// keputusan skema, dan migrasi baru hanya dari keputusan yang TERCATAT —
// jadi tombolnya berdiri dan MENYEBUT apa yang ditunggunya.
//
// ⚠️ Berdiri, bukan disembunyikan: tombol yang hilang membuat layar tampak
// lengkap padahal alurnya belum dapat dimulai.

import { useCallback, useEffect, useState } from 'react'

import { KOLOM_INBOX_POLIS, TOMBOL_POLIS } from '../../assets/labels.premiumlist'
import { BelumTersedia, Gagal, Kosong, Memuat } from '../../components/ui/dasar'
import { unduhXlsx, type KolomEksporXlsx } from '../../lib/exportXlsx'
import {
  ambilKotakMasukPolis,
  type BarisInboxPolis,
  type HalamanInboxPolis,
} from '../../services/api'

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
        {/* ⛔ Keduanya DINYATAKAN, bukan dihilangkan — lihat kepala berkas. */}
        <BelumTersedia apa={TOMBOL_POLIS.inputOffer} />{' '}
        <BelumTersedia apa={TOMBOL_POLIS.inputPremium} />{' '}
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
