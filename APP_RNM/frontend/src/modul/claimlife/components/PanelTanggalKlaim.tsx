// Tiga tanggal klaim dialog Edit Date — sensus Claim Life §3.1 (28-09-2026).
//
// Meniru `Section/EditDateClaimLife_Section.xml` selain DOL:
//
//   CLAIM RECEIVED DATE     b1069 → `.CLAIM_RECEIVED_DATE` b1076
//   DOCUMENT COMPLETE DATE  b1381 → `.COMPLETE_DATE` b1387
//   CONFIRMATION DATE       b1619 → `.CONFIRMATION_DATE` b1626
//   Save                    b1910 → `UpdateDateClaimLife_Act` b1929
//
// ⛔ SATU tombol untuk ketiganya, bukan simpan-saat-berubah seperti DOL: di
// Pega ketiganya hanya `postValue` sampai `Save` ditekan. DOL berbeda karena
// perubahannya memicu `ValidasiDOL_Act` saat itu juga (b837).
//
// ⚠️ Terbuka-tidaknya kotak (`bolehUbahTanggalKlaim`) HANYA tampilan; gerbang
// sebenarnya — Admin, tahap Outstanding, kasus terbuka — di backend.

import { useState } from 'react'

import { EDIT_DATE } from '../labels'
import { pesanGalat } from '../../../inti/klien'
import { ubahTanggalKlaim, type Peserta, type TanggalKlaim } from '../api'

/**
 * Teks tanggal backend → nilai `<input type="date">` — MURNI.
 *
 * ⛔ Backend memulangkan `YYYY-MM-DD HH24:MI:SS` (`fmtTanggalOracle`), dan
 * kotak tanggal HTML MENOLAK bentuk itu tanpa galat: kotaknya tampak kosong,
 * sehingga "sudah diisi" terbaca "belum diisi". Hanya bagian tanggalnya yang
 * dipakai; bentuk lain menjadi kosong, bukan tebakan.
 */
export function keIsianTanggal(teks: string): string {
  const cocok = /^(\d{4}-\d{2}-\d{2})/.exec(teks.trim())
  return cocok?.[1] ?? ''
}

/** Isi awal dialog dari peserta — MURNI. */
export function tanggalAwal(p: Peserta): TanggalKlaim {
  return {
    tanggalTerimaKlaim: keIsianTanggal(p.tanggalTerimaKlaim),
    tanggalDokumenLengkap: keIsianTanggal(p.tanggalDokumenLengkap),
    tanggalKonfirmasi: keIsianTanggal(p.tanggalKonfirmasi),
  }
}

/** Ketiga isian, berurut seperti di section-nya. */
const MEDAN: Array<[keyof TanggalKlaim, string]> = [
  ['tanggalTerimaKlaim', EDIT_DATE.terimaKlaim],
  ['tanggalDokumenLengkap', EDIT_DATE.dokumenLengkap],
  ['tanggalKonfirmasi', EDIT_DATE.konfirmasi],
]

export function PanelTanggalKlaim({
  klaimID,
  peserta,
  boleh,
  sesudahSimpan,
}: {
  klaimID: string
  peserta: Peserta
  boleh: boolean
  sesudahSimpan?: () => void
}) {
  const [isi, setIsi] = useState<TanggalKlaim>(() => tanggalAwal(peserta))
  const [sibuk, setSibuk] = useState(false)
  const [galat, setGalat] = useState<string | null>(null)

  async function simpan(): Promise<void> {
    if (sibuk) return
    setSibuk(true)
    setGalat(null)
    try {
      await ubahTanggalKlaim(klaimID, peserta.id, isi)
      sesudahSimpan?.()
    } catch (e) {
      setGalat(pesanGalat(e) ?? 'Tanggal klaim gagal disimpan.')
    } finally {
      setSibuk(false)
    }
  }

  return (
    <p>
      {MEDAN.map(([kunci, label]) => (
        <label key={kunci}>
          {label}{' '}
          <input
            type="date"
            value={isi[kunci]}
            disabled={!boleh || sibuk}
            onChange={(e) => {
              const nilai = e.target.value
              setIsi((lama) => ({ ...lama, [kunci]: nilai }))
            }}
          />{' '}
        </label>
      ))}
      <button type="button" disabled={!boleh || sibuk} onClick={() => void simpan()}>
        {sibuk ? 'Menyimpan…' : EDIT_DATE.simpan}
      </button>
      {galat !== null && <span role="alert"> {galat}</span>}
    </p>
  )
}
