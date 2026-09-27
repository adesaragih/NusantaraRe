// Layar Outstanding Claim — A3 kelompok Outstanding.
//
// Meniru `FlowAction/OSClaimLife.xml` → `Section/InputOSClaimLife.xml`.
//
// ⚠️ Kesebelas medan polisnya SAMA dengan layar Register dan terikat
// `.PolicyDataLife.*` pula - keduanya menampilkan polis yang sama pada tahap
// yang berbeda. Karena itu keduanya memakai `PanelDataPolis` yang sama, dan
// keduanya menunggu modul PremiumList Life (butir av). Menyalin panelnya
// berarti dua bentuk yang harus berubah bersama.
//
// ⛔ BUTIR aw - dua tombol perpindahan. Di Pega keduanya TIDAK menulis apa
// pun pada posisi Admin: `SendtoAdmin_Act` dan `SendtoAdmin_Act1` sama-sama
// berprasyarat `pyPosition=="ReasLifeMedicalAdvisor"` (WhenFalse=3 = lewati),
// padahal layar ini dipegang `ReasLifeAdmin`. Itu cacat rule warisan, dan
// yang ditiru MAKSUDnya - lihat `OQ-untuk-tim.md` OQ-C.

import { useState } from 'react'

import { OUTSTANDING, TOMBOL_OS } from '../assets/labels'
import { PanelDataPolis } from '../components/PanelDataPolis'
import { Gagal } from '../components/ui/dasar'
import { pesanGalat, pindahTahap, TAHAP_JALUR, type TahapJalur } from '../services/api'

export interface OutstandingProps {
  /** Pengenal work kasus yang sedang dibuka. */
  klaimID: string
  /** Dipanggil sesudah kasus berpindah — layar kembali ke Inbox. */
  onPindah: () => void
  /** Membuka panel detail (`ViewClaimDetailLifeGCNM`). */
  onDetail: () => void
}

/** Kedua tombol perpindahan, dengan tahap tujuannya. */
const PERPINDAHAN: ReadonlyArray<{ label: string; tujuan: TahapJalur }> = [
  // b21404 → `pyLocalAction SendtoAdmin` 21433.
  { label: TOMBOL_OS.kembaliKeRegister, tujuan: TAHAP_JALUR.inputRegister },
  // b21839 → `SendtoAdmin_Act1` 21863.
  { label: TOMBOL_OS.kirimKeMedis, tujuan: TAHAP_JALUR.medicalCheck },
]

export default function OutstandingClaimLife({
  klaimID,
  onPindah,
  onDetail,
}: OutstandingProps) {
  const [galat, setGalat] = useState<unknown>(null)
  const [sibuk, setSibuk] = useState(false)

  async function pindah(tujuan: TahapJalur): Promise<void> {
    setSibuk(true)
    setGalat(null)
    try {
      await pindahTahap(klaimID, tujuan)
      onPindah()
    } catch (e) {
      setGalat(e)
    } finally {
      setSibuk(false)
    }
  }

  return (
    <section className="os">
      <header className="os__kepala">
        <h2 className="os__judul">
          {OUTSTANDING.nomorKlaim}: {klaimID}
        </h2>
        <button type="button" className="os__detail" onClick={onDetail}>
          Detail klaim
        </button>
      </header>

      {/* Panel yang SAMA dengan layar Register - polis yang sama, tahap
          berbeda. */}
      <PanelDataPolis />

      <div className="os__aksi">
        {PERPINDAHAN.map((p) => (
          <button
            key={p.tujuan}
            type="button"
            className="os__tombol"
            disabled={sibuk}
            onClick={() => {
              void pindah(p.tujuan)
            }}
          >
            {p.label}
          </button>
        ))}
      </div>

      {galat !== null && (
        <>
          <Gagal galat={galat} />
          {/* Pesan server diteruskan apa adanya: 403 dan 409 masing-masing
              menyebut sebabnya, dan kalimat itu yang berguna bagi pemakai. */}
          {pesanGalat(galat) !== undefined && <p role="alert">{pesanGalat(galat)}</p>}
        </>
      )}
    </section>
  )
}
