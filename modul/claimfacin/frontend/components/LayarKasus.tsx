// Layar satu kasus Claim Fac In - FlowAction InputRegister (Assignment1), InputEstimasi (Assignment7), atau
// InputSurveyor / Choose Surveyor (Assignment3). Isinya pohon tata server; aksi dikirim bersama nilai semua medan
// terbuka. Panel baris grid masterDetail (objek -> item -> estimasi / adjustment) di bawah barisnya, rekursif
// (`TataView` + `rincian.ts`); local action / harness (`Layar.modal`) sebagai Modal; pop-up pemilih / tampilan yang
// isinya bukan tata halaman sebagai `Popup` (`sesudahAksi.ts`). Pola
// `modul/claimnonprop/frontend/components/LayarKasus.tsx` (disalin, bukan impor).

import { useCallback, useEffect, useMemo, useRef, useState } from 'react'

import { Gagal, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { ApiFailure } from '../../../../inti/frontend/klien'
import {
  aksiKasus,
  ambilAcuan,
  bukaKasus,
  GalatValidasiAksi,
  pilihanKasus,
  type AcuanStatis,
  type BarisPolisCari,
  type Halaman,
  type Layar,
  type Pilihan,
  type Tata,
} from '../api'
import { putuskanAksi, type AksiDiminta } from '../antreAksi'
import { CFI } from '../labels'
import { ambil, kunciOpsi, masukan, modeKirim, semuaTata, setel, sumberServer } from '../nilai'
import PanelLampiran from './PanelLampiran'
import Popup, { HasilPolis } from './Popup'
import { judulModal, lanjutanAksi, MODAL_POLIS, modalSesudah, type JenisPopup } from './sesudahAksi'
import { pisahKakiModal } from './susun'
import TataView, { konteksDi, LayarTata, Tombol, type KonteksTata } from './TataView'

/** Daftar FacRetroList halaman polis (RetroList_Harnness). */
const DAFTAR_RETRO = 'OfferFacIn.FacRetroList'

/** Tombol aktif beraksi `aksi` di layout utama (tab Lampiran: Save = tombol Save layar, pola Claim Prop). */
function adaTombol(ts: readonly Tata[], aksi: string): boolean {
  return ts.some((t) => (t.jenis === 'tombol' && t.aksi === aksi && !t.nonaktif) || adaTombol(t.anak ?? [], aksi))
}

/** Aksi tombol Save Input Register (click:save). */
const AKSI_SIMPAN = 'Simpan'

/** Satu aksi layar beserta ubahan medan yang memicunya. */
interface Permintaan extends AksiDiminta {
  ubahan: Record<string, string>
  param: string
  /** Tombol kaki modal: modal ditutup bila aksinya berhasil tanpa pesan. */
  tutupModal: boolean
}

export default function LayarKasus({
  id,
  pelaku,
  onKembali,
  hanyaLihat = false,
}: {
  id: string
  pelaku: string
  onKembali: () => void
  /** Tampilan saja walau pemegang. */
  hanyaLihat?: boolean
}) {
  const [layar, setLayar] = useState<Layar | null>(null)
  const [h, setH] = useState<Halaman | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [sibuk, setSibuk] = useState(false)
  const [info, setInfo] = useState<string | null>(null)
  const [modal, setModal] = useState<string | null>(null)
  const [popup, setPopup] = useState<{ jenis: JenisPopup; indeks: number } | null>(null)
  const [acuan, setAcuan] = useState<AcuanStatis | null>(null)
  const [opsiPeta, setOpsiPeta] = useState<Record<string, Pilihan[]>>({})
  const [versi, setVersi] = useState(0)
  const [hasilPolis, setHasilPolis] = useState<BarisPolisCari[] | null>(null)

  const terima = useCallback((l: Layar) => {
    setLayar(l)
    setH(l.halaman)
    setGalat(null)
    setVersi((v) => v + 1)
  }, [])

  useEffect(() => {
    bukaKasus(id, hanyaLihat).then(terima, (g: unknown) => setGalat(g))
    ambilAcuan().then(setAcuan, (g: unknown) => setGalat(g))
  }, [id, hanyaLihat, terima])

  // Grid hasil Choose Polis (SearchPolis_act): Search Type / Search Name yang baru dikirim server.
  const cariPolis = useCallback(
    (hk: Halaman) => {
      pilihanKasus<BarisPolisCari[]>(id, 'polis', {
        jenisCari: hk.nilai.SearchType ?? '',
        cari: hk.nilai.SearchName ?? '',
      }).then(
        (d) => setHasilPolis(d ?? []),
        (g: unknown) => setGalat(g),
      )
    },
    [id],
  )

  // Aksi yang sedang berjalan dan satu aksi yang menunggu (`antreAksi.ts`).
  const berjalan = useRef<Permintaan | null>(null)
  const menunggu = useRef<Permintaan | null>(null)

  const jalankan = useCallback(
    function jalan(p: Permintaan, ly: Layar, hKini: Halaman) {
      const { aksi, indeks, konteks, ubahan, param, tutupModal } = p
      let h2 = hKini
      for (const [j, v] of Object.entries(ubahan)) h2 = setel(h2, j, v)
      setH(h2)
      const semua = semuaTata(ly.tata, [ly.panel ?? {}, ly.modal ?? {}])
      berjalan.current = p
      setSibuk(true)
      setInfo(null)
      aksiKasus(id, {
        aksi,
        konteks,
        indeks,
        param,
        tahap: ly.kasus.tahap,
        masukan: masukan(h2, semua),
        mode: modeKirim(ly.mode, h2),
      }).then(
        (l) => {
          terima(l)
          if (l.info) setInfo(l.info)
          setModal((m) => modalSesudah(m, l, tutupModal))
          const lanjut = lanjutanAksi(aksi, indeks)
          setPopup(lanjut && 'popup' in lanjut ? { jenis: lanjut.popup, indeks: lanjut.indeks } : null)
          if (lanjut && 'cariPolis' in lanjut) cariPolis(h2)
          berjalan.current = null
          const berikut = menunggu.current
          menunggu.current = null
          if (berikut) jalan(berikut, l, l.halaman)
          else setSibuk(false)
        },
        (g: unknown) => {
          berjalan.current = null
          menunggu.current = null
          setSibuk(false)
          if (g instanceof GalatValidasiAksi) {
            // Validasi layar gagal: aksi dibatalkan server; isian lokal TETAP (h tidak diganti), pesan halaman dan pesan
            // per medan dari layar validasinya.
            setLayar((lama) => {
              const dasar = g.layar ?? lama
              return dasar ? { ...dasar, pesan: g.pesan } : dasar
            })
            setGalat(null)
          } else setGalat(g)
        },
      )
    },
    [id, terima, cariPolis],
  )

  const kirim = useCallback(
    (aksi: string, indeks = 0, ubahan: Record<string, string> = {}, konteks = '', param = '', tutupModal = false) => {
      if (!layar || !h) return
      const p: Permintaan = { aksi, indeks, konteks, ubahan, param, tutupModal }
      const putusan = putuskanAksi(berjalan.current, p)
      if (putusan === 'abaikan') return
      if (putusan === 'antre') {
        menunggu.current = p
        return
      }
      jalankan(p, layar, h)
    },
    [layar, h, jalankan],
  )

  // Pilihan server: autocomplete ber-saringan dan dropdown bersumber server, disimpan per sumber + panel + baris.
  const saran = useCallback(
    (sumber: string, konteks: string, indeks: number, cari: string) => {
      if (!sumberServer(sumber)) return
      const kunci = kunciOpsi(sumber, konteks, indeks)
      pilihanKasus<Pilihan[]>(id, sumber, { konteks, indeks, cari }).then(
        (d) => setOpsiPeta((o) => ({ ...o, [kunci]: d ?? [] })),
        () => undefined,
      )
    },
    [id],
  )

  const opsi = useCallback(
    (sumber: string, konteks: string, indeks: number): Pilihan[] => {
      if (sumber === 'mataUang') return acuan?.mataUang ?? []
      if (sumber === 'jenisReas') return acuan?.jenisReas ?? []
      if (sumber.startsWith('kode:')) {
        const p = sumber.slice(5)
        return (acuan?.kode[p] ?? []).map((k) => ({ nilai: k, label: acuan?.labelKode?.[p]?.[k] ?? k }))
      }
      return opsiPeta[kunciOpsi(sumber, konteks, indeks)] ?? []
    },
    [acuan, opsiPeta],
  )

  const k: KonteksTata | null = useMemo(() => {
    if (!layar || !h) return null
    return {
      h,
      ubah: (j, v) => setH((x) => (x ? setel(x, j, v) : x)),
      aksi: (nama, indeks, ubahan, konteks, param) => kirim(nama, indeks, ubahan, konteks, param ?? ''),
      opsi,
      saran,
      pesanMedan: layar.pesanMedan ?? {},
      sibuk,
      konteks: '',
      panel: layar.panel ?? {},
      versi,
    }
  }, [layar, h, kirim, opsi, saran, sibuk, versi])

  // Grid hasil Choose Polis hidup selama modalnya terbuka.
  useEffect(() => {
    if (modal !== MODAL_POLIS) setHasilPolis(null)
  }, [modal])

  // Pesan halaman tampil di atas layar; aksi dari bagian bawah menggulir ke kotak pesan supaya penolakannya terlihat.
  const kotakPesan = useRef<HTMLDivElement>(null)
  const pesanLayar = layar?.pesan
  useEffect(() => {
    if ((pesanLayar ?? []).length > 0) kotakPesan.current?.scrollIntoView({ block: 'nearest', behavior: 'smooth' })
  }, [pesanLayar])

  if (!layar || !h || !k) {
    return (
      <section className="inbox claimfacin__akar">
        {galat ? <Gagal galat={galat} /> : <Memuat pesan={CFI.memuat} />}
      </section>
    )
  }

  const pesanGalat = galat instanceof ApiFailure ? galat.detail.message : null
  const pesan = layar.pesan ?? []
  // Pesan dan galat aksi tampil di atas layar DAN di modal / pop-up yang sedang terbuka (aksi dari dalamnya).
  const kotakGalat =
    galat !== null && (pesanGalat ? <div className="alert alert--error">{pesanGalat}</div> : <Gagal galat={galat} />)
  const kotakPesanModal = pesan.length > 0 && (
    <div className="alert alert--error">
      <ul>
        {pesan.map((p, i) => (
          <li key={i}>{p}</li>
        ))}
      </ul>
    </div>
  )
  // Tombol di akhir isi modal ke kaki `Modal`; tombol penutup section = tombol batal `Modal` (bukan dua Cancel).
  const tutupModal = () => setModal(null)
  const isiModal = modal ? pisahKakiModal(layar.modal?.[modal] ?? []) : null
  const kModal: KonteksTata | null = modal ? konteksDi(k, modal, tutupModal) : null
  const kKaki: KonteksTata | null = kModal
    ? {
        ...kModal,
        aksi: (nama, indeks, ubahan, konteks, param) => kirim(nama, indeks, ubahan, konteks, param ?? '', true),
      }
    : null
  return (
    <section className="inbox claimfacin__akar">
      <header className="inbox__kepala">
        <button type="button" className="btn btn--ghost btn--sm" onClick={onKembali}>
          {CFI.kembali}
        </button>
        <h2 className="inbox__judul">
          {layar.kasus.id} - {layar.label || layar.kasus.statusWork}
        </h2>
      </header>
      {!layar.bolehKerja && <div className="alert">{CFI.hanyaLihat}</div>}
      {kotakGalat}
      {pesan.length > 0 && (
        <div ref={kotakPesan} className="alert alert--error">
          <ul>
            {pesan.map((p, i) => (
              <li key={i}>{p}</li>
            ))}
          </ul>
        </div>
      )}
      {info && (
        <div className="alert alert--ok" role="status">
          {info}
        </div>
      )}
      <LayarTata tata={layar.tata} k={k} />
      <PanelLampiran
        id={id}
        hanyaLihat={hanyaLihat}
        bolehSimpan={layar.bolehKerja && adaTombol(layar.tata, AKSI_SIMPAN)}
        onSimpan={() => kirim(AKSI_SIMPAN)}
      />
      {modal && isiModal && kModal && kKaki && (
        <Modal
          judul={judulModal(modal)}
          onTutup={tutupModal}
          lebar
          labelBatal={isiModal.batal}
          aksi={isiModal.kaki.map((t, i) => (
            <Tombol key={i} t={t} k={kKaki} utama />
          ))}
        >
          {kotakGalat}
          {kotakPesanModal}
          <TataView tata={isiModal.isi} k={kModal} />
          {modal === MODAL_POLIS && hasilPolis !== null && (
            <HasilPolis
              baris={hasilPolis}
              sibuk={sibuk}
              onPilih={(param) => kirim('CopyNB', 0, {}, MODAL_POLIS, param)}
            />
          )}
        </Modal>
      )}
      {popup && (
        // Satu Popup per jenis: data jenis lama tidak pernah dirender dengan bentuk jenis baru. Layar baru (Save
        // katastrofe baru) memuat ulang isinya.
        <Popup
          key={`${popup.jenis}:${popup.indeks}:${versi}`}
          id={id}
          jenis={popup.jenis}
          indeks={popup.indeks}
          pelaku={pelaku}
          sts={{ sts: ambil(h, 'ClaimData.StsKatastrofe'), non: ambil(h, 'ClaimData.NonKatastrofeType') }}
          retro={h.daftar[DAFTAR_RETRO] ?? []}
          peringatan={
            <>
              {kotakGalat}
              {kotakPesanModal}
            </>
          }
          onPilih={(a, p) => kirim(a, 0, {}, '', p)}
          onTutup={() => setPopup(null)}
        />
      )}
    </section>
  )
}
