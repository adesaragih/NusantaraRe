// Panel lampiran - `InboxProductName` wadah b64133 (PARITAS §6), grid `TempData.AttachmentList`.
//
//  `Add attachment` b64747 → `ProductNameAttachContent` (submit `Attach` b24, `Cancel` b22) → `ProductNameSaveAttachment`
//  - banyak berkas sekaligus dan seret-lepas (permintaan work owner 03-10-2026): diunggah satu per satu lewat rute yang ada
//  `Refresh` b65270 → `LoadAttachmentProdName`;  `Download All` b67657 (zip lampiran produk ini, RALAT R15)
//  tautan nama berkas b68903 → `DownloadAttProdName_Act`;  `View Office Online` b69291 → URL bertanda tangan dibuka di
//  penampil kantor (`DownloadAttProdName_Act` 7 b1103; keputusan work owner 03-10-2026): form GET tersembunyi
//  ber-`action` tetap ke penampil (`penampilOffice.ts`) dikirim ke bingkai DI POPUP;  `View` pdf / gambar (permintaan
//  work owner 03-10-2026 "dari popup atau windows baru (bukan tab baru)"): isi dari rute unduh yang ada, objek URL
//  lokal di popup. Jendela peramban terpisah menuntut pembukaan jendela lewat skrip - dilarang penjaga lintas-modul
//  `unduhdokumen.test.ts`.
//  `Delete` b69714 → `DeleteAttacProdName_act`.  `Download` b67376 (`OTHER FALSE`) mati - tidak dirender.
//
// ⛔ Lampiran melekat pada produk TERSIMPAN (tiket 08: produk dulu, lampiran menyusul) - panel ini dirender
// hanya untuk produk ber-ID. Status per lampiran (terunggah / gagal / belum) dan kirim ulang: tiket 08–09.
// Mode lihat (keputusan work owner 03-10-2026 "jika view tidak tambah/edit/delete"): `Add attachment`, kirim ulang, dan
// `Delete` tersembunyi; `Refresh`, `Download All`, unduh berkas, `View Office Online`, dan `View` tetap.

import { useCallback, useEffect, useRef, useState } from 'react'

