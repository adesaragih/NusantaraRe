// Close Claim — `CloseClaim_Section.xml` b1081 → `ProtectCloseClaim_act` b1101 DAN `closeContainer` b1129.
//
// ⛔ DIPINDAH dari `pages/KlaimLife.tsx` (10-10-2026) supaya layar Outstanding (`InputOSClaimLife` b22808, modal
// CloseClaim) dan layar Detail memakai SATU logika: periksa kesiapan, konfirmasi VERBATIM b499, tutup. Gerbang yang
// sebenarnya diperiksa LAGI di server; yang di sini hanya menentukan apakah tombolnya ditawarkan.

import { useState } from 'react'

import { DETAIL } from '../labels'
import {
  bolehTutupDiLayar,
  kasusTertutup,
  penghalangDariGalat,
  periksaBolehTutup,
  tutupKlaim,
  type HasilPeriksaTutup,
  type Klaim,
} from '../api'
import { pesanGalat } from '../../../../inti/frontend/klien'

export interface PanelTutupKlaimProps {
  klaim: Klaim
  /** Sesudah kasus ditutup (padanan `closeContainer`): layar melepas klaimnya. Diberi nomor klaim yang ditutup. */
  sesudahTutup: (nomor: string) => void
  /** Sembunyikan catatan "hanya tersedia pada tahap …" (layar yang sudah menyatakan tahapnya sendiri). */
  tanpaCatatan?: boolean
}

export function PanelTutupKlaim({ klaim, sesudahTutup, tanpaCatatan = false }: PanelTutupKlaimProps) {
  // Hasil gerbang tutup. null = belum pernah diperiksa.
  const [gerbangTutup, setGerbangTutup] = useState<HasilPeriksaTutup | null>(null)
  const [tutupSibuk, setTutupSibuk] = useState(false)
  // Sudah dikonfirmasi pemakai? Konfirmasinya VERBATIM b499.
  const [mintaKonfirmasi, setMintaKonfirmasi] = useState(false)

  /** Memeriksa kesiapan tutup; TIDAK menutup apa pun. */
  async function periksaTutup(): Promise<void> {
    if (tutupSibuk) return
    setTutupSibuk(true)
    try {
      setGerbangTutup(await periksaBolehTutup(klaim.id))
    } catch (e) {
      setGerbangTutup({
        boleh: false,
        penghalang: [{ urutan: 0, nomorSertifikat: '', pesan: pesanGalat(e) ?? 'Close readiness could not be checked.' }],
      })
    } finally {
      setTutupSibuk(false)
    }
  }

  /** Menutup kasus (butir bb). ⛔ Gerbangnya dijalankan LAGI di server. */
  async function jalankanTutup(): Promise<void> {
    if (tutupSibuk) return
    const nomor = klaim.nomorKlaim !== '' ? klaim.nomorKlaim : klaim.id
    setTutupSibuk(true)
    try {
      await tutupKlaim(klaim.id)
      // Padanan `closeContainer` b1129: jendelanya DITUTUP - bukan dimuat ulang (sesudah tutup, TAHAP dikosongkan dan
      // kasusnya hilang dari kotak masuk).
      setMintaKonfirmasi(false)
      setGerbangTutup(null)
      sesudahTutup(nomor)
    } catch (e) {
      const halangan = penghalangDariGalat(e)
      setGerbangTutup({
        boleh: false,
        penghalang:
          halangan.length > 0
            ? halangan
            : [{ urutan: 0, nomorSertifikat: '', pesan: pesanGalat(e) ?? 'The claim cannot be closed.' }],
      })
      setMintaKonfirmasi(false)
    } finally {
      setTutupSibuk(false)
    }
  }

  // ⛔ Hanya pada tahap yang di Pega menawarkannya: `pyLocalAction>CloseClaim` ada di TEPAT DUA section -
  // InputOSClaimLife (b22837, b22988) dan InputAkseptasiClaimLife (b21457, b21602).
  if (!bolehTutupDiLayar(klaim)) {
    if (tanpaCatatan) return null
    return (
      <p role="note">
        {kasusTertutup(klaim)
          ? 'This case is already closed.'
          : `${DETAIL.tutupKlaim} is only available at the Outstanding Claim and Claim Analis stages.`}
      </p>
    )
  }

  return (
    <div className="tutup">
      <span className="tutup__tombol">
        <button type="button" className="btn btn--ghost" disabled={tutupSibuk} onClick={() => void periksaTutup()}>
          Check readiness: {DETAIL.tutupKlaim}
        </button>
        <button
          type="button"
          className="btn btn--ghost"
          disabled={tutupSibuk}
          onClick={() => {
            setMintaKonfirmasi(true)
          }}
        >
          {DETAIL.tutupKlaim}
        </button>
      </span>

      {/* ⛔ Konfirmasinya VERBATIM `CloseClaim_Section.xml` b499 → `pyCaption` b1499. */}
      {mintaKonfirmasi && (
        <div role="alertdialog" aria-label={DETAIL.konfirmasiTutup} className="alert alert--warn">
          <p>{DETAIL.konfirmasiTutup}</p>
          <button type="button" className="btn btn--primary btn--sm" disabled={tutupSibuk} onClick={() => void jalankanTutup()}>
            Yes, close claim
          </button>{' '}
          <button
            type="button"
            className="btn btn--ghost btn--sm"
            disabled={tutupSibuk}
            onClick={() => {
              setMintaKonfirmasi(false)
            }}
          >
            Cancel
          </button>
        </div>
      )}

      {gerbangTutup !== null && gerbangTutup.boleh && <p role="status">All insured have been accepted.</p>}
      {gerbangTutup !== null && !gerbangTutup.boleh && (
        <div role="alert" className="alert alert--warn">
          <p>The claim cannot be closed yet:</p>
          <ul>
            {gerbangTutup.penghalang.map((p) => (
              <li key={`${p.urutan}:${p.nomorSertifikat}`}>{p.pesan}</li>
            ))}
          </ul>
        </div>
      )}
    </div>
  )
}
