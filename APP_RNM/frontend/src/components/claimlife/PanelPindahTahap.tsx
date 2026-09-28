// Tombol perpindahan tahap per LAYAR — A3 kelompok Medis dan Akseptasi.
//
// ⛔ Tiap tahap menawarkan persis perpindahan yang section-nya sendiri punya.
// Itu bukan tata letak melainkan cacah berkas:
//
//   Outstanding Claim  `InputOSClaimLife.xml`
//        b21404 `Send Back to Register` -> Input Register
//        b21839 `Send to Medical Check` -> Medical Check
//   Medical Check      `MedicalCheckClaimLife.xml`
//        b20256 `Send Back to Admin`    -> Outstanding
//        b21151 `Send to Claim Analyst` -> Claim Analis
//   Claim Analis       `InputAkseptasiClaimLife.xml`
//        b20221 `Send Back to Admin`    -> Outstanding
//        b20467 `Send Back to Medical`  -> Medical Check
//
// Menyalin tombol antarlayar membuka jalur yang di sistem lama tidak ada —
// dan `models.SerahTerimaSah` di backend akan menolaknya dengan 409, sehingga
// yang lahir hanya tombol yang selalu gagal.
//
// ⛔ Layar Input Register TIDAK ada di sini: satu-satunya perpindahannya
// (`InputRegisterClaimLife` -> Outstanding) terjadi lewat `Submit`
// pendaftaran, bukan lewat tombol perpindahan.

import { useState } from 'react'

import { TOMBOL_AKSEPTASI, TOMBOL_MEDIS, TOMBOL_OS } from '../../assets/labels.claimlife'
import {
  pesanGalat,
  pindahTahap,
  TAHAP_JALUR,
  type TahapJalur,
} from '../../services/api'

/** Satu tombol perpindahan: label VERBATIM dan tahap tujuannya. */
export interface TombolPindah {
  label: string
  tujuan: TahapJalur
}

/**
 * Perpindahan yang DITAWARKAN sebuah tahap.
 *
 * ⚠️ Dipisah dari komponennya supaya dapat diuji tanpa DOM, dan supaya
 * daftarnya dapat dikunci dari kedua arah: tombol yang HILANG dari layar
 * sama merusaknya dengan tombol yang ditambahkan, dan yang pertama tidak
 * berbunyi sendiri.
 *
 * Tahap yang tidak dikenal mendapat daftar KOSONG — gagal tertutup.
 */
export function tombolPindahTahap(tahap: string): TombolPindah[] {
  switch (tahap) {
    case 'Outstanding Claim':
      return [
        { label: TOMBOL_OS.kembaliKeRegister, tujuan: TAHAP_JALUR.inputRegister },
        { label: TOMBOL_OS.kirimKeMedis, tujuan: TAHAP_JALUR.medicalCheck },
      ]
    case 'Medical Check':
      return [
        { label: TOMBOL_MEDIS.kembaliKeAdmin, tujuan: TAHAP_JALUR.outstanding },
        { label: TOMBOL_MEDIS.kirimKeAnalis, tujuan: TAHAP_JALUR.claimAnalis },
      ]
    case 'Claim Analis':
      return [
        { label: TOMBOL_AKSEPTASI.kembaliKeAdmin, tujuan: TAHAP_JALUR.outstanding },
        { label: TOMBOL_AKSEPTASI.kembaliKeMedis, tujuan: TAHAP_JALUR.medicalCheck },
      ]
    default:
      return []
  }
}

export function PanelPindahTahap({
  klaimID,
  tahap,
  sesudahPindah,
}: {
  klaimID: string
  tahap: string
  sesudahPindah?: () => void
}) {
  const [sibuk, setSibuk] = useState<string | null>(null)
  const [galat, setGalat] = useState<string | null>(null)

  const tombol = tombolPindahTahap(tahap)
  if (tombol.length === 0) return null

  async function pindah(tujuan: TahapJalur): Promise<void> {
    if (sibuk !== null) return
    setSibuk(tujuan)
    setGalat(null)
    try {
      await pindahTahap(klaimID, tujuan)
      sesudahPindah?.()
    } catch (e) {
      setGalat(pesanGalat(e) ?? 'Perpindahan tahap gagal.')
    } finally {
      setSibuk(null)
    }
  }

  return (
    <section className="tahap">
      <h4 className="tahap__judul">Perpindahan dari {tahap}</h4>
      <p>
        {tombol.map((t) => (
          <button
            key={t.tujuan}
            type="button"
            disabled={sibuk !== null}
            onClick={() => void pindah(t.tujuan)}
          >
            {sibuk === t.tujuan ? 'Memindahkan…' : t.label}
          </button>
        ))}
      </p>
      {galat !== null && <p role="alert">{galat}</p>}
    </section>
  )
}
