// Lampiran berkas - section `AttachmentsBdx` (`InputBordereaux` b27972: sesudah Close/Save, sebelum History), di atas
// penyimpanan bersama `inti/backend/penyimpanan` (keputusan work owner 08-10-2026, "kaya XML nya").
//
//  grid `KategoriDocument.pxResults` (deferred `GetKategotyDocBdx`): Category b3039, Count b3192, Upload File b3342,
//  View File b3490; `Refresh` b1800. `Download All` b1566 tidak dirender (`pyCondition never` b1656, tanpa aksi).
//  `Upload File` b4112 = modal `BordereauxAttach` (Submit b21 / Cancel b20): banyak berkas, diunggah satu per satu -
//  `AttachDocBdx_Post` mengulang `dragDropFileUpload.pxResults`.
//  `View File` b4618 = popup `AttachmentDetailBdx`: File Name b1529 (unduh), View Office Online b2562 (xls…pptx, b2876),
//  Type b1796 (kategori), Delete b3185 (tanpa konfirmasi, seperti Pega); `View` pdf / gambar (permintaan work owner
//  08-10-2026, seperti Product Name Life). View dan View Office Online membuka popup penampil layar penuh.
//  Upload File dan Delete hanya bila `boleh` = berkas dibuka Edit oleh yang berhak (Pega `ViewStage != 1 || IT
//  Developer`, b4459 / b3376; pengecualian IT Developer dibuang - keputusan work owner 08-10-2026).
//
// ⛔ Unduh dan View lewat fetch beridentitas (penjaga `unduhdokumen.test.ts`); View Office Online = form GET
// tersembunyi ke penampil di bingkai popup penampil (`penampilOffice.ts`) - Pega membuka jendela popup
// (`openUrlInWindow` b2217).
// ⛔ Lampiran melekat pada berkas TERSIMPAN (BDX_ID lahir saat Save) - panel ini hanya dirender untuk berkas ber-ID.

import { useCallback, useEffect, useRef, useState } from 'react'

