// Dialog konfirmasi tombol bawah form (PARITAS §3.3).
//
//  `Save` b59041 → `localAction SaveProductName_Confirm` b59259: FlowAction submit `Save` b34, `Cancel` b33;
//     section b519 `Do you want to save the data?`, area teks b1025 `Comment` (`ProductName.Comment`) - dipakai
//     `AddCommentList_Act` (`SaveProductName_Act` 7 b1515).
//  `Edit` b59489 → `localAction EditProductName_Confirm` b59663: FlowAction submit `Edit` b19, `Cancel` b18;
//     section b496 `Do you want to Edit the data?`; transform `SetViewEdit` b71 → 1 b151 `IsView := "false"`.

import { Area, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { EDIT_MPNL, SIMPAN_MPNL, TOMBOL_MPNL } from '../labels'

export function DialogSimpan({
  komentar,
  onKomentar,
  menyimpan,
  onSimpan,
  onTutup,
}: {
  komentar: string
  onKomentar: (v: string) => void
  menyimpan: boolean
  onSimpan: () => void
  onTutup: () => void
}) {
  return (
    <Modal
      judul={TOMBOL_MPNL.save}
      onTutup={onTutup}
      labelBatal={SIMPAN_MPNL.cancel}
      aksi={
        <button type="button" className="btn btn--primary" disabled={menyimpan} onClick={onSimpan}>
          {SIMPAN_MPNL.save}
        </button>
      }
    >
      <p>{SIMPAN_MPNL.tanya}</p>
      <Area label={SIMPAN_MPNL.comment} value={komentar} onChange={onKomentar} />
    </Modal>
  )
}

export function DialogEdit({ onEdit, onTutup }: { onEdit: () => void; onTutup: () => void }) {
  return (
    <Modal
      judul={TOMBOL_MPNL.edit}
      onTutup={onTutup}
      labelBatal={EDIT_MPNL.cancel}
      aksi={
        <button type="button" className="btn btn--primary" onClick={onEdit}>
          {EDIT_MPNL.edit}
        </button>
      }
    >
      <p>{EDIT_MPNL.tanya}</p>
    </Modal>
  )
}