import { Gagal, Kosong, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { pesanGalat } from '../../../../inti/frontend/klien'
import {
  ambilIsiLampiran,
  ambilLampiran,
  hapusLampiran,
  lihatOffice,
  ulangiLampiran,
  unduhLampiran,
  unduhSemuaLampiran,
  unggahLampiran,
  type Lampiran,
} from '../api'
import { gabungBerkas, jenisViewOnline, mimeViewOnline, tampilViewOffice, unggahBerurutan } from '../bentuk'
import { GRID_MPNL, LAIN_MPNL, LAMPIRAN_MPNL } from '../labels'
import { BINGKAI_PENAMPIL, PARAM_PENAMPIL, PENAMPIL_OFFICE } from '../penampilOffice'

function teksStatus(l: Lampiran): string {
  if (l.status === 'terunggah') return LAIN_MPNL.terunggah
  if (l.status === 'gagal') return LAIN_MPNL.gagal
  return LAIN_MPNL.belum
}

export default function PanelLampiran({ produkId, lihat }: { produkId: string; lihat: boolean }) {
  const [daftar, setDaftar] = useState<Lampiran[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [galatAksi, setGalatAksi] = useState<unknown>(null)
  const [unggah, setUnggah] = useState(false)
  // `Add attachment`: berkas terpilih, berkas yang sedang diunggah ("2/5 nama"), kegagalan per berkas, seret di atas kotak.
  const [berkas, setBerkas] = useState<File[]>([])
  const [proses, setProses] = useState<string | null>(null)
  const [gagalUnggah, setGagalUnggah] = useState<{ berkas: File; galat: unknown }[]>([])
  const [seret, setSeret] = useState(false)
  const [sibuk, setSibuk] = useState(false)
  // `View Office Online`: form GET tersembunyi ke penampil dan input `src`-nya.
  const formOffice = useRef<HTMLFormElement>(null)
  const urlOffice = useRef<HTMLInputElement>(null)
  // Popup penampil: berkas, jenis isi, dan apakah isinya sudah terpasang.
  const [penampil, setPenampil] = useState<{ nama: string; jenis: 'office' | 'pdf' | 'gambar'; siap: boolean } | null>(null)
  const bingkai = useRef<HTMLIFrameElement>(null)
  const gambar = useRef<HTMLImageElement>(null)
  // Objek URL yang sedang tampil (dicabut saat popup ditutup) dan giliran buka (jawaban popup yang sudah ditutup dibuang).
  const objekAktif = useRef<string | null>(null)
  const giliran = useRef(0)

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

  function tutupPenampil(): void {
    giliran.current++
    if (objekAktif.current !== null) URL.revokeObjectURL(objekAktif.current)
    objekAktif.current = null
    setPenampil(null)
  }

  /** `View Office Online`: URL bertanda tangan → form GET ke bingkai popup. */
  function bukaOffice(l: Lampiran): void {
    if (sibuk) return
    const ke = ++giliran.current
    setPenampil({ nama: l.fileName, jenis: 'office', siap: false })
    void jalankan(async () => {
      try {
        const url = await lihatOffice(produkId, l.id)
        if (ke !== giliran.current) return
        if (formOffice.current !== null && urlOffice.current !== null) {
          urlOffice.current.value = url
          formOffice.current.submit()
        }
        setPenampil((p) => (p === null ? p : { ...p, siap: true }))
      } catch (err) {
        if (ke === giliran.current) tutupPenampil()
        throw err
      }
    }, false)
  }

  /** `View` pdf / gambar: isi dari rute unduh yang ada → objek URL lokal di popup. */
  function bukaOnline(l: Lampiran): void {
    const jenis = jenisViewOnline(l.fileMimeType)
    if (sibuk || jenis === null) return
    const ke = ++giliran.current
    setPenampil({ nama: l.fileName, jenis, siap: false })
    void jalankan(async () => {
      try {
        const isi = await ambilIsiLampiran(produkId, l.id)
        if (ke !== giliran.current) return
        // Tipe dari ekstensi, bukan dari jawaban: hanya pdf / gambar raster yang dirender.
        const objekURL = URL.createObjectURL(new Blob([isi], { type: mimeViewOnline(l.fileMimeType) }))
        objekAktif.current = objekURL
        if (jenis === 'gambar' && gambar.current !== null) gambar.current.src = objekURL
        if (jenis === 'pdf' && bingkai.current !== null) bingkai.current.src = objekURL
        setPenampil((p) => (p === null ? p : { ...p, siap: true }))
      } catch (err) {
        if (ke === giliran.current) tutupPenampil()
        throw err
      }
    }, false)
  }

  return (
    <section className="mpnl-bagian">
      <div className="aksi-baris">
        {!lihat && (
          <button
            type="button"
            className="btn btn--primary btn--sm"
            onClick={() => {
              setBerkas([])
              setGagalUnggah([])
              setUnggah(true)
            }}
          >
            {LAMPIRAN_MPNL.add}
          </button>
        )}{' '}
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
        <div className="mpnl-tabel">
          <table className="inbox__tabel">
            <thead>
              <tr>
                <th>{LAMPIRAN_MPNL.fileName}</th>
                <th />
                {!lihat && <th className="table__actions" />}
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
                          bukaOffice(l)
                        }}
                      >
                        {LAMPIRAN_MPNL.viewOffice}
                      </a>
                    )}
                    {jenisViewOnline(l.fileMimeType) !== null && (
                      <a
                        href="#"
                        className="mpnl-tautan"
                        onClick={(e) => {
                          e.preventDefault()
                          bukaOnline(l)
                        }}
                      >
                        {GRID_MPNL.view}
                      </a>
                    )}
                  </td>
                  {!lihat && (
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
                  )}
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {/* `View Office Online` b1103: `<penampil>?src=<URL bertanda tangan>` di bingkai popup. */}
      <form ref={formOffice} method="get" action={PENAMPIL_OFFICE} target={BINGKAI_PENAMPIL} hidden>
        <input ref={urlOffice} type="hidden" name={PARAM_PENAMPIL} />
      </form>

      {penampil !== null && (
        <Modal judul={penampil.nama} onTutup={tutupPenampil} labelBatal={LAMPIRAN_MPNL.cancel} penuh>
          {!penampil.siap && <Memuat />}
          {penampil.jenis === 'gambar' ? (
            <img ref={gambar} alt={penampil.nama} className="mpnl-penampil__gambar" />
          ) : (
            <iframe ref={bingkai} name={BINGKAI_PENAMPIL} title={penampil.nama} className="mpnl-penampil__bingkai" />
          )}
        </Modal>
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
                  if (berkas.length === 0) {
                    // Kalimat backend VERBATIM: "Tidak ada file yg diattach" (`ProductNameSaveAttachment` 1 b292).
                    await unggahLampiran(produkId, null)
                    return
                  }
                  setGagalUnggah([])
                  const gagal = await unggahBerurutan(
                    berkas,
                    (f) => unggahLampiran(produkId, f),
                    (f, i) => {
                      setProses(`${i + 1}/${berkas.length} ${f.name}`)
                    },
                  )
                  setProses(null)
                  // Yang berhasil keluar dari pilihan (tampil di grid); yang gagal tinggal untuk diperbaiki/diulang.
                  setGagalUnggah(gagal)
                  setBerkas(gagal.map((g) => g.berkas))
                  if (gagal.length === 0) setUnggah(false)
                })
              }
            >
              {LAMPIRAN_MPNL.attach}
            </button>
          }
        >
          {galatAksi !== null && <Gagal galat={galatAksi} />}
          <label
            className={'mpnl-unggah' + (seret ? ' mpnl-unggah--seret' : '')}
            onDragOver={(e) => {
              e.preventDefault()
              setSeret(true)
            }}
            onDragLeave={() => {
              setSeret(false)
            }}
            onDrop={(e) => {
              e.preventDefault()
              setSeret(false)
              const jatuh = Array.from(e.dataTransfer.files)
              setBerkas((b) => gabungBerkas(b, jatuh))
            }}
          >
            <span>{LAIN_MPNL.seretBerkas}</span>
            <input
              type="file"
              multiple
              aria-label={LAMPIRAN_MPNL.fileName}
              disabled={sibuk}
              onChange={(e) => {
                const dipilih = Array.from(e.target.files ?? [])
                setBerkas((b) => gabungBerkas(b, dipilih))
                // Dikosongkan: memilih berkas yang sama lagi tetap memicu onChange.
                e.target.value = ''
              }}
            />
          </label>
          {berkas.length > 0 && (
            <ul className="mpnl-unggah__daftar">
              {berkas.map((f) => (
                <li key={f.name}>
                  <span>{f.name}</span>
                  <button
                    type="button"
                    className="btn btn--ghost btn--sm"
                    disabled={sibuk}
                    onClick={() => {
                      setBerkas((b) => b.filter((x) => x !== f))
                    }}
                  >
                    {LAIN_MPNL.buangPilihan}
                  </button>
                </li>
              ))}
            </ul>
          )}
          {proses !== null && (
            <p className="muted">
              {LAIN_MPNL.mengunggah} {proses}
            </p>
          )}
          {gagalUnggah.length > 0 && (
            <ul className="mpnl-unggah__gagal">
              {gagalUnggah.map((g) => (
                <li key={g.berkas.name}>
                  {g.berkas.name}: {pesanGalat(g.galat) ?? (g.galat instanceof Error ? g.galat.message : LAIN_MPNL.gagal)}
                </li>
              ))}
            </ul>
          )}
        </Modal>
      )}
    </section>
  )
}