import { Gagal, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { pesanGalat } from '../../../../inti/frontend/klien'
import {
  ambilIsiLampiran,
  ambilKategoriLampiran,
  ambilLampiran,
  hapusLampiran,
  tautanOffice,
  unduhLampiran,
  unggahLampiran,
  type KategoriLampiran,
  type Lampiran,
} from '../api'
import { gabungBerkas } from '../aturan'
import { LAMPIRAN_BDX } from '../labels'
import {
  BINGKAI_PENAMPIL,
  PARAM_PENAMPIL,
  PENAMPIL_OFFICE,
  bisaViewOffice,
  jenisViewOnline,
  mimeViewOnline,
} from '../penampilOffice'

export default function PanelLampiran({ bdxId, boleh }: { bdxId: string; boleh: boolean }) {
  const [kategori, setKategori] = useState<KategoriLampiran[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [unggahKe, setUnggahKe] = useState<KategoriLampiran | null>(null)
  const [lihatKe, setLihatKe] = useState<KategoriLampiran | null>(null)

  const muat = useCallback(() => {
    setGalat(null)
    ambilKategoriLampiran(bdxId).then(
      (r) => setKategori(r.daftar),
      (g: unknown) => setGalat(g),
    )
  }, [bdxId])

  useEffect(() => {
    muat()
  }, [muat])

  return (
    <section className="bordereaux__kartu">
      <div className="bordereaux__lampiran-kepala">
        <h3 className="bordereaux__judul-kartu">{LAMPIRAN_BDX.judul}</h3>
        <span className="toolbar__spacer" />
        <button type="button" className="btn btn--ghost btn--sm" onClick={muat}>
          {LAMPIRAN_BDX.refresh}
        </button>
      </div>
      {galat !== null && <Gagal galat={galat} />}
      {kategori === null && galat === null && <Memuat pesan={LAMPIRAN_BDX.memuat} />}
      {kategori !== null &&
        (kategori.length === 0 ? (
          <p className="muted">{LAMPIRAN_BDX.kosong}</p>
        ) : (
          <table className="inbox__tabel bordereaux__grid">
            <thead>
              <tr>
                <th>{LAMPIRAN_BDX.category}</th>
                <th className="bordereaux__angka">{LAMPIRAN_BDX.count}</th>
                <th>{LAMPIRAN_BDX.uploadFile}</th>
                <th>{LAMPIRAN_BDX.viewFile}</th>
              </tr>
            </thead>
            <tbody>
              {kategori.map((k) => (
                <tr key={k.id} className="inbox__baris">
                  <td>{k.nama}</td>
                  <td className="bordereaux__angka">{k.cacah}</td>
                  <td>
                    {boleh && (
                      <button type="button" className="btn btn--ghost btn--sm" onClick={() => setUnggahKe(k)}>
                        {LAMPIRAN_BDX.uploadFile}
                      </button>
                    )}
                  </td>
                  <td>
                    <button type="button" className="btn btn--ghost btn--sm" onClick={() => setLihatKe(k)}>
                      {LAMPIRAN_BDX.viewFile}
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        ))}

      {unggahKe !== null && (
        <DialogUnggah
          bdxId={bdxId}
          kategori={unggahKe}
          onTutup={(berubah) => {
            setUnggahKe(null)
            if (berubah) muat()
          }}
        />
      )}
      {lihatKe !== null && (
        <DialogLihat
          bdxId={bdxId}
          kategori={lihatKe}
          boleh={boleh}
          onTutup={(berubah) => {
            setLihatKe(null)
            if (berubah) muat()
          }}
        />
      )}
    </section>
  )
}

/**
 * `BordereauxAttach` - `pyAttachmentScreen` (`dragDropFileUpload`): seret-lepas atau klik untuk memilih, berkas digabung
 * ke daftar; Submit mengunggah satu per satu, yang gagal tetap terpilih.
 */
function DialogUnggah({
  bdxId,
  kategori,
  onTutup,
}: {
  bdxId: string
  kategori: KategoriLampiran
  onTutup: (berubah: boolean) => void
}) {
  const [berkas, setBerkas] = useState<File[]>([])
  const [proses, setProses] = useState<string | null>(null)
  const [gagal, setGagal] = useState<{ nama: string; galat: unknown }[]>([])
  const [pesan, setPesan] = useState<string | null>(null)
  const [masuk, setMasuk] = useState(false)
  const [seret, setSeret] = useState(false)

  const tambah = (baru: readonly File[]) => {
    if (proses !== null || baru.length === 0) return
    setBerkas((b) => gabungBerkas(b, baru))
    setGagal([])
    setPesan(null)
  }

  const kirim = async () => {
    if (proses !== null) return
    if (berkas.length === 0) {
      setPesan(LAMPIRAN_BDX.tanpaBerkas)
      return
    }
    setPesan(null)
    const tertolak: { berkas: File; galat: unknown }[] = []
    let ada = masuk
    for (const [i, f] of berkas.entries()) {
      setProses(LAMPIRAN_BDX.mengunggah(i + 1, berkas.length, f.name))
      try {
        await unggahLampiran(bdxId, kategori.id, f)
        ada = true
      } catch (g: unknown) {
        tertolak.push({ berkas: f, galat: g })
      }
    }
    setProses(null)
    setMasuk(ada)
    if (tertolak.length === 0) {
      onTutup(true)
      return
    }
    setBerkas(tertolak.map((t) => t.berkas))
    setGagal(tertolak.map((t) => ({ nama: t.berkas.name, galat: t.galat })))
  }

  return (
    <Modal
      judul={LAMPIRAN_BDX.judulUnggah(kategori.nama)}
      onTutup={() => onTutup(masuk)}
      onKirim={() => void kirim()}
      labelBatal={LAMPIRAN_BDX.cancel}
      aksi={
        <button type="submit" className="btn btn--primary" disabled={proses !== null}>
          {LAMPIRAN_BDX.submit}
        </button>
      }
    >
      <label
        className={'bordereaux__lampiran-seret' + (seret ? ' bordereaux__lampiran-seret--aktif' : '')}
        onDragOver={(e) => {
          e.preventDefault()
          e.dataTransfer.dropEffect = proses === null ? 'copy' : 'none'
          setSeret(true)
        }}
        onDragLeave={(e) => {
          // Pindah ke anak kotak (teks, input) bukan meninggalkan kotak.
          if (!e.currentTarget.contains(e.relatedTarget as Node | null)) setSeret(false)
        }}
        onDrop={(e) => {
          e.preventDefault()
          setSeret(false)
          tambah(Array.from(e.dataTransfer.files))
        }}
      >
        <span>{LAMPIRAN_BDX.seretBerkas}</span>
        <input
          type="file"
          multiple
          aria-label={LAMPIRAN_BDX.pilihBerkas}
          disabled={proses !== null}
          onChange={(e) => {
            const dipilih = Array.from(e.target.files ?? [])
            // Dikosongkan: memilih berkas yang sama lagi tetap memicu onChange.
            e.target.value = ''
            tambah(dipilih)
          }}
        />
      </label>
      {berkas.length > 0 && (
        <ul className="bordereaux__lampiran-pilihan">
          <li className="muted">{LAMPIRAN_BDX.berkasTerpilih(berkas.length)}</li>
          {berkas.map((f) => (
            <li key={f.name + f.size}>
              <span>{f.name}</span>
              <button
                type="button"
                className="btn btn--ghost btn--sm"
                disabled={proses !== null}
                onClick={() => setBerkas((b) => b.filter((x) => x !== f))}
              >
                {LAMPIRAN_BDX.buangPilihan}
              </button>
            </li>
          ))}
        </ul>
      )}
      {proses !== null && <p className="muted">{proses}</p>}
      {pesan !== null && <p className="bordereaux__lampiran-galat">{pesan}</p>}
      {gagal.length > 0 && (
        <div className="bordereaux__lampiran-galat">
          <p>{LAMPIRAN_BDX.gagalUnggah}</p>
          <ul>
            {gagal.map((g) => (
              <li key={g.nama}>
                {g.nama}: {pesanGalat(g.galat) ?? String(g.galat)}
              </li>
            ))}
          </ul>
        </div>
      )}
    </Modal>
  )
}

/**
 * `AttachmentDetailBdx` - lampiran satu kategori. File Name mengunduh; `View` (pdf / gambar) dan View Office Online
 * membuka POPUP PENAMPIL layar penuh (permintaan work owner 08-10-2026, seperti Product Name Life). Selama penampil
 * terbuka, popup daftar diganti olehnya lalu kembali saat ditutup - modal bertumpuk tertutup bersama oleh Escape.
 */
function DialogLihat({
  bdxId,
  kategori,
  boleh,
  onTutup,
}: {
  bdxId: string
  kategori: KategoriLampiran
  boleh: boolean
  onTutup: (berubah: boolean) => void
}) {
  const [daftar, setDaftar] = useState<Lampiran[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [sibuk, setSibuk] = useState(false)
  const [berubah, setBerubah] = useState(false)
  // Popup penampil: berkas, jenis isi, dan apakah isinya sudah terpasang.
  const [penampil, setPenampil] = useState<{ nama: string; jenis: 'office' | 'pdf' | 'gambar'; siap: boolean } | null>(
    null,
  )
  const formOffice = useRef<HTMLFormElement>(null)
  const urlOffice = useRef<HTMLInputElement>(null)
  const bingkai = useRef<HTMLIFrameElement>(null)
  const gambar = useRef<HTMLImageElement>(null)
  // Objek URL yang sedang tampil (dicabut saat penampil ditutup) dan giliran buka (jawaban penampil tertutup dibuang).
  const objekAktif = useRef<string | null>(null)
  const giliran = useRef(0)

  const muat = useCallback(() => {
    ambilLampiran(bdxId, kategori.id).then(
      (r) => setDaftar(r.daftar),
      (g: unknown) => setGalat(g),
    )
  }, [bdxId, kategori.id])

  useEffect(() => {
    muat()
  }, [muat])

  useEffect(
    () => () => {
      if (objekAktif.current !== null) URL.revokeObjectURL(objekAktif.current)
    },
    [],
  )

  const jalankan = (kerja: () => Promise<void>) => {
    if (sibuk) return
    setSibuk(true)
    setGalat(null)
    kerja()
      .catch((g: unknown) => setGalat(g))
      .finally(() => setSibuk(false))
  }

  const tutupPenampil = () => {
    giliran.current++
    if (objekAktif.current !== null) URL.revokeObjectURL(objekAktif.current)
    objekAktif.current = null
    setPenampil(null)
  }

  /** View Office Online: URL bertanda tangan -> form GET ke bingkai popup penampil. */
  const bukaOffice = (l: Lampiran) => {
    if (sibuk) return
    const ke = ++giliran.current
    setPenampil({ nama: l.fileName, jenis: 'office', siap: false })
    jalankan(async () => {
      try {
        const { url } = await tautanOffice(bdxId, l)
        if (ke !== giliran.current) return
        if (formOffice.current !== null && urlOffice.current !== null) {
          urlOffice.current.value = url
          formOffice.current.submit()
        }
        setPenampil((p) => (p === null ? p : { ...p, siap: true }))
      } catch (g: unknown) {
        if (ke === giliran.current) tutupPenampil()
        throw g
      }
    })
  }

  /** View pdf / gambar: isi dari rute unduh beridentitas -> objek URL lokal di popup penampil. */
  const bukaOnline = (l: Lampiran) => {
    const jenis = jenisViewOnline(l.ekstensi)
    if (sibuk || jenis === null) return
    const ke = ++giliran.current
    setPenampil({ nama: l.fileName, jenis, siap: false })
    jalankan(async () => {
      try {
        const isi = await ambilIsiLampiran(bdxId, l)
        if (ke !== giliran.current) return
        // Tipe dari ekstensi, bukan dari jawaban: hanya pdf / gambar raster yang dirender.
        const objekURL = URL.createObjectURL(new Blob([isi], { type: mimeViewOnline(l.ekstensi) }))
        objekAktif.current = objekURL
        if (jenis === 'gambar' && gambar.current !== null) gambar.current.src = objekURL
        if (jenis === 'pdf' && bingkai.current !== null) bingkai.current.src = objekURL
        setPenampil((p) => (p === null ? p : { ...p, siap: true }))
      } catch (g: unknown) {
        if (ke === giliran.current) tutupPenampil()
        throw g
      }
    })
  }

  // ⛔ `key` BERBEDA untuk popup penampil dan popup daftar: keduanya `<Modal>` di posisi yang sama, sehingga tanpa
  // `key` React memakai ulang instans yang sama - popup daftar mewarisi status "sedang menutup" milik penampil
  // (backdrop keluar menutupi layar, timer tutup sudah terpakai: Close / X / Escape mati). Bug 08-10-2026.
  if (penampil !== null) {
    return (
      <Modal key="penampil" judul={penampil.nama} onTutup={tutupPenampil} labelBatal={LAMPIRAN_BDX.close} penuh>
        {!penampil.siap && <Memuat />}
        {/* `View Office Online` b388: `<penampil>?src=<URL bertanda tangan>` di bingkai popup ini. */}
        <form ref={formOffice} method="get" action={PENAMPIL_OFFICE} target={BINGKAI_PENAMPIL} hidden>
          <input ref={urlOffice} type="hidden" name={PARAM_PENAMPIL} />
        </form>
        {penampil.jenis === 'gambar' ? (
          <img ref={gambar} alt={penampil.nama} className="bordereaux__penampil-gambar" />
        ) : (
          <iframe
            ref={bingkai}
            name={BINGKAI_PENAMPIL}
            title={penampil.nama}
            className="bordereaux__penampil-bingkai"
          />
        )}
      </Modal>
    )
  }

  return (
    <Modal
      key="daftar"
      judul={LAMPIRAN_BDX.judulLihat(kategori.nama)}
      onTutup={() => onTutup(berubah)}
      labelBatal={LAMPIRAN_BDX.close}
      lebar
    >
      {galat !== null && <Gagal galat={galat} />}
      {daftar === null && galat === null && <Memuat pesan={LAMPIRAN_BDX.memuat} />}
      {daftar !== null &&
        (daftar.length === 0 ? (
          <p className="muted">{LAMPIRAN_BDX.kosong}</p>
        ) : (
          <table className="inbox__tabel bordereaux__grid">
            <thead>
              <tr>
                <th>{LAMPIRAN_BDX.fileName}</th>
                <th aria-label={LAMPIRAN_BDX.view} />
                <th>{LAMPIRAN_BDX.type}</th>
                <th aria-label={LAMPIRAN_BDX.delete} />
              </tr>
            </thead>
            <tbody>
              {daftar.map((l) => (
                <tr key={l.id} className="inbox__baris">
                  <td>
                    <button
                      type="button"
                      className="bordereaux__lampiran-nama"
                      disabled={sibuk}
                      onClick={() => jalankan(() => unduhLampiran(bdxId, l))}
                    >
                      {l.fileName}
                    </button>
                  </td>
                  <td>
                    <div className="bordereaux__lampiran-lihat">
                      {jenisViewOnline(l.ekstensi) !== null && (
                        <button
                          type="button"
                          className="btn btn--ghost btn--sm"
                          disabled={sibuk}
                          onClick={() => bukaOnline(l)}
                        >
                          {LAMPIRAN_BDX.view}
                        </button>
                      )}
                      {bisaViewOffice(l.ekstensi) && (
                        <button
                          type="button"
                          className="btn btn--ghost btn--sm"
                          disabled={sibuk}
                          onClick={() => bukaOffice(l)}
                        >
                          {LAMPIRAN_BDX.viewOffice}
                        </button>
                      )}
                    </div>
                  </td>
                  <td>{l.kategori}</td>
                  <td>
                    {boleh && (
                      <button
                        type="button"
                        className="btn btn--danger btn--sm"
                        disabled={sibuk}
                        onClick={() =>
                          jalankan(async () => {
                            await hapusLampiran(bdxId, l)
                            setBerubah(true)
                            muat()
                          })
                        }
                      >
                        {LAMPIRAN_BDX.delete}
                      </button>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        ))}
    </Modal>
  )
}
