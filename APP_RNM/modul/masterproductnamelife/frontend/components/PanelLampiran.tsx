// Panel lampiran - `InboxProductName` wadah b64133 (PARITAS §6), grid `TempData.AttachmentList`.
//
//  `Add attachment` b64747 → `ProductNameAttachContent` (submit `Attach` b24, `Cancel` b22) → `ProductNameSaveAttachment`
//  `Refresh` b65270 → `LoadAttachmentProdName`;  `Download All` b67657 (zip lampiran produk ini, RALAT R15)
//  tautan nama berkas b68903 → `DownloadAttProdName_Act`;  `View Office Online` b69291 (stub 503, OQ-MPNL-11)
//  `Delete` b69714 → `DeleteAttacProdName_act`.  `Download` b67376 (`OTHER FALSE`) mati - tidak dirender.
//
// ⛔ Lampiran melekat pada produk TERSIMPAN (tiket 08: produk dulu, lampiran menyusul) - panel ini dirender
// hanya untuk produk ber-ID. Status per lampiran (terunggah / gagal / belum) dan kirim ulang: tiket 08–09.

import { useCallback, useEffect, useState } from 'react'

import { Gagal, Kosong, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import {
  ambilLampiran,
  hapusLampiran,
  lihatOffice,
  ulangiLampiran,
  unduhLampiran,
  unduhSemuaLampiran,
  unggahLampiran,
  type Lampiran,
} from '../api'
import { tampilViewOffice } from '../bentuk'
import { LAIN_MPNL, LAMPIRAN_MPNL } from '../labels'

function teksStatus(l: Lampiran): string {
  if (l.status === 'terunggah') return LAIN_MPNL.terunggah
  if (l.status === 'gagal') return LAIN_MPNL.gagal
  return LAIN_MPNL.belum
}

export default function PanelLampiran({ produkId }: { produkId: string }) {
  const [daftar, setDaftar] = useState<Lampiran[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [galatAksi, setGalatAksi] = useState<unknown>(null)
  const [unggah, setUnggah] = useState(false)
  const [berkas, setBerkas] = useState<File | null>(null)
  const [sibuk, setSibuk] = useState(false)

  const muat = useCallback(async () => {
    try {
      setDaftar((await ambilLampiran(produkId)).daftar)
      setGalat(null)
    } catch (e) {
      setGalat(e)
    }
  }, [produkId])

  useEffect(() => {
    void muat()
  }, [muat])

  async function jalankan(aksi: () => Promise<unknown>, muatUlang = true): Promise<void> {
    if (sibuk) return
    setSibuk(true)
    setGalatAksi(null)
    try {
      await aksi()
      if (muatUlang) await muat()
    } catch (e) {
      setGalatAksi(e)
    } finally {
      setSibuk(false)
    }
  }

  return (
    <section className="mpnl-bagian">
      <div className="aksi-baris">
        <button
          type="button"
          className="btn btn--primary btn--sm"
          onClick={() => {
            setBerkas(null)
            setUnggah(true)
          }}
        >
          {LAMPIRAN_MPNL.add}
        </button>{' '}
        <button type="button" className="btn btn--ghost btn--sm" onClick={() => void muat()}>
          {LAMPIRAN_MPNL.refresh}
        </button>
      </div>
      <p className="mpnl-peringatan">{LAMPIRAN_MPNL.peringatanNama}</p>
      <p className="mpnl-peringatan">{LAMPIRAN_MPNL.peringatanGanti}</p>
      <div className="aksi-baris">
        <button type="button" className="btn btn--ghost btn--sm" onClick={() => void jalankan(() => unduhSemuaLampiran(produkId), false)}>
          {LAMPIRAN_MPNL.downloadAll}
        </button>
      </div>
      {galatAksi !== null && <Gagal galat={galatAksi} />}
      {daftar === null && galat === null && <Memuat />}
      {galat !== null && <Gagal galat={galat} />}
      {daftar !== null && daftar.length === 0 && <Kosong pesan={LAIN_MPNL.kosong} />}
      {daftar !== null && daftar.length > 0 && (
        <table className="inbox__tabel">
          <thead>
            <tr>
              <th>{LAMPIRAN_MPNL.fileName}</th>
              <th />
              <th className="table__actions" />
            </tr>
          </thead>
          <tbody>
            {daftar.map((l) => (
              <tr key={l.id} className="inbox__baris">
                <td>
                  <a
                    href="#"
                    className="mpnl-tautan"
                    onClick={(e) => {
                      e.preventDefault()
                      void jalankan(() => unduhLampiran(produkId, l), false)
                    }}
                  >
                    {l.fileName}
                  </a>{' '}
                  <span className={'mpnl-status' + (l.status === 'gagal' ? ' mpnl-status--gagal' : l.status === 'belum' ? ' mpnl-status--belum' : '')}>
                    {teksStatus(l)}
                  </span>
                  {l.galat !== undefined && l.galat !== '' && <div className="muted">{l.galat}</div>}
                </td>
                <td>
                  {tampilViewOffice(l.fileMimeType) && (
                    <a
                      href="#"
                      className="mpnl-tautan"
                      onClick={(e) => {
                        e.preventDefault()
                        void jalankan(() => lihatOffice(produkId, l.id), false)
                      }}
                    >
                      {LAMPIRAN_MPNL.viewOffice}
                    </a>
                  )}
                </td>
                <td className="table__actions">
                  {l.status !== 'terunggah' && (
                    <button type="button" className="btn btn--ghost btn--sm" onClick={() => void jalankan(() => ulangiLampiran(produkId, l.id))}>
                      {LAIN_MPNL.ulangi}
                    </button>
                  )}{' '}
                  <button type="button" className="btn btn--ghost btn--sm" onClick={() => void jalankan(() => hapusLampiran(produkId, l.id))}>
                    {LAMPIRAN_MPNL.delete}
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}

      {unggah && (
        <Modal
          judul={LAMPIRAN_MPNL.add}
          onTutup={() => {
            setUnggah(false)
          }}
          labelBatal={LAMPIRAN_MPNL.cancel}
          aksi={
            <button
              type="button"
              className="btn btn--primary"
              disabled={sibuk}
              onClick={() =>
                void jalankan(async () => {
                  await unggahLampiran(produkId, berkas)
                  setUnggah(false)
                })
              }
            >
              {LAMPIRAN_MPNL.attach}
            </button>
          }
        >
          {galatAksi !== null && <Gagal galat={galatAksi} />}
          <input
            type="file"
            aria-label={LAMPIRAN_MPNL.fileName}
            onChange={(e) => {
              setBerkas(e.target.files?.[0] ?? null)
            }}
          />
        </Modal>
      )}
    </section>
  )
}
