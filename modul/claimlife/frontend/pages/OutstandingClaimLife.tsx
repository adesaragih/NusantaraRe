// Layar Outstanding Claim — A3 kelompok Outstanding.
//
// Meniru `FlowAction/OSClaimLife.xml` → `Section/InputOSClaimLife.xml`.
//
// ⚠️ Kesebelas medan polisnya SAMA dengan layar Register dan terikat
// `.PolicyDataLife.*` pula - keduanya menampilkan polis yang sama pada tahap
// yang berbeda. Karena itu keduanya memakai `PanelDataPolis` yang sama, yang
// sejak butir av membacanya dari modul PremiumList Life. Menyalin panelnya
// berarti dua bentuk yang harus berubah bersama.
//
// ⛔ BUTIR aw - dua tombol perpindahan. Di Pega keduanya TIDAK menulis apa
// pun pada posisi Admin: `SendtoAdmin_Act` dan `SendtoAdmin_Act1` sama-sama
// berprasyarat `pyPosition=="ReasLifeMedicalAdvisor"` (WhenFalse=3 = lewati),
// padahal layar ini dipegang `ReasLifeAdmin`. Itu cacat rule warisan, dan
// yang ditiru MAKSUDnya - lihat `OQ-untuk-tim.md` OQ-C.

import { useEffect, useState } from 'react'

import { OUTSTANDING, TAHAP, TOMBOL_OS } from '../labels'
import { PanelDataPolis } from '../components/PanelDataPolis'
import { PanelPindahTahap } from '../components/PanelPindahTahap'
import { Gagal } from '../../../../inti/frontend/components/ui/dasar'
import {
  ambilDataPolis,
  ambilKlaimLife,
  simpanKeRNM,
  type HasilSimpanRNM,
  type PolicyDataLife,
} from '../api'
import { pesanGalat } from '../../../../inti/frontend/klien'

export interface OutstandingProps {
  /** Pengenal work kasus yang sedang dibuka. */
  klaimID: string
  /** Dipanggil sesudah kasus berpindah — layar kembali ke Inbox. */
  onPindah: () => void
  /** Membuka panel detail (`ViewClaimDetailLifeGCNM`). */
  onDetail: () => void
}

export default function OutstandingClaimLife({
  klaimID,
  onPindah,
  onDetail,
}: OutstandingProps) {
  const [galat, setGalat] = useState<unknown>(null)
  const [sibuk, setSibuk] = useState(false)
  const [tersimpan, setTersimpan] = useState<HasilSimpanRNM | null>(null)
  // av-2: panelnya kini BERISI - sebelumnya `<PanelDataPolis />` tanpa polis,
  // sehingga seluruh medan tampil kosong di layar ini.
  const [polis, setPolis] = useState<PolicyDataLife | null>(null)
  // ⛔ Tahap kasus yang SEBENARNYA (GILIRAN-11). Kotak masuk membuka layar ini
  // untuk kasus tahap mana pun; tanpa ini layar menawarkan tombol Outstanding
  // kepada kasus Medical Check - dan setiap kliknya ditolak 409.
  const [tahapKasus, setTahapKasus] = useState<string | null>(null)
  // ⛔ Klaim yang GAGAL dimuat dinyatakan, bukan ditelan: tanpa ini layar
  // tampil kosong tanpa satu tombol pun dan tanpa sebab (temuan /code-review).
  const [galatMuat, setGalatMuat] = useState<unknown>(null)

  useEffect(() => {
    let batal = false
    void (async () => {
      setGalatMuat(null)
      let nomorPolis: string
      try {
        const k = await ambilKlaimLife(klaimID)
        if (batal) return
        setTahapKasus(k.tahap)
        nomorPolis = k.nomorPolis
      } catch (e) {
        if (!batal) setGalatMuat(e)
        return
      }
      try {
        const p = await ambilDataPolis(nomorPolis)
        if (!batal) setPolis(p)
      } catch {
        // Polis yang belum ada di PremiumList Life BUKAN kerusakan layar ini:
        // panelnya menyatakan "belum terbaca", Save to RNM menolaknya terang.
        if (!batal) setPolis(null)
      }
    })()
    return () => {
      batal = true
    }
  }, [klaimID])

  /** `Save to RNM` b21102 — gerbangnya seluruhnya di server. */
  async function simpan(): Promise<void> {
    setSibuk(true)
    setGalat(null)
    try {
      setTersimpan(await simpanKeRNM(klaimID))
    } catch (e) {
      setTersimpan(null)
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
      <PanelDataPolis polis={polis} />

      {galatMuat !== null && (
        <>
          <Gagal galat={galatMuat} />
          {pesanGalat(galatMuat) !== undefined && <p role="alert">{pesanGalat(galatMuat)}</p>}
        </>
      )}

      {tahapKasus !== null && tahapKasus !== TAHAP.outstandingClaim && (
        <p role="note">
          Kasus ini berada di tahap {tahapKasus === '' ? 'yang tidak diketahui' : tahapKasus}, bukan{' '}
          {TAHAP.outstandingClaim}: tombol layar ini tidak berlaku untuknya. Buka
          dari layar Detail klaim.
        </p>
      )}

      {tahapKasus === TAHAP.outstandingClaim && (
      <div className="os__aksi">
        <button
          type="button"
          className="os__tombol"
          disabled={sibuk}
          onClick={() => {
            void simpan()
          }}
        >
          {sibuk ? 'Menyimpan…' : TOMBOL_OS.simpanRNM}
        </button>
      </div>
      )}
      {tersimpan !== null && (
        <p role="status">
          Tersimpan ke RNM — nomor klaim {tersimpan.nomorKlaim}
          {tersimpan.nomorBaru ? ' (baru diterbitkan)' : ''}; {tersimpan.barisDitandai} baris
          ditandai Outstanding; Arasapas: {tersimpan.arasapas}.
        </p>
      )}

      {/* ⛔ Perpindahan TIDAK disalin di sini lagi (GILIRAN-11): salinan lokal
          melewatkan konfirmasi `Send Back to Admin?` (SendtoAdmin_Section
          b566) yang dibuka local action `SendtoAdmin` b21433. Satu daftar,
          satu dialog — `PanelPindahTahap`. */}
      {tahapKasus === TAHAP.outstandingClaim && (
        <PanelPindahTahap klaimID={klaimID} tahap={TAHAP.outstandingClaim} sesudahPindah={onPindah} />
      )}

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
