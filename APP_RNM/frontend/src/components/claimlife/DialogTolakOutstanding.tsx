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

import { REJECT_OS, TOMBOL } from '../../assets/labels.claimlife'
import { pelakuStub } from '../../store/sesi'
import { Modal } from '../ui/dasar'

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
}: {
  sibuk: boolean
  galat: string | null
  onKirim: (remarks: string) => void
  onBatal: () => void
}) {
  const [remarks, setRemarks] = useState('')
  const [saat] = useState(() => new Date())
  const pic = pelakuStub()?.akunID ?? '—'
  const sah = remarksSah(remarks)

  return (
    <Modal
      judul={TOMBOL.tolakOutstanding}
      onTutup={onBatal}
      labelBatal={REJECT_OS.batal}
      onKirim={() => {
        if (sah && !sibuk) onKirim(remarks.trim())
      }}
      aksi={
        <button type="submit" className="btn btn--primary" disabled={!sah || sibuk}>
          {sibuk ? 'Menolak…' : REJECT_OS.kirim}
        </button>
      }
    >
      <p>
        {REJECT_OS.tanggal}: <output>{saat.toLocaleString('id-ID')}</output>
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
      {galat !== null && <p role="alert">{galat}</p>}
    </Modal>
  )
}
