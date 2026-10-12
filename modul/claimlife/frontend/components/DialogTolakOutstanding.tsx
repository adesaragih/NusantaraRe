// Dialog Reject Outstanding — OQ-M5 ditutup 29-09-2026 (GILIRAN-17).
//
// Meniru `Section/RejectOSClaimLife_Sec.xml` (flow action `RejectOSClaimLife.xml`):
//
//   Date     b783 → `.DataCommitteeTreaty.TanggalComitee` b790, bawaan `@CurrentDateTime()` b893
//   PIC      b969 → `OperatorID.pyUserName` b975, baca-saja b925/b933
//   Remarks  b1680 → `.DataCommitteeTreaty.Remarks` b1687, pxTextArea b1690, WAJIB b1695
//   Submit   b3098 → `RejectOSClaimLife_Act` b3117
//
// ⚠️ Date DITAMPILKAN baca-saja, bukan kotak isian. Di Pega kotaknya dapat
// diubah, tetapi `RejectOSClaimLife_Act` tidak pernah menulisnya. Kotak yang
// dapat diubah lalu dibuang menipu pemakai, yang mengira tanggal pilihannya
// tercatat. Waktu penolakan yang tersimpan adalah jam server di baris jejak.
// PIC = akun pelaku sesi, padanan `OperatorID.pyUserName`.

import { useState } from 'react'

import { REJECT_OS, TOMBOL } from '../labels'
import { pelakuStub } from '../../../../inti/frontend/store/sesi'
import { Modal } from '../../../../inti/frontend/components/ui/dasar'

/** Kalimat konfirmasi Reject Outstanding (keputusan work owner 11-10-2026: "diberi pesan peringatan"). */
export const KONFIRMASI_TOLAK = {
  tanya: (subjek: string) => `Are you sure you want to reject this outstanding claim for ${subjek}?`,
  peringatan: "This action cannot be undone. The participant's claim will be marked as Rejected.",
  ya: 'Yes, Reject',
  proses: 'Rejecting…',
} as const

/** Lebar `T_CLAIMLF_JEJAK.KOMENTAR` VARCHAR2(4000) dalam BYTE (migrasi 021). */
export const BATAS_REMARKS_BYTE = 4000

/** Apakah `Remarks` dapat dikirim — MURNI: wajib, dan muat di kolomnya. */
export function remarksSah(teks: string): boolean {
  const t = teks.trim()
  return t !== '' && new TextEncoder().encode(t).length <= BATAS_REMARKS_BYTE
}

export function DialogTolakOutstanding({
  sibuk,
  galat,
  onKirim,
  onBatal,
  subjek,
}: {
  sibuk: boolean
  galat: string | null
  onKirim: (remarks: string) => void
  onBatal: () => void
  /** Nama peserta / Certificate No yang ditolak — disebut di pertanyaan konfirmasi. */
  subjek?: string
}) {
  const [remarks, setRemarks] = useState('')
  const [saat] = useState(() => new Date())
  const pic = pelakuStub()?.akunID ?? '—'
  const sah = remarksSah(remarks)

  return (
    <Modal
      judul={TOMBOL.tolakOutstanding}
      onTutup={onBatal}
      tanpaTutup={sibuk}
      labelBatal={REJECT_OS.batal}
      onKirim={() => {
        if (sah && !sibuk) onKirim(remarks.trim())
      }}
      aksi={
        <button type="submit" className="btn btn--danger" disabled={!sah || sibuk}>
          {sibuk ? KONFIRMASI_TOLAK.proses : KONFIRMASI_TOLAK.ya}
        </button>
      }
    >
      {/* Konfirmasi + peringatan (11-10-2026); isian Remarks Pega (RejectOSClaimLife_Sec) dipertahankan di modal yang sama. */}
      <p>{KONFIRMASI_TOLAK.tanya(subjek !== undefined && subjek.trim() !== '' ? subjek.trim() : 'this insured')}</p>
      <p role="note" className="alert alert--warn">
        {KONFIRMASI_TOLAK.peringatan}
      </p>
      <p>
        {REJECT_OS.tanggal}: <output>{saat.toLocaleString('en-GB')}</output>
      </p>
      <p>
        {REJECT_OS.pic}: <output>{pic}</output>
      </p>
      <label>
        {REJECT_OS.alasan}
        <textarea
          required
          value={remarks}
          disabled={sibuk}
          onChange={(e) => {
            setRemarks(e.target.value)
          }}
        />
      </label>
      {remarks.trim() !== '' && !sah && (
        <p role="alert">
          {REJECT_OS.alasan} exceeds {BATAS_REMARKS_BYTE} bytes.
        </p>
      )}
      {galat !== null && <p role="alert">{galat}</p>}
    </Modal>
  )
}
