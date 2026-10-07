// Layar satu kasus endorsemen - SATU section untuk semua posisi (flow action `InboxPolicyTreatyInAddendum` admin dan
// `DeptHeadTreatyIn_UWAddendum` atasan sama-sama memakai `Section/GeneralPolicyTreatyInAddendum` ->
// `DetailPolicyTreatyInAddendum` using page `.PolicyTreatyIn`). Asal pola: `modul/nbtreatyin/frontend/pages/
// LayarKasus.tsx` (06-10-2026: kepala + Back, tata letak dua kolom, popup tolak / nomor polis, kaki aksi).
//
//   S2 "General" (NOHEADER)   kolom kiri S3-S6 / kanan S7-S9 (`medan.ts`); tombol Choose Business (S6, PN = ADM)
//   S10 Remark                disabled bila PN != ADM
//   S11 `.IsNewPolicyNonProp = 1`          `components/AddPremi.tsx`
//   S17 `.IsNewPolicyNonProp = 0 || ''`    tab Old Data / New Data / Value Difference (`components/TabData.tsx`)
//   ListSuggestEDM            `components/Usulan.tsx`
//   S18 / S19                 Save (admin) + Submit menurut `Layar.tombol` (backend: syarat IsApproved per posisi)
//
// ⛔ Setiap refresh berhitung dikirim ke backend (`POST .../hitung`); action set yang memuat lebih dari satu
// refresh dikirim SEKALI sebagai `urutan`. Medan wajib dari backend (`medanWajib`). Layar hanya-baca bila
// `!bolehKerja`.

import { Fragment, useCallback, useRef, useState } from 'react'

import { Gagal, Memuat, Modal, StripTab, type Opsi } from '../../../../inti/frontend/components/ui/dasar'
import { useAmbilBatal } from '../ambil'
import {
  POLIS,
  POSISI_ADMIN,
  ambilAcuan,
  bukaKasus,
  daftar,
  hitung,
  kirimKasus,
  nilai,
  pilihBisnis,
  setel,
  setelDaftar,
  simpanKasus,
  type Acuan,
  type Halaman,
  type Layar,
} from '../api'
import AddPremi from '../components/AddPremi'
import KotakMedan from '../components/KotakMedan'
import PilihBisnis from '../components/PilihBisnis'
import TabData from '../components/TabData'
import Usulan, { JALUR_APPROVAL } from '../components/Usulan'
import Wadah from '../components/Wadah'
import { buatPenjagaIsian, type PenjagaIsian, type UbahHalaman } from '../penjagaIsian'
import {
  JUDUL,
  JUDUL_POSISI,
  KONFIRMASI_TOLAK,
  NOMOR_DIAKSEP,
  PESAN,
  PORTAL,
  TAB_DATA,
  TOMBOL,
  type TabData as NamaTab,
} from '../labels'
import {
  MEDAN_ANGSURAN_NP,
  MEDAN_KANAN,
  MEDAN_KIRI,
  deretLayer,
  deretQ,
  medanRemark,
  medanTampil,
  nonPropBaru,
  tampilTabData,
  tataTab,
  varianTab,
  type Aksi,
  type Medan,
  type SumberAcuan,
} from '../medan'

function opsi(p: { nilai: string; label: string }[] | null | undefined): Opsi[] {
  return (p ?? []).map((x) => ({ value: x.nilai, label: x.label }))
}

export default function LayarKasus({
  id,
  onKembali,
  awal,
}: {
  id: string
  onKembali: (pesan?: string) => void
  /** Keadaan awal yang sudah dimuat (uji render); tanpa ini layar memuat sendiri. */
  awal?: { layar: Layar; acuan: Acuan }
}) {
  const [layar, setLayar] = useState<Layar | null>(awal?.layar ?? null)
  const [h, setH] = useState<Halaman | null>(awal?.layar.halaman ?? null)
  const [acuan, setAcuan] = useState<Acuan | null>(awal?.acuan ?? null)
  const [galat, setGalat] = useState<unknown>(null)
  // cacah permintaan berjalan / mengantre (penjaga isian menjalankannya berurutan)
  const [jalan, setJalan] = useState(0)
  const sibuk = jalan > 0
  const [info, setInfo] = useState('')
  const [tab, setTab] = useState<NamaTab>(TAB_DATA[0])
  const [popupBisnis, setPopupBisnis] = useState(false)
  const [konfirmasi, setKonfirmasi] = useState(false)
  const [nomor, setNomor] = useState(false)

  // Tinjauan kode 06-10-2026: jawaban server mengganti seluruh halaman - isian yang diketik selama permintaan
  // berjalan diterapkan ulang di atasnya, dan permintaan berangkat berurutan dengan halaman terbaru (`penjagaIsian`).
  const penjagaRef = useRef<PenjagaIsian | null>(null)
  if (penjagaRef.current === null) penjagaRef.current = buatPenjagaIsian(awal?.layar.halaman ?? null)
  const penjaga = penjagaRef.current

  const terima = useCallback(
    (ly: Layar) => {
      penjaga.pasang(ly.halaman)
      setLayar(ly)
      setH(ly.halaman)
    },
    [penjaga],
  )

  useAmbilBatal(
    () => (awal ? Promise.resolve([awal.layar, awal.acuan] as const) : Promise.all([bukaKasus(id), ambilAcuan()])),
    ([ly, a]) => {
      terima(ly)
      setAcuan(a)
    },
    setGalat,
    [id, terima],
  )

  if (galat !== null && layar === null) return <Gagal galat={galat} />
  if (layar === null || h === null)
    return (
      <div className="inbox edmt__akar">
        <Memuat />
      </div>
    )

  const posisi = layar.kasus.positionNote
  const admin = posisi === POSISI_ADMIN
  const boleh = layar.bolehKerja
  const wajib = new Set(layar.medanWajib ?? [])
  const daftarAcuan: Record<SumberAcuan, Opsi[]> = {
    mataUang: opsi(acuan?.mataUang),
    mo: opsi(acuan?.mo),
    jenisEdm: opsi(acuan?.jenisEdm),
  }

  /** Satu permintaan lewat antrean penjaga: `minta` menerima halaman terbaru saat berangkat. */
  async function jalankan<T>(
    minta: (hh: Halaman) => Promise<T>,
    halamanDari: (x: T) => Halaman | null,
    lanjut: (x: T) => void,
  ) {
    setJalan((n) => n + 1)
    setGalat(null)
    setInfo('')
    try {
      lanjut(await penjaga.kirim(minta, halamanDari))
    } catch (e: unknown) {
      setGalat(e)
    } finally {
      setJalan((n) => n - 1)
    }
  }
  /** Jawaban berhalaman: kasus dari server, halaman dari penjaga (jawaban + isian susulan). */
  const terimaLayar = (ly: Layar) => {
    setLayar(ly)
    setH(penjaga.kini())
  }
  const ubahH = (f: UbahHalaman) => {
    penjaga.ubah(f)
    setH(penjaga.kini())
  }

  const ubah = (jalur: string, v: string) => ubahH((x) => setel(x, jalur, v))
  /** Action set satu sel / tombol - SATU bentuk permintaan: `urutan` satu refresh atau lebih. */
  const refresh = (urutan: Aksi[], indeks?: number) => {
    if (urutan.length === 0 || !boleh) return
    void jalankan(
      (hh) => hitung(id, { urutan, indeks, halaman: hh }),
      (ly) => ly.halaman,
      terimaLayar,
    )
  }
  const selesai = (m: Medan, v: string) => {
    ubahH((x) => setel(x, m.jalur, v))
    if (m.aksi) refresh(m.aksi)
  }
  const ubahBaris = (jalur: string, i: number, kunci: string, v: string) =>
    ubahH((x) =>
      setelDaftar(
        x,
        jalur,
        daftar(x, jalur).map((b, j) => (j === i ? { ...b, [kunci]: v } : b)),
      ),
    )

  const simpan = () =>
    void jalankan(
      (hh) => simpanKasus(id, hh),
      (ly) => ly.halaman,
      (ly) => {
        terimaLayar(ly)
        setInfo(PESAN.tersimpan)
      },
    )
  const kirim = () =>
    void jalankan(
      (hh) => kirimKasus(id, hh),
      () => null,
      (hasil) => onKembali(hasil.pesanKonversi ? `${PESAN.terkirim} ${hasil.pesanKonversi}` : PESAN.terkirim),
    )

  const kotak = (m: Medan, i: number) => (
    <KotakMedan
      key={`${m.jalur}-${m.label}-${i}`}
      medan={m}
      halaman={h}
      wajib={wajib.has(m.jalur)}
      hanyaBaca={!boleh || (m.hanyaAdmin === true && !admin)}
      opsi={daftarAcuan}
      onUbah={ubah}
      onSelesai={selesai}
    />
  )
  /** Satu kolom label-kiri: deret "Q / U/Y" dan Layer satu baris. */
  const kolom = (ms: Medan[]) => (
    <div className="edmt__kolom">
      {deretQ(medanTampil(ms, h)).map((x, i) =>
        Array.isArray(x) ? (
          <div key={`deret-${i}`} className={deretLayer(x) ? 'edmt__deret edmt__deret--layer' : 'edmt__deret'}>
            {x.map(kotak)}
          </div>
        ) : (
          <Fragment key={`${x.jalur}-${i}`}>{kotak(x, i)}</Fragment>
        ),
      )}
    </div>
  )

  const tombolKirim = () => {
    if (!boleh) return null
    const klik: Record<string, () => void> = {
      // S19 [IsApproved==1] SetDueTo_act -> finishAssignment; S18 [IsApproved='0'] / Sec Head -> finishAssignment
      kirim,
      // S19 [IsApproved==0] -> localAction PolicyTreatyInDeclineConfirm
      'konfirmasi-tolak': () => setKonfirmasi(true),
      // S18 Dept Head [IsApproved==1] -> localAction ShowPolicyNoTreaty (PolicyNo sudah terisi)
      'nomor-polis': () => setNomor(true),
    }
    const f = klik[layar.tombol]
    if (f === undefined) return null
    return (
      <button type="button" className="btn btn--primary" disabled={sibuk} onClick={f}>
        {TOMBOL.submit}
      </button>
    )
  }

  const tata = tataTab(varianTab(tab, h, admin))

  return (
    <div className="inbox edmt__layar edmt__akar">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">
          {JUDUL_POSISI[posisi] ?? JUDUL.portal} — {layar.kasus.id}
        </h2>
        <button type="button" className="btn btn--ghost" onClick={() => onKembali()}>
          {TOMBOL.kembali}
        </button>
      </header>
      {!boleh && <div className="alert alert--info">{PORTAL.hanyaBaca}</div>}
      {galat !== null && <Gagal galat={galat} />}
      {info && <div className="alert alert--ok">{info}</div>}
      {(layar.pesan ?? []).length > 0 && (
        <div className="alert alert--warn">
          {(layar.pesan ?? []).map((p) => (
            <div key={p}>{p}</div>
          ))}
        </div>
      )}

      {/* S2 (pyTitle "General", NOHEADER: tanpa judul) */}
      <Wadah>
        {admin && boleh && (
          <div className="edmt__aksi">
            {/* S6 `pyWorkPage.PositionNote = 'ReasTreatyInAdmin'`: refresh (preDT SetOldMasterNo) -> showHarness popup */}
            <button type="button" className="btn" disabled={sibuk} onClick={() => setPopupBisnis(true)}>
              {TOMBOL.chooseBusiness}
            </button>
          </div>
        )}
        <div className="edmt__dua-kolom">
          {kolom(MEDAN_KIRI)}
          {kolom(MEDAN_KANAN)}
        </div>
        <div className="edmt__baris-penuh">{kolom([medanRemark(admin)])}</div>
      </Wadah>

      {nonPropBaru(h) && (
        <AddPremi
          halaman={h}
          suntingAngsuran={admin && boleh}
          opsiJenisReas={acuan?.jenisReas ?? []}
          onUbah={ubah}
          onRefreshAngsuran={() => refresh(MEDAN_ANGSURAN_NP.aksi ?? [])}
        />
      )}

      {tampilTabData(h) && (
        // `DetailPolicyTreatyInAddGeneralEditable` S2108: layout group Tab
        <section className="edmt__tab-data">
          <StripTab tab={TAB_DATA} aktif={tab} onPilih={setTab} />
          <TabData
            tata={tata}
            halaman={h}
            wajib={wajib}
            boleh={boleh}
            sibuk={sibuk}
            opsi={daftarAcuan}
            opsiJenisReas={acuan?.jenisReas ?? []}
            onUbah={ubah}
            onSelesai={selesai}
            onUbahBaris={ubahBaris}
            onRefresh={(urutan, indeks) => refresh(urutan, indeks)}
          />
        </section>
      )}

      <Usulan
        halaman={h}
        boleh={boleh}
        onUbah={ubah}
        onApproval={(v) => {
          ubah(JALUR_APPROVAL, v)
          // change -> SetDueTo_act -> Protection_Act (CekLimitTreatyAcc_Act tidak ada di aksi backend)
          refresh([{ aksi: 'SetDueTo' }, { aksi: 'Protection' }])
        }}
      />

      {boleh && (
        <div className="edmt__kaki">
          {admin && (
            // S19 Save
            <button type="button" className="btn" disabled={sibuk} onClick={simpan}>
              {TOMBOL.save}
            </button>
          )}
          {tombolKirim()}
        </div>
      )}

      {popupBisnis && (
        <PilihBisnis
          id={id}
          halaman={h}
          sibuk={sibuk}
          onTutup={() => setPopupBisnis(false)}
          onPilih={() => {
            setPopupBisnis(false)
            void jalankan(
              (hh) => pilihBisnis(id, hh),
              (ly) => ly.halaman,
              terimaLayar,
            )
          }}
        />
      )}
      {konfirmasi && (
        <Modal
          judul={JUDUL.tolak}
          onTutup={() => setKonfirmasi(false)}
          labelBatal={TOMBOL.no}
          aksi={
            <button
              type="button"
              className="btn btn--primary"
              onClick={() => {
                setKonfirmasi(false)
                kirim()
              }}
            >
              {TOMBOL.yes}
            </button>
          }
        >
          <p>{KONFIRMASI_TOLAK}</p>
        </Modal>
      )}
      {nomor && (
        // Tanpa Cancel / X / Escape / klik luar (pola NB, perintah work owner 06-10-2026: "selalu maju"): OK = kirim.
        <Modal
          judul={JUDUL.nomorPolis}
          onTutup={() => setNomor(false)}
          tanpaTutup
          aksi={
            <button
              type="button"
              className="btn btn--primary"
              onClick={() => {
                setNomor(false)
                kirim()
              }}
            >
              {TOMBOL.ok}
            </button>
          }
        >
          {/* Section ShowPolicyNoTreaty_SC: pyWorkPage.pyID, LABEL "telah diaksep menjadi", PolicyTreatyIn.PolicyNo */}
          <div className="edmt__nomor">
            <p className="edmt__nomor-kasus">{layar.kasus.id}</p>
            <p className="edmt__nomor-teks">{NOMOR_DIAKSEP}</p>
            <p className="edmt__nomor-polis">{nilai(h, POLIS + 'PolicyNo')}</p>
          </div>
        </Modal>
      )}
    </div>
  )
}
